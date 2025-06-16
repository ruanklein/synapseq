//
//  SynapSeq - Synapse-Sequenced Brainwave Generator
//
//  (c) 2025 Ruan <https://ruan.sh/>
//
//  Originally based on SBaGen (Sequenced Binaural Beat Generator) v1.4.4 by Jim Peters.
//  SBaGen homepage: http://uazu.net/sbagen/
//  License: GNU General Public License v2 (GPLv2)
//
//  SynapSeq is released under the GNU GPL version 2.
//  Use at your own risk.
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
//
//	OGG decoding using libvorbis
//
//        (c) 1999-2004 Jim Peters <jim@uazu.net>.  All Rights Reserved.
//        For latest version see http://sbagen.sf.net/ or
//        http://uazu.net/sbagen/.  Released under the GNU GPL version 2.
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <errno.h>
#include <time.h>
#include <math.h>
#include <ctype.h>
#include <sys/time.h>

#include <vorbis/codec.h>
#include <vorbis/vorbisfile.h>

extern FILE *mix_in;
extern void *Alloc(size_t);
extern void error(char *fmt, ...);
extern void warn(char *fmt, ...);
extern void debug(char *fmt, ...);
extern int out_rate, out_rate_def;
extern void inbuf_start(int (*rout)(int *, int), int len);

#define ALLOC_ARR(cnt, type) ((type *)Alloc((cnt) * sizeof(type)))

// Function declarations
void ogg_init();
void ogg_term();
int ogg_read(int *dst, int dlen);

// Static variables for OGG decoding
static OggVorbis_File oggfile;
static short *ogg_buf0, *ogg_buf1, *ogg_rd, *ogg_end;
static int ogg_mult;

void ogg_init() {
   vorbis_info *vi;
   vorbis_comment *vc;
   int len = 2048;
   int a;

   // Setup OGG decoder
   if (0 > ov_open(mix_in, &oggfile, NULL, 0)) 
      error("Input does not appear to be an Ogg bitstream");

   // Check for ReplayGain
   vc = ov_comment(&oggfile, -1);
   ogg_mult = 16;
   for (a = 0; a < vc->comments; a++) {
      char *str = vc->user_comments[a];
      if (0 == memcmp(str, "REPLAYGAIN_TRACK_GAIN=", 22)) {
         char *end;
         double val = strtod(str += 22, &end);
         if (end == str) 
            warn("Ignoring bad REPLAYGAIN_TRACK_GAIN: %s", str);
         else {
            val -= 3;		// Adjust vorbisgain's 89dB to 86dB
            ogg_mult = (int)(floor(0.5 + 16 * pow(10, val/20)));
            //warn("ReplayGain setting detected, Ogg scaling by %.2f", ogg_mult/16.0);
         }
      }
   }

   // Pick up sampling rate and override default if -r not used
   vi = ov_info(&oggfile, -1);
   if (out_rate_def) out_rate = vi->rate;
   out_rate_def = 0;

   // Setup buffer so that we can be consistent in our calls to ov_read
   ogg_buf0 = ALLOC_ARR(len, short);
   ogg_buf1 = ogg_buf0 + len;
   ogg_rd = ogg_end = ogg_buf0;

   // Start a thread to handle generation of the mix stream
   inbuf_start(ogg_read, 256*1024);	// 1024K buffer: 3s@44.1kHz
}

void ogg_term() {
   ov_clear(&oggfile);
   if (ogg_buf0) { 
      free(ogg_buf0); 
      ogg_buf0 = NULL; 
   }
}

int ogg_read(int *dst, int dlen) {
   int *dst0 = dst;
   int *dst1 = dst + dlen;

   while (dst < dst1) {
      int rv, sect;

      // Copy data from buffer
      if (ogg_rd != ogg_end) {
         while (ogg_rd != ogg_end && dst != dst1)
            *dst++ = *ogg_rd++ * ogg_mult;
         continue;
      }

      int samples_available = ogg_buf1 - ogg_buf0;
      int bytes_to_read = samples_available * sizeof(short);
      
      rv = ov_read(&oggfile, (char*)ogg_buf0, bytes_to_read, 0, 2, 1, &sect);
      
      if (rv < 0) {
         if (rv == OV_HOLE) {
            warn("Recoverable gap in Ogg stream");
            continue;
         } else if (rv == OV_EREAD) {
            warn("Read error in Ogg stream - possible file corruption");
            continue;
         } else {
            warn("Recoverable error in Ogg stream (code: %d)", rv);
            continue;
         }
      }
      
      if (rv == 0) // EOF
         return dst-dst0;
         
      if (rv % sizeof(short) != 0) {
         warn("UNEXPECTED: ov_read() returned unaligned byte count: %d", rv);
         rv = (rv / sizeof(short)) * sizeof(short);
      }
      
      ogg_rd = ogg_buf0;
      ogg_end = ogg_buf0 + (rv / sizeof(short));
   }
   return dst-dst0;
}

// END //