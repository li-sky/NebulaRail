<script setup>
import { ref, onMounted } from 'vue'
import Login from './components/Login.vue'
// 导入 Naive UI Provider 组件
import { NMessageProvider, NDialogProvider, NNotificationProvider, NLoadingBarProvider, NConfigProvider, NGlobalStyle } from 'naive-ui'

const isLoggedIn = ref(false)

function checkAuthStatus() {
  const token = localStorage.getItem('authToken')
  const expiry = localStorage.getItem('tokenExpiry')
  if (token && expiry && Date.now() < parseInt(expiry)) {
    isLoggedIn.value = true
  } else {
    localStorage.removeItem('authToken')
    localStorage.removeItem('tokenExpiry')
    isLoggedIn.value = false
  }
}

function handleLoginSuccess() {
  isLoggedIn.value = true
}

function handleLogout() {
  localStorage.removeItem('authToken')
  localStorage.removeItem('tokenExpiry')
  isLoggedIn.value = false
}

onMounted(() => {
  checkAuthStatus()
})

// Naive UI 主题配置 (可选)
const themeOverrides = {
  // common: {
  //   primaryColor: '#FF0000'
  // }
}
</script>

<template>
  <!-- 配置 Naive UI 全局配置和样式 -->
  <n-config-provider :theme-overrides="themeOverrides">
    <n-global-style />
    <!-- 配置消息提示 Provider -->
    <n-message-provider>
      <!-- 配置对话框 Provider -->
      <n-dialog-provider>
        <!-- 配置通知 Provider -->
        <n-notification-provider>
          <!-- 配置加载条 Provider -->
          <n-loading-bar-provider>
            <!-- 应用主体内容 -->
            <div id="app-content">
              <div v-if="isLoggedIn">
                <header>
                  <h1>欢迎回来！</h1>
                  <n-button @click="handleLogout" type="warning" ghost>登出</n-button>
                </header>
                <main>
                  <p>这里是 NebulaRail 的主应用界面。</p>
                </main>
              </div>
              <div v-else>
                <Login @login-success="handleLoginSuccess" />
              </div>
            </div>
          </n-loading-bar-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped>
/* 移除或调整之前的样式，因为 Naive UI 会接管大部分样式 */
#app-content { /* 给应用内容一个容器 */
  min-height: 100vh; /* 确保内容至少占满屏幕高度 */
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

header h1 {
  margin: 0;
  font-size: 1.5rem;
}

main {
  padding: 1rem;
}
</style>
