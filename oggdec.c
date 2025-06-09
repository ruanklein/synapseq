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
//	OGG decoding using Tremor (libvorbisidec), with optional
//	looping.
//
//        (c) 1999-2004 Jim Peters <jim@uazu.net>.  All Rights Reserved.
//        For latest version see http://sbagen.sf.net/ or
//        http://uazu.net/sbagen/.  Released under the GNU GPL version 2.
//
//	See comments lower down for OGG looper parameters.
//

// This file is a fork of the original oggdec.c file from SBaGen.
// SBAGEN_LOOPER is replaced with SYNAPSEQ_LOOPER.

#include "libs/ivorbiscodec.h"
#include "libs/ivorbisfile.h"

extern FILE *mix_in;
extern int mix_cnt;
extern void *Alloc(size_t);
extern void error(char *fmt, ...);
extern int out_rate, out_rate_def;

void ogg_init() ;
void ogg_term() ;
int ogg_read(int *dst, int dlen) ;
void looper_init() ;
void looper_term() ;
int looper_read(int *dst, int dlen) ;
static void looper_sched() ;
static void looper_sched2() ;

static OggVorbis_File oggfile;
static short *ogg_buf0, *ogg_buf1, *ogg_rd, *ogg_end;
static int ogg_mult;

void 
ogg_init() {
   vorbis_info *vi;
   vorbis_comment *vc;
   int len= 2048;
   int a;

   // Setup OGG decoder
   if (0 > ov_open(mix_in, &oggfile, NULL, 0)) 
      error("Input does not appear to be an Ogg bitstream");

   // Check for ReplayGain
   vc= ov_comment(&oggfile, -1);
   ogg_mult= 16;
   for (a= 0; a<vc->comments; a++) {
      char *str= vc->user_comments[a];
      if (0 == memcmp(str, "REPLAYGAIN_TRACK_GAIN=", 22)) {
	 char *end;
	 double val= strtod(str += 22, &end);
	 if (end == str) 
	    warn("Ignoring bad REPLAYGAIN_TRACK_GAIN: %s", str);
	 else {
	    val -= 3;		// Adjust vorbisgain's 89dB to 86dB
	    ogg_mult= (int)(floor(0.5 + 16 * pow(10, val/20)));
	    warn("ReplayGain setting detected, Ogg scaling by %.2f", ogg_mult/16.0);
	 }
      }
   }

   // Check to see is this is a looping OGG
   for (a= 0; a<vc->comments; a++) {
      if (0 == memcmp(vc->user_comments[a], "SYNAPSEQ_LOOPER=", 14)) {
	 vc= 0; break;
      }
   }
   if (!vc) {
      // This is a looping OGG.  Load the OGG into memory instead;
      // handle with looper_* code.  ov_clear is called after because
      // it closes mix_in.
      looper_init();
      ov_clear(&oggfile);
      return;
   }

   // Pick up sampling rate and override default if -r not used
   vi= ov_info(&oggfile, -1);
   if (out_rate_def) out_rate= vi->rate;
   out_rate_def= 0;

   // Setup buffer so that we can be consistent in our calls to ov_read
   ogg_buf0= ALLOC_ARR(len, short);
   ogg_buf1= ogg_buf0 + len;
   ogg_rd= ogg_end= ogg_buf0;

   // Start a thread to handle generation of the mix stream
   inbuf_start(ogg_read, 256*1024);	// 1024K buffer: 3s@44.1kHz
}

void 
ogg_term() {
   ov_clear(&oggfile);
   if (ogg_buf0) { free(ogg_buf0); ogg_buf0= 0; }
   looper_term();
}

int 
ogg_read(int *dst, int dlen) {
   int *dst0= dst;
   int *dst1= dst + dlen;

   while (dst < dst1) {
      int rv, sect;

      // Copy data from buffer
      if (ogg_rd != ogg_end) {
	 while (ogg_rd != ogg_end && dst != dst1)
	    *dst++= *ogg_rd++ * ogg_mult;
	 continue;
      }

      // Refill buffer
      rv= ov_read(&oggfile, (char*)ogg_buf0, (ogg_buf1-ogg_buf0)*sizeof(short), &sect);
      //debug("ov_read %d/%d", rv, (ogg_buf1-ogg_buf0)*sizeof(short));
      if (rv < 0) {
	 warn("Recoverable error in Ogg stream  ");
	 continue;
      }
      if (rv == 0) 	// EOF
	 return dst-dst0;
      if (rv & 3)
	 error("UNEXPECTED: ov_read() returned a partial sample count: %d", rv);
      ogg_rd= ogg_buf0;
      ogg_end= ogg_buf0 + (rv/2);
   }
   return dst-dst0;
}

