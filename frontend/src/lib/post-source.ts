import { type CollectionEntry, getCollection } from 'astro:content'

import type { PostDetail, PostSummary } from './post-types'

type BlogEntry = CollectionEntry<'blog'>
type PostSourceKind = 'cms' | 'fallback'

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
  source: PostSourceKind
}

const DEFAULT_CMS_ORIGIN = 'http://127.0.0.1:8080'
const FALLBACK_HERO_IMAGE = '/picture/封面.avif'
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

function slugifyTag(name: string): string {
  return name
    .trim()
    .toLowerCase()
    .replace(/[^\p{Letter}\p{Number}\u4e00-\u9fff\s-]/gu, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

function toISODate(value?: Date): string {
  return (value ?? new Date()).toISOString()
}

function mapFallbackEntryToSummary(entry: BlogEntry): PostSummary {
  return {
    id: entry.id,
    slug: entry.id,
    title: entry.data.title,
    excerpt: entry.data.description,
    hero_image_url: entry.data.heroImage?.src || FALLBACK_HERO_IMAGE,
    category: entry.data.category || '技术随笔',
    published_at: toISODate(entry.data.pubDate),
    view_count: 0,
    like_count: 0,
    comment_count: 0,
    created_at: toISODate(entry.data.pubDate),
    tags: (entry.data.tags || []).map((name) => ({
      name,
      slug: slugifyTag(name),
    })),
  }
}

function mapFallbackEntryToDetail(entry: BlogEntry): PostDetail {
  const summary = mapFallbackEntryToSummary(entry)

  return {
    ...summary,
    content_markdown: entry.body || '',
    status: 'published',
    updated_at: toISODate(entry.data.updatedDate || entry.data.pubDate),
    liked: false,
  }
}

async function getFallbackEntries(): Promise<BlogEntry[]> {
  return (await getCollection('blog')).sort(
    (left, right) => right.data.pubDate.valueOf() - left.data.pubDate.valueOf(),
  )
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

    return {
      data: items,
      source: 'cms',
    }
  } catch {
    const fallbackEntries = await getFallbackEntries()
    return {
      data: fallbackEntries.map(mapFallbackEntryToSummary),
      source: 'fallback',
    }
  }
}

async function loadPublishedPostDetail(slug: string): Promise<SourceResult<PostDetail | null>> {
  try {
    const detail = await fetchCMS<PostDetail>(`/posts/${encodeURIComponent(slug)}`)
    return {
      data: detail,
      source: 'cms',
    }
  } catch {
    const fallbackEntries = await getFallbackEntries()
    const entry = fallbackEntries.find((item) => item.id === slug)
    return {
      data: entry ? mapFallbackEntryToDetail(entry) : null,
      source: 'fallback',
    }
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
