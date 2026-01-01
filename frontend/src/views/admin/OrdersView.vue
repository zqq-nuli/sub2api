<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Statistics Cards -->
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary-100 dark:bg-primary-900/30">
              <svg class="h-5 w-5 text-primary-600 dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25z" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.totalOrders') }}</p>
              <p class="text-xl font-bold text-gray-900 dark:text-white">{{ statistics.total_orders }}</p>
            </div>
          </div>
        </div>

        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-yellow-100 dark:bg-yellow-900/30">
              <svg class="h-5 w-5 text-yellow-600 dark:text-yellow-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.pendingOrders') }}</p>
              <p class="text-xl font-bold text-yellow-600 dark:text-yellow-400">{{ statistics.pending_orders }}</p>
            </div>
          </div>
        </div>

        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-100 dark:bg-emerald-900/30">
              <svg class="h-5 w-5 text-emerald-600 dark:text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.paidOrders') }}</p>
              <p class="text-xl font-bold text-emerald-600 dark:text-emerald-400">{{ statistics.paid_orders }}</p>
            </div>
          </div>
        </div>

        <div class="card p-4">
          <div class="flex items-center gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-100 dark:bg-emerald-900/30">
              <svg class="h-5 w-5 text-emerald-600 dark:text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.totalAmount') }}</p>
              <p class="text-xl font-bold text-emerald-600 dark:text-emerald-400">{{ statistics.total_amount.toFixed(2) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="card">
        <div class="p-4">
          <div class="flex flex-wrap items-center gap-4">
            <div class="flex-1 min-w-[200px]">
              <input
                v-model="filters.order_no"
                type="text"
                :placeholder="t('admin.orders.searchOrderNo')"
                class="input"
                @input="debouncedSearch"
              />
            </div>
            <select v-model="filters.status" class="input w-auto" @change="fetchOrders">
              <option value="">{{ t('admin.orders.allStatus') }}</option>
              <option value="pending">{{ t('admin.orders.status.pending') }}</option>
              <option value="paid">{{ t('admin.orders.status.paid') }}</option>
              <option value="failed">{{ t('admin.orders.status.failed') }}</option>
              <option value="expired">{{ t('admin.orders.status.expired') }}</option>
            </select>
            <select v-model="filters.payment_method" class="input w-auto" @change="fetchOrders">
              <option value="">{{ t('admin.orders.allPaymentMethods') }}</option>
              <option value="alipay">{{ t('admin.orders.alipay') }}</option>
              <option value="wxpay">{{ t('admin.orders.wechat') }}</option>
              <option value="usdt">USDT</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Orders Table -->
      <div class="card">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50 dark:bg-dark-800">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.orderNo') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.user') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.product') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.amount') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.paymentMethod') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.statusColumn') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.createdAt') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.actions') }}
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
              <tr v-if="loading" class="text-center">
                <td colspan="8" class="px-4 py-8">
                  <svg class="mx-auto h-6 w-6 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </td>
              </tr>
              <tr v-else-if="orders.length === 0" class="text-center">
                <td colspan="8" class="px-4 py-8 text-gray-500 dark:text-dark-400">
                  {{ t('admin.orders.noOrders') }}
                </td>
              </tr>
              <tr v-for="order in orders" :key="order.id" class="hover:bg-gray-50 dark:hover:bg-dark-800">
                <td class="px-4 py-3">
                  <span class="font-mono text-sm">{{ order.order_no }}</span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-gray-900 dark:text-white">{{ order.user_email || `ID: ${order.user_id}` }}</span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-gray-900 dark:text-white">{{ order.product_name }}</span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm font-medium text-gray-900 dark:text-white">{{ order.amount.toFixed(2) }}</span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-gray-500 dark:text-dark-400">{{ getPaymentMethodLabel(order.payment_method) }}</span>
                </td>
                <td class="px-4 py-3">
                  <span :class="['inline-flex rounded-full px-2 py-1 text-xs font-medium', getStatusClass(order.status)]">
                    {{ t(`admin.orders.status.${order.status}`) }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(order.created_at) }}</span>
                </td>
                <td class="px-4 py-3">
                  <button
                    @click="viewOrder(order)"
                    class="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                  >
                    {{ t('admin.orders.view') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="flex items-center justify-between border-t border-gray-100 px-4 py-3 dark:border-dark-700">
          <p class="text-sm text-gray-500 dark:text-dark-400">
            {{ t('admin.orders.showing', { from: (page - 1) * limit + 1, to: Math.min(page * limit, total), total }) }}
          </p>
          <div class="flex gap-2">
            <button
              @click="prevPage"
              :disabled="page === 1"
              class="btn btn-secondary btn-sm"
            >
              {{ t('common.previous') }}
            </button>
            <button
              @click="nextPage"
              :disabled="page * limit >= total"
              class="btn btn-secondary btn-sm"
            >
              {{ t('common.next') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Order Detail Modal -->
      <div v-if="selectedOrder" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="selectedOrder = null">
        <div class="mx-4 max-h-[80vh] w-full max-w-lg overflow-y-auto rounded-lg bg-white shadow-xl dark:bg-dark-800">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('admin.orders.orderDetail') }}
              </h3>
              <button @click="selectedOrder = null" class="text-gray-400 hover:text-gray-500">
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
          <div class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.orderNo') }}</p>
                <p class="font-mono text-sm text-gray-900 dark:text-white">{{ selectedOrder.order_no }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.statusColumn') }}</p>
                <span :class="['inline-flex rounded-full px-2 py-1 text-xs font-medium', getStatusClass(selectedOrder.status)]">
                  {{ t(`admin.orders.status.${selectedOrder.status}`) }}
                </span>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.user') }}</p>
                <p class="text-sm text-gray-900 dark:text-white">{{ selectedOrder.user_email || `ID: ${selectedOrder.user_id}` }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.product') }}</p>
                <p class="text-sm text-gray-900 dark:text-white">{{ selectedOrder.product_name }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.paidAmount') }}</p>
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ selectedOrder.amount.toFixed(2) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.balance') }}</p>
                <p class="text-sm font-medium text-emerald-600 dark:text-emerald-400">${{ selectedOrder.actual_amount.toFixed(2) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.paymentMethod') }}</p>
                <p class="text-sm text-gray-900 dark:text-white">{{ getPaymentMethodLabel(selectedOrder.payment_method) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.tradeNo') }}</p>
                <p class="font-mono text-sm text-gray-900 dark:text-white">{{ selectedOrder.trade_no || '-' }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.createdAt') }}</p>
                <p class="text-sm text-gray-900 dark:text-white">{{ formatDateTime(selectedOrder.created_at) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.paidAt') }}</p>
                <p class="text-sm text-gray-900 dark:text-white">{{ selectedOrder.paid_at ? formatDateTime(selectedOrder.paid_at) : '-' }}</p>
              </div>
            </div>
            <div v-if="selectedOrder.notes">
              <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('admin.orders.notes') }}</p>
              <p class="text-sm text-gray-900 dark:text-white">{{ selectedOrder.notes }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { adminAPI } from '@/api/admin'
import type { AdminOrder, OrderStatistics } from '@/api/admin/orders'
import AppLayout from '@/components/layout/AppLayout.vue'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const appStore = useAppStore()

// State
const orders = ref<AdminOrder[]>([])
const statistics = reactive<OrderStatistics>({
  total_orders: 0,
  pending_orders: 0,
  paid_orders: 0,
  failed_orders: 0,
  expired_orders: 0,
  total_amount: 0,
  total_balance: 0
})
const selectedOrder = ref<AdminOrder | null>(null)
const loading = ref(false)
const page = ref(1)
const limit = ref(20)
const total = ref(0)

const filters = reactive({
  order_no: '',
  status: '',
  payment_method: ''
})

let searchTimeout: ReturnType<typeof setTimeout> | null = null

// Methods
const getStatusClass = (status: string) => {
  switch (status) {
    case 'paid':
      return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400'
    case 'pending':
      return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
    case 'failed':
      return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
    case 'expired':
      return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400'
    default:
      return 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400'
  }
}

const getPaymentMethodLabel = (method: string) => {
  switch (method) {
    case 'alipay':
      return t('admin.orders.alipay')
    case 'wxpay':
      return t('admin.orders.wechat')
    case 'usdt':
      return 'USDT'
    default:
      return method
  }
}

const fetchStatistics = async () => {
  try {
    const stats = await adminAPI.orders.getStatistics()
    Object.assign(statistics, stats)
  } catch (error) {
    console.error('Failed to fetch statistics:', error)
  }
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: page.value,
      limit: limit.value
    }
    if (filters.order_no) params.order_no = filters.order_no
    if (filters.status) params.status = filters.status
    if (filters.payment_method) params.payment_method = filters.payment_method

    const result = await adminAPI.orders.listOrders(params)
    orders.value = result.orders
    total.value = result.total
  } catch (error) {
    console.error('Failed to fetch orders:', error)
    appStore.showError(t('admin.orders.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchOrders()
  }, 300)
}

const viewOrder = (order: AdminOrder) => {
  selectedOrder.value = order
}

const prevPage = () => {
  if (page.value > 1) {
    page.value--
    fetchOrders()
  }
}

const nextPage = () => {
  if (page.value * limit.value < total.value) {
    page.value++
    fetchOrders()
  }
}

onMounted(() => {
  fetchStatistics()
  fetchOrders()
})
</script>
