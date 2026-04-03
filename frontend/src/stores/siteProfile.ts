import { atom } from 'nanostores'

import { api } from '../lib/api'
import {
  getDefaultSiteProfileSettings,
  normalizeSiteProfileSettings,
  toSiteProfilePayload,
  type SiteProfileSettings,
  type SiteProfileSettingsPayload,
} from '../lib/site-profile'

export const siteProfileSettings = atom<SiteProfileSettings>(getDefaultSiteProfileSettings())

let hydrated = false
let hydrationPromise: Promise<SiteProfileSettings> | null = null

export async function hydrateSiteProfileSettings(force = false): Promise<SiteProfileSettings> {
  if (typeof window === 'undefined') {
    return siteProfileSettings.get()
  }

  if (!force && hydrated) {
    return siteProfileSettings.get()
  }

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<SiteProfileSettingsPayload | undefined>('/site-settings/site-profile')
    .then((data) => {
      const normalized = normalizeSiteProfileSettings(data)
      siteProfileSettings.set(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      const fallback = siteProfileSettings.get()
      siteProfileSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveSiteProfileSettings(next: SiteProfileSettings): Promise<SiteProfileSettings> {
  const normalized = normalizeSiteProfileSettings(next)
  const saved = await api.put<SiteProfileSettingsPayload>('/admin/site-settings/site-profile', toSiteProfilePayload(normalized))
  const finalSettings = normalizeSiteProfileSettings(saved)
  siteProfileSettings.set(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetSiteProfileSettings(): Promise<SiteProfileSettings> {
  return saveSiteProfileSettings(getDefaultSiteProfileSettings())
}
