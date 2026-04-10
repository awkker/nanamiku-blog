<template>
  <section
    aria-hidden="true"
    class="relative isolate overflow-hidden rounded-[32px] border border-white/70 bg-white/45 p-4 shadow-[0_24px_80px_rgba(15,23,42,0.10)] backdrop-blur-xl sm:p-6 lg:h-[calc(100vh-3rem)] lg:min-h-[720px] lg:p-8"
  >
    <div
      class="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(57,197,187,0.20),transparent_34%),radial-gradient(circle_at_top_right,rgba(192,132,252,0.18),transparent_24%),linear-gradient(180deg,rgba(255,255,255,0.82),rgba(255,255,255,0.38))]"
    />
    <div
      class="absolute inset-0 bg-[linear-gradient(to_right,rgba(148,163,184,0.16)_1px,transparent_1px),linear-gradient(to_bottom,rgba(148,163,184,0.16)_1px,transparent_1px)] bg-[size:78px_78px] opacity-55"
    />
    <div class="absolute -left-12 top-8 h-44 w-44 rounded-full bg-miku/18 blur-3xl" />
    <div class="absolute -right-14 bottom-6 h-56 w-56 rounded-full bg-[#c084fc]/18 blur-3xl" />

    <div
      ref="sceneRef"
      class="relative h-full min-h-[460px] touch-none select-none overflow-hidden rounded-[28px] border border-white/60 bg-white/12 shadow-[inset_0_1px_0_rgba(255,255,255,0.72)] sm:min-h-[560px]"
    >
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.24),transparent_58%)]" />
      <div class="absolute left-1/2 top-[53%] h-[60%] w-[60%] -translate-x-1/2 -translate-y-1/2 rounded-full bg-miku/8 blur-3xl" />

      <div class="pointer-events-none absolute left-1/2 top-[49%] w-[78%] min-w-[18rem] max-w-[30rem] -translate-x-1/2 -translate-y-1/2 sm:w-[70%] lg:w-[72%]">
        <div class="absolute inset-0 rounded-[40px] bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.66),transparent_56%)] blur-2xl" />
        <div class="relative overflow-hidden rounded-[34px] border border-white/75 bg-white/36 shadow-[0_28px_70px_rgba(15,23,42,0.14)]">
          <div class="aspect-[4/5]">
            <img src="/picture/login.jpg" alt="" draggable="false" class="h-full w-full object-cover object-center" />
          </div>
          <div class="absolute inset-0 bg-[linear-gradient(180deg,rgba(255,255,255,0.12),transparent_40%,rgba(255,255,255,0.18))]" />
        </div>
      </div>

      <div
        v-for="card in cards"
        :key="card.id"
        class="absolute left-0 top-0 will-change-transform"
        :style="cardStyle(card)"
      >
        <div
          class="flex h-full w-full cursor-grab items-center justify-center rounded-[26px] border backdrop-blur-xl transition-transform duration-200 active:cursor-grabbing"
          :class="dragState.cardId === card.id ? 'scale-[1.03]' : 'scale-100'"
          :style="cardSurfaceStyle(card)"
          @pointerdown.stop.prevent="handleCardPointerDown(card.id, $event)"
        >
          <div class="absolute inset-[10px] rounded-[20px] border border-white/45 bg-white/18" />
          <div class="absolute inset-x-4 top-3 h-px bg-[linear-gradient(90deg,transparent,rgba(255,255,255,0.72),transparent)]" />

          <div
            class="relative z-10 flex items-center justify-center"
            :class="card.assetSrc ? 'h-[3.1rem] w-[4.15rem] sm:h-[3.6rem] sm:w-[4.9rem]' : 'h-12 w-12 sm:h-14 sm:w-14'"
            :style="card.assetSrc ? undefined : { color: card.iconPrimary }"
          >
            <img
              v-if="card.assetSrc"
              :src="card.assetSrc"
              alt=""
              draggable="false"
              class="h-full w-full object-contain"
              :class="card.assetClass"
            />
            <svg v-else :viewBox="card.viewBox" class="h-full w-full" aria-hidden="true">
              <path
                v-for="(segment, index) in card.segments"
                :key="`${card.id}-${index}`"
                :d="segment.d"
                :fill="segment.fill ?? 'none'"
                :stroke="segment.stroke ?? 'none'"
                :stroke-width="segment.strokeWidth"
                :stroke-linecap="segment.lineCap"
                :stroke-linejoin="segment.lineJoin"
                :fill-rule="segment.fillRule"
                :clip-rule="segment.clipRule"
                :opacity="segment.opacity"
              />
            </svg>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue'

