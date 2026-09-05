package transform

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/discord/lilliput"
)

type Result struct {
	Body        []byte
	ContentType string
	VaryAccept  bool
}

type Processor struct {
	mu  sync.Mutex
	ops *lilliput.ImageOps
}

func NewProcessor() *Processor {
	return &Processor{ops: lilliput.NewImageOps(4096)}
}

func (p *Processor) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ops.Close()
}

func (p *Processor) Process(source []byte, params Params, accept string) (Result, error) {
	decoder, err := lilliput.NewDecoder(source)
	if err != nil {
		return Result{}, fmt.Errorf("decode image: %w", err)
	}
	defer decoder.Close()

	header, err := decoder.Header()
	if err != nil {
		return Result{}, fmt.Errorf("read image header: %w", err)
	}

	format, contentType := outputFormat(params.Format, accept, decoder.Description())
	width, height := outputSize(params, header.Width(), header.Height())
	options := &lilliput.ImageOptions{
		FileType:              "." + format,
		Width:                 width,
		Height:                height,
		ResizeMethod:          resizeMethod(params, width, height, header.Width(), header.Height()),
		NormalizeOrientation:  true,
		DisableAnimatedOutput: !params.Animated,
		EncodeTimeout:         30 * time.Second,
		EncodeOptions:         encodeOptions(format, params.Quality),
	}

	buffer := make([]byte, max(16<<20, len(source)*2))
	p.mu.Lock()
	result, err := p.ops.Transform(decoder, options, buffer)
	p.mu.Unlock()
	if err != nil {
		return Result{}, fmt.Errorf("transform image: %w", err)
	}

	return Result{Body: result, ContentType: contentType, VaryAccept: params.Format == FormatAuto}, nil
}

func outputFormat(requested Format, accept, source string) (string, string) {
	if requested != FormatAuto {
		return string(requested), "image/" + string(requested)
	}
	if strings.Contains(accept, "image/webp") {
		return "webp", "image/webp"
	}
	if strings.Contains(accept, "image/avif") {
		return "avif", "image/avif"
	}

	format := strings.ToLower(source)
	switch format {
	case "jpeg", "png", "webp", "avif", "gif":
		return format, "image/" + format
	default:
		return "jpeg", "image/jpeg"
	}
}

func outputSize(params Params, sourceWidth, sourceHeight int) (int, int) {
	width, height := params.Width, params.Height
	if width == 0 && height == 0 {
		return sourceWidth, sourceHeight
	}
	if width == 0 {
		width = max(1, sourceWidth*height/sourceHeight)
	}
	if height == 0 {
		height = max(1, sourceHeight*width/sourceWidth)
	}

	if params.Fit == FitContain {
		scale := min(1.0, min(float64(width)/float64(sourceWidth), float64(height)/float64(sourceHeight)))
		width = max(1, int(float64(sourceWidth)*scale))
		height = max(1, int(float64(sourceHeight)*scale))
	}

	return min(width, sourceWidth), min(height, sourceHeight)
}

func resizeMethod(params Params, width, height, sourceWidth, sourceHeight int) lilliput.ImageOpsSizeMethod {
	if params.Fit == FitCover && params.Width != 0 && params.Height != 0 {
		return lilliput.ImageOpsFit
	}
	if width == sourceWidth && height == sourceHeight {
		return lilliput.ImageOpsNoResize
	}
	return lilliput.ImageOpsResize
}

func encodeOptions(format string, quality int) map[int]int {
	switch format {
	case "jpeg":
		return map[int]int{lilliput.JpegQuality: quality}
	case "webp":
		return map[int]int{lilliput.WebpQuality: quality}
	case "avif":
		return map[int]int{lilliput.AvifQuality: quality, lilliput.AvifSpeed: 10}
	default:
		return nil
	}
}
