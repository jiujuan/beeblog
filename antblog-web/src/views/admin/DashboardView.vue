<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { FileText, MessageSquare, Eye, Heart } from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { articleApi } from '@/api/article.api'
import { commentApi } from '@/api/comment.api'
import { Skeleton } from '@/components/ui/skeleton'

const stats = ref({
  articles: 0,
  comments: 0,
  totalViews: 0,
  totalLikes: 0,
})
const recentArticles = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const [arts, cmts] = await Promise.all([
      articleApi.adminList({ page: 1, page_size: 5 }),
      commentApi.adminList({ page: 1, page_size: 1 }),
    ])
    recentArticles.value = arts.list
    stats.value.articles = arts.total
    stats.value.comments = cmts.total
    stats.value.totalViews = arts.list.reduce((s: number, a: any) => s + a.view_count, 0)
    stats.value.totalLikes = arts.list.reduce((s: number, a: any) => s + a.like_count, 0)
  } finally {
    loading.value = false
  }
})

const statCards = [
  { label: '文章总数',  key: 'articles', icon: FileText,      color: 'text-blue-500'   },
  { label: '评论总数',  key: 'comments', icon: MessageSquare, color: 'text-green-500'  },
  { label: '总阅读量',  key: 'totalViews', icon: Eye,         color: 'text-amber-500'  },
  { label: '总点赞数',  key: 'totalLikes', icon: Heart,       color: 'text-rose-500'   },
]
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="font-serif text-2xl font-semibold">仪表盘</h1>
      <p class="text-muted-foreground text-sm mt-0.5">博客数据概览</p>
    </div>

    <!-- Stat cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <Card v-for="s in statCards" :key="s.key">
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-normal text-muted-foreground flex items-center gap-2">
            <component :is="s.icon" class="h-4 w-4" :class="s.color" />
            {{ s.label }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-16" />
          <p v-else class="text-3xl font-semibold font-serif">{{ (stats as any)[s.key] }}</p>
        </CardContent>
      </Card>
    </div>

    <!-- Recent articles -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <CardTitle>最近文章</CardTitle>
          <RouterLink to="/admin/articles" class="text-sm text-primary hover:underline">查看全部</RouterLink>
        </div>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="space-y-3">
          <Skeleton v-for="i in 5" :key="i" class="h-10 w-full" />
        </div>
        <table v-else class="w-full text-sm">
          <thead>
            <tr class="border-b border-border text-left text-muted-foreground">
              <th class="pb-2 font-medium">标题</th>
              <th class="pb-2 font-medium text-right">阅读</th>
              <th class="pb-2 font-medium text-right">点赞</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-for="art in recentArticles" :key="art.id">
              <td class="py-2.5">
                <RouterLink :to="`/admin/articles/${art.id}`" class="hover:text-primary transition-colors line-clamp-1">
                  {{ art.title }}
                </RouterLink>
              </td>
              <td class="py-2.5 text-right text-muted-foreground">{{ art.view_count }}</td>
              <td class="py-2.5 text-right text-muted-foreground">{{ art.like_count }}</td>
            </tr>
          </tbody>
        </table>
      </CardContent>
    </Card>
  </div>
</template>
