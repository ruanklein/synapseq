#ifndef TIME_MANAGER_H
#define TIME_MANAGER_H

#include "../types/time_types.h"
#include "../types/common_types.h"

#include <stdbool.h>
#include <stddef.h>

ResultCode time_parse_string(const char* time_str, time_ms_t* result);
ResultCode time_format_string(time_ms_t time_ms, char* buffer, size_t buffer_size);
ResultCode time_components_to_ms(const TimeComponents* components, time_ms_t* result);
ResultCode time_ms_to_components(time_ms_t time_ms, TimeComponents* components);

time_ms_t time_calculate_duration(time_ms_t start, time_ms_t end);
time_ms_t time_calculate_midpoint(time_ms_t start, time_ms_t end);
bool time_is_chronological(time_ms_t start, time_ms_t end);

bool time_is_valid(time_ms_t time_ms);
bool time_is_valid_string(const char* time_str);

bool time_sequence_is_valid(const time_ms_t* times, int count);
int time_find_period_index(const time_ms_t* times, int count, time_ms_t current_time);

#endif // TIME_MANAGER_H