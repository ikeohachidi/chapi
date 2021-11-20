import Vue from 'vue'
import Vuex from 'vuex'

import project from './modules/project';
import route from './modules/route';
import user from './modules/user';
import query from './modules/query';
import header from './modules/header';
import permOrigin from './modules/perm-origin';
import StoreState from './storeState';

Vue.use(Vuex)

export default new Vuex.Store<StoreState>({
  modules: {
    user,
    project,
    route,
    query,
    header,
    permOrigin
  }
})
