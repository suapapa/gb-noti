package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/pkg/errors"
	"github.com/suapapa/gb-noti/draw"
	"github.com/suapapa/gb-noti/receipt"
)

type chat struct {
	Msg        string `json:"msg"`
	From       string `json:"from"`
	RemoteAddr string `json:"remoteAddr"`
}

func print(c *chat) error {
	lineCnt := 0
	defer func() {
		if lineCnt > 0 {
			rp.CutPaper()
		}
	}()

	mFF, err := draw.GetHanuFont(32)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines := draw.FitLines(mFF, receipt.MaxWidth, c.Msg)
	for _, l := range lines {
		if img, err := draw.Txt2Img(mFF, l); err != nil {
			return errors.Wrap(err, "fail to print")
		} else {
			rp.PrintImage8bitDouble(img)
		}
	}

	fFF, err := draw.GetHanuFont(16)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines = draw.FitLines(fFF, receipt.MaxWidth, fmt.Sprintf("%s(%s)", c.From, c.RemoteAddr))
	for _, l := range lines {
		if img, err := draw.Txt2Img(fFF, l); err != nil {
			return errors.Wrap(err, "fail to print")
		} else {
			rp.PrintImage8bitDouble(img)
		}
	}
	return nil
}

var i = 0

func debug(c *chat) error {
	lineCnt := 0
	defer func() {
		if lineCnt > 0 {
			// rp.CutPaper()
		}
	}()

	mFF, err := draw.GetHanuFont(40)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines := draw.FitLines(mFF, receipt.MaxWidth, c.Msg)
	for _, l := range lines {
		if img, err := draw.Txt2Img(mFF, l); err != nil {
			return errors.Wrap(err, "fail to print")
		} else {
			// rp.PrintImage8bitDouble(img)
			f, err := os.Create(fmt.Sprintf("img-%d.png", i))
			if err != nil {
				panic(err)
			}
			png.Encode(f, img)
			i += 1
			f.Close()
		}
	}

	fFF, err := draw.GetHanuFont(20)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines = draw.FitLines(fFF, receipt.MaxWidth, fmt.Sprintf("%s(%s)", c.From, c.RemoteAddr))
	for _, l := range lines {
		if img, err := draw.Txt2Img(fFF, l); err != nil {
			return errors.Wrap(err, "fail to print")
		} else {
			// rp.PrintImage8bitDouble(img)
			f, err := os.Create(fmt.Sprintf("img-%d.png", i))
			if err != nil {
				panic(err)
			}
			png.Encode(f, img)
			i += 1
			f.Close()
		}
	}
	return nil
}
