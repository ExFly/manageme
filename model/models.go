package model

import (
	time "time"

	"github.com/globalsign/mgo/bson"
)

type Mood struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User    string        `json:"user" bson:"user"`
	Score   int           `json:"score" bson:"score"`
	Comment *string       `json:"comment" bson:"comment"`
	Time    time.Time     `json:"time" bson:"time"`
}

type User struct {
	ID       bson.ObjectId `json:"id_" bson:"_id"`
	Sex      Sex           `json:"sex" bson:"sex"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
	Moods    []string      `json:"moods" bson:"moods"`
}
