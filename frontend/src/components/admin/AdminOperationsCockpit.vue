<template>
  <section class="space-y-4">
    <AdminPlainCard max-width="100%" padding="0">
      <div class="rounded-[16px] bg-[linear-gradient(135deg,rgba(57,197,187,0.16),rgba(192,132,252,0.12),rgba(255,255,255,0.98))] p-5">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div class="max-w-3xl">
            <p class="inline-flex rounded-full border border-white/70 bg-white/72 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.24em] text-slate-600">
              {{ copy.badge }}
            </p>
            <h2 class="mt-3 text-2xl font-semibold text-slate-900">{{ copy.title }}</h2>
            <p class="mt-2 max-w-2xl text-sm leading-6 text-slate-600">{{ copy.subtitle }}</p>
          </div>

          <div class="min-w-[220px] rounded-2xl border border-white/80 bg-white/72 px-4 py-3 text-sm text-slate-600 shadow-[0_12px_36px_rgba(15,23,42,0.06)]">
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">{{ copy.windowPrefix }}</p>
            <p class="mt-2 font-semibold text-slate-900">{{ analyticsWindowLabel }}</p>
            <p class="mt-1 text-xs text-slate-500">{{ analyticsWindowRange }}</p>
          </div>
        </div>

        <div class="mt-5 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <div
            v-for="card in overviewCards"
            :key="card.key"
            class="rounded-2xl border px-4 py-4 shadow-[0_14px_36px_rgba(15,23,42,0.05)]"
            :class="overviewCardClass(card.tone)"
          >
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">{{ card.label }}</p>
            <p class="mt-3 font-mono text-3xl font-semibold text-slate-900">{{ card.value }}</p>
            <p class="mt-2 text-sm font-medium text-slate-700">{{ card.detail }}</p>
            <p class="mt-1 text-xs text-slate-500">{{ card.hint }}</p>
          </div>
        </div>
      </div>
    </AdminPlainCard>

    <AdminPlainCard max-width="100%" padding="20px">
      <div class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-lg font-semibold text-slate-900">{{ copy.alerts.title }}</h3>
          <p class="text-sm text-slate-600">{{ copy.subtitle }}</p>
        </div>
      </div>

      <div class="mt-4 grid gap-3 lg:grid-cols-2 xl:grid-cols-4">
        <div
          v-for="alert in operationsAlerts"
          :key="alert.key"
          class="rounded-2xl border px-4 py-3"
          :class="alertCardClass(alert.tone)"
        >
          <p class="text-sm font-semibold">{{ alert.title }}</p>
          <p class="mt-1 text-xs leading-5 opacity-85">{{ alert.detail }}</p>
        </div>
      </div>
    </AdminPlainCard>

    <div class="grid gap-4 xl:grid-cols-[1.25fr_0.95fr]">
      <AdminPlainCard max-width="100%" padding="20px">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold text-slate-900">{{ copy.performance.title }}</h3>
            <p class="text-sm text-slate-600">{{ copy.performance.subtitle }}</p>
          </div>
        </div>

        <div v-if="loading && topPosts.length === 0" class="mt-4 flex h-[220px] items-center justify-center text-sm text-slate-400">
          {{ baseCopy.common.loading }}
        </div>
        <div v-else-if="topPosts.length === 0" class="mt-4 rounded-2xl border border-slate-200/80 bg-white/45 px-4 py-8 text-center text-sm text-slate-500">
          {{ copy.performance.empty }}
        </div>
        <div v-else class="mt-4 space-y-3">
          <div
            v-for="post in topPosts"
            :key="post.id"
            class="rounded-2xl border border-slate-200/80 bg-[linear-gradient(180deg,rgba(255,255,255,0.92),rgba(248,250,252,0.92))] px-4 py-3"
          >
            <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <p class="truncate text-sm font-semibold text-slate-900">{{ post.title }}</p>
                  <span
                    v-if="isRecentPublished(post.published_at)"
                    class="inline-flex rounded-full border border-miku/20 bg-miku/10 px-2 py-0.5 text-[11px] font-semibold text-miku"
                  >
                    {{ copy.performance.recentBadge }}
                  </span>
                </div>
                <p class="mt-1 text-xs text-slate-500">{{ formatDateTime(post.published_at || post.created_at) }}</p>
              </div>

              <div class="grid grid-cols-4 gap-2 text-right lg:min-w-[288px]">
                <div class="rounded-xl bg-slate-50 px-2 py-2">
                  <p class="text-[11px] uppercase tracking-[0.14em] text-slate-400">{{ copy.performance.columnViews }}</p>
                  <p class="mt-1 text-sm font-semibold text-slate-900">{{ formatInteger(post.view_count) }}</p>
                </div>
                <div class="rounded-xl bg-slate-50 px-2 py-2">
                  <p class="text-[11px] uppercase tracking-[0.14em] text-slate-400">{{ copy.performance.columnLikes }}</p>
                  <p class="mt-1 text-sm font-semibold text-slate-900">{{ formatInteger(post.like_count) }}</p>
                </div>
                <div class="rounded-xl bg-slate-50 px-2 py-2">
                  <p class="text-[11px] uppercase tracking-[0.14em] text-slate-400">{{ copy.performance.columnComments }}</p>
                  <p class="mt-1 text-sm font-semibold text-slate-900">{{ formatInteger(post.comment_count) }}</p>
                </div>
                <div class="rounded-xl bg-[linear-gradient(135deg,rgba(57,197,187,0.14),rgba(192,132,252,0.12))] px-2 py-2">
                  <p class="text-[11px] uppercase tracking-[0.14em] text-slate-500">{{ copy.performance.columnScore }}</p>
                  <p class="mt-1 text-sm font-semibold text-slate-900">{{ formatInteger(postScore(post)) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </AdminPlainCard>

      <AdminPlainCard max-width="100%" padding="20px">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold text-slate-900">{{ copy.rhythm.title }}</h3>
            <p class="text-sm text-slate-600">{{ copy.rhythm.subtitle }}</p>
          </div>
        </div>

        <div class="mt-4 grid gap-3 sm:grid-cols-2">
          <div
            v-for="item in rhythmStats"
            :key="item.key"
            class="rounded-2xl border border-slate-200/80 bg-white/70 px-4 py-3"
          >
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">{{ item.label }}</p>
            <p class="mt-2 font-mono text-2xl font-semibold text-slate-900">{{ item.value }}</p>
          </div>
        </div>

        <div class="mt-5 grid gap-4 lg:grid-cols-2">
          <div>
            <p class="text-sm font-semibold text-slate-900">{{ copy.rhythm.recentTitle }}</p>
            <div v-if="recentPublishedPosts.length === 0" class="mt-3 rounded-2xl border border-slate-200/80 bg-white/45 px-4 py-6 text-center text-sm text-slate-500">
              {{ copy.rhythm.emptyRecent }}
            </div>
            <div v-else class="mt-3 space-y-2">
              <div
                v-for="post in recentPublishedPosts"
                :key="post.id"
                class="rounded-2xl border border-slate-200/80 bg-white/72 px-3 py-3"
              >
                <p class="truncate text-sm font-semibold text-slate-900">{{ post.title }}</p>
                <p class="mt-1 text-xs text-slate-500">{{ formatDateTime(post.published_at || post.created_at) }}</p>
              </div>
            </div>
          </div>

          <div>
            <p class="text-sm font-semibold text-slate-900">{{ copy.rhythm.queueTitle }}</p>
            <div v-if="queueFocusItems.length === 0" class="mt-3 rounded-2xl border border-slate-200/80 bg-white/45 px-4 py-6 text-center text-sm text-slate-500">
              {{ copy.rhythm.emptyQueue }}
            </div>
            <div v-else class="mt-3 space-y-2">
              <div
                v-for="item in queueFocusItems"
                :key="item.key"
                class="rounded-2xl border border-slate-200/80 bg-white/72 px-3 py-3"
              >
                <div class="flex items-center justify-between gap-2">
                  <p class="truncate text-sm font-semibold text-slate-900">{{ item.title }}</p>
                  <span class="inline-flex rounded-full px-2 py-0.5 text-[11px] font-semibold" :class="queuePillClass(item.tone)">
                    {{ item.label }}
                  </span>
                </div>
                <p class="mt-1 text-xs text-slate-500">{{ item.meta }}</p>
              </div>
            </div>
          </div>
        </div>
      </AdminPlainCard>
    </div>

    <div class="grid gap-4 xl:grid-cols-[1.05fr_0.95fr]">
      <AdminPlainCard max-width="100%" padding="20px">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold text-slate-900">{{ copy.review.title }}</h3>
            <p class="text-sm text-slate-600">{{ copy.review.subtitle }}</p>
          </div>
        </div>

        <div class="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <div
            v-for="item in reviewMetrics"
            :key="item.key"
            class="rounded-2xl border border-slate-200/80 bg-white/72 px-4 py-3"
          >
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">{{ item.label }}</p>
            <p class="mt-2 font-mono text-2xl font-semibold text-slate-900">{{ item.value }}</p>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap gap-2">
          <span
            v-for="item in reviewQueueHints"
            :key="item.key"
            class="inline-flex rounded-full border border-slate-200/80 bg-slate-50 px-3 py-1 text-xs font-medium text-slate-600"
          >
            {{ item.label }}
          </span>
        </div>

        <div v-if="activityLoading" class="mt-5 flex h-[160px] items-center justify-center text-sm text-slate-400">
          {{ baseCopy.common.loading }}
        </div>
        <div v-else-if="recentReviewActivities.length === 0" class="mt-5 rounded-2xl border border-slate-200/80 bg-white/45 px-4 py-8 text-center text-sm text-slate-500">
          {{ copy.review.empty }}
        </div>
        <div v-else class="mt-5 space-y-2">
          <div
            v-for="item in recentReviewActivities"
            :key="item.id"
            class="rounded-2xl border border-slate-200/80 bg-white/72 px-3 py-3"
          >
            <div class="flex items-start justify-between gap-2">
              <p class="text-sm text-slate-700">{{ formatActivity(item) }}</p>
              <span class="shrink-0 text-xs text-slate-500">{{ formatRelativeTime(item.created_at) }}</span>
            </div>
            <p v-if="item.admin_username" class="mt-1 text-xs text-slate-500">{{ baseCopy.common.operatorPrefix }}{{ item.admin_username }}</p>
          </div>
        </div>
      </AdminPlainCard>

      <AdminPlainCard max-width="100%" padding="20px">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold text-slate-900">{{ copy.dependencies.title }}</h3>
            <p class="text-sm text-slate-600">{{ copy.dependencies.subtitle }}</p>
          </div>
        </div>

        <div class="mt-4 grid gap-3 sm:grid-cols-2">
          <div
            v-for="item in dependencyCards"
            :key="item.key"
            class="rounded-2xl border px-4 py-3"
            :class="dependencyCardClass(item.state)"
          >
            <div class="flex items-center justify-between gap-2">
              <p class="text-sm font-semibold text-slate-900">{{ item.label }}</p>
              <span class="inline-flex rounded-full px-2 py-0.5 text-[11px] font-semibold" :class="dependencyPillClass(item.state)">
                {{ dependencyStateLabel(item.state) }}
              </span>
            </div>
            <p class="mt-2 text-sm text-slate-700">{{ item.message }}</p>
            <p class="mt-1 text-xs leading-5 text-slate-500">{{ item.detail }}</p>
          </div>
        </div>

        <div class="mt-5 rounded-2xl border border-slate-200/80 bg-[linear-gradient(180deg,rgba(255,255,255,0.92),rgba(241,245,249,0.9))] px-4 py-4">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.dependencies.widgetsTitle }}</p>
              <p class="mt-1 text-xs leading-5 text-slate-500">{{ copy.dependencies.widgetsHint }}</p>
            </div>
            <p class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">
              {{ enabledWidgetsCount }}/3 {{ copy.breakdown.activeWidgets }}
            </p>
          </div>

          <div class="mt-4 flex flex-wrap gap-2">
            <span
              v-for="widget in widgetStates"
              :key="widget.key"
              class="inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-semibold"
              :class="widgetChipClass(widget.enabled)"
            >
              <span>{{ widget.label }}</span>
              <span>{{ widget.enabled ? copy.dependencies.stateHealthy : copy.dependencies.stateDisabled }}</span>
            </span>
          </div>
        </div>
      </AdminPlainCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { adminCopy } from '../../content/copy'
