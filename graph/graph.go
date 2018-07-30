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

func (r *Resolver) Mood_user(ctx context.Context, obj *model.Mood) (model.User, error) {

	result, err := r.datasource.FindOneUser(bson.M{"_id": obj.User})

	return *result, err
}

func (r *Resolver) Mutation_CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {

	u := model.User{Username: user.Username, Password: user.Password}
	r.datasource.CreateUser(&u)
	return &u, nil
}

func (r *Resolver) Mutation_CreateMood(ctx context.Context, mood model.MoodInput) (*model.Mood, error) {

	entity := model.Mood{User: mood.Userid, Score: mood.Score, Comment: *mood.Comment, Time: time.Now()}

	_, err := r.datasource.CreateMood(&entity)
	return &entity, err

}

func (r *Resolver) Query_user(ctx context.Context, id string) (*model.User, error) {

	result, err := r.datasource.FindOneUser(bson.M{"_id": id})
	return result, err
}

func (r *Resolver) Query_Users(ctx context.Context) ([]model.User, error) {
	result, err := r.datasource.FindUsers(bson.M{})
	return result, err
}

func (r *Resolver) User_moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {

	result, err := r.datasource.FindMoods(bson.M{"user": obj.ID})

	return result, err
}

func (r *Resolver) Mood() MoodResolver {
	return &moodResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
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
	return r.Resolver.Mutation_CreateUser(ctx, user)
}
func (r *mutationResolver) CreateMood(ctx context.Context, mood model.MoodInput) (*model.Mood, error) {
	return r.Resolver.Mutation_CreateMood(ctx, mood)
}

// QueryResolver implementer
type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.Resolver.Query_user(ctx, id)
}
func (r *queryResolver) Users(ctx context.Context) ([]model.User, error) {
	return r.Resolver.Query_Users(ctx)
}

// UserResolver implementer
type userResolver struct{ *Resolver }

func (r *userResolver) Moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	return r.Resolver.User_moods(ctx, obj)
}
