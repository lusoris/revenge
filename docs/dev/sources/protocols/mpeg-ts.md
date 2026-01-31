# MPEG Transport Stream

> Source: https://en.wikipedia.org/wiki/MPEG_transport_stream
> Fetched: 2026-01-30T23:52:49.543089+00:00
> Content-Hash: a442862878f7d4a8
> Type: html

---

From Wikipedia, the free encyclopedia

Digital video format used for storage network transmission

"MPEG-2 Systems" redirects here. For the program stream technology also specified in the MPEG-2 Systems standard, see

MPEG program stream

.

".ts" redirects here. For the programming language associated with the .ts file extension, see

TypeScript

. For the translation source files in Qt, see

Qt Linguist

.

"SPTS" redirects here. For the American television production and distribution arm of

Sony

, see

Sony Pictures Television

.

Not to be confused with

.m2ts

.

MPEG Transport Stream

Filename extension

.ts, .tsv, .tsa, .m2t

[

1

]

Internet media type

video/MP2T

[

2

]

Uniform Type Identifier (UTI)

public.mpeg-2-transport-stream

[

3

]

Developed by

MPEG

Initial release

10 July 1995

; 30 years ago

(

1995-07-10

)

[

4

]

Latest release

ISO/IEC 13818-1:2022

September 2022

; 3 years ago

(

2022-09

)

Type of format

Container format

Container for

Audio, video, data

Extended to

M2TS

,

TOD

Standard

ISO/IEC 13818-1, ITU-T Recommendation H.222.0

[

4

]

Open format

?

Yes

Free format

?

Yes

[

5

]

MPEG transport stream

(

MPEG

-TS

,

MTS

) or simply

transport stream

(

TS

) is a standard

digital container format

for transmission and storage of

audio

,

video

, and

Program and System Information Protocol

(PSIP) data.

[

6

]

It is used in broadcast systems such as

DVB

,

ATSC

and

IPTV

.

Transport stream specifies a container format encapsulating

packetized elementary streams

, with

error correction

and

synchronization pattern

features for maintaining transmission integrity when the

communication channel

carrying the stream is

degraded

.

Transport streams differ from the similarly named

MPEG program stream

in several important ways: program streams are designed for reasonably reliable media, such as discs (like

DVDs

), while transport streams are designed for less

reliable

transmission, namely

terrestrial

or

satellite broadcast

. Further, a transport stream may carry multiple programs.

Transport stream is specified in

MPEG-2

Part 1, Systems

, formally known as

ISO/IEC

standard 13818-1

or

ITU-T

Rec. H.222.0

.

[

4

]

Overview

[

edit

]

Multiple MPEG programs are combined then sent to a transmitting antenna. The receiver parses and decodes one of the streams.

A transport stream encapsulates a number of other substreams, often

packetized elementary streams

(PESs) which in turn wrap the

main data stream

using the MPEG codec or any number of non-MPEG codecs (such as

AC3

or

DTS

audio, and

MJPEG

or

JPEG 2000

video), text and pictures for subtitles, tables identifying the streams, and even broadcaster-specific information such as an

electronic program guide

. Many streams are often mixed together, such as several different television channels, or multiple

angles

of a movie.

Each stream is chopped into (at most) 188-byte sections and interleaved together.  Due to the tiny packet size, streams can be interleaved with less latency and greater error resilience compared to

program streams

and other common containers such as

AVI

,

MOV

/

MP4

, and

MKV

, which generally wrap each frame into one packet. This is particularly important for videoconferencing, where large frames may introduce unacceptable audio delay.

Transport streams tend to be broadcast as

constant bitrate

(CBR) and filled with padding bytes when not enough data exists.

[

a

]

Elements

[

edit

]

Packet

[

edit

]

A

network packet

is the basic unit of data in a transport stream, and a transport stream is merely a sequence of packets. Each packet starts with a

sync byte

and a

header

, which may be followed by optional additional headers; the rest of the packet consists of

