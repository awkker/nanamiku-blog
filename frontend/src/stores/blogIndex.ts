import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultBlogIndexSettings,
  normalizeBlogIndexSettings,
  toBlogIndexPayload,
  type BlogIndexSettings,
  type BlogIndexSettingsPayload,
} from '../lib/blog-index'

export const blogIndexSettings = atom<BlogIndexSettings>(getDefaultBlogIndexSettings())

let hydrated = false
let hydrationPromise: Promise<BlogIndexSettings> | null = null

export async function hydrateBlogIndexSettings(force = false): Promise<BlogIndexSettings> {
  if (typeof window === 'undefined') {
    return blogIndexSettings.get()
  }

  if (!force && hydrated) {
    return blogIndexSettings.get()
  }

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<BlogIndexSettingsPayload | undefined>('/site-settings/blog-index')
    .then((data) => {
      const normalized = normalizeBlogIndexSettings(data)
      blogIndexSettings.set(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      const fallback = blogIndexSettings.get()
      blogIndexSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveBlogIndexSettings(next: BlogIndexSettings): Promise<BlogIndexSettings> {
  const normalized = normalizeBlogIndexSettings(next)
  const saved = await api.put<BlogIndexSettingsPayload>('/admin/site-settings/blog-index', toBlogIndexPayload(normalized))
  const finalSettings = normalizeBlogIndexSettings(saved)
  blogIndexSettings.set(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetBlogIndexSettings(): Promise<BlogIndexSettings> {
  return saveBlogIndexSettings(getDefaultBlogIndexSettings())
}
