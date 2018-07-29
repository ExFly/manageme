package graph

import (
	context "context"
	"log"

	model "github.com/exfly/manageme/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Resolver struct {
	session *mgo.Session
}

func NewResolver() *Resolver {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	application := Resolver{session: session}

	return &application
}

func (r *Resolver) Mood_user(ctx context.Context, obj *model.Mood) (model.User, error) {
	return model.User{ID: obj.ID}, nil
}
func (r *Resolver) Mutation_CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {
	session := r.session.Clone()
	defer session.Close()
	c := session.DB("test").C("user")
	u := model.User{ID: bson.NewObjectId(), Username: user.Username, Password: user.Password}
	err := c.Insert(u)
	if err != nil {
		log.Fatal(err)
	}
	return &u, nil
}
func (r *Resolver) Query_user(ctx context.Context, id string) (*model.User, error) {
	return &model.User{ID: bson.ObjectId(id)}, nil
}

func (r *Resolver) User_moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	return make([]model.Mood, 0), nil
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

type moodResolver struct{ *Resolver }

func (r *moodResolver) User(ctx context.Context, obj *model.Mood) (model.User, error) {
	return model.User{}, nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {
	return r.Resolver.Mutation_CreateUser(ctx, user)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return &model.User{ID: bson.ObjectIdHex(id)}, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	return make([]model.Mood, 0), nil
}
