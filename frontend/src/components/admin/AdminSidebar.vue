<template>
  <!-- Desktop sidebar -->
  <aside class="hidden w-56 shrink-0 lg:block">
    <div class="sticky top-6 flex h-[calc(100vh-48px)] flex-col border-r border-slate-200/70 pr-5">
      <nav class="flex-1 space-y-6 overflow-y-auto">
        <div v-for="group in navGroups" :key="group.key">
          <p class="px-2 pb-2 text-xs uppercase tracking-[0.18em] text-slate-400">{{ group.section }}</p>

          <div v-if="group.key !== 'settings'" class="space-y-0.5">
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

          <div v-else class="space-y-2">
            <div class="flex items-center gap-1.5">
              <a
                href="/admin/settings"
                class="flex min-w-0 flex-1 items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm transition duration-200"
                :class="activeKey === 'settings'
                  ? 'bg-miku-soft font-medium text-miku'
                  : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'"
                :aria-label="sb.nav.settingsCenter"
                @click="handleSettingsHomeClick"
              >
                <span>{{ sb.nav.settingsCenter }}</span>
              </a>
              <button
                type="button"
                class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-400 transition hover:bg-slate-100 hover:text-miku"
                :aria-label="settingsMenuOpen ? settingsCopy.sections.navCollapseButton : settingsCopy.sections.navExpandButton"
                @click="toggleSettingsMenu"
              >
                <svg
                  viewBox="0 0 20 20"
                  class="h-4 w-4 fill-none stroke-current stroke-[1.8] transition"
                  :class="settingsMenuOpen ? 'rotate-180' : ''"
                  aria-hidden="true"
                >
                  <path d="M5 8l5 5 5-5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>

            <div
              v-show="settingsMenuOpen"
              class="ml-3 space-y-3 rounded-[18px] border border-white/70 bg-white/72 px-3 py-3 shadow-[0_10px_22px_rgba(15,23,42,0.06)]"
            >
              <div
                v-for="groupItem in settingsMenuGroups"
                :key="groupItem.key"
                class="space-y-1.5"
              >
                <p class="px-1 text-[10px] uppercase tracking-[0.2em] text-slate-400">{{ groupItem.title }}</p>
                <a
                  v-for="item in groupItem.items"
                  :key="item.key"
                  :href="item.href"
                  class="block rounded-[14px] px-3 py-2 text-[13px] leading-5 transition duration-200"
                  :class="isSettingsChildActive(item.key)
                    ? 'bg-[linear-gradient(135deg,rgba(57,197,187,0.14),rgba(255,255,255,0.96))] font-medium text-miku shadow-[0_10px_18px_rgba(57,197,187,0.12)]'
                    : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'"
                  :aria-label="item.label"
                  @click="handleSettingsSectionClick(item.key)"
                >
                  {{ item.label }}
                </a>
              </div>
            </div>
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
        <div v-for="group in navGroups" :key="`m-${group.key}`">
          <p class="px-2 pb-2 text-xs uppercase tracking-[0.18em] text-slate-400">{{ group.section }}</p>

          <div v-if="group.key !== 'settings'" class="space-y-0.5">
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

          <div v-else class="space-y-2">
            <div class="flex items-center gap-1.5">
              <a
                href="/admin/settings"
                class="flex min-w-0 flex-1 items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm transition duration-200"
                :class="activeKey === 'settings'
                  ? 'bg-miku-soft font-medium text-miku'
                  : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'"
                :aria-label="sb.nav.settingsCenter"
                @click="handleSettingsHomeClick"
              >
                <span>{{ sb.nav.settingsCenter }}</span>
              </a>
              <button
                type="button"
                class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-400 transition hover:bg-slate-100 hover:text-miku"
                :aria-label="settingsMenuOpen ? settingsCopy.sections.navCollapseButton : settingsCopy.sections.navExpandButton"
                @click="toggleSettingsMenu"
              >
                <svg
                  viewBox="0 0 20 20"
                  class="h-4 w-4 fill-none stroke-current stroke-[1.8] transition"
                  :class="settingsMenuOpen ? 'rotate-180' : ''"
                  aria-hidden="true"
                >
                  <path d="M5 8l5 5 5-5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>

            <div
              v-show="settingsMenuOpen"
              class="ml-3 space-y-3 rounded-[18px] border border-white/70 bg-white/74 px-3 py-3 shadow-[0_10px_22px_rgba(15,23,42,0.06)]"
            >
              <div
                v-for="groupItem in settingsMenuGroups"
                :key="`mobile-group-${groupItem.key}`"
                class="space-y-1.5"
              >
                <p class="px-1 text-[10px] uppercase tracking-[0.2em] text-slate-400">{{ groupItem.title }}</p>
                <a
                  v-for="item in groupItem.items"
                  :key="`mobile-section-${item.key}`"
                  :href="item.href"
                  class="block rounded-[14px] px-3 py-2 text-[13px] leading-5 transition duration-200"
                  :class="isSettingsChildActive(item.key)
                    ? 'bg-[linear-gradient(135deg,rgba(57,197,187,0.14),rgba(255,255,255,0.96))] font-medium text-miku shadow-[0_10px_18px_rgba(57,197,187,0.12)]'
                    : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'"
                  :aria-label="item.label"
                  @click="handleSettingsSectionClick(item.key)"
                >
                  {{ item.label }}
                </a>
              </div>
            </div>
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
import { computed, onMounted, ref } from 'vue'

