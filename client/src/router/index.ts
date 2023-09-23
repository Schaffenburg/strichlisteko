import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path:"/checkout",
    name:"checkout",
    component: () => import('../views/CheckoutView.vue'),
  },
  {
    path:"/storage",
    name:"storage",
    component: () => import('../views/StorageView.vue'),
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
