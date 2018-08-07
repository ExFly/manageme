<template>
  <div class="home">
    <MoodList />
  </div>
</template>

<script>
// @ is an alias to /src
import MoodList from "@/components/MoodList.vue";

export default {
  name: "home",
  components: {
    MoodList
  }
};
</script>
