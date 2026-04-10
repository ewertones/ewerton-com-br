(function () {
  "use strict";

  // ======================== THEME TOGGLE ========================
  var themeToggleBtn = document.getElementById("themeToggleBtn");
  var sunIcon = document.getElementById("iconSun");
  var moonIcon = document.getElementById("iconMoon");

  function applyTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
    if (sunIcon && moonIcon) {
      if (theme === "dark") {
        sunIcon.style.display = "block";
        moonIcon.style.display = "none";
      } else {
        sunIcon.style.display = "none";
        moonIcon.style.display = "block";
      }
    }
  }

  // Sync icon state on load
  (function () {
    var saved = localStorage.getItem("theme");
    var theme = saved || (window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light");
    applyTheme(theme);
  })();

  if (themeToggleBtn) {
    themeToggleBtn.addEventListener("click", function () {
      var current = document.documentElement.getAttribute("data-theme");
      applyTheme(current === "dark" ? "light" : "dark");
    });
  }

  // ======================== LANGUAGE SWITCHER ========================
  var langSwitcher = document.getElementById("langSwitcher");
  var langBtn = document.getElementById("langBtn");
  var langDropdown = document.getElementById("langDropdown");
  var langChevron = document.getElementById("langChevron");

  if (langBtn && langDropdown) {
    langBtn.addEventListener("click", function (e) {
      e.stopPropagation();
      var isOpen = langDropdown.classList.contains("open");
      if (isOpen) {
        langDropdown.classList.remove("open");
        if (langChevron) langChevron.classList.remove("open");
      } else {
        langDropdown.classList.add("open");
        if (langChevron) langChevron.classList.add("open");
      }
    });

    document.addEventListener("click", function () {
      langDropdown.classList.remove("open");
      if (langChevron) langChevron.classList.remove("open");
    });

    langDropdown.querySelectorAll(".langOption").forEach(function (opt) {
      opt.addEventListener("click", function () {
        var locale = opt.getAttribute("data-locale");
        if (locale) {
          window.location.href = "/" + locale;
        }
      });
    });
  }

  // ======================== TYPEWRITER ========================
  var HELLO_WORDS = [
    "Hello",
    "Olá",
    "Hola",
    "Bonjour",
    "Ciao",
    "Hej",
    "Hallo",
    "Salut",
    "Merhaba",
    "こんにちは",
    "Привет",
    "مرحبا",
  ];

  var helloWordEl = document.getElementById("helloWord");

  if (helloWordEl) {
    var helloIndex = 0;
    var displayedHello = HELLO_WORDS[0];
    var isDeleting = false;
    var tw_timer = null;

    function typewriterTick() {
      var currentWord = HELLO_WORDS[helloIndex];

      if (!isDeleting && displayedHello === currentWord) {
        // Fully typed — pause then start deleting
        tw_timer = setTimeout(function () {
          isDeleting = true;
          typewriterTick();
        }, 2000);
        return;
      }

      if (isDeleting && displayedHello === "") {
        // Fully deleted — move to next word
        tw_timer = setTimeout(function () {
          isDeleting = false;
          helloIndex = (helloIndex + 1) % HELLO_WORDS.length;
          typewriterTick();
        }, 200);
        return;
      }

      var delay = isDeleting ? 30 : 60;
      tw_timer = setTimeout(function () {
        if (isDeleting) {
          displayedHello = displayedHello.slice(0, -1);
        } else {
          displayedHello = currentWord.slice(0, displayedHello.length + 1);
        }
        helloWordEl.textContent = displayedHello;
        typewriterTick();
      }, delay);
    }

    typewriterTick();
  }

  // ======================== SKILL SEARCH ========================
  var techSearch = document.getElementById("techSearch");
  var techGrid = document.getElementById("techGrid");

  if (techSearch && techGrid) {
    techSearch.addEventListener("input", function () {
      var query = techSearch.value.toLowerCase().trim();
      var categories = techGrid.querySelectorAll(".techCategory");

      categories.forEach(function (cat) {
        var catName = (cat.getAttribute("data-category-name") || "").toLowerCase();
        var cards = cat.querySelectorAll(".techSkillCard");
        var anyVisible = false;

        if (!query) {
          cards.forEach(function (c) {
            c.style.display = "";
          });
          cat.style.display = "";
          return;
        }

        // If category name matches, show all skills
        if (catName.includes(query)) {
          cards.forEach(function (c) {
            c.style.display = "";
          });
          cat.style.display = "";
          return;
        }

        cards.forEach(function (card) {
          var name = (card.getAttribute("data-name") || "").toLowerCase();
          var keywords = (card.getAttribute("data-keywords") || "").toLowerCase();
          var match = name.includes(query) || keywords.includes(query);
          card.style.display = match ? "" : "none";
          if (match) anyVisible = true;
        });

        cat.style.display = anyVisible ? "" : "none";
      });
    });
  }

  // ======================== PDF VIEWER MODAL ========================
  var pdfOverlay = document.getElementById("pdfViewerOverlay");
  var pdfFrame = document.getElementById("pdfViewerFrame");
  var viewResumeBtn = document.getElementById("viewResumeBtn");
  var RESUME_PATH = "/static/documents/resume-ewerton.pdf";

  function openPdfViewer() {
    if (!pdfOverlay || !pdfFrame) return;
    pdfFrame.src = RESUME_PATH;
    pdfOverlay.classList.add("open");
    document.body.style.overflow = "hidden";
  }

  function closePdfViewer(e) {
    if (e && e.target !== pdfOverlay) return;
    if (!pdfOverlay || !pdfFrame) return;
    pdfOverlay.classList.remove("open");
    pdfFrame.src = "";
    document.body.style.overflow = "";
  }

  // expose globally for onclick="" in template
  window.closePdfViewer = closePdfViewer;

  if (viewResumeBtn) {
    viewResumeBtn.addEventListener("click", openPdfViewer);
  }

  // ======================== COMPANY MODAL ========================
  var companyOverlay = document.getElementById("companyModalOverlay");

  function openCompanyModal(slug) {
    if (!companyOverlay) return;
    // Hide all company sections
    companyOverlay.querySelectorAll(".companyRoleSection").forEach(function (el) {
      el.style.display = "none";
    });
    // Show matching section
    var section = document.getElementById("modal-" + slug);
    if (section) section.style.display = "block";
    companyOverlay.classList.add("open");
    document.body.style.overflow = "hidden";
  }

  function closeCompanyModal(e) {
    if (e && e.target !== companyOverlay) return;
    if (!companyOverlay) return;
    companyOverlay.classList.remove("open");
    document.body.style.overflow = "";
  }

  // expose globally
  window.openCompanyModal = openCompanyModal;
  window.closeCompanyModal = closeCompanyModal;

  // Escape key closes any open modal
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") {
      if (pdfOverlay && pdfOverlay.classList.contains("open")) {
        pdfOverlay.classList.remove("open");
        if (pdfFrame) pdfFrame.src = "";
        document.body.style.overflow = "";
      }
      if (companyOverlay && companyOverlay.classList.contains("open")) {
        companyOverlay.classList.remove("open");
        document.body.style.overflow = "";
      }
    }
  });
})();
