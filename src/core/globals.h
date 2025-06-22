// core/globals.h
#ifndef SYNAPSEQ_GLOBALS_H
#define SYNAPSEQ_GLOBALS_H

#include <stdio.h>

#include "types.h"
#include "constants.h"

// Options
extern int opt_Q; // Quiet mode
extern int opt_D; // Test mode
extern char *opt_o; // Output file
extern char *opt_m; // Background file
extern int opt_W; // WAV output
extern int opt_V; // Global volume level

// Audio output
extern int out_rate; // Sample rate
extern int out_rate_def; // Sample rate is default value, not set by user
extern int out_mode; // Output mode: 0 unsigned char[2], 1 short[2], 2 swapped short[2]
extern int out_prate; // Rate of parameter change (for file and pipe output only)
extern int out_fd; // Output file descriptor
extern int out_bsiz; // Output buffer size (bytes)
extern int out_blen; // Output buffer length (samples) (1.0* or 0.5* out_bsiz)
extern int out_bps; // Output bytes per sample (2 or 4)
extern int out_buf_ms; // Time to output a buffer-ful in ms
extern int out_buf_lo; // Time to output a buffer-ful, fine-tuning in ms/0x10000
extern int *tmp_buf; // Temporary buffer for 20-bit mix values
extern short *out_buf; // Output buffer

// Fade interval
extern int fade_int; // Fade interval (ms)

// Parser
extern int in_lin; // Current input line
extern char buf[4096]; // Buffer for current line
extern char buf_copy[4096]; // Used to keep unmodified copy of line
extern char *lin; // Input line (uses buf[])
extern char *lin_copy; // Copy of input line
extern char saved_buf[4096]; // Buffer for saved line from readNameDef
extern char saved_buf_copy[4096]; // Copy of saved line for error messages
extern int saved_lin_num; // Line number of saved line (-1 = no saved line)

// Input file
extern FILE *in; // Input sequence file

// Spin configuration
extern double spin_carr_max; // Maximum 'carrier' value for spin (really max width in us)

// Noise table
extern int ns_tbl[1 << NS_BIT];
extern int ns_off;

// Time sequence configuration
extern int fast_tim0; // First time mentioned in the sequence file (for -q and -S option)
extern int fast_tim1; // Last time mentioned in the sequence file (for -E option)
extern int fast_mult; // 0 to sync to clock (adjusting as necessary), or else sync to output rate, with the multiplier indicated
extern S64 byte_count; // Number of bytes left to output, or -1 if unlimited
extern int tty_erase; // Chars to erase from current line (for ESC[K emulation)

// Background configuration
extern FILE *mix_in; // Input stream for mix sound data, or 0
extern int mix_cnt; // Version number from mix filename (#<digits>), or -1
extern int bigendian; // Is this platform Big-endian?
extern int mix_flag; // Has 'mix/*' been used in the sequence?
extern double *mix_amp; // Amplitude of mix sound data to use with effect. Default is 100%

// Miscellaneous
extern char *pdir; // Program directory (used as second place to look for background files)
extern double opt_bg_reduction_db; // Default background reduction in 12dB
extern double bg_gain_factor; // Background gain factor (0.0-1.0)
extern int wav_bits_per_sample; // Default 16-bit
extern int wav_channels; // Default 2 channels

extern Channel chan[N_CH]; // Current channel states
extern int now; // Current time in ms
extern Period *per; // Current period
extern NameDef *nlist; // Full list of name definitions

// Input data buffering
extern int *inbuf; // Buffer for input data (as 20-bit samples)
extern int ib_len; // Length of input buffer (in ints)
extern volatile int ib_rd; // Read-offset in inbuf
extern volatile int ib_wr; // Write-offset in inbuf
extern volatile int ib_eof; // End of file flag
extern int ib_cycle; // Time in ms for a complete loop through the buffer
extern int (*ib_read)(int *, int); // Routine to refill buffer

extern int *sin_tables[4]; // Sine tables for each waveform
extern char *waveform_name[]; // To be used for messages

#endif