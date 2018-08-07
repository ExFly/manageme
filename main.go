package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	db "github.com/exfly/manageme/database"
	"github.com/exfly/manageme/graph"
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/exfly/manageme/oauth"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlgen/graphql"
	"github.com/vektah/gqlgen/handler"
)

func isValidToken(token string) (*model.User, bool) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// from github.com/dhrijalva/jwt-go/hmac.go we should return a []byte
		// as we only use one single key, we just return it
		return oauth.HmacSecret, nil
	})
	if err != nil {
		return nil, false
	}
	claim, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}
	userID, ok := claim["userId"]
	if !ok {
		return nil, false
	}
	user, err := db.FindOneUser(bson.M{"_id": userID})
	if err != nil {
		return nil, false
	}
	return user, true
}

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenCookie, err := r.Cookie("jwt-token")
		if err == nil && tokenCookie != nil {
			// TODO
			user, ok := isValidToken(tokenCookie.Value)
			if ok && user != nil {
				mlog.DEBUG("set session context")
				ctx = context.WithValue(ctx, "user", user)
				ctx = context.WithValue(ctx, "userId", user.ID)
			}
		}
		mlog.DEBUG("session middleware:%v", tokenCookie)
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
	router := mux.NewRouter()
	router.Use(AllowOriginMiddleware)
	router.Use(sessionMiddleware)
	application := graph.NewResolver()
	db.SetupDataSource()

	graphqlHttpHandler := handler.GraphQL(graph.NewExecutableSchema(application),
		handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			rc := graphql.GetResolverContext(ctx)
			fmt.Println("Entered", rc.Object, rc.Field.Name)
			res, err = next(ctx)
			fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
			return res, err
		}),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// FIXME: do real check
				return true
			},
		}),
	)

	router.Handle("/", handler.Playground("manage_me", "/query"))
	router.Handle("/query", graphqlHttpHandler)
	router.Handle("/loginas", http.HandlerFunc(loginHandler))

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	mlog.DEBUG("Start server @ %s", addr)
	mlog.DEBUG("%v", srv.ListenAndServe())

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	un, oku := params["user"]
	pwd, okp := params["pwd"]
	mlog.DEBUG("%v %v", oku, okp)
	if !(oku && okp) {
		mlog.DEBUG("login faild")
		w.Write([]byte("con't login, maybe you done have param user or pwd"))
		return
	}
	username := un[0]
	password := pwd[0]
	mlog.DEBUG("%v", username)
	user, ok := db.FindOneUser(bson.M{"username": username, "password": password})
	if ok != nil || user == nil {
		mlog.DEBUG("dont have the user:%v", username)
		w.Write([]byte("dont have the user:" + username))
		return
	}
	jwtToken, err := oauth.CreateJWT(user.ID)
	if err != nil {
		mlog.DEBUG("%v", err)
		return
	}
	cookie := http.Cookie{Name: "jwt-token", Value: jwtToken, Path: "/", Expires: time.Now().Add(3600000 * time.Second)}
	http.SetCookie(w, &cookie)
	toWrite := fmt.Sprintf("user:%s    token:%s", user.ID, jwtToken)
	w.Write([]byte(toWrite))
	mlog.DEBUG("%v %v", "login as", user.Username)
}
