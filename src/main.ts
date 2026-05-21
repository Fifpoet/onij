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
