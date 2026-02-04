# CGO + FFmpeg Setup for go-astiav - Windows

**Date**: 2026-02-04
**Goal**: Enable CGO and install FFmpeg development libraries for go-astiav compilation
**Status**: 🔄 IN PROGRESS

---

## Problem Statement

go-astiav requires:
1. **CGO enabled** (currently: CGO_ENABLED=0)
2. **GCC/MinGW compiler** (currently: not found)
3. **FFmpeg development libraries** (headers + .lib/.a files)

Without these, tests fail with:
```
# github.com/asticode/go-astiav
undefined: MediaType
undefined: BuffersinkFlag
...
```

---

## Current Environment

```
Go Version: go1.25.6 windows/amd64
CGO_ENABLED: 0
CC: gcc (not installed)
Platform: Windows

Existing /mingw64/bin: EXISTS but contains .NET/Avalonia DLLs (not MSYS2)
MSYS2: NOT INSTALLED
vcpkg: NOT INSTALLED
GCC: NOT FOUND
FFmpeg: NOT FOUND
```

**Conclusion**: Clean installation needed

---

## Installation Steps

### Step 1: Install MinGW-w64 (GCC Compiler)

**Options**:
- **A) MSYS2** (recommended) - Full Unix-like environment
- **B) WinLibs standalone** - Minimal, portable
- **C) TDM-GCC** - Easy installer

**Chosen**: [PENDING USER CHOICE]

**Installation Command**: [TBD]

**Verification**:
```bash
gcc --version
# Expected: gcc (GCC) X.X.X
```

---

### Step 2: Install FFmpeg Development Libraries

**Options**:
- **A) MSYS2 pacman** - `pacman -S mingw-w64-x86_64-ffmpeg`
- **B) vcpkg** - `vcpkg install ffmpeg[core,avcodec,avformat,avutil]:x64-windows`
- **C) Manual download** - From ffmpeg.org (shared + dev packages)

**Chosen**: [PENDING USER CHOICE]

**Required Components**:
- libavcodec (headers + libs)
- libavformat (headers + libs)
- libavutil (headers + libs)
- libavfilter (headers + libs)
- libswscale (headers + libs)
- libswresample (headers + libs)

**Installation Command**: [TBD]

---

### Step 3: Configure CGO Environment

**Environment Variables to Set**:

```bash
# Enable CGO
CGO_ENABLED=1

# Point to MinGW GCC
CC=gcc

# FFmpeg pkg-config paths (depends on installation method)
PKG_CONFIG_PATH=/mingw64/lib/pkgconfig  # MSYS2 example

# Or manually set flags
CGO_CFLAGS=-I/path/to/ffmpeg/include
CGO_LDFLAGS=-L/path/to/ffmpeg/lib -lavcodec -lavformat -lavutil -lavfilter
```

**Verification**:
```bash
go env CGO_ENABLED  # Should output: 1
gcc --version       # Should show version
pkg-config --cflags libavcodec  # Should show include paths
```

---

### Step 4: Test go-astiav Compilation

**Test Command**:
```bash
go build ./internal/content/movie/
```

**Expected Output**:
```
# Success - no errors
```

**If Fails**: Document error in "Bugs" section below

---

### Step 5: Run Full Test Suite

**Command**:
```bash
go test -coverprofile=coverage.out ./...
```

**Expected**: All tests pass, coverage report generated

---

## Progress Tracker

| Step | Status | Time | Notes |
|------|--------|------|-------|
| 1. MinGW-w64 | ✅ COMPLETE | 3min | winget install MSYS2.MSYS2 |
| 2. FFmpeg Libs | ✅ COMPLETE | 8min | pacman -S mingw-w64-x86_64-toolchain mingw-w64-x86_64-ffmpeg (127 packages, 1.8GB) |
| 3. CGO Config | ✅ COMPLETE | 1min | go env -w CGO_ENABLED=1, Go 1.25.5 → 1.25.6 upgrade |
| 4. Test Compilation | ✅ COMPLETE | <1min | go build ./internal/content/movie/ SUCCESS |
| 5. Test Suite | 🔄 RUNNING | - | go test -coverprofile=coverage.out ./... (background) |

