<template>
  <div class="absolute bottom-2 left-1/2 z-20 w-[calc(100vw-0.35rem)] max-w-[900px] -translate-x-1/2 px-0.5 sm:bottom-2.5 sm:px-0">
    <div
      ref="dockRef"
      class="relative mx-auto flex w-fit max-w-full justify-center overflow-visible"
      @pointermove="onPointerMove"
      @pointerleave="onPointerLeave"
    >
      <div class="relative grid w-fit max-w-full overflow-visible">
        <LiquidGlassCard
          width="100%"
          maxWidth="100%"
          padding="0"
          :border-radius="28"
          class="dock-shell pointer-events-none col-start-1 row-start-1 h-full self-stretch min-h-[3.95rem] sm:min-h-[5.3rem]"
        >
          <div aria-hidden="true" class="h-full w-full" />
        </LiquidGlassCard>

        <div class="relative z-[1] col-start-1 row-start-1 flex w-fit max-w-full overflow-visible px-[8px] pb-[7px] pt-[8px]">
          <ul class="relative flex items-end justify-center gap-0.5 pt-0.5 sm:gap-1.5 sm:pt-1">
            <li
              v-for="(entry, index) in dockEntries"
              :key="entry.id"
              :ref="(el) => setItemRef(el, index)"
              class="relative flex shrink-0 justify-center"
              :style="getItemStyle(index)"
            >
              <div
                v-if="activeIndex === index"
                class="pointer-events-none absolute bottom-full left-1/2 z-20 mb-3 w-max max-w-none -translate-x-1/2 whitespace-nowrap rounded-2xl border border-white/80 bg-white/90 px-3 py-1.5 text-[11px] font-medium tracking-[0.01em] text-slate-600 shadow-[0_14px_32px_rgba(15,23,42,0.14)] backdrop-blur-md sm:mb-3.5"
                :style="getTooltipStyle(index)"
              >
                <span>{{ entry.label }}</span>
                <span class="absolute left-1/2 top-full h-2.5 w-2.5 -translate-x-1/2 -translate-y-1/2 rotate-45 border-b border-r border-white/75 bg-white/88" />
              </div>

              <a
                v-if="entry.kind === 'link'"
                :href="entry.href"
                class="group relative flex w-[2.3rem] touch-manipulation flex-col items-center justify-end gap-1 rounded-[1.15rem] px-0.5 pb-0.5 pt-0.5 outline-none transition-[filter] duration-200 focus-visible:ring-2 focus-visible:ring-miku/65 focus-visible:ring-offset-2 focus-visible:ring-offset-white/55 sm:w-[4.45rem] sm:gap-1.5"
                :title="entry.title || entry.label"
                :aria-label="entry.label"
                :aria-current="entry.href === currentPath ? 'page' : undefined"
                @focus="focusedId = entry.id"
                @blur="focusedId = null"
              >
                <span class="sr-only">{{ entry.label }}</span>
                <span class="dock-icon relative flex h-[2.1rem] w-[2.1rem] items-center justify-center rounded-[0.92rem] border sm:h-[3.1rem] sm:w-[3.1rem] sm:rounded-[1.08rem]" :style="getIconStyle(index, entry)">
                  <span class="pointer-events-none absolute inset-[1px] rounded-[inherit] bg-[linear-gradient(180deg,rgba(255,255,255,0.16),rgba(255,255,255,0))]" />
                  <span class="pointer-events-none absolute left-1/2 top-[10%] h-[18%] w-[72%] -translate-x-1/2 rounded-full bg-white/50 blur-[4px] sm:blur-[6px]" />
                  <svg viewBox="0 0 24 24" class="relative z-[1] h-[1rem] w-[1rem] fill-none stroke-current stroke-[1.75] sm:h-[1.38rem] sm:w-[1.38rem]">
                    <path :d="getIconPath(entry.icon)" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </span>
                <span class="pointer-events-none absolute bottom-[0.5rem] left-1/2 h-[0.34rem] w-[1.1rem] -translate-x-1/2 rounded-full bg-slate-900/12 blur-[5px] sm:bottom-[0.8rem] sm:h-[0.42rem] sm:w-[1.65rem] sm:blur-[7px]" :style="getReflectionStyle(index)" />
                <span class="relative z-[1] h-1 w-1 rounded-full sm:h-1.5 sm:w-1.5" :style="getIndicatorStyle(index, entry)" />
              </a>

              <button
                v-else
                type="button"
                class="group relative flex w-[2.3rem] touch-manipulation flex-col items-center justify-end gap-1 rounded-[1.15rem] px-0.5 pb-0.5 pt-0.5 outline-none transition-[filter] duration-200 focus-visible:ring-2 focus-visible:ring-miku/65 focus-visible:ring-offset-2 focus-visible:ring-offset-white/55 sm:w-[4.45rem] sm:gap-1.5"
                :title="entry.title || entry.label"
                :aria-label="entry.label"
                @click="handleActionClick(entry)"
                @focus="focusedId = entry.id"
                @blur="focusedId = null"
              >
                <span class="sr-only">{{ entry.label }}</span>
                <span class="dock-icon relative flex h-[2.1rem] w-[2.1rem] items-center justify-center rounded-[0.92rem] border sm:h-[3.1rem] sm:w-[3.1rem] sm:rounded-[1.08rem]" :style="getIconStyle(index, entry)">
                  <span class="pointer-events-none absolute inset-[1px] rounded-[inherit] bg-[linear-gradient(180deg,rgba(255,255,255,0.16),rgba(255,255,255,0))]" />
                  <span class="pointer-events-none absolute left-1/2 top-[10%] h-[18%] w-[72%] -translate-x-1/2 rounded-full bg-white/50 blur-[4px] sm:blur-[6px]" />
                  <svg
                    viewBox="0 0 24 24"
                    class="relative z-[1] h-[1rem] w-[1rem] fill-none stroke-current stroke-[1.75] transition-transform duration-500 sm:h-[1.38rem] sm:w-[1.38rem]"
                    :class="{ 'rotate-180': shuffleSpinning }"
                  >
                    <path :d="getIconPath(entry.icon)" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </span>
                <span class="pointer-events-none absolute bottom-[0.5rem] left-1/2 h-[0.34rem] w-[1.1rem] -translate-x-1/2 rounded-full bg-slate-900/12 blur-[5px] sm:bottom-[0.8rem] sm:h-[0.42rem] sm:w-[1.65rem] sm:blur-[7px]" :style="getReflectionStyle(index)" />
                <span class="relative z-[1] h-1 w-1 rounded-full sm:h-1.5 sm:w-1.5" :style="getIndicatorStyle(index, entry)" />
              </button>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { siteCopy } from '../../content/copy'
