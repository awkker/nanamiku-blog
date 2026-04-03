import { expect, test, type Page } from '@playwright/test'

const adminIdentifier = process.env.SMOKE_ADMIN_IDENTIFIER || ''
const adminPassword = process.env.SMOKE_ADMIN_PASSWORD || ''
const hasAdminCredentials = adminIdentifier !== '' && adminPassword !== ''

async function loginAsAdmin(page: Page) {
  await page.goto('/login')
  await expect(page.getByRole('button', { name: '登录' })).toBeVisible()
  await page.getByLabel('用户名或邮箱').fill(adminIdentifier)
  await page.getByLabel('密码').fill(adminPassword)

  await Promise.all([
    page.waitForURL('**/admin'),
    page.getByRole('button', { name: '登录' }).click(),
  ])
}

test.describe('admin smoke', () => {
  test.beforeEach(async ({ context }) => {
    await context.clearCookies()
  })

  test('guest user is redirected to login', async ({ page }) => {
    await page.goto('/admin')

    await expect(page).toHaveURL(/\/login$/)
    await expect(page.getByText('欢迎回来')).toBeVisible()
    await expect(page.getByRole('button', { name: '登录' })).toBeVisible()
  })

  test('admin can login, visit key pages, and logout', async ({ page }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running authenticated smoke.')

    await loginAsAdmin(page)

    await expect(page).toHaveURL(/\/admin$/)
    await expect(page.getByRole('heading', { name: '仪表盘' })).toBeVisible()

    await page.getByRole('link', { name: '友链管理' }).click()
    await expect(page).toHaveURL(/\/admin\/friends$/)
    await expect(page.getByText('申请队列')).toBeVisible()
    await expect(page.getByText('正式友链')).toBeVisible()

    await page.getByRole('link', { name: '数据备份' }).click()
    await expect(page).toHaveURL(/\/admin\/backup$/)
    await expect(page.getByRole('button', { name: '导出 JSON' })).toBeVisible()
    await expect(page.getByRole('button', { name: '导出 SQL' })).toBeVisible()

    await Promise.all([
      page.waitForURL('**/login'),
      page.getByRole('button', { name: '退出登录' }).click(),
    ])

    await page.goto('/admin')
    await expect(page).toHaveURL(/\/login$/)
  })
})
