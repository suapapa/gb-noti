package main

import (
	"testing"
)

const (
	gbMsg = `{"type":"gb","data":{"from":"이호민","content":"잘 보고 갑니다","ts":"2022-10-03T09:05:39+09:00"}}`
)

func TestGBMsg(t *testing.T) {
	gb, err := getGBFromMsg([]byte(gbMsg))
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s-%s", gb.Content, gb.From)
}
