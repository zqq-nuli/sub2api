<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
        @click.self="handleClose"
      >
        <div
          class="w-full max-w-lg rounded-lg bg-white shadow-xl dark:bg-dark-800"
          @click.stop
        >
          <!-- Header -->
          <div class="flex items-center justify-between border-b border-gray-200 px-6 py-4 dark:border-dark-600">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ title }}
            </h3>
            <button
              @click="handleClose"
              class="text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-dark-300"
            >
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="px-6 py-4">
            <div v-if="description" class="mb-4 text-sm text-gray-500 dark:text-dark-400">
              {{ description }}
            </div>

            <!-- Input based on type -->
            <div class="space-y-4">
              <!-- Text/URL Input -->
              <div v-if="inputType === 'text' || inputType === 'url'">
                <label class="input-label">{{ label }}</label>
                <input
                  v-model="localValue"
                  :type="inputType"
                  :placeholder="placeholder"
                  class="input w-full"
                  :required="required"
                />
              </div>

              <!-- Password Input -->
              <div v-else-if="inputType === 'password'">
                <label class="input-label">{{ label }}</label>
                <input
                  v-model="localValue"
                  type="password"
                  :placeholder="placeholder || '留空则不修改'"
                  class="input w-full"
                />
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">
                  密码已加密存储。留空则保持原密码不变。
                </p>
              </div>

              <!-- Number Input -->
              <div v-else-if="inputType === 'number'">
                <label class="input-label">{{ label }}</label>
                <input
                  v-model.number="localValue"
                  type="number"
                  :placeholder="placeholder"
                  :min="min"
                  :max="max"
                  class="input w-full"
                  :required="required"
                />
              </div>

              <!-- Array Input (comma-separated) -->
              <div v-else-if="inputType === 'array'">
                <label class="input-label">{{ label }}</label>
                <input
                  v-model="arrayInput"
                  type="text"
                  :placeholder="placeholder || '多个值用逗号分隔'"
                  class="input w-full"
                />
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">
                  多个值请用逗号分隔，例如：example.com, test.com
                </p>
                <div v-if="arrayPreview.length > 0" class="mt-2 flex flex-wrap gap-1">
                  <span
                    v-for="(item, index) in arrayPreview"
                    :key="index"
                    class="inline-flex items-center rounded-full bg-primary-100 px-2.5 py-0.5 text-xs font-medium text-primary-800 dark:bg-primary-900/30 dark:text-primary-300"
                  >
                    {{ item }}
                  </span>
                </div>
              </div>

              <!-- Textarea for long text -->
              <div v-else-if="inputType === 'textarea'">
                <label class="input-label">{{ label }}</label>
                <textarea
                  v-model="localValue"
                  :placeholder="placeholder"
                  rows="4"
                  class="input w-full resize-none"
                  :required="required"
                ></textarea>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="flex items-center justify-end gap-3 border-t border-gray-200 px-6 py-4 dark:border-dark-600">
            <button
              @click="handleClose"
              class="btn border-2 border-gray-300 bg-white text-gray-700 hover:bg-gray-50 dark:border-dark-600 dark:bg-dark-700 dark:text-dark-200 dark:hover:bg-dark-600"
            >
              取消
            </button>
            <button
              @click="handleSave"
              :disabled="isSaving"
              class="btn btn-primary"
            >
              <svg
                v-if="isSaving"
                class="mr-2 h-4 w-4 animate-spin"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isSaving ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'

const props = defineProps<{
  show: boolean
  title: string
  label: string
  description?: string
  inputType: 'text' | 'password' | 'number' | 'url' | 'array' | 'textarea'
  value: any
  placeholder?: string
  required?: boolean
  min?: number
  max?: number
}>()

const emit = defineEmits<{
  close: []
  save: [value: any]
}>()

const localValue = ref<any>(props.value)
const arrayInput = ref<string>('')
const isSaving = ref(false)

// Array preview
const arrayPreview = computed(() => {
  if (!arrayInput.value) return []
  return arrayInput.value.split(',').map(s => s.trim()).filter(s => s)
})

// Watch for value changes
watch(() => props.value, (newValue) => {
  if (props.inputType === 'array') {
    arrayInput.value = Array.isArray(newValue) ? newValue.join(', ') : ''
  } else {
    localValue.value = newValue
  }
}, { immediate: true })

watch(() => props.show, (newShow) => {
  if (newShow) {
    // Reset when dialog opens
    if (props.inputType === 'array') {
      arrayInput.value = Array.isArray(props.value) ? props.value.join(', ') : ''
    } else {
      localValue.value = props.value
    }
    isSaving.value = false
  }
})

function handleClose() {
  if (!isSaving.value) {
    emit('close')
  }
}

function handleSave() {
  isSaving.value = true

  let valueToSave: any
  if (props.inputType === 'array') {
    valueToSave = arrayPreview.value
  } else {
    valueToSave = localValue.value
  }

  emit('save', valueToSave)
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.95);
}
</style>
