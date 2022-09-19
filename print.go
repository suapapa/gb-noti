package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/suapapa/gb-noti/draw"
	"github.com/suapapa/gb-noti/receipt"
	"github.com/suapapa/site-gb/msg"
	"golang.org/x/image/font"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

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
func printToReceipt(c *msg.GuestBook) error {
	mFF, err := getFont(48)
	if err != nil {
		return errors.Wrap(err, "fail to print")
	}
	lines := draw.FitToLines(mFF, receipt.MaxWidth, c.Content)
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
			w := uint(img.Bounds().Dx())
			h := uint(img.Bounds().Dy())
			if !flagHQ {
				h /= 3
			}

			h = 8 * ((h + 7) / 8)
			img = resize.Resize(w, h, img, resize.Lanczos3)

			if flagHQ {
				if err := rp.PrintImage24bitDouble(img); err != nil {
					return errors.Wrap(err, "fail to print")
				}
			} else {
				if err := rp.PrintImage8bitDouble(img); err != nil {
					return errors.Wrap(err, "fail to print")
				}
			}
		}
	}

	return nil
}

func getFont(size float64) (font.Face, error) {
	if flagFontPath != "" {
		data, err := os.ReadFile(flagFontPath)
		if err != nil {
			return nil, errors.Wrap(err, "fail to load font")
		}
		f, err := truetype.Parse(data)
		if err != nil {
			return nil, errors.Wrap(err, "fail to load font")
		}

		nface := truetype.NewFace(f, &truetype.Options{
			Size:    size,
			Hinting: font.HintingFull,
			// Hinting: font.HintingNone,
		})
		return nface, nil
	}

	return draw.GetFont(size)
}
