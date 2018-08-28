package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/99designs/gqlgen/handler"
	"github.com/buger/jsonparser"
	"github.com/exfly/manageme/config"
	db "github.com/exfly/manageme/database"
	"github.com/exfly/manageme/graph"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"golang.org/x/net/publicsuffix"
)

var (
	AppTestServer *httptest.Server // CAN NOT CLOSE THE AppTestServer!!!!
)

type GqlClient struct {
	client  *http.Client
	Baseurl string
}

func NewGqlClient() *GqlClient {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	return &GqlClient{
		client:  &http.Client{Jar: jar},
		Baseurl: AppTestServer.URL + "/query",
	}
}

func (gqlc *GqlClient) Get(gqlquery string) ([]byte, error) {
	u, _ := url.Parse(gqlc.Baseurl)
	q := u.Query()
	q.Set("query", gqlquery)
	u.RawQuery = q.Encode()
	// defer res.Body.Close()
	client := gqlc.client
	res, err := client.Get(u.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	return body, err
}
func (gqlc *GqlClient) Login() error {
	client := gqlc.client
	_, err := client.Get(AppTestServer.URL + "/loginas?user=username&pwd=password")
	return err
}
func (gqlc *GqlClient) Logout() error {
	panic("not implement")
}

func init() {
	config.LoadConfig("../config.yml")

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
	db.SetupDataSource()

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

	router.Handle("/", handler.Playground("manage_me", "/query"))
	router.Handle("/query", graphqlHttpHandler)
	router.Handle("/loginas", http.HandlerFunc(LoginHandler))
	router.Handle("/logout", http.HandlerFunc(LogoutHandler))
	AppTestServer = httptest.NewServer(router)
}

func Test_GqlClientWithLogin(t *testing.T) {
	client := NewGqlClient()
	client.Login()
	data, err := client.Get(`{me{id}}`)
	t.Log(string(data))
	if err != nil {
		t.Error(err)
	}
	res, err := jsonparser.GetString(data, "data", "me", "id")
	if err != nil {
		t.Error(err)
	}
	if res == "" {
		t.Errorf("error response %v", res)
	}
}
