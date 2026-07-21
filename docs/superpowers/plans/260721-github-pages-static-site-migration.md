# GitHub Pages Static Site Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development
> (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use
> checkbox (`- [ ]`) syntax for tracking.

**Goal:** Convert the existing `ohchanwu/ohchanwu.dev` repository from a Go/EC2 service into a
dependency-free static site deployed from `static/` by GitHub Pages, without changing the site's
design beyond removing the footer technology line.

**Architecture:** GitHub Actions runs after changes land on `main`, uploads only `static/` as the
Pages artifact, and deploys it through GitHub's `github-pages` environment. GitHub Pages terminates
HTTPS. Cloudflare remains the DNS provider in DNS-only mode. The repository owner performs GitHub
Settings, DNS, and former-server shutdown operations after the local changes are reviewed.

**Tech Stack:** HTML, CSS, browser JavaScript, GitHub Actions, GitHub Pages, Cloudflare DNS, Python
standard-library preview server, gstack `/browse`, Gitleaks.

## Global Constraints

- Work only from the repository root on `chore/github-pages` unless the repository owner
  explicitly changes the branch.
- Keep the existing repository and remote. Do not create another repository, deployment branch,
  package manifest, generated site directory, `CNAME`, `.nojekyll`, or custom `404.html`.
- Keep the deployable artifact limited to `static/`. Repository docs and metadata must not be
  published as site content.
- Preserve the current design, bilingual behavior, theme behavior, scroll restoration, and media
  fade. The only intended visual change is removal of the two technology/copyright footer lines.
- Use Ponytail judgment throughout: delete dead runtime material and dead footer code; add no build
  tool or abstraction for this small static site.
- Never push, create a pull request, change GitHub Settings or DNS, or stop EC2 from this plan's
  local implementation phase. Commit locally at the named checkpoints.
- Run at most two or three frontend verification agents/browser sessions at once. One `/browse`
  session is sufficient for this plan.
- Keep screenshots and raw browser evidence under `/tmp`, not in Git.
- Before every documentation commit, inspect the complete staged diff, run Gitleaks with redaction,
  and manually check for credentials, personal data, and unnecessary production details.
- If a user-facing check fails or a real product decision appears, stop before committing and report
  it. Do not weaken the check to obtain a pass.

---

## Task 1: Remove the obsolete footer copy and its dead JavaScript

**Files:**

- Modify: `static/index.html:514`
- Modify: `static/script.js:1`

- [ ] **Step 1: Confirm the intended branch, clean tree, and guidance symlink**

  Run:

  ```sh
  git status --short --branch
  git branch --show-current
  test -L CLAUDE.md
  test "$(readlink CLAUDE.md)" = "AGENTS.md"
  cmp AGENTS.md CLAUDE.md
  ```

  Expected: branch `chore/github-pages`, no changed paths, and all symlink checks exit `0`. Stop if
  the tree is dirty because this plan assumes the approved spec/plan commit is the baseline.

- [ ] **Step 2: Run the failing footer-removal acceptance check**

  Run:

  ```sh
  ! rg -n 'built in go|Go로 개발|getElementsByClassName\("year"\)' \
    static/index.html static/script.js
  ```

  Expected before the edit: exit `1`, with matches for both footer lines and the year updater. This
  proves the acceptance check detects the old implementation.

- [ ] **Step 3: Remove only the obsolete footer paragraphs and year updater**

  In `static/index.html`, leave the shell-style exit line in place and delete these two paragraphs:

  ```html
  <p class="copy" lang="en">
    &copy; <span class="year"></span> &middot; built in go
  </p>
  <p class="copy" lang="ko">
    &copy; <span class="year"></span> &middot; Go로 개발
  </p>
  ```

  In `static/script.js`, delete the now-unreferenced opening block:

  ```js
  Array.from(document.getElementsByClassName("year")).forEach((el) => {
    el.textContent = new Date().getFullYear();
  });
  ```

  Do not alter the scroll-restoration, media-fade, language, or theme code.

- [ ] **Step 4: Run focused static checks**

  Run:

  ```sh
  ! rg -n 'built in go|Go로 개발|getElementsByClassName\("year"\)' \
    static/index.html static/script.js
  node --check static/script.js
  git diff --check -- static/index.html static/script.js
  git diff -- static/index.html static/script.js
  ```

  Expected: no obsolete matches, valid JavaScript, no whitespace errors, and a diff containing only
  the two footer paragraphs and the three-line year updater.

- [ ] **Step 5: Start or reuse the local static preview without opening the user's browser**

  Confirm port `4173` is free:

  ```sh
  lsof -nP -iTCP:4173 -sTCP:LISTEN
  ```

  Expected: no listener. Start this in a persistent terminal session and leave it running for human
  inspection after the task:

  ```sh
  python3 -m http.server 4173 --bind 127.0.0.1 --directory static
  ```

  Preview URL: `http://127.0.0.1:4173/`.

- [ ] **Step 6: Invoke `frontend-qa` and gstack `/browse`, then capture visual evidence**

  Read and follow both skills in full before browser work:

  ```text
  ~/.codex/skills/frontend-qa/SKILL.md
  ~/.agents/skills/gstack/browse/SKILL.md
  ```

  Run the `/browse` preamble and setup check from its skill. Use one headless session. If the browse
  binary reports `NEEDS_SETUP`, follow the skill's approval gate rather than opening another browser.

  Resolve `B` as specified by `/browse`, then run:

  ```sh
  $B goto http://127.0.0.1:4173/
  $B wait --load
  $B viewport 1440x900
  $B screenshot /tmp/ohchanwu-pages-desktop-light-en.png
  $B scroll footer
  $B screenshot /tmp/ohchanwu-pages-desktop-footer.png --selector footer
  $B click '#lang-toggle'
  $B screenshot /tmp/ohchanwu-pages-desktop-light-ko.png
  $B click '#theme-toggle'
  $B screenshot /tmp/ohchanwu-pages-desktop-dark-ko.png
  $B click '#lang-toggle'
  $B screenshot /tmp/ohchanwu-pages-desktop-dark-en.png
  $B click '#theme-toggle'
  $B viewport 390x844
  $B screenshot /tmp/ohchanwu-pages-mobile-light-en.png
  $B click '#lang-toggle'
  $B screenshot /tmp/ohchanwu-pages-mobile-light-ko.png
  $B click '#theme-toggle'
  $B screenshot /tmp/ohchanwu-pages-mobile-dark-ko.png
  $B click '#lang-toggle'
  $B screenshot /tmp/ohchanwu-pages-mobile-dark-en.png
  ```

  Open every PNG with the available image viewer. Verify the footer contains only the exit prompt;
  there is no overlap, clipping, unexplained blank band, horizontal overflow, or broken media; the
  portrait crop remains intentional; and language/theme controls remain legible and usable.

- [ ] **Step 7: Run browser behavior, DOM, and user-path checks**

  Return to desktop and run:

  ```sh
  $B viewport 1440x900
  $B reload
  $B js "JSON.stringify({overflow: document.documentElement.scrollWidth <= window.innerWidth + 2, footerCopy: document.querySelectorAll('footer .copy').length, fadeEnabled: document.documentElement.classList.contains('media-fade-enabled'), images: [...document.images].map(i => ({src: i.getAttribute('src'), complete: i.complete, width: i.naturalWidth, loaded: i.classList.contains('is-media-loaded'), error: i.classList.contains('is-media-error')}))})"
  $B console --errors
  $B network
  ```

  Expected: `overflow` and `fadeEnabled` are `true`, `footerCopy` is `0`, the portrait has a positive
  natural width plus its loaded class and no error class, no console errors, and no failed
  first-party requests. If Google Fonts is unreachable, record that external limitation explicitly;
  do not misreport it as a first-party pass.

  Exercise reload scroll restoration through the user path:

  ```sh
  $B js "window.scrollTo(0, 800); Math.round(window.scrollY)"
  $B reload
  $B wait --load
  $B js "Math.round(window.scrollY)"
  ```

  Expected: the pre- and post-reload values match within two pixels.

  From the page itself, click and inspect both resume links. Then click each unique representative
  external destination below, switch to the opened tab, and confirm the visible destination content
  matches the clicked item rather than merely returning HTTP `200`:

  - `resume-en.pdf`: browser PDF viewer shows the English resume.
  - `resume-kr.pdf`: browser PDF viewer shows the Korean resume.
  - `https://github.com/ohchanwu`: profile identifies `ohchanwu`.
  - `https://github.com/ohchanwu/job-scraper`: destination identifies the intended project, including
    any legitimate GitHub rename redirect.
  - `https://www.chaphysicaltherapy.com/`: destination identifies Cha Physical Therapy.
  - `https://ohchanwu.github.io/alix-chaput-architecture/`: destination identifies the architecture
    project.
  - `https://github.com/freeCodeCamp/freeCodeCamp/pull/66981`: destination identifies pull request
    `#66981` in `freeCodeCamp/freeCodeCamp`.

  Re-run `$B console --errors` and `$B network` on the local site after returning to its tab. Keep all
  screenshots under `/tmp`.

- [ ] **Step 8: Commit the verified frontend reduction**

  Run:

  ```sh
  git add static/index.html static/script.js
  git diff --cached --check
  git diff --cached -- static/index.html static/script.js
  git commit -m "refactor: remove server-specific footer copy"
  ```

  Expected: one local commit containing only the two verified frontend edits. Do not push.

---

## Task 2: Replace the Go/Docker runtime with a static Pages workflow

**Files:**

- Create: `.github/workflows/pages.yml`
- Modify: `.gitignore`
- Delete: `.dockerignore`
- Delete: `Dockerfile`
- Delete: `go.mod`
- Delete: `go.sum`
- Delete: `main.go`
- Delete: `main_test.go`

- [ ] **Step 1: Run the failing static-runtime acceptance gate**

  Run:

  ```sh
  test -f .github/workflows/pages.yml && \
    test ! -e .dockerignore && \
    test ! -e Dockerfile && \
    test ! -e go.mod && \
    test ! -e go.sum && \
    test ! -e main.go && \
    test ! -e main_test.go
  ```

  Expected before the edit: exit `1` because the workflow is absent and the old runtime files exist.

- [ ] **Step 2: Add the minimal official Pages workflow**

  Create `.github/workflows/pages.yml` with exactly:

  ```yaml
  name: Deploy static site to Pages

  on:
    push:
      branches: ["main"]
    workflow_dispatch:

  permissions:
    contents: read
    pages: write
    id-token: write

  concurrency:
    group: "pages"
    cancel-in-progress: false

  jobs:
    deploy:
      environment:
        name: github-pages
        url: ${{ steps.deployment.outputs.page_url }}
      runs-on: ubuntu-latest
      steps:
        - name: Checkout
          uses: actions/checkout@v4
        - name: Setup Pages
          uses: actions/configure-pages@v5
        - name: Upload static site
          uses: actions/upload-pages-artifact@v3
          with:
            path: static
        - name: Deploy to GitHub Pages
          id: deployment
          uses: actions/deploy-pages@v5
  ```

  This is intentionally a direct artifact upload. Do not add a build job, cache, package install, or
  generated branch.

- [ ] **Step 3: Delete the retired runtime and simplify ignores**

  Delete `.dockerignore`, `Dockerfile`, `go.mod`, `go.sum`, `main.go`, and `main_test.go` with the
  patch tool so the deletions remain reviewable.

  Replace `.gitignore` with:

  ```gitignore
  # Local environment and credentials
  .env
  .env.local
  *.pem
  *.key

  # Editor and operating-system files
  .DS_Store
  .idea/
  .vscode/

  # Local agent and Gas Town state
  .beads/
  .runtime/
  .claude/
  .gstack/
  CLAUDE.local.md

  # Temporary implementation evidence
  .superpowers/sdd/
  .superpowers/archive/
  ```

  Do not delete any historical commit or rewrite history; Git remains the recovery mechanism for the
  retired server.

- [ ] **Step 4: Validate YAML syntax and the Pages artifact boundary**

  Run:

  ```sh
  ruby -e 'require "yaml"; YAML.safe_load(File.read(".github/workflows/pages.yml"), permitted_classes: [], permitted_symbols: [], aliases: true); puts "workflow YAML parsed"'
  rg -n 'actions/checkout@v4|actions/configure-pages@v5|actions/upload-pages-artifact@v3|actions/deploy-pages@v5|path: static' \
    .github/workflows/pages.yml
  test "$(rg -c '^          path: static$' .github/workflows/pages.yml)" = "1"
  ! rg -n 'CNAME|\.nojekyll|404\.html|npm|go build|docker' .github/workflows/pages.yml
  ```

  Expected: `workflow YAML parsed`; one line for each required action and `path: static`; exactly one
  upload path; and no excluded deployment mechanisms.

- [ ] **Step 5: Validate the static-only active tree**

  Run:

  ```sh
  test -f .github/workflows/pages.yml
  test ! -e .dockerignore
  test ! -e Dockerfile
  test ! -e go.mod
  test ! -e go.sum
  test ! -e main.go
  test ! -e main_test.go
  test -f static/index.html
  test -f static/style.css
  test -f static/script.js
  test -f static/resume-en.pdf
  test -f static/resume-kr.pdf
  test -L CLAUDE.md
  test "$(readlink CLAUDE.md)" = "AGENTS.md"
  git diff --check
  ```

  Expected: all commands exit `0`. Reload `http://127.0.0.1:4173/` once through `/browse` and confirm
  it still renders, because the user-visible static files have not moved.

- [ ] **Step 6: Review and commit the deployment boundary**

  Run:

  ```sh
  git add .github/workflows/pages.yml .gitignore .dockerignore Dockerfile \
    go.mod go.sum main.go main_test.go
  git diff --cached --check
  git diff --cached --stat
  git diff --cached -- .github/workflows/pages.yml .gitignore .dockerignore Dockerfile \
    go.mod go.sum main.go main_test.go
  git commit -m "chore: deploy portfolio with GitHub Pages"
  ```

  Expected: one local commit adding the workflow, simplifying ignores, and deleting the retired
  runtime. Do not push.

---

## Task 3: Rewrite active documentation for the static deployment model

**Files:**

- Modify: `README.md`
- Modify: `README.ko.md`
- Modify: `AGENTS.md`
- Verify only: `CLAUDE.md`
- Modify: `docs/README.md`

- [ ] **Step 1: Run the failing stale-documentation check**

  Run:

  ```sh
  ! rg -n 'autocert|ACME_|Docker|EC2|main\.go|http\.FileServer|single binary|단일 바이너리' \
    README.md README.ko.md AGENTS.md docs/README.md
  ```

  Expected before the rewrite: exit `1`, with matches describing the retired runtime.

- [ ] **Step 2: Replace the English README with the static-site instructions**

  Replace `README.md` with:

  ````markdown
  # ohchanwu.dev

  _[한국어로 보기 🇰🇷](README.ko.md)_

  Personal portfolio site for Chanwu (Tyler) Oh.

  The site is plain HTML, CSS, and JavaScript published from `static/` by GitHub Pages.

  Project documentation is indexed in [docs/README.md](docs/README.md).

  ## Stack

  - HTML
  - CSS
  - JavaScript
  - GitHub Actions
  - GitHub Pages

  ## Local preview

  ```sh
  python3 -m http.server 4173 --bind 127.0.0.1 --directory static
  ```

  Open `http://127.0.0.1:4173/`.

  ## Deployment

  Pushes to `main` deploy `static/` through `.github/workflows/pages.yml`.

  The repository owner manages the GitHub Pages custom-domain setting and Cloudflare DNS records.

  ## License

  [MIT](LICENSE)
  ````

- [ ] **Step 3: Replace the Korean README with equivalent factual guidance**

  Replace `README.ko.md` with:

  ````markdown
  # ohchanwu.dev

  _[Read this in English 🇬🇧](README.md)_

  오찬우(Chanwu, Tyler Oh)의 개인 포트폴리오 사이트입니다.

  `static/`의 순수 HTML, CSS, JavaScript를 GitHub Pages로 배포합니다.

  프로젝트 문서는 [docs/README.md](docs/README.md)에 정리되어 있습니다.

  ## 기술 스택

  - HTML
  - CSS
  - JavaScript
  - GitHub Actions
  - GitHub Pages

  ## 로컬 미리 보기

  ```sh
  python3 -m http.server 4173 --bind 127.0.0.1 --directory static
  ```

  `http://127.0.0.1:4173/`을 여십시오.

  ## 배포

  `main`에 푸시된 변경 사항은 `.github/workflows/pages.yml`을 통해 `static/`을 배포합니다.

  저장소 소유자가 GitHub Pages 사용자 지정 도메인 설정과 Cloudflare DNS 레코드를 관리합니다.

  ## 라이선스

  [MIT](LICENSE)
  ````

- [ ] **Step 4: Replace repository guidance with the static architecture**

  Replace `AGENTS.md` with:

  ```markdown
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
  ```

  Leave `CLAUDE.md` as the relative symlink to `AGENTS.md`; do not replace it with a copied file.

- [ ] **Step 5: Update the documentation index**

  Replace `docs/README.md` with:

  ```markdown
  # Documentation

  This directory indexes maintained project documentation. Treat all tracked content as public.

  ## Contents

  - [Implementation knowledge](superpowers/README.md): active plans and specifications, durable
    decisions, and selected completed records.
  - Root [README](../README.md) and [Korean README](../README.ko.md): project overview, local preview,
    and deployment guidance.
  - Root [AGENTS.md](../AGENTS.md): canonical repository and agent guidance.
  - Root [CLAUDE.md](../CLAUDE.md): relative symlink to `AGENTS.md` for compatible tooling.

  Executable and conventional project material stays at the repository root: the Pages workflow,
  `static/`, localized READMEs, `AGENTS.md`, its `CLAUDE.md` symlink, and `LICENSE`.

  ## Publication note

  Tracked documentation must not contain credentials, private keys, cookies, personal exports, or
  unnecessary production-specific identifiers. Temporary raw implementation evidence belongs in
  ignored `.superpowers/` paths and must still be sanitized.
  ```

- [ ] **Step 6: Run documentation, symlink, and publication checks**

  Run:

  ```sh
  ! rg -n 'autocert|ACME_|Docker|EC2|main\.go|http\.FileServer|single binary|단일 바이너리' \
    README.md README.ko.md AGENTS.md docs/README.md
  rg -n 'GitHub Pages|static/' README.md README.ko.md AGENTS.md docs/README.md
  test -e docs/README.md
  test -e docs/superpowers/README.md
  test -L CLAUDE.md
  test "$(readlink CLAUDE.md)" = "AGENTS.md"
  cmp AGENTS.md CLAUDE.md
  git diff --check
  ```

  Expected: no retired-runtime terms in active docs; static/Pages guidance in every maintained file;
  a valid relative guidance symlink; and no whitespace errors.

  Stage, inspect, and scan before committing:

  ```sh
  git add README.md README.ko.md AGENTS.md CLAUDE.md docs/README.md
  git diff --cached --stat
  git diff --cached
  gitleaks git --staged --redact=100 --no-banner .
  ```

  Manually confirm the staged diff contains no credentials, personal/healthcare data, private
  topology, raw logs, or unnecessary infrastructure identifiers. Gitleaks must report no leaks.

- [ ] **Step 7: Commit the maintained documentation**

  Run:

  ```sh
  git commit -m "docs: document GitHub Pages hosting"
  ```

  Expected: one local documentation commit. Do not push.

---

## Task 4: Run the cumulative local verification gate and hand off

**Files:**

- Verify: all files changed since this implementation plan was committed
- Keep running: local preview at `http://127.0.0.1:4173/`

- [ ] **Step 1: Re-read the complete implementation diff**

  Resolve the baseline from the commit that introduced this plan:

  ```sh
  BASE="$(git log -1 --format=%H --fixed-strings \
    --grep='docs: add GitHub Pages migration plan')"
  test -n "$BASE"
  git diff --check "$BASE"..HEAD
  git diff --stat "$BASE"..HEAD
  git diff "$BASE"..HEAD
  ```

  Expected: exactly the frontend reduction, Pages workflow, runtime deletions, ignore cleanup, and
  maintained-doc rewrites specified above. Investigate any unrelated file.

- [ ] **Step 2: Re-run the complete local machine-verifiable suite**

  Run:

  ```sh
  BASE="$(git log -1 --format=%H --fixed-strings \
    --grep='docs: add GitHub Pages migration plan')"
  test -n "$BASE"
  ruby -e 'require "yaml"; YAML.safe_load(File.read(".github/workflows/pages.yml"), permitted_classes: [], permitted_symbols: [], aliases: true); puts "workflow YAML parsed"'
  node --check static/script.js
  test -f static/index.html
  test -f static/style.css
  test -f static/script.js
  test ! -e main.go
  test ! -e main_test.go
  test ! -e go.mod
  test ! -e go.sum
  test ! -e Dockerfile
  test ! -e .dockerignore
  test -L CLAUDE.md
  cmp AGENTS.md CLAUDE.md
  ! rg -n 'built in go|Go로 개발|getElementsByClassName\("year"\)' \
    static/index.html static/script.js
  ! rg -n 'autocert|ACME_|Docker|EC2|main\.go|http\.FileServer|single binary|단일 바이너리' \
    README.md README.ko.md AGENTS.md docs/README.md
  gitleaks git --log-opts="$BASE..HEAD" --redact=100 --no-banner .
  git status --short --branch
  ```

  Expected: every command passes, Gitleaks reports no leaks, and the worktree is clean.

- [ ] **Step 3: Re-run the final browser smoke check**

  With the same single `/browse` session and live preview:

  ```sh
  $B goto http://127.0.0.1:4173/
  $B viewport 1440x900
  $B reload
  $B screenshot /tmp/ohchanwu-pages-final-desktop.png
  $B viewport 390x844
  $B screenshot /tmp/ohchanwu-pages-final-mobile.png
  $B console --errors
  $B network
  ```

  Inspect both images. Reconfirm no footer technology line, no horizontal overflow, and no failed
  first-party requests. Keep the preview server running and include its URL in the handoff report.

- [ ] **Step 4: Report the local result and the production gate**

  Report:

  - Local commit SHAs and subjects.
  - Preview URL `http://127.0.0.1:4173/`.
  - Desktop/mobile, English/Korean, light/dark, reload, media, resume, external-link, console, and
    network evidence.
  - Confirmation that nothing was pushed and no live setting changed.
  - The owner checklist in Task 5 and the fact that EC2 must remain running until custom-domain QA
    passes.

---

## Task 5: Owner-operated Pages and DNS cutover

**Files:**

- No local code changes until production verification succeeds
- Later move: `docs/superpowers/specs/260721-github-pages-static-site-migration.md`
- Later move: `docs/superpowers/plans/260721-github-pages-static-site-migration.md`
- Later modify: `docs/superpowers/README.md`

This task is intentionally gated on repository-owner actions. The implementation agent must stop
after Task 4 unless the owner explicitly asks it to assist with these live operations.

- [ ] **Step 1: Owner enables GitHub Actions as the Pages source and lands the branch**

  In repository **Settings → Pages**, select **GitHub Actions** as the build and deployment source.
  Review and merge or push the local commits to `main`. Confirm the `Deploy static site to Pages`
  workflow succeeds and record its generated Pages URL.

- [ ] **Step 2: Verify the generated Pages URL before touching DNS**

  Run the Task 1 browser matrix against the generated Pages URL: desktop/mobile, both languages,
  both themes, reload restoration, media, both PDFs, representative external links, console, and
  first-party network checks. Do not proceed on a partial pass.

- [ ] **Step 3: Owner verifies the domain and configures the custom domain**

  Verify `ohchanwu.dev` in the GitHub account using the DNS record GitHub requests. Then set
  `ohchanwu.dev` as this repository's Pages custom domain before repointing public traffic. Do not add
  a repository `CNAME` file; custom Actions deployments use the repository setting.

- [ ] **Step 4: Owner changes Cloudflare DNS using GitHub's current official values**

  At cutover time, re-check GitHub's official custom-domain documentation rather than copying stale
  IP addresses from this plan. Replace the old apex origin records with GitHub's then-current
  documented apex records, set `www` as a CNAME to `ohchanwu.github.io`, and keep every Pages record
  **DNS only** in Cloudflare. Preserve the old origin values in the private operator record for rapid
  rollback; do not commit them.

- [ ] **Step 5: Verify both public hostnames and HTTPS before retiring EC2**

  Wait for DNS propagation and GitHub certificate issuance, then enable **Enforce HTTPS**. Through
  `/browse`, verify:

  - `https://ohchanwu.dev/` serves the expected portfolio.
  - `https://www.ohchanwu.dev/` reaches the expected canonical hostname.
  - Both languages, both themes, local assets, both PDFs, and representative links still work.
  - There are no console errors, first-party request failures, certificate warnings, or mixed content.

  If certificate issuance is pending, keep the old origin available and wait. Do not weaken HTTPS or
  restore manual ACME handling.

- [ ] **Step 6: Owner retires the former server only after the full pass**

  Once Task 5 Step 5 passes, the owner may stop and remove the former EC2 workload. Before that point,
  rollback is a DNS change to the privately recorded old origin. After shutdown, repository rollback
  is a revert and redeploy; recreating the old server additionally requires Git history and new
  infrastructure.

- [ ] **Step 7: Archive the completed migration records**

  Only after production verification, read `docs/superpowers/archive/AGENTS.md`, then move the
  completed spec and plan into a dated tracked archive location and update
  `docs/superpowers/README.md` so no completed migration remains under active specs/plans. Inspect the
  complete staged documentation diff, run:

  ```sh
  gitleaks git --staged --redact=100 --no-banner .
  ```

  Commit the archival update locally with:

  ```sh
  git commit -m "docs: archive GitHub Pages migration"
  ```

  Do not push unless the owner separately authorizes it.
