# MPEG Transport Stream

> Source: https://en.wikipedia.org/wiki/MPEG_transport_stream
> Fetched: 2026-02-01T11:46:49.728524+00:00
> Content-Hash: 5a2f003d90ce04ea
> Type: html

---

From Wikipedia, the free encyclopedia

Digital video format used for storage network transmission

"MPEG-2 Systems" redirects here. For the program stream technology also specified in the MPEG-2 Systems standard, see [MPEG program stream](/wiki/MPEG_program_stream "MPEG program stream").

".ts" redirects here. For the programming language associated with the .ts file extension, see [TypeScript](/wiki/TypeScript "TypeScript"). For the translation source files in Qt, see [Qt Linguist](/wiki/Qt_Linguist "Qt Linguist").

"SPTS" redirects here. For the American television production and distribution arm of [Sony](/wiki/Sony "Sony"), see [Sony Pictures Television](/wiki/Sony_Pictures_Television "Sony Pictures Television").

Not to be confused with [.m2ts](/wiki/.m2ts ".m2ts").

MPEG Transport Stream  
---  
[Filename extension](/wiki/Filename_extension "Filename extension")|  .ts, .tsv, .tsa, .m2t[1]  
[Internet media type](/wiki/Media_type "Media type")|  video/MP2T[2]  
[Uniform Type Identifier (UTI)](/wiki/Uniform_Type_Identifier "Uniform Type Identifier")| public.mpeg-2-transport-stream[3]  
Developed by| [MPEG](/wiki/MPEG "MPEG")  
Initial release| 10 July 1995; 30 years ago [1995-07-10](4)  
[Latest release](/wiki/Software_release_life_cycle "Software release life cycle")| ISO/IEC 13818-1:2022  
September 2022; 3 years ago (2022-09)  
Type of format| [Container format](/wiki/Container_format "Container format")  
[Container for](/wiki/Container_format "Container format")| Audio, video, data  
Extended to| [M2TS](/wiki/M2TS "M2TS"), [TOD](/wiki/MOD_and_TOD_\(video_format\) "MOD and TOD \(video format\)")  
[Standard](/wiki/International_standard "International standard")| ISO/IEC 13818-1, ITU-T Recommendation H.222.0[4]  
[Open format](/wiki/Open_file_format "Open file format")?| Yes  
[Free format](/wiki/Open_standard#Comparison_of_definitions "Open standard")?| Yes[5]  
  
**MPEG transport stream** (**[MPEG](/wiki/MPEG "MPEG")-TS**, **MTS**) or simply **transport stream** (**TS**) is a standard [digital container format](/wiki/Digital_container_format "Digital container format") for transmission and storage of [audio](/wiki/Digital_audio "Digital audio"), [video](/wiki/Digital_video "Digital video"), and [Program and System Information Protocol](/wiki/Program_and_System_Information_Protocol "Program and System Information Protocol") (PSIP) data.[6] It is used in broadcast systems such as [DVB](/wiki/DVB "DVB"), [ATSC](/wiki/ATSC "ATSC") and [IPTV](/wiki/IPTV "IPTV").

Transport stream specifies a container format encapsulating [packetized elementary streams](/wiki/Packetized_elementary_stream "Packetized elementary stream"), with [error correction](/wiki/Error_correction "Error correction") and [synchronization pattern](/wiki/Synchronization_pattern "Synchronization pattern") features for maintaining transmission integrity when the [communication channel](/wiki/Communication_channel "Communication channel") carrying the stream is [degraded](/wiki/Degradation_\(telecommunications\) "Degradation \(telecommunications\)").

Transport streams differ from the similarly named [MPEG program stream](/wiki/MPEG_program_stream "MPEG program stream") in several important ways: program streams are designed for reasonably reliable media, such as discs (like [DVDs](/wiki/DVD "DVD")), while transport streams are designed for less [reliable](/wiki/Reliability_\(computer_networking\) "Reliability \(computer networking\)") transmission, namely [terrestrial](/wiki/Terrestrial_television "Terrestrial television") or [satellite broadcast](/wiki/Satellite_television "Satellite television"). Further, a transport stream may carry multiple programs.

Transport stream is specified in _[MPEG-2](/wiki/MPEG-2 "MPEG-2") Part 1, Systems_, formally known as _[ISO/IEC](/wiki/ISO/IEC "ISO/IEC") standard 13818-1_ or _[ITU-T](/wiki/ITU-T "ITU-T") Rec. H.222.0_.[4]

## Overview

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=1 "Edit section: Overview")]

[](/wiki/File:MPEG_Transport_Stream_HL.svg)Multiple MPEG programs are combined then sent to a transmitting antenna. The receiver parses and decodes one of the streams.

A transport stream encapsulates a number of other substreams, often [packetized elementary streams](/wiki/Packetized_elementary_stream "Packetized elementary stream") (PESs) which in turn wrap the [main data stream](/wiki/Elementary_stream "Elementary stream") using the MPEG codec or any number of non-MPEG codecs (such as [AC3](/wiki/Dolby_Digital "Dolby Digital") or [DTS](/wiki/DTS_\(sound_system\) "DTS \(sound system\)") audio, and [MJPEG](/wiki/MJPEG "MJPEG") or [JPEG 2000](/wiki/JPEG_2000 "JPEG 2000") video), text and pictures for subtitles, tables identifying the streams, and even broadcaster-specific information such as an [electronic program guide](/wiki/Electronic_program_guide "Electronic program guide"). Many streams are often mixed together, such as several different television channels, or multiple [angles](/wiki/DVD-Video#Chapters_and_angles "DVD-Video") of a movie.

Each stream is chopped into (at most) 188-byte sections and interleaved together. Due to the tiny packet size, streams can be interleaved with less latency and greater error resilience compared to [program streams](/wiki/MPEG_program_stream "MPEG program stream") and other common containers such as [AVI](/wiki/Audio_Video_Interleave "Audio Video Interleave"), [MOV](/wiki/QuickTime_File_Format "QuickTime File Format")/[MP4](/wiki/MP4 "MP4"), and [MKV](/wiki/Matroska "Matroska"), which generally wrap each frame into one packet. This is particularly important for videoconferencing, where large frames may introduce unacceptable audio delay.

Transport streams tend to be broadcast as [constant bitrate](/wiki/Constant_bitrate "Constant bitrate") (CBR) and filled with padding bytes when not enough data exists.[a]

## Elements

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=2 "Edit section: Elements")]

### Packet

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=3 "Edit section: Packet")]

