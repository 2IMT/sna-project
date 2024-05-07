package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
	"sync"
)

type Bot struct {
	api *tgbotapi.BotAPI
	gameQueue []int64
	connectedPlayers map[int64]int64
	mu sync.Mutex
	player1          int64
	player2          int64
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

	result.connectedPlayers = make(map[int64]int64)

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
	go handleUpdates(b.api, result, b.connectedPlayers)
	return result
}

func (b *Bot) AddToGameQueue(chatId int64){
	b.mu.Lock()
	defer b.mu.Unlock()

	b.gameQueue = append(b.gameQueue, chatId)
	if len(b.gameQueue)>=2{
		b.player1=b.gameQueue[0]
		b.player2=b.gameQueue[1]
		b.gameQueue=b.gameQueue[2:]

		b.connectedPlayers[b.player1]=b.player2
		b.connectedPlayers[b.player2]=b.player1

		b.SendMessage(b.player1, "You are now connected with aeaeaea")
		b.SendMessage(b.player2, "You are now connected with sdkdkja")
	}
}

func(b *Bot) ForwardMessage(chatId int64, text string){
	b.mu.Lock()
	defer b.mu.Unlock()

	if opponent, ok := b.connectedPlayers[chatId]; ok {
		b.SendMessage(opponent, text)
	}
}

func handleUpdates(api *tgbotapi.BotAPI, msgChan chan Message, connectedPlayers map[int64]int64) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updatesChan := api.GetUpdatesChan(updateConfig)

	for update := range updatesChan {
		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID
		//username := update.Message.Chat.UserName

		var text string
		if len(update.Message.Text) != 0 {
			text = update.Message.Text
		} else if len(update.Message.Caption) != 0 {
			text = update.Message.Caption
		} else {
			continue
		}

		if text=="/play"{
			B.AddToGameQueue(chatId)
		}else {
            B.ForwardMessage(chatId, text)
        }
	}
}
