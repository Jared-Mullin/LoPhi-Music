import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
const VueCookies = require('vue-cookies');

Vue.config.productionTip = false

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')

Vue.use(VueCookies)