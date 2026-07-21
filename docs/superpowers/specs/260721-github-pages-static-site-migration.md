# GitHub Pages Static Site Migration

**Status:** Approved for implementation
**Created:** 2026-07-21

## Objective

Serve `ohchanwu.dev` as a static GitHub Pages site from the existing
`ohchanwu/ohchanwu.dev` repository. Remove the Go server, Docker image, manual ACME certificate
handling, and EC2 deployment model from the active tree while preserving their history in Git.

## Constraints

- Keep the existing GitHub repository and remote. Do not create another repository.
- Keep `main` as the source branch and use the repository's normal pull-request workflow.
- Publish only the contents of `static/`; repository documentation and source metadata must not
  become web-accessible artifacts.
- Keep the site framework-free and dependency-free: plain HTML, CSS, JavaScript, images, and PDFs.
- Preserve the current visual design and behavior except for removing the obsolete footer
  technology line.
- The repository owner performs GitHub Settings, DNS, and infrastructure shutdown operations.

## Chosen architecture

A custom GitHub Actions workflow publishes `static/` directly as the GitHub Pages artifact after
changes land on `main`. The workflow uses GitHub's Pages configuration, artifact upload, and Pages
deployment actions with the permissions and `github-pages` environment required by the official
static-site workflow.

This approach keeps the deployable website isolated without moving its files or maintaining a
generated branch. The following alternatives are rejected:

- A `gh-pages` branch would duplicate the site and require synchronization.
- Publishing the root of `main` could expose unrelated repository files.
- Publishing `main:/docs` would mix the public website with project documentation.
- A static-site generator or frontend build tool would add work without changing the current site.

The workflow publishes on pushes to `main` and supports manual dispatch for a deliberate redeploy.
Deployment concurrency prevents overlapping Pages releases while allowing the latest queued release
to proceed.

## Repository changes

Delete the active server implementation and its build inputs:

- `main.go`
- `main_test.go`
- `go.mod`
- `go.sum`
- `Dockerfile`
- `.dockerignore`

Prune Go, Docker, Terraform, and server-runtime entries from `.gitignore` while retaining local
editor, environment, credential, Gas Town, gstack, and temporary SDD exclusions.

Add `.github/workflows/pages.yml` with `static/` as the only uploaded path. Do not add a `CNAME`
file: GitHub ignores it for custom Actions deployments and stores the custom domain in repository
settings. Do not add `.nojekyll`, a generated output directory, or a custom `404.html`; none is
needed for this single-page site.

Update `README.md`, `README.ko.md`, `AGENTS.md`, and `docs/README.md` so they describe the static
source, Pages deployment, local preview, and owner-operated cutover. Keep `CLAUDE.md` as the symlink
to `AGENTS.md`. Remove the footer's entire technology line in both languages without replacing it
with hosting copy, and remove the now-unused JavaScript that populated its year.

## Runtime behavior

`static/index.html` is the document root of the deployed artifact. Its relative stylesheet, script,
image, favicon, and PDF URLs remain valid without rewriting. Theme, language, media fade, and scroll
state remain browser-side features backed by the existing JavaScript and local storage. There are no
API calls or forms that require a server.

GitHub Pages owns HTTP serving, redirects, caching, certificates, and standard response headers.
The repository no longer configures ACME, origin TLS, HSTS, server timeouts, PDF response headers,
or graceful shutdown. Resume links continue to open in a new tab; whether a browser displays or
downloads a PDF follows the browser and Pages response behavior.

## Domain and cutover flow

The canonical custom domain remains `ohchanwu.dev`, with `www.ohchanwu.dev` configured as the
redirecting companion hostname.

1. Merge and verify the Pages deployment at its generated GitHub URL before changing DNS.
2. Verify ownership of the domain in the GitHub account using GitHub's requested DNS record.
3. Set `ohchanwu.dev` as the repository's Pages custom domain before pointing public DNS at Pages.
4. In Cloudflare DNS, replace the old origin records with GitHub's then-current documented apex
   records and point `www` to `ohchanwu.github.io`.
5. Keep those records DNS-only during and after the migration. The Cloudflare proxy and its second
   TLS layer are unnecessary for this static site.
6. Wait for DNS resolution and GitHub certificate issuance, then enable **Enforce HTTPS**.
7. Verify the apex and `www` hostnames, their redirect behavior, assets, and certificate before
   retiring the former server.

The old origin remains running during the cutover. Before it is retired, rollback is a DNS change
back to the prior origin records. After retirement, repository rollback is a revert of the migration
pull request followed by an intentional redeploy; restoring the old runtime would additionally
require recreating its infrastructure from Git history.

## Verification

Before merge:

- Validate the workflow syntax and confirm its artifact contains only `static/`.
- Serve `static/` through a local static HTTP server.
- Check the page at desktop and mobile viewport widths.
- Exercise English and Korean modes, light and dark themes, reload scroll restoration, and media
  loading.
- Open both resume PDFs and representative external links in the browser and verify their expected
  destinations.
- Confirm there are no browser console errors or failed first-party network requests.
- Confirm deleted Go and ACME terms no longer appear in active documentation.

After merge but before DNS changes, repeat the browser checks against the generated Pages URL. After
the owner completes the cutover, repeat them against both custom hostnames and confirm HTTPS and the
expected canonical redirect.

## Failure handling

- A failed workflow must leave the last successful Pages deployment serving traffic.
- DNS must not move until the generated Pages deployment passes verification.
- The old origin must not be terminated until the custom-domain verification passes.
- If GitHub certificate issuance is pending, keep the old origin available and wait rather than
  weakening HTTPS or reintroducing manual certificate handling.

## Out of scope

- Visual redesign or portfolio-copy rewrite
- New analytics, forms, APIs, or serverless functions
- Cloudflare proxying, Workers, Pages, or other edge logic
- A static-site generator, bundler, dependency manager, or deployment branch
- Automated EC2 teardown or DNS mutation

## References

- [Configuring a publishing source for GitHub Pages](https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site)
- [Managing a custom domain for GitHub Pages](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site/managing-a-custom-domain-for-your-github-pages-site)
