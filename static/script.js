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

  const updateAriaLabel = () => {
    const next = getEffectiveTheme() === "dark" ? "light" : "dark";
    toggle.setAttribute("aria-label", `Switch to ${next} mode`);
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
})();
