<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Save, Eye, Upload } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Select } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import ArticleContent from '@/components/article/ArticleContent.vue'
import { articleApi } from '@/api/article.api'
import { useCategoryStore } from '@/stores/category.store'
import { useTagStore } from '@/stores/tag.store'
import { useToast } from '@/components/ui/toast'
import { useUpload } from '@/composables/useUpload'

const route = useRoute()
const router = useRouter()
const { toast } = useToast()
const { upload, selectFile } = useUpload()
const catStore = useCategoryStore()
const tagStore = useTagStore()

const isNew = computed(() => route.params.id === undefined || route.path.endsWith('/new'))
const saving = ref(false)
const activeTab = ref('write')

function createInitialForm() {
  return {
    title: '',
    slug: '',
    summary: '',
    content: '',
    cover: '',
    category_id: undefined as number | undefined,
    tag_ids: [] as number[],
    is_top: false,
    is_featured: false,
    allow_comment: true,
    status: 1,
  }
}

const form = ref(createInitialForm())

function toOptionalNumber(value: unknown) {
  if (value === undefined || value === null || value === '' || value === 'undefined') return undefined
  const num = Number(value)
  return Number.isNaN(num) ? undefined : num
}

function normalizePayload() {
  return {
    title: form.value.title,
    slug: form.value.slug,
    summary: form.value.summary,
    content: form.value.content,
    cover: form.value.cover,
    category_id: toOptionalNumber(form.value.category_id),
    tag_ids: form.value.tag_ids.map((id) => Number(id)).filter((id) => !Number.isNaN(id)),
    is_top: form.value.is_top,
    is_featured: form.value.is_featured,
    allow_comment: form.value.allow_comment,
  }
}

async function load() {
  await Promise.all([catStore.fetchAll(), tagStore.fetchAll()])
  if (isNew.value) {
    form.value = createInitialForm()
    activeTab.value = 'write'
    return
  }
  const art = await articleApi.adminDetail(Number(route.params.id))
  form.value = {
    title: art.title,
    slug: art.slug,
    summary: art.summary,
    content: art.content,
    cover: art.cover,
    category_id: art.category_id ?? undefined,
    tag_ids: art.tags.map((t) => t.id),
    is_top: art.is_top,
    is_featured: art.is_featured,
    allow_comment: art.allow_comment,
    status: art.status,
  }
}

async function save(publish = false) {
  saving.value = true
  try {
    const payload = normalizePayload()
    const targetStatus = publish ? 2 : Number(form.value.status)
    if (isNew.value) {
      const art = await articleApi.create({
        ...payload,
        status: targetStatus,
      })
      toast({ title: '创建成功', variant: 'success' })
      router.replace(`/admin/articles/${art.id}`)
    } else {
      const id = Number(route.params.id)
      await articleApi.update(id, payload)
      await articleApi.updateStatus(id, { status: targetStatus })
      form.value.status = targetStatus
      toast({ title: '保存成功', variant: 'success' })
    }
  } catch (e: any) {
    toast({ title: '保存失败', description: e.message, variant: 'destructive' })
  } finally {
    saving.value = false
  }
}

async function handleCoverUpload() {
  const file = await selectFile()
  if (!file) return
  const media = await upload(file)
  if (media) form.value.cover = media.url
}

watch(() => route.fullPath, () => { void load() }, { immediate: true })
</script>

<template>
  <div class="space-y-5 max-w-5xl">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="font-serif text-2xl font-semibold">{{ isNew ? '新建文章' : '编辑文章' }}</h1>
      <div class="flex items-center gap-2">
        <Button variant="outline" :disabled="saving" @click="save(false)">
          <Save class="h-3.5 w-3.5 mr-1" /> 保存草稿
        </Button>
        <Button :disabled="saving" @click="save(true)">
          <Eye class="h-3.5 w-3.5 mr-1" /> {{ Number(form.status) === 2 ? '更新发布' : '发布' }}
        </Button>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-[1fr_240px] gap-5">
      <!-- Editor -->
      <div class="space-y-4">
        <Input v-model="form.title" placeholder="文章标题…" class="text-lg font-serif h-11" />
        <Input v-model="form.summary" placeholder="摘要（可选，不填则自动截取正文）" />

        <Tabs v-model="activeTab">
          <TabsList>
            <TabsTrigger value="write">编写</TabsTrigger>
            <TabsTrigger value="preview">预览</TabsTrigger>
          </TabsList>
          <TabsContent value="write">
            <Textarea
              v-model="form.content"
              placeholder="用 Markdown 写作…"
              :rows="24"
              class="font-mono text-sm resize-none"
            />
          </TabsContent>
          <TabsContent value="preview">
            <div class="min-h-[480px] rounded-md border border-border p-5">
              <ArticleContent :content="form.content" />
            </div>
          </TabsContent>
        </Tabs>
      </div>

      <!-- Meta sidebar -->
      <div class="space-y-5">
        <!-- Status -->
        <div class="rounded-xl border border-border p-4 space-y-3">
          <p class="text-sm font-medium">发布状态</p>
          <Select v-model="form.status">
            <option :value="1">草稿</option>
            <option :value="2">已发布</option>
            <option :value="3">已归档</option>
          </Select>
        </div>

        <!-- Category -->
        <div class="rounded-xl border border-border p-4 space-y-3">
          <p class="text-sm font-medium">分类</p>
          <Select v-model="form.category_id">
            <option :value="undefined">无分类</option>
            <option v-for="cat in catStore.categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </Select>
        </div>

        <!-- Tags -->
        <div class="rounded-xl border border-border p-4 space-y-3">
          <p class="text-sm font-medium">标签</p>
          <div class="flex flex-wrap gap-2">
            <label
              v-for="tag in tagStore.tags"
              :key="tag.id"
              class="flex items-center gap-1.5 cursor-pointer text-sm"
            >
              <input
                type="checkbox"
                :value="tag.id"
                v-model="form.tag_ids"
                class="rounded border-border"
              />
              <span :style="{ color: tag.color }">{{ tag.name }}</span>
            </label>
          </div>
        </div>

        <!-- Cover -->
        <div class="rounded-xl border border-border p-4 space-y-3">
          <p class="text-sm font-medium">封面图</p>
          <img v-if="form.cover" :src="form.cover" class="w-full rounded-md object-cover h-28" />
          <div class="flex gap-2">
            <Input v-model="form.cover" placeholder="图片 URL" class="flex-1" />
            <Button type="button" variant="outline" size="icon" @click="handleCoverUpload">
              <Upload class="h-4 w-4" />
            </Button>
          </div>
        </div>

        <!-- Options -->
        <div class="rounded-xl border border-border p-4 space-y-3">
          <p class="text-sm font-medium">其他选项</p>
          <label class="flex items-center justify-between text-sm">
            <span class="text-muted-foreground">置顶</span>
            <Switch v-model="form.is_top" />
          </label>
          <label class="flex items-center justify-between text-sm">
            <span class="text-muted-foreground">精选</span>
            <Switch v-model="form.is_featured" />
          </label>
          <label class="flex items-center justify-between text-sm">
            <span class="text-muted-foreground">允许评论</span>
            <Switch v-model="form.allow_comment" />
          </label>
        </div>

        <!-- Slug -->
        <div class="rounded-xl border border-border p-4 space-y-3">
          <p class="text-sm font-medium">URL Slug</p>
          <Input v-model="form.slug" placeholder="自动生成（可留空）" class="font-mono text-xs" />
        </div>
      </div>
    </div>
  </div>
</template>
