import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import VueRouter from 'vue-router';
import Artists from './components/Artists';
import Tracks from './components/Tracks';
import Genres from './components/Genres';


Vue.config.productionTip = false

Vue.use(VueRouter);
const routes = [
  { path: '/tracks', component: Tracks, meta: {title: "Tracks"} },
  { path: '/artists', component: Artists, meta: {title: "Artists"} },
  { path: '/genres', component: Genres, meta: {title: "Genres"}},
]
const router = new VueRouter({
  routes
})

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app');
