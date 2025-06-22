// utils/memory.c
#include "memory.h"
#include "error.h"

#include <stdlib.h>
#include <string.h>

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