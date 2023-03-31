package database

import (
	"context"
	"time"

	api "subscription-bot/pkg/api"
	"subscription-bot/pkg/config"
	"subscription-bot/pkg/container"
	"subscription-bot/pkg/domain/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "weathersubscribers"
	collectionName = "subscribers"
)

type MongoStorageService struct {
	container container.BotContainer
	request   api.RequestWeatherService
}

func NewMongoStorageService(c container.BotContainer, r api.RequestWeatherService) MongoStorageService {
	return MongoStorageService{container: c, request: r}
}

func (m *MongoStorageService) initMongoClient(cfg *config.Config) (*mongo.Client, error) {
	logger := m.container.GetLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.NewClient(opt)
	if err != nil {
		logger.Errorf("error creating new mongo client: %s", err)
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		logger.Errorf("unable to create mongo client connection: %s ", err)
		return nil, err
	}
	logger.Info("Successfully connected to mongoDB")

	return client, nil
}

func openCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)

	return collection
}

func (m *MongoStorageService) createUser(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message, cfg *config.Config) error {
	logger := m.container.GetLogger()

	col := openCollection(client)

	timezone, err := m.request.RequestTimezone(msg, cfg)
	if err != nil {
		logger.Errorf("faild to get user's timezone: %s", err)
		return err
	}

	user := models.NewUser(msg.From.ID, msg.Location.Latitude, msg.Location.Longitude, timezone)

	if _, err := col.InsertOne(ctx, user); err != nil {
		logger.Errorf("faild to insert new user: %s", err)
	}

	return nil
}

func (m *MongoStorageService) updateUser(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message, cfg *config.Config) error {
	logger := m.container.GetLogger()

	col := openCollection(client)

	timezone, err := m.request.RequestTimezone(msg, cfg)
	if err != nil {
		logger.Errorf("error inserting user's timezone: %s", err)
		return err
	}

	filter := bson.D{{Key: "chatid", Value: msg.Chat.ID}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "latitude", Value: msg.Location.Latitude},
		{Key: "longitude", Value: msg.Location.Longitude},
		{Key: "timezone", Value: timezone},
	}}}

	if _, err = col.UpdateOne(ctx, filter, update); err != nil {
		logger.Errorf("faild to update user's location data: %s", err)
		return err
	}
	logger.Info("Subscriber's location data successfully updated")

	return nil
}

func (m *MongoStorageService) checkUserExist(ctx context.Context, client *mongo.Client, message *tgbotapi.Message) (bool, error) {
	logger := m.container.GetLogger()

	var u models.User

	col := openCollection(client)

	filter := bson.D{{Key: "chatid", Value: message.Chat.ID}}

	userCursor := col.FindOne(ctx, filter)

	if err := userCursor.Decode(&u); err != nil {
		logger.Errorf("error decoding document into result: %s", err)
		return false, nil
	}

	if u.ChatID != message.From.ID {
		return false, nil
	}

	return true, nil
}

func (m *MongoStorageService) updateUserTime(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message, userTime string) error {
	logger := m.container.GetLogger()

	col := openCollection(client)

	filter := bson.D{{Key: "chatid", Value: msg.Chat.ID}}

	upd := bson.D{{Key: "$set", Value: bson.D{{Key: "time", Value: userTime + ":00"}}}}

	if _, err := col.UpdateOne(ctx, filter, upd); err != nil {
		logger.Errorf("faild to update user's time: %s", err)
		return err
	}
	logger.Info("Subscriber's time successfully updated")

	return nil
}

func (m *MongoStorageService) getAllUsers(client *mongo.Client) ([]models.User, error) {
	logger := m.container.GetLogger()

	var users []models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := openCollection(client)

	usersCursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		logger.Errorf("faild to find documents in collection: %s", err)
		return nil, err
	}

	if err = usersCursor.All(ctx, &users); err != nil {
		logger.Errorf("error decoding documents into result: %s", err)
		return nil, err
	}

	return users, nil
}
