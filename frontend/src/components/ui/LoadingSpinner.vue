<template>
  <svg
    class="rocket-spinner text-current"
    :class="sizeClass"
    viewBox="0 0 24 24"
    fill="none"
    aria-hidden="true"
  >
    <circle class="rocket-spinner__halo" cx="12" cy="12" r="9" />
    <g class="rocket-spinner__ship">
      <path
        d="M14.7 4.2c1.8.9 3.2 2.3 4.1 4.1L13.2 14l-2.8.4.4-2.8 3.9-7.4z"
        fill="currentColor"
        fill-opacity="0.9"
      />
      <path d="M10.6 11.6L8 14.2a1.8 1.8 0 000 2.6l.8.8a1.8 1.8 0 002.6 0l2.6-2.6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
      <circle cx="14.4" cy="8.7" r="1.1" fill="white" fill-opacity="0.9" />
    </g>
    <path class="rocket-spinner__flame" d="M8.3 17.5l-2.1 2.6 3.3-1.3z" fill="#f59e0b" />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
})

const sizeClass = computed(() => {
  if (props.size === 'sm') {
    return 'h-4 w-4'
  }

  if (props.size === 'lg') {
    return 'h-8 w-8'
  }

  return 'h-6 w-6'
})
</script>

<style scoped>
.rocket-spinner__halo {
  stroke: currentColor;
  stroke-opacity: 0.2;
  stroke-width: 1.4;
  stroke-dasharray: 11 38;
  transform-origin: 50% 50%;
  animation: rocket-halo 1.2s linear infinite;
}

.rocket-spinner__ship {
  transform-origin: 50% 50%;
  animation: rocket-bob 0.95s ease-in-out infinite;
}

.rocket-spinner__flame {
  transform-origin: 8px 18px;
  animation: rocket-flame 0.6s ease-in-out infinite;
}

@keyframes rocket-halo {
  to {
    transform: rotate(360deg);
  }
}

@keyframes rocket-bob {
  0%, 100% {
    transform: translateY(0) rotate(-2deg);
  }
  50% {
    transform: translateY(-0.8px) rotate(1deg);
  }
}

@keyframes rocket-flame {
  0%, 100% {
    transform: scale(1, 0.92);
    opacity: 0.72;
  }
  50% {
    transform: scale(1.06, 1.22);
    opacity: 1;
  }
}
</style>
