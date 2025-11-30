// Initialize Lucide icons first
lucide.createIcons();

// Initialize dayjs with relativeTime plugin
dayjs.extend(dayjs_plugin_relativeTime);

let synapseq = null;
let progressInterval = null;
let lastSequenceData = null;
let savedSequences = JSON.parse(
  localStorage.getItem("saved-sequences") || "[]"
);
let saveDebounceTimer = null;
let saveTimeUpdateInterval = null;
let errorLine = null; // Store the line number with syntax error
let activePresetLines = null; // Store the line range of active preset
let isPlaying = false; // Track if sequence is playing

// Load last sequence from localStorage
try {
  const stored = localStorage.getItem("last-sequence");
  if (stored) {
    const parsed = JSON.parse(stored);
    if (parsed && typeof parsed === "object" && parsed.sequence) {
      lastSequenceData = parsed;
    } else {
      // Old format - migrate to new format
      lastSequenceData = {
        sequence: typeof parsed === "string" ? parsed : stored,
        date: new Date().toISOString(),
      };
    }
  }
} catch (e) {
  // Old format - migrate to new format
  const stored = localStorage.getItem("last-sequence") || "";
  if (stored) {
    lastSequenceData = {
      sequence: stored,
      date: new Date().toISOString(),
    };
  }
}

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
  const savedTheme = localStorage.getItem("synapseq-theme") || "light";
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
              updateSyntaxHighlight();

              const now = new Date().toISOString();
              lastSequenceData = {
                sequence: newSpsq,
                date: now,
              };
              localStorage.setItem(
                "last-sequence",
                JSON.stringify(lastSequenceData)
              );

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
  if (
    hubSequenceID === null &&
    lastSequenceData &&
    lastSequenceData.sequence.length > 0
  ) {
    document.getElementById("spsqEditor").value = lastSequenceData.sequence;
    updateLineNumbers();
    updateSyntaxHighlight();
  }

  // Start time updater if there's a saved sequence
  if (lastSequenceData && lastSequenceData.date) {
    startSaveTimeUpdater();
  }

  function updateThemeIcon(icon) {
    // Remove the old icon
    const iconElement = themeToggle.querySelector("i, svg");
    if (iconElement) {
      iconElement.remove();
    }

    // Create new icon element
    const newIcon = document.createElement("i");
    newIcon.setAttribute("data-lucide", icon);
    newIcon.setAttribute("style", "width: 1.25rem; height: 1.25rem");
    themeToggle.appendChild(newIcon);

    // Render the icon
    lucide.createIcons();
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
  const highlightedLines = lines.map((line, lineIndex) => {
    let highlighted = line;
    const lineNumber = lineIndex + 1;
    const hasError = errorLine === lineNumber;
    const isActive =
      activePresetLines &&
      lineNumber >= activePresetLines.start &&
      lineNumber <= activePresetLines.end;

    // Comments (must be first to avoid double-highlighting)
    if (/^\s*#/.test(line)) {
      const result = `<span class="syntax-comment">${line}</span>`;
      if (hasError) return `<span class="error-line">${result}</span>`;
      if (isActive) return `<span class="active-preset-line">${result}</span>`;
      return result;
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

    // Preset names (lines starting with word character, must be before keywords)
    if (/^[a-zA-Z]/.test(line) && !highlighted.includes("<span")) {
      highlighted = highlighted.replace(
        /^([a-zA-Z][a-zA-Z0-9_-]*)(\s+(?:as|from)\s+(?:template|[a-zA-Z][a-zA-Z0-9_-]*))?(.*)$/,
        (match, preset, templatePart, rest) => {
          let result = `<span class="syntax-preset">${preset}</span>`;
          if (templatePart) {
            // Parse the template part to preserve all spaces
            const templateMatch = templatePart.match(
              /(\s+)(as|from)(\s+)(template|[a-zA-Z][a-zA-Z0-9_-]*)/
            );
            if (templateMatch) {
              const [, space1, keyword, space2, target] = templateMatch;
              result +=
                space1 +
                `<span class="syntax-template">${keyword}</span>` +
                space2 +
                `<span class="syntax-${
                  keyword === "as" ? "template" : "preset"
                }">${target}</span>`;
            }
          }
          result += rest;
          return result;
        }
      );
    }

    // Timeline entries (time followed by preset name or silence)
    highlighted = highlighted.replace(
      /\b(\d{2}:\d{2}:\d{2})(\s+)([a-zA-Z][a-zA-Z0-9_-]*)(\s+(?:steady|ease-out|ease-in|smooth))?/g,
      (match, time, space1, preset, rampPart) => {
        let result = `<span class="syntax-time">${time}</span>${space1}<span class="syntax-timeline-preset">${preset}</span>`;
        if (rampPart) {
          // Parse the ramp part to preserve all spaces
          const rampMatch = rampPart.match(
            /(\s+)(steady|ease-out|ease-in|smooth)/
          );
          if (rampMatch) {
            const [, space2, ramp] = rampMatch;
            result += space2 + `<span class="syntax-ramp">${ramp}</span>`;
          }
        }
        return result;
      }
    );

    // Keywords (tone, noise, background, etc) - must be preceded by exactly 2 spaces
    highlighted = highlighted.replace(
      /^  (tone|noise|background|track|waveform|amplitude|rate|intensity)\b/g,
      (match, keyword) => `  <span class="syntax-keyword">${keyword}</span>`
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

    // Wrap with active preset or error styling
    if (hasError) {
      return `<span class="error-line">${highlighted}</span>`;
    }
    if (isActive) {
      return `<span class="active-preset-line">${highlighted}</span>`;
    }
    return highlighted;
  });

  return highlightedLines.join("\n");
}

// Find which preset is active based on current time
function findActivePreset(code, currentTimeSeconds) {
  const lines = code.split("\n");
  const timeline = [];

  // Parse timeline entries
  lines.forEach((line, index) => {
    const timeMatch = line.match(
      /^\s*(\d{2}):(\d{2}):(\d{2})\s+([a-zA-Z][a-zA-Z0-9_-]*)/
    );
    if (timeMatch) {
      const hours = parseInt(timeMatch[1], 10);
      const minutes = parseInt(timeMatch[2], 10);
      const seconds = parseInt(timeMatch[3], 10);
      const presetName = timeMatch[4];
      const timeInSeconds = hours * 3600 + minutes * 60 + seconds;

      timeline.push({
        time: timeInSeconds,
        preset: presetName,
        lineIndex: index,
      });
    }
  });

  // Find active preset (last timeline entry before or at current time)
  let activePreset = null;
  for (let i = timeline.length - 1; i >= 0; i--) {
    if (timeline[i].time <= currentTimeSeconds) {
      activePreset = timeline[i].preset;
      break;
    }
  }

  if (!activePreset) return null;

  // Find preset definition lines
  const presetStartLine = lines.findIndex((line) => {
    const match = line.match(/^([a-zA-Z][a-zA-Z0-9_-]*)/);
    return match && match[1] === activePreset;
  });

  if (presetStartLine === -1) return null; // Preset not in editor (from @presetlist)

  // Find end of preset (next preset definition or empty line or directive)
  let presetEndLine = presetStartLine;
  for (let i = presetStartLine + 1; i < lines.length; i++) {
    const line = lines[i];
    // Stop at next preset, directive, timeline, or significant empty section
    if (
      /^[a-zA-Z]/.test(line) ||
      /^@/.test(line) ||
      /^\d{2}:\d{2}:\d{2}/.test(line)
    ) {
      break;
    }
    // Include lines starting with 2 spaces (tracks) or comments
    if (/^  /.test(line) || /^\s*#/.test(line) || line.trim() === "") {
      presetEndLine = i;
    } else {
      break;
    }
  }

  return {
    start: presetStartLine + 1, // 1-indexed
    end: presetEndLine + 1,
  };
}

function updateSyntaxHighlight() {
  const textarea = document.getElementById("spsqEditor");
  const highlight = document.getElementById("syntaxHighlight");
  const code = textarea.value;

  // Apply syntax highlighting directly (no need to escape for SPSQ syntax)
  highlight.innerHTML = highlightSyntax(code);
}

// Sync scroll between line numbers and textarea
document.getElementById("spsqEditor").addEventListener("scroll", (e) => {
  document.getElementById("lineNumbers").scrollTop = e.target.scrollTop;
  const highlight = document.getElementById("syntaxHighlight");
  highlight.scrollTop = e.target.scrollTop;
  highlight.scrollLeft = e.target.scrollLeft;
});

// Update line numbers on input
document.getElementById("spsqEditor").addEventListener("input", () => {
  updateLineNumbers();
  updateSyntaxHighlight();

  // Clear error highlight when user starts typing
  if (errorLine !== null) {
    errorLine = null;
    hideAlert("error");
  }

  // Debounced save to lastSequence
  saveCurrentSequenceDebounced();
});

// Debounced save function
function saveCurrentSequenceDebounced() {
  // Clear existing timer
  if (saveDebounceTimer) {
    clearTimeout(saveDebounceTimer);
  }

  // Set new timer
  saveDebounceTimer = setTimeout(() => {
    const content = document.getElementById("spsqEditor").value;
    const now = new Date().toISOString();
    lastSequenceData = {
      sequence: content,
      date: now,
    };
    localStorage.setItem("last-sequence", JSON.stringify(lastSequenceData));

    // Update alert with current time
    updateSaveTimeDisplay();
    showAlert(
      "success",
      "Saved",
      `Sequence auto-saved ${dayjs(lastSequenceData.date).fromNow()}`
    );
  }, 1000);
}

// Show alert - uses separate containers for error and success
function showAlert(type, title, message, help = null) {
  if (type === "error") {
    const container = document.getElementById("errorAlertContainer");
    const alertTitle = document.getElementById("errorAlertTitle");
    const alertSubtitle = document.getElementById("errorAlertSubtitle");
    const alertMessage = document.getElementById("errorAlertMessage");
    const alertHelp = document.getElementById("errorAlertHelp");

    alertTitle.textContent = title;
    alertSubtitle.textContent = message;
    alertMessage.textContent = "";

    if (help) {
      alertHelp.style.display = "flex";
    } else {
      alertHelp.style.display = "none";
    }

    container.className = "alert-container show error";
    lucide.createIcons();
  } else if (type === "success") {
    const container = document.getElementById("successAlertContainer");
    const alertTitle = document.getElementById("successAlertTitle");
    const alertSubtitle = document.getElementById("successAlertSubtitle");

    alertTitle.textContent = title;
    alertSubtitle.textContent = message;

    container.className = "alert-container show success";
    lucide.createIcons();
  }

  // Force sync overlay scroll after layout shift
  setTimeout(() => {
    const textarea = document.getElementById("spsqEditor");
    const highlight = document.getElementById("syntaxHighlight");
    if (textarea && highlight) {
      highlight.scrollTop = textarea.scrollTop;
      highlight.scrollLeft = textarea.scrollLeft;
    }
  }, 100);
}

function hideAlert(type = null) {
  if (type === "error" || type === null) {
    const container = document.getElementById("errorAlertContainer");
    if (container) container.classList.remove("show");
  }

  if (type === "success" || type === null) {
    const container = document.getElementById("successAlertContainer");
    if (container) container.classList.remove("show");
  }

  // Force sync overlay scroll after layout shift
  setTimeout(() => {
    const textarea = document.getElementById("spsqEditor");
    const highlight = document.getElementById("syntaxHighlight");
    if (textarea && highlight) {
      highlight.scrollTop = textarea.scrollTop;
      highlight.scrollLeft = textarea.scrollLeft;
    }
  }, 100);
}

// Update save time display
function updateSaveTimeDisplay() {
  if (!lastSequenceData || !lastSequenceData.date) return;

  const alertSubtitle = document.getElementById("successAlertSubtitle");
  const container = document.getElementById("successAlertContainer");

  if (alertSubtitle && container && container.classList.contains("show")) {
    alertSubtitle.textContent = `Sequence auto-saved ${dayjs(
      lastSequenceData.date
    ).fromNow()}`;
  }
}

// Start continuous time updater (runs always in background)
function startSaveTimeUpdater() {
  stopSaveTimeUpdater();

  saveTimeUpdateInterval = setInterval(() => {
    updateSaveTimeDisplay();
  }, 1000);
}

function stopSaveTimeUpdater() {
  if (saveTimeUpdateInterval) {
    clearInterval(saveTimeUpdateInterval);
    saveTimeUpdateInterval = null;
  }
}

// Save sequence to history
function saveSequenceToHistory(content) {
  if (!content || !content.trim()) return;

  const timestamp = new Date()
    .toISOString()
    .replace(/:/g, "-")
    .replace(/\./g, "-")
    .substring(0, 19);
  const name = `synapseq-${timestamp}`;

  const sequence = {
    name: name,
    content: content,
    timestamp: Date.now(),
  };

  // Remove if already exists (by content)
  savedSequences = savedSequences.filter((s) => s.content !== content);

  // Add to beginning
  savedSequences.unshift(sequence);

  // Keep only last 10
  if (savedSequences.length > 10) {
    savedSequences = savedSequences.slice(0, 10);
  }

  // Save to localStorage
  localStorage.setItem("saved-sequences", JSON.stringify(savedSequences));

  // Update UI
  renderSequenceHistory();
}

// Render sequence history
function renderSequenceHistory() {
  const historyList = document.getElementById("sequenceHistoryList");
  const emptyState = document.getElementById("sequenceHistoryEmpty");

  if (!historyList || !emptyState) return;

  if (savedSequences.length === 0) {
    historyList.style.display = "none";
    emptyState.style.display = "block";
  } else {
    historyList.style.display = "block";
    emptyState.style.display = "none";

    historyList.innerHTML = savedSequences
      .map((seq, index) => {
        const date = new Date(seq.timestamp);
        const timeAgo = getTimeAgo(seq.timestamp);

        return `
        <div class="history-item" data-index="${index}">
          <div class="history-item-header">
            <span class="history-item-name">${seq.name}</span>
            <button class="history-item-delete" data-index="${index}" aria-label="Delete sequence">
              <i data-lucide="trash-2" style="width: 0.875rem; height: 0.875rem"></i>
            </button>
          </div>
          <div class="history-item-time">${timeAgo}</div>
        </div>
      `;
      })
      .join("");

    lucide.createIcons();

    // Add click handlers
    historyList.querySelectorAll(".history-item").forEach((item) => {
      item.addEventListener("click", (e) => {
        if (e.target.closest(".history-item-delete")) return;

        const index = parseInt(item.dataset.index);
        loadSequenceFromHistory(index);
      });
    });

    // Add delete handlers
    historyList.querySelectorAll(".history-item-delete").forEach((btn) => {
      btn.addEventListener("click", (e) => {
        e.stopPropagation();
        const index = parseInt(btn.dataset.index);
        deleteSequenceFromHistory(index);
      });
    });
  }
}

// Load sequence from history
function loadSequenceFromHistory(index) {
  if (index < 0 || index >= savedSequences.length) return;

  const seq = savedSequences[index];
  document.getElementById("spsqEditor").value = seq.content;
  updateLineNumbers();
  updateSyntaxHighlight();

  // Update lastSequenceData immediately
  const now = new Date().toISOString();
  lastSequenceData = {
    sequence: seq.content,
    date: now,
  };
  localStorage.setItem("last-sequence", JSON.stringify(lastSequenceData));

  showAlert("success", "Loaded", "Sequence loaded from history");
  setTimeout(() => hideAlert(), 1000);
}

// Delete sequence from history
function deleteSequenceFromHistory(index) {
  if (index < 0 || index >= savedSequences.length) return;

  savedSequences.splice(index, 1);
  localStorage.setItem("saved-sequences", JSON.stringify(savedSequences));
  renderSequenceHistory();
}

// Clear all history
function clearAllHistory() {
  if (savedSequences.length === 0) return;

  if (confirm("Are you sure you want to clear all saved sequences?")) {
    savedSequences = [];
    localStorage.setItem("saved-sequences", JSON.stringify(savedSequences));
    renderSequenceHistory();
  }
}

// Time ago helper
function getTimeAgo(timestamp) {
  const seconds = Math.floor((Date.now() - timestamp) / 1000);

  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;

  return new Date(timestamp).toLocaleDateString();
}

// Initialize line numbers
updateLineNumbers();
updateSyntaxHighlight();

// Initialize history UI
document.addEventListener("DOMContentLoaded", () => {
  renderSequenceHistory();

  // Clear history button
  const clearHistoryBtn = document.getElementById("clearHistoryBtn");
  if (clearHistoryBtn) {
    clearHistoryBtn.addEventListener("click", clearAllHistory);
  }
});

// Error handling
function showError(message) {
  // Try to extract line number from error message
  const lineMatch = message.match(/line\s+(\d+)/i);
  if (lineMatch) {
    errorLine = parseInt(lineMatch[1], 10);
    updateSyntaxHighlight(); // Re-render to show error line
  }

  showAlert(
    "error",
    "Syntax Error",
    "There's an issue with your sequence",
    "Check the documentation for syntax help"
  );
  const alertMessage = document.getElementById("errorAlertMessage");
  alertMessage.textContent = message;
}

function hideError() {
  errorLine = null; // Clear error line
  updateSyntaxHighlight(); // Re-render to remove error highlight
  hideAlert();
}

// Close alert button handlers
document.addEventListener("DOMContentLoaded", () => {
  const errorAlertClose = document.getElementById("errorAlertClose");
  if (errorAlertClose) {
    errorAlertClose.addEventListener("click", () => hideAlert("error"));
  }

  const successAlertClose = document.getElementById("successAlertClose");
  if (successAlertClose) {
    successAlertClose.addEventListener("click", () => hideAlert("success"));
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

    // Update active preset highlight
    if (isPlaying) {
      const code = document.getElementById("spsqEditor").value;
      const newActiveLines = findActivePreset(code, currentTime);

      // Only update if changed
      if (
        JSON.stringify(newActiveLines) !== JSON.stringify(activePresetLines)
      ) {
        activePresetLines = newActiveLines;
        updateSyntaxHighlight();

        // Auto-scroll to active preset
        if (activePresetLines) {
          scrollToActiveLine();
        }
      }
    }
  }
}

// Auto-scroll to active preset line
function scrollToActiveLine() {
  if (!activePresetLines) return;

  const textarea = document.getElementById("spsqEditor");
  const lineHeight = parseFloat(getComputedStyle(textarea).lineHeight);
  const containerHeight = textarea.clientHeight;

  // Scroll to center the active preset
  const targetScrollTop =
    (activePresetLines.start - 1) * lineHeight - containerHeight / 3;

  textarea.scrollTop = Math.max(0, targetScrollTop);

  // Sync overlay scroll
  const highlight = document.getElementById("syntaxHighlight");
  highlight.scrollTop = textarea.scrollTop;
  highlight.scrollLeft = textarea.scrollLeft;

  const lineNumbers = document.getElementById("lineNumbers");
  lineNumbers.scrollTop = textarea.scrollTop;
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
      // Sequence loaded successfully
    };

    synapseq.ongenerating = () => {
      setStatus("Generating audio...");
      showLoading();
      document.getElementById("playBtn").disabled = true;
      document.getElementById("fileMenuBtn").disabled = true;
      document.getElementById("spsqEditor").disabled = true;
    };

    synapseq.onplaying = () => {
      hideLoading();
      setStatus("Playing...");
      document.getElementById("playBtn").disabled = true;
      document.getElementById("pauseBtn").disabled = false;
      document.getElementById("stopBtn").disabled = false;
      document.getElementById("fileMenuBtn").disabled = true;
      document.getElementById("spsqEditor").disabled = true;
      isPlaying = true;
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
      document.getElementById("fileMenuBtn").disabled = true;
      isPlaying = false;
      activePresetLines = null;
      updateSyntaxHighlight();
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
      document.getElementById("fileMenuBtn").disabled = false;
      document.getElementById("spsqEditor").disabled = false;
      isPlaying = false;
      activePresetLines = null;
      updateSyntaxHighlight();
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
      document.getElementById("fileMenuBtn").disabled = false;
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
    if (
      lastSequenceData &&
      spsq != lastSequenceData.sequence &&
      synapseq.isLoaded()
    ) {
      synapseq.stop();
    }

    await synapseq.load(spsq);

    const now = new Date().toISOString();
    lastSequenceData = {
      sequence: spsq,
      date: now,
    };
    localStorage.setItem("last-sequence", JSON.stringify(lastSequenceData));

    // Save to history
    saveSequenceToHistory(spsq);

    // Play the sequence
    await synapseq.play();
  } catch (error) {
    hideLoading();
    showError(error.message);
    setStatus("Error");
    document.getElementById("playBtn").disabled = false;
    document.getElementById("fileMenuBtn").disabled = false;
    document.getElementById("spsqEditor").disabled = false;
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
        updateSyntaxHighlight();
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
