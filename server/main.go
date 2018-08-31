package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/exfly/manageme/config"
	db "github.com/exfly/manageme/database"
	"github.com/exfly/manageme/graph"
	mlog "github.com/exfly/manageme/log"
	"github.com/exfly/manageme/model"
	"github.com/exfly/manageme/oauth"
	"github.com/exfly/manageme/util"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func DataloaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := graph.NewLoader()
		ctx := context.WithValue(r.Context(), graph.LOADERKEY, loader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
	user, err := db.FindOneUser(context.Background(), util.M{"_id": userID})
	if err != nil {
		return nil, false
	}
	return user, true
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenCookie, err := r.Cookie("jwt-token")
		if err == nil && tokenCookie != nil {
			// TODO
			user, ok := isValidToken(tokenCookie.Value)
			if ok && user != nil {
				ctx = context.WithValue(ctx, "user", user)
				ctx = context.WithValue(ctx, "userId", user.ID)
			}
		}
		lenoftoken := len(tokenCookie.String())
		if lenoftoken > 15 {
			mlog.DEBUG("session middleware: %v...%v", tokenCookie.String()[0:15], tokenCookie.String()[lenoftoken-10:lenoftoken])
		} else {
			mlog.DEBUG("session middleware: %v", tokenCookie)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AllowOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
		// w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
		next.ServeHTTP(w, r)
	})
}
func BeginAndEndRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mlog.DEBUG("-----------------start--------------")
		mlog.INFO("method:%7v url: %v", r.Method, r.URL)
		next.ServeHTTP(w, r)
		mlog.DEBUG("------------------end---------------")
	})
}
func serverFactory(configName string) *http.Server {
	config.LoadConfig("../config.yml")
	config.LoadConfig("config.yml")
	util.DoInit()
	router := mux.NewRouter()
	router.Use(BeginAndEndRequest)
	// router.Use(AllowOriginMiddleware)
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            viper.GetBool("server.debug"),
	}).Handler)
	router.Use(SessionMiddleware)
	router.Use(DataloaderMiddleware)

	// application := graph.Config{Resolvers: &graph.Resolver{}}
	// db.SetupDataSource()

	graphqlHttpHandler := handler.GraphQL(graph.NewExecutableSchema(ResolverFactory()),
		// handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
		// 	rc := graphql.GetResolverContext(ctx)
		// 	mlog.DEBUG("Entered %v %v", rc.Object, rc.Field.Name)
		// 	res, err = next(ctx)
		// 	mlog.DEBUG("Left %v, %v => %v %v", rc.Object, rc.Field.Name, res, err)
		// 	return res, err
		// }),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// FIXME: do real check
				return true
			},
		}),
	)
	if viper.GetBool("server.debug") {
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

		// Manually add support for paths linked to by index page at /debug/pprof/
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
	}

	router.Handle("/", handler.Playground("manage_me", "/query"))
	router.Handle("/query", graphqlHttpHandler)
	router.Handle("/loginas", http.HandlerFunc(LoginHandler))
	router.Handle("/logout", http.HandlerFunc(LogoutHandler))
	port := os.Getenv("port")
	if port == "" {
		port = viper.GetString("server.graphql.port")
	}
	addr := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	mlog.INFO("generate server @ %s", addr)
	return &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

}
func main() {
	mlog.ERROR("%v", serverFactory("me").ListenAndServe())

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	un, oku := params["user"]
	pwd, okp := params["pwd"]
	if !(oku && okp) {
		mlog.DEBUG("login faild")
		w.Write([]byte("con't login, maybe you done have param user or pwd"))
		return
	}
	username := un[0]
	password := pwd[0]
	user, ok := db.FindOneUser(r.Context(), util.M{"username": username, "password": password})
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
	cookie := http.Cookie{Name: "jwt-token", Value: jwtToken, Path: "/", Expires: time.Now().Add(3600 * time.Second)}
	http.SetCookie(w, &cookie)
	toWrite := fmt.Sprintf("user:%s    token:%s", user.ID, jwtToken)
	w.Write([]byte(toWrite))
	mlog.DEBUG("%v %v", "login as", user.Username)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "jwt-token", Value: "", Path: "/", Expires: time.Now()}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Loginouted"))
}
