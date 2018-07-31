import Vue from 'vue'
import App from './App.vue'
import { createProvider } from './vue-apollo'
import 'bootstrap/scss/bootstrap.scss';

Vue.config.productionTip = false

new Vue({
  provide: createProvider().provide(),
  render: h => h(App)
}).$mount('#app')
