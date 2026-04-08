/**
 * 博客文章数据源（单一真源：后端 CMS API）
 *
 * 所有构建期与运行时的文章列表、详情、RSS 等均由此模块驱动，
 * 不再 fallback 到 Astro Content Collections。
 */

import type { PostDetail, PostSummary } from './post-types'

interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
}

interface PagedData<T = unknown> {
  items: T[]
  total: number
  page: number
  size: number
}

interface SourceResult<T> {
  data: T
  source: 'cms'
}

const DEFAULT_CMS_ORIGIN = 'http://127.0.0.1:8080'
const CMS_PAGE_SIZE = 100

const cmsOrigin = (
  import.meta.env.CMS_API_ORIGIN
  || import.meta.env.PUBLIC_CMS_API_ORIGIN
  || process.env.CMS_API_ORIGIN
  || DEFAULT_CMS_ORIGIN
).replace(/\/$/, '')

const cmsAPIBase = `${cmsOrigin}/api/v1`

let publishedPostsCache: Promise<SourceResult<PostSummary[]>> | null = null
const postDetailCache = new Map<string, Promise<SourceResult<PostDetail | null>>>()

async function fetchCMS<T>(path: string): Promise<T> {
  const res = await fetch(`${cmsAPIBase}${path}`, {
    headers: {
      Accept: 'application/json',
    },
  })

  if (!res.ok) {
    throw new Error(`CMS request failed: ${res.status}`)
  }

  const body = await res.json() as ApiResponse<T>
  if (!body || body.code !== 0 || typeof body.data === 'undefined') {
    throw new Error(body?.message || 'CMS request failed')
  }

  return body.data
}

async function loadPublishedPostSummaries(): Promise<SourceResult<PostSummary[]>> {
  try {
    const firstPage = await fetchCMS<PagedData<PostSummary>>(`/posts?page=1&size=${CMS_PAGE_SIZE}`)
    const items = [...(firstPage.items || [])]
    const total = Number(firstPage.total || items.length)
    const totalPages = Math.max(1, Math.ceil(total / CMS_PAGE_SIZE))

    for (let currentPage = 2; currentPage <= totalPages; currentPage += 1) {
      const nextPage = await fetchCMS<PagedData<PostSummary>>(`/posts?page=${currentPage}&size=${CMS_PAGE_SIZE}`)
      items.push(...(nextPage.items || []))
    }

    return { data: items, source: 'cms' }
  } catch (err) {
    console.warn('[post-source] CMS unreachable, returning empty post list:', (err as Error).message)
    return { data: [], source: 'cms' }
  }
}

async function loadPublishedPostDetail(slug: string): Promise<SourceResult<PostDetail | null>> {
  try {
    const detail = await fetchCMS<PostDetail>(`/posts/${encodeURIComponent(slug)}`)
    return { data: detail, source: 'cms' }
  } catch (err) {
    console.warn(`[post-source] CMS unreachable for slug "${slug}":`, (err as Error).message)
    return { data: null, source: 'cms' }
  }
}

export async function getPublishedPostSummaries(): Promise<SourceResult<PostSummary[]>> {
  if (!publishedPostsCache) {
    publishedPostsCache = loadPublishedPostSummaries()
  }

  return publishedPostsCache
}

export async function getPublishedPostCount(): Promise<number> {
  const { data } = await getPublishedPostSummaries()
  return data.length
}

export async function getPublishedPostDetail(slug: string): Promise<SourceResult<PostDetail | null>> {
  if (!postDetailCache.has(slug)) {
    postDetailCache.set(slug, loadPublishedPostDetail(slug))
  }

  return postDetailCache.get(slug)!
}

export async function getRelatedPublishedPosts(slug: string, limit = 3): Promise<PostSummary[]> {
  const { data } = await getPublishedPostSummaries()
  return data
    .filter((item) => item.slug !== slug)
    .slice(0, limit)
}
