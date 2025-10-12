# SynapSeq Integration Examples

This directory contains Python scripts demonstrating how to integrate SynapSeq into your applications. These examples show how to use structured formats (JSON) and build automated workflows for brainwave entrainment audio generation.

## Prerequisites

- SynapSeq installed and available in PATH
- Python 3.7 or later
- curl (for streaming server examples)
- sox or ffmpeg (optional, for audio playback)

---

## Script 1: Programmatic JSON Generation

**File:** `generate_json_sequence.py`

### What it does

Generates SynapSeq JSON sequences programmatically with customizable parameters. Perfect for automating session creation based on user input or business logic.

### How it works

1. Accepts command-line arguments (start frequency, end frequency, duration)
2. Creates a JSON sequence with progressive frequency transitions
3. Saves the JSON file to disk
4. Invokes SynapSeq CLI to generate the WAV audio file

### Usage

```bash
python3 generate_json_sequence.py <start_freq> <end_freq> <duration_minutes>
```

### Examples

```bash
# Generate 10-minute alpha-to-theta transition (relaxation → deep meditation)
python3 generate_json_sequence.py 10.0 6.0 10

# Generate 15-minute beta-to-alpha transition (focus → relaxation)
python3 generate_json_sequence.py 14.0 10.0 15

# Generate 20-minute alpha-to-theta transition (extended meditation)
python3 generate_json_sequence.py 8.0 4.0 20
```

### Output

- `generated/progressive_<start>Hz-<end>Hz_<duration>min.json` - JSON sequence
- `generated/progressive_<start>Hz-<end>Hz_<duration>min.wav` - Audio file

### Use Cases

- **Automated workflow**: Generate sessions based on user preferences
- **Batch processing**: Create multiple sessions with different parameters
- **Dynamic generation**: Integrate with databases or APIs
- **Template system**: Build reusable sequence templates

---

## Script 2: HTTP Streaming Server

**File:** `streaming_server.py`

### What it does

Creates an HTTP server that generates and streams brainwave entrainment audio in real-time. Includes a web interface for easy testing and integration.

### How it works

1. HTTP server receives requests with audio parameters via query string
2. Generates JSON sequence on-the-fly based on parameters
3. Invokes SynapSeq with RAW output to stdout
4. Streams audio chunks to client via HTTP
5. Client receives and plays/processes audio in real-time

### Usage

Start the server:

```bash
python3 streaming_server.py
```

Access the web interface:

```
http://localhost:8000/
```

### API Endpoint

```
GET /stream?freq=<hz>&duration=<min>&mode=<type>&noise=<type>&carrier=<hz>&amplitude=<0-100>
```

### Examples

**Web browser:**

- Open `http://localhost:8000/` and use the interactive form

**Command-line with curl + sox:**

```bash
# Stream and play 5-minute alpha session
curl -N "http://localhost:8000/stream?freq=10&duration=5" | \
    play -t raw -r 44100 -e signed-integer -b 24 -c 2 -

# Stream and save to RAW file
curl -N "http://localhost:8000/stream?freq=6&duration=10&mode=isochronic" > meditation.raw

# Stream and convert to WAV
curl -N "http://localhost:8000/stream?freq=8&duration=5" | \
    ffmpeg -f s24le -ar 44100 -ac 2 -i - output.wav
```

### Query Parameters

| Parameter   | Type   | Default  | Description                              |
| ----------- | ------ | -------- | ---------------------------------------- |
| `freq`      | float  | 10       | Resonance frequency (Hz)                 |
| `duration`  | int    | 5        | Session duration (minutes)               |
| `mode`      | string | binaural | Tone mode (binaural/monaural/isochronic) |
| `noise`     | string | pink     | Background noise (white/pink/brown)      |
| `carrier`   | int    | 300      | Carrier frequency (Hz)                   |
| `amplitude` | int    | 15       | Tone amplitude (0-100)                   |

### Use Cases

- **Web applications**: Integrate with HTML5 audio players
- **Mobile apps**: Stream audio via HTTP APIs
- **Microservices**: Provide audio generation as a service
- **On-demand generation**: No need to pre-generate or store audio files
- **Real-time processing**: Client receives audio as it's generated

---

## Technical Details

### Audio Format

All examples use SynapSeq's structured JSON format for sequence definition:

- **Advantages**: Easy to generate programmatically, validate, and version control
- **Format**: Standard JSON with sequences, tracks, tones, and noises
- **Output**: 24-bit stereo WAV (44100Hz sample rate)

### RAW Streaming Format

The streaming server uses RAW PCM audio format:

- **Type**: RAW (PCM)
- **Sample Rate**: 44100 Hz
- **Bit Depth**: 24 bits
- **Channels**: 2 (stereo)
- **Encoding**: Signed Integer (little-endian)

### Subprocess Integration

Both scripts use Python's `subprocess` module to invoke SynapSeq CLI:

```python
# JSON generation example
subprocess.run(['synapseq', '-json', 'input.json', 'output.wav'])

# RAW streaming example
subprocess.Popen(['synapseq', '-json', 'input.json', '-'], stdout=PIPE)
```

---

## Learn More

- [SynapSeq Documentation](../docs/USAGE.md) - Complete syntax reference
- [Structured Format Samples](../samples/structured/) - JSON/XML/YAML examples
- [FAQ](../docs/FAQ.md) - Common questions and integration patterns
- [GitHub Repository](https://github.com/ruanklein/synapseq) - Source code and issues

---

## Contributing

Have ideas for new integration examples? Found a bug? Want to improve the scripts?

- Submit issues: [GitHub Issues](https://github.com/ruanklein/synapseq/issues)
- Discuss integrations: [GitHub Discussions](https://github.com/ruanklein/synapseq/discussions)