//
//	Looping OGG decoding using Tremor (libvorbisidec)
//
//        (c) 1999-2004 Jim Peters <jim@uazu.net>.  All Rights Reserved.
//        For latest version see http://sbagen.sf.net/ or
//        http://uazu.net/sbagen/.  Released under the GNU GPL version 2.
//
//	This code preloads the OGG file into memory, and loops and
//	cross-fades the audio in segments to provide an endless output
//	stream.
//
//	The type of cross-fading and looping for any particular OGG
//	file is controlled by a spec-string embedded in the file with
//	the tag "SYNAPSEQ_LOOPER".  This may be set using the standard
//	OGG 'vorbiscomment' tool, for example:
//
//	  vorbiscomment -a -t "SYNAPSEQ_LOOPER=s4-16 f0.2" in.ogg out.ogg
//
//	The tag contains entries as follows, optionally separated by
//	white space:
//
//	  s<dur>	Set segment size to given duration in seconds
//	  s<min>-<max>	Set segment size to randomly vary within given range
//	  f<dur>	Set duration for cross-fades
//	  c<cnt>	Number of channels: 1 or 2
//	  w<bool>	Swap stereo on second channel? 0 no, or 1 yes
//	  d<min>-<max>	Set part of OGG file data to use as source audio (in seconds)
//	  #<digits>	Following settings apply only to section <digits>
//	  
//	The default is something like "s99999999f1c1w1d0-99999999",
//	i.e. segments the full length of the audio, taken out of the
//	whole of the audio file, repeated forever with a 1 second
//	cross-fade at start/end.
//
//	The # part works as follows.  The .ogg filename on the command
//	line may be followed by #<digits>.  If this is missing, it is
//	equivalent to #0.  This allows different groups of settings to
//	be selected out of the SYNAPSEQ_LOOPER string.  As an example
//	(spaced to make it more readable):
//
//	  SYNAPSEQ_LOOPER=s4-16  #0 f0.2  #1 f0.5  #2 f1
//
//	The initial part is read in all cases, and then the #0, #1 and
//	#2 select different cross-fade times.
//

static char *data;	// OGG data
static char *data_end;	// OGG data end +1
static int datlen;	// Length of OGG data in bytes (data_end-data)
static int datcnt;	// Length of OGG data in samples
static int datcnt0;	// Original datcnt value, before 'd' modifications
static int datbase;	// Base sample offset for seeking
static int datrate;	// Sampling rate of file, even if the output rate has been overridden

static int fade_cnt;	// Count for fade-in or fade-out in samples
static int seg0, seg1;	// Segment size range in samples (min/max) including fade in/out
static int ch2;		// Channel 2 mode: 0 off, 1 on
static int ch2_swap;	// Channel 2 swapped-stereo?  0 off, 1 on
static uint del_amp;	// Fade in/out delta in amplitude per sample

typedef struct {
   OggVorbis_File ogg;	// File
   char *rd;		// Current read position in 'data'
   int off;		// Sample-offset that we started playing at for this segment
   int chan;		// Which channel do we below to?  0 or 1; if 1 then ch2_swap applies
   int mode;		// Current mode: 0 inactive, 1 waiting, 2 fade-in, 3 hold, 4 fade-out
   int cnt;		// Current count-down
   int cnt_all;		// Total sample count for fade-in/hold/fade-out 
   uint amp;		// Current amplitude 0-0xFFFFFFFF
   uint del;		// Current amplitude delta (i.e. increment)
   short *buf0, *buf1;	// Buffer for samples, buf0 to buf1-1
   short *b_rd, *b_end;	// Buffer read position, end pos +1
} AStream;
AStream str[3];	// Three read streams from the data

