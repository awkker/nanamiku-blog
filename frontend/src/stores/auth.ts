import { atom, computed } from 'nanostores'

import { api, ApiError } from '../lib/api'
import {
  emitAuthStateChanged,
  ensureSessionUser,
  type SessionUserPayload,
} from '../lib/auth-session'

// Centralized authentication state:
// - hydrate from HttpOnly cookie sessions
// - real backend login/logout flows
// - shared auth status for route guards and admin islands

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
  user: AuthUser | null
}

interface LoginResponse {
  expires_at: number
}

interface MeResponse extends SessionUserPayload {}

interface AccountResponse extends MeResponse {
  session_revoked?: boolean
}

export interface AccountUpdateResult {
  user: AuthUser
  sessionRevoked: boolean
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
  user: null,
})

export const isAuthenticated = computed(authState, (state) => state.status === 'authenticated')

let hydratePromise: Promise<AuthUser | null> | null = null
let hydrateVersion = 0

function isBrowser() {
  return typeof window !== 'undefined'
}

function setGuestState() {
  hydrateVersion += 1
  authState.set({ status: 'guest', user: null })
  emitAuthStateChanged()
}

function setAuthenticatedState(user: AuthUser) {
  hydrateVersion += 1
  authState.set({ status: 'authenticated', user })
  emitAuthStateChanged()
}

function commitHydrationGuest(version: number): AuthUser | null {
  if (version !== hydrateVersion) {
    return authState.get().user
  }

  authState.set({ status: 'guest', user: null })
  emitAuthStateChanged()
  return null
}

function commitHydrationUser(version: number, user: AuthUser): AuthUser {
  if (version !== hydrateVersion) {
    return authState.get().user ?? user
  }

  authState.set({ status: 'authenticated', user })
  emitAuthStateChanged()
  return user
}

export async function hydrateAuth(force = false): Promise<AuthUser | null> {
  if (!isBrowser()) {
    return null
  }

  const current = authState.get()
  if (!force && current.status === 'authenticated' && current.user) {
    return current.user
  }
  if (!force && hydratePromise) {
    return hydratePromise
  }

  const version = hydrateVersion + 1
  hydrateVersion = version
  authState.set({ status: 'checking', user: current.user })

  hydratePromise = (async () => {
    try {
      const me = await ensureSessionUser()
      if (!me) {
        return commitHydrationGuest(version)
      }

      return commitHydrationUser(version, toAuthUser(me))
    } catch {
      return commitHydrationGuest(version)
    } finally {
      if (version === hydrateVersion) {
        hydratePromise = null
      }
    }
  })()

  return hydratePromise
}

export async function loginWithPassword(identifier: string, password: string) {
  try {
    await api.post<LoginResponse>('/auth/login', {
      identifier: identifier.trim(),
      password: password.trim(),
    })

    const user = await hydrateAuth(true)
    if (!user) {
      throw new Error('会话初始化失败')
    }
    return user
  } catch (err) {
    if (err instanceof ApiError && err.status === 401) {
      throw new Error('账号或密码错误，请检查后重试。')
    }
    throw new Error('登录失败，请稍后重试。')
  }
}

export async function logout() {
  if (!isBrowser()) {
    return
  }

  try {
    await api.post('/auth/logout')
  } catch {
    // ignore logout API errors
  }

  setGuestState()
}

export async function updateMyProfile(displayName: string, avatarURL: string) {
  const me = await api.put<MeResponse>('/auth/me', {
    display_name: displayName,
    avatar_url: avatarURL,
  })

  const user = toAuthUser(me)
  setAuthenticatedState(user)
  return user
}

export async function updateMyAccount(username: string, email: string, newPassword: string) {
  const me = await api.put<AccountResponse>('/auth/account', {
    username,
    email,
    new_password: newPassword,
  })

  const user = toAuthUser(me)
  const sessionRevoked = Boolean(me.session_revoked)

  if (sessionRevoked) {
    setGuestState()
  } else {
    setAuthenticatedState(user)
  }

  return {
    user,
    sessionRevoked,
  } satisfies AccountUpdateResult
}
