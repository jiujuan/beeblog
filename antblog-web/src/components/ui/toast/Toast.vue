<script setup lang="ts">
import { useToast } from './useToast'
import { cn } from '@/utils/cn'
const { toasts, dismiss } = useToast()
const variantClass = (v?: string) => {
  if (v === 'destructive') return 'border-destructive bg-destructive text-destructive-foreground'
  if (v === 'success')     return 'border-green-500 bg-green-50 text-green-900 dark:bg-green-900/20 dark:text-green-300'
  return 'border-border bg-card text-card-foreground'
}
</script>
<template>
  <Teleport to="body">
    <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 w-80">
      <TransitionGroup name="toast">
        <div
          v-for="t in toasts"
          :key="t.id"
          :class="cn('rounded-lg border p-4 shadow-lg cursor-pointer', variantClass(t.variant))"
          @click="dismiss(t.id)"
        >
          <p class="font-semibold text-sm">{{ t.title }}</p>
          <p v-if="t.description" class="text-xs mt-0.5 opacity-80">{{ t.description }}</p>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from { transform: translateX(100%); opacity: 0; }
.toast-leave-to   { transform: translateX(100%); opacity: 0; }
</style>
