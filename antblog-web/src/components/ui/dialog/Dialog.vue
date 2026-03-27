<script setup lang="ts">
import { cn } from '@/utils/cn'
interface Props { open?: boolean; title?: string; class?: string }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:open': [v: boolean] }>()
</script>
<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('update:open', false)" />
        <div :class="cn('relative z-10 w-full max-w-md rounded-xl border bg-card p-6 shadow-xl', props.class)">
          <div v-if="title" class="mb-4 flex items-center justify-between">
            <h2 class="font-serif text-lg font-semibold">{{ title }}</h2>
            <button type="button" class="text-muted-foreground hover:text-foreground" @click="emit('update:open', false)">✕</button>
          </div>
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
<style scoped>
.dialog-enter-active, .dialog-leave-active { transition: opacity 0.2s; }
.dialog-enter-from, .dialog-leave-to { opacity: 0; }
</style>
