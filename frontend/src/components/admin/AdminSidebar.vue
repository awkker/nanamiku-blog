<template>
  <!-- Desktop sidebar -->
  <aside class="hidden w-56 shrink-0 lg:block">
    <div class="sticky top-6 flex h-[calc(100vh-48px)] flex-col border-r border-slate-200/70 pr-5">
      <nav class="flex-1 space-y-6 overflow-y-auto">
        <div v-for="group in navGroups" :key="group.section">
          <p class="px-2 pb-2 text-xs uppercase tracking-[0.18em] text-slate-400">{{ group.section }}</p>
          <div class="space-y-0.5">
            <a
              v-for="item in group.items"
              :key="item.key"
              :href="item.href"
              class="flex items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm transition duration-200"
              :class="item.key === activeKey
                ? 'bg-miku-soft font-medium text-miku'
                : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'"
              :aria-label="item.label"
            >
              <span>{{ item.label }}</span>
            </a>
          </div>
        </div>
      </nav>

      <div class="border-t border-slate-200/70 pt-4">
        <a
          href="/"
          class="flex items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm text-slate-600 transition duration-200 hover:bg-slate-100 hover:text-slate-900"
          :aria-label="sb.footer.homeAria"
        >
          {{ sb.footer.home }}
        </a>
        <button
          type="button"
          class="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm text-slate-600 transition duration-200 hover:bg-slate-100 hover:text-slate-900"
          :aria-label="sb.footer.logoutAria"
          @click="handleLogout"
        >
          {{ sb.footer.logout }}
        </button>
      </div>
    </div>
  </aside>

  <!-- Mobile overlay -->
  <Transition name="drawer-fade">
    <div
      v-if="mobileOpen"
      class="fixed inset-0 z-40 bg-slate-950/45 backdrop-blur-sm lg:hidden"
      @click="setSidebarOpen(false)"
    />
  </Transition>

  <!-- Mobile drawer -->
  <Transition name="drawer-slide">
    <aside
      v-if="mobileOpen"
      class="fixed left-0 top-0 z-50 h-full w-64 bg-slate-50 px-4 py-5 shadow-xl lg:hidden"
    >
      <div class="mb-5 flex items-center justify-between">
        <p class="text-sm font-semibold text-slate-900">{{ sb.mobile.brandTitle }}</p>
        <button
          type="button"
          class="rounded-lg p-1.5 text-slate-500 transition hover:bg-slate-100 hover:text-slate-900"
          :aria-label="sb.mobile.closeAria"
          @click="setSidebarOpen(false)"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
            <path d="M6 6l12 12M18 6L6 18" />
          </svg>
        </button>
      </div>

      <nav class="space-y-5">
        <div v-for="group in navGroups" :key="`m-${group.section}`">
          <p class="px-2 pb-2 text-xs uppercase tracking-[0.18em] text-slate-400">{{ group.section }}</p>
          <div class="space-y-0.5">
            <a
              v-for="item in group.items"
              :key="`mobile-${item.key}`"
              :href="item.href"
              class="flex items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm transition duration-200"
              :class="item.key === activeKey
                ? 'bg-miku-soft font-medium text-miku'
                : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'"
              :aria-label="item.label"
              @click="setSidebarOpen(false)"
            >
              <span>{{ item.label }}</span>
            </a>
          </div>
        </div>
      </nav>

      <div class="mt-6 border-t border-slate-200/70 pt-4">
        <a
          href="/"
          class="flex items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
          :aria-label="sb.footer.homeAria"
          @click="setSidebarOpen(false)"
        >
          {{ sb.footer.home }}
        </a>
        <button
          type="button"
          class="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm text-slate-600 transition hover:bg-slate-100"
          :aria-label="sb.footer.logoutAria"
          @click="handleLogout"
        >
          {{ sb.footer.logout }}
        </button>
      </div>
    </aside>
  </Transition>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'

import { logout } from '../../stores/auth'
import { sidebarOpen, setSidebarOpen } from '../../stores/ui'
import { adminCopy } from '../../content/copy'

const sb = adminCopy.sidebar

interface Props {
  activeKey?: string
}

withDefaults(defineProps<Props>(), {
  activeKey: 'dashboard',
})

const navGroups = [
  {
    section: sb.sections.overview,
    items: [
      { key: 'dashboard', label: sb.nav.dashboard, href: '/admin' },
    ],
  },
  {
    section: sb.sections.content,
    items: [
      { key: 'posts', label: sb.nav.posts, href: '/admin/posts' },
      { key: 'moments', label: sb.nav.moments, href: '/admin/moments' },
      { key: 'comments', label: sb.nav.comments, href: '/admin/comments' },
    ],
  },
  {
    section: sb.sections.settings,
    items: [
      { key: 'settings', label: sb.nav.settingsCenter, href: '/admin/settings' },
    ],
  },
  {
    section: sb.sections.system,
    items: [
      { key: 'footer', label: sb.nav.footer, href: '/admin/footer' },
      { key: 'profile', label: sb.nav.profile, href: '/admin/profile' },
      { key: 'friends', label: sb.nav.friends, href: '/admin/friends' },
      { key: 'backup', label: sb.nav.backup, href: '/admin/backup' },
    ],
  },
]

const mobileOpen = useStore(sidebarOpen)

async function handleLogout() {
  await logout()
  window.location.replace('/login')
}
</script>

<style scoped>
.drawer-fade-enter-active,
.drawer-fade-leave-active {
  transition: opacity 0.3s ease;
}

.drawer-fade-enter-from,
.drawer-fade-leave-to {
  opacity: 0;
}

.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: transform 0.35s ease, opacity 0.35s ease;
}

.drawer-slide-enter-from,
.drawer-slide-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}
</style>