import { api, type PagedData } from '../../lib/api'
import {
  getDefaultSiteIntegrationsSettings,
  normalizeSiteIntegrationsSettings,
  type SiteIntegrationsSettings,
  type SiteIntegrationsSettingsPayload,
} from '../../lib/site-integrations'
import AdminPlainCard from '../ui/AdminPlainCard.vue'

interface MetricValue {
  value: number
  change: number
}

interface AnalyticsOverview {
  window: {
    label: string
    start: string
    end: string
  }
  summary: {
    visitors: MetricValue
    views: MetricValue
  }
  location: {
    countries: Array<{ code: string; visitors: number }>
  }
}

interface DashboardStats {
  total_posts: number
  total_likes: number
  pending_comments: number
  friend_count: number
  draft_count: number
}

interface AuditLogItem {
  id: string
  action: string
  target_type: string
  admin_username?: string
  created_at: string
}

interface AdminPostItem {
  id: string
  title: string
  status: string
  published_at?: string
  scheduled_at?: string
  created_at: string
  updated_at: string
  view_count: number
  like_count: number
  comment_count: number
}

interface AdminMomentItem {
  id: string
  content: string
  publish_status: string
  published_at?: string
  scheduled_at?: string
  created_at: string
}

interface AdminFriendItem {
  id: string
  name: string
  health_status: string
}

