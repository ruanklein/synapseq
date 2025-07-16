#ifndef WAVE_GENERATOR_H
#define WAVE_GENERATOR_H

#include "../types/audio_types.h"
#include "../types/common_types.h"
#include "../core/math_utils.h"

#define PHASE_ACCUMULATOR_BITS 16
#define PHASE_ACCUMULATOR_SIZE (1 << PHASE_ACCUMULATOR_BITS)
#define PHASE_ACCUMULATOR_MASK (PHASE_ACCUMULATOR_SIZE - 1)

typedef struct WaveGenerator {
    WaveTables* wave_tables;
    int sample_rate;
    bool initialized;
} WaveGenerator;

typedef struct WaveChannel {
    WaveType waveform;
    double frequency;
    double amplitude;
    double phase;
    double phase_increment;
    uint32_t phase_accumulator;
    bool active;
} WaveChannel;

WaveGenerator* wave_generator_create(int sample_rate);
void wave_generator_destroy(WaveGenerator* generator);
ResultCode wave_generator_initialize(WaveGenerator* generator);

WaveChannel* wave_channel_create(WaveType waveform, double frequency, double amplitude);
void wave_channel_destroy(WaveChannel* channel);
void wave_channel_set_frequency(WaveChannel* channel, double frequency, int sample_rate);
void wave_channel_set_amplitude(WaveChannel* channel, double amplitude);
void wave_channel_set_phase(WaveChannel* channel, double phase);
void wave_channel_reset_phase(WaveChannel* channel);

int16_t wave_generator_generate_sample(WaveGenerator* generator, WaveChannel* channel);
void wave_generator_generate_samples(WaveGenerator* generator, WaveChannel* channel, int16_t* buffer, int sample_count);

int16_t wave_generator_get_sine_sample(WaveGenerator* generator, uint32_t phase);
int16_t wave_generator_get_square_sample(WaveGenerator* generator, uint32_t phase);
int16_t wave_generator_get_triangle_sample(WaveGenerator* generator, uint32_t phase);
int16_t wave_generator_get_sawtooth_sample(WaveGenerator* generator, uint32_t phase);

void wave_generator_mix_channels(WaveGenerator* generator, WaveChannel* channels[], int channel_count, int16_t* left_buffer, int16_t* right_buffer, int sample_count);

int32_t wave_generator_generate_sample_32bit(WaveGenerator* generator, WaveChannel* channel);

#endif // WAVE_GENERATOR_H