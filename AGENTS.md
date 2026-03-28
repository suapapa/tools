# Agent Workspace Instructions

This file provides guidance for AI agents (like Gemini CLI) interacting with this codebase.

## Project Structure

- `cmd/`: Each subdirectory contains a standalone command-line tool.
- `pkg/`: Reusable packages intended for use by multiple tools.
- `internal/`: Private packages for specific tools that should not be exposed externally.

## Core Conventions

- **Idiomatic Go**: Follow standard Go patterns (Go 1.21+). Use `context.Context` for all I/O and external process execution.
- **Error Handling**: Prefer wrapping errors with `%w`. Avoid `panic` and `log.Fatal` in library code or `run()` functions; handle them in `main()`.
- **Modern Standards**: Use `io` and `os` packages instead of the deprecated `io/ioutil`.
- **Modularity**: When adding shared logic, place it in `pkg/` if it's generally useful, or `internal/` if it's tool-specific.

## Maintenance Notes

This project was systematically refactored in March 2026 to modernize the structure and dependencies. When adding new tools, ensure they follow the pattern of having a `run() error` function called by `main()`.