type IconId = 'github' | 'arch' | 'vue' | 'go' | 'rust' | 'docker'

interface IconSegment {
  d: string
  fill?: string
  stroke?: string
  strokeWidth?: number
  lineCap?: 'round' | 'butt' | 'square'
  lineJoin?: 'round' | 'miter' | 'bevel'
  fillRule?: 'nonzero' | 'evenodd'
  clipRule?: 'nonzero' | 'evenodd'
  opacity?: number
}

interface CardBlueprint {
  id: IconId
  size: number
  viewBox: string
  borderColor: string
  surfaceColor: string
  shadowColor: string
  iconPrimary: string
  iconSecondary: string
  assetClass?: string
  assetSrc?: string
  segments: IconSegment[]
}

interface PhysicsCard extends CardBlueprint {
  x: number
  y: number
  vx: number
  vy: number
  radius: number
}

interface DragState {
  active: boolean
  cardId: IconId | null
  pointerId: number
  offsetX: number
  offsetY: number
  velocityX: number
  velocityY: number
  lastTime: number
}

const sceneRef = ref<HTMLElement | null>(null)

const sceneBounds = reactive({
  width: 0,
  height: 0,
})

const cards = reactive<PhysicsCard[]>([])

const dragState = reactive<DragState>({
  active: false,
  cardId: null,
  pointerId: -1,
  offsetX: 0,
  offsetY: 0,
  velocityX: 0,
  velocityY: 0,
  lastTime: 0,
})

const lowerArcAngles = [160, 132, 104, 76, 48, 20]

