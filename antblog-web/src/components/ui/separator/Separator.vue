<script setup lang="ts">
import { cn } from '@/utils/cn'
const props = defineProps<{ class?: string; orientation?: 'horizontal' | 'vertical' }>()
</script>
<template>
  <div
    :class="cn(
      'shrink-0 bg-border',
      props.orientation === 'vertical' ? 'w-px h-full' : 'h-px w-full',
      props.class,
    )"
  />
</template>
