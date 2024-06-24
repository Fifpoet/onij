import { createApp, h } from 'vue'
import './style.css'
import App from './App.vue'
import 'virtual:uno.css'


import { router } from './router'
import { pinia } from './store'

const app = createApp({
    render: ()=>h(App)
})
app.use(router)
app.use(pinia)
app.mount('#app')
