<script setup lang="ts">
import { onMounted } from 'vue'
import ArticleList from '@/components/article/ArticleList.vue'
import AppSidebar from '@/components/common/AppSidebar.vue'
import { useArticleList } from '@/composables/useArticleList'

const { articles, loading, page, pageCount, fetchList, goTo } = useArticleList()

onMounted(fetchList)
</script>

<template>
  <div class="container py-8">
    <div class="grid grid-cols-1 lg:grid-cols-[1fr_260px] gap-10">
      <!-- Main content -->
      <main>
        <div class="mb-8">
          <h1 class="font-serif text-3xl font-semibold">最新文章</h1>
          <p class="text-muted-foreground mt-1 text-sm">记录思考，分享技术</p>
        </div>
        <ArticleList
          :articles="articles"
          :loading="loading"
          :page="page"
          :page-count="pageCount"
          @update:page="goTo"
        />
      </main>

      <!-- Sidebar -->
      <aside class="hidden lg:block">
        <div class="sticky top-20">
          <AppSidebar />
        </div>
      </aside>
    </div>
  </div>
</template>
