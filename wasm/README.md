# SynapSeq.js - WebAssembly Library

A JavaScript wrapper library for SynapSeq WASM, providing an elegant object-oriented API for generating and playing SynapSeq sequences directly in the browser.

## Features

- Real-time streaming audio generation
- Generate binaural/monaural/isochronic tones from SPSQ sequences
- Support for both text (SPSQ) and JSON formats
- WebAssembly-powered for high performance
- AudioWorklet-based streaming for low-latency playback
- Integrated Web Worker for non-blocking audio generation (no external worker file needed!)
- Single-file library with embedded worker
- Support for local and remote WASM files
- Full JSDoc documentation
- Promise-based async operations
- Built-in audio playback controls

## Installation

Simply include the `synapseq.js` file in your HTML:

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
const synapseq = new SynapSeq();

// Load a sequence in text format (SPSQ)
const spsqCode = `
# Presets
alpha
  noise pink amplitude 30
  tone 250 isochronic 10 amplitude 15

# Timeline
00:00:00 silence
00:00:15 alpha
00:04:45 alpha
00:05:00 silence
`;

async function play() {
  // Load from string (text format is default)
  await synapseq.load(spsqCode);
  await synapseq.play();

  // Or load JSON format
  // await synapseq.load(jsonString, "json");

  // Or load from File object (format auto-detected)
  // const fileInput = document.getElementById("fileInput");
  // await synapseq.load(fileInput.files[0]);
}

play();
```

### Using Custom Paths

You can specify custom paths for WASM files, enabling CDN usage or custom directory structures:

```javascript
// Local custom paths
const synapseq = new SynapSeq({
  wasmPath: "./dist/synapseq.wasm",
  wasmExecPath: "./dist/wasm_exec.js",
});

// Remote CDN
const synapseq = new SynapSeq({
  wasmPath: "https://cdn.example.com/synapseq/synapseq.wasm",
  wasmExecPath: "https://cdn.example.com/synapseq/wasm_exec.js",
});

