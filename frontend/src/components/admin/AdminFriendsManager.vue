<template>
  <section class="space-y-5">
    <AdminPlainCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">{{ copy.page.title }}</h1>
          <p class="mt-1 text-sm text-slate-600">{{ copy.page.subtitle }}</p>
        </div>
        <MikuButton variant="solid" :aria-label="copy.page.addButtonAria" @click="toggleCreateForm">
          {{ copy.page.addButton }}
        </MikuButton>
      </div>
    </AdminPlainCard>

    <AdminPlainCard v-if="showCreateForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">{{ copy.forms.createTitle }}</h2>
      <form class="space-y-3" @submit.prevent="createFriend">
        <div class="grid gap-3 md:grid-cols-2">
          <input v-model="newFriend.name" type="text" :placeholder="copy.forms.namePlaceholder" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
          <input v-model="newFriend.url" type="text" :placeholder="copy.forms.urlPlaceholder" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        </div>
        <input v-model="newFriend.avatar_url" type="text" :placeholder="copy.forms.avatarPlaceholder" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="newFriend.description" type="text" :placeholder="copy.forms.descriptionPlaceholder" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <MikuButton type="submit" variant="solid" :disabled="creatingFriend">{{ creatingFriend ? copy.forms.creatingButton : copy.forms.createButton }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeCreateForm">{{ copy.forms.cancelButton }}</button>
        </div>
      </form>
    </AdminPlainCard>

    <AdminPlainCard v-if="showEditForm" padding="24px">
      <h2 class="mb-4 text-lg font-semibold text-slate-900">{{ copy.forms.editTitle }}</h2>
      <form class="space-y-3" @submit.prevent="updateFriend">
        <div class="grid gap-3 md:grid-cols-2">
          <input v-model="editFriend.name" type="text" :placeholder="copy.forms.namePlaceholder" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
          <input v-model="editFriend.url" type="text" :placeholder="copy.forms.urlPlaceholder" class="rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        </div>
        <input v-model="editFriend.avatar_url" type="text" :placeholder="copy.forms.avatarPlaceholder" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model="editFriend.description" type="text" :placeholder="copy.forms.descriptionPlaceholder" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <input v-model.number="editFriend.sort_order" type="number" min="0" :placeholder="copy.forms.sortPlaceholder" class="w-full rounded-xl border border-slate-200/80 bg-white/60 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-1 focus:ring-miku/30" />
        <div class="flex items-center gap-3">
          <MikuButton type="submit" variant="solid" :disabled="updatingFriend">{{ updatingFriend ? copy.forms.updatingButton : copy.forms.updateButton }}</MikuButton>
          <button type="button" class="text-sm text-slate-500 hover:text-slate-700" @click="closeEditForm">{{ copy.forms.cancelButton }}</button>
        </div>
      </form>
    </AdminPlainCard>

    <div class="grid gap-4 sm:grid-cols-3">
      <AdminPlainCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">{{ copy.stats.approved }}</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-slate-900">{{ approvedCount }}</p>
      </AdminPlainCard>
      <AdminPlainCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">{{ copy.stats.pending }}</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-amber-600">{{ pendingCount }}</p>
      </AdminPlainCard>
      <AdminPlainCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">{{ copy.stats.down }}</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-red-600">{{ downCount }}</p>
      </AdminPlainCard>
    </div>

    <AdminPlainCard padding="0px">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-slate-200/60 px-5 py-4">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">{{ copy.applications.title }}</h2>
          <p class="mt-1 text-sm text-slate-600">{{ copy.applications.subtitle }}</p>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.applications.columns.site }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.applications.columns.contact }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.applications.columns.description }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.applications.columns.status }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.applications.columns.createdAt }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">{{ copy.applications.columns.actions }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="application in applications"
              :key="application.id"
              class="border-b border-slate-100/60 align-top transition hover:bg-white/40"
            >
              <td class="px-5 py-3.5">
                <div class="min-w-[220px]">
                  <p class="font-medium text-slate-900">{{ application.siteName }}</p>
                  <a
                    :href="application.siteUrl"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="mt-1 inline-flex text-xs text-miku underline decoration-miku/30 transition hover:decoration-miku"
                  >
                    {{ application.siteUrl }}
                  </a>
                </div>
              </td>
              <td class="px-5 py-3.5 text-slate-600">
                <div class="min-w-[180px] space-y-1">
                  <p>{{ application.contactEmail }}</p>
                  <p v-if="application.contactNote" class="text-xs leading-relaxed text-slate-400">{{ application.contactNote }}</p>
                </div>
              </td>
              <td class="px-5 py-3.5 text-slate-600">
                <p class="max-w-sm leading-relaxed">{{ application.description }}</p>
                <p v-if="application.reviewedAt" class="mt-2 text-xs text-slate-400">
                  {{ copy.applications.reviewedPrefix }}{{ application.reviewedAt }}
                </p>
                <p v-if="application.reviewNote" class="mt-1 text-xs text-slate-400">
                  {{ copy.applications.notePrefix }}{{ application.reviewNote }}
                </p>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium" :class="statusClass(application.status)">
                  {{ statusLabel(application.status) }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-slate-600">{{ application.createdAt }}</td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    v-if="application.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-emerald-200/80 bg-white/50 px-2.5 py-1 text-xs text-emerald-600 transition hover:bg-emerald-50 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="processingApplicationId === application.id"
                    @click="approveApplication(application)"
                  >
                    {{ processingApplicationId === application.id ? copy.applications.processingButton : copy.applications.approveButton }}
                  </button>
                  <button
                    v-if="application.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="processingApplicationId === application.id"
                    @click="rejectApplication(application)"
                  >
                    {{ copy.applications.rejectButton }}
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="applications.length === 0 && !loadingApplications">
              <td colspan="6" class="px-5 py-8 text-center text-sm text-slate-500">{{ copy.applications.empty }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </AdminPlainCard>

    <AdminPlainCard padding="0px">
      <div class="flex items-center justify-between gap-3 border-b border-slate-200/60 px-5 py-4">
        <h2 class="text-lg font-semibold text-slate-900">{{ copy.links.tableTitle }}</h2>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200/60">
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.links.columns.name }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.links.columns.url }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.links.columns.status }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.links.columns.health }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700">{{ copy.links.columns.createdAt }}</th>
              <th class="px-5 py-3.5 font-semibold text-slate-700 text-center">{{ copy.links.columns.actions }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="link in friends"
              :key="link.id"
              class="border-b border-slate-100/60 transition hover:bg-white/40"
            >
              <td class="px-5 py-3.5 font-medium text-slate-900">{{ link.name }}</td>
              <td class="px-5 py-3.5">
                <a
                  :href="link.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-miku underline decoration-miku/30 transition hover:decoration-miku"
                >
                  {{ link.url }}
                </a>
              </td>
              <td class="px-5 py-3.5">
                <span class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium" :class="statusClass(link.status)">
                  {{ statusLabel(link.status) }}
                </span>
              </td>
              <td class="px-5 py-3.5">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="link.health === 'ok' ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-600'"
                >
                  {{ link.health === 'ok' ? copy.links.healthOk : copy.links.healthDown }}
                </span>
              </td>
              <td class="px-5 py-3.5 text-slate-600">{{ link.createdAt }}</td>
              <td class="px-5 py-3.5 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button
                    v-if="link.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-emerald-200/80 bg-white/50 px-2.5 py-1 text-xs text-emerald-600 transition hover:bg-emerald-50 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="processingFriendId === link.id"
                    @click="approveFriend(link)"
                  >
                    {{ copy.links.approveButton }}
                  </button>
                  <button
                    v-if="link.status === 'pending'"
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-60"
                    :disabled="processingFriendId === link.id"
                    @click="rejectFriend(link)"
                  >
                    {{ copy.links.rejectButton }}
                  </button>
                  <button
                    v-if="link.status !== 'pending'"
                    type="button"
                    class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-700 transition hover:border-miku/40 hover:text-miku"
                    @click="startEditFriend(link)"
                  >
                    {{ copy.links.editButton }}
                  </button>
                  <button
                    v-if="link.status !== 'pending'"
                    type="button"
                    class="rounded-xl border border-red-200/80 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                    @click="deleteFriend(link.id)"
                  >
                    {{ copy.links.deleteButton }}
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="friends.length === 0 && !loadingFriends">
              <td colspan="6" class="px-5 py-8 text-center text-sm text-slate-500">{{ copy.links.empty }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </AdminPlainCard>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { adminCopy } from '../../content/copy'
import { api, type PagedData } from '../../lib/api'
import { showToast } from '../../stores/ui'
import AdminPlainCard from '../ui/AdminPlainCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface ApiFriendLink {
  id: string
  name: string
  description: string
  url: string
  domain: string
  avatar_url: string
  status: string
  health_status: string
  sort_order: number
  created_at: string
}

interface ApiFriendApplication {
  id: string
  site_name: string
  site_url: string
  avatar_url: string
  description: string
  contact_email: string
  contact_note: string
  status: string
  created_at: string
  reviewed_at?: string
  review_note: string
}

interface FriendLink {
  id: string
  name: string
  description: string
  url: string
  domain: string
  avatar_url: string
  sortOrder: number
  status: 'approved' | 'pending' | 'rejected'
  health: 'ok' | 'down'
  createdAt: string
}

interface FriendApplication {
  id: string
  siteName: string
  siteUrl: string
  avatarUrl: string
  description: string
  contactEmail: string
  contactNote: string
  status: 'approved' | 'pending' | 'rejected'
  createdAt: string
  reviewedAt: string
  reviewNote: string
}

interface FriendForm {
  name: string
  url: string
  avatar_url: string
  description: string
  sort_order: number
}

const copy = adminCopy.friendsManager

function mapStatus(s: string): 'approved' | 'pending' | 'rejected' {
  if (s === 'approved') return 'approved'
  if (s === 'rejected') return 'rejected'
  return 'pending'
}

function formatDate(iso: string): string {
  try {
    const date = new Date(iso)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    })
  } catch {
    return iso
  }
}

