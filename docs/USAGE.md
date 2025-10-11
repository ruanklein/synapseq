# SynapSeq Usage Guide

## Table of Contents

- [Introduction to SynapSeq](#introduction-to-synapseq)
- [Introduction to Brainwave Entrainment](#introduction-to-brainwave-entrainment)
- [Understanding the Syntax](#understanding-the-syntax)
- [Command Line](#command-line)
- [Notes](#notes)

## Introduction to SynapSeq

This guide will help you get started with SynapSeq, a text-based tool for generating tones to support meditation, relaxation, and altered states of consciousness.

SynapSeq is built on the philosophy of doing one thing well. Providing a concise, human-readable syntax that lets you focus on what you want to achieve, not on how to express it.

### Design Principles

- One way to do it
- Clarity over cleverness
- Intention over syntax
- Less options, more focus
- You write tones, not code

## Introduction to Brainwave Entrainment

Brainwave entrainment is a technique that uses rhythmic stimuli - such as pulsing sounds or flashing lights - to encourage the brain to synchronize its brainwaves with the frequency of the external stimulus. This phenomenon is known as the 'frequency following response', where the brain tends to align its electrical activity to the rhythm it perceives. Depending on the frequency used, this can promote states such as relaxation, focus, or deep sleep.

SynapSeq supports a few types of brainwave entrainment:

1. **Binaural Beats**: When two slightly different frequencies are played in each ear, the brain perceives a third "beat" frequency equal to the difference between the two tones. For example, if 200Hz is played in one ear and 210Hz in the other, the brain perceives a 10Hz beat. Binaural beats **require headphones** to be effective.

2. **Monaural Beats**: In this mode, two binaural beats with mirrored frequency offsets are played in opposite channels. This configuration creates a perceived rhythmic pulse that can be effective even without headphones. The auditory system interprets the resulting interference as a monaural beat, promoting brainwave entrainment similarly to traditional methods. Monaural beats can be used with or without headphones.

3. **Isochronic Tones**: These are single tones that are turned on and off at regular intervals. The brain responds to this rhythmic stimulation and begins to resonate with the frequency. Isochronic tones can be effective with or without headphones.

### Brainwave Frequency Bands

Different frequency ranges correspond to different mental states:

- **Delta (0.5-4 Hz)**: Deep sleep, healing, deep meditation
- **Theta (4-8 Hz)**: Dreaming sleep, meditation, creativity
- **Alpha (8-13 Hz)**: Relaxed alertness, calm, learning
- **Beta (13-30 Hz)**: Active thinking, focus, alertness
- **Gamma (30+ Hz)**: Higher mental activity, peak concentration

## Understanding the Syntax

SynapSeq uses a specific syntax to create a sequence. The syntax is based on the following elements:

- **tone**: A `tone` is a single frequency (known as carrier frequency).
- **noise**: A `noise` is a random signal and is used to create a background sound.
- **background**: A `background` is a user-defined sound that is played in the background.
- **waveform**: A `waveform` is a shape of the tone.
- **silence**: A `silence` is a period of time with no sound.

#### `tone`

The `tone` syntax for Binaural/Monaural/Isochronic is:

```
tone [carrier frequency] [type of brainwave entrainment] [frequency offset] amplitude [amplitude value]
```

Examples:

```
# for binaural beats
tone 400 binaural 10 amplitude 10
# for monaural beats
tone 300 monaural 10 amplitude 10
# for isochronic tones
tone 200 isochronic 10 amplitude 10
# for carrier only (no method)
tone 100 amplitude 10
```

#### `noise`

The `noise` syntax is:

```
noise [type of noise] amplitude [amplitude value]
```

Examples:

```
noise white amplitude 5
noise pink amplitude 15
noise brown amplitude 30
```

#### `background`

The `background` syntax is:

```
background amplitude [amplitude value]
```

Examples:

```
background amplitude 50
```

Also, you can use with effects like a `spin` and `pulse`:

```
background spin 300 rate 7.5 intensity 30 amplitude 50
# Or
background pulse 7.5 intensity 40 amplitude 50
```

The **"spin"** effect creates a sensation of circular sound movement between your ears.

The **"pulse"** effect creates a rhythmic pulsing sensation by modulating the amplitude of the sound.

Each preset can have only one background.

#### `waveform`

The `waveform` could be `sine`, `square`, `triangle`, `sawtooth`.

Waveform is used in: `tone` and `background` (with effects).

Examples:

```
waveform square tone 400 isochronic 10 amplitude 2.5
waveform triangle background spin 500 rate 5.5 intensity 30 amplitude 10
waveform sawtooth background pulse 6 intensity 80 amplitude 50
```

For default, the `waveform` is `sine`.

#### Presets

To create your sequence, you need to define presets before. The presets could be a word to define a state of mind or a any name you want.

The preset syntax is:

```
[preset name]
  [elements] ...
  ...
```

The preset need start with a letter and can contain letters, numbers, underscores, and hyphens.

Examples:

```
alpha
  noise brown amplitude 40
  tone 100 binaural 8 amplitude 15
```

The "alpha" word is a custom preset defined to play a brown noise and a 100Hz binaural beat with 8Hz offset and 15% amplitude.

You can create a many presets as you want.

```
alpha1
  tone 300 isochronic 10 amplitude 15
alpha2
  tone 300 isochronic 8 amplitude 15
alpha3
  tone 250 isochronic 9 amplitude 10
```

Your custom presets can have several tones and other elements:

```
preparation
  noise pink amplitude 30
  tone 100 binaural 7 amplitude 15
  tone 150 binaural 9 amplitude 5
  ...
```

The rules for the presets are:

- Same word can be used only once.
- The elements are separated by a newline.
- The elements starts with 2 indentations after preset line.

#### Timeline

The timeline is a sequence of presets controlled by time. Timeline is defined in end of the file, after the all presets.

Timeline have a start and end time. The start time is 00:00:00 and the end time is the total duration of the sequence.

The timeline syntax is:

```
hh:mm:ss [preset name]
```

Where `hh:mm:ss` is the time in **h**ours, **m**inutes, and **s**econds.

To insert a fade in/out in your timeline, you can use the reserved word **`silence`**.

Examples:

```
00:00:00 silence
00:00:15 alpha1
00:01:00 alpha1
00:02:00 alpha2
00:03:00 silence
```

This sequence starts with a fade-in from `silence` to the `alpha1` preset over 15 seconds. It then maintains the `alpha1` preset until the 1-minute mark. At 2 minutes, it transitions (slides) smoothly from `alpha1` to `alpha2`, and finally fades out to `silence` for the last minute. The total duration of this sequence is 3 minutes.

In SynapSeq V3+, transitions ("slides") between presets are always smooth for elements of the same type (e.g., tones, noise, background) and with matching parameters. However, automatic fade-in/fade-out between different types of tones, waveforms, or background effects is no longer performed by default.

If you want to create a fade between different tone types (for example, from binaural to isochronic, or pink noise to white noise, ...), you must explicitly define how the transition should occur. There are two main ways to do this:

1. **Insert a silence preset in the timeline**
   This creates a clear fade-out and fade-in between different sound types.

For example:

```
# Presets
ps1
  # Track 1
  tone 300 binaural 10 amplitude 10
ps2
  # Track 1
  tone 200 isochronic 5 amplitude 8
ps3
  ...

# Timeline
00:00:00 silence
00:00:15 ps1
00:02:00 ps1
# Fade-out
00:02:30 silence
# Fade-in
00:03:00 ps2
...
```

2. **Manually control amplitudes in your presets**
   For each track, set the amplitude to zero in the preset where it should be silent, and gradually increase or decrease the amplitude over time using the timeline.

For example, to transition from a binaural tone to an isochronic tone without abrupt changes:

Example:

```
# Presets
alpha1
  # Track 1 of alpha1
  tone 300 binaural 10 amplitude 10
  # Track 2 of alpha1
  tone 300 isochronic 10 amplitude 0
alpha2
  # Track 1 of alpha2
  tone 300 binaural 10 amplitude 0
  # Track 2 of alpha2
  tone 300 isochronic 10 amplitude 10
```

Here, the amplitude of each track is explicitly set to zero when it should be silent, ensuring a smooth transition. This approach applies to all types of tones, waveforms, and background effects.

#### Transitions

Transitions control how audio parameters change between two presets over time. When you move from one preset to another in the timeline, SynapSeq smoothly interpolates numerical parameters like amplitude, carrier frequency, resonance, and intensity.

The transition syntax is:

```
hh:mm:ss [preset name] [transition type]
```

Where `[transition type]` can be: `steady`, `ease-out`, `ease-in`, or `smooth`.

**What changes with transitions:**

- Numerical parameters: amplitude, carrier frequency, resonance, intensity
- What doesn't crossfade: tone type (binaural/monaural/isochronic) and waveform shape

These properties switch instantly according to the destination preset.

##### Transition Types

SynapSeq offers four types of transitions, each designed to support different phases of brainwave entrainment:

###### 1. STEADY (default)

The steady transition provides a constant rate of change throughout the entire transition period.

```
Progress:  0% ──── 25% ──── 50% ──── 75% ──── 100%
                 (uniform change)

Transition wave:
20Hz ════════════════════════════════════════════════ 5Hz
     constant rate throughout
```

**Characteristics:**

- Uniform progression from start to finish
- Neutral, mechanical feel
- Predictable and consistent

**Best for:**

- Testing and debugging sequences
- When you want predictable, linear changes
- Technical or experimental sequences

**Brainwave entrainment benefit:** Provides a steady, predictable stimulus that works well as a baseline reference, though it may be more noticeable to the brain than natural transitions.

###### 2. EASE-OUT (logarithmic)

The ease-out transition starts with rapid change and gradually slows down, creating a smooth landing.

```
Progress:  0% ──── 60% ──── 80% ──── 90% ──── 100%
             (fast start, gentle ending)

Transition wave:
20Hz ══════════════════════════════════════════ 5Hz
     fast change ──── gradually slower ──── very gentle

     0%    15%        35%            60%         85%    100%
```

**Characteristics:**

- Most change happens early in the transition
- Gradually decelerates as it approaches the target
- Like a car smoothly braking to a stop

**Best for:**

- Transitioning from high to low frequencies (e.g., beta → theta)
- Relaxation and meditation entry
- Any transition toward slower, deeper states

**Brainwave entrainment benefit:** Mimics the natural process of falling asleep or entering relaxation. The rapid initial change captures the brain's attention and begins the shift, while the gentle ending allows the nervous system to stabilize comfortably in the new state without resistance.

###### 3. EASE-IN (exponential)

The ease-in transition starts gently and accelerates toward the end, creating a smooth departure.

```
Progress:  0% ──── 10% ──── 20% ──── 40% ──── 100%
             (gentle start, fast ending)

Transition wave:
20Hz ══════════════════════════════════════════ 5Hz
     very gentle ──── gradually faster ──── rapid change

     0%    15%            40%          65%        85%   100%
```

**Characteristics:**

- Starts slowly and gradually accelerates
- Most change happens near the end
- Like a car smoothly accelerating from a stop

**Best for:**

- Transitioning from low to high frequencies (e.g., theta → beta)
- Awakening and activation sequences
- Any transition toward faster, alert states

**Brainwave entrainment benefit:** Mirrors the natural awakening process. The gentle start avoids shocking the nervous system when emerging from deep states, while the accelerating finish firmly establishes the new alert state without leaving residual drowsiness.

###### 4. SMOOTH (sigmoid)

The smooth transition provides the most natural feeling, with gentle starts and endings, and faster change in the middle.

```
Progress:  0% ──── 20% ──── 50% ──── 80% ──── 100%
             (slow → fast → slow, S-shaped)

Transition wave:
20Hz ═╗                                    ╔═ 5Hz
      ║           ╱──────────╲            ║
      ║          ╱            ╲           ║
      ╚═════════════════════════════════════╝
      gentle    rapid      rapid    gentle

      0%   5%      25%       50%      75%   95%  100%
```

**Characteristics:**

- Starts slowly (ease-in)
- Accelerates in the middle
- Ends slowly (ease-out)
- S-shaped curve, most natural and organic

**Best for:**

- General-purpose transitions in any direction
- Maximum comfort and minimal perception
- Therapeutic and meditative sessions
- When you want the smoothest possible transition

**Brainwave entrainment benefit:** Provides the most comfortable and natural-feeling transition. The gentle start and end minimize perception of change, while the middle acceleration ensures the transition completes smoothly. This approach mimics natural processes and feels organic to the nervous system, making it ideal for therapeutic and meditative applications. Based on principles of neural adaptation, this is considered the most effective transition for brainwave entrainment, though more research is needed to quantify these benefits.

##### Transition Examples

```
# Presets
awake
  tone 250 binaural 14 amplitude 40  # Low beta (relaxed alertness)

deep
  tone 200 binaural 4 amplitude 20   # Theta (deep meditation)

# Timeline
00:00:00 silence
00:00:30 awake

# Entering meditation - use EASE-OUT
# Rapid initial shift, gentle stabilization in deep state
00:01:00 awake ease-out
00:06:00 deep

# Maintain deep meditation
00:10:00 deep

# Awakening - use EASE-IN
# Gentle emergence, firm arrival in alert state
00:10:30 deep ease-in
00:15:00 awake

00:16:00 silence
```

```
# Natural, comfortable transition - use SMOOTH
00:00:00 silence
00:00:30 preset1
00:05:00 preset1 smooth
00:10:00 preset2
00:12:00 silence
```

If no transition type is specified, **steady** (linear) is used by default.

#### Comments

Comments are only valid if they occupy an entire line by themselves; inline comments (on the same line as other elements) are not allowed and will cause a syntax error.

Comments are ignored by SynapSeq during processing.
The comments syntax is:

```
# [comment]
```

Examples:

```
# A simple comment
alpha
  ...
```

If you use two `#` in the same line, your comment will be printed in the output. Example:

```
## This is a comment that will be printed in the output
alpha
  ...
  ...
```

In SynapSeq execution, the output will be:

```
> This is a comment that will be printed in the output
...
```

Comments are useful to explain your sequence and to help you to remember what you want to achieve.

More examples:

```
preset1
  # This is a comment for preset1
  tone 440 binaural 8 amplitude 5
preset2
  # This is a comment for preset2
  tone 440 binaural 10 amplitude 10
```

### Global Options

SynapSeq has a global options that can be set in the top of the file. All options starts with `@`.

#### `@background`

The `@background` option is used to set the background sound.

The syntax is:

```
@background [path of the background sound]
```

Examples:

```
@background /path/to/background.wav
```

You can use `~` to import audio from your home directory:

```
@background ~/Downloads/rain.wav
```

You can also load background audio directly from the web using HTTP or HTTPS URLs:

```
@background https://example.com/sounds/rain.wav
```

**Background Audio Requirements:**

- SynapSeq supports `.wav` files with 24 Bit and 2 Channels
- The sample rate must match the sequence sample rate (set with `@samplerate` option)
- SynapSeq automatically creates a looping effect for background sounds

The amplitude and optional spin/pulse effects of the background is controlled by the `background` element in the sequence.

For information about file size limits and Content-Type validation when using HTTP/HTTPS URLs, see the [Notes](#notes) section.

#### `@gainlevel`

This option is used to set the gain level of the `@background` sound.

The syntax is:

```
@gainlevel [level]
```

The levels are:

- `verylow`: set the gain to the -20db
- `low`: set the gain to the -16db
- `medium`: set the gain to the -12db
- `high`: set the gain to the -6db
- `veryhigh`: set the gain to the 0db

The `medium` level is the default and is applied to the `@background` sound to avoid any distortion. If you don't want any gain level, you can set the `@gainlevel` to `veryhigh`, this is normal gain of the background sound.

#### `@volume`

This option is used to set the volume of the output.

The syntax is:

```
@volume [volume value]
```

The volume is a value between 0 and 100. The default is 100.

#### `@samplerate`

This option is used to set the sample rate of the output.

The syntax is:

```
@samplerate [samplerate value]
```

The default is 44100.

#### `@presetlist`

The `@presetlist` option allows you to import presets from external files, enabling preset reuse across multiple sequences. This is particularly useful for creating modular, reusable session components.

The syntax is:

```
@presetlist [path to preset file]
```

**Basic Usage Examples:**

Local file:

```
@presetlist /path/to/my-presets.spsq
```

From home directory:

```
@presetlist ~/sequences/relaxation-presets.spsq
```

From the web (HTTP/HTTPS):

```
@presetlist https://example.com/presets/focus-presets.spsq
```

See the `samples/` directory for practical examples of preset file usage, including `presets-relax.spsq`, `presets-focus.spsq`, and their usage in session files like `sample-genesis.spsq`.

**Important Notes:**

- Preset names must be unique across all imported files and local definitions
- Only preset definitions are imported from preset files
- Timeline sections, global options (e.g., `@background`, `@samplerate`), and background elements in preset files will trigger a syntax error
- Background audio must be defined in the main sequence file using the `@background` option, not within preset files

## Command Line

The command line syntax is:

```
synapseq [options] [path of the sequence file] [path of the output file]
```

Example:

```
synapseq sample-binaural.spsq sample-binaural.wav
```

You can use `-` to open sequence from stdin:

```
cat example.spsq | synapseq - output.wav
```

You can also use HTTP or HTTPS URLs to load sequences directly from the web:

```
synapseq https://example.com/sequences/my-sequence.spsq output.wav
```

On \*nix systems, you can also play a sequence in RAW format using other audio tools, such as ffplay or the play command from the sox package. Example:

```
synapseq sample-binaural.spsq - | play -t raw -r 44100 -e signed-integer -b 24 -c 2 -
```

If you want to use another tool to process the output, keep in mind that the audio is emitted in RAW format with the following parameters:

- **Type**: RAW
- **Sample Rate**: 44100 Hz (default, but can be changed using the `@samplerate` option in the session file)
- **Encoding**: Signed Integer
- **Bit Depth**: 24
- **Channels**: 2 (stereo)

Any software used to handle the output must be explicitly configured with these parameters to correctly interpret the audio stream.

#### `-json`

Parse the input file as JSON format.

```
synapseq -json sequence.json output.wav
```

You can also use stdin or HTTP/HTTPS URLs:

```
cat sequence.json | synapseq -json - output.wav
synapseq -json https://example.com/sequence.json output.wav
```

#### `-xml`

Parse the input file as XML format.

```
synapseq -xml sequence.xml output.wav
```

You can also use stdin or HTTP/HTTPS URLs:

```
cat sequence.xml | synapseq -xml - output.wav
synapseq -xml https://example.com/sequence.xml output.wav
```

#### `-yaml`

Parse the input file as YAML format.

```
synapseq -yaml sequence.yaml output.wav
```

You can also use stdin or HTTP/HTTPS URLs:

```
cat sequence.yaml | synapseq -yaml - output.wav
synapseq -yaml https://example.com/sequence.yaml output.wav
```

#### `-help`

Show the help and exit.

#### `-quiet`

Quiet mode. Used to hide terminal output. Errors and comments will be displayed.

#### `-debug`

Debug mode. Used to check file syntax without having to generate the wav file.

#### `-version`

Show the version.

## Notes

### File Size Limits

SynapSeq enforces different file size limits depending on the file type:

- **Text format (`.spsq`)**: Maximum **32 KB** per file

  - Applies to: sequence files and preset files loaded with `@presetlist`
  - Files larger than 32 KB will be truncated

- **Structured formats (JSON, XML, YAML)**: Maximum **128 KB** per file

  - Applies to: files loaded with `-json`, `-xml`, or `-yaml` flags
  - Files larger than 128 KB will be rejected

- **Background audio files (`.wav`)**: Maximum **10 MB** per file
  - Applies to: files loaded with `@background` option
  - Files larger than 10 MB will be read up to the 10 MB limit; the rest will be ignored

### Channel Limits

The total number of tones and noises per timestamp cannot exceed **16 channels**. This limit applies to all formats (text and structured).

### Content-Type Validation for HTTP/HTTPS URLs

When loading files from web URLs, SynapSeq validates the `Content-Type` header returned by the server. If the Content-Type does not match the expected format, the request will be rejected.

#### Text Format Files (.spsq)

For sequence files and preset files loaded via HTTP/HTTPS, the server must return:

- `text/plain`

#### Structured Format Files

**JSON files** must return one of:

- `application/json`
- `text/json`
- Any Content-Type ending with `+json` (e.g., `application/vnd.api+json`)

**XML files** must return one of:

- `application/xml`
- `text/xml`
- Any Content-Type ending with `+xml` (e.g., `application/atom+xml`)

**YAML files** must return one of:

- `application/x-yaml`
- `application/yaml`
- `text/yaml`
- `text/x-yaml`
- Any Content-Type ending with `+yaml` or `+yml`

#### Background Audio Files (.wav)

For background audio files loaded via HTTP/HTTPS, the server must return one of:

- `audio/wav`
- `audio/x-wav`
- `audio/wave`
