#include "math_utils.h"
#include "../utils/memory_utils.h"
#include "../types/audio_types.h"
#include "wave_generator.h"

#include <math.h>
#include <string.h>

#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

#define RANDOM_MULTIPLIER 75

WaveTables* math_wave_tables_create(void) {
    WaveTables* tables = memory_alloc(sizeof(WaveTables));
    tables->sine_table = NULL;
    tables->square_table = NULL;
    tables->triangle_table = NULL;
    tables->sawtooth_table = NULL;
    tables->initialized = false;
    return tables;
}

void math_wave_tables_destroy(WaveTables* tables) {
    if (!tables) return;
    
    memory_free(tables->sine_table);
    memory_free(tables->square_table);
    memory_free(tables->triangle_table);
    memory_free(tables->sawtooth_table);
    memory_free(tables);
}

ResultCode math_wave_tables_initialize(WaveTables* tables) {
    if (!tables) return RESULT_ERROR_INVALID_ARGUMENT;
    
    tables->sine_table = memory_alloc_array(MATH_WAVE_TABLE_SIZE, sizeof(int));
    tables->square_table = memory_alloc_array(MATH_WAVE_TABLE_SIZE, sizeof(int));
    tables->triangle_table = memory_alloc_array(MATH_WAVE_TABLE_SIZE, sizeof(int));
    tables->sawtooth_table = memory_alloc_array(MATH_WAVE_TABLE_SIZE, sizeof(int));
    
    for (int i = 0; i < MATH_WAVE_TABLE_SIZE; i++) {
        double phase = (i * 2.0 * M_PI) / MATH_WAVE_TABLE_SIZE;
        
        tables->sine_table[i] = (int)(sin(phase) * MATH_WAVE_AMPLITUDE);
        
        tables->square_table[i] = (sin(phase) >= 0) ? MATH_WAVE_AMPLITUDE : -MATH_WAVE_AMPLITUDE;
        
        if (phase < M_PI) {
            tables->triangle_table[i] = (int)(((2.0 * phase / M_PI) - 1.0) * MATH_WAVE_AMPLITUDE);
        } else {
            tables->triangle_table[i] = (int)((3.0 - (2.0 * phase / M_PI)) * MATH_WAVE_AMPLITUDE);
        }
        
        tables->sawtooth_table[i] = (int)(((2.0 * phase / (2.0 * M_PI)) - 1.0) * MATH_WAVE_AMPLITUDE);
    }
    
    tables->initialized = true;
    return RESULT_SUCCESS;
}

NoiseGenerator* math_noise_generator_create(void) {
    NoiseGenerator* generator = memory_alloc(sizeof(NoiseGenerator));
    
    memset(generator->bands, 0, sizeof(generator->bands));
    generator->offset = 0;
    memset(generator->noise_buffer, 0, sizeof(generator->noise_buffer));
    generator->buffer_index = 0;
    generator->brown_last_value = 0;
    generator->random_seed = 2;
    
    return generator;
}

void math_noise_generator_destroy(NoiseGenerator* generator) {
    memory_free(generator);
}

void math_noise_generator_set_seed(NoiseGenerator* generator, uint32_t seed) {
    if (!generator) return;
    generator->random_seed = seed;
}

int math_wave_table_get_sample(const WaveTables* tables, WaveType type, int phase) {
    if (!tables || !tables->initialized) return 0;
    
    int index = (phase >> 16) & (MATH_WAVE_TABLE_SIZE - 1);
    
    switch (type) {
        case WAVE_TYPE_SINE:
            return tables->sine_table[index];
        case WAVE_TYPE_SQUARE:
            return tables->square_table[index];
        case WAVE_TYPE_TRIANGLE:
            return tables->triangle_table[index];
        case WAVE_TYPE_SAWTOOTH:
            return tables->sawtooth_table[index];
        default:
            return tables->sine_table[index];
    }
}

int math_noise_pink_generate(NoiseGenerator* generator) {
    if (!generator) return 0;
    
    int total = 0;
    int offset = generator->offset++;
    int count = 1;
    NoiseBand* band = generator->bands;
    NoiseBand* band_end = generator->bands + MATH_NOISE_BANDS;
    
    total = ((generator->random_seed = generator->random_seed * RANDOM_MULTIPLIER % 131074) - 65535) *
            (MATH_NOISE_AMPLITUDE / 65535 / (MATH_NOISE_BANDS + 1));
    
    while ((count & offset) && band < band_end) {
        int value = ((generator->random_seed = generator->random_seed * RANDOM_MULTIPLIER % 131074) - 65535) *
                   (MATH_NOISE_AMPLITUDE / 65535 / (MATH_NOISE_BANDS + 1));
        total += band->value += band->increment = (value - band->value) / (count += count);
        band++;
    }
    
    while (band < band_end) {
        total += (band->value += band->increment);
        band++;
    }
    
    return generator->noise_buffer[generator->buffer_index++] = (total >> MATH_NOISE_AMPLITUDE_SHIFT);
}

