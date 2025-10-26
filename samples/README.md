# SynapSeq Sequence Files

This directory contains example sequence files demonstrating various features and capabilities of SynapSeq. Each sample showcases different aspects of brainwave entrainment technology and audio synthesis.

## Text Format

The `sample-*` files is written in the SynapSeq text format, illustrating various techniques and presets. Below is a brief description of each sample file.

### sample-binaural.spsq

**Demonstrates:** Binaural beats technology with progressive alpha frequencies

**Description:** This sequence uses binaural beats with a 250Hz carrier tone, progressively moving through alpha frequencies (8-11Hz). It includes pink noise as a background element that gradually decreases in amplitude to emphasize the binaural tones.

**Audio Intent:** Deep relaxation, enhanced focus, and improved learning. Ideal for study sessions, meditation, or creative work. Requires stereo headphones for optimal effect.

**Timeline:** 7 minutes total, starting with 8Hz (deep relaxation), transitioning through 9Hz (enhanced focus), 10Hz (classic alpha state), and ending with 11Hz (mental clarity).

### sample-monaural.spsq

**Demonstrates:** Monaural beats technology with alpha frequencies

**Description:** Features monaural beats at 250Hz carrier with 8-11Hz entrainment frequencies. Unlike binaural beats, monaural beats are physically mixed before reaching the ears, creating an audible beating pattern.

**Audio Intent:** Deep relaxation and creative state induction. More direct entrainment effect compared to binaural beats. Works with speakers or headphones.

**Timeline:** 7 minutes total, progressing through the same alpha range as the binaural sample but with physically audible beating patterns.

### sample-isochronic.spsq

**Demonstrates:** Isochronic tones with distinct pulses

**Description:** Uses isochronic tones at 250Hz carrier with 8-11Hz pulsing frequencies. Creates clear, distinct pulses of sound by turning the tone on and off at regular intervals.

**Audio Intent:** Most effective entrainment method due to clear, rhythmic pulses. Ideal for active meditation, study sessions, and focus enhancement. No headphones required.

**Timeline:** 7 minutes total with progressively increasing pulse frequencies through the alpha range.

### sample-noise.spsq

**Demonstrates:** Different noise types (brown, pink, white)

**Description:** Comparison sequence of the three main noise types. Brown noise has more energy in lower frequencies, pink noise has equal energy per octave, and white noise has equal energy across all frequencies.

**Audio Intent:** Educational comparison to help users identify which noise type works best for their specific needs: relaxation, focus, sleep, or sound masking.

**Timeline:** Short 1-minute sequence with 10-second samples of each noise type.

### sample-waveform.spsq

**Demonstrates:** Four fundamental waveform types (sine, square, triangle, sawtooth)

**Description:** Showcases how different waveforms affect the sonic character and therapeutic application of tones. Each waveform has distinct harmonic content and psychological effects.

**Audio Intent:** Educational demonstration of waveform characteristics. Helps users understand which waveform is best suited for different mental states: sine for relaxation, square for alertness, triangle for balanced stimulation, and sawtooth for complex entrainment.

**Timeline:** 1 minute 20 seconds with brief demonstrations of each waveform type.

### sample-background-spin.spsq

**Demonstrates:** Spin effect with background audio

**Description:** Creates a sensation of circular sound movement between the ears using pink noise background. The spin effect modulates the amplitude between left and right channels to create spatial rotation.

**Audio Intent:** Relaxation and focus enhancement through spatial audio manipulation. Combined with binaural tones for deeper entrainment effect.

**Timeline:** 10 minutes with varying spin rates and widths to demonstrate the effect's versatility.

**Requirements:** Pink noise WAV file at `sounds/pink-noise.wav`

### sample-background-pulse.spsq

**Demonstrates:** Pulse effect with background audio

**Description:** Creates rhythmic pulsing by modulating the amplitude of the pink noise background. Works similarly to isochronic tones but applied to background sound.

**Audio Intent:** Relaxation and focus through rhythmic audio stimulation. Combined with monaural tones for enhanced effect.

**Timeline:** 10 minutes with different pulse rates and intensities.

**Requirements:** Pink noise WAV file at `sounds/pink-noise.wav`

### sample-transitions.spsq

**Demonstrates:** Four transition types (steady, ease-out, ease-in, smooth)

**Description:** Comprehensive demonstration of all available transition types between high (12Hz binaural) and low (2Hz binaural) frequency states. Each transition type has different progression curves suited for specific entrainment scenarios.

**Audio Intent:** Educational demonstration of transition mechanics. Shows how different transition curves affect the listening experience and entrainment effectiveness.

**Timeline:** 12 minutes showcasing steady (linear), ease-out (logarithmic), ease-in (exponential), and smooth (sigmoid) transitions.

### sample-relax-one.spsq

**Demonstrates:** Preset reuse with `@presetlist` option

