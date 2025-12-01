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
let currentSuggestion = null; // Store current autocomplete suggestion

// Detect if device is mobile - prioritize user agent for real mobile devices
function checkIsMobile() {
  // Check mobile user agent first (most reliable for real devices)
  const mobileUA =
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      navigator.userAgent
    );

  // If it's a mobile device, always return true
  if (mobileUA) {
    return true;
  }

  // For desktop, check if window is resized to mobile size
  const smallScreen = window.innerWidth <= 768;
  return smallScreen;
}

let isMobileDevice = checkIsMobile();

// Update mobile detection on resize
window.addEventListener("resize", () => {
  const wasMobile = isMobileDevice;
  isMobileDevice = checkIsMobile();

  // If changed from mobile to desktop or vice versa, hide current autocomplete
  if (wasMobile !== isMobileDevice) {
    hideAutocomplete();
  }
});

// Autocomplete keyword definitions with descriptions
const autocompleteKeywords = {
  // Global options (start with @)
  "@samplerate": { desc: "Sample rate in Hz", requiresNumber: true },
  "@volume": { desc: "Global volume (0-100)", requiresNumber: true },
  "@presetlist": { desc: "URL to preset list", requiresURL: true },
  "@background": { desc: "URL to background audio", requiresURL: true },
  "@gainlevel": {
    desc: "Gain level for background audio",
    options: ["off", "high", "medium", "low"],
  },

  // Regular keywords
  tone: { desc: "Carrier tone", next: ["binaural", "monaural", "isochronic"] },
  noise: { desc: "Noise type", options: ["white", "pink", "brown"] },
  background: {
    desc: "Background audio control",
    options: ["amplitude", "pulse", "spin"],
  },
  waveform: {
    desc: "Track waveform",
    options: ["sine", "square", "triangle", "sawtooth"],
    next: ["tone", "noise", "background"],
  },
  track: {
    desc: "Modify preset track parameter",
    requiresInteger: true,
    next: [
      "binaural",
      "monaural",
      "isochronic",
      "amplitude",
      "spin",
      "rate",
      "pulse",
      "tone",
    ],
  },
  binaural: { desc: "Binaural beat" },
  monaural: { desc: "Monaural beat" },
  isochronic: { desc: "Isochronic tone" },
  white: { desc: "White noise" },
  pink: { desc: "Pink noise" },
  brown: { desc: "Brown noise" },
  amplitude: { desc: "Volume level" },
  pulse: { desc: "Pulse effect" },
  spin: { desc: "Spin effect" },
  rate: { desc: "Rate parameter" },
  intensity: { desc: "Intensity level" },
  sine: { desc: "Sine wave" },
  square: { desc: "Square wave" },
  triangle: { desc: "Triangle wave" },
  sawtooth: { desc: "Sawtooth wave" },
  steady: { desc: "Linear transition" },
  "ease-in": { desc: "Ease-in transition" },
  "ease-out": { desc: "Ease-out transition" },
  smooth: { desc: "Smooth transition" },
};

let selectedOptionIndex = 0; // Track selected option in autocomplete menu

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

  // Apply syntax highlighting with proper scroll height matching
  highlight.innerHTML = `<div style="min-height: ${
    textarea.scrollHeight
  }px;">${highlightSyntax(code)}</div>`;
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

  // Check for autocomplete
  checkAutocomplete();

  // Debounced save to lastSequence
  saveCurrentSequenceDebounced();
});

// Handle Tab key for autocomplete
document.getElementById("spsqEditor").addEventListener("keydown", (e) => {
  // ESC always closes autocomplete and error balloons
  if (e.key === "Escape") {
    hideAutocomplete();
    return;
  }

  if (!currentSuggestion) return;

  if (e.key === "Tab") {
    e.preventDefault();
    applyAutocomplete();
  } else if (e.key === "ArrowDown" || e.key === "PageDown") {
    e.preventDefault();
    navigateAutocomplete(1);
  } else if (e.key === "ArrowUp" || e.key === "PageUp") {
    e.preventDefault();
    navigateAutocomplete(-1);
  }
});

// Autocomplete functions
function isValidURL(string) {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
}

function parseGlobalOptionContext(line) {
  // Global options format: @option <value>
  const trimmed = line.trim();
  const tokens = trimmed.split(/\s+/).filter((t) => t.length > 0);

  return {
    tokens,
    lastToken: tokens[tokens.length - 1] || "",
    isComplete:
      line.endsWith(" ") || (tokens.length === 1 && line.endsWith(tokens[0])),
    raw: line,
  };
}

function parseLineContext(line) {
  // Parse tokens from current line (after initial 2 spaces for keywords)
  const afterSpaces = line.substring(2);
  const trimmed = afterSpaces.trim();
  const tokens = trimmed.split(/\s+/).filter((t) => t.length > 0);

  return {
    tokens,
    lastToken: tokens[tokens.length - 1] || "",
    isComplete: afterSpaces.endsWith(" ") || afterSpaces === "", // Complete if ends with space or empty
    raw: line,
  };
}

// Parse timeline context (lines starting with timestamp)
function parseTimelineContext(line) {
  // Timeline format: <hh>:<mm>:<ss> <preset> [ramp:<slide>]
  const timestampRegex = /^(\d{2}:\d{2}:\d{2})\s*/;
  const match = line.match(timestampRegex);

  if (!match) {
    return { isTimeline: false };
  }

  const timestamp = match[1];
  const afterTimestamp = line.substring(match[0].length);
  const trimmed = afterTimestamp.trim();
  const tokens = trimmed.split(/\s+/).filter((t) => t.length > 0);

  return {
    isTimeline: true,
    timestamp,
    tokens,
    lastToken: tokens[tokens.length - 1] || "",
    isComplete: afterTimestamp.endsWith(" ") || afterTimestamp === "",
    raw: line,
  };
}

// Get available presets from the document
function getAvailablePresets() {
  const textarea = document.getElementById("spsqEditor");
  const text = textarea.value;
  const lines = text.split("\n");
  const presets = [];

  for (const line of lines) {
    // Match preset definitions: lines that start with word characters (no indentation, no @)
    // and are not comments (don't start with #) and are not timestamps
    const trimmed = line.trim();
    if (
      trimmed &&
      !trimmed.startsWith("#") &&
      !trimmed.match(/^\d{2}:\d{2}:\d{2}/)
    ) {
      // Check if line starts at column 0 (no indentation) and is just a name
      if (line.match(/^[\w-]+$/)) {
        presets.push(line.trim());
      }
    }
  }

  // Add built-in preset "silence"
  if (!presets.includes("silence")) {
    presets.push("silence");
  }

  return presets;
}

