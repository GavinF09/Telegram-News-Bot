package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	// set a single user to remind
	remindUser int64 = 0
)

type userToRemind struct {
	ChatID int64 `bson:"chat_id"`
	Remind bool  `bson:"reminder"`
}

// this file is a test to send a message to a user every 1 min - to test out sending messages that is not a reply or command

func ToggleReminder(userid int64) int {
	// see if user already enabled echo
	// user is already registered
	if userid == remindUser {
		remindUser = 0
		return 0
	} else {
		remindUser = userid
		return 1
	}
}

func SendReminder(bot *tgbotapi.BotAPI) (int, error) {
	tick := time.Tick(1 * time.Minute)
	for range tick {
		log.Println("Running reminder")
		if remindUser == 0 {
			return 0, nil
		}
		msg := tgbotapi.NewMessage(remindUser, "")
		msg.Text = "Reminder"
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return -1, err
		}
	}
	return 1, nil
}

func ToggleReminderDB(userid int64, dbClient *mongo.Client) (int, error) {
	// see if user already enabled echo
	coll := dbClient.Database("telegram-bot").Collection("news-users")
	filter := bson.D{{Key: "chat_id", Value: userid}}
	var result userToRemind
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Println(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("Not found: %d\n", userid)
			// user not registered in any way - insert them in
			user := userToRemind{ChatID: userid, Remind: true}
			result, err := coll.InsertOne(context.TODO(), user)
			fmt.Println(result.InsertedID)
			if err != nil {
				return -1, err
			}
			return 1, nil
		}
		panic(err)
	}
	// user is already registered - toggle result
	var newStatus int
	if result.Remind {
		// will be disabled
		newStatus = 0
	} else {
		// will be enabled
		newStatus = 1
	}
	// update user and change remind value
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "reminder", Value: !result.Remind}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return -1, err
	}

	return newStatus, nil
}

func SendReminderDB(bot *tgbotapi.BotAPI, dbClient *mongo.Client) (int, error) {
	tick := time.Tick(1 * time.Minute)
	for range tick {
		log.Println("Running reminder")
		// get back all users in database
		coll := dbClient.Database("telegram-bot").Collection("news-users")
		filter := bson.D{{}}
		// retrieve documents
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		// unpack cursor into slice
		var results []userToRemind
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		for _, user := range results {
			fmt.Println(user)
			if user.ChatID == 0 {
				continue
			}
			if user.Remind {
				msg := tgbotapi.NewMessage(user.ChatID, "")
				msg.Text = "Reminder"
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					return -1, err
				}
			}
		}
	}
	return 1, nil
}