**Description:** First part of a relaxation session focusing on preparation and alpha state. Uses presets defined in `presets-relax.spsq`.

**Audio Intent:** Gentle relaxation and stress reduction. Preparation phase for deeper meditation work.

**Timeline:** 8 minutes of preparation and alpha-phase alternation.

**Requirements:** `presets-relax.spsq` file

### sample-relax-two.spsq

**Demonstrates:** Advanced preset usage with deeper states

**Description:** Second part of relaxation session focusing on theta frequencies and gentle closing. Continues from the alpha work in part one.

**Audio Intent:** Deeper relaxation and meditative states. Suitable for extended meditation sessions or pre-sleep relaxation.

**Timeline:** 9 minutes including theta-phase work and gentle closing.

**Requirements:** `presets-relax.spsq` file

### sample-focus-one.spsq

**Demonstrates:** Focus-oriented session structure

**Description:** First part of focus session with gentle activation into focused state. Uses presets from `presets-focus.spsq`.

**Audio Intent:** Gradual awakening and focus enhancement. Ideal for starting work sessions or study periods.

**Timeline:** 9 minutes of activation and focus phases.

**Requirements:** `presets-focus.spsq` file

### sample-focus-two.spsq

**Demonstrates:** Sustained focus with higher frequencies

**Description:** Second part of focus session maintaining and deepening alertness with higher beta frequencies and isochronic tones.

**Audio Intent:** Sustained concentration and increased alertness. Suitable for demanding cognitive tasks.

**Timeline:** 10 minutes alternating between focus and deep-focus states.

**Requirements:** `presets-focus.spsq` file

### sample-genesis.spsq

**Demonstrates:** Multi-phase session combining multiple preset files

**Description:** Complete session combining both relaxation and focus presets. Starts with deep relaxation, then transitions to active focus and alertness.

**Audio Intent:** Full-spectrum session for complete mental state transformation. From deep relaxation through to high alertness.

**Timeline:** 30 minutes divided into relaxation phase (0-12 minutes) and focus phase (12-30 minutes).

**Requirements:** Both `presets-relax.spsq` and `presets-focus.spsq` files

### sample-test.spsq

**Demonstrates:** Pure tone generation (no beats)

**Description:** Simple test sequence using pure tone at 440Hz without any beat frequencies. Used for testing and validating the pure tone generation feature.

**Audio Intent:** Technical testing and audio system validation.

**Timeline:** 1 minute with 45 seconds of pure 440Hz tone.

### presets-relax.spsq

**Purpose:** Reusable presets for relaxation-oriented sessions

**Contains:** Preparation, alpha-phase, theta-phase, and closing presets using brown noise and binaural tones at various frequencies (5-12Hz).

**Usage:** Include in other sequences using `@presetlist presets-relax.spsq`

### presets-focus.spsq

**Purpose:** Reusable presets for focus-oriented sessions

**Contains:** Activation, focus, alert, and deep-focus presets using brown noise and isochronic tones at higher frequencies (14-16Hz).

**Usage:** Include in other sequences using `@presetlist presets-focus.spsq`

### sample-ladder.spsq

**Demonstrates:** Template preset and binaural ladder structure

**Description:** This example uses the template preset feature to create a "binaural ladder", a sequence of tones with a fixed binaural beat (8 Hz) and progressively higher carrier frequencies. The `ladder-base` template defines all steps (tones) and the pink noise background. Each derived preset (`ladder-one`, `ladder-two`, etc.) activates a specific step by increasing the amplitude of that tone and keeping the others silent. This allows for a clear and modular progression without code repetition.

**Audio Intent:** Demonstrates how to create progressive brain activation sequences, keeping the beat frequency constant and varying the carrier. Ideal for studies on binaural beat perception and for gradual activation sessions.

**Timeline:** 16 minutes, each step lasts about 1.5 to 2 minutes, starting from the lowest tone to the highest, with pink noise background.

**Requirements:** Headphones are required to perceive the binaural effect.

## Structured Format

The `structured/` subdirectory contains samples in JSON, XML, and YAML formats demonstrating how to use SynapSeq with structured data formats. See [structured/README.md](structured/README.md) for details.

## Background Audio Files

The `sounds/` directory contains:

- `pink-noise.wav` - Pink noise background used by spin and pulse effect samples

## Usage Tips

1. Start with basic samples (`sample-binaural.spsq`, `sample-noise.spsq`) to understand core concepts
2. Explore effect samples (`sample-background-spin.spsq`, `sample-background-pulse.spsq`) to learn spatial audio features
3. Study transition sample (`sample-transitions.spsq`) to understand smooth state changes
4. Use session examples as templates for creating your own sequences
5. Always use headphones for binaural beats; other methods work with speakers too

## Creating Your Own Sequences

Use these samples as starting points. Modify frequencies, amplitudes, timelines, and effects to suit your specific needs. See [../docs/USAGE.md](../docs/USAGE.md) for complete syntax reference.
