import Vue from 'vue';
import Router from 'vue-router';

import MoodList from '@/components/MoodList.vue';
import Login from '@/components/Login.vue';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'moodlist',
      component: MoodList,
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
    },
  ],
});
