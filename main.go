package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/exfly/manageme/graph"
	"github.com/vektah/gqlgen/graphql"
	"github.com/vektah/gqlgen/handler"
)

func main() {

	// http.Handle("/", handler.Playground("manage_me", "/query"))
	// http.Handle("/query", handler.GraphQL(graph.MakeExecutableSchema(&graph.App{}),
	// 	handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	// 		rc := graphql.GetResolverContext(ctx)
	// 		fmt.Println("Entered", rc.Object, rc.Field.Name)
	// 		res, err = next(ctx)
	// 		fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
	// 		return res, err
	// 	}),
	// ))

	// log.Fatal(http.ListenAndServe(":8080", nil))

	application := graph.NewResolver()

	http.Handle("/", handler.Playground("manage_me", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(application),
		handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
			rc := graphql.GetResolverContext(ctx)
			fmt.Println("Entered", rc.Object, rc.Field.Name)
			res, err = next(ctx)
			fmt.Println("Left", rc.Object, rc.Field.Name, "=>", res, err)
			return res, err
		}),
	))
	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// http.Handle("/", handler.Playground("Todo", "/query"))
	// http.Handle("/query", handler.GraphQL(
	// 	graph.NewExecutableSchema(&graph.App{})))

	// fmt.Println("Listening on :8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
