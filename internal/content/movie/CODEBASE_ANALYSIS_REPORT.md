# Codebase Analysis & Design Compliance Report

**Date**: 2026-02-04
**Status**: v0.3.0 MVP (Deep Dive)

## 1. Codebase & Design Alignment Analysis

### 1.1 Design Document Compliance
Based on the implementation status versus the referenced design documents in `TODO_COMPREHENSIVE_v0.3.0.md`:

*   **MFA Design (`/docs/dev/design/services/MFA.md`)**:
    *   **Status**: Fully Implemented.
    *   **Compliance**: The implementation follows the "Unified MFA Manager" pattern described in the design, orchestrating TOTP, WebAuthn, and Backup Codes via `internal/service/mfa/manager.go`.
    *   **Security**: Secrets are stored using AES-256-GCM (as seen in `internal/crypto/encryption.go`), adhering to the security design requirements.

*   **Movie Module Design (`/docs/dev/design/features/video/MOVIE_MODULE.md`)**:
    *   **Status**: Backend Complete.
    *   **Compliance**:
        *   **Entity Layer**: `Movie`, `MovieFile` structs match the schema defined in migrations.
        *   **Library Scanning**: The implementation has shifted from a CLI-based approach to using `go-astiav` (FFmpeg bindings) in `mediainfo.go`, aligning with the requirement for high-performance metadata extraction.
        *   **Jobs**: River jobs (`MovieMetadataRefreshJob`, `MovieLibraryScanJob`) implement the asynchronous processing model defined in the design.

*   **Search Design (`/docs/dev/design/services/SEARCH.md`)**:
    *   **Status**: Implemented (Typesense).
    *   **Compliance**: The schema (`movie_schema.go`) includes specific facets (genres, year, resolution) and sortable fields (popularity, added_at) as specified in the design.

### 1.2 Source Code Analysis: `internal/content/movie/mediainfo.go`
A detailed review of the provided source code reveals:

*   **Technology Choice**: The code uses `github.com/asticode/go-astiav` for direct CGO bindings to FFmpeg libraries. This is a significant architectural decision that impacts build requirements (CGO enabled, dev libs required).
*   **Implementation Detail**:
    *   **Prober Interface**: The code defines a `Prober` interface, which is crucial for the testing strategy (mocking FFmpeg calls).
    *   **HDR Detection**: The `detectDynamicRange` function explicitly checks for `ColorTransferCharacteristicSmptest2084` (PQ) and `AribStdB67` (HLG) to identify HDR10/Dolby Vision content.
    *   **Memory Management**: Explicit `formatCtx.Free()` and `formatCtx.CloseInput()` calls are present, mitigating memory leaks common in CGO.
    *   **Stream Handling**: It iterates through streams to separate Video, Audio, and Subtitle data, extracting granular details like `ColorSpace`, `ColorPrimaries`, and `IsForced` flags for subtitles.

## 2. Gap Analysis vs. Current Code

### 2.1 Identified Discrepancies
*   **File Watching**: The design implies real-time monitoring. While `fsnotify` is in `go.mod`, the code currently relies on River jobs (`MovieLibraryScanJob`) or Radarr Webhooks. This is a deviation from a "reactive" design to a "polling/event-driven" implementation for v0.3.0.
*   **Session Handling**: The bug report (`session_service_bugs.md`) indicates a mismatch between the `pgx` driver error handling and the repository layer (`sql.ErrNoRows` vs `pgx.ErrNoRows`). The design likely assumes standard `database/sql` interfaces, but the implementation leaks driver-specific behavior.
*   **User Filtering**: The `User Service` bug (`user_service_bugs.md`) shows that the SQL generation logic (likely `sqlc`) treats `nil` pointers as default values (`false`) for booleans, breaking the "optional filter" design pattern.

## 3. Architecture & Technology Stack
*   **Backend**: Go (Modular Monolith)
*   **Database**: PostgreSQL (using `pgx` driver, not standard `lib/pq`)
*   **Caching**: Hybrid L1 (Otter) / L2 (Dragonfly/Redis)
*   **Search**: Typesense (Schema defined in code)
*   **Job Queue**: River (PostgreSQL-based)
*   **Media Processing**: FFmpeg via CGO (`go-astiav`)

## 4. Critical Action Items

1.  **Fix Driver Abstraction**: In `internal/service/session/repository.go`, normalize `pgx` errors to domain errors (`ErrNotFound`) to satisfy the interface contract.
2.  **Fix User Filtering**: Modify the SQL query or the Go wrapper to handle `nil` filters correctly (do not default to `false`).
3.  **Build Pipeline**: Ensure the Docker build pipeline includes `libavformat-dev`, `libavcodec-dev`, `libavutil-dev`, and `libswscale-dev` to support the `go-astiav` dependency found in `mediainfo.go`.
4.  **Test Coverage**: The current 5.8% coverage is critically low given the complexity of the `MediaInfo` parsing logic.
