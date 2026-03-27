import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { articleApi } from '@/api/article.api'
import { useArticleStore } from '@/stores/article.store'
import { usePagination } from './usePagination'
import type { ListArticleQuery } from '@/types/article.types'

export function useArticleList(initialQuery?: Partial<ListArticleQuery>) {
  const store = useArticleStore()
  const { list } = storeToRefs(store)
  const { page, pageSize, total, pageCount, goTo, reset } = usePagination(10)

  const query = ref<ListArticleQuery>({
    page: 1,
    page_size: 10,
    ...initialQuery,
  })
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchList() {
    loading.value = true
    error.value = null
    try {
      const res = await articleApi.list({ ...query.value, page: page.value, page_size: pageSize.value })
      store.setList(res.list, res.total)
      total.value = res.total
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  function updateQuery(q: Partial<ListArticleQuery>) {
    query.value = { ...query.value, ...q }
    reset()
  }

  watch(page, fetchList)

  return {
    articles: list,
    total,
    page,
    pageSize,
    pageCount,
    loading,
    error,
    query,
    fetchList,
    updateQuery,
    goTo,
    reset,
  }
}
