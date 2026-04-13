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

// `homeHeroSettings` 是首页主视觉文案的统一状态。
// 组件本身不直接请求接口，而是只消费 store，保持模板层尽量简单。
const settings = useStore(homeHeroSettings)

onMounted(() => {
  // 首页进入时先用本地缓存秒开，避免标题/副标题闪动。
  primeHomeHeroSettingsFromCache()
  // 然后再异步请求后台最新配置，把缓存值校正成服务端版本。
  void hydrateHomeHeroSettings()
})
</script>
