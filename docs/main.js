// Initialize Lucide icons
document.addEventListener("DOMContentLoaded", () => {
  lucide.createIcons();

  // Mobile menu toggle
  initMobileMenu();

  // Sidebar navigation
  initSidebarNav();

  // Copy buttons
  initCopyButtons();

  // Expandable examples
  initExpandableExamples();

  // Smooth scroll
  initSmoothScroll();

  // Scroll spy for active navigation
  initScrollSpy();

  // Load Google Charts
  loadGoogleCharts();
});

// Load and draw transition charts
function loadGoogleCharts() {
  google.charts.load("current", { packages: ["corechart"] });
  google.charts.setOnLoadCallback(drawTransitionCharts);
}

function drawTransitionCharts() {
  // Check if chart containers exist before drawing
  if (!document.getElementById("chart-steady")) {
    return; // Exit early if containers don't exist
  }

  // Chart options with dark theme
  const darkOptions = {
    backgroundColor: "transparent",
    legend: { position: "none" },
    hAxis: {
      title: "Time (minutes)",
      titleTextStyle: { color: "#cbd5e1", fontSize: 11 },
      textStyle: { color: "#94a3b8", fontSize: 10 },
      gridlines: { color: "#334155", count: 5 },
      minorGridlines: { color: "transparent" },
      baselineColor: "#475569",
    },
    vAxis: {
      title: "Frequency (Hz)",
      titleTextStyle: { color: "#cbd5e1", fontSize: 11 },
      textStyle: { color: "#94a3b8", fontSize: 10 },
      gridlines: { color: "#334155" },
      minorGridlines: { color: "transparent" },
      baselineColor: "#475569",
      viewWindow: { min: 0, max: 25 },
    },
    chartArea: {
      left: 60,
      top: 20,
      right: 20,
      bottom: 50,
      backgroundColor: "transparent",
    },
    curveType: "function",
    lineWidth: 3,
    pointSize: 0,
    colors: ["#22d3ee"],
    animation: {
      duration: 1000,
      easing: "out",
      startup: true,
    },
  };

  // Steady transition (linear)
  const steadyData = google.visualization.arrayToDataTable([
    ["Time", "Frequency"],
    [0, 20],
    [1, 17.5],
    [2, 15],
    [3, 12.5],
    [4, 10],
    [5, 7.5],
    [6, 5],
  ]);

  const steadyChart = new google.visualization.LineChart(
    document.getElementById("chart-steady")
  );
  steadyChart.draw(steadyData, darkOptions);

  // Ease-out transition (logarithmic)
  const easeOutData = google.visualization.arrayToDataTable([
    ["Time", "Frequency"],
    [0, 20],
    [0.5, 14],
    [1, 10.5],
    [2, 8],
    [3, 6.5],
    [4, 5.8],
    [5, 5.3],
    [6, 5],
  ]);

  const easeOutChart = new google.visualization.LineChart(
    document.getElementById("chart-ease-out")
  );
  const easeOutOptions = { ...darkOptions, colors: ["#60a5fa"] };
  easeOutChart.draw(easeOutData, easeOutOptions);

  // Ease-in transition (exponential)
  const easeInData = google.visualization.arrayToDataTable([
    ["Time", "Frequency"],
    [0, 20],
    [1, 19.5],
    [2, 18.5],
    [3, 16],
    [4, 12],
    [5, 8],
    [6, 5],
  ]);

  const easeInChart = new google.visualization.LineChart(
    document.getElementById("chart-ease-in")
  );
  const easeInOptions = { ...darkOptions, colors: ["#4ade80"] };
  easeInChart.draw(easeInData, easeInOptions);

  // Smooth transition (sigmoid)
  const smoothData = google.visualization.arrayToDataTable([
    ["Time", "Frequency"],
    [0, 20],
    [0.5, 19.5],
    [1, 18.5],
    [2, 15],
    [3, 10],
    [4, 6.5],
    [5, 5.5],
    [6, 5],
  ]);

  const smoothChart = new google.visualization.LineChart(
    document.getElementById("chart-smooth")
  );
  const smoothOptions = { ...darkOptions, colors: ["#a78bfa"] };
  smoothChart.draw(smoothData, smoothOptions);

  // Redraw charts on window resize
  window.addEventListener("resize", () => {
    steadyChart.draw(steadyData, darkOptions);
    easeOutChart.draw(easeOutData, easeOutOptions);
    easeInChart.draw(easeInData, easeInOptions);
    smoothChart.draw(smoothData, smoothOptions);
  });
}

