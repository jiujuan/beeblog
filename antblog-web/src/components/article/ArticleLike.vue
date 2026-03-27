<script setup lang="ts">
import { Heart } from 'lucide-vue-next'
import { formatCount } from '@/utils/format'
import { cn } from '@/utils/cn'

const props = defineProps<{ liked?: boolean; count: number; loading?: boolean }>()
const emit = defineEmits<{ click: [] }>()
</script>

<template>
  <button
    type="button"
    :disabled="loading"
    :class="cn(
      'group inline-flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-sm transition-all',
      'hover:border-rose-400/50 hover:text-rose-500',
      liked
        ? 'border-rose-400/50 text-rose-500 bg-rose-50 dark:bg-rose-900/20'
        : 'border-border text-muted-foreground',
    )"
    @click="emit('click')"
  >
    <Heart
      class="h-4 w-4 transition-transform group-hover:scale-110"
      :class="liked ? 'fill-rose-500' : ''"
    />
    <span>{{ formatCount(count) }}</span>
  </button>
</template>
