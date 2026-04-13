<template>
  <div
    ref="containerRef"
    class="absolute inset-0 overflow-hidden"
    style="perspective: 1200px"
  >
    <div
      ref="layerRef"
      class="parallax-layer absolute"
      :style="{ transform: layerTransform }"
    >
      <img
        v-for="(src, i) in currentHeroImages"
        :key="src"
        :src="src"
        :alt="`${hp.coverAltPrefix}${i + 1}`"
        class="parallax-img absolute inset-0 h-full w-full object-cover"
        :class="mounted && i === currentIndex ? 'is-active' : ''"
        :loading="i === currentIndex ? 'eager' : 'lazy'"
        :fetchpriority="i === currentIndex ? 'high' : 'low'"
        decoding="async"
        draggable="false"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { siteCopy } from '../../content/copy'
import {
  hydrateHomeAssetsSettings,
  primeHomeAssetsSettingsFromCache,
} from '../../stores/homeAssets'

const hp = siteCopy.components.heroParallax

import { heroImages, heroIndex, shuffleHeroImage } from '../../stores/heroImage'

// `heroImages` / `heroIndex` 是首页背景图的共享状态：
// - `HeroParallax` 负责把当前图片显示出来
// - `HeroShuffleBtn` 负责切换索引
// 两边通过同一个 store 协作，所以不会出现按钮切了图、背景没变的情况。
const $heroIndex = useStore(heroIndex)
const $heroImages = useStore(heroImages)
const currentIndex = ref($heroIndex.value)
const mounted = ref(false)
const currentHeroImages = computed(() => $heroImages.value)

watch($heroIndex, (v) => {
  // store 里当前图片索引一旦变化，立刻同步到本组件的展示状态。
  currentIndex.value = v
})

watch($heroImages, (images) => {
  // 如果后台把图片列表改短了，旧索引可能越界，这里要兜底归零。
  if (currentIndex.value >= images.length) {
    currentIndex.value = 0
  }
})

const containerRef = ref<HTMLElement | null>(null)
const layerRef = ref<HTMLElement | null>(null)

const RANGE = 10
const SHIFT = 18

// 这里维护两组值：
// - rx/ry/tx/ty：当前真正应用到 DOM 上的值
// - targetRx/...：鼠标移动后“目标应该去到哪里”
// 再用 `lerp()` 在每一帧慢慢逼近，从而做出柔和的缓动效果。
let rx = 0
let ry = 0
let tx = 0
let ty = 0
let targetRx = 0
let targetRy = 0
let targetTx = 0
let targetTy = 0
let rafId = 0

const layerTransform = ref('rotateX(0deg) rotateY(0deg) translateX(0px) translateY(0px)')

function onMouseMove(e: MouseEvent) {
  // 把鼠标位置归一化成 -1 到 1 的比例，
  // 这样无论屏幕多大，视差公式都能复用。
  const w = window.innerWidth
  const h = window.innerHeight
  const xRatio = (e.clientX / w - 0.5) * 2
  const yRatio = (e.clientY / h - 0.5) * 2

  targetRy = xRatio * RANGE
  targetRx = -yRatio * RANGE
  targetTx = -xRatio * SHIFT
  targetTy = -yRatio * SHIFT
}

function onMouseLeave() {
  // 鼠标离开后，把目标值重置到中心点，背景就会缓慢回正。
  targetRx = 0
  targetRy = 0
  targetTx = 0
  targetTy = 0
}

function lerp(a: number, b: number, t: number) {
  // 线性插值：每帧走一小步，而不是瞬间跳到目标值。
  return a + (b - a) * t
}

function loop() {
  rx = lerp(rx, targetRx, 0.06)
  ry = lerp(ry, targetRy, 0.06)
  tx = lerp(tx, targetTx, 0.06)
  ty = lerp(ty, targetTy, 0.06)

  layerTransform.value =
    `rotateX(${rx.toFixed(2)}deg) rotateY(${ry.toFixed(2)}deg) translateX(${tx.toFixed(1)}px) translateY(${ty.toFixed(1)}px)`

  rafId = requestAnimationFrame(loop)
}

onMounted(() => {
  // 背景图资源同样先从缓存恢复，再异步同步后台配置。
  primeHomeAssetsSettingsFromCache()
  void hydrateHomeAssetsSettings()
  mounted.value = true
  // 首页第一次进入时随机挑一张图，避免每次都从同一张开始。
  shuffleHeroImage()
  currentIndex.value = heroIndex.get()

  const el = containerRef.value
  if (el) {
    el.addEventListener('mousemove', onMouseMove, { passive: true })
    el.addEventListener('mouseleave', onMouseLeave)
  }
  rafId = requestAnimationFrame(loop)
})

onUnmounted(() => {
  // 清理事件和动画帧，避免页面切走后还在后台持续执行。
  const el = containerRef.value
  if (el) {
    el.removeEventListener('mousemove', onMouseMove)
    el.removeEventListener('mouseleave', onMouseLeave)
  }
  if (rafId) cancelAnimationFrame(rafId)
})
</script>

<style scoped>
.parallax-layer {
  inset: -30px;
  transform-origin: 50% 50%;
  transform-style: preserve-3d;
  will-change: transform;
}

.parallax-img {
  opacity: 0;
  transition: opacity 1.2s cubic-bezier(0.4, 0, 0.2, 1);
  will-change: opacity;
}

.parallax-img.is-active {
  opacity: 1;
}
</style>
