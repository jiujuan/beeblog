import { ref } from 'vue'
import { articleApi } from '@/api/article.api'
import { useArticleStore } from '@/stores/article.store'
import { useAuthStore } from '@/stores/auth.store'
import type { Article } from '@/types/article.types'

export function useArticleDetail() {
  const store = useArticleStore()
  const authStore = useAuthStore()
  const article = ref<Article | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchDetail(slug: string) {
    // Use cache if available
    const cached = store.getCachedDetail(slug)
    if (cached) { article.value = cached; return }

    loading.value = true
    error.value = null
    try {
      const data = await articleApi.detail(slug)
      article.value = data
      store.cacheDetail(slug, data)
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function toggleLike() {
    if (!article.value || !authStore.isLoggedIn) return
    const id = article.value.id
    const liked = article.value.liked
    // Optimistic update
    article.value.liked = !liked
    article.value.like_count += liked ? -1 : 1
    store.updateInteraction(id, 'liked', !liked)
    try {
      liked ? await articleApi.unlike(id) : await articleApi.like(id)
    } catch {
      // Rollback
      article.value.liked = liked
      article.value.like_count += liked ? 1 : -1
      store.updateInteraction(id, 'liked', !!liked)
    }
  }

  async function toggleBookmark() {
    if (!article.value || !authStore.isLoggedIn) return
    const id = article.value.id
    const bookmarked = article.value.bookmarked
    article.value.bookmarked = !bookmarked
    article.value.bookmark_count += bookmarked ? -1 : 1
    store.updateInteraction(id, 'bookmarked', !bookmarked)
    try {
      bookmarked ? await articleApi.unbookmark(id) : await articleApi.bookmark(id)
    } catch {
      article.value.bookmarked = bookmarked
      article.value.bookmark_count += bookmarked ? 1 : -1
      store.updateInteraction(id, 'bookmarked', !!bookmarked)
    }
  }

  return { article, loading, error, fetchDetail, toggleLike, toggleBookmark }
}
