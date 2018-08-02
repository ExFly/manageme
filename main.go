package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/exfly/manageme/graph"
	mlog "github.com/exfly/manageme/log"
	"github.com/vektah/gqlgen/graphql"
	"github.com/vektah/gqlgen/handler"
)

func AllowOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {

	application := graph.NewResolver()

	graphqlHttpHandler := handler.GraphQL(graph.NewExecutableSchema(application),
		handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			rc := graphql.GetResolverContext(ctx)
			fmt.Println("Entered", rc.Object, rc.Field.Name)
			res, err = next(ctx)
			fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
			return res, err
		}),
	)

	http.Handle("/", handler.Playground("manage_me", "/query"))
	http.Handle("/query", AllowOriginMiddleware(graphqlHttpHandler))

	mlog.LOG("INFO", "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
