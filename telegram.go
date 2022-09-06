package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

var (
	tgBot     *tgbotapi.BotAPI
	tgAPIToke = os.Getenv("TELEGRAM_APITOKEN")
	tgRoomID  = 1236355825
)

func sendMsgToTelegram(msg []byte) error {
	var err error
	if tgBot == nil {
		tgBot, err = tgbotapi.NewBotAPI(tgAPIToke)
		if err != nil {
			return errors.Wrap(err, "fail to send msg to telegram")
		}
		// tgBot.Debug = true
	}

	// TODO: parse pretty msg from incomming json msg bytes
	c := tgbotapi.NewMessage(int64(tgRoomID), string(msg))
	if _, err := tgBot.Send(c); err != nil {
		return errors.Wrap(err, "fail to send msg to telegram")
	}
	return nil
}
