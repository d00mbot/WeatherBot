package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscriber struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    int64              `bson:"chatid"`
	Latitude  float64            `bson:"latitude"`
	Longitude float64            `bson:"longitude"`
	Timezone  string             `bson:"timezone"`
	Time      string             `bson:"time"`
}