interface GitHubProbeResponse {
  profile?: {
    name?: string
    total_repos?: number
    followers?: number
  }
}

interface WeatherProbeResponse {
  temp?: string
  desc?: string
  location?: string
}

type AlertTone = 'miku' | 'amber' | 'rose' | 'sky'
type SummaryTone = 'miku' | 'lavender' | 'amber' | 'slate'
type QueueTone = 'miku' | 'amber' | 'lavender'
type DependencyState = 'ok' | 'warning' | 'error' | 'idle'

const props = defineProps<{
  analytics: AnalyticsOverview | null
  dashboardStats: DashboardStats | null
  activityItems: AuditLogItem[]
  activityLoading: boolean
}>()

const baseCopy = adminCopy.dashboard
const copy = adminCopy.dashboard.operations

const loading = ref(true)
const posts = ref<AdminPostItem[]>([])
const moments = ref<AdminMomentItem[]>([])
const pendingCommentsTotal = ref(0)
const pendingGuestbookTotal = ref(0)
const pendingFriendApplicationsTotal = ref(0)
const friends = ref<AdminFriendItem[]>([])
const integrations = ref<SiteIntegrationsSettings>(getDefaultSiteIntegrationsSettings())
const githubDependency = ref<{ state: DependencyState; message: string; detail: string }>({
  state: 'idle',
  message: copy.dependencies.stateIdle,
  detail: copy.dependencies.githubMissing,
})
const weatherDependency = ref<{ state: DependencyState; message: string; detail: string }>({
  state: 'idle',
  message: copy.dependencies.stateIdle,
  detail: copy.dependencies.weatherDefault,
})

