import { siteCopy } from '../content/copy'

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
  return {
    githubUsername: trimText(siteCopy.aboutPage.githubUsername),
    weatherLocation: '',
    showWeather: true,
    showMusic: true,
    showClock: true,
  }
}

export function normalizeSiteIntegrationsSettings(input: unknown): SiteIntegrationsSettings {
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
