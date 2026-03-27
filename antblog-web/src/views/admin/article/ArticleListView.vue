<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { Plus, Pencil, Trash2, Eye } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Pagination } from '@/components/ui/pagination'
import { Dialog } from '@/components/ui/dialog'
import { articleApi } from '@/api/article.api'
import { useToast } from '@/components/ui/toast'
import { formatDate } from '@/utils/format'
import { usePagination } from '@/composables/usePagination'
import type { ArticleListItem } from '@/types/article.types'

const { toast } = useToast()
const articles = ref<ArticleListItem[]>([])
const loading = ref(false)
const keyword = ref('')
const statusFilter = ref('')
const deleteTarget = ref<ArticleListItem | null>(null)
const { page, pageSize, total, pageCount, goTo, reset } = usePagination()

const statusMap: Record<number, { label: string; variant: 'default' | 'secondary' | 'outline' | 'destructive' }> = {
  1: { label: '草稿',   variant: 'secondary' },
  2: { label: '已发布', variant: 'default'   },
  3: { label: '已归档', variant: 'outline'   },
}

async function fetchList() {
  loading.value = true
  try {
    const res = await articleApi.adminList({
      keyword: keyword.value || undefined,
      status: statusFilter.value ? Number(statusFilter.value) : undefined,
      page: page.value,
      page_size: pageSize.value,
    })
    articles.value = res.list
    total.value = res.total
  } finally {
    loading.value = false }
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  try {
    await articleApi.delete(deleteTarget.value.id)
    toast({ title: '删除成功', variant: 'success' })
    deleteTarget.value = null
    fetchList()
  } catch (e: any) {
    toast({ title: '删除失败', description: e.message, variant: 'destructive' })
  }
}

onMounted(fetchList)
watch(page, fetchList)

function applyFilter() { reset(); fetchList() }
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="font-serif text-2xl font-semibold">文章管理</h1>
        <p class="text-muted-foreground text-sm mt-0.5">共 {{ total }} 篇文章</p>
      </div>
      <RouterLink to="/admin/articles/new">
        <Button><Plus class="h-4 w-4 mr-1" /> 新建文章</Button>
      </RouterLink>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-2">
      <Input v-model="keyword" placeholder="搜索标题…" class="w-52" @keyup.enter="applyFilter" />
      <Select v-model="statusFilter" class="w-32">
        <option value="">全部状态</option>
        <option value="1">草稿</option>
        <option value="2">已发布</option>
        <option value="3">已归档</option>
      </Select>
      <Button variant="outline" @click="applyFilter">筛选</Button>
    </div>

    <!-- Table -->
    <div class="rounded-xl border border-border overflow-hidden">
      <div v-if="loading" class="p-4 space-y-3">
        <Skeleton v-for="i in 8" :key="i" class="h-12 w-full" />
      </div>
      <table v-else class="w-full text-sm">
        <thead class="bg-muted/50 border-b border-border">
          <tr class="text-left text-muted-foreground">
            <th class="px-4 py-3 font-medium">标题</th>
            <th class="px-4 py-3 font-medium hidden md:table-cell">状态</th>
            <th class="px-4 py-3 font-medium hidden lg:table-cell">发布时间</th>
            <th class="px-4 py-3 font-medium text-right hidden sm:table-cell">阅读</th>
            <th class="px-4 py-3 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr
            v-for="art in articles"
            :key="art.id"
            class="hover:bg-muted/30 transition-colors"
          >
            <td class="px-4 py-3">
              <span class="font-medium line-clamp-1">{{ art.title }}</span>
            </td>
            <td class="px-4 py-3 hidden md:table-cell">
              <Badge :variant="statusMap[art.status]?.variant">
                {{ statusMap[art.status]?.label }}
              </Badge>
            </td>
            <td class="px-4 py-3 text-muted-foreground hidden lg:table-cell">
              {{ formatDate(art.published_at) || '—' }}
            </td>
            <td class="px-4 py-3 text-right text-muted-foreground hidden sm:table-cell">
              {{ art.view_count }}
            </td>
            <td class="px-4 py-3">
              <div class="flex items-center justify-end gap-1">
                <a
                  :href="`/articles/${art.slug}`"
                  target="_blank"
                  class="h-8 w-8 flex items-center justify-center rounded hover:bg-accent transition-colors text-muted-foreground hover:text-foreground"
                >
                  <Eye class="h-3.5 w-3.5" />
                </a>
                <RouterLink
                  :to="`/admin/articles/${art.id}`"
                  class="h-8 w-8 flex items-center justify-center rounded hover:bg-accent transition-colors text-muted-foreground hover:text-foreground"
                >
                  <Pencil class="h-3.5 w-3.5" />
                </RouterLink>
                <button
                  type="button"
                  class="h-8 w-8 flex items-center justify-center rounded hover:bg-destructive/10 transition-colors text-muted-foreground hover:text-destructive"
                  @click="deleteTarget = art"
                >
                  <Trash2 class="h-3.5 w-3.5" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="pageCount > 1" class="flex justify-center">
      <Pagination :page="page" :page-count="pageCount" @update:page="goTo" />
    </div>

    <!-- Delete confirm dialog -->
    <Dialog :open="!!deleteTarget" title="确认删除" @update:open="deleteTarget = null">
      <p class="text-sm text-muted-foreground mb-4">
        确定要删除「{{ deleteTarget?.title }}」吗？此操作不可撤销。
      </p>
      <div class="flex gap-2 justify-end">
        <Button variant="outline" @click="deleteTarget = null">取消</Button>
        <Button variant="destructive" @click="confirmDelete">删除</Button>
      </div>
    </Dialog>
  </div>
</template>
