import {
  getDefaultSiteProfileSettings,
  getDefaultSiteProfileTitleTokens,
  normalizeSiteProfileSettings,
  type SiteProfileSettings,
  type SiteProfileSettingsPayload,
} from './site-profile'

const SITE_PROFILE_ENDPOINT = '/api/v1/site-settings/site-profile'
interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
}

declare global {
  interface Window {
    __mikuSiteProfile?: SiteProfileSettings
    __mikuSiteProfilePromise?: Promise<SiteProfileSettings>
  }
}

function resolvePublicURL(pathOrURL: string, siteURL: string): string {
  if (!pathOrURL) return ''
  try {
    return new URL(pathOrURL, siteURL).toString()
  } catch {
    return pathOrURL
  }
}

async function fetchSiteProfile(): Promise<SiteProfileSettings> {
  if (typeof window === 'undefined') {
    return getDefaultSiteProfileSettings()
  }

  if (window.__mikuSiteProfile) {
    return window.__mikuSiteProfile
  }

  if (window.__mikuSiteProfilePromise) {
    return window.__mikuSiteProfilePromise
  }

  window.__mikuSiteProfilePromise = fetch(SITE_PROFILE_ENDPOINT, {
    method: 'GET',
    credentials: 'include',
  })
    .then(async (response) => {
      if (!response.ok) {
        return getDefaultSiteProfileSettings()
      }

      const contentType = response.headers.get('content-type') || ''
      if (!contentType.includes('application/json')) {
        return getDefaultSiteProfileSettings()
      }

      const body = await response.json() as ApiResponse<SiteProfileSettingsPayload | null>
      return normalizeSiteProfileSettings(body.data)
    })
    .catch(() => getDefaultSiteProfileSettings())
    .finally(() => {
      window.__mikuSiteProfilePromise = undefined
    })

  const settings = await window.__mikuSiteProfilePromise
  window.__mikuSiteProfile = settings
  return settings
}

function updateTextContent(selector: string, value: string) {
  document.querySelectorAll<HTMLElement>(selector).forEach((node) => {
    node.textContent = value
  })
}

function updateAttribute(selector: string, attr: string, value: string) {
  document.querySelectorAll<HTMLElement>(selector).forEach((node) => {
    node.setAttribute(attr, value)
  })
}

function replaceTitleTokens(original: string, siteTitle: string): string {
  if (!original) return siteTitle

  const tokens = getDefaultSiteProfileTitleTokens()
  let updated = original
  let replaced = false

  for (const token of tokens) {
    if (!token || !updated.includes(token)) continue
    updated = updated.split(token).join(siteTitle)
    replaced = true
  }

  return replaced ? updated : original
}

function updateTitleLikeElement(selector: string, siteTitle: string) {
  document.querySelectorAll(selector).forEach((node) => {
    const original = node.getAttribute('data-site-original') || ''
    const next = replaceTitleTokens(original, siteTitle)

    if (node.tagName.toLowerCase() === 'title') {
      node.textContent = next
      return
    }

    node.setAttribute('content', next)
  })
}

function updateDescriptionLikeElement(selector: string, description: string) {
  document.querySelectorAll(selector).forEach((node) => {
    if (node.getAttribute('data-site-uses-default') !== 'true') return

    if (node.tagName.toLowerCase() === 'title') {
      node.textContent = description
      return
    }

    node.setAttribute('content', description)
  })
}

function updateImageMeta(selector: string, imageURL: string) {
  document.querySelectorAll(selector).forEach((node) => {
    if (node.getAttribute('data-site-uses-default') !== 'true') return
    node.setAttribute('content', imageURL)
  })
}

function updatePageURLs(siteURL: string) {
  const pageURL = resolvePublicURL(`${window.location.pathname}${window.location.search}`, siteURL)
  const rssURL = resolvePublicURL('/rss.xml', siteURL)

  updateAttribute('[data-site-canonical]', 'href', pageURL)
  updateAttribute('[data-site-og-url]', 'content', pageURL)
  updateAttribute('[data-site-twitter-url]', 'content', pageURL)
  updateAttribute('[data-site-rss-link]', 'href', rssURL)
}

function applySiteProfile(settings: SiteProfileSettings) {
  const socialImageURL = resolvePublicURL(settings.defaultSocialImage, settings.siteUrl)

  updateTextContent('[data-site-brand-text]', settings.brandText)
  updateAttribute('[data-site-logo-alt]', 'alt', settings.logoAlt)
  updateTitleLikeElement('title[data-site-original]', settings.siteTitle)
  updateTitleLikeElement('[data-site-title-meta]', settings.siteTitle)
  updateDescriptionLikeElement('[data-site-description-meta]', settings.defaultDescription)
  updateImageMeta('[data-site-image-meta]', socialImageURL)
  updatePageURLs(settings.siteUrl)
  updateAttribute('[data-site-rss-link]', 'title', settings.siteTitle)
}

export async function syncSiteProfileIntoDocument(): Promise<SiteProfileSettings> {
  const settings = await fetchSiteProfile()

  const apply = () => applySiteProfile(settings)
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', apply, { once: true })
  } else {
    apply()
  }

  return settings
}
