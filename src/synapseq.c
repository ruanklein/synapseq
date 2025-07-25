//
//  SynapSeq - Synapse-Sequenced Brainwave Generator
//
//  (c) 2025 Ruan <https://ruan.sh/>
//
//  Originally based on SBaGen (Sequenced Binaural Beat Generator).
//  SBaGen homepage: http://uazu.net/sbagen/
//  License: GNU General Public License v2 (GPLv2)
//
//  This program is free software; you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, version 2. 
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU General Public License for more details.
//
//  See the file COPYING.txt for the full license text.

#define VERSION "2.1.1"

//  T_WIN32	To build for Windows using MinGW and Win32 calls
//  T_POSIX	To build for any Posix-compliant OS
//
// Ogg and MP3 support is handled separately from the T_* macros.

// Define UNIX_TIME to use UNIX calls for getting time
// Define WIN_TIME to use Win32 calls for getting time
// Define UNIX_MISC to use UNIX calls for various miscellaneous things
// Define WIN_MISC to use Windows calls for various miscellaneous things

// Define OGG_DECODE to include OGG support code
// Define MP3_DECODE to include MP3 support code

#ifdef T_WIN32
#define WIN_TIME
#define WIN_MISC
#endif

#ifdef T_POSIX
#define UNIX_TIME
#define UNIX_MISC
#endif

// Make sure one of the _TIME macros is set
#ifndef UNIX_TIME
#ifndef WIN_TIME
#error UNIX_TIME or WIN_TIME not defined.  Maybe you did not define one of T_WIN32/T_POSIX ?
#endif
#endif

// Make sure one of the _MISC macros is set
#ifndef UNIX_MISC
#ifndef WIN_MISC
#error UNIX_MISC or WIN_MISC not defined.  Maybe you did not define one of T_WIN32/T_POSIX ?
#endif
#endif

#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <math.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <time.h>

#include <sys/time.h>
#include <unistd.h>
typedef long long S64;

#ifdef T_WIN32
#define vsnprintf _vsnprintf
#include <windows.h>
#endif

#ifdef UNIX_TIME
#include <sys/ioctl.h>
#include <sys/times.h>
#endif
#ifdef UNIX_MISC
#include <pthread.h>
#endif

typedef struct Channel Channel;
typedef struct Voice Voice;
typedef struct Period Period;
typedef struct NameDef NameDef;
typedef unsigned char uchar;

int inbuf_loop(void *vp);
int inbuf_read(int *dst, int dlen);
void inbuf_start(int (*rout)(int *, int), int len);
int t_per24(int t0, int t1);
int t_per0(int t0, int t1);
int t_mid(int t0, int t1);
int main(int argc, char **argv);
void status(char *);
void dispCurrPer(FILE *);
void init_sin_table();
void init_builtin_namedefs();
void debug(char *fmt, ...);
void warn(char *fmt, ...);
void *Alloc(size_t len);
char *StrDup(char *str);
int calcNow();
void loop();
void outChunk();
void corrVal(int);
int readLine();
char *getWord();
void badSeq();
void readSeq(int ac, char **av);
void correctPeriods();
void setup_device(void);
void readNameDef();
void readTimeLine();
int voicesEq(Voice *, Voice *);
void error(char *fmt, ...);
int sprintTime(char *, int);
int sprintVoice(char *, Voice *, Voice *, int);
int readTime(char *, int *);
void writeWAV();
void writeOut(char *, int);
void find_wav_data_start(FILE *in);
int raw_mix_in(int *dst, int dlen);
void scanOptions(int *acp, char ***avp);
void handleOptionInSequence(char *p);
extern int out_rate, out_rate_def;
void normalizeAmplitude(Voice *voices, int numChannels, const char *line, int lineNum);
void checkBackgroundInSequence(NameDef *nd); // Check if background amplitude is specified
void create_noise_spin_effect(int typ, int amp, int spin_position, int *left, int *right); // Create a spin effect
int restart_background_file(); // Restart the background file for loop
void calculate_bg_gain_factor(); // Calculate the background gain factor
int is_stdout_redirected(); // Check if stdout is redirected

#define ALLOC_ARR(cnt, type) ((type *)Alloc((cnt) * sizeof(type)))
#define uint unsigned int

#ifdef OGG_DECODE
#include "oggdec.c"
#endif
#ifdef MP3_DECODE
#include "mp3dec.c"
#endif

#define NL "\n"

void help() {
  printf(
      "SynapSeq - Synapse-Sequenced Brainwave Generator, version " VERSION NL
      "(c) 2025 Ruan, https://ruan.sh/" NL
      "Released under the GNU GPL v2. See file COPYING." NL NL
      "Usage: synapseq [options] <sequence-file> ..." NL NL
      "Options:  --help                  show this help and exit" NL
      "          --quiet                 quiet mode" NL
      "          --output file           write output to file" NL
      "          --raw                   write raw data instead of WAV" NL
      "          --version               show version and exit" NL NL
      "Example: synapseq --output output.wav sequence-file" NL
      "         synapseq --output - sequence-file | play -" NL
      );
  exit(0);
}

void usage() {
  error("SynapSeq - Synapse-Sequenced Brainwave Generator, version " VERSION NL
        "(c) 2025 Ruan, https://ruan.sh/" NL
        "Released under the GNU GPL v2. See file COPYING." NL NL
        "Usage: synapseq [options] <sequence-file> ..." NL
        "Type 'synapseq --help' for full usage help." NL
      );
}

#define N_CH 16            // Number of channels

struct Voice {
  int typ; // Voice type: 0 off, 1 binaural, 2 pink noise, 3 monaural beats, 4 spin pink, 5
           // background, 6 effect spin, 7 effect pulse, 8 isochronic, 9 white noise, 10 brown
           // noise, 11 spin brown, 12 spin white
  double amp;   // Amplitude level (0-4096 for 0-100%)
  double carr;  // Carrier freq (for binaural/monaural/isochronic), width (for spin)
  double res;   // Resonance freq (-ve or +ve) (for binaural/spin/isochronic)
  int waveform; // 0=sine, 1=square, 2=triangle, 3=sawtooth
};

struct Channel {
  Voice v; // Current voice setting (updated from current period)
  int typ; // Voice type: 0 off, 1 binaural, 2 pink noise, 3 monaural beats, 4 spin pink, 5
           // background, 6 effect spin, 7 effect pulse, 8 isochronic, 9 white noise, 10 brown
           // noise, 11 spin brown, 12 spin white
  int amp, amp2;  // Current state, according to current type
  int inc1, off1; //  ::  (for binaural tones, offset + increment into sine
  int inc2, off2; //  ::   table * 65536)
};

struct Period {
  Period *nxt, *prv;        // Next/prev in chain
  int tim;                  // Start time (end time is ->nxt->tim)
  Voice v0[N_CH], v1[N_CH]; // Start and end voices
  int fi, fo;               // Temporary: Fade-in, fade-out modes
};

struct NameDef {
  NameDef *nxt;
  char *name;     // Name of definition
  Voice vv[N_CH]; // Voice-set for it (unless a block definition)
};

#define ST_AMP 0x7FFFF // Amplitude of wave in sine-table
#define NS_ADJ 12 // Noise is generated internally with amplitude ST_AMP<<NS_ADJ
#define NS_DITHER 16 // How many bits right to shift the noise for dithering
#define NS_AMP (ST_AMP << NS_ADJ)
#define ST_SIZ 16384 // Number of elements in sine-table (power of 2)
int *sin_tables[4];  // 0=sine, 1=square, 2=triangle, 3=sawtooth
#define AMP_DA(pc) (40.96 * (pc))   // Display value (%age) to ->amp value
#define AMP_AD(amp) ((amp) / 40.96) // Amplitude value to display %age

Channel chan[N_CH]; // Current channel states
int now;            // Current time (milliseconds from midnight)
Period *per = 0;    // Current period
NameDef *nlist;     // Full list of name definitions

int *tmp_buf;         // Temporary buffer for 20-bit mix values
short *out_buf;       // Output buffer
int out_bsiz;         // Output buffer size (bytes)
int out_blen;         // Output buffer length (samples) (1.0* or 0.5* out_bsiz)
int out_bps;          // Output bytes per sample (2 or 4)
int out_buf_ms;       // Time to output a buffer-ful in ms
int out_buf_lo;       // Time to output a buffer-ful, fine-tuning in ms/0x10000
int out_fd;       // Output file descriptor
int out_rate = 44100; // Sample rate
int out_rate_def = 1; // Sample rate is default value, not set by user
int out_mode = 1; // Output mode: 0 unsigned char[2], 1 short[2], 2 swapped short[2]
int out_prate = 10; // Rate of parameter change (for file and pipe output only)
int fade_int = 60000;      // Fade interval (ms)
FILE *in;                  // Input sequence file
int in_lin;                // Current input line
char buf[4096];            // Buffer for current line
char buf_copy[4096];       // Used to keep unmodified copy of line
char *lin;                 // Input line (uses buf[])
char *lin_copy;            // Copy of input line
char saved_buf[4096];      // Buffer for saved line from readNameDef
char saved_buf_copy[4096]; // Copy of saved line for error messages
int saved_lin_num = -1;    // Line number of saved line (-1 = no saved line)

double
    spin_carr_max; // Maximum 'carrier' value for spin (really max width in us)

#define NS_BIT 10
int ns_tbl[1 << NS_BIT];
int ns_off = 0;

int fast_tim0 =
    -1; // First time mentioned in the sequence file (for -q and -S option)
int fast_tim1 = -1; // Last time mentioned in the sequence file (for -E option)
int fast_mult =
    1; // 0 to sync to clock (adjusting as necessary), or else sync to
       //  output rate, with the multiplier indicated
S64 byte_count = -1; // Number of bytes left to output, or -1 if unlimited
int tty_erase;       // Chars to erase from current line (for ESC[K emulation)

int opt_Q; // Quiet mode
int opt_D = 1; // Test mode
char *opt_o, *opt_m;
int opt_W = 1; // Default to WAV output
int opt_V = 100; // Global volume level (default 100%)
char *waveform_name[] = {"sine", "square", "triangle",
                         "sawtooth"}; // To be used for messages

FILE *mix_in;           // Input stream for mix sound data, or 0
int mix_cnt;            // Version number from mix filename (#<digits>), or -1
int bigendian;          // Is this platform Big-endian?
int mix_flag = 0;       // Has 'mix/*' been used in the sequence?
double *mix_amp = NULL; // Amplitude of mix sound data to use with
                        // effect. Default is 100%

char *pdir; // Program directory (used as second place to look for background files)
double opt_bg_reduction_db = 12.0; // Default background reduction in 12dB
double bg_gain_factor = 0.25; // Background gain factor (0.0-1.0)
static int wav_bits_per_sample = 16;  // Default 16-bit
static int wav_channels = 2;  // Default 2 channels

//
//	Delay for a short period of time (in ms)
//

#ifdef UNIX_MISC
void delay(int ms) {
  struct timespec ts;
  ts.tv_sec = ms / 1000;
  ts.tv_nsec = (ms % 1000) * 1000000;
  nanosleep(&ts, 0);
}
#endif
#ifdef WIN_MISC
void delay(int ms) { Sleep(ms); }
#endif

//
//	WAV/OGG/MP3 input data buffering
//

int *inbuf;                 // Buffer for input data (as 20-bit samples)
int ib_len;                 // Length of input buffer (in ints)
volatile int ib_rd;         // Read-offset in inbuf
volatile int ib_wr;         // Write-offset in inbuf
volatile int ib_eof;        // End of file flag
int ib_cycle = 100;         // Time in ms for a complete loop through the buffer
int (*ib_read)(int *, int); // Routine to refill buffer

int inbuf_loop(void *vp) {
  int now = -1;
  int waited = 0; // Used to bail out if the main thread dies for some reason
  int a;

  while (1) {
    int rv;
    int rd = ib_rd;
    int wr = ib_wr;
    int cnt = (rd - 1 - wr) & (ib_len - 1);
    if (cnt > ib_len - wr)
      cnt = ib_len - wr;
    if (cnt > ib_len / 8)
      cnt = ib_len / 8;

    // Choose to only work in ib_len/8 units, although this is not
    // 100% necessary
    if (cnt < ib_len / 8) {
      // Wait a little while for the buffer to empty (minimum 1ms)
      if (waited > 10000 + ib_cycle)
        error("Mix stream halted for more than 10 seconds; aborting");
      delay(a = 1 + ib_cycle / 4);
      waited += a;
      continue;
    }
    waited = 0;

    rv = ib_read(inbuf + wr, cnt);
    // debug("ib_read %d-%d (%d) -> %d", wr, wr+cnt-1, cnt, rv);
    if (rv != cnt) {
      // Instead of immediately setting EOF, try restarting
      if (restart_background_file()) {
        // Try reading again after restarting
        int remaining = cnt - rv;
        int new_rv = ib_read(inbuf + wr + rv, remaining);
        rv += new_rv;
        
        if (rv != cnt) {
          ib_eof = 1;
          return 0;
        }
      } else {
        ib_eof = 1;
        return 0;
      }
    }

    ib_wr = (wr + rv) & (ib_len - 1);

    // Whenever we roll over, recalculate 'ib_cycle'
    if (ib_wr < wr) {
      int prev = now;
      now = calcNow();
      if (prev >= 0 && now > prev)
        ib_cycle = now - prev;
      // debug("Input buffer cycle duration is now %dms", ib_cycle);
    }
  }
  return 0;
}

