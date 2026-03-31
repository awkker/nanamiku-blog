import { atom } from 'nanostores'

export type ThemeMode = 'light' | 'night'

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

  const saved = normalizeThemeMode(window.localStorage.getItem(STORAGE_KEY))
  setThemeMode(saved, false)

  if (hydrated) {
    return
  }
  hydrated = true

  if (!storageBound) {
    storageBound = true
    window.addEventListener('storage', (event) => {
      if (event.key !== STORAGE_KEY) {
        return
      }
      setThemeMode(normalizeThemeMode(event.newValue), false)
    })
  }
}
