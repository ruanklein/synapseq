// core/constants.h
#ifndef SYNAPSEQ_CONSTANTS_H
#define SYNAPSEQ_CONSTANTS_H

#include "config.h"

#ifdef T_WIN32
#define WIN_TIME // Windows time functions
#define WIN_MISC // Windows miscellaneous functions
#define vsnprintf _vsnprintf
#include <windows.h>
#endif

#ifdef T_POSIX
#define UNIX_TIME // UNIX time functions
#define UNIX_MISC // UNIX miscellaneous functions
#endif

#ifdef UNIX_TIME
#include <sys/ioctl.h>
#include <sys/times.h>
#endif
#ifdef UNIX_MISC
#include <pthread.h>
#endif

#define ALLOC_ARR(cnt, type) ((type *)Alloc((cnt) * sizeof(type))) // Allocate an array
#define uint unsigned int // 32-bit unsigned integer

#define ST_AMP 0x7FFFF // Amplitude of wave in sine-table
#define NS_ADJ 12 // Noise is generated internally with amplitude ST_AMP<<NS_ADJ
#define NS_DITHER 16 // How many bits right to shift the noise for dithering
#define NS_AMP (ST_AMP << NS_ADJ)
#define ST_SIZ 16384 // Number of elements in sine-table (power of 2)

#define AMP_DA(pc) (40.96 * (pc))   // Display value (%age) to ->amp value
#define AMP_AD(amp) ((amp) / 40.96) // Amplitude value to display %age

#define NS_BIT 10 // Number of bits in noise-table

#define H24 (86400000) // 24 hours
#define H12 (43200000) // 12 hours

#define RAND_MULT 75 // Random number generator multiplier

#define NS_BANDS 9 // Number of bands in noise-table

// Add a 4-byte unsigned integer to the buffer
#define addU4(xx, p)                                                             \
{                                                                            \
  int a = xx;                                                                \
  *p++ = a;                                                                  \
  *p++ = (a >>= 8);                                                          \
  *p++ = (a >>= 8);                                                          \
  *p++ = (a >>= 8);                                                          \
}
// Add a string to the buffer
#define addStr(xx, p)                                                             \
{                                                                            \
  char *q = xx;                                                              \
  *p++ = *q++;                                                               \
  *p++ = *q++;                                                               \
  *p++ = *q++;                                                               \
  *p++ = *q++;                                                               \
}

#define NL "\n" // New line

#endif