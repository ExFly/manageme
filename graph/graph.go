package graph

import (
	context "context"
	"log"

	mlog "github.com/exfly/manageme/log"
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
	mlog.DEBUG("", application)
	return &application
}

func (r *Resolver) Mood_user(ctx context.Context, obj *model.Mood) (model.User, error) {
	mlog.DEBUG("Mood_usr", obj)
	return model.User{ID: obj.ID}, nil
}
func (r *Resolver) Mutation_CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {
	mlog.DEBUG("start")
	session := r.session.Clone()
	mlog.DEBUG("session clone")
	defer session.Close()
	c := session.DB("test").C("user")
	u := model.User{ID: bson.NewObjectId(), Username: user.Username, Password: user.Password}
	err := c.Insert(u)
	if err != nil {
		log.Fatal(err)
	}
	mlog.DEBUG("end")
	return &u, nil
}
func (r *Resolver) Query_user(ctx context.Context, id string) (*model.User, error) {
	mlog.DEBUG(id)
	return &model.User{ID: bson.ObjectId(id)}, nil
}

func (r *Resolver) User_moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	mlog.DEBUG("", obj)
	return make([]model.Mood, 0), nil
}

func (r *Resolver) Mood() MoodResolver {
	mlog.DEBUG("")
	return &moodResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	mlog.DEBUG("")
	return &queryResolver{r}
}
func (r *Resolver) User() UserResolver {
	mlog.DEBUG("")
	return &userResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	mlog.DEBUG("")
	return &mutationResolver{r}
}

type moodResolver struct{ *Resolver }

func (r *moodResolver) User(ctx context.Context, obj *model.Mood) (model.User, error) {
	mlog.DEBUG("", obj)
	return model.User{}, nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, user model.UserInput) (*model.User, error) {
	mlog.DEBUG("", user)
	return r.Resolver.Mutation_CreateUser(ctx, user)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	mlog.DEBUG(id)
	return &model.User{ID: bson.ObjectIdHex(id)}, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) Moods(ctx context.Context, obj *model.User) ([]model.Mood, error) {
	mlog.DEBUG("", obj)
	return make([]model.Mood, 0), nil
}
