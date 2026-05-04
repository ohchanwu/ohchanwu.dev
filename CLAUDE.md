# ohchanwu.dev

Personal portfolio website for Chanwu (Tyler) Oh.

## Stack

- Go (backend, single binary)
- Raw HTML/CSS/JS (no framework)
- Docker / docker-compose for deployment
- Cloudflare DNS proxied, Full (strict) TLS mode
- Deployed on EC2 t3.micro (Amazon Linux 2023)

## Conventions

- Match existing code style; do not introduce new dependencies without asking.
- Write commit messages in conventional commit format: `feat:`, `fix:`, `chore:`, etc.
- Do not add comments that restate what code does; only comments that explain _why_.
- Prefer standard library over third-party packages when reasonable.
- Tests live alongside the code they test (`foo.go` + `foo_test.go`).

## Workflow

- Work happens on feature branches, never directly on `main`.
- Branch naming: `feat/`, `fix/`, `chore/` prefixes.
- All changes go through PRs; CodeRabbit reviews automatically.

## Out of scope

- Do not generate AI-voice prose for the README or site copy. I write those myself.
- Do not run deployment commands; I handle deploys.

## Working with Claude Code

- Default to small, scoped changes. If a task touches more than ~3 files or ~150 lines, propose breaking it into multiple PRs before starting.
- Show diffs before writing files when changes are non-trivial. Wait for confirmation if I haven't pre-approved the approach.
- After making changes, do not run `git commit` or `git push` unless explicitly asked. I handle commits.
- When you encounter ambiguity (e.g., "the user said X but the code suggests Y"), stop and ask. Don't guess and proceed.
- If you suggest using a third-party library, name the alternative stdlib approach and let me decide.

## Project structure

- Entry point: `main.go` at repo root.
- HTTP handlers: in handler files at repo root for now (e.g., `handlers.go`). Move to `internal/` only when it's clearly justified.
- Templates: `templates/` directory, served via `html/template`.
- Static assets: `static/` directory (CSS, JS, images), served via `http.FileServer`.
- Tests live next to the code they test.

## What "production-ready" means here

This is a personal portfolio site, not a production service. Skip:

- Elaborate error handling for impossible cases
- Defensive coding for inputs the site doesn't take
- Premature abstraction (interfaces with one implementation, etc.)
- Microservice patterns

Prefer simplicity that someone can read and understand in 30 seconds.
