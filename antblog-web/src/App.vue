<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import AppHeader from '@/components/common/AppHeader.vue'
import AppFooter from '@/components/common/AppFooter.vue'
import { Toast } from '@/components/ui/toast'
import { useTheme } from '@/composables/useTheme'
import { computed } from 'vue'

// Init theme on app mount
useTheme()

const route = useRoute()
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <!-- Admin layout handles its own layout -->
    <template v-if="isAdminRoute">
      <RouterView />
    </template>

    <!-- Blog layout -->
    <template v-else>
      <AppHeader />
      <main class="flex-1">
        <RouterView />
      </main>
      <AppFooter />
    </template>

    <!-- Global toast notifications -->
    <Toast />
  </div>
</template>
