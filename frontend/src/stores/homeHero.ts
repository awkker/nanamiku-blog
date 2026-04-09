import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultHomeHeroSettings,
  normalizeHomeHeroSettings,
  toHomeHeroPayload,
  type HomeHeroSettings,
  type HomeHeroSettingsPayload,
} from '../lib/home-hero'

const STORAGE_KEY = 'miku_home_hero_settings'

function readCachedHomeHeroSettings(): HomeHeroSettings | null {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    return normalizeHomeHeroSettings(JSON.parse(raw))
  } catch {
    return null
  }
}

function writeCachedHomeHeroSettings(settings: HomeHeroSettings) {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(toHomeHeroPayload(settings)))
  } catch {
    // ignore local cache failures
  }
}

export const homeHeroSettings = atom<HomeHeroSettings>(getDefaultHomeHeroSettings())

let hydrated = false
let hydrationPromise: Promise<HomeHeroSettings> | null = null

export function primeHomeHeroSettingsFromCache(): HomeHeroSettings {
  const cached = readCachedHomeHeroSettings()
  if (cached) {
    homeHeroSettings.set(cached)
    return cached
  }

  return homeHeroSettings.get()
}

export async function hydrateHomeHeroSettings(force = false): Promise<HomeHeroSettings> {
  if (typeof window === 'undefined') {
    return homeHeroSettings.get()
  }

  if (!force && hydrated) {
    return homeHeroSettings.get()
  }

  primeHomeHeroSettingsFromCache()

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<HomeHeroSettingsPayload | undefined>('/site-settings/home-hero')
    .then((data) => {
      const normalized = normalizeHomeHeroSettings(data)
      homeHeroSettings.set(normalized)
      writeCachedHomeHeroSettings(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      const fallback = homeHeroSettings.get()
      homeHeroSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveHomeHeroSettings(next: HomeHeroSettings): Promise<HomeHeroSettings> {
  const normalized = normalizeHomeHeroSettings(next)
  const saved = await api.put<HomeHeroSettingsPayload>('/admin/site-settings/home-hero', toHomeHeroPayload(normalized))
  const finalSettings = normalizeHomeHeroSettings(saved)
  homeHeroSettings.set(finalSettings)
  writeCachedHomeHeroSettings(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetHomeHeroSettings(): Promise<HomeHeroSettings> {
  return saveHomeHeroSettings(getDefaultHomeHeroSettings())
}
