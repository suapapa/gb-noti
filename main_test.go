package main

import (
	"encoding/json"
	"testing"
)

func TestPrint(t *testing.T) {
	jsonStr := `{"from":"","msg":"우리의 소원은 통일\r\n꿈에도 소원은 통일\r\n이 정성 다해서 통일\r\n통일을 이루자\r\n\r\n이 겨레 살리는 통일\r\n이 나라 살리는 통일\r\n통일이여 어서오라\r\n통일이여 오라","remoteAddr":"10.128.0.7:42213","timestamp":"2022-09-09T13:40:12Z"}`
	var c chat
	json.Unmarshal([]byte(jsonStr), &c)
	if err := debug(&c); err != nil {
		t.Error(err)
	}
}
