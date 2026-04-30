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
