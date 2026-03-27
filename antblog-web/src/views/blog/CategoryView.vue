<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import ArticleList from '@/components/article/ArticleList.vue'
import AppSidebar from '@/components/common/AppSidebar.vue'
import { categoryApi } from '@/api/category.api'
import { articleApi } from '@/api/article.api'
import { usePagination } from '@/composables/usePagination'
import type { Category } from '@/types/category.types'
import type { ArticleListItem } from '@/types/article.types'

const route = useRoute()
const category = ref<Category | null>(null)
const articles = ref<ArticleListItem[]>([])
const loading = ref(false)
const { page, pageSize, total, pageCount, goTo, reset } = usePagination()

async function load() {
  loading.value = true
  try {
    const slug = route.params.slug as string
    category.value = await categoryApi.getBySlug(slug)
    const res = await articleApi.list({ category_id: category.value.id, page: page.value, page_size: pageSize.value })
    articles.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => route.params.slug, () => { reset(); load() })
watch(page, load)
</script>

<template>
  <div class="container py-8">
    <div class="grid grid-cols-1 lg:grid-cols-[1fr_260px] gap-10">
      <main>
        <div class="mb-8">
          <p class="text-sm text-muted-foreground mb-1">分类</p>
          <h1 class="font-serif text-3xl font-semibold">{{ category?.name ?? '…' }}</h1>
          <p v-if="category?.description" class="text-muted-foreground mt-1 text-sm">{{ category.description }}</p>
        </div>
        <ArticleList :articles="articles" :loading="loading" :page="page" :page-count="pageCount" @update:page="goTo" />
      </main>
      <aside class="hidden lg:block">
        <div class="sticky top-20"><AppSidebar /></div>
      </aside>
    </div>
  </div>
</template>
