<template>
  <div ref="frameRef" class="liquid-glass-frame" :class="{ 'is-night': isNight }" :style="frameStyle">
    <div class="liquid-glass-content">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { createLiquidGlass, type LiquidGlassController } from '../../utils/liquidGlass'
import { hydrateThemeMode } from '../../stores/theme'

interface Props {
  width?: string
  maxWidth?: string
  padding?: string
  borderRadius?: number
  cornerSoftness?: number
  displacementStrength?: number
  edgeRefractionStrength?: number
  blur?: number
  contrast?: number
  brightness?: number
  saturate?: number
  interactive?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  maxWidth: '800px',
  padding: '30px 50px',
  borderRadius: 24,
  cornerSoftness: 0.12,
  displacementStrength: 1,
  edgeRefractionStrength: 0.75,
  blur: 0.3,
  contrast: 1.14,
  brightness: 1.04,
  saturate: 1.08,
  interactive: true,
})

const frameRef = ref<HTMLElement | null>(null)
let controller: LiquidGlassController | null = null
const isNight = ref(false)
let themeObserver: MutationObserver | null = null

const frameStyle = computed(() => ({
  width: props.width,
  maxWidth: props.maxWidth,
  padding: props.padding,
  borderRadius: `${props.borderRadius}px`,
}))

const liquidOptions = computed(() => ({
  borderRadius: props.borderRadius,
  cornerSoftness: props.cornerSoftness,
  displacementStrength: props.displacementStrength,
  edgeRefractionStrength: props.edgeRefractionStrength,
  blur: props.blur,
  contrast: props.contrast,
  brightness: props.brightness,
  saturate: props.saturate,
  interactive: props.interactive,
}))

function destroyController() {
  controller?.destroy()
  controller = null
}

function syncThemeFromDocument() {
  if (typeof document === 'undefined') {
    isNight.value = false
    return
  }

  const root = document.documentElement
  const byClass = root.classList.contains('theme-night')
  const byDataAttr = root.getAttribute('data-theme') === 'night'
  isNight.value = byClass || byDataAttr
}

function handleThemeChange() {
  syncThemeFromDocument()
}

function initController() {
  if (!frameRef.value || controller || isNight.value) {
    return
  }

  controller = createLiquidGlass(frameRef.value, liquidOptions.value)
}

onMounted(() => {
  hydrateThemeMode()
  syncThemeFromDocument()
  window.addEventListener('miku-theme-change', handleThemeChange)

  if (typeof MutationObserver !== 'undefined') {
    themeObserver = new MutationObserver(() => syncThemeFromDocument())
    themeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class', 'data-theme'],
    })
  }

  initController()
})

watch(
  liquidOptions,
  // 外部 props 变化时同步更新液态玻璃参数
  (next) => {
    if (isNight.value) {
      return
    }

    if (!controller) {
      initController()
      return
    }

    controller.update(next)
  },
)

watch(
  isNight,
  (night) => {
    if (night) {
      destroyController()
      return
    }

    initController()
  },
)

onBeforeUnmount(() => {
  // 组件卸载时释放监听与滤镜资源
  window.removeEventListener('miku-theme-change', handleThemeChange)
  themeObserver?.disconnect()
  themeObserver = null
  destroyController()
})
</script>

<style scoped>
.liquid-glass-frame {
  position: relative;
  box-sizing: border-box;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.62);
  background: linear-gradient(140deg, rgba(255, 255, 255, 0.45) 0%, rgba(255, 255, 255, 0.12) 100%);
  box-shadow:
    0 10px 36px rgba(10, 18, 34, 0.12),
    inset 0 1px 1px rgba(255, 255, 255, 0.75),
    inset 0 -10px 26px rgba(255, 255, 255, 0.1);
}

.liquid-glass-frame::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(120deg, rgba(255, 255, 255, 0.56) 5%, rgba(255, 255, 255, 0) 40%);
  pointer-events: none;
}

.liquid-glass-content {
  position: relative;
  z-index: 1;
}

.liquid-glass-frame.is-night {
  border: 1px solid rgba(71, 85, 105, 0.55);
  background: rgba(15, 23, 42, 0.92);
  box-shadow: 0 14px 36px rgba(2, 6, 23, 0.38);
}

.liquid-glass-frame.is-night::before {
  display: none;
}
</style>
