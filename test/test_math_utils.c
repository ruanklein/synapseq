#include <stdio.h>
#include <assert.h>
#include <stdlib.h>
#include <math.h>
#include "../src/core/math_utils.h"

#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

void test_wave_tables_creation() {
    printf("Testing wave tables creation...\n");
    
    WaveTables* tables = math_wave_tables_create();
    assert(tables != NULL);
    assert(tables->initialized == false);
    
    ResultCode result = math_wave_tables_initialize(tables);
    assert(result == RESULT_SUCCESS);
    assert(tables->initialized == true);
    
    math_wave_tables_destroy(tables);
    printf("✓ Wave tables creation tests passed\n");
}

void test_wave_table_samples() {
    printf("Testing wave table samples...\n");
    
    WaveTables* tables = math_wave_tables_create();
    math_wave_tables_initialize(tables);
    
    // Test sine wave at π/2 (quarter table) should be positive
    int phase_quarter = (MATH_WAVE_TABLE_SIZE / 4) << 16;
    int sine_quarter = math_wave_table_get_sample(tables, WAVE_TYPE_SINE, phase_quarter);
    assert(sine_quarter > 0);
    
    // Test sine wave at 0 should be close to 0
    int phase_zero = 0;
    int sine_zero = math_wave_table_get_sample(tables, WAVE_TYPE_SINE, phase_zero);
    assert(abs(sine_zero) < 1000);
    
    // Test square wave
    int square_quarter = math_wave_table_get_sample(tables, WAVE_TYPE_SQUARE, phase_quarter);
    assert(square_quarter > 0);
    
    math_wave_tables_destroy(tables);
    printf("✓ Wave table samples tests passed\n");
}

void test_noise_generation() {
    printf("Testing noise generation...\n");
    
    NoiseGenerator* generator = math_noise_generator_create();
    assert(generator != NULL);
    
    math_noise_generator_set_seed(generator, 12345);
    
    int white_noise = math_noise_white_generate(generator);
    int pink_noise = math_noise_pink_generate(generator);
    int brown_noise = math_noise_brown_generate(generator);
    
    // Noise should be different values
    assert(white_noise != pink_noise || pink_noise != brown_noise);
    
    math_noise_generator_destroy(generator);
    printf("✓ Noise generation tests passed\n");
}

void test_amplitude_conversions() {
    printf("Testing amplitude conversions...\n");
    
    // Test display to internal
    double internal = math_amplitude_display_to_internal(50.0);
    assert(internal == 2048.0);
    
    // Test internal to display
    double display = math_amplitude_internal_to_display(2048);
    assert(fabs(display - 50.0) < 0.1);
    
    printf("✓ Amplitude conversion tests passed\n");
}

void test_noise_spin() {
    printf("Testing noise spin effect...\n");
    
    NoiseGenerator* generator = math_noise_generator_create();
    NoiseSpin* spin = math_noise_spin_create(generator, 0);
    
    int left = 0, right = 0;
    math_noise_spin_apply(spin, 1000, 0, &left, &right);
    
    // With spin position 0, both channels should get some value
    assert(left != 0 || right != 0);
    
    math_noise_spin_destroy(spin);
    math_noise_generator_destroy(generator);
    printf("✓ Noise spin tests passed\n");
}

void test_random_generator() {
    printf("Testing random generator...\n");
    
    uint32_t seed = 12345;
    uint32_t random1 = math_random_next(&seed);
    uint32_t random2 = math_random_next(&seed);
    
    assert(random1 != random2);
    
    // Test with NULL seed
    uint32_t null_result = math_random_next(NULL);
    assert(null_result == 0);
    
    printf("✓ Random generator tests passed\n");
}

int main() {
    printf("=== Testing Math Utils Module ===\n");
    
    test_wave_tables_creation();
    test_wave_table_samples();
    test_noise_generation();
    test_amplitude_conversions();
    test_noise_spin();
    test_random_generator();
    
    printf("✓ All math utils tests passed!\n");
    return 0;
}