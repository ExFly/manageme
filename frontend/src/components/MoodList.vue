<template>
  <div>
    <h2>This is a mood list</h2>
    <ul>
      <li v-for="mood of moods" :key="mood.id"/>
    </ul>
  </div>
</template>

<script>
import gql from 'graphql-tag'

export default {
  data() {
    return {
      moods: []
    }
  },
  apollo: {
    moods: {
      query: gql`
      query GetMood($userID: ID!){
        User(id: $userID) {
          moods {
            id
            time
            comment
            score
          }
        }
      }`
    }
  },
  methods: {
    addMood() {
      this.$apollo.mutate({
        mutation: gql`
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
        }`,
        variables: {
          user: this.$currentUser,
          score: this.score,
          comment: this.comment,
          time: new Date()
        },
      })
    }
  }
}
</script>

