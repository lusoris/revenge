package transcode

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/asticode/go-astiav"
)

// codecTagHVC1 is the FourCC codec tag for 'hvc1' HEVC sample entry.
// MKTAG('h','v','c','1') = 'h' | ('v'<<8) | ('c'<<16) | ('1'<<24) = 0x31637668.
// Forces VPS/SPS/PPS into the init segment (hvcC box) instead of inline in segments,
// which is required for HLS fMP4 and MSE/SourceBuffer compatibility.
const codecTagHVC1 astiav.CodecTag = 0x31637668

// streamMapping maps an input stream to its decoder, encoder, filter graph, and output stream.
type streamMapping struct {
	inputStream   *astiav.Stream
	outputStream  *astiav.Stream
	decCodecCtx   *astiav.CodecContext
	encCodecCtx   *astiav.CodecContext
	filterGraph   *astiav.FilterGraph
	buffersrcCtx  *astiav.BuffersrcFilterContext
	buffersinkCtx *astiav.BuffersinkFilterContext
	filterFrame   *astiav.Frame
	encPkt        *astiav.Packet
	bsfCtx        *astiav.BitStreamFilterContext // optional BSF for remux (e.g. DV NAL stripping)
	needsDecode   bool
	mediaType     astiav.MediaType
}

// TranscodeJob represents an in-process astiav transcode/remux job that replaces
// the subprocess-based FFmpegProcess. All FFmpeg work is done through libavcodec/
// libavformat/libavfilter C libraries via go-astiav bindings — no child processes.
type TranscodeJob struct {
	// Configuration
	InputFile  string
	OutputDir  string
	OutputFile string // playlist path (e.g. index.m3u8)

	// Job identity
	SessionID string
	Profile   string

	// Transcode settings
	VideoCodec   string // "copy" or "libx264"
	AudioCodec   string // "copy" or "aac" or "" (disabled)
	Width        int    // target width (0 = keep)
	Height       int    // target height (0 = keep)
	VideoBitrate int    // kbps (0 = no limit)
	AudioBitrate int    // kbps (0 = default)
	CRF          int    // constant rate factor (0 = default 23)
	Preset       string // encoding preset (empty = "veryfast")

	// HLS settings
	SegmentDuration int // seconds per segment
	SegmentPattern  string

	// Stream selection
	VideoStreamIndex int // -1 to disable video
	AudioStreamIndex int // -1 to disable audio
	SeekSeconds      int

	// DV handling
	StripDolbyVision bool // strip DV RPU NALs + patch hvcC for non-DV clients

	// Lifecycle
	Done        chan struct{}
	Err         error
	IsTranscode bool // true if encoding (not copy)

	// Cancellation
	cancel     context.CancelFunc
	interrupter *astiav.IOInterrupter

	mu      sync.Mutex
	stopped bool
}

// TranscodeJobConfig holds the configuration for creating a TranscodeJob.
type TranscodeJobConfig struct {
	InputFile         string
	OutputDir         string
	SessionID         string
	Profile           string
	VideoCodec        string
	AudioCodec        string
	Width             int
	Height            int
	VideoBitrate      int
	AudioBitrate      int
	CRF               int
	Preset            string
	SegmentDuration   int
	VideoStreamIndex  int // -1 to disable
	AudioStreamIndex  int // -1 to disable
	SeekSeconds       int
	StripDolbyVision  bool // strip DV RPU NALs + patch hvcC for non-DV clients
}

// NewTranscodeJob creates a new transcode job from the given config.
func NewTranscodeJob(cfg TranscodeJobConfig) *TranscodeJob {
	preset := cfg.Preset
	if preset == "" {
		preset = "veryfast"
	}
	crf := cfg.CRF
	if crf == 0 {
		crf = 23
	}
	segDur := cfg.SegmentDuration
	if segDur == 0 {
		segDur = 6
	}

	outputFile := filepath.Join(cfg.OutputDir, "index.m3u8")
	segPattern := filepath.Join(cfg.OutputDir, "seg-%05d.m4s")

	isTranscode := (cfg.VideoCodec != "copy" && cfg.VideoCodec != "") || (cfg.AudioCodec != "copy" && cfg.AudioCodec != "")

	return &TranscodeJob{
		InputFile:        cfg.InputFile,
		OutputDir:        cfg.OutputDir,
		OutputFile:       outputFile,
		SessionID:        cfg.SessionID,
		Profile:          cfg.Profile,
		VideoCodec:       cfg.VideoCodec,
		AudioCodec:       cfg.AudioCodec,
		Width:            cfg.Width,
		Height:           cfg.Height,
		VideoBitrate:     cfg.VideoBitrate,
		AudioBitrate:     cfg.AudioBitrate,
		CRF:              crf,
		Preset:           preset,
		SegmentDuration:  segDur,
		SegmentPattern:   segPattern,
		VideoStreamIndex: cfg.VideoStreamIndex,
		AudioStreamIndex: cfg.AudioStreamIndex,
		SeekSeconds:      cfg.SeekSeconds,
		StripDolbyVision: cfg.StripDolbyVision,
		Done:             make(chan struct{}),
		IsTranscode:      isTranscode,
	}
}

