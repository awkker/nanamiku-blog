import { siteCopy } from '../content/copy'
import { DEFAULT_PUBLIC_AUTHOR_AVATAR_URL, DEFAULT_PUBLIC_AVATAR_URL } from './default-assets'

const MAX_SKILLS = 8
const MAX_NOW_ITEMS = 8
const MAX_SOCIAL_LINKS = 6
const MAX_CONTACT_LINKS = 6
const LOCAL_ASSET_FILE_RE = /\.(avif|gif|jpe?g|png|svg|webp)$/i

export interface AuthorSocialLink {
  label: string
  href: string
  iconKey: string
}

export interface AuthorContactLink {
  label: string
  href: string
}

export interface AuthorProfileSettings {
  displayName: string
  avatarUrl: string
  role: string
  bio: string
  aboutDescription: string
  location: string
  since: string
  skills: string[]
  nowItems: string[]
  quote: string
  contactEmail: string
  socialLinks: AuthorSocialLink[]
  contactLinks: AuthorContactLink[]
}

export interface AuthorProfileSettingsPayload {
  display_name?: string
  avatar_url?: string
  role?: string
  bio?: string
  about_description?: string
  location?: string
  since?: string
  skills?: string[]
  now_items?: string[]
  quote?: string
  contact_email?: string
  social_links?: Array<{
    label?: string
    href?: string
    icon_key?: string
  }>
  contact_links?: Array<{
    label?: string
    href?: string
  }>
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

function normalizeAssetURL(input: unknown): string {
  const trimmed = trimText(input)
  if (!trimmed) return ''
  if (trimmed.startsWith('/')) return trimmed
  if (trimmed.startsWith('./') || trimmed.startsWith('../')) {
    return `/${trimmed.replace(/^(\.\.\/|\.\/)+/, '').replace(/^\/+/, '')}`
  }
  if (looksLikeLocalAssetPath(trimmed)) {
    return `/${trimmed.replace(/^\/+/, '')}`
  }
  if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) return trimmed
  return `https://${trimmed.replace(/^\/+/, '')}`
}

function looksLikeLocalAssetPath(input: string): boolean {
  if (!input.includes('/') || !LOCAL_ASSET_FILE_RE.test(input)) {
    return false
  }

  const [head] = input.split('/')
  return Boolean(head) && !head.includes('.')
}

