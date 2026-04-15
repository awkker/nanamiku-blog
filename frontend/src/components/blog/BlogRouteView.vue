<template>
  <div>
    <div v-if="mode === 'loading'" class="mx-auto w-full max-w-[1320px] px-4 pb-20 pt-8 sm:px-6 lg:px-8">
      <div class="space-y-6">
        <div class="animate-pulse rounded-[32px] border border-white/65 bg-white/70 p-6 shadow-[0_18px_42px_rgba(15,23,42,0.1)]">
          <div class="h-6 w-32 rounded-full bg-slate-200/70" />
          <div class="mt-4 h-12 w-2/3 rounded-2xl bg-slate-200/70" />
          <div class="mt-3 h-4 w-1/2 rounded-full bg-slate-100/80" />
          <div class="mt-8 grid gap-5 lg:grid-cols-[minmax(0,1fr)_320px]">
            <div class="space-y-4">
              <div class="h-72 rounded-[28px] bg-slate-200/60" />
              <div class="h-72 rounded-[28px] bg-slate-200/55" />
            </div>
            <div class="space-y-4">
              <div class="h-44 rounded-[28px] bg-slate-200/55" />
              <div class="h-56 rounded-[28px] bg-slate-200/50" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <template v-else-if="mode === 'archive'">
      <div class="blog-index-background pointer-events-none fixed inset-0 -z-10 bg-[radial-gradient(circle_at_top,rgba(57,197,187,0.22),transparent_46%),radial-gradient(circle_at_85%_18%,rgba(192,132,252,0.16),transparent_36%),linear-gradient(180deg,#f8fafc_0%,#eef2ff_100%)]" />

      <BlogHeroSection />

      <main id="content-start" class="mx-auto mt-6 grid w-full max-w-[1320px] grid-cols-1 gap-8 px-4 pb-20 lg:grid-cols-[minmax(0,1fr)_320px]">
        <section id="archive" class="scroll-mt-24">
          <div id="latest-posts" class="scroll-mt-24" />
          <BlogFeed />
        </section>

        <aside class="space-y-6">
          <BlogAuthorSidebar />

          <div class="glass-layer rounded-2xl p-4">
            <h3 class="text-sm font-semibold text-slate-900">{{ copy.playlistTitle }}</h3>
            <ul class="mt-3 space-y-2 text-xs text-slate-600">
              <li
                v-for="song in copy.playlist"
                :key="song"
                class="rounded-lg border border-white/60 bg-white/66 px-3 py-2"
              >
                {{ song }}
              </li>
            </ul>
          </div>

          <div class="relative z-0 mb-8 overflow-hidden rounded-2xl border border-white/35 bg-white/20 p-4 shadow-md backdrop-blur-xl">
            <SiteTrend />
          </div>

          <div class="relative z-0">
            <LiquidGlassCard padding="20px">
              <LatestMoments />
            </LiquidGlassCard>
          </div>
        </aside>
      </main>
    </template>

    <main v-else class="relative min-h-screen px-4 pb-20 pt-6 sm:px-6 lg:px-8">
      <div class="blog-post-view-background pointer-events-none fixed inset-0 -z-10 bg-[radial-gradient(circle_at_top,rgba(57,197,187,0.12),transparent_50%),radial-gradient(circle_at_85%_18%,rgba(192,132,252,0.08),transparent_40%),linear-gradient(180deg,#f8fafc_0%,#eef2ff_100%)]" />
      <BlogPostView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { siteCopy } from '../../content/copy'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'
import LatestMoments from '../moments/LatestMoments.vue'
import BlogAuthorSidebar from './BlogAuthorSidebar.vue'
import BlogFeed from './BlogFeed.vue'
import BlogHeroSection from './BlogHeroSection.vue'
import BlogPostView from './BlogPostView.vue'
import SiteTrend from './SiteTrend.vue'

type BlogRouteMode = 'loading' | 'archive' | 'detail'

const copy = siteCopy.blogIndex
const mode = ref<BlogRouteMode>('loading')

function resolveMode() {
  if (typeof window === 'undefined') {
    mode.value = 'loading'
    return
  }

  const segments = window.location.pathname.split('/').filter(Boolean)
  mode.value = segments[0] === 'blog' && segments.length > 1 ? 'detail' : 'archive'
}

onMounted(() => {
  resolveMode()
})
</script>
