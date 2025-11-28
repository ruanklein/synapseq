// Initialize Lucide icons first
lucide.createIcons();

let synapseq = null;
let progressInterval = null;
let lastSequence = localStorage.getItem("last-sequence") || "";

// Wait for DOM to be fully loaded
document.addEventListener("DOMContentLoaded", function () {
  // Theme toggle
  const themeToggle = document.getElementById("themeToggle");
  const body = document.body;

  if (!themeToggle) {
    console.error("Theme toggle button not found");
    return;
  }

  // Load saved theme or default to dark
  const savedTheme = localStorage.getItem("synapseq-theme") || "dark";
  if (savedTheme === "light") {
    body.classList.remove("dark");
    updateThemeIcon("moon");
  } else {
    body.classList.add("dark");
    updateThemeIcon("lightbulb");
  }

  // Check for hub sequence ID in URL
  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  let hubSequenceID = urlParams.get("id");

  // Load sequence from hub if ID is present
  if (hubSequenceID) {
    fetch("https://synapseq-hub.ruan.sh/manifest.json")
      .then((response) => {
        if (!response.ok) {
          console.error("Failed to fetch manifest:", response.statusText);
          return;
        }

        response.json().then((data) => {
          const sequence = data.entries.find((e) => e.id === hubSequenceID);
          if (!sequence) {
            console.error("Sequence ID not found in manifest");
            return;
          }

          fetch(sequence.download_url).then((res) => {
            if (!res.ok) {
              console.error("Failed to fetch sequence:", res.statusText);
              return;
            }

            res.text().then((spsq) => {
              let newSpsq = spsq;
              sequence.dependencies.forEach((dep) => {
                if (dep.type === "background") {
                  newSpsq = newSpsq.replace(
                    new RegExp(`^@background ${dep.name}.wav\\b`, "m"),
                    `@background ${dep.download_url}`
                  );
                }

                if (dep.type === "presetlist") {
                  newSpsq = newSpsq.replace(
                    new RegExp(`^@presetlist ${dep.name}.spsq\\b`, "m"),
                    `@presetlist ${dep.download_url}`
                  );
                }
              });

              document.getElementById("spsqEditor").value = newSpsq;
              document.getElementById("editorContainer").focus();
              updateLineNumbers();
              //updateSyntaxHighlight();

              lastSequence = newSpsq;
              localStorage.setItem("last-sequence", lastSequence);

              // Show success message and spotlight with sequence name
              showHubLoadSuccess(sequence.name);
            });
          });
        });
      })
      .catch((error) => {
        console.error("Failed to fetch manifest:", error);
      });
  }
  if (hubSequenceID === null && lastSequence.length > 0) {
    document.getElementById("spsqEditor").value = lastSequence;
    updateLineNumbers();
    //updateSyntaxHighlight();
  }

  function updateThemeIcon(icon) {
    const iconElement = themeToggle.querySelector("i");
    if (iconElement) {
      iconElement.setAttribute("data-lucide", icon);
      // Force icon update by removing and re-adding
      iconElement.outerHTML = `<i data-lucide="${icon}" style="width: 1.25rem; height: 1.25rem"></i>`;
      lucide.createIcons();
    }
  }

  themeToggle.addEventListener("click", () => {
    const isDark = body.classList.contains("dark");

    if (isDark) {
      body.classList.remove("dark");
      localStorage.setItem("synapseq-theme", "light");
      updateThemeIcon("moon");
    } else {
      body.classList.add("dark");
      localStorage.setItem("synapseq-theme", "dark");
      updateThemeIcon("lightbulb");
    }
  });

  // Mobile menu toggle
  const hamburger = document.getElementById("hamburger");
  const navLinks = document.querySelector(".nav-links");

  if (hamburger) {
    hamburger.addEventListener("click", () => {
      hamburger.classList.toggle("active");
      navLinks.classList.toggle("active");
    });

    // Close menu when clicking on a link
    navLinks.querySelectorAll("a").forEach((link) => {
      link.addEventListener("click", () => {
        hamburger.classList.remove("active");
        navLinks.classList.remove("active");
      });
    });
  }
});

// Line numbers functionality
function updateLineNumbers() {
  const textarea = document.getElementById("spsqEditor");
  const lineNumbers = document.getElementById("lineNumbers");
  const lines = textarea.value.split("\n").length;

  const numbers = [];
  for (let i = 1; i <= lines; i++) {
    numbers.push(i);
  }
  lineNumbers.textContent = numbers.join("\n");
}

