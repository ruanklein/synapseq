#include "audio_processor.h"
#include "../utils/memory_utils.h"

#include <string.h>
#include <math.h>

AudioProcessor* audio_processor_create(int sample_rate) {
    AudioProcessor* processor = memory_alloc(sizeof(AudioProcessor));
    processor->wave_generator = NULL;
    processor->sample_rate = sample_rate;
    processor->background_buffer = NULL;
    processor->background_buffer_size = 0;
    processor->background_gain = 1.0;
    processor->initialized = false;
    
    // Initialize all voices as inactive
    for (int i = 0; i < MAX_CHANNELS; i++) {
        audio_processor_reset_voice(processor, i);
    }
    
    return processor;
}

void audio_processor_destroy(AudioProcessor* processor) {
    if (!processor) return;
    
    // Clean up voice resources
    for (int i = 0; i < MAX_CHANNELS; i++) {
        audio_processor_reset_voice(processor, i);
    }
    
    if (processor->wave_generator) {
        wave_generator_destroy(processor->wave_generator);
    }
    
    memory_free(processor);
}

ResultCode audio_processor_initialize(AudioProcessor* processor) {
    if (!processor) return RESULT_ERROR_INVALID_ARGUMENT;
    
    processor->wave_generator = wave_generator_create(processor->sample_rate);
    if (!processor->wave_generator) {
        return RESULT_ERROR_MEMORY;
    }
    
    ResultCode result = wave_generator_initialize(processor->wave_generator);
    if (result != RESULT_SUCCESS) {
        wave_generator_destroy(processor->wave_generator);
        processor->wave_generator = NULL;
        return result;
    }
    
    processor->initialized = true;
    return RESULT_SUCCESS;
}

void audio_processor_reset_voice(AudioProcessor* processor, int channel) {
    if (!processor || channel < 0 || channel >= MAX_CHANNELS) return;
    
    VoiceProcessor* voice = &processor->voices[channel];
    
    // Clean up existing resources
    for (int i = 0; i < 2; i++) {
        if (voice->channels[i]) {
            wave_channel_destroy(voice->channels[i]);
            voice->channels[i] = NULL;
        }
    }
    
    if (voice->noise_generator) {
        math_noise_generator_destroy(voice->noise_generator);
        voice->noise_generator = NULL;
    }
    
    if (voice->noise_spin) {
        math_noise_spin_destroy(voice->noise_spin);
        voice->noise_spin = NULL;
    }
    
    // Reset state
    voice->type = VOICE_TYPE_OFF;
    voice->carrier_frequency = 0.0;
    voice->modulation_frequency = 0.0;
    voice->amplitude = 0.0;
    voice->waveform = WAVE_TYPE_SINE;
    voice->phase_accumulator[0] = 0;
    voice->phase_accumulator[1] = 0;
    voice->phase_increment[0] = 0;
    voice->phase_increment[1] = 0;
    voice->active = false;
}

