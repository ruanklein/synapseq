#include "../src/io/audio_output.h"
#include "../src/types/audio_types.h"

#include <stdio.h>
#include <assert.h>
#include <stdlib.h>
#include <math.h>

void test_wav_output() {
    printf("Testing WAV output...\n");
    
    AudioConfig config = {
        .sample_rate = 44100,
        .channels = 2,
        .bits_per_sample = 16,
        .is_raw_output = false,
        .output_filename = "test_output.wav",
        .quiet_mode = false,
        .volume_percentage = 100
    };
    
    AudioOutput* output = audio_output_create_wav(&config);
    assert(output != NULL);
    
    const int duration_samples = config.sample_rate;
    int16_t* samples = malloc(duration_samples * config.channels * sizeof(int16_t));
    
    for (int i = 0; i < duration_samples; i++) {
        double t = (double)i / config.sample_rate;
        double wave = sin(2.0 * M_PI * 440.0 * t);
        int16_t sample = (int16_t)(wave * 16000);
        
        samples[i * 2] = sample;
        samples[i * 2 + 1] = sample;
    }
    
    ResultCode result = audio_output_write_samples(output, samples, duration_samples * config.channels);
    assert(result == RESULT_SUCCESS);
    
    assert(audio_output_get_samples_written(output) == duration_samples * config.channels);
    
    result = audio_output_finalize(output);
    assert(result == RESULT_SUCCESS);
    
    audio_output_destroy(output);
    free(samples);
    
    printf("✓ WAV output test passed\n");
}

void test_raw_output() {
    printf("Testing RAW output...\n");
    
    AudioConfig config = {
        .sample_rate = 44100,
        .channels = 2,
        .bits_per_sample = 16,
        .is_raw_output = true,
        .output_filename = "test_output.raw",
        .quiet_mode = false,
        .volume_percentage = 100
    };
    
    AudioOutput* output = audio_output_create_raw(&config);
    assert(output != NULL);
    
    const int duration_samples = config.sample_rate / 2;
    int16_t* samples = malloc(duration_samples * config.channels * sizeof(int16_t));
    
    for (int i = 0; i < duration_samples; i++) {
        double t = (double)i / config.sample_rate;
        double wave = sin(2.0 * M_PI * 880.0 * t);
        int16_t sample = (int16_t)(wave * 8000);
        
        samples[i * 2] = sample;
        samples[i * 2 + 1] = sample;
    }
    
    ResultCode result = audio_output_write_samples(output, samples, duration_samples * config.channels);
    assert(result == RESULT_SUCCESS);
    
    result = audio_output_finalize(output);
    assert(result == RESULT_SUCCESS);
    
    audio_output_destroy(output);
    free(samples);
    
    printf("✓ RAW output test passed\n");
}

int main() {
    printf("=== Testing Audio Output Module ===\n");
    
    test_wav_output();
    test_raw_output();
    
    printf("\n✓ All audio output tests passed!\n");
    printf("You can test the generated WAV file with: file test_output.wav\n");
    printf("Or play it with: play test_output.wav\n");
    
    return 0;
}