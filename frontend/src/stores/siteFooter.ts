import { atom } from 'nanostores'

import { siteCopy } from '../content/copy'

const STORAGE_KEY = 'miku_blog_site_footer'
const MAX_CUSTOM_TEXTS = 8

export interface SiteFooterSettings {
  icpText: string
  icpLink: string
  customTexts: string[]
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
  const source = input as Partial<SiteFooterSettings>

  return {
    icpText: trimText(source.icpText) || defaults.icpText,
    icpLink: trimText(source.icpLink) || defaults.icpLink,
    customTexts: sanitizeCustomTexts(source.customTexts),
  }
}

export const siteFooterSettings = atom<SiteFooterSettings>(getDefaultSettings())

let hydrated = false

export function hydrateSiteFooterSettings() {
  if (hydrated || typeof window === 'undefined') {
    return
  }

  hydrated = true

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return
    }

    const parsed = JSON.parse(raw)
    siteFooterSettings.set(normalizeSettings(parsed))
  } catch {
    siteFooterSettings.set(getDefaultSettings())
  }
}

export function saveSiteFooterSettings(next: SiteFooterSettings) {
  const normalized = normalizeSettings(next)
  siteFooterSettings.set(normalized)

  if (typeof window === 'undefined') {
    return
  }

  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(normalized))
}

export function resetSiteFooterSettings() {
  const defaults = getDefaultSettings()
  siteFooterSettings.set(defaults)

  if (typeof window === 'undefined') {
    return
  }

  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(defaults))
}
