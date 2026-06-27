Array.from(document.getElementsByClassName("year")).forEach((el) => {
  el.textContent = new Date().getFullYear();
});

/* Native scroll restoration runs before fonts and media settle, then scroll
   anchoring can re-save a slightly shifted value on each reload. Take over so
   reload/back-forward restore the exact saved position after layout settles. */
(function restoreScroll() {
  if (!("scrollRestoration" in history)) return;
  history.scrollRestoration = "manual";

  const key = `ohchanwu:scroll:${location.pathname}${location.search}`;

  const save = () => {
    try {
      sessionStorage.setItem(key, String(Math.round(window.pageYOffset)));
    } catch (e) {}
  };

  window.addEventListener("pagehide", save);
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") save();
  });

  const entry =
    (performance.getEntriesByType &&
      performance.getEntriesByType("navigation")[0]) ||
    null;
  const navType = entry
    ? entry.type
    : performance.navigation && performance.navigation.type === 1
      ? "reload"
      : "navigate";
  if (navType !== "reload" && navType !== "back_forward") return;

  let saved = 0;
  try {
    saved = parseInt(sessionStorage.getItem(key), 10) || 0;
  } catch (e) {}
  if (saved <= 0) return;

  const apply = () => {
    window.scrollTo(0, saved);
  };

  apply();
  window.addEventListener("DOMContentLoaded", apply);
  window.addEventListener("load", () => {
    requestAnimationFrame(apply);
    window.setTimeout(apply, 800);
  });
  if (document.fonts && document.fonts.ready) {
    document.fonts.ready.then(() => requestAnimationFrame(apply));
  }
})();

(function () {
  function init() {
    document.documentElement.classList.add("media-fade-enabled");

    const reveal = (el) => {
      el.classList.add("is-media-loaded");
    };
    const fail = (el) => {
      el.classList.add("is-media-error");
    };

    const settleImage = (img) => {
      const show = () => {
        if (img.decode) {
          img.decode().then(
            () => reveal(img),
            () => reveal(img),
          );
        } else {
          reveal(img);
        }
      };

      if (img.complete && img.naturalWidth > 0) {
        show();
        return;
      }
      img.addEventListener("load", show, { once: true });
      img.addEventListener("error", () => fail(img), { once: true });
    };

    const settleVideo = (video) => {
      if (video.readyState >= 2) {
        reveal(video);
        return;
      }
      video.addEventListener("loadeddata", () => reveal(video), { once: true });
      video.addEventListener("error", () => fail(video), { once: true });
    };

    document.querySelectorAll(".content-media").forEach((media) => {
      if (media.tagName === "IMG") {
        settleImage(media);
      } else if (media.tagName === "VIDEO") {
        settleVideo(media);
      }
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();

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
