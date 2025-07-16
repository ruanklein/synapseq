#include "sequence_parser.h"
#include "../utils/memory_utils.h"
#include "../utils/string_utils.h"
#include "../core/time_manager.h"

#include <stdlib.h>
#include <string.h>
#include <ctype.h>

static SequenceParseContext* context_create(const char* filename);
static void context_destroy(SequenceParseContext* context);
static ResultCode context_read_line(SequenceParseContext* context);
static char* context_get_word(SequenceParseContext* context);
static void context_error(SequenceParseContext* context, const char* message);

static ResultCode parse_preset_definition(SequenceParseContext* context);
static ResultCode parse_timeline_entry(SequenceParseContext* context);
static ResultCode parse_option_line(SequenceParseContext* context);
static ResultCode parse_voice_command(SequenceParseContext* context, Voice* voice);

SequenceParseResult* sequence_parser_parse_file(const char* filename) {
    if (!filename) return NULL;
    
    SequenceParseResult* result = memory_alloc(sizeof(SequenceParseResult));
    SequenceParseContext* context = context_create(filename);
    
    if (!context) {
        result->result = RESULT_ERROR_FILE_NOT_FOUND;
        return result;
    }
    
    result->periods = NULL;
    result->name_definitions = NULL;
    result->first_time = TIME_INVALID;
    result->last_time = TIME_INVALID;
    result->line_count = 0;
    result->result = RESULT_SUCCESS;
    
    // Initialize built-in presets
    NameDefinition* silence = sequence_parser_create_preset("silence");
    silence->next = result->name_definitions;
    result->name_definitions = silence;
    
    context->name_definitions = result->name_definitions;
    
    // Main parsing loop
    while (context_read_line(context) == RESULT_SUCCESS) {
        result->line_count++;
        
        // Skip empty lines
        if (!context->line_position || *context->line_position == '\0') {
            continue;
        }
        
        // Comments starting with ##
        if (context->line_position[0] == '#' && context->line_position[1] == '#') {
            // Print comment to stderr if not quiet
            continue;
        }
        
        // Skip regular comments
        if (context->line_position[0] == '#') {
            continue;
        }
        
        // Options (@commands)
        if (context->line_position[0] == '@') {
            if (parse_option_line(context) != RESULT_SUCCESS) {
                result->result = RESULT_ERROR_INVALID_FORMAT;
                break;
            }
            continue;
        }
        
        // Check if it's a preset definition (starts with letter, no colon)
        if (isalpha(context->line_position[0])) {
            char* p = context->line_position;
            while (isalnum(*p) || *p == '_' || *p == '-') p++;
            
            // Skip spaces
            while (isspace(*p)) p++;
            
            if (*p == '\0') {
                // Preset definition
                if (parse_preset_definition(context) != RESULT_SUCCESS) {
                    result->result = RESULT_ERROR_INVALID_FORMAT;
                    break;
                }
                continue;
            }
        }
        
        // Timeline entry
        if (parse_timeline_entry(context) != RESULT_SUCCESS) {
            result->result = RESULT_ERROR_INVALID_FORMAT;
            break;
        }
    }
    
    // Copy results from context
    result->periods = context->periods;
    result->name_definitions = context->name_definitions;
    result->options = context->options;
    result->first_time = context->first_time;
    result->last_time = context->last_time;
    
    context_destroy(context);
    
    // Validate sequence
    if (result->result == RESULT_SUCCESS) {
        result->result = sequence_parser_validate_sequence(result->periods);
    }
    
    return result;
}

void sequence_parser_result_destroy(SequenceParseResult* result) {
    if (!result) return;
    
    sequence_parser_periods_destroy(result->periods);
    
    NameDefinition* nd = result->name_definitions;
    while (nd) {
        NameDefinition* next = nd->next;
        sequence_parser_preset_destroy(nd);
        nd = next;
    }
    
    memory_free(result->options.background_file);
    memory_free(result);
}

static SequenceParseContext* context_create(const char* filename) {
    SequenceParseContext* context = memory_alloc(sizeof(SequenceParseContext));
    
    context->file = fopen(filename, "r");
    if (!context->file) {
        memory_free(context);
        return NULL;
    }
    
    context->filename = memory_duplicate_string(filename);
    context->current_line = 0;
    context->line_position = NULL;
    context->name_definitions = NULL;
    context->periods = NULL;
    context->first_time = TIME_INVALID;
    context->last_time = TIME_INVALID;
    context->has_background_effect = false;
    
    // Initialize options with defaults
    context->options.background_file = NULL;
    context->options.volume_level = 100;
    context->options.sample_rate = 44100;
    context->options.gain_level_db = 12.0;
    context->options.sample_rate_user_set = false;
    
    return context;
}