// Stop interrupts the running job. Safe to call multiple times.
func (j *TranscodeJob) Stop() {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.stopped {
		return
	}
	j.stopped = true
	if j.interrupter != nil {
		j.interrupter.Interrupt()
	}
	if j.cancel != nil {
		j.cancel()
	}
}

// Run executes the transcode/remux job synchronously. Call this in a goroutine.
// It uses astiav's in-process FFmpeg libraries — no subprocess is spawned.
func (j *TranscodeJob) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	j.cancel = cancel
	defer cancel()

	// Create interrupter for cancellation of blocking I/O operations
	j.interrupter = astiav.NewIOInterrupter()
	defer j.interrupter.Free()

	// Monitor context cancellation and trigger interrupter
	go func() {
		<-ctx.Done()
		j.interrupter.Interrupt()
	}()

	// --- Open input ---
	inputFmtCtx := astiav.AllocFormatContext()
	if inputFmtCtx == nil {
		return errors.New("failed to allocate input format context")
	}
	defer inputFmtCtx.Free()
	inputFmtCtx.SetIOInterrupter(j.interrupter)

	if err := inputFmtCtx.OpenInput(j.InputFile, nil, nil); err != nil {
		return fmt.Errorf("failed to open input %q: %w", j.InputFile, err)
	}
	defer inputFmtCtx.CloseInput()

	if err := inputFmtCtx.FindStreamInfo(nil); err != nil {
		return fmt.Errorf("failed to find stream info: %w", err)
	}

	// --- Seek if needed ---
	if j.SeekSeconds > 0 {
		// Seek to timestamp in AV_TIME_BASE units
		ts := int64(j.SeekSeconds) * int64(astiav.TimeBase)
		if err := inputFmtCtx.SeekFrame(-1, ts, astiav.NewSeekFlags(astiav.SeekFlagBackward)); err != nil {
			return fmt.Errorf("failed to seek to %ds: %w", j.SeekSeconds, err)
		}
	}

	// --- Find input streams ---
	streams := make(map[int]*streamMapping)

	// Identify which input streams to process
	for _, is := range inputFmtCtx.Streams() {
		cp := is.CodecParameters()
		mt := cp.MediaType()

		switch mt {
		case astiav.MediaTypeVideo:
			if j.VideoStreamIndex == -1 {
				continue // video disabled
			}
			if is.Index() != j.VideoStreamIndex && j.VideoStreamIndex >= 0 {
				continue
			}
			// Use first video stream if VideoStreamIndex == 0
			if _, exists := findStreamByType(streams, astiav.MediaTypeVideo); exists {
				continue // already have one
			}
			streams[is.Index()] = &streamMapping{
				inputStream: is,
				needsDecode: j.VideoCodec != "copy",
				mediaType:   mt,
			}

		case astiav.MediaTypeAudio:
			if j.AudioStreamIndex == -1 {
				continue // audio disabled
			}
			// AudioStreamIndex is the relative audio index, find matching absolute stream
			audioIdx := countStreamsBefore(inputFmtCtx, is.Index(), astiav.MediaTypeAudio)
			if audioIdx != j.AudioStreamIndex {
				continue
			}
			streams[is.Index()] = &streamMapping{
				inputStream: is,
				needsDecode: j.AudioCodec != "copy",
				mediaType:   mt,
			}
		}
	}

	if len(streams) == 0 {
		return errors.New("no matching input streams found")
	}

	// --- Set up decoders for streams that need transcoding ---
	var cleanups []func()
	defer func() {
		for i := len(cleanups) - 1; i >= 0; i-- {
			cleanups[i]()
		}
	}()

	for _, sm := range streams {
		if !sm.needsDecode {
			continue
		}
		codec := astiav.FindDecoder(sm.inputStream.CodecParameters().CodecID())
		if codec == nil {
			return fmt.Errorf("decoder not found for codec %s", sm.inputStream.CodecParameters().CodecID().Name())
		}
		sm.decCodecCtx = astiav.AllocCodecContext(codec)
		if sm.decCodecCtx == nil {
			return errors.New("failed to allocate decoder codec context")
		}
		cleanups = append(cleanups, sm.decCodecCtx.Free)

		if err := sm.inputStream.CodecParameters().ToCodecContext(sm.decCodecCtx); err != nil {
			return fmt.Errorf("failed to copy codec params to decoder: %w", err)
		}
		if err := sm.decCodecCtx.Open(codec, nil); err != nil {
			return fmt.Errorf("failed to open decoder: %w", err)
		}
	}

	// --- Open output (HLS muxer) ---
	outputFmtCtx, err := astiav.AllocOutputFormatContext(nil, "hls", j.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to allocate output format context: %w", err)
	}
	if outputFmtCtx == nil {
		return errors.New("output format context is nil")
	}
	defer outputFmtCtx.Free()

	// Set HLS muxer options via private data
	if pd := outputFmtCtx.PrivateData(); pd != nil {
		opts := pd.Options()
		searchFlags := astiav.NewOptionSearchFlags()
		_ = opts.Set("hls_time", strconv.Itoa(j.SegmentDuration), searchFlags)
		_ = opts.Set("hls_playlist_type", "event", searchFlags)
		_ = opts.Set("hls_segment_filename", j.SegmentPattern, searchFlags)
		_ = opts.Set("start_number", "0", searchFlags)
		// Use fMP4 segments instead of MPEG-TS for HEVC/AV1/modern codec support.
		// fMP4 (fragmented MP4) is required by HLS spec for H.265, AV1, and
		// provides better seeking, codec flexibility, and lower overhead.
		_ = opts.Set("hls_segment_type", "fmp4", searchFlags)
		_ = opts.Set("hls_fmp4_init_filename", "init.mp4", searchFlags)
	}

	// --- Create output streams and encoders ---
	for _, is := range inputFmtCtx.Streams() {
		sm, ok := streams[is.Index()]
		if !ok {
			continue
		}

		if sm.needsDecode {
			// Transcoding path: create encoder
			if err := j.setupEncoder(sm, outputFmtCtx, &cleanups); err != nil {
				return fmt.Errorf("failed to setup encoder for stream %d: %w", is.Index(), err)
			}
		} else {
			// Remux/copy path: copy codec parameters
			sm.outputStream = outputFmtCtx.NewStream(nil)
			if sm.outputStream == nil {
				return errors.New("failed to create output stream")
			}
			if err := is.CodecParameters().Copy(sm.outputStream.CodecParameters()); err != nil {
				return fmt.Errorf("failed to copy codec parameters: %w", err)
			}

			outCP := sm.outputStream.CodecParameters()
			// For HEVC video in fMP4/HLS, force the 'hvc1' sample entry tag.
			// FFmpeg defaults to 'hev1' which stores parameter sets (VPS/SPS/PPS)
			// inline in segments. 'hvc1' stores them only in the init segment,
			// which is required by Apple HLS spec and for MSE/SourceBuffer compat.
			// Without this, browsers reject fMP4 data with bufferAppendingError
			// because the CODECS string says 'hvc1' but the actual box is 'hev1'.
			if outCP.CodecID() == astiav.CodecIDHevc && sm.mediaType == astiav.MediaTypeVideo {
				// Force the 'hvc1' sample entry tag for HEVC in fMP4/HLS.
				// FFmpeg defaults to 'hev1', but 'hvc1' stores parameter sets
				// only in init.mp4 which is required for MSE/SourceBuffer compat.
				outCP.SetCodecTag(codecTagHVC1)

				if j.StripDolbyVision {
					// Patch the hvcC extradata BEFORE WriteHeader to clean DV artifacts.
					// DV Profile 8 uses constraint byte 0x90 which Chrome rejects.
					// We set the general_non_packed_constraint_flag (bit 5) to convert
					// 0x90 → 0xB0 which is standard HEVC Main 10 with PQ.
					patchHvcCConstraints(outCP)

					// Strip DOVI_CONF side data to prevent 'dby1' ftyp brand.
					// FFmpeg's fMP4 muxer checks for DOVI_CONF on the codec params
					// and adds 'dby1' as a compatible brand when present — Chrome
					// sees this and refuses the stream.
					stripDOVIConfSideData(outCP)

					// Strip Dolby Vision RPU NALs from HEVC segments.
					// Uses dovi_rpu BSF with strip=1 which removes DV RPU data blocks
					// (NAL types 62/63) — required for Chrome/Firefox MSE compat.
					if bsfCtx, err := setupDVRemovalFilter(is); err != nil {
						slog.Warn("failed to setup DV removal BSF, continuing without", "error", err)
					} else if bsfCtx != nil {
						sm.bsfCtx = bsfCtx
						cleanups = append(cleanups, bsfCtx.Free)
						slog.Info("DV stripping enabled for HEVC remux",
							"session", j.SessionID, "profile", j.Profile, "bsf", "dovi_rpu")
					}
				}
			} else {
				outCP.SetCodecTag(0)
			}

			sm.outputStream.SetTimeBase(is.TimeBase())
		}
	}

	// --- Set up filter graphs for transcoded streams ---
	for _, sm := range streams {
		if !sm.needsDecode {
			continue
		}
		if err := j.setupFilters(sm, &cleanups); err != nil {
			return fmt.Errorf("failed to setup filters: %w", err)
		}
	}

	// --- Open output IO ---
	if !outputFmtCtx.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		ioCtx, err := astiav.OpenIOContext(j.OutputFile, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), j.interrupter, nil)
		if err != nil {
			return fmt.Errorf("failed to open output IO context: %w", err)
		}
		cleanups = append(cleanups, func() { _ = ioCtx.Close() })
		outputFmtCtx.SetPb(ioCtx)
	}

	// --- Write header ---
	if err := outputFmtCtx.WriteHeader(nil); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// --- Main read/process/write loop ---
	decFrame := astiav.AllocFrame()
	if decFrame == nil {
		return errors.New("failed to allocate decode frame")
	}
	defer decFrame.Free()

	pkt := astiav.AllocPacket()
	if pkt == nil {
		return errors.New("failed to allocate packet")
	}
	defer pkt.Free()

	for {
		// Check cancellation
		if ctx.Err() != nil {
			break
		}

		if err := inputFmtCtx.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			// Interrupted I/O — check if we were cancelled
			if ctx.Err() != nil {
				break
			}
			return fmt.Errorf("failed to read frame: %w", err)
		}

		sm, ok := streams[pkt.StreamIndex()]
		if !ok {
			pkt.Unref()
			continue
		}

		if sm.needsDecode {
			// Decode → filter → encode → write
			if err := j.decodeFilterEncode(sm, decFrame, pkt, outputFmtCtx); err != nil {
				pkt.Unref()
				if ctx.Err() != nil {
					break
				}
				return fmt.Errorf("decode/filter/encode failed: %w", err)
			}
		} else {
			// Remux: optionally apply BSF (e.g. DV NAL stripping), then write
			pkt.SetStreamIndex(sm.outputStream.Index())
			pkt.RescaleTs(sm.inputStream.TimeBase(), sm.outputStream.TimeBase())
			pkt.SetPos(-1)

			if sm.bsfCtx != nil {
				// Run packet through bitstream filter
				if err := sm.bsfCtx.SendPacket(pkt); err != nil {
					pkt.Unref()
					if ctx.Err() != nil {
						break
					}
					return fmt.Errorf("bsf send packet failed: %w", err)
				}
				for {
					if err := sm.bsfCtx.ReceivePacket(pkt); err != nil {
						break // EAGAIN or EOF — done with this packet
					}
					if err := outputFmtCtx.WriteInterleavedFrame(pkt); err != nil {
						if ctx.Err() != nil {
							break
						}
						return fmt.Errorf("failed to write interleaved frame: %w", err)
					}
				}
			} else {
				if err := outputFmtCtx.WriteInterleavedFrame(pkt); err != nil {
					if ctx.Err() != nil {
						break
					}
					return fmt.Errorf("failed to write interleaved frame: %w", err)
				}
			}
		}

		pkt.Unref()
	}

	// --- Flush encoders ---
	for _, sm := range streams {
		if !sm.needsDecode || sm.encCodecCtx == nil {
			continue
		}
		if err := j.flushEncoder(sm, outputFmtCtx); err != nil {
			// Ignore flush errors on cancellation
			if ctx.Err() != nil {
				break
			}
			return fmt.Errorf("failed to flush encoder: %w", err)
		}
	}

	// --- Write trailer ---
	if err := outputFmtCtx.WriteTrailer(); err != nil {
		// Ignore trailer errors on cancellation
		if ctx.Err() == nil {
			return fmt.Errorf("failed to write trailer: %w", err)
		}
	}

	return ctx.Err()
}

