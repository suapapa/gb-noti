package draw

import (
	"embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	//go:embed font/Hoengseong_Hanu.ttf
	efs embed.FS

	hanuFont map[float64]font.Face
)

func init() {
	hanuFont = make(map[float64]font.Face)
}

func GetHanuFont(points float64) (font.Face, error) {
	if ff, ok := hanuFont[points]; ok {
		return ff, nil
	}

	fontName := "font/Hoengseong_Hanu.ttf"
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
	hanuFont[points] = nface
	return nface, nil
}
