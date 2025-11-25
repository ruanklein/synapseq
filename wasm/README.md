# SynapSeq.js - WebAssembly Library

A JavaScript wrapper library for SynapSeq WASM, providing an elegant object-oriented API for generating and playing SynapSeq sequences directly in the browser.

## Features

- Generate binaural/monaural/isochronic tones from SPSQ sequences
- WebAssembly-powered for high performance
- Integrated Web Worker for non-blocking audio generation (no external worker file needed!)
- Single-file library with embedded worker
- Support for local and remote WASM files
- Full JSDoc documentation
- Promise-based async operations
- Built-in audio playback controls

## Installation

Simply include the `synapseq.js` file in your HTML. No need for separate worker files!

```html
<!DOCTYPE html>
<html>
  <head>
    <script src="synapseq.js"></script>
  </head>
  <body>
    <script>
      // Your code here
    </script>
  </body>
</html>
```

## Quick Start

### Basic Usage

```javascript
// Create a new SynapSeq instance
const synapse = new SynapSeq();

// Load a sequence
const spsqCode = `
# Presets
alpha
  noise pink amplitude 30
  tone 250 isochronic 10 amplitude 15

# Timeline
00:00:00 alpha
00:05:00 silence
`;

async function play() {
  await synapse.load(spsqCode);
  await synapse.play();
}

play();
```

### Using Custom Paths

You can specify custom paths for WASM files, enabling CDN usage or custom directory structures:

```javascript
// Local custom paths
const synapse = new SynapSeq({
  wasmPath: "./dist/synapseq.wasm",
  wasmExecPath: "./dist/wasm_exec.js",
});

// Remote CDN
const synapse = new SynapSeq({
  wasmPath: "https://cdn.example.com/synapseq/synapseq.wasm",
  wasmExecPath: "https://cdn.example.com/synapseq/wasm_exec.js",
});

// Mixed (local lib, remote WASM)
const synapse = new SynapSeq({
  wasmPath: "https://cdn.example.com/synapseq.wasm",
  // wasmExecPath defaults to local 'wasm_exec.js'
});
```

## API Reference

### Constructor

#### `new SynapSeq(options)`

Creates a new SynapSeq instance and initializes the embedded Web Worker.

**Parameters:**

- `options` (Object) - Optional configuration object
  - `wasmPath` (string) - Path or URL to the WASM file (default: `'synapseq.wasm'`)
  - `wasmExecPath` (string) - Path or URL to the wasm_exec.js file (default: `'wasm_exec.js'`)

**Returns:** `SynapSeq` instance

**Examples:**

```javascript
// Use default paths (files in same directory)
const synapse = new SynapSeq();

// Use custom local paths
const synapse = new SynapSeq({
  wasmPath: "./assets/synapseq.wasm",
  wasmExecPath: "./assets/wasm_exec.js",
});

// Use remote CDN
const synapse = new SynapSeq({
  wasmPath: "https://cdn.jsdelivr.net/npm/synapseq/synapseq.wasm",
  wasmExecPath: "https://cdn.jsdelivr.net/npm/synapseq/wasm_exec.js",
});
```

---

```javascript
// Load from string
await synapse.load("# Presets\nalpha\n  tone 250 isochronic 10");

// Load from File object
const fileInput = document.getElementById("fileInput");
await synapse.load(fileInput.files[0]);
```

---

#### `play()`

Plays the loaded sequence. Generates audio if not already generated, or resumes if paused.

**Returns:** `Promise<void>`

**Throws:** Error if no sequence is loaded

**Example:**

```javascript
await synapse.play();
```

---

#### `pause()`

Pauses the currently playing sequence.

**Returns:** `void`

**Throws:** Error if no audio is playing

**Example:**

```javascript
synapse.pause();
```

---

#### `stop()`

Stops the currently playing sequence and resets playback position.

**Returns:** `void`

**Example:**

```javascript
synapse.stop();
```

---

#### `getCurrentTime()`

Gets the current playback position in seconds.

**Returns:** `number` - Current time in seconds (0 if not playing)

**Example:**

```javascript
const currentTime = synapse.getCurrentTime();
console.log(`Current position: ${currentTime}s`);
```

---

#### `getDuration()`

Gets the total duration of the loaded audio in seconds.

**Returns:** `number` - Duration in seconds (0 if no audio loaded)

**Example:**

```javascript
const duration = synapse.getDuration();
console.log(`Total duration: ${duration}s`);
```

---

#### `getState()`

Gets the current playback state.

**Returns:** `string` - One of: `'idle'`, `'generating'`, `'playing'`, `'paused'`, `'stopped'`

**Example:**

```javascript
const state = synapse.getState();
console.log(`Current state: ${state}`);
```

---

#### `isLoaded()`

Checks if a sequence is currently loaded.

**Returns:** `boolean` - True if a sequence is loaded

**Example:**

```javascript
if (synapse.isLoaded()) {
  await synapse.play();
}
```

---

#### `isReady()`

Checks if the Web Worker is initialized and ready.

**Returns:** `boolean` - True if worker is ready

**Example:**

```javascript
if (synapse.isReady()) {
  await synapse.load(sequence);
}
```

---

#### `download(filename)`

Downloads the generated WAV file.

**Parameters:**

