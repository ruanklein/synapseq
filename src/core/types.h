// core/types.h
#ifndef SYNAPSEQ_TYPES_H
#define SYNAPSEQ_TYPES_H

#include "config.h"

// Forward declarations
typedef struct Channel Channel;
typedef struct Voice Voice; 
typedef struct Period Period;
typedef struct NameDef NameDef;
typedef struct Noise Noise;

typedef long long S64;
typedef unsigned char uchar;

// Estruturas completas
struct Voice {
  int typ;      // Voice type: 0 off, 1 binaural, 2 pink noise, etc.
  double amp;   // Amplitude level (0-4096 for 0-100%)
  double carr;  // Carrier freq (for binaural/monaural/isochronic), width (for spin)
  double res;   // Resonance freq (-ve or +ve) (for binaural/spin/isochronic)
  int waveform; // 0=sine, 1=square, 2=triangle, 3=sawtooth
};

struct Channel {
  Voice v;        // Current voice setting (updated from current period)
  int typ;        // Voice type
  int amp, amp2;  // Current state, according to current type
  int inc1, off1; // For binaural tones, offset + increment into sine table * 65536
  int inc2, off2; // For binaural tones, offset + increment into sine table * 65536
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
  Voice vv[N_CH]; // Voice-set for this definition
};

struct Noise {
  int val; // Current output value
  int inc; // Increment
};

#endif // SYNAPSEQ_TYPES_H