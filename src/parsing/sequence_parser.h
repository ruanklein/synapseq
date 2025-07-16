#ifndef SEQUENCE_PARSER_H
#define SEQUENCE_PARSER_H

#include "../types/sequence_types.h"
#include "../types/common_types.h"
#include "../types/time_types.h"

#include <stdio.h>

SequenceParseResult* sequence_parser_parse_file(const char* filename);
SequenceParseResult* sequence_parser_parse_files(const char* filenames[], int count);
void sequence_parser_result_destroy(SequenceParseResult* result);

NameDefinition* sequence_parser_find_preset(NameDefinition* definitions, const char* name);
NameDefinition* sequence_parser_create_preset(const char* name);
void sequence_parser_preset_destroy(NameDefinition* preset);

Period* sequence_parser_create_period(time_ms_t start_time);
void sequence_parser_period_destroy(Period* period);
void sequence_parser_periods_destroy(Period* periods);

ResultCode sequence_parser_validate_sequence(Period* periods);
ResultCode sequence_parser_normalize_voices(Voice* voices, int channel_count);

ResultCode sequence_parser_parse_time_string(const char* time_str, time_ms_t* result);
ResultCode sequence_parser_parse_voice_line(const char* line, Voice* voices, int* channel_count);
ResultCode sequence_parser_parse_timeline(const char* line, time_ms_t* time, char** preset_name);

const char* sequence_parser_get_voice_type_name(VoiceType type);
const char* sequence_parser_get_wave_type_name(WaveType type);

#endif // SEQUENCE_PARSER_H