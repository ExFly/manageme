package model

import (
	time "time"
)

type Mood struct {
	ID      string    `json:"id" bson:"_id"`
	User    string    `json:"user" bson:"user"`
	Score   int       `json:"score" bson:"score"`
	Comment *string   `json:"comment" bson:"comment"`
	Time    time.Time `json:"time" bson:"time"`
}

type User struct {
	ID       string   `json:"id" bson:"_id"`
	Sex      Sex      `json:"sex" bson:"sex"`
	Username string   `json:"username" bson:"username"`
	Password string   `json:"password" bson:"password"`
	Moods    []string `json:"moods" bson:"moods"`
}
