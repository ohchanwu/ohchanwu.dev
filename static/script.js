Array.from(document.getElementsByClassName("year")).forEach((el) => {
  el.textContent = new Date().getFullYear();
});

(function () {
  const root = document.documentElement;
  const toggle = document.getElementById("theme-toggle");
  if (!toggle) return;

  const getEffectiveTheme = () => {
    const explicit = root.getAttribute("data-theme");
    if (explicit === "light" || explicit === "dark") return explicit;
    return window.matchMedia("(prefers-color-scheme: light)").matches
      ? "light"
      : "dark";
  };

  const getEffectiveLang = () => {
    const explicit = root.getAttribute("data-lang");
    return explicit === "ko" ? "ko" : "en";
  };

  const updateAriaLabel = () => {
    const next = getEffectiveTheme() === "dark" ? "light" : "dark";
    const lang = getEffectiveLang();
    let label;
    if (lang === "ko") {
      label = next === "light" ? "라이트 모드로 변경" : "다크 모드로 변경";
    } else {
      label = `Switch to ${next} mode`;
    }
    toggle.setAttribute("aria-label", label);
  };

  updateAriaLabel();

  toggle.addEventListener("click", () => {
    const next = getEffectiveTheme() === "dark" ? "light" : "dark";

    root.classList.add("theme-transition");
    root.setAttribute("data-theme", next);
    try {
      localStorage.setItem("theme", next);
    } catch (e) {}

    updateAriaLabel();

    window.setTimeout(() => {
      root.classList.remove("theme-transition");
    }, 250);
  });

  window
    .matchMedia("(prefers-color-scheme: light)")
    .addEventListener("change", () => {
      if (!root.hasAttribute("data-theme")) updateAriaLabel();
    });

  document.addEventListener("langchange", updateAriaLabel);
})();

(function () {
  const root = document.documentElement;
  const toggle = document.getElementById("lang-toggle");
  if (!toggle) return;

  const getEffectiveLang = () => {
    const explicit = root.getAttribute("data-lang");
    return explicit === "ko" ? "ko" : "en";
  };

  const updateAriaLabel = () => {
    toggle.setAttribute(
      "aria-label",
      getEffectiveLang() === "en" ? "언어를 한국어로 변경" : "Switch to English",
    );
  };

  updateAriaLabel();

  toggle.addEventListener("click", () => {
    const next = getEffectiveLang() === "en" ? "ko" : "en";

    root.setAttribute("data-lang", next);
    root.setAttribute("lang", next);
    try {
      localStorage.setItem("lang", next);
    } catch (e) {}

    const titleEl = document.querySelector("title");
    const descEl = document.querySelector('meta[name="description"]');
    const key = next === "ko" ? "i18nKo" : "i18nEn";
    if (titleEl && titleEl.dataset[key]) {
      titleEl.textContent = titleEl.dataset[key];
    }
    if (descEl && descEl.dataset[key]) {
      descEl.setAttribute("content", descEl.dataset[key]);
    }

    updateAriaLabel();
    document.dispatchEvent(new CustomEvent("langchange"));
  });
})();