static void context_destroy(SequenceParseContext* context) {
    if (!context) return;
    
    if (context->file) {
        fclose(context->file);
    }
    
    memory_free(context->filename);
    memory_free(context);
}

static ResultCode context_read_line(SequenceParseContext* context) {
    if (!fgets(context->line_buffer, MAX_LINE_LENGTH, context->file)) {
        return RESULT_ERROR_FILE_NOT_FOUND; // EOF
    }
    
    context->current_line++;
    strcpy(context->line_copy, context->line_buffer);
    
    // Remove comments and trailing whitespace
    char* comment = strchr(context->line_buffer, '#');
    if (comment) *comment = '\0';
    
    // Trim whitespace
    context->line_position = string_trim(context->line_buffer);
    
    return RESULT_SUCCESS;
}

static char* context_get_word(SequenceParseContext* context) {
    if (!context->line_position) return NULL;
    
    // Skip whitespace
    while (isspace(*context->line_position)) {
        context->line_position++;
    }
    
    if (*context->line_position == '\0') return NULL;
    
    char* start = context->line_position;
    
    // Find end of word
    while (*context->line_position && !isspace(*context->line_position)) {
        context->line_position++;
    }
    
    if (*context->line_position) {
        *context->line_position = '\0';
        context->line_position++;
    }
    
    return start;
}

NameDefinition* sequence_parser_create_preset(const char* name) {
    NameDefinition* preset = memory_alloc(sizeof(NameDefinition));
    preset->name = memory_duplicate_string(name);
    preset->next = NULL;
    
    // Initialize all voices to off
    for (int i = 0; i < MAX_CHANNELS; i++) {
        preset->voices[i].type = VOICE_TYPE_OFF;
        preset->voices[i].amplitude = 0.0;
        preset->voices[i].carrier = 0.0;
        preset->voices[i].resonance = 0.0;
        preset->voices[i].waveform = WAVE_TYPE_SINE;
    }
    
    return preset;
}

void sequence_parser_preset_destroy(NameDefinition* preset) {
    if (!preset) return;
    memory_free(preset->name);
    memory_free(preset);
}

Period* sequence_parser_create_period(time_ms_t start_time) {
    Period* period = memory_alloc(sizeof(Period));
    period->start_time = start_time;
    period->next = NULL;
    period->prev = NULL;
    period->fade_in_mode = 1;
    period->fade_out_mode = 1;
    
    // Initialize voices to off
    for (int i = 0; i < MAX_CHANNELS; i++) {
        period->start_voices[i].type = VOICE_TYPE_OFF;
        period->start_voices[i].amplitude = 0.0;
        period->start_voices[i].carrier = 0.0;
        period->start_voices[i].resonance = 0.0;
        period->start_voices[i].waveform = WAVE_TYPE_SINE;
        
        period->end_voices[i] = period->start_voices[i];
    }
    
    return period;
}

void sequence_parser_period_destroy(Period* period) {
    if (!period) return;
    memory_free(period);
}

void sequence_parser_periods_destroy(Period* periods) {
    if (!periods) return;
    
    Period* current = periods;
    Period* first = periods;
    
    do {
        Period* next = current->next;
        sequence_parser_period_destroy(current);
        current = next;
    } while (current && current != first);
}

