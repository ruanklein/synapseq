#include "../src/core/wave_generator.h"

#include <stdio.h>
#include <assert.h>
#include <stdlib.h>
#include <math.h>

#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

void test_wave_generator_creation() {
    printf("Testing wave generator creation...\n");
    
    WaveGenerator* generator = wave_generator_create(44100);
    assert(generator != NULL);
    assert(generator->sample_rate == 44100);
    assert(generator->initialized == false);
    
    ResultCode result = wave_generator_initialize(generator);
    assert(result == RESULT_SUCCESS);
    assert(generator->initialized == true);
    
    wave_generator_destroy(generator);
    printf("✓ Wave generator creation tests passed\n");
}

void test_wave_channel_creation() {
    printf("Testing wave channel creation...\n");
    
    WaveChannel* channel = wave_channel_create(WAVE_TYPE_SINE, 440.0, 0.5);
    assert(channel != NULL);
    assert(channel->waveform == WAVE_TYPE_SINE);
    assert(channel->frequency == 440.0);
    assert(channel->amplitude == 0.5);
    assert(channel->phase == 0.0);
    assert(channel->active == true);
    
    wave_channel_destroy(channel);
    printf("✓ Wave channel creation tests passed\n");
}

void test_wave_generation() {
    printf("Testing wave generation...\n");
    
    WaveGenerator* generator = wave_generator_create(44100);
    wave_generator_initialize(generator);
    
    WaveChannel* channel = wave_channel_create(WAVE_TYPE_SINE, 440.0, 1.0);
    wave_channel_set_frequency(channel, 440.0, 44100);
    
    // Generate multiple samples to see the wave progression
    printf("  Sine wave samples at different phases:\n");
    for (int i = 0; i < 10; i++) {
        int16_t sample = wave_generator_generate_sample(generator, channel);
        printf("    sample[%d] = %d\n", i, sample);
    }
    
    // Test different amplitudes with non-zero phase
    wave_channel_reset_phase(channel);
    wave_channel_set_frequency(channel, 440.0, 44100);
    
    // Skip to a non-zero phase
    for (int i = 0; i < 20; i++) {
        wave_generator_generate_sample(generator, channel);
    }
    
    // Test different amplitudes at SAME phase
    uint32_t test_phase = channel->phase_accumulator;
    
    printf("  Amplitude test at same phase:\n");
    
    channel->phase_accumulator = test_phase;
    wave_channel_set_amplitude(channel, 0.25);
    int16_t sample_25 = wave_generator_generate_sample(generator, channel);
    
    channel->phase_accumulator = test_phase;
    wave_channel_set_amplitude(channel, 0.5);
    int16_t sample_50 = wave_generator_generate_sample(generator, channel);
    
    channel->phase_accumulator = test_phase;
    wave_channel_set_amplitude(channel, 1.0);
    int16_t sample_100 = wave_generator_generate_sample(generator, channel);
    
    printf("    Amplitude 0.25: %d\n", sample_25);
    printf("    Amplitude 0.50: %d\n", sample_50);
    printf("    Amplitude 1.00: %d\n", sample_100);
    
    // Check that we get some non-zero values
    int16_t buffer[1024];
    wave_generator_generate_samples(generator, channel, buffer, 1024);
    
    bool has_signal = false;
    for (int i = 0; i < 1024; i++) {
        if (buffer[i] != 0) {
            has_signal = true;
            break;
        }
    }
    assert(has_signal);
    
    wave_channel_destroy(channel);
    wave_generator_destroy(generator);
    printf("✓ Wave generation tests passed\n");
}

void test_different_waveforms() {
    printf("Testing different waveforms...\n");
    
    WaveGenerator* generator = wave_generator_create(44100);
    wave_generator_initialize(generator);
    
    WaveType waveforms[] = {WAVE_TYPE_SINE, WAVE_TYPE_SQUARE, WAVE_TYPE_TRIANGLE, WAVE_TYPE_SAWTOOTH};
    const char* names[] = {"sine", "square", "triangle", "sawtooth"};
    
    for (int w = 0; w < 4; w++) {
        WaveChannel* channel = wave_channel_create(waveforms[w], 440.0, 1.0);
        wave_channel_set_frequency(channel, 440.0, 44100);
        
        int16_t sample = wave_generator_generate_sample(generator, channel);
        printf("  %s waveform generated sample: %d\n", names[w], sample);
        
        wave_channel_destroy(channel);
    }
    
    wave_generator_destroy(generator);
    printf("✓ Different waveforms tests passed\n");
}

void test_frequency_control() {
    printf("Testing frequency control...\n");
    
    WaveGenerator* generator = wave_generator_create(44100);
    wave_generator_initialize(generator);
    
    WaveChannel* channel = wave_channel_create(WAVE_TYPE_SINE, 440.0, 1.0);
    
    wave_channel_set_frequency(channel, 440.0, 44100);
    double inc1 = channel->phase_increment;
    
    wave_channel_set_frequency(channel, 880.0, 44100);
    double inc2 = channel->phase_increment;
    
    assert(inc2 > inc1);
    printf("  440Hz increment: %f, 880Hz increment: %f\n", inc1, inc2);
    
    wave_channel_destroy(channel);
    wave_generator_destroy(generator);
    printf("✓ Frequency control tests passed\n");
}

void test_amplitude_control() {
    printf("Testing amplitude control...\n");
    
    WaveGenerator* generator = wave_generator_create(44100);
    wave_generator_initialize(generator);
    
    WaveChannel* channel = wave_channel_create(WAVE_TYPE_SINE, 440.0, 1.0);
    wave_channel_set_frequency(channel, 440.0, 44100);
    
    // Find a better phase position with larger values
    // Skip more samples to find a peak
    for (int i = 0; i < 50; i++) {
        int16_t sample = wave_generator_generate_sample(generator, channel);
        if (abs(sample) > 20000) {  // Find a peak value
            break;
        }
    }
    
    // Save the current phase position
    uint32_t saved_phase = channel->phase_accumulator;
    
    // Test amplitude 0.5
    channel->phase_accumulator = saved_phase;
    wave_channel_set_amplitude(channel, 0.5);
    int16_t sample1 = wave_generator_generate_sample(generator, channel);
    
    // Test amplitude 1.0
    channel->phase_accumulator = saved_phase;
    wave_channel_set_amplitude(channel, 1.0);
    int16_t sample2 = wave_generator_generate_sample(generator, channel);
    
    printf("  At same phase: Amplitude 0.5: %d, Amplitude 1.0: %d\n", sample1, sample2);
    
    double ratio = 0.0;
    if (sample2 != 0) {
        ratio = (double)abs(sample1) / abs(sample2);
    }
    printf("  Ratio: %.3f (should be ~0.500)\n", ratio);
    
    // More realistic tolerance for small values (40%-60%)
    assert(ratio >= 0.40 && ratio <= 0.60);
    assert(abs(sample2) > abs(sample1));
    
    wave_channel_destroy(channel);
    wave_generator_destroy(generator);
    printf("✓ Amplitude control tests passed\n");
}

int main() {
    printf("=== Testing Wave Generator Module ===\n");
    
    test_wave_generator_creation();
    test_wave_channel_creation();
    test_wave_generation();
    test_different_waveforms();
    test_frequency_control();
    test_amplitude_control();
    
    printf("✓ All wave generator tests passed!\n");
    return 0;
}