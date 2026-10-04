import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import router from './router'
import App from './App.vue'
import './index.css'
import './theme.css'

document.documentElement.classList.add('dark')

createApp(App).use(router).use(ElementPlus).mount('#app')