payload

. All header fields are read as

big-endian

. Packets are 188 bytes in length, but the communication medium may add additional information.

[

b

]

The 188-byte packet size was originally chosen for compatibility with

Asynchronous Transfer Mode (ATM) systems

.

[

8

]

[

9

]

Programs

[

edit

]

Transport stream has a concept of

programs

. Every program is described by a program map table (PMT). The elementary streams associated with that program have PIDs listed in the PMT. Another PID is associated with the PMT itself. For instance, a transport stream used in digital television might contain three programs to represent three television channels. Suppose each channel consists of one video stream, one or two audio streams, and any necessary metadata. A

receiver

wishing to decode one of the three channels merely has to decode the payloads of each PID associated with its program. It can discard the contents of all other PIDs. A transport stream with more than one program is referred to as a multi-program transport stream (MPTS). A single program transport stream is referred to as a single-program transport stream (SPTS).

Program specific information

[

edit

]

Main article:

Program-specific information

There are 4 program-specific information (PSI) tables: program association (PAT), program map (PMT), conditional access (CAT), and network information (NIT). The MPEG-2 specification does not specify the format of the CAT and NIT.

PCR

[

edit

]

To enable a decoder to present synchronized content, such as audio tracks matching the associated video, at least once each 100 ms, a

program clock reference

(PCR) is transmitted in the adaptation field of an MPEG-2 transport stream packet. The PID with the PCR for an MPEG-2 program is identified by the

pcr_pid

value in the associated PMT. The value of the PCR, when properly used, is employed to generate a

system_timing_clock

in the decoder. The system time clock (STC) decoder, when properly implemented, provides a highly accurate time base that is used to synchronize audio and video elementary streams. Timing in MPEG-2 references this clock. For example, the

presentation time stamp

(PTS) is intended to be relative to the PCR. The first 33 bits are based on a 90 kHz clock. The last 9 bits are based on a 27 MHz clock. The maximum jitter permitted for the PCR is

+/- 500 ns

.

Null packets

[

edit

]

Some transmission schemes, such as those in

ATSC

and

DVB

, impose strict constant bitrate requirements on the transport stream. In order to ensure that the stream maintains a constant bitrate, a multiplexer may need to insert some additional packets. The PID 0x1FFF is reserved for this purpose. The null packets have a payload that is filled with 0xFF, and the receiver is expected to ignore its contents.

[

10

]

M2TS

[

edit

]

Main article:

.m2ts

Transport Stream was originally designed for broadcast. Later, it was adapted for use with digital video cameras, recorders and players by adding a 4-byte timecode (TC) field to the standard 188-byte packets, resulting in a 192-byte packet.

[

11

]

[

12

]

This is what is informally called

M2TS

stream,  commonly found in

HDV

cameras. The timecode allows quick access to any part of the stream, either from a media player or from a non-linear video editing system.

[

13

]

Use in digital video cameras

[

edit

]

JVC called M2TS "

TOD

"

[

c

]

when used in HDD-based camcorders like

GZ-HD7

.

[

14

]

[

15

]

It is also used to synchronize video streams from several cameras in a

multiple-camera setup

.

Use in Blu-ray

[

edit

]

Blu-ray Disc video titles authored with menu support are in the

Blu-ray Disc Movie

(BDMV) format and contain audio, video, and other streams in a BDAV container, which is based on the M2TS format.

[

16

]

[

17

]

The

Blu-ray Disc Association

calls it "

BDAV MPEG-2 transport stream

".

[

11

]

Blu-ray Disc video uses these modified MPEG-2 transport streams, compared to DVD's program streams that don't have the extra transport overhead.

There is also the BDAV (Blu-ray Disc Audio/Visual) format, the consumer-oriented alternative to the BDMV format used for movie releases. The BDAV format is used on

Blu-ray Disc recordable

for audio/video recording.

[

17

]

[

d

]

