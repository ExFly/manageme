package graph

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"
)

func TestStarwars(t *testing.T) {

	application := ResolverFactory()
	srv := httptest.NewServer(handler.GraphQL(NewExecutableSchema(application)))
	c := client.New(srv.URL)

	// create one user
	var user struct {
		CreateUser struct{ ID string }
	}
	c.Post(`mutation CreateUser {
		CreateUser(user: {sex: UNKNOWN, username: "username", password: "password"}) {
		  id
		}
	  }`, &user)
	userid := user.CreateUser.ID

	queryString := fmt.Sprintf(`mutation CreateMood {
		CreateMood(mood: {userid: "%v", score: 5, comment: "mycommon"}) {
		  id
		  user {
			id
		  }
	  }}`, userid)
	// create one mood
	// t.Error(queryString)
	var mood struct {
		CreateMood struct {
			ID   string
			User struct {
				ID string
			}
		}
	}
	c.Post(queryString, &mood)
	if mood.CreateMood.ID == "" {
		t.Error(mood)
	}
}
