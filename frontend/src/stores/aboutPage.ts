import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultAboutPageSettings,
  normalizeAboutPageSettings,
  toAboutPagePayload,
  type AboutPageSettings,
  type AboutPageSettingsPayload,
} from '../lib/about-page'

const STORAGE_KEY = 'miku_about_page_settings'

function readCachedAboutPageSettings(): AboutPageSettings | null {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    return normalizeAboutPageSettings(JSON.parse(raw))
  } catch {
    return null
  }
}

function writeCachedAboutPageSettings(settings: AboutPageSettings) {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(toAboutPagePayload(settings)))
  } catch {
    // ignore local cache failures
  }
}

export const aboutPageSettings = atom<AboutPageSettings>(getDefaultAboutPageSettings())

let hydrated = false
let hydrationPromise: Promise<AboutPageSettings> | null = null

export function primeAboutPageSettingsFromCache(): AboutPageSettings {
  const cached = readCachedAboutPageSettings()
  if (cached) {
    aboutPageSettings.set(cached)
    return cached
  }

  return aboutPageSettings.get()
}

export async function hydrateAboutPageSettings(force = false): Promise<AboutPageSettings> {
  if (typeof window === 'undefined') {
    return aboutPageSettings.get()
  }

  if (!force && hydrated) {
    return aboutPageSettings.get()
  }

  primeAboutPageSettingsFromCache()

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<AboutPageSettingsPayload | undefined>('/site-settings/about-page')
    .then((data) => {
      const normalized = normalizeAboutPageSettings(data)
      aboutPageSettings.set(normalized)
      writeCachedAboutPageSettings(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      const fallback = aboutPageSettings.get()
      aboutPageSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveAboutPageSettings(next: AboutPageSettings): Promise<AboutPageSettings> {
  const normalized = normalizeAboutPageSettings(next)
  const saved = await api.put<AboutPageSettingsPayload>('/admin/site-settings/about-page', toAboutPagePayload(normalized))
  const finalSettings = normalizeAboutPageSettings(saved)
  aboutPageSettings.set(finalSettings)
  writeCachedAboutPageSettings(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetAboutPageSettings(): Promise<AboutPageSettings> {
  return saveAboutPageSettings(getDefaultAboutPageSettings())
}
