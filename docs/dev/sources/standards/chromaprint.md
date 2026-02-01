# Chromaprint/AcoustID

> Source: https://acoustid.org/chromaprint
> Fetched: 2026-02-01T11:51:27.352416+00:00
> Content-Hash: 8b4cac6e789bf17c
> Type: html

---

## Chromaprint

Chromaprint is the core component of the AcoustID project. It's a client-side library that implements a custom algorithm for extracting fingerprints from any audio source. Overview of the fingerprint extraction process can be found in the blog post ["How does Chromaprint work?"](https://oxygene.sk/2011/01/how-does-chromaprint-work/).

### Download

Latest release — 1.6.0 (2025-08-28)

- [chromaprint-1.6.0.tar.gz](https://github.com/acoustid/chromaprint/releases/download/v1.6.0/chromaprint-1.6.0.tar.gz) (1.6 MB)
- [chromaprint-fpcalc-1.6.0-linux-arm64.tar.gz](https://github.com/acoustid/chromaprint/releases/download/v1.6.0/chromaprint-fpcalc-1.6.0-linux-arm64.tar.gz) (1.4 MB)
- [chromaprint-fpcalc-1.6.0-linux-x86_64.tar.gz](https://github.com/acoustid/chromaprint/releases/download/v1.6.0/chromaprint-fpcalc-1.6.0-linux-x86_64.tar.gz) (2.4 MB)
- [chromaprint-fpcalc-1.6.0-macos-arm64.tar.gz](https://github.com/acoustid/chromaprint/releases/download/v1.6.0/chromaprint-fpcalc-1.6.0-macos-arm64.tar.gz) (1.3 MB)
- [chromaprint-fpcalc-1.6.0-macos-universal.tar.gz](https://github.com/acoustid/chromaprint/releases/download/v1.6.0/chromaprint-fpcalc-1.6.0-macos-universal.tar.gz) (2.6 MB)
- [chromaprint-fpcalc-1.6.0-macos-x86_64.tar.gz](https://github.com/acoustid/chromaprint/releases/download/v1.6.0/chromaprint-fpcalc-1.6.0-macos-x86_64.tar.gz) (1.4 MB)
- [chromaprint-fpcalc-1.6.0-windows-x86_64.zip](https://github.com/acoustid/chromaprint/releases/download/v1.6.0/chromaprint-fpcalc-1.6.0-windows-x86_64.zip) (1.8 MB)

Most Linux distributions also have their own packages for Chromaprint.

You can find downloads for older releases on [GitHub](https://github.com/acoustid/chromaprint/releases).

### Usage

The library exposes a simple C API. The documentation for the C API can be found in [chromaprint.h](https://github.com/acoustid/chromaprint/blob/master/src/chromaprint.h).

Note that the library only calculates audio fingerprints from the provided raw uncompressed audio data. It does not deal with audio file formats in any way. Your application needs to find a way to decode audio files (MP3, MP4, FLAC, etc.) and feed the uncompressed data to Chromaprint.

You can use [pyacoustid](https://pypi.python.org/pypi/pyacoustid) to interact with the library from Python. It provides a direct wrapper around the library, but also higher-level functions for generating fingerprints from audio files.

You can also use the fpcalc utility programatically. It can produce JSON output, which should be easy to parse in any language. This is the recommended way to use Chromaprint if all you need is generate fingerprints for AcoustID.

### Development

You can dowload the development version of the source code from [GitHub](https://github.com/acoustid/chromaprint). Either you can use [Git](http://git-scm.com/) to clone the repository or download a zip/tar.gz file with the latest version.

You will need a C++ compiler and [CMake](https://cmake.org/) to build the library. [FFmpeg](https://ffmpeg.org/) is required to build the fpcalc tool.

    git clone https://github.com/acoustid/chromaprint.git
    cd chromaprint
    cmake .
    make
    

See the [README](https://github.com/acoustid/chromaprint/blob/master/README.md) file for more details on building the library.

Software created by [Lukáš Lalinský](https://oxygene.sk/), hosted by [AcoustID OÜ](https://acoustid.biz), data crowd-sourced by [thousands of contributors](/stats).

[Contact](/contact) | [Blog](https://blog.acoustid.org/) | [Twitter](https://twitter.com/acoustid) | [Facebook](https://www.facebook.com/acoustid)
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