// Mixed (local lib, remote WASM)
const synapseq = new SynapSeq({
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

---

#### `load(input, format)`

Loads a sequence from a string or File object. Supports both text (SPSQ) and JSON formats.

**Parameters:**

- `input` (string|File) - Sequence content as string or File object
- `format` (string) - Optional format: `'text'` (default) or `'json'`

**Returns:** `Promise<void>`

**Throws:** Error if input is invalid

**Examples:**

```javascript
// Load from SPSQ text string (default format)
const spsqCode = `
# Presets
alpha
  tone 250 isochronic 10

# Timeline
00:00:00 alpha
`;
await synapseq.load(spsqCode);
// or explicitly:
await synapseq.load(spsqCode, "text");

// Load from File object (.spsq file)
const fileInput = document.getElementById("fileInput");
await synapseq.load(fileInput.files[0], "text");

// Load from File object (.json file)
await synapseq.load(fileInput.files[0], "json");
```

---

#### `play()`

Plays the loaded sequence. Streams and generates audio in real-time.

**Returns:** `Promise<void>`

**Throws:** Error if no sequence is loaded

**Example:**

```javascript
await synapseq.play();
```

---

#### `stop()`

Stops the currently playing sequence and resets playback position.

**Returns:** `void`

**Example:**

```javascript
synapseq.stop();
```

---

#### `getCurrentTime()`

Gets the current playback position in seconds since playback started.

**Returns:** `number` - Current time in seconds (0 if not playing)

**Example:**

```javascript
const currentTime = synapseq.getCurrentTime();
console.log(`Current position: ${currentTime}s`);
```

---

#### `getSampleRate()`

Gets the sample rate of the loaded sequence.

**Returns:** `number` - Sample rate in Hz

**Example:**

```javascript
const sampleRate = synapseq.getSampleRate();
console.log(`Sample rate: ${sampleRate}Hz`);
```

---

#### `getState()`

Gets the current playback state.

**Returns:** `string` - One of: `'idle'`, `'playing'`, `'stopped'`

**Example:**

```javascript
const state = synapseq.getState();
console.log(`Current state: ${state}`);
```

---

#### `isLoaded()`

Checks if a sequence is currently loaded.

**Returns:** `boolean` - True if a sequence is loaded

**Example:**

```javascript
if (synapseq.isLoaded()) {
  await synapseq.play();
}
```

---

#### `isReady()`

Checks if the Web Worker is initialized and ready.

**Returns:** `boolean` - True if worker is ready

**Example:**

```javascript
if (synapseq.isReady()) {
  await synapseq.load(sequence);
}
```

---

#### `getVersion()`

Gets the SynapSeq version from the WASM module.

**Returns:** `Promise<string>` - The version string

**Example:**

```javascript
const version = await synapseq.getVersion();
console.log(`SynapSeq Version: ${version}`);
```

---

#### `getBuildDate()`

Gets the build date of the SynapSeq WASM module.

**Returns:** `Promise<string>` - The build date string

**Example:**

```javascript
const buildDate = await synapseq.getBuildDate();
console.log(`Build Date: ${buildDate}`);
```

---

#### `getHash()`

Gets the hash of the SynapSeq WASM build.

**Returns:** `Promise<string>` - The hash string

**Example:**

```javascript
const hash = await synapseq.getHash();
console.log(`Build Hash: ${hash}`);
```

---

#### `destroy()`

Cleans up resources and terminates the Web Worker.

**Returns:** `void`

**Example:**

```javascript
synapseq.destroy();
```

---

### Event Handlers

All event handlers are optional callback functions that can be assigned to handle different states.

#### `onloaded`

Called when a sequence is successfully loaded.

```javascript
synapseq.onloaded = () => {
  console.log("Sequence loaded and ready to play");
};
```

---

#### `ongenerating`

Called when audio generation starts.

```javascript
synapseq.ongenerating = () => {
  console.log("Generating audio...");
};
```

---

#### `onplaying`

Called when playback starts.

```javascript
synapseq.onplaying = () => {
  console.log("Now playing");
};
```

---

#### `onstopped`

Called when playback is stopped.

```javascript
synapseq.onstopped = () => {
  console.log("Playback stopped");
};
```

---

#### `onended`

Called when playback ends naturally.

```javascript
synapseq.onended = () => {
  console.log("Playback finished");
};
```

---

#### `onerror`

Called when an error occurs.

**Parameters:**

- `detail` (Object) - Contains `error` property with the Error object

```javascript
synapseq.onerror = (detail) => {
  console.error("Error occurred:", detail.error);
};
```

---

## Complete Example

```javascript
// Create instance
const synapseq = new SynapSeq();

// Setup event handlers
synapseq.onloaded = () => {
  console.log("Sequence loaded!");
  document.getElementById("playBtn").disabled = false;
};

synapseq.ongenerating = () => {
  console.log("Generating audio...");
  document.getElementById("status").textContent = "Generating...";
};

synapseq.onplaying = () => {
  console.log("Playing");
  document.getElementById("status").textContent = "Playing";
  document.getElementById("playBtn").disabled = true;
  document.getElementById("stopBtn").disabled = false;
};

synapseq.onstopped = () => {
  console.log("Stopped");
  document.getElementById("status").textContent = "Stopped";
  document.getElementById("playBtn").disabled = false;
  document.getElementById("stopBtn").disabled = true;
};

synapseq.onended = () => {
  console.log("Finished");
  document.getElementById("status").textContent = "Finished";
  document.getElementById("playBtn").disabled = false;
  document.getElementById("stopBtn").disabled = true;
};

synapseq.onerror = (detail) => {
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
00:03:00 alpha-two
00:04:00 silence
  `;

  try {
    await synapseq.load(spsqCode);
    await synapseq.play();
  } catch (error) {
    console.error("Failed to start:", error);
  }
}

// Control functions
function stop() {
  synapseq.stop();
}

// Progress tracking
setInterval(() => {
  const current = synapseq.getCurrentTime();
  const sampleRate = synapseq.getSampleRate();
  if (current > 0) {
    document.getElementById("time").textContent = `Time: ${current.toFixed(
      1
    )}s @ ${sampleRate}Hz`;
  }
}, 100);
```

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

## License

GNU GPL v2 - See [COPYING.txt](../COPYING.txt) for details.

## Links

- [SynapSeq GitHub Repository](https://github.com/ruanklein/synapseq)
- [Full Documentation](../README.md)
- [Usage Guide](../docs/USAGE.md)

---

**SynapSeq** - Synapse-Sequenced Brainwave Generator

Copyright (c) 2025 [ruan.sh](https://ruan.sh)
