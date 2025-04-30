<script setup>
import { ref } from 'vue'
// useMessage 会被 unplugin-auto-import 自动导入
// import { useMessage } from 'naive-ui' // 如果自动导入不生效，可以手动导入

const password = ref('')
// const errorMessage = ref('') // 不再需要，使用 message API
const isLoading = ref(false)
const message = useMessage() // 获取 message API 实例

// 定义 emit
const emit = defineEmits(['login-success'])

async function handleLogin() {
  isLoading.value = true
  // errorMessage.value = '' // 清除旧错误
  try {
    const response = await fetch('/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ password: password.value }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || '登录失败')
    }

    // 登录成功
    localStorage.setItem('authToken', data.token)
    localStorage.setItem('tokenExpiry', Date.now() + data.expires_in * 1000)

    message.success('登录成功！') // 使用 Naive UI 的 message 提示

    // 触发登录成功事件
    emit('login-success')

  } catch (error) {
    console.error('Login error:', error)
    // errorMessage.value = error.message || '发生未知错误'
    message.error(error.message || '发生未知错误') // 使用 Naive UI 的 message 提示错误
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-view">
    <n-card class="login-card" :bordered="false">
      <n-h1 style="text-align: center; margin-bottom: 2rem;">
        <n-text type="primary">
          登录 NebulaRail
        </n-text>
      </n-h1>
      <n-form @submit.prevent="handleLogin">
        <n-form-item path="password" label="密码" label-props="{ for: 'password-input' }">
          <n-input
            id="password-input"
            v-model:value="password"
            type="password"
            placeholder="密码"
            show-password-on="click"
            :input-props="{ autocomplete: 'current-password' }"
            @keydown.enter="handleLogin"
            :disabled="isLoading"
          />
        </n-form-item>
        <n-button
          type="primary"
          block
          attr-type="submit"
          :loading="isLoading"
          :disabled="isLoading"
        >
          登录
        </n-button>
        <!-- 错误信息现在通过 message API 显示，不再需要 n-alert -->
        <!-- <n-alert title="错误" type="error" v-if="errorMessage" style="margin-top: 1rem;">
          {{ errorMessage }}
        </n-alert> -->
      </n-form>
    </n-card>
  </div>
</template>

<style scoped>
.login-view {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.login-card {
  width: 100%;
  max-width: 380px; /* 可以根据需要调整宽度 */
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1); /* 添加一点阴影 */
}

/* 可以移除大部分之前的样式，因为 Naive UI 接管了 */
</style>
