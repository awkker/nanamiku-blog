import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultAuthorProfileSettings,
  normalizeAuthorProfileSettings,
  toAuthorProfilePayload,
  type AuthorProfileSettings,
  type AuthorProfileSettingsPayload,
} from '../lib/author-profile'

export const authorProfileSettings = atom<AuthorProfileSettings>(getDefaultAuthorProfileSettings())

let hydrated = false
let hydrationPromise: Promise<AuthorProfileSettings> | null = null

export async function hydrateAuthorProfileSettings(force = false): Promise<AuthorProfileSettings> {
  if (typeof window === 'undefined') {
    return authorProfileSettings.get()
  }

  if (!force && hydrated) {
    return authorProfileSettings.get()
  }

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<AuthorProfileSettingsPayload | undefined>('/site-settings/author-profile')
    .then((data) => {
      const normalized = normalizeAuthorProfileSettings(data)
      authorProfileSettings.set(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      const fallback = authorProfileSettings.get()
      authorProfileSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveAuthorProfileSettings(next: AuthorProfileSettings): Promise<AuthorProfileSettings> {
  const normalized = normalizeAuthorProfileSettings(next)
  const saved = await api.put<AuthorProfileSettingsPayload>('/admin/site-settings/author-profile', toAuthorProfilePayload(normalized))
  const finalSettings = normalizeAuthorProfileSettings(saved)
  authorProfileSettings.set(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetAuthorProfileSettings(): Promise<AuthorProfileSettings> {
  return saveAuthorProfileSettings(getDefaultAuthorProfileSettings())
}
