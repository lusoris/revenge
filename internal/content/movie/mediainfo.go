package movie

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astiav"
)

// MediaInfoProber extracts media information from files using FFmpeg (Unix/Linux/macOS only)
type MediaInfoProber struct{}

// Ensure MediaInfoProber implements Prober
var _ Prober = (*MediaInfoProber)(nil)

// NewMediaInfoProber creates a new media info prober
func NewMediaInfoProber() *MediaInfoProber {
	return &MediaInfoProber{}
}

// Probe extracts media information from the given file path
func (p *MediaInfoProber) Probe(filePath string) (*MediaInfo, error) {
	// Allocate format context
	formatCtx := astiav.AllocFormatContext()
	if formatCtx == nil {
		return nil, fmt.Errorf("failed to allocate format context")
	}
	defer formatCtx.Free()

	// Open input file
	if err := formatCtx.OpenInput(filePath, nil, nil); err != nil {
		return nil, fmt.Errorf("failed to open input: %w", err)
	}
	defer formatCtx.CloseInput()

	// Find stream info
	if err := formatCtx.FindStreamInfo(nil); err != nil {
		return nil, fmt.Errorf("failed to find stream info: %w", err)
	}

	// Build MediaInfo
	info := &MediaInfo{
		FilePath:        filePath,
		Container:       getContainerFormat(filePath, formatCtx),
		DurationSeconds: float64(formatCtx.Duration()) / float64(astiav.TimeBase),
		BitrateKbps:     formatCtx.BitRate() / 1000,
	}

	// Process streams
	for _, stream := range formatCtx.Streams() {
		codecParams := stream.CodecParameters()

		switch codecParams.MediaType() {
		case astiav.MediaTypeVideo:
			if info.VideoCodec == "" { // Take first video stream
				p.processVideoStream(stream, info)
			}
		case astiav.MediaTypeAudio:
			audioInfo := p.processAudioStream(stream, len(info.AudioStreams))
			info.AudioStreams = append(info.AudioStreams, audioInfo)
		case astiav.MediaTypeSubtitle:
			subInfo := p.processSubtitleStream(stream, len(info.SubtitleStreams))
			info.SubtitleStreams = append(info.SubtitleStreams, subInfo)
		}
	}

	return info, nil
}

// processVideoStream extracts video stream information
func (p *MediaInfoProber) processVideoStream(stream *astiav.Stream, info *MediaInfo) {
	codecParams := stream.CodecParameters()

	// Codec info - use CodecID.Name() which is available
	codecID := codecParams.CodecID()
	info.VideoCodec = codecID.Name()
	info.VideoCodecLong = codecID.String()

	// Profile - convert to string representation
	profile := codecParams.Profile()
	info.VideoProfile = getProfileName(codecID, profile)

	// Resolution
	info.Width = codecParams.Width()
	info.Height = codecParams.Height()
	info.Resolution = fmt.Sprintf("%dx%d", info.Width, info.Height)
	info.ResolutionLabel = getResolutionLabel(info.Width, info.Height)

	// Framerate
	avgFramerate := stream.AvgFrameRate()
	if avgFramerate.Den() > 0 {
		info.Framerate = float64(avgFramerate.Num()) / float64(avgFramerate.Den())
	}

	// Pixel format
	info.PixelFormat = codecParams.PixelFormat().String()

	// Color info - use the type values directly
	info.ColorSpace = getColorSpaceName(codecParams.ColorSpace())
	info.ColorPrimaries = getColorPrimariesName(codecParams.ColorPrimaries())
	info.ColorTransfer = getColorTransferName(codecParams.ColorTransferCharacteristic())
	info.ColorRange = getColorRangeName(codecParams.ColorRange())

	// Determine dynamic range
	info.DynamicRange = detectDynamicRange(codecParams)

	// Build RFC 6381 CODECS string from extradata (hvcC / avcC)
	info.VideoCodecString = buildVideoCodecString(codecID, codecParams)

	// Bitrate
	if codecParams.BitRate() > 0 {
		info.VideoBitrateKbps = codecParams.BitRate() / 1000
	}
}

