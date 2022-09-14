package main

import (
	"fmt"
	"log"
	"time"

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
	mFF, err := draw.GetHandWritingFont(36)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines := draw.FitToLines(mFF, receipt.MaxWidth, c.Msg)
	if len(lines) == 0 {
		return fmt.Errorf("no content")
	}

	tsStr, err := makeKSTstr(c.TimeStamp)
	if err != nil {
		log.Printf("WARN: %v", err)
	}
	fromCP949, _, err := transform.String(korean.EUCKR.NewEncoder(), c.From)
	if err != nil {
		log.Printf("WARN: %v", errors.Wrap(err, "my printer only works with CP949 string"))
		fromCP949 = "UNKNOWN"
	}
	fromUTF := fmt.Sprintf("%s\n%s", tsStr, fromCP949)
	rp.WriteString(fromUTF)

	defer rp.CutPaper()
	lines = append(lines, " ") // TODO: 마지막 줄이 잘려서 패딩 라인 붙임
	for _, l := range lines {
		if img, err := draw.Txt2Img(mFF, receipt.MaxWidth, l); err != nil {
			return errors.Wrap(err, "fail to print")
		} else {
			rp.PrintImage8bitDouble(img)
		}
	}

	return nil
}

func makeKSTstr(timestamp string) (string, error) {
	ts, err := time.Parse(timestamp, time.RFC3339)
	if err != nil {
		return timestamp, errors.Wrap(err, "fail to make kst str")
	}
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return timestamp, errors.Wrap(err, "fail to make kst str")
	}
	kstTS := ts.In(loc)
	return kstTS.Format(time.RFC3339), nil
}
