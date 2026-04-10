<template>
  <div class="space-y-4 sm:space-y-5">
    <section class="relative mx-auto max-w-[1320px] px-4 pt-12 sm:px-6 lg:px-8 lg:pt-16">
      <div class="grid gap-6 lg:grid-cols-[1.1fr_360px]">
        <LiquidGlassCard padding="30px" maxWidth="100%">
          <div class="space-y-5">
            <p class="inline-flex rounded-full border border-white/60 bg-white/58 px-4 py-1 text-xs tracking-[0.18em] text-slate-500">
              {{ copy.heroBadge }}
            </p>
            <h1 class="bg-gradient-to-r from-slate-900 via-[#39c5bb] to-[#c084fc] bg-clip-text text-[clamp(2rem,4.8vw,3.8rem)] font-black tracking-tight text-transparent">
              {{ heroTitle }}
            </h1>
            <p class="max-w-3xl text-base leading-relaxed text-slate-600 sm:text-lg">
              {{ aboutDescription }}
            </p>

            <div class="flex flex-wrap gap-2">
              <span
                v-for="tag in identityTags"
                :key="tag"
                class="rounded-full border border-slate-200 bg-white/70 px-3 py-1 text-xs font-medium text-slate-600"
              >
                {{ tag }}
              </span>
            </div>

            <div class="flex flex-wrap gap-2.5 pt-1">
              <a
                v-for="(action, index) in heroActions"
                :key="`${action.href}-${index}`"
                :href="action.href"
                :class="heroActionClasses[index] ?? heroActionClasses[0]"
              >
                {{ action.label }}
              </a>
            </div>

            <div class="grid gap-2.5 pt-2 sm:grid-cols-3">
              <div class="rounded-xl border border-white/65 bg-white/62 px-3 py-2.5">
                <p class="text-xs text-slate-400">{{ copy.stats.postsLabel }}</p>
                <p class="mt-1 text-lg font-bold text-slate-900">{{ postCount }}</p>
              </div>
              <div class="rounded-xl border border-white/65 bg-white/62 px-3 py-2.5">
                <p class="text-xs text-slate-400">{{ copy.stats.focusLabel }}</p>
                <p class="mt-1 text-sm font-semibold text-slate-800">{{ copy.stats.focusValue }}</p>
              </div>
              <div class="rounded-xl border border-white/65 bg-white/62 px-3 py-2.5">
                <p class="text-xs text-slate-400">{{ copy.stats.statusLabel }}</p>
                <p class="mt-1 text-sm font-semibold text-slate-800">{{ copy.stats.statusValue }}</p>
              </div>
            </div>
          </div>
        </LiquidGlassCard>

        <LiquidGlassCard padding="24px" maxWidth="100%">
          <div class="space-y-4">
            <div class="flex items-center gap-3">
              <div class="relative h-14 w-14">
                <div
                  v-if="!avatarReady"
                  class="absolute inset-0 flex items-center justify-center rounded-xl border border-slate-200/75 bg-white/82"
                >
                  <svg viewBox="0 0 24 24" class="h-4.5 w-4.5 text-miku" fill="none" aria-hidden="true">
                    <path d="M14.7 4.2c1.8.9 3.2 2.3 4.1 4.1L13.2 14l-2.8.4.4-2.8 3.9-7.4z" fill="currentColor" fill-opacity="0.92" />
                    <path d="M10.6 11.6L8 14.2a1.8 1.8 0 000 2.6l.8.8a1.8 1.8 0 002.6 0l2.6-2.6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                    <path d="M8.3 17.5l-2.1 2.6 3.3-1.3z" fill="#f59e0b" />
                  </svg>
                </div>
                <img
                  :src="displayAvatarUrl"
                  :alt="copy.profileCard.avatarAlt"
                  class="h-14 w-14 rounded-xl border border-slate-200 object-cover transition-opacity duration-300"
                  :class="avatarReady ? 'opacity-100' : 'opacity-0'"
                  @load="handleAvatarLoaded"
                  @error="handleAvatarError"
                />
              </div>
              <div>
                <p class="text-base font-semibold text-slate-900">{{ settings.displayName }}</p>
                <p class="text-xs text-slate-500">{{ settings.role }}</p>
              </div>
            </div>
            <div class="rounded-xl border border-white/65 bg-white/58 px-3 py-2.5">
              <p class="text-xs tracking-[0.14em] text-miku">{{ copy.profileCard.nowTitle }}</p>
              <ul class="mt-2 space-y-2 text-xs text-slate-600">
                <li
                  v-for="item in settings.nowItems"
                  :key="item"
                  class="flex items-start gap-2"
                >
                  <span class="mt-1 h-1.5 w-1.5 rounded-full bg-miku" />
                  <span>{{ item }}</span>
                </li>
              </ul>
            </div>
            <blockquote class="rounded-xl border border-white/65 bg-white/58 px-3 py-2.5 text-sm leading-relaxed text-slate-600">
              {{ settings.quote }}
            </blockquote>
          </div>
        </LiquidGlassCard>
      </div>
    </section>

    <section v-if="githubUsername" id="github" class="relative mx-auto max-w-[1320px] px-4 sm:px-6 lg:px-8">
      <AboutGithubProfile :github-username="githubUsername" />
    </section>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref, watch } from 'vue'

