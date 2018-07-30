package database

import (
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type DataSource struct {
	session *mgo.Session
}

var DATABASENAME string = "test"

type Collection string

const (
	//U serCollection 用户表
	CollectionUser Collection = "user"
	// MoodCollection 评价表
	CollectionMood Collection = "mood"
)

func NewDataSource() *DataSource {
	session, err := mgo.Dial("localhost")
	if err != nil {
		mlog.DEBUG("db error %v", err)
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return &DataSource{session}
}

func genarateID() string {
	return bson.NewObjectId().Hex()
}

func (d *DataSource) C(name Collection) *mgo.Collection {
	return d.session.DB(DATABASENAME).C(string(name))
}

func (d *DataSource) Close() {
	d.session.Close()
}

func (d *DataSource) CreateUser(entity *model.User) (string, error) {
	if entity.ID == "" {
		entity.ID = genarateID()
	}
	err := d.C(CollectionUser).Insert(entity)
	if err != nil {
		// todo: if insert error
		mlog.ERROR("%v : Insert error", err)
	}
	return entity.ID, nil
}

func (d *DataSource) FindUsers(query bson.M) (ret []model.User, err error) {
	err = d.C(CollectionUser).Find(query).All(&ret)
	return
}

func (d *DataSource) FindOneUser(query bson.M) (ret *model.User, err error) {
	err = d.C(CollectionUser).Find(query).Limit(1).One(&ret)
	return
}
