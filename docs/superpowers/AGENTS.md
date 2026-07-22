# Implementation knowledge instructions

- Keep only active plans in `plans/` and only active or awaiting-implementation requirements in `specs/`.
- Record durable choices in `decisions/`; update maintained documentation when a decision changes.
- When implementation finishes, distill useful conclusions, move only lasting records to `archive/`, and update this directory's index.
- Prefer Git history over repetitive or obsolete reports.
- Sanitize every tracked record for credentials, personal data, private topology, raw logs, and machine-local paths before committing.
- Do not load `archive/` during ordinary work. Follow `archive/AGENTS.md` when historical evidence is specifically relevant.