import { AUTH_STATE_CHANGED_EVENT, ensureSessionUser } from '../../lib/auth-session'
import { heroImages, heroIndex, shuffleHeroImage } from '../../stores/heroImage'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

type DockKind = 'link' | 'action'

interface DockEntry {
  id: string
  label: string
  icon: string
  kind: DockKind
  href?: string
  title?: string
}

interface DockPalette {
  from: string
  to: string
  shadow: string
  stroke: string
  glyph: string
}

const homeCopy = siteCopy.home
const shuffleCopy = siteCopy.components.heroShuffleBtn

const $heroIndex = useStore(heroIndex)
const $heroImages = useStore(heroImages)

const dockRef = ref<HTMLElement | null>(null)
const dockGlassFrameRef = ref<HTMLElement | null>(null)
const itemRefs = ref<(HTMLElement | null)[]>([])
const itemCenters = ref<number[]>([])
const pointerX = ref<number | null>(null)
const currentPath = ref('/')
const isAuthed = ref(false)
const authResolved = ref(false)
const canTrackPointer = ref(false)
const focusedId = ref<string | null>(null)
const shuffleSpinning = ref(false)

let resizeObserver: ResizeObserver | null = null
let shuffleTimer: number | null = null

const ICON_PATHS: Record<string, string> = {
  person:
    'M12 12c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm0 2c-3.87 0-7 2.01-7 4.5V21h14v-2.5c0-2.49-3.13-4.5-7-4.5z',
  book: 'M18 2H8a3 3 0 0 0-3 3v14a3 3 0 0 1 3-3h10V2zm2 16H8a1 1 0 0 0-1 1v1h13v-2z',
  moments: 'M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z',
  message: 'M4 5h16a2 2 0 012 2v9a2 2 0 01-2 2H8l-4 3v-3H4a2 2 0 01-2-2V7a2 2 0 012-2zm4 5h8M8 13h6',
  link: 'M10 13a5 5 0 007.07 0l2.12-2.12a5 5 0 10-7.07-7.07L10 5M14 11a5 5 0 00-7.07 0L4.81 13.12a5 5 0 007.07 7.07L14 19',
  login: 'M10 17l5-5-5-5M15 12H3M17 4h2a2 2 0 012 2v12a2 2 0 01-2 2h-2',
  admin: 'M4 5h16v14H4V5zm4 3h8M8 11h8M8 15h5',
  shuffle: 'M16 3h5v5M4 20L21 3M21 16v5h-5M4 4l17 17',
}

