<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArticleList from '@/components/article/ArticleList.vue'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Search } from 'lucide-vue-next'
import { articleApi } from '@/api/article.api'
import { usePagination } from '@/composables/usePagination'
import type { ArticleListItem } from '@/types/article.types'

const route = useRoute()
const router = useRouter()
const keyword = ref((route.query.q as string) ?? '')
const articles = ref<ArticleListItem[]>([])
const loading = ref(false)
const { page, pageSize, total, pageCount, goTo, reset } = usePagination()

async function search() {
  if (!keyword.value.trim()) return
  loading.value = true
  try {
    const res = await articleApi.list({ keyword: keyword.value, page: page.value, page_size: pageSize.value })
    articles.value = res.list
    total.value = res.total
  } finally {
    loading.value = false }
}

function submitSearch() {
  reset()
  router.replace({ query: { q: keyword.value } })
  search()
}

onMounted(() => { if (keyword.value) search() })
watch(() => route.query.q, (q) => { keyword.value = q as string ?? ''; if (keyword.value) search() })
watch(page, search)
</script>

<template>
  <div class="container py-8 max-w-3xl mx-auto">
    <h1 class="font-serif text-3xl font-semibold mb-6">搜索</h1>

    <form class="flex gap-2 mb-8" @submit.prevent="submitSearch">
      <div class="relative flex-1">
        <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input v-model="keyword" class="pl-10" placeholder="输入关键字搜索文章…" />
      </div>
      <Button type="submit">搜索</Button>
    </form>

    <p v-if="keyword && !loading" class="text-sm text-muted-foreground mb-4">
      "{{ keyword }}" 的搜索结果，共 {{ total }} 篇
    </p>

    <ArticleList :articles="articles" :loading="loading" :page="page" :page-count="pageCount" @update:page="goTo" />
  </div>
</template>
