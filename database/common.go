package database

import (
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DataSource like the name
type DataSource struct {
	session *mgo.Session
}

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
func NewDataSource() *DataSource {
	session, err := mgo.Dial("localhost")
	if err != nil {
		mlog.DEBUG("db error %v", err)
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	mlog.INFO("%v", session)
	return &DataSource{session}
}

func genarateID() string {
	return bson.NewObjectId().Hex()
}

// C get the collection by the name
func (d *DataSource) C(name Collection) *mgo.Collection {
	return d.session.DB(DATABASENAME).C(string(name))
}

// Close close the sessiion
func (d *DataSource) Close() {
	mlog.DEBUG("")
	d.session.Close()
}

// CreateUser like the name
func (d *DataSource) CreateUser(entity *model.User) (string, error) {
	if entity.ID == "" {
		entity.ID = genarateID()
	}
	err := d.C(CollectionUser).Insert(entity)
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
func (d *DataSource) FindUsers(query bson.M) (ret []model.User, err error) {
	err = d.C(CollectionUser).Find(query).All(&ret)
	mlog.DEBUG("")
	return
}

// FindOneUser find one user
func (d *DataSource) FindOneUser(query bson.M) (ret *model.User, err error) {
	err = d.C(CollectionUser).Find(query).Limit(1).One(&ret)
	mlog.DEBUG("%v", ret.ID)
	return
}

// CreateMood like the name
func (d *DataSource) CreateMood(entity *model.Mood) (string, error) {
	if entity.ID == "" {
		entity.ID = genarateID()
	}
	err := d.C(CollectionMood).Insert(entity)
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
func (d *DataSource) FindMoods(query bson.M) (ret []model.Mood, err error) {
	err = d.C(CollectionMood).Find(query).All(&ret)
	mlog.DEBUG("")
	return
}

// FindOneMood like the name
func (d *DataSource) FindOneMood(query bson.M) (ret *model.Mood, err error) {
	err = d.C(CollectionMood).Find(query).Limit(1).One(&ret)
	mlog.DEBUG("%v", ret.ID)
	return
}

func (d *DataSource) DeleteMood(query bson.M) (err error) {
	err = d.C(CollectionMood).Remove(query)
	return
}
