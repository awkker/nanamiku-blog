<template>
  <div class="flex items-center gap-1 text-[11px] min-[380px]:gap-1.5 sm:gap-3">
    <MusicPlayer v-if="settings.showMusic" />
    <span v-if="settings.showWeather" class="hidden min-[360px]:inline-flex">
      <WeatherDisplay :location="weatherLocation" />
    </span>
    <span v-if="settings.showClock" class="hidden min-[430px]:inline-flex">
      <SystemClock />
    </span>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted } from 'vue'

import {
  hydrateSiteIntegrationsSettings,
  primeSiteIntegrationsSettingsFromCache,
  siteIntegrationsSettings,
} from '../../stores/siteIntegrations'
import MusicPlayer from './MusicPlayer.vue'
import SystemClock from './SystemClock.vue'
import WeatherDisplay from './WeatherDisplay.vue'

const settings = useStore(siteIntegrationsSettings)

const weatherLocation = computed(() => settings.value.weatherLocation || '')

onMounted(() => {
  primeSiteIntegrationsSettingsFromCache()
  void hydrateSiteIntegrationsSettings()
})
</script>