// Simple/crude 16-bit pseudo-random number generator, from ZX Spectrum
static unsigned short zxrand_seed;	// Current random seed
static int 
zxrand_0_65536() {
   zxrand_seed= (1 + (int)zxrand_seed) * 75 % 65537 - 1;
   return zxrand_seed; 
}
static int 	// Returns crude pseudo-random value from 0 to mult-1
zxrand_0(int mult) {
   long long tmp= mult;
   tmp *= zxrand_0_65536();
   tmp >>= 16;
   return (int)tmp;
}
// Returns crude pseudo-random value from r0 to r1-1, or r0 if the range is invalid
static int
zxrand(int r0, int r1) {
   if (r1 <= r0) return r0;
   return r0 + zxrand_0(r1-r0);
}

// Return random value out of multiple ranges.  Each range is
// represented by two values (lo, hi), and it represents a range
// lo->(hi-1).  It is fine if any range is invalid (i.e. 0 or negative
// length), as it will be ignored.  However, if all ranges are
// invalid, then the default value 'def' is returned instead.
// 'fmt' is a string of format characters, indicating the arguments 
// to follow:
//	'r' indicates a range (int lo, int hi), from lo->(hi-1)
//	'o' indicates a range (int lo, int hi) that will be used to
//	      limit all the following 'r' ranges
static int
zxrandM(int def, char *fmt, ...) {
   va_list ap;
   int cnt= 0;
   int val;
   char *p;
   int olo, ohi;	// Overall limiting range

   va_start(ap, fmt); 
   olo= 0x80000000; ohi= 0x7FFFFFFF;
   for (p= fmt; *p; p++) {
      int lo= va_arg(ap, int);
      int hi= va_arg(ap, int);
      if (*p == 'o') { olo= lo; ohi= hi; }
      else if (*p == 'r') {
	 if (lo < olo) lo= olo;
	 if (hi > ohi) hi= ohi;
	 if (hi-lo > 0) cnt += hi-lo;
      } else 
	 error("Bad zxrandM format: %s", fmt);
   }
   va_end(ap);

   if (!cnt) return def;
   val= zxrand_0(cnt);
   
   va_start(ap, fmt);
   olo= 0x80000000; ohi= 0x7FFFFFFF;
   for (p= fmt; *p; p++) {
      int lo= va_arg(ap, int);
      int hi= va_arg(ap, int);
      if (*p == 'o') { olo= lo; ohi= hi; }
      else if (*p == 'r') {
	 if (lo < olo) lo= olo;
	 if (hi > ohi) hi= ohi;
	 cnt= hi-lo;
	 if (cnt > 0) {
	    if (val < cnt) return lo + val;
	    val -= cnt;
	 }
      }
   }
   va_end(ap);
   return def;
}

// File buffer operations
size_t
oc_read(void *ptr, size_t size, size_t nmemb, void *datasource) {
   AStream *ss= (AStream*)datasource;
   int len= size * nmemb;
   if (len > data_end - ss->rd) 
      len= data_end - ss->rd;
   if (len) {
      memcpy(ptr, ss->rd, len);
      ss->rd += len;
   }
   return len;
}

int
oc_seek(void *datasource, ogg_int64_t offset, int whence) {
   AStream *ss= (AStream*)datasource;
   int pos= ss->rd - data;
   switch (whence) {
    case SEEK_SET:
       pos= offset; break;
    case SEEK_CUR:
       pos += offset; break;
    case SEEK_END:
       pos= datlen += offset; break;
   }
   if (pos < 0) pos= 0;
   if (pos > datlen) pos= datlen;
   ss->rd= data + pos;
   return pos;
}

int
oc_close(void *datasource) {
   return 0;
}

long 
oc_tell(void *datasource) {
   AStream *ss= (AStream*)datasource;
   return ss->rd - data;
}

