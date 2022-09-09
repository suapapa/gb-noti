package draw

import (
	"embed"
	"sort"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	//go:embed font/NotoSans-Medium.ttf
	//go:embed font/NotoSans-Light.ttf
	efs embed.FS

	fonts map[float64]*font.Face
)

func init() {
	fonts = make(map[float64]*font.Face)
}

func FitLines(ff font.Face, maxWidth int, origTxt string) []string {
	origTxt = strings.ReplaceAll(origTxt, "\r\n", "\n")
	// origTxt = strings.ReplaceAll(origTxt, "\r\n", "\n")
	lines := strings.Split(origTxt, "\n")
	var lineOut []string

	for _, line := range lines {
		rl := []rune(line)
		rlLen := len(rl)
		rlStart := 0
		for rlStart < rlLen {
			rlSub := rl[rlStart:]
			rlSubLen := len(rlSub)

			i := sort.Search(rlSubLen, func(i int) bool {
				w, _ := MeasureTxt(ff, string(rlSub[:i]))
				return w > maxWidth
			})

			if i > rlSubLen {
				i = rlSubLen
			}

			w, _ := MeasureTxt(ff, string(rlSub[:i]))
			for w > maxWidth {
				i -= 1
				w, _ = MeasureTxt(ff, string(rlSub[:i]))
			}

			sl := string(rl[rlStart : rlStart+i])
			lineOut = append(lineOut, sl)

			rlStart += i
		}
	}
	return lineOut
}

func MeasureTxt(ff font.Face, txt string) (w, h int) {
	d := &font.Drawer{
		Face: ff,
	}
	w = int(d.MeasureString(txt) >> 6)
	h = int(ff.Metrics().Height >> 6)
	return
}

func GetFont(fontName string, points float64) (font.Face, error) {
	fontName = "font/" + fontName + ".ttf"
	if ff, ok := fonts[points]; ok {
		return *ff, nil
	}

	data, err := efs.ReadFile("font/NotoSans-Medium.ttf")
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
	fonts[points] = &nface
	return nface, nil
}
