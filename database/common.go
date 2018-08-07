package database

import (
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DataSource like the name
var session *mgo.Session

// DATABASENAME the name of the database
const DATABASENAME string = "test"

// Collection enum Collection, Prevent spelling mistakes
type Collection string

const (
	//CollectionUser 用户表
	CollectionUser Collection = "user"
	// CollectionMood 评价表
	CollectionMood Collection = "mood"
)

// NewDataSource the constructor of the data source
func SetupDataSource() {
	ses, err := mgo.Dial("localhost/test")
	if err != nil {
		mlog.DEBUG("db error %v", err)
		panic(err)
	}
	ses.SetMode(mgo.Monotonic, true)
	session = ses
	mlog.INFO("%v", session)
}

func genarateID() string {
	return bson.NewObjectId().Hex()
}

// C get the collection by the name
func C(name Collection) *mgo.Collection {
	return session.DB(DATABASENAME).C(string(name))
}

// Close close the sessiion
func Close() {
	mlog.DEBUG("")
	session.Close()
}

// CreateUser like the name
func CreateUser(entity *model.User) (string, error) {
	if entity.ID == "" {
		entity.ID = genarateID()
	}
	err := C(CollectionUser).Insert(entity)
	if err != nil {
		// todo: if insert error
		mlog.ERROR("%v : Insert error", err)
	}
	if err == nil {
		mlog.DEBUG("%v", entity.ID)
	}
	return entity.ID, nil
}

// FindUsers query the Users
func FindUsers(query bson.M) (ret []model.User, err error) {
	err = C(CollectionUser).Find(query).All(&ret)
	mlog.DEBUG("")
	return
}

// FindOneUser find one user
func FindOneUser(query bson.M) (ret *model.User, err error) {
	err = C(CollectionUser).Find(query).Limit(1).One(&ret)
	mlog.DEBUG("%v", ret)
	return
}

// CreateMood like the name
func CreateMood(entity *model.Mood) (string, error) {
	if entity.ID == "" {
		entity.ID = genarateID()
	}
	err := C(CollectionMood).Insert(entity)
	if err != nil {
		// todo: if insert error
		mlog.ERROR("%v : Insert error", err)
	}
	if err == nil {
		mlog.DEBUG("%v", entity.ID)
	}
	return entity.ID, nil
}

// FindMoods like the name
func FindMoods(query bson.M) (ret []model.Mood, err error) {
	err = C(CollectionMood).Find(query).All(&ret)
	mlog.DEBUG("")
	return
}

// FindOneMood like the name
func FindOneMood(query bson.M) (ret *model.Mood, err error) {
	err = C(CollectionMood).Find(query).Limit(1).One(&ret)
	mlog.DEBUG("%v", ret.ID)
	return
}

func DeleteMood(query bson.M) (err error) {
	err = C(CollectionMood).Remove(query)
	return
}