function mapFriend(item: ApiFriendLink): FriendLink {
  return {
    id: item.id,
    name: item.name,
    description: item.description || '',
    url: item.url,
    domain: item.domain || '',
    avatar_url: item.avatar_url || '',
    sortOrder: Number(item.sort_order) || 0,
    status: mapStatus(item.status),
    health: item.health_status === 'ok' ? 'ok' : 'down',
    createdAt: formatDate(item.created_at),
  }
}

function mapApplication(item: ApiFriendApplication): FriendApplication {
  return {
    id: item.id,
    siteName: item.site_name,
    siteUrl: item.site_url,
    avatarUrl: item.avatar_url || '',
    description: item.description || '',
    contactEmail: item.contact_email || '',
    contactNote: item.contact_note || '',
    status: mapStatus(item.status),
    createdAt: formatDate(item.created_at),
    reviewedAt: item.reviewed_at ? formatDate(item.reviewed_at) : '',
    reviewNote: item.review_note || '',
  }
}

function createEmptyFriendForm(): FriendForm {
  return {
    name: '',
    url: '',
    avatar_url: '',
    description: '',
    sort_order: 0,
  }
}

function deriveDomainFromURL(rawURL: string): string {
  const input = rawURL.trim()
  if (!input) return ''
  try {
    const withScheme = /^[a-zA-Z][a-zA-Z\d+\-.]*:\/\//.test(input) ? input : `https://${input}`
    return new URL(withScheme).hostname.replace(/^www\./, '')
  } catch {
    return ''
  }
}