// Mobile Menu
function initMobileMenu() {
  const toggle = document.querySelector(".mobile-menu-toggle");
  const sidebar = document.querySelector(".sidebar");

  if (!toggle || !sidebar) return;

  toggle.addEventListener("click", () => {
    sidebar.classList.toggle("active");
  });

  // Close sidebar when clicking outside
  document.addEventListener("click", (e) => {
    if (!sidebar.contains(e.target) && !toggle.contains(e.target)) {
      sidebar.classList.remove("active");
    }
  });

  // Close sidebar when clicking a nav item (mobile)
  const navItems = sidebar.querySelectorAll(".nav-item");
  navItems.forEach((item) => {
    item.addEventListener("click", () => {
      if (window.innerWidth <= 768) {
        sidebar.classList.remove("active");
      }
    });
  });
}

// Sidebar Navigation
function initSidebarNav() {
  const navItems = document.querySelectorAll('.nav-item[href^="#"]');

  navItems.forEach((item) => {
    item.addEventListener("click", (e) => {
      e.preventDefault();
      const targetId = item.getAttribute("href");
      const targetSection = document.querySelector(targetId);

      if (targetSection) {
        const offset = 100; // Account for fixed header
        const targetPosition = targetSection.offsetTop - offset;

        window.scrollTo({
          top: targetPosition,
          behavior: "smooth",
        });

        // Update active state
        navItems.forEach((navItem) => navItem.classList.remove("active"));
        item.classList.add("active");
      }
    });
  });
}

// Copy Buttons
function initCopyButtons() {
  const copyButtons = document.querySelectorAll(".copy-btn");

  copyButtons.forEach((button) => {
    button.addEventListener("click", async () => {
      const code = button.getAttribute("data-code");
      const expandableId = button.getAttribute("data-copy-expandable");

      let textToCopy = "";

      if (expandableId) {
        // Copy from expandable example code block
        const codeBlock = document.getElementById(expandableId);
        if (codeBlock) {
          textToCopy = codeBlock.textContent;
        }
      } else if (code) {
        // Decode HTML entities from data-code attribute
        textToCopy = code
          .replace(/&#10;/g, "\n")
          .replace(/&lt;/g, "<")
          .replace(/&gt;/g, ">")
          .replace(/&amp;/g, "&");
      }

      if (!textToCopy) return;

      try {
        await navigator.clipboard.writeText(textToCopy);

        // Visual feedback
        const icon = button.querySelector("i");
        const originalIcon = icon.getAttribute("data-lucide");

        // Change to check icon
        icon.setAttribute("data-lucide", "check");
        lucide.createIcons();

        // Reset after 2 seconds
        setTimeout(() => {
          icon.setAttribute("data-lucide", originalIcon || "copy");
          lucide.createIcons();
        }, 2000);
      } catch (err) {
        console.error("Failed to copy:", err);
      }
    });
  });
}

// Expandable Examples
function initExpandableExamples() {
  const expandables = document.querySelectorAll(".expandable-example");

  expandables.forEach((expandable) => {
    const header = expandable.querySelector(".expandable-header");

    header.addEventListener("click", () => {
      const isActive = expandable.classList.contains("active");

      // Close all other expandables
      expandables.forEach((exp) => {
        exp.classList.remove("active");
      });

      // Toggle current expandable
      if (!isActive) {
        expandable.classList.add("active");
        // Reinitialize icons after DOM change
        setTimeout(() => lucide.createIcons(), 100);
      }
    });
  });
}

// Smooth Scroll
function initSmoothScroll() {
  document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
    anchor.addEventListener("click", function (e) {
      const href = this.getAttribute("href");

      // Skip if it's just "#"
      if (href === "#") return;

      e.preventDefault();

      const target = document.querySelector(href);
      if (target) {
        const offset = 100;
        const targetPosition = target.offsetTop - offset;

        window.scrollTo({
          top: targetPosition,
          behavior: "smooth",
        });
      }
    });
  });
}

// Scroll Spy
function initScrollSpy() {
  const sections = document.querySelectorAll(".doc-section[id]");
  const navItems = document.querySelectorAll('.nav-item[href^="#"]');

  const observerOptions = {
    root: null,
    rootMargin: "-20% 0px -70% 0px",
    threshold: 0,
  };

  const observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        const id = entry.target.getAttribute("id");
        const activeNavItem = document.querySelector(
          `.nav-item[href="#${id}"]`
        );

        if (activeNavItem) {
          // Remove active from all nav items
          navItems.forEach((item) => item.classList.remove("active"));

          // Add active to current item
          activeNavItem.classList.add("active");

          // Scroll nav item into view if needed
          const sidebar = document.querySelector(".sidebar");
          const sidebarRect = sidebar.getBoundingClientRect();
          const itemRect = activeNavItem.getBoundingClientRect();

          if (
            itemRect.top < sidebarRect.top ||
            itemRect.bottom > sidebarRect.bottom
          ) {
            activeNavItem.scrollIntoView({
              behavior: "smooth",
              block: "nearest",
            });
          }
        }
      }
    });
  }, observerOptions);

  sections.forEach((section) => {
    observer.observe(section);
  });
}