Blu-ray Disc employs the MPEG-2 transport stream recording method. This enables transport streams of a BDAV converted digital broadcast to be recorded as they are with minimal alteration of the packets.

[

12

]

It also enables simple stream cut style editing of a BDAV converted digital broadcast that is recorded as is and where the data can be edited just by discarding unwanted packets from the stream. Although it is quite natural, a function for high-speed and easy-to-use retrieval is built in.

[

12

]

[

19

]

See also

[

edit

]

MPEG media transport

(MMT)

Real-time Transport Protocol

(RTP)

Unidirectional Lightweight Encapsulation

(ULE)

Notes

[

edit

]

^

The

Blu-ray

format does not require CBR.

^

Forward error correction

is added by

ISDB

and

DVB

(16 bytes) and

ATSC

(20 bytes),

[

7

]

while the

M2TS

format prefixes packets with a 4-byte copyright and timestamp tag.

^

Possibly an abbreviation for "Transport stream on disc".

^

Filename extension

.m2ts

is used on Blu-ray Disc video files which contain an incompatible BDAV MPEG-2 transport stream due to the four additional octets added to every packet.

[

11

]

[

18

]

References

[

edit

]

^

"TVNT.net - Le forum de la TNT • [Topic Unique] Akira DHB-B31HDR - Double tuner enregistreur TNT HD - MKV - DIVX - DTS : Les adaptateurs pour recevoir la TNT gratuite en SD ou HD"

.

www.tvnt.net

.

^

MIME Type Registration of RTP Payload Formats

. July 2003.

doi

:

10.17487/RFC3555

.

RFC

3555

.

^

"mpeg2TransportStream"

.

Apple Developer Documentation

.

Apple Inc

.

^

a

b

c

ITU-T (October 2014).

"Recommendation H.222.0 (10/14)"

.

^

MPEG-2 Encoding Family

(Full draft). Sustainability of Digital Formats. Washington, D.C.: Library of Congress. 14 February 2012

. Retrieved

13 December

2021

.

Licenses pertain to tools and not to streams or files per se.

^

"MPEG-2 Transport Stream"

.

AfterDawn.com

. Retrieved

8 June

2010

.

^

"ATSC transmission"

.

Broadcastengineering.com

. 20 June 2005

. Retrieved

17 May

2012

.

^

"MPEG Systems FAQ"

.

Mpeg.chiariglione.org

. Archived from

the original

on 9 May 2012

. Retrieved

17 May

2012

.

^

"ATSC MPEG Transport Stream Monitor"

.

Tek.com

. Archived from

the original

on 3 December 2012

. Retrieved

17 May

2012

.

^

A Guide to MPEG Fundamentals and Protocol Analysis

(PDF)

, Tektronix, p. 37

, retrieved

23 April

2020

^

a

b

c

BD ROM – Audio Visual Application Format Specifications

(PDF)

, Blu-ray Disc Association, March 2005, pp.

15–

16, archived from

the original

(PDF)

on 3 November 2020

, retrieved

26 July

2009

^

a

b

c

BD-RE – Audiovisual Application Format Specification for BD-RE 2.1

(PDF)

, Blu-ray Disc Association, March 2008, archived from

the original

(PDF)

on 6 February 2009

^

"How MPEG-TS works"

.

Forum.videohelp.com

. Retrieved

17 May

2012

.

[

self-published source?

]

^

"Steve Mullen, M2TS primer"

.

Dvinfo.net

.

^

"Working with JVC Everio MOD & TOD files"

. Archived from the original on 23 October 2008.

^

Afterdawn.com

Glossary – BD-MV (Blu-ray Movie) and BDAV container

Archived

18 February 2009 at the

Wayback Machine

, Retrieved on 26 July 2009

^

a

b

Afterdawn.com

Glossary – BDAV container

, Retrieved on 26 July 2009

^

Videohelp.com

What is Blu-ray Disc and HD DVD?

Archived

24 December 2009 at the

Wayback Machine

, Retrieved on 26 July 2009

^

