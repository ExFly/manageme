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
	ErrNotLogined   = errors.New("not logined")
	ErrNoPermission = errors.New("no permission")
	ErrBadRequest   = errors.New("Bad Request")
)

func getUser(ctx context.Context) *model.User {
	user, ok := ctx.Value("user").(*model.User)
	if !ok {
		return nil
	}
	return user
}

// Resolver implement Resolvers and ResolverRoot
type Resolver struct {
}

// NewResolver the Resolver's constructor
func NewResolver() *Resolver {
	application := Resolver{}
	mlog.DEBUG("%v", application)
	return &application
}

// Mood_user how to get the user in model.Mood
func (r *Resolver) Mood_user(ctx context.Context, obj *model.Mood) (model.User, error) {

	result, err := db.FindOneUser(bson.M{"_id": obj.User})

	return *result, err
}

// Mutation_CreateUser like the name
func (r *Resolver) Mutation_CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {

	u := model.User{Username: user.Username, Password: user.Password}
	db.CreateUser(&u)
	return &u, nil
}

// Mutation_CreateMood like the name
func (r *Resolver) Mutation_CreateMood(ctx context.Context, mood model.MoodInput) (*model.Mood, error) {

	user := getUser(ctx)
	if user == nil {
		return nil, ErrNotLogined
	}
	entity := model.Mood{User: user.ID, Score: mood.Score, Comment: *mood.Comment, Time: time.Now()}

	_, err := db.CreateMood(&entity)
	return &entity, err

}

func (r *Resolver) Mutation_DeleteMood(ctx context.Context, id string) (bool, error) {
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

func (r *Resolver) Query_me(ctx context.Context) (*model.User, error) {
	return getUser(ctx), nil
}

// Query_user like the name
func (r *Resolver) Query_user(ctx context.Context, id string) (*model.User, error) {

	result, err := db.FindOneUser(bson.M{"_id": id})
	return result, err
}

// Query_Users like the name
func (r *Resolver) Query_Users(ctx context.Context) ([]model.User, error) {
	result, err := db.FindUsers(bson.M{})
	return result, err
}

// User_moods like the name
func (r *Resolver) User_moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {

	result, err := db.FindMoods(bson.M{"user": obj.ID})

	return result, err
}

// Mood l
func (r *Resolver) Mood() MoodResolver {
	return &moodResolver{r}
}

// Query l
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// User l
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

// Mutation l
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// MoodResolver implementer
type moodResolver struct{ *Resolver }

func (r *moodResolver) User(ctx context.Context, obj *model.Mood) (model.User, error) {
	return r.Resolver.Mood_user(ctx, obj)
}

// MutationResolver implementer
type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {
	mlog.DEBUG("usr:%v", user)
	return r.Resolver.Mutation_CreateUser(ctx, user)
}
func (r *mutationResolver) CreateMood(ctx context.Context, mood model.MoodInput) (*model.Mood, error) {
	mlog.DEBUG("sr:%v", mood)
	return r.Resolver.Mutation_CreateMood(ctx, mood)
}

func (r *mutationResolver) DeleteMood(ctx context.Context, id string) (bool, error) {
	mlog.DEBUG("%v", id)
	return r.Resolver.Mutation_DeleteMood(ctx, id)
}

// QueryResolver implementer
type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	mlog.DEBUG("usr:%v", id)
	return r.Resolver.Query_user(ctx, id)
}
func (r *queryResolver) Users(ctx context.Context) ([]model.User, error) {
	mlog.DEBUG("")
	return r.Resolver.Query_Users(ctx)
}
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return r.Resolver.Query_me(ctx)
}

// UserResolver implementer
type userResolver struct{ *Resolver }

func (r *userResolver) Moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	mlog.DEBUG("usr:%v", *obj)
	return r.Resolver.User_moods(ctx, obj)
}
