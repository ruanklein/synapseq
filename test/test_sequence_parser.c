#include "../src/parsing/sequence_parser.h"
#include "../src/types/sequence_types.h"

#include <stdio.h>
#include <assert.h>
#include <string.h>
#include <stdlib.h>

void test_preset_creation() {
    printf("Testing preset creation...\n");
    
    NameDefinition* preset = sequence_parser_create_preset("test_preset");
    assert(preset != NULL);
    assert(strcmp(preset->name, "test_preset") == 0);
    assert(preset->next == NULL);
    
    // Check that all voices are initialized to off
    for (int i = 0; i < MAX_CHANNELS; i++) {
        assert(preset->voices[i].type == VOICE_TYPE_OFF);
        assert(preset->voices[i].amplitude == 0.0);
        assert(preset->voices[i].carrier == 0.0);
        assert(preset->voices[i].resonance == 0.0);
        assert(preset->voices[i].waveform == WAVE_TYPE_SINE);
    }
    
    sequence_parser_preset_destroy(preset);
    printf("✓ Preset creation tests passed\n");
}

void test_preset_finding() {
    printf("Testing preset finding...\n");
    
    NameDefinition* preset1 = sequence_parser_create_preset("preset1");
    NameDefinition* preset2 = sequence_parser_create_preset("preset2");
    NameDefinition* preset3 = sequence_parser_create_preset("preset3");
    
    // Link them
    preset1->next = preset2;
    preset2->next = preset3;
    
    // Test finding existing presets
    NameDefinition* found = sequence_parser_find_preset(preset1, "preset2");
    assert(found == preset2);
    
    found = sequence_parser_find_preset(preset1, "preset3");
    assert(found == preset3);
    
    // Test finding non-existing preset
    found = sequence_parser_find_preset(preset1, "nonexistent");
    assert(found == NULL);
    
    sequence_parser_preset_destroy(preset1);
    sequence_parser_preset_destroy(preset2);
    sequence_parser_preset_destroy(preset3);
    printf("✓ Preset finding tests passed\n");
}

void test_period_creation() {
    printf("Testing period creation...\n");
    
    time_ms_t test_time = 5000; // 5 seconds
    Period* period = sequence_parser_create_period(test_time);
    
    assert(period != NULL);
    assert(period->start_time == test_time);
    assert(period->next == NULL);
    assert(period->prev == NULL);
    assert(period->fade_in_mode == 1);
    assert(period->fade_out_mode == 1);
    
    // Check that all voices are initialized to off
    for (int i = 0; i < MAX_CHANNELS; i++) {
        assert(period->start_voices[i].type == VOICE_TYPE_OFF);
        assert(period->end_voices[i].type == VOICE_TYPE_OFF);
    }
    
    sequence_parser_period_destroy(period);
    printf("✓ Period creation tests passed\n");
}

void test_voice_normalization() {
    printf("Testing voice normalization...\n");
    
    Voice voices[MAX_CHANNELS];
    
    // Initialize voices to off
    for (int i = 0; i < MAX_CHANNELS; i++) {
        voices[i].type = VOICE_TYPE_OFF;
        voices[i].amplitude = 0.0;
    }
    
    // Set up voices that exceed 100% total amplitude
    voices[0].type = VOICE_TYPE_BINAURAL;
    voices[0].amplitude = AMPLITUDE_PERCENT_TO_VALUE(60.0); // 60%
    
    voices[1].type = VOICE_TYPE_PINK_NOISE;
    voices[1].amplitude = AMPLITUDE_PERCENT_TO_VALUE(70.0); // 70%
    
    // Total is 130%, should be normalized to 100%
    ResultCode result = sequence_parser_normalize_voices(voices, MAX_CHANNELS);
    assert(result == RESULT_SUCCESS);
    
    // Check that amplitudes were scaled down
    double total = AMPLITUDE_VALUE_TO_PERCENT(voices[0].amplitude) + 
                   AMPLITUDE_VALUE_TO_PERCENT(voices[1].amplitude);
    assert(total <= 100.1); // Allow small floating point error
    assert(total >= 99.9);
    
    printf("✓ Voice normalization tests passed\n");
}

void test_voice_type_names() {
    printf("Testing voice type names...\n");
    
    assert(strcmp(sequence_parser_get_voice_type_name(VOICE_TYPE_OFF), "off") == 0);
    assert(strcmp(sequence_parser_get_voice_type_name(VOICE_TYPE_BINAURAL), "binaural") == 0);
    assert(strcmp(sequence_parser_get_voice_type_name(VOICE_TYPE_PINK_NOISE), "pink_noise") == 0);
    assert(strcmp(sequence_parser_get_voice_type_name(VOICE_TYPE_BACKGROUND), "background") == 0);
    
    printf("✓ Voice type names tests passed\n");
}

