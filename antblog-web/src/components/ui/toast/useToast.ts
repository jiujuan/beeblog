import { ref } from 'vue'

export interface ToastMessage {
  id: string
  title: string
  description?: string
  variant?: 'default' | 'destructive' | 'success'
  duration?: number
}

const toasts = ref<ToastMessage[]>([])

export function useToast() {
  function toast(msg: Omit<ToastMessage, 'id'>) {
    const id = Math.random().toString(36).slice(2)
    toasts.value.push({ ...msg, id })
    setTimeout(() => dismiss(id), msg.duration ?? 3000)
  }

  function dismiss(id: string) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  return { toasts, toast, dismiss }
}
