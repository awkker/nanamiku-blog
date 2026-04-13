<template>
  <header class="relative overflow-hidden border-b border-white/60">
    <div class="absolute inset-0 bg-[url('/picture/bloglist.jpg')] bg-cover bg-center bg-no-repeat" />
    <div class="absolute inset-0 bg-gradient-to-br from-slate-950/42 via-slate-900/30 to-white/86" />
    <div class="pointer-events-none absolute -left-16 top-6 h-56 w-56 rounded-full bg-miku/24 blur-3xl" />
    <div class="pointer-events-none absolute -right-20 top-16 h-52 w-52 rounded-full bg-[#c084fc]/26 blur-3xl" />

    <div class="relative z-10 mx-auto max-w-[1320px] px-4 pb-12 pt-20 md:pt-24">
      <div class="grid gap-7 lg:grid-cols-[minmax(0,1fr)_330px] lg:items-end">
        <div class="space-y-5">
          <p class="inline-flex animate-fade-up rounded-full border border-white/40 bg-white/20 px-4 py-1 text-xs tracking-[0.24em] text-[#99f6e4] backdrop-blur">
            {{ settings.heroBadge }}
          </p>
          <h1 class="animate-fade-up text-[clamp(2.3rem,6.1vw,4.3rem)] font-black tracking-[0.08em] text-[#39c5bb] drop-shadow-lg">
            {{ settings.heroTitle }}
          </h1>
          <p class="max-w-2xl animate-fade-up text-base leading-relaxed text-[#b8fff8] md:text-lg">
            {{ settings.heroDescription }}
          </p>

          <div class="animate-fade-up flex flex-wrap gap-2.5">
            <a
              v-for="(action, index) in settings.heroActions"
              :key="`${action.label}-${action.href}`"
              :href="action.href"
              :class="heroActionClasses[index] ?? heroActionClasses[0]"
            >
              {{ action.label }}
            </a>
          </div>

          <div class="grid gap-3 animate-fade-up sm:grid-cols-3">
            <div
              v-for="stat in settings.quickStats"
              :key="`${stat.label}-${stat.value}`"
              class="rounded-2xl border border-white/45 bg-white/16 p-3 text-[#b8fff8] shadow-[0_10px_24px_rgba(15,23,42,0.18)] backdrop-blur-xl"
            >
              <p class="text-[11px] tracking-[0.16em] text-[#99f6e4]">{{ stat.label }}</p>
              <p class="mt-1 text-sm font-semibold">{{ stat.value }}</p>
            </div>
          </div>
        </div>

        <div class="hero-focus-card animate-fade-up rounded-2xl border border-white/50 bg-white/22 p-4 text-[#b8fff8] shadow-[0_18px_40px_rgba(15,23,42,0.22)] backdrop-blur-xl">
          <p class="text-xs tracking-[0.2em] text-[#99f6e4]">{{ settings.focusCard.badge }}</p>
          <h2 class="mt-2 text-lg font-bold text-[#39c5bb]">{{ settings.focusCard.title }}</h2>
          <p class="mt-2 text-sm leading-relaxed text-[#b8fff8]">
            {{ settings.focusCard.description }}
          </p>
          <div class="mt-4 rounded-xl border border-white/30 bg-white/12 px-3 py-2.5 text-xs text-[#99f6e4]">
            {{ settings.focusCard.footnote }}
          </div>
        </div>
      </div>
    </div>

    <a
      href="#content-start"
      class="group absolute bottom-4 left-1/2 z-20 -translate-x-1/2 rounded-full border border-white/40 bg-white/18 px-4 py-1.5 text-xs tracking-[0.18em] text-[#99f6e4] backdrop-blur transition hover:bg-white/24"
      :aria-label="settings.scrollCue.ariaLabel"
    >
      <span class="inline-flex items-center gap-2">
        {{ settings.scrollCue.label }}
        <span class="scroll-cue-dot h-1.5 w-1.5 rounded-full bg-[#39c5bb]" />
      </span>
    </a>
  </header>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useStore } from '@nanostores/vue'

import { blogIndexSettings, hydrateBlogIndexSettings } from '../../stores/blogIndex'

const store = useStore(blogIndexSettings)

const heroActionClasses = [
  'inline-flex items-center rounded-xl border border-[#c084fc]/45 bg-[#f3e8ff] px-4 py-2.5 text-sm font-semibold text-[#8b5cf6] transition hover:-translate-y-0.5 hover:border-[#c084fc]/70',
  'inline-flex items-center rounded-xl border border-white/50 bg-white/18 px-4 py-2.5 text-sm font-semibold text-[#99f6e4] transition hover:-translate-y-0.5 hover:bg-white/24',
  'inline-flex items-center rounded-xl border border-white/50 bg-white/12 px-4 py-2.5 text-sm font-semibold text-[#99f6e4] transition hover:-translate-y-0.5 hover:bg-white/18',
]

const settings = computed(() => store.value)

onMounted(() => {
  void hydrateBlogIndexSettings()
})
</script>