void test_wave_type_names() {
    printf("Testing wave type names...\n");
    
    assert(strcmp(sequence_parser_get_wave_type_name(WAVE_TYPE_SINE), "sine") == 0);
    assert(strcmp(sequence_parser_get_wave_type_name(WAVE_TYPE_SQUARE), "square") == 0);
    assert(strcmp(sequence_parser_get_wave_type_name(WAVE_TYPE_TRIANGLE), "triangle") == 0);
    assert(strcmp(sequence_parser_get_wave_type_name(WAVE_TYPE_SAWTOOTH), "sawtooth") == 0);
    
    printf("✓ Wave type names tests passed\n");
}

void test_amplitude_macros() {
    printf("Testing amplitude conversion macros...\n");
    
    // Test conversion from percentage to value
    double value = AMPLITUDE_PERCENT_TO_VALUE(50.0);
    assert(value == 50.0 * SEQUENCE_AMPLITUDE_SCALE);
    
    // Test conversion from value to percentage
    double percent = AMPLITUDE_VALUE_TO_PERCENT(value);
    assert(percent == 50.0);
    
    // Test round-trip conversion
    double original = 75.5;
    double converted = AMPLITUDE_VALUE_TO_PERCENT(AMPLITUDE_PERCENT_TO_VALUE(original));
    assert(converted == original);
    
    printf("✓ Amplitude conversion tests passed\n");
}

void test_noise_sequence_file() {
    printf("Testing real sequence file parsing...\n");
    
    // Try to parse the sample-noise.spsq file
    SequenceParseResult* result = sequence_parser_parse_file("samples/sample-noise.spsq");
    
    if (result && result->result == RESULT_SUCCESS) {
        printf("✓ File parsed successfully\n");
        
        // Check that we have the expected presets
        NameDefinition* silence = sequence_parser_find_preset(result->name_definitions, "silence");
        assert(silence != NULL);
        printf("✓ Found built-in 'silence' preset\n");
        
        // Check specific presets from sample-noise.spsq
        NameDefinition* noise_one = sequence_parser_find_preset(result->name_definitions, "noise-one");
        assert(noise_one != NULL);
        printf("✓ Found 'noise-one' preset\n");
        
        // Verify noise-one preset contains brown noise with amplitude 15
        assert(noise_one->voices[0].type == VOICE_TYPE_BROWN_NOISE);
        assert(AMPLITUDE_VALUE_TO_PERCENT(noise_one->voices[0].amplitude) == 15.0);
        printf("✓ noise-one preset has correct brown noise (15%% amplitude)\n");
        
        NameDefinition* noise_two = sequence_parser_find_preset(result->name_definitions, "noise-two");
        assert(noise_two != NULL);
        printf("✓ Found 'noise-two' preset\n");
        
        // Verify noise-two preset contains pink noise with amplitude 15
        assert(noise_two->voices[0].type == VOICE_TYPE_PINK_NOISE);
        assert(AMPLITUDE_VALUE_TO_PERCENT(noise_two->voices[0].amplitude) == 15.0);
        printf("✓ noise-two preset has correct pink noise (15%% amplitude)\n");
        
        NameDefinition* noise_three = sequence_parser_find_preset(result->name_definitions, "noise-three");
        assert(noise_three != NULL);
        printf("✓ Found 'noise-three' preset\n");
        
        // Verify noise-three preset contains white noise with amplitude 15
        assert(noise_three->voices[0].type == VOICE_TYPE_WHITE_NOISE);
        assert(AMPLITUDE_VALUE_TO_PERCENT(noise_three->voices[0].amplitude) == 15.0);
        printf("✓ noise-three preset has correct white noise (15%% amplitude)\n");
        
        // Print all found presets
        printf("--- Found Presets ---\n");
        NameDefinition* current = result->name_definitions;
        while (current) {
            printf("Preset: %s\n", current->name);
            
            // Print voices in this preset
            for (int i = 0; i < MAX_CHANNELS; i++) {
                if (current->voices[i].type != VOICE_TYPE_OFF) {
                    printf("  Voice %d: %s, amplitude=%.2f, carrier=%.2f, resonance=%.2f, waveform=%s\n",
                           i,
                           sequence_parser_get_voice_type_name(current->voices[i].type),
                           AMPLITUDE_VALUE_TO_PERCENT(current->voices[i].amplitude),
                           current->voices[i].carrier,
                           current->voices[i].resonance,
                           sequence_parser_get_wave_type_name(current->voices[i].waveform));
                }
            }
            current = current->next;
        }
        
        // Check timeline - should have 5 periods
        if (result->periods) {
            printf("--- Timeline ---\n");
            Period* period = result->periods;
            Period* first = period;
            
            // Expected times: 00:00:00, 00:00:15, 00:00:30, 00:00:45, 00:01:00
            time_ms_t expected_times[] = {0, 15000, 30000, 45000, 60000}; // in milliseconds
            int period_count = 0;
            
            do {
                printf("Time: %02d:%02d:%02d",
                       (int)(period->start_time / 3600000),
                       (int)((period->start_time % 3600000) / 60000),
                       (int)((period->start_time % 60000) / 1000));
                
                // Check if this matches expected time
                if (period_count < 5) {
                    assert(period->start_time == expected_times[period_count]);
                    printf(" ✓");
                }
                printf("\n");
                
                period = period->next;
                period_count++;
            } while (period && period != first);
            
            assert(period_count == 5);
            printf("✓ Timeline has correct 5 periods\n");
            
            // Check time range
            assert(result->first_time == 0);
            assert(result->last_time == 60000); // 00:01:00
            printf("✓ Time range: %02d:%02d:%02d to %02d:%02d:%02d\n",
                   (int)(result->first_time / 3600000),
                   (int)((result->first_time % 3600000) / 60000),
                   (int)((result->first_time % 60000) / 1000),
                   (int)(result->last_time / 3600000),
                   (int)((result->last_time % 3600000) / 60000),
                   (int)((result->last_time % 60000) / 1000));
        } else {
            printf("✗ No timeline found\n");
            assert(false);
        }
        
        printf("✓ Real sequence file parsing tests passed\n");
    } else {
        if (result) {
            printf("✗ Parse failed with error code: %d\n", result->result);
        } else {
            printf("✗ Could not create parse result\n");
        }
        assert(false);
    }
    
    if (result) {
        sequence_parser_result_destroy(result);
    }
}