// 这里统一收口为初音绿系。
// 后面如果你想让图标再浅一点或更通透，优先改 `from / to / shadow`。
const MIKU_DOCK_PALETTE: DockPalette = {
  from: '#83efe5',
  to: '#39c5bb',
  shadow: 'rgba(57, 197, 187, 0.34)',
  stroke: 'rgba(238, 255, 252, 0.88)',
  glyph: 'rgba(255, 255, 255, 0.97)',
}

// 下面这几个常量是 Dock 手感的核心旋钮。
// 如果你想让效果更像 macOS，可以优先从这里调：
// - RANGE 决定“鼠标离多远还会影响相邻图标”
// - SCALE_BOOST 决定中心图标最多能放大多少
// - LIFT 决定图标向上抬起的高度
const DOCK_EFFECT_RANGE = 168
const DOCK_MAX_SCALE_BOOST = 0.88
const DOCK_MAX_LIFT = 22
const DOCK_LABEL_RANGE = 86

const shuffleTitle = computed(() => {
  const total = Math.max($heroImages.value.length, 1)
  const current = Math.min($heroIndex.value + 1, total)
  return `${shuffleCopy.titlePrefix} (${current}/${total})`
})

const dockEntries = computed<DockEntry[]>(() => {
  const visibleLinks = homeCopy.dockItems
    .filter((item) => {
      if (item.auth === 'authed') {
        return authResolved.value && isAuthed.value
      }
      if (item.auth === 'guest') {
        return authResolved.value && !isAuthed.value
      }
      return true
    })
    .map((item) => ({
      id: `dock:${item.href}`,
      label: item.name,
      icon: item.icon,
      kind: 'link' as const,
      href: item.href,
    }))

  return [
    ...visibleLinks,
    {
      id: 'dock:shuffle',
      label: shuffleCopy.label,
      icon: 'shuffle',
      kind: 'action',
      title: shuffleTitle.value,
    },
  ]
})

const activeIndex = computed(() => {
  if (focusedId.value) {
    return dockEntries.value.findIndex((entry) => entry.id === focusedId.value)
  }

  if (!canTrackPointer.value || pointerX.value === null || itemCenters.value.length === 0) {
    return -1
  }

  let nearestIndex = -1
  let nearestDistance = Number.POSITIVE_INFINITY

  itemCenters.value.forEach((center, index) => {
    const distance = Math.abs(pointerX.value! - center)
    if (distance < nearestDistance) {
      nearestDistance = distance
      nearestIndex = index
    }
  })

  return nearestDistance <= DOCK_LABEL_RANGE ? nearestIndex : -1
})

const dockMetrics = computed(() =>
  dockEntries.value.map((_, index) => {
    const center = itemCenters.value[index]

    if (!canTrackPointer.value || pointerX.value === null || typeof center !== 'number') {
      return {
        intensity: 0,
        scale: 1,
        lift: 0,
        tilt: 0,
      }
    }

    const offset = pointerX.value - center
    const distance = Math.abs(offset)
    const normalized = Math.max(0, 1 - distance / DOCK_EFFECT_RANGE)

    // 这里用 smoothstep 让曲线更顺。
    // 如果直接线性放大，图标在靠近与离开时会显得比较“硬”。
    const intensity = normalized * normalized * (3 - 2 * normalized)

    return {
      intensity,
      scale: 1 + intensity * DOCK_MAX_SCALE_BOOST,
      lift: intensity * DOCK_MAX_LIFT,
      tilt: (-offset / DOCK_EFFECT_RANGE) * 6 * intensity,
    }
  }),
)

