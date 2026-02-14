// Package subtitle provides subtitle extraction from media files to WebVTT format.
package subtitle

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/asticode/go-astiav"
)

// ExtractToWebVTT extracts a subtitle track from a media file to WebVTT format
// using astiav's in-process FFmpeg libraries (no subprocess).
// trackIndex is the subtitle stream index (0-based within subtitle streams).
func ExtractToWebVTT(ctx context.Context, inputFile, outputDir string, trackIndex int) (string, error) {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return "", fmt.Errorf("failed to create subtitle output dir: %w", err)
	}

	outputFile := filepath.Join(outputDir, strconv.Itoa(trackIndex)+".vtt")

	// Open input
	inputFmtCtx := astiav.AllocFormatContext()
	if inputFmtCtx == nil {
		return "", errors.New("failed to allocate input format context")
	}
	defer inputFmtCtx.Free()

	if err := inputFmtCtx.OpenInput(inputFile, nil, nil); err != nil {
		return "", fmt.Errorf("failed to open input %q: %w", inputFile, err)
	}
	defer inputFmtCtx.CloseInput()

	if err := inputFmtCtx.FindStreamInfo(nil); err != nil {
		return "", fmt.Errorf("failed to find stream info: %w", err)
	}

	// Find the target subtitle stream (by relative index within subtitle streams)
	var subStream *astiav.Stream
	subIdx := 0
	for _, s := range inputFmtCtx.Streams() {
		if s.CodecParameters().MediaType() == astiav.MediaTypeSubtitle {
			if subIdx == trackIndex {
				subStream = s
				break
			}
			subIdx++
		}
	}
	if subStream == nil {
		return "", fmt.Errorf("subtitle track %d not found", trackIndex)
	}

	// Find decoder
	decoder := astiav.FindDecoder(subStream.CodecParameters().CodecID())
	if decoder == nil {
		return "", fmt.Errorf("decoder not found for subtitle codec %s", subStream.CodecParameters().CodecID().Name())
	}

	decCtx := astiav.AllocCodecContext(decoder)
	if decCtx == nil {
		return "", errors.New("failed to allocate decoder codec context")
	}
	defer decCtx.Free()

	if err := subStream.CodecParameters().ToCodecContext(decCtx); err != nil {
		return "", fmt.Errorf("failed to copy codec params to decoder: %w", err)
	}
	if err := decCtx.Open(decoder, nil); err != nil {
		return "", fmt.Errorf("failed to open subtitle decoder: %w", err)
	}

	// Open output (WebVTT muxer)
	outputFmtCtx, err := astiav.AllocOutputFormatContext(nil, "webvtt", outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to allocate output format context: %w", err)
	}
	if outputFmtCtx == nil {
		return "", errors.New("output format context is nil")
	}
	defer outputFmtCtx.Free()

	// Create output stream — for webvtt we copy the subtitle codec params
	outStream := outputFmtCtx.NewStream(nil)
	if outStream == nil {
		return "", errors.New("failed to create output stream")
	}
	if err := subStream.CodecParameters().Copy(outStream.CodecParameters()); err != nil {
		return "", fmt.Errorf("failed to copy codec parameters: %w", err)
	}
	outStream.CodecParameters().SetCodecTag(0)

	// Open IO context
	if !outputFmtCtx.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		ioCtx, err := astiav.OpenIOContext(outputFile, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
		if err != nil {
			return "", fmt.Errorf("failed to open output IO context: %w", err)
		}
		defer func() { _ = ioCtx.Close() }()
		outputFmtCtx.SetPb(ioCtx)
	}

	// Write header
	if err := outputFmtCtx.WriteHeader(nil); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	// Read packets and remux subtitle
	pkt := astiav.AllocPacket()
	if pkt == nil {
		return "", errors.New("failed to allocate packet")
	}
	defer pkt.Free()

	for {
		if ctx.Err() != nil {
			break
		}

		if err := inputFmtCtx.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			return "", fmt.Errorf("failed to read frame: %w", err)
		}

		// Only process packets from our subtitle stream
		if pkt.StreamIndex() != subStream.Index() {
			pkt.Unref()
			continue
		}

		// Remap stream index and rescale timestamps
		pkt.SetStreamIndex(outStream.Index())
		pkt.RescaleTs(subStream.TimeBase(), outStream.TimeBase())
		pkt.SetPos(-1)

		if err := outputFmtCtx.WriteInterleavedFrame(pkt); err != nil {
			pkt.Unref()
			return "", fmt.Errorf("failed to write subtitle packet: %w", err)
		}
		pkt.Unref()
	}

	// Write trailer
	if err := outputFmtCtx.WriteTrailer(); err != nil {
		return "", fmt.Errorf("failed to write trailer: %w", err)
	}

	return outputFile, nil
}
