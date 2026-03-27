<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Pencil, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Dialog } from '@/components/ui/dialog'
import { Skeleton } from '@/components/ui/skeleton'
import { categoryApi } from '@/api/category.api'
import { useToast } from '@/components/ui/toast'
import { useCategoryStore } from '@/stores/category.store'
import type { Category } from '@/types/category.types'

const { toast } = useToast()
const catStore = useCategoryStore()
const loading = ref(false)
const dialogOpen = ref(false)
const editing = ref<Category | null>(null)
const deleteTarget = ref<Category | null>(null)

const form = ref({ name: '', slug: '', description: '', cover: '' })

function openCreate() { editing.value = null; form.value = { name: '', slug: '', description: '', cover: '' }; dialogOpen.value = true }
function openEdit(cat: Category) {
  editing.value = cat
  form.value = { name: cat.name, slug: cat.slug, description: cat.description, cover: cat.cover }
  dialogOpen.value = true
}

async function save() {
  try {
    if (editing.value) {
      await categoryApi.update(editing.value.id, form.value)
      toast({ title: '更新成功', variant: 'success' })
    } else {
      await categoryApi.create(form.value)
      toast({ title: '创建成功', variant: 'success' })
    }
    dialogOpen.value = false
    catStore.invalidate()
    await catStore.fetchAll()
  } catch (e: any) {
    toast({ title: '操作失败', description: e.message, variant: 'destructive' })
  }
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  try {
    await categoryApi.delete(deleteTarget.value.id)
    toast({ title: '删除成功', variant: 'success' })
    deleteTarget.value = null
    catStore.invalidate()
    await catStore.fetchAll()
  } catch (e: any) {
    toast({ title: '删除失败', description: e.message, variant: 'destructive' })
  }
}

onMounted(async () => {
  loading.value = true
  catStore.invalidate()
  await catStore.fetchAll()
  loading.value = false
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="font-serif text-2xl font-semibold">分类管理</h1>
      <Button @click="openCreate"><Plus class="h-4 w-4 mr-1" /> 新建分类</Button>
    </div>

    <div class="rounded-xl border border-border overflow-hidden">
      <div v-if="loading" class="p-4 space-y-3">
        <Skeleton v-for="i in 4" :key="i" class="h-12 w-full" />
      </div>
      <table v-else class="w-full text-sm">
        <thead class="bg-muted/50 border-b border-border">
          <tr class="text-left text-muted-foreground">
            <th class="px-4 py-3 font-medium">名称</th>
            <th class="px-4 py-3 font-medium hidden md:table-cell">Slug</th>
            <th class="px-4 py-3 font-medium text-right hidden sm:table-cell">文章数</th>
            <th class="px-4 py-3 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-for="cat in catStore.categories" :key="cat.id" class="hover:bg-muted/30">
            <td class="px-4 py-3 font-medium">{{ cat.name }}</td>
            <td class="px-4 py-3 text-muted-foreground font-mono text-xs hidden md:table-cell">{{ cat.slug }}</td>
            <td class="px-4 py-3 text-right text-muted-foreground hidden sm:table-cell">{{ cat.article_count }}</td>
            <td class="px-4 py-3">
              <div class="flex items-center justify-end gap-1">
                <button type="button" class="h-8 w-8 flex items-center justify-center rounded hover:bg-accent text-muted-foreground hover:text-foreground" @click="openEdit(cat)"><Pencil class="h-3.5 w-3.5" /></button>
                <button type="button" class="h-8 w-8 flex items-center justify-center rounded hover:bg-destructive/10 text-muted-foreground hover:text-destructive" @click="deleteTarget = cat"><Trash2 class="h-3.5 w-3.5" /></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Edit Dialog -->
    <Dialog :open="dialogOpen" :title="editing ? '编辑分类' : '新建分类'" @update:open="dialogOpen = $event">
      <div class="space-y-3">
        <div class="space-y-1.5"><label class="text-sm font-medium">名称</label><Input v-model="form.name" /></div>
        <div class="space-y-1.5"><label class="text-sm font-medium">Slug（可留空）</label><Input v-model="form.slug" class="font-mono text-xs" /></div>
        <div class="space-y-1.5"><label class="text-sm font-medium">描述</label><Textarea v-model="form.description" :rows="2" /></div>
        <div class="flex gap-2 justify-end pt-2">
          <Button variant="outline" @click="dialogOpen = false">取消</Button>
          <Button @click="save">{{ editing ? '保存' : '创建' }}</Button>
        </div>
      </div>
    </Dialog>

    <!-- Delete confirm -->
    <Dialog :open="!!deleteTarget" title="确认删除" @update:open="deleteTarget = null">
      <p class="text-sm text-muted-foreground mb-4">确定删除「{{ deleteTarget?.name }}」？</p>
      <div class="flex gap-2 justify-end">
        <Button variant="outline" @click="deleteTarget = null">取消</Button>
        <Button variant="destructive" @click="confirmDelete">删除</Button>
      </div>
    </Dialog>
  </div>
</template>