//
//	Read a chunk of int data from the input buffer.  This will
//	always return enough data unless we have hit the end of the
//	file, in which case it returns a lower number or 0.  If not
//	enough data has been read by the input thread, then this
//	thread pauses until data is ready -- but this should hopefully
//	never happen.
//

int inbuf_read(int *dst, int dlen) {
  int rv = 0;
  int waited =
      0; // As a precaution, bail out if other thread hangs for some reason
  int a;

  while (dlen > 0) {
    int rd = ib_rd;
    int wr = ib_wr;
    int avail = (wr - rd) & (ib_len - 1);
    int toend = ib_len - rd;
    if (avail > toend)
      avail = toend;
    if (avail > dlen)
      avail = dlen;

    if (avail == 0) {
      if (ib_eof)
        return rv;

      // Necessary to wait for incoming mix data.  This should
      // never happen in normal running, though, unless we are
      // outputting to a file
      if (waited > 10000)
        error("Mix stream problem; waited more than 10 seconds for data; "
              "aborting");
      // debug("Waiting for input thread (%d)", ib_eof);
      delay(a = ib_cycle / 4 > 100 ? 100 : 1 + ib_cycle / 4);
      waited += a;
      continue;
    }
    waited = 0;

    memcpy(dst, inbuf + rd, avail * sizeof(int));
    dst += avail;
    dlen -= avail;
    rv += avail;
    ib_rd = (rd + avail) & (ib_len - 1);
  }
  return rv;
}

//
//	Start off the thread that fills the buffer
//

void inbuf_start(int (*rout)(int *, int), int len) {
  if (0 != (len & (len - 1)))
    error("inbuf_start() called with length not a power of two");

  ib_read = rout;
  ib_len = len;
  inbuf = ALLOC_ARR(ib_len, int);
  ib_rd = 0;
  ib_wr = 0;
  ib_eof = 0;

  // Preload 75% of the buffer -- or at least attempt to do so;
  // errors/eof/etc will be picked up in the inbuf_loop() routine
  ib_wr = ib_read(inbuf, ib_len * 3 / 4);

  // Start the thread off
#ifdef UNIX_MISC
  {
    pthread_t thread;
    if (0 != pthread_create(&thread, NULL, (void *)&inbuf_loop, NULL))
      error("Failed to start input buffering thread");
  }
#endif
#ifdef WIN_MISC
  {
    DWORD tmp;
    if (0 ==
        CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)&inbuf_loop, 0, 0, &tmp))
      error("Failed to start input buffering thread");
  }
#endif
}

//
//	Time-keeping functions
//

#define H24 (86400000) // 24 hours
#define H12 (43200000) // 12 hours

int t_per24(int t0, int t1) { // Length of period starting at t0, ending at t1.
  int td = t1 - t0;           // NB for t0==t1 this gives 24 hours, *NOT 0*
  return td > 0 ? td : td + H24;
}
int t_per0(int t0, int t1) { // Length of period starting at t0, ending at t1.
  int td = t1 - t0;          // NB for t0==t1 this gives 0 hours
  return td >= 0 ? td : td + H24;
}
int t_mid(int t0, int t1) { // Midpoint of period from t0 to t1
  return ((t1 < t0) ? (H24 + t0 + t1) / 2 : (t0 + t1) / 2) % H24;
}

int restart_background_file() {
  if (!mix_in) return 0;

  char *ext = strrchr(opt_m, '.');
  if (!ext) {
    warn("No extension found in background file: %s", opt_m);
    return 0;
  }

#ifdef OGG_DECODE
  if (opt_m && 0 == strcmp(ext, "ogg")) {
    if (ov_raw_seek(&oggfile, 0) != 0) {
      warn("Could not seek to beginning of Ogg file");
      return 0;
    }

    ogg_rd = ogg_end = ogg_buf0;
    ib_eof = 0;
    return 1;
  }
#endif
  // Go back to the beginning of the file
  if (fseek(mix_in, 0, SEEK_SET) != 0) {
    if (!opt_Q)
      warn("Warning: Could not seek to beginning of background file for loop");
    return 0;
  }

  if (opt_m && 0 == strcmp(ext, "wav")) {
    find_wav_data_start(mix_in);
  }
  
  // Reset input buffer flags
  ib_eof = 0;
  return 1;
}

// Calculate the background gain factor based on the reduction in dB
void calculate_bg_gain_factor() {
  bg_gain_factor = pow(10.0, -opt_bg_reduction_db / 20.0);
}

int is_stdout_redirected() {
  return isatty(STDOUT_FILENO) == 0;
}

//
//	M A I N
//

int main(int argc, char **argv) {
  short test = 0x1100;
  char *p;

  pdir = StrDup(argv[0]);
  p = strchr(pdir, 0);
  while (p > pdir && p[-1] != '/' && p[-1] != '\\')
    *--p = 0;

  argc--;
  argv++;
  bigendian = ((char *)&test)[0] != 0;

  // Process all the options
  scanOptions(&argc, &argv);

  calculate_bg_gain_factor();

  if (argc < 1) usage();

  init_builtin_namedefs();
  readSeq(argc, argv);
  init_sin_table();

  mix_in = 0;
  if (opt_m) {
    char *p;
    char tmp[4];
    int raw = 1;
    if (opt_m) {
      // Pick up #<digits> on end of filename
      p = strchr(opt_m, 0);
      mix_cnt = -1;
      if (p > opt_m && isdigit(p[-1])) {
        mix_cnt = 0;
        while (p > opt_m && isdigit(p[-1]))
          mix_cnt = mix_cnt * 10 + *--p - '0';
        if (p > opt_m && p[-1] == '#')
          *--p = 0;
        else {
          p = strchr(opt_m, 0);
          mix_cnt = -1;
        }
      }
      // p points to end of filename (NUL)

      // Open file
      mix_in = fopen(opt_m, "rb");
      if (!mix_in && opt_m[0] != '/') {
        int len = strlen(opt_m) + strlen(pdir) + 1;
        char *tmp = ALLOC_ARR(len, char);
        strcpy(tmp, pdir);
        strcat(tmp, opt_m);
        mix_in = fopen(tmp, "rb");
        free(tmp);
      }
      if (!mix_in)
        error("Can't open background input file: %s", opt_m);

      // Pick up extension
      if (p - opt_m >= 4 && p[-4] == '.') {
        tmp[0] = tolower(p[-3]);
        tmp[1] = tolower(p[-2]);
        tmp[2] = tolower(p[-1]);
        tmp[3] = 0;
      }
    }
    if (0 == strcmp(tmp, "wav")) // Skip header on WAV files
      find_wav_data_start(mix_in);
    if (0 == strcmp(tmp, "ogg")) {
#ifdef OGG_DECODE
      ogg_init();
      raw = 0;
#else
      error("Sorry: Ogg support wasn't compiled into this executable");
#endif
    }
    if (0 == strcmp(tmp, "mp3")) {
#ifdef MP3_DECODE
      mp3_init();
      raw = 0;
#else
      error("Sorry: MP3 support wasn't compiled into this executable");
#endif
    }
    // If this is a raw/wav data stream, setup a 256*1024-int
    // buffer (3s@44.1kHz)
    if (raw)
      inbuf_start(raw_mix_in, 256 * 1024);
  }

  loop();
  return 0;
}

//
//	Scan options.  Returns a flag indicating what is expected to
//	interpret the rest of the arguments: 0 normal, 'i' immediate
//	(-i option), 'p' -p option.
//

void scanOptions(int *acp, char ***avp) {
  int argc = *acp;
  char **argv = *avp;

  // Scan options
  while (argc > 0 && argv[0][0] == '-' && argv[0][1] == '-') {
    char *current_opt = argv[0] + 2;
    argc--;
    argv++;
    
    if (strcmp(current_opt, "help") == 0) {
      help();
    }
    else if (strcmp(current_opt, "quiet") == 0) {
      opt_Q = 1;
    }
    else if (strcmp(current_opt, "output") == 0) {
      if (argc < 1)
        error("--output expects output file name");

      opt_o = argv[0];
      opt_D = 0;

      if (opt_o[0] == '-' && opt_o[1] == 0) {
        out_fd = 1; // stdout
        opt_o = NULL;
      }

      argc--;
      argv++;
    }
    else if (strcmp(current_opt, "raw") == 0) {
      opt_W = 0;
    }
    else if (strcmp(current_opt, "version") == 0) {
      printf("SynapSeq version %s\n", VERSION);
      printf("\nAudio format support for the @background option:\n");
#ifdef OGG_DECODE
      printf("  ✓ OGG support: enabled\n");
#else 
      printf("  ✗ OGG support: disabled\n");
#endif
#ifdef MP3_DECODE
      printf("  ✓ MP3 support: enabled\n");
#else
      printf("  ✗ MP3 support: disabled\n");
#endif
      printf("  ✓ WAV support: enabled (built-in)\n");
      exit(0);
    }
    else {
      error("Invalid option: --%s. Type 'synapseq --help' for help.", current_opt);
    }
  }

  *acp = argc;
  *avp = argv;
}

//
// Handle options in the sequence file
//

void handleOptionInSequence(char *p) {
  char *option = getWord();

  if (strcmp(option, "@background") == 0) {
    if (opt_m) {
      error("Background file already set at line %d: %s", in_lin, lin_copy);
    }
    char *file_name = getWord();
    if (file_name) {
      opt_m = StrDup(file_name);
    } else {
      error("File name expected at line %d: %s", in_lin, lin_copy);
    }
  } else if (strcmp(option, "@gainlevel") == 0) {
    char *gain_level = getWord();
    if (!gain_level) {
      error("Gain level expected at line %d: %s", in_lin, lin_copy);
    }
    if (strcmp(gain_level, "verylow") == 0) {
      opt_bg_reduction_db = 20.0;
    } else if (strcmp(gain_level, "low") == 0) {
      opt_bg_reduction_db = 16.0;
    } else if (strcmp(gain_level, "medium") == 0) {
      opt_bg_reduction_db = 12.0;
    } else if (strcmp(gain_level, "high") == 0) {
      opt_bg_reduction_db = 6.0;
    } else if (strcmp(gain_level, "veryhigh") == 0) {
      opt_bg_reduction_db = 0.0;
    } else {
      error("Invalid gain level at line %d: %s", in_lin, lin_copy);
    }
    calculate_bg_gain_factor();
  } 
  else if (strcmp(option, "@volume") == 0) {
    char *volume_str = getWord();
    if (1 != sscanf(volume_str, "%d", &opt_V)) {
      error("Invalid volume value at line %d: %s", in_lin, lin_copy);
    }
    if (opt_V < 0 || opt_V > 100) {
      error("Volume value must be between 0 and 100 at line %d: %s", in_lin,
            lin_copy);
    }
  } else if (strcmp(option, "@samplerate") == 0) {
    char *samplerate_str = getWord();
    if (1 != sscanf(samplerate_str, "%d", &out_rate)) {
      error("Invalid samplerate value at line %d: %s", in_lin, lin_copy);
    }
    out_rate_def = 0;
  }
  else {
    error("Invalid option at line %d: %s", in_lin, lin_copy);
  }
}

//
//	If this is a WAV file we've been given, skip forward to the
//	'data' section.  Don't bother checking any of the 'fmt '
//	stuff.  If they didn't give us a valid 16-bit stereo file at
//	the right rate, then tough!
//

void find_wav_data_start(FILE *in) {
  unsigned char buf[32];
  
  if (1 != fread(buf, 12, 1, in))
    goto bad;
  if (0 != memcmp(buf, "RIFF", 4))
    goto bad;
  if (0 != memcmp(buf + 8, "WAVE", 4))
    goto bad;

  while (1) {
    int len;
    if (1 != fread(buf, 8, 1, in))
      goto bad;
    if (0 == memcmp(buf, "data", 4)) {
      return;
    }
    
    len = buf[4] + (buf[5] << 8) + (buf[6] << 16) + (buf[7] << 24);
    if (len & 1)
      len++;
      
    if (out_rate_def && 0 == memcmp(buf, "fmt ", 4)) {
      if (1 != fread(buf, len > 24 ? 24 : len, 1, in))
        goto bad;
      wav_channels = buf[2] + (buf[3] << 8);
      out_rate = buf[4] + (buf[5] << 8) + (buf[6] << 16) + (buf[7] << 24);
      wav_bits_per_sample = buf[14] + (buf[15] << 8);
      out_rate_def = 0;
      
      if (len > 24 && 0 != fseek(in, len - 24, SEEK_CUR))
        goto bad;
    } else {
      if (0 != fseek(in, len, SEEK_CUR))
        goto bad;
    }
  }

bad:
  warn("WARNING: Not a valid WAV file, treating as RAW");
  rewind(in);
  wav_bits_per_sample = 16;
  wav_channels = 2;
}

//
//	Input raw audio data from the 'mix_in' stream, and convert to
//	32-bit values (max 'dlen')
//

