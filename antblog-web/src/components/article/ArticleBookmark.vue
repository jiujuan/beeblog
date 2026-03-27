<script setup lang="ts">
import { Bookmark } from 'lucide-vue-next'
import { formatCount } from '@/utils/format'
import { cn } from '@/utils/cn'

const props = defineProps<{ bookmarked?: boolean; count: number; loading?: boolean }>()
const emit = defineEmits<{ click: [] }>()
</script>

<template>
  <button
    type="button"
    :disabled="loading"
    :class="cn(
      'group inline-flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-sm transition-all',
      'hover:border-amber-400/50 hover:text-amber-600',
      bookmarked
        ? 'border-amber-400/50 text-amber-600 bg-amber-50 dark:bg-amber-900/20'
        : 'border-border text-muted-foreground',
    )"
    @click="emit('click')"
  >
    <Bookmark
      class="h-4 w-4 transition-transform group-hover:scale-110"
      :class="bookmarked ? 'fill-amber-500' : ''"
    />
    <span>{{ formatCount(count) }}</span>
  </button>
</template>
