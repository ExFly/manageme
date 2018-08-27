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

// GetGqlQuery query url, query
func GetGqlQuery(baseurl string, jar *cookiejar.Jar, gqlquery string) ([]byte, error) {
	u, _ := url.Parse(baseurl)
	q := u.Query()
	q.Set("query", gqlquery)
	u.RawQuery = q.Encode()
	client := &http.Client{}
	if jar != nil {
		client = &http.Client{
			Jar: jar,
		}
	}
	res, err := client.Get(u.String())
	// defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	return body, err
}

func LoginT(urls string, jar *cookiejar.Jar) {
	client := &http.Client{
		Jar: jar,
	}
	if _, err := client.Get(urls); err != nil {
		log.Fatal(err)
	}
}
func LogoutT(jar *cookiejar.Jar) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
}

func Test_HttptestServer(t *testing.T) {
	data, err := GetGqlQuery(AppTestServer.URL+"/query", nil, `
	{
		me{
			id
		}
	}`)
	if err != nil {
		t.Error(err)
	}
	res, err := jsonparser.GetString(data, "errors", "[0]", "message")
	if err != nil {
		t.Error(err)
	}
	target := "Not Logined"
	if res != target {
		t.Errorf("%v != %v", res, target)
	}
}

func Test_HttptestServerWithLogined(t *testing.T) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	LoginT(AppTestServer.URL+"/loginas?user=username&pwd=password", jar)
	data, err := GetGqlQuery(AppTestServer.URL+"/query", jar, `
	{
		me{
			id
		}
	}`)
	if err != nil {
		t.Error(err)
	}
	res, err := jsonparser.GetString(data, "data", "me", "id")
	if err != nil {
		t.Error(err)
	}
	target := "5b5ff11d2816453fe932f3b3"
	if res != target {
		t.Errorf("data:%v||| %v != %v", string(data), res, target)
	}
}
