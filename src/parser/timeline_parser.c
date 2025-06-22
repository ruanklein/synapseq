// parser/timeline_parser.c
#include "parser.h"
#include "../core/globals.h"
#include "../utils/error.h"
#include "../utils/memory.h"

#include <string.h>
#include <stdlib.h>
#include <ctype.h>

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

// Read time
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