// Initialise all
void 
looper_init() {
   vorbis_info *vi;
   int a;
   ov_callbacks oc;
   vorbis_comment *vc;
   char *looper;
   int prev_flag= 0;
   int on;

   // Init random seed
   zxrand_seed= 0xFFFF & time(NULL);

   // Find the length of the OGG file and load it into memory
   if (0 != fseek(mix_in, 0, SEEK_END) ||
       0 > (datlen= ftell(mix_in)))
      error("Can't determine length of OGG file: %s", strerror(errno));
   
   data= ALLOC_ARR(datlen, char);
   data_end= data + datlen;
   if (0 != fseek(mix_in, 0, SEEK_SET) ||
       1 != fread(data, datlen, 1, mix_in))
      error("Can't read loopable OGG file into memory: %s", strerror(errno));

   // Open each of the three OGG streamers
   oc.read_func= oc_read;
   oc.seek_func= oc_seek;
   oc.close_func= oc_close;
   oc.tell_func= oc_tell;
   for (a= 0; a<3; a++) {
      AStream *aa= &str[a];
      aa->rd= data;
      aa->mode= 0;
      aa->buf0= ALLOC_ARR(2048, short);
      aa->buf1= aa->buf0 + 2048;
      aa->b_rd= aa->b_end= aa->buf0;
      if (0 > ov_open_callbacks(aa, &aa->ogg, NULL, 0, oc)) 
	 error("Problem opening OGG bitstream");
   }

   // Find the total length of the file
   datcnt= ov_pcm_total(&str[0].ogg, -1);
   datcnt0= datcnt;
   datbase= 0;
   
   // Pick up sampling rate and override default if -r not used
   vi= ov_info(&str[0].ogg, -1);
   datrate= vi->rate;
   if (out_rate_def) out_rate= vi->rate;
   out_rate_def= 0;

   // Find the SYNAPSEQ_LOOPER tag
   looper= "";
   vc= ov_comment(&str[0].ogg, -1);
   for (a= 0; a<vc->comments; a++) 
      if (0 == memcmp(vc->user_comments[a], "SYNAPSEQ_LOOPER=", 14)) 
	 looper= vc->user_comments[a] + 14;
   
   // Setup info from SYNAPSEQ_LOOPER tag
   seg0= seg1= datcnt;
   fade_cnt= datrate;
   ch2= 0;
   ch2_swap= 1;
   on= 1;
   if (mix_cnt < 0) mix_cnt= 0;
   while (*looper) {
      char flag, *p;
      double val;

      flag= *looper++;
      if (isspace(flag)) continue;
      if (!strchr("s-fcwd#", flag)) {
	 warn("Bad SYNAPSEQ_LOOPER flag: %c", flag);
	 continue;
      }
      if (flag == '-') switch (prev_flag) {
       case 's': flag= 'S'; break;
       case 'd': flag= 'D'; break;
       default:
	 warn("SYNAPSEQ_LOOPER '-' found not in form s<val>-<val>");
	 continue;
      }
      prev_flag= flag;

      val= strtod(looper, &p);
      if (p == looper) {
	 warn("Bad SYNAPSEQ_LOOPER value for flag '%c': %s", flag, p);
	 continue;
      }
      looper= p;

      if (flag == '#')
	 on= (val == mix_cnt);
      else if (on) switch (flag) {
       case 's': seg0= seg1= val * datrate; break;
       case 'S': seg1= val * datrate; break;
       case 'd': datbase= val * datrate; datcnt= datcnt0 - datbase; break;	
       case 'D': datcnt= val * datrate - datbase; break;
       case 'f': fade_cnt= val * datrate; break;
       case 'c': ch2= val > 1.5; break;
       case 'w': ch2_swap= val > 0.5; break;
      }
   }

   // Tidy up LOOPER settings
   if (fade_cnt < datrate/50) fade_cnt= datrate/50;	// 20ms min fade
   if (datcnt + datbase > datcnt0) datcnt= datcnt0-datbase;
   if (datcnt < 0) 
      error("Source data range invalid in SYNAPSEQ_LOOPER settings");
   if (datcnt <= 3 * fade_cnt) 
      error("Length of source data 'd' too short for fade-length of %gs\n"
	    " in SYNAPSEQ_LOOPER settings", fade_cnt * 1.0 / datrate);
   if (seg0 > datcnt) seg0= datcnt;
   if (seg1 > datcnt) seg1= datcnt;
   if (seg0 > seg1) seg0= seg1;
   if (seg0 < 3 * fade_cnt) { 
      seg0= 3 * fade_cnt;
      warn("SYNAPSEQ_LOOPER segment size too short for fade-length of %gs; adjusted.",
	   fade_cnt * 1.0 / datrate);
   }
   if (seg1 < seg0) seg1= seg0;

   // Calculate delta to use for fades
   del_amp= 0xFFFFFFFFU/fade_cnt;		// Rely on rounding down here
   if (del_amp * (uint)fade_cnt < 0xF0000000)	// Paranoid check
      error("Internal rounding error in calculating amplitude delta");
   if (ch2) del_amp >>= 1;

   //   debug("Segment range %d-%d, fade %d, amp delta per sample %d", 
   //	      seg0, seg1, fade_cnt, del_amp);

   // Init three streams and start off
   looper_sched();

   // Start a thread to handle generation of the mix stream
   inbuf_start(looper_read, 256*1024);	// 1024K buffer: 3s@44.1kHz
}