A [network packet](/wiki/Network_packet "Network packet") is the basic unit of data in a transport stream, and a transport stream is merely a sequence of packets. Each packet starts with a [sync byte](/wiki/Sync_byte "Sync byte") and a [header](/wiki/Header_\(computing\) "Header \(computing\)"), which may be followed by optional additional headers; the rest of the packet consists of [payload](/wiki/Payload_\(computing\) "Payload \(computing\)"). All header fields are read as [big-endian](/wiki/Endianness#Big-endian "Endianness"). Packets are 188 bytes in length, but the communication medium may add additional information.[b] The 188-byte packet size was originally chosen for compatibility with [Asynchronous Transfer Mode (ATM) systems](/wiki/Asynchronous_Transfer_Mode "Asynchronous Transfer Mode").[8][9]

### Programs

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=4 "Edit section: Programs")]

Transport stream has a concept of _programs_. Every program is described by a program map table (PMT). The elementary streams associated with that program have PIDs listed in the PMT. Another PID is associated with the PMT itself. For instance, a transport stream used in digital television might contain three programs to represent three television channels. Suppose each channel consists of one video stream, one or two audio streams, and any necessary metadata. A [receiver](/wiki/ATSC_tuner "ATSC tuner") wishing to decode one of the three channels merely has to decode the payloads of each PID associated with its program. It can discard the contents of all other PIDs. A transport stream with more than one program is referred to as a multi-program transport stream (MPTS). A single program transport stream is referred to as a single-program transport stream (SPTS).

### Program specific information

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=5 "Edit section: Program specific information")]

Main article: [Program-specific information](/wiki/Program-specific_information "Program-specific information")

There are 4 program-specific information (PSI) tables: program association (PAT), program map (PMT), conditional access (CAT), and network information (NIT). The MPEG-2 specification does not specify the format of the CAT and NIT.

### PCR

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=6 "Edit section: PCR")]

To enable a decoder to present synchronized content, such as audio tracks matching the associated video, at least once each 100 ms, a _program clock reference_ (PCR) is transmitted in the adaptation field of an MPEG-2 transport stream packet. The PID with the PCR for an MPEG-2 program is identified by the _pcr_pid_ value in the associated PMT. The value of the PCR, when properly used, is employed to generate a _system_timing_clock_ in the decoder. The system time clock (STC) decoder, when properly implemented, provides a highly accurate time base that is used to synchronize audio and video elementary streams. Timing in MPEG-2 references this clock. For example, the [presentation time stamp](/wiki/Presentation_time_stamp "Presentation time stamp") (PTS) is intended to be relative to the PCR. The first 33 bits are based on a 90 kHz clock. The last 9 bits are based on a 27 MHz clock. The maximum jitter permitted for the PCR is +/- 500 ns.

### Null packets

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=7 "Edit section: Null packets")]

Some transmission schemes, such as those in [ATSC](/wiki/ATSC "ATSC") and [DVB](/wiki/DVB "DVB"), impose strict constant bitrate requirements on the transport stream. In order to ensure that the stream maintains a constant bitrate, a multiplexer may need to insert some additional packets. The PID 0x1FFF is reserved for this purpose. The null packets have a payload that is filled with 0xFF, and the receiver is expected to ignore its contents.[10]

## M2TS

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=8 "Edit section: M2TS")]

Main article: [.m2ts](/wiki/.m2ts ".m2ts")

Transport Stream was originally designed for broadcast. Later, it was adapted for use with digital video cameras, recorders and players by adding a 4-byte timecode (TC) field to the standard 188-byte packets, resulting in a 192-byte packet.[11][12] This is what is informally called [M2TS](/wiki/M2TS "M2TS") stream, commonly found in [HDV](/wiki/HDV "HDV") cameras. The timecode allows quick access to any part of the stream, either from a media player or from a non-linear video editing system.[13]

### Use in digital video cameras

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=9 "Edit section: Use in digital video cameras")]

JVC called M2TS "[TOD](/wiki/MOD_and_TOD "MOD and TOD")"[c] when used in HDD-based camcorders like [GZ-HD7](/wiki/JVC_GZ-HD7 "JVC GZ-HD7").[14][15] It is also used to synchronize video streams from several cameras in a [multiple-camera setup](/wiki/Multiple-camera_setup "Multiple-camera setup").

### Use in Blu-ray

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=10 "Edit section: Use in Blu-ray")]

Blu-ray Disc video titles authored with menu support are in the [Blu-ray Disc Movie](/wiki/Blu-ray_Disc_Movie "Blu-ray Disc Movie") (BDMV) format and contain audio, video, and other streams in a BDAV container, which is based on the M2TS format.[16][17] The [Blu-ray Disc Association](/wiki/Blu-ray_Disc_Association "Blu-ray Disc Association") calls it "[BDAV MPEG-2 transport stream](/wiki/.m2ts ".m2ts")".[11] Blu-ray Disc video uses these modified MPEG-2 transport streams, compared to DVD's program streams that don't have the extra transport overhead.

There is also the BDAV (Blu-ray Disc Audio/Visual) format, the consumer-oriented alternative to the BDMV format used for movie releases. The BDAV format is used on [Blu-ray Disc recordable](/wiki/Blu-ray_Disc_recordable "Blu-ray Disc recordable") for audio/video recording.[17][d] Blu-ray Disc employs the MPEG-2 transport stream recording method. This enables transport streams of a BDAV converted digital broadcast to be recorded as they are with minimal alteration of the packets.[12] It also enables simple stream cut style editing of a BDAV converted digital broadcast that is recorded as is and where the data can be edited just by discarding unwanted packets from the stream. Although it is quite natural, a function for high-speed and easy-to-use retrieval is built in.[12][19]

## See also

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=11 "Edit section: See also")]

- [MPEG media transport](/wiki/MPEG_media_transport "MPEG media transport") (MMT)
- [Real-time Transport Protocol](/wiki/Real-time_Transport_Protocol "Real-time Transport Protocol") (RTP)
- [Unidirectional Lightweight Encapsulation](/wiki/Unidirectional_Lightweight_Encapsulation "Unidirectional Lightweight Encapsulation") (ULE)

## Notes

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=12 "Edit section: Notes")]

  1. **^** The [Blu-ray](/wiki/Blu-ray "Blu-ray") format does not require CBR.
  2. **^** [Forward error correction](/wiki/Forward_error_correction "Forward error correction") is added by [ISDB](/wiki/ISDB "ISDB") and [DVB](/wiki/DVB "DVB") (16 bytes) and [ATSC](/wiki/ATSC "ATSC") (20 bytes),[7] while the [M2TS](/wiki/M2TS "M2TS") format prefixes packets with a 4-byte copyright and timestamp tag.
  3. **^** Possibly an abbreviation for "Transport stream on disc".
  4. **^** Filename extension [.m2ts](/wiki/.m2ts ".m2ts") is used on Blu-ray Disc video files which contain an incompatible BDAV MPEG-2 transport stream due to the four additional octets added to every packet.[11][18]

