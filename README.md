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
  CreateMood(mood: {userid: "5b5ff11d2816453fe932f3b3", score: 5, comment: "mycommon"}) {
    id
    user {
      id
    }
    score
    comment
    time
  }
}
mutation DeleteMood {
  DeleteMood(id:"5b5ff62e28164548b2930030")
}

query User {
  User(id: "5b5ff11d2816453fe932f3b3") {
    id
    sex
    username
    password
    moods {
      id
      time
      comment
      score
    }
  }
}

query Users {
  Users {
    id
    sex
    username
    moods {
      id
      score
      comment
      time
    }
  }
}
```

## 资源
* [vue](https://cn.vuejs.org/index.html)
* [vue graphql client](https://akryum.github.io/vue-apollo/guide/apollo/queries.html#simple-query)
* [calendar-google-vue](https://github.com/FlowzPlatform/calendar-google-vue)
