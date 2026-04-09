<template>
  <div class="space-y-6">
    <LiquidGlassCard padding="0px">
      <div class="relative overflow-hidden">
        <div class="h-20 bg-gradient-to-r from-[#39c5bb]/30 via-[#c084fc]/20 to-[#39c5bb]/25" />
        <div class="relative -mt-10 px-5 pb-5">
          <div class="author-avatar-ring mx-auto h-20 w-20 rounded-full p-[3px]">
            <div class="relative h-full w-full overflow-hidden rounded-full border-2 border-white">
              <div
                v-if="!avatarReady"
                class="absolute inset-0 flex items-center justify-center bg-white/86"
              >
                <svg viewBox="0 0 24 24" class="h-4.5 w-4.5 text-miku" fill="none" aria-hidden="true">
                  <path d="M14.7 4.2c1.8.9 3.2 2.3 4.1 4.1L13.2 14l-2.8.4.4-2.8 3.9-7.4z" fill="currentColor" fill-opacity="0.92" />
                  <path d="M10.6 11.6L8 14.2a1.8 1.8 0 000 2.6l.8.8a1.8 1.8 0 002.6 0l2.6-2.6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                  <path d="M8.3 17.5l-2.1 2.6 3.3-1.3z" fill="#f59e0b" />
                </svg>
              </div>
              <img
                :src="settings.avatarUrl"
                :alt="copy.authorCard.avatarAlt"
                class="h-full w-full rounded-full object-cover transition-[opacity,transform] duration-700 hover:rotate-[360deg]"
                :class="avatarReady ? 'opacity-100' : 'opacity-0'"
                @load="avatarReady = true"
                @error="avatarReady = true"
              />
            </div>
          </div>

          <div class="mt-3 text-center">
            <h3 class="text-xl font-bold text-slate-900">{{ settings.displayName }}</h3>
            <p class="mt-0.5 text-xs text-miku/70">{{ settings.role }}</p>
            <p class="mt-2 text-sm leading-relaxed text-slate-500">{{ settings.bio }}</p>
          </div>

          <div class="mt-3 flex items-center justify-center gap-3 text-[11px] text-slate-400">
            <span class="inline-flex items-center gap-1">
              <svg viewBox="0 0 24 24" class="h-3 w-3 fill-none stroke-current stroke-[2]" aria-hidden="true"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0118 0z" /><circle cx="12" cy="10" r="3" /></svg>
              {{ settings.location }}
            </span>
            <span class="h-0.5 w-0.5 rounded-full bg-slate-300" />
            <span class="inline-flex items-center gap-1">
              <svg viewBox="0 0 24 24" class="h-3 w-3 fill-none stroke-current stroke-[2]" aria-hidden="true"><rect x="3" y="4" width="18" height="18" rx="2" ry="2" /><line x1="16" y1="2" x2="16" y2="6" /><line x1="8" y1="2" x2="8" y2="6" /><line x1="3" y1="10" x2="21" y2="10" /></svg>
              {{ settings.since }}
            </span>
          </div>

          <div class="mt-3 flex flex-wrap justify-center gap-1.5">
            <span
              v-for="skill in settings.skills"
              :key="skill"
              class="rounded-md border border-slate-200 bg-slate-50/80 px-2 py-0.5 text-[10px] font-medium text-slate-500"
            >
              {{ skill }}
            </span>
          </div>

          <div class="mt-4 flex items-center justify-center gap-2">
            <a
              v-for="link in settings.socialLinks"
              :key="`${link.label}-${link.href}`"
              :href="link.href"
              target="_blank"
              rel="noopener noreferrer"
              class="flex h-8 w-8 items-center justify-center rounded-lg border border-slate-200 bg-white/60 text-slate-400 transition duration-300 hover:border-miku/40 hover:text-miku hover:shadow-sm"
              :aria-label="link.label"
            >
              <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
                <path :d="iconPath(link.iconKey)" />
              </svg>
            </a>
          </div>

          <div class="mt-4">
            <AuthorStats />
          </div>

          <a
            href="/about"
            class="mt-4 inline-flex w-full items-center justify-center rounded-xl border border-miku/35 bg-miku-soft px-3 py-2 text-xs font-semibold text-miku transition hover:border-miku/55"
          >
            {{ copy.authorCard.aboutCta }}
          </a>
        </div>
      </div>
    </LiquidGlassCard>

    <div id="recent-updates" class="glass-layer scroll-mt-24 rounded-2xl p-4">
      <h3 class="text-sm font-semibold text-slate-900">{{ copy.nowTitle }}</h3>
      <ul class="mt-3 space-y-2">
        <li
          v-for="item in settings.nowItems"
          :key="item"
          class="rounded-lg border border-white/60 bg-white/66 px-3 py-2 text-xs leading-relaxed text-slate-600"
        >
          {{ item }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref, watch } from 'vue'

import { siteCopy } from '../../content/copy'
import {
  authorProfileSettings,
  hydrateAuthorProfileSettings,
} from '../../stores/authorProfile'
import AuthorStats from './AuthorStats.vue'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

const copy = siteCopy.blogIndex
const store = useStore(authorProfileSettings)
const avatarReady = ref(false)

const settings = computed(() => ({
  displayName: store.value.displayName,
  avatarUrl: store.value.avatarUrl,
  role: store.value.role,
  bio: store.value.bio,
  location: store.value.location,
  since: store.value.since,
  skills: store.value.skills,
  socialLinks: store.value.socialLinks,
  nowItems: store.value.nowItems,
}))

watch(() => settings.value.avatarUrl, () => {
  avatarReady.value = false
}, { immediate: true })

function iconPath(iconKey: string) {
  const normalized = (iconKey || '').trim().toLowerCase()
  const iconMap: Record<string, string> = {
    github: 'M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 00-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0020 4.77 5.07 5.07 0 0019.91 1S18.73.65 16 2.48a13.38 13.38 0 00-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 005 4.77a5.44 5.44 0 00-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 009 18.13V22',
    x: 'M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z',
    mail: 'M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2zm16 2l-8 5-8-5',
    qq: 'M12 3c3.87 0 7 2.98 7 6.65 0 2.12-1.04 4.01-2.67 5.23.13.52.42 1.35.97 2.28.19.32-.09.72-.45.65-1.45-.28-2.51-.79-3.13-1.14-.86.23-1.77.35-2.72.35-3.87 0-7-2.98-7-6.65S8.13 3 12 3z',
    bilibili: 'M5 6.5h14a2.5 2.5 0 012.5 2.5v6A2.5 2.5 0 0119 17.5H5A2.5 2.5 0 012.5 15V9A2.5 2.5 0 015 6.5zm4.5 4v3l2.5-1.5-2.5-1.5zM8 4 6.5 6M16 4l1.5 2',
    link: 'M10 13a5 5 0 007.07 0l2.12-2.12a5 5 0 10-7.07-7.07L10 5M14 11a5 5 0 00-7.07 0L4.81 13.12a5 5 0 007.07 7.07L14 19',
  }

  return iconMap[normalized] || iconMap.link
}

onMounted(() => {
  void hydrateAuthorProfileSettings()
})
</script>

<style scoped>
.author-avatar-ring {
  background: conic-gradient(#39c5bb, #c084fc, #39c5bb);
  animation: ring-rotate 6s linear infinite;
}

@keyframes ring-rotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
