import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
// 引入全局样式
import './assets/global.less'
import App from './App.vue'
import router from './router'
import '@/access.ts'
const app = createApp(App)
import VueCropper from 'vue-cropper'
import 'vue-cropper/dist/index.css'

app.use(VueCropper)

app.use(createPinia())
app.use(router)
app.use(Antd)

app.mount('#app')
