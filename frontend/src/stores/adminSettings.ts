import { atom } from 'nanostores'

export type AdminSettingsSectionKey =
  | 'site-profile'
  | 'home-hero'
  | 'home-assets'
  | 'blog-index'
  | 'author-profile'
  | 'site-integrations'
  | 'site-footer'
  | 'admin-profile'

export const adminSettingsSectionKeys: AdminSettingsSectionKey[] = [
  'site-profile',
  'home-hero',
  'home-assets',
  'blog-index',
  'author-profile',
  'site-integrations',
  'site-footer',
  'admin-profile',
]

export const adminSettingsSection = atom<AdminSettingsSectionKey>('site-profile')

export function isAdminSettingsSectionKey(value: string | null): value is AdminSettingsSectionKey {
  return typeof value === 'string' && adminSettingsSectionKeys.includes(value as AdminSettingsSectionKey)
}

export function readAdminSettingsSectionFromLocation(search: string | URLSearchParams): AdminSettingsSectionKey | null {
  const params = typeof search === 'string'
    ? new URLSearchParams(search.startsWith('?') ? search.slice(1) : search)
    : search

  const section = params.get('section')
  return isAdminSettingsSectionKey(section) ? section : null
}

export function hydrateAdminSettingsSectionFromWindow(): AdminSettingsSectionKey {
  if (typeof window === 'undefined') {
    return adminSettingsSection.get()
  }

  const next = readAdminSettingsSectionFromLocation(window.location.search) ?? 'site-profile'
  adminSettingsSection.set(next)
  return next
}

export function setAdminSettingsSection(section: AdminSettingsSectionKey) {
  adminSettingsSection.set(section)
}