// Back to top on logo click (if needed)
const backHome = document.querySelector(".back-home");
if (backHome) {
  backHome.addEventListener("click", (e) => {
    // If we're on the docs page, allow normal navigation
    // This is just ensuring the link works as expected
  });
}

// Dynamic page loading system
const pageMap = {
  introduction: "pages/introduction.html",
  "brainwave-entrainment": "pages/brainwave-entrainment.html",
  syntax: "pages/syntax.html",
  // All elements are in syntax.html
  tone: "pages/syntax.html",
  noise: "pages/syntax.html",
  background: "pages/syntax.html",
  waveform: "pages/syntax.html",
  presets: "pages/syntax.html",
  timeline: "pages/syntax.html",
  transitions: "pages/syntax.html",
  comments: "pages/syntax.html",
  "global-options": "pages/global-options.html",
  "command-line": "pages/command-line.html",
  compilation: "pages/compilation.html",
  programming: "pages/programming.html",
  notes: "pages/notes.html",
  // Sub-sections will fall back to their parent sections
  "option-background": "pages/global-options.html",
  "option-gainlevel": "pages/global-options.html",
  "option-volume": "pages/global-options.html",
  "option-samplerate": "pages/global-options.html",
  "option-presetlist": "pages/global-options.html",
  "hub-commands": "pages/command-line.html",
  extract: "pages/command-line.html",
  "playing-mode": "pages/command-line.html",
  "export-options": "pages/command-line.html",
  "structured-formats": "pages/command-line.html",
  "ffmpeg-paths": "pages/command-line.html",
  "other-options": "pages/command-line.html",
  "windows-options": "pages/command-line.html",
  "installing-ffmpeg": "pages/command-line.html",
  "synapseq-js": "pages/programming.html",
  "go-library": "pages/programming.html",
};

async function loadPage(pageId) {
  const pageUrl = pageMap[pageId];
  if (!pageUrl) return false;

  const contentContainer = document.getElementById("dynamic-content");
  if (!contentContainer) return false;

  try {
    const response = await fetch(pageUrl);
    if (!response.ok) return false;

    const html = await response.text();

    // Extract base page name from URL for section ID
    const basePageId = pageUrl.split("/").pop().replace(".html", "");

    // Wrap content in section with doc-section class for proper styling
    contentContainer.innerHTML = `<section id="${basePageId}" class="doc-section">${html}</section>`;

    // Re-initialize Lucide icons for new content
    lucide.createIcons();

    // Re-initialize copy buttons for new content
    initCopyButtons();

    // Re-initialize expandable examples for new content
    initExpandableExamples();

    // Re-initialize Google Charts if needed (only if chart containers exist)
    if (typeof drawTransitionCharts === "function") {
      setTimeout(() => {
        if (document.getElementById("chart-steady")) {
          drawTransitionCharts();
        }
      }, 100);
    }

    // Scroll to specific element if it exists (for sub-sections)
    setTimeout(() => {
      const element = document.getElementById(pageId);
      if (element) {
        element.scrollIntoView({ behavior: "smooth", block: "start" });
      } else {
        // If no specific element, scroll to top
        window.scrollTo({ top: 0, behavior: "smooth" });
      }
    }, 100);

    return true;
  } catch (error) {
    console.error(`Failed to load page: ${pageId}`, error);
    return false;
  }
}

// Handle navigation clicks
function initDynamicNavigation() {
  const navLinks = document.querySelectorAll(".sidebar-nav a[href^='#']");

  navLinks.forEach((link) => {
    link.addEventListener("click", async (e) => {
      const href = link.getAttribute("href");
      const pageId = href.substring(1); // Remove '#'

      // Try to load dynamic page
      const loaded = await loadPage(pageId);

      if (loaded) {
        e.preventDefault();
        // Update URL without triggering page reload
        history.pushState(null, "", href);

        // Update active nav item
        navLinks.forEach((l) => l.classList.remove("active"));
        link.classList.add("active");

        // Scroll to top
        window.scrollTo({ top: 0, behavior: "smooth" });
      }
      // If not loaded, allow default anchor behavior
    });
  });

  // Handle browser back/forward
  window.addEventListener("popstate", () => {
    const hash = window.location.hash.substring(1);
    if (hash) {
      loadPage(hash);
    }
  });

  // Load initial page from URL hash
  const initialHash = window.location.hash.substring(1);
  if (initialHash && pageMap[initialHash]) {
    loadPage(initialHash);
  }
}

// Initialize dynamic navigation when DOM is ready
if (document.getElementById("dynamic-content")) {
  initDynamicNavigation();
}
