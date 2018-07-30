package graph

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/vektah/gqlgen/client"
	"github.com/vektah/gqlgen/handler"
)

func TestStarwars(t *testing.T) {
	srv := httptest.NewServer(handler.GraphQL(NewExecutableSchema(NewResolver())))
	c := client.New(srv.URL)

	// create one user
	var user struct {
		CreateUser struct{ Id string }
	}
	c.Post(`mutation CreateUser {
		CreateUser(user: {sex: UNKNOWN, username: "username", password: "password"}) {
		  id
		}
	  }`, &user)
	userid := user.CreateUser.Id

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
