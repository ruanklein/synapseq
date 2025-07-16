#include "audio_output.h"
#include "wav_header.h"
#include "../utils/memory_utils.h"

#include <string.h>
#include <unistd.h>

typedef struct {
    int64_t data_size;
    bool header_written;
} WavPrivateData;

static ResultCode wav_write_samples(AudioOutput* output, const int16_t* samples, int sample_count) {
    if (!output || !samples || sample_count <= 0) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    WavPrivateData* wav_data = (WavPrivateData*)output->private_data;
    
    if (!wav_data->header_written) {
        WavHeader header;
        wav_header_create(&header, &output->config, 0);
        wav_header_write(output->file, &header);
        wav_data->header_written = true;
    }
    
    size_t bytes_written = fwrite(samples, sizeof(int16_t), (size_t)sample_count, output->file);
    if (bytes_written != (size_t)sample_count) {
        return RESULT_ERROR_AUDIO_DEVICE;
    }
    
    output->total_samples_written += sample_count;
    wav_data->data_size += sample_count * sizeof(int16_t);
    
    return RESULT_SUCCESS;
}

static ResultCode wav_finalize(AudioOutput* output) {
    if (!output || output->is_finalized) {
        return RESULT_SUCCESS;
    }
    
    WavPrivateData* wav_data = (WavPrivateData*)output->private_data;
    
    if (output->file && !audio_output_is_stdout(output)) {
        wav_header_update_size(output->file, (uint32_t)wav_data->data_size);
    }
    
    output->is_finalized = true;
    return RESULT_SUCCESS;
}

static void wav_destroy(AudioOutput* output) {
    if (!output) return;
    
    if (output->file && output->file != stdout) {
        fclose(output->file);
    }
    
    memory_free(output->private_data);
    memory_free(output);
}

static ResultCode raw_write_samples(AudioOutput* output, const int16_t* samples, int sample_count) {
    if (!output || !samples || sample_count <= 0) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    size_t bytes_written = fwrite(samples, sizeof(int16_t), (size_t)sample_count, output->file);
    if (bytes_written != (size_t)sample_count) {
        return RESULT_ERROR_AUDIO_DEVICE;
    }
    
    output->total_samples_written += sample_count;
    return RESULT_SUCCESS;
}

static ResultCode raw_finalize(AudioOutput* output) {
    if (!output || output->is_finalized) {
        return RESULT_SUCCESS;
    }
    
    output->is_finalized = true;
    return RESULT_SUCCESS;
}

static void raw_destroy(AudioOutput* output) {
    if (!output) return;
    
    if (output->file && output->file != stdout) {
        fclose(output->file);
    }
    
    memory_free(output);
}

AudioOutput* audio_output_create_wav(const AudioConfig* config) {
    if (!config) return NULL;
    
    AudioOutput* output = memory_alloc(sizeof(AudioOutput));
    WavPrivateData* wav_data = memory_alloc(sizeof(WavPrivateData));
    
    output->interface.write_samples = wav_write_samples;
    output->interface.finalize = wav_finalize;
    output->interface.destroy = wav_destroy;
    
    output->config = *config;
    output->total_samples_written = 0;
    output->is_finalized = false;
    output->private_data = wav_data;
    
    wav_data->data_size = 0;
    wav_data->header_written = false;
    
    if (strcmp(config->output_filename, "-") == 0) {
        output->file = stdout;
    } else {
        output->file = fopen(config->output_filename, "wb");
        if (!output->file) {
            memory_free(wav_data);
            memory_free(output);
            return NULL;
        }
    }
    
    return output;
}

AudioOutput* audio_output_create_raw(const AudioConfig* config) {
    if (!config) return NULL;
    
    AudioOutput* output = memory_alloc(sizeof(AudioOutput));
    
    output->interface.write_samples = raw_write_samples;
    output->interface.finalize = raw_finalize;
    output->interface.destroy = raw_destroy;
    
    output->config = *config;
    output->total_samples_written = 0;
    output->is_finalized = false;
    output->private_data = NULL;
    
    if (strcmp(config->output_filename, "-") == 0) {
        output->file = stdout;
    } else {
        output->file = fopen(config->output_filename, "wb");
        if (!output->file) {
            memory_free(output);
            return NULL;
        }
    }
    
    return output;
}

ResultCode audio_output_write_samples(AudioOutput* output, const int16_t* samples, int sample_count) {
    if (!output) return RESULT_ERROR_INVALID_ARGUMENT;
    return output->interface.write_samples(output, samples, sample_count);
}

ResultCode audio_output_finalize(AudioOutput* output) {
    if (!output) return RESULT_ERROR_INVALID_ARGUMENT;
    return output->interface.finalize(output);
}

void audio_output_destroy(AudioOutput* output) {
    if (!output) return;
    output->interface.destroy(output);
}

int64_t audio_output_get_samples_written(const AudioOutput* output) {
    return output ? output->total_samples_written : 0;
}

bool audio_output_is_stdout(const AudioOutput* output) {
    return output && output->file == stdout;
}