const analyticsWindowLabel = computed(() => props.analytics?.window.label || baseCopy.common.windowFallback)

const analyticsWindowRange = computed(() => {
  if (!props.analytics?.window.start || !props.analytics?.window.end) {
    return baseCopy.common.noData
  }
  return `${formatDateTime(props.analytics.window.start)} ~ ${formatDateTime(props.analytics.window.end)}`
})

const publishedPosts = computed(() => posts.value.filter((item) => item.status === 'published'))
const scheduledPosts = computed(() => posts.value.filter((item) => item.status === 'scheduled'))
const draftPosts = computed(() => posts.value.filter((item) => item.status === 'draft'))

const topPosts = computed(() => {
  return publishedPosts.value
    .slice()
    .sort((a, b) => postScore(b) - postScore(a))
    .slice(0, 5)
})

const recentPublishedPosts = computed(() => {
  return publishedPosts.value
    .slice()
    .sort((a, b) => compareTime(b.published_at || b.created_at, a.published_at || a.created_at))
    .slice(0, 4)
})

const staleDraftItems = computed(() => {
  return draftPosts.value
    .filter((item) => ageInDays(item.updated_at || item.created_at) >= 14)
    .slice()
    .sort((a, b) => compareTime(a.updated_at || a.created_at, b.updated_at || b.created_at))
})

const scheduledMoments = computed(() => moments.value.filter((item) => item.publish_status === 'scheduled'))

const publishedPosts7d = computed(() => publishedPosts.value.filter((item) => ageInDays(item.published_at || item.created_at) <= 7).length)

const activeMoments7d = computed(() => {
  return moments.value.filter((item) => {
    if (item.publish_status !== 'published') return false
    return ageInDays(item.published_at || item.created_at) <= 7
  }).length
})

const totalPendingReview = computed(() => pendingCommentsTotal.value + pendingGuestbookTotal.value + pendingFriendApplicationsTotal.value)

const totalCommentStock = computed(() => publishedPosts.value.reduce((sum, item) => sum + Number(item.comment_count || 0), 0))

const enabledWidgetsCount = computed(() => {
  let total = 0
  if (integrations.value.showWeather) total += 1
  if (integrations.value.showMusic) total += 1
  if (integrations.value.showClock) total += 1
  return total
})

