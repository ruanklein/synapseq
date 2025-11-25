/**
 * SynapSeq - Synapse-Sequenced Brainwave Generator
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
   */
  constructor() {
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
     * @type {Promise<void>|null}
     */
    this._initPromise = null;

    this._initializeWorker();
  }

  /**
   * Initializes the Web Worker for WASM processing
   * @private
   * @returns {Promise<void>}
   */
  _initializeWorker() {
    this._initPromise = new Promise((resolve, reject) => {
      try {
        this._worker = new Worker("worker.js");

        this._worker.onmessage = (e) => {
          const data = e.data;

          if (data.type === "ready") {
            this._workerReady = true;
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

      if (this._audio) {
        this._audio.pause();
        this._audio = null;
      }

      this._audio = new Audio(url);

      this._audio.addEventListener("ended", () => {
        this._dispatchEvent("ended");
      });

      this._audio.addEventListener("error", (e) => {
        this._handleError(new Error("Audio playback error"));
      });

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
   * Downloads the generated WAV file
   * @param {string} filename - Name for the downloaded file (default: 'synapseq.wav')
   * @throws {Error} If no audio has been generated
   * @example
   * synapse.download('my-sequence.wav');
   */
  download(filename = "synapseq.wav") {
    if (!this._audioBlob) {
      throw new Error("No audio has been generated yet");
    }

    const url = URL.createObjectURL(this._audioBlob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
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
