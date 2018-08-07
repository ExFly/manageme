package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	db "github.com/exfly/manageme/database"
	"github.com/exfly/manageme/graph"
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/exfly/manageme/oauth"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlgen/graphql"
	"github.com/vektah/gqlgen/handler"
)

func isValidToken(token string, payload string) bool {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// from github.com/dhrijalva/jwt-go/hmac.go we should return a []byte
		// as we only use one single key, we just return it
		return oauth.HmacSecret, nil
	})
	if err != nil {
		return false
	}
	claim, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	userID, ok := claim["userId"]
	if !ok {
		return false
	}
	if userID == payload {
		return true
	}
	return false
}

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenCookie, err := r.Cookie("jwt-token")
		if err == nil && tokenCookie != nil {
			// TODO
			user, ok := &model.User{}, true //isValidToken(tokenCookie.Value)
			if ok && user != nil {
				ctx = context.WithValue(ctx, "user", user)
				ctx = context.WithValue(ctx, "userId", user.ID)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AllowOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {

	application := graph.NewResolver()
	db.SetupDataSource()

	graphqlHttpHandler := handler.GraphQL(graph.NewExecutableSchema(application),
		handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			rc := graphql.GetResolverContext(ctx)
			fmt.Println("Entered", rc.Object, rc.Field.Name)
			res, err = next(ctx)
			fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
			return res, err
		}), handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// FIXME: do real check
				return true
			},
		}),
	)

	http.Handle("/", handler.Playground("manage_me", "/query"))
	http.Handle("/query", AllowOriginMiddleware(graphqlHttpHandler))

	mlog.LOG("INFO", "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

/*
handler.WebsocketUpgrader(websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// FIXME: do real check
		return true
	},
})
*/
