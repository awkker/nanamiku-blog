<template>
  <div class="space-y-7">
    <!-- Loading -->
    <div v-if="loading" class="space-y-6">
      <div v-for="i in 3" :key="i" class="animate-pulse rounded-3xl border border-white/55 bg-white/72 p-6">
        <div class="h-52 rounded-2xl bg-slate-200/60" />
        <div class="mt-4 h-5 w-2/3 rounded bg-slate-200/60" />
        <div class="mt-2 h-4 w-1/2 rounded bg-slate-100/80" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="rounded-2xl border border-red-200/60 bg-red-50/60 p-8 text-center">
      <p class="text-sm text-red-600">{{ error }}</p>
      <button type="button" class="mt-3 rounded-xl border border-red-200 bg-white px-4 py-2 text-xs text-red-600 transition hover:bg-red-50" @click="load">
        {{ copy.retry }}
      </button>
    </div>

    <template v-else>
      <section ref="searchPanelRef" class="relative z-40 isolate overflow-visible rounded-3xl border border-white/65 bg-white/78 p-5 shadow-[0_12px_30px_rgba(15,23,42,0.1)] backdrop-blur">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <p class="text-xs font-semibold tracking-[0.16em] text-miku">{{ copy.searchLabel }}</p>
            <p class="mt-1 text-xs text-slate-500">{{ copy.searchHint }}</p>
          </div>
          <div
            v-if="searchIndexLoading"
            class="inline-flex items-center gap-1.5 rounded-full border border-slate-200 bg-white/75 px-2.5 py-1 text-[11px] text-slate-500"
          >
            <LoadingSpinner size="sm" />
            <span>{{ copy.searchLoadingHint }}</span>
          </div>
        </div>

        <div class="relative mt-3">
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="copy.searchPlaceholder"
            class="w-full rounded-2xl border border-slate-200 bg-white/86 px-4 py-2.5 pr-16 text-sm text-slate-700 outline-none transition focus:border-miku/45 focus:ring-2 focus:ring-miku/15"
            @focus="handleSearchFocus"
            @input="handleSearchInput"
            @keydown="handleSearchKeydown"
          />
          <button
            v-if="searchQuery.trim().length > 0"
            type="button"
            class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg border border-slate-200 bg-white px-2.5 py-1 text-xs text-slate-500 transition hover:border-miku/35 hover:text-miku"
            @click="clearSearch"
          >
            {{ copy.searchClear }}
          </button>
        </div>

        <div
          v-if="showSuggestions"
          class="absolute left-5 right-5 top-[122px] z-[90] rounded-2xl border border-slate-200 bg-white/95 p-2 shadow-[0_14px_34px_rgba(15,23,42,0.12)] backdrop-blur"
        >
          <p class="px-2 py-1 text-[11px] font-semibold tracking-[0.12em] text-slate-500">{{ copy.searchSuggestionTitle }}</p>
          <ul class="space-y-1">
            <li v-for="(item, index) in suggestionItems" :key="item.id">
              <button
                type="button"
                class="flex w-full items-center justify-between rounded-xl px-2 py-2 text-left transition"
                :class="index === activeSuggestionIndex ? 'bg-miku-soft text-miku' : 'hover:bg-slate-50 text-slate-600'"
                @mousedown.prevent
                @click="selectSuggestion(item)"
              >
                <span class="truncate text-sm">{{ item.title }}</span>
                <span class="ml-2 shrink-0 text-[11px] text-slate-400">{{ item.category || '--' }}</span>
              </button>
            </li>
          </ul>
          <p class="mt-2 border-t border-slate-100 px-2 pt-2 text-[11px] text-slate-400">
            {{ copy.searchSuggestionEnterHint }}
          </p>
        </div>
      </section>

      <div v-if="isSearchMode" class="rounded-2xl border border-miku/25 bg-miku-soft/60 px-4 py-2 text-sm text-slate-600">
        {{ copy.searchResultPrefix }}
        <span class="mx-1 font-semibold text-miku">{{ visiblePosts.length }}</span>
        {{ copy.searchResultSuffix }}
      </div>

      <!-- Empty -->
      <div v-if="visiblePosts.length === 0" class="rounded-2xl border border-slate-200/60 bg-white/60 p-12 text-center">
        <p class="text-sm text-slate-500">{{ isSearchMode ? copy.searchEmpty : copy.empty }}</p>
      </div>

      <template v-else>
        <!-- Featured post -->
        <article v-if="visibleFeatured" class="blog-card group relative z-0 overflow-hidden rounded-3xl border border-white/65 bg-white/78 shadow-[0_20px_46px_rgba(15,23,42,0.14)] backdrop-blur">
          <a :href="`/blog/post?slug=${visibleFeatured.slug}`" class="block lg:grid lg:grid-cols-[1.2fr_1fr]">
            <div class="relative min-h-[260px] overflow-hidden lg:min-h-full">
              <img
                v-if="visibleFeatured.hero_image_url"
                :src="visibleFeatured.hero_image_url"
                :alt="visibleFeatured.title"
                class="h-full w-full object-cover transition duration-700 group-hover:scale-105"
                loading="lazy"
              />
              <div v-else class="flex h-full min-h-[260px] items-center justify-center bg-gradient-to-br from-miku/10 to-[#c084fc]/10">
                <svg viewBox="0 0 24 24" class="h-16 w-16 fill-none stroke-miku/30 stroke-[1]"><path d="M4 19.5A2.5 2.5 0 016.5 17H20" /><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" /></svg>
              </div>
              <div class="absolute inset-0 bg-gradient-to-t from-slate-900/35 via-transparent to-transparent" />
              <div class="absolute left-4 top-4 rounded-full border border-white/45 bg-white/22 px-3 py-1 text-xs font-semibold text-white backdrop-blur">
                {{ copy.featuredBadge }}
              </div>
            </div>

            <div class="space-y-4 p-5 md:p-6">
              <div class="flex flex-wrap items-center gap-2">
                <span v-if="visibleFeatured.category" class="inline-flex items-center gap-1 rounded-full border border-miku/35 bg-miku-soft px-3 py-1 text-xs font-semibold text-miku">
                  {{ visibleFeatured.category }}
                </span>
                <span v-if="visibleFeatured.tags && visibleFeatured.tags.length > 0" v-for="tag in visibleFeatured.tags.slice(0, 2)" :key="tag.name" class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] text-slate-500">
                  {{ tag.name }}
                </span>
              </div>

              <h2 class="text-2xl font-black leading-tight text-slate-900 transition group-hover:text-miku">
                {{ visibleFeatured.title }}
              </h2>
              <p class="text-sm leading-relaxed text-slate-600 md:text-base">
                {{ visibleFeatured.excerpt }}
              </p>

              <div class="grid gap-2 text-xs text-slate-500 sm:grid-cols-3">
                <span class="rounded-lg border border-slate-200 bg-white/75 px-2.5 py-1.5">
                  {{ visibleFeatured.view_count.toLocaleString() }}{{ copy.readSuffix }}
                </span>
                <span class="rounded-lg border border-slate-200 bg-white/75 px-2.5 py-1.5">
                  {{ visibleFeatured.like_count }}{{ copy.likeSuffix }}
                </span>
                <span class="rounded-lg border border-slate-200 bg-white/75 px-2.5 py-1.5">
                  {{ copy.publishedPrefix }}{{ formatDate(visibleFeatured.published_at || visibleFeatured.created_at) }}
                </span>
              </div>
            </div>
          </a>
        </article>

        <!-- Grid -->
        <div class="grid grid-cols-1 gap-7 md:grid-cols-2">
          <article
            v-for="(post, index) in visibleRestPosts"
            :key="post.id"
            :class="[
              'blog-card group relative z-0 overflow-hidden rounded-3xl border border-white/55 bg-white/72 shadow-lg backdrop-blur transition duration-500 hover:border-miku/25 hover:shadow-[0_16px_40px_rgba(57,197,187,0.12)]',
              index > 0 && index % 4 === 0 && 'md:col-span-2',
            ]"
          >
            <a
              :href="`/blog/post?slug=${post.slug}`"
              :class="[
                'block h-full',
                index > 0 && index % 4 === 0 && 'md:grid md:grid-cols-[320px_1fr] md:items-stretch',
              ]"
            >
              <div :class="['relative overflow-hidden', index > 0 && index % 4 === 0 ? 'h-56 md:h-full' : 'h-52']">
                <img
                  v-if="post.hero_image_url"
                  :src="post.hero_image_url"
                  :alt="post.title"
                  class="h-full w-full object-cover transition duration-700 group-hover:scale-105"
                  loading="lazy"
                />
                <div v-else class="flex h-full items-center justify-center bg-gradient-to-br from-miku/8 to-[#c084fc]/8">
                  <svg viewBox="0 0 24 24" class="h-12 w-12 fill-none stroke-miku/25 stroke-[1]"><path d="M4 19.5A2.5 2.5 0 016.5 17H20" /><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" /></svg>
                </div>
                <div class="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent opacity-0 transition duration-500 group-hover:opacity-100" />
              </div>

              <div class="space-y-3 px-5 py-5">
                <div class="flex flex-wrap items-center gap-2">
                  <span v-if="post.category" class="inline-flex items-center gap-1 rounded-full border border-miku/35 bg-miku-soft px-3 py-1 text-xs font-semibold text-miku">
                    {{ post.category }}
                  </span>
                  <span v-for="tag in (post.tags || []).slice(0, 2)" :key="tag.name" class="inline-flex rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] text-slate-500">
                    {{ tag.name }}
                  </span>
                </div>
                <h2 class="text-xl font-bold text-slate-900 transition duration-300 group-hover:text-miku">{{ post.title }}</h2>
                <p class="overflow-hidden text-sm leading-relaxed text-slate-500 [display:-webkit-box] [-webkit-box-orient:vertical] [-webkit-line-clamp:2]">
                  {{ post.excerpt }}
                </p>
                <div class="flex flex-wrap items-center gap-3 text-xs text-slate-400">
                  <span class="inline-flex items-center gap-1">
                    <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
                      <path d="M20 21V7a2 2 0 00-2-2h-3V3H9v2H6a2 2 0 00-2 2v14M16 11h-8M16 15h-8" />
                    </svg>
                    {{ formatDate(post.published_at || post.created_at) }}
                  </span>
                  <span class="h-0.5 w-0.5 rounded-full bg-slate-300" />
                  <span class="inline-flex items-center gap-1 text-slate-500">
                    <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.5]" aria-hidden="true"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" /><circle cx="12" cy="12" r="3" /></svg>
                    {{ post.view_count.toLocaleString() }}
                  </span>
                  <span class="h-0.5 w-0.5 rounded-full bg-slate-300" />
                  <span>{{ post.like_count }}{{ copy.shortLikeSuffix }}</span>
                </div>
              </div>
            </a>
          </article>
        </div>

        <!-- Pagination -->
        <div v-if="showPagination" class="flex items-center justify-center gap-2 pt-4">
          <button
            v-for="p in totalPages"
            :key="p"
            type="button"
            :class="[
              'rounded-xl border px-3.5 py-2 text-sm transition',
              p === page
                ? 'border-miku/45 bg-miku-soft text-miku'
                : 'border-slate-200 bg-white/60 text-slate-500 hover:border-miku/30 hover:text-miku',
            ]"
            @click="goToPage(p)"
          >
            {{ p }}
          </button>
        </div>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { api, type PagedData } from '../../lib/api'