ResultCode sequence_parser_validate_sequence(Period* periods) {
    if (!periods) {
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Check that sequence starts at 00:00:00
    if (periods->start_time != 0) {
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Check chronological order
    Period* current = periods;
    Period* first = periods;
    
    do {
        if (current->next != first && current->start_time >= current->next->start_time) {
            return RESULT_ERROR_INVALID_FORMAT;
        }
        current = current->next;
    } while (current != first);
    
    return RESULT_SUCCESS;
}

static void context_error(SequenceParseContext* context, const char* message) {
    fprintf(stderr, "Error in %s at line %d: %s\n", 
            context->filename, context->current_line, message);
    fprintf(stderr, "Line: %s\n", context->line_copy);
}

NameDefinition* sequence_parser_find_preset(NameDefinition* definitions, const char* name) {
    NameDefinition* current = definitions;
    while (current) {
        if (strcmp(current->name, name) == 0) {
            return current;
        }
        current = current->next;
    }
    return NULL;
}

static ResultCode parse_option_line(SequenceParseContext* context) {
    char* option = context_get_word(context);
    if (!option) return RESULT_ERROR_INVALID_FORMAT;
    
    // Skip the '@' character
    if (option[0] == '@') option++;
    
    if (strcmp(option, "background") == 0) {
        char* filename = context_get_word(context);
        if (!filename) {
            context_error(context, "Background filename expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        context->options.background_file = memory_duplicate_string(filename);
    }
    else if (strcmp(option, "volume") == 0) {
        char* volume_str = context_get_word(context);
        if (!volume_str) {
            context_error(context, "Volume value expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        int volume = atoi(volume_str);
        if (volume < 0 || volume > 100) {
            context_error(context, "Volume must be between 0 and 100");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        context->options.volume_level = volume;
    }
    else if (strcmp(option, "samplerate") == 0) {
        char* rate_str = context_get_word(context);
        if (!rate_str) {
            context_error(context, "Sample rate value expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        int rate = atoi(rate_str);
        if (rate <= 0) {
            context_error(context, "Invalid sample rate");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        context->options.sample_rate = rate;
        context->options.sample_rate_user_set = true;
    }
    else if (strcmp(option, "gainlevel") == 0) {
        char* gain_str = context_get_word(context);
        if (!gain_str) {
            context_error(context, "Gain level expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        if (strcmp(gain_str, "verylow") == 0) {
            context->options.gain_level_db = 20.0;
        } else if (strcmp(gain_str, "low") == 0) {
            context->options.gain_level_db = 16.0;
        } else if (strcmp(gain_str, "medium") == 0) {
            context->options.gain_level_db = 12.0;
        } else if (strcmp(gain_str, "high") == 0) {
            context->options.gain_level_db = 6.0;
        } else if (strcmp(gain_str, "veryhigh") == 0) {
            context->options.gain_level_db = 0.0;
        } else {
            context_error(context, "Invalid gain level");
            return RESULT_ERROR_INVALID_FORMAT;
        }
    }
    else {
        context_error(context, "Unknown option");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    return RESULT_SUCCESS;
}

static ResultCode parse_preset_definition(SequenceParseContext* context) {
    // Get preset name from current line
    char* preset_name = context_get_word(context);
    if (!preset_name) {
        context_error(context, "Preset name expected");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Check if preset already exists
    if (sequence_parser_find_preset(context->name_definitions, preset_name)) {
        context_error(context, "Preset already defined");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Create new preset
    NameDefinition* preset = sequence_parser_create_preset(preset_name);
    if (!preset) {
        context_error(context, "Failed to create preset");
        return RESULT_ERROR_MEMORY;
    }
    
    // Add to list
    preset->next = context->name_definitions;
    context->name_definitions = preset;
    
    int voice_count = 0;
    
    // Save current file position
    long saved_pos = ftell(context->file);
    int saved_line = context->current_line;
    
    // Read indented lines
    while (context_read_line(context) == RESULT_SUCCESS) {
        // Check if line is indented (starts with 2 spaces)
        if (strlen(context->line_buffer) < 2 || 
            context->line_buffer[0] != ' ' || 
            context->line_buffer[1] != ' ') {
            // Not indented, restore position and break
            fseek(context->file, saved_pos, SEEK_SET);
            context->current_line = saved_line;
            break;
        }
        
        // Update saved position after each successful read
        saved_pos = ftell(context->file);
        saved_line = context->current_line;
        
        // Skip the indentation
        context->line_position = context->line_buffer + 2;
        context->line_position = string_trim(context->line_position);
        
        // Skip empty lines
        if (!context->line_position || *context->line_position == '\0') {
            continue;
        }
        
        // Parse voice command
        if (voice_count >= MAX_CHANNELS) {
            context_error(context, "Too many voices in preset");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        Voice* voice = &preset->voices[voice_count];
        if (parse_voice_command(context, voice) != RESULT_SUCCESS) {
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        voice_count++;
    }
    
    if (voice_count == 0) {
        context_error(context, "Empty preset definition");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Normalize voice amplitudes
    sequence_parser_normalize_voices(preset->voices, voice_count);
    
    return RESULT_SUCCESS;
}

static ResultCode parse_timeline_entry(SequenceParseContext* context) {
    // Parse time
    char* time_str = context_get_word(context);
    if (!time_str) {
        context_error(context, "Time expected");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    time_ms_t time;
    if (sequence_parser_parse_time_string(time_str, &time) != RESULT_SUCCESS) {
        context_error(context, "Invalid time format");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Parse preset name
    char* preset_name = context_get_word(context);
    if (!preset_name) {
        context_error(context, "Preset name expected");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Find preset
    NameDefinition* preset = sequence_parser_find_preset(context->name_definitions, preset_name);
    if (!preset) {
        context_error(context, "Preset not found");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    // Update time tracking
    if (context->first_time == TIME_INVALID) {
        context->first_time = time;
    }
    context->last_time = time;
    
    // Create period
    Period* period = sequence_parser_create_period(time);
    if (!period) {
        context_error(context, "Failed to create period");
        return RESULT_ERROR_MEMORY;
    }
    
    // Copy voices from preset
    for (int i = 0; i < MAX_CHANNELS; i++) {
        period->start_voices[i] = preset->voices[i];
        period->end_voices[i] = preset->voices[i];
    }
    
    // Add to linked list
    if (!context->periods) {
        // First period
        context->periods = period;
        period->next = period;
        period->prev = period;
    } else {
        // Insert at end
        period->next = context->periods;
        period->prev = context->periods->prev;
        context->periods->prev->next = period;
        context->periods->prev = period;
    }
    
    return RESULT_SUCCESS;
}

static ResultCode parse_voice_command(SequenceParseContext* context, Voice* voice) {
    char* command = context_get_word(context);
    if (!command) {
        context_error(context, "Voice command expected");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    if (strcmp(command, "noise") == 0) {
        // Parse: noise <type> amplitude <value>
        char* type = context_get_word(context);
        if (!type) {
            context_error(context, "Noise type expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        char* amp_word = context_get_word(context);
        if (!amp_word || strcmp(amp_word, "amplitude") != 0) {
            context_error(context, "Expected 'amplitude' keyword");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        char* amp_value = context_get_word(context);
        if (!amp_value) {
            context_error(context, "Amplitude value expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        double amplitude = atof(amp_value);
        if (amplitude < 0 || amplitude > 100) {
            context_error(context, "Amplitude must be between 0 and 100");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        // Set voice type
        if (strcmp(type, "brown") == 0) {
            voice->type = VOICE_TYPE_BROWN_NOISE;
        } else if (strcmp(type, "pink") == 0) {
            voice->type = VOICE_TYPE_PINK_NOISE;
        } else if (strcmp(type, "white") == 0) {
            voice->type = VOICE_TYPE_WHITE_NOISE;
        } else {
            context_error(context, "Unknown noise type");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        voice->amplitude = AMPLITUDE_PERCENT_TO_VALUE(amplitude);
        voice->carrier = 0.0;
        voice->resonance = 0.0;
        voice->waveform = WAVE_TYPE_SINE;
        
    } else if (strcmp(command, "tone") == 0) {
        // Parse: tone <freq> <type> <value> amplitude <amp>
        char* freq_str = context_get_word(context);
        if (!freq_str) {
            context_error(context, "Frequency expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        char* type = context_get_word(context);
        if (!type) {
            context_error(context, "Tone type expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        char* value_str = context_get_word(context);
        if (!value_str) {
            context_error(context, "Tone value expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        char* amp_word = context_get_word(context);
        if (!amp_word || strcmp(amp_word, "amplitude") != 0) {
            context_error(context, "Expected 'amplitude' keyword");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        char* amp_value = context_get_word(context);
        if (!amp_value) {
            context_error(context, "Amplitude value expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        double freq = atof(freq_str);
        double value = atof(value_str);
        double amplitude = atof(amp_value);
        
        if (freq < 0 || value < 0 || amplitude < 0 || amplitude > 100) {
            context_error(context, "Invalid tone parameters");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        // Set voice type
        if (strcmp(type, "binaural") == 0) {
            voice->type = VOICE_TYPE_BINAURAL;
        } else if (strcmp(type, "monaural") == 0) {
            voice->type = VOICE_TYPE_MONAURAL;
        } else if (strcmp(type, "isochronic") == 0) {
            voice->type = VOICE_TYPE_ISOCHRONIC;
        } else {
            context_error(context, "Unknown tone type");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        voice->amplitude = AMPLITUDE_PERCENT_TO_VALUE(amplitude);
        voice->carrier = freq;
        voice->resonance = value;
        voice->waveform = WAVE_TYPE_SINE;
        
    } else if (strcmp(command, "background") == 0) {
        // Parse: background amplitude <value>
        char* amp_word = context_get_word(context);
        if (!amp_word || strcmp(amp_word, "amplitude") != 0) {
            context_error(context, "Expected 'amplitude' keyword");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        char* amp_value = context_get_word(context);
        if (!amp_value) {
            context_error(context, "Amplitude value expected");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        double amplitude = atof(amp_value);
        if (amplitude < 0 || amplitude > 100) {
            context_error(context, "Amplitude must be between 0 and 100");
            return RESULT_ERROR_INVALID_FORMAT;
        }
        
        voice->type = VOICE_TYPE_BACKGROUND;
        voice->amplitude = AMPLITUDE_PERCENT_TO_VALUE(amplitude);
        voice->carrier = 0.0;
        voice->resonance = 0.0;
        voice->waveform = WAVE_TYPE_SINE;
        
    } else {
        context_error(context, "Unknown voice command");
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    return RESULT_SUCCESS;
}

ResultCode sequence_parser_parse_time_string(const char* time_str, time_ms_t* result) {
    return time_parse_string(time_str, result);
}

const char* sequence_parser_get_voice_type_name(VoiceType type) {
    switch (type) {
        case VOICE_TYPE_OFF: return "off";
        case VOICE_TYPE_BINAURAL: return "binaural";
        case VOICE_TYPE_MONAURAL: return "monaural";
        case VOICE_TYPE_ISOCHRONIC: return "isochronic";
        case VOICE_TYPE_WHITE_NOISE: return "white_noise";
        case VOICE_TYPE_PINK_NOISE: return "pink_noise";
        case VOICE_TYPE_BROWN_NOISE: return "brown_noise";
        case VOICE_TYPE_SPIN_WHITE: return "spin_white";
        case VOICE_TYPE_SPIN_PINK: return "spin_pink";
        case VOICE_TYPE_SPIN_BROWN: return "spin_brown";
        case VOICE_TYPE_BACKGROUND: return "background";
        case VOICE_TYPE_EFFECT_SPIN: return "effect_spin";
        case VOICE_TYPE_EFFECT_PULSE: return "effect_pulse";
        default: return "unknown";
    }
}

const char* sequence_parser_get_wave_type_name(WaveType type) {
    switch (type) {
        case WAVE_TYPE_SINE: return "sine";
        case WAVE_TYPE_SQUARE: return "square";
        case WAVE_TYPE_TRIANGLE: return "triangle";
        case WAVE_TYPE_SAWTOOTH: return "sawtooth";
        default: return "unknown";
    }
}

SequenceParseResult* sequence_parser_parse_files(const char* filenames[], int count) {
    // TODO: Implement multiple file parsing
    return NULL;
}

ResultCode sequence_parser_parse_voice_line(const char* line, Voice* voices, int* channel_count) {
    // TODO: Implement voice line parsing
    return RESULT_ERROR_INVALID_FORMAT;
}

ResultCode sequence_parser_parse_timeline(const char* line, time_ms_t* time, char** preset_name) {
    // TODO: Implement timeline parsing
    return RESULT_ERROR_INVALID_FORMAT;
}

ResultCode sequence_parser_normalize_voices(Voice* voices, int channel_count) {
    double total_amplitude = 0.0;
    
    // Calculate total amplitude (excluding effects)
    for (int i = 0; i < channel_count; i++) {
        if (voices[i].type != VOICE_TYPE_OFF && 
            voices[i].type != VOICE_TYPE_EFFECT_SPIN &&
            voices[i].type != VOICE_TYPE_EFFECT_PULSE) {
            total_amplitude += AMPLITUDE_VALUE_TO_PERCENT(voices[i].amplitude);
        }
    }
    
    // Normalize if over 100%
    if (total_amplitude > 100.0) {
        double factor = 100.0 / total_amplitude;
        for (int i = 0; i < channel_count; i++) {
            if (voices[i].type != VOICE_TYPE_OFF && 
                voices[i].type != VOICE_TYPE_EFFECT_SPIN &&
                voices[i].type != VOICE_TYPE_EFFECT_PULSE) {
                voices[i].amplitude *= factor;
            }
        }
    }
    
    return RESULT_SUCCESS;
}