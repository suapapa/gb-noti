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

// 전체 메시지를 통 이미지로 만들어 출력
/*
func printToReceipt(c *chat) error {
	mFF, err := draw.GetHanuFont(36)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	img, err := draw.DrawLines(mFF, c.Msg, receipt.MaxWidth)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	defer rp.CutPaper()

	rp.WriteString(fmt.Sprintf("%s(%s)", c.From, c.RemoteAddr))
	rp.PrintImage8bitDouble(img)

	return nil
}
*/

// 각 줄을 이미지로 만들어 출력
func printToReceipt(c *chat) error {
	defer rp.CutPaper()
	rp.WriteString(fmt.Sprintf("%s(%s)", c.From, c.RemoteAddr))

	mFF, err := draw.GetHanuFont(36)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines := draw.FitLines(mFF, receipt.MaxWidth, c.Msg)
	for _, l := range lines {
		if img, err := draw.Txt2Img(mFF, receipt.MaxWidth, l); err != nil {
			return errors.Wrap(err, "fail to print")
		} else {
			rp.PrintImage8bitDouble(img)
		}
	}

	return nil
}
