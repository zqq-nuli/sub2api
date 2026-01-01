<template>
  <AppLayout>
    <div class="mx-auto max-w-lg">
      <div class="card">
        <div class="p-8 text-center">
          <!-- Loading State -->
          <div v-if="loading" class="space-y-4">
            <div class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-primary-100 dark:bg-primary-900/30">
              <svg class="h-10 w-10 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            </div>
            <p class="text-lg font-medium text-gray-900 dark:text-white">
              {{ t('paymentResult.checking') }}
            </p>
            <p class="text-sm text-gray-500 dark:text-dark-400">
              {{ t('paymentResult.pleaseWait') }}
            </p>
          </div>

          <!-- Success State -->
          <div v-else-if="order && order.status === 'paid'" class="space-y-4">
            <div class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-emerald-100 dark:bg-emerald-900/30">
              <svg
                class="h-10 w-10 text-emerald-600 dark:text-emerald-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h2 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ t('paymentResult.success') }}
            </h2>
            <p class="text-gray-500 dark:text-dark-400">
              {{ t('paymentResult.balanceAdded', { amount: order.actual_amount.toFixed(2) }) }}
            </p>
            <div class="mt-6 rounded-lg bg-gray-50 p-4 dark:bg-dark-800">
              <div class="flex justify-between text-sm">
                <span class="text-gray-500 dark:text-dark-400">{{ t('paymentResult.orderNo') }}</span>
                <span class="font-mono text-gray-900 dark:text-white">{{ order.order_no }}</span>
              </div>
              <div class="mt-2 flex justify-between text-sm">
                <span class="text-gray-500 dark:text-dark-400">{{ t('paymentResult.paidAmount') }}</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ order.amount.toFixed(2) }}</span>
              </div>
              <div class="mt-2 flex justify-between text-sm">
                <span class="text-gray-500 dark:text-dark-400">{{ t('paymentResult.balanceReceived') }}</span>
                <span class="font-medium text-emerald-600 dark:text-emerald-400">${{ order.actual_amount.toFixed(2) }}</span>
              </div>
            </div>
          </div>

          <!-- Pending State -->
          <div v-else-if="order && order.status === 'pending'" class="space-y-4">
            <div class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-yellow-100 dark:bg-yellow-900/30">
              <svg
                class="h-10 w-10 text-yellow-600 dark:text-yellow-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h2 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ t('paymentResult.pending') }}
            </h2>
            <p class="text-gray-500 dark:text-dark-400">
              {{ t('paymentResult.waitingForPayment') }}
            </p>
            <p class="text-sm text-gray-400 dark:text-dark-500">
              {{ t('paymentResult.autoRefresh') }}
            </p>
          </div>

          <!-- Failed/Expired State -->
          <div v-else-if="order" class="space-y-4">
            <div class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
              <svg
                class="h-10 w-10 text-red-600 dark:text-red-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <h2 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ order.status === 'expired' ? t('paymentResult.expired') : t('paymentResult.failed') }}
            </h2>
            <p class="text-gray-500 dark:text-dark-400">
              {{ order.status === 'expired' ? t('paymentResult.orderExpired') : t('paymentResult.paymentFailed') }}
            </p>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="space-y-4">
            <div class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
              <svg
                class="h-10 w-10 text-red-600 dark:text-red-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <h2 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ t('paymentResult.error') }}
            </h2>
            <p class="text-gray-500 dark:text-dark-400">
              {{ error }}
            </p>
          </div>

          <!-- Actions -->
          <div class="mt-8 space-y-3">
            <router-link
              to="/recharge"
              class="btn btn-primary w-full"
            >
              {{ t('paymentResult.backToRecharge') }}
            </router-link>
            <router-link
              to="/dashboard"
              class="btn btn-secondary w-full"
            >
              {{ t('paymentResult.backToDashboard') }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ordersAPI, type Order } from '@/api'
import AppLayout from '@/components/layout/AppLayout.vue'

const { t } = useI18n()
const route = useRoute()
const authStore = useAuthStore()

const order = ref<Order | null>(null)
const loading = ref(true)
const error = ref('')
let pollInterval: ReturnType<typeof setInterval> | null = null
let pollCount = 0
const MAX_POLL_COUNT = 100 // 最大轮询次数（100次 * 3秒 = 5分钟）

const fetchOrder = async () => {
  // 尝试多种方式获取订单号：
  // 1. 易支付标准参数 out_trade_no
  // 2. 自定义参数 order_no
  // 3. localStorage 备选
  let orderNo = (route.query.out_trade_no || route.query.order_no) as string

  if (!orderNo) {
    // 尝试从 localStorage 获取
    orderNo = localStorage.getItem('pending_order_no') || ''
  }

  if (!orderNo) {
    error.value = t('paymentResult.noOrderNo')
    loading.value = false
    return
  }

  try {
    order.value = await ordersAPI.getOrderByNo(orderNo)

    // If paid, refresh user data and clear localStorage
    if (order.value.status === 'paid') {
      localStorage.removeItem('pending_order_no')
      await authStore.refreshUser()
      stopPolling()
    }

    // If not pending, stop polling
    if (order.value.status !== 'pending') {
      stopPolling()
    }
  } catch (err: any) {
    error.value = err.message || t('paymentResult.fetchFailed')
    stopPolling()
  } finally {
    loading.value = false
  }
}

const startPolling = () => {
  pollCount = 0 // 重置计数器
  // Poll every 3 seconds
  pollInterval = setInterval(() => {
    pollCount++
    if (pollCount >= MAX_POLL_COUNT) {
      // 超过最大轮询次数，停止轮询
      stopPolling()
      if (order.value && order.value.status === 'pending') {
        error.value = t('paymentResult.pollingTimeout')
      }
      return
    }
    fetchOrder()
  }, 3000)
}

const stopPolling = () => {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
}

onMounted(() => {
  fetchOrder()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>
