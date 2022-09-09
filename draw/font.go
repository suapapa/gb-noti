package draw

import (
	"log"
	"sort"
	"strings"

	"golang.org/x/image/font"
)

func FitLines(ff font.Face, maxWidth int, origTxt string) []string {
	origTxt = strings.ReplaceAll(origTxt, "\r\n", "\n")
	log.Println(origTxt)
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
	log.Println(outLines)
	return outLines
}

func MeasureTxt(ff font.Face, txt string) (w, h int) {
	d := &font.Drawer{
		Face: ff,
	}
	w = int(d.MeasureString(txt) >> 6)
	h = int(ff.Metrics().Height >> 6)
	return
}
