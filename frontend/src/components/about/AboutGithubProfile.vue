<template>
  <div class="space-y-4">
    <div v-if="loading" class="space-y-4">
      <LiquidGlassCard padding="18px" maxWidth="100%">
        <div class="flex items-center gap-3">
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl border border-white/70 bg-white/70 text-miku">
            <LoadingSpinner size="lg" />
          </span>
          <div>
            <p class="text-sm font-semibold text-slate-800">{{ copy.loadingTitle }}</p>
            <p class="mt-0.5 text-xs text-slate-500">{{ copy.loadingHint }}</p>
          </div>
        </div>
      </LiquidGlassCard>

      <div class="grid gap-4 lg:grid-cols-2">
        <div class="h-56 animate-pulse rounded-2xl border border-white/40 bg-white/25 p-4">
          <div class="h-5 w-1/3 rounded bg-white/55" />
          <div class="mt-4 flex items-center gap-3">
            <div class="h-14 w-14 rounded-xl bg-white/55" />
            <div class="flex-1 space-y-2">
              <div class="h-4 w-2/3 rounded bg-white/55" />
              <div class="h-3 w-4/5 rounded bg-white/45" />
            </div>
          </div>
          <div class="mt-4 grid grid-cols-3 gap-2">
            <div class="h-12 rounded-lg bg-white/50" />
            <div class="h-12 rounded-lg bg-white/50" />
            <div class="h-12 rounded-lg bg-white/50" />
          </div>
        </div>
        <div class="h-56 animate-pulse rounded-2xl border border-white/40 bg-white/25 p-4">
          <div class="h-5 w-1/3 rounded bg-white/55" />
          <div class="mt-3 h-[168px] rounded-xl bg-white/50" />
        </div>
      </div>
      <div class="animate-pulse rounded-2xl border border-white/40 bg-white/25 p-4">
        <div class="h-4 w-1/5 rounded bg-white/55" />
        <div class="mt-3 flex flex-wrap gap-2">
          <div class="h-7 w-20 rounded-full bg-white/50" />
          <div class="h-7 w-24 rounded-full bg-white/50" />
          <div class="h-7 w-16 rounded-full bg-white/50" />
          <div class="h-7 w-[5.5rem] rounded-full bg-white/50" />
        </div>
      </div>
      <div class="rounded-2xl border border-white/40 bg-white/25 p-4">
        <div class="mb-3 flex items-center gap-2 text-slate-500">
          <LoadingSpinner size="sm" />
          <span class="text-xs">{{ copy.loadingLinks }}</span>
        </div>
        <div class="grid gap-2.5 md:grid-cols-2 lg:grid-cols-3">
          <div v-for="i in 6" :key="i" class="h-20 animate-pulse rounded-xl border border-white/55 bg-white/45" />
        </div>
      </div>
    </div>

    <LiquidGlassCard v-else-if="error" padding="20px" maxWidth="100%">
      <p class="text-sm text-slate-500">{{ copy.loadError }}</p>
    </LiquidGlassCard>

    <template v-else>
      <!-- Row 1: Profile Card + Activity Chart -->
      <div class="grid gap-4 lg:grid-cols-2">
        <LiquidGlassCard padding="24px" maxWidth="100%">
          <div class="flex items-start gap-5">
            <div class="relative h-20 w-20 shrink-0 overflow-hidden rounded-2xl border-2 border-[#39c5bb]/25 shadow-md">
              <div
                v-if="!avatarReady"
                class="absolute inset-0 z-10 flex items-center justify-center bg-white/80 text-miku"
              >
                <LoadingSpinner size="sm" />
              </div>
              <img
                v-if="profile.avatarUrl"
                :src="profile.avatarUrl"
                :alt="profile.name"
                class="h-full w-full object-cover transition-opacity duration-300"
                :class="avatarReady ? 'opacity-100' : 'opacity-0'"
                @load="handleAvatarLoaded"
                @error="handleAvatarLoaded"
              />
              <div v-else class="flex h-full w-full items-center justify-center bg-white/72 text-[#39c5bb]">
                <svg viewBox="0 0 24 24" class="h-9 w-9 fill-none stroke-current stroke-[1.6]" aria-hidden="true">
                  <circle cx="12" cy="8" r="4" />
                  <path d="M4.5 20a7.5 7.5 0 0115 0" />
                </svg>
              </div>
            </div>
            <div class="min-w-0">
              <h2 class="text-lg font-bold text-slate-900">{{ profile.name }}</h2>
              <p class="mt-0.5 text-sm leading-relaxed text-slate-500">{{ profile.bio }}</p>
              <a
                :href="profile.htmlUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="mt-2 inline-flex items-center gap-1.5 text-xs font-medium text-[#39c5bb] transition hover:text-[#2ba89f]"
              >
                <svg class="h-3.5 w-3.5" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
                </svg>
                {{ copy.viewProfile }}
              </a>
            </div>
          </div>
          <div class="mt-5 grid grid-cols-3 gap-2.5">
            <div class="rounded-xl border border-white/65 bg-white/62 px-3 py-2.5 text-center">
              <p class="text-lg font-bold text-slate-900">{{ profile.totalRepos }}</p>
              <p class="text-xs text-slate-400">{{ copy.stats.repos }}</p>
            </div>
            <div class="rounded-xl border border-white/65 bg-white/62 px-3 py-2.5 text-center">
              <p class="text-lg font-bold text-slate-900">{{ profile.totalStars }}</p>
              <p class="text-xs text-slate-400">{{ copy.stats.stars }}</p>
            </div>
            <div class="rounded-xl border border-white/65 bg-white/62 px-3 py-2.5 text-center">
              <p class="text-lg font-bold text-slate-900">{{ profile.followers }}</p>
              <p class="text-xs text-slate-400">{{ copy.stats.followers }}</p>
            </div>
          </div>
        </LiquidGlassCard>

        <LiquidGlassCard padding="24px" maxWidth="100%">
          <h3 class="text-base font-semibold text-slate-900">{{ copy.activityTitle }}</h3>
          <p class="mt-0.5 text-xs text-slate-400">{{ copy.activitySubtitle }}</p>
          <div class="mt-3" style="height: 184px">
            <v-chart :option="chartOption" autoresize />
          </div>
        </LiquidGlassCard>
      </div>

      <!-- Tech Stack -->
      <LiquidGlassCard padding="20px" maxWidth="100%">
        <h3 class="text-base font-semibold text-slate-900">{{ copy.techStackTitle }}</h3>
        <p class="mt-0.5 text-xs text-slate-400">{{ copy.techStackSubtitle }}</p>
        <div class="mt-3 flex flex-wrap gap-2">
          <span
            v-for="(tech, index) in techStack"
            :key="tech.name"
            class="rounded-full border px-3 py-1 text-xs font-medium"
            :class="
              index % 3 === 0
                ? 'border-[#39c5bb]/35 bg-[#39c5bb]/8 text-[#2a9d8f]'
                : index % 3 === 1
                  ? 'border-[#c084fc]/35 bg-[#c084fc]/8 text-[#8b5cf6]'
                  : 'border-slate-200 bg-white/70 text-slate-600'
            "
          >
            {{ tech.name }}
          </span>
        </div>
      </LiquidGlassCard>

      <!-- Recent Repos -->
      <LiquidGlassCard padding="24px" maxWidth="100%">
        <h3 class="text-base font-semibold text-slate-900">{{ copy.recentReposTitle }}</h3>
        <p class="mt-0.5 text-xs text-slate-400">{{ copy.recentReposSubtitle }}</p>
        <div class="mt-3 grid gap-2.5 md:grid-cols-2 lg:grid-cols-3">
          <a
            v-for="repo in recentRepos"
            :key="repo.name"
            :href="repo.htmlUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="group block rounded-xl border border-white/65 bg-white/62 p-4 transition hover:-translate-y-0.5 hover:border-[#39c5bb]/30 hover:shadow-md"
          >
            <div class="flex items-center justify-between">
              <h4 class="truncate text-sm font-semibold text-[#39c5bb] group-hover:text-[#2ba89f]">{{ repo.name }}</h4>
              <span class="ml-2 shrink-0 text-xs text-slate-400">&#9733; {{ repo.stars }}</span>
            </div>
            <p v-if="repo.description" class="mt-1.5 line-clamp-2 text-xs leading-relaxed text-slate-500">{{ repo.description }}</p>
            <div class="mt-2 flex items-center gap-3 text-xs text-slate-400">
              <span v-if="repo.language" class="flex items-center gap-1">
                <span class="inline-block h-2 w-2 rounded-full" :style="{ backgroundColor: getLangColor(repo.language) }" />
                {{ repo.language }}
              </span>
              <span>{{ formatDate(repo.pushedAt) }}</span>
            </div>
          </a>
        </div>
      </LiquidGlassCard>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import LoadingSpinner from '../ui/LoadingSpinner.vue'
import { siteCopy } from '../../content/copy'
import { api } from '../../lib/api'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent])

