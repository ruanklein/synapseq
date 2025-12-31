// Initialize Lucide icons
document.addEventListener("DOMContentLoaded", () => {
  lucide.createIcons();

  // Mobile menu toggle
  initMobileMenu();

  // Sidebar navigation
  initSidebarNav();

  // Copy buttons
  initCopyButtons();

  // Smooth scroll
  initSmoothScroll();

  // Scroll spy for active navigation
  initScrollSpy();
});

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

      if (!code) return;

      try {
        // Decode HTML entities
        const decodedCode = code
          .replace(/&#10;/g, "\n")
          .replace(/&lt;/g, "<")
          .replace(/&gt;/g, ">")
          .replace(/&amp;/g, "&");

        await navigator.clipboard.writeText(decodedCode);

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
