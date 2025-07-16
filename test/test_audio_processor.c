#include "../src/core/audio_processor.h"
#include "../src/io/audio_output.h"
#include "../src/utils/memory_utils.h"
#include "../src/types/audio_types.h"
#include "../src/types/sequence_types.h"


#include <stdio.h>
#include <assert.h>
#include <stdlib.h>
#include <math.h>

#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

void test_audio_processor_creation() {
    printf("Testing audio processor creation...\n");
    
    AudioProcessor* processor = audio_processor_create(44100);
    assert(processor != NULL);
    assert(processor->sample_rate == 44100);
    assert(processor->initialized == false);
    
    ResultCode result = audio_processor_initialize(processor);
    assert(result == RESULT_SUCCESS);
    assert(processor->initialized == true);
    
    audio_processor_destroy(processor);
    printf("✓ Audio processor creation tests passed\n");
}

void test_binaural_voice_setup() {
    printf("Testing binaural voice setup...\n");
    
    AudioProcessor* processor = audio_processor_create(44100);
    audio_processor_initialize(processor);
    
    // Create a binaural voice: 250Hz carrier, 8Hz beat
    Voice binaural_voice = {
        .type = VOICE_TYPE_BINAURAL,
        .amplitude = AMPLITUDE_PERCENT_TO_VALUE(50.0),  // 50% amplitude
        .carrier = 250.0,     // 250Hz carrier frequency
        .resonance = 8.0,     // 8Hz binaural beat frequency
        .waveform = WAVE_TYPE_SINE
    };
    
    ResultCode result = audio_processor_setup_voice(processor, 0, &binaural_voice);
    assert(result == RESULT_SUCCESS);
    
    VoiceProcessor* voice = &processor->voices[0];
    assert(voice->active == true);
    assert(voice->type == VOICE_TYPE_BINAURAL);
    assert(voice->carrier_frequency == 250.0);
    assert(voice->modulation_frequency == 8.0);
    assert(voice->channels[0] != NULL);  // Left channel
    assert(voice->channels[1] != NULL);  // Right channel
    
    // Check that frequencies are correctly set
    // Left: 250 - 8/2 = 246Hz, Right: 250 + 8/2 = 254Hz
    assert(voice->channels[0]->frequency == 246.0);
    assert(voice->channels[1]->frequency == 254.0);
    
    printf("  Left frequency: %.1fHz, Right frequency: %.1fHz\n", 
           voice->channels[0]->frequency, voice->channels[1]->frequency);
    printf("  Beat frequency: %.1fHz\n", 
           voice->channels[1]->frequency - voice->channels[0]->frequency);
    
    audio_processor_destroy(processor);
    printf("✓ Binaural voice setup tests passed\n");
}

void test_binaural_sample_generation() {
    printf("Testing binaural sample generation...\n");
    
    AudioProcessor* processor = audio_processor_create(44100);
    audio_processor_initialize(processor);
    
    // Create a simple binaural voice
    Voice binaural_voice = {
        .type = VOICE_TYPE_BINAURAL,
        .amplitude = AMPLITUDE_PERCENT_TO_VALUE(25.0),  // 25% amplitude
        .carrier = 440.0,     // 440Hz carrier (A4 note)
        .resonance = 10.0,    // 10Hz binaural beat
        .waveform = WAVE_TYPE_SINE
    };
    
    audio_processor_setup_voice(processor, 0, &binaural_voice);
    
    // Generate some samples
    int16_t left_buffer[1024];
    int16_t right_buffer[1024];
    
    audio_processor_generate_samples(processor, left_buffer, right_buffer, 1024);
    
    // Check that we have signal in both channels
    bool left_has_signal = false;
    bool right_has_signal = false;
    bool channels_different = false;
    
    for (int i = 0; i < 1024; i++) {
        if (left_buffer[i] != 0) left_has_signal = true;
        if (right_buffer[i] != 0) right_has_signal = true;
        if (left_buffer[i] != right_buffer[i]) channels_different = true;
    }
    
    assert(left_has_signal);
    assert(right_has_signal);
    assert(channels_different);  // L/R should be different due to frequency difference
    
    printf("  Generated %d samples\n", 1024);
    printf("  Left channel has signal: %s\n", left_has_signal ? "yes" : "no");
    printf("  Right channel has signal: %s\n", right_has_signal ? "yes" : "no");
    printf("  Channels are different: %s\n", channels_different ? "yes" : "no");
    
    // Print first few samples for inspection
    printf("  First 10 samples:\n");
    for (int i = 0; i < 10; i++) {
        printf("    [%d] L: %6d, R: %6d, diff: %6d\n", 
               i, left_buffer[i], right_buffer[i], right_buffer[i] - left_buffer[i]);
    }
    
    audio_processor_destroy(processor);
    printf("✓ Binaural sample generation tests passed\n");
}

