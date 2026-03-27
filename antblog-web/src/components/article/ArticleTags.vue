<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { TagSimple } from '@/types/article.types'

defineProps<{ tags: TagSimple[] }>()
</script>

<template>
  <div class="flex flex-wrap gap-1.5">
    <RouterLink
      v-for="tag in tags"
      :key="tag.id"
      :to="`/tags/${tag.slug}`"
      class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-opacity hover:opacity-75"
      :style="{ backgroundColor: tag.color + '22', color: tag.color, borderColor: tag.color + '55' }"
      style="border-width: 1px;"
    >
      {{ tag.name }}
    </RouterLink>
  </div>
</template>
