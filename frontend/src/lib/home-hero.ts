import { siteCopy } from '../content/copy'

// 首页主视觉设置的“前端使用结构”。
export interface HomeHeroSettings {
  heroTitle: string
  heroSubtitle: string
}

// 后端接口使用 snake_case，这里单独声明 payload 类型，方便做字段转换。
export interface HomeHeroSettingsPayload {
  hero_title?: string
  hero_subtitle?: string
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

export function getDefaultHomeHeroSettings(): HomeHeroSettings {
  // 默认值来自 copy，这样即使后端没配置，首页也仍然能直接工作。
  return {
    heroTitle: trimText(siteCopy.home.heroTitle),
    heroSubtitle: trimText(siteCopy.home.heroSubtitle),
  }
}

export function normalizeHomeHeroSettings(input: unknown): HomeHeroSettings {
  // normalize 的目的是把“各种可能长相的数据”整理成稳定的前端结构：
  // - 允许 camelCase / snake_case 混用
  // - 缺字段时回退到 copy 默认值
  // - 文本自动 trim
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
  // 保存到后端前，再统一转换回接口要求的 snake_case。
  return {
    hero_title: settings.heroTitle,
    hero_subtitle: settings.heroSubtitle,
  }
}
