import { ref, computed } from 'vue'

export function usePagination(defaultPageSize = 10) {
  const page = ref(1)
  const pageSize = ref(defaultPageSize)
  const total = ref(0)

  const pageCount = computed(() => Math.ceil(total.value / pageSize.value) || 1)
  const hasNext = computed(() => page.value < pageCount.value)
  const hasPrev = computed(() => page.value > 1)

  function goTo(p: number) {
    page.value = Math.max(1, Math.min(p, pageCount.value))
  }
  function next() { if (hasNext.value) page.value++ }
  function prev() { if (hasPrev.value) page.value-- }
  function reset() { page.value = 1 }

  return { page, pageSize, total, pageCount, hasNext, hasPrev, goTo, next, prev, reset }
}
