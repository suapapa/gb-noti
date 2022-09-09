package draw

import (
	"image"
	"image/color"
	"log"
	"sort"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

func FitLines(ff font.Face, maxWidth int, origTxt string) []string {
	origTxt = strings.ReplaceAll(origTxt, "\r\n", "\n")
	lines := strings.Split(origTxt, "\n")
	var outLines []string

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
			outLines = append(outLines, sl)

			rlStart += i
		}
	}
	return outLines
}

func MeasureTxt(ff font.Face, txt string) (w, h int) {
	d := &font.Drawer{
		Face: ff,
	}
	w = int(d.MeasureString(txt) >> 6)
	h = int(ff.Metrics().Height >> 6)
	w = ((w + 7) / 8) * 8 // set w to multiple 8
	h = ((h + 7) / 8) * 8 // set h to multiple 8

	log.Println(w, h)
	return
}

// var (
// 	i = 0
// )

func Txt2Img(ff font.Face, txt string) (image.Image, error) {
	w, h := MeasureTxt(ff, txt)
	dc := gg.NewContext(w, h)
	dc.SetColor(color.White)
	dc.Clear()
	dc.SetColor(color.Black)
	dc.SetFontFace(ff)
	dc.DrawStringAnchored(txt, float64(w)/2, float64(h)/2, 0.5, 0.3)
	// if err := dc.SavePNG(fmt.Sprintf("out_%d.png", i)); err != nil {
	// 	return nil, errors.Wrap(err, "fail to print")
	// }
	// i += 1
	return dc.Image(), nil
}
