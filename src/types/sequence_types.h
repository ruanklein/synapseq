#ifndef SEQUENCE_TYPES_H
#define SEQUENCE_TYPES_H

#include "audio_types.h"
#include "time_types.h"
#include "common_types.h"

#include <stdio.h>

#define SEQUENCE_AMPLITUDE_SCALE 40.96
#define SEQUENCE_FADE_INTERVAL_MS 60000

#define AMPLITUDE_PERCENT_TO_VALUE(percent) ((percent) * SEQUENCE_AMPLITUDE_SCALE)
#define AMPLITUDE_VALUE_TO_PERCENT(value) ((value) / SEQUENCE_AMPLITUDE_SCALE)

typedef struct Voice {
    VoiceType type;
    double amplitude;
    double carrier;
    double resonance;
    WaveType waveform;
} Voice;

typedef struct Period {
    struct Period* next;
    struct Period* prev;
    time_ms_t start_time;
    Voice start_voices[MAX_CHANNELS];
    Voice end_voices[MAX_CHANNELS];
    int fade_in_mode;
    int fade_out_mode;
} Period;

typedef struct NameDefinition {
    struct NameDefinition* next;
    char* name;
    Voice voices[MAX_CHANNELS];
} NameDefinition;

typedef struct SequenceOptions {
    char* background_file;
    int volume_level;
    int sample_rate;
    double gain_level_db;
    bool sample_rate_user_set;
} SequenceOptions;

typedef struct SequenceParseResult {
    Period* periods;
    NameDefinition* name_definitions;
    SequenceOptions options;
    time_ms_t first_time;
    time_ms_t last_time;
    int line_count;
    ResultCode result;
} SequenceParseResult;

typedef struct SequenceParseContext {
    FILE* file;
    char* filename;
    int current_line;
    char line_buffer[MAX_LINE_LENGTH];
    char line_copy[MAX_LINE_LENGTH];
    char* line_position;
    NameDefinition* name_definitions;
    Period* periods;
    SequenceOptions options;
    time_ms_t first_time;
    time_ms_t last_time;
    bool has_background_effect;
} SequenceParseContext;

#endif // SEQUENCE_TYPES_H