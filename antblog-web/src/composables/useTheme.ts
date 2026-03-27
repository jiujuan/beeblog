import { ref, watchEffect } from 'vue'
import { storage } from '@/utils/storage'

type Theme = 'light' | 'dark' | 'system'

const theme = ref<Theme>(storage.get<Theme>('theme') ?? 'system')

function applyTheme(t: Theme) {
  const root = document.documentElement
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const isDark = t === 'dark' || (t === 'system' && prefersDark)
  root.classList.toggle('dark', isDark)
}

export function useTheme() {
  watchEffect(() => {
    applyTheme(theme.value)
    storage.set('theme', theme.value)
  })

  function setTheme(t: Theme) {
    theme.value = t
  }

  function toggleDark() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
  }

  return { theme, setTheme, toggleDark }
}
