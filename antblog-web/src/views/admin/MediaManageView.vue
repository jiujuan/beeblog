<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Upload, Trash2, Link } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Pagination } from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import { Dialog } from '@/components/ui/dialog'
import { mediaApi } from '@/api/media.api'
import { useToast } from '@/components/ui/toast'
import { useUpload } from '@/composables/useUpload'
import { usePagination } from '@/composables/usePagination'
import type { Media } from '@/types/media.types'

const { toast } = useToast()
const { upload, selectFile, uploading } = useUpload()
const media = ref<Media[]>([])
const loading = ref(false)
const deleteTarget = ref<Media | null>(null)
const { page, pageSize, total, pageCount, goTo } = usePagination(20)

async function fetchList() {
  loading.value = true
  try {
    const res = await mediaApi.adminList({ page: page.value, page_size: pageSize.value })
    media.value = res.list
    total.value = res.total
  } finally { loading.value = false }
}

async function handleUpload() {
  const file = await selectFile()
  if (!file) return
  const result = await upload(file)
  if (result) { toast({ title: '上传成功', variant: 'success' }); fetchList() }
}

async function copyUrl(url: string) {
  await navigator.clipboard.writeText(url)
  toast({ title: '已复制 URL', variant: 'success' })
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  try { await mediaApi.delete(deleteTarget.value.id); toast({ title: '删除成功', variant: 'success' }); deleteTarget.value = null; fetchList() }
  catch (e: any) { toast({ title: '删除失败', description: e.message, variant: 'destructive' }) }
}

onMounted(fetchList)
watch(page, fetchList)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="font-serif text-2xl font-semibold">媒体管理</h1>
        <p class="text-muted-foreground text-sm mt-0.5">共 {{ total }} 个文件</p>
      </div>
      <Button :disabled="uploading" @click="handleUpload">
        <Upload class="h-4 w-4 mr-1" /> {{ uploading ? '上传中…' : '上传图片' }}
      </Button>
    </div>

    <!-- Grid -->
    <div v-if="loading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
      <Skeleton v-for="i in 15" :key="i" class="aspect-square rounded-lg" />
    </div>
    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
      <div
        v-for="item in media"
        :key="item.id"
        class="group relative aspect-square rounded-lg overflow-hidden border border-border bg-muted"
      >
        <img :src="item.url" :alt="item.original_name" class="h-full w-full object-cover" loading="lazy" />
        <!-- Overlay -->
        <div class="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-end">
          <div class="w-full p-2 flex items-center justify-between gap-1">
            <p class="text-white text-xs truncate flex-1">{{ item.original_name }}</p>
            <button type="button" class="h-7 w-7 flex items-center justify-center rounded bg-white/20 hover:bg-white/40 text-white transition-colors" @click.stop="copyUrl(item.url)"><Link class="h-3.5 w-3.5" /></button>
            <button type="button" class="h-7 w-7 flex items-center justify-center rounded bg-white/20 hover:bg-red-500/80 text-white transition-colors" @click.stop="deleteTarget = item"><Trash2 class="h-3.5 w-3.5" /></button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="pageCount > 1" class="flex justify-center">
      <Pagination :page="page" :page-count="pageCount" @update:page="goTo" />
    </div>

    <Dialog :open="!!deleteTarget" title="确认删除" @update:open="deleteTarget = null">
      <div v-if="deleteTarget" class="space-y-4">
        <img :src="deleteTarget.url" class="w-full h-40 object-contain rounded-md bg-muted" />
        <p class="text-sm text-muted-foreground">确定删除「{{ deleteTarget.original_name }}」？物理文件也将一并删除。</p>
        <div class="flex gap-2 justify-end">
          <Button variant="outline" @click="deleteTarget = null">取消</Button>
          <Button variant="destructive" @click="confirmDelete">删除</Button>
        </div>
      </div>
    </Dialog>
  </div>
</template>
