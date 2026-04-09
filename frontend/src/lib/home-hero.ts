import { siteCopy } from '../content/copy'

export interface HomeHeroSettings {
  heroTitle: string
  heroSubtitle: string
}

export interface HomeHeroSettingsPayload {
  hero_title?: string
  hero_subtitle?: string
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

export function getDefaultHomeHeroSettings(): HomeHeroSettings {
  return {
    heroTitle: trimText(siteCopy.home.heroTitle),
    heroSubtitle: trimText(siteCopy.home.heroSubtitle),
  }
}

export function normalizeHomeHeroSettings(input: unknown): HomeHeroSettings {
  const defaults = getDefaultHomeHeroSettings()
  if (!input || typeof input !== 'object') {
    return defaults
  }

  const source = input as HomeHeroSettingsPayload & Partial<HomeHeroSettings> & Record<string, unknown>

  return {
    heroTitle: trimText(source.heroTitle ?? source.hero_title) || defaults.heroTitle,
    heroSubtitle: trimText(source.heroSubtitle ?? source.hero_subtitle) || defaults.heroSubtitle,
  }
}

export function toHomeHeroPayload(settings: HomeHeroSettings): HomeHeroSettingsPayload {
  return {
    hero_title: settings.heroTitle,
    hero_subtitle: settings.heroSubtitle,
  }
}
