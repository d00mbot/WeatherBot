package database

import (
	"context"
	"time"

	api "subscription-bot/pkg/api"
	domain "subscription-bot/pkg/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoClient(mongoURI string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client().ApplyURI(mongoURI)

	client, err := mongo.NewClient(opt)
	if err != nil {
		log.Errorf("error creating new mongo client: %s", err)
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Errorf("unable to create mongo client connection: %s ", err)
		return nil, err
	}
	log.Infof("Successfully connected to mongoDB")

	return client, nil
}

func openCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database("weathersubscribers").Collection("subscribers")

	return collection
}

func InsertSubscriber(ctx context.Context, client *mongo.Client, message *tgbotapi.Message, weatherToken string) error {
	col := openCollection(client)

	timezone, err := api.GetUserTimezone(message, weatherToken)
	if err != nil {
		log.Errorf("error inserting user's timezone: %s", err)
		return err
	}

	filter := bson.D{{Key: "chatid", Value: message.Chat.ID}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "latitude", Value: message.Location.Latitude},
		{Key: "longitude", Value: message.Location.Longitude},
		{Key: "timezone", Value: timezone},
		{Key: "time", Value: "08:00"},
	}}}

	opt := options.Update().SetUpsert(true)

	if _, err = col.UpdateOne(ctx, filter, update, opt); err != nil {
		log.Errorf("error inserting subscriber's data: %s", err)
		return err
	}
	log.Info("Subscriber's data successfully saved")

	return nil
}

func CheckUserExist(ctx context.Context, client *mongo.Client, message *tgbotapi.Message) (bool, error) {
	sub := domain.Subscriber{}

	col := openCollection(client)

	filter := bson.D{{Key: "chatid", Value: message.Chat.ID}}

	subCursor := col.FindOne(ctx, filter)

	if err := subCursor.Decode(&sub); err != nil {
		log.Errorf("error decoding document into result: %s", err)
		return false, nil
	}

	if sub.ChatID != message.From.ID {
		return false, nil
	}

	return true, nil
}

func UpdateUserTime(ctx context.Context, client *mongo.Client, message *tgbotapi.Message, userTime string) error {
	col := openCollection(client)

	filter := bson.D{{Key: "chatid", Value: message.Chat.ID}}

	upd := bson.D{{Key: "$set", Value: bson.D{{Key: "time", Value: userTime + ":00"}}}}

	if _, err := col.UpdateOne(ctx, filter, upd); err != nil {
		log.Errorf("error inseting user's time: %s", err)
		return err
	}
	log.Info("Subscriber's time successfully updated")

	return nil
}

func GetAllSubscribers(client *mongo.Client) ([]domain.Subscriber, error) {
	subs := []domain.Subscriber{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := openCollection(client)

	subCursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		log.Errorf("error trying to find matching documents from collection: %s", err)
		return nil, err
	}

	if err = subCursor.All(ctx, &subs); err != nil {
		log.Errorf("error decoding documents into result: %s", err)
		return nil, err
	}

	return subs, nil
}
