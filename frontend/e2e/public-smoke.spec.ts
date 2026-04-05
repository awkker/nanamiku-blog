import { expect, test } from '@playwright/test'

import {
  createSmokeText,
  getPublishedSmokePost,
  hasAdminCredentials,
  loginAsAdmin,
  requireBackend,
} from './smoke-helpers'

test.describe('public smoke', () => {
  test.beforeEach(async ({ context }) => {
    await context.clearCookies()
  })

  test('blog detail loads and comment submission succeeds', async ({ page, request }) => {
    await requireBackend(request)
    const post = await getPublishedSmokePost(request)

    await page.goto(`/blog/${encodeURIComponent(post.slug)}`)

    await expect(page.getByRole('heading', { level: 1 })).toContainText(post.title)
    await expect(page.getByRole('heading', { name: '发布评论' })).toBeVisible()

    await page.getByLabel('评论昵称').fill(createSmokeText('blog-comment-author'))
    await page.getByLabel('评论邮箱').fill('smoke@example.com')
    await page.getByLabel('评论内容').fill(createSmokeText('blog-comment'))

    const submitResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/posts/')
      && response.url().includes('/comments')
      && response.request().method() === 'POST',
    )

    await page.getByRole('button', { name: '发送评论' }).click()

    await expect(await submitResponse).toBeOK()
    await expect(page.getByText('评论已提交，审核通过后会显示在列表中。')).toBeVisible()
  })

  test('blog comment can be approved in admin and becomes visible publicly', async ({ page, request }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running moderation smoke.')
    await requireBackend(request)
    const post = await getPublishedSmokePost(request)
    const author = createSmokeText('moderation-author')
    const commentContent = createSmokeText('moderation-comment')

    await page.goto(`/blog/${encodeURIComponent(post.slug)}`)
    await expect(page.getByRole('heading', { level: 1 })).toContainText(post.title)

    await page.getByLabel('评论昵称').fill(author)
    await page.getByLabel('评论邮箱').fill('smoke@example.com')
    await page.getByLabel('评论内容').fill(commentContent)

    const submitResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/posts/')
      && response.url().includes('/comments')
      && response.request().method() === 'POST',
    )

    await page.getByRole('button', { name: '发送评论' }).click()

    await expect(await submitResponse).toBeOK()
    await expect(page.getByText('评论已提交，审核通过后会显示在列表中。')).toBeVisible()

    await loginAsAdmin(page)
    await page.goto('/admin/comments')

    const pendingRow = page.locator('tbody tr').filter({ hasText: commentContent }).first()
    await expect(pendingRow).toBeVisible()
    await expect(pendingRow.getByText('待审核')).toBeVisible()

    const approveResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/admin/comments/')
      && response.url().includes('/approve')
      && response.request().method() === 'POST',
    )

    await pendingRow.getByRole('button', { name: '通过评论' }).click()

    await expect(await approveResponse).toBeOK()
    await expect(pendingRow.getByText('已通过')).toBeVisible()

    await page.goto(`/blog/${encodeURIComponent(post.slug)}`)
    const publicComment = page.locator('#post-comments').getByText(commentContent)
    await expect(publicComment).toBeVisible()

    await page.goto('/admin/comments')
    const approvedRow = page.locator('tbody tr').filter({ hasText: commentContent }).first()
    await expect(approvedRow).toBeVisible()

    const deleteResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/admin/comments/')
      && response.request().method() === 'DELETE',
    )

    await approvedRow.getByRole('button', { name: '删除评论' }).click()

    await expect(await deleteResponse).toBeOK()
    await expect(page.locator('tbody tr').filter({ hasText: commentContent })).toHaveCount(0)
  })

  test('guestbook submission succeeds', async ({ page, request }) => {
    await requireBackend(request)

    await page.goto('/guestbook')

    await expect(page.getByRole('heading', { name: '发布留言' })).toBeVisible()

    await page.getByLabel('昵称').fill('Playwright Guest')
    await page.getByLabel('留言内容').fill(createSmokeText('guestbook-message'))

    const submitResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/guestbook/messages')
      && response.request().method() === 'POST',
    )

    await page.getByRole('button', { name: '发布' }).click()

    await expect(await submitResponse).toBeOK()
    await expect(page.getByRole('status').filter({ hasText: '留言已提交，等待审核' })).toBeVisible()
  })
})