import { adminCopy } from '../../content/copy'
import { logout } from '../../stores/auth'
import {
  adminSettingsSection,
  hydrateAdminSettingsSectionFromWindow,
  setAdminSettingsSection,
  type AdminSettingsSectionKey,
} from '../../stores/adminSettings'
import { sidebarOpen, setSidebarOpen } from '../../stores/ui'

const sb = adminCopy.sidebar
const settingsCopy = adminCopy.settingsCenter

interface Props {
  activeKey?: string
}

const props = withDefaults(defineProps<Props>(), {
  activeKey: 'dashboard',
})

const navGroups = [
  {
    key: 'overview',
    section: sb.sections.overview,
    items: [
      { key: 'dashboard', label: sb.nav.dashboard, href: '/admin' },
    ],
  },
  {
    key: 'content',
    section: sb.sections.content,
    items: [
      { key: 'posts', label: sb.nav.posts, href: '/admin/posts' },
      { key: 'moments', label: sb.nav.moments, href: '/admin/moments' },
      { key: 'comments', label: sb.nav.comments, href: '/admin/comments' },
    ],
  },
  {
    key: 'settings',
    section: sb.sections.settings,
    items: [
      { key: 'settings', label: sb.nav.settingsCenter, href: '/admin/settings' },
    ],
  },
  {
    key: 'system',
    section: sb.sections.system,
    items: [
      { key: 'friends', label: sb.nav.friends, href: '/admin/friends' },
      { key: 'backup', label: sb.nav.backup, href: '/admin/backup' },
    ],
  },
]

const settingsMenuGroups = computed(() => ([
  {
    key: 'frontstage',
    title: settingsCopy.sections.groups.frontstage.title,
    items: [
      { key: 'site-profile' as AdminSettingsSectionKey, label: settingsCopy.sections.siteProfile.title, href: '/admin/settings?section=site-profile' },
      { key: 'home-hero' as AdminSettingsSectionKey, label: settingsCopy.sections.homeHero.title, href: '/admin/settings?section=home-hero' },
      { key: 'home-assets' as AdminSettingsSectionKey, label: settingsCopy.sections.homeAssets.title, href: '/admin/settings?section=home-assets' },
    ],
  },
  {
    key: 'identity',
    title: settingsCopy.sections.groups.identity.title,
    items: [
      { key: 'author-profile' as AdminSettingsSectionKey, label: settingsCopy.sections.authorProfile.title, href: '/admin/settings?section=author-profile' },
      { key: 'site-integrations' as AdminSettingsSectionKey, label: settingsCopy.sections.siteIntegrations.title, href: '/admin/settings?section=site-integrations' },
    ],
  },
  {
    key: 'operations',
    title: settingsCopy.sections.groups.operations.title,
    items: [
      { key: 'site-footer' as AdminSettingsSectionKey, label: settingsCopy.sections.siteFooter.title, href: '/admin/settings?section=site-footer' },
      { key: 'admin-profile' as AdminSettingsSectionKey, label: settingsCopy.sections.adminProfile.title, href: '/admin/settings?section=admin-profile' },
    ],
  },
]))

const activeKey = computed(() => props.activeKey)
const settingsMenuOpen = ref(props.activeKey === 'settings')
const mobileOpen = useStore(sidebarOpen)
const currentSettingsSection = useStore(adminSettingsSection)

function toggleSettingsMenu() {
  settingsMenuOpen.value = !settingsMenuOpen.value
}

function isSettingsChildActive(section: AdminSettingsSectionKey) {
  return activeKey.value === 'settings' && currentSettingsSection.value === section
}

function handleSettingsHomeClick() {
  setAdminSettingsSection('site-profile')
  if (mobileOpen.value) {
    setSidebarOpen(false)
  }
}

function handleSettingsSectionClick(section: AdminSettingsSectionKey) {
  setAdminSettingsSection(section)
  if (mobileOpen.value) {
    setSidebarOpen(false)
  }
}

async function handleLogout() {
  await logout()
  window.location.replace('/login')
}

onMounted(() => {
  if (activeKey.value === 'settings') {
    hydrateAdminSettingsSectionFromWindow()
    settingsMenuOpen.value = true
  }
})
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