// processAudioStream extracts audio stream information
func (p *MediaInfoProber) processAudioStream(stream *astiav.Stream, index int) AudioStreamInfo {
	codecParams := stream.CodecParameters()

	audioInfo := AudioStreamInfo{
		Index:      index,
		SampleRate: codecParams.SampleRate(),
	}

	// Codec info
	codecID := codecParams.CodecID()
	audioInfo.Codec = codecID.Name()
	audioInfo.CodecLong = codecID.String()

	// Channel layout
	channelLayout := codecParams.ChannelLayout()
	audioInfo.Channels = channelLayout.Channels()
	audioInfo.Layout = getChannelLayoutName(audioInfo.Channels)

	// Bitrate
	if codecParams.BitRate() > 0 {
		audioInfo.BitrateKbps = codecParams.BitRate() / 1000
	}

	// Metadata (language, title)
	if metadata := stream.Metadata(); metadata != nil {
		if entry := metadata.Get("language", nil, astiav.NewDictionaryFlags()); entry != nil {
			audioInfo.Language = entry.Value()
		}
		if entry := metadata.Get("title", nil, astiav.NewDictionaryFlags()); entry != nil {
			audioInfo.Title = entry.Value()
		}
	}

	// Default flag
	audioInfo.IsDefault = stream.DispositionFlags().Has(astiav.DispositionFlagDefault)

	return audioInfo
}

// processSubtitleStream extracts subtitle stream information
func (p *MediaInfoProber) processSubtitleStream(stream *astiav.Stream, index int) SubtitleStreamInfo {
	codecParams := stream.CodecParameters()

	subInfo := SubtitleStreamInfo{
		Index: index,
	}

	// Codec info
	codecID := codecParams.CodecID()
	subInfo.Codec = codecID.Name()

	// Metadata
	if metadata := stream.Metadata(); metadata != nil {
		if entry := metadata.Get("language", nil, astiav.NewDictionaryFlags()); entry != nil {
			subInfo.Language = entry.Value()
		}
		if entry := metadata.Get("title", nil, astiav.NewDictionaryFlags()); entry != nil {
			subInfo.Title = entry.Value()
		}
	}

	// Flags
	disposition := stream.DispositionFlags()
	subInfo.IsForced = disposition.Has(astiav.DispositionFlagForced)
	subInfo.IsDefault = disposition.Has(astiav.DispositionFlagDefault)

	return subInfo
}

// getContainerFormat returns the container format name
func getContainerFormat(filePath string, formatCtx *astiav.FormatContext) string {
	// Try to get from format context
	if inputFmt := formatCtx.InputFormat(); inputFmt != nil {
		return inputFmt.Name()
	}
	// Fallback to extension
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	return strings.ToLower(ext)
}

// getProfileName returns a human-readable profile name
func getProfileName(codecID astiav.CodecID, profile astiav.Profile) string {
	// Different codecs have overlapping profile values, so we need to check codec first
	switch codecID {
	case astiav.CodecIDH264:
		switch profile {
		case astiav.ProfileH264Baseline:
			return "Baseline"
		case astiav.ProfileH264ConstrainedBaseline:
			return "Constrained Baseline"
		case astiav.ProfileH264Main:
			return "Main"
		case astiav.ProfileH264Extended:
			return "Extended"
		case astiav.ProfileH264High:
			return "High"
		case astiav.ProfileH264High10:
			return "High 10"
		case astiav.ProfileH264High422:
			return "High 4:2:2"
		case astiav.ProfileH264High444Predictive:
			return "High 4:4:4 Predictive"
		}
	case astiav.CodecIDHevc:
		switch profile {
		case astiav.ProfileHevcMain:
			return "Main"
		case astiav.ProfileHevcMain10:
			return "Main 10"
		case astiav.ProfileHevcMainStillPicture:
			return "Main Still Picture"
		}
	case astiav.CodecIDAv1:
		switch profile {
		case astiav.ProfileAv1Main:
			return "Main"
		case astiav.ProfileAv1High:
			return "High"
		case astiav.ProfileAv1Professional:
			return "Professional"
		}
	}
	// Unknown codec/profile combo
	if profile >= 0 {
		return fmt.Sprintf("Profile %d", int(profile))
	}
	return "Unknown"
}