const geoipDependency = computed(() => {
  const countryCount = props.analytics?.location?.countries?.length || 0
  if (!props.analytics) {
    return {
      key: 'geoip',
      label: copy.dependencies.geoip,
      state: 'idle' as DependencyState,
      message: copy.dependencies.stateIdle,
      detail: copy.dependencies.geoipWarningDetail,
    }
  }
  if (countryCount > 0) {
      return {
        key: 'geoip',
        label: copy.dependencies.geoip,
        state: 'ok' as DependencyState,
        message: `${formatInteger(countryCount)} ${copy.breakdown.geoipNodesSuffix}`,
        detail: copy.dependencies.geoipHealthyDetail,
      }
  }
  return {
    key: 'geoip',
    label: copy.dependencies.geoip,
    state: 'warning' as DependencyState,
    message: copy.dependencies.stateWarning,
    detail: copy.dependencies.geoipWarningDetail,
  }
})

const friendDependency = computed(() => {
  if (friends.value.length === 0) {
    return {
      key: 'friends',
      label: copy.dependencies.friends,
      state: 'idle' as DependencyState,
      message: copy.dependencies.stateIdle,
      detail: copy.dependencies.friendsUnknownDetail,
    }
  }

  const downCount = friends.value.filter((item) => item.health_status === 'down').length
  const unknownCount = friends.value.filter((item) => item.health_status === 'unknown').length
  if (downCount > 0) {
      return {
        key: 'friends',
        label: copy.dependencies.friends,
        state: 'error' as DependencyState,
        message: `${copy.dependencies.friendsWarningPrefix} ${formatInteger(downCount)}`,
        detail: `${formatInteger(downCount)} / ${formatInteger(friends.value.length)} ${copy.dependencies.friendsDownDetail}`,
      }
    }
  if (unknownCount > 0) {
    return {
      key: 'friends',
      label: copy.dependencies.friends,
      state: 'warning' as DependencyState,
      message: `${formatInteger(unknownCount)} ${copy.dependencies.friendsPendingDetail}`,
      detail: copy.dependencies.friendsUnknownDetail,
    }
  }
  return {
    key: 'friends',
    label: copy.dependencies.friends,
    state: 'ok' as DependencyState,
    message: `${formatInteger(friends.value.length)} / ${formatInteger(friends.value.length)} ${copy.breakdown.dependenciesSuffix}`,
    detail: copy.dependencies.friendsHealthyDetail,
  }
})

const dependencyCards = computed(() => [
  {
    key: 'github',
    label: copy.dependencies.github,
    state: githubDependency.value.state,
    message: githubDependency.value.message,
    detail: githubDependency.value.detail,
  },
  {
    key: 'weather',
    label: copy.dependencies.weather,
    state: weatherDependency.value.state,
    message: weatherDependency.value.message,
    detail: weatherDependency.value.detail,
  },
  geoipDependency.value,
  friendDependency.value,
])

const healthyDependenciesCount = computed(() => dependencyCards.value.filter((item) => item.state === 'ok').length)

const widgetStates = computed(() => [
  { key: 'weather', label: copy.dependencies.widgetWeather, enabled: integrations.value.showWeather },
  { key: 'music', label: copy.dependencies.widgetMusic, enabled: integrations.value.showMusic },
  { key: 'clock', label: copy.dependencies.widgetClock, enabled: integrations.value.showClock },
])

const overviewCards = computed(() => [
  {
    key: 'published',
    label: copy.cards.published,
    value: formatInteger(Number(props.dashboardStats?.total_posts || 0)),
    detail: `${formatInteger(publishedPosts.value.length)} ${copy.breakdown.publishedStatus} · ${formatInteger(scheduledPosts.value.length)} ${copy.breakdown.scheduledStatus}`,
    hint: copy.cards.publishedHint,
    tone: 'miku' as SummaryTone,
  },
  {
    key: 'interactions',
    label: copy.cards.interactions,
    value: formatInteger(Number(props.dashboardStats?.total_likes || 0)),
    detail: `${formatInteger(totalCommentStock.value)} ${copy.breakdown.commentsStock}`,
    hint: copy.cards.interactionsHint,
    tone: 'lavender' as SummaryTone,
  },
  {
    key: 'review',
    label: copy.cards.reviewQueue,
    value: formatInteger(totalPendingReview.value),
    detail: `${formatInteger(pendingCommentsTotal.value)} ${copy.breakdown.commentsShort} · ${formatInteger(pendingGuestbookTotal.value)} ${copy.breakdown.guestbookShort} · ${formatInteger(pendingFriendApplicationsTotal.value)} ${copy.breakdown.friendsShort}`,
    hint: copy.cards.reviewQueueHint,
    tone: 'amber' as SummaryTone,
  },
  {
    key: 'dependencies',
    label: copy.cards.dependencies,
    value: `${healthyDependenciesCount.value}/4`,
    detail: `${enabledWidgetsCount.value}/3 ${copy.breakdown.activeWidgets}`,
    hint: copy.cards.dependenciesHint,
    tone: 'slate' as SummaryTone,
  },
])