void 
looper_term() {
   int a;
   for (a= 0; a<3; a++) {
      ov_clear(&str[a].ogg);
      free(str[a].buf0);
   }
}

int 
looper_read(int *dst, int dlen) {
   int *dst0= dst;
   int *dst1= dst + dlen;

   // Clear the whole thing
   memset(dst0, 0, (char*)dst1 - (char*)dst0);

   // Go through in chunks
   while (dst0 < dst1) {
      int a;
      int len= (dst1-dst0)/2;
      int resched= 0;

      for (a= 0; a<3; a++) {
	 AStream *aa= &str[a];
	 if (aa->mode && aa->cnt < len)
	    len= aa->cnt;
      }

      // Process 'len' samples in each channel
      for (a= 0; a<3; a++) {
	 AStream *aa= &str[a];
	 int cnt= len;
	 
	 if (!aa->mode) continue;

	 dst= dst0;
	 while (cnt > 0) {
	    // Refill buffer if necessary
	    if (aa->mode > 1 && 
		aa->b_rd == aa->b_end) {
	       int sect;
	       char *buf= (char*)aa->buf0;	
	       int len= (aa->buf1-aa->buf0)*sizeof(aa->buf0[0]);
	       int rv= ov_read(&aa->ogg, buf, len, &sect);
	       if (rv < 0) {
		  warn("Recoverable error in Ogg stream  ");
		  continue;
	       }
	       if (rv == 0) {    // EOF; internal error?
		  warn("Hit EOF in looping OGG stream, filling with zeros");
		  memset(buf, 0, len);
		  rv= len;
	       }
	       if (rv & 3)
		  error("UNEXPECTED: ov_read() returned a partial sample count: %d", rv);
	       aa->b_rd= aa->buf0;
	       aa->b_end= aa->buf0 + (rv/2);
	    }
	    
	    // Waiting to start playing
	    if (aa->mode == 1) {
	       aa->cnt -= cnt;
	       cnt= 0;
	       continue;
	    }
	    
	    // Playing fade-in / main part / fade-out
	    if (aa->chan && ch2_swap) {
	       // Output with swapped L+R
	       while (cnt > 0 && aa->b_rd != aa->b_end) {
		  uint amp= (~aa->amp) >> 16; amp= (~(amp*amp))>>21; amp *= ogg_mult;
		  cnt--;
		  *dst++ += ((int)(aa->b_rd[1] * amp)) >> 11;
		  *dst++ += ((int)(aa->b_rd[0] * amp)) >> 11;
		  aa->b_rd += 2;
		  aa->amp += aa->del;
		  aa->cnt--;
	       }
	    } else {
	       // Output with normal L+R 
	       while (cnt > 0 && aa->b_rd != aa->b_end) {
		  uint amp= (~aa->amp) >> 16; amp= (~(amp*amp))>>21; amp *= ogg_mult;
		  cnt--;
		  *dst++ += ((int)(*aa->b_rd++ * amp)) >> 11;
		  *dst++ += ((int)(*aa->b_rd++ * amp)) >> 11;
		  aa->amp += aa->del;
		  aa->cnt--;
	       }
	    }
	 }

	 if (!aa->cnt) switch (aa->mode) {
	  case 1:
	     aa->mode= 2;
	     aa->cnt= fade_cnt;
	     aa->del= del_amp;
	     break;
	  case 2:
	     aa->mode= 3;
	     aa->cnt= aa->cnt_all - 2 * fade_cnt;
	     aa->del= 0;
	     break;
	  case 3:
	     aa->mode= 4;
	     aa->cnt= fade_cnt;
	     aa->del= -del_amp;
	     break;
	  case 4:
	     aa->mode= 0;
	     resched= 1;
	     break;
	 }
      }

      dst0 += len*2;
      
      // Schedule in another section if necessary
      if (resched) looper_sched();
   }

   return dlen;		// Always returns a full buffer
}

