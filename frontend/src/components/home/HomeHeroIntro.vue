<template>
  <div class="flex flex-col items-center">
    <h1 class="hero-title animate-fade-up text-[clamp(2.9rem,7.4vw,5.6rem)] tracking-[0.03em]">
      <span class="hero-text">
        <HeroTitle :text="settings.heroTitle" />
      </span>
    </h1>
    <TypewriterSubtitle :text="settings.heroSubtitle" />
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { onMounted } from 'vue'

import {
  hydrateHomeHeroSettings,
  homeHeroSettings,
  primeHomeHeroSettingsFromCache,
} from '../../stores/homeHero'
import HeroTitle from './HeroTitle.vue'
import TypewriterSubtitle from './TypewriterSubtitle.vue'

const settings = useStore(homeHeroSettings)

onMounted(() => {
  primeHomeHeroSettingsFromCache()
  void hydrateHomeHeroSettings()
})
</script>
