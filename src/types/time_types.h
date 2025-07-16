#ifndef TIME_TYPES_H
#define TIME_TYPES_H

#include "common_types.h"

typedef uint32_t time_ms_t;

typedef struct {
    int hours;      // 0-23
    int minutes;    // 0-59
    int seconds;    // 0-59
} TimeComponents;

#define TIME_MS_PER_SECOND 1000
#define TIME_MS_PER_MINUTE (60 * TIME_MS_PER_SECOND)
#define TIME_MS_PER_HOUR   (60 * TIME_MS_PER_MINUTE)
#define TIME_MS_PER_DAY    (24 * TIME_MS_PER_HOUR)

#define TIME_INVALID       0xFFFFFFFF

#endif // TIME_TYPES_H