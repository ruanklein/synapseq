// parser/lexer.c
#include "parser.h"
#include "../core/globals.h"
#include "../utils/error.h"

#include <string.h>
#include <ctype.h>

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