void test_multiple_binaural_voices() {
    printf("Testing multiple binaural voices...\n");
    
    AudioProcessor* processor = audio_processor_create(44100);
    audio_processor_initialize(processor);
    
    // Create two different binaural voices
    Voice alpha_voice = {
        .type = VOICE_TYPE_BINAURAL,
        .amplitude = AMPLITUDE_PERCENT_TO_VALUE(30.0),
        .carrier = 250.0,
        .resonance = 8.0,      // 8Hz alpha wave
        .waveform = WAVE_TYPE_SINE
    };
    
    Voice beta_voice = {
        .type = VOICE_TYPE_BINAURAL,
        .amplitude = AMPLITUDE_PERCENT_TO_VALUE(20.0),
        .carrier = 150.0,
        .resonance = 20.0,     // 20Hz beta wave
        .waveform = WAVE_TYPE_SINE
    };
    
    audio_processor_setup_voice(processor, 0, &alpha_voice);
    audio_processor_setup_voice(processor, 1, &beta_voice);
    
    // Generate samples with both voices active
    int16_t left_buffer[512];
    int16_t right_buffer[512];
    
    audio_processor_generate_samples(processor, left_buffer, right_buffer, 512);
    
    // Check that both voices contribute to the output
    bool has_signal = false;
    int32_t max_amplitude = 0;
    
    for (int i = 0; i < 512; i++) {
        if (left_buffer[i] != 0 || right_buffer[i] != 0) {
            has_signal = true;
        }
        
        int32_t amplitude = abs(left_buffer[i]) + abs(right_buffer[i]);
        if (amplitude > max_amplitude) {
            max_amplitude = amplitude;
        }
    }
    
    assert(has_signal);
    
    printf("  Two binaural voices mixed successfully\n");
    printf("  Maximum combined amplitude: %d\n", max_amplitude);
    
    audio_processor_destroy(processor);
    printf("✓ Multiple binaural voices tests passed\n");
}

