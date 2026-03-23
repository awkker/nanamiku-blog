<template>
  <div class="border border-slate-200 bg-white text-slate-800 shadow-[0_8px_24px_rgba(15,23,42,0.06)]" :style="cardStyle">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  width?: string
  maxWidth?: string
  padding?: string
  borderRadius?: number
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  maxWidth: '100%',
  padding: '20px',
  borderRadius: 16,
})

const cardStyle = computed(() => ({
  width: props.width,
  maxWidth: props.maxWidth,
  padding: props.padding,
  borderRadius: `${props.borderRadius}px`,
}))
</script>
