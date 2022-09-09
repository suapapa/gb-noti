package main

import (
	"fmt"

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
	mFF, err := draw.GetHanuFont(40)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines := draw.FitLines(mFF, receipt.MaxWidth, c.Msg)
	for _, l := range lines {
		if err := draw.Txt2Img(mFF, l); err != nil {
			return errors.Wrap(err, "fail to print")
		}
	}

	fFF, err := draw.GetHanuFont(20)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines = draw.FitLines(fFF, receipt.MaxWidth, fmt.Sprintf("%s(%s)", c.From, c.RemoteAddr))
	for _, l := range lines {
		if err := draw.Txt2Img(fFF, l); err != nil {
			return errors.Wrap(err, "fail to print")
		}
	}
	return nil
}
