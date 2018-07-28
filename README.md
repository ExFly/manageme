# manage me
使用 `go graphql vue vuex` 做的管理自己日常todo+心情记录的项目

## 命令
* gqlgen
## todo
* [ ] 学习graphql，实现mongo数据查询
* [ ] 

## 发现一个bug，当vendor中有相关包的时候，会出现
```sh
manageme go run main.go
# command-line-arguments./main.go:28:29: cannot use graphm.NewExecutableSchema(&graphm.App literal) (type "github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".ExecutableSchema) as type "github.com/vektah/gqlgen/graphql".ExecutableSchema in argument to handler.GraphQL:        "github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".ExecutableSchema does not implement "github.com/vektah/gqlgen/graphql".ExecutableSchema (wrong type for Mutation method)                have Mutation(context.Context, *"github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/neelance/query".Operation) *"github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".Response                want Mutation(context.Context, *"github.com/vektah/gqlgen/neelance/query".Operation) *"github.com/vektah/gqlgen/graphql".Response
```
