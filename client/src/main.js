import Vue from 'vue'
import App from './App.vue'
import './plugins/element.js'
import * as http from "./api/api"

Vue.config.productionTip = false
Vue.prototype.$api = http

new Vue({
  render: h => h(App),
}).$mount('#app')
