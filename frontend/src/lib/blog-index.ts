import { siteCopy } from '../content/copy'

const MAX_HERO_ACTIONS = 3
const MAX_QUICK_STATS = 3

export interface BlogIndexHeroAction {
  label: string
  href: string
}

export interface BlogIndexQuickStat {
  label: string
  value: string
}

export interface BlogIndexFocusCard {
  badge: string
  title: string
  description: string
  footnote: string
}

export interface BlogIndexScrollCue {
  label: string
  ariaLabel: string
}

export interface BlogIndexSettings {
  heroBadge: string
  heroTitle: string
  heroDescription: string
  heroActions: BlogIndexHeroAction[]
  quickStats: BlogIndexQuickStat[]
  focusCard: BlogIndexFocusCard
  scrollCue: BlogIndexScrollCue
}

export interface BlogIndexSettingsPayload {
  hero_badge?: string
  hero_title?: string
  hero_description?: string
  hero_actions?: Array<{
    label?: string
    href?: string
  }>
  quick_stats?: Array<{
    label?: string
    value?: string
  }>
  focus_card?: {
    badge?: string
    title?: string
    description?: string
    footnote?: string
  }
  scroll_cue?: {
    label?: string
    aria_label?: string
  }
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

function normalizePublicLink(input: unknown): string {
  const trimmed = trimText(input)
  if (!trimmed) return ''
  if (
    trimmed.startsWith('/') ||
    trimmed.startsWith('#') ||
    trimmed.startsWith('http://') ||
    trimmed.startsWith('https://') ||
    trimmed.startsWith('mailto:') ||
    trimmed.startsWith('tel:')
  ) {
    return trimmed
  }
  return `https://${trimmed.replace(/^\/+/, '')}`
}

function sanitizeHeroActions(input: unknown): BlogIndexHeroAction[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: BlogIndexHeroAction[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const label = trimText(source.label)
    const href = normalizePublicLink(source.href)
    if (!label || !href) continue
    if (result.some((current) => current.label === label && current.href === href)) continue
    result.push({ label, href })
    if (result.length >= MAX_HERO_ACTIONS) break
  }

  return result
}

function sanitizeQuickStats(input: unknown): BlogIndexQuickStat[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: BlogIndexQuickStat[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const label = trimText(source.label)
    const value = trimText(source.value)
    if (!label || !value) continue
    if (result.some((current) => current.label === label && current.value === value)) continue
    result.push({ label, value })
    if (result.length >= MAX_QUICK_STATS) break
  }

  return result
}

export function getDefaultBlogIndexSettings(): BlogIndexSettings {
  return {
    heroBadge: trimText(siteCopy.blogIndex.heroBadge),
    heroTitle: trimText(siteCopy.blogIndex.heroTitle),
    heroDescription: trimText(siteCopy.blogIndex.heroDescription),
    heroActions: sanitizeHeroActions(siteCopy.blogIndex.heroActions),
    quickStats: sanitizeQuickStats(siteCopy.blogIndex.quickStats),
    focusCard: {
      badge: trimText(siteCopy.blogIndex.focusCard.badge),
      title: trimText(siteCopy.blogIndex.focusCard.title),
      description: trimText(siteCopy.blogIndex.focusCard.description),
      footnote: trimText(siteCopy.blogIndex.focusCard.footnote),
    },
    scrollCue: {
      label: trimText(siteCopy.blogIndex.scrollCue.label),
      ariaLabel: trimText(siteCopy.blogIndex.scrollCue.ariaLabel) || trimText(siteCopy.blogIndex.scrollCue.label),
    },
  }
}

export function normalizeBlogIndexSettings(input: unknown): BlogIndexSettings {
  const defaults = getDefaultBlogIndexSettings()
  if (!input || typeof input !== 'object') {
    return defaults
  }

  const source = input as BlogIndexSettingsPayload & Partial<BlogIndexSettings> & Record<string, unknown>
  const focusCardSource = (source.focusCard ?? source.focus_card ?? {}) as Record<string, unknown>
  const scrollCueSource = (source.scrollCue ?? source.scroll_cue ?? {}) as Record<string, unknown>
  const heroActions = sanitizeHeroActions(source.heroActions ?? source.hero_actions)
  const quickStats = sanitizeQuickStats(source.quickStats ?? source.quick_stats)
  const scrollCueLabel = trimText(scrollCueSource.label) || defaults.scrollCue.label

  return {
    heroBadge: trimText(source.heroBadge ?? source.hero_badge) || defaults.heroBadge,
    heroTitle: trimText(source.heroTitle ?? source.hero_title) || defaults.heroTitle,
    heroDescription: trimText(source.heroDescription ?? source.hero_description) || defaults.heroDescription,
    heroActions: heroActions.length > 0 ? heroActions : defaults.heroActions,
    quickStats: quickStats.length > 0 ? quickStats : defaults.quickStats,
    focusCard: {
      badge: trimText(focusCardSource.badge) || defaults.focusCard.badge,
      title: trimText(focusCardSource.title) || defaults.focusCard.title,
      description: trimText(focusCardSource.description) || defaults.focusCard.description,
      footnote: trimText(focusCardSource.footnote) || defaults.focusCard.footnote,
    },
    scrollCue: {
      label: scrollCueLabel,
      ariaLabel: trimText(scrollCueSource.ariaLabel ?? scrollCueSource.aria_label) || scrollCueLabel || defaults.scrollCue.ariaLabel,
    },
  }
}

export function toBlogIndexPayload(settings: BlogIndexSettings): BlogIndexSettingsPayload {
  return {
    hero_badge: trimText(settings.heroBadge),
    hero_title: trimText(settings.heroTitle),
    hero_description: trimText(settings.heroDescription),
    hero_actions: sanitizeHeroActions(settings.heroActions).map((item) => ({
      label: item.label,
      href: item.href,
    })),
    quick_stats: sanitizeQuickStats(settings.quickStats).map((item) => ({
      label: item.label,
      value: item.value,
    })),
    focus_card: {
      badge: trimText(settings.focusCard.badge),
      title: trimText(settings.focusCard.title),
      description: trimText(settings.focusCard.description),
      footnote: trimText(settings.focusCard.footnote),
    },
    scroll_cue: {
      label: trimText(settings.scrollCue.label),
      aria_label: trimText(settings.scrollCue.ariaLabel) || trimText(settings.scrollCue.label),
    },
  }
}
