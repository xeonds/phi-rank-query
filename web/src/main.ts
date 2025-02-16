import { createApp } from 'vue'
import router from './router'
import App from './App.vue'
import './index.css'
import VueLazyload from 'vue-lazyload'

createApp(App).use(router).use(VueLazyload).mount('#app')
