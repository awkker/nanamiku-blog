import { expect, test, type APIRequestContext, type Page, type Response } from '@playwright/test'

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

interface SmokeMoment {
  id: string
  content: string
}

interface SmokeResource<T> {
  item: T
  cleanup: () => Promise<void>
}

const noopCleanup = async () => {}
let cachedAdminCookies: Awaited<ReturnType<ReturnType<Page['context']>['cookies']>> | null = null

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
  if (cachedAdminCookies && cachedAdminCookies.length > 0) {
    await page.context().addCookies(cachedAdminCookies)
    await page.goto('/admin')

    const hasSession = await page.evaluate(async () => {
      try {
        const response = await fetch('/api/v1/auth/me', {
          method: 'GET',
          credentials: 'include',
        })
        if (!response.ok) {
          return false
        }

        const body = await response.json().catch(() => null)
        return Boolean(body && body.code === 0 && body.data)
      } catch {
        return false
      }
    })

    if (hasSession) {
      await page.waitForURL(/\/admin(?:\/)?$/)
      return
    }

    cachedAdminCookies = null
    await page.context().clearCookies()
  }

  await page.goto('/login')
  await page.locator('input[autocomplete="username"]').fill(adminIdentifier)
  await page.locator('input[autocomplete="current-password"]').fill(adminPassword)

  const loginResponse = page.waitForResponse((response) =>
    response.url().includes('/api/v1/auth/login')
    && response.request().method() === 'POST',
  )

  await page.getByRole('button', { name: '登录' }).click()
  await expectPageResponseOK(loginResponse)

  await page.goto('/admin')
  await page.waitForURL(/\/admin(?:\/)?$/)
  cachedAdminCookies = await page.context().cookies()
}

export async function rememberAdminCookies(page: Page) {
  cachedAdminCookies = await page.context().cookies()
}

export function forgetCachedAdminCookies() {
  cachedAdminCookies = null
}

export async function expectPageResponseOK(responseOrPromise: Response | Promise<Response>) {
  const response = await responseOrPromise
  expect(
    response.ok(),
    `${response.request().method()} ${response.url()} -> ${response.status()}`,
  ).toBeTruthy()
  return response
}

export function createSmokeText(prefix: string) {
  return `${prefix}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

async function listPublishedSmokePosts(request: APIRequestContext) {
  const response = await request.get(`${backendURL}/api/v1/posts?page=1&size=10`, {
    timeout: 8_000,
  })

  if (!response.ok()) {
    throw new Error(`Unable to load published posts from backend: ${backendURL}`)
  }

  const body = await response.json() as ApiEnvelope<PagedData<SmokePost>>
  return Array.isArray(body?.data?.items) ? body.data.items : []
}

async function loginAdminAPI(request: APIRequestContext) {
  test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running authenticated smoke.')

  const response = await request.post(`${backendURL}/api/v1/auth/login`, {
    data: {
      identifier: adminIdentifier,
      password: adminPassword,
    },
    timeout: 8_000,
  })

  if (!response.ok()) {
    throw new Error(`Unable to login via admin API: ${backendURL}`)
  }
}

async function deleteAdminResource(request: APIRequestContext, path: string) {
  const response = await request.delete(`${backendURL}/api/v1${path}`, {
    timeout: 8_000,
  })

  if (!response.ok() && response.status() !== 404) {
    throw new Error(`Unable to cleanup smoke resource via ${path}: ${response.status()}`)
  }
}

export async function cleanupAdminResource(request: APIRequestContext, path: string) {
  await loginAdminAPI(request)
  await deleteAdminResource(request, path)
}

export async function acquirePublishedSmokePost(request: APIRequestContext): Promise<SmokeResource<SmokePost>> {
  const items = await listPublishedSmokePosts(request)
  const existing = items.find((item) => uuidPattern.test(item.id || '') && String(item.slug || '').trim() !== '')

  if (existing) {
    return { item: existing, cleanup: noopCleanup }
  }

  await loginAdminAPI(request)

  const slug = createSmokeText('playwright-post')
  const title = createSmokeText('playwright-post-title')
  const response = await request.post(`${backendURL}/api/v1/admin/posts`, {
    data: {
      slug,
      title,
      excerpt: 'Playwright smoke post',
      content_markdown: `# ${title}\n\nThis post is created automatically for smoke validation.`,
      hero_image_url: '/picture/logo-64.webp',
      category: 'Smoke',
      status: 'published',
      tags: ['smoke', 'playwright'],
    },
    timeout: 8_000,
  })

  if (!response.ok()) {
    throw new Error(`Unable to create published smoke post via admin API: ${backendURL}`)
  }

  const body = await response.json() as ApiEnvelope<{ id: string }>
  const post = {
    id: String(body?.data?.id || ''),
    slug,
    title,
  }

  if (!uuidPattern.test(post.id)) {
    throw new Error('Admin API did not return a valid smoke post id.')
  }

  return {
    item: post,
    // Keep the seeded post so retries and later tests can reuse the same routed slug.
    cleanup: noopCleanup,
  }
}

export async function createSmokeMoment(request: APIRequestContext): Promise<SmokeResource<SmokeMoment>> {
  const content = createSmokeText('playwright-moment-content')
  const response = await request.post(`${backendURL}/api/v1/moments`, {
    data: {
      author_name: createSmokeText('playwright-moment-author'),
      content,
      image_urls: [],
    },
    timeout: 8_000,
  })

  if (!response.ok()) {
    throw new Error(`Unable to create public smoke moment: ${backendURL}`)
  }

  const body = await response.json() as ApiEnvelope<{ id: string }>
  const moment = {
    id: String(body?.data?.id || ''),
    content,
  }

  if (!uuidPattern.test(moment.id)) {
    throw new Error('Public moments API did not return a valid smoke moment id.')
  }

  if (!hasAdminCredentials) {
    return { item: moment, cleanup: noopCleanup }
  }

  return {
    item: moment,
    cleanup: async () => {
      await cleanupAdminResource(request, `/admin/moments/${moment.id}`)
    },
  }
}