import { siteCopy } from '../../content/copy'
import { getAuthorAvatarFallbackChain } from '../../lib/author-profile'
import {
  authorProfileSettings,
  hydrateAuthorProfileSettings,
} from '../../stores/authorProfile'
import {
  hydrateSiteIntegrationsSettings,
  siteIntegrationsSettings,
} from '../../stores/siteIntegrations'
import AboutGithubProfile from './AboutGithubProfile.vue'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

interface Props {
  postCount: number
}

const props = defineProps<Props>()
const copy = siteCopy.aboutPage

const heroActionClasses = [
  'inline-flex items-center rounded-xl border border-miku/40 bg-miku-soft px-4 py-2.5 text-sm font-semibold text-miku transition hover:-translate-y-0.5 hover:border-miku/60',
  'inline-flex items-center rounded-xl border border-slate-200 bg-white/70 px-4 py-2.5 text-sm font-semibold text-slate-600 transition hover:-translate-y-0.5 hover:border-[#c084fc]/45 hover:text-[#8b5cf6]',
  'inline-flex items-center rounded-xl border border-slate-200 bg-white/70 px-4 py-2.5 text-sm font-semibold text-slate-600 transition hover:-translate-y-0.5 hover:border-miku/40 hover:text-miku',
]

const authorStore = useStore(authorProfileSettings)
const integrationsStore = useStore(siteIntegrationsSettings)
const avatarReady = ref(false)
const avatarFallbackIndex = ref(0)

const settings = computed(() => authorStore.value)
const postCount = computed(() => props.postCount)
const heroTitle = computed(() => settings.value.displayName ? `我是 ${settings.value.displayName}` : copy.heroTitle)
const aboutDescription = computed(() => settings.value.aboutDescription || copy.heroDescription)
const identityTags = computed(() => settings.value.skills.length > 0 ? settings.value.skills : copy.identityTags)
const githubUsername = computed(() => (integrationsStore.value.githubUsername || '').trim())
const contactEmail = computed(() => (settings.value.contactEmail || '').trim())
const heroActions = computed(() => {
  const actions = [
    copy.heroActions[0],
    copy.heroActions[1],
  ]

  if (contactEmail.value) {
    actions.push({
      ...copy.heroActions[2],
      href: `mailto:${contactEmail.value}?subject=${encodeURIComponent('来自博客的联系')}`,
    })
  }

  return actions
})

const avatarSources = computed(() => getAuthorAvatarFallbackChain(settings.value.avatarUrl))
const displayAvatarUrl = computed(() => {
  return avatarSources.value[avatarFallbackIndex.value] || avatarSources.value[0] || ''
})

watch(() => avatarSources.value.join('||'), () => {
  avatarReady.value = false
  avatarFallbackIndex.value = 0
}, { immediate: true })

function handleAvatarLoaded() {
  avatarReady.value = true
}

function handleAvatarError() {
  if (avatarFallbackIndex.value < avatarSources.value.length - 1) {
    avatarFallbackIndex.value += 1
    avatarReady.value = false
    return
  }

  avatarReady.value = true
}

onMounted(() => {
  void Promise.all([
    hydrateAuthorProfileSettings(),
    hydrateSiteIntegrationsSettings(),
  ])
})
</script>
