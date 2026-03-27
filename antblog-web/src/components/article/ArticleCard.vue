<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { Pin } from 'lucide-vue-next'
import ArticleMeta from './ArticleMeta.vue'
import ArticleTags from './ArticleTags.vue'
import type { ArticleListItem } from '@/types/article.types'

defineProps<{ article: ArticleListItem }>()
</script>

<template>
  <article class="group relative flex flex-col gap-3 rounded-xl border border-border/60 bg-card p-5
                  transition-all hover:border-border hover:shadow-md hover:-translate-y-0.5">
    <!-- Top pin badge -->
    <div v-if="article.is_top" class="absolute right-4 top-4 flex items-center gap-1 text-xs text-primary">
      <Pin class="h-3 w-3 fill-current" />
      <span class="font-medium">置顶</span>
    </div>

    <!-- Cover image -->
    <RouterLink v-if="article.cover" :to="`/articles/${article.slug}`" class="block overflow-hidden rounded-lg">
      <img
        :src="article.cover"
        :alt="article.title"
        class="h-44 w-full object-cover transition-transform duration-500 group-hover:scale-[1.02]"
        loading="lazy"
      />
    </RouterLink>

    <!-- Content -->
    <div class="flex flex-col gap-2">
      <ArticleMeta
        :published-at="article.published_at"
        :view-count="article.view_count"
        :word-count="article.word_count"
      />

      <RouterLink :to="`/articles/${article.slug}`">
        <h2 class="font-serif text-lg font-semibold leading-snug text-foreground
                   group-hover:text-primary transition-colors line-clamp-2">
          {{ article.title }}
        </h2>
      </RouterLink>

      <p v-if="article.summary" class="text-sm text-muted-foreground line-clamp-3 leading-relaxed">
        {{ article.summary }}
      </p>

      <div class="flex items-center justify-between pt-1">
        <ArticleTags :tags="article.tags" />
        <div class="flex items-center gap-2 text-xs text-muted-foreground shrink-0 ml-2">
          <span>❤ {{ article.like_count }}</span>
          <span>🔖 {{ article.bookmark_count }}</span>
        </div>
      </div>
    </div>
  </article>
</template>
