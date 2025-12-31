// Initialize Lucide icons
lucide.createIcons();

// Mobile menu toggle
const menuButton = document.getElementById("menuButton");
const mobileMenu = document.getElementById("mobileMenu");

if (menuButton && mobileMenu) {
  menuButton.addEventListener("click", () => {
    mobileMenu.classList.toggle("active");

    // Recreate icons after DOM change
    setTimeout(() => lucide.createIcons(), 0);
  });

  // Close mobile menu when clicking a link
  const mobileLinks = document.querySelectorAll(".mobile-link");
  mobileLinks.forEach((link) => {
    link.addEventListener("click", () => {
      mobileMenu.classList.remove("active");
    });
  });

  // Close mobile menu when clicking outside
  document.addEventListener("click", (e) => {
    if (!menuButton.contains(e.target) && !mobileMenu.contains(e.target)) {
      mobileMenu.classList.remove("active");
    }
  });
}

// FAQ Accordion
const faqQuestions = document.querySelectorAll(".faq-question");

faqQuestions.forEach((question) => {
  question.addEventListener("click", () => {
    const faqItem = question.parentElement;
    const isActive = faqItem.classList.contains("active");

    // Close all other FAQ items
    document.querySelectorAll(".faq-item").forEach((item) => {
      item.classList.remove("active");
    });

    // Toggle current item
    if (!isActive) {
      faqItem.classList.add("active");
    }

    // Recreate icons after DOM change
    setTimeout(() => lucide.createIcons(), 0);
  });
});

// Copy to clipboard functionality
const copyButtons = document.querySelectorAll(".copy-btn");

copyButtons.forEach((button) => {
  button.addEventListener("click", async () => {
    const textToCopy = button.getAttribute("data-copy");

    try {
      await navigator.clipboard.writeText(textToCopy);

      // Visual feedback
      button.classList.add("copied");
      const icon = button.querySelector("[data-lucide]");
      if (icon) {
        icon.setAttribute("data-lucide", "check");
        lucide.createIcons();
      }

      // Reset after 2 seconds
      setTimeout(() => {
        button.classList.remove("copied");
        if (icon) {
          icon.setAttribute("data-lucide", "copy");
          lucide.createIcons();
        }
      }, 2000);
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  });
});

// Smooth scroll with offset for fixed navbar
document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
  anchor.addEventListener("click", function (e) {
    e.preventDefault();
    const targetId = this.getAttribute("href");

    if (targetId === "#") return;

    const targetElement = document.querySelector(targetId);
    if (targetElement) {
      const offset = 100; // Navbar height + padding
      const targetPosition = targetElement.offsetTop - offset;

      window.scrollTo({
        top: targetPosition,
        behavior: "smooth",
      });
    }
  });
});

// Add active state to navigation links based on scroll position
const sections = document.querySelectorAll("section[id]");
const navLinks = document.querySelectorAll(".nav-link");

function updateActiveLink() {
  const scrollPosition = window.scrollY + 150;

  sections.forEach((section) => {
    const sectionTop = section.offsetTop;
    const sectionHeight = section.offsetHeight;
    const sectionId = section.getAttribute("id");

    if (
      scrollPosition >= sectionTop &&
      scrollPosition < sectionTop + sectionHeight
    ) {
      navLinks.forEach((link) => {
        link.style.color = "";
        if (link.getAttribute("href") === `#${sectionId}`) {
          link.style.color = "var(--color-cyan-400)";
        }
      });
    }
  });
}

window.addEventListener("scroll", updateActiveLink);
updateActiveLink(); // Run on page load

// Navbar background on scroll
const navBar = document.querySelector(".nav-container");
let lastScroll = 0;

window.addEventListener("scroll", () => {
  const currentScroll = window.scrollY;

  if (currentScroll > 50) {
    navBar.style.background = "rgba(15, 23, 42, 0.95)";
    navBar.style.borderColor = "rgba(51, 65, 85, 0.8)";
  } else {
    navBar.style.background = "rgba(15, 23, 42, 0.8)";
    navBar.style.borderColor = "rgba(51, 65, 85, 0.5)";
  }

  lastScroll = currentScroll;
});

// Intersection Observer for fade-in animations
const observerOptions = {
  threshold: 0.1,
  rootMargin: "0px 0px -50px 0px",
};

const observer = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      entry.target.style.opacity = "1";
      entry.target.style.transform = "translateY(0)";
    }
  });
}, observerOptions);

// Observe cards for animation
const cards = document.querySelectorAll(
  ".hub-card, .download-card, .resource-card"
);
cards.forEach((card) => {
  card.style.opacity = "0";
  card.style.transform = "translateY(20px)";
  card.style.transition = "opacity 0.6s ease, transform 0.6s ease";
  observer.observe(card);
});

// Log console message
console.log(
  "%cSynapSeq",
  "color: #60a5fa; font-size: 24px; font-weight: bold;"
);
console.log(
  "%cGuiding brainwave states through sound",
  "color: #94a3b8; font-size: 14px;"
);
console.log(
  "%c🧠 Explore the code: https://github.com/ruanklein/synapseq",
  "color: #22d3ee; font-size: 12px;"
);
