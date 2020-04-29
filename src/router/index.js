import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/components',
    name: 'components',
    component: () => import('../views/Components.vue')
  },
  {
    path: '/component-info',
    name: 'component-info',
    component: () => import('../views/ComponentInfo.vue')
  },
  {
    path: '/component-new',
    name: 'component-new',
    component: () => import('../views/ComponentNew.vue')
  },
  {
    path: '/services',
    name: 'services',
    component: () => import('../views/Services.vue')
  },
  {
    path: '/service-info',
    name: 'service-info',
    component: () => import('../views/ServiceInfo.vue')
  },
  {
    path: '/service-new',
    name: 'service-new',
    component: () => import('../views/ServiceNew.vue')
  }
]

const router = new VueRouter({
  routes
})
export default router
