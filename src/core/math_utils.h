#ifndef MATH_UTILS_H
#define MATH_UTILS_H

#include "../types/audio_types.h"
#include "../types/common_types.h"

#include <stdint.h>

#define MATH_WAVE_TABLE_SIZE 16384
#define MATH_WAVE_AMPLITUDE 0x7FFFF  // Original ST_AMP = 524287
#define MATH_NOISE_AMPLITUDE_SHIFT 12
#define MATH_NOISE_AMPLITUDE (MATH_WAVE_AMPLITUDE << MATH_NOISE_AMPLITUDE_SHIFT)
#define MATH_NOISE_DITHER_SHIFT 16
#define MATH_NOISE_BANDS 9


typedef struct {
    int* sine_table;
    int* square_table;
    int* triangle_table;
    int* sawtooth_table;
    bool initialized;
} WaveTables;

typedef struct {
    int value;
    int increment;
} NoiseBand;

typedef struct {
    NoiseBand bands[MATH_NOISE_BANDS];
    int offset;
    int noise_buffer[256];
    uint8_t buffer_index;
    int brown_last_value;
    uint32_t random_seed;
} NoiseGenerator;

typedef struct {
    NoiseGenerator* noise_generator;
    int noise_type;
} NoiseSpin;

WaveTables* math_wave_tables_create(void);
void math_wave_tables_destroy(WaveTables* tables);
ResultCode math_wave_tables_initialize(WaveTables* tables);

NoiseGenerator* math_noise_generator_create(void);
void math_noise_generator_destroy(NoiseGenerator* generator);
void math_noise_generator_set_seed(NoiseGenerator* generator, uint32_t seed);

NoiseSpin* math_noise_spin_create(NoiseGenerator* generator, int noise_type);
void math_noise_spin_destroy(NoiseSpin* spin);
void math_noise_spin_apply(NoiseSpin* spin, int amplitude, int spin_position, int* left, int* right);

int math_wave_table_get_sample(const WaveTables* tables, WaveType type, int phase);

int math_noise_pink_generate(NoiseGenerator* generator);
int math_noise_white_generate(NoiseGenerator* generator);
int math_noise_brown_generate(NoiseGenerator* generator);

double math_amplitude_display_to_internal(double percentage);
double math_amplitude_internal_to_display(int amplitude);

int32_t math_wave_table_get_sample_32bit(WaveTables* tables, WaveType waveform, uint32_t phase);

uint32_t math_random_next(uint32_t* seed);

#endif // MATH_UTILS_H