package main

import (
	"errors"
	"log"
	"os"

	// "time"

	// "github.com/mmcdole/gofeed"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	// bot info
	bot *tgbotapi.BotAPI

	// fake database
	// echoUsers = make([]int, 0)
	// the val of the map is not important
	echoUsers = make(map[int64]bool)
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// make a channel to listen for updates
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// MessageConfig to send message to user, not to be confused with Update.Message
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		// use a pointer for functions
		pMsg := &msg

		if update.Message.IsCommand() {
			handleCommand(update.Message, pMsg)
		} else {
			// check if echo is enabled for user
			_, ok := echoUsers[update.Message.Chat.ID]
			if ok {
				msg.Text = update.Message.Text
			}
		}
		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}

}

func handleCommand(message *tgbotapi.Message, msg *tgbotapi.MessageConfig) error {
	log.Printf("Received command: %s", message.Command())

	switch message.Command() {
	// START global commands
	case "start":
		msg.Text = "Welcome to my first test bot"
	case "help":
		msg.Text = "This is a testing bot work in progress"
	case "settings":
		msg.Text = "This is a testing bot work in progress"
	// END  global commands
	case "echo":
		status, err := toggleEcho(message.Chat.ID)
		if err != nil {
			msg.Text = "Something failed"
			return nil
		}
		if status == 0 {
			msg.Text = "Stopped echo"
		}
		if status == 1 {
			msg.Text = "Enabled echo"
		}

	default:
		msg.Text = "Unknown command"
	}

	return nil
}

func toggleEcho(userid int64) (echoStatus int, err error) {
	// see if user already enabled echo
	_, ok := echoUsers[userid]
	// user is already registered
	if ok {
		delete(echoUsers, userid)
		return 0, nil
	} else {
		echoUsers[userid] = true
		return 1, nil
	}

	return -1, errors.New("something failed")

}
