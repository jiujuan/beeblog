<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Pencil, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog } from '@/components/ui/dialog'
import { Skeleton } from '@/components/ui/skeleton'
import { tagApi } from '@/api/tag.api'
import { useToast } from '@/components/ui/toast'
import { useTagStore } from '@/stores/tag.store'
import type { Tag } from '@/types/tag.types'

const { toast } = useToast()
const tagStore = useTagStore()
const loading = ref(false)
const dialogOpen = ref(false)
const editing = ref<Tag | null>(null)
const deleteTarget = ref<Tag | null>(null)
const form = ref({ name: '', slug: '', color: '#6B7280' })

function openCreate() { editing.value = null; form.value = { name: '', slug: '', color: '#6B7280' }; dialogOpen.value = true }
function openEdit(tag: Tag) { editing.value = tag; form.value = { name: tag.name, slug: tag.slug, color: tag.color }; dialogOpen.value = true }

async function save() {
  try {
    editing.value ? await tagApi.update(editing.value.id, form.value) : await tagApi.create(form.value)
    toast({ title: editing.value ? '更新成功' : '创建成功', variant: 'success' })
    dialogOpen.value = false
    tagStore.invalidate(); await tagStore.fetchAll()
  } catch (e: any) { toast({ title: '操作失败', description: e.message, variant: 'destructive' }) }
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  try {
    await tagApi.delete(deleteTarget.value.id)
    toast({ title: '删除成功', variant: 'success' })
    deleteTarget.value = null
    tagStore.invalidate(); await tagStore.fetchAll()
  } catch (e: any) { toast({ title: '删除失败', description: e.message, variant: 'destructive' }) }
}

onMounted(async () => { loading.value = true; tagStore.invalidate(); await tagStore.fetchAll(); loading.value = false })
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="font-serif text-2xl font-semibold">标签管理</h1>
      <Button @click="openCreate"><Plus class="h-4 w-4 mr-1" /> 新建标签</Button>
    </div>

    <div class="rounded-xl border border-border overflow-hidden">
      <div v-if="loading" class="p-4 space-y-3"><Skeleton v-for="i in 6" :key="i" class="h-10 w-full" /></div>
      <table v-else class="w-full text-sm">
        <thead class="bg-muted/50 border-b border-border">
          <tr class="text-left text-muted-foreground">
            <th class="px-4 py-3 font-medium">名称</th>
            <th class="px-4 py-3 font-medium hidden md:table-cell">颜色</th>
            <th class="px-4 py-3 font-medium text-right hidden sm:table-cell">文章数</th>
            <th class="px-4 py-3 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr v-for="tag in tagStore.tags" :key="tag.id" class="hover:bg-muted/30">
            <td class="px-4 py-3">
              <span class="inline-flex items-center gap-2">
                <span class="h-2.5 w-2.5 rounded-full" :style="{ backgroundColor: tag.color }" />
                <span :style="{ color: tag.color }" class="font-medium">{{ tag.name }}</span>
              </span>
            </td>
            <td class="px-4 py-3 font-mono text-xs text-muted-foreground hidden md:table-cell">{{ tag.color }}</td>
            <td class="px-4 py-3 text-right text-muted-foreground hidden sm:table-cell">{{ tag.article_count }}</td>
            <td class="px-4 py-3">
              <div class="flex items-center justify-end gap-1">
                <button type="button" class="h-8 w-8 flex items-center justify-center rounded hover:bg-accent text-muted-foreground" @click="openEdit(tag)"><Pencil class="h-3.5 w-3.5" /></button>
                <button type="button" class="h-8 w-8 flex items-center justify-center rounded hover:bg-destructive/10 text-muted-foreground hover:text-destructive" @click="deleteTarget = tag"><Trash2 class="h-3.5 w-3.5" /></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Dialog :open="dialogOpen" :title="editing ? '编辑标签' : '新建标签'" @update:open="dialogOpen = $event">
      <div class="space-y-3">
        <div class="space-y-1.5"><label class="text-sm font-medium">名称</label><Input v-model="form.name" /></div>
        <div class="space-y-1.5"><label class="text-sm font-medium">颜色</label>
          <div class="flex gap-2"><input type="color" v-model="form.color" class="h-9 w-12 cursor-pointer rounded border border-input" /><Input v-model="form.color" class="font-mono text-xs flex-1" /></div>
        </div>
        <div class="flex gap-2 justify-end pt-2"><Button variant="outline" @click="dialogOpen = false">取消</Button><Button @click="save">{{ editing ? '保存' : '创建' }}</Button></div>
      </div>
    </Dialog>

    <Dialog :open="!!deleteTarget" title="确认删除" @update:open="deleteTarget = null">
      <p class="text-sm text-muted-foreground mb-4">确定删除「{{ deleteTarget?.name }}」？</p>
      <div class="flex gap-2 justify-end"><Button variant="outline" @click="deleteTarget = null">取消</Button><Button variant="destructive" @click="confirmDelete">删除</Button></div>
    </Dialog>
  </div>
</template>