Blu-ray Disc Association (August 2004)

Blu-ray Disc Format, White paper

Archived

12 June 2009 at the

Wayback Machine

(PDF) Page 22, Retrieved on 28 July 2009

External links

[

edit

]

ITU-T H.222.0 | ISO/IEC 13818-1 Systems Spec Documents

Latest free copy of the spec, August 2023

MPEG-4 Systems FAQ

Archived

27 October 2020 at the

Wayback Machine

TSDuck

– Free open-source tool to manipulate MPEG transport streams.

v

t

e

Multimedia

compression

and

container

formats

Video

compression

ISO

,

IEC

,

MPEG

DV

MJPEG

Motion JPEG 2000

MPEG-1

MPEG-2

Part 2

MPEG-4

Part 2 / ASP

Part 10 / AVC

Part 33 / IVC

MPEG-H

Part 2 / HEVC

MPEG-I

Part 3 / VVC

MPEG-5

Part 1 / EVC

Part 2 / LCEVC

ITU-T

,

VCEG

H.120

H.261

H.262

H.263

H.264 / AVC

H.265 / HEVC

H.266 / VVC

H.267 / Enhanced Compression Model

SMPTE

VC-1

VC-2

VC-3

VC-5

VC-6

TrueMotion

and AOMedia

TrueMotion S

VP3

VP6

VP7

VP8

VP9

AV1

AV2

Chinese Standard

AVS1 P2/AVS+

(GB/T 20090.2/16)

AVS2 P2

(GB/T 33475.2,GY/T 299.1)

HDR Vivid(GY/T 358)

AVS3 P2(GY/T 368)

Others

Apple Video

AVS

Bink

Cinepak

Daala

DVI

FFV1

Huffyuv

Indeo

Lagarith

Microsoft Video 1

MSU Lossless

OMS Video

Pixlet

ProRes

422

4444

QuickTime

Animation

Graphics

RealVideo

RTVideo

SheerVideo

Smacker

Sorenson Video/Spark

Theora

Thor

Ut

WMV

XEB

YULS

Audio

compression

ISO

,

IEC

,

MPEG

MPEG-1 Layer II

Multichannel

MPEG-1 Layer I

MPEG-1 Layer III (MP3)

AAC

HE-AAC

AAC-LD

MPEG Surround

MPEG-4 ALS

MPEG-4 SLS

MPEG-4 DST

MPEG-4 HVXC

MPEG-4 CELP

MPEG-D USAC

MPEG-H 3D Audio

ITU-T

G.711

A-law

µ-law

G.718

G.719

G.722

G.722.1

G.722.2

G.723

G.723.1

G.726

G.728

G.729

G.729.1

IETF

Opus

iLBC

Speex

Vorbis

FLAC

3GPP

AMR

AMR-WB

AMR-WB+

EVRC

EVRC-B

EVS

GSM-HR

GSM-FR

GSM-EFR

ETSI

AC-3

AC-4

DTS

Bluetooth SIG

SBC

LC3

Chinese Standard

AVS1 P10

(GB/T 20090.10)

AVS2 P3

(GB/T 33475.3)

Audio Vivid

(GY/T 363)

DRA

(GB/T 22726)

ExAC(SJ/T 11299.4)

Others

ACELP

ALAC

Asao

ATRAC

CELT

Codec 2

iSAC

Lyra

MELP

Monkey's Audio

MT9

Musepack

OptimFROG

OSQ

QCELP

RCELP

RealAudio

SD2

SHN

SILK

Siren

SMV

SVOPC

TTA

True Audio

TwinVQ

VMR-WB

VSELP

WavPack

WMA

MQA

aptX

aptX HD

aptX Low Latency

aptX Adaptive

LDAC

LHDC

LLAC

TrueHD

Image

compression

IEC

,

ISO

,

IETF

,

W3C

,

ITU-T

,

JPEG

CCITT Group 4

GIF

HEIC / HEIF

HEVC

JBIG