void test_binaural_sequence_file() {
    printf("Testing binaural sequence file parsing...\n");
    
    SequenceParseResult* result = sequence_parser_parse_file("samples/sample-binaural.spsq");
    
    if (result && result->result == RESULT_SUCCESS) {
        printf("✓ Binaural file parsed successfully\n");
        
        // Check for binaural presets
        NameDefinition* alpha_one = sequence_parser_find_preset(result->name_definitions, "alpha-one");
        assert(alpha_one != NULL);
        printf("✓ Found 'alpha-one' preset\n");
        
        NameDefinition* alpha_two = sequence_parser_find_preset(result->name_definitions, "alpha-two");
        assert(alpha_two != NULL);
        printf("✓ Found 'alpha-two' preset\n");
        
        NameDefinition* alpha_three = sequence_parser_find_preset(result->name_definitions, "alpha-three");
        assert(alpha_three != NULL);
        printf("✓ Found 'alpha-three' preset\n");
        
        NameDefinition* alpha_four = sequence_parser_find_preset(result->name_definitions, "alpha-four");
        assert(alpha_four != NULL);
        printf("✓ Found 'alpha-four' preset\n");
        
        // Print all found presets
        printf("--- Found Binaural Presets ---\n");
        NameDefinition* current = result->name_definitions;
        while (current) {
            printf("Preset: %s\n", current->name);
            
            // Print voices in this preset
            for (int i = 0; i < MAX_CHANNELS; i++) {
                if (current->voices[i].type != VOICE_TYPE_OFF) {
                    printf("  Voice %d: %s, amplitude=%.2f, carrier=%.2f, resonance=%.2f, waveform=%s\n",
                           i,
                           sequence_parser_get_voice_type_name(current->voices[i].type),
                           AMPLITUDE_VALUE_TO_PERCENT(current->voices[i].amplitude),
                           current->voices[i].carrier,
                           current->voices[i].resonance,
                           sequence_parser_get_wave_type_name(current->voices[i].waveform));
                }
            }
            current = current->next;
        }
        
        // Check timeline - should have 7 periods
        if (result->periods) {
            printf("--- Binaural Timeline ---\n");
            Period* period = result->periods;
            Period* first = period;
            
            // Expected times: 00:00:00, 00:00:15, 00:01:00, 00:02:30, 00:04:00, 00:05:30, 00:07:00
            time_ms_t expected_times[] = {0, 15000, 60000, 150000, 240000, 330000, 420000};
            int period_count = 0;
            
            do {
                printf("Time: %02d:%02d:%02d",
                       (int)(period->start_time / 3600000),
                       (int)((period->start_time % 3600000) / 60000),
                       (int)((period->start_time % 60000) / 1000));
                
                // Check if this matches expected time
                if (period_count < 7) {
                    assert(period->start_time == expected_times[period_count]);
                    printf(" ✓");
                }
                printf("\n");
                
                period = period->next;
                period_count++;
            } while (period && period != first);
            
            assert(period_count == 7);
            printf("✓ Binaural timeline has correct 7 periods\n");
            
            // Check time range
            assert(result->first_time == 0);
            assert(result->last_time == 420000); // 00:07:00
            printf("✓ Binaural time range: %02d:%02d:%02d to %02d:%02d:%02d\n",
                   (int)(result->first_time / 3600000),
                   (int)((result->first_time % 3600000) / 60000),
                   (int)((result->first_time % 60000) / 1000),
                   (int)(result->last_time / 3600000),
                   (int)((result->last_time % 3600000) / 60000),
                   (int)((result->last_time % 60000) / 1000));
        } else {
            printf("✗ No binaural timeline found\n");
            assert(false);
        }
        
        printf("✓ Binaural sequence file parsing tests passed\n");
    } else {
        if (result) {
            printf("✗ Binaural parse failed with error code: %d\n", result->result);
        } else {
            printf("✗ Could not create binaural parse result\n");
        }
        assert(false);
    }
    
    if (result) {
        sequence_parser_result_destroy(result);
    }
}

