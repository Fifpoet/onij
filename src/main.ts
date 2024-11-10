import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import 'virtual:uno.css'
import piniaPluginPersistedstate  from "pinia-plugin-persistedstate";

import { router } from './router'
import { pinia } from './store'

const app = createApp(App)
app.use(router)
app.use(pinia.use(piniaPluginPersistedstate))
app.mount('#app')
