import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import Toast from "vue-toastification";

import 'remixicon/fonts/remixicon.css';
import "vue-toastification/dist/index.css";

import { Route, NavigationGuardNext } from 'vue-router';

import './index.css';

Vue.config.productionTip = false
Vue.use(Toast);

router.beforeEach(async (to: Route, from: Route, next: NavigationGuardNext) => {
  if (to.meta!.requiresAuth) {
    if (store.state.user.user) {
      next();
    } else {
      next({ path: from.path || '' })
    }

    return
  }

  next();
})

new Vue({
  router,
  store,
  render: h => h(App)
})
.$mount('#app')