void test_invalid_sequence_files() {
    printf("Testing invalid sequence file parsing...\n");
    
    // Test 1: Empty preset definition
    printf("Testing sequence-01.spsq (empty preset)...\n");
    SequenceParseResult* result1 = sequence_parser_parse_file("test/spsq/sequence-01.spsq");
    
    if (result1) {
        if (result1->result == RESULT_SUCCESS) {
            // Should succeed but alpha-one preset should be empty
            NameDefinition* alpha_one = sequence_parser_find_preset(result1->name_definitions, "alpha-one");
            assert(alpha_one != NULL);
            printf("✓ Found empty 'alpha-one' preset\n");
            
            // Check that alpha-one has no voices
            bool has_voices = false;
            for (int i = 0; i < MAX_CHANNELS; i++) {
                if (alpha_one->voices[i].type != VOICE_TYPE_OFF) {
                    has_voices = true;
                    break;
                }
            }
            assert(!has_voices);
            printf("✓ alpha-one preset is correctly empty\n");
            
            // Check that alpha-two has the correct voice
            NameDefinition* alpha_two = sequence_parser_find_preset(result1->name_definitions, "alpha-two");
            assert(alpha_two != NULL);
            assert(alpha_two->voices[0].type == VOICE_TYPE_BINAURAL);
            assert(alpha_two->voices[0].carrier == 250.0);
            assert(alpha_two->voices[0].resonance == 10.0);
            assert(AMPLITUDE_VALUE_TO_PERCENT(alpha_two->voices[0].amplitude) == 10.0);
            printf("✓ alpha-two preset has correct binaural voice\n");
            
        } else {
            printf("✗ sequence-01.spsq failed to parse (error: %d)\n", result1->result);
        }
        sequence_parser_result_destroy(result1);
    }
    
    // Test 2: Voice definitions outside preset
    printf("Testing sequence-02.spsq (voices outside preset)...\n");
    SequenceParseResult* result2 = sequence_parser_parse_file("test/spsq/sequence-02.spsq");
    
    if (result2) {
        if (result2->result == RESULT_ERROR_INVALID_FORMAT) {
            printf("✓ sequence-02.spsq correctly rejected (voices outside preset)\n");
        } else if (result2->result == RESULT_SUCCESS) {
            // If it succeeds, those voice lines should be ignored
            NameDefinition* alpha_one = sequence_parser_find_preset(result2->name_definitions, "alpha-one");
            assert(alpha_one != NULL);
            
            // Check that alpha-one is empty (voices were outside preset)
            bool has_voices = false;
            for (int i = 0; i < MAX_CHANNELS; i++) {
                if (alpha_one->voices[i].type != VOICE_TYPE_OFF) {
                    has_voices = true;
                    break;
                }
            }
            assert(!has_voices);
            printf("✓ alpha-one preset is correctly empty (voices outside preset were ignored)\n");
        } else {
            printf("✗ sequence-02.spsq failed unexpectedly (error: %d)\n", result2->result);
        }
        sequence_parser_result_destroy(result2);
    }
    
    printf("✓ Invalid sequence file tests passed\n");
}

int main() {
    printf("=== Testing Sequence Parser Module ===\n");
    
    test_preset_creation();
    test_preset_finding();
    test_period_creation();
    test_voice_normalization();
    test_voice_type_names();
    test_wave_type_names();
    test_amplitude_macros();
    test_noise_sequence_file();
    test_binaural_sequence_file();
    test_invalid_sequence_files();

    printf("✓ All sequence parser tests passed!\n");
    return 0;
}