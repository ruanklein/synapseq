#include "wav_header.h"

#include <string.h>

void wav_header_create(WavHeader* header, const AudioConfig* config, uint32_t data_size) {
    if (!header || !config) return;
    
    memset(header, 0, sizeof(WavHeader));
    
    memcpy(header->riff_id, "RIFF", 4);
    header->file_size = data_size + 36;
    memcpy(header->wave_id, "WAVE", 4);
    
    memcpy(header->fmt_id, "fmt ", 4);
    header->fmt_size = 16;
    header->audio_format = 1;
    header->channels = config->channels;
    header->sample_rate = config->sample_rate;
    header->bits_per_sample = config->bits_per_sample;
    header->byte_rate = config->sample_rate * config->channels * (config->bits_per_sample / 8);
    header->block_align = config->channels * (config->bits_per_sample / 8);
    
    memcpy(header->data_id, "data", 4);
    header->data_size = data_size;
}

void wav_header_write(FILE* file, const WavHeader* header) {
    if (!file || !header) return;
    fwrite(header, sizeof(WavHeader), 1, file);
}

void wav_header_update_size(FILE* file, uint32_t data_size) {
    if (!file) return;
    
    long current_pos = ftell(file);
    
    fseek(file, 4, SEEK_SET);
    uint32_t file_size = data_size + 36;
    fwrite(&file_size, sizeof(uint32_t), 1, file);
    
    fseek(file, 40, SEEK_SET);
    fwrite(&data_size, sizeof(uint32_t), 1, file);
    
    fseek(file, current_pos, SEEK_SET);
}