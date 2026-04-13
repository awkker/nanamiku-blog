import { siteCopy } from '../content/copy'

// 站点级资料：品牌名、SEO 标题、默认描述、默认分享图等。
export interface SiteProfileSettings {
  brandText: string
  logoAlt: string
  siteTitle: string
  siteUrl: string
  defaultDescription: string
  defaultSocialImage: string
}

export interface SiteProfileSettingsPayload {
  brand_text?: string
  logo_alt?: string
  site_title?: string
  site_url?: string
  default_description?: string
  default_social_image?: string
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

function normalizeSiteURL(input: unknown): string {
  // 统一把站点地址整理成“可公开访问”的完整 URL。
  const trimmed = trimText(input)
  if (!trimmed) return ''
  if (trimmed.startsWith('//')) return `https:${trimmed}`
  if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) return trimmed.replace(/\/+$/, '')
  return `https://${trimmed.replace(/^\/+/, '').replace(/\/+$/, '')}`
}

function normalizeSiteAssetURL(input: unknown): string {
  // 站点资源允许写相对路径（走本站）或绝对 URL（走外部）。
  const trimmed = trimText(input)
  if (!trimmed) return ''
  if (trimmed.startsWith('/')) return trimmed
  return normalizeSiteURL(trimmed)
}

export function getDefaultSiteProfileSettings(): SiteProfileSettings {
  // 默认值全部来自 copy，保证仓库开箱即用。
  return {
    brandText: trimText(siteCopy.brand.text),
    logoAlt: trimText(siteCopy.brand.logoAlt),
    siteTitle: trimText(siteCopy.seo.siteTitle),
    siteUrl: normalizeSiteURL(siteCopy.seo.siteUrl),
    defaultDescription: trimText(siteCopy.seo.defaultDescription),
    defaultSocialImage: normalizeSiteAssetURL(siteCopy.seo.defaultSocialImage),
  }
}

export function getDefaultSiteProfileTitleTokens(): string[] {
  // 这些 token 会被运行时脚本拿去替换成后台新的站点标题。
  return Array.from(new Set([
    trimText(siteCopy.brand.text),
    trimText(siteCopy.seo.siteTitle),
    trimText(siteCopy.home.metaTitle),
    trimText(siteCopy.blogIndex.metaTitle),
  ].filter((item) => item.length > 0)))
}

export function normalizeSiteProfileSettings(input: unknown): SiteProfileSettings {
  // 统一兼容后台 snake_case、前端 camelCase 以及不完整输入。
  const defaults = getDefaultSiteProfileSettings()
  if (!input || typeof input !== 'object') {
    return defaults
  }

  const source = input as SiteProfileSettingsPayload & Partial<SiteProfileSettings> & Record<string, unknown>
  const brandText = trimText(source.brandText ?? source.brand_text) || defaults.brandText
  const siteTitle = trimText(source.siteTitle ?? source.site_title) || defaults.siteTitle

  return {
    brandText,
    logoAlt: trimText(source.logoAlt ?? source.logo_alt) || trimText(defaults.logoAlt) || `${brandText} logo`,
    siteTitle,
    siteUrl: normalizeSiteURL(source.siteUrl ?? source.site_url) || defaults.siteUrl,
    defaultDescription: trimText(source.defaultDescription ?? source.default_description) || defaults.defaultDescription,
    defaultSocialImage: normalizeSiteAssetURL(source.defaultSocialImage ?? source.default_social_image) || defaults.defaultSocialImage,
  }
}

export function toSiteProfilePayload(settings: SiteProfileSettings): SiteProfileSettingsPayload {
  return {
    brand_text: settings.brandText,
    logo_alt: settings.logoAlt,
    site_title: settings.siteTitle,
    site_url: settings.siteUrl,
    default_description: settings.defaultDescription,
    default_social_image: settings.defaultSocialImage,
  }
}
