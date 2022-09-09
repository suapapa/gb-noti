package main

import (
	"encoding/json"
	"testing"
)

func TestPrint(t *testing.T) {
	jsonStr := `{"from":"","msg":"Hello world.\r\n네가없는 거리에서\r\n내가할일이없어서 마냥걷다거다보면","remoteAddr":"10.128.0.7:40172","timestamp":"2022-09-09T12:16:30Z"}`
	var c map[string]string
	json.Unmarshal([]byte(jsonStr), &c)
	if err := print(c); err != nil {
		t.Error(err)
	}
}
