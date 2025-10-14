import { createRouter, createWebHistory } from 'vue-router'

import NProgress from 'nprogress'
import './styles/nprogress.css'


const basicRoutes = [
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/views/SearchPage.vue')
  },
  {
    path: '/artist',
    name: 'Artist',
    component: () => import('@/views/ArtistDetail.vue')
  }
]


export const router = createRouter({
  history: createWebHistory('/'),
  routes: basicRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

NProgress.configure({ showSpinner: false })

router.beforeEach((to, from, next) => {
  NProgress.start()
  for (let i = 0; i < 5; i++) NProgress.inc()
  setTimeout(() => NProgress.done(), 300)
  next()
})
