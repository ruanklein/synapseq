#include "parser.h"

// Bad sequence file
void badSeq() {
  error("Error in sequence file at line: %d\n  %s", in_lin, lin_copy);
}

// Bad time
void badTime(char *tim) {
  error("Invalid time \"%s\", line %d:\n  %s", tim, in_lin, lin_copy);
}

// Normalize amplitude of voices
void normalizeAmplitude(Voice *voices, int numChannels, const char *line,
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

// Check if background is specified in sequence
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

// Handle options in sequence
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

// Check if voices are equal
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

// Correct periods
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