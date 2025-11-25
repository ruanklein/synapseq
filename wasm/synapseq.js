/**
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 *
 * JavaScript wrapper for SynapSeq WASM
 *
 * @class SynapSeq
 * @example
 * const synapse = new SynapSeq();
 * await synapse.load(spsqContent);
 * await synapse.play();
 */
class SynapSeq {
  /**
   * Creates a new SynapSeq instance
   * @constructor
   * @param {Object} options - Configuration options
   * @param {string} [options.wasmPath='synapseq.wasm'] - Path or URL to the WASM file
   * @param {string} [options.wasmExecPath='wasm_exec.js'] - Path or URL to the wasm_exec.js file
   * @example
   * // Use local files (default)
   * const synapse = new SynapSeq();
   *
   * // Use custom paths
   * const synapse = new SynapSeq({
   *   wasmPath: './dist/synapseq.wasm',
   *   wasmExecPath: './dist/wasm_exec.js'
   * });
   *
   * // Use remote CDN
   * const synapse = new SynapSeq({
   *   wasmPath: 'https://cdn.example.com/synapseq.wasm',
   *   wasmExecPath: 'https://cdn.example.com/wasm_exec.js'
   * });
   */
  constructor(options = {}) {
    /**
     * @private
     * @type {string}
     */
    this._wasmPath = options.wasmPath || "synapseq.wasm";

    /**
     * @private
     * @type {string}
     */
    this._wasmExecPath = options.wasmExecPath || "wasm_exec.js";

    /**
     * @private
     * @type {string|null}
     */
    this._sequence = null;

    /**
     * @private
     * @type {Worker|null}
     */
    this._worker = null;

    /**
     * @private
     * @type {HTMLAudioElement|null}
     */
    this._audio = null;

    /**
     * @private
     * @type {Blob|null}
     */
    this._audioBlob = null;

    /**
     * @private
     * @type {boolean}
     */
    this._workerReady = false;

    /**
     * @private
     * @type {string}
     */
    this._version = "unknown";

    /**
     * @private
     * @type {string}
     */
    this._buildDate = "";

    /**
     * @private
     * @type {string}
     */
    this._hash = "";

    /**
     * @private
     * @type {Promise<void>|null}
     */
    this._initPromise = null;

    this._initializeWorker();
  }

