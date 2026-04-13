<template>
  <div class="flex items-center justify-between gap-4">
    <div class="flex items-center gap-3">
      <button
        type="button"
        class="inline-flex h-9 w-9 items-center justify-center rounded-lg text-slate-500 transition hover:bg-slate-100 hover:text-slate-900 lg:hidden"
        :aria-label="ch.openMenuAria"
        @click="toggleSidebar"
      >
        <svg viewBox="0 0 24 24" class="h-5 w-5 fill-none stroke-current stroke-[1.8]">
          <path d="M4 7h16M4 12h16M4 17h16" />
        </svg>
      </button>

      <h1 class="text-xl font-semibold text-slate-900">{{ pageTitle }}</h1>
    </div>

    <div class="flex items-center gap-3">
      <div class="hidden text-right sm:block">
        <p class="text-sm font-medium text-slate-900">{{ userName }}</p>
        <p class="text-xs text-slate-500">{{ ch.roleLabel }}</p>
      </div>
      <div
        class="relative h-9 w-9 overflow-hidden rounded-full border border-white/70 bg-miku/10 shadow-[0_10px_24px_rgba(57,197,187,0.18)]"
        :aria-label="userName"
      >
        <img
          v-if="displayAvatarUrl"
          :src="displayAvatarUrl"
          :alt="userName"
          class="h-full w-full object-cover"
          @error="handleAvatarError"
        />
        <div
          v-else
          class="flex h-full w-full items-center justify-center text-sm font-semibold text-miku"
          aria-hidden="true"
        >
          {{ avatarLetter }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useStore } from '@nanostores/vue'

import { authState, hydrateAuth } from '../../stores/auth'
import { getAuthorAvatarFallbackChain } from '../../lib/author-profile'
import {
  authorProfileSettings,
  hydrateAuthorProfileSettings,
} from '../../stores/authorProfile'
import { toggleSidebar } from '../../stores/ui'
import { adminCopy } from '../../content/copy'

const ch = adminCopy.contentHeader

interface Props {
  pageTitle?: string
}

withDefaults(defineProps<Props>(), {
  pageTitle: adminCopy.contentHeader.defaultPageTitle,
})

const auth = useStore(authState)
const authorProfile = useStore(authorProfileSettings)
const mounted = ref(false)
const avatarFallbackIndex = ref(0)

const userName = computed(() => {
  if (!mounted.value) return 'Admin'
  return auth.value.user?.name ?? 'Admin'
})

const avatarLetter = computed(() => {
  const name = userName.value
  return name.charAt(0).toUpperCase()
})

const avatarSources = computed(() => getAuthorAvatarFallbackChain(authorProfile.value.avatarUrl))
const displayAvatarUrl = computed(() => {
  return avatarSources.value[avatarFallbackIndex.value] || avatarSources.value[0] || ''
})

watch(() => avatarSources.value.join('||'), () => {
  avatarFallbackIndex.value = 0
}, { immediate: true })

function handleAvatarError() {
  if (avatarFallbackIndex.value < avatarSources.value.length - 1) {
    avatarFallbackIndex.value += 1
    return
  }
  avatarFallbackIndex.value = avatarSources.value.length
}

onMounted(() => {
  mounted.value = true
  void hydrateAuth()
  void hydrateAuthorProfileSettings()
})
</script>