// Calculate sample-count before end of playback on this channel
#define CNT_TO_END(aa) \
(aa->mode == 1 ? aa->cnt_all + aa->cnt : \
 aa->mode == 2 ? aa->cnt_all - fade_cnt + aa->cnt : \
 aa->mode == 3 ? fade_cnt + aa->cnt : \
 aa->mode == 4 ? aa->cnt : 0)

//
//	Schedule in the next segment or segments
//

static void 
looper_sched() {
   // int update= 0;

   // Handle 2-channel separately
   if (ch2) {
      looper_sched2();
      return;
   }
   
   while (1) {
      AStream *aa= &str[0];
      AStream *bb= &str[1];
      int cnt_all, rv;

      // Complete when both are set up ready
      if (aa->mode && bb->mode) break;
      // update= 1;

      // Make 'bb' the one we are filling in
      if (bb->mode && !aa->mode) {
	 AStream *tmp= aa; aa= bb; bb= tmp;
      }

      // Make 'aa' zero if it is not active
      if (!aa->mode) aa= 0;

      // Setup a new segment for 'bb'
      bb->off= -1;
      bb->mode= 1;
      bb->cnt= 0;
      bb->amp= 0;
      bb->b_rd= bb->b_end= bb->buf0;
      
      // Position segment to follow on from segment currently playing
      // on 'aa'
      if (aa) {
	 bb->cnt= CNT_TO_END(aa) - fade_cnt;
	 if (bb->cnt < 0) bb->cnt= 0;
      }

      // Select the segment length
      bb->cnt_all= cnt_all= zxrand(seg0, seg1+1);
      
      // Look for a segment starting-point which doesn't overlap the
      // segment playing on the other channel.  First range below is
      // before 'aa' segment, second is after 'aa' segment.
      if (aa) {
	 bb->off= zxrandM(-1, "rr", 	// Default -1 if no range is available
			  0, aa->off - cnt_all,
			  aa->off + aa->cnt_all, datcnt - cnt_all);
      }

      // If we haven't been able to work around 'aa', or 'aa' is not
      // active, then just randomly choose a position
      if (bb->off < 0)
	 bb->off= zxrand(0, datcnt - cnt_all);

      // Reseek the OGG stream
      rv= ov_pcm_seek(&bb->ogg, datbase + bb->off);
      if (rv) error("UNEXPECTED: Can't reseek the looping OGG stream: %d", rv);
   }

   //   //DEBUG
   //   if (update) {
   //      int a;
   //      for (a= 0; a<2; a++) {
   //	 AStream *aa= &str[a];
   //	 printf("STR%d ", a);
   //	 if (!aa->mode) { printf("off\n"); continue; }
   //	 printf("mode %d, cnt %d, cnt_all %d\n",
   //		aa->mode, aa->cnt, aa->cnt_all);
   //      }
   //   }
}

