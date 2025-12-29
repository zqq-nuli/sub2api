<template>
  <div class="flex items-center">
    <!-- Boolean Toggle -->
    <div v-if="type === 'boolean'" class="flex items-center gap-3">
      <button
        v-if="editable"
        @click="toggleValue"
        :class="[
          'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 dark:focus:ring-offset-dark-800',
          modelValue ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
        ]"
      >
        <span
          :class="[
            'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
            modelValue ? 'translate-x-6' : 'translate-x-1'
          ]"
        ></span>
      </button>
      <span v-else class="text-sm" :class="modelValue ? 'text-green-600 dark:text-green-400' : 'text-gray-500 dark:text-dark-400'">
        {{ modelValue ? '已启用' : '已禁用' }}
      </span>
      <span class="text-sm font-medium" :class="modelValue ? 'text-green-600 dark:text-green-400' : 'text-gray-500 dark:text-dark-400'">
        {{ modelValue ? '开启' : '关闭' }}
      </span>
    </div>

    <!-- Text -->
    <span v-else-if="type === 'text'" class="text-sm text-gray-900 dark:text-white">
      {{ displayValue || '-' }}
    </span>

    <!-- Password (masked) -->
    <div v-else-if="type === 'password'" class="flex items-center gap-2">
      <span class="text-sm text-gray-500 dark:text-dark-400">
        {{ displayValue ? '••••••••' : '未设置' }}
      </span>
      <span v-if="displayValue" class="text-xs text-gray-400 dark:text-dark-500">
        (已加密)
      </span>
    </div>

    <!-- Number -->
    <span v-else-if="type === 'number'" class="text-sm font-mono text-gray-900 dark:text-white">
      {{ displayValue ?? '-' }}
    </span>

    <!-- URL -->
    <a
      v-else-if="type === 'url' && displayValue"
      :href="displayValue"
      target="_blank"
      rel="noopener noreferrer"
      class="flex items-center gap-1 text-sm text-primary-600 hover:underline dark:text-primary-400"
    >
      {{ truncateUrl(displayValue) }}
      <svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
      </svg>
    </a>
    <span v-else-if="type === 'url'" class="text-sm text-gray-500 dark:text-dark-400">
      未设置
    </span>

    <!-- Array -->
    <div v-else-if="type === 'array'" class="flex flex-wrap gap-1">
      <span
        v-if="Array.isArray(displayValue) && displayValue.length > 0"
        v-for="(item, index) in displayValue"
        :key="index"
        class="inline-flex items-center rounded-full bg-primary-100 px-2.5 py-0.5 text-xs font-medium text-primary-800 dark:bg-primary-900/30 dark:text-primary-300"
      >
        {{ item }}
      </span>
      <span v-else class="text-sm text-gray-500 dark:text-dark-400">
        无限制
      </span>
    </div>

    <!-- Default -->
    <span v-else class="text-sm text-gray-500 dark:text-dark-400">
      {{ displayValue || '-' }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  type: 'boolean' | 'text' | 'password' | 'number' | 'url' | 'array'
  value: any
  editable?: boolean
}>()

const emit = defineEmits<{
  change: [value: any]
}>()

const modelValue = computed(() => props.value)

const displayValue = computed(() => {
  if (props.value === null || props.value === undefined) {
    return null
  }
  return props.value
})

function toggleValue() {
  emit('change', !modelValue.value)
}

function truncateUrl(url: string, maxLength: number = 50): string {
  if (url.length <= maxLength) return url
  return url.substring(0, maxLength) + '...'
}
</script>
