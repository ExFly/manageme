package database

import (
	"fmt"
	"testing"

	model "github.com/exfly/manageme/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func mongot() {
	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("user")

	err = c.Insert(model.User{ID: bson.NewObjectId().Hex(), Username: "mgoUsername"})
	if err != nil {
		panic(err)
	}

	result := model.User{}
	err = c.Find(bson.M{"username": "mgoUsername"}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.ID)
}

func Test_mongot(t *testing.T) {
	mongot()
}
