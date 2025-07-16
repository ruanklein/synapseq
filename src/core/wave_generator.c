#include "wave_generator.h"
#include "../utils/memory_utils.h"

#include <math.h>

WaveGenerator* wave_generator_create(int sample_rate) {
    WaveGenerator* generator = memory_alloc(sizeof(WaveGenerator));
    generator->wave_tables = NULL;
    generator->sample_rate = sample_rate;
    generator->initialized = false;
    return generator;
}

void wave_generator_destroy(WaveGenerator* generator) {
    if (!generator) return;
    
    if (generator->wave_tables) {
        math_wave_tables_destroy(generator->wave_tables);
    }
    
    memory_free(generator);
}

ResultCode wave_generator_initialize(WaveGenerator* generator) {
    if (!generator) return RESULT_ERROR_INVALID_ARGUMENT;
    
    generator->wave_tables = math_wave_tables_create();
    if (!generator->wave_tables) {
        return RESULT_ERROR_MEMORY;
    }
    
    ResultCode result = math_wave_tables_initialize(generator->wave_tables);
    if (result != RESULT_SUCCESS) {
        math_wave_tables_destroy(generator->wave_tables);
        generator->wave_tables = NULL;
        return result;
    }
    
    generator->initialized = true;
    return RESULT_SUCCESS;
}

WaveChannel* wave_channel_create(WaveType waveform, double frequency, double amplitude) {
    WaveChannel* channel = memory_alloc(sizeof(WaveChannel));
    channel->waveform = waveform;
    channel->frequency = frequency;
    channel->amplitude = amplitude;
    channel->phase = 0.0;
    channel->phase_increment = 0.0;
    channel->phase_accumulator = 0;
    channel->active = true;
    return channel;
}

void wave_channel_destroy(WaveChannel* channel) {
    if (!channel) return;
    memory_free(channel);
}

void wave_channel_set_frequency(WaveChannel* channel, double frequency, int sample_rate) {
    if (!channel) return;
    
    channel->frequency = frequency;
    channel->phase_increment = (frequency * MATH_WAVE_TABLE_SIZE * (1 << PHASE_ACCUMULATOR_BITS)) / sample_rate;
}

void wave_channel_set_amplitude(WaveChannel* channel, double amplitude) {
    if (!channel) return;
    channel->amplitude = amplitude;
}

void wave_channel_set_phase(WaveChannel* channel, double phase) {
    if (!channel) return;
    channel->phase = phase;
    channel->phase_accumulator = (uint32_t)(phase * PHASE_ACCUMULATOR_SIZE / (2.0 * M_PI));
}

void wave_channel_reset_phase(WaveChannel* channel) {
    if (!channel) return;
    channel->phase = 0.0;
    channel->phase_accumulator = 0;
}

int16_t wave_generator_generate_sample(WaveGenerator* generator, WaveChannel* channel) {
    if (!generator || !channel || !generator->initialized || !channel->active) {
        return 0;
    }
    
    int16_t sample = 0;
    
    switch (channel->waveform) {
        case WAVE_TYPE_SINE:
            sample = math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_SINE, channel->phase_accumulator);
            break;
        case WAVE_TYPE_SQUARE:
            sample = math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_SQUARE, channel->phase_accumulator);
            break;
        case WAVE_TYPE_TRIANGLE:
            sample = math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_TRIANGLE, channel->phase_accumulator);
            break;
        case WAVE_TYPE_SAWTOOTH:
            sample = math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_SAWTOOTH, channel->phase_accumulator);
            break;
    }
    
    channel->phase_accumulator += (uint32_t)channel->phase_increment;
    channel->phase_accumulator &= ((MATH_WAVE_TABLE_SIZE << PHASE_ACCUMULATOR_BITS) - 1);
    
    return (int16_t)(sample * channel->amplitude);
}

void wave_generator_generate_samples(WaveGenerator* generator, WaveChannel* channel, int16_t* buffer, int sample_count) {
    if (!generator || !channel || !buffer || !generator->initialized) return;
    
    for (int i = 0; i < sample_count; i++) {
        buffer[i] = wave_generator_generate_sample(generator, channel);
    }
}

int16_t wave_generator_get_sine_sample(WaveGenerator* generator, uint32_t phase) {
    if (!generator || !generator->wave_tables || !generator->initialized) return 0;
    
    // Pass the original phase, not the index
    return math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_SINE, phase << PHASE_ACCUMULATOR_BITS);
}

int16_t wave_generator_get_square_sample(WaveGenerator* generator, uint32_t phase) {
    if (!generator || !generator->wave_tables || !generator->initialized) return 0;
    
    return math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_SQUARE, phase << PHASE_ACCUMULATOR_BITS);
}

int16_t wave_generator_get_triangle_sample(WaveGenerator* generator, uint32_t phase) {
    if (!generator || !generator->wave_tables || !generator->initialized) return 0;
    
    return math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_TRIANGLE, phase << PHASE_ACCUMULATOR_BITS);
}

int16_t wave_generator_get_sawtooth_sample(WaveGenerator* generator, uint32_t phase) {
    if (!generator || !generator->wave_tables || !generator->initialized) return 0;
    
    return math_wave_table_get_sample(generator->wave_tables, WAVE_TYPE_SAWTOOTH, phase << PHASE_ACCUMULATOR_BITS);
}

void wave_generator_mix_channels(WaveGenerator* generator, WaveChannel* channels[], int channel_count, int16_t* left_buffer, int16_t* right_buffer, int sample_count) {
    if (!generator || !channels || !left_buffer || !right_buffer || channel_count <= 0 || sample_count <= 0) return;
    
    for (int i = 0; i < sample_count; i++) {
        left_buffer[i] = 0;
        right_buffer[i] = 0;
    }
    
    for (int ch = 0; ch < channel_count; ch++) {
        if (!channels[ch] || !channels[ch]->active) continue;
        
        for (int i = 0; i < sample_count; i++) {
            int16_t sample = wave_generator_generate_sample(generator, channels[ch]);
            left_buffer[i] += sample;
            right_buffer[i] += sample;
        }
    }
}

int32_t wave_generator_generate_sample_32bit(WaveGenerator* generator, WaveChannel* channel) {
    if (!generator || !channel || !generator->initialized || !channel->active) {
        return 0;
    }
    
    int32_t sample = 0;
    
    switch (channel->waveform) {
        case WAVE_TYPE_SINE:
            sample = math_wave_table_get_sample_32bit(generator->wave_tables, WAVE_TYPE_SINE, channel->phase_accumulator);
            break;
        case WAVE_TYPE_SQUARE:
            sample = math_wave_table_get_sample_32bit(generator->wave_tables, WAVE_TYPE_SQUARE, channel->phase_accumulator);
            break;
        case WAVE_TYPE_TRIANGLE:
            sample = math_wave_table_get_sample_32bit(generator->wave_tables, WAVE_TYPE_TRIANGLE, channel->phase_accumulator);
            break;
        case WAVE_TYPE_SAWTOOTH:
            sample = math_wave_table_get_sample_32bit(generator->wave_tables, WAVE_TYPE_SAWTOOTH, channel->phase_accumulator);
            break;
    }
    
    channel->phase_accumulator += (uint32_t)channel->phase_increment;
    channel->phase_accumulator &= ((MATH_WAVE_TABLE_SIZE << PHASE_ACCUMULATOR_BITS) - 1);
    
    return sample; // Return raw sample from wave tables without amplitude multiplication
}