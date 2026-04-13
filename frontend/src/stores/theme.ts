import { atom } from 'nanostores'

export type ThemeMode = 'light' | 'night'

// 主题模式同时要影响三处：
// 1. store 本身，供 Vue 组件读取
// 2. documentElement，上屏时立刻生效
// 3. localStorage，保证刷新后记住用户选择
const STORAGE_KEY = 'miku_blog_theme_mode'
const NIGHT_CLASS = 'theme-night'

export const themeMode = atom<ThemeMode>('light')

let hydrated = false
let storageBound = false

function normalizeThemeMode(input: string | null | undefined): ThemeMode {
  return input === 'night' ? 'night' : 'light'
}

function applyThemeToDocument(mode: ThemeMode) {
  if (typeof document === 'undefined') {
    return
  }

  const root = document.documentElement
  root.classList.toggle(NIGHT_CLASS, mode === 'night')
  root.setAttribute('data-theme', mode)
}

export function setThemeMode(next: ThemeMode, persist = true) {
  const mode = normalizeThemeMode(next)
  themeMode.set(mode)
  applyThemeToDocument(mode)

  if (typeof window !== 'undefined') {
    // 自定义事件用于通知像 `LiquidGlassCard` 这类需要立即响应主题变化的组件。
    window.dispatchEvent(new CustomEvent('miku-theme-change', {
      detail: { mode },
    }))
  }

  if (!persist || typeof window === 'undefined') {
    return
  }

  window.localStorage.setItem(STORAGE_KEY, mode)
}

export function toggleThemeMode() {
  const next: ThemeMode = themeMode.get() === 'night' ? 'light' : 'night'
  setThemeMode(next)
}

export function hydrateThemeMode() {
  if (typeof window === 'undefined') {
    return
  }

  // 首次挂载时先从 localStorage 恢复主题，再同步到 DOM 和 store。
  const saved = normalizeThemeMode(window.localStorage.getItem(STORAGE_KEY))
  setThemeMode(saved, false)

  if (hydrated) {
    return
  }
  hydrated = true

  if (!storageBound) {
    storageBound = true
    window.addEventListener('storage', (event) => {
      // 多标签页切换主题时，其他标签页也能跟着同步。
      if (event.key !== STORAGE_KEY) {
        return
      }
      setThemeMode(normalizeThemeMode(event.newValue), false)
    })
  }
}
