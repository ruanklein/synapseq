// core/globals.c
#include "globals.h"

// Options
int opt_Q = 0; // Quiet mode
int opt_D = 1; // Test mode
char *opt_o = NULL; // Output file
char *opt_m = NULL; // Background file
int opt_W = 1; // WAV output
int opt_V = 100; // Global volume level

// Audio output
int out_rate = 44100; // Sample rate
int out_rate_def = 1; // Sample rate is default value, not set by user
int out_mode = 1; // Output mode: 0 unsigned char[2], 1 short[2], 2 swapped short[2]
int out_prate = 10; // Rate of parameter change (for file and pipe output only)
int out_fd = 0; // Output file descriptor
int out_bsiz = 0; // Output buffer size (bytes)
int out_blen = 0; // Output buffer length (samples) (1.0* or 0.5* out_bsiz)
int out_bps = 0; // Output bytes per sample (2 or 4)
int out_buf_ms = 0; // Time to output a buffer-ful in ms
int out_buf_lo = 0; // Time to output a buffer-ful, fine-tuning in ms/0x10000
int *tmp_buf = NULL; // Temporary buffer for 20-bit mix values
short *out_buf = NULL; // Output buffer

// Fade interval
int fade_int = 60000; // Fade interval (ms)

// Parser
int in_lin = 0; // Current input line
char buf[4096] = {0}; // Buffer for current line
char buf_copy[4096] = {0}; // Used to keep unmodified copy of line
char *lin = NULL; // Input line (uses buf[])
char *lin_copy = NULL; // Copy of input line
char saved_buf[4096] = {0}; // Buffer for saved line from readNameDef
char saved_buf_copy[4096] = {0}; // Copy of saved line for error messages
int saved_lin_num = -1; // Line number of saved line (-1 = no saved line)

// Input file
FILE *in = NULL; // Input sequence file

// Spin configuration
double spin_carr_max = 0.0; // Maximum 'carrier' value for spin (really max width in us)

// Noise table
int ns_tbl[1 << NS_BIT] = {0};
int ns_off = 0;

// Time sequence configuration
int fast_tim0 = -1; // First time mentioned in the sequence file
int fast_tim1 = -1; // Last time mentioned in the sequence file
int fast_mult = 1; // 0 to sync to clock (adjusting as necessary), or else sync to output rate, with the multiplier indicated
S64 byte_count = -1; // Number of bytes left to output, or -1 if unlimited
int tty_erase = 0; // Chars to erase from current line (for ESC[K emulation)

// Background configuration
FILE *mix_in = NULL; // Input stream for mix sound data, or 0
int mix_cnt = -1; // Version number from mix filename (#<digits>), or -1
int bigendian = 0; // Is this platform Big-endian?
int mix_flag = 0; // Has 'mix/*' been used in the sequence?
double *mix_amp = NULL; // Amplitude of mix sound data to use with effect. Default is 100%

// Miscellaneous
char *pdir = NULL; // Program directory (used as second place to look for background files)
double opt_bg_reduction_db = 12.0; // Default background reduction in 12dB
double bg_gain_factor = 0.25; // Background gain factor (0.0-1.0)
int wav_bits_per_sample = 16; // Default 16-bit
int wav_channels = 2; // Default 2 channels

Channel chan[N_CH] = {0}; // Current channel states
int now = 0; // Current time in ms
Period *per = 0; // Current period
NameDef *nlist = NULL; // Full list of name definitions

// Input data buffering
int *inbuf = NULL; // Buffer for input data (as 20-bit samples)
int ib_len = 0; // Length of input buffer (in ints)
volatile int ib_rd = 0; // Read-offset in inbuf
volatile int ib_wr = 0; // Write-offset in inbuf
volatile int ib_eof = 0; // End of file flag
int ib_cycle = 100; // Time in ms for a complete loop through the buffer
int (*ib_read)(int *, int) = NULL; // Routine to refill buffer

int *sin_tables[4] = {NULL, NULL, NULL, NULL}; // Sine tables for each waveform

// To be used for messages
char *waveform_name[] = {
  "sine",
  "square",
  "triangle",
  "sawtooth"
};