int raw_mix_in(int *dst, int dlen) {
  int a, rv;

   if (wav_bits_per_sample == 16) {
    short *tmp = (short *)(dst + dlen / 2);
    rv = fread(tmp, 2, dlen, mix_in);
    if (rv == 0) {
      if (feof(mix_in)) return 0;
      error("Read error on mix input:\n  %s", strerror(errno));
    }
    
    for (a = 0; a < rv; a++) {
      int sample = *tmp++;
      *dst++ = sample << 4;  // 16-bit → 20-bit
    }
    
  } else if (wav_bits_per_sample == 24) {
    unsigned char *tmp = (unsigned char *)Alloc(dlen * 3);
    rv = fread(tmp, 3, dlen, mix_in);
    if (rv == 0) {
      if (feof(mix_in)) {
        free(tmp);
        return 0;
      }
      error("Read error on mix input:\n  %s", strerror(errno));
    }
    
    unsigned char *ptr = tmp;
    for (a = 0; a < rv; a++) {
      int sample = ptr[0] + (ptr[1] << 8) + (ptr[2] << 16);
      if (sample & 0x800000) sample |= 0xFF000000;
      *dst++ = sample >> 4;
      ptr += 3;
    }
    
    free(tmp);
    
  } else {
    error("Unsupported WAV format: %d bits per sample", wav_bits_per_sample);
  }

  return rv;
}

//
//	Update a status line
//

void status(char *err) {
  int a;
  int nch = N_CH;
  char *p = buf, *p0, *p1;

  if (opt_Q)
    return;

  p0 = p; // Start of line
  *p++ = ' ';
  *p++ = ' ';
  p += sprintTime(p, now);
  while (nch > 1 && chan[nch - 1].v.typ == 0)
    nch--;
  for (a = 0; a < nch; a++)
    p += sprintVoice(p, &chan[a].v, 0, 0);
  if (err)
    p += sprintf(p, " %s", err);
  p1 = p; // End of line

  while (tty_erase > p - p0)
    *p++ = ' ';

  tty_erase = p1 - p0; // Characters that will need erasing

  fprintf(stderr, "%s\r", buf);
  fflush(stderr);
}

void // Display current period details
dispCurrPer(FILE *fp) {
  int a;
  Voice *v0, *v1;
  char *p0, *p1;
  int len0, len1;
  int nch = N_CH;

  if (opt_Q)
    return;

  p0 = buf;
  p1 = buf_copy;

  p0 += sprintf(p0, "- ");
  p0 += sprintTime(p0, per->tim);
  p1 += sprintf(p1, "  ");
  p1 += sprintTime(p1, per->nxt->tim);

  v0 = per->v0;
  v1 = per->v1;
  while (nch > 1 && v0[nch - 1].typ == 0)
    nch--;
  for (a = 0; a < nch; a++, v0++, v1++) {
    p0 += len0 = sprintVoice(p0, v0, 0, 1);
    p1 += len1 = sprintVoice(p1, v1, v0, 1);
    while (len0 < len1) {
      *p0++ = ' ';
      len0++;
    }
    while (len1 < len0) {
      *p1++ = ' ';
      len1++;
    }
  }
  *p0 = 0;
  *p1 = 0;
  fprintf(fp, "%s\n%s\n", buf, buf_copy);
  fflush(fp);
}

int sprintTime(char *p, int tim) {
  return sprintf(p, "%02d:%02d:%02d", tim % 86400000 / 3600000,
                 tim % 3600000 / 60000, tim % 60000 / 1000);
}

