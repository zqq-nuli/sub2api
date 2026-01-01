<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl space-y-6">
      <!-- Payment Polling Modal -->
      <div v-if="pollingStatus === 'polling'" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60">
        <div class="mx-4 w-full max-w-md rounded-lg bg-white p-8 text-center shadow-xl dark:bg-dark-800">
          <div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary-100 dark:bg-primary-900/30">
            <svg class="h-8 w-8 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('recharge.waitingForPayment') }}</h3>
          <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">{{ t('recharge.paymentWindowOpened') }}</p>
          <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">{{ t('recharge.orderNo') }}: {{ pollingOrderNo }}</p>
          <button @click="cancelPolling" class="btn btn-secondary mt-6">{{ t('recharge.cancelPayment') }}</button>
        </div>
      </div>

      <!-- Current Balance Card -->
      <div class="card overflow-hidden">
        <div class="bg-gradient-to-br from-primary-500 to-primary-600 px-6 py-8 text-center">
          <div
            class="mb-4 inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 backdrop-blur-sm"
          >
            <svg
              class="h-8 w-8 text-white"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z"
              />
            </svg>
          </div>
          <p class="text-sm font-medium text-primary-100">{{ t('recharge.currentBalance') }}</p>
          <p class="mt-2 text-4xl font-bold text-white">
            ${{ user?.balance?.toFixed(2) || '0.00' }}
          </p>
        </div>
      </div>

      <!-- Recharge Products -->
      <div class="card">
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('recharge.selectPackage') }}
          </h2>
        </div>
        <div class="p-6">
          <!-- Loading State -->
          <div v-if="loadingProducts" class="flex items-center justify-center py-8">
            <svg class="h-6 w-6 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
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

          <!-- Products Grid -->
          <div v-else-if="products.length > 0" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="product in products"
              :key="product.id"
              @click="selectProduct(product)"
              :class="[
                'relative cursor-pointer rounded-xl border-2 p-4 transition-all',
                selectedProduct?.id === product.id
                  ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-gray-200 hover:border-gray-300 dark:border-dark-600 dark:hover:border-dark-500'
              ]"
            >
              <!-- Hot Badge -->
              <div
                v-if="product.is_hot"
                class="absolute -top-2 -right-2 rounded-full bg-red-500 px-2 py-0.5 text-xs font-bold text-white"
              >
                {{ t('recharge.hot') }}
              </div>

              <!-- Discount Label -->
              <div
                v-if="product.discount_label"
                class="absolute -top-2 left-2 rounded-full bg-emerald-500 px-2 py-0.5 text-xs font-bold text-white"
              >
                {{ product.discount_label }}
              </div>

              <div class="text-center">
                <p class="text-lg font-semibold text-gray-900 dark:text-white">
                  {{ product.name }}
                </p>
                <p class="mt-2 text-3xl font-bold text-primary-600 dark:text-primary-400">
                  {{ product.amount.toFixed(2) }}
                </p>
                <p class="text-sm text-gray-500 dark:text-dark-400">
                  {{ t('recharge.getBalance', { balance: product.total_balance.toFixed(2) }) }}
                </p>
                <p v-if="product.bonus_balance > 0" class="mt-1 text-xs text-emerald-600 dark:text-emerald-400">
                  {{ t('recharge.bonus', { bonus: product.bonus_balance.toFixed(2) }) }}
                </p>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="empty-state py-8">
            <p class="text-sm text-gray-500 dark:text-dark-400">
              {{ t('recharge.noProducts') }}
            </p>
          </div>
        </div>
      </div>

      <!-- Payment Method -->
      <div v-if="selectedProduct" class="card">
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('recharge.selectPaymentMethod') }}
          </h2>
        </div>
        <div class="p-6">
          <!-- Loading State -->
          <div v-if="loadingChannels" class="flex items-center justify-center py-4">
            <svg class="h-5 w-5 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
          <!-- No Channels -->
          <div v-else-if="paymentChannels.length === 0" class="text-center py-4 text-gray-500 dark:text-dark-400">
            {{ t('recharge.noPaymentMethods') }}
          </div>
          <!-- Channels List -->
          <div v-else class="flex flex-wrap gap-4">
            <button
              v-for="channel in paymentChannels"
              :key="channel.key"
              @click="selectedPaymentMethod = channel.key"
              :class="[
                'flex items-center gap-2 rounded-lg border-2 px-4 py-3 transition-all',
                selectedPaymentMethod === channel.key
                  ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-gray-200 hover:border-gray-300 dark:border-dark-600 dark:hover:border-dark-500'
              ]"
            >
              <span class="text-2xl">{{ getChannelIcon(channel.icon) }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ channel.display_name }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Submit Button -->
      <div v-if="selectedProduct && selectedPaymentMethod" class="card">
        <div class="p-6">
          <div class="mb-4 rounded-lg bg-gray-50 p-4 dark:bg-dark-800">
            <div class="flex justify-between text-sm">
              <span class="text-gray-500 dark:text-dark-400">{{ t('recharge.package') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ selectedProduct.name }}</span>
            </div>
            <div class="mt-2 flex justify-between text-sm">
              <span class="text-gray-500 dark:text-dark-400">{{ t('recharge.paymentAmount') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ selectedProduct.amount.toFixed(2) }}</span>
            </div>
            <div class="mt-2 flex justify-between text-sm">
              <span class="text-gray-500 dark:text-dark-400">{{ t('recharge.balanceToAdd') }}</span>
              <span class="font-medium text-emerald-600 dark:text-emerald-400">${{ selectedProduct.total_balance.toFixed(2) }}</span>
            </div>
          </div>

          <button
            @click="handleRecharge"
            :disabled="submitting"
            class="btn btn-primary w-full py-3"
          >
            <svg
              v-if="submitting"
              class="-ml-1 mr-2 h-5 w-5 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
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
            {{ submitting ? t('recharge.processing') : t('recharge.payNow') }}
          </button>
        </div>
      </div>

      <!-- Order History -->
      <div class="card">
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('recharge.orderHistory') }}
          </h2>
        </div>
        <div class="p-6">
          <!-- Loading State -->
          <div v-if="loadingOrders" class="flex items-center justify-center py-8">
            <svg class="h-6 w-6 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
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

          <!-- Orders List -->
          <div v-else-if="orders.length > 0" class="space-y-3">
            <div
              v-for="order in orders"
              :key="order.id"
              class="flex items-center justify-between rounded-xl bg-gray-50 p-4 dark:bg-dark-800"
            >
              <div class="flex items-center gap-4">
                <div
                  :class="[
                    'flex h-10 w-10 items-center justify-center rounded-xl',
                    getStatusClass(order.status)
                  ]"
                >
                  <svg
                    v-if="order.status === 'paid'"
                    class="h-5 w-5 text-emerald-600 dark:text-emerald-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                  <svg
                    v-else-if="order.status === 'pending'"
                    class="h-5 w-5 text-yellow-600 dark:text-yellow-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <svg
                    v-else
                    class="h-5 w-5 text-red-600 dark:text-red-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ order.product_name }}
                  </p>
                  <p class="text-xs text-gray-500 dark:text-dark-400">
                    {{ formatDateTime(order.created_at) }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm font-semibold text-gray-900 dark:text-white">
                  {{ order.amount.toFixed(2) }}
                </p>
                <p :class="['text-xs', getStatusTextClass(order.status)]">
                  {{ t(`recharge.status.${order.status}`) }}
                </p>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="empty-state py-8">
            <div
              class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gray-100 dark:bg-dark-800"
            >
              <svg
                class="h-8 w-8 text-gray-400 dark:text-dark-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <p class="text-sm text-gray-500 dark:text-dark-400">
              {{ t('recharge.noOrders') }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { rechargeAPI, ordersAPI, type RechargeProduct, type Order, type PaymentChannel } from '@/api'
import AppLayout from '@/components/layout/AppLayout.vue'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const user = computed(() => authStore.user)

// State
const products = ref<RechargeProduct[]>([])
const orders = ref<Order[]>([])
const selectedProduct = ref<RechargeProduct | null>(null)
const selectedPaymentMethod = ref<string>('')
const loadingProducts = ref(false)
const loadingOrders = ref(false)
const submitting = ref(false)

// Payment channels (dynamic from backend)
const paymentChannels = ref<PaymentChannel[]>([])
const loadingChannels = ref(false)

// Methods
const selectProduct = (product: RechargeProduct) => {
  selectedProduct.value = product
}

const getChannelIcon = (icon: string) => {
  const iconMap: Record<string, string> = {
    'alipay': '',
    'wechat': '',
    'usdt': '',
    'wxpay': ''
  }
  return iconMap[icon] || ''
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'paid':
      return 'bg-emerald-100 dark:bg-emerald-900/30'
    case 'pending':
      return 'bg-yellow-100 dark:bg-yellow-900/30'
    default:
      return 'bg-red-100 dark:bg-red-900/30'
  }
}

const getStatusTextClass = (status: string) => {
  switch (status) {
    case 'paid':
      return 'text-emerald-600 dark:text-emerald-400'
    case 'pending':
      return 'text-yellow-600 dark:text-yellow-400'
    default:
      return 'text-red-600 dark:text-red-400'
  }
}

const fetchProducts = async () => {
  loadingProducts.value = true
  try {
    products.value = await rechargeAPI.getProducts()
  } catch (error) {
    console.error('Failed to fetch products:', error)
    appStore.showError(t('recharge.fetchProductsFailed'))
  } finally {
    loadingProducts.value = false
  }
}

const fetchPaymentChannels = async () => {
  loadingChannels.value = true
  try {
    paymentChannels.value = await rechargeAPI.getPaymentChannels()
    // 默认选择第一个支付渠道
    if (paymentChannels.value.length > 0 && !selectedPaymentMethod.value) {
      selectedPaymentMethod.value = paymentChannels.value[0].key
    }
  } catch (error) {
    console.error('Failed to fetch payment channels:', error)
  } finally {
    loadingChannels.value = false
  }
}

const fetchOrders = async () => {
  loadingOrders.value = true
  try {
    const result = await ordersAPI.getOrders(1, 10)
    orders.value = result.orders
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loadingOrders.value = false
  }
}

// 轮询订单状态
const pollingOrderNo = ref<string>('')
const pollingStatus = ref<'idle' | 'polling' | 'success' | 'failed'>('idle')
let pollTimer: ReturnType<typeof setInterval> | null = null
let pollCount = 0
const MAX_POLL_COUNT = 100 // 最大轮询次数（100次 * 3秒 = 5分钟）

const handleRecharge = async () => {
  if (!selectedProduct.value || !selectedPaymentMethod.value) {
    return
  }

  submitting.value = true
  try {
    const result = await ordersAPI.createOrder({
      product_id: selectedProduct.value.id,
      payment_method: selectedPaymentMethod.value,
      return_url: window.location.origin + '/recharge/result'
    })

    // 保存订单号并开始轮询
    pollingOrderNo.value = result.order.order_no
    pollingStatus.value = 'polling'
    startPolling()

    // 打开新窗口进行支付
    if (result.payment_params) {
      // 创建隐藏表单并在新窗口POST提交
      const form = document.createElement('form')
      form.method = 'POST'
      form.action = result.payment_url
      form.target = '_blank'  // 新窗口打开
      form.style.display = 'none'

      for (const [key, value] of Object.entries(result.payment_params)) {
        const input = document.createElement('input')
        input.type = 'hidden'
        input.name = key
        input.value = value as string
        form.appendChild(input)
      }

      document.body.appendChild(form)
      form.submit()
      document.body.removeChild(form)
    } else {
      window.open(result.payment_url, '_blank')
    }
  } catch (error: any) {
    console.error('Failed to create order:', error)
    appStore.showError(error.message || t('recharge.createOrderFailed'))
    pollingStatus.value = 'idle'
  } finally {
    submitting.value = false
  }
}

const startPolling = () => {
  if (pollTimer) return

  pollCount = 0 // 重置计数器
  pollTimer = setInterval(async () => {
    if (!pollingOrderNo.value) {
      stopPolling()
      return
    }

    pollCount++
    if (pollCount >= MAX_POLL_COUNT) {
      // 超过最大轮询次数，停止轮询
      pollingStatus.value = 'idle'
      stopPolling()
      appStore.showError(t('recharge.pollingTimeout'))
      await fetchOrders()
      return
    }

    try {
      const order = await ordersAPI.getOrderByNo(pollingOrderNo.value)

      if (order.status === 'paid') {
        pollingStatus.value = 'success'
        stopPolling()
        appStore.showSuccess(t('recharge.paymentSuccess'))
        // 刷新用户余额和订单列表
        await authStore.refreshUser()
        await fetchOrders()
      } else if (order.status === 'failed' || order.status === 'expired') {
        pollingStatus.value = 'failed'
        stopPolling()
        await fetchOrders()
      }
    } catch (error) {
      console.error('Failed to poll order status:', error)
    }
  }, 3000) // 每3秒轮询一次
}

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  pollingOrderNo.value = ''
}

const cancelPolling = () => {
  pollingStatus.value = 'idle'
  stopPolling()
}

// 检查是否有待支付订单（页面加载时）
const checkPendingOrder = async () => {
  try {
    const result = await ordersAPI.getOrders(1, 1)
    if (result.orders.length > 0) {
      const latestOrder = result.orders[0]
      // 如果最近的订单是待支付状态且未过期
      if (latestOrder.status === 'pending' && new Date(latestOrder.expired_at) > new Date()) {
        // 再次查询最新状态，确保订单确实是待支付
        const freshOrder = await ordersAPI.getOrderByNo(latestOrder.order_no)
        if (freshOrder.status === 'pending') {
          pollingOrderNo.value = freshOrder.order_no
          pollingStatus.value = 'polling'
          startPolling()
        } else if (freshOrder.status === 'paid') {
          // 订单已支付，刷新用户余额
          await authStore.refreshUser()
          await fetchOrders()
        }
      }
    }
  } catch (error) {
    console.error('Failed to check pending order:', error)
  }
}

onMounted(() => {
  fetchProducts()
  fetchPaymentChannels()
  fetchOrders()
  // 页面加载时检查是否有待支付订单
  checkPendingOrder()
})

onUnmounted(() => {
  stopPolling()
})
</script>
