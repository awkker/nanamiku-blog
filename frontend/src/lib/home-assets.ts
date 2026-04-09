import { siteCopy } from '../content/copy'

const MAX_HERO_IMAGES = 8

export interface HomeAssetsSettings {
  heroImages: string[]
}

export interface HomeAssetsSettingsPayload {
  hero_images?: string[]
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

function normalizeAssetURL(input: unknown): string {
  const trimmed = trimText(input)
  if (!trimmed) return ''
  if (trimmed.startsWith('/')) return trimmed
  if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) return trimmed
  return `https://${trimmed.replace(/^\/+/, '')}`
}

function sanitizeHeroImages(input: unknown): string[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: string[] = []
  const seen = new Set<string>()

  for (const item of input) {
    const normalized = normalizeAssetURL(item)
    if (!normalized || seen.has(normalized)) continue
    seen.add(normalized)
    result.push(normalized)
    if (result.length >= MAX_HERO_IMAGES) break
  }

  return result
}

export function getDefaultHomeAssetsSettings(): HomeAssetsSettings {
  return {
    heroImages: sanitizeHeroImages(siteCopy.home.heroImages),
  }
}

export function normalizeHomeAssetsSettings(input: unknown): HomeAssetsSettings {
  const defaults = getDefaultHomeAssetsSettings()
  if (!input || typeof input !== 'object') {
    return defaults
  }

  const source = input as HomeAssetsSettingsPayload & Partial<HomeAssetsSettings> & Record<string, unknown>
  const hasHeroImages = Object.prototype.hasOwnProperty.call(source, 'heroImages')
    || Object.prototype.hasOwnProperty.call(source, 'hero_images')
  const heroImages = hasHeroImages ? sanitizeHeroImages(source.heroImages ?? source.hero_images) : defaults.heroImages

  return {
    heroImages: heroImages.length > 0 ? heroImages : defaults.heroImages,
  }
}

export function toHomeAssetsPayload(settings: HomeAssetsSettings): HomeAssetsSettingsPayload {
  return {
    hero_images: sanitizeHeroImages(settings.heroImages),
  }
}