---

## Bugs / Issues

### Issue #30: Go Version Mismatch (RESOLVED ✅)
**Status**: RESOLVED
**Description**: go.mod requires go >= 1.25.6 but Go 1.25.5 was installed, causing "package X is not in std" errors
**Resolution**:
- Upgraded Go via `winget upgrade --id GoLang.Go` (1.25.5 → 1.25.6)
- Reset GOTOOLCHAIN to auto
- Compilation now successful

### Issue #31: CGO_ENABLED=0 by default (RESOLVED ✅)
**Status**: RESOLVED
**Description**: CGO was disabled, go-astiav couldn't compile
**Resolution**: `go env -w CGO_ENABLED=1`

---

## Questions / Decisions Needed

### Q1: Which MinGW-w64 distribution?
**Options**:
- MSYS2 (full toolchain, 1GB+, includes package manager)
- WinLibs (minimal, portable, 200MB)
- TDM-GCC (easy installer, older versions)

**Recommendation**: **MSYS2** - All-in-one solution with pacman for easy FFmpeg installation

**Decision**: **MSYS2** (proceeding with Option A)

---

### Q2: Which FFmpeg installation method?
**Options**:
- MSYS2 pacman (easiest with MSYS2)
- vcpkg (cross-platform, integrates with Visual Studio)
- Manual download (most control, most complex)

**Recommendation**: MSYS2 pacman if using MSYS2, otherwise vcpkg

**Decision**: [PENDING]

---

### Q3: System-wide or project-local installation?
**Options**:
- System-wide (add to PATH permanently)
- Project-local (portable, CI-friendly)

**Recommendation**: System-wide for development, document for CI/Docker

**Decision**: [PENDING]

---

## Installation Commands (Ready to Execute)

### Option A: MSYS2 (Full Toolchain)

```bash
# 1. Download MSYS2 installer from https://www.msys2.org/
# 2. Run installer, install to C:\msys64
# 3. Open MSYS2 MinGW 64-bit terminal
# 4. Update packages
pacman -Syu

# 5. Install MinGW-w64 toolchain
pacman -S mingw-w64-x86_64-toolchain

# 6. Install FFmpeg with all components
pacman -S mingw-w64-x86_64-ffmpeg

# 7. Add to Windows PATH (permanent)
# Add: C:\msys64\mingw64\bin

# 8. Verify in PowerShell/CMD
gcc --version
pkg-config --cflags libavcodec

# 9. Set Go environment
go env -w CGO_ENABLED=1
go env -w CC=gcc
```

### Option B: vcpkg (Cross-Platform)

```bash
# 1. Install vcpkg
git clone https://github.com/Microsoft/vcpkg.git C:\vcpkg
cd C:\vcpkg
.\bootstrap-vcpkg.bat

# 2. Install FFmpeg
.\vcpkg install ffmpeg[core,avcodec,avformat,avutil]:x64-windows

# 3. Install MinGW-w64 (WinLibs standalone)
# Download from: https://winlibs.com/
# Extract to C:\mingw64
# Add C:\mingw64\bin to PATH

# 4. Set environment variables
$env:CGO_ENABLED=1
$env:CGO_CFLAGS="-IC:\vcpkg\installed\x64-windows\include"
$env:CGO_LDFLAGS="-LC:\vcpkg\installed\x64-windows\lib -lavcodec -lavformat -lavutil"
```

---

## Next Actions

1. **USER DECISION**: Choose installation method (A or B)
2. Execute installation commands
3. Verify setup with test compilation
4. Document any issues encountered
5. Proceed with full test suite

---

**Last Updated**: 2026-02-04
**Updated By**: Claude Code Agent
