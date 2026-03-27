<script setup lang="ts">
import { inject, computed } from 'vue'
const props = defineProps<{ value: string }>()
const tabs = inject<any>('tabs')
const isActive = computed(() => tabs?.active.value === props.value)
</script>
<template>
  <div v-show="isActive" class="mt-4 focus-visible:outline-none"><slot /></div>
</template>
