package database

import (
	"context"
	"subscription-bot/container"
	"subscription-bot/internal/api"
	"subscription-bot/internal/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbSubscribers         = "weathersubscribers"
	collectionSubscribers = "subscribers"
	defaultTime           = "08:00"
)

type mongoStorageService struct {
	container      container.Container
	weatherService api.WeatherService
}

func NewMongoStorageService(c container.Container, ws api.WeatherService) mongoStorageService {
	return mongoStorageService{container: c, weatherService: ws}
}

func NewMongoClient(c container.Container) (*mongo.Client, error) {
	logger := c.GetLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client().ApplyURI(c.GetConfig().MongoURI)

	client, err := mongo.NewClient(opt)
	if err != nil {
		logger.Errorf("error creating new mongo client: %s", err)
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		logger.Errorf("unable to create mongo client connection: %s ", err)
		return nil, err
	}
	logger.Info("Successfully connected to mongoDB")

	return client, nil
}

func openCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database(dbSubscribers).Collection(collectionSubscribers)

	return collection
}

func (ms *mongoStorageService) createUser(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message) error {
	col := openCollection(client)

	_, timezone, err := ms.weatherService.GetForecast(msg.Location.Latitude, msg.Location.Longitude)
	if err != nil {
		return err
	}

	_, err = col.InsertOne(ctx, models.User{
		ChatID:    msg.From.ID,
		Latitude:  msg.Location.Latitude,
		Longitude: msg.Location.Longitude,
		Timezone:  timezone,
		Time:      defaultTime,
	})
	if err != nil {
		ms.container.GetLogger().Errorf("error inserting new user:\n'%v'", err)
		return err
	}
	ms.container.GetLogger().Info("Subscriber successfully created")

	return nil
}

func (ms *mongoStorageService) updateUser(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message) error {
	col := openCollection(client)

	_, timezone, err := ms.weatherService.GetForecast(msg.Location.Latitude, msg.Location.Longitude)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "chatid", Value: msg.Chat.ID}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "latitude", Value: msg.Location.Latitude},
		{Key: "longitude", Value: msg.Location.Longitude},
		{Key: "timezone", Value: timezone},
	}}}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		ms.container.GetLogger().Errorf("error updating user's location data:\n'%v'", err)
		return err
	}
	ms.container.GetLogger().Info("Subscriber's location data successfully updated")

	return nil
}

func (ms *mongoStorageService) checkUserExist(ctx context.Context, client *mongo.Client, message *tgbotapi.Message) (bool, error) {
	var u models.User

	col := openCollection(client)

	filter := bson.D{{Key: "chatid", Value: message.Chat.ID}}

	userCursor := col.FindOne(ctx, filter)

	if err := userCursor.Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			ms.container.GetLogger().Infof("faild to find matching document: %s", err)
			return false, nil
		}
		ms.container.GetLogger().Errorf("error decoding document into result:\n'%v'", err)
		return false, err
	}

	if u.ChatID != message.From.ID {
		return false, nil
	}

	return true, nil
}

func (ms *mongoStorageService) updateUserTime(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message, userTime string) error {
	col := openCollection(client)

	filter := bson.D{{Key: "chatid", Value: msg.Chat.ID}}

	upd := bson.D{{Key: "$set", Value: bson.D{{Key: "time", Value: userTime + ":00"}}}}

	_, err := col.UpdateOne(ctx, filter, upd)
	if err != nil {
		ms.container.GetLogger().Errorf("error updating user's time:\n'%v'", err)
		return err
	}
	ms.container.GetLogger().Info("Subscriber's time successfully updated")

	return nil
}

func (ms *mongoStorageService) getAllUsers(client *mongo.Client) ([]models.User, error) {
	var users []models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := openCollection(client)

	usersCursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		ms.container.GetLogger().Errorf("faild to find documents in collection: %s", err)
		return nil, err
	}

	if err := usersCursor.All(ctx, &users); err != nil {
		ms.container.GetLogger().Errorf("decode documents into result:\n'%v'", err)
		return nil, err
	}

	return users, nil
}