// setupEncoder creates an output stream with the appropriate encoder for a transcoded stream.
func (j *TranscodeJob) setupEncoder(sm *streamMapping, outputFmtCtx *astiav.FormatContext, cleanups *[]func()) error {
	var codecID astiav.CodecID
	switch sm.mediaType {
	case astiav.MediaTypeVideo:
		codecID = resolveVideoCodecID(j.VideoCodec)
	case astiav.MediaTypeAudio:
		codecID = resolveAudioCodecID(j.AudioCodec)
	default:
		return fmt.Errorf("unsupported media type for encoding: %s", sm.mediaType)
	}

	encCodec := astiav.FindEncoder(codecID)
	if encCodec == nil {
		return fmt.Errorf("encoder not found for codec ID %s", codecID.Name())
	}

	sm.encCodecCtx = astiav.AllocCodecContext(encCodec)
	if sm.encCodecCtx == nil {
		return errors.New("failed to allocate encoder codec context")
	}
	*cleanups = append(*cleanups, sm.encCodecCtx.Free)

	switch sm.mediaType {
	case astiav.MediaTypeVideo:
		// Determine output dimensions
		height := sm.decCodecCtx.Height()
		width := sm.decCodecCtx.Width()
		if j.Height > 0 {
			height = j.Height
			// Maintain aspect ratio: width = -2 equivalent (even number)
			if sm.decCodecCtx.Height() > 0 {
				width = sm.decCodecCtx.Width() * j.Height / sm.decCodecCtx.Height()
				// Round to nearest even number
				if width%2 != 0 {
					width++
				}
			}
		}

		sm.encCodecCtx.SetWidth(width)
		sm.encCodecCtx.SetHeight(height)

		// Use encoder's supported pixel format or fall back to source
		if fmts := encCodec.SupportedPixelFormats(); len(fmts) > 0 {
			sm.encCodecCtx.SetPixelFormat(fmts[0])
		} else {
			sm.encCodecCtx.SetPixelFormat(sm.decCodecCtx.PixelFormat())
		}

		sm.encCodecCtx.SetSampleAspectRatio(sm.decCodecCtx.SampleAspectRatio())

		// Use the stream's average frame rate to derive a correct time_base.
		// The decoder's TimeBase is often wrong for H.264 (e.g. 1/60 for 30fps
		// due to ticks_per_frame=2), causing "Invalid argument" from the encoder.
		fr := sm.inputStream.AvgFrameRate()
		if fr.Num() > 0 && fr.Den() > 0 {
			sm.encCodecCtx.SetFramerate(fr)
			sm.encCodecCtx.SetTimeBase(astiav.NewRational(fr.Den(), fr.Num()))
		} else {
			sm.encCodecCtx.SetTimeBase(sm.decCodecCtx.TimeBase())
		}

		// GOP size: one keyframe per segment for clean HLS splits
		if fr.Num() > 0 && fr.Den() > 0 {
			fps := fr.Num() / fr.Den()
			if fps <= 0 {
				fps = 30
			}
			sm.encCodecCtx.SetGopSize(fps * j.SegmentDuration)
		} else {
			sm.encCodecCtx.SetGopSize(30 * j.SegmentDuration) // fallback
		}

		// Set encoding options via private data
		if pd := sm.encCodecCtx.PrivateData(); pd != nil {
			opts := pd.Options()
			searchFlags := astiav.NewOptionSearchFlags()
			_ = opts.Set("preset", j.Preset, searchFlags)
			_ = opts.Set("crf", strconv.Itoa(j.CRF), searchFlags)
		}

		// Set bitrate limits if specified
		if j.VideoBitrate > 0 {
			sm.encCodecCtx.SetBitRate(int64(j.VideoBitrate) * 1000)
			sm.encCodecCtx.SetRateControlMaxRate(int64(j.VideoBitrate) * 1000)
			sm.encCodecCtx.SetRateControlBufferSize(j.VideoBitrate * 2000)
		}

	case astiav.MediaTypeAudio:
		// Channel layout
		if layouts := encCodec.SupportedChannelLayouts(); len(layouts) > 0 {
			sm.encCodecCtx.SetChannelLayout(layouts[0])
		} else {
			sm.encCodecCtx.SetChannelLayout(sm.decCodecCtx.ChannelLayout())
		}

		sm.encCodecCtx.SetSampleRate(sm.decCodecCtx.SampleRate())

		if fmts := encCodec.SupportedSampleFormats(); len(fmts) > 0 {
			sm.encCodecCtx.SetSampleFormat(fmts[0])
		} else {
			sm.encCodecCtx.SetSampleFormat(sm.decCodecCtx.SampleFormat())
		}

		sm.encCodecCtx.SetTimeBase(astiav.NewRational(1, sm.encCodecCtx.SampleRate()))

		if j.AudioBitrate > 0 {
			sm.encCodecCtx.SetBitRate(int64(j.AudioBitrate) * 1000)
		}
	}

	// Global header flag for muxers that need it
	if outputFmtCtx.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
		sm.encCodecCtx.SetFlags(sm.encCodecCtx.Flags().Add(astiav.CodecContextFlagGlobalHeader))
	}

	if err := sm.encCodecCtx.Open(encCodec, nil); err != nil {
		return fmt.Errorf("failed to open encoder: %w", err)
	}

	// Create output stream
	sm.outputStream = outputFmtCtx.NewStream(nil)
	if sm.outputStream == nil {
		return errors.New("failed to create output stream")
	}

	if err := sm.outputStream.CodecParameters().FromCodecContext(sm.encCodecCtx); err != nil {
		return fmt.Errorf("failed to copy encoder params to output stream: %w", err)
	}
	sm.outputStream.SetTimeBase(sm.encCodecCtx.TimeBase())

	return nil
}

