<template>
  <span class="jump-letter__container">
    <span
      v-for="(char, i) in chars"
      :key="i"
      class="jump-letter__alphabet"
      @mouseenter="onEnter"
      @animationend="onAnimEnd"
    >{{ char }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  text: string
}

const props = defineProps<Props>()

// 把整句标题拆成单个字符，后面模板才能逐字绑定 hover 动画。
const chars = computed(() => [...props.text])

function onEnter(e: Event) {
  const el = e.target as HTMLElement
  // 动画播放期间不重复加类，避免频繁 hover 时动画抖动或重叠。
  if (el.classList.contains('is-active')) return
  el.classList.add('is-active')
}

function onAnimEnd(e: AnimationEvent) {
  // 动画结束后移除激活类，下次 hover 才能重新触发。
  ;(e.target as HTMLElement).classList.remove('is-active')
}
</script>

<style scoped>
.jump-letter__container {
  display: inline-flex;
}

.jump-letter__alphabet {
  display: inline-block;
  transform-origin: center bottom;
  cursor: default;
}

.jump-letter__alphabet.is-active {
  animation: jump-letter 1.5s 0s cubic-bezier(0.165, 0.85, 0.45, 1);
}

@keyframes jump-letter {
  0%, 100% {
    transform: translateY(0) scale(1, 1);
    text-shadow: none;
  }
  25%, 75% {
    transform: translateY(0) scale(1.2, 0.8);
    text-shadow: 0 -3px 2px rgba(192, 132, 252, 0.5);
  }
  50% {
    transform: translateY(-24px) scale(0.8, 1.2);
    text-shadow: 0 -8px 4px rgba(192, 132, 252, 0.4);
    color: #c084fc;
  }
}
</style>