function updatePointerCapability() {
  if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') {
    canTrackPointer.value = false
    return
  }

  canTrackPointer.value = window.matchMedia('(hover: hover) and (pointer: fine)').matches
  if (!canTrackPointer.value) {
    pointerX.value = null
  }
}

function setItemRef(el: Element | null, index: number) {
  itemRefs.value[index] = el as HTMLElement | null
}

function measureItemCenters() {
  if (!dockRef.value) {
    itemCenters.value = []
    return
  }

  const dockRect = dockRef.value.getBoundingClientRect()
  itemCenters.value = itemRefs.value.map((el) => {
    if (!el) return 0
    const rect = el.getBoundingClientRect()
    return rect.left - dockRect.left + rect.width / 2
  })
}

function ensureDockGlassFrame() {
  if (!dockGlassFrameRef.value && dockRef.value) {
    dockGlassFrameRef.value = dockRef.value.querySelector('.dock-shell')
  }

  return dockGlassFrameRef.value
}

async function syncAuthState() {
  isAuthed.value = Boolean(await ensureSessionUser())
  authResolved.value = true
  await nextTick()
  measureItemCenters()
}

function onPointerMove(event: PointerEvent) {
  if (!canTrackPointer.value || !dockRef.value) {
    return
  }

  const dockRect = dockRef.value.getBoundingClientRect()
  pointerX.value = event.clientX - dockRect.left

  const glassFrame = ensureDockGlassFrame()
  if (glassFrame) {
    glassFrame.dispatchEvent(
      new MouseEvent('mousemove', {
        clientX: event.clientX,
        clientY: event.clientY,
      }),
    )
  }
}

function onPointerLeave() {
  pointerX.value = null
}

function handleActionClick(entry: DockEntry) {
  if (entry.icon !== 'shuffle') {
    return
  }

  shuffleSpinning.value = true
  shuffleHeroImage()

  if (shuffleTimer) {
    window.clearTimeout(shuffleTimer)
  }

  shuffleTimer = window.setTimeout(() => {
    shuffleSpinning.value = false
    shuffleTimer = null
  }, 520)
}

function getPalette(icon: string): DockPalette {
  return MIKU_DOCK_PALETTE
}

function getIconPath(icon: string) {
  return ICON_PATHS[icon] || ICON_PATHS.shuffle
}

function getItemStyle(index: number) {
  const metrics = dockMetrics.value[index]
  const baseZIndex = activeIndex.value === index ? 80 : 30

  return {
    zIndex: String(baseZIndex + Math.round(metrics.intensity * 30)),
  }
}

function getIconStyle(index: number, entry: DockEntry) {
  const metrics = dockMetrics.value[index]
  const palette = getPalette(entry.icon)
  const shadowSpread = 22 + metrics.intensity * 18
  const shadowLift = 12 + metrics.intensity * 18

  // 真正决定“macOS 味道”的是这里：
  // - transform 负责放大、上抬、轻微倾斜
  // - boxShadow 负责把中心图标做得更像从底座里弹出来
  return {
    transform: `translate3d(0, ${(-metrics.lift).toFixed(1)}px, 0) scale(${metrics.scale.toFixed(3)}) rotate(${metrics.tilt.toFixed(2)}deg)`,
    background: `linear-gradient(180deg, ${palette.from} 0%, ${palette.to} 100%)`,
    borderColor: palette.stroke,
    color: palette.glyph,
    boxShadow: `0 ${shadowLift.toFixed(1)}px ${shadowSpread.toFixed(1)}px ${palette.shadow}, inset 0 1px 0 rgba(255,255,255,0.78), inset 0 -10px 18px rgba(15,23,42,0.16)`,
  }
}

