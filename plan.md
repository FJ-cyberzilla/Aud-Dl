# Refactoring Plan: Resolving Architectural Duplication

## 1. Goal
Consolidate duplicated types and resolve namespace collisions in `internal/strategy` and `internal/vault` to enable a successful project build and maintain structural integrity.

## 2. Strategy
- **Surgical Consolidation:** Eliminate duplicated type definitions.
- **Architectural Alignment:** Use explicit composition over complex inheritance if needed, but primarily focus on making existing types unique and logically placed.
- **Thermo-Nuclear Audit:** Identify "God classes" or over-complex files that caused the duplication and split them if necessary.

## 3. Execution Steps

### Phase 1: Strategy Package Refactoring
- **Action:** Delete `internal/strategy/conductor_deps.go` entirely as it contains stubbed, duplicated types that violate the Zero Pseudo-Code policy.
- **Verification:** Run `go build ./internal/strategy/...` to ensure all necessary types are present and correctly referenced.

### Phase 2: Vault Package Refactoring
- **Action:**
  - Rename `duplicate_detector.go`'s `AudioFingerprint` to `AcousticFingerprint`.
  - Rename `persistent_detector.go`'s `AudioFingerprint` to `HashFingerprint`.
  - Update all references in `internal/vault` to use the new, distinct names.
  - Create `internal/vault/models.go` to store these shared type definitions if they are used across multiple files.
- **Verification:** Run `go build ./internal/vault/...`.

### Phase 3: Final Build and Validation
- **Action:** Run the full project build and execute existing test suites.
- **Validation:** Ensure 0 build errors and passing tests.

## 4. Constraints
- Strict adherence to `internal/` package rules (no circular imports).
- Preservation of all existing functional logic.
- Mypy/Type check compliant Go code (idiomatic interfaces).
