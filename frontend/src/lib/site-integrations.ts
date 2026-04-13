import { siteCopy } from '../content/copy'

// 首页右上角小组件设置。
export interface SiteIntegrationsSettings {
  githubUsername: string
  weatherLocation: string
  showWeather: boolean
  showMusic: boolean
  showClock: boolean
}

export interface SiteIntegrationsSettingsPayload {
  github_username?: string
  weather_location?: string
  show_weather?: boolean
  show_music?: boolean
  show_clock?: boolean
}

function trimText(input: unknown): string {
  return typeof input === 'string' ? input.trim() : ''
}

function normalizeBoolean(input: unknown, fallback: boolean): boolean {
  return typeof input === 'boolean' ? input : fallback
}

export function getDefaultSiteIntegrationsSettings(): SiteIntegrationsSettings {
  // GitHub 用户名默认沿用 About 页配置，其余显示开关给出比较完整的初始体验。
  return {
    githubUsername: trimText(siteCopy.aboutPage.githubUsername),
    weatherLocation: '',
    showWeather: true,
    showMusic: true,
    showClock: true,
  }
}

export function normalizeSiteIntegrationsSettings(input: unknown): SiteIntegrationsSettings {
  // 同时兼容后端返回值、表单值和本地缓存值，最终统一成一种结构。
  const defaults = getDefaultSiteIntegrationsSettings()
  if (!input || typeof input !== 'object') {
    return defaults
  }

  const source = input as SiteIntegrationsSettingsPayload & Partial<SiteIntegrationsSettings> & Record<string, unknown>

  return {
    githubUsername: trimText(source.githubUsername ?? source.github_username) || defaults.githubUsername,
    weatherLocation: trimText(source.weatherLocation ?? source.weather_location),
    showWeather: normalizeBoolean(source.showWeather ?? source.show_weather, defaults.showWeather),
    showMusic: normalizeBoolean(source.showMusic ?? source.show_music, defaults.showMusic),
    showClock: normalizeBoolean(source.showClock ?? source.show_clock, defaults.showClock),
  }
}

export function toSiteIntegrationsPayload(settings: SiteIntegrationsSettings): SiteIntegrationsSettingsPayload {
  return {
    github_username: trimText(settings.githubUsername),
    weather_location: trimText(settings.weatherLocation),
    show_weather: settings.showWeather,
    show_music: settings.showMusic,
    show_clock: settings.showClock,
  }
}
