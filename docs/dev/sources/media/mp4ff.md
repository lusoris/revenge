# mp4ff (MP4 parsing)

> Auto-fetched from [https://pkg.go.dev/github.com/Eyevinn/mp4ff](https://pkg.go.dev/github.com/Eyevinn/mp4ff)
> Last Updated: 2026-01-29T20:14:41.395091+00:00

---

Overview
¶
Command Line Tools
Example code
Packages
Specifications
Module mp4ff implements MP4 media file parsing and writing for AVC and HEVC video, AAC and AC-3 audio, stpp and wvtt subtitles, and
timed metadata tracks.
It is focused on fragmented files as used for streaming in MPEG-DASH, MSS and HLS fMP4, but can also decode and encode all
boxes needed for progressive MP4 files.
Command Line Tools
¶
Some useful command line tools are available in [cmd](cmd) directory.
mp4ff-info
prints a tree of the box hierarchy of a mp4 file with information
about the boxes.
mp4ff-pslister
extracts and displays SPS and PPS for AVC or HEVC in a mp4 or a bytestream (Annex B) file.
Partial information is printed for HEVC.
mp4ff-nallister
lists NALUs and picture types for video in progressive or fragmented file
mp4ff-subslister
lists details of wvtt or stpp (WebVTT or TTML in ISOBMFF) subtitle samples
mp4ff-crop
crops a **progressive** mp4 file to a specified duration
mp4ff-encrypt
encrypts a fragmented file using cenc or cbcs Common Encryption scheme
mp4ff-decrypt
decrypts a fragmented file encrypted using cenc or cbcs Common Encryption scheme
You can install these tools by going to their respective directory and run `go install .` or directly from the repo with
go install github.com/Eyevinn/mp4ff/cmd/mp4ff-info@latest
go install github.com/Eyevinn/mp4ff/cmd/mp4ff-encrypt@latests
for each individual tool.
Example code
¶
Example code for some common use cases is available in the [examples](examples) directory.
The examples and their functions are:
initcreator
creates typical init segments (ftyp + moov) for different video and
audio codecs
resegmenter
reads a segmented file (CMAF track) and resegments it with other
segment durations using `FullSample`
segmenter
takes a progressive mp4 file and creates init and media segments from it.
This tool has been extended to support generation of segments with multiple tracks as well
as reading and writing `mdat` in lazy mode
multitrack
parses a fragmented file with multiple tracks
combine-segs
combines single-track init and media segments into multi-track segments
add-sidx
adds a top-level sidx box describing the segments of a fragmented files.
Packages
¶
The top-level packages in the mp4ff module are
mp4
provides support for for parsing (called Decode) and writing (Encode) a plethor of mp4 boxes.
It also contains helper functions for extracting, encrypting, dectrypting samples and a lot more.
avc
deals with AVC (aka H.264) video in the `mp4ff/avc` package including parsing of SPS and PPS,
and finding start-codes in Annex B byte streams.
hevc
provides structures and functions for dealing with HEVC video and its packaging
sei
provides support for handling  Supplementary Enhancement Information (SEI) such as timestamps
for AVC and HEVC video.
av1
provides basic support for AV1 video packaging
aac
provides support for AAC audio. This includes handling ADTS headers which is common
for AAC inside MPEG-2 TS streams.
bits
provides bit-wise and byte-wise readers and writers used by the other packages.
Specifications
¶
The main specification for the MP4 file format is the ISO Base Media File Format (ISOBMFF) standard
ISO/IEC 14496-12 7th edition 2021. Some boxes are specified in other standards, as should be commented
in the code.