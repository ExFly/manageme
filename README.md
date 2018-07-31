# manage me
使用 `go graphql vue vuex` 做的管理自己日常todo+心情记录的项目

## 命令
* gqlgen
* go generate ./...

## todo
* [ ] 学习graphql，实现mongo数据查询
* [ ] 

## 发现一个bug，当vendor中有相关包的时候，会出现
```sh
go run main.go
# command-line-arguments./main.go:28:29: cannot use graphm.NewExecutableSchema(&graphm.App literal) (type "github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".ExecutableSchema) as type "github.com/vektah/gqlgen/graphql".ExecutableSchema in argument to handler.GraphQL:        "github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".ExecutableSchema does not implement "github.com/vektah/gqlgen/graphql".ExecutableSchema (wrong type for Mutation method)                have Mutation(context.Context, *"github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/neelance/query".Operation) *"github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".Response                want Mutation(context.Context, *"github.com/vektah/gqlgen/neelance/query".Operation) *"github.com/vektah/gqlgen/graphql".Response
```

## 检索
```graphql
mutation CreateUser {
		CreateUser(user: {sex: UNKNOWN, username: "username", password: "password"}) {
		  id
		}
}

mutation CreateMood {
		CreateMood(mood: {userid: "5b5fb8cb2816450c0605bda9", score: 5, comment: "mycommon"}) {
		  id
		  user {
				id
      }
  }
}

query Users {
  Users{
    id
    moods{
      id
      user{
        id
        moods{
          id
        }
      }
      score
      comment
    }
  }
}
```

## 资源
* [vue](https://cn.vuejs.org/index.html)
* [calendar-google-vue](https://github.com/FlowzPlatform/calendar-google-vue)
