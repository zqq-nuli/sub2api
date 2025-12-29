<template>
  <div class="overflow-hidden rounded-lg border border-gray-200 dark:border-dark-600">
    <table class="w-full table-auto">
      <thead class="bg-gray-50 dark:bg-dark-700">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">
            配置项
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">
            当前值
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">
            描述
          </th>
          <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">
            操作
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-200 bg-white dark:divide-dark-600 dark:bg-dark-800">
        <tr
          v-for="item in items"
          :key="item.key"
          class="transition-colors hover:bg-gray-50 dark:hover:bg-dark-700/50"
        >
          <!-- 配置项名称 -->
          <td class="whitespace-nowrap px-6 py-4">
            <div class="flex items-center">
              <div>
                <div class="text-sm font-medium text-gray-900 dark:text-white">
                  {{ item.label }}
                </div>
                <div v-if="item.required" class="mt-0.5 text-xs text-red-500">
                  * 必填
                </div>
              </div>
            </div>
          </td>

          <!-- 当前值 -->
          <td class="px-6 py-4">
            <SettingValueDisplay
              :type="item.type"
              :value="item.value"
              :editable="item.editable && item.type === 'boolean'"
              @change="handleValueChange(item.key, $event)"
            />
          </td>

          <!-- 描述 -->
          <td class="px-6 py-4">
            <div class="text-sm text-gray-500 dark:text-dark-400">
              {{ item.description }}
            </div>
          </td>

          <!-- 操作 -->
          <td class="whitespace-nowrap px-6 py-4 text-right text-sm">
            <button
              v-if="item.editable && item.type !== 'boolean'"
              @click="$emit('edit', item.key)"
              class="text-primary-600 transition-colors hover:text-primary-900 dark:text-primary-400 dark:hover:text-primary-300"
            >
              编辑
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import SettingValueDisplay from './SettingValueDisplay.vue'

export interface SettingItem {
  key: string           // 配置key
  label: string         // 显示名称
  value: any           // 当前值
  type: 'boolean' | 'text' | 'password' | 'number' | 'url' | 'array'
  description: string  // 描述
  editable: boolean    // 是否可编辑
  required?: boolean   // 是否必填
}

defineProps<{
  items: SettingItem[]
}>()

const emit = defineEmits<{
  change: [key: string, value: any]
  edit: [key: string]
}>()

function handleValueChange(key: string, value: any) {
  emit('change', key, value)
}
</script>
