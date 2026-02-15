package transcode

/*
#include <libavcodec/avcodec.h>
#include <libavcodec/packet.h>

// removeDOVIConf strips AV_PKT_DATA_DOVI_CONF from coded_side_data.
// This prevents FFmpeg's fMP4 muxer from writing the 'dby1' compatible brand
// in the ftyp box, which triggers Chrome MSE to reject the init segment.
static int remove_dovi_conf_side_data(AVCodecParameters *par) {
    if (!par || !par->coded_side_data || par->nb_coded_side_data == 0) {
        return 0;
    }
    av_packet_side_data_remove(par->coded_side_data, &par->nb_coded_side_data,
                               AV_PKT_DATA_DOVI_CONF);
    return 1;
}
*/
import "C"

import (
	"log/slog"
	"unsafe"

	"github.com/asticode/go-astiav"
)

// stripDOVIConfSideData removes AV_PKT_DATA_DOVI_CONF from the codec parameters'
// coded_side_data. Without this, FFmpeg's fMP4 muxer writes 'dby1' as a compatible
// brand in the ftyp box of the init segment. Chrome recognizes 'dby1' as Dolby Vision
// and rejects the stream with bufferAppendingError since it can't decode DV.
//
// This MUST be called before WriteHeader for the output format context, because
// WriteHeader is when the muxer inspects side data to determine ftyp brands.
//
// Uses unsafe.Pointer to reach the private C struct inside go-astiav's CodecParameters.
// The struct layout is: type CodecParameters struct { c *C.AVCodecParameters }
func stripDOVIConfSideData(cp *astiav.CodecParameters) {
	// go-astiav.CodecParameters is: struct { c *C.AVCodecParameters }
	// Extract the *C.AVCodecParameters via unsafe pointer arithmetic.
	type codecParametersLayout struct {
		c *C.AVCodecParameters
	}
	layout := (*codecParametersLayout)(unsafe.Pointer(cp))
	if layout.c == nil {
		return
	}

	result := C.remove_dovi_conf_side_data(layout.c)
	if result == 1 {
		slog.Info("stripped DOVI_CONF side data from codec parameters")
	}
}
