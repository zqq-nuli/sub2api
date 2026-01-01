<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('admin.rechargeProducts.title') }}
        </h1>
        <button @click="openCreateModal" class="btn btn-primary">
          <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          {{ t('admin.rechargeProducts.create') }}
        </button>
      </div>

      <!-- Products Table -->
      <div class="card">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50 dark:bg-dark-800">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.name') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.amount') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.balance') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.bonus') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.sortOrder') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.status') }}
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.actions') }}
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
              <tr v-if="loading" class="text-center">
                <td colspan="7" class="px-4 py-8">
                  <svg class="mx-auto h-6 w-6 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </td>
              </tr>
              <tr v-else-if="products.length === 0" class="text-center">
                <td colspan="7" class="px-4 py-8 text-gray-500 dark:text-dark-400">
                  {{ t('admin.rechargeProducts.noProducts') }}
                </td>
              </tr>
              <tr v-for="product in products" :key="product.id" class="hover:bg-gray-50 dark:hover:bg-dark-800">
                <td class="px-4 py-3">
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-medium text-gray-900 dark:text-white">{{ product.name }}</span>
                    <span v-if="product.is_hot" class="rounded-full bg-red-500 px-2 py-0.5 text-xs text-white">
                      {{ t('admin.rechargeProducts.hot') }}
                    </span>
                    <span v-if="product.discount_label" class="rounded-full bg-emerald-500 px-2 py-0.5 text-xs text-white">
                      {{ product.discount_label }}
                    </span>
                  </div>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-gray-900 dark:text-white">{{ product.amount.toFixed(2) }}</span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-gray-900 dark:text-white">${{ product.balance.toFixed(2) }}</span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-emerald-600 dark:text-emerald-400">+${{ product.bonus_balance.toFixed(2) }}</span>
                </td>
                <td class="px-4 py-3">
                  <span class="text-sm text-gray-500 dark:text-dark-400">{{ product.sort_order }}</span>
                </td>
                <td class="px-4 py-3">
                  <span :class="['inline-flex rounded-full px-2 py-1 text-xs font-medium', product.status === 'active' ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400' : 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400']">
                    {{ product.status === 'active' ? t('admin.rechargeProducts.active') : t('admin.rechargeProducts.inactive') }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button @click="openEditModal(product)" class="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400">
                      {{ t('common.edit') }}
                    </button>
                    <button @click="confirmDelete(product)" class="text-sm text-red-600 hover:text-red-700 dark:text-red-400">
                      {{ t('common.delete') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Create/Edit Modal -->
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="closeModal">
        <div class="card mx-4 max-h-[80vh] w-full max-w-lg overflow-y-auto">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ editingProduct ? t('admin.rechargeProducts.edit') : t('admin.rechargeProducts.create') }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-500">
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
          <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
            <div>
              <label class="input-label">{{ t('admin.rechargeProducts.name') }}</label>
              <input v-model="form.name" type="text" class="input" required />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="input-label">{{ t('admin.rechargeProducts.amount') }} (CNY)</label>
                <input v-model.number="form.amount" type="number" step="0.01" min="0" class="input" required />
              </div>
              <div>
                <label class="input-label">{{ t('admin.rechargeProducts.balance') }} (USD)</label>
                <input v-model.number="form.balance" type="number" step="0.01" min="0" class="input" required />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="input-label">{{ t('admin.rechargeProducts.bonus') }} (USD)</label>
                <input v-model.number="form.bonus_balance" type="number" step="0.01" min="0" class="input" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.rechargeProducts.sortOrder') }}</label>
                <input v-model.number="form.sort_order" type="number" min="0" class="input" />
              </div>
            </div>
            <div>
              <label class="input-label">{{ t('admin.rechargeProducts.descriptionLabel') }}</label>
              <textarea v-model="form.description" rows="2" class="input"></textarea>
            </div>
            <div>
              <label class="input-label">{{ t('admin.rechargeProducts.discountLabel') }}</label>
              <input v-model="form.discount_label" type="text" class="input" :placeholder="t('admin.rechargeProducts.discountLabelPlaceholder')" />
            </div>
            <div class="flex items-center gap-6">
              <label class="flex items-center gap-2">
                <input v-model="form.is_hot" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                <span class="text-sm text-gray-700 dark:text-dark-300">{{ t('admin.rechargeProducts.markAsHot') }}</span>
              </label>
              <label v-if="editingProduct" class="flex items-center gap-2">
                <input v-model="form.is_active" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                <span class="text-sm text-gray-700 dark:text-dark-300">{{ t('admin.rechargeProducts.active') }}</span>
              </label>
            </div>
            <div class="flex justify-end gap-3 pt-4">
              <button type="button" @click="closeModal" class="btn btn-secondary">
                {{ t('common.cancel') }}
              </button>
              <button type="submit" :disabled="submitting" class="btn btn-primary">
                <svg v-if="submitting" class="-ml-1 mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ editingProduct ? t('common.save') : t('common.create') }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Delete Confirmation Modal -->
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="showDeleteModal = false">
        <div class="card mx-4 w-full max-w-md">
          <div class="p-6 text-center">
            <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
              <svg class="h-6 w-6 text-red-600 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.rechargeProducts.confirmDelete') }}
            </h3>
            <p class="mb-6 text-sm text-gray-500 dark:text-dark-400">
              {{ t('admin.rechargeProducts.deleteWarning', { name: deletingProduct?.name }) }}
            </p>
            <div class="flex justify-center gap-3">
              <button @click="showDeleteModal = false" class="btn btn-secondary">
                {{ t('common.cancel') }}
              </button>
              <button @click="handleDelete" :disabled="submitting" class="btn btn-danger">
                <svg v-if="submitting" class="-ml-1 mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ t('common.delete') }}
              </button>
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
import type { RechargeProduct } from '@/api/admin/recharge-products'
import AppLayout from '@/components/layout/AppLayout.vue'

const { t } = useI18n()
const appStore = useAppStore()

// State
const products = ref<RechargeProduct[]>([])
const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingProduct = ref<RechargeProduct | null>(null)
const deletingProduct = ref<RechargeProduct | null>(null)

const form = reactive({
  name: '',
  amount: 0,
  balance: 0,
  bonus_balance: 0,
  description: '',
  sort_order: 0,
  is_hot: false,
  discount_label: '',
  is_active: true
})

// Methods
const fetchProducts = async () => {
  loading.value = true
  try {
    products.value = await adminAPI.rechargeProducts.listProducts()
  } catch (error) {
    console.error('Failed to fetch products:', error)
    appStore.showError(t('admin.rechargeProducts.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  editingProduct.value = null
  Object.assign(form, {
    name: '',
    amount: 0,
    balance: 0,
    bonus_balance: 0,
    description: '',
    sort_order: 0,
    is_hot: false,
    discount_label: '',
    is_active: true
  })
  showModal.value = true
}

const openEditModal = (product: RechargeProduct) => {
  editingProduct.value = product
  Object.assign(form, {
    name: product.name,
    amount: product.amount,
    balance: product.balance,
    bonus_balance: product.bonus_balance,
    description: product.description,
    sort_order: product.sort_order,
    is_hot: product.is_hot,
    discount_label: product.discount_label,
    is_active: product.status === 'active'
  })
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingProduct.value = null
}

const confirmDelete = (product: RechargeProduct) => {
  deletingProduct.value = product
  showDeleteModal.value = true
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (editingProduct.value) {
      await adminAPI.rechargeProducts.updateProduct(editingProduct.value.id, {
        name: form.name,
        amount: form.amount,
        balance: form.balance,
        bonus_balance: form.bonus_balance,
        description: form.description,
        sort_order: form.sort_order,
        is_hot: form.is_hot,
        discount_label: form.discount_label,
        status: form.is_active ? 'active' : 'inactive'
      })
      appStore.showSuccess(t('admin.rechargeProducts.updateSuccess'))
    } else {
      await adminAPI.rechargeProducts.createProduct({
        name: form.name,
        amount: form.amount,
        balance: form.balance,
        bonus_balance: form.bonus_balance,
        description: form.description,
        sort_order: form.sort_order,
        is_hot: form.is_hot,
        discount_label: form.discount_label
      })
      appStore.showSuccess(t('admin.rechargeProducts.createSuccess'))
    }
    closeModal()
    fetchProducts()
  } catch (error: any) {
    console.error('Failed to save product:', error)
    appStore.showError(error.message || t('admin.rechargeProducts.saveFailed'))
  } finally {
    submitting.value = false
  }
}

const handleDelete = async () => {
  if (!deletingProduct.value) return

  submitting.value = true
  try {
    await adminAPI.rechargeProducts.deleteProduct(deletingProduct.value.id)
    appStore.showSuccess(t('admin.rechargeProducts.deleteSuccess'))
    showDeleteModal.value = false
    deletingProduct.value = null
    fetchProducts()
  } catch (error: any) {
    console.error('Failed to delete product:', error)
    appStore.showError(error.message || t('admin.rechargeProducts.deleteFailed'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchProducts()
})
</script>