// Get timeline suggestions based on context
function getTimelineSuggestions(context) {
  const { tokens, isComplete, lastToken } = context;

  // After timestamp, suggest presets
  if (tokens.length === 0) {
    if (isComplete) {
      const presets = getAvailablePresets();
      if (presets.length > 0) {
        return presets.map((name) => ({ keyword: name, desc: "Preset" }));
      }
    }
    return null;
  }

  // While typing preset name, filter presets
  if (tokens.length === 1 && !isComplete) {
    const presets = getAvailablePresets();
    if (presets.length === 0) return null;
    const partial = lastToken.toLowerCase();
    const filtered = presets
      .filter((name) => name.toLowerCase().startsWith(partial))
      .map((name) => ({ keyword: name, desc: "Preset" }));
    // Return filtered results, or all presets if no match (user might be typing a custom preset)
    return filtered.length > 0
      ? filtered
      : presets.map((name) => ({ keyword: name, desc: "Preset" }));
  }

  // After preset name, suggest ramp/slide options
  if (tokens.length === 1 && isComplete) {
    return [
      { keyword: "steady", desc: autocompleteKeywords.steady.desc },
      { keyword: "ease-in", desc: autocompleteKeywords["ease-in"].desc },
      { keyword: "ease-out", desc: autocompleteKeywords["ease-out"].desc },
      { keyword: "smooth", desc: autocompleteKeywords.smooth.desc },
    ];
  }

  // While typing ramp, filter options
  if (tokens.length === 2 && !isComplete) {
    const partial = lastToken.toLowerCase();
    const rampOptions = [
      { keyword: "steady", desc: autocompleteKeywords.steady.desc },
      { keyword: "ease-in", desc: autocompleteKeywords["ease-in"].desc },
      { keyword: "ease-out", desc: autocompleteKeywords["ease-out"].desc },
      { keyword: "smooth", desc: autocompleteKeywords.smooth.desc },
    ];

    const filtered = rampOptions.filter((opt) =>
      opt.keyword.startsWith(partial)
    );
    return filtered.length > 0 ? filtered : null;
  }

  return null;
}

function getGlobalOptionSuggestions(context) {
  const { tokens, lastToken, isComplete } = context;

  // If typing @ or @something (first token), suggest all global options
  if (tokens.length === 1) {
    const partial = lastToken.toLowerCase();
    const options = [
      {
        keyword: "@samplerate",
        desc: autocompleteKeywords["@samplerate"].desc,
      },
      { keyword: "@volume", desc: autocompleteKeywords["@volume"].desc },
      {
        keyword: "@presetlist",
        desc: autocompleteKeywords["@presetlist"].desc,
      },
      {
        keyword: "@background",
        desc: autocompleteKeywords["@background"].desc,
      },
      {
        keyword: "@gainlevel",
        desc: autocompleteKeywords["@gainlevel"].desc,
      },
    ];

    // If just typed @, show all
    if (partial === "@") {
      return options;
    }

    // Filter by what's typed
    const filtered = options.filter((opt) => opt.keyword.startsWith(partial));
    return filtered.length > 0 ? filtered : null;
  }

  // If user typed @gainlevel, suggest the 4 level options
  if (tokens.length >= 1 && tokens[0] === "@gainlevel") {
    const partial = tokens.length === 2 ? lastToken.toLowerCase() : "";
    const levels = [
      { keyword: "off", desc: "0dB (no attenuation)" },
      { keyword: "high", desc: "-3dB" },
      { keyword: "medium", desc: "-9dB" },
      { keyword: "low", desc: "-18dB" },
    ];

    // If just typed space after @gainlevel, show all options
    if (isComplete && tokens.length === 1) {
      return levels;
    }

    // If typing a partial level name, filter
    if (tokens.length === 2) {
      const filtered = levels.filter((level) =>
        level.keyword.startsWith(partial)
      );
      return filtered.length > 0 ? filtered : null;
    }
  }

  // No suggestions after the option name for numbers/URLs
  return null;
}

