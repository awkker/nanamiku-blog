import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultHomeAssetsSettings,
  normalizeHomeAssetsSettings,
  toHomeAssetsPayload,
  type HomeAssetsSettings,
  type HomeAssetsSettingsPayload,
} from '../lib/home-assets'

// 首页背景资源配置 store。
// 当前只存 `heroImages`，但结构独立出来是为了以后继续扩展首页素材项。
const STORAGE_KEY = 'miku_home_assets_settings'

function readCachedHomeAssetsSettings(): HomeAssetsSettings | null {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    return normalizeHomeAssetsSettings(JSON.parse(raw))
  } catch {
    return null
  }
}

function writeCachedHomeAssetsSettings(settings: HomeAssetsSettings) {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(toHomeAssetsPayload(settings)))
  } catch {
    // ignore local cache failures
  }
}

export const homeAssetsSettings = atom<HomeAssetsSettings>(getDefaultHomeAssetsSettings())

let hydrated = false
let hydrationPromise: Promise<HomeAssetsSettings> | null = null

export function primeHomeAssetsSettingsFromCache(): HomeAssetsSettings {
  const cached = readCachedHomeAssetsSettings()
  if (cached) {
    homeAssetsSettings.set(cached)
    return cached
  }

  return homeAssetsSettings.get()
}

export async function hydrateHomeAssetsSettings(force = false): Promise<HomeAssetsSettings> {
  if (typeof window === 'undefined') {
    return homeAssetsSettings.get()
  }

  if (!force && hydrated) {
    return homeAssetsSettings.get()
  }

  primeHomeAssetsSettingsFromCache()

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<HomeAssetsSettingsPayload | undefined>('/site-settings/home-assets')
    .then((data) => {
      // 这里的 normalize 会顺便做图片链接合法化和去重。
      const normalized = normalizeHomeAssetsSettings(data)
      homeAssetsSettings.set(normalized)
      writeCachedHomeAssetsSettings(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      // 背景资源拉取失败时，直接保留默认图或缓存图即可。
      const fallback = homeAssetsSettings.get()
      homeAssetsSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveHomeAssetsSettings(next: HomeAssetsSettings): Promise<HomeAssetsSettings> {
  const normalized = normalizeHomeAssetsSettings(next)
  const saved = await api.put<HomeAssetsSettingsPayload>('/admin/site-settings/home-assets', toHomeAssetsPayload(normalized))
  const finalSettings = normalizeHomeAssetsSettings(saved)
  homeAssetsSettings.set(finalSettings)
  writeCachedHomeAssetsSettings(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetHomeAssetsSettings(): Promise<HomeAssetsSettings> {
  return saveHomeAssetsSettings(getDefaultHomeAssetsSettings())
}
