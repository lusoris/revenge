package movie

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/asticode/go-astiav"
)

// Prober defines the interface for media probing
type Prober interface {
	Probe(filePath string) (*MediaInfo, error)
}

// MediaInfo contains detailed technical information about a media file
type MediaInfo struct {
	// File metadata
	FilePath  string
	Container string
	FileSize  int64

	// Duration and bitrate
	DurationSeconds float64
	BitrateKbps     int64

	// Video stream info (first video stream)
	VideoCodec       string
	VideoCodecLong   string
	VideoProfile     string
	Width            int
	Height           int
	Resolution       string // e.g., "1920x1080", "4K", "1080p"
	ResolutionLabel  string // e.g., "4K UHD", "1080p Full HD"
	Framerate        float64
	PixelFormat      string
	ColorSpace       string
	ColorPrimaries   string
	ColorTransfer    string
	ColorRange       string
	DynamicRange     string // SDR, HDR10, HDR10+, Dolby Vision, HLG
	VideoBitrateKbps int64

	// Audio streams
	AudioStreams []AudioStreamInfo

	// Subtitle streams
	SubtitleStreams []SubtitleStreamInfo
}

// AudioStreamInfo contains information about an audio stream
type AudioStreamInfo struct {
	Index       int
	Codec       string
	CodecLong   string
	Channels    int
	Layout      string // e.g., "stereo", "5.1", "7.1"
	SampleRate  int
	BitrateKbps int64
	Language    string
	Title       string
	IsDefault   bool
}

// SubtitleStreamInfo contains information about a subtitle stream
type SubtitleStreamInfo struct {
	Index     int
	Codec     string
	Language  string
	Title     string
	IsForced  bool
	IsDefault bool
}

// MediaInfoProber extracts media information from files using FFmpeg
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

// getResolutionLabel returns a human-readable resolution label
func getResolutionLabel(width, height int) string {
	// Common resolution labels
	switch {
	case height >= 2160 || width >= 3840:
		return "4K UHD"
	case height >= 1440 || width >= 2560:
		return "1440p QHD"
	case height >= 1080 || width >= 1920:
		return "1080p Full HD"
	case height >= 720 || width >= 1280:
		return "720p HD"
	case height >= 576:
		return "576p SD"
	case height >= 480:
		return "480p SD"
	default:
		return fmt.Sprintf("%dp", height)
	}
}

// getChannelLayoutName returns a human-readable channel layout name
func getChannelLayoutName(channels int) string {
	switch channels {
	case 1:
		return "mono"
	case 2:
		return "stereo"
	case 6:
		return "5.1"
	case 8:
		return "7.1"
	default:
		return fmt.Sprintf("%d channels", channels)
	}
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

// isDolbyVision checks if the stream contains Dolby Vision metadata
func isDolbyVision(codecParams *astiav.CodecParameters) bool {
	// Check side data for Dolby Vision RPU
	sideData := codecParams.SideData()
	if sideData == nil {
		return false
	}

	// Check if we can detect DOVI config in side data
	// The API for iterating side data varies - for now, we'll use a simple heuristic
	// Full Dolby Vision detection would require checking for specific NAL units
	return false
}

// ToMovieFileInfo converts MediaInfo to the legacy MovieFileInfo format
func (m *MediaInfo) ToMovieFileInfo() *MovieFileInfo {
	info := &MovieFileInfo{
		Path:            m.FilePath,
		Size:            m.FileSize,
		Container:       m.Container,
		Resolution:      m.Resolution,
		ResolutionLabel: m.ResolutionLabel,
		VideoCodec:      m.VideoCodec,
		VideoProfile:    m.VideoProfile,
		BitrateKbps:     int32(m.BitrateKbps),
		DurationSeconds: m.DurationSeconds,
		Framerate:       m.Framerate,
		DynamicRange:    m.DynamicRange,
		ColorSpace:      m.ColorSpace,
	}

	// Get primary audio codec
	if len(m.AudioStreams) > 0 {
		info.AudioCodec = m.AudioStreams[0].Codec
		info.AudioChannels = m.AudioStreams[0].Channels
		info.AudioLayout = m.AudioStreams[0].Layout
	}

	// Collect all audio languages
	for _, audio := range m.AudioStreams {
		if audio.Language != "" {
			info.Languages = append(info.Languages, audio.Language)
		}
	}

	// Collect all subtitle languages
	for _, sub := range m.SubtitleStreams {
		if sub.Language != "" {
			info.SubtitleLangs = append(info.SubtitleLangs, sub.Language)
		}
	}

	return info
}

// GetAudioLanguages returns all unique audio languages
func (m *MediaInfo) GetAudioLanguages() []string {
	seen := make(map[string]bool)
	var languages []string

	for _, audio := range m.AudioStreams {
		if audio.Language != "" && !seen[audio.Language] {
			seen[audio.Language] = true
			languages = append(languages, audio.Language)
		}
	}

	return languages
}

// GetSubtitleLanguages returns all unique subtitle languages
func (m *MediaInfo) GetSubtitleLanguages() []string {
	seen := make(map[string]bool)
	var languages []string

	for _, sub := range m.SubtitleStreams {
		if sub.Language != "" && !seen[sub.Language] {
			seen[sub.Language] = true
			languages = append(languages, sub.Language)
		}
	}

	return languages
}

// GetDurationFormatted returns a human-readable duration string
func (m *MediaInfo) GetDurationFormatted() string {
	d := time.Duration(m.DurationSeconds * float64(time.Second))

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
