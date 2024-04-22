package main

import "log"

var Env Environment
var B Bot

func main() {
    log.Printf("[INFO] loading environment")

    var err error
	Env, err = LoadEnvironment()
	if err != nil {
        log.Fatalf("[ERROR] Failed to load environment: %s\n", err)
	}

    log.Printf("[INFO] starting the bot")

    B, err = NewBot()
    if err != nil {
        log.Fatalf("[ERROR] failed to start bot: %s\n", err)
    }

    log.Printf("[INFO] listening for updates...")

    msgChan := B.GetMessageChan()
    for msg := range msgChan {
        B.SendMessage(msg.ChatId, msg.Text)
    }
}
