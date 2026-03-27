<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'
interface Props { page: number; pageCount: number; class?: string }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:page': [p: number] }>()
const pages = computed(() => {
  const total = props.pageCount
  const cur   = props.page
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
  if (cur <= 4)   return [1, 2, 3, 4, 5, '...', total]
  if (cur >= total - 3) return [1, '...', total-4, total-3, total-2, total-1, total]
  return [1, '...', cur-1, cur, cur+1, '...', total]
})
const btn = (p: number | string) => cn(
  'inline-flex h-8 w-8 items-center justify-center rounded-md text-sm transition-colors',
  typeof p === 'number' && p === props.page
    ? 'bg-primary text-primary-foreground font-semibold'
    : 'hover:bg-accent hover:text-accent-foreground',
  typeof p === 'string' ? 'cursor-default text-muted-foreground' : 'cursor-pointer',
)
</script>
<template>
  <nav :class="cn('flex items-center gap-1', props.class)">
    <button type="button" :disabled="page <= 1" :class="btn(0)" @click="emit('update:page', page-1)">‹</button>
    <button
      v-for="p in pages"
      :key="p"
      type="button"
      :class="btn(p)"
      :disabled="typeof p === 'string'"
      @click="typeof p === 'number' && emit('update:page', p)"
    >{{ p }}</button>
    <button type="button" :disabled="page >= pageCount" :class="btn(0)" @click="emit('update:page', page+1)">›</button>
  </nav>
</template>
