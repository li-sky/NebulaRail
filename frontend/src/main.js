
import { createApp } from 'vue'
import App from './App.vue'

// 引入字体 (确保路径正确)
import 'vfonts/Lato.css' // 通用字体
import 'vfonts/FiraCode.css' // 等宽字体

const app = createApp(App)

app.mount('#app')
