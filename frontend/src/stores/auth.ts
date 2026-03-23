import { atom, computed } from 'nanostores'

import { api, ApiError } from '../lib/api'

// Centralized authentication state:
// - hydrate from localStorage
// - real backend login/logout flows
// - shared auth status for route guards and admin islands
const TOKEN_KEY = 'miku_blog_token'
const REFRESH_KEY = 'miku_blog_refresh'
const USER_KEY = 'miku_blog_user'

export interface AuthUser {
  id: string
  name: string
  username: string
  email: string
  role: 'admin'
  avatar: string
}

export interface AuthState {
  status: 'checking' | 'guest' | 'authenticated'
  token: string | null
  user: AuthUser | null
}

interface TokenPair {
  access_token: string
  refresh_token: string
  expires_at: number
}

interface MeResponse {
  id: string
  username: string
  display_name?: string
  avatar_url?: string
  email: string
  role: string
}

function normalizeAvatarURL(url?: string): string {
  return (url || '').trim() || '/picture/author.jpg'
}

function toAuthUser(me: MeResponse): AuthUser {
  const name = (me.display_name || '').trim() || me.username
  return {
    id: me.id,
    name,
    username: me.username,
    email: me.email,
    role: 'admin',
    avatar: normalizeAvatarURL(me.avatar_url),
  }
}

export const authState = atom<AuthState>({
  status: 'checking',
  token: null,
  user: null,
})

export const isAuthenticated = computed(authState, (state) => state.status === 'authenticated')

function isBrowser() {
  return typeof window !== 'undefined'
}

function persistAuth(token: string, refreshToken: string, user: AuthUser) {
  if (!isBrowser()) {
    return
  }

  try {
    window.localStorage.setItem(TOKEN_KEY, token)
    window.localStorage.setItem(REFRESH_KEY, refreshToken)
    window.localStorage.setItem(USER_KEY, JSON.stringify(user))
  } catch {
    throw new Error('浏览器阻止了本地存储，无法保持登录状态。请关闭无痕模式后重试。')
  }
}

function clearPersistedAuth() {
  if (!isBrowser()) {
    return
  }

  window.localStorage.removeItem(TOKEN_KEY)
  window.localStorage.removeItem(REFRESH_KEY)
  window.localStorage.removeItem(USER_KEY)
}

export function hydrateAuth() {
  if (!isBrowser()) {
    return
  }

  let token: string | null = null
  let rawUser: string | null = null

  try {
    token = window.localStorage.getItem(TOKEN_KEY)
    rawUser = window.localStorage.getItem(USER_KEY)
  } catch {
    authState.set({ status: 'guest', token: null, user: null })
    return
  }

  if (!token || !rawUser) {
    authState.set({ status: 'guest', token: null, user: null })
    return
  }

  try {
    const parsed = JSON.parse(rawUser) as Partial<AuthUser> & { username?: string }
    const user: AuthUser = {
      id: parsed.id || '',
      name: parsed.name || parsed.username || 'Admin',
      username: parsed.username || parsed.name || 'admin',
      email: parsed.email || '',
      role: 'admin',
      avatar: normalizeAvatarURL(parsed.avatar),
    }
    authState.set({ status: 'authenticated', token, user })
  } catch {
    clearPersistedAuth()
    authState.set({ status: 'guest', token: null, user: null })
  }
}

export async function loginWithPassword(identifier: string, password: string) {
  try {
    const pair = await api.post<TokenPair>('/auth/login', {
      username: identifier.trim(),
      password: password.trim(),
    })

    const me = await fetchMe(pair.access_token)

    const user = toAuthUser(me)

    persistAuth(pair.access_token, pair.refresh_token, user)
    authState.set({ status: 'authenticated', token: pair.access_token, user })

    return user
  } catch (err) {
    if (err instanceof ApiError && err.status === 401) {
      throw new Error('账号或密码错误，请检查后重试。')
    }
    throw new Error('登录失败，请稍后重试。')
  }
}

export async function logout() {
  if (!isBrowser()) return

  const refreshToken = window.localStorage.getItem(REFRESH_KEY)
  try {
    await api.post('/auth/logout', { refresh_token: refreshToken || '' })
  } catch {
    // ignore logout API errors
  }

  clearPersistedAuth()
  authState.set({ status: 'guest', token: null, user: null })
}

export function getStoredToken() {
  if (!isBrowser()) {
    return null
  }

  return window.localStorage.getItem(TOKEN_KEY)
}

async function fetchMe(token: string): Promise<MeResponse> {
  const res = await fetch('/api/v1/auth/me', {
    headers: { Authorization: `Bearer ${token}` },
    credentials: 'include',
  })
  const body = await res.json()
  if (!res.ok || body.code !== 0) throw new Error('fetch me failed')
  return body.data as MeResponse
}

export async function updateMyProfile(displayName: string, avatarURL: string) {
  const me = await api.put<MeResponse>('/auth/me', {
    display_name: displayName,
    avatar_url: avatarURL,
  })

  const user = toAuthUser(me)
  const current = authState.get()
  authState.set({
    status: current.status === 'authenticated' ? 'authenticated' : 'guest',
    token: current.token,
    user: current.status === 'authenticated' ? user : null,
  })

  if (isBrowser()) {
    try {
      window.localStorage.setItem(USER_KEY, JSON.stringify(user))
    } catch {
      // ignore storage write error
    }
  }

  return user
}
