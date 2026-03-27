<script setup lang="ts">
import { cn } from '@/utils/cn'
interface Props { modelValue?: string | number; disabled?: boolean; class?: string }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:modelValue': [v: string] }>()
</script>
<template>
  <select
    :value="modelValue"
    :disabled="disabled"
    :class="cn(
      'flex h-9 w-full items-center rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm',
      'focus:outline-none focus:ring-2 focus:ring-ring',
      'disabled:cursor-not-allowed disabled:opacity-50',
      props.class,
    )"
    @change="emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
  >
    <slot />
  </select>
</template>