// setupFilters creates the filter graph for a transcoded stream.
// Video: handles pixel format conversion and optional scaling.
// Audio: handles sample format and channel layout conversion.
func (j *TranscodeJob) setupFilters(sm *streamMapping, cleanups *[]func()) error {
	sm.filterGraph = astiav.AllocFilterGraph()
	if sm.filterGraph == nil {
		return errors.New("failed to allocate filter graph")
	}
	*cleanups = append(*cleanups, sm.filterGraph.Free)

	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		return errors.New("failed to allocate filter outputs")
	}
	defer outputs.Free()

	inputs := astiav.AllocFilterInOut()
	if inputs == nil {
		return errors.New("failed to allocate filter inputs")
	}
	defer inputs.Free()

	// Buffersrc parameters
	bscp := astiav.AllocBuffersrcFilterContextParameters()
	defer bscp.Free()

	var buffersrc, buffersink *astiav.Filter
	var filterDesc string

	switch sm.mediaType {
	case astiav.MediaTypeVideo:
		buffersrc = astiav.FindFilterByName("buffer")
		buffersink = astiav.FindFilterByName("buffersink")

		bscp.SetWidth(sm.decCodecCtx.Width())
		bscp.SetHeight(sm.decCodecCtx.Height())
		bscp.SetPixelFormat(sm.decCodecCtx.PixelFormat())
		bscp.SetSampleAspectRatio(sm.decCodecCtx.SampleAspectRatio())
		bscp.SetTimeBase(sm.inputStream.TimeBase())

		// Build filter description
		var filters []string

		// Scale filter if target dimensions differ
		if j.Height > 0 && j.Height != sm.decCodecCtx.Height() {
			filters = append(filters, fmt.Sprintf("scale=-2:%d", j.Height))
		}

		// Pixel format conversion to match encoder
		filters = append(filters, fmt.Sprintf("format=pix_fmts=%s", sm.encCodecCtx.PixelFormat().Name()))

		filterDesc = strings.Join(filters, ",")

	case astiav.MediaTypeAudio:
		buffersrc = astiav.FindFilterByName("abuffer")
		buffersink = astiav.FindFilterByName("abuffersink")

		bscp.SetChannelLayout(sm.decCodecCtx.ChannelLayout())
		bscp.SetSampleFormat(sm.decCodecCtx.SampleFormat())
		bscp.SetSampleRate(sm.decCodecCtx.SampleRate())
		bscp.SetTimeBase(sm.decCodecCtx.TimeBase())

		// Audio filter chain for transcoding:
		// 1. aresample: resample + fix discontinuous timestamps (async=1)
		// 2. aformat: convert to encoder's sample format and channel layout
		// 3. asetnsamples: split into exactly 1024-sample frames for AAC encoder
		// 4. asetpts: regenerate clean monotonic PTS from sample count
		//    (fMP4 muxer is strict about monotonically increasing DTS)
		filterDesc = fmt.Sprintf(
			"aresample=async=1,aformat=sample_fmts=%s:channel_layouts=%s,asetnsamples=n=1024,asetpts=N/SR/TB",
			sm.encCodecCtx.SampleFormat().Name(),
			sm.encCodecCtx.ChannelLayout().String(),
		)
	}

	if buffersrc == nil || buffersink == nil {
		return errors.New("buffersrc or buffersink filter not found")
	}

	// Create buffersrc context
	var err error
	sm.buffersrcCtx, err = sm.filterGraph.NewBuffersrcFilterContext(buffersrc, "in")
	if err != nil {
		return fmt.Errorf("failed to create buffersrc context: %w", err)
	}
	if err = sm.buffersrcCtx.SetParameters(bscp); err != nil {
		return fmt.Errorf("failed to set buffersrc parameters: %w", err)
	}
	if err = sm.buffersrcCtx.Initialize(nil); err != nil {
		return fmt.Errorf("failed to initialize buffersrc: %w", err)
	}

	// Create buffersink context
	sm.buffersinkCtx, err = sm.filterGraph.NewBuffersinkFilterContext(buffersink, "out")
	if err != nil {
		return fmt.Errorf("failed to create buffersink context: %w", err)
	}

	// Wire up filter in/out pads
	outputs.SetName("in")
	outputs.SetFilterContext(sm.buffersrcCtx.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	inputs.SetName("out")
	inputs.SetFilterContext(sm.buffersinkCtx.FilterContext())
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// Parse and configure filter graph
	if err = sm.filterGraph.Parse(filterDesc, inputs, outputs); err != nil {
		return fmt.Errorf("failed to parse filter graph %q: %w", filterDesc, err)
	}
	if err = sm.filterGraph.Configure(); err != nil {
		return fmt.Errorf("failed to configure filter graph: %w", err)
	}

	// Allocate filter frame and encode packet
	sm.filterFrame = astiav.AllocFrame()
	if sm.filterFrame == nil {
		return errors.New("failed to allocate filter frame")
	}
	*cleanups = append(*cleanups, sm.filterFrame.Free)

	sm.encPkt = astiav.AllocPacket()
	if sm.encPkt == nil {
		return errors.New("failed to allocate encode packet")
	}
	*cleanups = append(*cleanups, sm.encPkt.Free)

	return nil
}