function getReflectionStyle(index: number) {
  const metrics = dockMetrics.value[index]
  const opacity = 0.18 + metrics.intensity * 0.28
  const scale = 0.85 + metrics.intensity * 0.36

  return {
    opacity: opacity.toFixed(3),
    transform: `translateX(-50%) scale(${scale.toFixed(3)})`,
  }
}

function getTooltipStyle(index: number) {
  const metrics = dockMetrics.value[index]
  const baseOffset = 16

  return {
    marginBottom: `${(baseOffset + metrics.lift).toFixed(1)}px`,
  }
}

function getIndicatorStyle(index: number, entry: DockEntry) {
  const metrics = dockMetrics.value[index]
  const isCurrent = Boolean(entry.href && entry.href === currentPath.value)
  const isActive = activeIndex.value === index || isCurrent
  const opacity = isActive ? 0.96 : metrics.intensity * 0.35
  const scale = isActive ? 1 : 0.76 + metrics.intensity * 0.38

  return {
    opacity: opacity.toFixed(3),
    transform: `scale(${scale.toFixed(3)})`,
    background: isCurrent ? 'rgba(57, 197, 187, 0.96)' : 'rgba(255, 255, 255, 0.92)',
    boxShadow: isActive
      ? '0 0 12px rgba(57, 197, 187, 0.42)'
      : '0 0 10px rgba(255, 255, 255, 0.22)',
  }
}

watch(
  dockEntries,
  async () => {
    // 登录态切换后，Dock 项数量会变化，所以这里要重新量一次中心点。
    itemRefs.value = itemRefs.value.slice(0, dockEntries.value.length)
    await nextTick()
    measureItemCenters()
  },
  { immediate: true },
)

onMounted(async () => {
  currentPath.value = window.location.pathname
  updatePointerCapability()
  await syncAuthState()

  window.addEventListener('resize', updatePointerCapability, { passive: true })
  window.addEventListener('resize', measureItemCenters, { passive: true })
  window.addEventListener(AUTH_STATE_CHANGED_EVENT, syncAuthState)

  if (dockRef.value && typeof ResizeObserver !== 'undefined') {
    resizeObserver = new ResizeObserver(() => measureItemCenters())
    resizeObserver.observe(dockRef.value)
  }

  await nextTick()
  ensureDockGlassFrame()
  measureItemCenters()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updatePointerCapability)
  window.removeEventListener('resize', measureItemCenters)
  window.removeEventListener(AUTH_STATE_CHANGED_EVENT, syncAuthState)
  resizeObserver?.disconnect()
  resizeObserver = null

  if (shuffleTimer) {
    window.clearTimeout(shuffleTimer)
    shuffleTimer = null
  }
})
</script>

<style scoped>
.dock-shell {
  transform: translateZ(0);
  border-color: rgba(225, 255, 251, 0.76);
  background: linear-gradient(
    145deg,
    rgba(241, 255, 253, 0.52) 0%,
    rgba(212, 252, 247, 0.28) 38%,
    rgba(57, 197, 187, 0.14) 100%
  );
  box-shadow:
    0 14px 34px rgba(9, 77, 73, 0.16),
    inset 0 1px 0 rgba(255, 255, 255, 0.76),
    inset 0 -10px 22px rgba(57, 197, 187, 0.08);
}

.dock-shell::before {
  background: linear-gradient(120deg, rgba(255, 255, 255, 0.48) 8%, rgba(255, 255, 255, 0) 40%);
}

.dock-shell::after {
  display: none;
}

.dock-icon {
  transform-origin: center bottom;
  transition:
    transform 180ms cubic-bezier(0.22, 1, 0.36, 1),
    box-shadow 180ms cubic-bezier(0.22, 1, 0.36, 1),
    filter 180ms ease;
  will-change: transform;
}

.dock-icon svg {
  filter: drop-shadow(0 1px 1px rgba(255, 255, 255, 0.2));
}

@media (prefers-reduced-motion: reduce) {
  .dock-icon,
  .dock-icon svg {
    transition-duration: 0.01ms;
  }
}
</style>
