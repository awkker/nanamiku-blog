import { siteCopy } from '../content/copy'

const MAX_INTRO_CARDS = 6
const MAX_MILESTONES = 8
const MAX_CAPABILITY_GROUPS = 6
const MAX_CAPABILITY_STACK_ITEMS = 8
const MAX_FEATURED_PROJECTS = 6
const MAX_MONTHLY_GOALS = 8
const MAX_LISTENING_ITEMS = 8

export interface AboutIntroCard {
  title: string
  description: string
}

export interface AboutMilestone {
  year: string
  title: string
  summary: string
  result: string
}

export interface AboutCapabilityGroup {
  title: string
  desc: string
  stack: string[]
}

export interface AboutFeaturedProject {
  name: string
  focus: string
  role: string
  metric: string
  href: string
}

export interface AboutSignatureSettings {
  description: string
  footer: string
}

export interface AboutPageSettings {
  introCards: AboutIntroCard[]
  milestones: AboutMilestone[]
  capabilityGroups: AboutCapabilityGroup[]
  featuredProjects: AboutFeaturedProject[]
  monthlyGoals: string[]
  listeningNow: string[]
  signature: AboutSignatureSettings
}

export interface AboutPageSettingsPayload {
  intro_cards?: Array<{
    title?: string
    description?: string
  }>
  milestones?: Array<{
    year?: string
    title?: string
    summary?: string
    result?: string
  }>
  capability_groups?: Array<{
    title?: string
    desc?: string
    stack?: string[]
  }>
  featured_projects?: Array<{
    name?: string
    focus?: string
    role?: string
    metric?: string
    href?: string
  }>
  monthly_goals?: string[]
  listening_now?: string[]
  signature?: {
    description?: string
    footer?: string
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

function sanitizeIntroCards(input: unknown): AboutIntroCard[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: AboutIntroCard[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const title = trimText(source.title)
    const description = trimText(source.description)
    if (!title || !description) continue
    if (result.some((current) => current.title === title && current.description === description)) continue
    result.push({ title, description })
    if (result.length >= MAX_INTRO_CARDS) break
  }
  return result
}

function sanitizeMilestones(input: unknown): AboutMilestone[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: AboutMilestone[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const year = trimText(source.year)
    const title = trimText(source.title)
    const summary = trimText(source.summary)
    const resultText = trimText(source.result)
    if (!year || !title || !summary || !resultText) continue
    if (result.some((current) => current.year === year && current.title === title)) continue
    result.push({
      year,
      title,
      summary,
      result: resultText,
    })
    if (result.length >= MAX_MILESTONES) break
  }
  return result
}

function sanitizeCapabilityGroups(input: unknown): AboutCapabilityGroup[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: AboutCapabilityGroup[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const title = trimText(source.title)
    const desc = trimText(source.desc)
    const stack = sanitizeTextItems(source.stack, MAX_CAPABILITY_STACK_ITEMS)
    if (!title || !desc || stack.length === 0) continue
    if (result.some((current) => current.title === title && current.desc === desc)) continue
    result.push({ title, desc, stack })
    if (result.length >= MAX_CAPABILITY_GROUPS) break
  }
  return result
}

function sanitizeFeaturedProjects(input: unknown): AboutFeaturedProject[] {
  if (!Array.isArray(input)) {
    return []
  }

  const result: AboutFeaturedProject[] = []
  for (const item of input) {
    if (!item || typeof item !== 'object') continue
    const source = item as Record<string, unknown>
    const name = trimText(source.name)
    const focus = trimText(source.focus)
    const role = trimText(source.role)
    const metric = trimText(source.metric)
    const href = normalizePublicLink(source.href)
    if (!name || !focus || !role || !metric || !href) continue
    if (result.some((current) => current.name === name && current.href === href)) continue
    result.push({ name, focus, role, metric, href })
    if (result.length >= MAX_FEATURED_PROJECTS) break
  }
  return result
}

function sanitizeSignature(input: unknown): AboutSignatureSettings {
  if (!input || typeof input !== 'object') {
    return { description: '', footer: '' }
  }

  const source = input as Record<string, unknown>
  return {
    description: trimText(source.description),
    footer: trimText(source.footer),
  }
}

export function getDefaultAboutPageSettings(): AboutPageSettings {
  return {
    introCards: sanitizeIntroCards(siteCopy.aboutPage.introCards),
    milestones: sanitizeMilestones(siteCopy.aboutPage.milestones),
    capabilityGroups: sanitizeCapabilityGroups(siteCopy.aboutPage.capabilityGroups),
    featuredProjects: sanitizeFeaturedProjects(siteCopy.aboutPage.featuredProjects),
    monthlyGoals: sanitizeTextItems(siteCopy.aboutPage.monthlyGoals, MAX_MONTHLY_GOALS),
    listeningNow: sanitizeTextItems(siteCopy.aboutPage.listeningNow, MAX_LISTENING_ITEMS),
    signature: sanitizeSignature(siteCopy.aboutPage.signature),
  }
}

export function normalizeAboutPageSettings(input: unknown): AboutPageSettings {
  const defaults = getDefaultAboutPageSettings()
  if (!input || typeof input !== 'object') {
    return defaults
  }

  const source = input as AboutPageSettingsPayload & Partial<AboutPageSettings> & Record<string, unknown>
  const introCards = sanitizeIntroCards(source.introCards ?? source.intro_cards)
  const milestones = sanitizeMilestones(source.milestones)
  const capabilityGroups = sanitizeCapabilityGroups(source.capabilityGroups ?? source.capability_groups)
  const featuredProjects = sanitizeFeaturedProjects(source.featuredProjects ?? source.featured_projects)
  const monthlyGoals = sanitizeTextItems(source.monthlyGoals ?? source.monthly_goals, MAX_MONTHLY_GOALS)
  const listeningNow = sanitizeTextItems(source.listeningNow ?? source.listening_now, MAX_LISTENING_ITEMS)
  const signature = sanitizeSignature(source.signature)

  return {
    introCards: introCards.length > 0 ? introCards : defaults.introCards,
    milestones: milestones.length > 0 ? milestones : defaults.milestones,
    capabilityGroups: capabilityGroups.length > 0 ? capabilityGroups : defaults.capabilityGroups,
    featuredProjects: featuredProjects.length > 0 ? featuredProjects : defaults.featuredProjects,
    monthlyGoals: monthlyGoals.length > 0 ? monthlyGoals : defaults.monthlyGoals,
    listeningNow: listeningNow.length > 0 ? listeningNow : defaults.listeningNow,
    signature: {
      description: signature.description || defaults.signature.description,
      footer: signature.footer || defaults.signature.footer,
    },
  }
}

export function toAboutPagePayload(settings: AboutPageSettings): AboutPageSettingsPayload {
  const normalized = normalizeAboutPageSettings(settings)
  return {
    intro_cards: normalized.introCards.map((item) => ({
      title: item.title,
      description: item.description,
    })),
    milestones: normalized.milestones.map((item) => ({
      year: item.year,
      title: item.title,
      summary: item.summary,
      result: item.result,
    })),
    capability_groups: normalized.capabilityGroups.map((item) => ({
      title: item.title,
      desc: item.desc,
      stack: sanitizeTextItems(item.stack, MAX_CAPABILITY_STACK_ITEMS),
    })),
    featured_projects: normalized.featuredProjects.map((item) => ({
      name: item.name,
      focus: item.focus,
      role: item.role,
      metric: item.metric,
      href: item.href,
    })),
    monthly_goals: sanitizeTextItems(normalized.monthlyGoals, MAX_MONTHLY_GOALS),
    listening_now: sanitizeTextItems(normalized.listeningNow, MAX_LISTENING_ITEMS),
    signature: {
      description: normalized.signature.description,
      footer: normalized.signature.footer,
    },
  }
}
