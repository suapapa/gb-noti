package main

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/suapapa/gb-noti/draw"
	"github.com/suapapa/gb-noti/receipt"
)

func print(c map[string]string) error {

	mFF, err := draw.GetHanuFont(40)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}

	lines := draw.FitLines(mFF, receipt.MaxWidth, c["msg"])
	for i, l := range lines {
		w, h := draw.MeasureTxt(mFF, l)

		dc := gg.NewContext(w+4, h+4)
		dc.SetColor(color.White)
		dc.Clear()
		dc.SetRGB(0, 0, 0)
		dc.SetFontFace(mFF)
		dc.DrawStringAnchored(l, 2, 2, 0, 0.8)
		if err := dc.SavePNG(fmt.Sprintf("out_%d.png", i)); err != nil {
			return errors.Wrap(err, "fail to print")
		}
	}
	return nil
}
