<script setup lang="ts">
import { inject, computed } from 'vue'
import { cn } from '@/utils/cn'
const props = defineProps<{ value: string; class?: string }>()
const tabs = inject<any>('tabs')
const isActive = computed(() => tabs?.active.value === props.value)
</script>
<template>
  <button
    type="button"
    :class="cn(
      'inline-flex items-center justify-center whitespace-nowrap rounded-md px-3 py-1 text-sm font-medium',
      'transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
      'disabled:pointer-events-none disabled:opacity-50',
      isActive ? 'bg-background text-foreground shadow' : 'hover:text-foreground',
      props.class,
    )"
    @click="tabs?.setActive(props.value)"
  >
    <slot />
  </button>
</template>
