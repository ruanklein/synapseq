#ifndef AUDIO_OUTPUT_H
#define AUDIO_OUTPUT_H

#include "../types/audio_types.h"
#include "../types/common_types.h"

#include <stdio.h>

typedef struct AudioOutput AudioOutput;

typedef struct {
    ResultCode (*write_samples)(AudioOutput* output, const int16_t* samples, int sample_count);
    ResultCode (*finalize)(AudioOutput* output);
    void (*destroy)(AudioOutput* output);
} AudioOutputInterface;

struct AudioOutput {
    AudioOutputInterface interface;
    AudioConfig config;
    FILE* file;
    int64_t total_samples_written;
    bool is_finalized;
    void* private_data;
};

AudioOutput* audio_output_create_wav(const AudioConfig* config);
AudioOutput* audio_output_create_raw(const AudioConfig* config);

ResultCode audio_output_write_samples(AudioOutput* output, const int16_t* samples, int sample_count);
ResultCode audio_output_finalize(AudioOutput* output);
void audio_output_destroy(AudioOutput* output);

int64_t audio_output_get_samples_written(const AudioOutput* output);
bool audio_output_is_stdout(const AudioOutput* output);

#endif // AUDIO_OUTPUT_H