// getColorSpaceName returns a human-readable color space name
func getColorSpaceName(cs astiav.ColorSpace) string {
	switch cs {
	case astiav.ColorSpaceBt709:
		return "BT.709"
	case astiav.ColorSpaceBt2020Ncl:
		return "BT.2020 NCL"
	case astiav.ColorSpaceBt2020Cl:
		return "BT.2020 CL"
	case astiav.ColorSpaceSmpte170M:
		return "SMPTE 170M"
	case astiav.ColorSpaceSmpte240M:
		return "SMPTE 240M"
	case astiav.ColorSpaceUnspecified:
		return "Unspecified"
	default:
		return fmt.Sprintf("ColorSpace(%d)", int(cs))
	}
}

// getColorPrimariesName returns a human-readable color primaries name
func getColorPrimariesName(cp astiav.ColorPrimaries) string {
	switch cp {
	case astiav.ColorPrimariesBt709:
		return "BT.709"
	case astiav.ColorPrimariesBt2020:
		return "BT.2020"
	case astiav.ColorPrimariesBt470M:
		return "BT.470M"
	case astiav.ColorPrimariesBt470Bg:
		return "BT.470BG"
	case astiav.ColorPrimariesSmpte170M:
		return "SMPTE 170M"
	case astiav.ColorPrimariesSmpte240M:
		return "SMPTE 240M"
	case astiav.ColorPrimariesUnspecified:
		return "Unspecified"
	default:
		return fmt.Sprintf("ColorPrimaries(%d)", int(cp))
	}
}

// getColorTransferName returns a human-readable color transfer name
func getColorTransferName(ct astiav.ColorTransferCharacteristic) string {
	switch ct {
	case astiav.ColorTransferCharacteristicBt709:
		return "BT.709"
	case astiav.ColorTransferCharacteristicSmptest2084:
		return "SMPTE ST 2084 (PQ)"
	case astiav.ColorTransferCharacteristicAribStdB67:
		return "ARIB STD-B67 (HLG)"
	case astiav.ColorTransferCharacteristicLinear:
		return "Linear"
	case astiav.ColorTransferCharacteristicGamma22:
		return "Gamma 2.2"
	case astiav.ColorTransferCharacteristicGamma28:
		return "Gamma 2.8"
	case astiav.ColorTransferCharacteristicUnspecified:
		return "Unspecified"
	default:
		return fmt.Sprintf("ColorTransfer(%d)", int(ct))
	}
}

// getColorRangeName returns a human-readable color range name
func getColorRangeName(cr astiav.ColorRange) string {
	switch cr {
	case astiav.ColorRangeMpeg:
		return "Limited"
	case astiav.ColorRangeJpeg:
		return "Full"
	case astiav.ColorRangeUnspecified:
		return "Unspecified"
	default:
		return fmt.Sprintf("ColorRange(%d)", int(cr))
	}
}

// detectDynamicRange determines the dynamic range of the video
func detectDynamicRange(codecParams *astiav.CodecParameters) string {
	colorTransfer := codecParams.ColorTransferCharacteristic()
	colorPrimaries := codecParams.ColorPrimaries()

	// Check for HDR indicators
	switch {
	case colorTransfer == astiav.ColorTransferCharacteristicSmptest2084:
		// PQ transfer function indicates HDR10 or Dolby Vision
		// Check side data for Dolby Vision metadata
		if isDolbyVision(codecParams) {
			return "Dolby Vision"
		}
		return "HDR10"
	case colorTransfer == astiav.ColorTransferCharacteristicAribStdB67:
		// HLG transfer function
		return "HLG"
	case colorPrimaries == astiav.ColorPrimariesBt2020:
		// BT.2020 primaries but unknown transfer - could be HDR
		return "HDR"
	default:
		return "SDR"
	}
}

