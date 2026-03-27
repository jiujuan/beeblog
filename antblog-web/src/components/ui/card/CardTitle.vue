<script setup lang="ts">
import { cn } from '@/utils/cn'
const props = defineProps<{ class?: string }>()
</script>
<template>
  <h3 :class="cn('font-serif font-semibold leading-none tracking-tight', props.class)"><slot /></h3>
</template>
