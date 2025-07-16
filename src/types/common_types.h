#ifndef COMMON_TYPES_H
#define COMMON_TYPES_H

#include <stdbool.h>
#include <stdint.h>

#define MAX_CHANNELS 16
#define MAX_FILENAME_LENGTH 256
#define MAX_LINE_LENGTH 4096
#define MAX_PRESET_NAME_LENGTH 64

typedef unsigned char byte;
typedef uint64_t timestamp_ms;

typedef enum {
    RESULT_SUCCESS = 0,
    RESULT_ERROR_MEMORY,
    RESULT_ERROR_FILE_NOT_FOUND,
    RESULT_ERROR_INVALID_FORMAT,
    RESULT_ERROR_INVALID_ARGUMENT,
    RESULT_ERROR_AUDIO_DEVICE
} ResultCode;

typedef struct {
    ResultCode code;
    char message[256];
    int line_number;
    char filename[MAX_FILENAME_LENGTH];
} ErrorInfo;

#endif // COMMON_TYPES_H