const MIKU = '#39c5bb'

interface Props {
  githubUsername: string
}

const props = defineProps<Props>()

const loading = ref(true)
const error = ref(false)
const copy = siteCopy.components.aboutGithubProfile

interface ProfileData {
  avatarUrl: string
  name: string
  bio: string
  htmlUrl: string
  totalRepos: number
  totalStars: number
  followers: number
}

const profile = ref<ProfileData>({
  avatarUrl: '',
  name: '',
  bio: '',
  htmlUrl: '',
  totalRepos: 0,
  totalStars: 0,
  followers: 0,
})
const avatarReady = ref(false)

const techStack = ref<Array<{ name: string; count: number }>>([])

interface RepoItem {
  name: string
  description: string
  htmlUrl: string
  language: string
  stars: number
  pushedAt: string
}

const recentRepos = ref<RepoItem[]>([])
const activityData = ref<number[]>(new Array(12).fill(0))

interface GitHubProfilePayload {
  profile: {
    avatar_url: string
    name: string
    bio: string
    html_url: string
    total_repos: number
    total_stars: number
    followers: number
  }
  tech_stack: Array<{ name: string; count: number }>
  recent_repos: Array<{
    name: string
    description: string
    html_url: string
    language: string
    stars: number
    pushed_at: string
  }>
  activity_data: number[]
}

