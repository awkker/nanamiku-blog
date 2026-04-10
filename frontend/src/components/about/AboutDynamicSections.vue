<template>
  <div class="space-y-4 sm:space-y-5">
    <section class="relative mx-auto max-w-[1320px] px-4 sm:px-6 lg:px-8">
      <div class="grid gap-4 md:grid-cols-3">
        <LiquidGlassCard
          v-for="card in settings.introCards"
          :key="`${card.title}-${card.description}`"
          padding="20px"
          maxWidth="100%"
        >
          <h2 class="text-sm font-semibold text-slate-900">{{ card.title }}</h2>
          <p class="mt-2 text-sm leading-relaxed text-slate-600">{{ card.description }}</p>
        </LiquidGlassCard>
      </div>
    </section>

    <section id="timeline" class="relative mx-auto max-w-[1320px] px-4 sm:px-6 lg:px-8">
      <LiquidGlassCard padding="24px" maxWidth="100%">
        <div class="mb-4">
          <h2 class="text-xl font-bold text-slate-900">{{ copy.timeline.title }}</h2>
          <p class="mt-1 text-sm text-slate-500">{{ copy.timeline.subtitle }}</p>
        </div>
        <ol class="timeline-list relative space-y-4 pl-5">
          <li
            v-for="item in settings.milestones"
            :key="`${item.year}-${item.title}`"
            class="relative"
          >
            <span class="absolute -left-[21px] top-2 h-3 w-3 rounded-full border border-white bg-miku shadow-[0_0_0_3px_rgba(57,197,187,0.14)]" />
            <div class="rounded-xl border border-white/65 bg-white/62 p-4">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <h3 class="text-base font-semibold text-slate-900">{{ item.title }}</h3>
                <span class="rounded-full bg-miku-soft px-2.5 py-0.5 text-xs font-semibold text-miku">{{ item.year }}</span>
              </div>
              <p class="mt-2 text-sm leading-relaxed text-slate-600">{{ item.summary }}</p>
              <p class="mt-1.5 text-xs text-slate-500">{{ copy.timeline.resultPrefix }}{{ item.result }}</p>
            </div>
          </li>
        </ol>
      </LiquidGlassCard>
    </section>

    <section id="skills" class="relative mx-auto max-w-[1320px] px-4 sm:px-6 lg:px-8">
      <div class="grid gap-4 md:grid-cols-3">
        <LiquidGlassCard
          v-for="(group, index) in settings.capabilityGroups"
          :key="`${group.title}-${group.desc}`"
          padding="20px"
          maxWidth="100%"
        >
          <h3 class="text-base font-semibold text-slate-900">{{ group.title }}</h3>
          <p class="mt-2 text-sm leading-relaxed text-slate-600">{{ group.desc }}</p>
          <div class="mt-3 flex flex-wrap gap-1.5">
            <span
              v-for="item in group.stack"
              :key="`${group.title}-${item}`"
              class="rounded-full border px-2.5 py-1 text-xs"
              :class="capabilityChipClass(index)"
            >
              {{ item }}
            </span>
          </div>
        </LiquidGlassCard>
      </div>
    </section>

    <section id="projects" class="relative mx-auto max-w-[1320px] px-4 sm:px-6 lg:px-8">
      <div class="mb-3">
        <h2 class="text-xl font-bold text-slate-900">{{ copy.projectsSection.title }}</h2>
        <p class="mt-1 text-sm text-slate-500">{{ copy.projectsSection.subtitle }}</p>
      </div>
      <div class="grid gap-4 lg:grid-cols-3">
        <LiquidGlassCard
          v-for="project in settings.featuredProjects"
          :key="`${project.name}-${project.href}`"
          padding="20px"
          maxWidth="100%"
        >
          <h3 class="text-base font-semibold text-slate-900">{{ project.name }}</h3>
          <p class="mt-2 text-sm text-slate-600">{{ project.focus }}</p>
          <p class="mt-2 text-xs leading-relaxed text-slate-500">{{ project.role }}</p>
          <p class="mt-2 rounded-lg border border-white/70 bg-white/65 px-2.5 py-2 text-xs text-slate-600">{{ project.metric }}</p>
          <a
            :href="project.href"
            class="mt-3 inline-flex items-center rounded-lg border border-miku/35 bg-miku-soft px-3 py-1.5 text-xs font-semibold text-miku transition hover:border-miku/55"
          >
            {{ copy.projectsSection.cta }}
          </a>
        </LiquidGlassCard>
      </div>
    </section>

    <section class="relative mx-auto max-w-[1320px] px-4 sm:px-6 lg:px-8">
      <div class="grid gap-4 md:grid-cols-3">
        <LiquidGlassCard padding="20px" maxWidth="100%">
          <h3 class="text-sm font-semibold text-slate-900">{{ copy.monthlyGoalsTitle }}</h3>
          <ul class="mt-3 space-y-2">
            <li
              v-for="goal in settings.monthlyGoals"
              :key="goal"
              class="flex items-start gap-2 text-xs text-slate-600"
            >
              <span class="mt-1 h-1.5 w-1.5 rounded-full bg-miku" />
              <span>{{ goal }}</span>
            </li>
          </ul>
        </LiquidGlassCard>

        <LiquidGlassCard padding="20px" maxWidth="100%">
          <h3 class="text-sm font-semibold text-slate-900">{{ copy.listeningTitle }}</h3>
          <ul class="mt-3 space-y-2 text-xs text-slate-600">
            <li
              v-for="song in settings.listeningNow"
              :key="song"
              class="rounded-lg border border-white/65 bg-white/62 px-2.5 py-2"
            >
              {{ song }}
            </li>
          </ul>
        </LiquidGlassCard>

        <LiquidGlassCard padding="20px" maxWidth="100%">
          <h3 class="text-sm font-semibold text-slate-900">{{ copy.signature.title }}</h3>
          <p class="mt-3 text-sm leading-relaxed text-slate-600">
            {{ settings.signature.description }}
          </p>
          <p class="mt-2 text-xs text-slate-500">{{ settings.signature.footer }}</p>
        </LiquidGlassCard>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted } from 'vue'

import { siteCopy } from '../../content/copy'
import {
  aboutPageSettings,
  hydrateAboutPageSettings,
} from '../../stores/aboutPage'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

const aboutStore = useStore(aboutPageSettings)
const copy = siteCopy.aboutPage

const settings = computed(() => aboutStore.value)

function capabilityChipClass(index: number) {
  if (index === 0) {
    return 'border-miku/35 bg-miku-soft text-miku'
  }
  if (index === 1) {
    return 'border-slate-200 bg-white/70 text-slate-600'
  }
  return 'border-[#c084fc]/35 bg-[#f3e8ff] text-[#8b5cf6]'
}

onMounted(() => {
  void hydrateAboutPageSettings()
})
</script>

<style scoped>
  .timeline-list::before {
    content: '';
    position: absolute;
    left: 9px;
    top: 8px;
    bottom: 8px;
    width: 1px;
    pointer-events: none;
    background: linear-gradient(
      to bottom,
      rgba(57, 197, 187, 0.6),
      rgba(192, 132, 252, 0.55),
      rgba(57, 197, 187, 0.25)
    );
  }
</style>
