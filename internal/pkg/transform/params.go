package transform

import (
	"fmt"
	"strconv"
)

type Fit string

const (
	FitCover   Fit = "cover"
	FitContain Fit = "contain"
)

type Format string

const (
	FormatAuto Format = ""
	FormatWebP Format = "webp"
	FormatAVIF Format = "avif"
	FormatPNG  Format = "png"
	FormatJPEG Format = "jpeg"
	FormatGIF  Format = "gif"
)

type Params struct {
	Width    int
	Height   int
	Fit      Fit
	Format   Format
	Quality  int
	Animated bool
}

func ParseParams(queries map[string]string) (Params, error) {
	p := Params{Fit: FitCover, Quality: 80}
	var err error

	if value, ok := queries["size"]; ok {
		p.Width, err = parseInt("size", value)
		if err != nil || !isSize(p.Width) {
			return Params{}, fmt.Errorf("invalid size")
		}
		p.Height = p.Width
	}
	if value, ok := queries["width"]; ok && p.Width == 0 {
		p.Width, err = parseBoundedInt("width", value)
		if err != nil {
			return Params{}, err
		}
	}
	if value, ok := queries["height"]; ok && p.Height == 0 {
		p.Height, err = parseBoundedInt("height", value)
		if err != nil {
			return Params{}, err
		}
	}
	if value, ok := queries["fit"]; ok {
		p.Fit = Fit(value)
		if p.Fit != FitCover && p.Fit != FitContain {
			return Params{}, fmt.Errorf("invalid fit")
		}
	}
	if value, ok := queries["format"]; ok {
		switch Format(value) {
		case FormatWebP, FormatAVIF, FormatPNG, FormatJPEG, FormatGIF:
			p.Format = Format(value)
		default:
			return Params{}, fmt.Errorf("invalid format")
		}
	}
	if value, ok := queries["quality"]; ok {
		p.Quality, err = parseInt("quality", value)
		if err != nil || p.Quality < 1 || p.Quality > 100 {
			return Params{}, fmt.Errorf("invalid quality")
		}
	}
	if value, ok := queries["animated"]; ok {
		if value != "true" && value != "false" {
			return Params{}, fmt.Errorf("invalid animated")
		}
		p.Animated = value == "true"
	}

	return p, nil
}

func parseInt(name, value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return parsed, nil
}

func parseBoundedInt(name, value string) (int, error) {
	parsed, err := parseInt(name, value)
	if err != nil || parsed < 1 || parsed > 4096 {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return parsed, nil
}

func isSize(value int) bool {
	switch value {
	case 16, 32, 64, 128, 256, 512, 1024, 2048, 4096:
		return true
	default:
		return false
	}
}