const cardBlueprints: CardBlueprint[] = [
  {
    id: 'github',
    size: 80,
    viewBox: '0 0 16 16',
    borderColor: 'rgba(71, 85, 105, 0.28)',
    surfaceColor: 'rgba(241, 245, 249, 0.78)',
    shadowColor: 'rgba(15, 23, 42, 0.14)',
    iconPrimary: '#0f172a',
    iconSecondary: '#475569',
    segments: [
      {
        d: 'M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z',
        fill: 'var(--orbit-icon-primary)',
      },
    ],
  },
  {
    id: 'arch',
    size: 80,
    viewBox: '0 0 24 24',
    borderColor: 'rgba(56, 189, 248, 0.28)',
    surfaceColor: 'rgba(224, 242, 254, 0.82)',
    shadowColor: 'rgba(14, 165, 233, 0.18)',
    iconPrimary: '#0284c7',
    iconSecondary: '#7dd3fc',
    segments: [
      {
        d: 'M12 2.5c1.5 2.14 2.92 4.63 4.23 7.48 1.22 2.63 2.12 5.08 2.71 7.35-.86-.56-1.79-.84-2.8-.84-1.08 0-2.12.41-3.13 1.23 1.13-1.56 1.67-3.02 1.62-4.35-.04-1.22-.55-2.69-1.53-4.42-.36.86-.85 1.83-1.47 2.89-.61 1.01-.92 1.92-.92 2.73 0 1.1.44 2.14 1.32 3.15-.97-.8-2-1.2-3.09-1.2-.97 0-1.89.28-2.76.84.58-2.25 1.47-4.65 2.67-7.22C9.07 7.31 10.48 4.72 12 2.5z',
        fill: 'var(--orbit-icon-primary)',
      },
      {
        d: 'M12 7.35c.66 1.17.99 2.1 1.02 2.8.03.86-.4 1.9-1.28 3.1-.84-1.05-1.26-2.01-1.26-2.89 0-.74.5-1.74 1.52-3.01z',
        fill: 'var(--orbit-icon-secondary)',
      },
    ],
  },
  {
    id: 'vue',
    size: 80,
    viewBox: '0 0 24 24',
    borderColor: 'rgba(16, 185, 129, 0.28)',
    surfaceColor: 'rgba(220, 252, 231, 0.82)',
    shadowColor: 'rgba(5, 150, 105, 0.16)',
    iconPrimary: '#059669',
    iconSecondary: '#6ee7b7',
    segments: [
      {
        d: 'M2 4h5.8L12 11.1 16.2 4H22l-10 16L2 4z',
        fill: 'var(--orbit-icon-primary)',
      },
      {
        d: 'M6.3 4H9l3 5 3-5h2.7L12 13.4 6.3 4z',
        fill: 'var(--orbit-icon-secondary)',
      },
    ],
  },
  {
    id: 'go',
    size: 80,
    viewBox: '0 0 24 24',
    borderColor: 'rgba(56, 189, 248, 0.28)',
    surfaceColor: 'rgba(236, 254, 255, 0.84)',
    shadowColor: 'rgba(6, 182, 212, 0.16)',
    iconPrimary: '#0891b2',
    iconSecondary: '#67e8f9',
    segments: [
      {
        d: 'M2.5 8.7h5.5M1.5 12h4.5M2.5 15.3h5.5',
        stroke: 'var(--orbit-icon-secondary)',
        strokeWidth: 1.8,
        lineCap: 'round',
      },
      {
        d: 'M11.2 8.9a4.1 4.1 0 100 8.2h3v-2.2H11.9a1.9 1.9 0 110-3.8 1.73 1.73 0 011.77 1.45h2.2a3.93 3.93 0 00-4.68-3.66z',
        fill: 'var(--orbit-icon-primary)',
      },
      {
        d: 'M18 8.9a4.1 4.1 0 110 8.2 4.1 4.1 0 010-8.2zm0 2.24a1.86 1.86 0 100 3.72 1.86 1.86 0 000-3.72z',
        fill: 'var(--orbit-icon-primary)',
      },
    ],
  },
  {
    id: 'rust',
    size: 80,
    viewBox: '0 0 24 24',
    borderColor: 'rgba(217, 119, 6, 0.28)',
    surfaceColor: 'rgba(255, 247, 237, 0.84)',
    shadowColor: 'rgba(180, 83, 9, 0.16)',
    iconPrimary: '#b45309',
    iconSecondary: '#fdba74',
    assetClass: 'scale-[0.92] sm:scale-[0.96]',
    assetSrc: '/picture/rustacean-flat-noshadow.svg',
    segments: [],
  },
  {
    id: 'docker',
    size: 80,
    viewBox: '0 0 24 24',
    borderColor: 'rgba(59, 130, 246, 0.28)',
    surfaceColor: 'rgba(239, 246, 255, 0.84)',
    shadowColor: 'rgba(37, 99, 235, 0.16)',
    iconPrimary: '#2563eb',
    iconSecondary: '#93c5fd',
    segments: [
      {
        d: 'M3.9 9.1h2.7v2.7H3.9zm3.2 0h2.7v2.7H7.1zm3.2 0H13v2.7h-2.7zM5.5 12.3h2.7V15H5.5zm3.2 0h2.7V15H8.7zm3.2 0h2.7V15h-2.7z',
        fill: 'var(--orbit-icon-primary)',
      },
      {
        d: 'M20.2 12c.15-.78-.06-1.69-.64-2.31-.75.53-1.2 1.18-1.36 1.94h-1.27V15H6.2c.55 1.42 1.74 2.32 3.36 2.32h3.12c3.42 0 6.08-1.58 7.52-4.55z',
        fill: 'var(--orbit-icon-secondary)',
      },
    ],
  },
]

const scenePadding = 18
const stopVelocity = 0.01
const damping = 0.982
const bounceLoss = 0.94
let frameId = 0
let lastFrameTime = 0
let resizeObserver: ResizeObserver | null = null

function clamp(value: number, min: number, max: number) {
  return Math.min(max, Math.max(min, value))
}

function sceneCenter() {
  return {
    x: sceneBounds.width * 0.5,
    y: sceneBounds.height * 0.45,
  }
}

function arcRadiusX() {
  return Math.min(sceneBounds.width * 0.36, 320)
}

function arcRadiusY() {
  return Math.min(sceneBounds.height * 0.22, 160)
}

function syncSceneBounds() {
  if (!sceneRef.value) {
    return
  }

  const bounds = sceneRef.value.getBoundingClientRect()
  sceneBounds.width = bounds.width
  sceneBounds.height = bounds.height
}

function placeCardsOnLowerArc() {
  if (!sceneBounds.width || !sceneBounds.height) {
    return
  }

  const center = sceneCenter()
  const radiusX = arcRadiusX()
  const radiusY = arcRadiusY()

  cards.splice(
    0,
    cards.length,
    ...cardBlueprints.map((blueprint, index) => {
      const angle = (lowerArcAngles[index] * Math.PI) / 180
      return {
        ...blueprint,
        x: center.x + Math.cos(angle) * radiusX,
        y: center.y + Math.sin(angle) * radiusY,
        vx: 0,
        vy: 0,
        radius: blueprint.size * 0.46,
      }
    }),
  )

  resolveAllCollisions()
}

