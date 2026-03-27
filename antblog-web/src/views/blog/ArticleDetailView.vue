<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { ChevronLeft } from 'lucide-vue-next'
import ArticleContent from '@/components/article/ArticleContent.vue'
import ArticleMeta from '@/components/article/ArticleMeta.vue'
import ArticleTags from '@/components/article/ArticleTags.vue'
import ArticleLike from '@/components/article/ArticleLike.vue'
import ArticleBookmark from '@/components/article/ArticleBookmark.vue'
import CommentList from '@/components/comment/CommentList.vue'
import { Skeleton } from '@/components/ui/skeleton'
import { useArticleDetail } from '@/composables/useArticleDetail'
import { useAuthStore } from '@/stores/auth.store'
import { useToast } from '@/components/ui/toast'

const route = useRoute()
const authStore = useAuthStore()
const { toast } = useToast()
const { article, loading, error, fetchDetail, toggleLike, toggleBookmark } = useArticleDetail()

async function load() {
  await fetchDetail(route.params.slug as string)
}

onMounted(load)
watch(() => route.params.slug, load)

async function handleLike() {
  if (!authStore.isLoggedIn) { toast({ title: '请先登录', variant: 'destructive' }); return }
  await toggleLike()
}
async function handleBookmark() {
  if (!authStore.isLoggedIn) { toast({ title: '请先登录', variant: 'destructive' }); return }
  await toggleBookmark()
}
</script>

<template>
  <div class="container py-8 max-w-3xl mx-auto">
    <RouterLink to="/" class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-8 transition-colors">
      <ChevronLeft class="h-4 w-4" />
      返回列表
    </RouterLink>

    <!-- Loading skeleton -->
    <div v-if="loading" class="space-y-4">
      <Skeleton class="h-10 w-3/4" />
      <Skeleton class="h-4 w-48" />
      <Skeleton class="h-52 w-full rounded-xl" />
      <div class="space-y-2 pt-4">
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-full" />
        <Skeleton class="h-4 w-5/6" />
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="py-20 text-center text-muted-foreground">
      <p class="text-lg">文章加载失败 😔</p>
      <p class="text-sm mt-1">{{ error }}</p>
    </div>

    <!-- Article -->
    <article v-else-if="article">
      <!-- Cover -->
      <img
        v-if="article.cover"
        :src="article.cover"
        :alt="article.title"
        class="w-full rounded-xl object-cover max-h-72 mb-8"
      />

      <!-- Title -->
      <h1 class="font-serif text-3xl md:text-4xl font-bold leading-tight mb-4 animate-fade-up">
        {{ article.title }}
      </h1>

      <!-- Meta + tags -->
      <div class="flex flex-wrap items-center gap-3 mb-6 animate-fade-up" style="animation-delay:60ms">
        <ArticleMeta
          :published-at="article.published_at"
          :view-count="article.view_count"
          :word-count="article.word_count"
        />
        <ArticleTags :tags="article.tags" />
      </div>

      <!-- Body -->
      <div class="animate-fade-up" style="animation-delay:120ms">
        <ArticleContent :content="article.content" :html="article.content_html" />
      </div>

      <!-- Interactions -->
      <div class="flex items-center gap-3 mt-10 pt-8 border-t border-border">
        <ArticleLike
          :liked="article.liked"
          :count="article.like_count"
          @click="handleLike"
        />
        <ArticleBookmark
          :bookmarked="article.bookmarked"
          :count="article.bookmark_count"
          @click="handleBookmark"
        />
      </div>

      <!-- Comments -->
      <CommentList :article-id="article.id" />
    </article>
  </div>
</template>
