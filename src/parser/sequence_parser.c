// parser/sequence_parser.c
#include "parser.h"
#include "../core/globals.h"
#include "../utils/error.h"
#include "../utils/memory.h"

#include <string.h>
#include <stdlib.h>
#include <ctype.h>

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