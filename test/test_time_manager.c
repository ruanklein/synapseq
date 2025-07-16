#include "../src/core/time_manager.h"

#include <stdio.h>
#include <assert.h>
#include <string.h>

void test_time_parsing() {
    printf("Testing time parsing...\n");
    
    time_ms_t result;
    
    assert(time_parse_string("00:00:00", &result) == RESULT_SUCCESS);
    assert(result == 0);
    
    assert(time_parse_string("00:00:15", &result) == RESULT_SUCCESS);
    assert(result == 15 * TIME_MS_PER_SECOND);
    
    assert(time_parse_string("00:00:45", &result) == RESULT_SUCCESS);
    assert(result == 45 * TIME_MS_PER_SECOND);
    
    assert(time_parse_string("01:00:00", &result) == RESULT_SUCCESS);
    assert(result == 1 * TIME_MS_PER_HOUR);
    
    assert(time_parse_string("12:30:45", &result) == RESULT_SUCCESS);
    assert(result == 12 * TIME_MS_PER_HOUR + 30 * TIME_MS_PER_MINUTE + 45 * TIME_MS_PER_SECOND);
    
    assert(time_parse_string("23:59:59", &result) == RESULT_SUCCESS);
    assert(result == 23 * TIME_MS_PER_HOUR + 59 * TIME_MS_PER_MINUTE + 59 * TIME_MS_PER_SECOND);
    
    assert(time_parse_string("08:15", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:30", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:30:45:99", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12", &result) == RESULT_ERROR_INVALID_FORMAT);
    
    assert(time_parse_string("12-30-45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12.30.45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12 30 45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:30-45", &result) == RESULT_ERROR_INVALID_FORMAT);
    
    assert(time_parse_string("1a:30:45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:3b:45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:30:4c", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("invalid", &result) == RESULT_ERROR_INVALID_FORMAT);
    
    assert(time_parse_string("24:00:00", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:60:00", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:30:60", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("25:70:80", &result) == RESULT_ERROR_INVALID_FORMAT);
    
    assert(time_parse_string("12 :30:45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12: 30:45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:30 :45", &result) == RESULT_ERROR_INVALID_FORMAT);
    assert(time_parse_string("12:30: 45", &result) == RESULT_ERROR_INVALID_FORMAT);
    
    printf("✓ Time parsing tests passed\n");
}

void test_time_validation() {
    printf("Testing time validation...\n");
    
    assert(time_is_valid_string("00:00:00") == true);
    assert(time_is_valid_string("12:30:45") == true);
    assert(time_is_valid_string("23:59:59") == true);
    
    assert(time_is_valid_string("08:15") == false);
    assert(time_is_valid_string("12:30") == false);
    assert(time_is_valid_string("invalid") == false);
    assert(time_is_valid_string("25:00:00") == false);
    assert(time_is_valid_string("12:60:00") == false);
    assert(time_is_valid_string("12:30:60") == false);
    
    printf("✓ Time validation tests passed\n");
}

void test_sequence_validation() {
    printf("Testing sequence validation...\n");
    
    time_ms_t sequence[] = {
        0,
        15 * TIME_MS_PER_SECOND,
        45 * TIME_MS_PER_SECOND,
        1 * TIME_MS_PER_HOUR
    };
    
    int count = sizeof(sequence) / sizeof(sequence[0]);
    
    assert(time_sequence_is_valid(sequence, count) == true);
    
    time_ms_t invalid_sequence[] = {
        1 * TIME_MS_PER_HOUR,
        15 * TIME_MS_PER_SECOND,
        45 * TIME_MS_PER_SECOND
    };
    
    assert(time_sequence_is_valid(invalid_sequence, 3) == false);
    
    printf("✓ Sequence validation tests passed\n");
}

void test_duration_calculation() {
    printf("Testing duration calculation...\n");
    
    time_ms_t t1 = 0;
    time_ms_t t2 = 15 * TIME_MS_PER_SECOND;
    time_ms_t t3 = 45 * TIME_MS_PER_SECOND;
    time_ms_t t4 = 1 * TIME_MS_PER_HOUR;
    
    assert(time_calculate_duration(t1, t2) == 15 * TIME_MS_PER_SECOND);
    assert(time_calculate_duration(t2, t3) == 30 * TIME_MS_PER_SECOND);
    
    time_ms_t expected = 59 * TIME_MS_PER_MINUTE + 15 * TIME_MS_PER_SECOND;
    assert(time_calculate_duration(t3, t4) == expected);
    
    printf("✓ Duration calculation tests passed\n");
}

void test_period_finding() {
    printf("Testing period finding...\n");
    
    time_ms_t sequence[] = {
        0,
        15 * TIME_MS_PER_SECOND,
        45 * TIME_MS_PER_SECOND,
        1 * TIME_MS_PER_HOUR
    };
    
    int count = sizeof(sequence) / sizeof(sequence[0]);
    
    assert(time_find_period_index(sequence, count, 0) == 0);
    assert(time_find_period_index(sequence, count, 10 * TIME_MS_PER_SECOND) == 0);
    assert(time_find_period_index(sequence, count, 20 * TIME_MS_PER_SECOND) == 1);
    assert(time_find_period_index(sequence, count, 30 * TIME_MS_PER_MINUTE) == 2);
    
    printf("✓ Period finding tests passed\n");
}

void test_time_formatting() {
    printf("Testing time formatting...\n");
    
    char buffer[32];
    
    assert(time_format_string(0, buffer, sizeof(buffer)) == RESULT_SUCCESS);
    assert(strcmp(buffer, "00:00:00") == 0);
    
    assert(time_format_string(15 * TIME_MS_PER_SECOND, buffer, sizeof(buffer)) == RESULT_SUCCESS);
    assert(strcmp(buffer, "00:00:15") == 0);
    
    assert(time_format_string(45 * TIME_MS_PER_SECOND, buffer, sizeof(buffer)) == RESULT_SUCCESS);
    assert(strcmp(buffer, "00:00:45") == 0);
    
    assert(time_format_string(1 * TIME_MS_PER_HOUR, buffer, sizeof(buffer)) == RESULT_SUCCESS);
    assert(strcmp(buffer, "01:00:00") == 0);
    
    printf("✓ Time formatting tests passed\n");
}

int main() {
    printf("=== Testing Time Manager ===\n");
    
    test_time_parsing();
    test_time_validation();
    test_sequence_validation();
    test_duration_calculation();
    test_period_finding();
    test_time_formatting();
    
    printf("\n✓ All time manager tests passed!\n");
    
    return 0;
}