function getNextSuggestions(context) {
  const { tokens, lastToken, isComplete } = context;

  // If empty line or just spaces, suggest main keywords
  if (tokens.length === 0) {
    return [
      {
        keyword: "tone",
        desc: autocompleteKeywords.tone?.desc || "Carrier tone",
      },
      {
        keyword: "noise",
        desc: autocompleteKeywords.noise?.desc || "Noise type",
      },
      {
        keyword: "background",
        desc:
          autocompleteKeywords.background?.desc || "Background audio control",
      },
      {
        keyword: "waveform",
        desc: autocompleteKeywords.waveform?.desc || "Track waveform",
      },
      {
        keyword: "track",
        desc: autocompleteKeywords.track?.desc || "Track volume",
      },
    ];
  }

  const firstKeyword = tokens[0];

  // Handle tone keyword
  if (firstKeyword === "tone") {
    if (tokens.length === 1 && !isComplete) return null; // typing "tone", wait
    if (tokens.length === 1 && isComplete) {
      // After "tone ", expect number (no suggestions)
      return null;
    }
    if (tokens.length === 2 && !isComplete) return null; // typing number
    if (tokens.length === 2 && isComplete) {
      // After "tone <num> ", suggest binaural/monaural/isochronic
      return [
        { keyword: "binaural", desc: autocompleteKeywords.binaural.desc },
        { keyword: "monaural", desc: autocompleteKeywords.monaural.desc },
        { keyword: "isochronic", desc: autocompleteKeywords.isochronic.desc },
      ];
    }
    if (tokens.length === 3 && !isComplete) {
      // typing binaural/monaural/isochronic, offer suggestions
      const partial = lastToken.toLowerCase();
      const validTypes = ["binaural", "monaural", "isochronic"];
      // Don't suggest if already a valid complete keyword
      if (validTypes.includes(partial)) return null;
      return [
        { keyword: "binaural", desc: autocompleteKeywords.binaural.desc },
        { keyword: "monaural", desc: autocompleteKeywords.monaural.desc },
        { keyword: "isochronic", desc: autocompleteKeywords.isochronic.desc },
      ].filter((s) => s.keyword.startsWith(partial));
    }
    if (tokens.length === 3 && isComplete) return null; // expect number
    if (tokens.length === 4 && !isComplete) return null; // typing number
    if (tokens.length === 4 && isComplete) {
      // After "tone <num> <type> <num> ", suggest amplitude
      return [
        { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
      ];
    }
    if (tokens.length === 5 && !isComplete) {
      const partial = lastToken.toLowerCase();
      // Don't suggest if already complete
      if (partial === "amplitude") return null;
      if ("amplitude".startsWith(partial)) {
        return [
          { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
        ];
      }
    }
  }

  // Handle noise keyword
  if (firstKeyword === "noise") {
    if (tokens.length === 1 && !isComplete) return null;
    if (tokens.length === 1 && isComplete) {
      // After "noise ", suggest white/pink/brown
      return [
        { keyword: "white", desc: autocompleteKeywords.white.desc },
        { keyword: "pink", desc: autocompleteKeywords.pink.desc },
        { keyword: "brown", desc: autocompleteKeywords.brown.desc },
      ];
    }
    if (tokens.length === 2 && !isComplete) {
      const partial = lastToken.toLowerCase();
      const validTypes = ["white", "pink", "brown"];
      // Don't suggest if already a valid complete keyword
      if (validTypes.includes(partial)) return null;
      return [
        { keyword: "white", desc: autocompleteKeywords.white.desc },
        { keyword: "pink", desc: autocompleteKeywords.pink.desc },
        { keyword: "brown", desc: autocompleteKeywords.brown.desc },
      ].filter((s) => s.keyword.startsWith(partial));
    }
    // After "noise <type> ", suggest amplitude
    if (tokens.length === 2 && isComplete) {
      return [
        { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
      ];
    }
    if (tokens.length === 3 && !isComplete) {
      const partial = lastToken.toLowerCase();
      // Don't suggest if already complete
      if (partial === "amplitude") return null;
      if ("amplitude".startsWith(partial)) {
        return [
          { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
        ];
      }
    }
    if (tokens.length === 4 && !isComplete) {
      const partial = lastToken.toLowerCase();
      // Don't suggest if already complete
      if (partial === "amplitude") return null;
      if ("amplitude".startsWith(partial)) {
        return [
          { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
        ];
      }
    }
  }

  // Handle background keyword
  if (firstKeyword === "background") {
    if (tokens.length === 1 && !isComplete) return null;
    if (tokens.length === 1 && isComplete) {
      return [
        { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
        { keyword: "pulse", desc: autocompleteKeywords.pulse.desc },
        { keyword: "spin", desc: autocompleteKeywords.spin.desc },
      ];
    }
    if (tokens.length === 2 && !isComplete) {
      const partial = lastToken.toLowerCase();
      const validTypes = ["amplitude", "pulse", "spin"];
      // Don't suggest if already a valid complete keyword
      if (validTypes.includes(partial)) return null;
      return [
        { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
        { keyword: "pulse", desc: autocompleteKeywords.pulse.desc },
        { keyword: "spin", desc: autocompleteKeywords.spin.desc },
      ].filter((s) => s.keyword.startsWith(partial));
    }

    const secondKeyword = tokens[1];

    // background amplitude <num>
    if (secondKeyword === "amplitude") {
      return null; // only expects number
    }

    // background pulse <num> intensity <num> amplitude <num>
    if (secondKeyword === "pulse") {
      if (tokens.length === 2 && isComplete) return null; // expect number
      if (tokens.length === 3 && !isComplete) return null; // typing number
      if (tokens.length === 3 && isComplete) {
        return [
          { keyword: "intensity", desc: autocompleteKeywords.intensity.desc },
        ];
      }
      if (tokens.length === 4 && !isComplete) {
        const partial = lastToken.toLowerCase();
        // Don't suggest if already complete
        if (partial === "intensity") return null;
        if ("intensity".startsWith(partial)) {
          return [
            { keyword: "intensity", desc: autocompleteKeywords.intensity.desc },
          ];
        }
      }
      if (tokens.length === 4 && isComplete) return null; // expect number
      if (tokens.length === 5 && !isComplete) return null; // typing number
      if (tokens.length === 5 && isComplete) {
        return [
          { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
        ];
      }
      if (tokens.length === 6 && !isComplete) {
        const partial = lastToken.toLowerCase();
        // Don't suggest if already complete
        if (partial === "amplitude") return null;
        if ("amplitude".startsWith(partial)) {
          return [
            { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
          ];
        }
      }
    }

    // background spin <num> rate <num> intensity <num> amplitude <num>
    if (secondKeyword === "spin") {
      if (tokens.length === 2 && isComplete) return null; // expect number
      if (tokens.length === 3 && !isComplete) return null; // typing number
      if (tokens.length === 3 && isComplete) {
        return [{ keyword: "rate", desc: autocompleteKeywords.rate.desc }];
      }
      if (tokens.length === 4 && !isComplete) {
        const partial = lastToken.toLowerCase();
        // Don't suggest if already complete
        if (partial === "rate") return null;
        if ("rate".startsWith(partial)) {
          return [{ keyword: "rate", desc: autocompleteKeywords.rate.desc }];
        }
      }
      if (tokens.length === 4 && isComplete) return null; // expect number
      if (tokens.length === 5 && !isComplete) return null; // typing number
      if (tokens.length === 5 && isComplete) {
        return [
          { keyword: "intensity", desc: autocompleteKeywords.intensity.desc },
        ];
      }
      if (tokens.length === 6 && !isComplete) {
        const partial = lastToken.toLowerCase();
        // Don't suggest if already complete
        if (partial === "intensity") return null;
        if ("intensity".startsWith(partial)) {
          return [
            { keyword: "intensity", desc: autocompleteKeywords.intensity.desc },
          ];
        }
      }
      if (tokens.length === 6 && isComplete) return null; // expect number
      if (tokens.length === 7 && !isComplete) return null; // typing number
      if (tokens.length === 7 && isComplete) {
        return [
          { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
        ];
      }
      if (tokens.length === 8 && !isComplete) {
        const partial = lastToken.toLowerCase();
        // Don't suggest if already complete
        if (partial === "amplitude") return null;
        if ("amplitude".startsWith(partial)) {
          return [
            { keyword: "amplitude", desc: autocompleteKeywords.amplitude.desc },
          ];
        }
      }
    }
  }

  // Handle waveform keyword
  if (firstKeyword === "waveform") {
    if (tokens.length === 1 && !isComplete) return null;
    if (tokens.length === 1 && isComplete) {
      return [
        { keyword: "sine", desc: autocompleteKeywords.sine.desc },
        { keyword: "square", desc: autocompleteKeywords.square.desc },
        { keyword: "triangle", desc: autocompleteKeywords.triangle.desc },
        { keyword: "sawtooth", desc: autocompleteKeywords.sawtooth.desc },
      ];
    }
    if (tokens.length === 2 && !isComplete) {
      const partial = lastToken.toLowerCase();
      const validTypes = ["sine", "square", "triangle", "sawtooth"];
      // Don't suggest if already a valid complete keyword
      if (validTypes.includes(partial)) return null;
      return [
        { keyword: "sine", desc: autocompleteKeywords.sine.desc },
        { keyword: "square", desc: autocompleteKeywords.square.desc },
        { keyword: "triangle", desc: autocompleteKeywords.triangle.desc },
        { keyword: "sawtooth", desc: autocompleteKeywords.sawtooth.desc },
      ].filter((s) => s.keyword.startsWith(partial));
    }
    // After waveform type, can have tone or background (noise doesn't support waveform)
    if (tokens.length === 2 && isComplete) {
      return [
        { keyword: "tone", desc: autocompleteKeywords.tone.desc },
        { keyword: "background", desc: autocompleteKeywords.background.desc },
      ];
    }
    // Handle waveform <type> tone/noise/background (recursive)
    if (tokens.length >= 3) {
      const waveformType = tokens[1];
      const validWaveforms = ["sine", "square", "triangle", "sawtooth"];
      if (validWaveforms.includes(waveformType)) {
        // Create new context starting from the third token
        const recursiveTokens = tokens.slice(2);
        const recursiveContext = {
          tokens: recursiveTokens,
          lastToken: recursiveTokens[recursiveTokens.length - 1] || "",
          isComplete: isComplete,
          raw: recursiveTokens.join(" "),
        };
        return getNextSuggestions(recursiveContext);
      }
    }
  }

  // Handle track keyword
  if (firstKeyword === "track") {
    // If we have track + number, suggest the parameter keywords
    if (tokens.length >= 2) {
      const trackNumber = tokens[1];
      const isValidInteger = /^[1-9]\d*$/.test(trackNumber);

      if (isValidInteger && tokens.length === 2 && isComplete) {
        // Show parameter options only when user typed "track 1 "
        return [
          { keyword: "binaural", desc: "Binaural beat frequency" },
          { keyword: "monaural", desc: "Monaural beat frequency" },
          { keyword: "isochronic", desc: "Isochronic tone frequency" },
          { keyword: "amplitude", desc: "Track amplitude (0-100)" },
          { keyword: "spin", desc: "Spin effect parameter" },
          { keyword: "rate", desc: "Rate parameter" },
          { keyword: "pulse", desc: "Pulse effect parameter" },
          { keyword: "tone", desc: "Carrier tone frequency" },
        ];
      }

      // If typing a parameter keyword (partial), filter suggestions
      if (isValidInteger && tokens.length === 3 && !isComplete) {
        const partial = lastToken.toLowerCase();
        const validParams = [
          "binaural",
          "monaural",
          "isochronic",
          "amplitude",
          "spin",
          "rate",
          "pulse",
          "tone",
        ];

        // Only show suggestions if the third token is NOT a valid parameter yet
        if (!validParams.includes(tokens[2])) {
          const options = [
            { keyword: "binaural", desc: "Binaural beat frequency" },
            { keyword: "monaural", desc: "Monaural beat frequency" },
            { keyword: "isochronic", desc: "Isochronic tone frequency" },
            { keyword: "amplitude", desc: "Track amplitude (0-100)" },
            { keyword: "spin", desc: "Spin effect parameter" },
            { keyword: "rate", desc: "Rate parameter" },
            { keyword: "pulse", desc: "Pulse effect parameter" },
            { keyword: "tone", desc: "Carrier tone frequency" },
          ];
          const filtered = options.filter((opt) =>
            opt.keyword.startsWith(partial)
          );
          return filtered.length > 0 ? filtered : null;
        }
      }
    }
    return null;
  }

  return null;
}

// Validate if numbers are present where expected
function validateGlobalOption(context) {
  const { tokens, isComplete } = context;

  if (tokens.length === 0) return null;

  const option = tokens[0];
  const isValidNumber = (str) => /^-?\d+(\.\d+)?$/.test(str);
  const makeError = (msg) => ({ error: true, message: msg });

  // Validate @samplerate
  if (option === "@samplerate") {
    if (tokens.length === 1 && isComplete) {
      return makeError("Expected: number (sample rate in Hz)");
    }
    if (tokens.length === 2 && !isValidNumber(tokens[1])) {
      return makeError("Expected: number (sample rate in Hz)");
    }
  }

  // Validate @volume
  if (option === "@volume") {
    if (tokens.length === 1 && isComplete) {
      return makeError("Expected: number (volume 0-100)");
    }
    if (tokens.length === 2 && !isValidNumber(tokens[1])) {
      return makeError("Expected: number (volume 0-100)");
    }
  }

  // Validate @gainlevel
  if (option === "@gainlevel") {
    const validLevels = ["off", "high", "medium", "low"];
    // Only validate if user has typed something after @gainlevel
    if (tokens.length === 2 && !validLevels.includes(tokens[1])) {
      return makeError("Expected: off, high, medium, or low");
    }
    // Don't show error if just "@gainlevel " - let autocomplete show suggestions
  }

  // Validate @presetlist
  if (option === "@presetlist") {
    if (tokens.length === 1 && isComplete) {
      return makeError("Expected: URL to preset list");
    }
    if (tokens.length === 2 && !isValidURL(tokens[1])) {
      return makeError("Expected: valid URL");
    }
  }

  // Validate @background
  if (option === "@background") {
    if (tokens.length === 1 && isComplete) {
      return makeError("Expected: URL to background audio");
    }
    if (tokens.length === 2 && !isValidURL(tokens[1])) {
      return makeError("Expected: valid URL");
    }
  }

  return null;
}

function validateSyntax(context) {
  const { tokens, isComplete } = context;

  if (tokens.length === 0) return null;

  const firstKeyword = tokens[0];
  const isValidNumber = (str) => /^-?\d+(\.\d+)?$/.test(str);

  // Helper to create error object
  const makeError = (msg) => ({ error: true, message: msg });

  // Handle tone keyword: tone <num> [binaural|monaural|isochronic] <num> amplitude <num>
  if (firstKeyword === "tone") {
    // After "tone ", must have number
    if (tokens.length === 1 && isComplete) {
      return makeError("Expected: number (base frequency)");
    }
    // Check if user typed invalid keyword instead of number after "tone"
    if (tokens.length === 2 && !isValidNumber(tokens[1])) {
      const validTypes = ["binaural", "monaural", "isochronic"];
      if (validTypes.includes(tokens[1])) {
        return makeError("Expected: number before " + tokens[1]);
      }
    }
    // After "tone <num> <type> ", must have number
    if (tokens.length === 3 && isComplete) {
      const validTypes = ["binaural", "monaural", "isochronic"];
      if (isValidNumber(tokens[1]) && validTypes.includes(tokens[2])) {
        return makeError("Expected: number (beat frequency)");
      }
    }
    // Check if user typed amplitude without number after type
    if (tokens.length === 4 && !isValidNumber(tokens[3])) {
      if (tokens[3] === "amplitude") {
        return makeError("Expected: number before amplitude");
      }
    }
    // After "tone <num> <type> <num> amplitude ", must have number
    if (tokens.length === 5 && tokens[4] === "amplitude" && isComplete) {
      return makeError("Expected: number (amplitude 0-100)");
    }
  }

  // Handle noise keyword: noise [white|pink|brown] amplitude <num>
  if (firstKeyword === "noise") {
    // After "noise ", must have type
    if (tokens.length === 1 && isComplete) {
      return makeError("Expected: noise type (white, pink, or brown)");
    }
    // After "noise <type> ", must have amplitude
    if (tokens.length === 2 && isComplete) {
      const validTypes = ["white", "pink", "brown"];
      if (validTypes.includes(tokens[1])) {
        return makeError("Expected: amplitude");
      }
    }
    // After "noise <type> amplitude ", must have number
    if (tokens.length === 3 && tokens[2] === "amplitude" && isComplete) {
      return makeError("Expected: number (amplitude 0-100)");
    }
  }

  // Handle background keyword
  if (firstKeyword === "background") {
    if (tokens.length >= 2) {
      const secondKeyword = tokens[1];

      // background amplitude <num>
      if (secondKeyword === "amplitude") {
        if (tokens.length === 2 && isComplete) {
          return makeError("Expected: number (amplitude 0-100)");
        }
      }

      // background pulse <num> intensity <num> amplitude <num>
      else if (secondKeyword === "pulse") {
        if (tokens.length === 2 && isComplete) {
          return makeError("Expected: number (pulse frequency)");
        }
        if (tokens.length === 4 && tokens[3] === "intensity" && isComplete) {
          return makeError("Expected: number (intensity 0-100)");
        }
        if (tokens.length === 6 && tokens[5] === "amplitude" && isComplete) {
          return makeError("Expected: number (amplitude 0-100)");
        }
      }

      // background spin <num> rate <num> intensity <num> amplitude <num>
      else if (secondKeyword === "spin") {
        if (tokens.length === 2 && isComplete) {
          return makeError("Expected: number (spin frequency)");
        }
        if (tokens.length === 4 && tokens[3] === "rate" && isComplete) {
          return makeError("Expected: number (rotation rate)");
        }
        if (tokens.length === 6 && tokens[5] === "intensity" && isComplete) {
          return makeError("Expected: number (intensity 0-100)");
        }
        if (tokens.length === 8 && tokens[7] === "amplitude" && isComplete) {
          return makeError("Expected: number (amplitude 0-100)");
        }
      }
    }
  }

  // Handle waveform keyword
  if (firstKeyword === "waveform") {
    if (tokens.length >= 3) {
      const waveformType = tokens[1];
      const validWaveforms = ["sine", "square", "triangle", "sawtooth"];
      if (validWaveforms.includes(waveformType)) {
        // Recursively validate the tone/noise/background part
        const recursiveTokens = tokens.slice(2);
        const recursiveContext = {
          tokens: recursiveTokens,
          lastToken: recursiveTokens[recursiveTokens.length - 1] || "",
          isComplete: isComplete,
          raw: recursiveTokens.join(" "),
        };
        return validateSyntax(recursiveContext);
      }
    }
  }

  // Handle track keyword: track <num>
  if (firstKeyword === "track") {
    if (tokens.length === 1 && isComplete) {
      return makeError("Expected: track number (positive integer)");
    }
    if (tokens.length === 2) {
      const trackNum = tokens[1];
      const isValidInteger = /^[1-9]\d*$/.test(trackNum);
      if (!isValidInteger) {
        return makeError("Expected: positive integer (cannot be 0 or decimal)");
      }
    }
    if (tokens.length === 3) {
      const validParams = [
        "binaural",
        "monaural",
        "isochronic",
        "amplitude",
        "spin",
        "rate",
        "pulse",
        "tone",
      ];

      // Check if it's a valid parameter
      if (!validParams.includes(tokens[2])) {
        if (isComplete) {
          return makeError(
            "Expected: binaural, monaural, isochronic, amplitude, spin, rate, pulse, or tone"
          );
        }
      } else {
        // Valid parameter, now expect a number
        if (isComplete) {
          return makeError("Expected: number value for " + tokens[2]);
        }
      }
    }
    if (tokens.length === 4) {
      const param = tokens[2];
      const value = tokens[3];
      const isValidNumber = /^-?\d+(\.\d+)?$/.test(value);
      if (!isValidNumber) {
        return makeError("Expected: number value for " + param);
      }
    }
  }

  return null; // No error
}

function checkAutocomplete() {
  const textarea = document.getElementById("spsqEditor");
  const cursorPos = textarea.selectionStart;
  const text = textarea.value;

  // Get current line
  const beforeCursor = text.substring(0, cursorPos);
  const lastNewLine = beforeCursor.lastIndexOf("\n");
  const currentLine = beforeCursor.substring(lastNewLine + 1);

  // Check if it's a global option line (starts with @)
  if (/^@/.test(currentLine)) {
    const context = parseGlobalOptionContext(currentLine);

    // Check for validation errors
    const validationError = validateGlobalOption(context);
    if (validationError) {
      showValidationError(validationError.message, cursorPos);
      return;
    }

    // Show suggestions
    const suggestions = getGlobalOptionSuggestions(context);
    if (suggestions && suggestions.length > 0) {
      showAutocomplete(suggestions, cursorPos);
    } else {
      hideAutocomplete();
    }
    return;
  }

  // Check if it's a timeline line (starts with timestamp)
  const timelineContext = parseTimelineContext(currentLine);
  if (timelineContext.isTimeline) {
    const suggestions = getTimelineSuggestions(timelineContext);
    if (suggestions && suggestions.length > 0) {
      showAutocomplete(suggestions, cursorPos);
    } else {
      hideAutocomplete();
    }
    return;
  }

  // Only show autocomplete for lines starting with 2 spaces (keyword lines)
  if (!/^  /.test(currentLine)) {
    hideAutocomplete();
    return;
  }

  // Parse context
  const context = parseLineContext(currentLine);

  // Check for syntax errors first
  const validationError = validateSyntax(context);
  if (validationError) {
    showValidationError(validationError.message, cursorPos);
    return;
  }

  // If no errors, check for autocomplete suggestions
  const suggestions = getNextSuggestions(context);

  if (suggestions && suggestions.length > 0) {
    showAutocomplete(suggestions, cursorPos);
  } else {
    hideAutocomplete();
  }
}

function showValidationError(message, cursorPos) {
  const textarea = document.getElementById("spsqEditor");
  const balloon = document.getElementById("autocompleteBalloon");
  const optionsContainer = balloon.querySelector(".autocomplete-options");

  currentSuggestion = null;
  selectedOptionIndex = 0;

  // Clear and show error message
  optionsContainer.innerHTML = "";
  const errorDiv = document.createElement("div");
  errorDiv.className = "autocomplete-error";

  // Create icon element using Lucide
  const icon = document.createElement("i");
  icon.setAttribute("data-lucide", "alert-triangle");
  icon.style.width = "1rem";
  icon.style.height = "1rem";
  icon.style.display = "inline-block";
  icon.style.marginRight = "0.5rem";

  errorDiv.appendChild(icon);
  errorDiv.appendChild(document.createTextNode(message));
  optionsContainer.appendChild(errorDiv);

  // Initialize Lucide icons
  if (window.lucide) {
    window.lucide.createIcons();
  }

  // Show balloon first to calculate height
  balloon.classList.add("show");

  // Calculate position (above cursor with more offset)
  const coords = getCaretCoordinates(textarea, cursorPos);
  const balloonRect = balloon.getBoundingClientRect();
  const balloonHeight = balloonRect.height;
  const lineHeight = parseInt(getComputedStyle(textarea).lineHeight);

  balloon.style.left = coords.left + "px";
  balloon.style.top = coords.top - balloonHeight - lineHeight - 5 + "px";
}

function showAutocomplete(suggestions, cursorPos) {
  // Use mobile bar on mobile devices
  if (isMobileDevice) {
    showMobileAutocomplete(suggestions);
    return;
  }

  const textarea = document.getElementById("spsqEditor");
  const balloon = document.getElementById("autocompleteBalloon");
  const optionsContainer = balloon.querySelector(".autocomplete-options");

  currentSuggestion = suggestions;
  selectedOptionIndex = 0;

  // Clear and populate options
  optionsContainer.innerHTML = "";
  suggestions.forEach((suggestion, index) => {
    const div = document.createElement("div");
    div.className = "autocomplete-option" + (index === 0 ? " selected" : "");
    div.dataset.index = index;

    // Create keyword and description elements
    const keywordSpan = document.createElement("span");
    keywordSpan.className = "autocomplete-keyword";
    keywordSpan.textContent = suggestion.keyword;

    const descSpan = document.createElement("span");
    descSpan.className = "autocomplete-desc";
    descSpan.textContent = suggestion.desc;

    div.appendChild(keywordSpan);
    div.appendChild(descSpan);

    // Click handler for touch/mouse
    div.addEventListener("click", (e) => {
      selectedOptionIndex = parseInt(e.currentTarget.dataset.index);
      applyAutocomplete();
    });

    optionsContainer.appendChild(div);
  });

  // Show balloon first to calculate height
  balloon.classList.add("show");

  // Calculate position (above cursor with more offset)
  const coords = getCaretCoordinates(textarea, cursorPos);
  const balloonHeight = balloon.offsetHeight;
  const lineHeight = parseFloat(getComputedStyle(textarea).lineHeight);

  balloon.style.left = coords.left + "px";
  balloon.style.top = coords.top - balloonHeight - lineHeight - 5 + "px"; // More space above
}

function navigateAutocomplete(direction) {
  if (!currentSuggestion) return;

  const optionsContainer = document.querySelector(".autocomplete-options");
  const options = optionsContainer.querySelectorAll(".autocomplete-option");

  // Remove current selection
  options[selectedOptionIndex].classList.remove("selected");

  // Update index
  selectedOptionIndex += direction;
  if (selectedOptionIndex < 0) selectedOptionIndex = options.length - 1;
  if (selectedOptionIndex >= options.length) selectedOptionIndex = 0;

  // Add new selection
  options[selectedOptionIndex].classList.add("selected");

  // Scroll into view
  options[selectedOptionIndex].scrollIntoView({
    block: "nearest",
    behavior: "smooth",
  });
}

function hideAutocomplete() {
  const balloon = document.getElementById("autocompleteBalloon");
  balloon.classList.remove("show");

  // Always try to hide mobile bar too (in case of mode switch)
  const mobileBar = document.getElementById("mobileAutocompleteBar");
  if (mobileBar) {
    mobileBar.classList.remove("show");
  }

  currentSuggestion = null;
  selectedOptionIndex = 0;
}

function showMobileAutocomplete(suggestions) {
  // First, make sure desktop balloon is hidden
  const balloon = document.getElementById("autocompleteBalloon");
  balloon.classList.remove("show");

  const mobileBar = document.getElementById("mobileAutocompleteBar");
  const chipsContainer = mobileBar.querySelector(".mobile-autocomplete-chips");

  // Clear existing chips
  chipsContainer.innerHTML = "";

  // Create chips for each suggestion
  suggestions.forEach((suggestion) => {
    const chip = document.createElement("div");
    chip.className = "mobile-autocomplete-chip";

    const keyword = document.createElement("div");
    keyword.className = "mobile-autocomplete-chip-keyword";
    keyword.textContent = suggestion.keyword;
    chip.appendChild(keyword);

    if (suggestion.desc) {
      const desc = document.createElement("div");
      desc.className = "mobile-autocomplete-chip-desc";
      desc.textContent = suggestion.desc;
      chip.appendChild(desc);
    }

    // Add click handler
    chip.addEventListener("click", () => {
      applyMobileAutocomplete(suggestion.keyword);
    });

    chipsContainer.appendChild(chip);
  });

  // Show the bar
  mobileBar.classList.add("show");
}

function hideMobileAutocomplete() {
  const mobileBar = document.getElementById("mobileAutocompleteBar");
  mobileBar.classList.remove("show");
}

function applyMobileAutocomplete(keyword) {
  const textarea = document.getElementById("spsqEditor");
  const cursorPos = textarea.selectionStart;
  const text = textarea.value;

  // Get current line
  const beforeCursor = text.substring(0, cursorPos);
  const lastNewLine = beforeCursor.lastIndexOf("\n");
  const currentLine = beforeCursor.substring(lastNewLine + 1);
  const lineStart = lastNewLine + 1;

  // Check if global option line
  if (/^@/.test(currentLine)) {
    const newText =
      text.substring(0, lineStart) + keyword + " " + text.substring(cursorPos);
    textarea.value = newText;
    textarea.selectionStart = textarea.selectionEnd =
      lineStart + keyword.length + 1;

    // Hide mobile autocomplete and update display
    hideMobileAutocomplete();
    updateSyntaxHighlight();
    updateLineNumbers();
    checkAutocomplete();
    return;
  }

  // Check if timeline or keyword line
  const timelineContext = parseTimelineContext(currentLine);

  if (timelineContext.isTimeline) {
    // Timeline autocomplete
    const afterTimestamp = currentLine
      .substring(timelineContext.timestamp.length)
      .trimStart();
    const tokens = afterTimestamp.split(/\s+/).filter((t) => t);

    let newText;
    if (tokens.length === 0) {
      // No preset yet, insert preset name
      newText =
        text.substring(0, lineStart) +
        timelineContext.timestamp +
        " " +
        keyword +
        text.substring(cursorPos);
      textarea.value = newText;
      textarea.selectionStart = textarea.selectionEnd =
        lineStart + timelineContext.timestamp.length + 1 + keyword.length + 1;
    } else if (tokens.length === 1) {
      // Has preset, insert ramp
      newText =
        text.substring(0, lineStart) +
        timelineContext.timestamp +
        " " +
        tokens[0] +
        " " +
        keyword +
        text.substring(cursorPos);
      textarea.value = newText;
      textarea.selectionStart = textarea.selectionEnd =
        lineStart +
        timelineContext.timestamp.length +
        1 +
        tokens[0].length +
        1 +
        keyword.length +
        1;
    }
  } else {
    // Keyword line autocomplete
    const context = parseLineContext(currentLine);
    const tokens = context.tokens;
    const partial = context.lastToken;

    let newText;
    if (context.isComplete) {
      // Just append the keyword
      newText =
        text.substring(0, cursorPos) +
        keyword +
        " " +
        text.substring(cursorPos);
      textarea.value = newText;
      textarea.selectionStart = textarea.selectionEnd =
        cursorPos + keyword.length + 1;
    } else {
      // Replace partial token
      const tokenStartPos = cursorPos - partial.length;
      newText =
        text.substring(0, tokenStartPos) +
        keyword +
        " " +
        text.substring(cursorPos);
      textarea.value = newText;
      textarea.selectionStart = textarea.selectionEnd =
        tokenStartPos + keyword.length + 1;
    }
  }

  // Hide mobile autocomplete and update display
  hideMobileAutocomplete();
  updateSyntaxHighlight();
  updateLineNumbers();
  checkAutocomplete();
}

function applyAutocomplete() {
  if (!currentSuggestion) return;

  const textarea = document.getElementById("spsqEditor");
  const cursorPos = textarea.selectionStart;
  const text = textarea.value;

  // Get selected keyword
  const selectedKeyword = currentSuggestion[selectedOptionIndex].keyword;

  // Get current line
  const beforeCursor = text.substring(0, cursorPos);
  const lastNewLine = beforeCursor.lastIndexOf("\n");
  const currentLine = beforeCursor.substring(lastNewLine + 1);
  const lineStart = lastNewLine + 1;

  // Check if it's a global option line
  if (/^@/.test(currentLine)) {
    const context = parseGlobalOptionContext(currentLine);
    const { tokens, isComplete, lastToken } = context;
    let newLine;

    // If completing a @gainlevel option (tokens[0] is @gainlevel)
    if (tokens.length >= 1 && tokens[0] === "@gainlevel") {
      // Keep @gainlevel and add/replace the option
      if (tokens.length === 1 && isComplete) {
        // User typed "@gainlevel ", append the option
        newLine = currentLine + selectedKeyword + " ";
      } else if (tokens.length === 2) {
        // User is typing a partial option, replace it
        const optionStart = currentLine.indexOf(tokens[1]);
        newLine = currentLine.substring(0, optionStart) + selectedKeyword + " ";
      } else {
        // Default: append
        newLine = currentLine + selectedKeyword + " ";
      }
    } else {
      // For other global options (@samplerate, @volume, etc.)
      if (isComplete || !lastToken || lastToken === "@") {
        // Just append
        newLine = selectedKeyword + " ";
      } else {
        // Replace partial token
        newLine = selectedKeyword + " ";
      }
    }

    // Get rest of document after current line
    const afterCursor = text.substring(cursorPos);
    const nextNewLine = afterCursor.indexOf("\n");
    const restOfDoc =
      nextNewLine === -1 ? "" : afterCursor.substring(nextNewLine);

    // Build new text
    const newText = text.substring(0, lineStart) + newLine + restOfDoc;
    textarea.value = newText;

    // Set cursor after inserted keyword
    const newCursorPos = lineStart + newLine.length;
    textarea.setSelectionRange(newCursorPos, newCursorPos);

    // Update UI
    updateLineNumbers();
    updateSyntaxHighlight();
    hideAutocomplete();

    // Trigger autocomplete check again
    setTimeout(() => checkAutocomplete(), 50);
    return;
  }

  // Check if it's a timeline line
  const timelineContext = parseTimelineContext(currentLine);
  if (timelineContext.isTimeline) {
    const { tokens, isComplete, timestamp, lastToken } = timelineContext;
    let newLine;

    if (isComplete || !lastToken) {
      // Just append
      newLine = currentLine + selectedKeyword + " ";
    } else {
      // Replace partial token
      const afterTimestamp = currentLine
        .substring(timestamp.length)
        .trimStart();
      const existingPart =
        timestamp +
        " " +
        afterTimestamp.substring(0, afterTimestamp.lastIndexOf(lastToken));
      newLine = existingPart + selectedKeyword + " ";
    }

    // Get rest of document after current line
    const afterCursor = text.substring(cursorPos);
    const nextNewLine = afterCursor.indexOf("\n");
    const restOfDoc =
      nextNewLine === -1 ? "" : afterCursor.substring(nextNewLine);

    // Build new text
    const newText = text.substring(0, lineStart) + newLine + restOfDoc;
    textarea.value = newText;

    // Set cursor after inserted keyword
    const newCursorPos = lineStart + newLine.length;
    textarea.setSelectionRange(newCursorPos, newCursorPos);

    // Update UI
    updateLineNumbers();
    updateSyntaxHighlight();
    hideAutocomplete();

    // Trigger autocomplete check again
    setTimeout(() => checkAutocomplete(), 50);
    return;
  }

  // Original logic for keyword lines
  const context = parseLineContext(currentLine);
  const { tokens, isComplete } = context;

  let newLine;
  if (tokens.length === 0 || (tokens.length === 1 && !isComplete)) {
    // Replace entire content after "  "
    newLine = "  " + selectedKeyword + " ";
  } else {
    // Append to existing line
    if (isComplete || !context.lastToken) {
      // If line ends with space or no lastToken, just append
      newLine = currentLine + selectedKeyword + " ";
    } else {
      // Replace the partial token being typed
      const existingPart = currentLine.substring(
        0,
        currentLine.lastIndexOf(context.lastToken)
      );
      newLine = existingPart + selectedKeyword + " ";
    }
  }

  // Get rest of document after current line
  const afterCursor = text.substring(cursorPos);
  const nextNewLine = afterCursor.indexOf("\n");
  const restOfDoc =
    nextNewLine === -1 ? "" : afterCursor.substring(nextNewLine);

  // Build new text
  const newText = text.substring(0, lineStart) + newLine + restOfDoc;
  textarea.value = newText;

  // Set cursor after inserted keyword
  const newCursorPos = lineStart + newLine.length;
  textarea.setSelectionRange(newCursorPos, newCursorPos);

  // Update UI
  updateLineNumbers();
  updateSyntaxHighlight();
  hideAutocomplete();

  // Trigger autocomplete check again (for next suggestions)
  setTimeout(() => checkAutocomplete(), 50);
}

// Get caret coordinates (for positioning balloon)
function getCaretCoordinates(element, position) {
  const div = document.createElement("div");
  const style = getComputedStyle(element);

  // Copy styles
  [
    "fontFamily",
    "fontSize",
    "fontWeight",
    "letterSpacing",
    "lineHeight",
    "padding",
    "border",
  ].forEach((prop) => {
    div.style[prop] = style[prop];
  });

  div.style.position = "absolute";
  div.style.visibility = "hidden";
  div.style.whiteSpace = "pre-wrap";
  div.style.wordWrap = "break-word";
  div.style.width = element.clientWidth + "px";

  document.body.appendChild(div);

  const text = element.value.substring(0, position);
  div.textContent = text;

  const span = document.createElement("span");
  span.textContent = element.value.substring(position) || ".";
  div.appendChild(span);

  const coordinates = {
    top: span.offsetTop + element.offsetTop - element.scrollTop,
    left: span.offsetLeft + element.offsetLeft - element.scrollLeft,
  };

  document.body.removeChild(div);
  return coordinates;
}

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

// Extract total duration from SPSQ timeline
function getTimelineDuration() {
  const spsq = document.getElementById("spsqEditor").value;
  const lines = spsq.split("\n");

  let lastTimestamp = null;

  // Find the last timeline entry
  for (let i = lines.length - 1; i >= 0; i--) {
    const line = lines[i].trim();
    // Match timestamp format: hh:mm:ss
    const timestampMatch = line.match(/^(\d{2}):(\d{2}):(\d{2})\s+/);
    if (timestampMatch) {
      const hours = parseInt(timestampMatch[1], 10);
      const minutes = parseInt(timestampMatch[2], 10);
      const seconds = parseInt(timestampMatch[3], 10);
      lastTimestamp = hours * 3600 + minutes * 60 + seconds;
      break;
    }
  }

  return lastTimestamp || 0;
}

// Progress tracking
function updateProgress() {
  const currentTime = synapseq.getCurrentTime();
  const duration = getTimelineDuration();

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

    synapseq.onplaying = () => {
      setStatus("Playing...");
      document.getElementById("playBtn").disabled = true;
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
        navigator.mediaSession.setActionHandler("stop", () => synapseq.stop());
        navigator.mediaSession.playbackState = "playing";
      }
    };

    synapseq.onstopped = () => {
      setStatus("Stopped");
      document.getElementById("playBtn").disabled = false;
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
      console.error("Error:", detail.error);
      showError(detail.error.message || detail.error);
      setStatus("Error");
      document.getElementById("playBtn").disabled = false;
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
    showError(error.message);
    setStatus("Error");
    document.getElementById("playBtn").disabled = false;
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
