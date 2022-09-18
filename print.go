package main

import (
	"fmt"
	"log"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/suapapa/gb-noti/draw"
	"github.com/suapapa/gb-noti/receipt"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

type chat struct {
	Msg        string `json:"msg"`
	From       string `json:"from"`
	TimeStamp  string `json:"timestamp"`
	RemoteAddr string `json:"remoteAddr"`
	Pork       bool   `json:"pork,omitempty"`
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
	mFF, err := draw.GetFont(48)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines := draw.FitToLines(mFF, receipt.MaxWidth, c.Msg)
	if len(lines) == 0 {
		return fmt.Errorf("no content")
	}

	fromCP949, _, err := transform.String(korean.EUCKR.NewEncoder(), c.From)
	if err != nil {
		log.Printf("WARN: %v", errors.Wrap(err, "my printer only works with CP949 string"))
		fromCP949 = "UNKNOWN"
	}
	fromUTF := fmt.Sprintf("%s\n%s", c.TimeStamp, fromCP949)
	rp.WriteString(fromUTF)

	defer rp.CutPaper()
	lines = append(lines, " ") // TODO: 마지막 줄이 잘려서 패딩 라인 붙임
	for _, l := range lines {
		if img, err := draw.Txt2Img(mFF, receipt.MaxWidth, l); err != nil {
			return errors.Wrap(err, "fail to print")
		} else {
			w := img.Bounds().Dx()
			h := (img.Bounds().Dy() + 7) / 3
			img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
			rp.PrintImage8bitDouble(img)
		}
	}

	return nil
}