function handleAvatarLoaded() {
  avatarReady.value = true
}

watch(
  () => profile.value.avatarUrl,
  (nextAvatar) => {
    avatarReady.value = !nextAvatar
  },
  { immediate: true },
)

const monthLabels = computed(() => {
  const now = new Date()
  const labels: string[] = []
  for (let i = 11; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
    labels.push(`${d.getMonth() + 1}${copy.monthSuffix}`)
  }
  return labels
})

const chartOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    backgroundColor: 'rgba(255,255,255,0.92)',
    borderColor: 'rgba(0,0,0,0.06)',
    textStyle: { color: '#334155', fontSize: 12 },
  },
  grid: { top: 12, left: 8, right: 8, bottom: 0, containLabel: true },
  xAxis: {
    type: 'category',
    data: monthLabels.value,
    boundaryGap: false,
    axisLine: { lineStyle: { color: '#e2e8f0' } },
    axisLabel: { color: '#94a3b8', fontSize: 11 },
  },
  yAxis: {
    type: 'value',
    splitLine: { lineStyle: { color: '#f1f5f9' } },
    axisLabel: { color: '#94a3b8', fontSize: 11 },
  },
  series: [
    {
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 5,
      data: activityData.value,
      lineStyle: { color: MIKU, width: 2.5 },
      itemStyle: { color: MIKU },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(57,197,187,0.25)' },
            { offset: 1, color: 'rgba(57,197,187,0.02)' },
          ],
        },
      },
    },
  ],
}))

const LANG_COLORS: Record<string, string> = {
  JavaScript: '#f1e05a', TypeScript: '#3178c6', Python: '#3572A5',
  Go: '#00ADD8', Rust: '#dea584', Java: '#b07219',
  'C++': '#f34b7d', C: '#555555', 'C#': '#178600',
  Ruby: '#701516', Swift: '#F05138', Kotlin: '#A97BFF',
  PHP: '#4F5D95', HTML: '#e34c26', CSS: '#563d7c',
  SCSS: '#c6538c', Vue: '#41b883', Shell: '#89e051',
  Dart: '#00B4AB', Lua: '#000080', Dockerfile: '#384d54',
}

function getLangColor(lang: string): string {
  return LANG_COLORS[lang] || '#94a3b8'
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

async function fetchGitHubData() {
  const username = props.githubUsername
  if (!username) {
    error.value = false
    loading.value = false
    return
  }

  try {
    loading.value = true
    error.value = false
    const data = await api.get<GitHubProfilePayload>(`/github/profile?username=${encodeURIComponent(username)}`)

    profile.value = {
      avatarUrl: data.profile.avatar_url || '',
      name: data.profile.name || username,
      bio: data.profile.bio || '',
      htmlUrl: data.profile.html_url || `https://github.com/${username}`,
      totalRepos: Number(data.profile.total_repos) || 0,
      totalStars: Number(data.profile.total_stars) || 0,
      followers: Number(data.profile.followers) || 0,
    }

    techStack.value = (data.tech_stack || []).map((item) => ({
      name: item.name,
      count: Number(item.count) || 0,
    }))

    recentRepos.value = (data.recent_repos || []).map((repo) => ({
      name: repo.name,
      description: repo.description || '',
      htmlUrl: repo.html_url,
      language: repo.language || '',
      stars: Number(repo.stars) || 0,
      pushedAt: repo.pushed_at || '',
    }))

    activityData.value = Array.isArray(data.activity_data) && data.activity_data.length === 12
      ? data.activity_data
      : new Array(12).fill(0)

    loading.value = false
  } catch {
    error.value = true
    loading.value = false
  }
}

watch(() => props.githubUsername, (next, prev) => {
  if (next === prev) {
    return
  }
  void fetchGitHubData()
})

onMounted(() => { void fetchGitHubData() })
</script>
