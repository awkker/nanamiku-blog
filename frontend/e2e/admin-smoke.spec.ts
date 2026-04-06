import { expect, test } from '@playwright/test'

import {
  cleanupAdminResource,
  createSmokeText,
  expectPageResponseOK,
  hasAdminCredentials,
  loginAsAdmin,
  requireBackend,
} from './smoke-helpers'

function formatLocalDateTimeInput(date: Date) {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60 * 1000)
  return local.toISOString().slice(0, 16)
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

  test('admin can login, visit key pages, and logout', async ({ page, request }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running authenticated smoke.')
    await requireBackend(request)

    await loginAsAdmin(page)

    await expect(page).toHaveURL(/\/admin(?:\/)?$/)
    await expect(page.getByRole('heading', { name: '仪表盘' })).toBeVisible()

    await page.getByRole('link', { name: '文章管理' }).click()
    await expect(page).toHaveURL(/\/admin\/posts(?:\/)?$/)
    await expect(page.getByTestId('admin-post-create-toggle')).toBeVisible()

    await page.getByRole('link', { name: '友链管理' }).click()
    await expect(page).toHaveURL(/\/admin\/friends(?:\/)?$/)
    await expect(page.getByText('申请队列')).toBeVisible()
    await expect(page.getByText('正式友链')).toBeVisible()

    await page.getByRole('link', { name: '数据备份' }).click()
    await expect(page).toHaveURL(/\/admin\/backup(?:\/)?$/)
    await expect(page.getByRole('button', { name: '导出 JSON' })).toBeVisible()
    await expect(page.getByRole('button', { name: '导出 SQL' })).toBeVisible()

    await Promise.all([
      page.waitForURL(/\/login(?:\/)?$/),
      page.getByRole('button', { name: '退出登录' }).click(),
    ])

    await page.goto('/admin')
    await expect(page).toHaveURL(/\/login(?:\/)?$/)
  })

  test('admin can create, publish, edit, draft, and delete a post', async ({ page, request }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running authenticated smoke.')
    await requireBackend(request)

    const originalTitle = createSmokeText('admin-post')
    const updatedTitle = createSmokeText('admin-post-updated')
    const originalContent = createSmokeText('admin-post-content')
    const updatedContent = createSmokeText('admin-post-content-updated')
    let postID = ''

    try {
      await loginAsAdmin(page)
      await page.goto('/admin/posts')
      await expect(page.getByTestId('admin-post-create-toggle')).toBeVisible()

      await page.getByTestId('admin-post-create-toggle').click()
      await expect(page.getByTestId('admin-post-create-form')).toBeVisible()

      await page.getByTestId('admin-post-create-title').fill(originalTitle)
      await page.getByPlaceholder('技术 / 随笔 / 教程').fill('Smoke')
      await page.getByPlaceholder('简短描述文章内容').fill(createSmokeText('admin-post-excerpt'))
      await page.getByTestId('admin-post-create-content').fill(originalContent)
      await page.getByTestId('admin-post-create-status').selectOption('draft')

      const createResponse = page.waitForResponse((response) =>
        response.url().includes('/api/v1/admin/posts')
        && response.request().method() === 'POST',
      )
      await page.getByTestId('admin-post-create-submit').click()

      const created = await createResponse
      await expectPageResponseOK(created)
      postID = String((await created.json()).data?.id || '')

      const draftRow = page.getByTestId('admin-post-row').filter({ hasText: originalTitle }).first()
      await expect(draftRow).toBeVisible()
      await expect(draftRow.getByText('草稿')).toBeVisible()

      await page.goto('/blog')
      await expect(page.getByText(originalTitle)).toHaveCount(0)

      await page.goto('/admin/posts')
      const publishRow = page.getByTestId('admin-post-row').filter({ hasText: originalTitle }).first()
      await publishRow.hover()

      const publishResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/posts/${postID}/publish`)
        && response.request().method() === 'POST',
      )
      await publishRow.getByTestId('admin-post-publish-button').click()

      await expectPageResponseOK(publishResponse)
      await expect(publishRow.getByText('已发布')).toBeVisible()

      await page.goto('/blog')
      await expect(page.getByText(originalTitle)).toBeVisible()

      await page.goto('/admin/posts')
      const editRow = page.getByTestId('admin-post-row').filter({ hasText: originalTitle }).first()
      await editRow.hover()
      await editRow.getByTestId('admin-post-edit-button').click()

      await expect(page.getByTestId('admin-post-edit-form')).toBeVisible()
      await page.getByTestId('admin-post-edit-title').fill(updatedTitle)
      await page.getByTestId('admin-post-edit-content').fill(updatedContent)

      const updateResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/posts/${postID}`)
        && response.request().method() === 'PUT',
      )
      await page.getByTestId('admin-post-edit-submit').click()

      await expectPageResponseOK(updateResponse)
      const updatedRow = page.getByTestId('admin-post-row').filter({ hasText: updatedTitle }).first()
      await expect(updatedRow).toBeVisible()

      await page.goto('/blog')
      await expect(page.getByText(updatedTitle)).toBeVisible()

      await page.goto('/admin/posts')
      const unpublishRow = page.getByTestId('admin-post-row').filter({ hasText: updatedTitle }).first()
      await unpublishRow.hover()

      const unpublishResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/posts/${postID}/unpublish`)
        && response.request().method() === 'POST',
      )
      await unpublishRow.getByTestId('admin-post-unpublish-button').click()

      await expectPageResponseOK(unpublishResponse)
      await expect(unpublishRow.getByText('草稿')).toBeVisible()

      await page.goto('/blog')
      await expect(page.getByText(updatedTitle)).toHaveCount(0)

      await page.goto('/admin/posts')
      const deleteRow = page.getByTestId('admin-post-row').filter({ hasText: updatedTitle }).first()
      await deleteRow.hover()

      const deleteResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/posts/${postID}`)
        && response.request().method() === 'DELETE',
      )
      await deleteRow.getByTestId('admin-post-delete-button').click()

      await expectPageResponseOK(deleteResponse)
      await expect(page.getByTestId('admin-post-row').filter({ hasText: updatedTitle })).toHaveCount(0)
      postID = ''
    } finally {
      if (postID) {
        await cleanupAdminResource(request, `/admin/posts/${postID}`)
      }
    }
  })

  test('admin can schedule, publish, edit, draft, and delete a moment', async ({ page, request }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running authenticated smoke.')
    await requireBackend(request)

    const originalContent = createSmokeText('admin-moment')
    const updatedContent = createSmokeText('admin-moment-updated')
    const scheduledAt = formatLocalDateTimeInput(new Date(Date.now() + 45 * 60 * 1000))
    let momentID = ''

    try {
      await loginAsAdmin(page)
      await page.goto('/admin/moments')
      await expect(page.getByTestId('admin-moment-create-toggle')).toBeVisible()

      await page.getByTestId('admin-moment-create-toggle').click()
      await expect(page.getByTestId('admin-moment-create-form')).toBeVisible()

      await page.getByTestId('admin-moment-create-content').fill(originalContent)

      const createResponse = page.waitForResponse((response) =>
        response.url().includes('/api/v1/admin/moments')
        && response.request().method() === 'POST',
      )
      await page.getByTestId('admin-moment-create-submit').click()

      const created = await createResponse
      await expectPageResponseOK(created)
      momentID = String((await created.json()).data?.id || '')

      const draftRow = page.getByTestId('admin-moment-row').filter({ hasText: originalContent }).first()
      await expect(draftRow).toBeVisible()
      await expect(draftRow.getByText('草稿')).toBeVisible()

      await page.goto('/moments')
      await expect(page.getByText(originalContent)).toHaveCount(0)

      await page.goto('/admin/moments')
      const scheduleRow = page.getByTestId('admin-moment-row').filter({ hasText: originalContent }).first()
      await scheduleRow.hover()

      page.once('dialog', async (dialog) => {
        await dialog.accept(scheduledAt)
      })
      const scheduleResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/moments/${momentID}/schedule`)
        && response.request().method() === 'POST',
      )
      await scheduleRow.getByTestId('admin-moment-schedule-button').click()

      await expectPageResponseOK(scheduleResponse)
      await expect(scheduleRow.getByText('定时发布')).toBeVisible()

      await page.goto('/moments')
      await expect(page.getByText(originalContent)).toHaveCount(0)

      await page.goto('/admin/moments')
      const publishRow = page.getByTestId('admin-moment-row').filter({ hasText: originalContent }).first()
      await publishRow.hover()

      const publishResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/moments/${momentID}/publish`)
        && response.request().method() === 'POST',
      )
      await publishRow.getByTestId('admin-moment-publish-button').click()

      await expectPageResponseOK(publishResponse)
      await expect(publishRow.getByText('已发布')).toBeVisible()

      await page.goto('/moments')
      await expect(page.getByText(originalContent)).toBeVisible()

      await page.goto('/admin/moments')
      const editRow = page.getByTestId('admin-moment-row').filter({ hasText: originalContent }).first()
      await editRow.hover()
      await editRow.getByTestId('admin-moment-edit-button').click()

      await expect(page.getByTestId('admin-moment-edit-form')).toBeVisible()
      await page.getByTestId('admin-moment-edit-content').fill(updatedContent)

      const updateResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/moments/${momentID}`)
        && response.request().method() === 'PUT',
      )
      await page.getByTestId('admin-moment-edit-submit').click()

      await expectPageResponseOK(updateResponse)
      const updatedRow = page.getByTestId('admin-moment-row').filter({ hasText: updatedContent }).first()
      await expect(updatedRow).toBeVisible()

      await page.goto('/moments')
      await expect(page.getByText(updatedContent)).toBeVisible()

      await page.goto('/admin/moments')
      const unpublishRow = page.getByTestId('admin-moment-row').filter({ hasText: updatedContent }).first()
      await unpublishRow.hover()

      const unpublishResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/moments/${momentID}/unpublish`)
        && response.request().method() === 'POST',
      )
      await unpublishRow.getByTestId('admin-moment-unpublish-button').click()

      await expectPageResponseOK(unpublishResponse)
      await expect(unpublishRow.getByText('草稿')).toBeVisible()

      await page.goto('/moments')
      await expect(page.getByText(updatedContent)).toHaveCount(0)

      await page.goto('/admin/moments')
      const deleteRow = page.getByTestId('admin-moment-row').filter({ hasText: updatedContent }).first()
      await deleteRow.hover()

      page.once('dialog', async (dialog) => {
        await dialog.accept()
      })
      const deleteResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/admin/moments/${momentID}`)
        && response.request().method() === 'DELETE',
      )
      await deleteRow.getByTestId('admin-moment-delete-button').click()

      await expectPageResponseOK(deleteResponse)
      await expect(page.getByTestId('admin-moment-row').filter({ hasText: updatedContent })).toHaveCount(0)
      momentID = ''
    } finally {
      if (momentID) {
        await cleanupAdminResource(request, `/admin/moments/${momentID}`)
      }
    }
  })
})
