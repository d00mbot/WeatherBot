package models

const defaultTime = "08:00"

type User struct {
	ChatID    int64   `bson:"chatid"`
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	Timezone  string  `bson:"timezone"`
	Time      string  `bson:"time"`
}

func NewUser(chatID int64, lat float64, lon float64, timezone string) *User {
	return &User{
		ChatID:    chatID,
		Latitude:  lat,
		Longitude: lon,
		Timezone:  timezone,
		Time:      defaultTime,
	}
}
