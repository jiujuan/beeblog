import { ref } from 'vue'
import { mediaApi } from '@/api/media.api'
import type { Media } from '@/types/media.types'

const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
const MAX_SIZE_MB = 10

export function useUpload() {
  const uploading = ref(false)
  const progress = ref(0)
  const error = ref<string | null>(null)

  async function upload(file: File, articleId?: number): Promise<Media | null> {
    error.value = null

    if (!ALLOWED_TYPES.includes(file.type)) {
      error.value = `不支持的文件类型：${file.type}`
      return null
    }
    if (file.size > MAX_SIZE_MB * 1024 * 1024) {
      error.value = `文件大小不能超过 ${MAX_SIZE_MB}MB`
      return null
    }

    uploading.value = true
    progress.value = 0
    try {
      const media = await mediaApi.upload(file, articleId)
      progress.value = 100
      return media
    } catch (e: any) {
      error.value = e.message
      return null
    } finally {
      uploading.value = false
    }
  }

  function selectFile(): Promise<File | null> {
    return new Promise((resolve) => {
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = ALLOWED_TYPES.join(',')
      input.onchange = () => resolve(input.files?.[0] ?? null)
      input.click()
    })
  }

  return { uploading, progress, error, upload, selectFile }
}
