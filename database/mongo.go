package database

import (
	"context"
	"log"
	"math/rand"
	"net/url"
	"os"

	"git.in.chaitin.com/babysitter/man-month/server/utils"
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/exfly/manageme/util"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/spf13/viper"
)

var (
	Client          *mongo.Client
	defaultDatabase string
	UserCollection  *mongo.Collection
	MoodCollection  *mongo.Collection
)

func init() {
	util.RegisterInitFunction("SetupDB", func() {
		mongourl := os.Getenv("mongo")
		if mongourl == "" {
			mongourl = viper.GetString("db.url")
		}
		if mongourl == "" {
			mlog.ERROR("config not load")
			return
		}
		err := Client.Connect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("connected to database")
		u, _ := url.Parse(mongourl)
		defaultDatabase = u.Path[1:]
		mlog.INFO("connect: %v", mongourl)

		UserCollection = C("user")
		MoodCollection = C("mood")
	}, 1)
}

// UnmistakebleChars copied from https://github.com/meteor/meteor/blob/24865b28a0689de8b4949fb69ea1f95da647cd7a/packages/random/random.js#L88
const UnmistakebleChars = "23456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz"

func genarateID() string {
	// FIXME: port the Random.id()
	// return bson.NewObjectId().Hex()
	buf := make([]byte, 17)
	for i := 0; i < 17; i++ {
		buf[i] = UnmistakebleChars[rand.Intn(len(UnmistakebleChars))]
	}
	return string(buf)
}

// C get the collection by the name
// If the session is empty, it means that the database initialization failed and should fail quickly.
func C(name string) *mongo.Collection {
	return Client.Database(defaultDatabase).Collection(name)
}
func CreateUser(entity *model.User) (string, error) {
	if entity.ID == "" {
		entity.ID = genarateID()
	}
	_, err := UserCollection.InsertOne(context.Background(), entity)
	if err != nil {
		// todo: if insert error
		mlog.ERROR("%v : Insert error", err)
	}
	if err == nil {
		mlog.DEBUG("%v", entity.ID)
	}
	return entity.ID, nil
}

func CreateMood(entity *model.Mood) (string, error) {
	if entity.ID == "" {
		entity.ID = genarateID()
	}
	_, err := MoodCollection.InsertOne(context.Background(), entity)
	if err != nil {
		// todo: if insert error
		mlog.ERROR("%v : Insert error", err)
	}
	if err == nil {
		mlog.DEBUG("%v", entity.ID)
	}
	return entity.ID, nil
}
func DeleteMood(query util.M) (err error) {
	_, err = MoodCollection.DeleteOne(context.Background(), query)
	return
}

func FindUser(ctx context.Context, query interface{}, opts ...findopt.Find) (ret []model.User, err error) {
	cur, err := UserCollection.Find(ctx, query, opts...)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var _User model.User
		err = cur.Decode(&_User)
		if err != nil {
			return
		}
		ret = append(ret, _User)
	}
	return
}
func FindOneUser(ctx context.Context, query interface{}, opts ...findopt.One) (ret *model.User, err error) {
	cur := UserCollection.FindOne(ctx, query)
	err = cur.Decode(&ret)
	var rett model.User
	err = cur.Decode(&rett)
	ret = &rett
	return
}

func FindOneUserById(ctx context.Context, id interface{}) (ret *model.User, err error) {
	return FindOneUser(ctx, utils.M{"_id": id})
}

func FindMood(ctx context.Context, query interface{}, opts ...findopt.Find) (ret []model.Mood, err error) {
	cur, err := MoodCollection.Find(ctx, query, opts...)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var _Mood model.Mood
		err = cur.Decode(&_Mood)
		if err != nil {
			return
		}
		ret = append(ret, _Mood)
	}
	return
}
func FindOneMood(ctx context.Context, query interface{}, opts ...findopt.One) (ret *model.Mood, err error) {
	cur := MoodCollection.FindOne(ctx, query)
	err = cur.Decode(&ret)
	var rett model.Mood
	err = cur.Decode(&rett)
	ret = &rett
	return
}
func FindOneMoodById(ctx context.Context, id interface{}) (ret *model.Mood, err error) {
	return FindOneMood(ctx, utils.M{"_id": id})
}
