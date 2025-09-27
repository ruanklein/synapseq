# SynapSeq Frequently Asked Questions (FAQ)

## About SynapSeq

### What is SynapSeq?

SynapSeq is an efficient engine for brainwave entrainment, designed to generate audio sequences that guide brainwave states using a simple, human-readable text format. It can be used as a command-line tool (CLI) or as a library, allowing integration into other projects and workflows. SynapSeq helps users create custom soundscapes for meditation, relaxation, focus, and altered states of consciousness.

### What is the difference between SBaGen, SBaGen+, and SynapSeq?

- **SBaGen**: The original C-based brainwave generator, released in the early 2000s. It introduced text-based scripting for brainwave sequences but has a complex syntax and legacy codebase.
- **SBaGen+**: A modernized and extended fork of SBaGen, adding new features, many of which were requested by SBaGen users and implemented based on the original project's TODO list (such as isochronic tones, custom waveforms, and more). Maintained by the SynapSeq author as an intermediate step.
- **SynapSeq**: A complete rewrite in Go, focused on minimalism, explicit control, and ease of use. SynapSeq does **not** depend on SBaGen or SBaGen+ code and offers a simpler, more readable syntax, advanced features (custom waveforms, background audio, smooth transitions), and improved maintainability.

### How do I install SynapSeq?

You need Go (v1.25+) and make installed on your system. See [README](../README.md) for platform-specific installation and compilation instructions.

### How do I use SynapSeq?

Write your sequence in a `.spsq` file using the documented syntax, then run:

```
synapseq my-sequence.spsq output.wav
```

See [USAGE](USAGE.md) for detailed syntax and examples.

### How can I play my sequence?

SynapSeq does not support real-time playback. It generates a WAV file for offline listening. If you want to stream audio directly, you can use the RAW output by redirecting to stdout and piping to an external player (e.g., sox/play or ffplay):

```
synapseq my-sequence.spsq - | play -t raw -r 44100 -e signed-integer -b 24 -c 2 -
```

### What audio formats does SynapSeq support?

SynapSeq outputs 24-bit stereo WAV files by default. You can also pipe raw audio to other tools for playback or conversion.

### Can I use my own background sounds?

