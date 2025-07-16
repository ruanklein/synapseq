#ifndef AUDIO_TYPES_H
#define AUDIO_TYPES_H

#include "common_types.h"

typedef enum {
    WAVE_TYPE_SINE = 0,
    WAVE_TYPE_SQUARE = 1,
    WAVE_TYPE_TRIANGLE = 2,
    WAVE_TYPE_SAWTOOTH = 3
} WaveType;

typedef enum {
    VOICE_TYPE_OFF = 0,
    VOICE_TYPE_BINAURAL = 1,
    VOICE_TYPE_MONAURAL = 2,
    VOICE_TYPE_ISOCHRONIC = 3,
    VOICE_TYPE_WHITE_NOISE = 4,
    VOICE_TYPE_PINK_NOISE = 5,
    VOICE_TYPE_BROWN_NOISE = 6,
    VOICE_TYPE_SPIN_WHITE = 7,
    VOICE_TYPE_SPIN_PINK = 8,
    VOICE_TYPE_SPIN_BROWN = 9,
    VOICE_TYPE_BACKGROUND = 10,
    VOICE_TYPE_EFFECT_SPIN = 11,
    VOICE_TYPE_EFFECT_PULSE = 12,
} VoiceType;

typedef struct {
    VoiceType type;
    double amplitude;
    double carrier_frequency;
    double resonance_frequency;
    WaveType waveform;
} VoiceConfig;

typedef struct {
    int sample_rate;
    int channels;
    int bits_per_sample;
    bool is_raw_output;
    char output_filename[MAX_FILENAME_LENGTH];
    bool quiet_mode;
    int volume_percentage;
} AudioConfig;

typedef struct {
    double carrier_frequency;
    double resonance_frequency;
    double amplitude;
    WaveType waveform;
    double time_offset;
} WaveParameters;

#endif // AUDIO_TYPES_H