// Initialize Lucide icons
lucide.createIcons();

// Elements
const codeInput = document.getElementById("codeInput");
const lineNumbers = document.getElementById("lineNumbers");
const playBtn = document.getElementById("playBtn");
const stopBtn = document.getElementById("stopBtn");
const uploadBtn = document.getElementById("uploadBtn");
const shareBtn = document.getElementById("shareBtn");
const fileInput = document.getElementById("fileInput");
const statusIcon = document.getElementById("statusIcon");
const statusText = document.getElementById("statusText");
const progressFill = document.getElementById("progressFill");
const progressTime = document.getElementById("progressTime");
const errorMessage = document.getElementById("errorMessage");
const errorText = document.getElementById("errorText");

// SynapSeq instance
let synapseq = null;
let progressInterval = null;
let sequenceDuration = 0;
let wakeLock = null;

// Wake Lock functions
async function requestWakeLock() {
  try {
    if ("wakeLock" in navigator) {
      wakeLock = await navigator.wakeLock.request("screen");
      console.log("Wake Lock activated");
    }
  } catch (err) {
    console.warn("Wake Lock not available:", err);
  }
}

async function releaseWakeLock() {
  if (wakeLock) {
    try {
      await wakeLock.release();
      wakeLock = null;
      console.log("Wake Lock released");
    } catch (err) {
      console.warn("Failed to release Wake Lock:", err);
    }
  }
}

// Initialize SynapSeq
function initSynapSeq() {
  synapseq = new SynapSeq({
    wasmPath: "https://synapseq.org/lib/synapseq.wasm",
    wasmExecPath: "https://synapseq.org/lib/wasm_exec.js",
  });

  // Event handlers
  synapseq.onloaded = () => {
    setStatus("idle", "Ready to play");
    hideError();
  };

  synapseq.ongenerating = () => {
    setStatus("loading", "Generating audio...");
  };

  synapseq.onplaying = () => {
    setStatus("playing", "Playing");
    playBtn.disabled = true;
    stopBtn.disabled = false;
    uploadBtn.disabled = true;
    codeInput.readOnly = true;
    startProgressTracking();
    requestWakeLock();
  };

  synapseq.onstopped = () => {
    setStatus("stopped", "Stopped");
    playBtn.disabled = false;
    stopBtn.disabled = true;
    uploadBtn.disabled = false;
    codeInput.readOnly = false;
    stopProgressTracking();
    resetProgress();
    releaseWakeLock();
  };

  synapseq.onended = () => {
    setStatus("idle", "Finished");
    playBtn.disabled = false;
    stopBtn.disabled = true;
    uploadBtn.disabled = false;
    codeInput.readOnly = false;
    stopProgressTracking();
    resetProgress();
    releaseWakeLock();
  };

  synapseq.onerror = (detail) => {
    setStatus("idle", "Error");
    showError(detail.error.message);
    playBtn.disabled = false;
    stopBtn.disabled = true;
    uploadBtn.disabled = false;
    codeInput.readOnly = false;
    stopProgressTracking();
    releaseWakeLock();
  };
}

// Update line numbers
function updateLineNumbers() {
  const lines = codeInput.value.split("\n").length;
  lineNumbers.innerHTML = Array.from({ length: lines }, (_, i) => i + 1).join(
    "\n"
  );
}

// Sync scroll
codeInput.addEventListener("scroll", () => {
  lineNumbers.scrollTop = codeInput.scrollTop;
});

// Update line numbers on input
codeInput.addEventListener("input", updateLineNumbers);

// Set status
function setStatus(state, text) {
  statusIcon.className = `status-icon ${state}`;
  statusText.textContent = text;
}

// Show error
function showError(message) {
  errorText.textContent = message;
  errorMessage.classList.add("show");

  // Scroll to error on mobile
  setTimeout(() => {
    errorMessage.scrollIntoView({ behavior: "smooth", block: "nearest" });
  }, 100);
}

// Hide error
function hideError() {
  errorMessage.classList.remove("show");
}

