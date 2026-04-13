const API_BASE = '/api/v1'
const AUTH_ME_PATH = `${API_BASE}/auth/me`
const AUTH_REFRESH_PATH = `${API_BASE}/auth/refresh`

// 这个事件名是前端“登录态广播”的统一入口。
// 登录、退出或 token 刷新后，页面可以监听它来刷新自身 UI。
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
  // 所有鉴权相关请求都必须把 cookie 带上。
  return {
    ...init,
    credentials: 'include',
  }
}

function shouldSkipRefresh(input: string): boolean {
  // 这些接口本身就处在登录流程里，不应该再递归触发 refresh。
  return input.includes('/auth/login') || input.includes('/auth/refresh') || input.includes('/auth/logout')
}

async function parseJsonEnvelope<T>(response: Response): Promise<ApiResponse<T> | null> {
  // 后端采用统一响应包结构 `{ code, message, data }`。
  // 这里先安全解析，避免遇到非 JSON 响应直接抛异常。
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
    // 多个请求同时 401 时，只发一次 refresh 请求，其余请求复用同一个 Promise。
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
  // 所有 API 请求的统一策略：
  // 先正常请求；
  // 如果发现 401，就尝试 refresh；
  // refresh 成功后再重试一次原请求。
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
    // 首页只关心“当前是否存在登录用户”，
    // 所以这里返回用户对象或 null 即可。
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

  // 给任意页面/组件一个统一的通知渠道。
  window.dispatchEvent(new CustomEvent(AUTH_STATE_CHANGED_EVENT))
}
