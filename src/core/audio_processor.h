#ifndef AUDIO_PROCESSOR_H
#define AUDIO_PROCESSOR_H

#include "../types/audio_types.h"
#include "../types/sequence_types.h"
#include "../core/wave_generator.h"
#include "../core/math_utils.h"

#define PHASE_MASK ((MATH_WAVE_TABLE_SIZE << 16) - 1)
#define AMPLITUDE_SCALE_FACTOR 1.0  // 1/8192 for math_utils amplitude scaling

typedef struct VoiceProcessor {
    VoiceType type;
    WaveChannel* channels[2];          // Left/Right or dual channels
    NoiseGenerator* noise_generator;
    NoiseSpin* noise_spin;
    
    // Voice-specific state
    double carrier_frequency;
    double modulation_frequency;
    double amplitude;
    WaveType waveform;
    
    // Phase accumulators for complex voices
    uint32_t phase_accumulator[2];
    uint32_t phase_increment[2];
    
    bool active;
} VoiceProcessor;

typedef struct AudioProcessor {
    WaveGenerator* wave_generator;
    VoiceProcessor voices[MAX_CHANNELS];
    int sample_rate;
    
    // Background mixing
    int16_t* background_buffer;
    int background_buffer_size;
    double background_gain;
    
    bool initialized;
} AudioProcessor;

// Core functions
AudioProcessor* audio_processor_create(int sample_rate);
void audio_processor_destroy(AudioProcessor* processor);
ResultCode audio_processor_initialize(AudioProcessor* processor);

// Voice management
ResultCode audio_processor_setup_voice(AudioProcessor* processor, int channel, const Voice* voice);
void audio_processor_reset_voice(AudioProcessor* processor, int channel);
void audio_processor_update_voices(AudioProcessor* processor, const Voice voices[MAX_CHANNELS]);

// Audio generation
void audio_processor_generate_samples(AudioProcessor* processor, int16_t* left_buffer, int16_t* right_buffer, int sample_count);
void audio_processor_process_voice(AudioProcessor* processor, int channel, int16_t* left_sample, int16_t* right_sample);

// Specific voice processors
void audio_processor_process_binaural(VoiceProcessor* voice, WaveGenerator* generator, int16_t* left_sample, int16_t* right_sample);
void audio_processor_process_monaural(VoiceProcessor* voice, WaveGenerator* generator, int16_t* left_sample, int16_t* right_sample);
void audio_processor_process_isochronic(VoiceProcessor* voice, WaveGenerator* generator, int16_t* left_sample, int16_t* right_sample);
void audio_processor_process_noise(VoiceProcessor* voice, int16_t* left_sample, int16_t* right_sample);
void audio_processor_process_spin_noise(VoiceProcessor* voice, int16_t* left_sample, int16_t* right_sample);
void audio_processor_process_background(VoiceProcessor* voice, AudioProcessor* processor, int16_t* left_sample, int16_t* right_sample);

// Utility functions
void audio_processor_set_background_buffer(AudioProcessor* processor, int16_t* buffer, int buffer_size, double gain);
void audio_processor_normalize_output(int16_t* left_sample, int16_t* right_sample);

#endif