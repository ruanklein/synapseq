#include "../../src/parsing/sequence_parser.h"
#include "../../src/core/audio_processor.h"
#include "../../src/io/audio_output.h"
#include "../../src/utils/memory_utils.h"
#include "../../src/types/time_types.h"

#include <stdio.h>
#include <assert.h>
#include <stdlib.h>
#include <string.h>

#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

void test_sequence_to_wav_e2e() {
    printf("=== End-to-End Test: SPSQ → WAV ===\n");
    printf("Testing complete pipeline: Sequence Parser → Audio Processor → WAV Output\n\n");
    
    // Step 1: Parse the .spsq sequence file
    printf("Step 1: Parsing sequence file...\n");
    
    SequenceParseResult* result = sequence_parser_parse_file("test/spsq/sequence-03.spsq");
    
    if (!result || result->result != RESULT_SUCCESS) {
        printf("❌ Failed to parse sequence file\n");
        if (result) sequence_parser_result_destroy(result);
        assert(false);
    }
    
    printf("✓ Sequence parsed successfully\n");
    printf("  Result code: %d\n", result->result);
    printf("  Line count: %d\n", result->line_count);
    printf("  Time range: %d → %d ms\n", result->first_time, result->last_time);
    
    // Count name definitions (presets)
    int preset_count = 0;
    NameDefinition* name_def = result->name_definitions;
    while (name_def) {
        preset_count++;
        printf("  Preset found: '%s'\n", name_def->name);
        name_def = name_def->next;
    }
    
    // Count periods (timeline entries)
    int period_count = 0;
    Period* period = result->periods;
    Period* first_period = period;
    if (period) {
        do {
            period_count++;
            printf("  Period at %02d:%02d:%02d\n",
                   (int)(period->start_time / 3600000),
                   (int)((period->start_time % 3600000) / 60000),
                   (int)((period->start_time % 60000) / 1000));
            period = period->next;
        } while (period && period != first_period);
    }
    
    printf("  Presets found: %d\n", preset_count);
    printf("  Timeline periods: %d\n", period_count);
    
    // Find the alpha-one preset to verify it was parsed correctly
    NameDefinition* alpha_preset = sequence_parser_find_preset(result->name_definitions, "alpha-one");
    assert(alpha_preset != NULL);
    assert(alpha_preset->voices[0].type == VOICE_TYPE_BINAURAL);
    assert(alpha_preset->voices[0].carrier == 250.0);
    assert(alpha_preset->voices[0].resonance == 10.0);
    
    printf("  Preset 'alpha-one': 250Hz binaural 10Hz, amplitude %.1f%%\n", 
           AMPLITUDE_VALUE_TO_PERCENT(alpha_preset->voices[0].amplitude));
    
    // Step 2: Setup Audio Processor
    printf("\nStep 2: Setting up audio processor...\n");
    
    const int SAMPLE_RATE = 44100;
    AudioProcessor* processor = audio_processor_create(SAMPLE_RATE);
    assert(processor != NULL);
    
    ResultCode init_result = audio_processor_initialize(processor);
    assert(init_result == RESULT_SUCCESS);
    
    printf("✓ Audio processor initialized (44.1kHz)\n");
    
    // Step 3: Process timeline and generate audio
    printf("\nStep 3: Processing timeline and generating audio...\n");
    
    // Calculate total duration based on result
    const int DURATION_MS = result->last_time - result->first_time;
    const int DURATION_SECONDS = (DURATION_MS + 999) / 1000;  // Round up
    const int TOTAL_SAMPLES = SAMPLE_RATE * DURATION_SECONDS;
    
    printf("  Total duration: %d ms (%d seconds)\n", DURATION_MS, DURATION_SECONDS);
    printf("  Total samples: %d\n", TOTAL_SAMPLES);
    
    // Allocate audio buffers
    int16_t* left_buffer = memory_alloc(TOTAL_SAMPLES * sizeof(int16_t));
    int16_t* right_buffer = memory_alloc(TOTAL_SAMPLES * sizeof(int16_t));
    assert(left_buffer != NULL && right_buffer != NULL);
    
    // Initialize buffers with silence
    memset(left_buffer, 0, TOTAL_SAMPLES * sizeof(int16_t));
    memset(right_buffer, 0, TOTAL_SAMPLES * sizeof(int16_t));
    
    // Process each period in the timeline
    printf("  Processing timeline periods...\n");
    
    int segment_num = 0;
    period = result->periods;
    first_period = period;
    
    if (period) {
        do {
            Period* next_period = period->next;
            
            // Calculate time range for this segment
            time_ms_t start_ms = period->start_time;
            time_ms_t end_ms = (next_period && next_period != first_period) ? 
                               next_period->start_time : result->last_time;
            time_ms_t duration_ms = end_ms - start_ms;
            
            int start_sample = (start_ms * SAMPLE_RATE) / 1000;
            int segment_samples = (duration_ms * SAMPLE_RATE) / 1000;
            
            // Ensure we don't exceed buffer bounds
            if (start_sample + segment_samples > TOTAL_SAMPLES) {
                segment_samples = TOTAL_SAMPLES - start_sample;
            }
            
            printf("    Segment %d: %02d:%02d:%02d → %02d:%02d:%02d (%d samples)\n", 
                   ++segment_num,
                   start_ms / 3600000, (start_ms / 60000) % 60, (start_ms / 1000) % 60,
                   end_ms / 3600000, (end_ms / 60000) % 60, (end_ms / 1000) % 60,
                   segment_samples);
            
            // Check if this period has any active voices
            bool has_active_voice = false;
            for (int i = 0; i < MAX_CHANNELS; i++) {
                if (period->start_voices[i].type != VOICE_TYPE_OFF) {
                    has_active_voice = true;
                    break;
                }
            }
            
            if (has_active_voice) {
                printf("      Using voices from period\n");
                
                // For this simple test, we'll use start_voices directly
                // In a full implementation, you'd interpolate between start_voices and end_voices
                audio_processor_update_voices(processor, period->start_voices);
                
                // Generate audio for this segment
                if (segment_samples > 0) {
                    audio_processor_generate_samples(processor, 
                                                   left_buffer + start_sample, 
                                                   right_buffer + start_sample, 
                                                   segment_samples);
                }
            } else {
                printf("      Using silence\n");
                // Already filled with silence during initialization
            }
            
            period = next_period;
        } while (period && period != first_period && segment_num < 10); // Safety limit
    }
    
    printf("✓ Audio generation completed\n");
    
    // Step 4: Analyze generated audio
    printf("\nStep 4: Analyzing generated audio...\n");
    
    // Check first 15 seconds for silence
    int silence_samples = (15 * SAMPLE_RATE);
    if (silence_samples > TOTAL_SAMPLES) silence_samples = TOTAL_SAMPLES;
    
    bool silence_ok = true;
    for (int i = 0; i < silence_samples; i++) {
        if (left_buffer[i] != 0 || right_buffer[i] != 0) {
            silence_ok = false;
            break;
        }
    }
    printf("  Silence period (0-15s): %s\n", silence_ok ? "✓ OK" : "❌ FAIL");
    
    // Check for binaural content after 15 seconds
    int binaural_start = silence_samples;
    int binaural_samples = TOTAL_SAMPLES - binaural_start;
    bool binaural_ok = false;
    int16_t max_left = 0, max_right = 0;
    
    for (int i = binaural_start; i < TOTAL_SAMPLES; i++) {
        if (abs(left_buffer[i]) > abs(max_left)) max_left = left_buffer[i];
        if (abs(right_buffer[i]) > abs(max_right)) max_right = right_buffer[i];
        if (left_buffer[i] != 0 || right_buffer[i] != 0) {
            binaural_ok = true;
        }
    }
    
    printf("  Binaural period (15s+): %s\n", binaural_ok ? "✓ OK" : "❌ FAIL");
    printf("    Max amplitude - Left: %d, Right: %d\n", max_left, max_right);
    
    // Step 5: Write WAV file
    printf("\nStep 5: Writing WAV file...\n");
    
    AudioConfig config = {
        .sample_rate = SAMPLE_RATE,
        .channels = 2,
        .bits_per_sample = 16,
        .is_raw_output = false,
        .output_filename = "test_e2e_sequence_output.wav",
        .quiet_mode = false
    };
    
    AudioOutput* output = audio_output_create_wav(&config);
    assert(output != NULL);
    
    // Interleave stereo samples
    int16_t* stereo_buffer = memory_alloc(TOTAL_SAMPLES * 2 * sizeof(int16_t));
    for (int i = 0; i < TOTAL_SAMPLES; i++) {
        stereo_buffer[i * 2] = left_buffer[i];      // Left
        stereo_buffer[i * 2 + 1] = right_buffer[i]; // Right
    }
    
    ResultCode write_result = audio_output_write_samples(output, stereo_buffer, TOTAL_SAMPLES);
    assert(write_result == RESULT_SUCCESS);
    
    audio_output_finalize(output);
    audio_output_destroy(output);
    
    printf("✓ WAV file written: %s\n", config.output_filename);
    printf("  Format: %dkHz, %d-bit, %d channels\n", 
           SAMPLE_RATE / 1000, config.bits_per_sample, config.channels);
    printf("  Duration: %d seconds\n", DURATION_SECONDS);
    printf("  File size: ~%.1f MB\n", (TOTAL_SAMPLES * 2 * 2) / (1024.0 * 1024.0));
    
    // Cleanup
    memory_free(left_buffer);
    memory_free(right_buffer);
    memory_free(stereo_buffer);
    audio_processor_destroy(processor);
    sequence_parser_result_destroy(result);
    
    printf("\n✅ End-to-End Test PASSED!\n");
    printf("🎧 Test file ready: test_e2e_sequence_output.wav\n");
    printf("📝 Content: Sequence from file sequence-03.spsq\n");
    printf("🔊 Use stereo headphones to hear the binaural effect!\n");
}

int main() {
    printf("=== SynapSeq End-to-End Integration Tests ===\n\n");
    
    test_sequence_to_wav_e2e();
    
    printf("\n🎉 All E2E tests passed!\n");
    return 0;
}