// This is for two tracks playing constantly
static void 
looper_sched2() {
   while (1) {
      AStream *aa= &str[0];
      AStream *bb= &str[1];
      AStream *cc= &str[2];
      AStream *tmp;
      int cnt_all, rv;
      
      // Complete when all are ready
      if (aa->mode && bb->mode && cc->mode) break;
      
      // Put the used slots in order: aa first, then bb, then cc
      if (!aa->mode && bb->mode) {
 	 tmp= aa; aa= bb; bb= tmp;
      }
      if (!aa->mode && cc->mode) {
 	 tmp= aa; aa= cc; cc= tmp;
      }
      if (!bb->mode && cc->mode) {
 	 tmp= bb; bb= cc; cc= tmp;
      }
      
      // If there are only 0 or 1 in use, we can add a new one on the
      // other channel, so long as its fade-out time doesn't clash
      // with 'aa', if active.
      if (!bb->mode) {
	 // Setup a new segment for 'bb'
 	 bb->chan= aa ? !aa->chan : 0;
	 bb->off= -1;
	 bb->mode= 1;
	 bb->amp= 0;
	 bb->cnt= 0;
	 bb->b_rd= bb->b_end= bb->buf0;

	 // Select the segment length and start-point, avoiding a
	 // length that would cause a clash with 'aa's fade-out
	 if (!aa) {
	    cnt_all= zxrand(seg0, seg1+1);
	 } else {
	    int end= CNT_TO_END(aa);		// Count before aa fade-out end
	    cnt_all= zxrandM(-1, "orr",
			     seg0, seg1+1,	    // Overall limit to valid range seg0->seg1
			     seg0, end-fade_cnt,    // To fit before the 'aa' fade-out
			     end+fade_cnt, seg1+1); // To fit after

	    // If -1 then we can't fit it in with the provided
	    // segment-length range, so start late instead
	    if (cnt_all < 0) {
	       bb->cnt= end+fade_cnt - seg1;
	       cnt_all= seg1;
	    }
	 } 
	 bb->cnt_all= cnt_all;
	 if (bb->cnt < 0) bb->cnt= 0;

	 // Look for a segment starting-point which doesn't overlap
	 // the segment playing on the other channel, or if that
	 // didn't work, just choose a random position
	 if (aa) {
	    bb->off= zxrandM(-1, "rr", 	// Default -1 if no range is available
			     0, aa->off - cnt_all,
			     aa->off + aa->cnt_all, datcnt - cnt_all);
	 }
	 if (bb->off < 0)
	    bb->off= zxrand(0, datcnt - cnt_all);

	 // Reseek the OGG stream
	 rv= ov_pcm_seek(&bb->ogg, datbase + bb->off);
	 if (rv) error("UNEXPECTED: Can't reseek the looping OGG stream: %d", rv);
	 continue;
      }
	 
      // Okay, so we have 'aa' and 'bb' and we need to fill in 'cc'.
      // 'aa' and 'bb' should be on different channels.  We need to
      // create a 'cc' that continues where the shorter of 'aa' and
      // 'bb' leaves off.
      if (aa->chan == bb->chan)
	 error("UNEXPECTED: internal error, looper_sched2(), aa/bb on same chan");

      // Make 'aa' the shorter one
      if (CNT_TO_END(aa) > CNT_TO_END(bb)) {
 	 tmp= aa; aa= bb; bb= tmp;
      }

      // Setup a new segment for 'cc'
      cc->chan= aa->chan;
      cc->cnt= CNT_TO_END(aa) - fade_cnt;
      cc->off= -1;
      cc->mode= 1;
      cc->amp= 0;
      cc->b_rd= cc->b_end= cc->buf0;
      
      // Select the segment length and start-point, avoiding a
      // length that would cause a clash with 'bb's fade-out
      if (!bb) {
	 cnt_all= zxrand(seg0, seg1+1);
      } else {
	 int end= CNT_TO_END(bb);		// Count before bb fade-out end
	 end -= cc->cnt;			// Take account of the offset we are starting at
	 cnt_all= zxrandM(-1, "orr",
			  seg0, seg1+1,	    // Overall limit to valid range seg0->seg1
			  seg0, end-fade_cnt,    // To fit before the 'bb' fade-out
			  end+fade_cnt, seg1+1); // To fit after
	 
	 // If -1 then we can't fit it in with the provided
	 // segment-length range, so go shorter or longer
	 if (cnt_all < 0) {
	    if (end-fade_cnt > fade_cnt * 2) 
	       cnt_all= end-fade_cnt;
	    else 
	       cnt_all= end+fade_cnt;	// Might overflow datcnt, so this is last resort
	 }
      } 
      cc->cnt_all= cnt_all;

      // Look for a segment starting-point which doesn't overlap
      // either of the segments playing on the other two channels
      {
	 int r0= aa->off;
	 int r1= aa->off + aa->cnt_all;
	 int r2= bb->off;
	 int r3= bb->off + bb->cnt_all;

	 if (r0 > r2) {
	    int tmp;
	    tmp= r0; r0= r2; r2= tmp;
	    tmp= r1; r1= r3; r3= tmp;
	 }
	 
	 cc->off= zxrandM(-1, "rrr", 	// Default -1 if no range is available
			  0, r0 - cnt_all,
			  r1, r2 - cnt_all, 
			  r3, datcnt - cnt_all);
      }

      // Failing that, choose a segment at random
      if (cc->off < 0)
	 cc->off= zxrand(0, datcnt - cnt_all);
      
      // Reseek the OGG stream
      rv= ov_pcm_seek(&cc->ogg, datbase + cc->off);
      if (rv) error("UNEXPECTED: Can't reseek the looping OGG stream: %d", rv);
      continue;
   }
}

// END //
