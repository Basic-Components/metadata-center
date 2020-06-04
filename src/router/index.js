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
    path: '/component-info/:id',
    name: 'componentinfo',
    component: () => import('../views/ComponentInfo.vue'),
    props: true
  },
  {
    path: '/component-new',
    name: 'componentnew',
    component: () => import('../views/ComponentNew.vue')
  },
  {
    path: '/services',
    name: 'services',
    component: () => import('../views/Services.vue')
  },
  {
    path: '/service-info/:id',
    name: 'serviceinfo',
    component: () => import('../views/ServiceInfo.vue'),
    props: true
  },
  {
    path: '/service-new',
    name: 'servicenew',
    component: () => import('../views/ServiceNew.vue')
  }
]

const router = new VueRouter({
  routes
})
export default router
