<template>
  <button
    type="button"
    class="inline-flex items-center gap-1.5 rounded-xl border px-2 py-1.5 text-xs font-semibold transition sm:px-2.5"
    :class="themeClass"
    :aria-label="ariaLabel"
    @click="toggleThemeMode"
  >
    <svg v-if="mode === 'night'" viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
      <path d="M21 12.8A9 9 0 1111.2 3a7 7 0 009.8 9.8z" />
    </svg>
    <svg v-else viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
      <circle cx="12" cy="12" r="4" />
      <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
    </svg>
    <span class="hidden min-[360px]:inline">{{ modeLabel }}</span>
  </button>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted } from 'vue'

import { siteCopy } from '../../content/copy'
import { hydrateThemeMode, themeMode, toggleThemeMode } from '../../stores/theme'

const copy = siteCopy.blogTopNav.themeToggle
const mode = useStore(themeMode)

const modeLabel = computed(() => {
  return mode.value === 'night' ? copy.nightLabel : copy.lightLabel
})

const ariaLabel = computed(() => {
  return mode.value === 'night' ? copy.switchToLightAria : copy.switchToNightAria
})

const themeClass = computed(() => {
  if (mode.value === 'night') {
    return 'border-slate-500/65 bg-slate-900/90 text-slate-100 hover:border-miku/45 hover:text-miku'
  }

  return 'border-white/35 bg-white/30 text-slate-700 hover:border-miku/35 hover:text-miku'
})

onMounted(() => {
  hydrateThemeMode()
})
</script>
