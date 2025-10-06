# Structured Samples (JSON / XML / YAML)

This folder contains SynapSeq samples using structured formats (JSON, XML, YAML). These formats are best when:

- You generate sequences programmatically (pipelines, scripts, external tools).
- You need machine-readable, schema-friendly files for validation or integration.
- You want explicit control of channel states per timestamp (no implicit presets).

If you prefer hand-authoring or preset reuse, the text format (.spsq) is often more ergonomic.

## CLI usage

- JSON:
  synapseq -json sample-binaural.json ~/Downloads/output.wav

- XML:
  synapseq -xml sample-binaural.xml ~/Downloads/output.wav

- YAML:
  synapseq -yaml sample-binaural.yaml ~/Downloads/output.wav

## Format overview

Root object:

- description: array of strings (converted from lines starting with "##" in .spsq).
- options:
  - samplerate: integer (e.g., 44100)
  - volume: integer 0–100
- sequence: array of entries, each entry:
  - time: integer milliseconds (first must be 0, strictly increasing)
  - track:
    - tones: array of tone objects
      - mode: binaural | monaural | isochronic | pure
      - carrier: Hz (float)
      - resonance: Hz (float) — for beats
      - amplitude: 0–100 (float)
      - waveform: sine | square | triangle | sawtooth
    - noises: array of noise objects
      - mode: white | pink | brown
      - amplitude: 0–100 (float)

Notes:

- Silence is explicit: use amplitude 0 at start/end when you want fade-in/out.
- The total size of tones + noises arrays per entry must not exceed 16. This is an internal SynapSeq rule and also applies to .spsq text files.
