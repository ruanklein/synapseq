# <img style="border-radius: 15%;" src="../build/assets/synapseq.png" alt="SynapSeq Logo" width="32" height="32"> SynapSeq Usage Guide

## 📑 Table of Contents

- [Introduction to SynapSeq](#-introduction-to-synapseq)
- [Introduction to Brainwave Entrainment](#-introduction-to-brainwave-entrainment)
- [Understanding the Syntax](#-understanding-the-syntax)

## 📑 Introduction to SynapSeq

This guide will help you get started with SynapSeq, a text-based tool for generating tones to support meditation, relaxation, and altered states of consciousness.

SynapSeq is built on the philosophy of doing one thing well. Providing a concise, human-readable syntax that lets you focus on what you want to achieve, not on how to express it.

Instead of flexibility, it offers clarity. Instead of options, it offers direction.

### Design Principles

- One way to do it
- Clarity over cleverness
- Intention over syntax
- Less options, more focus
- You write tones, not code

## 🧠 Introduction to Brainwave Entrainment

Brainwave entrainment is a method to stimulate the brain into entering a specific state by using a pulsing sound, light, or electromagnetic field. The pulses elicit the brain's 'frequency following' response, encouraging the brainwaves to align to the frequency of the given beat.

SynapSeq supports a few types of brainwave entrainment:

1. **Binaural Beats**: When two slightly different frequencies are played in each ear, the brain perceives a third "beat" frequency equal to the difference between the two tones. For example, if 200Hz is played in one ear and 210Hz in the other, the brain perceives a 10Hz beat. Binaural beats **require headphones** to be effective.

2. **Monaural Beats**: In this mode, two binaural beats with mirrored frequency offsets are played in opposite channels. This configuration creates a perceived rhythmic pulse that can be effective even without headphones. The auditory system interprets the resulting interference as a monaural beat, promoting brainwave entrainment similarly to traditional methods. Monaural beats can be used with or without headphones.

3. **Isochronic Tones**: These are single tones that are turned on and off at regular intervals. The brain responds to this rhythmic stimulation and begins to resonate with the frequency. Isochronic tones can be effective with or without headphones.

### 📊 Brainwave Frequency Bands

Different frequency ranges correspond to different mental states:

- **Delta (0.5-4 Hz)**: Deep sleep, healing, deep meditation
- **Theta (4-8 Hz)**: Dreaming sleep, meditation, creativity
- **Alpha (8-13 Hz)**: Relaxed alertness, calm, learning
- **Beta (13-30 Hz)**: Active thinking, focus, alertness
- **Gamma (30+ Hz)**: Higher mental activity, peak concentration

## 📝 Understanding the Syntax

SynapSeq uses a specific syntax to create a sequence. The syntax is based on the following elements:

- **tone**: A `tone` is a single frequency (known as carrier frequency).
- **noise**: A `noise` is a random signal and is used to create a background sound.
- **background**: A `background` is a user-defined sound that is played in the background.
- **effect**: A `effect` is a background effect that is applied to the background sound.
- **spin**: A `spin` is a type of noise that creates a gentle, binaural-like pulsing between left and right ears.
- **waveform**: A `waveform` is a shape of the tone.
- **silence**: A `silence` is a period of time with no sound.

#### `tone`

The `tone` syntax is:

```
tone [carrier frequency] [type of brainwave entrainment] [frequency offset] amplitude [amplitude value]
```

Examples:

```
tone 400 binaural 10 amplitude 10 # for binaural beats
tone 300 monaural 10 amplitude 10 # for monaural beats
tone 200 isochronic 10 amplitude 10 # for isochronic tones
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

#### `effect`

_Note_: The `effect` is only valid with the `background`.

`effect` provides with a two types of effects: `spin` and `pulse`.

The `effect` syntax is:

```
effect spin width [width value] rate [rate value] intensity [intensity value]
effect pulse [pulse value] intensity [intensity value]
```

Examples:

```
effect spin width 500 rate 1.5 intensity 80
effect pulse 8.5 intensity 90
```

The `effect spin` creates a same `spin` effect, but applied to the background instead of the noise.

The `effect pulse` creates a pulse effect on the background.

#### `spin`

The `spin` syntax is:

```
spin [type of spin] width [width value] rate [rate value] amplitude [amplitude value]
```

Examples:

```
spin white width 400 rate 4.0 amplitude 10
spin pink width 300 rate 2.0 amplitude 25
spin brown width 200 rate 1.0 amplitude 40
```

#### `waveform`

The `waveform` could be `sine`, `square`, `triangle`, `sawtooth`.

Waveform is used in: `tone`, `spin`, and `effect`.

Examples:

```
waveform square tone 400 iso 10 amplitude 2.5
waveform triangle spin pink width 500 rate 5.5 amplitude 10
waveform sawtooth effect pulse 6 intensity 80
```

For default, the `waveform` for `tone`, `spin`, and `effect` is `sine`.

#### Notes

- If the sum of the amplitude of the `tone`, `noise`, `spin`, and `background` is greater than 100, the amplitude will be normalized to 100 without any warning.

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
  tone 100 bin 8 amplitude 15
```

The "alpha" word is a custom preset defined to play a brown noise and a 100Hz binaural beat with 8Hz offset and 15% amplitude.

You can create a many presets as you want.

```
alpha1
  tone 300 iso 10 amplitude 15
alpha2
  tone 300 iso 8 amplitude 15
alpha3
  tone 250 iso 9 amplitude 10
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

To insert a fade in/out or silence in your timeline, you can use the reserved word **`silence`**.

Examples:

```
00:00:00 silence
00:00:15 alpha1
00:01:00 alpha1
00:02:00 alpha2
00:03:00 silence
```

This creates a fade in to alpha1 preset for 15 seconds, and keeps it alpha1 preset until 1 minute, then slides to alpha2 preset for 1 minute, and fade out to silence for 1 minute. The total duration of this sequence is 3 minutes.

Another example:

```
00:00:00 silence
00:00:30 beta
00:01:00 beta
00:02:00 alpha
00:05:00 theta
00:15:00 delta
00:20:00 theta
00:25:00 alpha
00:29:00 alpha
00:30:00 silence
```

In SynapSeq, all is slide to the next preset. In other words, if you don't define a fixed time for a preset, it will slide to the next preset.

Slide is a smooth transition between presets. It is a default behavior in SynapSeq. It is valid for `tone`, `noise`, `spin`, `background`, and `effect`.

But, if your next preset is a different tone of brainwave entrainment, a different waveform, effect, spin or noise, it will not slide. It will create automatically a fade in/out to the next preset.

Example:

```
alpha1
  tone 300 binaural 10 amplitude 10 # Voice 1 of alpha1
alpha2
  tone 300 isochronic 10 amplitude 10 # Voice 1 of alpha2
```

In this example, the `alpha1` (with voice 1) preset has a binaural beat, and the `alpha2` (with voice 1) preset has an isochronic tone. Because the tone type is different in the same voices (voice 1), it will not slide. It will create automatically a fade in/out to the next preset.

Another example:

```
theta1
  noise brown amplitude 40 # Voice 1 of theta1
  tone 125 binaural 7.0 amplitude 10 # Voice 2 of theta1
theta2
  noise pink amplitude 40 # different noise, create a fade in/out from theta1 to theta2
  tone 125 binaural 7.0 amplitude 15 # same tone type, slide from theta1 to theta2
```

#### Comments

You can use comments in your sequence. Comments are ignored by SynapSeq in the processing.

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

If you use two `#` in the same line, your comment will printed in the output. Example:

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
  tone 440 binaural 10 amplitude 10 # This is a comment for preset2
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

The SynapSeq support `.wav`, `.ogg`, and `.mp3` files. For default, SynapSeq creates a looping for the background sound.

The amplitude of the background is controlled by the `background` element in the sequence.

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

#### `@waveform`

This option is used to set the waveform for all elements in the sequence.

The syntax is:

```
@waveform [waveform]
```

The waveform could be `sine`, `square`, `triangle`, `sawtooth`. `sine` is the default.

**Note**: If you set waveform manually in the sequence, the `@waveform` option will be ignored.

#### `@samplerate`

This option is used to set the sample rate of the output.

The syntax is:

```
@samplerate [samplerate value]
```

The default is 44100.

#### `@verbose`

This option is used to set the verbose mode of the output.

The syntax is:

```
@verbose
```

#### `@quiet`

This option is used to set the quiet mode of the output.

The syntax is:

```
@quiet
```

#### `@test`

This option is used to test the sequence.

The syntax is:

```
@test
```

The test mode is used to test the sequence and to see the output in the console.

## Command Line

The command line syntax is:

```
synapseq [options] [path of the sequence file]
```

You can open a sequence file through pipe with stdin. Example:

```
cat sequence.spsq | synapseq -
```

#### `--help`

Show the help and exit.

#### `--test`

Test sequence. Equivalent to `@test` option.

#### `--quiet`

Quiet mode. Equivalent to `@quiet` option.

#### `--verbose`

Verbose mode. Equivalent to `@verbose` option.

#### `--volume`

The syntax is:

```
--volume [value]
```

The value is a number between 0 and 100. The default is 100.
Equivalent to `@volume` option.

#### `--waveform`

Set the waveform for all elements in the sequence.

The syntax is:

```
--waveform [waveform]
```

The waveform could be `sine`, `square`, `triangle`, `sawtooth`. `sine` is the default.

#### `--sample-rate`

Set the sample rate of the output.

The syntax is:

```
--samplerate [samplerate value]
```

The default is 44100.

#### `--output-raw-file`

The syntax is:

```
--output-raw-file [path of the output file]
```

The output file is a raw audio data file.

#### `--output-wav-file`

The syntax is:

```
--output-wav-file [path of the output file]
```

The output file is a WAV file.

#### `--raw`

The syntax is:

```
--raw [path of the output file]
```

Print raw audio data to the console insted of playing the sequence.

#### `--wav`

The syntax is:

```
--wav [path of the output file]
```

Print WAV data to the console insted of playing the sequence.

#### `--buffer-size` (macOS)

The syntax is:

```
--buffer-size [buffer size]
```

The buffer size (in samples) is a number between 1024 and 4096. The default is 2048.

#### `--device` (Linux)

The syntax is:

```
--device [device name]
```

The ALSA device name is the name of the device to use. The default is `default`.
