import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 导入插件
import Components from 'unplugin-vue-components/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import AutoImport from 'unplugin-auto-import/vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    // 配置 unplugin-vue-components
    Components({
      resolvers: [NaiveUiResolver()] // 自动导入 Naive UI 组件
    }),
    // 配置 unplugin-auto-import
    AutoImport({
      imports: [
        'vue', // 自动导入 Vue 相关函数，如：ref, reactive, toRef 等
        {
          'naive-ui': [ // 自动导入 Naive UI 相关函数
            'useDialog',
            'useMessage',
            'useNotification',
            'useLoadingBar'
          ]
        }
      ]
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
