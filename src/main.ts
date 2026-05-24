import { createApp } from 'vue'
import 'virtual:uno.css'
import './style.css'
import App from './App.vue'
import piniaPluginPersistedstate  from "pinia-plugin-persistedstate";

import { router } from './router'
import { pinia } from './store'
import MateChat from '@matechat/core';

import '@devui-design/icons/icomoon/devui-icon.css';

const app = createApp(App)
app.use(MateChat)
app.use(router)
app.use(pinia.use(piniaPluginPersistedstate))
app.mount('#app')

// 路由表变更不会走组件 HMR，开发时改 router.ts 后自动整页刷新
if (import.meta.hot) {
  import.meta.hot.accept('./router.ts', () => {
    window.location.reload()
  })
}