const rhythmStats = computed(() => [
  { key: 'published7d', label: copy.rhythm.published7d, value: formatInteger(publishedPosts7d.value) },
  { key: 'scheduledPosts', label: copy.rhythm.scheduledPosts, value: formatInteger(scheduledPosts.value.length) },
  { key: 'scheduledMoments', label: copy.rhythm.scheduledMoments, value: formatInteger(scheduledMoments.value.length) },
  { key: 'activeMoments', label: copy.rhythm.activeMoments, value: formatInteger(activeMoments7d.value) },
])

const queueFocusItems = computed(() => {
  const items: Array<{ key: string; title: string; label: string; meta: string; tone: QueueTone }> = []

  scheduledPosts.value.slice(0, 3).forEach((item) => {
    items.push({
      key: `post-${item.id}`,
      title: item.title,
      label: copy.rhythm.dueLabel,
      meta: formatDateTime(item.scheduled_at || item.created_at),
      tone: 'miku',
    })
  })

  scheduledMoments.value.slice(0, 2).forEach((item) => {
    items.push({
      key: `moment-${item.id}`,
      title: trimText(item.content).slice(0, 28) || copy.rhythm.scheduledMoments,
      label: copy.rhythm.dueLabel,
      meta: formatDateTime(item.scheduled_at || item.created_at),
      tone: 'lavender',
    })
  })

  staleDraftItems.value.slice(0, 3).forEach((item) => {
    items.push({
      key: `draft-${item.id}`,
      title: item.title,
      label: copy.rhythm.staleLabel,
      meta: `${formatInteger(ageInDays(item.updated_at || item.created_at))}${copy.breakdown.daysSuffix}`,
      tone: 'amber',
    })
  })

  return items.slice(0, 6)
})

const reviewMetrics = computed(() => [
  { key: 'comments', label: copy.review.pendingComments, value: formatInteger(pendingCommentsTotal.value) },
  { key: 'guestbook', label: copy.review.pendingGuestbook, value: formatInteger(pendingGuestbookTotal.value) },
  { key: 'friends', label: copy.review.pendingApplications, value: formatInteger(pendingFriendApplicationsTotal.value) },
  { key: 'actions', label: copy.review.recentActions, value: formatInteger(recentReviewActivities.value.length) },
])

const reviewQueueHints = computed(() => [
  { key: 'comments', label: `${copy.review.commentQueueHint} · ${formatInteger(pendingCommentsTotal.value)}` },
  { key: 'guestbook', label: `${copy.review.guestbookQueueHint} · ${formatInteger(pendingGuestbookTotal.value)}` },
  { key: 'friends', label: `${copy.review.friendQueueHint} · ${formatInteger(pendingFriendApplicationsTotal.value)}` },
])

const recentReviewActivities = computed(() => {
  return (props.activityItems || [])
    .filter((item) => ['approve', 'reject', 'delete'].includes(item.action))
    .filter((item) => ['comment', 'guestbook', 'friend_link', 'friend_link_application'].includes(item.target_type))
    .slice(0, 6)
})

const operationsAlerts = computed(() => {
  const alerts: Array<{ key: string; title: string; detail: string; tone: AlertTone }> = []
  const visitorChange = Number(props.analytics?.summary?.visitors?.change || 0)
  const dependencyIssues = dependencyCards.value.filter((item) => item.state === 'error').length

  if (visitorChange <= -25) {
    alerts.push({
      key: 'traffic-drop',
      title: copy.alerts.trafficDropTitle,
      detail: copy.alerts.trafficDropDetail,
      tone: 'rose',
    })
  }

  if (totalPendingReview.value >= 6) {
    alerts.push({
      key: 'review-backlog',
      title: copy.alerts.reviewBacklogTitle,
      detail: copy.alerts.reviewBacklogDetail,
      tone: 'amber',
    })
  }

  if (staleDraftItems.value.length > 0) {
    alerts.push({
      key: 'stale-drafts',
      title: copy.alerts.draftStaleTitle,
      detail: copy.alerts.draftStaleDetail,
      tone: 'sky',
    })
  }

  if (dependencyIssues > 0) {
    alerts.push({
      key: 'dependencies',
      title: copy.alerts.dependencyTitle,
      detail: copy.alerts.dependencyDetail,
      tone: 'rose',
    })
  }

  if (alerts.length === 0) {
    alerts.push({
      key: 'stable',
      title: copy.alerts.stableTitle,
      detail: copy.alerts.stableDetail,
      tone: 'miku',
    })
  }

  return alerts.slice(0, 4)
})

