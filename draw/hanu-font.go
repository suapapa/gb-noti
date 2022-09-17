package draw

import (
	"embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	//go:embed font/NotoSans-Regular.ttf
	efs embed.FS

	handwritingFont map[float64]font.Face
)

func init() {
	handwritingFont = make(map[float64]font.Face)
}

func GetHandWritingFont(points float64) (font.Face, error) {
	if ff, ok := handwritingFont[points]; ok {
		return ff, nil
	}

	fontName := "font/NotoSans-Regular.ttf"
	data, err := efs.ReadFile(fontName)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	nface := truetype.NewFace(f, &truetype.Options{
		Size:    points,
		Hinting: font.HintingFull,
		// Hinting: font.HintingNone,
	})
	handwritingFont[points] = nface
	return nface, nil
}