// isDolbyVision checks if the stream contains Dolby Vision metadata.
//
// Detection heuristics in order of reliability:
// 1. Codec tag: dvh1/dvhe (HEVC DV) or dva1/dav1 (AV1 DV) — definitive.
// 2. hvcC constraint byte: DV Profile 8 content uses constraint byte 6 = 0x90
//    (general_non_packed_constraint_flag=0). Standard HEVC Main 10 uses 0xB0.
//    This catches DV-in-MKV where the codec tag may be just hev1/hvc1.
//
// Note: go-astiav v0.40.0 doesn't expose AV_PKT_DATA_DOVI_CONF side data accessors,
// so we can't check DOVI_CONF directly. The above heuristics are sufficient for
// Profile 8 (HEVC + HDR10 fallback) which is the most common DV streaming format.
func isDolbyVision(codecParams *astiav.CodecParameters) bool {
	// Check 1: codec tag (FourCC)
	tag := uint32(codecParams.CodecTag())
	switch tag {
	case
		0x31687664, // dvh1  (Dolby Vision HEVC, hvcC-style)
		0x65687664, // dvhe  (Dolby Vision HEVC, hevC-style)
		0x31617664, // dva1  (Dolby Vision AV1)
		0x31766164: // dav1  (Dolby Vision AV1 alt)
		return true
	}

	// Check 2: hvcC constraint flags for DV Profile 8
	// In the HEVCDecoderConfigurationRecord, byte 6 is the first byte of
	// general_constraint_indicator_flags. DV Profile 8 sets this to 0x90
	// (missing general_non_packed_constraint_flag, bit 5). Standard HEVC Main 10
	// with PQ transfer uses 0xB0 (bit 5 set). This is the same byte that
	// patchHvcCConstraints() fixes in the playback pipeline.
	if codecParams.CodecID() == astiav.CodecIDHevc {
		extradata := codecParams.ExtraData()
		if len(extradata) >= 13 {
			constraintByte := extradata[6]
			// 0x90 = 10010000: DV Profile 8 pattern
			// Bit 7 set (general_progressive_source_flag), bit 4 set (general_frame_only_constraint_flag)
			// Bit 5 NOT set (general_non_packed_constraint_flag) — this is the DV tell
			if constraintByte == 0x90 {
				return true
			}
		}
	}

	return false
}

// buildVideoCodecString builds an RFC 6381 CODECS string from codec extradata.
// For HEVC it parses the HEVCDecoderConfigurationRecord (ISO 14496-15 Section 8.3.3.1)
// and produces the correct hvc1.{A}.{B}.{C}.{D} string per Section E.3.
// Returns empty string if extradata is missing or unparseable — caller should fall back.
func buildVideoCodecString(codecID astiav.CodecID, codecParams *astiav.CodecParameters) string {
	switch codecID.Name() {
	case "hevc":
		return buildHEVCCodecString(codecParams.ExtraData())
	case "h264":
		return buildH264CodecString(codecParams.ExtraData())
	default:
		return ""
	}
}

