package model

import time "time"

type Mood struct {
	ID      string    `json:"id"`
	User    string    `json:"user"`
	Score   int       `json:"score"`
	Comment *string   `json:"comment"`
	Time    time.Time `json:"time"`
}

type User struct {
	ID       string   `json:"id"`
	Sex      Sex      `json:"sex"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Moods    []string `json:"moods"`
}
