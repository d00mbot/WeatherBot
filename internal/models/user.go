package models

type User struct {
	ChatID    int64   `bson:"chatid"`
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	Timezone  string  `bson:"timezone"`
	Time      string  `bson:"time"`
}
