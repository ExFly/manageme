package graph

import (
	context "context"

	model "github.com/exfly/manageme/model"
	"github.com/globalsign/mgo/bson"
)

type Resolver struct{}

func (r *Resolver) Mood_user(ctx context.Context, obj *model.Mood) (model.User, error) {
	return model.User{ID: obj.ID}, nil
}

func (r *Resolver) Query_user(ctx context.Context, id string) (*model.User, error) {
	return &model.User{ID: bson.ObjectId(id)}, nil
}

func (r *Resolver) User_moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	return make([]model.Mood, 0), nil
}

func (r *Resolver) Mutation_create(ctx context.Context, name *string) (*string, error) {
	t := "create test string"
	return &t, nil
}

type App struct{
	
}

func (r *App) Mood() MoodResolver {
	return &moodResolver{r}
}
func (r *App) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *App) User() UserResolver {
	return &userResolver{r}
}
func (r *App) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type moodResolver struct{ *App }

func (r *moodResolver) User(ctx context.Context, obj *model.Mood) (model.User, error) {
	return model.User{}, nil
}

type mutationResolver struct{ *App }

func (r *mutationResolver) Create(ctx context.Context, name *string) (*string, error) {
	t := "create test"
	return &t, nil
}

type queryResolver struct{ *App }

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return &model.User{ID: bson.ObjectIdHex(id)}, nil
}

type userResolver struct{ *App }

func (r *userResolver) Moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	return make([]model.Mood, 0), nil
}
