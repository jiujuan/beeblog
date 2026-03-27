<script setup lang="ts">
import { cn } from '@/utils/cn'

interface Props {
  class?: string
  placeholder?: string
  disabled?: boolean
  rows?: number
  modelValue?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
</script>

<template>
  <textarea
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :rows="rows ?? 4"
    :class="cn(
      'flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm',
      'transition-colors placeholder:text-muted-foreground resize-y min-h-[80px]',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
      'disabled:cursor-not-allowed disabled:opacity-50',
      props.class,
    )"
    @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
</template>
