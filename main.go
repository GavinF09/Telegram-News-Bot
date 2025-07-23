package main

import (
	"context"
	"log"
	"os"

	"tele-news-bot/commands"
	"tele-news-bot/connect_mongodb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	// bot info
	bot      *tgbotapi.BotAPI
	dbClient *mongo.Client
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

	// connect to db
	certFile := "./certs/cert.crt"
	keyFile := "./certs/key.pem"
	dbURI := os.Getenv("MONGODB_URI")
	dbClient = connect_mongodb.ConnectMongoDB(certFile, keyFile, dbURI)
	// coll := dbClient.Database("telegram-bot").Collection("news-users")
	defer func() {
		if err := dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	go commands.SendReminderDB(bot, dbClient)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// make a channel to listen for updates
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		receiveUpdates(update)
	}

}

func receiveUpdates(update tgbotapi.Update) {
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
	case "reminder":
		status, err := commands.ToggleReminderDB(message.Chat.ID, dbClient)
		if err != nil {
			log.Println(err)
			return err
		}
		if status == 0 {
			msg.Text = "Stopped reminder"
		}
		if status == 1 {
			msg.Text = "Enabled reminder"
		}

	default:
		msg.Text = "Unknown command"
	}

	return nil
}
