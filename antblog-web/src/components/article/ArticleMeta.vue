<script setup lang="ts">
import { Eye, Clock, Calendar } from 'lucide-vue-next'
import { formatDate, formatCount, readingTime } from '@/utils/format'

defineProps<{
  publishedAt?: string | null
  viewCount?: number
  wordCount?: number
  showDate?: boolean
}>()
</script>

<template>
  <div class="flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
    <span v-if="showDate !== false && publishedAt" class="flex items-center gap-1">
      <Calendar class="h-3 w-3" />
      {{ formatDate(publishedAt) }}
    </span>
    <span v-if="viewCount !== undefined" class="flex items-center gap-1">
      <Eye class="h-3 w-3" />
      {{ formatCount(viewCount) }}
    </span>
    <span v-if="wordCount" class="flex items-center gap-1">
      <Clock class="h-3 w-3" />
      {{ readingTime(wordCount) }} 分钟
    </span>
  </div>
</template>
