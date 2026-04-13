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

// 顶栏小组件的显示开关来自统一设置中心。
// 这样首页模板不用知道“音乐/天气/时钟是否该展示”，只负责渲染。
const settings = useStore(siteIntegrationsSettings)

// 天气组件只关心最终地点字符串，所以这里先做一次简单计算。
const weatherLocation = computed(() => settings.value.weatherLocation || '')

onMounted(() => {
  // 同样走“缓存优先、接口兜底”的首页设置读取策略。
  primeSiteIntegrationsSettingsFromCache()
  void hydrateSiteIntegrationsSettings()
})
</script>
