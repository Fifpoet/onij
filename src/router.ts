import { createRouter, createWebHistory } from 'vue-router'

import NProgress from 'nprogress'
import './styles/nprogress.css'


const basicRoutes = [
  // {
  //   name: 'Home',
  //   path: '/',
  //   component: () => import('@/views/home/index.vue'),
  // }
]


export const router = createRouter({
  history: createWebHistory('/'),
  routes: basicRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

router.afterEach((to) => {
  document.title = `${to.meta?.title ?? import.meta.env.VITE_APP_TITLE}`
})

NProgress.configure({ showSpinner: false })

router.beforeEach((to, from, next) => {
  NProgress.start()
  for (let i = 0; i < 5; i++) NProgress.inc()
  setTimeout(() => NProgress.done(), 300)
  next()
})