JBIG2

JPEG

JPEG 2000

JPEG-LS

JPEG XL

JPEG XR

JPEG XS

JPEG XT

PNG

APNG

TIFF

TIFF/EP

TIFF/IT

Others

AV1

AVIF

BPG

DjVu

EXR

FLIF

ICER

MNG

PGF

QOI

QTVR

WBMP

WebP

Containers

ISO

,

IEC

MPEG-ES

MPEG-PES

MPEG-PS

MPEG-TS

ISO/IEC base media file format

MPEG-4 Part 14

(MP4)

Motion JPEG 2000

MPEG-21 Part 9

MPEG media transport

ITU-T

H.222.0

T.802

IETF

RTP

Ogg

Matroska

SMPTE

GXF

MXF

Others

3GP and 3G2

AMV

ASF

AIFF

AVI

AU

BPG

Bink

Smacker

BMP

DivX Media Format

EVO

Flash Video

HEIF

IFF

M2TS

Matroska

WebM

QuickTime File Format

RatDVD

RealMedia

RIFF

WAV

MOD and TOD

VOB, IFO and BUP

Collaborations

NETVC

MPEG LA

Alliance for Open Media

Methods

Entropy

Arithmetic

Huffman

Modified

LPC

ACELP

CELP

LSP

WLPC

Lossless

Lossy

LZ

DEFLATE

LZW

PCM

A-law

µ-law

ADPCM

DPCM

Transforms

DCT

FFT

MDCT

Wavelet

Daubechies

DWT

Lists

Comparison of audio coding formats

Comparison of video codecs

List of codecs

See

Compression methods

for techniques and

Compression software

for codecs

v

t

e

MPEG (Moving Picture Experts Group)

MPEG-1

2

3

4

7

21

A

B

C

D

E

G

V

M

U

H

I

5

MPEG-1 Parts

Part 1: Systems

Program stream

Part 2: Video

based on

H.261

Part 3: Audio

Layer I

Layer II

Layer III

MPEG-2 Parts

Part 1: Systems (H.222.0)

Transport stream

Program stream

Part 2: Video (H.262)

Part 3: Audio

Layer I

Layer II

Layer III

MPEG Multichannel

Part 6: DSM CC

Part 7: Advanced Audio Coding

MPEG-4 Parts

Part 2: Video

based on

H.263

Part 3: Audio

Part 6: DMIF

Part 10: Advanced Video Coding (H.264)

Part 11: Scene description

Part 12: ISO base media file format

Part 14: MP4 file format

Part 17: Streaming text format

Part 20: LASeR

Part 22: Open Font Format

Part 33: Internet Video Coding

MPEG-7 Parts

Part 2: Description definition language

MPEG-21 Parts

Parts 2, 3 and 9: Digital Item

Part 5: Rights Expression Language

MPEG-D Parts

Part 1: MPEG Surround

Part 3: Unified Speech and Audio Coding

MPEG-G Parts

Part 1: Transport and Storage of Genomic Information

Part 2: Coding of Genomic Information

Part 3: APIs

Part 4: Reference Software

Part 5: Conformance

MPEG-H Parts

Part 1: MPEG media transport

Part 2: High Efficiency Video Coding (H.265)

Part 3: MPEG-H 3D Audio

Part 12: High Efficiency Image File Format

MPEG-I Parts

Part 3: Versatile Video Coding (H.266)

MPEG-5 Parts

Part 1: Essential Video Coding

Part 2: Low Complexity Enhancement Video Coding

Other

MPEG-DASH

Retrieved from "

https://en.wikipedia.org/w/index.php?title=MPEG_transport_stream&oldid=1324443299

"

Categories

:

ATSC

Digital container formats

MPEG-2

ITU-T recommendations

Hidden categories:

All articles with self-published sources

Articles with self-published sources from May 2012

CS1: unfit URL

Webarchive template wayback links

Articles with short description

Short description matches Wikidata

Use dmy dates from August 2019