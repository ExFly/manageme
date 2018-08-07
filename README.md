# manage me

使用 `go graphql vue` 做的管理自己日常todo+心情记录的项目

## How to use

```sh
git clone https://github.com/ExFly/manageme.git

# backend
dep ensure
docker-compose up -d
go run main.go

# frontend
cd frontend
yarn
yarn run serve
```

## Todo

- [ ] 学习graphql，实现mongo数据查询
- [ ] bootstrap

## Test

### CreateUser

```graphql
mutation CreateUser($user: UserInput!) {
  CreateUser(user: $user) {
    id
  }
}
# variables
{
  "user": {
    "sex": "MALE",
    "username": "huihui",
    "password": "123123"
  }
}
```

### CreateMood

```graphql
mutation CreateMood($mood: MoodInput!) {
  CreateMood(mood: $mood) {
    id
    user {
      id
    }
    score
    comment
    time
  }
}
# variables
{
  "mood": {
    "userid": "5b608e8d87bfadee3de41c8b",
    "score": 1,
    "comment": "fuck"
  }
}
```

### GetMood

```graphql
query GetMood($userID: ID!){
  User(id: $userID) {
    moods {
      id
      time
      comment
      score
    }
  }
}
# variables
{
  "userID": "5b608e8d87bfadee3de41c8b"
}
```

## login
* http://localhost:8080/loginas?user=username&pwd=password
* http://localhost:8080/query?query={me{id%20username}}
* http://localhost:8080/logout

## query create and del your mood
* http://localhost:8080/query?query=query%20queryme{%20me{%20id%20moods{%20id%20comment%20}%20}%20}
* http://localhost:8080/query?query=mutation%20CreateMood%20{%20CreateMood(mood:%20{score:%205,%20comment:%20%22mycommon%22})%20{%20id%20user%20{%20id%20}%20}%20}
* http://localhost:8080/query?query=mutation%20delMood{%20DeleteMood(id:%225b6972ae421aa9c59a31eefd%22)%20}

## Reference

- [vue](https://cn.vuejs.org/index.html)
- [vue graphql client](https://akryum.github.io/vue-apollo/guide/apollo/queries.html#simple-query)
- [calendar-google-vue](https://github.com/FlowzPlatform/calendar-google-vue)
- [GraphQL Guides](https://www.graphql.com/guides/)
- [gqlgen](https://gqlgen.com)

## Bug report

当vendor中有相关包的时候，会出现

```sh
go run main.go
# command-line-arguments./main.go:28:29: cannot use graphm.NewExecutableSchema(&graphm.App literal) (type "github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".ExecutableSchema) as type "github.com/vektah/gqlgen/graphql".ExecutableSchema in argument to handler.GraphQL:        "github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".ExecutableSchema does not implement "github.com/vektah/gqlgen/graphql".ExecutableSchema (wrong type for Mutation method)                have Mutation(context.Context, *"github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/neelance/query".Operation) *"github.com/exfly/manageme/vendor/github.com/vektah/gqlgen/graphql".Response                want Mutation(context.Context, *"github.com/vektah/gqlgen/neelance/query".Operation) *"github.com/vektah/gqlgen/graphql".Response
```
