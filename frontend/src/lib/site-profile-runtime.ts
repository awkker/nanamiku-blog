import {
  getDefaultSiteProfileSettings,
  getDefaultSiteProfileTitleTokens,
  normalizeSiteProfileSettings,
  type SiteProfileSettings,
  type SiteProfileSettingsPayload,
} from './site-profile'

/**
 * 运行时站点资料同步器
 *
 * `BaseHead.astro` 在构建期会先输出一套默认 SEO / 品牌信息，
 * 但后台设置可能在运行时把站点名、描述、分享图等改掉。
 *
 * 这个文件的作用就是：
 * 1. 去后台拿最新站点资料
 * 2. 在浏览器里把 document 上对应节点改成最新值
 *
 * 这样既保留了 SSG 的默认可用性，又支持后台动态配置。
 */
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
  // 支持把相对路径（如 `/favicon.png`）补成完整公开地址。
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
    // 已经拿过数据时直接复用缓存，避免每个页面都重复请求。
    return window.__mikuSiteProfile
  }

  if (window.__mikuSiteProfilePromise) {
    // 首次请求还在进行中时，其它调用直接等待同一个 Promise。
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
  // 统一封装 DOM 批量更新，便于下面复用。
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
  // 页面标题里可能预埋了默认站点名，比如：
  // `NanaMiku Blog | 某个页面`
  // 这里会把默认 token 替换成后台配置后的站点名。
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
    // 只替换“仍在使用默认描述”的节点，
    // 这样页面自己显式传入的 description 不会被站点默认值覆盖。
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
    // 同理，只覆盖默认分享图，不碰页面自己指定的文章封面图。
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
  // 这里统一把“站点级资料”同步到页面上所有打了 data-* 标记的节点。
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
    // 如果 DOM 还没准备好，就等 `DOMContentLoaded` 再写。
    document.addEventListener('DOMContentLoaded', apply, { once: true })
  } else {
    apply()
  }

  return settings
}
