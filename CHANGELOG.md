# Changelog

## [0.1.0] - 2026-08-25
### Added
- Implemented core audio processing pipeline in `internal/task/task_division.go`.
- Integrated `VaultAdmin` into `SystemBindings` in `internal/task/task_binder.go` for structured vault management.
- Enabled functional audio streaming orchestration: fetch (GateHub) -> transcode (PipeTranscoder) -> tag (EliteTagger) -> ingest (VaultAdmin).
- Enhanced `VaultAdmin` initialization in `TaskBinder`.
- Established atomic file movement with fallback copy in pipeline ingestion.

### Changed
- Refactored `TaskDivision.processManifest` to orchestrate multi-subsystem tasks.
- Updated `SystemBindings` struct to include `*vault.VaultAdmin`.
- Streamlined `BindSystemServices` in `TaskBinder` for seamless service dependency injection.

### Fixed
- Resolved missing `VaultAdmin` dependency in task orchestration.
