<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { articleApi } from '@/api/article.api'
import { Skeleton } from '@/components/ui/skeleton'
import Pagination from '@/components/ui/pagination/Pagination.vue'
import { MONTHS, formatDate } from '@/utils/format'
import { usePagination } from '@/composables/usePagination'
import type { ArchiveItem, ArticleListItem } from '@/types/article.types'

const archive = ref<ArchiveItem[]>([])
const articles = ref<ArticleListItem[]>([])
const loading = ref(false)
const titleLoading = ref(false)
const { page, pageSize, total, pageCount, goTo } = usePagination(50)

async function loadArchive() {
  loading.value = true
  try {
    archive.value = await articleApi.archive()
  }
  finally { loading.value = false }
}

async function loadArticles() {
  titleLoading.value = true
  try {
    const res = await articleApi.list({
      page: page.value,
      page_size: pageSize.value,
    })
    articles.value = res.list
    total.value = res.total
  } finally {
    titleLoading.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadArchive(), loadArticles()])
})

// Group by year
const grouped = computed(() => {
  const map = new Map<number, ArchiveItem[]>()
  for (const item of archive.value) {
    if (!map.has(item.year)) map.set(item.year, [])
    map.get(item.year)!.push(item)
  }
  return [...map.entries()].sort((a, b) => b[0] - a[0])
})

watch(page, () => { void loadArticles() })
</script>

<template>
  <div class="container py-8 max-w-2xl mx-auto">
    <h1 class="font-serif text-3xl font-semibold mb-8">归档</h1>

    <div v-if="loading" class="space-y-8">
      <div v-for="i in 3" :key="i" class="space-y-3">
        <Skeleton class="h-7 w-20" />
        <div class="grid grid-cols-3 gap-3">
          <Skeleton v-for="j in 6" :key="j" class="h-16 rounded-lg" />
        </div>
      </div>
    </div>

    <div v-else class="space-y-10">
      <section v-for="[year, items] in grouped" :key="year">
        <h2 class="font-serif text-2xl font-semibold mb-4 text-primary">{{ year }}</h2>
        <div class="grid grid-cols-3 sm:grid-cols-4 gap-3">
          <RouterLink
            v-for="item in items"
            :key="item.month"
            :to="{ name: 'home', query: { year, month: item.month } }"
            class="flex flex-col items-center justify-center rounded-xl border border-border p-3 text-center
                   hover:border-primary/50 hover:bg-accent transition-all"
          >
            <span class="text-sm font-medium text-foreground">{{ MONTHS[item.month - 1] }}</span>
            <span class="text-xs text-muted-foreground mt-0.5">{{ item.article_count }} 篇</span>
          </RouterLink>
        </div>
      </section>

      <section class="space-y-4">
        <div class="flex items-center justify-between gap-2">
          <h2 class="font-serif text-2xl font-semibold text-primary">全部文章标题</h2>
          <span class="text-sm text-muted-foreground">共 {{ total }} 篇</span>
        </div>

        <div v-if="titleLoading" class="space-y-2">
          <Skeleton v-for="i in 10" :key="i" class="h-10 rounded-md" />
        </div>

        <div v-else class="rounded-xl border border-border divide-y divide-border">
          <RouterLink
            v-for="item in articles"
            :key="item.id"
            :to="{ name: 'article-detail', params: { slug: item.slug } }"
            class="flex items-center justify-between gap-4 px-4 py-3 hover:bg-accent/60 transition-colors"
          >
            <span class="text-sm sm:text-base text-foreground line-clamp-1">{{ item.title }}</span>
            <span class="shrink-0 text-xs text-muted-foreground">
              {{ formatDate(item.published_at || item.created_at) }}
            </span>
          </RouterLink>
          <div v-if="articles.length === 0" class="px-4 py-6 text-sm text-muted-foreground text-center">
            暂无文章
          </div>
        </div>

        <div v-if="total > 50 && pageCount > 1" class="flex justify-center pt-2">
          <Pagination :page="page" :page-count="pageCount" @update:page="goTo" />
        </div>
      </section>
    </div>
  </div>
</template>
