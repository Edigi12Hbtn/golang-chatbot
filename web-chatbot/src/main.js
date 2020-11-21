import Vue from 'vue'
import App from './App.vue'
import Chat from 'vue-beautiful-chat'

Vue.config.productionTip = false
Vue.use(Chat)

new Vue({
  render: h => h(App),
}).$mount('#app')
