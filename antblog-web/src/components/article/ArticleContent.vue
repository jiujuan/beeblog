<script setup lang="ts">
import { computed } from 'vue'
import md from '@/utils/markdown'

const props = defineProps<{ content: string; html?: string }>()

function normalizeImageURLBlocks(markdown: string) {
  return markdown.replace(
    /^\s*(https?:\/\/[^\s]+\.(?:png|jpe?g|gif|webp|svg)(?:\?[^\s]*)?)\s*$/gim,
    '![]($1)',
  )
}

const rendered = computed(() => {
  const html = props.html?.trim() || ''
  if (html && /<\/?[a-z][\s\S]*>/i.test(html)) {
    return html
  }
  const source = props.content?.trim() || html
  return md.render(normalizeImageURLBlocks(source))
})
</script>

<template>
  <div
    class="markdown-body"
    v-html="rendered"
  />
</template>
