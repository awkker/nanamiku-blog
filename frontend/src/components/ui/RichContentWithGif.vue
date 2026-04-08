<template>
  <div :class="textClass">
    <template v-for="segment in segments" :key="segment.key">
      <span v-if="segment.type === 'text'">{{ segment.value }}</span>
      <img
        v-else
        :src="segment.src"
        :alt="`${altPrefix}${segment.label}`"
        :title="segment.label"
        loading="lazy"
        decoding="async"
        :class="imageClass"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import { parseRichContentWithGif } from '../../lib/gifEmotes'
import { siteCopy } from '../../content/copy'

interface Props {
  content: string
  textClass?: string
  imageClass?: string
  altPrefix?: string
}

const props = withDefaults(defineProps<Props>(), {
  textClass: 'whitespace-pre-wrap break-words',
  imageClass: 'my-1 mr-1 inline-block h-16 w-16 rounded-lg object-cover align-middle ring-1 ring-white/80',
  altPrefix: siteCopy.components.richContentWithGif.altPrefix,
})

const segments = computed(() => parseRichContentWithGif(props.content))
</script>