- `filename` (string) - Name for the downloaded file (default: `'synapseq.wav'`)

**Returns:** `void`

**Throws:** Error if no audio has been generated

**Example:**

```javascript
synapse.download("my-meditation.wav");
```

---

#### `destroy()`

Cleans up resources and terminates the Web Worker.

**Returns:** `void`

**Example:**

```javascript
synapse.destroy();
```

---

### Event Handlers

All event handlers are optional callback functions that can be assigned to handle different states.

#### `onloaded`

Called when a sequence is successfully loaded.

```javascript
synapse.onloaded = () => {
  console.log("Sequence loaded and ready to play");
};
```

---

#### `ongenerating`

Called when audio generation starts.

```javascript
synapse.ongenerating = () => {
  console.log("Generating audio...");
  showLoadingSpinner();
};
```

---

#### `onplaying`

Called when playback starts.

```javascript
synapse.onplaying = () => {
  console.log("Now playing");
  updateUIToPlaying();
};
```

---

#### `onpaused`

Called when playback is paused.

```javascript
synapse.onpaused = () => {
  console.log("Playback paused");
  updateUIToPaused();
};
```

---

#### `onstopped`

Called when playback is stopped.

```javascript
synapse.onstopped = () => {
  console.log("Playback stopped");
  updateUIToStopped();
};
```

---

#### `onended`

Called when playback ends naturally.

```javascript
synapse.onended = () => {
  console.log("Playback finished");
  showCompletionMessage();
};
```

---

#### `onerror`

Called when an error occurs.

**Parameters:**

- `detail` (Object) - Contains `error` property with the Error object

```javascript
synapse.onerror = (detail) => {
  console.error("Error occurred:", detail.error);
  showErrorMessage(detail.error.message);
};
```

---

## Complete Example

```javascript
// Create instance
const synapse = new SynapSeq();

// Setup event handlers
synapse.onloaded = () => {
  console.log("Sequence loaded!");
  document.getElementById("playBtn").disabled = false;
};

synapse.ongenerating = () => {
  console.log("Generating audio...");
  document.getElementById("status").textContent = "Generating...";
};

synapse.onplaying = () => {
  console.log("Playing");
  document.getElementById("status").textContent = "Playing";
  document.getElementById("playBtn").disabled = true;
  document.getElementById("pauseBtn").disabled = false;
};

synapse.onpaused = () => {
  console.log("Paused");
  document.getElementById("status").textContent = "Paused";
  document.getElementById("playBtn").disabled = false;
  document.getElementById("pauseBtn").disabled = true;
};

synapse.onstopped = () => {
  console.log("Stopped");
  document.getElementById("status").textContent = "Stopped";
  document.getElementById("playBtn").disabled = false;
  document.getElementById("pauseBtn").disabled = true;
};

synapse.onended = () => {
  console.log("Finished");
  document.getElementById("status").textContent = "Finished";
};

synapse.onerror = (detail) => {
  console.error("Error:", detail.error);
  alert("Error: " + detail.error.message);
};

// Load and play
async function start() {
  const spsqCode = `
# Presets
alpha-one
  noise pink amplitude 30
  tone 250 isochronic 8 amplitude 15

alpha-two
  noise pink amplitude 25
  tone 250 isochronic 10 amplitude 15

# Timeline
00:00:00 silence
00:00:15 alpha-one
00:02:00 alpha-two
00:04:00 silence
  `;

  try {
    await synapse.load(spsqCode);
    await synapse.play();
  } catch (error) {
    console.error("Failed to start:", error);
  }
}

// Control functions
function pause() {
  synapse.pause();
}

function stop() {
  synapse.stop();
}

function downloadAudio() {
  synapse.download("my-sequence.wav");
}

// Progress tracking
setInterval(() => {
  const current = synapse.getCurrentTime();
  const duration = synapse.getDuration();
  if (duration > 0) {
    const progress = (current / duration) * 100;
    document.getElementById("progressBar").style.width = progress + "%";
  }
}, 100);
```

## SynapSeq Syntax

SPSQ (SynapSeq Sequence) is the offical format used in SynapSeq. See [full documentation.](../docs/USAGE.md) for complete details.

## Browser Compatibility

- Modern browsers with WebAssembly support
- Web Workers support required
- Audio API support required

Tested on:

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Required Files

To use this library, you need these files:

### Minimal Setup (Local Files)

- `synapseq.js` - Main library (contains embedded worker)
- `synapseq.wasm` - Compiled WASM binary
- `wasm_exec.js` - Go WASM runtime

### CDN Setup

You only need to include `synapseq.js` locally or from CDN, then point to remote WASM files:

```html
<script src="synapseq.js"></script>
<script>
  const synapse = new SynapSeq({
    wasmPath: "https://your-cdn.com/synapseq.wasm",
    wasmExecPath: "https://your-cdn.com/wasm_exec.js",
  });
</script>
```

## License

GNU GPL v2 - See [COPYING.txt](../COPYING.txt) for details.

## Links

- [SynapSeq GitHub Repository](https://github.com/ruanklein/synapseq)
- [Full Documentation](../README.md)
- [Usage Guide](../docs/USAGE.md)

---

**SynapSeq** - Synapse-Sequenced Brainwave Generator
Copyright (c) 2025 [ruan.sh](https://ruan.sh)