## References

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=13 "Edit section: References")]

  1. **^** ["TVNT.net - Le forum de la TNT • [Topic Unique] Akira DHB-B31HDR - Double tuner enregistreur TNT HD - MKV - DIVX - DTS : Les adaptateurs pour recevoir la TNT gratuite en SD ou HD"](http://www.tvnt.net/forum/akira-dhb-b31hdr-double-tuner-enregistreur-tnt-hd-mkv-divx-dts-t26336.html). _www.tvnt.net_.
  2. **^** [_MIME Type Registration of RTP Payload Formats_](https://www.rfc-editor.org/rfc/rfc3555). July 2003. [doi](/wiki/Doi_\(identifier\) "Doi \(identifier\)"):[10.17487/RFC3555](https://doi.org/10.17487%2FRFC3555). [RFC](/wiki/Request_for_Comments "Request for Comments") [3555](https://datatracker.ietf.org/doc/html/rfc3555).
  3. **^** ["mpeg2TransportStream"](https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3551535-mpeg2transportstream). _Apple Developer Documentation_. [Apple Inc](/wiki/Apple_Inc "Apple Inc").
  4. ^ _**a**_ _**b**_ _**c**_ ITU-T (October 2014). ["Recommendation H.222.0 (10/14)"](http://www.itu.int/rec/T-REC-H.222.0-201410-I/en).
  5. **^** [_MPEG-2 Encoding Family_](https://www.loc.gov/preservation/digital/formats/fdd/fdd000335.shtml) (Full draft). Sustainability of Digital Formats. Washington, D.C.: Library of Congress. 14 February 2012. Retrieved 13 December 2021. "Licenses pertain to tools and not to streams or files per se."
  6. **^** ["MPEG-2 Transport Stream"](http://www.afterdawn.com/glossary/term.cfm/mpeg2_transport_stream). _AfterDawn.com_. Retrieved 8 June 2010.
  7. **^** ["ATSC transmission"](http://broadcastengineering.com/infrastructure/Atsc-transmission-digital-20050620/). _Broadcastengineering.com_. 20 June 2005. Retrieved 17 May 2012.
  8. **^** ["MPEG Systems FAQ"](https://web.archive.org/web/20120509115745/http://mpeg.chiariglione.org/faq/mp2-sys/mp2-sys.htm#mp2-12). _Mpeg.chiariglione.org_. Archived from [the original](http://mpeg.chiariglione.org/faq/mp2-sys/mp2-sys.htm#mp2-12) on 9 May 2012. Retrieved 17 May 2012.
  9. **^** ["ATSC MPEG Transport Stream Monitor"](https://web.archive.org/web/20121203054329/http://www.tek.com/datasheet/atsc-mpeg-transport-stream-monitor). _Tek.com_. Archived from [the original](http://www.tek.com/datasheet/atsc-mpeg-transport-stream-monitor) on 3 December 2012. Retrieved 17 May 2012.
  10. **^** [_A Guide to MPEG Fundamentals and Protocol Analysis_](http://www.img.lx.it.pt/~fp/cav/Additional_material/MPEG2_overview.pdf) (PDF), Tektronix, p. 37, retrieved 23 April 2020
  11. ^ _**a**_ _**b**_ _**c**_ [_BD ROM – Audio Visual Application Format Specifications_](https://web.archive.org/web/20201103093732/http://www.blu-raydisc.com/Assets/Downloadablefile/2b_bdrom_audiovisualapplication_0305-12955-15269.pdf) (PDF), Blu-ray Disc Association, March 2005, pp. 15–16, archived from [the original](http://www.blu-raydisc.com/Assets/Downloadablefile/2b_bdrom_audiovisualapplication_0305-12955-15269.pdf) (PDF) on 3 November 2020, retrieved 26 July 2009
  12. ^ _**a**_ _**b**_ _**c**_ [_BD-RE – Audiovisual Application Format Specification for BD-RE 2.1_](https://web.archive.org/web/20090206111829/http://www.blu-raydisc.com/Assets/Downloadablefile/BD-RE_Part3_V2.1_WhitePaper_080406-15271.pdf) (PDF), Blu-ray Disc Association, March 2008, archived from [the original](http://www.blu-raydisc.com/Assets/Downloadablefile/BD-RE_Part3_V2.1_WhitePaper_080406-15271.pdf) (PDF) on 6 February 2009
  13. **^** ["How MPEG-TS works"](http://forum.videohelp.com/threads/306126-HFS10-AVCHD-how-to-maintain-quality?p=1881643&viewfull=1#post1881643). _Forum.videohelp.com_. Retrieved 17 May 2012.[_[self-published source?](/wiki/Wikipedia:Verifiability#Self-published_sources "Wikipedia:Verifiability")_]
  14. **^** ["Steve Mullen, M2TS primer"](http://dvinfo.net/conf/showthread.php?t=105486). _Dvinfo.net_.
  15. **^** ["Working with JVC Everio MOD & TOD files"](https://web.archive.org/web/20081023045657/http://www.avchduser.com/articles/JVC_Everio_mod_files.jsp). Archived from the original on 23 October 2008.
  16. **^** Afterdawn.com [Glossary – BD-MV (Blu-ray Movie) and BDAV container](http://www.afterdawn.com/glossary/terms/bd-mv.cfm) [Archived](https://web.archive.org/web/20090218234755/http://www.afterdawn.com/glossary/terms/bd-mv.cfm) 18 February 2009 at the [Wayback Machine](/wiki/Wayback_Machine "Wayback Machine"), Retrieved on 26 July 2009
  17. ^ _**a**_ _**b**_ Afterdawn.com [Glossary – BDAV container](http://www.afterdawn.com/glossary/terms/bdav.cfm), Retrieved on 26 July 2009
  18. **^** Videohelp.com [What is Blu-ray Disc and HD DVD?](http://www.videohelp.com/hd) [Archived](https://web.archive.org/web/20091224035325/http://www.videohelp.com/hd) 24 December 2009 at the [Wayback Machine](/wiki/Wayback_Machine "Wayback Machine"), Retrieved on 26 July 2009
  19. **^** Blu-ray Disc Association (August 2004) [Blu-ray Disc Format, White paper](http://www.blu-raydisc.com/Assets/Downloadablefile/general_bluraydiscformat-15263.pdf) [Archived](https://web.archive.org/web/20090612042130/http://www.blu-raydisc.com/Assets/Downloadablefile/general_bluraydiscformat-15263.pdf) 12 June 2009 at the [Wayback Machine](/wiki/Wayback_Machine "Wayback Machine") (PDF) Page 22, Retrieved on 28 July 2009

## External links

[[edit](/w/index.php?title=MPEG_transport_stream&action=edit&section=14 "Edit section: External links")]

- [ITU-T H.222.0 | ISO/IEC 13818-1 Systems Spec Documents](http://www.itu.int/rec/T-REC-H.222.0)
- [Latest free copy of the spec, August 2023](https://www.itu.int/rec/T-REC-H.222.0-202308-S/en)
- [MPEG-4 Systems FAQ](http://mpeg.chiariglione.org/faq/mp4-sys/mp4-sys.htm) [Archived](https://web.archive.org/web/20201027160601/http://mpeg.chiariglione.org/faq/mp4-sys/mp4-sys.htm) 27 October 2020 at the [Wayback Machine](/wiki/Wayback_Machine "Wayback Machine")
- [TSDuck](https://tsduck.io/) – Free open-source tool to manipulate MPEG transport streams.

- [v](/wiki/Template:Compression_formats "Template:Compression formats")
- [t](/wiki/Template_talk:Compression_formats "Template talk:Compression formats")
- [e](/wiki/Special:EditPage/Template:Compression_formats "Special:EditPage/Template:Compression formats")

[Multimedia](/wiki/Multimedia "Multimedia") [compression](/wiki/Data_compression "Data compression") and [container](/wiki/Container_format_\(computing\) "Container format \(computing\)") formats  
---  
[Video  
compression](/wiki/Video_coding_format "Video coding format")| | [ISO](/wiki/International_Organization_for_Standardization "International Organization for Standardization"), [IEC](/wiki/International_Electrotechnical_Commission "International Electrotechnical Commission"),
[MPEG](/wiki/Moving_Picture_Experts_Group "Moving Picture Experts Group")|

- [DV](/wiki/DV_\(video_format\) "DV \(video format\)")
- [MJPEG](/wiki/Motion_JPEG "Motion JPEG")
- [Motion JPEG 2000](/wiki/Motion_JPEG_2000 "Motion JPEG 2000")
- [MPEG-1](/wiki/MPEG-1 "MPEG-1")
- [MPEG-2](/wiki/MPEG-2 "MPEG-2")
  - [Part 2](/wiki/H.262/MPEG-2_Part_2 "H.262/MPEG-2 Part 2")
- [MPEG-4](/wiki/MPEG-4 "MPEG-4")
  - [Part 2 / ASP](/wiki/MPEG-4_Part_2 "MPEG-4 Part 2")
  - [Part 10 / AVC](/wiki/H.264/MPEG-4_AVC "H.264/MPEG-4 AVC")
  - [Part 33 / IVC](/wiki/MPEG-4_IVC "MPEG-4 IVC")
- [MPEG-H](/wiki/MPEG-H "MPEG-H")
  - [Part 2 / HEVC](/wiki/High_Efficiency_Video_Coding "High Efficiency Video Coding")
- [MPEG-I](/w/index.php?title=MPEG-I&action=edit&redlink=1 "MPEG-I \(page does not exist\)")
  - [Part 3 / VVC](/wiki/Versatile_Video_Coding "Versatile Video Coding")
- [MPEG-5](/wiki/MPEG-5 "MPEG-5")
  - [Part 1 / EVC](/wiki/Essential_Video_Coding "Essential Video Coding")
  - [Part 2 / LCEVC](/wiki/LCEVC "LCEVC")

---|---  
[ITU-T](/wiki/ITU-T "ITU-T"), [VCEG](/wiki/Video_Coding_Experts_Group "Video Coding Experts Group")|

- [H.120](/wiki/H.120 "H.120")
- [H.261](/wiki/H.261 "H.261")
- [H.262](/wiki/H.262/MPEG-2_Part_2 "H.262/MPEG-2 Part 2")
- [H.263](/wiki/H.263 "H.263")
- [H.264 / AVC](/wiki/Advanced_Video_Coding "Advanced Video Coding")
- [H.265 / HEVC](/wiki/High_Efficiency_Video_Coding "High Efficiency Video Coding")
- [H.266 / VVC](/wiki/Versatile_Video_Coding "Versatile Video Coding")
- [H.267 / Enhanced Compression Model](/w/index.php?title=H.267&action=edit&redlink=1 "H.267 \(page does not exist\)")

[SMPTE](/wiki/Society_of_Motion_Picture_and_Television_Engineers "Society of Motion Picture and Television Engineers")|

- [VC-1](/wiki/VC-1 "VC-1")
- [VC-2](/wiki/Dirac_\(video_compression_format\) "Dirac \(video compression format\)")
- [VC-3](/wiki/Avid_DNxHD "Avid DNxHD")
- [VC-5](/wiki/CineForm "CineForm")
- [VC-6](/wiki/VC-6 "VC-6")

[TrueMotion](/wiki/On2_Technologies "On2 Technologies") and AOMedia|

- [TrueMotion S](/wiki/On2_Technologies#TrueMotion_S "On2 Technologies")
- [VP3](/wiki/VP3 "VP3")
- [VP6](/wiki/VP6 "VP6")
- [VP7](/wiki/VP7 "VP7")
- [VP8](/wiki/VP8 "VP8")
- [VP9](/wiki/VP9 "VP9")
- [AV1](/wiki/AV1 "AV1")
- [AV2](/wiki/AV2 "AV2")

Chinese Standard|

- [AVS1 P2/AVS+](/wiki/Audio_Video_Standard#First_generation "Audio Video Standard")(GB/T 20090.2/16)
- [AVS2 P2](/wiki/Audio_Video_Standard#Second_generation "Audio Video Standard")(GB/T 33475.2,GY/T 299.1)
  - HDR Vivid(GY/T 358)
- AVS3 P2(GY/T 368)

Others|

- [Apple Video](/wiki/Apple_Video "Apple Video")
- [AVS](/wiki/Audio_Video_Standard "Audio Video Standard")
- [Bink](/wiki/Bink_Video "Bink Video")
- [Cinepak](/wiki/Cinepak "Cinepak")
- [Daala](/wiki/Daala "Daala")
- [DVI](/wiki/Digital_Video_Interactive "Digital Video Interactive")
- [FFV1](/wiki/FFV1 "FFV1")
- [Huffyuv](/wiki/Huffyuv "Huffyuv")
- [Indeo](/wiki/Indeo "Indeo")
- [Lagarith](/wiki/Lagarith "Lagarith")
- [Microsoft Video 1](/wiki/Microsoft_Video_1 "Microsoft Video 1")
- [MSU Lossless](/wiki/MSU_Lossless_Video_Codec "MSU Lossless Video Codec")
- [OMS Video](/wiki/OMS_Video "OMS Video")
- [Pixlet](/wiki/Pixlet "Pixlet")
- [ProRes](/wiki/Apple_ProRes "Apple ProRes")
  - [422](/wiki/ProRes_422 "ProRes 422")
  - [4444](/wiki/ProRes_4444 "ProRes 4444")
- QuickTime
  - [Animation](/wiki/QuickTime_Animation "QuickTime Animation")
  - [Graphics](/wiki/QuickTime_Graphics "QuickTime Graphics")
- [RealVideo](/wiki/RealVideo "RealVideo")
- [RTVideo](/wiki/RTVideo "RTVideo")
- [SheerVideo](/wiki/SheerVideo "SheerVideo")
- [Smacker](/wiki/Smacker_video "Smacker video")
- [Sorenson Video/Spark](/wiki/Sorenson_Media "Sorenson Media")
- [Theora](/wiki/Theora "Theora")
- [Thor](/wiki/Thor_\(video_codec\) "Thor \(video codec\)")
- [Ut](/wiki/Ut_Video_Codec_Suite "Ut Video Codec Suite")
- [WMV](/wiki/Windows_Media_Video "Windows Media Video")
- [XEB](/wiki/RatDVD "RatDVD")
- [YULS](/wiki/YULS "YULS")

[Audio  
compression](/wiki/Audio_coding_format "Audio coding format")| | [ISO](/wiki/International_Organization_for_Standardization "International Organization for Standardization"), [IEC](/wiki/International_Electrotechnical_Commission "International Electrotechnical Commission"),  
[MPEG](/wiki/Moving_Picture_Experts_Group "Moving Picture Experts Group")|

- [MPEG-1 Layer II](/wiki/MPEG-1_Audio_Layer_II "MPEG-1 Audio Layer II")
  - [Multichannel](/wiki/MPEG_Multichannel "MPEG Multichannel")
- [MPEG-1 Layer I](/wiki/MPEG-1_Audio_Layer_I "MPEG-1 Audio Layer I")
- [MPEG-1 Layer III (MP3)](/wiki/MP3 "MP3")
- [AAC](/wiki/Advanced_Audio_Coding "Advanced Audio Coding")
  - [HE-AAC](/wiki/High-Efficiency_Advanced_Audio_Coding "High-Efficiency Advanced Audio Coding")
  - [AAC-LD](/wiki/AAC-LD "AAC-LD")
- [MPEG Surround](/wiki/MPEG_Surround "MPEG Surround")
- [MPEG-4 ALS](/wiki/Audio_Lossless_Coding "Audio Lossless Coding")
- [MPEG-4 SLS](/wiki/MPEG-4_SLS "MPEG-4 SLS")
- [MPEG-4 DST](/wiki/Super_Audio_CD#DST "Super Audio CD")
- [MPEG-4 HVXC](/wiki/Harmonic_Vector_Excitation_Coding "Harmonic Vector Excitation Coding")
- [MPEG-4 CELP](/wiki/Code-excited_linear_prediction "Code-excited linear prediction")
- [MPEG-D USAC](/wiki/Unified_Speech_and_Audio_Coding "Unified Speech and Audio Coding")
- [MPEG-H 3D Audio](/wiki/MPEG-H_3D_Audio "MPEG-H 3D Audio")

---|---  
[ITU-T](/wiki/ITU-T "ITU-T")|

- [G.711](/wiki/G.711 "G.711")
  - [A-law](/wiki/A-law_algorithm "A-law algorithm")
  - [µ-law](/wiki/%CE%9C-law_algorithm "Μ-law algorithm")
- [G.718](/wiki/G.718 "G.718")
- [G.719](/wiki/G.719 "G.719")
- [G.722](/wiki/G.722 "G.722")
- [G.722.1](/wiki/G.722.1 "G.722.1")
- [G.722.2](/wiki/Adaptive_Multi-Rate_Wideband "Adaptive Multi-Rate Wideband")
- [G.723](/wiki/G.723 "G.723")
- [G.723.1](/wiki/G.723.1 "G.723.1")
- [G.726](/wiki/G.726 "G.726")
- [G.728](/wiki/G.728 "G.728")
- [G.729](/wiki/G.729 "G.729")
- [G.729.1](/wiki/G.729.1 "G.729.1")

[IETF](/wiki/Internet_Engineering_Task_Force "Internet Engineering Task Force")|

- [Opus](/wiki/Opus_\(audio_format\) "Opus \(audio format\)")
- [iLBC](/wiki/Internet_Low_Bitrate_Codec "Internet Low Bitrate Codec")
- [Speex](/wiki/Speex "Speex")
- [Vorbis](/wiki/Vorbis "Vorbis")
- [FLAC](/wiki/FLAC "FLAC")

[3GPP](/wiki/3GPP "3GPP")|

- [AMR](/wiki/Adaptive_Multi-Rate_audio_codec "Adaptive Multi-Rate audio codec")
- [AMR-WB](/wiki/Adaptive_Multi-Rate_Wideband "Adaptive Multi-Rate Wideband")
- [AMR-WB+](/wiki/Extended_Adaptive_Multi-Rate_%E2%80%93_Wideband "Extended Adaptive Multi-Rate – Wideband")
- [EVRC](/wiki/Enhanced_Variable_Rate_Codec "Enhanced Variable Rate Codec")
- [EVRC-B](/wiki/Enhanced_Variable_Rate_Codec_B "Enhanced Variable Rate Codec B")
- [EVS](/wiki/Enhanced_Voice_Services "Enhanced Voice Services")
- [GSM-HR](/wiki/Half_Rate "Half Rate")
- [GSM-FR](/wiki/Full_Rate "Full Rate")
- [GSM-EFR](/wiki/Enhanced_full_rate "Enhanced full rate")

[ETSI](/wiki/ETSI "ETSI")|

- [AC-3](/wiki/Dolby_Digital "Dolby Digital")
- [AC-4](/wiki/Dolby_AC-4 "Dolby AC-4")
- [DTS](/wiki/DTS_\(sound_system\) "DTS \(sound system\)")

[Bluetooth SIG](/wiki/Bluetooth_Special_Interest_Group "Bluetooth Special Interest Group")|

- [SBC](/wiki/SBC_\(codec\) "SBC \(codec\)")
- [LC3](/wiki/LC3_\(codec\) "LC3 \(codec\)")

Chinese Standard|

- [AVS1 P10](/wiki/Audio_Video_Standard#First_generation "Audio Video Standard")(GB/T 20090.10)
- [AVS2 P3](/wiki/Audio_Video_Standard#Second_generation "Audio Video Standard")(GB/T 33475.3)
  - [Audio Vivid](/w/index.php?title=Audio_Vivid&action=edit&redlink=1 "Audio Vivid \(page does not exist\)")(GY/T 363)
- [DRA](/wiki/Dynamic_Resolution_Adaptation "Dynamic Resolution Adaptation")(GB/T 22726)
- ExAC(SJ/T 11299.4)

Others|

- [ACELP](/wiki/Algebraic_code-excited_linear_prediction "Algebraic code-excited linear prediction")
- [ALAC](/wiki/Apple_Lossless_Audio_Codec "Apple Lossless Audio Codec")
- [Asao](/wiki/Asao_\(codec\) "Asao \(codec\)")
- [ATRAC](/wiki/Adaptive_Transform_Acoustic_Coding "Adaptive Transform Acoustic Coding")
- [CELT](/wiki/CELT "CELT")
- [Codec 2](/wiki/Codec_2 "Codec 2")
- [iSAC](/wiki/Internet_Speech_Audio_Codec "Internet Speech Audio Codec")
- [Lyra](/wiki/Lyra_\(codec\) "Lyra \(codec\)")
- [MELP](/wiki/Mixed-excitation_linear_prediction "Mixed-excitation linear prediction")
- [Monkey's Audio](/wiki/Monkey%27s_Audio "Monkey's Audio")
- [MT9](/wiki/MT9 "MT9")
- [Musepack](/wiki/Musepack "Musepack")
- [OptimFROG](/wiki/OptimFROG "OptimFROG")
- [OSQ](/wiki/Original_Sound_Quality "Original Sound Quality")
- [QCELP](/wiki/Qualcomm_code-excited_linear_prediction "Qualcomm code-excited linear prediction")
- [RCELP](/wiki/Relaxed_code-excited_linear_prediction "Relaxed code-excited linear prediction")
- [RealAudio](/wiki/RealAudio "RealAudio")
- [SD2](/wiki/Avid_Audio#Sound_Designer_file_formats "Avid Audio")
- [SHN](/wiki/Shorten_file_format "Shorten file format")
- [SILK](/wiki/SILK "SILK")
- [Siren](/wiki/Siren_\(codec\) "Siren \(codec\)")
- [SMV](/wiki/Selectable_Mode_Vocoder "Selectable Mode Vocoder")
- [SVOPC](/wiki/SVOPC "SVOPC")
- TTA
  - True Audio
- [TwinVQ](/wiki/TwinVQ "TwinVQ")
- [VMR-WB](/wiki/Variable-Rate_Multimode_Wideband "Variable-Rate Multimode Wideband")
- [VSELP](/wiki/Vector_sum_excited_linear_prediction "Vector sum excited linear prediction")
- [WavPack](/wiki/WavPack "WavPack")
- [WMA](/wiki/Windows_Media_Audio "Windows Media Audio")
- [MQA](/wiki/Master_Quality_Authenticated "Master Quality Authenticated")
- [aptX](/wiki/AptX "AptX")
- [aptX HD](/wiki/AptX#aptX_HD "AptX")
- [aptX Low Latency](/wiki/AptX#aptX_Low_Latency "AptX")
- [aptX Adaptive](/wiki/AptX#aptX_Adaptive "AptX")
- [LDAC](/wiki/LDAC_\(codec\) "LDAC \(codec\)")
- [LHDC](/wiki/LHDC_\(codec\) "LHDC \(codec\)")
- [LLAC](/wiki/LHDC_\(codec\)#LLAC "LHDC \(codec\)")
- [TrueHD](/wiki/Dolby_TrueHD "Dolby TrueHD")

[Image  
compression](/wiki/Image_compression "Image compression")| | [IEC](/wiki/International_Electrotechnical_Commission "International Electrotechnical Commission"), [ISO](/wiki/International_Organization_for_Standardization "International Organization for Standardization"), [IETF](/wiki/Internet_Engineering_Task_Force "Internet Engineering Task Force"),
[W3C](/wiki/World_Wide_Web_Consortium "World Wide Web Consortium"), [ITU-T](/wiki/ITU-T "ITU-T"), [JPEG](/wiki/Joint_Photographic_Experts_Group "Joint Photographic Experts Group")|

- [CCITT Group 4](/wiki/Group_4_compression "Group 4 compression")
- [GIF](/wiki/GIF "GIF")
- [HEIC / HEIF](/wiki/High_Efficiency_Image_File_Format#HEIC:_HEVC_in_HEIF "High Efficiency Image File Format")
- [HEVC](/wiki/High_Efficiency_Video_Coding#Main_Still_Picture "High Efficiency Video Coding")
- [JBIG](/wiki/JBIG "JBIG")
- [JBIG2](/wiki/JBIG2 "JBIG2")
- [JPEG](/wiki/JPEG "JPEG")
- [JPEG 2000](/wiki/JPEG_2000 "JPEG 2000")
- [JPEG-LS](/wiki/JPEG-LS "JPEG-LS")
- [JPEG XL](/wiki/JPEG_XL "JPEG XL")
- [JPEG XR](/wiki/JPEG_XR "JPEG XR")
- [JPEG XS](/wiki/JPEG_XS "JPEG XS")
- [JPEG XT](/wiki/JPEG_XT "JPEG XT")
- [PNG](/wiki/PNG "PNG")
  - [APNG](/wiki/APNG "APNG")
- [TIFF](/wiki/TIFF "TIFF")
- [TIFF/EP](/wiki/TIFF/EP "TIFF/EP")
- [TIFF/IT](/wiki/TIFF/IT "TIFF/IT")

---|---  
Others|

- [AV1](/wiki/AV1 "AV1")
- [AVIF](/wiki/AVIF "AVIF")
- [BPG](/wiki/Better_Portable_Graphics "Better Portable Graphics")
- [DjVu](/wiki/DjVu "DjVu")
- [EXR](/wiki/OpenEXR "OpenEXR")
- [FLIF](/wiki/Free_Lossless_Image_Format "Free Lossless Image Format")
- [ICER](/wiki/ICER_\(file_format\) "ICER \(file format\)")
- [MNG](/wiki/Multiple-image_Network_Graphics "Multiple-image Network Graphics")
- [PGF](/wiki/Progressive_Graphics_File "Progressive Graphics File")
- [QOI](/wiki/QOI_\(image_format\) "QOI \(image format\)")
- [QTVR](/wiki/QuickTime_VR "QuickTime VR")
- [WBMP](/wiki/Wireless_Application_Protocol_Bitmap_Format "Wireless Application Protocol Bitmap Format")
- [WebP](/wiki/WebP "WebP")

[Containers](/wiki/Digital_container_format "Digital container format")| | [ISO](/wiki/International_Organization_for_Standardization "International Organization for Standardization"), [IEC](/wiki/International_Electrotechnical_Commission "International Electrotechnical Commission")|

- [MPEG-ES](/wiki/MPEG_elementary_stream "MPEG elementary stream")
  - [MPEG-PES](/wiki/Packetized_elementary_stream "Packetized elementary stream")
- [MPEG-PS](/wiki/MPEG_program_stream "MPEG program stream")
- MPEG-TS
- [ISO/IEC base media file format](/wiki/ISO/IEC_base_media_file_format "ISO/IEC base media file format")
- [MPEG-4 Part 14](/wiki/MPEG-4_Part_14 "MPEG-4 Part 14") (MP4)
- [Motion JPEG 2000](/wiki/Motion_JPEG_2000 "Motion JPEG 2000")
- [MPEG-21 Part 9](/wiki/MPEG-21 "MPEG-21")
- [MPEG media transport](/wiki/MPEG_media_transport "MPEG media transport")

---|---  
[ITU-T](/wiki/ITU-T "ITU-T")|

- [H.222.0](/wiki/MPEG-2#Systems "MPEG-2")
- [T.802](/wiki/Motion_JPEG_2000 "Motion JPEG 2000")

[IETF](/wiki/Internet_Engineering_Task_Force "Internet Engineering Task Force")|

- [RTP](/wiki/Real-time_Transport_Protocol "Real-time Transport Protocol")
- [Ogg](/wiki/Ogg "Ogg")
- [Matroska](/wiki/Matroska "Matroska")

[SMPTE](/wiki/Society_of_Motion_Picture_and_Television_Engineers "Society of Motion Picture and Television Engineers")|

- [GXF](/wiki/General_Exchange_Format "General Exchange Format")
- [MXF](/wiki/Material_Exchange_Format "Material Exchange Format")

Others|

- [3GP and 3G2](/wiki/3GP_and_3G2 "3GP and 3G2")
- [AMV](/wiki/AMV_video_format "AMV video format")
- [ASF](/wiki/Advanced_Systems_Format "Advanced Systems Format")
- [AIFF](/wiki/Audio_Interchange_File_Format "Audio Interchange File Format")
- [AVI](/wiki/Audio_Video_Interleave "Audio Video Interleave")
- [AU](/wiki/Au_file_format "Au file format")
- [BPG](/wiki/Better_Portable_Graphics "Better Portable Graphics")
- [Bink](/wiki/Bink_Video "Bink Video")
  - [Smacker](/wiki/Smacker_video "Smacker video")
- [BMP](/wiki/BMP_file_format "BMP file format")
- [DivX Media Format](/wiki/DivX#DivX_Media_Format_\(DMF\) "DivX")
- [EVO](/wiki/Enhanced_VOB "Enhanced VOB")
- [Flash Video](/wiki/Flash_Video "Flash Video")
- [HEIF](/wiki/High_Efficiency_Image_File_Format "High Efficiency Image File Format")
- [IFF](/wiki/Interchange_File_Format "Interchange File Format")
- [M2TS](/wiki/.m2ts ".m2ts")
- [Matroska](/wiki/Matroska "Matroska")
  - [WebM](/wiki/WebM "WebM")
- [QuickTime File Format](/wiki/QuickTime_File_Format "QuickTime File Format")
- [RatDVD](/wiki/RatDVD "RatDVD")
- [RealMedia](/wiki/RealMedia "RealMedia")
- [RIFF](/wiki/Resource_Interchange_File_Format "Resource Interchange File Format")
  - [WAV](/wiki/WAV "WAV")
- [MOD and TOD](/wiki/MOD_and_TOD "MOD and TOD")
- [VOB, IFO and BUP](/wiki/VOB "VOB")

Collaborations|

- [NETVC](/wiki/NETVC "NETVC")
- [MPEG LA](/wiki/MPEG_LA "MPEG LA")
- [Alliance for Open Media](/wiki/Alliance_for_Open_Media "Alliance for Open Media")

[Methods](/wiki/Data_compression "Data compression")|

- [Entropy](/wiki/Entropy_encoding "Entropy encoding")
  - [Arithmetic](/wiki/Arithmetic_coding "Arithmetic coding")
  - [Huffman](/wiki/Huffman_coding "Huffman coding")
  - [Modified](/wiki/Modified_Huffman_coding "Modified Huffman coding")
- [LPC](/wiki/Linear_predictive_coding "Linear predictive coding")
  - [ACELP](/wiki/Algebraic_code-excited_linear_prediction "Algebraic code-excited linear prediction")
  - [CELP](/wiki/Code-excited_linear_prediction "Code-excited linear prediction")
  - [LSP](/wiki/Line_spectral_pairs "Line spectral pairs")
  - [WLPC](/wiki/Warped_linear_predictive_coding "Warped linear predictive coding")
- [Lossless](/wiki/Lossless_compression "Lossless compression")
- [Lossy](/wiki/Lossy_compression "Lossy compression")
- [LZ](/wiki/LZ77_and_LZ78 "LZ77 and LZ78")
  - [DEFLATE](/wiki/DEFLATE "DEFLATE")
  - [LZW](/wiki/Lempel%E2%80%93Ziv%E2%80%93Welch "Lempel–Ziv–Welch")
- [PCM](/wiki/Pulse-code_modulation "Pulse-code modulation")
  - [A-law](/wiki/A-law_algorithm "A-law algorithm")
  - [µ-law](/wiki/%CE%9C-law_algorithm "Μ-law algorithm")
  - [ADPCM](/wiki/Adaptive_differential_pulse-code_modulation "Adaptive differential pulse-code modulation")
  - [DPCM](/wiki/Differential_pulse-code_modulation "Differential pulse-code modulation")
- [Transforms](/wiki/Transform_coding "Transform coding")
  - [DCT](/wiki/Discrete_cosine_transform "Discrete cosine transform")
  - [FFT](/wiki/Fast_Fourier_transform "Fast Fourier transform")
  - [MDCT](/wiki/Modified_discrete_cosine_transform "Modified discrete cosine transform")
  - [Wavelet](/wiki/Wavelet "Wavelet")
    - [Daubechies](/wiki/Daubechies_wavelet "Daubechies wavelet")
    - [DWT](/wiki/Discrete_wavelet_transform "Discrete wavelet transform")

Lists|

- [Comparison of audio coding formats](/wiki/Comparison_of_audio_coding_formats "Comparison of audio coding formats")
- [Comparison of video codecs](/wiki/Comparison_of_video_codecs "Comparison of video codecs")
- [List of codecs](/wiki/List_of_codecs "List of codecs")

See [Compression methods](/wiki/Template:Compression_methods "Template:Compression methods") for techniques and [Compression software](/wiki/Template:Compression_software "Template:Compression software") for codecs  
  
- [v](/wiki/Template:MPEG "Template:MPEG")
- [t](/wiki/Template_talk:MPEG "Template talk:MPEG")
- [e](/wiki/Special:EditPage/Template:MPEG "Special:EditPage/Template:MPEG")

[MPEG (Moving Picture Experts Group)](/wiki/Moving_Picture_Experts_Group "Moving Picture Experts Group")  
---  
  
- [MPEG-1](/wiki/MPEG-1 "MPEG-1")
- [2](/wiki/MPEG-2 "MPEG-2")
- [3](/wiki/MPEG-3 "MPEG-3")
- [4](/wiki/MPEG-4 "MPEG-4")
- [7](/wiki/MPEG-7 "MPEG-7")
- [21](/wiki/MPEG-21 "MPEG-21")
- [A](/wiki/MPEG-A "MPEG-A")
- B
- C
- [D](/wiki/MPEG-D "MPEG-D")
- E
- [G](/wiki/MPEG-G "MPEG-G")
- V
- M
- U
- [H](/wiki/MPEG-H "MPEG-H")
- [I](/w/index.php?title=MPEG-I&action=edit&redlink=1 "MPEG-I \(page does not exist\)")
- [5](/wiki/MPEG-5 "MPEG-5")

MPEG-1 Parts|

- Part 1: Systems
  - [Program stream](/wiki/MPEG_program_stream "MPEG program stream")
- Part 2: Video
  - based on [H.261](/wiki/H.261 "H.261")
- Part 3: Audio
  - [Layer I](/wiki/MPEG-1_Audio_Layer_I "MPEG-1 Audio Layer I")
  - [Layer II](/wiki/MPEG-1_Audio_Layer_II "MPEG-1 Audio Layer II")
  - [Layer III](/wiki/MP3 "MP3")

MPEG-2 Parts|

- Part 1: Systems (H.222.0)
  - Transport stream
  - [Program stream](/wiki/MPEG_program_stream "MPEG program stream")
- [Part 2: Video (H.262)](/wiki/H.262/MPEG-2_Part_2 "H.262/MPEG-2 Part 2")
- [Part 3: Audio](/wiki/MPEG-2_Part_3 "MPEG-2 Part 3")
  - [Layer I](/wiki/MPEG-1_Audio_Layer_I "MPEG-1 Audio Layer I")
  - [Layer II](/wiki/MPEG-2_Audio_Layer_II "MPEG-2 Audio Layer II")
  - [Layer III](/wiki/MP3 "MP3")
  - [MPEG Multichannel](/wiki/MPEG_Multichannel "MPEG Multichannel")
- [Part 6: DSM CC](/wiki/DSM_CC "DSM CC")
- [Part 7: Advanced Audio Coding](/wiki/Advanced_Audio_Coding "Advanced Audio Coding")

MPEG-4 Parts|

- [Part 2: Video](/wiki/MPEG-4_Part_2 "MPEG-4 Part 2")
  - based on [H.263](/wiki/H.263 "H.263")
- [Part 3: Audio](/wiki/MPEG-4_Part_3 "MPEG-4 Part 3")
- [Part 6: DMIF](/wiki/Delivery_Multimedia_Integration_Framework "Delivery Multimedia Integration Framework")
- [Part 10: Advanced Video Coding (H.264)](/wiki/Advanced_Video_Coding "Advanced Video Coding")
- [Part 11: Scene description](/wiki/MPEG-4_Part_11 "MPEG-4 Part 11")
- [Part 12: ISO base media file format](/wiki/ISO_base_media_file_format "ISO base media file format")
- [Part 14: MP4 file format](/wiki/MPEG-4_Part_14 "MPEG-4 Part 14")
- [Part 17: Streaming text format](/wiki/MPEG-4_Part_17 "MPEG-4 Part 17")
- [Part 20: LASeR](/wiki/MPEG-4_Part_20 "MPEG-4 Part 20")
- [Part 22: Open Font Format](/wiki/Open_Font_Format "Open Font Format")
- [Part 33: Internet Video Coding](/wiki/Internet_Video_Coding "Internet Video Coding")

MPEG-7 Parts|

- [Part 2: Description definition language](/wiki/Description_Definition_Language "Description Definition Language")

MPEG-21 Parts|

- [Parts 2, 3 and 9: Digital Item](/wiki/Digital_Item "Digital Item")
- [Part 5: Rights Expression Language](/wiki/Rights_Expression_Language "Rights Expression Language")

MPEG-D Parts|

- [Part 1: MPEG Surround](/wiki/MPEG_Surround "MPEG Surround")
- [Part 3: Unified Speech and Audio Coding](/wiki/Unified_Speech_and_Audio_Coding "Unified Speech and Audio Coding")

MPEG-G Parts|

- [Part 1: Transport and Storage of Genomic Information](/wiki/MPEG-G "MPEG-G")
- [Part 2: Coding of Genomic Information](/wiki/MPEG-G "MPEG-G")
- [Part 3: APIs](/wiki/MPEG-G "MPEG-G")
- [Part 4: Reference Software](/wiki/MPEG-G "MPEG-G")
- [Part 5: Conformance](/wiki/MPEG-G "MPEG-G")

MPEG-H Parts|

- [Part 1: MPEG media transport](/wiki/MPEG_media_transport "MPEG media transport")
- [Part 2: High Efficiency Video Coding (H.265)](/wiki/High_Efficiency_Video_Coding "High Efficiency Video Coding")
- [Part 3: MPEG-H 3D Audio](/wiki/MPEG-H_3D_Audio "MPEG-H 3D Audio")
- [Part 12: High Efficiency Image File Format](/wiki/High_Efficiency_Image_File_Format "High Efficiency Image File Format")

MPEG-I Parts|

- [Part 3: Versatile Video Coding (H.266)](/wiki/Versatile_Video_Coding "Versatile Video Coding")

MPEG-5 Parts|

- [Part 1: Essential Video Coding](/wiki/Essential_Video_Coding "Essential Video Coding")
- [Part 2: Low Complexity Enhancement Video Coding](/wiki/LCEVC "LCEVC")

Other| [MPEG-DASH](/wiki/Dynamic_Adaptive_Streaming_over_HTTP "Dynamic Adaptive Streaming over HTTP")  
  
Retrieved from "[https://en.wikipedia.org/w/index.php?title=MPEG_transport_stream&oldid=1324443299](https://en.wikipedia.org/w/index.php?title=MPEG_transport_stream&oldid=1324443299)"

[Categories](/wiki/Help:Category "Help:Category"):

- [ATSC](/wiki/Category:ATSC "Category:ATSC")
- [Digital container formats](/wiki/Category:Digital_container_formats "Category:Digital container formats")
- [MPEG-2](/wiki/Category:MPEG-2 "Category:MPEG-2")
- [ITU-T recommendations](/wiki/Category:ITU-T_recommendations "Category:ITU-T recommendations")

Hidden categories:

- [All articles with self-published sources](/wiki/Category:All_articles_with_self-published_sources "Category:All articles with self-published sources")
- [Articles with self-published sources from May 2012](/wiki/Category:Articles_with_self-published_sources_from_May_2012 "Category:Articles with self-published sources from May 2012")
- [CS1: unfit URL](/wiki/Category:CS1:_unfit_URL "Category:CS1: unfit URL")
- [Webarchive template wayback links](/wiki/Category:Webarchive_template_wayback_links "Category:Webarchive template wayback links")
- [Articles with short description](/wiki/Category:Articles_with_short_description "Category:Articles with short description")
- [Short description matches Wikidata](/wiki/Category:Short_description_matches_Wikidata "Category:Short description matches Wikidata")
- [Use dmy dates from August 2019](/wiki/Category:Use_dmy_dates_from_August_2019 "Category:Use dmy dates from August 2019")

  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