async function loadOperations() {
  loading.value = true

  const results = await Promise.allSettled([
    api.get<PagedData<AdminPostItem>>('/admin/posts?page=1&size=100'),
    api.get<PagedData<AdminMomentItem>>('/admin/moments?page=1&size=100'),
    api.get<PagedData<unknown>>('/admin/comments?status=pending&page=1&size=1'),
    api.get<PagedData<unknown>>('/admin/guestbook/messages?status=pending&page=1&size=1'),
    api.get<PagedData<AdminFriendItem>>('/admin/friends?page=1&size=100'),
    api.get<PagedData<unknown>>('/admin/friends/applications?status=pending&page=1&size=1'),
    api.get<SiteIntegrationsSettingsPayload | undefined>('/site-settings/site-integrations'),
  ])

  const [postsResult, momentsResult, commentsResult, guestbookResult, friendsResult, friendApplicationsResult, integrationsResult] = results

  posts.value = postsResult.status === 'fulfilled' ? postsResult.value.items || [] : []
  moments.value = momentsResult.status === 'fulfilled' ? momentsResult.value.items || [] : []
  pendingCommentsTotal.value = commentsResult.status === 'fulfilled' ? Number(commentsResult.value.total || 0) : 0
  pendingGuestbookTotal.value = guestbookResult.status === 'fulfilled' ? Number(guestbookResult.value.total || 0) : 0
  friends.value = friendsResult.status === 'fulfilled' ? friendsResult.value.items || [] : []
  pendingFriendApplicationsTotal.value = friendApplicationsResult.status === 'fulfilled' ? Number(friendApplicationsResult.value.total || 0) : 0
  integrations.value = integrationsResult.status === 'fulfilled'
    ? normalizeSiteIntegrationsSettings(integrationsResult.value)
    : getDefaultSiteIntegrationsSettings()

  await Promise.allSettled([probeGitHubStatus(), probeWeatherStatus()])
  loading.value = false
}

async function probeGitHubStatus() {
  const username = trimText(integrations.value.githubUsername)
  if (!username) {
    githubDependency.value = {
      state: 'idle',
      message: copy.dependencies.stateIdle,
      detail: copy.dependencies.githubMissing,
    }
    return
  }

  try {
    const data = await api.get<GitHubProbeResponse>(`/github/profile?username=${encodeURIComponent(username)}`)
    const repoCount = Number(data.profile?.total_repos || 0)
    const displayName = trimText(data.profile?.name) || username
    githubDependency.value = {
      state: 'ok',
      message: displayName,
      detail: repoCount > 0 ? `${formatInteger(repoCount)} ${copy.dependencies.githubReadableDetail}` : copy.dependencies.githubOkDetail,
    }
  } catch {
    githubDependency.value = {
      state: 'error',
      message: copy.dependencies.stateError,
      detail: copy.dependencies.githubErrorDetail,
    }
  }
}

async function probeWeatherStatus() {
  const location = trimText(integrations.value.weatherLocation)
  try {
    const query = location ? `?location=${encodeURIComponent(location)}` : ''
    const data = await api.get<WeatherProbeResponse>(`/weather${query}`)
    const displayLocation = trimText(data.location) || location || copy.dependencies.defaultLocation
    const temp = trimText(data.temp)
    weatherDependency.value = {
      state: 'ok',
      message: temp ? `${displayLocation} · ${temp}°C` : displayLocation,
      detail: trimText(data.desc) || (location ? copy.dependencies.weatherOkDetail : copy.dependencies.weatherDefault),
    }
  } catch {
    weatherDependency.value = {
      state: 'error',
      message: copy.dependencies.stateError,
      detail: copy.dependencies.weatherErrorDetail,
    }
  }
}

function postScore(item: Pick<AdminPostItem, 'view_count' | 'like_count' | 'comment_count'>): number {
  return Number(item.view_count || 0) + Number(item.like_count || 0) * 3 + Number(item.comment_count || 0) * 5
}