function toFriendPayload(form: FriendForm) {
  const url = form.url.trim()
  return {
    name: form.name.trim(),
    url,
    domain: deriveDomainFromURL(url),
    avatar_url: form.avatar_url.trim(),
    description: form.description.trim(),
    sort_order: Number(form.sort_order) || 0,
  }
}

const friends = ref<FriendLink[]>([])
const applications = ref<FriendApplication[]>([])
const loadingFriends = ref(false)
const loadingApplications = ref(false)
const showCreateForm = ref(false)
const creatingFriend = ref(false)
const showEditForm = ref(false)
const updatingFriend = ref(false)
const editingFriendID = ref<string | null>(null)
const processingApplicationId = ref<string | null>(null)
const processingFriendId = ref<string | null>(null)

const newFriend = ref<FriendForm>(createEmptyFriendForm())
const editFriend = ref<FriendForm>(createEmptyFriendForm())

function toggleCreateForm() {
  showCreateForm.value = !showCreateForm.value
  if (showCreateForm.value) {
    closeEditForm()
  }
}

function closeCreateForm() {
  showCreateForm.value = false
}

function closeEditForm() {
  showEditForm.value = false
  editingFriendID.value = null
  editFriend.value = createEmptyFriendForm()
}

async function loadFriends() {
  loadingFriends.value = true
  try {
    const data = await api.get<PagedData<ApiFriendLink>>('/admin/friends?size=100')
    friends.value = (data.items || []).map(mapFriend)
  } catch (err) {
    console.error('[AdminFriends] loadFriends failed:', err)
    showToast(copy.toasts.loadFriendsFailed, 'error')
    friends.value = []
  } finally {
    loadingFriends.value = false
  }
}

async function loadApplications() {
  loadingApplications.value = true
  try {
    const data = await api.get<PagedData<ApiFriendApplication>>('/admin/friends/applications?size=100')
    applications.value = (data.items || []).map(mapApplication)
  } catch (err) {
    console.error('[AdminFriends] loadApplications failed:', err)
    showToast(copy.toasts.loadApplicationsFailed, 'error')
    applications.value = []
  } finally {
    loadingApplications.value = false
  }
}

