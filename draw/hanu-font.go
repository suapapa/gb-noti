package draw

import (
	"embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	//go:embed font/KOTRA_SONGEULSSI.ttf
	//go:embed font/나눔손글씨_강부장님체.ttf
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

	fontName := "font/나눔손글씨_강부장님체.ttf"
	data, err := efs.ReadFile(fontName)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	nface := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingNone,
		// Hinting: font.HintingFull,
	})
	hanuFont[points] = nface
	return nface, nil
}