int math_noise_white_generate(NoiseGenerator* generator) {
    if (!generator) return 0;
    
    return ((generator->random_seed = generator->random_seed * RANDOM_MULTIPLIER % 131074) - 65535) * 
           (MATH_WAVE_AMPLITUDE / 65535);
}

int math_noise_brown_generate(NoiseGenerator* generator) {
    if (!generator) return 0;
    
    int random_value = ((generator->random_seed = generator->random_seed * RANDOM_MULTIPLIER % 131074) - 65535);
    
    generator->brown_last_value = (generator->brown_last_value + (random_value / 16)) * 0.9;
    
    if (generator->brown_last_value > 65535) {
        generator->brown_last_value = 65535;
    }
    if (generator->brown_last_value < -65535) {
        generator->brown_last_value = -65535;
    }
    
    return generator->brown_last_value * (MATH_WAVE_AMPLITUDE / 65535);
}

double math_amplitude_display_to_internal(double percentage) {
    return 40.96 * percentage;
}

double math_amplitude_internal_to_display(int amplitude) {
    return amplitude / 40.96;
}

uint32_t math_random_next(uint32_t* seed) {
    if (!seed) return 0;
    return (*seed = *seed * RANDOM_MULTIPLIER % 131074) - 65535;
}

NoiseSpin* math_noise_spin_create(NoiseGenerator* generator, int noise_type) {
    if (!generator) return NULL;
    
    NoiseSpin* spin = memory_alloc(sizeof(NoiseSpin));
    spin->noise_generator = generator;
    spin->noise_type = noise_type;
    return spin;
}

void math_noise_spin_destroy(NoiseSpin* spin) {
    memory_free(spin);
}

static int spin_noise_generate(NoiseSpin* spin) {
    switch (spin->noise_type) {
    case VOICE_TYPE_BROWN_NOISE:
        return math_noise_brown_generate(spin->noise_generator);
    case VOICE_TYPE_WHITE_NOISE:
        return math_noise_white_generate(spin->noise_generator);
    case VOICE_TYPE_PINK_NOISE:
    default:
        return math_noise_pink_generate(spin->noise_generator);
    }
}

static void spin_noise_calculate_stereo(int base_noise, int amplified_position, int* noise_l, int* noise_r) {
    int positive_position = amplified_position < 0 ? -amplified_position : amplified_position;
    
    if (amplified_position >= 0) {
        *noise_l = (base_noise * (128 - positive_position)) >> 7;
        *noise_r = base_noise + ((base_noise * positive_position) >> 7);
    } else {
        *noise_l = base_noise + ((base_noise * positive_position) >> 7);
        *noise_r = (base_noise * (128 - positive_position)) >> 7;
    }
}

void math_noise_spin_apply(NoiseSpin* spin, int amplitude, int spin_position, int* left, int* right) {
    if (!spin || !left || !right) return;
    
    int amplified_position = (int)(spin_position * 1.5);
    
    if (amplified_position > 127) amplified_position = 127;
    if (amplified_position < -128) amplified_position = -128;
    
    int base_noise = spin_noise_generate(spin);
    
    int noise_l, noise_r;
    spin_noise_calculate_stereo(base_noise, amplified_position, &noise_l, &noise_r);    
    
    *left = amplitude * noise_l;
    *right = amplitude * noise_r;
}

int32_t math_wave_table_get_sample_32bit(WaveTables* tables, WaveType waveform, uint32_t phase) {
    if (!tables) return 0;
    
    int table_index = (phase >> PHASE_ACCUMULATOR_BITS) & (MATH_WAVE_TABLE_SIZE - 1);
    
    switch (waveform) {
        case WAVE_TYPE_SINE:
            return tables->sine_table[table_index];
        case WAVE_TYPE_SQUARE:
            return tables->square_table[table_index];
        case WAVE_TYPE_TRIANGLE:
            return tables->triangle_table[table_index];
        case WAVE_TYPE_SAWTOOTH:
            return tables->sawtooth_table[table_index];
        default:
            return tables->sine_table[table_index];
    }
}