function updateSceneBounds() {
  syncSceneBounds()

  if (!cards.length) {
    placeCardsOnLowerArc()
    return
  }

  for (const card of cards) {
    card.x = clamp(card.x, card.radius + scenePadding, sceneBounds.width - card.radius - scenePadding)
    card.y = clamp(card.y, card.radius + scenePadding, sceneBounds.height - card.radius - scenePadding)
  }

  resolveAllCollisions()
}

function keepCardInsideBounds(card: PhysicsCard) {
  const minX = card.radius + scenePadding
  const maxX = sceneBounds.width - card.radius - scenePadding
  const minY = card.radius + scenePadding
  const maxY = sceneBounds.height - card.radius - scenePadding

  if (card.x < minX) {
    card.x = minX
    card.vx = Math.abs(card.vx) * bounceLoss
  } else if (card.x > maxX) {
    card.x = maxX
    card.vx = -Math.abs(card.vx) * bounceLoss
  }

  if (card.y < minY) {
    card.y = minY
    card.vy = Math.abs(card.vy) * bounceLoss
  } else if (card.y > maxY) {
    card.y = maxY
    card.vy = -Math.abs(card.vy) * bounceLoss
  }
}

function resolvePairCollision(a: PhysicsCard, b: PhysicsCard) {
  const dx = b.x - a.x
  const dy = b.y - a.y
  const distance = Math.hypot(dx, dy) || 0.001
  const minDistance = a.radius + b.radius + 5

  if (distance >= minDistance) {
    return
  }

  const nx = dx / distance
  const ny = dy / distance
  const overlap = minDistance - distance
  const aDragged = dragState.cardId === a.id
  const bDragged = dragState.cardId === b.id

  if (aDragged && !bDragged) {
    b.x += nx * overlap
    b.y += ny * overlap
  } else if (!aDragged && bDragged) {
    a.x -= nx * overlap
    a.y -= ny * overlap
  } else {
    const half = overlap * 0.5
    a.x -= nx * half
    a.y -= ny * half
    b.x += nx * half
    b.y += ny * half
  }

  const relativeVelocityX = b.vx - a.vx
  const relativeVelocityY = b.vy - a.vy
  const velocityAlongNormal = relativeVelocityX * nx + relativeVelocityY * ny

  if (velocityAlongNormal > 0) {
    return
  }

  const restitution = 0.84
  const impulse = -((1 + restitution) * velocityAlongNormal) / 2

  if (!aDragged) {
    a.vx -= impulse * nx
    a.vy -= impulse * ny
  }

  if (!bDragged) {
    b.vx += impulse * nx
    b.vy += impulse * ny
  }
}

function resolveAllCollisions() {
  for (let pass = 0; pass < 2; pass += 1) {
    for (const card of cards) {
      keepCardInsideBounds(card)
    }

    for (let index = 0; index < cards.length; index += 1) {
      for (let nestedIndex = index + 1; nestedIndex < cards.length; nestedIndex += 1) {
        resolvePairCollision(cards[index], cards[nestedIndex])
      }
    }
  }
}

function startSimulation() {
  if (frameId) {
    return
  }

  lastFrameTime = 0
  frameId = window.requestAnimationFrame(tick)
}

function stopSimulation() {
  if (!frameId) {
    return
  }

  window.cancelAnimationFrame(frameId)
  frameId = 0
  lastFrameTime = 0
}

function tick(now: number) {
  if (!lastFrameTime) {
    lastFrameTime = now
  }

  const delta = Math.min((now - lastFrameTime) / 16.6667, 2)
  lastFrameTime = now

  let moving = dragState.active

  for (const card of cards) {
    if (dragState.cardId === card.id) {
      continue
    }

    card.vx *= Math.pow(damping, delta)
    card.vy *= Math.pow(damping, delta)

    if (Math.abs(card.vx) < stopVelocity) {
      card.vx = 0
    }

    if (Math.abs(card.vy) < stopVelocity) {
      card.vy = 0
    }

    card.x += card.vx * delta
    card.y += card.vy * delta

    if (card.vx !== 0 || card.vy !== 0) {
      moving = true
    }
  }

  resolveAllCollisions()

  if (!moving) {
    stopSimulation()
    return
  }

  frameId = window.requestAnimationFrame(tick)
}

