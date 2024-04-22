package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

type Message struct {
	ChatId   int64
	Username string
	Text     string
}

func NewBot() (Bot, error) {
	var result Bot

	var err error
	result.api, err = tgbotapi.NewBotAPI(Env.BotToken)
	if err != nil {
		return Bot{}, err
	}

	return result, nil
}

func (b Bot) SendMessage(chatId int64, text string) {
    msg := tgbotapi.NewMessage(chatId, text)
    if _, err := b.api.Send(msg); err != nil {
        log.Printf("[ERROR] failed to send message to chat %d: %s\n", chatId, err);
    }
}

func (b Bot) GetMessageChan() chan Message {
	result := make(chan Message)
	go handleUpdates(b.api, result)
	return result
}

func handleUpdates(api *tgbotapi.BotAPI, msgChan chan Message) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updatesChan := api.GetUpdatesChan(updateConfig)

	for update := range updatesChan {
		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID
		username := update.Message.Chat.UserName

		var text string
		if len(update.Message.Text) != 0 {
			text = update.Message.Text
		} else if len(update.Message.Caption) != 0 {
			text = update.Message.Caption
		} else {
			continue
		}

		msgChan <- Message{chatId, username, text}
	}
}