Yes! Use the `@background` option at the top of your sequence file to specify a WAV file as background. See [USAGE](USAGE.md#background) for details.

---

## Sequence Creation

### How do I create a fade-in or fade-out?

Use the reserved word `silence` in your timeline to create fades between presets. See [USAGE](USAGE.md#timeline) for examples.

### Can I use multiple tones or noises in one preset?

Yes! Each preset can contain multiple elements (tones, noises, background effects). See [USAGE](USAGE.md#presets).

### How do I control amplitude and transitions?

Set amplitude values for each element in your presets. For smooth transitions, adjust amplitudes over time in your timeline.

### How do I define a preset?

Presets are named sections that group together tones, noises, and background effects. Each preset starts with a name followed by indented elements. For example:

```
alpha
  noise brown amplitude 40
  tone 100 binaural 8 amplitude 15
```

You can create as many presets as you want, and use them in your timeline.

### What is the timeline and how do I use it?

The timeline defines when each preset is played. It uses the format `hh:mm:ss [preset name]` at the end of your sequence file. Example:

```
00:00:00 silence
00:00:15 alpha
00:02:00 alpha
00:03:00 theta
```

### How do I add comments to my sequence?

Comments must be on their own line, starting with `#`. Inline comments are not allowed. Example:

```
# This is a comment
alpha
  tone 440 binaural 8 amplitude 5
```

### How do I set the output volume or sample rate?

Use the global options at the top of your file:

```
@volume 80
@samplerate 48000
```

Volume ranges from 0 to 100 (default is 100). Sample rate defaults to 44100 Hz.

### Can I use custom waveforms?

Yes! You can specify `waveform` before a tone or background effect. Supported waveforms: sine, square, triangle, sawtooth. Example:

```
waveform square tone 400 isochronic 10 amplitude 2.5
```

### How do I use a background WAV file?

Add the `@background` option at the top of your file:

```
@background /path/to/background.wav
```

You can also control its amplitude and apply effects like spin or pulse in your presets.

### What happens if I make a syntax error?

SynapSeq will show an error and refuse to generate audio. Always follow the syntax rules described in the [USAGE](USAGE.md) guide to avoid mistakes. Pay attention to indentation, comments, and the structure of presets and timeline.

---

## Brainwave Entrainment

### What are binaural beats?

Binaural beats occur when two slightly different frequencies are played in each ear. The brain perceives a third "beat" at the difference frequency, which can help entrain brainwaves to a desired state. **Headphones are required** for binaural beats to work.

### What are monaural beats?

Monaural beats are created by mixing two frequencies together before playback. The resulting beat is physically present in the audio and can be perceived without headphones, though headphones may enhance the effect.

### What are isochronic tones?

Isochronic tones are single tones that pulse on and off at regular intervals. This rhythmic stimulation can entrain brainwaves and is effective with or without headphones.

### Which method is best for meditation or focus?

- **Binaural beats**: Best for deep meditation, relaxation, and creativity.
- **Monaural beats**: Good for general entrainment, can be used without headphones.
- **Isochronic tones**: Effective for alertness, focus, and can be used with speakers.

### Do I need headphones?

- **Binaural beats**: Yes, headphones are required.
- **Monaural beats**: Optional, but recommended for best results.
- **Isochronic tones**: Not required.

### What frequencies should I use?

- **Delta (0.5–4 Hz)**: Deep sleep, healing
- **Theta (4–8 Hz)**: Meditation, creativity
- **Alpha (8–13 Hz)**: Relaxed alertness, learning
- **Beta (13–30 Hz)**: Focus, active thinking
- **Gamma (30+ Hz)**: Peak concentration

### Is it safe to use SynapSeq and brainwave entrainment?

For most people, brainwave entrainment is considered safe when used responsibly. However, some individuals (such as those with epilepsy, heart conditions, or neurological disorders) should consult a medical professional before using any entrainment software. General tips:

- Do not use while driving or operating machinery.
- Stop immediately if you feel discomfort, dizziness, or unusual sensations.
- Use moderate volume and avoid excessive session lengths.
- Brainwave entrainment is not a substitute for medical treatment.

### I listened to a brainwave audio I created but didn't feel anything. What could be wrong?

There are several factors that can affect your experience with brainwave entrainment:

- **Technical issues:**
  - Make sure you are using the correct method (e.g., headphones for binaural beats).
  - For binaural beats, the quality of your headphones is especially important; cheap or low-quality headphones may not reproduce the subtle frequency differences needed for the effect. For isochronic tones and monaural beats, both headphones and speakers can be used, but the equipment should have decent frequency response and minimal distortion. You don't need expensive or professional gear, just avoid very cheap or poor-quality devices.
  - Check that the frequencies used are appropriate for your desired state (see frequency guide above).
  - Ensure the audio is not too quiet or distorted; use proper amplitude and volume settings.
  - Avoid converting to low-quality formats (like low-bitrate MP3), which can reduce effectiveness.
- **Personal and environmental factors:**
  - Find a quiet, comfortable place free from distractions.
  - Allow yourself time to relax and focus; effects may take several minutes to be noticed.
  - Some people are more sensitive to brainwave entrainment than others; results can vary.
  - Avoid multitasking or using the audio while working or driving.
  - Try different frequencies, durations, or methods (binaural, monaural, isochronic) to see what works best for you.

**Remember**: brainwave entrainment is a subtle process and may require experimentation and patience. If you still don't notice effects, try adjusting your approach, your equipment, or consult resources on meditation and relaxation techniques.

### Why should I use noise (white, pink, brown) as background? Isn't it just annoying static?

While noise may sound like "annoying static" to some, it plays an important role in brainwave entrainment and audio sessions:

- **Masking distractions:** Noise helps mask external sounds, making it easier to focus or relax.
- **Enhancing entrainment:** A noise background can make tones and beats more effective by smoothing transitions and reducing abrupt changes.
- **Promoting relaxation:** Many people find that noise (especially pink or brown) creates a calming atmosphere, similar to rain or wind.

SynapSeq supports three types of noise:

- **White noise:** Equal energy across all frequencies; sounds like radio static. Good for masking and general use.
- **Pink noise:** More energy in lower frequencies; sounds softer, like steady rain. Often used for relaxation and sleep.
- **Brown noise:** Even more emphasis on low frequencies; sounds deep, like distant thunder or a waterfall. Useful for deep relaxation and meditation.

You can choose the noise type that best fits your session. Try different options to see which feels most comfortable and effective for you.

### Can brainwave entrainment simulate the effects of drugs or medications?

**No**. Despite claims from some companies, there is no scientific evidence that brainwave entrainment (using binaural beats, monaural beats, or isochronic tones) can reproduce or simulate the effects of drugs, medications, or any psychoactive substances. The mechanisms of action for substances are complex and involve biochemical processes in the body that cannot be replicated by audio stimulation alone.

Brainwave entrainment can help with relaxation, focus, meditation, and sleep, but it cannot induce states comparable to those produced by drugs or medications. Any suggestion otherwise is misleading and not supported by credible research. Always be skeptical of products or services that promise such effects.

### What does science say about brainwave entrainment? Is it proven to work?

Scientific research on brainwave entrainment (binaural beats, monaural beats, isochronic tones) has produced mixed results. Some studies suggest that these techniques can help with relaxation, focus, meditation, and sleep, but the effects are generally modest and vary from person to person. There is no broad scientific consensus that brainwave entrainment produces strong or universal effects.

Most evidence points to benefits for stress reduction, mood improvement, and sleep quality, but claims of dramatic cognitive enhancement or medical effects are not supported. If you are interested, look for peer-reviewed studies and reviews in scientific journals. Always approach claims with healthy skepticism and use brainwave entrainment as a complement to, not a replacement for, established health practices.

---

## Usage, Licensing, and Distribution

### Can I use SynapSeq as a library in my own project?

Yes, SynapSeq can be used as a library in your own projects. However, please note that SynapSeq is licensed under the GPL v2. This means that if you use SynapSeq as a library in your project, **your entire project must also comply with the GPL v2 license**. In particular, you cannot use SynapSeq as a library in closed-source or proprietary software. For commercial use, your project must be open source and distributed under a compatible license. Always review the GPL v2 terms before integrating SynapSeq into your application.

### Can I sell audio tracks generated with SynapSeq?

Yes, you can sell audio tracks generated with SynapSeq. However, you are responsible for ensuring that your sequence files and any background sounds used do not infringe on third-party copyrights. If you use your own original sequences and background audio, you are free to distribute or sell the resulting tracks. Always check the license terms of any third-party sounds you include.

### How do I convert the WAV output to MP3 or other formats? Will I lose quality?

WAV files generated by SynapSeq are uncompressed and high quality (24-bit stereo). Converting to MP3 or other lossy formats will reduce audio quality due to compression. For best results:

- Use a high bitrate (at least 256 kbps, ideally 320 kbps) when converting to MP3.
- Prefer lossless formats (e.g., FLAC) if you want to preserve all details.
- Use tools like `ffmpeg` or `sox` for conversion. Example:

```
ffmpeg -i output.wav -codec:a libmp3lame -b:a 320k output.mp3
```

**Technical note**: MP3 compression removes subtle details and may introduce artifacts, especially in brainwave audio. Always keep your original WAV files for best fidelity.

### The output audio is too quiet or distorted. How can I fix this?

- Adjust the `@volume` option in your sequence file. If you are using a background audio with `@background`, you can also tweak the `@gainlevel` option to control its gain.
- Make sure amplitude values are set appropriately for each element.

---

For more information, see [USAGE](USAGE.md) and [README](../README.md).