void test_binaural_wav_generation() {
    printf("Testing binaural WAV file generation...\n");
    
    AudioProcessor* processor = audio_processor_create(44100);
    audio_processor_initialize(processor);
    
    // Create a 10Hz alpha binaural beat at 250Hz carrier
    Voice alpha_binaural = {
        .type = VOICE_TYPE_BINAURAL,
        .amplitude = AMPLITUDE_PERCENT_TO_VALUE(25.0),  // 60% amplitude for clear hearing
        .carrier = 250.0,     // 250Hz carrier frequency  
        .resonance = 10.0,    // 10Hz alpha binaural beat
        .waveform = WAVE_TYPE_SINE
    };
    
    audio_processor_setup_voice(processor, 0, &alpha_binaural);
    
    // Generate 10 seconds of audio (44100 * 10 samples)
    int duration_seconds = 10;
    int sample_count = 44100 * duration_seconds;
    int16_t* left_buffer = memory_alloc(sample_count * sizeof(int16_t));
    int16_t* right_buffer = memory_alloc(sample_count * sizeof(int16_t));
    
    printf("  Generating %d samples (%.1f seconds) of binaural audio...\n", 
           sample_count, (float)duration_seconds);
    
    // Generate all samples
    audio_processor_generate_samples(processor, left_buffer, right_buffer, sample_count);

    printf("  Audio generation debug:\n");
    printf("    Duration: %d seconds\n", duration_seconds);
    printf("    Sample rate: 44100 Hz\n");
    printf("    Expected total samples: %d\n", duration_seconds * 44100);
    printf("    Actual sample_count: %d\n", sample_count);
    printf("    Stereo buffer size: %d elements\n", sample_count * 2);
    printf("    Expected file data size: %d bytes\n", sample_count * 2 * 2);  // * 2 channels * 2 bytes

    printf("  Analyzing generated samples...\n");

    int16_t min_left = 32767, max_left = -32768;
    int16_t min_right = 32767, max_right = -32768;
    int zero_count = 0;

    for (int i = 0; i < 1000; i++) {  // Analyze first 1000 samples
        if (left_buffer[i] < min_left) min_left = left_buffer[i];
        if (left_buffer[i] > max_left) max_left = left_buffer[i];
        if (right_buffer[i] < min_right) min_right = right_buffer[i];
        if (right_buffer[i] > max_right) max_right = right_buffer[i];
        if (left_buffer[i] == 0 && right_buffer[i] == 0) zero_count++;
    }

    printf("  Sample range analysis (first 1000 samples):\n");
    printf("    Left:  min=%d, max=%d\n", min_left, max_left);
    printf("    Right: min=%d, max=%d\n", min_right, max_right);
    printf("    Zero samples: %d/1000\n", zero_count);

    // Print first few samples
    printf("  First 10 samples:\n");
    for (int i = 0; i < 10; i++) {
        printf("    [%d] L:%6d R:%6d\n", i, left_buffer[i], right_buffer[i]);
    }
    
    // Create audio output configuration
    AudioConfig config = {
        .sample_rate = 44100,
        .channels = 2,
        .bits_per_sample = 16,
        .is_raw_output = false,
        .output_filename = "test_binaural_10hz_alpha.wav",
        .quiet_mode = false
    };
    
    // Create audio output and write WAV file
    AudioOutput* output = audio_output_create_wav(&config);
    assert(output != NULL);
    
    printf("  Writing to file: %s\n", config.output_filename);
    
    // Interleave the stereo samples (L, R, L, R, ...)
    int16_t* stereo_buffer = memory_alloc(sample_count * 2 * sizeof(int16_t));

    printf("  Byte order check (first sample):\n");
int16_t first_sample = stereo_buffer[0];
uint8_t* bytes = (uint8_t*)&first_sample;
printf("    Sample value: %d\n", first_sample);
printf("    Byte 0: 0x%02X, Byte 1: 0x%02X\n", bytes[0], bytes[1]);
printf("    Expected little-endian: low byte first\n");
    for (int i = 0; i < sample_count; i++) {
        stereo_buffer[i * 2] = left_buffer[i];      // Left channel
        stereo_buffer[i * 2 + 1] = right_buffer[i]; // Right channel
    }
    
    ResultCode result = audio_output_write_samples(output, stereo_buffer, sample_count);
    assert(result == RESULT_SUCCESS);
    
    audio_output_finalize(output);
    audio_output_destroy(output);
    
    // Cleanup
    memory_free(left_buffer);
    memory_free(right_buffer);
    memory_free(stereo_buffer);
    audio_processor_destroy(processor);
    
    printf("✓ WAV file generated successfully!\n");
    printf("  File: test_binaural_10hz_alpha.wav\n");
    printf("  Format: 44.1kHz, 16-bit, Stereo\n");
    printf("  Content: 250Hz carrier with 10Hz alpha binaural beat\n");
    printf("  Duration: %d seconds\n", duration_seconds);
    printf("  Instructions: Use STEREO HEADPHONES to hear the binaural effect!\n");
    printf("✓ Binaural WAV generation tests passed\n");
}