// Syntax highlighting
function highlightSyntax(code) {
  // Split into lines for better processing
  const lines = code.split("\n");
  const highlightedLines = lines.map((line) => {
    let highlighted = line;

    // Comments (must be first to avoid double-highlighting)
    if (/^\s*#/.test(line)) {
      return `<span class="syntax-comment">${line}</span>`;
    }

    // Directives and their arguments
    highlighted = highlighted.replace(
      /(@\w+)(\s+.+)?$/g,
      (match, directive, rest) => {
        let result = `<span class="syntax-directive">${directive}</span>`;
        if (rest) {
          result += `<span class="syntax-text">${rest}</span>`;
        }
        return result;
      }
    );

    // Timeline entries (time followed by preset name or silence)
    highlighted = highlighted.replace(
      /\b(\d{2}:\d{2}:\d{2})\s+(\w+)/g,
      (match, time, preset) => {
        return `<span class="syntax-time">${time}</span> <span class="syntax-timeline-preset">${preset}</span>`;
      }
    );

    // Keywords (tone, noise, background, etc)
    highlighted = highlighted.replace(
      /\b(tone|noise|background|waveform|amplitude|rate|intensity)\b/g,
      (match) => `<span class="syntax-keyword">${match}</span>`
    );

    // Types (binaural, monaural, etc)
    highlighted = highlighted.replace(
      /\b(binaural|monaural|isochronic|pulse|spin|white|pink|brown)\b/g,
      (match) => `<span class="syntax-type">${match}</span>`
    );

    // Waveforms (sine, square, etc)
    highlighted = highlighted.replace(
      /\b(sine|square|triangle|sawtooth)\b/g,
      (match) => `<span class="syntax-waveform">${match}</span>`
    );

    // Numbers (not part of time)
    highlighted = highlighted.replace(
      /\b(\d+(?:\.\d+)?)\b/g,
      (match) => `<span class="syntax-number">${match}</span>`
    );

    // Preset names (lines starting with word character, not already highlighted)
    if (/^[a-zA-Z]/.test(line) && !highlighted.includes("<span")) {
      highlighted = highlighted.replace(
        /^([a-zA-Z][a-zA-Z0-9_-]*)\b/,
        (match) => `<span class="syntax-preset">${match}</span>`
      );
    }

    return highlighted;
  });

  return highlightedLines.join("\n");
}

function updateSyntaxHighlight() {
  const textarea = document.getElementById("spsqEditor");
  const highlight = document.getElementById("syntaxHighlight");
  const code = textarea.value;

  // Escape HTML and apply syntax highlighting
  const escaped = code
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

  highlight.innerHTML = highlightSyntax(escaped);
}

// Sync scroll between line numbers and textarea
document.getElementById("spsqEditor").addEventListener("scroll", (e) => {
  document.getElementById("lineNumbers").scrollTop = e.target.scrollTop;
  document.getElementById("syntaxHighlight").scrollTop = e.target.scrollTop;
});

// Update line numbers on input
document.getElementById("spsqEditor").addEventListener("input", () => {
  updateLineNumbers();
  //updateSyntaxHighlight();
});

// Initialize line numbers
updateLineNumbers();
//updateSyntaxHighlight();

// Error handling
function showError(message) {
  const errorContainer = document.getElementById("errorContainer");
  const errorMessage = document.getElementById("errorMessage");
  errorMessage.textContent = message;
  errorContainer.classList.add("show");
  lucide.createIcons();

  // Scroll to error
  errorContainer.scrollIntoView({ behavior: "smooth", block: "nearest" });
}

function hideError() {
  const errorContainer = document.getElementById("errorContainer");
  errorContainer.classList.remove("show");
}

// Close error button handler
document.addEventListener("DOMContentLoaded", () => {
  const errorClose = document.getElementById("errorClose");
  if (errorClose) {
    errorClose.addEventListener("click", hideError);
  }
});

// Status message
function setStatus(message) {
  document.getElementById("statusMessage").textContent = message;
}

// Loading overlay
function showLoading() {
  const overlay = document.getElementById("loadingOverlay");
  overlay.classList.add("show");
  lucide.createIcons();
}

function hideLoading() {
  const overlay = document.getElementById("loadingOverlay");
  overlay.classList.remove("show");
}

