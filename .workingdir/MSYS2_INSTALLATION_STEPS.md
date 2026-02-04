# MSYS2 Installation - Manual Steps Required

**Status**: ⏸️ WAITING FOR USER ACTION
**Date**: 2026-02-04

---

## Manual Steps (Cannot be Automated)

### Step 1: Download MSYS2 Installer
1. Open browser: https://www.msys2.org/
2. Download: **msys2-x86_64-YYYYMMDD.exe** (latest version)
3. Save to: `C:\Users\ms\Downloads\`

---

### Step 2: Run Installer (Requires Admin)
1. Double-click `msys2-x86_64-YYYYMMDD.exe`
2. **Installation Directory**: `C:\msys64` (default, recommended)
3. Click "Next" → "Install"
4. Wait for extraction (2-3 minutes)
5. **Uncheck** "Run MSYS2 now" at the end
6. Click "Finish"

---

### Step 3: Initial MSYS2 Setup
1. Open **Start Menu** → Search "MSYS2 MinGW 64-bit"
2. Run as Administrator (right-click → "Run as administrator")
3. In the MSYS2 terminal, run:

```bash
# Update package database
pacman -Syu
```

4. If prompted to close terminal and restart, do so
5. Reopen "MSYS2 MinGW 64-bit" terminal
6. Run update again:

```bash
pacman -Su
```

---

### Step 4: Install Toolchain + FFmpeg (Run in MSYS2 terminal)

```bash
# Install MinGW-w64 GCC toolchain
pacman -S --needed mingw-w64-x86_64-toolchain

# Press Enter to select "all" when prompted
# Press "Y" to proceed with installation
# Wait 5-10 minutes for download + install

# Install FFmpeg with all libraries
pacman -S --needed mingw-w64-x86_64-ffmpeg

# Press "Y" to proceed
# Wait 2-3 minutes
```

---

### Step 5: Add to Windows PATH (Permanent)

**Option A - GUI Method**:
1. Press `Win + R` → type `sysdm.cpl` → Enter
2. Click "Advanced" tab → "Environment Variables"
3. Under "System variables", find "Path" → "Edit"
4. Click "New" → Add: `C:\msys64\mingw64\bin`
5. Click "OK" on all dialogs
6. **Restart all terminals and VS Code**

**Option B - PowerShell Method** (Admin required):
```powershell
# Run in Administrator PowerShell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\msys64\mingw64\bin", "Machine")
```

---

### Step 6: Verify Installation (Run in NEW terminal)

After restarting terminal, run these commands:

```bash
# Check GCC
gcc --version
# Expected: gcc.exe (Rev...) 14.2.0 or similar

# Check pkg-config
pkg-config --version
# Expected: 2.3.0 or similar

# Check FFmpeg libraries
pkg-config --cflags libavcodec
# Expected: -IC:/msys64/mingw64/include

pkg-config --libs libavcodec
# Expected: -LC:/msys64/mingw64/lib -lavcodec

# Check all required FFmpeg components
pkg-config --exists libavcodec libavformat libavutil libavfilter libswscale libswresample && echo "All FFmpeg libs found" || echo "Missing libs"
# Expected: All FFmpeg libs found
```

---

### Step 7: Configure Go for CGO

```bash
# Enable CGO permanently
go env -w CGO_ENABLED=1

# Verify
go env CGO_ENABLED
# Expected: 1

# Set C compiler (should auto-detect, but be explicit)
go env -w CC=gcc

# Verify Go can find GCC
go env CC
# Expected: gcc
```

---

## Verification Checklist

- [ ] MSYS2 installed to `C:\msys64`
- [ ] `C:\msys64\mingw64\bin` added to Windows PATH
- [ ] All terminals/IDEs restarted
- [ ] `gcc --version` works in CMD/PowerShell
- [ ] `pkg-config --version` works
- [ ] All FFmpeg libraries detected by pkg-config
- [ ] `go env CGO_ENABLED` returns `1`
- [ ] `go env CC` returns `gcc`

---

## After Installation Complete

**Run this command to verify go-astiav can compile**:

```bash
cd c:\Users\ms\dev\revenge
go build ./internal/content/movie/
```

**Expected**: No errors, successful compilation

**If successful, notify Claude**: "MSYS2 installation complete, ready to proceed"

---

## Troubleshooting

### Issue: gcc not found after PATH update
**Solution**: Restart ALL terminals, VS Code, and any IDEs

### Issue: pkg-config not found
**Solution**: Verify `C:\msys64\mingw64\bin` is in PATH (run `echo $env:Path` in PowerShell)

### Issue: FFmpeg libraries not found
**Solution**:
```bash
# Re-run in MSYS2 terminal
pacman -S --needed mingw-w64-x86_64-ffmpeg
```

### Issue: "compile: version mismatch"
**Solution**:
```bash
go clean -cache -modcache -testcache
```

---

**Estimated Total Time**: 20-30 minutes (download speed dependent)