// buildHEVCCodecString parses the HEVCDecoderConfigurationRecord and builds
// the RFC 6381 codec string: hvc1.{A}.{B}.{C}.{D}
//
// Record layout (ISO 14496-15 Section 8.3.3.1):
//
//	Byte 0:    configurationVersion (1)
//	Byte 1:    [general_profile_space(2) | general_tier_flag(1) | general_profile_idc(5)]
//	Bytes 2-5: general_profile_compatibility_flags (32 bits, MSB first)
//	Bytes 6-11: general_constraint_indicator_flags (48 bits / 6 bytes)
//	Byte 12:   general_level_idc
//
// CODECS string format (ISO 14496-15 Section E.3):
//
//	A = {profile_space_char}{general_profile_idc}  (space char: empty/A/B/C for 0/1/2/3)
//	B = general_profile_compatibility_flags as hex, with REVERSED bit order
//	C = {general_tier_flag: L or H}{general_level_idc}
//	D = constraint bytes as hex, dot-separated, trailing zero bytes omitted
func buildHEVCCodecString(extradata []byte) string {
	if len(extradata) < 13 {
		return ""
	}

	// Byte 1: profile_space (bits 7-6), tier_flag (bit 5), profile_idc (bits 4-0)
	profileSpace := (extradata[1] >> 6) & 0x03
	tierFlag := (extradata[1] >> 5) & 0x01
	profileIDC := extradata[1] & 0x1F

	// Bytes 2-5: general_profile_compatibility_flags (big-endian)
	compatFlags := uint32(extradata[2])<<24 | uint32(extradata[3])<<16 | uint32(extradata[4])<<8 | uint32(extradata[5])

	// Reverse the bit order of the 32-bit compatibility flags per spec
	reversed := reverseBits32(compatFlags)

	// Bytes 6-11: constraint_indicator_flags (6 bytes)
	constraints := extradata[6:12]

	// Byte 12: general_level_idc
	levelIDC := extradata[12]

	// Build component A: profile_space prefix + profile_idc
	spacePrefix := ""
	switch profileSpace {
	case 1:
		spacePrefix = "A"
	case 2:
		spacePrefix = "B"
	case 3:
		spacePrefix = "C"
	}

	// Build component C: tier + level
	tierChar := "L"
	if tierFlag == 1 {
		tierChar = "H"
	}

	// Build component D: constraint bytes, trailing zeros omitted, dot-separated
	constraintStr := buildConstraintString(constraints)

	codecStr := fmt.Sprintf("hvc1.%s%d.%X.%s%d", spacePrefix, profileIDC, reversed, tierChar, levelIDC)
	if constraintStr != "" {
		codecStr += "." + constraintStr
	}

	return codecStr
}

// buildH264CodecString parses the AVCDecoderConfigurationRecord and builds
// the RFC 6381 codec string: avc1.PPCCLL (profile, constraint, level as hex).
//
// Record layout (ISO 14496-15):
//
//	Byte 0: configurationVersion
//	Byte 1: AVCProfileIndication
//	Byte 2: profile_compatibility (constraint_set flags)
//	Byte 3: AVCLevelIndication
func buildH264CodecString(extradata []byte) string {
	if len(extradata) < 4 {
		return ""
	}
	return fmt.Sprintf("avc1.%02X%02X%02X", extradata[1], extradata[2], extradata[3])
}

// reverseBits32 reverses the bit order of a 32-bit value.
// The HEVC CODECS string spec requires the compatibility flags to be
// bit-reversed before hex encoding.
func reverseBits32(v uint32) uint32 {
	var r uint32
	for i := 0; i < 32; i++ {
		r = (r << 1) | (v & 1)
		v >>= 1
	}
	return r
}

// buildConstraintString encodes 6 constraint bytes as dot-separated hex,
// omitting trailing zero bytes.
func buildConstraintString(constraints []byte) string {
	// Find last non-zero byte
	last := -1
	for i := len(constraints) - 1; i >= 0; i-- {
		if constraints[i] != 0 {
			last = i
			break
		}
	}
	if last < 0 {
		return ""
	}
	parts := make([]string, 0, last+1)
	for i := 0; i <= last; i++ {
		parts = append(parts, fmt.Sprintf("%02X", constraints[i]))
	}
	return strings.Join(parts, ".")
}

// NOTE: ToMovieFileInfo, GetAudioLanguages, GetSubtitleLanguages, and GetDurationFormatted
// methods are defined in mediainfo_types.go (platform-independent)
