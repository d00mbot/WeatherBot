package database

import (
	"context"
	"subscription-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type userStorageService struct {
	storage mongoStorageService
}

type UserStorageService interface {
	Create(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message) error
	Update(ctx context.Context, client *mongo.Client, message *tgbotapi.Message) error
	CheckExist(ctx context.Context, client *mongo.Client, message *tgbotapi.Message) (bool, error)
	UpdateTime(ctx context.Context, client *mongo.Client, message *tgbotapi.Message, userTime string) error
	FindAll(client *mongo.Client) ([]models.User, error)
}

func NewUserStorageService(m mongoStorageService) UserStorageService {
	return &userStorageService{storage: m}
}

func (us *userStorageService) Create(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message) error {
	if err := us.storage.createUser(ctx, client, msg); err != nil {
		return err
	}

	return nil
}

func (us *userStorageService) Update(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message) error {
	if err := us.storage.updateUser(ctx, client, msg); err != nil {
		return err
	}

	return nil
}

func (us *userStorageService) CheckExist(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message) (bool, error) {
	exist, err := us.storage.checkUserExist(ctx, client, msg)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (us *userStorageService) UpdateTime(ctx context.Context, client *mongo.Client, msg *tgbotapi.Message, userTime string) error {
	if err := us.storage.updateUserTime(ctx, client, msg, userTime); err != nil {
		return err
	}

	return nil
}

func (us *userStorageService) FindAll(client *mongo.Client) ([]models.User, error) {
	users, err := us.storage.getAllUsers(client)
	if err != nil {
		return nil, err
	}

	return users, nil
}
