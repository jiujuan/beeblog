<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Trash2, CheckCircle, XCircle } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Select } from '@/components/ui/select'
import { Pagination } from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import { commentApi } from '@/api/comment.api'
import { useToast } from '@/components/ui/toast'
import { fromNow } from '@/utils/format'
import { usePagination } from '@/composables/usePagination'
import type { Comment } from '@/types/comment.types'

const { toast } = useToast()
const comments = ref<Comment[]>([])
const loading = ref(false)
const statusFilter = ref('')
const { page, pageSize, total, pageCount, goTo, reset } = usePagination(15)

const statusMap: Record<number, { label: string; variant: any }> = {
  1: { label: '待审核', variant: 'secondary' },
  2: { label: '已通过', variant: 'default' },
  3: { label: '已拒绝', variant: 'destructive' },
}

async function fetchList() {
  loading.value = true
  try {
    const res = await commentApi.adminList({ status: statusFilter.value ? Number(statusFilter.value) : undefined, page: page.value, page_size: pageSize.value })
    comments.value = res.list
    total.value = res.total
  } finally { loading.value = false }
}

async function review(id: number, status: number) {
  try { await commentApi.adminReview(id, status); toast({ title: '操作成功', variant: 'success' }); fetchList() }
  catch (e: any) { toast({ title: '操作失败', description: e.message, variant: 'destructive' }) }
}

async function remove(id: number) {
  try { await commentApi.adminDelete(id); toast({ title: '删除成功', variant: 'success' }); fetchList() }
  catch (e: any) { toast({ title: '删除失败', description: e.message, variant: 'destructive' }) }
}

onMounted(fetchList)
watch(page, fetchList)
function applyFilter() { reset(); fetchList() }
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="font-serif text-2xl font-semibold">评论管理</h1>
        <p class="text-muted-foreground text-sm mt-0.5">共 {{ total }} 条评论</p>
      </div>
      <Select v-model="statusFilter" class="w-32" @change="applyFilter">
        <option value="">全部</option>
        <option value="1">待审核</option>
        <option value="2">已通过</option>
        <option value="3">已拒绝</option>
      </Select>
    </div>

    <div class="space-y-3">
      <div v-if="loading" class="space-y-3"><Skeleton v-for="i in 8" :key="i" class="h-20 w-full rounded-xl" /></div>
      <div
        v-for="c in comments"
        :key="c.id"
        class="rounded-xl border border-border p-4 space-y-2 hover:bg-muted/20 transition-colors"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <span class="font-medium text-sm">{{ c.nickname || '匿名' }}</span>
              <Badge :variant="statusMap[c.status]?.variant" class="text-xs">{{ statusMap[c.status]?.label }}</Badge>
              <span class="text-xs text-muted-foreground">{{ fromNow(c.created_at) }}</span>
            </div>
            <p class="text-sm text-muted-foreground mt-1 line-clamp-2">{{ c.content }}</p>
          </div>
          <div class="flex items-center gap-1 shrink-0">
            <button v-if="c.status !== 2" type="button" class="h-8 w-8 flex items-center justify-center rounded hover:bg-green-100 dark:hover:bg-green-900/30 text-green-600 transition-colors" title="通过" @click="review(c.id, 2)"><CheckCircle class="h-4 w-4" /></button>
            <button v-if="c.status !== 3" type="button" class="h-8 w-8 flex items-center justify-center rounded hover:bg-orange-100 dark:hover:bg-orange-900/30 text-orange-600 transition-colors" title="拒绝" @click="review(c.id, 3)"><XCircle class="h-4 w-4" /></button>
            <button type="button" class="h-8 w-8 flex items-center justify-center rounded hover:bg-destructive/10 text-muted-foreground hover:text-destructive transition-colors" @click="remove(c.id)"><Trash2 class="h-4 w-4" /></button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="pageCount > 1" class="flex justify-center">
      <Pagination :page="page" :page-count="pageCount" @update:page="goTo" />
    </div>
  </div>
</template>