  /**
   * Creates Web Worker from function
   * @private
   * @returns {Worker}
   */
  _createWorker() {
    // Worker code as a standalone function
    function workerFunction() {
      let wasmReady = false;
      let wasmPath = "";
      let wasmExecPath = "";

      self.onmessage = async function (e) {
        if (e.data.type === "init") {
          wasmPath = e.data.wasmPath;
          wasmExecPath = e.data.wasmExecPath;

          try {
            const baseUrl = e.data.baseUrl || "";
            const absoluteExecPath =
              wasmExecPath.startsWith("http") || wasmExecPath.startsWith("/")
                ? wasmExecPath
                : baseUrl + wasmExecPath;
            const absoluteWasmPath =
              wasmPath.startsWith("http") || wasmPath.startsWith("/")
                ? wasmPath
                : baseUrl + wasmPath;

            importScripts(absoluteExecPath);

            const go = new Go();
            const result = await WebAssembly.instantiateStreaming(
              fetch(absoluteWasmPath),
              go.importObject
            );

            go.run(result.instance);
            wasmReady = true;

            self.postMessage({
              type: "ready",
              version: synapseqVersion,
              buildDate: synapseqBuildDate,
              hash: synapseqHash,
            });
          } catch (error) {
            self.postMessage({
              type: "error",
              error: "Failed to load WASM: " + error.message,
            });
          }
          return;
        }

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
            const result = await synapseqGenerate(spsqBytes);

            if (result.error) {
              self.postMessage({
                type: "error",
                error: result.error,
              });
              return;
            }

            self.postMessage(
              {
                type: "success",
                wav: result.wav,
              },
              [result.wav.buffer]
            );
          } catch (error) {
            self.postMessage({
              type: "error",
              error: error.message || error || "Unknown error occurred",
            });
          }
        }
      };
    }

    const workerCode = `(${workerFunction.toString()})();`;
    const blob = new Blob([workerCode], { type: "application/javascript" });
    return new Worker(URL.createObjectURL(blob));
  }

  /**
   * Initializes the Web Worker for WASM processing
   * @private
   * @returns {Promise<void>}
   */
  _initializeWorker() {
    this._initPromise = new Promise((resolve, reject) => {
      try {
        this._worker = this._createWorker();

        this._worker.onmessage = (e) => {
          const data = e.data;

          if (data.type === "ready") {
            this._workerReady = true;

            this._version = data.version || "unknown";
            this._buildDate = data.buildDate || "";
            this._hash = data.hash || "";

            resolve();
          } else if (data.type === "success") {
            this._handleAudioGenerated(data.wav);
          } else if (data.type === "error") {
            this._handleError(new Error(data.error));
          }
        };

        this._worker.onerror = (error) => {
          reject(new Error("Worker initialization failed: " + error.message));
        };

        // Send initialization message with paths
        const baseUrl =
          typeof window !== "undefined"
            ? window.location.href.substring(
                0,
                window.location.href.lastIndexOf("/") + 1
              )
            : "";

        this._worker.postMessage({
          type: "init",
          wasmPath: this._wasmPath,
          wasmExecPath: this._wasmExecPath,
          baseUrl: baseUrl,
        });
      } catch (error) {
        reject(new Error("Failed to create worker: " + error.message));
      }
    });

    return this._initPromise;
  }

  /**
   * Handles successful audio generation from worker
   * @private
   * @param {Uint8Array} wavBytes - Generated WAV file bytes
   */
  _handleAudioGenerated(wavBytes) {
    try {
      this._audioBlob = new Blob([wavBytes], { type: "audio/wav" });
      const url = URL.createObjectURL(this._audioBlob);

      // Reuse the audio element created during user interaction
      if (!this._audio) {
        this._audio = new Audio();

        this._audio.addEventListener("ended", () => {
          this._dispatchEvent("ended");
        });

        this._audio.addEventListener("error", (e) => {
          this._handleError(new Error("Audio playback error"));
        });
      }

      // Set the source and play
      this._audio.src = url;
      this._audio
        .play()
        .then(() => {
          this._dispatchEvent("playing");
        })
        .catch((error) => {
          this._handleError(error);
        });
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Handles errors during processing or playback
   * @private
   * @param {Error} error - The error that occurred
   */
  _handleError(error) {
    this._dispatchEvent("error", { error });
  }

  /**
   * Dispatches custom events
   * @private
   * @param {string} eventName - Name of the event
   * @param {Object} detail - Event detail data
   */
  _dispatchEvent(eventName, detail = {}) {
    if (typeof this[`on${eventName}`] === "function") {
      this[`on${eventName}`](detail);
    }
  }

  /**
   * Loads a SPSQ sequence from string or File object
   * @param {string|File} input - SPSQ sequence content or File object
   * @returns {Promise<void>}
   * @throws {Error} If input is invalid or worker is not ready
   * @example
   * // Load from string
   * await synapse.load('# Presets\nalpha\n  tone 250 isochronic 8');
   *
   * // Load from File object
   * const file = document.getElementById('fileInput').files[0];
   * await synapse.load(file);
   */
  async load(input) {
    // Wait for worker to be ready
    if (!this._workerReady) {
      await this._initPromise;
    }

    if (!input) {
      throw new Error("Input is required");
    }

    // Handle File object
    if (input instanceof File) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = (e) => {
          this._sequence = e.target.result;
          this._dispatchEvent("loaded");
          resolve();
        };

        reader.onerror = () => {
          reject(new Error("Failed to read file"));
        };

        reader.readAsText(input);
      });
    }

    // Handle string
    if (typeof input === "string") {
      this._sequence = input;
      this._dispatchEvent("loaded");
      return Promise.resolve();
    }

    throw new Error("Input must be a string or File object");
  }

  /**
   * Plays the loaded sequence
   * @returns {Promise<void>}
   * @throws {Error} If no sequence is loaded or worker is not ready
   * @example
   * await synapse.play();
   */
  async play() {
    if (!this._workerReady) {
      throw new Error("Worker is not ready. Please wait for initialization.");
    }

    if (!this._sequence) {
      throw new Error("No sequence loaded. Call load() first.");
    }

    // If audio exists and is paused, resume playback
    if (this._audio && this._audio.paused && this._audioBlob) {
      try {
        await this._audio.play();
        this._dispatchEvent("playing");
        return;
      } catch (error) {
        throw new Error("Failed to resume playback: " + error.message);
      }
    }

    // Create audio element immediately to satisfy mobile autoplay policies
    // This must happen during the user interaction (click event)
    if (this._audio) {
      this._audio.pause();
    }
    this._audio = new Audio();

    // Setup event listeners before generation
    this._audio.addEventListener("ended", () => {
      this._dispatchEvent("ended");
    });

    this._audio.addEventListener("error", (e) => {
      this._handleError(new Error("Audio playback error"));
    });

    // Generate new audio
    this._dispatchEvent("generating");

    const encoder = new TextEncoder();
    const spsqBytes = encoder.encode(this._sequence);

    this._worker.postMessage({
      type: "generate",
      spsqBytes: spsqBytes,
    });
  }

  /**
   * Pauses the currently playing sequence
   * @returns {void}
   * @throws {Error} If no audio is playing
   * @example
   * synapse.pause();
   */
  pause() {
    if (!this._audio) {
      throw new Error("No audio is playing");
    }

    if (!this._audio.paused) {
      this._audio.pause();
      this._dispatchEvent("paused");
    }
  }

  /**
   * Stops the currently playing sequence and resets playback position
   * @returns {void}
   * @example
   * synapse.stop();
   */
  stop() {
    if (this._audio) {
      this._audio.pause();
      this._audio.currentTime = 0;
      this._audio = null;
      this._audioBlob = null;
      this._dispatchEvent("stopped");
    }
  }

  /**
   * Gets the current playback position in seconds
   * @returns {number} Current time in seconds, or 0 if not playing
   * @example
   * const currentTime = synapse.getCurrentTime();
   */
  getCurrentTime() {
    return this._audio ? this._audio.currentTime : 0;
  }

  /**
   * Gets the total duration of the loaded audio in seconds
   * @returns {number} Duration in seconds, or 0 if no audio is loaded
   * @example
   * const duration = synapse.getDuration();
   */
  getDuration() {
    return this._audio ? this._audio.duration : 0;
  }

  /**
   * Gets the current playback state
   * @returns {string} One of: 'idle', 'generating', 'playing', 'paused', 'stopped'
   * @example
   * const state = synapse.getState();
   */
  getState() {
    if (!this._audio) {
      return "idle";
    }
    if (this._audio.src === "" || this._audio.readyState < 3) {
      return "generating";
    }
    if (this._audio.paused) {
      return this._audio.currentTime > 0 ? "paused" : "stopped";
    }
    return "playing";
  }

  /**
   * Checks if a sequence is currently loaded
   * @returns {boolean} True if a sequence is loaded
   * @example
   * if (synapse.isLoaded()) {
   *   await synapse.play();
   * }
   */
  isLoaded() {
    return this._sequence !== null;
  }

  /**
   * Checks if the worker is ready
   * @returns {boolean} True if worker is initialized and ready
   * @example
   * if (synapse.isReady()) {
   *   await synapse.load(sequence);
   * }
   */
  isReady() {
    return this._workerReady;
  }

  /**
   * Gets the generated audio blob
   * @returns {Blob|null} The audio blob or null if not generated
   * @example
   * const blob = synapse.getAudioBlob();
   * if (blob) {
   *   const url = URL.createObjectURL(blob);
   *   // Use the URL as needed
   * }
   */
  getAudioBlob() {
    return this._audioBlob;
  }

  /**
   * Gets the SynapSeq version
   * @returns {string} The version string
   * @example
   * const version = synapse.getVersion();
   * console.log('SynapSeq Version:', version);
   */
  async getVersion() {
    if (!this._workerReady) {
      await this._initPromise;
    }
    return this._version;
  }

  /**
   * Gets the build date of the SynapSeq WASM
   * @returns {string} The build date string
   * @example
   * const buildDate = synapse.getBuildDate();
   * console.log('SynapSeq Build Date:', buildDate);
   */
  async getBuildDate() {
    if (!this._workerReady) {
      await this._initPromise;
    }
    return this._buildDate;
  }

  /**
   * Gets the hash of the SynapSeq WASM build
   * @returns {string} The hash string
   * @example
   * const hash = synapse.getHash();
   * console.log('SynapSeq Hash:', hash);
   */
  async getHash() {
    if (!this._workerReady) {
      await this._initPromise;
    }
    return this._hash;
  }

  /**
   * Cleans up resources and terminates the worker
   * @example
   * synapse.destroy();
   */
  destroy() {
    this.stop();

    if (this._worker) {
      this._worker.terminate();
      this._worker = null;
    }

    this._workerReady = false;
    this._sequence = null;
    this._initPromise = null;
  }

  /**
   * Event handler called when sequence is loaded
   * @type {Function|null}
   * @example
   * synapse.onloaded = () => console.log('Sequence loaded');
   */
  onloaded = null;

  /**
   * Event handler called when audio generation starts
   * @type {Function|null}
   * @example
   * synapse.ongenerating = () => console.log('Generating audio...');
   */
  ongenerating = null;

  /**
   * Event handler called when playback starts
   * @type {Function|null}
   * @example
   * synapse.onplaying = () => console.log('Now playing');
   */
  onplaying = null;

  /**
   * Event handler called when playback is paused
   * @type {Function|null}
   * @example
   * synapse.onpaused = () => console.log('Paused');
   */
  onpaused = null;

  /**
   * Event handler called when playback is stopped
   * @type {Function|null}
   * @example
   * synapse.onstopped = () => console.log('Stopped');
   */
  onstopped = null;

  /**
   * Event handler called when playback ends naturally
   * @type {Function|null}
   * @example
   * synapse.onended = () => console.log('Playback finished');
   */
  onended = null;

  /**
   * Event handler called when an error occurs
   * @type {Function|null}
   * @example
   * synapse.onerror = (detail) => console.error('Error:', detail.error);
   */
  onerror = null;
}

// Export for use in modules or global scope
if (typeof module !== "undefined" && module.exports) {
  module.exports = SynapSeq;
}