import { siteCopy } from '../../content/copy'
import LoadingSpinner from '../ui/LoadingSpinner.vue'

interface TagItem {
  name: string
  slug: string
}

interface PostItem {
  id: string
  slug: string
  title: string
  excerpt: string
  hero_image_url: string
  category: string
  published_at?: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  tags?: TagItem[]
}

interface RankedPost {
  post: PostItem
  score: number
}

interface SearchDoc {
  title: string
  excerpt: string
  category: string
  tags: string
  merged: string
  words: string[]
}

const posts = ref<PostItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const error = ref('')
const copy = siteCopy.components.blogFeed

const searchQuery = ref('')
const searchFocused = ref(false)
const activeSuggestionIndex = ref(-1)
const searchIndexLoading = ref(false)
const searchIndexed = ref(false)
const searchPool = ref<PostItem[]>([])
const searchPanelRef = ref<HTMLElement | null>(null)

const SEARCH_INDEX_PAGE_SIZE = 60
const SEARCH_INDEX_MAX_PAGES = 20
const searchDocCache = new Map<string, SearchDoc>()

const normalizedSearchQuery = computed(() => normalizeForSearch(searchQuery.value))
const isSearchMode = computed(() => normalizedSearchQuery.value.length > 0)
const searchSourcePosts = computed(() => (searchPool.value.length > 0 ? searchPool.value : posts.value))
const rankedSearchPosts = computed<RankedPost[]>(() => {
  if (!isSearchMode.value) {
    return []
  }
  return rankPosts(searchSourcePosts.value, normalizedSearchQuery.value)
})

