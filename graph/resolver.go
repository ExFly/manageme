//go:generate gorunpkg github.com/99designs/gqlgen

package graph

import (
	"context"
	"time"

	db "github.com/exfly/manageme/database"
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/globalsign/mgo/bson"
)

// ResolverFactory Config Constructor with Resolver Directives and others
func ResolverFactory() Config {
	application := Config{
		Resolvers: &Resolver{},
		Directives: DirectiveRoot{
			Logined: Logined,
			Can:     RequirePermission,
		},
	}
	return application
}

// Resolver like the name
type Resolver struct{}

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
	sex := user.Sex
	if sex == "" {
		sex = model.SexUnknown
	}
	u := model.User{Sex: sex, Username: user.Username, Password: user.Password}
	db.CreateUser(&u)
	return &u, nil
}
func (r *mutationResolver) CreateMood(ctx context.Context, mood model.MoodInput) (*model.Mood, error) {
	user := getUser(ctx)
	entity := model.Mood{User: user.ID, Score: mood.Score, Comment: *mood.Comment, Time: time.Now()}
	_, err := db.CreateMood(&entity)
	return &entity, err
}
func (r *mutationResolver) UpdateMood(ctx context.Context, moodID string, score *int, Comment *string) (model.Mood, error) {
	user := getUser(ctx)
	if score == nil || Comment == nil || moodID == "" {
		return model.Mood{}, ErrBadRequest
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
	db.C(db.CollectionMood).Update(bson.M{"_id": moodID, "user": user.ID}, query)
	mood, err := db.FindOneMood(bson.M{"_id": moodID})
	// TODO: Update requires permission
	return *mood, err
}
func (r *mutationResolver) DeleteMood(ctx context.Context, id string) (bool, error) {
	user := getUser(ctx)
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
	result, err := db.FindMoods(bson.M{"user": getUser(ctx).ID})
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
		return nil, ErrBadRequest
	}
	result, err := db.FindMoods(bson.M{"user": obj.ID})
	return result, err
}
