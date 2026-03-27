<script setup lang="ts">
import ArticleCard from './ArticleCard.vue'
import { Pagination } from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import EmptyState from '@/components/common/EmptyState.vue'
import type { ArticleListItem } from '@/types/article.types'

defineProps<{
  articles: ArticleListItem[]
  loading?: boolean
  page: number
  pageCount: number
}>()

const emit = defineEmits<{ 'update:page': [p: number] }>()
</script>

<template>
  <div>
    <!-- Skeleton loading -->
    <div v-if="loading" class="grid gap-5 stagger">
      <div v-for="i in 6" :key="i" class="rounded-xl border border-border/60 p-5 space-y-3">
        <Skeleton class="h-44 w-full rounded-lg" />
        <Skeleton class="h-4 w-24" />
        <Skeleton class="h-6 w-3/4" />
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-2/3" />
      </div>
    </div>

    <!-- Article grid -->
    <div v-else-if="articles.length" class="grid gap-5 stagger">
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
    </div>

    <!-- Empty -->
    <EmptyState v-else icon="📝" title="暂无文章" description="还没有发布任何文章，稍后再来吧。" />

    <!-- Pagination -->
    <div v-if="pageCount > 1 && !loading" class="mt-10 flex justify-center">
      <Pagination
        :page="page"
        :page-count="pageCount"
        @update:page="emit('update:page', $event)"
      />
    </div>
  </div>
</template>
