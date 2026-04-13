<template>
  <button
    class="group flex w-12 flex-col items-center gap-1 rounded-2xl py-2 text-[10px] text-[#39c5bb]/90 transition duration-300 hover:scale-110 hover:bg-white/20 sm:w-16 sm:text-xs sm:hover:scale-125"
    :title="btnTitle"
    @click="onShuffle"
  >
    <span
      class="flex h-8 w-8 items-center justify-center rounded-xl border border-white/30 bg-white/15 transition-transform duration-500 sm:h-9 sm:w-9"
      :class="{ 'rotate-180': spinning }"
    >
      <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.5] sm:h-5 sm:w-5">
        <path d="M16 3h5v5M4 20L21 3M21 16v5h-5M4 4l17 17" />
      </svg>
    </span>
    <span>{{ hsb.label }}</span>
    <span class="h-1.5 w-1.5 rounded-full bg-white opacity-0" />
  </button>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref, watch } from 'vue'
import { siteCopy } from '../../content/copy'
import {
  hydrateHomeAssetsSettings,
  primeHomeAssetsSettingsFromCache,
} from '../../stores/homeAssets'

const hsb = siteCopy.components.heroShuffleBtn

import { heroImages, heroIndex, shuffleHeroImage } from '../../stores/heroImage'

// 这个按钮不自己维护图片列表，而是直接读取和 `HeroParallax` 相同的 store。
// 所以它只是“发出切换动作”，真正展示背景的仍然是背景组件。
const $heroIndex = useStore(heroIndex)
const $heroImages = useStore(heroImages)
const currentIdx = ref(0)
const mounted = ref(false)
const spinning = ref(false)
const total = computed(() => $heroImages.value.length || 1)

const btnTitle = computed(() =>
  // title 会展示“当前第几张 / 总共几张”，鼠标悬停时更容易理解按钮作用。
  mounted.value ? `${hsb.titlePrefix} (${currentIdx.value + 1}/${total.value})` : hsb.titlePrefix
)

watch($heroIndex, (v) => {
  currentIdx.value = v
})

onMounted(() => {
  primeHomeAssetsSettingsFromCache()
  void hydrateHomeAssetsSettings()
  mounted.value = true
  currentIdx.value = heroIndex.get()
})

function onShuffle() {
  // 点击时给图标一个短暂旋转反馈，但真正的换图动作还是交给 store 完成。
  spinning.value = true
  shuffleHeroImage()
  setTimeout(() => {
    spinning.value = false
  }, 500)
}
</script>