// decodeFilterEncode decodes a packet, passes frames through the filter graph,
// encodes the filtered frames, and writes them to the output.
func (j *TranscodeJob) decodeFilterEncode(sm *streamMapping, decFrame *astiav.Frame, pkt *astiav.Packet, outputFmtCtx *astiav.FormatContext) error {
	// Send packet to decoder
	if err := sm.decCodecCtx.SendPacket(pkt); err != nil {
		return fmt.Errorf("send packet to decoder failed: %w", err)
	}

	// Receive all decoded frames
	for {
		if err := sm.decCodecCtx.ReceiveFrame(decFrame); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				return nil
			}
			return fmt.Errorf("receive frame from decoder failed: %w", err)
		}

		// Push frame into filter graph
		if err := sm.buffersrcCtx.AddFrame(decFrame, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
			decFrame.Unref()
			return fmt.Errorf("add frame to filter failed: %w", err)
		}
		decFrame.Unref()

		// Pull all filtered frames
		for {
			if err := sm.buffersinkCtx.GetFrame(sm.filterFrame, astiav.NewBuffersinkFlags()); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					break
				}
				return fmt.Errorf("get frame from filter failed: %w", err)
			}

			sm.filterFrame.SetPictureType(astiav.PictureTypeNone)

			// Encode filtered frame
			if err := j.encodeWriteFrame(sm, sm.filterFrame, outputFmtCtx); err != nil {
				sm.filterFrame.Unref()
				return err
			}
			sm.filterFrame.Unref()
		}
	}
}

