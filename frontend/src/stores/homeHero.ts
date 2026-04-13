import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultHomeHeroSettings,
  normalizeHomeHeroSettings,
  toHomeHeroPayload,
  type HomeHeroSettings,
  type HomeHeroSettingsPayload,
} from '../lib/home-hero'

// 首页主视觉文案的本地缓存 key。
// 这个 store 采用“默认值 -> localStorage -> 后台接口”的三级回退策略。
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

// `hydrated` 表示当前标签页是否已经成功完成过一次“服务端同步”。
// `hydrationPromise` 用来做并发复用，防止多个组件同时请求同一个设置接口。
let hydrated = false
let hydrationPromise: Promise<HomeHeroSettings> | null = null

export function primeHomeHeroSettingsFromCache(): HomeHeroSettings {
  // 页面刚挂载时先做一次“预热”：
  // 有缓存就立刻用缓存，没有就先继续用默认值。
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

  // 默认只需要在首次真正请求后台；
  // 除非调用方显式传 `force=true` 才会重新拉取。
  if (!force && hydrated) {
    return homeHeroSettings.get()
  }

  primeHomeHeroSettingsFromCache()

  if (hydrationPromise) {
    // 已经有人发起请求时，直接复用同一个 Promise。
    return hydrationPromise
  }

  hydrationPromise = api
    .get<HomeHeroSettingsPayload | undefined>('/site-settings/home-hero')
    .then((data) => {
      // 后台数据可能是 snake_case，也可能缺字段，
      // 所以先统一 `normalize` 再写进 store 和缓存。
      const normalized = normalizeHomeHeroSettings(data)
      homeHeroSettings.set(normalized)
      writeCachedHomeHeroSettings(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      // 接口失败时不抛给首页，直接继续用当前 store 中已有的值。
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
  // 后台保存入口也统一走 normalize，保证发出去的是干净结构。
  const normalized = normalizeHomeHeroSettings(next)
  const saved = await api.put<HomeHeroSettingsPayload>('/admin/site-settings/home-hero', toHomeHeroPayload(normalized))
  const finalSettings = normalizeHomeHeroSettings(saved)
  homeHeroSettings.set(finalSettings)
  writeCachedHomeHeroSettings(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetHomeHeroSettings(): Promise<HomeHeroSettings> {
  // 重置本质上就是把“默认值”再保存一遍到后台。
  return saveHomeHeroSettings(getDefaultHomeHeroSettings())
}
