<template>
  <div>
    <h2>This is a mood list</h2>
    <ul>
      <li v-for="mood of moods" :key="mood.id"/>
    </ul>
  </div>
</template>

<script>
// import gql from 'graphql-tag'
import CREATE_MOOD from '@/graphql/CreateMood.gql'
import GET_MOOD from '@/graphql/GetMood.gql'

export default {
  data() {
    return {
      moods: []
    }
  },
  apollo: {
    moods: {
      query: GET_MOOD
    }
  },
  methods: {
    addMood() {
      this.$apollo.mutate({
        mutation: CREATE_MOOD,
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