// encodeWriteFrame encodes a frame and writes the resulting packets to the output.
func (j *TranscodeJob) encodeWriteFrame(sm *streamMapping, frame *astiav.Frame, outputFmtCtx *astiav.FormatContext) error {
	if err := sm.encCodecCtx.SendFrame(frame); err != nil {
		return fmt.Errorf("send frame to encoder failed: %w", err)
	}

	for {
		if err := sm.encCodecCtx.ReceivePacket(sm.encPkt); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				return nil
			}
			return fmt.Errorf("receive packet from encoder failed: %w", err)
		}

		sm.encPkt.SetStreamIndex(sm.outputStream.Index())
		sm.encPkt.RescaleTs(sm.encCodecCtx.TimeBase(), sm.outputStream.TimeBase())

		if err := outputFmtCtx.WriteInterleavedFrame(sm.encPkt); err != nil {
			sm.encPkt.Unref()
			return fmt.Errorf("write interleaved frame failed: %w", err)
		}
		sm.encPkt.Unref()
	}
}

// flushEncoder flushes remaining frames from the encoder.
func (j *TranscodeJob) flushEncoder(sm *streamMapping, outputFmtCtx *astiav.FormatContext) error {
	// Send nil frame to signal end of stream
	if err := sm.encCodecCtx.SendFrame(nil); err != nil {
		return fmt.Errorf("flush encoder send failed: %w", err)
	}

	for {
		if err := sm.encCodecCtx.ReceivePacket(sm.encPkt); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				return nil
			}
			return fmt.Errorf("flush encoder receive failed: %w", err)
		}

		sm.encPkt.SetStreamIndex(sm.outputStream.Index())
		sm.encPkt.RescaleTs(sm.encCodecCtx.TimeBase(), sm.outputStream.TimeBase())

		if err := outputFmtCtx.WriteInterleavedFrame(sm.encPkt); err != nil {
			sm.encPkt.Unref()
			return fmt.Errorf("flush write interleaved frame failed: %w", err)
		}
		sm.encPkt.Unref()
	}
}

