import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultSiteIntegrationsSettings,
  normalizeSiteIntegrationsSettings,
  toSiteIntegrationsPayload,
  type SiteIntegrationsSettings,
  type SiteIntegrationsSettingsPayload,
} from '../lib/site-integrations'

// 首页右上角的小组件开关（天气、音乐、时钟）统一放在这里管理。
const STORAGE_KEY = 'miku_site_integrations_settings'

function readCachedSiteIntegrationsSettings(): SiteIntegrationsSettings | null {
  if (typeof window === 'undefined') {
    return null
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    return normalizeSiteIntegrationsSettings(JSON.parse(raw))
  } catch {
    return null
  }
}

function writeCachedSiteIntegrationsSettings(settings: SiteIntegrationsSettings) {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(toSiteIntegrationsPayload(settings)))
  } catch {
    // ignore cache write failures
  }
}

export const siteIntegrationsSettings = atom<SiteIntegrationsSettings>(getDefaultSiteIntegrationsSettings())

// 与 `homeHero.ts` 同样的策略：
// - `hydrated` 避免重复请求
// - `hydrationPromise` 避免并发请求重复打到后端
let hydrated = false
let hydrationPromise: Promise<SiteIntegrationsSettings> | null = null

export function primeSiteIntegrationsSettingsFromCache(): SiteIntegrationsSettings {
  const cached = readCachedSiteIntegrationsSettings()
  if (cached) {
    siteIntegrationsSettings.set(cached)
    return cached
  }

  return siteIntegrationsSettings.get()
}

export async function hydrateSiteIntegrationsSettings(force = false): Promise<SiteIntegrationsSettings> {
  if (typeof window === 'undefined') {
    return siteIntegrationsSettings.get()
  }

  if (!force && hydrated) {
    return siteIntegrationsSettings.get()
  }

  primeSiteIntegrationsSettingsFromCache()

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<SiteIntegrationsSettingsPayload | undefined>('/site-settings/site-integrations')
    .then((data) => {
      // 后台配置最终会整理成前端固定使用的 camelCase 结构。
      const normalized = normalizeSiteIntegrationsSettings(data)
      siteIntegrationsSettings.set(normalized)
      writeCachedSiteIntegrationsSettings(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      // 小组件设置失败时不阻断页面，只保留当前已有状态。
      const fallback = siteIntegrationsSettings.get()
      siteIntegrationsSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveSiteIntegrationsSettings(next: SiteIntegrationsSettings): Promise<SiteIntegrationsSettings> {
  const normalized = normalizeSiteIntegrationsSettings(next)
  const saved = await api.put<SiteIntegrationsSettingsPayload>('/admin/site-settings/site-integrations', toSiteIntegrationsPayload(normalized))
  const finalSettings = normalizeSiteIntegrationsSettings(saved)
  siteIntegrationsSettings.set(finalSettings)
  writeCachedSiteIntegrationsSettings(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetSiteIntegrationsSettings(): Promise<SiteIntegrationsSettings> {
  return saveSiteIntegrationsSettings(getDefaultSiteIntegrationsSettings())
}
