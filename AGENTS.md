# ohchanwu.dev

Personal portfolio site for Chanwu (Tyler) Oh.

## Stack

- Plain HTML, CSS, and browser JavaScript under `static/`
- GitHub Actions and GitHub Pages for deployment
- Cloudflare for DNS only

## Conventions

- Keep the site framework-free and dependency-free unless the owner approves a real need.
- Preserve English and Korean content and test both language modes after frontend changes.
- Keep site asset URLs relative to the `static/` publishing root.
- Preserve keyboard access, reduced-motion behavior, theme behavior, scroll restoration, and the
  JavaScript-gated media fade.
- After frontend changes, start or reuse the local preview and run screenshot-led browser QA on
  desktop and mobile. Check both themes, both languages, console errors, first-party network
  failures, and representative links.
- Do not add AI-voice marketing prose to the site or READMEs.

## Workflow

- Work on a feature branch and keep commits focused.
- Run local checks before each commit and inspect the cumulative diff before handoff.
- Never push or change live deployment settings unless the owner explicitly asks.
- Site copy facts and the resume PDFs in `static/` come from the owner's jss-genr harness
  (`~/projects/jss-genr`). Refresh the PDFs with `make sync-site` in `jss-genr/resume/`
  instead of editing them here.

## Project structure

- `.github/workflows/pages.yml`: publishes only `static/` to GitHub Pages.
- `static/index.html`: site markup and bilingual content.
- `static/style.css`: responsive layout, themes, and presentation.
- `static/script.js`: language/theme controls, scroll restoration, and media loading behavior.
- `static/`: public images, icons, and resume PDFs.
- `docs/`: maintained architecture, plans, decisions, and selected historical records.

## Local preview

Run:

```sh
python3 -m http.server 4173 --bind 127.0.0.1 --directory static
```

Inspect `http://127.0.0.1:4173/` with the prescribed headless browser tooling. Do not open the
user's default browser.

## Documentation lifecycle

- Use `docs/README.md` as the documentation index.
- Keep active implementation plans and specifications under `docs/superpowers/`.
- Move useful completed records to the tracked archive only after the implementation and its
  production verification are complete.
- Keep raw local reports, screenshots, logs, diffs, and handoff evidence under ignored
  `.superpowers/` paths.
- Treat every tracked file as public. Never commit credentials, private keys, cookies, personal
  exports, or unnecessary production identifiers.

## What production-ready means here

- The Pages workflow publishes only `static/`.
- The page works at desktop and mobile widths in English and Korean, in light and dark themes.
- Local assets load without first-party network failures and the browser console stays clean.
- Both resume PDFs and representative external links reach the expected content.
- The custom domain serves over HTTPS after the owner-operated DNS cutover.