// --- Helpers ---

// resolveVideoCodecID maps a codec name to astiav CodecID.
func resolveVideoCodecID(name string) astiav.CodecID {
	switch strings.ToLower(name) {
	case "libx264", "h264":
		return astiav.CodecIDH264
	case "libx265", "hevc", "h265":
		return astiav.CodecIDHevc
	case "libvpx-vp9", "vp9":
		return astiav.CodecIDVp9
	case "libaom-av1", "av1":
		return astiav.CodecIDAv1
	default:
		return astiav.CodecIDH264
	}
}

// resolveAudioCodecID maps a codec name to astiav CodecID.
func resolveAudioCodecID(name string) astiav.CodecID {
	switch strings.ToLower(name) {
	case "aac":
		return astiav.CodecIDAac
	case "mp3":
		return astiav.CodecIDMp3
	case "ac3":
		return astiav.CodecIDAc3
	case "eac3":
		return astiav.CodecIDEac3
	case "opus":
		return astiav.CodecIDOpus
	default:
		return astiav.CodecIDAac
	}
}

// countStreamsBefore counts how many streams of the given type appear before the given index.
func countStreamsBefore(fmtCtx *astiav.FormatContext, beforeIndex int, mediaType astiav.MediaType) int {
	count := 0
	for _, s := range fmtCtx.Streams() {
		if s.Index() >= beforeIndex {
			break
		}
		if s.CodecParameters().MediaType() == mediaType {
			count++
		}
	}
	return count
}

