import { atom } from 'nanostores'

import { siteCopy } from '../content/copy'
import { api } from '../lib/api'

const MAX_CUSTOM_TEXTS = 8

export interface SiteFooterSettings {
  icpText: string
  icpLink: string
  customTexts: string[]
}

interface SiteFooterSettingsPayload {
  icp_text?: string
  icp_link?: string
  custom_texts?: string[]
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

function sanitizeCustomTexts(input: unknown): string[] {
  if (!Array.isArray(input)) {
    return []
  }

  return input
    .map((item) => trimText(item))
    .filter((item) => item.length > 0)
    .slice(0, MAX_CUSTOM_TEXTS)
}

function getDefaultSettings(): SiteFooterSettings {
  const defaults = siteCopy.siteFooter.defaults

  return {
    icpText: trimText(defaults.icpText),
    icpLink: trimText(defaults.icpLink),
    customTexts: sanitizeCustomTexts(defaults.customTexts),
  }
}

function normalizeSettings(input: unknown): SiteFooterSettings {
  if (!input || typeof input !== 'object') {
    return getDefaultSettings()
  }

  const defaults = getDefaultSettings()
  const source = input as Partial<SiteFooterSettings> &
    SiteFooterSettingsPayload &
    Record<string, unknown>
  const hasCustomTexts = Object.prototype.hasOwnProperty.call(source, 'customTexts')
    || Object.prototype.hasOwnProperty.call(source, 'custom_texts')

  return {
    icpText: trimText(source.icpText ?? source.icp_text) || defaults.icpText,
    icpLink: trimText(source.icpLink ?? source.icp_link) || defaults.icpLink,
    customTexts: hasCustomTexts ? sanitizeCustomTexts(source.customTexts ?? source.custom_texts) : defaults.customTexts,
  }
}

function toPayload(settings: SiteFooterSettings): SiteFooterSettingsPayload {
  return {
    icp_text: settings.icpText,
    icp_link: settings.icpLink,
    custom_texts: settings.customTexts,
  }
}

export const siteFooterSettings = atom<SiteFooterSettings>(getDefaultSettings())

let hydrated = false
let hydrationPromise: Promise<SiteFooterSettings> | null = null

export async function hydrateSiteFooterSettings(force = false): Promise<SiteFooterSettings> {
  if (typeof window === 'undefined') {
    return siteFooterSettings.get()
  }

  if (!force && hydrated) {
    return siteFooterSettings.get()
  }

  if (hydrationPromise) {
    return hydrationPromise
  }

  hydrationPromise = api
    .get<SiteFooterSettingsPayload | undefined>('/site-settings/footer')
    .then((data) => {
      const normalized = normalizeSettings(data)
      siteFooterSettings.set(normalized)
      hydrated = true
      return normalized
    })
    .catch(() => {
      const fallback = siteFooterSettings.get()
      siteFooterSettings.set(fallback)
      return fallback
    })
    .finally(() => {
      hydrationPromise = null
    })

  return hydrationPromise
}

export async function saveSiteFooterSettings(next: SiteFooterSettings): Promise<SiteFooterSettings> {
  const normalized = normalizeSettings(next)
  const saved = await api.put<SiteFooterSettingsPayload>('/admin/site-settings/footer', toPayload(normalized))
  const finalSettings = normalizeSettings(saved)
  siteFooterSettings.set(finalSettings)
  hydrated = true
  return finalSettings
}

export function resetSiteFooterSettings(): Promise<SiteFooterSettings> {
  return saveSiteFooterSettings(getDefaultSettings())
}
