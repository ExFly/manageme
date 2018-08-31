package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	db "github.com/exfly/manageme/database"
	"github.com/exfly/manageme/model"
	"github.com/mongodb/mongo-go-driver/bson"
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
	rand.Seed(time.Now().Unix())
	os.Setenv("mongo", "mongodb://localhost:27017/manageme_test")
	os.Setenv("port", string(10000+rand.Intn(50000)))
	log.Println(os.Getenv("port"))
	svr := serverFactory("me_test")
	AppTestServer = httptest.NewServer(svr.Handler)
	_, err := db.Client.Database("manageme_test").RunCommand(
		context.Background(),
		bson.NewDocument(bson.EC.Int32("dropDatabase", 1)),
	)
	log.Print(err)
	log.Print(db.UserCollection)
	db.UserCollection.InsertOne(context.Background(), model.User{ID: db.GenarateID(), Username: "username", Password: "password"})
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

func Test_GqlMood(t *testing.T) {
	client := NewGqlClient()
	client.Login()
	data, err := client.Get(`mutation mood{
		CreateMood(mood:{score:10,comment:"newcoment"}){comment}}`)
	t.Log(string(data))
	if err != nil {
		t.Error(err)
	}
	res, err := jsonparser.GetString(data, "data", "CreateMood", "comment")
	if err != nil {
		t.Error(err)
	}
	if res != "newcoment" {
		t.Errorf("error response %v", res)
	}
}