// patchHvcCConstraints fixes the general_constraint_indicator_flags in the
// HEVCDecoderConfigurationRecord (hvcC extradata) for Dolby Vision sources.
//
// DV content has constraint byte 0x90 at offset 6 (general_non_packed_constraint_flag=0).
// Chrome's MSE refuses to create a SourceBuffer when the hvcC declares these
// DV-specific constraints. Standard HEVC uses 0xB0 (general_non_packed_constraint_flag=1).
//
// This function patches THREE locations:
// 1. The outer hvcC header at byte 6
// 2. The VPS NALU's profile_tier_level constraint byte
// 3. The SPS NALU's profile_tier_level constraint byte
//
// Chrome's MSE decoder parses VPS/SPS NALUs inside hvcC, so all must be patched.
//
// HEVCDecoderConfigurationRecord layout (ISO/IEC 14496-15):
//
//	Byte 0:    configurationVersion
//	Byte 1:    profile_space(2) | tier_flag(1) | profile_idc(5)
//	Bytes 2-5: general_profile_compatibility_flags
//	Bytes 6-11: general_constraint_indicator_flags  ← byte 6 is patched
//	Byte 12:   general_level_idc
//	...
//	Byte 22:   numOfArrays
//	followed by arrays of VPS, SPS, PPS NALUs
func patchHvcCConstraints(cp *astiav.CodecParameters) {
	extradata := cp.ExtraData()
	if len(extradata) < 23 {
		return // not a valid hvcC record
	}

	patched := make([]byte, len(extradata))
	copy(patched, extradata)
	patchCount := 0

	// 1. Patch outer hvcC header constraint byte (offset 6)
	if patched[6]&0x20 == 0 {
		before := patched[6]
		patched[6] |= 0x20
		slog.Debug("patched hvcC header constraint byte",
			"before", fmt.Sprintf("0x%02X", before),
			"after", fmt.Sprintf("0x%02X", patched[6]))
		patchCount++
	}

	// 2. Patch VPS and SPS NALUs inside hvcC
	// Parse NALU arrays starting at byte 22
	numArrays := int(patched[22])
	pos := 23
	for i := 0; i < numArrays && pos+3 <= len(patched); i++ {
		nalType := patched[pos] & 0x3F
		numNALUs := int(patched[pos+1])<<8 | int(patched[pos+2])
		pos += 3

		for j := 0; j < numNALUs && pos+2 <= len(patched); j++ {
			naluLen := int(patched[pos])<<8 | int(patched[pos+1])
			naluStart := pos + 2 // start of actual NALU data
			naluEnd := naluStart + naluLen
			if naluEnd > len(patched) {
				break
			}

			// VPS (type 32) and SPS (type 33) both have profile_tier_level
			// The constraint byte is at offset 9 within the NALU (after NAL header + other fields)
			// VPS: NAL header(2) + vps_id/max_layers(2) + profile_tier_level starts
			// SPS: NAL header(2) + sps_video_parameter_set_id(4bits) + ... profile_tier_level
			// Both have constraint_indicator_flags at similar relative position
			if nalType == 32 || nalType == 33 { // VPS or SPS
				// profile_tier_level in VPS/SPS: constraint bytes start at byte 5 after syntax start
				// For VPS: starts around byte 6 of NALU
				// For SPS: starts around byte 5-6 of NALU (after nal_header + sps_video_parameter_set_id)
				// The pattern 02 20 00 00 03 00 XX where XX is constraint byte
				// We look for the pattern and patch byte after "00 03 00"
				for k := naluStart; k < naluEnd-6; k++ {
					// Look for profile_tier_level pattern:
					// profile_space(2) + tier(1) + profile_idc(5) = 0x02 for Main 10
					// then profile_compat_flags(32-bit) often 20 00 00 03 or similar
					// then constraint bytes starting with 00 90 (DV) or 00 B0 (standard)
					// The "03" is RBSP escape, followed by the original byte
					// Raw bytes: 00 03 00 90 means escaped 00 90
					if patched[k] == 0x00 && patched[k+1] == 0x03 && patched[k+2] == 0x00 {
						// Check if byte at k+3 is 0x90 (DV constraint)
						if patched[k+3] == 0x90 {
							patched[k+3] = 0xB0
							slog.Debug("patched NALU constraint byte",
								"nal_type", nalType, "offset", k+3,
								"before", "0x90", "after", "0xB0")
							patchCount++
						}
					}
				}
			}
			pos = naluEnd
		}
	}

	if patchCount > 0 {
		if err := cp.SetExtraData(patched); err != nil {
			slog.Warn("failed to patch hvcC constraint flags", "error", err)
		} else {
			slog.Info("patched hvcC constraint flags for DV removal",
				"patch_count", patchCount)
		}
	}
}

// setupDVRemovalFilter creates a bitstream filter to strip Dolby Vision metadata
// from HEVC streams. This uses dovi_rpu with strip=1 which removes both the DOVI
// configuration record and RPU data blocks (NAL unit types 62/63) from the stream.
// This is critical for Chrome/Firefox MSE compatibility — DV RPU NALs in segments
// cause bufferAppendingError because the browser can't decode them.
// Falls back to filter_units (NAL-only stripping) if dovi_rpu is unavailable.
func setupDVRemovalFilter(inputStream *astiav.Stream) (*astiav.BitStreamFilterContext, error) {
	// Primary: dovi_rpu with strip=1 (strips RPU NALs + cleans DOVI config record)
	bsfName := "dovi_rpu"
	optKey := "strip"
	optVal := "1"

	bsf := astiav.FindBitStreamFilterByName(bsfName)
	if bsf == nil {
		// Fallback: filter_units (strips NALs only, no config cleanup)
		bsfName = "filter_units"
		optKey = "remove_types"
		optVal = "62|63"
		bsf = astiav.FindBitStreamFilterByName(bsfName)
	}
	if bsf == nil {
		return nil, nil // no suitable BSF available, skip
	}

	slog.Info("setting up DV removal BSF", "bsf", bsfName)

	bsfCtx, err := astiav.AllocBitStreamFilterContext(bsf)
	if err != nil {
		return nil, fmt.Errorf("alloc BSF context (%s): %w", bsfName, err)
	}

	// Copy input codec parameters to the BSF
	if err := inputStream.CodecParameters().Copy(bsfCtx.InputCodecParameters()); err != nil {
		bsfCtx.Free()
		return nil, fmt.Errorf("copy codec params to BSF: %w", err)
	}
	bsfCtx.SetInputTimeBase(inputStream.TimeBase())

	// Configure the BSF option
	if pd := bsfCtx.PrivateData(); pd != nil {
		opts := pd.Options()
		searchFlags := astiav.NewOptionSearchFlags()
		_ = opts.Set(optKey, optVal, searchFlags)
	}

	if err := bsfCtx.Initialize(); err != nil {
		bsfCtx.Free()
		return nil, fmt.Errorf("init BSF (%s): %w", bsfName, err)
	}

	return bsfCtx, nil
}

// findStreamByType returns true if a stream of the given type already exists in the mapping.
func findStreamByType(streams map[int]*streamMapping, mt astiav.MediaType) (*streamMapping, bool) {
	for _, sm := range streams {
		if sm.mediaType == mt {
			return sm, true
		}
	}
	return nil, false
}