const visiblePosts = computed(() => (isSearchMode.value ? rankedSearchPosts.value.map((item) => item.post) : posts.value))
const visibleFeatured = computed(() => visiblePosts.value[0] || null)
const visibleRestPosts = computed(() => visiblePosts.value.slice(1))
const totalPages = computed(() => (isSearchMode.value ? 1 : Math.max(1, Math.ceil(total.value / pageSize))))
const showPagination = computed(() => !isSearchMode.value && totalPages.value > 1)
const suggestionItems = computed(() => rankedSearchPosts.value.slice(0, 6).map((item) => item.post))
const showSuggestions = computed(
  () => searchFocused.value
    && normalizedSearchQuery.value.length > 0
    && suggestionItems.value.length > 0,
)

function normalizeForSearch(input: string): string {
  return input
    .normalize('NFKC')
    .toLowerCase()
    .replace(/[^\p{L}\p{N}\s]/gu, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}

function tokenize(text: string): string[] {
  return normalizeForSearch(text).split(' ').filter(Boolean)
}

function parsePostTime(post: PostItem): number {
  const raw = post.published_at || post.created_at || ''
  const ts = Date.parse(raw)
  return Number.isNaN(ts) ? 0 : ts
}

function mergeUniquePosts(items: PostItem[]): PostItem[] {
  const merged = new Map<string, PostItem>()
  for (const item of items) {
    const key = item.id || item.slug || item.title
    if (!key) {
      continue
    }
    merged.set(key, item)
  }
  return Array.from(merged.values())
}

function getSearchDoc(post: PostItem): SearchDoc {
  const cacheKey = post.id || post.slug || post.title
  const cached = searchDocCache.get(cacheKey)
  if (cached) {
    return cached
  }

  const title = normalizeForSearch(post.title || '')
  const excerpt = normalizeForSearch(post.excerpt || '')
  const category = normalizeForSearch(post.category || '')
  const tags = normalizeForSearch((post.tags || []).map((tag) => tag.name).join(' '))
  const merged = normalizeForSearch(`${title} ${excerpt} ${category} ${tags} ${post.slug || ''}`)
  const words = Array.from(new Set(tokenize(merged))).slice(0, 90)
  const doc = { title, excerpt, category, tags, merged, words }
  searchDocCache.set(cacheKey, doc)
  return doc
}

function createBigrams(text: string): string[] {
  if (text.length <= 1) {
    return text ? [text] : []
  }
  const values: string[] = []
  for (let idx = 0; idx < text.length - 1; idx += 1) {
    values.push(text.slice(idx, idx + 2))
  }
  return values
}

function diceCoefficient(left: string, right: string): number {
  if (!left || !right) {
    return 0
  }

  if (left === right) {
    return 1
  }

  const leftPairs = createBigrams(left)
  const rightPairs = createBigrams(right)
  if (leftPairs.length === 0 || rightPairs.length === 0) {
    return 0
  }

  const bucket = new Map<string, number>()
  for (const pair of leftPairs) {
    bucket.set(pair, (bucket.get(pair) || 0) + 1)
  }

  let overlap = 0
  for (const pair of rightPairs) {
    const count = bucket.get(pair) || 0
    if (count > 0) {
      overlap += 1
      bucket.set(pair, count - 1)
    }
  }

  return (2 * overlap) / (leftPairs.length + rightPairs.length)
}

function boundedLevenshtein(a: string, b: string, maxDistance = 2): number {
  if (a === b) {
    return 0
  }

  if (Math.abs(a.length - b.length) > maxDistance) {
    return maxDistance + 1
  }

  const previous = Array.from({ length: b.length + 1 }, (_, idx) => idx)
  const current = new Array<number>(b.length + 1).fill(0)

  for (let i = 1; i <= a.length; i += 1) {
    current[0] = i
    let minInRow = current[0]

    for (let j = 1; j <= b.length; j += 1) {
      const cost = a[i - 1] === b[j - 1] ? 0 : 1
      const substitution = previous[j - 1] + cost
      const deletion = previous[j] + 1
      const insertion = current[j - 1] + 1
      const value = Math.min(substitution, deletion, insertion)
      current[j] = value
      if (value < minInRow) {
        minInRow = value
      }
    }

    if (minInRow > maxDistance) {
      return maxDistance + 1
    }

    for (let j = 0; j <= b.length; j += 1) {
      previous[j] = current[j]
    }
  }

  return previous[b.length]
}

function bestTokenSimilarity(token: string, words: string[]): number {
  let best = 0
  for (const word of words) {
    if (!word) {
      continue
    }
    const similarity = diceCoefficient(token, word)
    if (similarity > best) {
      best = similarity
    }
    if (best >= 0.95) {
      break
    }
  }
  return best
}

function hasTypoTolerantMatch(token: string, words: string[]): boolean {
  if (token.length < 4) {
    return false
  }
  for (const word of words) {
    if (Math.abs(word.length - token.length) > 2) {
      continue
    }
    if (boundedLevenshtein(token, word, 1) <= 1) {
      return true
    }
  }
  return false
}

function scorePost(post: PostItem, query: string, queryTokens: string[]): number {
  const doc = getSearchDoc(post)
  let score = 0
  let confidentHit = false
  let tokenHits = 0

  if (doc.title.includes(query)) {
    score += 140
    confidentHit = true
  }

  if (doc.title.startsWith(query)) {
    score += 64
    confidentHit = true
  }

  if (doc.merged.includes(query)) {
    score += 78
    confidentHit = true
  }

  for (const token of queryTokens) {
    if (!token) {
      continue
    }

    let tokenScore = 0
    if (doc.title.includes(token)) {
      tokenScore = Math.max(tokenScore, 36)
    }
    if (doc.tags.includes(token)) {
      tokenScore = Math.max(tokenScore, 32)
    }
    if (doc.category.includes(token)) {
      tokenScore = Math.max(tokenScore, 28)
    }
    if (doc.excerpt.includes(token)) {
      tokenScore = Math.max(tokenScore, 20)
    }

    if (tokenScore === 0) {
      const similarity = bestTokenSimilarity(token, doc.words)
      if (similarity >= 0.8) {
        tokenScore = Math.max(tokenScore, 18 + similarity * 16)
      } else if (similarity >= 0.68) {
        tokenScore = Math.max(tokenScore, 10 + similarity * 12)
      } else if (hasTypoTolerantMatch(token, doc.words)) {
        tokenScore = Math.max(tokenScore, 12)
      }
    }

    if (tokenScore > 0) {
      tokenHits += 1
      score += tokenScore
      if (tokenScore >= 18) {
        confidentHit = true
      }
    }
  }

  if (queryTokens.length > 0) {
    score += (tokenHits / queryTokens.length) * 36
  }

  const titleSimilarity = diceCoefficient(query, doc.title)
  if (titleSimilarity >= 0.55) {
    score += titleSimilarity * 48
    confidentHit = true
  }

  const mergedSimilarity = diceCoefficient(query, doc.merged)
  if (mergedSimilarity >= 0.45) {
    score += mergedSimilarity * 26
  }

  score += Math.log10((post.view_count || 0) + 1) * 3.2
  score += Math.log10((post.like_count || 0) + 1) * 2.3

  if (!confidentHit && score < 40) {
    return 0
  }

  return score
}

function rankPosts(pool: PostItem[], query: string): RankedPost[] {
  if (!query) {
    return []
  }

  const queryTokens = tokenize(query)
  const ranked: RankedPost[] = []

  for (const post of pool) {
    const score = scorePost(post, query, queryTokens)
    if (score > 0) {
      ranked.push({ post, score })
    }
  }

  ranked.sort((left, right) => {
    if (right.score !== left.score) {
      return right.score - left.score
    }
    return parsePostTime(right.post) - parsePostTime(left.post)
  })

  return ranked
}

function formatDate(iso?: string): string {
  if (!iso) return '--'
  try {
    const d = new Date(iso)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  } catch {
    return iso.slice(0, 10)
  }
}

async function ensureSearchPool() {
  if (searchIndexed.value || searchIndexLoading.value) {
    return
  }

  searchIndexLoading.value = true
  try {
    const firstPage = await api.get<PagedData<PostItem>>(`/posts?page=1&size=${SEARCH_INDEX_PAGE_SIZE}`)
    const allItems = [...(firstPage.items || [])]
    const totalCount = Number(firstPage.total || allItems.length)
    const totalPages = Math.max(1, Math.ceil(totalCount / SEARCH_INDEX_PAGE_SIZE))
    const cappedPages = Math.min(totalPages, SEARCH_INDEX_MAX_PAGES)
    const requests: Array<Promise<PagedData<PostItem>>> = []

    for (let current = 2; current <= cappedPages; current += 1) {
      requests.push(api.get<PagedData<PostItem>>(`/posts?page=${current}&size=${SEARCH_INDEX_PAGE_SIZE}`))
    }

    const restPages = requests.length > 0 ? await Promise.all(requests) : []
    for (const pageData of restPages) {
      allItems.push(...(pageData.items || []))
    }

    searchPool.value = mergeUniquePosts([...searchPool.value, ...allItems])
    searchIndexed.value = true
  } catch {
    if (searchPool.value.length === 0) {
      searchPool.value = mergeUniquePosts(posts.value)
    }
  } finally {
    searchIndexLoading.value = false
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const data = await api.get<PagedData<PostItem>>(`/posts?page=${page.value}&size=${pageSize}`)
    const loadedPosts = data.items || []
    posts.value = loadedPosts
    total.value = data.total || 0
    if (!searchIndexed.value) {
      searchPool.value = mergeUniquePosts([...searchPool.value, ...loadedPosts])
    }
  } catch {
    error.value = copy.loadError
    posts.value = []
  } finally {
    loading.value = false
  }
}

function goToPage(p: number) {
  page.value = p
  load()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function handleSearchFocus() {
  searchFocused.value = true
  void ensureSearchPool()
}

function handleSearchInput() {
  activeSuggestionIndex.value = -1
  if (!searchIndexed.value && normalizedSearchQuery.value.length > 1) {
    void ensureSearchPool()
  }
}

function clearSearch() {
  searchQuery.value = ''
  activeSuggestionIndex.value = -1
}

function selectSuggestion(item: PostItem) {
  searchQuery.value = item.title
  activeSuggestionIndex.value = -1
  searchFocused.value = false
}

function handleSearchKeydown(event: KeyboardEvent) {
  if (!showSuggestions.value) {
    if (event.key === 'Escape') {
      searchFocused.value = false
    }
    return
  }

  if (event.key === 'ArrowDown') {
    event.preventDefault()
    activeSuggestionIndex.value = Math.min(activeSuggestionIndex.value + 1, suggestionItems.value.length - 1)
    return
  }

  if (event.key === 'ArrowUp') {
    event.preventDefault()
    activeSuggestionIndex.value = Math.max(activeSuggestionIndex.value - 1, 0)
    return
  }

  if (event.key === 'Enter' && activeSuggestionIndex.value >= 0) {
    event.preventDefault()
    const item = suggestionItems.value[activeSuggestionIndex.value]
    if (item) {
      selectSuggestion(item)
    }
    return
  }

  if (event.key === 'Escape') {
    event.preventDefault()
    searchFocused.value = false
    activeSuggestionIndex.value = -1
  }
}

function handleOutsidePointer(event: PointerEvent) {
  const target = event.target
  if (!(target instanceof Node)) {
    return
  }
  if (searchPanelRef.value?.contains(target)) {
    return
  }
  searchFocused.value = false
  activeSuggestionIndex.value = -1
}

watch(
  normalizedSearchQuery,
  (nextQuery) => {
    if (nextQuery.length > 1 && !searchIndexed.value) {
      void ensureSearchPool()
    }
  },
)

defineExpose({ postCount: total })

onMounted(() => {
  load()
  window.addEventListener('pointerdown', handleOutsidePointer)
})

onBeforeUnmount(() => {
  window.removeEventListener('pointerdown', handleOutsidePointer)
})
</script>

<style scoped>
.blog-card {
  transform: translateY(0);
  transition: transform 0.5s cubic-bezier(0.22, 1, 0.36, 1), box-shadow 0.5s ease, border-color 0.5s ease;
}
.blog-card:hover {
  transform: translateY(-4px);
}
</style>
