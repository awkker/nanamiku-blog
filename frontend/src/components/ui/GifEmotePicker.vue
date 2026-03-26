<template>
  <div class="rounded-xl border border-slate-200/80 bg-white/60 p-2.5">
    <div class="flex items-center justify-between gap-2">
      <button
        type="button"
        class="inline-flex items-center gap-1.5 rounded-lg border border-miku/30 bg-miku/8 px-2.5 py-1 text-xs font-semibold text-miku transition hover:border-miku/50 hover:bg-miku/12"
        :aria-label="panelOpen ? copy.collapseAria : copy.expandAria"
        @click="togglePanel"
      >
        <svg viewBox="0 0 24 24" class="h-3.5 w-3.5 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
          <path d="M7 10h10M7 14h7" />
          <rect x="3" y="4" width="18" height="16" rx="3" />
        </svg>
        {{ panelOpen ? copy.closeButton : copy.openButton }}
      </button>
      <span class="text-[11px] text-slate-400">{{ copy.hint }}</span>
    </div>

    <div v-if="panelOpen" class="mt-2.5 space-y-2.5">
      <div class="inline-flex items-center gap-1 rounded-lg border border-slate-200 bg-white/80 p-1">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          type="button"
          class="rounded-md px-2.5 py-1 text-[11px] font-medium transition"
          :class="activePack === tab.key ? 'bg-miku text-white' : 'text-slate-500 hover:text-slate-700'"
          @click="activePack = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <div class="grid grid-cols-5 gap-2 sm:grid-cols-6">
        <button
          v-for="item in activeItems"
          :key="item.id"
          type="button"
          class="inline-flex items-center justify-center rounded-lg border border-slate-200/80 bg-white/80 p-1 transition hover:border-miku/40 hover:bg-miku/5"
          :title="item.label"
          :aria-label="`${copy.insertAriaPrefix}${item.label}`"
          @click="handleSelect(item.path)"
        >
          <img
            :src="item.src"
            :alt="`${copy.altPrefix}${item.label}`"
            loading="lazy"
            decoding="async"
            class="h-10 w-10 rounded-md object-cover sm:h-12 sm:w-12"
          />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

import { siteCopy } from '../../content/copy'
import { GIF_EMOTE_PACKS, GIF_PACK_ORDER, type GifPackKey } from '../../lib/gifEmotes'

const copy = siteCopy.components.gifEmotePicker

const emit = defineEmits<{
  (e: 'select', path: string): void
}>()

const panelOpen = ref(false)
const activePack = ref<GifPackKey>('miku')

const tabLabelMap = new Map(copy.tabs.map((item) => [item.key as GifPackKey, item.label]))
const tabs = computed(() =>
  GIF_PACK_ORDER.map((pack) => ({
    key: pack,
    label: tabLabelMap.get(pack) || pack,
  })),
)
const activeItems = computed(() => GIF_EMOTE_PACKS[activePack.value] || [])

function togglePanel() {
  panelOpen.value = !panelOpen.value
}

function handleSelect(path: string) {
  emit('select', path)
}
</script>