async function loadAll() {
  await Promise.all([loadFriends(), loadApplications()])
}

async function createFriend() {
  if (!newFriend.value.name.trim() || !newFriend.value.url.trim()) return
  creatingFriend.value = true
  try {
    await api.post('/admin/friends', toFriendPayload(newFriend.value))
    closeCreateForm()
    newFriend.value = createEmptyFriendForm()
    showToast(copy.toasts.createSuccess, 'success')
    await loadFriends()
  } catch (err) {
    console.error('[AdminFriends] createFriend failed:', err)
    showToast(copy.toasts.createFailed, 'error')
  } finally {
    creatingFriend.value = false
  }
}

function startEditFriend(link: FriendLink) {
  editingFriendID.value = link.id
  editFriend.value = {
    name: link.name,
    url: link.url,
    avatar_url: link.avatar_url,
    description: link.description,
    sort_order: link.sortOrder,
  }
  showEditForm.value = true
  showCreateForm.value = false
}

async function updateFriend() {
  if (!editingFriendID.value) return
  if (!editFriend.value.name.trim() || !editFriend.value.url.trim()) return
  updatingFriend.value = true
  try {
    await api.put(`/admin/friends/${editingFriendID.value}`, toFriendPayload(editFriend.value))
    showToast(copy.toasts.updateSuccess, 'success')
    closeEditForm()
    await loadFriends()
  } catch (err) {
    console.error('[AdminFriends] updateFriend failed:', err)
    showToast(copy.toasts.updateFailed, 'error')
  } finally {
    updatingFriend.value = false
  }
}

async function deleteFriend(id: string) {
  try {
    await api.delete(`/admin/friends/${id}`)
    friends.value = friends.value.filter((item) => item.id !== id)
    showToast(copy.toasts.deleteSuccess, 'success')
  } catch (err) {
    console.error('[AdminFriends] deleteFriend failed:', err)
    showToast(copy.toasts.deleteFailed, 'error')
  }
}

async function approveFriend(link: FriendLink) {
  processingFriendId.value = link.id
  try {
    await api.post(`/admin/friends/${link.id}/approve`)
    showToast(copy.toasts.approveFriendSuccess, 'success')
    await loadFriends()
  } catch (err) {
    console.error('[AdminFriends] approveFriend failed:', err)
    showToast(copy.toasts.loadFriendsFailed, 'error')
  } finally {
    processingFriendId.value = null
  }
}

async function rejectFriend(link: FriendLink) {
  processingFriendId.value = link.id
  try {
    await api.post(`/admin/friends/${link.id}/reject`)
    showToast(copy.toasts.rejectFriendSuccess, 'success')
    await loadFriends()
  } catch (err) {
    console.error('[AdminFriends] rejectFriend failed:', err)
    showToast(copy.toasts.loadFriendsFailed, 'error')
  } finally {
    processingFriendId.value = null
  }
}

async function approveApplication(application: FriendApplication) {
  processingApplicationId.value = application.id
  try {
    await api.post(`/admin/friends/applications/${application.id}/approve`)
    showToast(copy.toasts.approveApplicationSuccess, 'success')
    await loadAll()
  } catch (err) {
    console.error('[AdminFriends] approveApplication failed:', err)
    showToast(copy.toasts.loadApplicationsFailed, 'error')
  } finally {
    processingApplicationId.value = null
  }
}

async function rejectApplication(application: FriendApplication) {
  processingApplicationId.value = application.id
  try {
    await api.post(`/admin/friends/applications/${application.id}/reject`)
    showToast(copy.toasts.rejectApplicationSuccess, 'success')
    await loadApplications()
  } catch (err) {
    console.error('[AdminFriends] rejectApplication failed:', err)
    showToast(copy.toasts.loadApplicationsFailed, 'error')
  } finally {
    processingApplicationId.value = null
  }
}

onMounted(() => {
  loadAll()
})

const approvedCount = computed(() => friends.value.filter((item) => item.status === 'approved').length)
const pendingCount = computed(() => (
  friends.value.filter((item) => item.status === 'pending').length
  + applications.value.filter((item) => item.status === 'pending').length
))
const downCount = computed(() => friends.value.filter((item) => item.health === 'down').length)

function statusClass(status: FriendLink['status']) {
  if (status === 'approved') return 'bg-emerald-100 text-emerald-700'
  if (status === 'rejected') return 'bg-red-100 text-red-600'
  return 'bg-amber-100 text-amber-700'
}

function statusLabel(status: FriendLink['status']) {
  if (status === 'approved') return copy.status.approved
  if (status === 'rejected') return copy.status.rejected
  return copy.status.pending
}
</script>
