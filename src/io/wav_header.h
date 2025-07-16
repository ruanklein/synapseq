#ifndef WAV_HEADER_H
#define WAV_HEADER_H

#include "../types/audio_types.h"

#include <stdio.h>

typedef struct {
    char riff_id[4];        // "RIFF"
    uint32_t file_size;     // file size - 8
    char wave_id[4];        // "WAVE"
    char fmt_id[4];         // "fmt "
    uint32_t fmt_size;      // fmt section size (16)
    uint16_t audio_format;  // 1 = PCM
    uint16_t channels;      // number of channels
    uint32_t sample_rate;   // sample rate
    uint32_t byte_rate;     // byte rate
    uint16_t block_align;   // block align
    uint16_t bits_per_sample; // bits per sample
    char data_id[4];        // "data"
    uint32_t data_size;     // data size
} __attribute__((packed)) WavHeader;

void wav_header_create(WavHeader* header, const AudioConfig* config, uint32_t data_size);
void wav_header_write(FILE* file, const WavHeader* header);
void wav_header_update_size(FILE* file, uint32_t data_size);

#endif // WAV_HEADER_H