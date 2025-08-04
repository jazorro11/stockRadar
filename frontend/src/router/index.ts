import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Content from '@/components/Content.vue'
import StockDetails from '@/components/StockDetails.vue'

const routes: Array<RouteRecordRaw> = [
  { path: '/', component: Content },
  { path: '/:ticker', name: 'StockDetails', component: StockDetails }
//   { path: '/stock/:ticker', name: 'StockDetails', component: StockDetails }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router