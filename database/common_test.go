package database

import (
	"testing"
	"time"

	"github.com/exfly/manageme/model"
	"github.com/globalsign/mgo/bson"
)

func TestDataSource_CreateUser(t *testing.T) {
	var ds = NewDataSource()
	t.Run("t", func(t *testing.T) {
		user := model.User{"id", "", "username", "password", []string{}}
		id, _ := ds.CreateUser(&user)
		if id != "id" {
			t.Error(id)
		}

		user.ID = "tid"
		user.Username = "tname"
		ds.CreateUser(&user)
	})
	ds.Close()
}

func TestDataSource_FindUsers(t *testing.T) {
	var ds = NewDataSource()
	t.Run("findusers", func(t *testing.T) {
		data, _ := ds.FindUsers(bson.M{"_id": "id"})
		if data[0].ID != "id" {
			t.Error(data)
		}
	})
	ds.Close()
}

func TestDataSource_FindOneUser(t *testing.T) {
	var ds = NewDataSource()
	t.Run("findusers", func(t *testing.T) {
		data, err := ds.FindOneUser(bson.M{"_id": "id"})
		if err != nil {
			t.Error()
		}
		if data.ID != "id" {
			t.Error(data)
		}
		ds.Close()
	})
	ds.Close()
}

func TestDataSource_CreateMood(t *testing.T) {
	var ds = NewDataSource()
	t.Run("t", func(t *testing.T) {
		entity := model.Mood{ID: "moodid", User: "5b5ec1f1281645941e262c78", Score: 10, Comment: "testcomment", Time: time.Now()}
		id, _ := ds.CreateMood(&entity)
		if id != "moodid" {
			t.Error(id)
		}

		entity = model.Mood{ID: "moodid1", User: "5b5ec1f1281645941e262c78", Score: 10, Comment: "testcomment1", Time: time.Now()}
		ds.CreateMood(&entity)
		entity = model.Mood{ID: "moodid2", User: "5b5ec1f1281645941e262c78", Score: 10, Comment: "testcomment2", Time: time.Now()}
		ds.CreateMood(&entity)
	})
	ds.Close()
}

func TestDataSource_FindMoods(t *testing.T) {
	var ds = NewDataSource()
	ds.CreateMood(&model.Mood{ID: "moodid", User: "5b5eba322816458579c7c56d", Score: 10, Comment: "testcomment", Time: time.Now()})

	t.Run("findusers", func(t *testing.T) {
		data, _ := ds.FindMoods(bson.M{"_id": "moodid"})
		if data[0].ID != "moodid" {
			t.Error(data)
		}
	})
	ds.Close()
}

func TestDataSource_FindOneMood(t *testing.T) {
	var ds = NewDataSource()
	ds.CreateMood(&model.Mood{ID: "moodid", User: "5b5eba322816458579c7c56d", Score: 10, Comment: "testcomment", Time: time.Now()})

	t.Run("findusers", func(t *testing.T) {
		data, err := ds.FindOneMood(bson.M{"_id": "moodid"})
		if err != nil {
			t.Error()
		}
		if data.ID != "moodid" {
			t.Error(data)
		}
		ds.Close()
	})
	ds.Close()
}