ResultCode audio_processor_setup_voice(AudioProcessor* processor, int channel, const Voice* voice) {
    if (!processor || !voice || channel < 0 || channel >= MAX_CHANNELS) {
        return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    // Reset the voice first
    audio_processor_reset_voice(processor, channel);
    
    VoiceProcessor* voice_proc = &processor->voices[channel];
    
    // Copy basic properties
    voice_proc->type = voice->type;
    voice_proc->carrier_frequency = voice->carrier;
    voice_proc->modulation_frequency = voice->resonance;
    voice_proc->amplitude = voice->amplitude * AMPLITUDE_SCALE_FACTOR;
    voice_proc->waveform = voice->waveform;
    
    // Setup voice-specific resources based on type
    switch (voice->type) {
        case VOICE_TYPE_OFF:
            voice_proc->active = false;
            break;
            
        case VOICE_TYPE_BINAURAL:
            // Create two wave channels for L/R with slightly different frequencies
             voice_proc->channels[0] = wave_channel_create(voice->waveform, 
                voice->carrier - voice->resonance / 2.0, voice_proc->amplitude);
            voice_proc->channels[1] = wave_channel_create(voice->waveform, 
                voice->carrier + voice->resonance / 2.0, voice_proc->amplitude);
                
            wave_channel_set_frequency(voice_proc->channels[0], 
                voice->carrier - voice->resonance / 2.0, processor->sample_rate);
            wave_channel_set_frequency(voice_proc->channels[1], 
                voice->carrier + voice->resonance / 2.0, processor->sample_rate);
                
            voice_proc->active = true;
            break;
            
        case VOICE_TYPE_MONAURAL:
            // Setup phase increments for two frequencies
            voice_proc->phase_increment[0] = (uint32_t)((voice->carrier + voice->resonance / 2.0) 
                / processor->sample_rate * MATH_WAVE_TABLE_SIZE * 65536);
            voice_proc->phase_increment[1] = (uint32_t)((voice->carrier - voice->resonance / 2.0) 
                / processor->sample_rate * MATH_WAVE_TABLE_SIZE * 65536);
            voice_proc->active = true;
            break;
            
        case VOICE_TYPE_ISOCHRONIC:
            // Setup carrier and modulator frequencies
            voice_proc->phase_increment[0] = (uint32_t)(voice->carrier 
                / processor->sample_rate * MATH_WAVE_TABLE_SIZE * 65536);
            voice_proc->phase_increment[1] = (uint32_t)(voice->resonance 
                / processor->sample_rate * MATH_WAVE_TABLE_SIZE * 65536);
            voice_proc->active = true;
            break;
            
        case VOICE_TYPE_WHITE_NOISE:
        case VOICE_TYPE_PINK_NOISE:
        case VOICE_TYPE_BROWN_NOISE:
            voice_proc->noise_generator = math_noise_generator_create();
            voice_proc->active = true;
            break;
            
        case VOICE_TYPE_SPIN_WHITE:
        case VOICE_TYPE_SPIN_PINK:
        case VOICE_TYPE_SPIN_BROWN:
            voice_proc->noise_generator = math_noise_generator_create();
            voice_proc->noise_spin = math_noise_spin_create(voice_proc->noise_generator, voice->type);
            voice_proc->phase_increment[0] = (uint32_t)(voice->resonance 
                / processor->sample_rate * MATH_WAVE_TABLE_SIZE * 65536);
            voice_proc->active = true;
            break;
            
        case VOICE_TYPE_BACKGROUND:
            voice_proc->active = true;
            break;
            
        default:
            return RESULT_ERROR_INVALID_ARGUMENT;
    }
    
    return RESULT_SUCCESS;
}

void audio_processor_update_voices(AudioProcessor* processor, const Voice voices[MAX_CHANNELS]) {
    if (!processor || !voices) return;
    
    for (int i = 0; i < MAX_CHANNELS; i++) {
        audio_processor_setup_voice(processor, i, &voices[i]);
    }
}

void audio_processor_generate_samples(AudioProcessor* processor, int16_t* left_buffer, int16_t* right_buffer, int sample_count) {
    if (!processor || !left_buffer || !right_buffer || !processor->initialized) return;
    
    for (int sample = 0; sample < sample_count; sample++) {
        int32_t left_total = 0;
        int32_t right_total = 0;
        
        // Process each voice
        for (int voice = 0; voice < MAX_CHANNELS; voice++) {
            if (processor->voices[voice].active) {
                int16_t left_sample = 0;
                int16_t right_sample = 0;
                
                audio_processor_process_voice(processor, voice, &left_sample, &right_sample);
                
                left_total += left_sample;
                right_total += right_sample;
            }
        }
        
        // Clamp to 16-bit range
        if (left_total > 32767) left_total = 32767;
        if (left_total < -32768) left_total = -32768;
        if (right_total > 32767) right_total = 32767;
        if (right_total < -32768) right_total = -32768;
        
        left_buffer[sample] = (int16_t)left_total;
        right_buffer[sample] = (int16_t)right_total;
    }
}

void audio_processor_process_voice(AudioProcessor* processor, int channel, int16_t* left_sample, int16_t* right_sample) {
    if (!processor || !left_sample || !right_sample || channel < 0 || channel >= MAX_CHANNELS) {
        *left_sample = 0;
        *right_sample = 0;
        return;
    }
    
    VoiceProcessor* voice = &processor->voices[channel];
    
    switch (voice->type) {
        case VOICE_TYPE_BINAURAL:
            audio_processor_process_binaural(voice, processor->wave_generator, left_sample, right_sample);
            break;
        case VOICE_TYPE_MONAURAL:
            audio_processor_process_monaural(voice, processor->wave_generator, left_sample, right_sample);
            break;
        case VOICE_TYPE_ISOCHRONIC:
            audio_processor_process_isochronic(voice, processor->wave_generator, left_sample, right_sample);
            break;
        case VOICE_TYPE_WHITE_NOISE:
        case VOICE_TYPE_PINK_NOISE:
        case VOICE_TYPE_BROWN_NOISE:
            audio_processor_process_noise(voice, left_sample, right_sample);
            break;
        case VOICE_TYPE_SPIN_WHITE:
        case VOICE_TYPE_SPIN_PINK:
        case VOICE_TYPE_SPIN_BROWN:
            audio_processor_process_spin_noise(voice, left_sample, right_sample);
            break;
        case VOICE_TYPE_BACKGROUND:
            audio_processor_process_background(voice, processor, left_sample, right_sample);
            break;
        default:
            *left_sample = 0;
            *right_sample = 0;
            break;
    }
}

void audio_processor_process_binaural(VoiceProcessor* voice, WaveGenerator* generator, int16_t* left_sample, int16_t* right_sample) {
    if (!voice || !generator || !left_sample || !right_sample || !voice->active) {
        *left_sample = 0;
        *right_sample = 0;
        return;
    }
    
    // Binaural uses two wave channels with slightly different frequencies
    // Left channel: carrier - resonance/2 (lower frequency)
    // Right channel: carrier + resonance/2 (higher frequency)
    // The brain perceives the difference as a "beat" frequency
    
    if (voice->channels[0] && voice->channels[1]) {
        // Get samples from wave tables (in scale 0x7FFFF)
        int32_t left_wave = wave_generator_generate_sample_32bit(generator, voice->channels[0]);
        int32_t right_wave = wave_generator_generate_sample_32bit(generator, voice->channels[1]);
        
        // Apply amplitude (as in original: amp * wave_table_value)
        int32_t left_total = (left_wave * (int32_t)voice->amplitude) >> 16;  // Shift 16 bits as original
        int32_t right_total = (right_wave * (int32_t)voice->amplitude) >> 16;
        
        // Clamp to 16-bit
        *left_sample = (left_total > 32767) ? 32767 : ((left_total < -32768) ? -32768 : (int16_t)left_total);
        *right_sample = (right_total > 32767) ? 32767 : ((right_total < -32768) ? -32768 : (int16_t)right_total);
    } else {
        *left_sample = 0;
        *right_sample = 0;
    }
}

void audio_processor_process_monaural(VoiceProcessor* voice, WaveGenerator* generator, int16_t* left_sample, int16_t* right_sample) {
    (void)voice;
    (void)generator;
    
    // TODO: Implement monaural beats
    *left_sample = 0;
    *right_sample = 0;
}

void audio_processor_process_isochronic(VoiceProcessor* voice, WaveGenerator* generator, int16_t* left_sample, int16_t* right_sample) {
    (void)voice;
    (void)generator;
    
    // TODO: Implement isochronic tones
    *left_sample = 0;
    *right_sample = 0;
}

void audio_processor_process_noise(VoiceProcessor* voice, int16_t* left_sample, int16_t* right_sample) {
    (void)voice;
    
    // TODO: Implement noise generation
    *left_sample = 0;
    *right_sample = 0;
}

void audio_processor_process_spin_noise(VoiceProcessor* voice, int16_t* left_sample, int16_t* right_sample) {
    (void)voice;
    
    // TODO: Implement spinning noise
    *left_sample = 0;
    *right_sample = 0;
}

void audio_processor_process_background(VoiceProcessor* voice, AudioProcessor* processor, int16_t* left_sample, int16_t* right_sample) {
    (void)voice;
    (void)processor;
    
    // TODO: Implement background mixing
    *left_sample = 0;
    *right_sample = 0;
}

void audio_processor_set_background_buffer(AudioProcessor* processor, int16_t* buffer, int buffer_size, double gain) {
    if (!processor) return;
    
    processor->background_buffer = buffer;
    processor->background_buffer_size = buffer_size;
    processor->background_gain = gain;
}

void audio_processor_normalize_output(int16_t* left_sample, int16_t* right_sample) {
    if (!left_sample || !right_sample) return;
    
    // Basic clipping protection
    if (*left_sample > 32767) *left_sample = 32767;
    if (*left_sample < -32768) *left_sample = -32768;
    if (*right_sample > 32767) *right_sample = 32767;
    if (*right_sample < -32768) *right_sample = -32768;
}