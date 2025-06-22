// parser/preset_parser.c
#include "parser.h"
#include "../core/globals.h"
#include "../utils/error.h"
#include "../utils/memory.h"

#include <string.h>
#include <stdlib.h>
#include <ctype.h>

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