function isRecentPublished(value?: string) {
  if (!value) return false
  return ageInDays(value) <= 7
}

function compareTime(a?: string, b?: string) {
  return parseTime(a).getTime() - parseTime(b).getTime()
}

function parseTime(value?: string) {
  const date = new Date(value || '')
  if (Number.isNaN(date.getTime())) {
    return new Date(0)
  }
  return date
}

function ageInDays(value?: string) {
  const target = parseTime(value)
  if (target.getTime() === 0) return 999
  return Math.max(0, Math.floor((Date.now() - target.getTime()) / (24 * 60 * 60 * 1000)))
}

function trimText(value?: string) {
  return typeof value === 'string' ? value.trim() : ''
}

function formatInteger(value: number) {
  return Number(value || 0).toLocaleString('en-US')
}

function formatDateTime(iso: string) {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

function formatRelativeTime(iso: string) {
  const target = parseTime(iso)
  if (target.getTime() === 0) return iso

  const diffMs = Date.now() - target.getTime()
  const mins = Math.floor(diffMs / 60000)
  if (mins < 1) return baseCopy.common.justNow
  if (mins < 60) return `${mins} ${baseCopy.common.minutesAgo}`

  const hoursDiff = Math.floor(mins / 60)
  if (hoursDiff < 24) return `${hoursDiff} ${baseCopy.common.hoursAgo}`

  const days = Math.floor(hoursDiff / 24)
  if (days < 7) return `${days} ${baseCopy.common.daysAgo}`

  return formatDateTime(iso)
}

function formatActivity(item: AuditLogItem) {
  const action = baseCopy.actions[item.action as keyof typeof baseCopy.actions] || item.action
  const target = baseCopy.targets[item.target_type as keyof typeof baseCopy.targets] || item.target_type
  return `${action}${baseCopy.common.activityTargetPrefix}${target}`
}

function overviewCardClass(tone: SummaryTone) {
  if (tone === 'miku') {
    return 'border-miku/25 bg-white/82'
  }
  if (tone === 'lavender') {
    return 'border-[#c084fc]/25 bg-white/82'
  }
  if (tone === 'amber') {
    return 'border-amber-200/80 bg-white/82'
  }
  return 'border-slate-200/80 bg-white/82'
}

function alertCardClass(tone: AlertTone) {
  if (tone === 'miku') {
    return 'border-miku/20 bg-miku/10 text-miku'
  }
  if (tone === 'amber') {
    return 'border-amber-200/80 bg-amber-50/80 text-amber-800'
  }
  if (tone === 'rose') {
    return 'border-rose-200/80 bg-rose-50/80 text-rose-800'
  }
  return 'border-sky-200/80 bg-sky-50/80 text-sky-800'
}

function dependencyCardClass(state: DependencyState) {
  if (state === 'ok') {
    return 'border-miku/20 bg-[linear-gradient(135deg,rgba(57,197,187,0.12),rgba(255,255,255,0.95))]'
  }
  if (state === 'warning') {
    return 'border-amber-200/80 bg-amber-50/80'
  }
  if (state === 'error') {
    return 'border-rose-200/80 bg-rose-50/80'
  }
  return 'border-slate-200/80 bg-slate-50/80'
}

function dependencyPillClass(state: DependencyState) {
  if (state === 'ok') {
    return 'border border-miku/20 bg-white/80 text-miku'
  }
  if (state === 'warning') {
    return 'border border-amber-200 bg-white/70 text-amber-700'
  }
  if (state === 'error') {
    return 'border border-rose-200 bg-white/70 text-rose-700'
  }
  return 'border border-slate-200 bg-white/70 text-slate-500'
}

function dependencyStateLabel(state: DependencyState) {
  if (state === 'ok') return copy.dependencies.stateHealthy
  if (state === 'warning') return copy.dependencies.stateWarning
  if (state === 'error') return copy.dependencies.stateError
  return copy.dependencies.stateIdle
}

function widgetChipClass(enabled: boolean) {
  return enabled
    ? 'border-miku/20 bg-miku/10 text-miku'
    : 'border-slate-200 bg-slate-50 text-slate-500'
}

function queuePillClass(tone: QueueTone) {
  if (tone === 'miku') {
    return 'bg-miku/10 text-miku'
  }
  if (tone === 'lavender') {
    return 'bg-[#c084fc]/12 text-[#8b5cf6]'
  }
  return 'bg-amber-100 text-amber-700'
}

onMounted(async () => {
  await loadOperations()
})
</script>
