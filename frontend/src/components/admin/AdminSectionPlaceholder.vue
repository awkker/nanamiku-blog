<template>
  <section class="space-y-5">
    <AdminPlainCard padding="24px">
      <h1 class="text-2xl font-semibold text-slate-900">{{ title }}</h1>
      <p class="mt-2 text-sm text-slate-700">{{ description }}</p>
    </AdminPlainCard>

    <div class="grid gap-4 lg:grid-cols-3">
      <AdminPlainCard
        v-for="item in cards"
        :key="item.title"
        padding="20px"
      >
        <h2 class="text-base font-semibold text-slate-900">{{ item.title }}</h2>
        <p class="mt-2 text-sm text-slate-700">{{ item.desc }}</p>
      </AdminPlainCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import AdminPlainCard from '../ui/AdminPlainCard.vue'
import { adminCopy } from '../../content/copy'

interface Props {
  title: string
  description: string
  module: 'posts' | 'comments' | 'friends'
}

const props = defineProps<Props>()

const cards = computed(() => adminCopy.sectionPlaceholder[props.module])
</script>
