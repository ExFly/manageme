package graph

import (
	"context"
	"time"

	db "github.com/exfly/manageme/database"
	"github.com/exfly/manageme/model"
	"github.com/globalsign/mgo/bson"
)

// define your loader here
const LOADERKEY = "LOADER"

type Loader struct {
	User UserLoader
	Mood MoodLoader
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
	l := Loader{}
	wait := 10 * time.Millisecond
	maxBatch := 1000

	l.User = UserLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) ([]*model.User, []error) {
			entities, err := db.FindUsers(bson.M{"_id": bson.M{"$in": keys}})
			if err != nil {
				return nil, dupError(err, len(keys))
			}
			ret := make([]*model.User, 0, len(keys))
			cache := make(map[string]*model.User)
			for _, entity := range entities {
				item := entity
				cache[entity.ID] = &item
			}
			for _, id := range keys {
				ret = append(ret, cache[id])
			}
			return ret, nil
		},
	}
	l.Mood = MoodLoader{
		wait:     wait,
		maxBatch: maxBatch,
		fetch: func(keys []string) ([]*model.Mood, []error) {
			entities, err := db.FindMoods(bson.M{"_id": bson.M{"$in": keys}})
			if err != nil {
				return nil, dupError(err, len(keys))
			}
			ret := make([]*model.Mood, 0, len(keys))
			cache := make(map[string]*model.Mood)
			for _, entity := range entities {
				item := entity
				cache[entity.ID] = &item
			}
			for _, id := range keys {
				ret = append(ret, cache[id])
			}
			return ret, nil
		},
	}

	return &l
}
