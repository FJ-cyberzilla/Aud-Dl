# Aud-Dl (AudioCommandCenter)

*Powered by FJ™ Cybertronic Systems*

Aud-Dl is a professional-grade audio downloading and management utility designed for efficiency, resilience, and advanced performance, adhering to strict idiomatic Go standards.

## Branding & Contact

*   **Developer:** FJ™ Cybertronic Systems
*   **GitHub Repository:** [https://github.com/FJ-cyberzilla](https://github.com/FJ-cyberzilla)
*   **Contact Email:** cyberzilla.systems@gmail.com

## Architecture & Modules

The application is built on a strictly layered architecture to ensure separation of concerns, utilizing `SystemBindings` as the central dependency injection pattern for `task` orchestration.

*   **`cmd/`**: Entry point for the application (Composition root).
*   **`internal/config/`**: Centralized application configuration.
*   **`internal/task/`**: Orchestration & Lifecycle Subsystem (leveraging `SystemBindings`).
*   **`internal/gate/`**: Ingress/Egress & Network Security Subsystem.
*   **`internal/strategy/`**: Resiliency & Request Harmonics Suite.
*   **`internal/search/`**: Audio Search Engine & Discovery Package.
*   **`internal/cache/`**: Performance & In-Memory Caching Engine (`CacheManager` implementation with automated eviction).
*   **`internal/processor/`**: Transcoding, Tagging & Recovery Engine.
*   **`internal/vault/`**: Storage, Indexing & Deduplication Subsystem.
*   **`internal/cli/`**: CLI Menus, Junctions & Skin Facade.
*   **`internal/ui/`**: Terminal Rendering & UX Engine (including `DisplayCompiler`, `DisplayAssist`, `HealthQuest`).

## Project Structure

```text
audio-command-center/
├── cmd/
│   └── app/
│       └── main.go                  # Entry point (Composition root)
├── internal/
│   ├── config/                      # Application Configuration Management
│   ├── task/                        # Orchestration & Lifecycle Subsystem
│   ├── gate/                        # Ingress/Egress & Network Security Subsystem
│   ├── strategy/                    # Resiliency & Request Harmonics Suite
│   ├── search/                      # Audio Search Engine & Discovery Package
│   ├── cache/                       # Performance & In-Memory Caching Engine
│   ├── processor/                   # Transcoding, Tagging & Recovery Engine
│   ├── vault/                       # Storage, Indexing & Deduplication Subsystem
│   ├── cli/                         # CLI Menus, Junctions & Skin Facade
│   └── ui/                          # Terminal Rendering & UX Engine
├── config/                          # Configuration files directory
├── music_vault/                     # Default Output Music Directory
├── go.mod
└── go.sum
```

## Standards

This project adheres to the following core mandates:
*   **Zero Pseudo-Code Policy:** All code must be production-ready.
*   **Layered Architecture:** Unidirectional data flow from `cmd` to `vault/processor`.
*   **Internationalization:** Full support for Unicode Normalization (NFC).
*   **Signal Handling:** Robust context propagation for graceful teardown.