void test_multiple_binaural_wav_generation() {
    printf("Testing multiple binaural beats WAV generation...\n");
    
    AudioProcessor* processor = audio_processor_create(44100);
    audio_processor_initialize(processor);
    
    // Create a complex binaural session with multiple frequencies
    Voice alpha_voice = {
        .type = VOICE_TYPE_BINAURAL,
        .amplitude = AMPLITUDE_PERCENT_TO_VALUE(25.0),
        .carrier = 250.0,
        .resonance = 8.0,      // 8Hz alpha
        .waveform = WAVE_TYPE_SINE
    };
    
    Voice theta_voice = {
        .type = VOICE_TYPE_BINAURAL,
        .amplitude = AMPLITUDE_PERCENT_TO_VALUE(15.0),
        .carrier = 150.0,
        .resonance = 6.0,      // 6Hz theta  
        .waveform = WAVE_TYPE_SINE
    };
    
    audio_processor_setup_voice(processor, 0, &alpha_voice);
    audio_processor_setup_voice(processor, 1, &theta_voice);
    
    // Generate 5 seconds of complex binaural audio
    int duration_seconds = 5;
    int sample_count = 44100 * duration_seconds;
    int16_t* left_buffer = memory_alloc(sample_count * sizeof(int16_t));
    int16_t* right_buffer = memory_alloc(sample_count * sizeof(int16_t));
    
    printf("  Generating complex binaural session...\n");
    printf("  Voice 1: 250Hz carrier, 8Hz alpha beat\n");
    printf("  Voice 2: 150Hz carrier, 6Hz theta beat\n");
    
    audio_processor_generate_samples(processor, left_buffer, right_buffer, sample_count);

    printf("  Audio generation debug:\n");
    printf("    Duration: %d seconds\n", duration_seconds);
    printf("    Sample rate: 44100 Hz\n");
    printf("    Expected total samples: %d\n", duration_seconds * 44100);
    printf("    Actual sample_count: %d\n", sample_count);
    printf("    Stereo buffer size: %d elements\n", sample_count * 2);
    printf("    Expected file data size: %d bytes\n", sample_count * 2 * 2);  // * 2 channels * 2 bytes
    
    // Save as WAV
    AudioConfig config = {
        .sample_rate = 44100,
        .channels = 2,
        .bits_per_sample = 16,
        .is_raw_output = false,
        .output_filename = "test_binaural_complex.wav",
        .quiet_mode = false
    };
    
    AudioOutput* output = audio_output_create_wav(&config);
    
    // Interleave stereo samples
    int16_t* stereo_buffer = memory_alloc(sample_count * 2 * sizeof(int16_t));
    for (int i = 0; i < sample_count; i++) {
        stereo_buffer[i * 2] = left_buffer[i];
        stereo_buffer[i * 2 + 1] = right_buffer[i];
    }

    printf("  Verifying stereo interleaving (first 6 samples):\n");
    for (int i = 0; i < 3; i++) {
        printf("    Frame %d: L=%d, R=%d -> Stereo[%d]=%d, Stereo[%d]=%d\n", 
            i, left_buffer[i], right_buffer[i],
            i*2, stereo_buffer[i*2], i*2+1, stereo_buffer[i*2+1]);
    }
    
    ResultCode result = audio_output_write_samples(output, stereo_buffer, sample_count);
    assert(result == RESULT_SUCCESS);
    
    audio_output_finalize(output);
    audio_output_destroy(output);
    
    // Cleanup
    memory_free(left_buffer);
    memory_free(right_buffer);
    memory_free(stereo_buffer);
    audio_processor_destroy(processor);
    
    printf("✓ Complex binaural WAV generated: test_binaural_complex.wav\n");
    printf("✓ Multiple binaural WAV generation tests passed\n");
}

void test_debug_simple_wav() {
    printf("Testing debug simple WAV with known values...\n");
    
    // Generate a simple 1-second, 1000Hz sine wave manually
    int sample_rate = 44100;
    int duration = 1;  // 1 second
    int sample_count = sample_rate * duration;
    
    int16_t* left_buffer = memory_alloc(sample_count * sizeof(int16_t));
    int16_t* right_buffer = memory_alloc(sample_count * sizeof(int16_t));
    
    // Generate pure 1000Hz sine wave with known amplitude
    for (int i = 0; i < sample_count; i++) {
        double t = (double)i / sample_rate;
        double frequency = 1000.0;  // 1kHz
        double amplitude = 8000.0;  // Moderate amplitude
        
        int16_t sample = (int16_t)(amplitude * sin(2.0 * M_PI * frequency * t));
        left_buffer[i] = sample;
        right_buffer[i] = sample;  // Same in both channels for now
    }
    
    printf("  Generated pure 1000Hz sine wave\n");
    printf("  First 5 samples: ");
    for (int i = 0; i < 5; i++) {
        printf("%d ", left_buffer[i]);
    }
    printf("\n");
    
    // Save as WAV
    AudioConfig config = {
        .sample_rate = sample_rate,
        .channels = 2,
        .bits_per_sample = 16,
        .is_raw_output = false,
        .output_filename = "debug_1000hz_pure.wav",
        .quiet_mode = false
    };
    
    AudioOutput* output = audio_output_create_wav(&config);
    
    // Interleave
    int16_t* stereo_buffer = memory_alloc(sample_count * 2 * sizeof(int16_t));
    for (int i = 0; i < sample_count; i++) {
        stereo_buffer[i * 2] = left_buffer[i];
        stereo_buffer[i * 2 + 1] = right_buffer[i];
    }
    
    audio_output_write_samples(output, stereo_buffer, sample_count * 2);
    audio_output_finalize(output);
    audio_output_destroy(output);
    
    memory_free(left_buffer);
    memory_free(right_buffer);
    memory_free(stereo_buffer);
    
    printf("✓ Pure sine wave generated: debug_1000hz_pure.wav\n");
    printf("  This should sound like a clear 1000Hz tone\n");
}

int main() {
    printf("=== Testing Audio Processor Module ===\n");
    
    test_audio_processor_creation();
    test_binaural_voice_setup();
    test_binaural_sample_generation();
    test_multiple_binaural_voices();
    test_binaural_wav_generation();
    test_multiple_binaural_wav_generation();
    test_debug_simple_wav();
    
    printf("✓ All audio processor tests passed!\n");
    return 0;
}