const API_BASE = '/api/v1'
const AUTH_ME_PATH = `${API_BASE}/auth/me`
const AUTH_REFRESH_PATH = `${API_BASE}/auth/refresh`

export const AUTH_STATE_CHANGED_EVENT = 'miku-auth-state-changed'

interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
}

export interface SessionUserPayload {
  id: string
  username: string
  display_name?: string
  avatar_url?: string
  email: string
  role: string
}

let refreshPromise: Promise<boolean> | null = null

function withCredentials(init: RequestInit = {}): RequestInit {
  return {
    ...init,
    credentials: 'include',
  }
}

function shouldSkipRefresh(input: string): boolean {
  return input.includes('/auth/login') || input.includes('/auth/refresh') || input.includes('/auth/logout')
}

async function parseJsonEnvelope<T>(response: Response): Promise<ApiResponse<T> | null> {
  const contentType = response.headers.get('content-type') || ''
  if (!contentType.includes('application/json')) {
    return null
  }

  try {
    return await response.json() as ApiResponse<T>
  } catch {
    return null
  }
}

export async function refreshSessionCookies(): Promise<boolean> {
  if (typeof window === 'undefined') {
    return false
  }

  if (refreshPromise) {
    return refreshPromise
  }

  refreshPromise = (async () => {
    try {
      const response = await fetch(AUTH_REFRESH_PATH, withCredentials({
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: '{}',
      }))
      if (!response.ok) {
        return false
      }

      const body = await parseJsonEnvelope<unknown>(response)
      return Boolean(body && body.code === 0)
    } catch {
      return false
    }
  })()

  try {
    return await refreshPromise
  } finally {
    refreshPromise = null
  }
}

export async function fetchWithSessionRetry(
  input: string,
  init: RequestInit = {},
  retryOnUnauthorized = true,
): Promise<Response> {
  const response = await fetch(input, withCredentials(init))

  if (!retryOnUnauthorized || response.status !== 401 || shouldSkipRefresh(input)) {
    return response
  }

  const refreshed = await refreshSessionCookies()
  if (!refreshed) {
    return response
  }

  return fetch(input, withCredentials(init))
}

export async function ensureSessionUser(): Promise<SessionUserPayload | null> {
  try {
    const response = await fetchWithSessionRetry(AUTH_ME_PATH, { method: 'GET' })
    if (!response.ok) {
      return null
    }

    const body = await parseJsonEnvelope<SessionUserPayload>(response)
    if (!body || body.code !== 0 || !body.data) {
      return null
    }

    return body.data
  } catch {
    return null
  }
}

export function emitAuthStateChanged() {
  if (typeof window === 'undefined') {
    return
  }

  window.dispatchEvent(new CustomEvent(AUTH_STATE_CHANGED_EVENT))
}
