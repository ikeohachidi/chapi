import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  { 
    path: '/', 
    name: 'Home', 
    component: () => import(/* webpackChunkName: "home" */ '../views/Home.vue')
  },
  { 
    path: '/dashboard', 
    name: 'Dashboard', 
    meta: { requiresAuth: true },
    component: () => import(/* webpackChunkName: "dashboard" */ '../views/Dashboard.vue'),
    children: [
      {
        path: 'route',
        name: 'Route',
        meta: { requiresAuth: true },
        component: () => import(/* webpackChunkName: "dashboard" */ '../views/Route.vue'),
      },
      {
        path: 'routes-list',
        name: 'Routes List',
        meta: { requiresAuth: true },
        component: () => import(/* webpackChunkName: "dashboard" */ '../views/RouteList.vue'),
      }
    ]
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
