package main

import (
	"log"
	"os"

	// "time"

	// "github.com/mmcdole/gofeed"

	"tele-news-bot/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	// bot info
	bot *tgbotapi.BotAPI
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
		} else if commands.EchoMessage(update.Message.Chat.ID) {
			msg.Text = update.Message.Text
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
		status := commands.ToggleEcho(message.Chat.ID)
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