function findCard(cardId: IconId | null) {
  return cards.find((card) => card.id === cardId) ?? null
}

function handleCardPointerDown(cardId: IconId, event: PointerEvent) {
  if (dragState.active || !sceneRef.value) {
    return
  }

  const card = findCard(cardId)
  if (!card) {
    return
  }

  const bounds = sceneRef.value.getBoundingClientRect()
  dragState.active = true
  dragState.cardId = cardId
  dragState.pointerId = event.pointerId
  dragState.offsetX = event.clientX - bounds.left - card.x
  dragState.offsetY = event.clientY - bounds.top - card.y
  dragState.velocityX = 0
  dragState.velocityY = 0
  dragState.lastTime = performance.now()
  card.vx = 0
  card.vy = 0
  startSimulation()
}

function handleWindowPointerMove(event: PointerEvent) {
  if (!dragState.active || dragState.pointerId !== event.pointerId || !sceneRef.value) {
    return
  }

  const card = findCard(dragState.cardId)
  if (!card) {
    return
  }

  const bounds = sceneRef.value.getBoundingClientRect()
  const nextX = event.clientX - bounds.left - dragState.offsetX
  const nextY = event.clientY - bounds.top - dragState.offsetY
  const clampedX = clamp(nextX, card.radius + scenePadding, sceneBounds.width - card.radius - scenePadding)
  const clampedY = clamp(nextY, card.radius + scenePadding, sceneBounds.height - card.radius - scenePadding)
  const now = performance.now()
  const elapsed = Math.max(now - dragState.lastTime, 8)
  const releaseSpeedLimit = 16

  dragState.velocityX = clamp(((clampedX - card.x) / elapsed) * 16.6667, -releaseSpeedLimit, releaseSpeedLimit)
  dragState.velocityY = clamp(((clampedY - card.y) / elapsed) * 16.6667, -releaseSpeedLimit, releaseSpeedLimit)
  dragState.lastTime = now

  card.x = clampedX
  card.y = clampedY
  card.vx = dragState.velocityX
  card.vy = dragState.velocityY

  keepCardInsideBounds(card)
  resolveAllCollisions()
}

function stopDragging(pointerId: number) {
  if (!dragState.active || dragState.pointerId !== pointerId) {
    return
  }

  const card = findCard(dragState.cardId)
  if (card) {
    card.vx = dragState.velocityX
    card.vy = dragState.velocityY
  }

  dragState.active = false
  dragState.cardId = null
  dragState.pointerId = -1
  dragState.velocityX = 0
  dragState.velocityY = 0
  startSimulation()
}

function handleWindowPointerUp(event: PointerEvent) {
  stopDragging(event.pointerId)
}

function handleWindowPointerCancel(event: PointerEvent) {
  stopDragging(event.pointerId)
}

function cardStyle(card: PhysicsCard) {
  return {
    width: `${card.size}px`,
    height: `${card.size}px`,
    transform: `translate3d(${(card.x - card.size / 2).toFixed(1)}px, ${(card.y - card.size / 2).toFixed(1)}px, 0)`,
    zIndex: dragState.cardId === card.id ? '40' : '20',
  }
}

function cardSurfaceStyle(card: PhysicsCard) {
  return {
    borderColor: card.borderColor,
    background: `linear-gradient(150deg, rgba(255,255,255,0.86), ${card.surfaceColor})`,
    boxShadow: `0 20px 42px ${card.shadowColor}`,
    '--orbit-icon-primary': card.iconPrimary,
    '--orbit-icon-secondary': card.iconSecondary,
  }
}

onMounted(() => {
  updateSceneBounds()

  if (sceneRef.value) {
    resizeObserver = new ResizeObserver(() => updateSceneBounds())
    resizeObserver.observe(sceneRef.value)
  }

  window.addEventListener('pointermove', handleWindowPointerMove)
  window.addEventListener('pointerup', handleWindowPointerUp)
  window.addEventListener('pointercancel', handleWindowPointerCancel)
})

onBeforeUnmount(() => {
  stopSimulation()
  window.removeEventListener('pointermove', handleWindowPointerMove)
  window.removeEventListener('pointerup', handleWindowPointerUp)
  window.removeEventListener('pointercancel', handleWindowPointerCancel)
  resizeObserver?.disconnect()
  resizeObserver = null
})
</script>