// Format time
function formatTime(seconds) {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, "0")}`;
}

// Start progress tracking
function startProgressTracking() {
  progressInterval = setInterval(() => {
    const current = synapseq.getCurrentTime();

    // Update time display
    if (sequenceDuration > 0) {
      progressTime.textContent = `${formatTime(current)} / ${formatTime(
        sequenceDuration
      )}`;
    } else {
      progressTime.textContent = formatTime(current);
    }

    // Update progress bar
    if (sequenceDuration > 0) {
      const percent = (current / sequenceDuration) * 100;
      progressFill.style.width = `${Math.min(percent, 100)}%`;
    } else {
      // Fallback: just show it's playing
      progressFill.style.width = "50%";
    }
  }, 100);
}

// Stop progress tracking
function stopProgressTracking() {
  if (progressInterval) {
    clearInterval(progressInterval);
    progressInterval = null;
  }
}

// Reset progress
function resetProgress() {
  progressFill.style.width = "0%";
  if (sequenceDuration > 0) {
    progressTime.textContent = `0:00 / ${formatTime(sequenceDuration)}`;
  } else {
    progressTime.textContent = "0:00";
  }
}

// Play button handler
playBtn.addEventListener("click", async () => {
  const code = codeInput.value.trim();

  if (!code) {
    showError("Please write some sequence first");
    return;
  }

  hideError();
  playBtn.disabled = true;

  try {
    // Reset duration before parsing
    sequenceDuration = 0;

    // Parse sequence to find all timestamps at line start (hh:mm:ss format)
    const times = [...code.matchAll(/^(\d{2}):(\d{2}):(\d{2})/gm)];
    if (times.length > 0) {
      // Get the last timestamp as sequence duration
      const lastTime = times[times.length - 1];
      sequenceDuration =
        parseInt(lastTime[1]) * 3600 +
        parseInt(lastTime[2]) * 60 +
        parseInt(lastTime[3]);
    }

    await synapseq.load(code);
    await synapseq.play();
  } catch (error) {
    showError(error.message);
    playBtn.disabled = false;
  }
});

// Stop button handler
stopBtn.addEventListener("click", () => {
  synapseq.stop();
});

// Upload button handler
uploadBtn.addEventListener("click", () => {
  fileInput.click();
});

// File input handler
fileInput.addEventListener("change", async (e) => {
  const file = e.target.files[0];
  if (!file) return;

  try {
    const text = await file.text();
    codeInput.value = text;
    updateLineNumbers();
    hideError();

    // Reset duration so it recalculates on next play
    sequenceDuration = 0;
    resetProgress();
  } catch (error) {
    showError("Failed to read file: " + error.message);
  }

  // Reset input
  fileInput.value = "";
});

// Share button - Copy URL with base64 encoded sequence
shareBtn.addEventListener("click", () => {
  const sequence = codeInput.value.trim();

  if (!sequence) {
    showError("Cannot share an empty sequence");
    return;
  }

  try {
    // Encode sequence to base64 and make it URL-safe
    const encodedSequence = encodeURIComponent(btoa(sequence));

    // Create shareable URL
    const url = new URL(window.location.href);
    // Remove any existing hash
    url.hash = "";
    url.searchParams.set("sequence", encodedSequence);

    // Copy to clipboard
    navigator.clipboard
      .writeText(url.toString())
      .then(() => {
        // Visual feedback - change button text temporarily
        const originalHTML = shareBtn.innerHTML;
        shareBtn.innerHTML = '<i data-lucide="check"></i><span>Copied!</span>';
        lucide.createIcons();

        setTimeout(() => {
          shareBtn.innerHTML = originalHTML;
          lucide.createIcons();
        }, 2000);
      })
      .catch((error) => {
        showError("Failed to copy link to clipboard");
      });
  } catch (error) {
    showError("Failed to create shareable link");
  }
});

// Initialize
initSynapSeq();
updateLineNumbers();

// Load sequence from URL if present
function loadSequenceFromURL() {
  const urlParams = new URLSearchParams(window.location.search);
  const sequenceBase64 = urlParams.get("sequence");

  if (sequenceBase64) {
    try {
      // Decode from URL-safe base64
      const decodedSequence = atob(decodeURIComponent(sequenceBase64));
      codeInput.value = decodedSequence;
      updateLineNumbers();
      hideError();
    } catch (error) {
      showError("Failed to load sequence from URL: Invalid format");
    }
  }
}

// Load sequence on page load
loadSequenceFromURL();
