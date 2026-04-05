import { test, type APIRequestContext, type Page } from '@playwright/test'

export const adminIdentifier = process.env.SMOKE_ADMIN_IDENTIFIER || ''
export const adminPassword = process.env.SMOKE_ADMIN_PASSWORD || ''
export const hasAdminCredentials = adminIdentifier !== '' && adminPassword !== ''

export const backendURL = process.env.SMOKE_BACKEND_URL || 'http://127.0.0.1:8080'

const uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i

interface ApiEnvelope<T> {
  code: number
  message: string
  data: T
}

interface PagedData<T> {
  items: T[]
  total: number
  page: number
  size: number
}

interface SmokePost {
  id: string
  slug: string
  title: string
}

export async function requireBackend(request: APIRequestContext) {
  let available = false

  try {
    const response = await request.get(`${backendURL}/api/v1/health`, {
      timeout: 5_000,
    })
    available = response.ok()
  } catch {
    available = false
  }

  test.skip(!available, `Start backend before running smoke: ${backendURL}`)
}

export async function loginAsAdmin(page: Page) {
  await page.goto('/login')
  await page.getByLabel('用户名或邮箱').fill(adminIdentifier)
  await page.getByLabel('密码').fill(adminPassword)

  await Promise.all([
    page.waitForURL('**/admin'),
    page.getByRole('button', { name: '登录' }).click(),
  ])
}

export function createSmokeText(prefix: string) {
  return `${prefix}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

export async function getPublishedSmokePost(request: APIRequestContext): Promise<SmokePost> {
  const response = await request.get(`${backendURL}/api/v1/posts?page=1&size=10`, {
    timeout: 8_000,
  })

  test.skip(!response.ok(), `Unable to load published posts from backend: ${backendURL}`)

  const body = await response.json() as ApiEnvelope<PagedData<SmokePost>>
  const items = Array.isArray(body?.data?.items) ? body.data.items : []
  const post = items.find((item) => uuidPattern.test(item.id || '') && String(item.slug || '').trim() !== '')

  test.skip(!post, 'Seed at least one published CMS post with a UUID id before running blog smoke.')

  return post!
}
