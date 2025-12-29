<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 dark:bg-gray-900">
    <div class="w-full max-w-md px-6">
      <div class="rounded-lg bg-white p-8 shadow-lg dark:bg-gray-800">
        <div class="text-center">
          <div v-if="loading" class="space-y-4">
            <div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">正在处理SSO登录...</h2>
            <p class="text-sm text-gray-600 dark:text-gray-400">请稍候</p>
          </div>

          <div v-else-if="error" class="space-y-4">
            <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/20">
              <svg class="h-6 w-6 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </div>
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">登录失败</h2>
            <p class="text-sm text-red-600 dark:text-red-400">{{ error }}</p>
            <button
              @click="goToLogin"
              class="w-full rounded-lg bg-primary-600 px-4 py-2 text-white transition hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
            >
              返回登录页面
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { apiClient } from '@/api/client'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const code = route.query.code as string
    const state = route.query.state as string
    const sessionID = sessionStorage.getItem('sso_session_id')

    if (!code || !state || !sessionID) {
      error.value = '缺少必要的回调参数'
      loading.value = false
      return
    }

    // 调用后端完成SSO登录
    const response = await apiClient.get<{
      access_token: string
      token_type: string
      user: {
        id: number
        email: string
        role: string
        balance: number
        concurrency: number
        status: string
        created_at: string
      }
      is_new_user: boolean
    }>('/auth/sso/callback', {
      params: {
        code,
        state,
        session_id: sessionID
      }
    })

    // 清除session storage
    sessionStorage.removeItem('sso_session_id')

    // 保存token
    authStore.token = response.data.access_token
    localStorage.setItem('auth_token', response.data.access_token)

    // 保存用户信息（使用as any绕过类型检查，因为SSO返回的user可能缺少某些字段）
    authStore.user = response.data.user as any
    localStorage.setItem('auth_user', JSON.stringify(response.data.user))

    // 跳转到仪表板
    router.push('/dashboard')
  } catch (err: any) {
    console.error('SSO callback error:', err)
    error.value = err.response?.data?.message || err.message || 'SSO登录失败'
    loading.value = false
  }
})

function goToLogin() {
  router.push('/login')
}
</script>