function normalizeAuthorAvatarURL(input: unknown): string {
  const normalized = normalizeAssetURL(input)
  if (/^https?:\/\/picture\//.test(normalized) && LOCAL_ASSET_FILE_RE.test(normalized)) {
    return normalized.replace(/^https?:\/\/picture\//, '/picture/')
  }
  return normalized
}

export function resolveAuthorAvatarURL(input: unknown): string {
  return normalizeAuthorAvatarURL(input) || DEFAULT_PUBLIC_AUTHOR_AVATAR_URL
}

export function getAuthorAvatarFallbackChain(input: unknown): string[] {
  return [normalizeAuthorAvatarURL(input), DEFAULT_PUBLIC_AUTHOR_AVATAR_URL, DEFAULT_PUBLIC_AVATAR_URL]
    .filter((item, index, list): item is string => Boolean(item) && list.indexOf(item) === index)
}

function normalizePublicLink(input: unknown): string {
  const trimmed = trimText(input)
  if (!trimmed) return ''
  if (
    trimmed.startsWith('/') ||
    trimmed.startsWith('http://') ||
    trimmed.startsWith('https://') ||
    trimmed.startsWith('mailto:') ||
    trimmed.startsWith('tel:')
  ) {
    return trimmed
  }
  return `https://${trimmed.replace(/^\/+/, '')}`
}

function sanitizeTextItems(input: unknown, limit: number): string[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: string[] = []
  for (const item of input) {
    const normalized = trimText(item)
    if (!normalized || result.includes(normalized)) continue
    result.push(normalized)
    if (result.length >= limit) break
  }

  return result
}

function sanitizeSocialLinks(input: unknown): AuthorSocialLink[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: AuthorSocialLink[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const label = trimText(source.label)
    const href = normalizePublicLink(source.href)
    if (!label || !href) continue
    result.push({
      label,
      href,
      iconKey: trimText(source.iconKey ?? source.icon_key),
    })
    if (result.length >= MAX_SOCIAL_LINKS) break
  }
  return result
}

function sanitizeContactLinks(input: unknown): AuthorContactLink[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: AuthorContactLink[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const label = trimText(source.label)
    const href = normalizePublicLink(source.href)
    if (!label || !href) continue
    result.push({ label, href })
    if (result.length >= MAX_CONTACT_LINKS) break
  }
  return result
}

function deriveSocialIconKey(label: string): string {
  const normalized = label.trim().toLowerCase()
  if (normalized.includes('github')) return 'github'
  if (normalized === 'x' || normalized.includes('twitter')) return 'x'
  if (normalized.includes('mail') || normalized.includes('email')) return 'mail'
  if (normalized.includes('qq')) return 'qq'
  if (normalized.includes('bilibili')) return 'bilibili'
  return 'link'
}

function defaultContactEmail() {
  const fromHref = trimText(siteCopy.aboutPage.contactSection.emailHref).replace(/^mailto:/i, '')
  return fromHref.split('?')[0] || ''
}

export function extractGitHubUsernameFromURL(input: string): string {
  const href = trimText(input)
  if (!href) {
    return ''
  }

  try {
    const url = new URL(href.startsWith('http://') || href.startsWith('https://') ? href : `https://${href}`)
    if (!/(^|\.)github\.com$/i.test(url.hostname)) {
      return ''
    }

    const [username] = url.pathname
      .split('/')
      .map((segment) => segment.trim())
      .filter(Boolean)

    return username || ''
  } catch {
    return ''
  }
}

export function resolveGitHubUsername(preferred: string, socialLinks: AuthorSocialLink[]): string {
  const direct = trimText(preferred)
  if (direct) {
    return direct
  }

  for (const link of socialLinks) {
    const iconKey = trimText(link.iconKey).toLowerCase()
    const label = trimText(link.label).toLowerCase()
    if (iconKey !== 'github' && !label.includes('github')) {
      continue
    }

    const derived = extractGitHubUsernameFromURL(link.href)
    if (derived) {
      return derived
    }
  }

  return ''
}

export function getDefaultAuthorProfileSettings(): AuthorProfileSettings {
  return {
    displayName: trimText(siteCopy.blogIndex.authorCard.name) || trimText(siteCopy.aboutPage.profileCard.name),
    avatarUrl: resolveAuthorAvatarURL(''),
    role: trimText(siteCopy.blogIndex.authorCard.role) || trimText(siteCopy.aboutPage.profileCard.role),
    bio: trimText(siteCopy.blogIndex.authorCard.bio),
    aboutDescription: trimText(siteCopy.aboutPage.heroDescription),
    location: trimText(siteCopy.blogIndex.authorCard.location),
    since: trimText(siteCopy.blogIndex.authorCard.since),
    skills: sanitizeTextItems(siteCopy.blogIndex.authorCard.skills ?? siteCopy.aboutPage.identityTags, MAX_SKILLS),
    nowItems: sanitizeTextItems(siteCopy.blogIndex.nowItems, MAX_NOW_ITEMS),
    quote: trimText(siteCopy.aboutPage.profileCard.quote),
    contactEmail: defaultContactEmail(),
    socialLinks: sanitizeSocialLinks(
      siteCopy.blogIndex.socialLinks.map((item) => ({
        label: item.label,
        href: item.href,
        icon_key: deriveSocialIconKey(item.label),
      })),
    ),
    contactLinks: sanitizeContactLinks(siteCopy.aboutPage.socialLinks),
  }
}

export function normalizeAuthorProfileSettings(input: unknown): AuthorProfileSettings {
  const defaults = getDefaultAuthorProfileSettings()
  if (!input || typeof input !== 'object') {
    return defaults
  }

  const source = input as AuthorProfileSettingsPayload & Partial<AuthorProfileSettings> & Record<string, unknown>
  const socialLinks = sanitizeSocialLinks(source.socialLinks)
  const contactLinks = sanitizeContactLinks(source.contactLinks)

  return {
    displayName: trimText(source.displayName ?? source.display_name) || defaults.displayName,
    avatarUrl: resolveAuthorAvatarURL(source.avatarUrl ?? source.avatar_url),
    role: trimText(source.role) || defaults.role,
    bio: trimText(source.bio) || defaults.bio,
    aboutDescription: trimText(source.aboutDescription ?? source.about_description) || defaults.aboutDescription,
    location: trimText(source.location) || defaults.location,
    since: trimText(source.since) || defaults.since,
    skills: sanitizeTextItems(source.skills, MAX_SKILLS).length > 0 ? sanitizeTextItems(source.skills, MAX_SKILLS) : defaults.skills,
    nowItems: sanitizeTextItems(source.nowItems ?? source.now_items, MAX_NOW_ITEMS).length > 0
      ? sanitizeTextItems(source.nowItems ?? source.now_items, MAX_NOW_ITEMS)
      : defaults.nowItems,
    quote: trimText(source.quote) || defaults.quote,
    contactEmail: trimText(source.contactEmail ?? source.contact_email) || defaults.contactEmail,
    socialLinks: socialLinks.length > 0 ? socialLinks : defaults.socialLinks,
    contactLinks: contactLinks.length > 0 ? contactLinks : defaults.contactLinks,
  }
}

export function toAuthorProfilePayload(settings: AuthorProfileSettings): AuthorProfileSettingsPayload {
  return {
    display_name: settings.displayName,
    avatar_url: settings.avatarUrl,
    role: settings.role,
    bio: settings.bio,
    about_description: settings.aboutDescription,
    location: settings.location,
    since: settings.since,
    skills: sanitizeTextItems(settings.skills, MAX_SKILLS),
    now_items: sanitizeTextItems(settings.nowItems, MAX_NOW_ITEMS),
    quote: settings.quote,
    contact_email: trimText(settings.contactEmail),
    social_links: sanitizeSocialLinks(settings.socialLinks).map((item) => ({
      label: item.label,
      href: item.href,
      icon_key: item.iconKey,
    })),
    contact_links: sanitizeContactLinks(settings.contactLinks),
  }
}
