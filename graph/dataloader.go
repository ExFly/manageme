package graph

import (
	"context"
	"time"

	db "github.com/exfly/manageme/database"
	"github.com/exfly/manageme/model"
	"github.com/exfly/manageme/util"
)

// define your loader here
const LOADERKEY = "LOADER"

type Loader struct {
	User *UserLoader
	Mood *MoodLoader
}

func GetLoader(ctx context.Context) *Loader {
	loader, ok := ctx.Value(LOADERKEY).(*Loader)
	if ok {
		return loader
	} else {
		return nil
	}
}

func dupError(err error, len int) []error {
	ret := make([]error, 0, len)
	for i := 0; i < len; i++ {
		ret = append(ret, err)
	}
	return ret
}

func NewLoader() *Loader {
	loader := Loader{}
	wait := 10 * time.Millisecond
	maxBatch := 1000

	loader.User = &UserLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) (ret []*model.User, errs []error) {
			cur, err := db.UserCollection.Find(context.Background(), util.M{"_id": util.M{"$in": keys}})
			if err != nil {
				return nil, dupError(err, len(keys))
			}
			defer cur.Close(context.Background())
			for cur.Next(context.Background()) {
				var User_ *model.User
				err := cur.Decode(&User_)
				ret = append(ret, User_)
				errs = append(errs, err)
			}
			return
		},
	}
	loader.Mood = &MoodLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) (ret []*model.Mood, errs []error) {
			cur, err := db.MoodCollection.Find(context.Background(), util.M{"_id": util.M{"$in": keys}})
			if err != nil {
				return nil, dupError(err, len(keys))
			}
			defer cur.Close(context.Background())
			for cur.Next(context.Background()) {
				var Mood_ *model.Mood
				err := cur.Decode(&Mood_)
				ret = append(ret, Mood_)
				errs = append(errs, err)
			}
			return
		},
	}

	return &loader
}
