<template>
  <footer class="border-t border-slate-200/70 bg-white/72 py-5 backdrop-blur-md">
    <div class="mx-auto flex w-full max-w-[1320px] flex-col items-center gap-2 px-4 text-center text-xs text-slate-500 sm:px-6 lg:px-8">
      <p class="font-medium text-slate-600">{{ copyrightLine }}</p>

      <div class="flex flex-wrap items-center justify-center gap-x-2 gap-y-1">
        <a
          v-if="settings.icpText && settings.icpLink"
          :href="settings.icpLink"
          target="_blank"
          rel="noopener noreferrer"
          class="font-medium text-miku transition hover:text-[#2ea99f] hover:underline"
          :aria-label="`${copy.display.icpAriaPrefix}${settings.icpText}`"
        >
          {{ settings.icpText }}
        </a>
        <span v-else-if="settings.icpText" class="font-medium text-slate-600">{{ settings.icpText }}</span>
        <span v-else>{{ copy.display.emptyIcpText }}</span>
      </div>

      <ul
        v-if="settings.customTexts.length > 0"
        class="space-y-1 text-[11px] leading-relaxed text-slate-500"
      >
        <li v-for="(line, index) in settings.customTexts" :key="`footer-line-${index}`">
          {{ line }}
        </li>
      </ul>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted } from 'vue'

import { siteCopy } from '../../content/copy'
import { hydrateSiteFooterSettings, siteFooterSettings } from '../../stores/siteFooter'

const copy = siteCopy.siteFooter
const settings = useStore(siteFooterSettings)

const currentYear = new Date().getFullYear()

const copyrightLine = computed(() => {
  return `${copy.display.copyrightPrefix}${currentYear} ${copy.display.siteName} ${copy.display.rightsText}`
})

onMounted(() => {
  hydrateSiteFooterSettings()
})
</script>