// Time formatting
function formatTime(seconds) {
  if (!isFinite(seconds)) return "00:00";
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins.toString().padStart(2, "0")}:${secs
    .toString()
    .padStart(2, "0")}`;
}

// Progress tracking
function updateProgress() {
  const currentTime = synapseq.getCurrentTime();
  const duration = synapseq.getDuration();

  if (duration > 0) {
    const progress = (currentTime / duration) * 100;
    document.getElementById("progressBar").style.width = progress + "%";
    document.getElementById("currentTime").textContent =
      formatTime(currentTime);
    document.getElementById("totalTime").textContent = formatTime(duration);
  }
}

function startProgressTracking() {
  if (progressInterval) clearInterval(progressInterval);
  progressInterval = setInterval(updateProgress, 100);
}

function stopProgressTracking() {
  if (progressInterval) {
    clearInterval(progressInterval);
    progressInterval = null;
  }
}

// Initialize SynapSeq
async function initSynapSeq() {
  try {
    synapseq = new SynapSeq({
      wasmPath: "wasm/synapseq.wasm",
      wasmExecPath: "wasm/wasm_exec.js",
    });

    // Event handlers
    synapseq.onloaded = () => {
      console.log("Sequence loaded");
    };

    synapseq.ongenerating = () => {
      setStatus("Generating audio...");
      showLoading();
      document.getElementById("playBtn").disabled = true;
      document.getElementById("uploadBtn").disabled = true;
      document.getElementById("spsqEditor").disabled = true;
    };

    synapseq.onplaying = () => {
      hideLoading();
      setStatus("Playing...");
      document.getElementById("playBtn").disabled = true;
      document.getElementById("pauseBtn").disabled = false;
      document.getElementById("stopBtn").disabled = false;
      document.getElementById("uploadBtn").disabled = true;
      document.getElementById("spsqEditor").disabled = true;
      startProgressTracking();

      if ("mediaSession" in navigator) {
        navigator.mediaSession.metadata = new MediaMetadata({
          title: "SynapSeq Sequence",
          artist: "Synapse-Sequenced Generator",
          album: "Brainwave Session",
          artwork: [
            {
              src: "/assets/icon-512.png",
              sizes: "512x512",
              type: "image/png",
            },
          ],
        });
        navigator.mediaSession.setActionHandler("play", () => synapseq.play());
        navigator.mediaSession.setActionHandler("pause", () =>
          synapseq.pause()
        );
        navigator.mediaSession.setActionHandler("stop", () => synapseq.stop());
        navigator.mediaSession.playbackState = "playing";
      }
    };

    synapseq.onpaused = () => {
      setStatus("Paused");
      document.getElementById("playBtn").disabled = false;
      document.getElementById("pauseBtn").disabled = true;
      document.getElementById("uploadBtn").disabled = true;
      stopProgressTracking();

      if ("mediaSession" in navigator) {
        navigator.mediaSession.playbackState = "paused";
      }
    };

    synapseq.onstopped = () => {
      setStatus("Stopped");
      document.getElementById("playBtn").disabled = false;
      document.getElementById("pauseBtn").disabled = true;
      document.getElementById("stopBtn").disabled = true;
      document.getElementById("uploadBtn").disabled = false;
      document.getElementById("spsqEditor").disabled = false;
      stopProgressTracking();
      document.getElementById("progressBar").style.width = "0%";
      document.getElementById("currentTime").textContent = "00:00";
      document.getElementById("totalTime").textContent = "00:00";
      if ("mediaSession" in navigator) {
        navigator.mediaSession.playbackState = "none";
      }
    };

    synapseq.onended = () => {
      setStatus("Playback finished");
      synapseq.stop();
      if ("mediaSession" in navigator) {
        navigator.mediaSession.playbackState = "none";
      }
    };

    synapseq.onerror = (detail) => {
      hideLoading();
      console.error("Error:", detail.error);
      showError(detail.error.message || detail.error);
      setStatus("Error");
      document.getElementById("playBtn").disabled = false;
      document.getElementById("uploadBtn").disabled = false;
      document.getElementById("spsqEditor").disabled = false;
    };

    // Wait for worker to be ready
    const checkReady = setInterval(() => {
      if (synapseq.isReady()) {
        clearInterval(checkReady);
        setStatus("Ready");
        document.getElementById("playBtn").disabled = false;
      }
    }, 100);
  } catch (error) {
    console.error("Failed to initialize:", error);
    setStatus("Failed to initialize WASM");
    showError("Failed to load SynapSeq. Please refresh the page.");
  }
}

// Button handlers
document.getElementById("playBtn").addEventListener("click", async () => {
  if (!synapseq || !synapseq.isReady()) {
    setStatus("WASM not initialized yet");
    return;
  }

  hideError();

  try {
    const spsq = document.getElementById("spsqEditor").value;

    if (!spsq.trim()) {
      showError("Please enter SPSQ code");
      setStatus("Ready");
      return;
    }

    // Always reload the sequence to get the latest content
    if (spsq != lastSequence && synapseq.isLoaded()) {
      synapseq.stop();
    }

    await synapseq.load(spsq);

    lastSequence = spsq;
    localStorage.setItem("last-sequence", lastSequence);

    // Play the sequence
    await synapseq.play();
  } catch (error) {
    showError(error.message);
    setStatus("Error");
  }
});

document.getElementById("pauseBtn").addEventListener("click", () => {
  if (synapseq) {
    try {
      synapseq.pause();
    } catch (error) {
      showError(error.message);
    }
  }
});

document.getElementById("stopBtn").addEventListener("click", () => {
  if (synapseq) {
    synapseq.stop();
  }
});

// File menu dropdown handler
const fileMenuBtn = document.getElementById("fileMenuBtn");
const fileMenu = document.getElementById("fileMenu");
let menuTimeout;

fileMenuBtn.addEventListener("click", (e) => {
  e.stopPropagation();
  const isOpen = fileMenu.classList.contains("show");

  if (isOpen) {
    closeFileMenu();
  } else {
    openFileMenu();
  }
});

function openFileMenu() {
  fileMenu.classList.add("show");
  fileMenuBtn.classList.add("menu-open");
  lucide.createIcons();
}

function closeFileMenu() {
  fileMenu.classList.remove("show");
  fileMenuBtn.classList.remove("menu-open");
}

// Close menu when clicking outside
document.addEventListener("click", (e) => {
  if (!fileMenu.contains(e.target) && !fileMenuBtn.contains(e.target)) {
    closeFileMenu();
  }
});

// Upload menu item
document.getElementById("uploadMenuItem").addEventListener("click", () => {
  closeFileMenu();
  document.getElementById("fileInput").click();
});

// Save menu item
document.getElementById("saveMenuItem").addEventListener("click", () => {
  closeFileMenu();
  saveSequenceToFile();
});

// Save sequence to file
function saveSequenceToFile() {
  const content = document.getElementById("spsqEditor").value;

  if (!content.trim()) {
    showError("Cannot save empty sequence");
    return;
  }

  // Generate timestamp-based filename
  const now = new Date();
  const timestamp = now
    .toISOString()
    .replace(/:/g, "-")
    .replace(/\./g, "-")
    .substring(0, 19);
  const filename = `synapseq-${timestamp}.spsq`;

  // Create blob and download
  const blob = new Blob([content], { type: "text/plain" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);

  setStatus("Sequence saved: " + filename);
}

// File upload handler
document
  .getElementById("fileInput")
  .addEventListener("change", async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    // Validate file extension
    if (!file.name.toLowerCase().endsWith(".spsq")) {
      showError("Invalid file type. Please upload a .spsq file");
      event.target.value = "";
      return;
    }

    hideError();

    try {
      // Read file content
      const reader = new FileReader();
      reader.onload = (e) => {
        document.getElementById("spsqEditor").value = e.target.result;
        updateLineNumbers();
        //updateSyntaxHighlight();
        setStatus("File loaded: " + file.name);
      };
      reader.onerror = () => {
        showError("Failed to read file");
      };
      reader.readAsText(file);
    } catch (error) {
      showError("Failed to load file: " + error.message);
    }

    // Reset input
    event.target.value = "";
  });

// Show hub load success
function showHubLoadSuccess(sequenceName) {
  // Wait for WASM to be ready before showing spotlight
  const checkReady = setInterval(() => {
    if (synapseq && synapseq.isReady()) {
      clearInterval(checkReady);

      // Show spotlight after a short delay
      setTimeout(() => {
        const spotlightOverlay = document.getElementById("spotlightOverlay");
        const spotlightClose = document.getElementById("spotlightClose");
        const spotlightPlayBtn = document.getElementById("spotlightPlayBtn");
        const spotlightSequenceName = document.getElementById(
          "spotlightSequenceName"
        );

        // Set sequence name
        if (sequenceName && spotlightSequenceName) {
          spotlightSequenceName.textContent = sequenceName;
        }

        spotlightOverlay.classList.add("show");
        lucide.createIcons();

        // Remove spotlight function
        const removeSpotlight = () => {
          spotlightOverlay.classList.remove("show");
        };

        // Close button handler
        spotlightClose.addEventListener("click", removeSpotlight);

        // Spotlight play button handler
        spotlightPlayBtn.addEventListener("click", async () => {
          removeSpotlight();

          // Trigger the main play button click
          const mainPlayBtn = document.getElementById("playBtn");
          mainPlayBtn.click();
        });

        // Close when clicking on backdrop
        const backdrop = spotlightOverlay.querySelector(".spotlight-backdrop");
        backdrop.addEventListener("click", removeSpotlight);
      }, 500);
    }
  }, 100);
} // Initialize on load
setStatus("Loading WASM...");
initSynapSeq();
