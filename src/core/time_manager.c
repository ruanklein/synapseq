#include "time_manager.h"
#include "../utils/memory_utils.h"
#include "../utils/string_utils.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

ResultCode time_parse_string(const char* time_str, time_ms_t* result) {
    if (!time_str || !result) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    *result = TIME_INVALID;
    
    char* str_copy = memory_duplicate_string(time_str);
    char* trimmed = string_trim(str_copy);
    
    int len = strlen(trimmed);
    if (len != 8) {
        memory_free(str_copy);
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    if (trimmed[2] != ':' || trimmed[5] != ':') {
        memory_free(str_copy);
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    for (int i = 0; i < 8; i++) {
        if (i == 2 || i == 5) continue;
        if (trimmed[i] < '0' || trimmed[i] > '9') {
            memory_free(str_copy);
            return RESULT_ERROR_INVALID_FORMAT;
        }
    }
    
    int hours = 0, minutes = 0, seconds = 0;
    int parsed = sscanf(trimmed, "%2d:%2d:%2d", &hours, &minutes, &seconds);
    
    memory_free(str_copy);
    
    if (parsed != 3) {
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    if (hours < 0 || hours >= 24 || 
        minutes < 0 || minutes >= 60 || 
        seconds < 0 || seconds >= 60) {
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    TimeComponents components = {hours, minutes, seconds};
    return time_components_to_ms(&components, result);
}

ResultCode time_format_string(time_ms_t time_ms, char* buffer, size_t buffer_size) {
    if (!buffer || buffer_size < 9) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    if (!time_is_valid(time_ms)) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    TimeComponents components;
    ResultCode result = time_ms_to_components(time_ms, &components);
    if (result != RESULT_SUCCESS) {
        return result;
    }
    
    int written = snprintf(buffer, buffer_size, "%02d:%02d:%02d", 
                          components.hours, components.minutes, components.seconds);
    
    if (written < 0 || written >= (int)buffer_size) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    return RESULT_SUCCESS;
}

ResultCode time_components_to_ms(const TimeComponents* components, time_ms_t* result) {
    if (!components || !result) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    if (components->hours < 0 || components->hours >= 24 ||
        components->minutes < 0 || components->minutes >= 60 ||
        components->seconds < 0 || components->seconds >= 60) {
        return RESULT_ERROR_INVALID_FORMAT;
    }
    
    *result = (time_ms_t)(components->hours * TIME_MS_PER_HOUR +
                          components->minutes * TIME_MS_PER_MINUTE +
                          components->seconds * TIME_MS_PER_SECOND);
    
    return RESULT_SUCCESS;
}

ResultCode time_ms_to_components(time_ms_t time_ms, TimeComponents* components) {
    if (!components || !time_is_valid(time_ms)) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    time_ms %= TIME_MS_PER_DAY;
    
    components->hours = time_ms / TIME_MS_PER_HOUR;
    time_ms %= TIME_MS_PER_HOUR;
    
    components->minutes = time_ms / TIME_MS_PER_MINUTE;
    time_ms %= TIME_MS_PER_MINUTE;
    
    components->seconds = time_ms / TIME_MS_PER_SECOND;
    
    return RESULT_SUCCESS;
}

time_ms_t time_calculate_duration(time_ms_t start, time_ms_t end) {
    if (!time_is_valid(start) || !time_is_valid(end)) {
        return TIME_INVALID;
    }
    
    if (end >= start) {
        return end - start;
    }
    
    return TIME_INVALID;
}

time_ms_t time_calculate_midpoint(time_ms_t start, time_ms_t end) {
    if (!time_is_valid(start) || !time_is_valid(end) || end < start) {
        return TIME_INVALID;
    }
    
    return start + (end - start) / 2;
}

bool time_is_chronological(time_ms_t start, time_ms_t end) {
    return time_is_valid(start) && time_is_valid(end) && end >= start;
}

bool time_is_valid(time_ms_t time_ms) {
    return time_ms != TIME_INVALID && time_ms < TIME_MS_PER_DAY;
}

bool time_is_valid_string(const char* time_str) {
    if (!time_str) return false;
    
    time_ms_t dummy;
    return time_parse_string(time_str, &dummy) == RESULT_SUCCESS;
}

bool time_sequence_is_valid(const time_ms_t* times, int count) {
    if (!times || count < 2) return false;
    
    for (int i = 1; i < count; i++) {
        if (!time_is_chronological(times[i-1], times[i])) {
            return false;
        }
    }
    
    return true;
}

int time_find_period_index(const time_ms_t* times, int count, time_ms_t current_time) {
    if (!times || count < 2 || !time_is_valid(current_time)) {
        return -1;
    }
    
    if (current_time < times[0]) {
        return 0;
    }
    
    if (current_time >= times[count-1]) {
        return count - 2;
    }
    
    for (int i = 0; i < count - 1; i++) {
        if (current_time >= times[i] && current_time < times[i+1]) {
            return i;
        }
    }
    
    return -1;
}