import Vue from 'vue'
import Vuex from 'vuex'
import menu from './modules/menu'
Vue.use(Vuex)
const debug = process.env.NODE_ENV !== 'production'
export default new Vuex.Store({
  state: {
  },
  mutations: {
  },
  actions: {
  },
  modules: {
    menu
  },
  strict: debug
})
