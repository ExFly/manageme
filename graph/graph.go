package graph

import (
	context "context"
	"time"

	"github.com/exfly/manageme/database"
	mlog "github.com/exfly/manageme/log"
	model "github.com/exfly/manageme/model"
	"github.com/globalsign/mgo/bson"
)

// Resolver implement Resolvers and ResolverRoot
type Resolver struct {
	datasource *database.DataSource
}

// NewResolver the Resolver's constructor
func NewResolver() *Resolver {
	datasource := database.NewDataSource()
	application := Resolver{datasource: datasource}
	mlog.DEBUG("", application)
	return &application
}

// Mood_user how to get the user in model.Mood
func (r *Resolver) Mood_user(ctx context.Context, obj *model.Mood) (model.User, error) {

	result, err := r.datasource.FindOneUser(bson.M{"_id": obj.User})

	return *result, err
}

// Mutation_CreateUser like the name
func (r *Resolver) Mutation_CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {

	u := model.User{Username: user.Username, Password: user.Password}
	r.datasource.CreateUser(&u)
	return &u, nil
}

// Mutation_CreateMood like the name
func (r *Resolver) Mutation_CreateMood(ctx context.Context, mood model.MoodInput) (*model.Mood, error) {

	entity := model.Mood{User: mood.Userid, Score: mood.Score, Comment: *mood.Comment, Time: time.Now()}

	_, err := r.datasource.CreateMood(&entity)
	return &entity, err

}

func (r *Resolver) Mutation_DeleteMood(ctx context.Context, id string) (bool, error) {
	err := r.datasource.DeleteMood(bson.M{"_id": id})
	if err != nil {
		mlog.ERROR("%v", err)
		return false, err
	}
	return true, err
}

// Query_user like the name
func (r *Resolver) Query_user(ctx context.Context, id string) (*model.User, error) {

	result, err := r.datasource.FindOneUser(bson.M{"_id": id})
	return result, err
}

// Query_Users like the name
func (r *Resolver) Query_Users(ctx context.Context) ([]model.User, error) {
	result, err := r.datasource.FindUsers(bson.M{})
	return result, err
}

// User_moods like the name
func (r *Resolver) User_moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {

	result, err := r.datasource.FindMoods(bson.M{"user": obj.ID})

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

// UserResolver implementer
type userResolver struct{ *Resolver }

func (r *userResolver) Moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	mlog.DEBUG("usr:%v", *obj)
	return r.Resolver.User_moods(ctx, obj)
}
