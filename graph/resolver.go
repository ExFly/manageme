//go:generate gorunpkg github.com/99designs/gqlgen

package graph

import (
	context "context"
	"errors"
	"time"

	db "github.com/exfly/manageme/database"
	mlog "github.com/exfly/manageme/log"
	model "github.com/exfly/manageme/model"
	"github.com/globalsign/mgo/bson"
)

var (
	// ErrNotLogined like the name
	ErrNotLogined = errors.New("not logined")
	// ErrNoPermission like the name
	ErrNoPermission = errors.New("no permission")
	// ErrBadRequest like the name
	ErrBadRequest = errors.New("Bad Request")
	// ErrParamIsNil like the name
	ErrParamIsNil = errors.New("Param nil err")
)

func getUser(ctx context.Context) *model.User {
	user, ok := ctx.Value("user").(*model.User)
	if !ok {
		return nil
	}
	return user
}

// Resolver like the name
type Resolver struct{}

// NewResolver like the name
func NewResolver() *Resolver {
	application := Resolver{}
	return &application
}

// Mood like the name
func (r *Resolver) Mood() MoodResolver {
	return &moodResolver{r}
}

// Mutation like the name
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query like the name
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// User like the name
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type moodResolver struct{ *Resolver }

func (r *moodResolver) User(ctx context.Context, obj *model.Mood) (model.User, error) {
	result, err := db.FindOneUser(bson.M{"_id": obj.User})
	return *result, err
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {
	u := model.User{Username: user.Username, Password: user.Password}
	db.CreateUser(&u)
	return &u, nil
}
func (r *mutationResolver) CreateMood(ctx context.Context, mood model.MoodInput) (*model.Mood, error) {
	user := getUser(ctx)
	if user == nil {
		return nil, ErrNotLogined
	}
	entity := model.Mood{User: user.ID, Score: mood.Score, Comment: *mood.Comment, Time: time.Now()}
	_, err := db.CreateMood(&entity)
	return &entity, err
}
func (r *mutationResolver) UpdateMood(ctx context.Context, moodID string, score *int, Comment *string) (model.Mood, error) {
	if score == nil || Comment == nil || moodID == "" {
		return model.Mood{}, ErrParamIsNil
	}
	query := bson.M{}
	if *score >= 0 {
		query["score"] = score
	}
	if *Comment != "" {
		query["comment"] = *Comment
	}
	query = bson.M{"$set": query}
	mlog.DEBUG("%v", query)
	db.C(db.CollectionMood).Update(bson.M{"_id": moodID}, query)
	mood, err := db.FindOneMood(bson.M{"_id": moodID})
	return *mood, err
}
func (r *mutationResolver) DeleteMood(ctx context.Context, id string) (bool, error) {
	user := getUser(ctx)
	if user == nil {
		return false, ErrNotLogined
	}
	err := db.DeleteMood(bson.M{"_id": id, "user": user.ID})
	if err != nil {
		mlog.ERROR("%v", err)
		return false, err
	}
	return true, err
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return getUser(ctx), nil
}
func (r *queryResolver) Moods(ctx context.Context) ([]model.Mood, error) {
	user := getUser(ctx)
	result, err := db.FindMoods(bson.M{"user": user.ID})
	return result, err
}
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	result, err := db.FindOneUser(bson.M{"_id": id})
	return result, err
}
func (r *queryResolver) Users(ctx context.Context) ([]model.User, error) {
	result, err := db.FindUsers(bson.M{})
	return result, err
}

type userResolver struct{ *Resolver }

func (r *userResolver) Moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	if obj == nil {
		return nil, ErrParamIsNil
	}
	result, err := db.FindMoods(bson.M{"user": obj.ID})
	return result, err
}
