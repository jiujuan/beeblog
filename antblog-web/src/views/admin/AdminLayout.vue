<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import {
  LayoutDashboard, FileText, Folder, Tag, MessageSquare, Image,
  LogOut, Menu, X, ChevronRight, FilePlus2
} from 'lucide-vue-next'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import { useAuthStore } from '@/stores/auth.store'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const sidebarOpen = ref(true)

const navItems = [
  { icon: LayoutDashboard, label: '仪表盘',  to: '/admin/dashboard' },
  { icon: FileText,        label: '文章',     to: '/admin/articles', quickCreateTo: '/admin/articles/new' },
  { icon: Folder,          label: '分类',     to: '/admin/categories' },
  { icon: Tag,             label: '标签',     to: '/admin/tags' },
  { icon: MessageSquare,   label: '评论',     to: '/admin/comments' },
  { icon: Image,           label: '媒体',     to: '/admin/media' },
]

const pageTitle = computed(() => {
  const item = navItems.find((n) => route.path.startsWith(n.to))
  return item?.label ?? '后台管理'
})

async function handleLogout() {
  await authStore.logout()
  router.push('/')
}

function handleQuickCreate(path: string) {
  if (route.path === path) {
    void router.replace({ path, query: { ...route.query, _reset: Date.now().toString() } })
    return
  }
  void router.push(path)
}
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-background">
    <!-- Sidebar -->
    <aside
      :class="[
        'flex flex-col border-r border-border bg-card transition-all duration-300 shrink-0',
        sidebarOpen ? 'w-56' : 'w-14',
      ]"
    >
      <!-- Logo -->
      <div class="flex h-14 items-center border-b border-border px-3 gap-2 overflow-hidden">
        <RouterLink to="/" class="flex items-center gap-2 font-serif font-semibold text-primary whitespace-nowrap">
          <span class="text-lg">✦</span>
          <span v-if="sidebarOpen" class="text-sm">AntBlog Admin</span>
        </RouterLink>
      </div>

      <!-- Nav -->
      <nav class="flex-1 py-4 px-2 space-y-0.5 overflow-hidden">
        <div
          v-for="item in navItems"
          :key="item.to"
          class="flex items-center gap-1"
        >
          <RouterLink
            :to="item.to"
            :class="[
              'flex min-w-0 flex-1 items-center gap-3 rounded-md px-2 py-2 text-sm transition-colors',
              route.path.startsWith(item.to)
                ? 'bg-accent text-foreground font-medium'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent/60',
            ]"
          >
            <component :is="item.icon" class="h-4 w-4 shrink-0" />
            <span v-if="sidebarOpen" class="whitespace-nowrap">{{ item.label }}</span>
          </RouterLink>
          <RouterLink
            v-if="sidebarOpen && item.quickCreateTo"
            :to="item.quickCreateTo"
            class="h-8 w-8 shrink-0 rounded-md text-muted-foreground hover:text-primary hover:bg-accent/60 transition-colors flex items-center justify-center"
            title="新建文章"
            @click.prevent="handleQuickCreate(item.quickCreateTo)"
          >
            <FilePlus2 class="h-4 w-4" />
          </RouterLink>
        </div>
      </nav>

      <!-- Bottom -->
      <div class="border-t border-border p-2 space-y-1 overflow-hidden">
        <button
          type="button"
          :class="[
            'flex w-full items-center gap-3 rounded-md px-2 py-2 text-sm',
            'text-muted-foreground hover:text-destructive hover:bg-accent/60 transition-colors',
          ]"
          @click="handleLogout"
        >
          <LogOut class="h-4 w-4 shrink-0" />
          <span v-if="sidebarOpen">退出登录</span>
        </button>
      </div>
    </aside>

    <!-- Main area -->
    <div class="flex flex-1 flex-col overflow-hidden">
      <!-- Topbar -->
      <header class="flex h-14 shrink-0 items-center justify-between border-b border-border px-4">
        <div class="flex items-center gap-3">
          <button
            type="button"
            class="h-8 w-8 flex items-center justify-center rounded-md hover:bg-accent transition-colors"
            @click="sidebarOpen = !sidebarOpen"
          >
            <Menu v-if="!sidebarOpen" class="h-4 w-4" />
            <X v-else class="h-4 w-4" />
          </button>
          <div class="flex items-center gap-1.5 text-sm text-muted-foreground">
            <span>后台管理</span>
            <ChevronRight class="h-3.5 w-3.5" />
            <span class="text-foreground font-medium">{{ pageTitle }}</span>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <ThemeToggle />
          <span class="text-sm text-muted-foreground">{{ authStore.user?.nickname }}</span>
        </div>
      </header>

      <!-- Content -->
      <main class="flex-1 overflow-y-auto p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>
