package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
	"sync"
	"net/http"
	"fmt"
	"sort"
)

type Bot struct {
	api *tgbotapi.BotAPI
	gameQueue []Player
	connectedPlayers map[int64]int64
	mu sync.Mutex
	player1          Player
	player2          Player
	previousLastLetter byte
	currentPlayer int
	players []Player
}

type Message struct {
	ChatId   int64
	Username string
	Text     string
}

type Player struct {
	PlayerID int64
	PlayerUsername string
	Score          int
}

func NewPlayer(id int64, username string) Player {
	return Player{
		PlayerID:       id,
		PlayerUsername: username,
	}
}

func (p *Player) Increment() {
    p.Score++
}

func (b *Bot) IncrementScore(playerID int64) {
    b.mu.Lock()
    defer b.mu.Unlock()

    for i := range b.players {
        if b.players[i].PlayerID == playerID {
            b.players[i].Increment()
            return
        }
    }
}

func (b *Bot) FindPlayer(playerID int64) Player{
	for i := range b.players {
        if b.players[i].PlayerID == playerID {
            return b.players[i]
        }
    }
	return Player{}
}
func DisplayLeaderboard() string {
    sort.Slice(B.players, func(i, j int) bool {
        return B.players[i].Score > B.players[j].Score
    })

    var leaderboard string
    leaderboard += "Leaderboard:\n"
    for i, player := range B.players {
        leaderboard += fmt.Sprintf("%d. %s - Score: %d\n", i+1, player.PlayerUsername, player.Score)
    }
    return leaderboard
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

func (b *Bot) AddToGameQueue(chatId int64, username string){
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, player := range b.players {
        if player.PlayerID == chatId {
			return
        }
    }

	newPlayer := NewPlayer(chatId, username)
	b.gameQueue = append(b.gameQueue, newPlayer)

	if len(b.gameQueue)>=2{
		b.player1=b.gameQueue[0]
		b.player2=b.gameQueue[1]
		b.gameQueue=b.gameQueue[2:]

		b.connectedPlayers[b.player1.PlayerID]=b.player2.PlayerID
		b.connectedPlayers[b.player2.PlayerID]=b.player1.PlayerID

    	b.SendMessage(b.player1.PlayerID, "You are now connected with "+b.player2.PlayerUsername)
     	b.SendMessage(b.player2.PlayerID, "You are now connected with "+b.player1.PlayerUsername)

		b.currentPlayer = 1
		b.SendMessage(b.player1.PlayerID, "Turn of "+b.player1.PlayerUsername)
		b.SendMessage(b.player2.PlayerID,"Turn of "+b.player1.PlayerUsername)
	} else {
		b.SendMessage(chatId, "Looking for available players...")
	}
}

func WordExists(word string) (bool, error) {
	resp, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func(b *Bot) ForwardMessage(chatId int64, text string){
	b.mu.Lock()
	defer b.mu.Unlock()

	if opponent, ok := b.connectedPlayers[chatId]; ok {
		exist, _ := WordExists(text)
		if b.currentPlayer == 1{
			b.SendMessage(chatId, "Turn of "+b.player2.PlayerUsername)
			b.SendMessage(opponent,"Turn of "+b.player2.PlayerUsername)
		}else if b.currentPlayer == 2{
			b.SendMessage(chatId, "Turn of "+b.player1.PlayerUsername)
			b.SendMessage(opponent,"Turn of "+b.player1.PlayerUsername)
		}
		if !exist{
			b.SendMessage(chatId, "Such word does not exist. You lose!")
			b.SendMessage(opponent, "You won!")
			b.IncrementScore(opponent)
			return
		} else if b.previousLastLetter != 0 && text[0] != b.previousLastLetter {
			b.SendMessage(chatId, "The word does not start with the previous word's last letter. You lose!")
			b.SendMessage(opponent, "You won!")
			b.IncrementScore(opponent)
			return
		} else if b.currentPlayer == 1 && chatId != b.player1.PlayerID {
            b.SendMessage(chatId, "It's not your turn. You lose!")
			b.SendMessage(opponent, "You won!")
			b.IncrementScore(opponent)
            return
        } else if b.currentPlayer == 2 && chatId != b.player2.PlayerID {
            b.SendMessage(chatId, "It's not your turn. You lose!")
			b.SendMessage(opponent, "You won!")
			b.IncrementScore(opponent)
            return
        } else { 
			b.SendMessage(opponent, text) 
			b.previousLastLetter = text[len(text)-1]
			b.currentPlayer = 3 - b.currentPlayer
		}
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
		username := update.Message.Chat.UserName

		var text string
		if len(update.Message.Text) != 0 {
			text = update.Message.Text
		} else if len(update.Message.Caption) != 0 {
			text = update.Message.Caption
		} else {
			continue
		}

		if text=="/play"{
			B.AddToGameQueue(chatId, username)
		} else if text=="/leaderbord"{
			B.SendMessage(chatId, DisplayLeaderboard())
		}else {
            B.ForwardMessage(chatId, text)
        }
	}
}
