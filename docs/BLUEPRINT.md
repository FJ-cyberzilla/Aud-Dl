# Architectural Blueprint: System Refactor

## Overview
This blueprint outlines the migration of the application to a highly modular, clean architecture. It centralizes control in orchestrators and separates concerns into domain-specific subsystems.

### Core Architectural Layers
The system follows a strict unidirectional dependency graph:
`cmd/app` -> `cli/` -> `task/` -> `vault/` / `processor/`

### Centralized Orchestration: SystemBindings
The `internal/task` package utilizes `SystemBindings` as a centralized dependency provider. This pattern ensures that all orchestration subsystems (`TaskDivision`, `TaskController`) are correctly initialized with required dependencies during the application bootstrap phase.

## Current Project Structure

```text
/data/data/com.termux/files/home/Aud-Dl/
├───BLUEPRINT.md
├───GEMINI.md
├───go.mod
├───go.sum
├───plan.md
├───tree.txt
├───cmd/
│   └───app/
│       └───main.go
├───config/
├───internal/
│   ├───cache/
│   ├───cli/
│   │   ├───cli_test.go
│   │   ├───interactive_menu.go
│   │   ├───interactive_style.go
│   │   ├───junction_renderer.go
│   │   ├───keyword.go
│   │   ├───menu_crust.go
│   │   ├───menu_junction.go
│   │   ├───menu_template.go
│   │   ├───renderer_admin.go
│   │   └───styling.go
│   ├───config/
│   │   ├───app_config_test.go
│   │   └───app_config.go
│   ├───gate/
│   │   ├───hub_controller.go
│   │   ├───hub_validator.go
│   │   └───system_shield.go
│   ├───processor/
│   │   ├───elite_tagger.go
│   │   ├───errors.go
│   │   ├───pipe_transcoder.go
│   │   └───processor_test.go
│   ├───search/
│   ├───strategy/
│   ├───task/
│   │   ├───task_admin.go
│   │   ├───task_binder.go
│   │   ├───task_controller.go
│   │   └───task_division.go
│   ├───ui/
│   │   ├───assist_dic.go
│   │   ├───barchart.go
│   │   ├───display_assist.go
│   │   ├───display_colors.go
│   │   ├───display_paginated.go
│   │   ├───dividers.go
│   │   ├───platform_display.go
│   │   ├───refresh.go
│   │   ├───table.go
│   │   ├───theme.go
│   │   └───ui_test.go
│   └───vault/
│       ├───duplicate_detector.go
│       ├───fingerprint_db.go
│       ├───persistent_detector.go
│       ├───playlist.go
│       ├───vault_admin.go
│       ├───vault_resolver.go
│       ├───vault_test.go
│       └───vault.go
└───music_vault/
```

## Subsystem Architecture

### 1. Gate Subsystem (`internal/gate`)
The gateway manages ingress/egress, security, and updates.
- **`hub_controller.go`**: Flow orchestrator.
- **`hub_validator.go`**: Path/Query sanitizer.

### 2. Task Subsystem (`internal/task`)
Orchestrates high-level application workflows.
- **`task_controller.go`**: Main loop for CLI interaction and lifecycle management.
- **`task_division.go`**: Logic for manifest processing and stream handling. It contains the `processManifest` pipeline, which coordinates fetching via `GateHub`, transcoding via `MemoryPipedTranscoder`, metadata tagging via `EliteTagger`, and ingestion into the vault using `VaultAdmin`.
- **`task_binder.go`**: Defines `SystemBindings` and `TaskBinder`, which manages the initialization and injection of all subsystems, including the `VaultAdmin` service, to ensure correct service wiring at runtime.

### 3. Vault Subsystem (`internal/vault`)
Handles storage, indexing, and deduplication.
- **`vault_admin.go`**: Facade for ingestion, index management, and storage lifecycle.
- **`vault_resolver.go`**: Logic for path resolution and directory structure strategy.
- **`persistent_detector.go`**: Acoustic fingerprinting and database interaction.

### 5. Cache Subsystem (`internal/cache`)
Handles high-performance in-memory caching.
- **`cache.go`**: Contains `CacheManager`, which provides thread-safe key-value storage with TTL-based expiration and an automated background janitor for memory management.

### 4. CLI & UI Subsystem (`internal/cli` / `internal/ui`)
Manages terminal user interaction and rendering.
- **`cli/interactive_menu.go`**: Input evaluation and menu navigation.
- **`ui/display_components.go`**: Advanced UI rendering components (`DisplayCompiler`, `DisplayAssist`, `HealthQuest`).
- **`ui/theme.go`**: Terminal styling and color management.
