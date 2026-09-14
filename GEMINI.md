# GEMINI.md: Production AI Collaboration Guide

## 1. Core Directives & Standards
- **Zero Pseudo-Code Policy:** Never output placeholder comments like `// TODO: implement logic` or partial snippets. Every code block provided must be fully syntactically valid, compilable, and production-ready.
- **Language & Runtime:** Go (Strictly version 1.22+), adhering strictly to idiomatic patterns, explicit error management, and zero global application state outside the composition root.
- **Architecture Enforcement:** Maintain strict unidirectional layering: `cmd/app` -> `cli/` -> `task/` -> `vault/` / `processor/`. Presentation packages must never import task orchestration logic directly.

## 2. Mandatory Technical Implementations
- **Error Handling:** Avoid generic error strings. Use custom typed error structs implementing `error` and wrapping root issues via `errors.As` with explicit remediation tips.
- **International Character Support:** All filenames and metadata paths must pass through Unicode Normalization Canonical Composition (`golang.org/x/text/unicode/norm` - NFC) to guarantee multi-byte safety for Chinese, Arabic, and Cyrillic scripts.
- **Context Propagation:** All long-running worker loops, network operations, and ticker UI routines must accept and respect `context.Context` for immediate cancellation/graceful teardown.

## 3. Terminal & UI Interface Rules
- **Pagination Safety:** Paged outputs must cleanly clear buffers using ANSI escape handlers and draw structured structural boundaries via responsive dividers.
- **Signal Control:** Standardize inputs through a keyword broker mapping user keystrokes (`q` for exit, `esc` for cancel, `s` for stop, `r` for resume) to strongly typed commands.
