<template>
  <LiquidGlassCard padding="26px" maxWidth="100%">
    <div class="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-lg font-bold text-slate-900">{{ copy.contactSection.title }}</h2>
        <p class="mt-1 text-sm text-slate-500">{{ copy.contactSection.subtitle }}</p>
      </div>
      <div class="flex flex-wrap gap-2">
        <a
          :href="emailHref"
          class="inline-flex items-center rounded-xl border border-miku/40 bg-miku-soft px-4 py-2 text-sm font-semibold text-miku transition hover:border-miku/60"
        >
          {{ copy.contactSection.emailButton }}
        </a>
        <a
          v-if="githubHref"
          :href="githubHref"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center rounded-xl border border-slate-200 bg-white/75 px-4 py-2 text-sm font-semibold text-slate-600 transition hover:border-[#c084fc]/40 hover:text-[#8b5cf6]"
        >
          {{ copy.contactSection.githubButton }}
        </a>
        <a
          v-for="link in contactLinks"
          :key="`${link.label}-${link.href}`"
          :href="link.href"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center rounded-xl border border-slate-200 bg-white/75 px-4 py-2 text-sm font-semibold text-slate-600 transition hover:-translate-y-0.5 hover:border-[#c084fc]/40 hover:text-[#8b5cf6]"
        >
          {{ link.label }}
        </a>
      </div>
    </div>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted } from 'vue'

import { siteCopy } from '../../content/copy'
import {
  authorProfileSettings,
  hydrateAuthorProfileSettings,
} from '../../stores/authorProfile'
import {
  hydrateSiteIntegrationsSettings,
  siteIntegrationsSettings,
} from '../../stores/siteIntegrations'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

const copy = siteCopy.aboutPage
const authorStore = useStore(authorProfileSettings)
const integrationsStore = useStore(siteIntegrationsSettings)

const emailHref = computed(() => {
  const email = (authorStore.value.contactEmail || '').trim()
  return email ? `mailto:${email}?subject=${encodeURIComponent('合作交流')}` : copy.contactSection.emailHref
})

const githubHref = computed(() => {
  const username = (integrationsStore.value.githubUsername || '').trim()
  return username ? `https://github.com/${username}` : ''
})

const contactLinks = computed(() => {
  return authorStore.value.contactLinks.length > 0 ? authorStore.value.contactLinks : copy.socialLinks
})

onMounted(() => {
  void Promise.all([
    hydrateAuthorProfileSettings(),
    hydrateSiteIntegrationsSettings(),
  ])
})
</script>