int sprintVoice(char *p, Voice *vp, Voice *dup, int multiline) {
  switch (vp->typ) {
  case 0:
    return sprintf(p, " -");
  case 1:
    if (dup && vp->carr == dup->carr && vp->res == dup->res &&
        vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(
          p, "\n\twaveform %s tone %.2f binaural %.2f amplitude %.2f",
          waveform_name[vp->waveform], vp->carr, vp->res, AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (tone:%.2f binaural:%.2f amplitude:%.2f)", vp->carr,
                     vp->res, AMP_AD(vp->amp));
    }
  case 2:
    if (dup && vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(p, "\n\tnoise pink amplitude %.2f", AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (noise:%.2f)", AMP_AD(vp->amp));
    }
  case 3:
    if (dup && vp->carr == dup->carr && vp->res == dup->res &&
        vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(p, "\n\twaveform %s tone %.2f monaural %.2f amplitude %.2f",
                     waveform_name[vp->waveform], vp->carr, vp->res,
                     AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (tone:%.2f monaural:%.2f amplitude:%.2f)", vp->carr,
                     vp->res, AMP_AD(vp->amp));
    }
  case 9:
    if (dup && vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(p, "\n\tnoise white amplitude %.2f", AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (noise:%.2f)", AMP_AD(vp->amp));
    }
  case 10:
    if (dup && vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(p, "\n\tnoise brown amplitude %.2f", AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (noise:%.2f)", AMP_AD(vp->amp));
    }
  case 4:
    if (dup && vp->carr == dup->carr && vp->res == dup->res &&
        vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(
          p,
          "\n\twaveform %s spin pink width %.2f rate %.2f amplitude %.2f",
          waveform_name[vp->waveform], vp->carr, vp->res, AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (width:%.2f rate:%.2f amplitude:%.2f)", vp->carr,
                     vp->res, AMP_AD(vp->amp));
    }
  case 5:
    if (dup && vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(p, "\n\tbackground amplitude %.2f", AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (amplitude:%.2f)", AMP_AD(vp->amp));
    }
  case 8: // Isochronic tones
    if (dup && vp->carr == dup->carr && vp->res == dup->res &&
        vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(
          p, "\n\twaveform %s tone %.2f isochronic %.2f amplitude %.2f",
          waveform_name[vp->waveform], vp->carr, vp->res, AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (tone:%.2f isochronic:%.2f amplitude:%.2f)", vp->carr,
                     vp->res, AMP_AD(vp->amp));
    }
  case 6: // Mixspin - spinning mix stream
    if (dup && vp->carr == dup->carr && vp->res == dup->res &&
        vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(p,
                     "\n\twaveform %s effect spin width %.2f rate %.2f "
                     "intensity %.2f",
                     waveform_name[vp->waveform], vp->carr, vp->res,
                     AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (width:%.2f rate:%.2f intensity:%.2f)", vp->carr,
                     vp->res, AMP_AD(vp->amp));
    }
  case 7: // Mixpulse - mix stream with pulse effect
    if (dup && vp->res == dup->res && vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(p, "\n\twaveform %s effect pulse %.2f intensity %.2f",
                     waveform_name[vp->waveform], vp->res, AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (pulse:%.2f intensity:%.2f)", vp->res,
                     AMP_AD(vp->amp));
    }
  case 11: // Bspin - spinning brown noise
    if (dup && vp->carr == dup->carr && vp->res == dup->res &&
        vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(
          p,
          "\n\twaveform %s spin brown width %.2f rate %.2f amplitude %.2f",
          waveform_name[vp->waveform], vp->carr, vp->res, AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (width:%.2f rate:%.2f amplitude:%.2f)", vp->carr,
                     vp->res, AMP_AD(vp->amp));
    }
  case 12: // Wspin - spinning white noise
    if (dup && vp->amp == dup->amp)
      return sprintf(p, "  --");
    if (multiline) {
      return sprintf(
          p,
          "\n\twaveform %s spin white width %.2f rate %.2f amplitude %.2f",
          waveform_name[vp->waveform], vp->carr, vp->res, AMP_AD(vp->amp));
    } else {
      return sprintf(p, " (width:%.2f rate:%.2f amplitude:%.2f)", vp->carr,
                     vp->res, AMP_AD(vp->amp));
    }
  default:
    return sprintf(p, " ???");
  }
}

void init_sin_table() {
  int a, waveform;

  // Initialize all 4 waveform tables
  for (waveform = 0; waveform < 4; waveform++) {
    int *arr = (int *)Alloc(ST_SIZ * sizeof(int));

    for (a = 0; a < ST_SIZ; a++) {
      double phase = (a * 2.0 * 3.14159265358979323846) / ST_SIZ;
      double val;

      switch (waveform) {
      case 0: // Sine
        val = sin(phase);
        break;
      case 1: // Square
        val = (sin(phase) >= 0) ? 1.0 : -1.0;
        break;
      case 2: // Triangle
        if (phase < 3.14159265358979323846)
          val = (2.0 * phase / 3.14159265358979323846) - 1.0;
        else
          val = 3.0 - (2.0 * phase / 3.14159265358979323846);
        break;
      case 3: // Sawtooth
        val = (2.0 * phase / (2.0 * 3.14159265358979323846)) - 1.0;
        break;
      default:
        val = sin(phase);
      }

      arr[a] = (int)(ST_AMP * val);
    }

    sin_tables[waveform] = arr;
  }
}

//
//  Initialize built-in name definitions (like 'silence')
//
void init_builtin_namedefs() {
  NameDef *nd;

  // Create the pre-defined nameDef "silence"
  nd = (NameDef *)Alloc(sizeof(NameDef));
  nd->name = StrDup("silence");

  // Initialize all voices as "off" (typ=0, amp=0)
  for (int i = 0; i < N_CH; i++) {
    nd->vv[i].typ = 0;
    nd->vv[i].amp = 0;
    nd->vv[i].carr = 0;
    nd->vv[i].res = 0;
    nd->vv[i].waveform = 0;
  }

  // Add to the list
  nd->nxt = nlist;
  nlist = nd;
}

void error(char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  vfprintf(stderr, fmt, ap);
  fprintf(stderr, "\n");

  exit(1);
}

void debug(char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  vfprintf(stderr, fmt, ap);
  fprintf(stderr, "\n");
}

void warn(char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  vfprintf(stderr, fmt, ap);
  fprintf(stderr, "\n");
}

void *Alloc(size_t len) {
  void *p = calloc(1, len);
  if (!p)
    error("Out of memory");
  return p;
}

char *StrDup(char *str) {
  char *rv = strdup(str);
  if (!rv)
    error("Out of memory");
  return rv;
}

#ifdef UNIX_TIME
// Precalculate a reference timestamp to accelerate calcNow().  This
// can be any recent time.  We recalculate it every 10 minutes.  The
// only reason for doing this is to cope with clocks going forwards or
// backwards when entering or leaving summer time so that people wake
// up on time on these two dates; an hour of the sequence will be
// repeated or skipped.  The 'time_ref*' variables will be initialised
// on the first call to calcNow().

static int time_ref_epoch = 0; // Reference time compared to UNIX epoch
static int time_ref_ms; // Reference time in synapseq 24-hour milliseconds

void setupRefTime() {
  struct tm *tt;
  time_t tim = time(0);
  tt = localtime(&tim);
  time_ref_epoch = tim;
  time_ref_ms = 1000 * tt->tm_sec + 60000 * tt->tm_min + 3600000 * tt->tm_hour;
}

inline int calcNow() {
  struct timeval tv;
  if (0 != gettimeofday(&tv, 0))
    error("Can't get current time");
  if (tv.tv_sec - time_ref_epoch > 600)
    setupRefTime();
  return (time_ref_ms + (tv.tv_sec - time_ref_epoch) * 1000 +
          tv.tv_usec / 1000) %
         H24;
}
#endif

#ifdef WIN_TIME
inline int calcNow() {
  SYSTEMTIME st;
  GetLocalTime(&st);
  return st.wMilliseconds + 1000 * st.wSecond + 60000 * st.wMinute +
         3600000 * st.wHour;
}
#endif

//
//	Simple random number generator.  Generates a repeating
//	sequence of 65536 odd numbers in the range -65535->65535.
//
//	Based on ZX Spectrum random number generator:
//	  seed= (seed+1) * 75 % 65537 - 1
//

#define RAND_MULT 75

static int seed = 2;

//
//	Generate next sample for simulated pink noise, with same
//	scaling as the sin_table[].  This version uses an inlined
//	random number generator, and smooths the lower frequency bands
//	as well.
//

#define NS_BANDS 9
typedef struct Noise Noise;
struct Noise {
  int val; // Current output value
  int inc; // Increment
};
Noise ntbl[NS_BANDS];
int nt_off;
int noise_buf[256];
uchar noise_off = 0;

static inline int noise2() {
  int tot;
  int off = nt_off++;
  int cnt = 1;
  Noise *ns = ntbl;
  Noise *ns1 = ntbl + NS_BANDS;

  tot = ((seed = seed * RAND_MULT % 131074) - 65535) *
        (NS_AMP / 65535 / (NS_BANDS + 1));

  while ((cnt & off) && ns < ns1) {
    int val = ((seed = seed * RAND_MULT % 131074) - 65535) *
              (NS_AMP / 65535 / (NS_BANDS + 1));
    tot += ns->val += ns->inc = (val - ns->val) / (cnt += cnt);
    ns++;
  }

  while (ns < ns1) {
    tot += (ns->val += ns->inc);
    ns++;
  }

  return noise_buf[noise_off++] = (tot >> NS_ADJ);
}

//
//	Generate next sample for white noise, with same
//	scaling as the sin_table[].
//
static inline int white_noise() {
  // White noise is simply a random value without filtering
  return ((seed = seed * RAND_MULT % 131074) - 65535) * (ST_AMP / 65535);
}

//
//	Generate next sample for brown noise, with same
//	scaling as the sin_table[].
//	Brown noise has more energy in lower frequencies,
//	implemented as a random walk (integration of white noise)
//
static int brown_last = 0;
static inline int brown_noise() {
  // Generate a random value
  int random = ((seed = seed * RAND_MULT % 131074) - 65535);

  // Integrate the random value with a decay factor to avoid overflow
  brown_last = (brown_last + (random / 16)) * 0.9;

  // Limit the value to avoid overflow
  if (brown_last > 65535)
    brown_last = 65535;
  if (brown_last < -65535)
    brown_last = -65535;

  // Scale to the same level as the sin_table
  return brown_last * (ST_AMP / 65535);
}

// Create a spin effect for noise based on the spin_position

void create_noise_spin_effect(int typ, int amp, int spin_position, int *left,
                              int *right) {
  // Apply intensity factor to rotation value
  int amplified_val = (int)(spin_position * 1.5);

  // Limit value between -128 and 127
  if (amplified_val > 127)
    amplified_val = 127;
  if (amplified_val < -128)
    amplified_val = -128;

  // Use absolute value for calculations
  int pos_val = amplified_val < 0 ? -amplified_val : amplified_val;
  int noise_l, noise_r;

  int base_noise;
  switch (typ) {
  case 11:
    base_noise = brown_noise();
    break;
  case 12:
    base_noise = white_noise();
    break;
  default:
    // Default is pink noise
    base_noise = noise_buf[(uchar)(noise_off + 128)];
    break;
  }

  // When val is close to 0, channels are played normally
  // When val approaches +/-128, channels are swapped or muted
  if (amplified_val >= 0) {
    // Rotation to the right: left channel decreases, right channel receives
    // part of the left channel
    noise_l = (base_noise * (128 - pos_val)) >> 7;
    noise_r = base_noise + ((base_noise * pos_val) >> 7);
  } else {
    // Rotation to the left: right channel decreases, left channel receives part
    // of the right channel
    noise_l = base_noise + ((base_noise * pos_val) >> 7);
    noise_r = (base_noise * (128 - pos_val)) >> 7;
  }

  // Apply noise to the left and right channels
  *left = amp * noise_l;
  *right = amp * noise_r;
}

//
//	Play loop
//

void loop() {
  int c, cnt;
  int err; // Error to add to 'now' until next cnt==0
  int fast = fast_mult != 0;
  int vfast = fast_mult > 20; // Very fast - update status line often
  int now_lo = 0; // Low-order 16 bits of 'now' (fractional)
  int err_lo = 0;
  int ms_inc;

  setup_device();
  spin_carr_max = 127.0 / 1E-6 / out_rate;
  cnt = 1 + 1999 / out_buf_ms; // Update every 2 seconds or so
  now = fast_tim0;
  err = fast ? out_buf_ms * (fast_mult - 1) : 0;
  // Calculate the correct duration based on the last time in the sequence file
  int duration = t_per0(fast_tim0, fast_tim1);
  byte_count =
      out_bps * (S64)(duration * 0.001 * out_rate / (fast ? fast_mult : 1));

  // Do byte-swapping if bigendian and outputting to a file or stream
  if (out_mode == 1 && bigendian)
    out_mode = 2;

  if (opt_W) {
    writeWAV();
  }

  corrVal(0);          // Get into correct period
  dispCurrPer(stderr); // Display
  status(0);

  while (1) {
    for (c = 0; c < cnt; c++) {
      corrVal(1);
      outChunk();
      ms_inc = out_buf_ms + err;
      now_lo += out_buf_lo + err_lo;
      if (now_lo >= 0x10000) {
        ms_inc += now_lo >> 16;
        now_lo &= 0xFFFF;
      }
      now += ms_inc;
      if (now > H24)
        now -= H24;
      if (vfast && (c & 1))
        status(0);
    }
  }
}

//
//	Output a chunk of sound (a buffer-ful), then return
//
//	Note: Optimised for 16-bit output.  Eight-bit output is
//	slower, but then it probably won't have to run at as high a
//	sample rate.
//

int rand0, rand1;

void outChunk() {
  int off = 0;

  if (mix_in) {
    int rv = inbuf_read(tmp_buf, out_blen);
    if (rv == 0) {
      // Instead of ending, restart the background file
      if (mix_in && !feof(mix_in)) {
        // If it's not really EOF, it might be an error
        if (!opt_Q)
          warn("\nBackground sound: end of input audio stream");
        exit(0);
      }
      
      // Restart the background file for loop
      if (restart_background_file()) {
        // Try reading again after restarting
        rv = inbuf_read(tmp_buf, out_blen);
        if (rv == 0) {
          // If still can't read, fill with silence
          memset(tmp_buf, 0, out_blen * sizeof(int));
          rv = out_blen;
        }
      } else {
        // If can't restart, fill with silence
        memset(tmp_buf, 0, out_blen * sizeof(int));
        rv = out_blen;
      }
    }
    while (rv < out_blen)
      tmp_buf[rv++] = 0;
  }

  while (off < out_blen) {
    int ns = noise2(); // Use same pink noise source for everything
    int tot1, tot2;    // Left and right channels
    int mix1, mix2;    // Incoming mix signals
    int val, a;
    Channel *ch;

    mix1 = tmp_buf[off];
    mix2 = tmp_buf[off + 1];

    // Do default mixing at 100% if no mix/* stuff is present
    if (!mix_flag) {
      // tot1 = mix1 << 12;
      // tot2 = mix2 << 12;
      tot1 = (int)(mix1 * bg_gain_factor) << 12;
      tot2 = (int)(mix2 * bg_gain_factor) << 12;
    } else {
      tot1 = tot2 = 0;
    }

    ch = &chan[0];
    for (a = 0; a < N_CH; a++, ch++)
      switch (ch->typ) {
      case 0:
        break;
      case 1: // Binaural tones
        ch->off1 += ch->inc1;
        ch->off1 &= (ST_SIZ << 16) - 1;
        // tot1 += ch->amp * sin_table[ch->off1 >> 16];
        tot1 += ch->amp * sin_tables[ch->v.waveform][ch->off1 >> 16];
        ch->off2 += ch->inc2;
        ch->off2 &= (ST_SIZ << 16) - 1;
        // tot2 += ch->amp2 * sin_table[ch->off2 >> 16];
        tot2 += ch->amp2 * sin_tables[ch->v.waveform][ch->off2 >> 16];
        break;
      case 2: // Pink noise
        val = ns * ch->amp;
        tot1 += val;
        tot2 += val;
        break;
      case 9: // White noise
        val = white_noise() * ch->amp;
        tot1 += val;
        tot2 += val;
        break;
      case 10: // Brown noise
        val = brown_noise() * ch->amp;
        tot1 += val;
        tot2 += val;
        break;
      case 3: // Monaural beats
        {
          // Advance phases for both frequency components
          ch->off1 += ch->inc1;
          ch->off1 &= (ST_SIZ << 16) - 1;
          ch->off2 += ch->inc2;
          ch->off2 &= (ST_SIZ << 16) - 1;
          
          // Generate the two frequency components
          int freq_high = sin_tables[ch->v.waveform][ch->off1 >> 16]; // carr + res/2
          int freq_low = sin_tables[ch->v.waveform][ch->off2 >> 16];  // carr - res/2
          
          // Traditional monaural: opposite phase mixing
          // Left channel: high + low frequencies 
          // Right channel: low + high frequencies (same content, creates monaural effect)
          int left_mix = freq_high + freq_low;
          int right_mix = freq_low + freq_high; // Same as left, but keeping structure clear
          
          // Apply amplitude and divide by 2 to prevent clipping (only for audio output)
          // User still sees the full amplitude value in sprintVoice
          int half_amp = ch->amp / 2;
          
          tot1 += half_amp * left_mix;
          tot2 += half_amp * right_mix;
        }
        break;
      case 4: // Spinning pink noise
        ch->off1 += ch->inc1;
        ch->off1 &= (ST_SIZ << 16) - 1;
        // val= (ch->inc2 * sin_table[ch->off1 >> 16]) >> 24;
        val = (ch->inc2 * sin_tables[ch->v.waveform][ch->off1 >> 16]) >> 24;
        {
          int spin_left, spin_right;
          create_noise_spin_effect(4, ch->amp, val, &spin_left, &spin_right);
          tot1 += spin_left;
          tot2 += spin_right;
        }
        break;
      case 5: // Mix level
        // tot1 += mix1 * ch->amp;
        // tot2 += mix2 * ch->amp;
        {
          int bg_amp = (int)(ch->amp * bg_gain_factor);
          tot1 += mix1 * bg_amp;
          tot2 += mix2 * bg_amp;
        }
        break;
      case 6: // Mixspin - spinning mix stream
        ch->off1 += ch->inc1;
        ch->off1 &= (ST_SIZ << 16) - 1;
        // val= (ch->inc2 * sin_table[ch->off1 >> 16]) >> 24;
        val = (ch->inc2 * sin_tables[ch->v.waveform][ch->off1 >> 16]) >> 24;

        // Mixspin intensity control
        {
          // Calculate intensity factor based on amplitude
          // Amplitude varies from 0 to 4096 (0-100%)
          // Normalize to a factor between 0.5 and 4.0
          double intensity_factor = 0.5 + (ch->amp / 4096.0) * 3.5;

          // Apply intensity factor to rotation value
          int amplified_val = (int)(val * intensity_factor);

          // Limit value between -128 and 127
          if (amplified_val > 127)
            amplified_val = 127;
          if (amplified_val < -128)
            amplified_val = -128;

          // Use absolute value for calculations
          int pos_val = amplified_val < 0 ? -amplified_val : amplified_val;
          int mix_l, mix_r;

          // When val is close to 0, channels are played normally
          // When val approaches +/-128, channels are swapped or muted
          if (amplified_val >= 0) {
            // Rotation to the right: left channel decreases, right channel
            // receives part of the left channel
            mix_l = (mix1 * (128 - pos_val)) >> 7;
            mix_r = mix2 + ((mix1 * pos_val) >> 7);
          } else {
            // Rotation to the left: right channel decreases, left channel
            // receives part of the right channel
            mix_l = mix1 + ((mix2 * pos_val) >> 7);
            mix_r = (mix2 * (128 - pos_val)) >> 7;
          }

          int bg_mix_l = (int)(mix_l * bg_gain_factor);
          int bg_mix_r = (int)(mix_r * bg_gain_factor);

          // Apply base volume (using 70% of amplitude for volume)
          int base_amp = (int)(mix_amp != NULL ? *mix_amp : 4096.0) * 0.7 * bg_gain_factor;
          tot1 += base_amp * bg_mix_l;
          tot2 += base_amp * bg_mix_r;
        }
        break;
      case 7:                 // Mixpulse - mix stream with pulse effect
        ch->off2 += ch->inc2; // Modulator (pulse frequency)
        ch->off2 &= (ST_SIZ << 16) - 1;

        // Create the isochronic pulse effect in the audio stream
        {
          int mod_val = sin_tables[ch->v.waveform][ch->off2 >> 16];
          // Apply a threshold to create distinct pulses with space between them
          double mod_factor = 0.0;

          // Use only the positive part of the sine wave and apply a threshold
          // to create a space between pulses
          if (mod_val > ST_AMP * 0.3) { // Threshold of 30% of the maximum value
            // Normalize from ST_AMP*0.3 to ST_AMP to 0 to 1
            mod_factor = (mod_val - (ST_AMP * 0.3)) / (double)(ST_AMP * 0.7);
            // Smooth the edges of the pulse to avoid clicks
            mod_factor = mod_factor * mod_factor *
                         (3 - 2 * mod_factor); // Cubic smoothing
          }

          int bg_mix_1 = (int)(mix1 * bg_gain_factor);
          int bg_mix_2 = (int)(mix2 * bg_gain_factor);

          // Apply the modulation factor to the audio stream
          // Use base_amp as in mixspin (70% of amplitude for volume)
          int base_amp = (int)((mix_amp != NULL ? *mix_amp : 4096.0) * 0.7 * bg_gain_factor);

          // Calculate the intensity of the effect based on amplitude (ch->amp)
          // When ch->amp is 0, there is no effect (only the original audio)
          // When ch->amp is 4096 (100%), the effect is maximum
          double effect_intensity = (ch->amp / 4096.0) * 1.5;

          // Apply the modulation to the audio stream with variable intensity
          // When effect_intensity is 0, only the original audio is reproduced
          // When effect_intensity is 1, the audio is fully modulated by the
          // pulse
          double gain =
              (1.0 - effect_intensity) + (effect_intensity * mod_factor);
          tot1 += base_amp * bg_mix_1 * gain;
          tot2 += base_amp * bg_mix_2 * gain;
        }
        break;
      case 8:                 // Isochronic tones
        ch->off1 += ch->inc1; // Carrier (tone frequency)
        ch->off1 &= (ST_SIZ << 16) - 1;
        ch->off2 += ch->inc2; // Modulator (pulse frequency)
        ch->off2 &= (ST_SIZ << 16) - 1;

        // Change the modulator to create a true isochronic tone with space
        // between pulses
        {
          int mod_val = sin_tables[ch->v.waveform][ch->off2 >> 16];
          // Apply a threshold to create distinct pulses with space between them
          // Only produce sound when the modulation value is above a certain
          // threshold
          double mod_factor = 0.0;

          // Use only the positive part of the sine wave and apply a threshold
          // to create a space between pulses
          if (mod_val > ST_AMP * 0.3) { // Threshold of 30% of the maximum value
            // Normalize from ST_AMP*0.3 to ST_AMP to 0 to 1
            mod_factor = (mod_val - (ST_AMP * 0.3)) / (double)(ST_AMP * 0.7);
            // Smooth the edges of the pulse to avoid clicks
            mod_factor = mod_factor * mod_factor *
                         (3 - 2 * mod_factor); // Cubic smoothing
          }

          // Apply the modulation to the carrier
          // val = ch->amp * sin_table[ch->off1 >> 16] * mod_factor;
          val =
              ch->amp * sin_tables[ch->v.waveform][ch->off1 >> 16] * mod_factor;
          tot1 += val;
          tot2 += val;
        }
        break;
      case 11: // Bspin - spinning brown noise
        ch->off1 += ch->inc1;
        ch->off1 &= (ST_SIZ << 16) - 1;
        // val= (ch->inc2 * sin_table[ch->off1 >> 16]) >> 24;
        val = (ch->inc2 * sin_tables[ch->v.waveform][ch->off1 >> 16]) >> 24;

        {
          int spin_left, spin_right;
          create_noise_spin_effect(11, ch->amp, val, &spin_left, &spin_right);
          tot1 += spin_left;
          tot2 += spin_right;
        }
        break;
      case 12: // Wspin - spinning white noise
        ch->off1 += ch->inc1;
        ch->off1 &= (ST_SIZ << 16) - 1;
        // val= (ch->inc2 * sin_table[ch->off1 >> 16]) >> 24;
        val = (ch->inc2 * sin_tables[ch->v.waveform][ch->off1 >> 16]) >> 24;

        {
          int spin_left, spin_right;
          create_noise_spin_effect(12, ch->amp, val, &spin_left, &spin_right);
          tot1 += spin_left;
          tot2 += spin_right;
        }
        break;
      default:
        break;
      }

    // Apply volume level
    if (opt_V != 100) {
      tot1 = ((long long)tot1 * opt_V + 50) / 100;
      tot2 = ((long long)tot2 * opt_V + 50) / 100;
    }

    // White noise dither; you could also try (rand0-rand1) for a
    // dither with more high frequencies
    rand0 = rand1;
    rand1 = (rand0 * 0x660D + 0xF35F) & 0xFFFF;
    if (tot1 <= 0x7FFF0000)
      tot1 += rand0;
    if (tot2 <= 0x7FFF0000)
      tot2 += rand0;

    out_buf[off++] = tot1 >> 16;
    out_buf[off++] = tot2 >> 16;
  }

  // Rewrite buffer for 8-bit mode
  if (out_mode == 0) {
    short *sp = out_buf;
    short *end = out_buf + out_blen;
    char *cp = (char *)out_buf;
    while (sp < end)
      *cp++ = (*sp++ >> 8) + 128;
  }

  // Rewrite buffer for 16-bit byte-swapping
  if (out_mode == 2) {
    char *cp = (char *)out_buf;
    char *end = (char *)(out_buf + out_blen);
    while (cp < end) {
      char tmp = *cp++;
      cp[-1] = cp[0];
      *cp++ = tmp;
    }
  }

  // Check and update the byte count if necessary
  if (byte_count > 0) {
    if (byte_count <= out_bsiz) {
      writeOut((char *)out_buf, byte_count);
      exit(0); // All done
    } else {
      writeOut((char *)out_buf, out_bsiz);
      byte_count -= out_bsiz;
    }
  } else
    writeOut((char *)out_buf, out_bsiz);
}

void writeOut(char *buf, int siz) {
  int rv;

  while (-1 != (rv = write(out_fd, buf, siz))) {
    if (0 == (siz -= rv))
      return;
    buf += rv;
  }
  error("Output error");
}

//
//	Correct channel values and types according to current period,
//	and current time
//

void corrVal(int running) {
  int a;
  int t0 = per->tim;
  int t1 = per->nxt->tim;
  Channel *ch;
  Voice *v0, *v1, *vv;
  double rat0, rat1;

  // Move to the correct period
  while ((now >= t0) ^ (now >= t1) ^ (t1 > t0)) {
    per = per->nxt;
    t0 = per->tim;
    t1 = per->nxt->tim;
    if (running) {
      if (tty_erase) {
        fprintf(stderr, "%*s\r", tty_erase, "");
        tty_erase = 0;
      }
      dispCurrPer(stderr);
      status(0);
    }
  }

  // Run through to calculate voice settings for current time
  rat1 = t_per0(t0, now) / (double)t_per24(t0, t1);
  rat0 = 1 - rat1;
  for (a = 0; a < N_CH; a++) {
    ch = &chan[a];
    v0 = &per->v0[a];
    v1 = &per->v1[a];
    vv = &ch->v;

    // Pointer to the amplitude of the mix to use with mixspin/mixpulse
    if (vv->typ == 5 && mix_amp == NULL)
      mix_amp = &vv->amp;

    if (vv->typ != v0->typ) {
      switch (vv->typ = ch->typ = v0->typ) {
      case 1:
        ch->off1 = ch->off2 = 0;
        break;
      case 2:
        break;
      case 3: // Monaural beats
        ch->off1 = ch->off2 = 0;
        break;
      case 4:
        ch->off1 = ch->off2 = 0;
        break;
      case 5:
        break;
      case 8: // Isochronic tones
        ch->off1 = ch->off2 = 0;
        break;
      case 6: // Mixspin
        ch->off1 = ch->off2 = 0;
        break;
      case 7: // Mixpulse
        ch->off1 = ch->off2 = 0;
        break;
      case 11: // Bspin - spinning brown noise
        ch->off1 = ch->off2 = 0;
        break;
      case 12: // Wspin - spinning white noise
        ch->off1 = ch->off2 = 0;
        break;
      default:
        ch->off1 = ch->off2 = 0;
        break;
      }
    }

    // Setup vv->*
    switch (vv->typ) {
    case 1:
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      vv->waveform = v0->waveform;
      break;
    case 2:
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->waveform = v0->waveform;
      break;
    case 3:
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      vv->waveform = v0->waveform;
      break;
    case 4:
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      if (vv->carr > spin_carr_max)
        vv->carr = spin_carr_max; // Clipping sweep width
      if (vv->carr < -spin_carr_max)
        vv->carr = -spin_carr_max;
      vv->waveform = v0->waveform;
      break;
    case 5:
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->waveform = v0->waveform;
      break;
    case 8: // Isochronic tones
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      vv->waveform = v0->waveform;
      break;
    case 6: // Mixspin
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      if (vv->carr > spin_carr_max)
        vv->carr = spin_carr_max; // Clipping sweep width
      if (vv->carr < -spin_carr_max)
        vv->carr = -spin_carr_max;
      vv->waveform = v0->waveform;
      break;
    case 7: // Mixpulse
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      vv->waveform = v0->waveform;
      break;
    case 11: // Bspin - spinning brown noise
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      if (vv->carr > spin_carr_max)
        vv->carr = spin_carr_max; // Clipping sweep width
      if (vv->carr < -spin_carr_max)
        vv->carr = -spin_carr_max;
      vv->waveform = v0->waveform;
      break;
    case 12: // Wspin - spinning white noise
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      if (vv->carr > spin_carr_max)
        vv->carr = spin_carr_max; // Clipping sweep width
      if (vv->carr < -spin_carr_max)
        vv->carr = -spin_carr_max;
      vv->waveform = v0->waveform;
      break;
    default: // Waveform based binaural
      vv->amp = rat0 * v0->amp + rat1 * v1->amp;
      vv->carr = rat0 * v0->carr + rat1 * v1->carr;
      vv->res = rat0 * v0->res + rat1 * v1->res;
      break;
    }
  }

  // Setup Channel data from Voice data
  for (a = 0; a < N_CH; a++) {
    ch = &chan[a];
    vv = &ch->v;

    // Setup ch->* from vv->*
    switch (vv->typ) {
      double freq1, freq2;
    case 1:
      freq1 = vv->carr + vv->res / 2;
      freq2 = vv->carr - vv->res / 2;
      ch->amp = ch->amp2 = (int)vv->amp;
      ch->inc1 = (int)(freq1 / out_rate * ST_SIZ * 65536);
      ch->inc2 = (int)(freq2 / out_rate * ST_SIZ * 65536);
      break;
    case 2:
    case 9:
    case 10:
      ch->amp = (int)vv->amp;
      break;
    case 3:
      ch->amp = (int)vv->amp;
      ch->inc1 = (int)((vv->carr + vv->res / 2) / out_rate * ST_SIZ * 65536); // freq high
      ch->inc2 = (int)((vv->carr - vv->res / 2) / out_rate * ST_SIZ * 65536); // freq low
      break;
    case 4:
      ch->amp = (int)vv->amp;
      ch->inc1 = (int)(vv->res / out_rate * ST_SIZ * 65536);
      ch->inc2 = (int)(vv->carr * 1E-6 * out_rate * (1 << 24) / ST_AMP);
      break;
    case 5:
      ch->amp = (int)vv->amp;
      break;
    case 8: // Isochronic tones
      ch->amp = (int)vv->amp;
      // Carrier (tone frequency)
      ch->inc1 = (int)(vv->carr / out_rate * ST_SIZ * 65536);
      // Modulator (pulse frequency)
      ch->inc2 = (int)(vv->res / out_rate * ST_SIZ * 65536);
      break;
    case 6: // Effect spin
      ch->amp = (int)vv->amp;
      ch->inc1 = (int)(vv->res / out_rate * ST_SIZ * 65536);
      ch->inc2 = (int)(vv->carr * 1E-6 * out_rate * (1 << 24) / ST_AMP);
      break;
    case 7: // Effect pulse
      ch->amp = (int)vv->amp;
      // Modulator (pulse frequency)
      ch->inc2 = (int)(vv->res / out_rate * ST_SIZ * 65536);
      break;
    case 11: // Spin - spinning brown noise
      ch->amp = (int)vv->amp;
      ch->inc1 = (int)(vv->res / out_rate * ST_SIZ * 65536);
      ch->inc2 = (int)(vv->carr * 1E-6 * out_rate * (1 << 24) / ST_AMP);
      break;
    case 12: // Spin - spinning white noise
      ch->amp = (int)vv->amp;
      ch->inc1 = (int)(vv->res / out_rate * ST_SIZ * 65536);
      ch->inc2 = (int)(vv->carr * 1E-6 * out_rate * (1 << 24) / ST_AMP);
      break;
    default:
      break;
    }
  }
}

//
//	Setup audio device
//

void setup_device(void) {
  // Handle output to files and pipes
  if (opt_o) {
    FILE *out;
    if (!(out = fopen(opt_o, "wb")))
      error("Can't open output file: %s", opt_o);
    out_fd = fileno(out);
  }

  // Handle output to files and pipes
  out_blen = out_rate * 2 / out_prate; // 10 fragments a second by default
  while (out_blen & (out_blen - 1))
    out_blen &= out_blen - 1; // Make power of two
  out_bsiz = out_blen * (out_mode ? 2 : 1);
  out_bps = out_mode ? 4 : 2;
  out_buf = (short *)Alloc(out_blen * sizeof(short));
  out_buf_lo = (int)(0x10000 * 1000.0 * 0.5 * out_blen / out_rate);
  out_buf_ms = out_buf_lo >> 16;
  out_buf_lo &= 0xFFFF;
  tmp_buf = (int *)Alloc(out_blen * sizeof(int));
}

//
//	Write a WAV header, and setup out_mode if byte-swapping is
//	required.  'byte_count' should have been set up by this point.
//

#define addU4(xx)                                                              \
  {                                                                            \
    int a = xx;                                                                \
    *p++ = a;                                                                  \
    *p++ = (a >>= 8);                                                          \
    *p++ = (a >>= 8);                                                          \
    *p++ = (a >>= 8);                                                          \
  }
#define addStr(xx)                                                             \
  {                                                                            \
    char *q = xx;                                                              \
    *p++ = *q++;                                                               \
    *p++ = *q++;                                                               \
    *p++ = *q++;                                                               \
    *p++ = *q++;                                                               \
  }

void writeWAV() {
  char buf[44], *p = buf;

  if (byte_count + 36 != (int)(byte_count + 36)) {
    int tmp;
    byte_count = 0xFFFFFFF8 - 36;
    tmp = byte_count / out_bps / out_rate;
    warn("WARNING: Selected length is too long for the WAV format; truncating "
         "to %dh%02dm%02ds",
         tmp / 3600, tmp / 60 % 60, tmp % 60);
  }

  addStr("RIFF");
  addU4(byte_count + 36);
  addStr("WAVE");
  addStr("fmt ");
  addU4(16);
  addU4(0x00020001);
  addU4(out_rate);
  addU4(out_rate * out_bps);
  addU4(0x0004 + 0x10000 * (out_bps * 4)); // 2,4 -> 8,16 - always assume stereo
  addStr("data");
  addU4(byte_count);
  writeOut(buf, 44);
}

//
//	Read a line, discarding blank lines and comments.  Rets:
//	Another line?  Comments starting with '##' are displayed on
//	stderr.
//

int readLine() {
  char *p;

  // Check if we have a saved line from readNameDef
  if (saved_lin_num >= 0) {
    // Use the saved line
    strcpy(buf, saved_buf);
    strcpy(buf_copy, saved_buf_copy);
    lin = buf;
    lin_copy = buf_copy;
    in_lin = saved_lin_num;
    saved_lin_num = -1; // Clear the saved line
    return 1;
  }

  lin = buf;

  while (1) {
    if (!fgets(lin, sizeof(buf), in)) {
      if (feof(in))
        return 0;
      error("Read error on sequence file");
    }

    in_lin++;

    while (isspace(*lin))
      lin++;
    p = strchr(lin, '#');
    if (p && p[1] == '#')
      fprintf(stderr, "> %s", p + 2);
    p = p ? p : strchr(lin, 0);
    while (p > lin && isspace(p[-1]))
      p--;
    if (p != lin)
      break;
  }
  *p = 0;
  lin_copy = buf_copy;
  strcpy(lin_copy, lin);
  return 1;
}

//
//	Get next word at '*lin', moving lin onwards, or return 0
//

char *getWord() {
  char *rv, *end;
  while (isspace(*lin))
    lin++;
  if (!*lin)
    return 0;

  rv = lin;
  while (*lin && !isspace(*lin))
    lin++;
  end = lin;
  if (*lin)
    lin++;
  *end = 0;

  return rv;
}

//
//	Bad sequence file
//

void badSeq() {
  error("Error in sequence file at line: %d\n  %s", in_lin, lin_copy);
}

//
//	Read a list of sequence files, and generate a list of Period
//	structures
//

void readSeq(int ac, char **av) {
  // Setup a 'now' value to use for NOW in the sequence file
  now = calcNow();

  while (ac-- > 0) {
    char *fnam = *av++;
    int start = 1;

    in = (0 == strcmp("-", fnam)) ? stdin : fopen(fnam, "r");
    if (!in)
      error("Error opening sequence file: %s", fnam);

    in_lin = 0;

    while (readLine()) {
      char *p = lin;

      // Blank lines
      if (!*p)
        continue;

      // Look for options
      if (*p == '@') {
        if (!start) {
          error("Options are only permitted at start of sequence file:\n  %s",
                p);
        }
        handleOptionInSequence(p);
        continue;
      }

      // Check if it's a name definition (no colon, just name)
      start = 0;
      if (isalpha(*p)) { // new syntax
        while (isalnum(*p) || *p == '_' || *p == '-')
          p++;
        // If line ends with just the name (no colon), it's a new syntax
        // definition
        if (*p == 0) {
          // Pure name definition - new syntax
          readNameDef();
        } else if (isspace(*p)) {
          // Name followed by spaces - check if there's more content
          while (isspace(*p))
            p++;
          if (*p != 0) {
            // There's content after the name - this is invalid for new syntax
            error("Invalid syntax at line %d: %s", in_lin, lin_copy);
          }
          // Just the name with trailing spaces - new syntax
          readNameDef();
        } else {
          // Must be a timeline entry (contains colon or other chars)
          readTimeLine();
        }
      } else {
        readTimeLine();
      }
    }

    if (in != stdin)
      fclose(in);
  }

  correctPeriods();
}

//
//	Fill in all the correct information for the Periods, assuming
//	they have just been loaded using readTimeLine()
//

void correctPeriods() {
  // Get times all correct
  {
    Period *pp = per;
    do {
      if (pp->fi == -2) {
        pp->tim = pp->nxt->tim;
        pp->fi = -1;
      }

      pp = pp->nxt;
    } while (pp != per);
  }

  // Make sure that the transitional periods each have enough time
  {
    Period *pp = per;
    do {
      if (pp->fi == -1) {
        int per = t_per0(pp->tim, pp->nxt->tim);
        if (per < fade_int) {
          int adj = (fade_int - per) / 2, adj0, adj1;
          adj0 = t_per0(pp->prv->tim, pp->tim);
          adj0 = (adj < adj0) ? adj : adj0;
          adj1 = t_per0(pp->nxt->tim, pp->nxt->nxt->tim);
          adj1 = (adj < adj1) ? adj : adj1;
          pp->tim = (pp->tim - adj0 + H24) % H24;
          pp->nxt->tim = (pp->nxt->tim + adj1) % H24;
        }
      }

      pp = pp->nxt;
    } while (pp != per);
  }

  // Fill in all the voice arrays, and sort out details of
  // transitional periods
  {
    Period *pp = per;
    do {
      if (pp->fi < 0) {
        int fo, fi;
        int a;
        int midpt = 0;

        Period *qq = (Period *)Alloc(sizeof(*qq));
        qq->prv = pp;
        qq->nxt = pp->nxt;
        qq->prv->nxt = qq->nxt->prv = qq;

        qq->tim = t_mid(pp->tim, qq->nxt->tim);

        memcpy(pp->v0, pp->prv->v1, sizeof(pp->v0));
        memcpy(qq->v1, qq->nxt->v0, sizeof(qq->v1));

        fo = pp->prv->fo;
        fi = qq->nxt->fi;

        // Special handling for -> slides:
        //   always slide, and stretch slide if possible
        if (pp->fi == -3) {
          fo = fi = 2; // Force slides for ->
          for (a = 0; a < N_CH; a++) {
            Voice *vp = &pp->v0[a];
            Voice *vq = &qq->v1[a];
            if (vp->typ == 0 && vq->typ != 0) {
              memcpy(vp, vq, sizeof(*vp));
              vp->amp = 0;
            } else if (vp->typ != 0 && vq->typ == 0) {
              memcpy(vq, vp, sizeof(*vq));
              vq->amp = 0;
            }
          }
        }

        memcpy(pp->v1, pp->v0, sizeof(pp->v1));
        memcpy(qq->v0, qq->v1, sizeof(qq->v0));

        for (a = 0; a < N_CH; a++) {
          Voice *vp = &pp->v1[a];
          Voice *vq = &qq->v0[a];
          if ((fo == 0 || fi == 0) ||           // Fade in/out to silence
              (vp->typ != vq->typ) ||           // Different types
              (vp->waveform != vq->waveform) || // Different waveforms
              ((fo == 1 || fi == 1) && // Fade thru, but different pitches
               (vp->typ == 1 || vp->typ < 0) &&
               (vp->carr != vq->carr || vp->res != vq->res))) {
            vp->amp = vq->amp = 0; // To silence
            midpt = 1;             // Definitely need the mid-point
          } else { // Else smooth transition
            vp->amp = vq->amp = (vp->amp + vq->amp) / 2;
            if (vp->typ == 1 || vp->typ == 4 || vp->typ < 0) {
              vp->carr = vq->carr = (vp->carr + vq->carr) / 2;
              vp->res = vq->res = (vp->res + vq->res) / 2;
            }
          }
        }

        // If we don't really need the mid-point, then get rid of it
        if (!midpt) {
          memcpy(pp->v1, qq->v1, sizeof(pp->v1));
          qq->prv->nxt = qq->nxt;
          qq->nxt->prv = qq->prv;
          free(qq);
        } else
          pp = qq;
      }

      pp = pp->nxt;
    } while (pp != per);
  }

  // NEW VALIDATION: Mandatory linear sequences - MUST RUN BEFORE CLEANUP
  {
    // Must have at least 2 periods (start and end)
    if (per->nxt == per) {
      error("Sequence must have at least a start and end time.");
    }

    // Validate strict chronological order
    if (fast_tim0 >= 0 && fast_tim1 >= 0) {
      if (fast_tim0 >= fast_tim1) {
        error("Times out of chronological order.\n"
              "First time: %02d:%02d:%02d\n"
              "Last time: %02d:%02d:%02d\n"
              "Last time must be greater than first time.",
              fast_tim0 / 3600000, (fast_tim0 / 60000) % 60,
              (fast_tim0 / 1000) % 60, fast_tim1 / 3600000,
              (fast_tim1 / 60000) % 60, (fast_tim1 / 1000) % 60);
      }
    }

    // Collect all periods into an array for validation
    Period *periods[1000]; // Max 1000 periods
    int period_count = 0;
    Period *pp = per;

    // Collect all periods in ORIGINAL order (to check chronological order)
    do {
      if (period_count >= 1000) {
        error("Too many periods (max 1000)");
      }
      periods[period_count++] = pp;
      pp = pp->nxt;
    } while (pp != per);

    // First validate that original order is chronological
    // Remove duplicates but preserve original order for validation
    Period *unique_original[1000];
    int unique_original_count = 0;

    for (int i = 0; i < period_count; i++) {
      // Check if this time already exists in unique list
      int already_exists = 0;
      for (int j = 0; j < unique_original_count; j++) {
        if (unique_original[j]->tim == periods[i]->tim) {
          already_exists = 1;
          break;
        }
      }
      if (!already_exists) {
        unique_original[unique_original_count++] = periods[i];
      }
    }

    // Validate original chronological order
    for (int i = 1; i < unique_original_count; i++) {
      if (unique_original[i]->tim <= unique_original[i - 1]->tim) {
        if (unique_original[i]->tim == unique_original[i - 1]->tim) {
          error("Duplicate time found: %02d:%02d:%02d\n"
                "Each time in sequence must be unique.",
                unique_original[i]->tim / 3600000,
                (unique_original[i]->tim / 60000) % 60,
                (unique_original[i]->tim / 1000) % 60);
        } else {
          error("Times out of chronological order: %02d:%02d:%02d comes after "
                "%02d:%02d:%02d\n"
                "Times in sequence file must be written in ascending "
                "chronological order.",
                unique_original[i]->tim / 3600000,
                (unique_original[i]->tim / 60000) % 60,
                (unique_original[i]->tim / 1000) % 60,
                unique_original[i - 1]->tim / 3600000,
                (unique_original[i - 1]->tim / 60000) % 60,
                (unique_original[i - 1]->tim / 1000) % 60);
        }
      }
    }

    // Now sort periods by time for further processing
    for (int i = 0; i < period_count - 1; i++) {
      for (int j = 0; j < period_count - i - 1; j++) {
        if (periods[j]->tim > periods[j + 1]->tim) {
          Period *temp = periods[j];
          periods[j] = periods[j + 1];
          periods[j + 1] = temp;
        }
      }
    }

    // Remove duplicate times - keep only unique times for validation
    int unique_count = 0;
    Period *unique_periods[1000];
    unique_periods[0] = periods[0];
    unique_count = 1;

    for (int i = 1; i < period_count; i++) {
      // Only add if time is different from previous unique time
      if (periods[i]->tim != unique_periods[unique_count - 1]->tim) {
        unique_periods[unique_count++] = periods[i];
      }
    }

    // Update arrays to use unique periods only
    period_count = unique_count;
    for (int i = 0; i < unique_count; i++) {
      periods[i] = unique_periods[i];
    }

    // Validate that sequence starts at 00:00:00
    if (periods[0]->tim != 0) {
      error("Sequence must start at 00:00:00.\n"
            "First time found: %02d:%02d:%02d\n"
            "Add a period starting at 00:00:00 to your sequence.",
            periods[0]->tim / 3600000, (periods[0]->tim / 60000) % 60,
            (periods[0]->tim / 1000) % 60);
    }

    // Validate chronological order and minimum intervals
    for (int i = 1; i < period_count; i++) {
      int time_diff = periods[i]->tim - periods[i - 1]->tim;

      if (time_diff < 0) {
        error("Time %02d:%02d:%02d cannot come after previous time.\n"
              "All times must be in chronological ascending order.",
              periods[i]->tim / 3600000, (periods[i]->tim / 60000) % 60,
              (periods[i]->tim / 1000) % 60);
      }
    }
  }

  // Clear out zero length sections, and duplicate sections (AFTER validation)
  {
    Period *pp;
    while (per != per->nxt) {
      pp = per;
      do {
        if (voicesEq(pp->v0, pp->v1) && voicesEq(pp->v0, pp->nxt->v0) &&
            voicesEq(pp->v0, pp->nxt->v1))
          pp->nxt->tim = pp->tim;

        if (pp->tim == pp->nxt->tim) {
          if (per == pp)
            per = per->prv;
          pp->prv->nxt = pp->nxt;
          pp->nxt->prv = pp->prv;
          free(pp);
          pp = 0;
          break;
        }
        pp = pp->nxt;
      } while (pp != per);
      if (pp)
        break;
    }
  }

  // Print the whole lot out
  if (opt_D) {
    Period *pp;
    if (per->nxt != per)
      while (per->prv->tim < per->tim)
        per = per->nxt;

    pp = per;

    fprintf(stderr, "\n*** This is a test mode. Use --output to generate the audio file. ***\n\n");

    do {
      dispCurrPer(stdout);
      per = per->nxt;
    } while (per != pp);
    printf("\n");

    exit(0); // All done
  }
}

int voicesEq(Voice *v0, Voice *v1) {
  int a = N_CH;

  while (a-- > 0) {
    if (v0->typ != v1->typ)
      return 0;
    switch (v0->typ) {
    case 1:
    case 4:
    default:
      if (v0->amp != v1->amp || v0->carr != v1->carr || v0->res != v1->res)
        return 0;
      break;
    case 2:
    case 5:
      if (v0->amp != v1->amp)
        return 0;
      break;
    }
    v0++;
    v1++;
  }
  return 1;
}

//
//	Read a name definition
//

void readNameDef() {
  char *p, *q;
  NameDef *nd;
  int ch = 0;

  // Get the name (no colon required)
  if (!(p = getWord()))
    badSeq();

  // Validate name characters
  for (q = p; *q; q++)
    if (!isalnum(*q) && *q != '-' && *q != '_')
      error("Invalid name \"%s\" in preset, line %d:\n  %s", p, in_lin,
            lin_copy);

  if (strcmp(p, "silence") == 0)
    error("Cannot redefine built-in name 'silence' at line %d.", in_lin);

  // Create NameDef
  nd = (NameDef *)Alloc(sizeof(NameDef));
  nd->name = StrDup(p);

  // Initialize all voices to off
  for (int i = 0; i < N_CH; i++) {
    nd->vv[i].typ = 0;
    nd->vv[i].amp = 0;
    nd->vv[i].carr = 0;
    nd->vv[i].res = 0;
    nd->vv[i].waveform = 0;
  }

  // Add to list immediately after creating the basic structure
  nd->nxt = nlist;
  nlist = nd;

  // Read multi-line definition in new syntax
  // Use a special version of readLine that preserves indentation
  int lines_processed = 0; // Count of indented lines processed

  while (1) {
    // Read line preserving indentation
    if (!fgets(buf, sizeof(buf), in)) {
      // End of file - check if we processed any lines
      if (lines_processed == 0) {
        error("Empty definition for '%s' at end of file. Name definitions must "
              "have at least one indented line.\n",
              nd->name);
      }
      break; // End of file
    }
    in_lin++;

    // Remove trailing whitespace and comments
    char *p = strchr(buf, '#');
    if (p && p[1] == '#')
      fprintf(stderr, "> %s", p + 2);
    p = p ? p : strchr(buf, 0);
    while (p > buf && isspace(p[-1]))
      p--;
    *p = 0;

    // Check if it's an empty line
    char *temp = buf;
    while (isspace(*temp))
      temp++;
    if (!*temp)
      continue; // Skip empty lines

    // Check indentation
    if (buf[0] == ' ') {
      // Line starts with space - validate indentation
      if (buf[1] != ' ') {
        error("Invalid indentation at line %d. Definition lines must have "
              "exactly 2 spaces.\n",
              in_lin);
      }
      if (buf[2] == ' ') {
        error("Invalid indentation at line %d. Definition lines must have "
              "exactly 2 spaces.\n",
              in_lin);
      }
      // Valid 2-space indentation - continue processing below
    } else if (buf[0] != ' ') {
      // Line doesn't start with space, save it for later processing
      strcpy(saved_buf, buf);
      strcpy(saved_buf_copy, buf_copy);
      saved_lin_num = in_lin;

      // Check if we processed any indented lines
      if (lines_processed == 0) {
        error("Empty definition for '%s' at line %d. Name definitions must "
              "have at least one indented line.\n",
              nd->name, in_lin, nd->name);
      }

      // Normalize amplitude before finishing
      normalizeAmplitude(nd->vv, N_CH, buf_copy, in_lin);

      return; // Exit function, saved line will be processed by main loop
    }

    // Skip the two spaces and parse the command
    lin = buf + 2;
    lin_copy = buf_copy;
    strcpy(lin_copy, lin);

    // Get the command type
    char *cmd = getWord();
    if (!cmd)
      continue; // Skip empty lines

    if (ch >= N_CH) {
      error("Too many voice definitions in '%s' (max %d)", nd->name, N_CH);
    }

    if (strcmp(cmd, "noise") == 0) {
      // Parse: noise <type> amplitude <value>
      char *type = getWord();
      char *amp_str = getWord();
      char *amp_value = getWord();

      if (!type || !amp_str || !amp_value ||
          strcmp(amp_str, "amplitude") != 0) {
        error("Invalid noise syntax at line %d. Expected: noise <type> "
              "amplitude <value>\n  %s",
              in_lin, lin_copy);
      }

      double amp;

      if (sscanf(amp_value, "%lf", &amp) != 1) {
        error("Invalid noise amplitude at line %d.\n  %s", in_lin, lin_copy);
      }

      if (strcmp(type, "pink") == 0) {
        nd->vv[ch].typ = 2;
      } else if (strcmp(type, "white") == 0) {
        nd->vv[ch].typ = 9;
      } else if (strcmp(type, "brown") == 0) {
        nd->vv[ch].typ = 10;
      } else {
        error("Unknown noise type '%s' at line %d. Use: pink, white, brown",
              type, in_lin);
      }

      // Check if there are extra words after the command
      char *next_word = getWord();
      if (next_word) {
        error("Invalid syntax at line %d. Expected: noise <type> amplitude "
              "<value>\n  %s",
              in_lin, lin_copy);
        free(next_word);
      }

      if (amp < 0 || amp > 100) {
        error("Invalid noise amplitude at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
      }

      nd->vv[ch].amp = AMP_DA(amp);
      ch++;
      lines_processed++; // Increment count of processed lines

    } else if (strcmp(cmd, "tone") == 0) {
      // Parse: tone <freq> <type> <value> amplitude <amp>
      char *freq_str = getWord();
      char *type = getWord();
      char *value_str = getWord();
      char *amp_str = getWord();
      char *amp_value = getWord();

      if (!freq_str || !type || !value_str || !amp_str || !amp_value ||
          strcmp(amp_str, "amplitude") != 0) {
        error("Invalid tone syntax at line %d. Expected: tone <freq> <type> "
              "<value> amplitude <amp>\n  %s",
              in_lin, lin_copy);
      }

      double freq, value, amp;

      if (sscanf(freq_str, "%lf", &freq) != 1) {
        error("Invalid tone frequency at line %d.\n  %s", in_lin, lin_copy);
      }
      if (sscanf(value_str, "%lf", &value) != 1) {
        error("Invalid tone value at line %d.\n  %s", in_lin, lin_copy);
      }
      if (sscanf(amp_value, "%lf", &amp) != 1) {
        error("Invalid tone amplitude at line %d.\n  %s", in_lin, lin_copy);
      }

      if (strcmp(type, "binaural") == 0) {
        // Binaural beat: freq+value/amp
        nd->vv[ch].typ = 1;
        nd->vv[ch].carr = freq;
        nd->vv[ch].res = value;
      } else if (strcmp(type, "isochronic") == 0) {
        // Isochronic pulse: freq@value/amp
        nd->vv[ch].typ = 8;
        nd->vv[ch].carr = freq;
        nd->vv[ch].res = value;
      } else if (strcmp(type, "monaural") == 0) {
        nd->vv[ch].typ = 3;
        nd->vv[ch].carr = freq;
        nd->vv[ch].res = value;
      }
      else {
        error("Unknown tone type '%s' at line %d. Use: binaural, "
              "monaural, isochronic",
              type, in_lin);
      }

      char *next_word = getWord();
      if (next_word) {
        error("Invalid syntax at line %d. Expected: tone <freq> <type> "
              "<value> amplitude <amp>\n  %s",
              in_lin, lin_copy);
        free(next_word);
      }

      if (freq < 0) {
        error("Invalid tone frequency at line %d.\n  %s", in_lin, lin_copy);
      }
      if (value < 0) {
        error("Invalid tone value at line %d.\n  %s", in_lin, lin_copy);
      }
      if (amp < 0 || amp > 100) {
        error("Invalid tone amplitude at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
      }

      nd->vv[ch].amp = AMP_DA(amp);
      ch++;
      lines_processed++; // Increment count of processed lines
    } else if (strcmp(cmd, "waveform") == 0) {
      // Parse: waveform <type> <value> amplitude <amp>
      char *waveform_type = getWord();

      if (strcmp(waveform_type, "sine") == 0) {
        nd->vv[ch].waveform = 0;
      } else if (strcmp(waveform_type, "square") == 0) {
        nd->vv[ch].waveform = 1;
      } else if (strcmp(waveform_type, "triangle") == 0) {
        nd->vv[ch].waveform = 2;
      } else if (strcmp(waveform_type, "sawtooth") == 0) {
        nd->vv[ch].waveform = 3;
      } else {
        error("Unknown waveform type '%s' at line %d. Use: sine, square, "
              "triangle, sawtooth",
              waveform_type, in_lin);
      }

      char *typ_nxt = getWord();
      if (strcmp(typ_nxt, "tone") == 0) {
        char *freq_str = getWord();
        char *type = getWord();
        char *value_str = getWord();
        char *amp_str = getWord();
        char *amp_value = getWord();

        if (!freq_str || !type || !value_str || !amp_str || !amp_value ||
            strcmp(amp_str, "amplitude") != 0) {
          error("Invalid tone syntax at line %d. Expected: tone <freq> <type> "
                "<value> amplitude <amp>\n  %s",
                in_lin, lin_copy);
        }

        double freq, value, amp;

        if (sscanf(freq_str, "%lf", &freq) != 1) {
          error("Invalid tone frequency at line %d.\n  %s", in_lin, lin_copy);
        }
        if (sscanf(value_str, "%lf", &value) != 1) {
          error("Invalid tone value at line %d.\n  %s", in_lin, lin_copy);
        }
        if (sscanf(amp_value, "%lf", &amp) != 1) {
          error("Invalid tone amplitude at line %d.\n  %s", in_lin, lin_copy);
        }
        if (strcmp(type, "binaural") == 0) {
          // Binaural beat: freq+value/amp
          nd->vv[ch].typ = 1;
          nd->vv[ch].carr = freq;
          nd->vv[ch].res = value;
        } else if (strcmp(type, "isochronic") == 0) {
          // Isochronic pulse: freq@value/amp
          nd->vv[ch].typ = 8;
          nd->vv[ch].carr = freq;
          nd->vv[ch].res = value;
        } else if (strcmp(type, "monaural") == 0) {
          nd->vv[ch].typ = 3;
          nd->vv[ch].carr = freq;
          nd->vv[ch].res = value;
        }
        else {
          error("Unknown tone type '%s' at line %d. Use: binaural, "
                "monaural, isochronic",
                type, in_lin);
        }

        char *next_word = getWord();
        if (next_word) {
          error("Invalid syntax at line %d. Expected: tone <freq> <type> "
                "<value> amplitude <amp>\n  %s",
                in_lin, lin_copy);
          free(next_word);
        }

        if (freq < 0) {
          error("Invalid tone frequency at line %d.\n  %s", in_lin, lin_copy);
        }
        if (value < 0) {
          error("Invalid tone value at line %d.\n  %s", in_lin, lin_copy);
        }
        if (amp < 0 || amp > 100) {
          error("Invalid tone amplitude at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
        }

        nd->vv[ch].amp = AMP_DA(amp);
        ch++;
        lines_processed++; // Increment count of processed lines
      } else if (strcmp(typ_nxt, "spin") == 0) {
        // Spin: <type> width <width> beat <beat> amplitude <amp>
        char *type = getWord();
        char *width_str = getWord();
        char *width_value = getWord();
        char *rate_str = getWord();
        char *rate_value = getWord();
        char *amp_str = getWord();
        char *amp_value = getWord();

        if (!type || !width_str || !width_value || !rate_str || !rate_value || !amp_str || !amp_value ||
            strcmp(amp_str, "amplitude") != 0 ||
            strcmp(rate_str, "rate") != 0 ||
            strcmp(width_str, "width") != 0) {
          error("Invalid spin syntax at line %d. Expected: spin <type> width "
                "<width> rate <rate> amplitude <amp>\n  %s",
                in_lin, lin_copy);
        }

        double width, rate, amp;

        if (sscanf(width_value, "%lf", &width) != 1) {
          error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
        }
        if (sscanf(rate_value, "%lf", &rate) != 1) {
          error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
        }
        if (sscanf(amp_value, "%lf", &amp) != 1) {
          error("Invalid spin amplitude at line %d.\n  %s", in_lin, lin_copy);
        }

        if (strcmp(type, "pink") == 0) {
          nd->vv[ch].typ = 4; // Spin
        } else if (strcmp(type, "white") == 0) {
          nd->vv[ch].typ = 12; // Bspin
        } else if (strcmp(type, "brown") == 0) {
          nd->vv[ch].typ = 11;
        } else {
          error("Unknown spin type '%s' at line %d. Use: pink, white, brown",
                type, in_lin);
        }

        char *next_word = getWord();
        if (next_word) {
          error("Invalid syntax at line %d. Expected: spin <type> width "
                "<width> rate <rate> amplitude <amp>\n  %s",
                in_lin, lin_copy);
          free(next_word);
        }

        if (width < 0) {
          error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
        }
        if (rate < 0) {
          error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
        }
        if (amp < 0 || amp > 100) {
          error("Invalid spin amplitude at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
        }

        nd->vv[ch].carr = width;
        nd->vv[ch].res = rate;
        nd->vv[ch].amp = AMP_DA(amp);
        ch++;
        lines_processed++; // Increment count of processed lines
      } else if (strcmp(typ_nxt, "effect") == 0) {
        checkBackgroundInSequence(nd);

        char *effect_type = getWord();
        if (strcmp(effect_type, "pulse") == 0) {
          // Pulse: <pulse> intensity <intensity>
          char *pulse_value = getWord();
          char *intensity_str = getWord();
          char *intensity_value = getWord();

          if (!pulse_value || !intensity_str || !intensity_value ||
              strcmp(intensity_str, "intensity") != 0) {
            error("Invalid pulse syntax at line %d. Expected: pulse <pulse> "
                  "intensity <intensity>",
                  in_lin);
          }

          double pulse, intensity;
          if (sscanf(pulse_value, "%lf", &pulse) != 1) {
            error("Invalid pulse at line %d.\n  %s", in_lin, lin_copy);
          }
          if (sscanf(intensity_value, "%lf", &intensity) != 1) {
            error("Invalid intensity at line %d.\n  %s", in_lin, lin_copy);
          }

          char *next_word = getWord();
          if (next_word) {
            error("Invalid syntax at line %d. Expected: pulse <pulse> "
                  "intensity <intensity>\n  %s",
                  in_lin, lin_copy);
            free(next_word);
          }

          if (pulse < 0) {
            error("Invalid pulse at line %d.\n  %s", in_lin, lin_copy);
          }
          if (intensity < 0 || intensity > 100) {
            error("Invalid intensity at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
          }

          nd->vv[ch].typ = 7;
          nd->vv[ch].res = pulse;
          nd->vv[ch].amp = AMP_DA(intensity);
          ch++;
          lines_processed++; // Increment count of processed lines

        } else if (strcmp(effect_type, "spin") == 0) {
          // Spin: width <width> frequency <frequency> intensity <intensity>
          char *width_str = getWord();
          char *width_value = getWord();
          char *rate_str = getWord();
          char *rate_value = getWord();
          char *intensity_str = getWord();
          char *intensity_value = getWord();

          if (!width_str || !width_value || !rate_str ||
              !rate_value || !intensity_str || !intensity_value ||
              strcmp(intensity_str, "intensity") != 0 ||
              strcmp(rate_str, "rate") != 0 ||
              strcmp(width_str, "width") != 0) {
            error("Invalid spin syntax at line %d. Expected: spin width "
                  "<width> rate <rate> intensity <intensity>\n  %s",
                  in_lin, lin_copy);
          }

          double width, rate, intensity;
          if (sscanf(width_value, "%lf", &width) != 1) {
            error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
          }
          if (sscanf(rate_value, "%lf", &rate) != 1) {
            error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
          }
          if (sscanf(intensity_value, "%lf", &intensity) != 1) {
            error("Invalid intensity at line %d.\n  %s", in_lin, lin_copy);
          }

          char *next_word = getWord();
          if (next_word) {
            error("Invalid syntax at line %d. Expected: effect spin width "
                  "<width> rate <rate> intensity <intensity>\n  %s",
                  in_lin, lin_copy);
            free(next_word);
          }

          if (width < 0) {
            error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
          }
          if (rate < 0) {
            error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
          }
          if (intensity < 0 || intensity > 100) {
            error("Invalid intensity at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
          }

          nd->vv[ch].typ = 6;
          nd->vv[ch].carr = width;
          nd->vv[ch].res = rate;
          nd->vv[ch].amp = AMP_DA(intensity);
          ch++;
          lines_processed++; // Increment count of processed lines
        } else {
          error("Unknown effect type '%s' at line %d. Use: pulse, spin",
                effect_type, in_lin);
        }
      } else {
        error("Waveform not valid for '%s' at line %d. Use: tone", typ_nxt,
              in_lin);
      }
    } else if (strcmp(cmd, "background") == 0) {
      // Parse: background amplitude <amp>
      char *amp_str = getWord();
      char *amp_value = getWord();

      if (!amp_str || !amp_value || strcmp(amp_str, "amplitude") != 0) {
        error("Invalid background syntax at line %d. Expected: background "
              "amplitude <amp>\n  %s",
              in_lin, lin_copy);
      }

      double amp;
      if (sscanf(amp_value, "%lf", &amp) != 1) {
        error("Invalid background amplitude at line %d.\n  %s", in_lin, lin_copy);
      }

      char *next_word = getWord();
      if (next_word) {
        error("Invalid syntax at line %d. Expected: background amplitude "
              "<amp>\n  %s",
              in_lin, lin_copy);
        free(next_word);
      }

      if (amp < 0 || amp > 100) {
        error("Invalid background amplitude at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
      }

      nd->vv[ch].typ = 5;
      nd->vv[ch].amp = AMP_DA(amp);
      mix_flag = 1;
      ch++;
      lines_processed++; // Increment count of processed lines
    } else if (strcmp(cmd, "spin") == 0) {
      // Parse: spin <type> width <width> frequency
      char *type = getWord();
      char *width_str = getWord();
      char *width_value = getWord();
      char *rate_str = getWord();
      char *rate_value = getWord();
      char *amp_str = getWord();
      char *amp_value = getWord();

      if (!type || !width_str || !width_value || !rate_str ||
          !rate_value || !amp_str || !amp_value ||
          strcmp(amp_str, "amplitude") != 0 ||
          strcmp(rate_str, "rate") != 0 ||
          strcmp(width_str, "width") != 0) {
        error("Invalid spin syntax at line %d. Expected: spin <type> width "
              "<width> rate <rate> amplitude <amp>\n  %s",
              in_lin, lin_copy);
      }

      double width, rate, amp;
      if (sscanf(width_value, "%lf", &width) != 1) {
        error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
      }
      if (sscanf(rate_value, "%lf", &rate) != 1) {
        error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
      }
      if (sscanf(amp_value, "%lf", &amp) != 1) {
        error("Invalid spin amplitude at line %d.\n  %s", in_lin, lin_copy);
      }

      if (strcmp(type, "pink") == 0) {
        nd->vv[ch].typ = 4;
      } else if (strcmp(type, "white") == 0) {
        nd->vv[ch].typ = 12;
      } else if (strcmp(type, "brown") == 0) {
        nd->vv[ch].typ = 11;
      } else {
        error("Unknown spin type '%s' at line %d. Use: pink, white, brown",
              type, in_lin);
      }

      char *next_word = getWord();
      if (next_word) {
        error("Invalid syntax at line %d. Expected: spin <type> width "
              "<width> rate <rate> amplitude <amp>\n  %s",
              in_lin, lin_copy);
        free(next_word);
      }

      if (width < 0) {
        error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
      }

      if (rate < 0) {
        error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
      }
      if (amp < 0 || amp > 100) {
        error("Invalid spin amplitude at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
      }

      nd->vv[ch].carr = width;
      nd->vv[ch].res = rate;
      nd->vv[ch].amp = AMP_DA(amp);
      ch++;
      lines_processed++; // Increment count of processed lines
    } else if (strcmp(cmd, "effect") == 0) {
      // Parse: effect <type> <value> amplitude <amp>
      char *type = getWord();

      if (strcmp(type, "pulse") == 0) {
        // Pulse: <pulse> intensity <intensity>
        char *pulse_value = getWord();
        char *intensity_str = getWord();
        char *intensity_value = getWord();

        if (!pulse_value || !intensity_str || !intensity_value ||
            strcmp(intensity_str, "intensity") != 0) {
          error("Invalid pulse syntax at line %d. Expected: pulse <pulse> "
                "intensity <intensity>",
                in_lin);
        }

        double pulse, intensity;
        if (sscanf(pulse_value, "%lf", &pulse) != 1) {
          error("Invalid pulse at line %d.\n  %s", in_lin, lin_copy);
        }
        if (sscanf(intensity_value, "%lf", &intensity) != 1) {
          error("Invalid intensity at line %d.\n  %s", in_lin, lin_copy);
        }

        char *next_word = getWord();
        if (next_word) {
          error("Invalid syntax at line %d. Expected: effect pulse <pulse> "
                "intensity <intensity>\n  %s",
                in_lin, lin_copy);
          free(next_word);
        }

        if (pulse < 0) {
          error("Invalid pulse at line %d.\n  %s", in_lin, lin_copy);
        }

        if (intensity < 0 || intensity > 100) {
          error("Invalid intensity at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
        }

        nd->vv[ch].typ = 7;
        nd->vv[ch].res = pulse;
        nd->vv[ch].amp = AMP_DA(intensity);
        ch++;
        lines_processed++; // Increment count of processed lines
      } else if (strcmp(type, "spin") == 0) {
        // Spin: width <width> rate <rate> intensity <intensity>
        char *width_str = getWord();
        char *width_value = getWord();
        char *rate_str = getWord();
        char *rate_value = getWord();
        char *intensity_str = getWord();
        char *intensity_value = getWord();

        if (!width_str || !width_value || !rate_str || !rate_value ||
            !intensity_str || !intensity_value ||
            strcmp(intensity_str, "intensity") != 0 ||
            strcmp(rate_str, "rate") != 0 ||
            strcmp(width_str, "width") != 0) {
          error("Invalid spin syntax at line %d. Expected: spin width "
                "<width> rate <rate> intensity <intensity>\n  %s",
                in_lin, lin_copy);
        }

        double width, rate, intensity;
        if (sscanf(width_value, "%lf", &width) != 1) {
          error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
        }
        if (sscanf(rate_value, "%lf", &rate) != 1) {
          error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
        }
        if (sscanf(intensity_value, "%lf", &intensity) != 1) {
          error("Invalid intensity at line %d.\n  %s", in_lin, lin_copy);
        }

        char *next_word = getWord();
        if (next_word) {
          error("Invalid syntax at line %d. Expected: effect spin width "
                "<width> rate <rate> intensity <intensity>\n  %s",
                in_lin, lin_copy);
          free(next_word);
        }

        if (width < 0) {
          error("Invalid spin width at line %d.\n  %s", in_lin, lin_copy);
        }

        if (rate < 0) {
          error("Invalid spin rate at line %d.\n  %s", in_lin, lin_copy);
        }

        if (intensity < 0 || intensity > 100) {
          error("Invalid intensity at line %d.\nSupported range: 0 to 100.\n  %s", in_lin, lin_copy);
        }

        nd->vv[ch].typ = 6;
        nd->vv[ch].carr = width;
        nd->vv[ch].res = rate;
        nd->vv[ch].amp = AMP_DA(intensity);
        ch++;
        lines_processed++; // Increment count of processed lines
      }
    } else {
      error("Unknown command '%s' at line %d. Use: noise, tone, waveform, "
            "spin, effect, background",
            cmd, in_lin);
    }
  }
  // Normalize total amplitude
  normalizeAmplitude(nd->vv, N_CH, lin_copy, in_lin);
}

//
//	Bad time
//

void badTime(char *tim) {
  error("Invalid time \"%s\", line %d:\n  %s", tim, in_lin, lin_copy);
}

//
//	Read a time-line of either type
//

void readTimeLine() {
  char *p, *tim_p;
  int nn;
  int fo, fi;
  Period *pp;
  NameDef *nd;
  int tim, rtim = 0;

  if (!(p = getWord()))
    badSeq();
  tim_p = p;

  // Read the time represented
  tim = -1;

  while (*p) {
    if (tim != -1)
      badTime(tim_p);

    if (0 == (nn = readTime(p, &rtim)))
      badTime(tim_p);
    p += nn;

    tim = (tim == -1) ? rtim : (tim + rtim) % H24;
  }

  if (fast_tim0 < 0)
    fast_tim0 = tim; // First time
  fast_tim1 = tim;   // Last time

  if (!(p = getWord()))
    badSeq();

  fi = fo = 1;

  for (nd = nlist; nd && 0 != strcmp(p, nd->name); nd = nd->nxt)
    ;
  if (!nd)
    error("Preset \"%s\" not defined, line %d:\n  %s", p, in_lin, lin_copy);

  // Normal name-def
  pp = (Period *)Alloc(sizeof(*pp));
  pp->tim = tim;
  pp->fi = fi;
  pp->fo = fo;

  memcpy(pp->v0, nd->vv, N_CH * sizeof(Voice));
  memcpy(pp->v1, nd->vv, N_CH * sizeof(Voice));

  if (!per)
    per = pp->nxt = pp->prv = pp;
  else {
    pp->nxt = per;
    pp->prv = per->prv;
    pp->prv->nxt = pp->nxt->prv = pp;
  }

  // Automatically add a transitional period
  pp = (Period *)Alloc(sizeof(*pp));
  pp->fi = -2; // Unspecified transition
  pp->nxt = per;
  pp->prv = per->prv;
  pp->prv->nxt = pp->nxt->prv = pp;

  if (0 != (p = getWord()))
    badSeq();

  pp->fi = -3;
  pp->tim = tim;
}

int readTime(char *p, int *timp) { // Rets chars consumed, or 0 error
  int nn, hh, mm, ss;

  if (3 > sscanf(p, "%2d:%2d:%2d%n", &hh, &mm, &ss, &nn)) {
    ss = 0;
  }

  if (hh < 0 || hh >= 24 || mm < 0 || mm >= 60 || ss < 0 || ss >= 60)
    return 0;

  *timp = ((hh * 60 + mm) * 60 + ss) * 1000;
  return nn;
}

//
// Function to normalize the total amplitude of voices
//
void
normalizeAmplitude(Voice *voices, int numChannels, const char *line,
                        int lineNum) {
  double totalAmplitude = 0.0;

  // Calculate the total amplitude of all voices (excluding effect spin/pulse)
  for (int ch = 0; ch < numChannels; ch++) {
    if (voices[ch].typ != 0 && voices[ch].typ != 6 && voices[ch].typ != 7) {
      double ampPercentage = voices[ch].amp / 40.96;
      totalAmplitude += ampPercentage;
    }
  }

  // If total amplitude exceeds 100%, normalize all active voices
  if (totalAmplitude > 100.0) {
    double normalizationFactor = 100.0 / totalAmplitude;
    // Apply normalization to all active voices (except effect spin/pulse)
    for (int ch = 0; ch < numChannels; ch++) {
      if (voices[ch].typ != 0 && voices[ch].typ != 6 && voices[ch].typ != 7) {
        voices[ch].amp *= normalizationFactor;
      }
    }
  }
}

void checkBackgroundInSequence(NameDef *nd) {
  int mix_exists = 0;

  for (int ch = 0; ch < N_CH; ch++) {
    if (nd->vv[ch].typ == 5) { // background
      mix_exists = 1;
      break;
    }
  }

  if (!mix_exists)
    error("effect spin/pulse without file amplitude specified, line %d:\n  %s",
          in_lin, lin_copy);
}

// END //
