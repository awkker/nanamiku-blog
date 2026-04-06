import { expect, test, type Page } from '@playwright/test'

import {
  acquirePublishedSmokePost,
  createSmokeMoment,
  createSmokeText,
  expectPageResponseOK,
  hasAdminCredentials,
  loginAsAdmin,
  requireBackend,
} from './smoke-helpers'

async function openSmokePost(page: Page, slug: string, title: string) {
  const commentsLoadResponse = page.waitForResponse((response) =>
    response.url().includes('/api/v1/posts/')
    && response.url().includes('/comments?page=1&size=20')
    && response.request().method() === 'GET',
  )
  await page.goto(`/blog/${encodeURIComponent(slug)}`)
  await expect(page.locator('header').getByRole('heading', { level: 1 }).first()).toContainText(title)
  await expectPageResponseOK(commentsLoadResponse)
}

test.describe('public smoke', () => {
  test.beforeEach(async ({ context }) => {
    await context.clearCookies()
  })

  test('blog detail loads and comment submission succeeds', async ({ page, request }) => {
    await requireBackend(request)
    const smokePost = await acquirePublishedSmokePost(request)
    const { item: post, cleanup } = smokePost

    try {
      await openSmokePost(page, post.slug, post.title)
      const commentsSection = page.locator('#post-comments')
      const commentForm = commentsSection.locator('form')
      await expect(commentsSection.getByRole('heading', { name: '发布评论' })).toBeVisible()

      await commentForm.getByLabel('评论昵称').fill(createSmokeText('blog-comment-author'))
      await commentForm.getByLabel('评论邮箱').fill('smoke@example.com')
      await commentForm.getByLabel('评论内容').fill(createSmokeText('blog-comment'))
      await expect(commentForm.getByRole('button', { name: '发送评论' })).toBeEnabled()

      const submitResponse = page.waitForResponse((response) =>
        response.url().includes('/api/v1/posts/')
        && response.url().includes('/comments')
        && response.request().method() === 'POST',
      )

      await commentForm.getByRole('button', { name: '发送评论' }).click()

      await expectPageResponseOK(submitResponse)
      await expect(page.getByText('评论已提交，审核通过后会显示在列表中。')).toBeVisible()
    } finally {
      await cleanup()
    }
  })

  test('blog comment can be approved in admin and becomes visible publicly', async ({ page, request }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running moderation smoke.')
    await requireBackend(request)
    const smokePost = await acquirePublishedSmokePost(request)
    const { item: post, cleanup } = smokePost
    const author = createSmokeText('moderation-author')
    const commentContent = createSmokeText('moderation-comment')

    try {
      await openSmokePost(page, post.slug, post.title)
      const commentsSection = page.locator('#post-comments')
      const commentForm = commentsSection.locator('form')

      await commentForm.getByLabel('评论昵称').fill(author)
      await commentForm.getByLabel('评论邮箱').fill('smoke@example.com')
      await commentForm.getByLabel('评论内容').fill(commentContent)
      await expect(commentForm.getByRole('button', { name: '发送评论' })).toBeEnabled()

      const submitResponse = page.waitForResponse((response) =>
        response.url().includes('/api/v1/posts/')
        && response.url().includes('/comments')
        && response.request().method() === 'POST',
      )

      await commentForm.getByRole('button', { name: '发送评论' }).click()

      await expectPageResponseOK(submitResponse)
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

      await expectPageResponseOK(approveResponse)
      await expect(pendingRow.getByText('已通过')).toBeVisible()

      await openSmokePost(page, post.slug, post.title)
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

      await expectPageResponseOK(deleteResponse)
      await expect(page.locator('tbody tr').filter({ hasText: commentContent })).toHaveCount(0)
    } finally {
      await cleanup()
    }
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

    await page.getByRole('button', { name: '发送留言' }).click()

    await expectPageResponseOK(submitResponse)
    await expect(page.getByRole('status').filter({ hasText: '留言已提交，等待审核' })).toBeVisible()
  })

  test('guestbook message can be approved in admin and becomes visible publicly', async ({ page, request }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running moderation smoke.')
    await requireBackend(request)

    const nickname = createSmokeText('guestbook-author')
    const content = createSmokeText('guestbook-message')

    await page.goto('/guestbook')
    await expect(page.getByRole('heading', { name: '发布留言' })).toBeVisible()

    await page.getByLabel('昵称').fill(nickname)
    await page.getByLabel('留言内容').fill(content)

    const submitResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/guestbook/messages')
      && response.request().method() === 'POST',
    )

    await page.getByRole('button', { name: '发送留言' }).click()

    await expectPageResponseOK(submitResponse)
    await expect(page.getByRole('status').filter({ hasText: '留言已提交，等待审核' })).toBeVisible()

    await loginAsAdmin(page)
    await page.goto('/admin/comments')
    await page.getByRole('button', { name: '留言板留言' }).click()

    const pendingRow = page.locator('tbody tr').filter({ hasText: content }).first()
    await expect(pendingRow).toBeVisible()
    await expect(pendingRow.getByText('待审核')).toBeVisible()

    const approveResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/admin/guestbook/messages/')
      && response.url().includes('/approve')
      && response.request().method() === 'POST',
    )

    await pendingRow.getByRole('button', { name: '通过评论' }).click()

    await expectPageResponseOK(approveResponse)
    await expect(pendingRow.getByText('已通过')).toBeVisible()

    await page.goto('/guestbook')
    await page.getByRole('button', { name: '最新' }).click()
    await expect(page.getByText(content)).toBeVisible()

    await page.goto('/admin/comments')
    await page.getByRole('button', { name: '留言板留言' }).click()
    const approvedRow = page.locator('tbody tr').filter({ hasText: content }).first()
    await expect(approvedRow).toBeVisible()

    const deleteResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/admin/guestbook/messages/')
      && response.request().method() === 'DELETE',
    )

    await approvedRow.getByRole('button', { name: '删除评论' }).click()

    await expectPageResponseOK(deleteResponse)
    await expect(page.locator('tbody tr').filter({ hasText: content })).toHaveCount(0)
  })

  test('friend application can be approved in admin and becomes visible publicly', async ({ page, request }) => {
    test.skip(!hasAdminCredentials, 'Set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD before running moderation smoke.')
    await requireBackend(request)

    const siteName = createSmokeText('friend-site')
    const siteURL = `https://${createSmokeText('friend')}.example.com`

    await page.goto('/friends')
    await page.getByRole('button', { name: '申请交换友链' }).click()

    await page.getByLabel('站点名称').fill(siteName)
    await page.getByLabel('站点地址').fill(siteURL)
    await page.getByLabel('联系邮箱').fill('smoke@example.com')
    await page.getByLabel('站点简介').fill(createSmokeText('friend-description'))

    const submitResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/friends/applications')
      && response.request().method() === 'POST',
    )

    await page.getByRole('button', { name: '提交申请' }).click()

    await expectPageResponseOK(submitResponse)
    await expect(page.getByText('友链申请已提交，审核通过后会出现在友链墙。')).toBeVisible()

    await loginAsAdmin(page)
    await page.goto('/admin/friends')

    const applicationRow = page.locator('tbody tr').filter({ hasText: siteURL }).first()
    await expect(applicationRow).toBeVisible()
    await expect(applicationRow.getByText('待审核')).toBeVisible()

    const approveResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/admin/friends/applications/')
      && response.url().includes('/approve')
      && response.request().method() === 'POST',
    )

    await applicationRow.getByRole('button', { name: '通过' }).click()

    await expectPageResponseOK(approveResponse)

    await page.goto('/friends')
    await expect(page.getByRole('heading', { name: siteName })).toBeVisible()

    await page.goto('/admin/friends')
    const friendRow = page.locator('table').last().locator('tbody tr').filter({ hasText: siteURL }).first()
    await expect(friendRow).toBeVisible()

    const deleteResponse = page.waitForResponse((response) =>
      response.url().includes('/api/v1/admin/friends/')
      && !response.url().includes('/applications/')
      && response.request().method() === 'DELETE',
    )

    await friendRow.getByRole('button', { name: '删除' }).click()

    await expectPageResponseOK(deleteResponse)
    await expect(page.locator('table').last().locator('tbody tr').filter({ hasText: siteURL })).toHaveCount(0)
  })

  test('moments interactions succeed publicly', async ({ page, request }) => {
    await requireBackend(request)
    const smokeMoment = await createSmokeMoment(request)
    const { item: moment, cleanup } = smokeMoment
    const commentContent = createSmokeText('moment-comment')

    try {
      await page.goto('/moments')

      const card = page.locator('article').filter({ hasText: moment.content }).first()
      await expect(card).toBeVisible()

      const likeResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/moments/${moment.id}/like`)
        && response.request().method() === 'POST',
      )
      await card.getByTestId('moment-like-button').click()
      await expectPageResponseOK(likeResponse)
      await expect(card.getByTestId('moment-like-button')).toContainText('1')

      const repostResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/moments/${moment.id}/repost`)
        && response.request().method() === 'POST',
      )
      await card.getByTestId('moment-repost-button').click()
      await expectPageResponseOK(repostResponse)
      await expect(card.getByTestId('moment-repost-button')).toContainText('1')

      await card.getByTestId('moment-comment-toggle').click()
      await card.getByPlaceholder('你的昵称').fill(createSmokeText('moment-comment-author'))
      await card.getByPlaceholder('写下评论...').fill(commentContent)

      const commentResponse = page.waitForResponse((response) =>
        response.url().includes(`/api/v1/moments/${moment.id}/comments`)
        && response.request().method() === 'POST',
      )
      await card.getByTestId('moment-comment-submit').click()
      await expectPageResponseOK(commentResponse)

      const commentItem = card.getByTestId('moment-comment-item').filter({ hasText: commentContent }).first()
      await expect(commentItem).toBeVisible()

      const commentLikeResponse = page.waitForResponse((response) =>
        response.url().includes('/api/v1/moments/comments/')
        && response.request().method() === 'POST',
      )
      await commentItem.getByTestId('moment-comment-like-button').click()
      await expectPageResponseOK(commentLikeResponse)
      await expect(commentItem.getByTestId('moment-comment-like-button')).toContainText('1')
    } finally {
      await cleanup()
    }
  })
})
