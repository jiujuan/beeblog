<script setup lang="ts">
import { cn } from '@/utils/cn'
const props = defineProps<{ modelValue?: boolean; disabled?: boolean; class?: string }>()
const emit = defineEmits<{ 'update:modelValue': [v: boolean] }>()
</script>
<template>
  <button
    type="button"
    role="switch"
    :aria-checked="modelValue"
    :disabled="disabled"
    :class="cn(
      'relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent',
      'transition-colors duration-200 ease-in-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
      'disabled:cursor-not-allowed disabled:opacity-50',
      modelValue ? 'bg-primary' : 'bg-input',
      props.class,
    )"
    @click="emit('update:modelValue', !modelValue)"
  >
    <span
      :class="cn(
        'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-background shadow transition duration-200 ease-in-out',
        modelValue ? 'translate-x-4' : 'translate-x-0',
      )"
    />
  </button>
</template>
