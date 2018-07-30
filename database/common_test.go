package database

import (
	"testing"

	"github.com/exfly/manageme/model"
	"github.com/globalsign/mgo/bson"
)

func TestDataSource_CreateUser(t *testing.T) {
	var ds = NewDataSource()
	t.Run("t", func(t *testing.T) {
		user := model.User{"id", "", "username", "password", []string{}}
		id := ds.CreateUser(&user)
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
