// SynapSeq WASM Worker
importScripts("wasm_exec.js");

let wasmReady = false;
const go = new Go();

// Initialize WASM
WebAssembly.instantiateStreaming(fetch("synapseq.wasm"), go.importObject)
  .then((result) => {
    go.run(result.instance);
    wasmReady = true;
    self.postMessage({ type: "ready" });
  })
  .catch((error) => {
    self.postMessage({
      type: "error",
      error: "Failed to load WASM: " + error.message,
    });
  });

// Listen for messages from main thread
self.onmessage = async function (e) {
  if (e.data.type === "generate") {
    if (!wasmReady) {
      self.postMessage({
        type: "error",
        error: "WASM not initialized yet",
      });
      return;
    }

    try {
      const spsqBytes = e.data.spsqBytes;

      // Call WASM function (now returns a Promise)
      const result = await synapseqGenerate(spsqBytes);

      if (result.error) {
        self.postMessage({
          type: "error",
          error: result.error,
        });
        return;
      }

      // Send the WAV bytes back to main thread
      self.postMessage(
        {
          type: "success",
          wav: result.wav,
        },
        [result.wav.buffer]
      ); // Transfer ownership for better performance
    } catch (error) {
      self.postMessage({
        type: "error",
        error: error.message || error || "Unknown error occurred",
      });
    }
  }
};
