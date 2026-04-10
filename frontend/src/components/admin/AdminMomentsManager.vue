<template>
  <section class="space-y-5">
    <!-- Header -->
    <AdminPlainCard padding="24px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">{{ mm.header.title }}</h1>
          <p class="mt-1 text-sm text-slate-600">{{ mm.header.subtitle }}</p>
        </div>
        <MikuButton variant="solid" :aria-label="mm.header.title" data-testid="admin-moment-create-toggle" @click="toggleCreateForm">{{ mm.createButton }}</MikuButton>
      </div>
    </AdminPlainCard>

    <!-- ===== Compose Card (Create) ===== -->
    <AdminPlainCard v-if="showCreateForm" padding="0px">
      <form data-testid="admin-moment-create-form" @submit.prevent="createMoment">
        <!-- Content zone -->
        <div class="compose-content-zone">
          <textarea
            v-model="newMoment.content"
            data-testid="admin-moment-create-content"
            :placeholder="mm.createForm.contentPlaceholder"
            class="compose-editor"
          />
        </div>

        <!-- Image preview -->
        <div v-if="createImagePreviews.length > 0" class="flex gap-2 px-7 pb-3">
          <div v-for="(url, idx) in createImagePreviews" :key="idx" data-testid="admin-moment-create-image-preview" class="relative h-20 w-20 shrink-0 overflow-hidden rounded-xl border border-slate-200/50">
            <img :src="url" alt="" class="h-full w-full object-cover" @error="($event.target as HTMLImageElement).style.display='none'" />
          </div>
        </div>

        <!-- Metadata drawer -->
        <div class="border-t border-slate-200/40">
          <button
            type="button"
            class="flex w-full items-center gap-2 px-7 py-2.5 text-left text-xs font-medium uppercase tracking-widest text-slate-400 transition hover:text-slate-600"
            @click="showCreateMeta = !showCreateMeta"
          >
            <svg class="h-3.5 w-3.5 transition-transform" :class="showCreateMeta ? 'rotate-90' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6" /></svg>
            {{ mm.createForm.metaToggle }}
          </button>
          <div v-show="showCreateMeta" class="compose-meta-drawer">
            <p class="rounded-lg border border-slate-200/80 bg-white/70 px-3 py-2 text-xs text-slate-500">
              {{ copy.profile.autoHint }}
            </p>
            <div class="compose-meta-row">
              <label class="compose-meta-label">{{ mm.createForm.statusLabel }}</label>
              <select v-model="newMoment.publish_status" data-testid="admin-moment-create-status" class="compose-meta-input">
                <option value="draft">{{ mm.createForm.statusDraft }}</option>
                <option value="published">{{ mm.createForm.statusPublished }}</option>
                <option value="scheduled">{{ mm.createForm.statusScheduled }}</option>
              </select>
            </div>
            <div v-if="newMoment.publish_status === 'scheduled'" class="compose-meta-row">
              <label class="compose-meta-label">{{ mm.createForm.scheduledLabel }}</label>
              <input v-model="newMoment.scheduled_at" data-testid="admin-moment-create-scheduled-at" type="datetime-local" class="compose-meta-input" />
            </div>
            <div class="compose-meta-row">
              <label class="compose-meta-label">{{ mm.createForm.imageLabel }}</label>
              <input v-model="newMoment.image_urls" data-testid="admin-moment-create-image-input" type="text" :placeholder="mm.createForm.imagePlaceholder" class="compose-meta-input" />
            </div>
          </div>
        </div>

        <!-- Bottom bar -->
        <div class="flex items-center justify-between border-t border-slate-200/40 px-7 py-3.5">
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span class="inline-block h-1.5 w-1.5 rounded-full" :class="newMoment.content.length > 0 ? 'bg-miku' : 'bg-slate-300'" />
            {{ charCount(newMoment.content) }}
          </div>
          <div class="flex items-center gap-3">
            <button type="button" class="rounded-xl px-4 py-2 text-sm text-slate-400 transition hover:bg-slate-100/50 hover:text-slate-600" @click="closeCreateForm">{{ mm.createForm.cancelButton }}</button>
            <MikuButton type="submit" variant="solid" data-testid="admin-moment-create-submit" :disabled="creating">{{ creating ? mm.createForm.publishLoading : mm.createForm.publishIdle }}</MikuButton>
          </div>
        </div>
      </form>
    </AdminPlainCard>

    <!-- ===== Compose Card (Edit) ===== -->
    <AdminPlainCard v-if="showEditForm" padding="0px">
      <form data-testid="admin-moment-edit-form" @submit.prevent="updateMoment">
        <div class="compose-content-zone">
          <p class="mb-1 text-xs tracking-wide text-slate-400">{{ mm.createForm.editingBadge }}</p>
          <textarea
            v-model="editMoment.content"
            data-testid="admin-moment-edit-content"
            :placeholder="mm.createForm.contentPlaceholder"
            class="compose-editor"
          />
        </div>

        <div v-if="editImagePreviews.length > 0" class="flex gap-2 px-7 pb-3">
          <div v-for="(url, idx) in editImagePreviews" :key="idx" data-testid="admin-moment-edit-image-preview" class="relative h-20 w-20 shrink-0 overflow-hidden rounded-xl border border-slate-200/50">
            <img :src="url" alt="" class="h-full w-full object-cover" @error="($event.target as HTMLImageElement).style.display='none'" />
          </div>
        </div>

        <div class="border-t border-slate-200/40">
          <button
            type="button"
            class="flex w-full items-center gap-2 px-7 py-2.5 text-left text-xs font-medium uppercase tracking-widest text-slate-400 transition hover:text-slate-600"
            @click="showEditMeta = !showEditMeta"
          >
            <svg class="h-3.5 w-3.5 transition-transform" :class="showEditMeta ? 'rotate-90' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6" /></svg>
            {{ mm.createForm.metaToggle }}
          </button>
          <div v-show="showEditMeta" class="compose-meta-drawer">
            <p class="rounded-lg border border-slate-200/80 bg-white/70 px-3 py-2 text-xs text-slate-500">
              {{ copy.profile.autoHint }}
            </p>
            <div class="compose-meta-row">
              <label class="compose-meta-label">{{ mm.createForm.statusLabel }}</label>
              <select v-model="editMoment.publish_status" data-testid="admin-moment-edit-status" class="compose-meta-input">
                <option value="draft">{{ mm.createForm.statusDraft }}</option>
                <option value="published">{{ mm.createForm.statusPublished }}</option>
                <option value="scheduled">{{ mm.createForm.statusScheduled }}</option>
              </select>
            </div>
            <div v-if="editMoment.publish_status === 'scheduled'" class="compose-meta-row">
              <label class="compose-meta-label">{{ mm.createForm.scheduledLabel }}</label>
              <input v-model="editMoment.scheduled_at" data-testid="admin-moment-edit-scheduled-at" type="datetime-local" class="compose-meta-input" />
            </div>
            <div class="compose-meta-row">
              <label class="compose-meta-label">{{ mm.createForm.imageLabel }}</label>
              <input v-model="editMoment.image_urls" data-testid="admin-moment-edit-image-input" type="text" :placeholder="mm.createForm.imagePlaceholder" class="compose-meta-input" />
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between border-t border-slate-200/40 px-7 py-3.5">
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span class="inline-block h-1.5 w-1.5 rounded-full" :class="editMoment.content.length > 0 ? 'bg-miku' : 'bg-slate-300'" />
            {{ charCount(editMoment.content) }}
          </div>
          <div class="flex items-center gap-3">
            <button type="button" class="rounded-xl px-4 py-2 text-sm text-slate-400 transition hover:bg-slate-100/50 hover:text-slate-600" @click="closeEditForm">{{ mm.createForm.cancelButton }}</button>
            <MikuButton type="submit" variant="solid" data-testid="admin-moment-edit-submit" :disabled="editing">{{ editing ? mm.createForm.saveLoading : mm.createForm.saveIdle }}</MikuButton>
          </div>
        </div>
      </form>
    </AdminPlainCard>

    <!-- Stats -->
    <div class="grid gap-4 sm:grid-cols-3">
      <AdminPlainCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">{{ mm.stats.totalLabel }}</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-slate-900">{{ momentsList.length }}</p>
      </AdminPlainCard>
      <AdminPlainCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">{{ mm.stats.likesLabel }}</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-miku">{{ totalLikes }}</p>
      </AdminPlainCard>
      <AdminPlainCard padding="16px">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-600">{{ mm.stats.commentsLabel }}</p>
        <p class="mt-1 font-mono text-2xl font-semibold text-[#c084fc]">{{ totalComments }}</p>
      </AdminPlainCard>
    </div>

    <!-- List -->
    <AdminPlainCard padding="0px">
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="loading-dot" /><div class="loading-dot delay-1" /><div class="loading-dot delay-2" />
      </div>
      <div v-else-if="momentsList.length === 0" class="py-20 text-center">
        <p class="text-base text-slate-400">{{ mm.list.emptyMessage }}</p>
        <button type="button" class="mt-3 text-sm text-miku/80 transition hover:text-miku" @click="toggleCreateForm">{{ mm.list.firstMomentCta }} &rarr;</button>
      </div>
      <div v-else class="divide-y divide-slate-100/60">
        <div
          v-for="item in momentsList"
          :key="item.id"
          data-testid="admin-moment-row"
          class="group px-5 py-4 transition-colors hover:bg-white/40"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <img :src="item.authorAvatar" :alt="item.author" class="h-6 w-6 rounded-full border border-slate-200 object-cover" />
                <span class="text-sm font-medium text-slate-900">{{ item.author }}</span>
                <span class="rounded-full px-2 py-0.5 text-[11px] font-medium" :class="statusClass(item.publishStatus)">
                  {{ statusLabel(item.publishStatus) }}
                </span>
                <span class="text-xs text-slate-400">{{ item.displayTime }}</span>
              </div>
              <p class="mt-1.5 whitespace-pre-wrap text-sm leading-relaxed text-slate-700">{{ item.content }}</p>
              <div v-if="item.images.length > 0" class="mt-2.5 flex gap-2">
                <img
                  v-for="(img, idx) in item.images.slice(0, 4)"
                  :key="idx"
                  :src="img"
                  :alt="`${idx + 1}`"
                  data-testid="admin-moment-row-image"
                  class="h-16 w-16 rounded-xl border border-slate-100/80 object-cover transition-transform hover:scale-105"
                  loading="lazy"
                />
              </div>
              <div class="mt-2 flex items-center gap-4 text-xs text-slate-400">
                <span>{{ item.likes }} ❤</span>
                <span>{{ item.reposts }} ↻</span>
                <span>{{ item.comments }} ✉</span>
              </div>
            </div>
            <div class="shrink-0 opacity-0 transition group-hover:opacity-100">
              <div class="flex items-center gap-2">
                <button
                  type="button"
                  data-testid="admin-moment-edit-button"
                  class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-700 transition hover:border-miku/40 hover:text-miku"
                  :aria-label="mm.listActions.editAria"
                  @click="startEditMoment(item)"
                >
                  {{ mm.listActions.editButton }}
                </button>
                <button
                  v-if="item.publishStatus !== 'published'"
                  type="button"
                  data-testid="admin-moment-publish-button"
                  class="rounded-xl border border-emerald-200/80 bg-white/50 px-2.5 py-1 text-xs text-emerald-600 transition hover:bg-emerald-50"
                  :aria-label="mm.listActions.publishAria"
                  @click="publishMoment(item.id)"
                >
                  {{ mm.listActions.publishButton }}
                </button>
                <button
                  type="button"
                  data-testid="admin-moment-schedule-button"
                  class="rounded-xl border border-[#e9d5ff]/80 bg-white/50 px-2.5 py-1 text-xs text-[#9333ea] transition hover:bg-[#faf5ff]"
                  :aria-label="mm.listActions.scheduleAria"
                  @click="scheduleMoment(item.id)"
                >
                  {{ mm.listActions.scheduleButton }}
                </button>
                <button
                  v-if="item.publishStatus !== 'draft'"
                  type="button"
                  data-testid="admin-moment-unpublish-button"
                  class="rounded-xl border border-slate-200/80 bg-white/50 px-2.5 py-1 text-xs text-slate-600 transition hover:bg-slate-50"
                  :aria-label="mm.listActions.unpublishAria"
                  @click="unpublishMoment(item.id)"
                >
                  {{ mm.listActions.unpublishButton }}
                </button>
                <button
                  type="button"
                  data-testid="admin-moment-delete-button"
                  class="rounded-xl border border-red-200/85 bg-white/50 px-2.5 py-1 text-xs text-red-600 transition hover:bg-red-50"
                  :aria-label="copy.actions.deleteAria"
                  @click="deleteMoment(item.id)"
                >
                  {{ copy.actions.deleteButton }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminPlainCard>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { api, ApiError, type PagedData } from '../../lib/api'
import { DEFAULT_PUBLIC_AVATAR_URL } from '../../lib/default-assets'
import { adminCopy } from '../../content/copy'
import { showToast } from '../../stores/ui'
import AdminPlainCard from '../ui/AdminPlainCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface ApiMoment {
  id: string
  author_name: string
  author_avatar_url?: string
  content: string
  image_urls: string[]
  like_count: number
  repost_count: number
  comment_count: number
  publish_status?: 'draft' | 'published' | 'scheduled'
  published_at?: string
  scheduled_at?: string
  created_at: string
}

interface MomentItem {
  id: string
  author: string
  authorAvatar: string
  content: string
  images: string[]
  likes: number
  reposts: number
  comments: number
  publishStatus: 'draft' | 'published' | 'scheduled'
  publishedAt: string
  scheduledAt: string
  scheduledAtISO: string
  displayTime: string
  createdAt: string
}

interface MomentForm {
  content: string
  image_urls: string
  publish_status: 'draft' | 'published' | 'scheduled'
  scheduled_at: string
}

function formatDate(iso: string): string {
  try {
    const d = new Date(iso)
    return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
  } catch {
    return iso
  }
}

function mapMoment(item: ApiMoment): MomentItem {
  const publishStatus = item.publish_status || 'published'
  const publishedAt = item.published_at ? formatDate(item.published_at) : '--'
  const scheduledAt = item.scheduled_at ? formatDate(item.scheduled_at) : '--'
  return {
    id: item.id,
    author: item.author_name,
    authorAvatar: item.author_avatar_url || DEFAULT_PUBLIC_AVATAR_URL,
    content: item.content,
    images: item.image_urls || [],
    likes: Number(item.like_count) || 0,
    reposts: Number(item.repost_count) || 0,
    comments: Number(item.comment_count) || 0,
    publishStatus,
    publishedAt,
    scheduledAt,
    scheduledAtISO: item.scheduled_at || '',
    displayTime: publishStatus === 'scheduled' ? scheduledAt : publishedAt,
    createdAt: formatDate(item.created_at),
  }
}

function createEmptyMomentForm(): MomentForm {
  return {
    content: '',
    image_urls: '',
    publish_status: 'draft',
    scheduled_at: '',
  }
}

function toImageURLs(input: string): string[] {
  return input.split(',').map((u: string) => u.trim()).filter(Boolean)
}

function formatDateInputLocal(date: Date): string {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60 * 1000)
  return local.toISOString().slice(0, 16)
}

function localInputToRFC3339(value: string): string {
  return new Date(value).toISOString()
}

const momentsList = ref<MomentItem[]>([])
const loading = ref(false)
const showCreateForm = ref(false)
const creating = ref(false)
const showEditForm = ref(false)
const editing = ref(false)
const editingMomentID = ref<string | null>(null)

const newMoment = ref<MomentForm>(createEmptyMomentForm())
const editMoment = ref<MomentForm>(createEmptyMomentForm())
const showCreateMeta = ref(false)
const showEditMeta = ref(false)
const mm = adminCopy.momentsManager
const copy = mm

const createImagePreviews = computed(() => {
  return newMoment.value.image_urls.split(',').map((u: string) => u.trim()).filter(Boolean).slice(0, 4)
})
const editImagePreviews = computed(() => {
  return editMoment.value.image_urls.split(',').map((u: string) => u.trim()).filter(Boolean).slice(0, 4)
})

function charCount(text: string): string {
  const len = text.length
  if (len === 0) return mm.charCount.zero
  return `${len}${mm.charCount.suffix}`
}

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
  editingMomentID.value = null
  editMoment.value = createEmptyMomentForm()
}

const totalLikes = computed(() => momentsList.value.reduce((sum, m) => sum + m.likes, 0))
const totalComments = computed(() => momentsList.value.reduce((sum, m) => sum + m.comments, 0))

async function loadMoments() {
  loading.value = true
  try {
    const data = await api.get<PagedData<ApiMoment>>('/admin/moments?size=50')
    momentsList.value = (data.items || []).map(mapMoment)
  } catch (err) {
    console.error('[AdminMoments] loadMoments failed:', err)
    showToast(mm.toast.loadFailed, 'error')
    momentsList.value = []
  } finally {
    loading.value = false
  }
}

async function createMoment() {
  if (!newMoment.value.content.trim()) return
  if (newMoment.value.publish_status === 'scheduled' && !newMoment.value.scheduled_at) {
    showToast(mm.toast.scheduledAtRequired, 'error')
    return
  }
  creating.value = true
  try {
    await api.post('/admin/moments', {
      content: newMoment.value.content.trim(),
      image_urls: toImageURLs(newMoment.value.image_urls),
      publish_status: newMoment.value.publish_status,
      scheduled_at: newMoment.value.publish_status === 'scheduled'
        ? localInputToRFC3339(newMoment.value.scheduled_at)
        : undefined,
    })
    closeCreateForm()
    newMoment.value = createEmptyMomentForm()
    showToast(mm.toast.createSuccess, 'success')
    await loadMoments()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : mm.toast.createFailed
    console.error('[AdminMoments] createMoment failed:', err)
    showToast(msg, 'error')
  } finally {
    creating.value = false
  }
}

function startEditMoment(item: MomentItem) {
  editingMomentID.value = item.id
  editMoment.value = {
    content: item.content,
    image_urls: (item.images || []).join(', '),
    publish_status: item.publishStatus,
    scheduled_at: item.publishStatus === 'scheduled' && item.scheduledAtISO
      ? formatDateInputLocal(new Date(item.scheduledAtISO))
      : '',
  }
  showEditForm.value = true
  showCreateForm.value = false
}

async function updateMoment() {
  if (!editingMomentID.value) return
  if (!editMoment.value.content.trim()) return
  if (editMoment.value.publish_status === 'scheduled' && !editMoment.value.scheduled_at) {
    showToast(mm.toast.scheduledAtRequired, 'error')
    return
  }
  editing.value = true
  try {
    await api.put(`/admin/moments/${editingMomentID.value}`, {
      content: editMoment.value.content.trim(),
      image_urls: toImageURLs(editMoment.value.image_urls),
      publish_status: editMoment.value.publish_status,
      scheduled_at: editMoment.value.publish_status === 'scheduled'
        ? localInputToRFC3339(editMoment.value.scheduled_at)
        : undefined,
    })
    showToast(mm.toast.updateSuccess, 'success')
    closeEditForm()
    await loadMoments()
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : mm.toast.updateFailed
    console.error('[AdminMoments] updateMoment failed:', err)
    showToast(msg, 'error')
  } finally {
    editing.value = false
  }
}

async function publishMoment(id: string) {
  try {
    await api.post(`/admin/moments/${id}/publish`)
    await loadMoments()
    showToast(mm.toast.publishSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : mm.toast.publishFailed
    showToast(msg, 'error')
  }
}

async function scheduleMoment(id: string) {
  const defaultTime = formatDateInputLocal(new Date(Date.now() + 30 * 60 * 1000))
  const next = window.prompt(mm.listActions.schedulePrompt, defaultTime)
  if (!next) return

  try {
    await api.post(`/admin/moments/${id}/schedule`, {
      scheduled_at: localInputToRFC3339(next),
    })
    await loadMoments()
    showToast(mm.toast.scheduleSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : mm.toast.scheduleFailed
    showToast(msg, 'error')
  }
}

async function unpublishMoment(id: string) {
  try {
    await api.post(`/admin/moments/${id}/unpublish`)
    await loadMoments()
    showToast(mm.toast.unpublishSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : mm.toast.unpublishFailed
    showToast(msg, 'error')
  }
}

async function deleteMoment(id: string) {
  if (!window.confirm(copy.actions.deleteConfirm)) {
    return
  }

  try {
    await api.delete(`/admin/moments/${id}`)
    await loadMoments()
    showToast(copy.actions.deleteSuccess, 'success')
    if (editingMomentID.value === id) {
      closeEditForm()
    }
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : copy.actions.deleteFailed
    showToast(msg, 'error')
  }
}

onMounted(() => {
  loadMoments()
})

function statusClass(status: MomentItem['publishStatus']) {
  if (status === 'published') return 'bg-emerald-100 text-emerald-700'
  if (status === 'draft') return 'bg-slate-100 text-slate-600'
  return 'bg-[#f3e8ff] text-[#9333ea]'
}

function statusLabel(status: MomentItem['publishStatus']) {
  if (status === 'published') return mm.status.published
  if (status === 'draft') return mm.status.draft
  return mm.status.scheduled
}
</script>

<style scoped>
/* ---- Compose Content Zone ---- */
.compose-content-zone {
  padding: 28px 28px 12px;
}

.compose-editor {
  width: 100%;
  min-height: 140px;
  resize: vertical;
  border: none;
  outline: none;
  background: transparent;
  font-size: 0.9375rem;
  line-height: 1.75;
  color: #1e293b;
  caret-color: rgb(57, 197, 187);
}

.compose-editor::placeholder {
  color: #cbd5e1;
}

.compose-editor::selection {
  background: rgba(57, 197, 187, 0.18);
}

/* ---- Compose Metadata ---- */
.compose-meta-drawer {
  padding: 0 28px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.compose-meta-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.compose-meta-label {
  flex-shrink: 0;
  width: 40px;
  font-size: 0.6875rem;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #94a3b8;
  text-align: right;
}

.compose-meta-input {
  flex: 1;
  border-radius: 10px;
  border: 1px solid rgba(203, 213, 225, 0.5);
  background: rgba(255, 255, 255, 0.4);
  padding: 7px 12px;
  font-size: 0.8125rem;
  color: #1e293b;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.compose-meta-input:focus {
  border-color: rgba(57, 197, 187, 0.4);
  box-shadow: 0 0 0 2px rgba(57, 197, 187, 0.08);
}

.compose-meta-input::placeholder {
  color: #b0bec5;
}

/* ---- Loading Dots ---- */
.loading-dot {
  width: 6px;
  height: 6px;
  margin: 0 4px;
  border-radius: 50%;
  background: rgb(57, 197, 187);
  opacity: 0.35;
  animation: dot-pulse 1.2s ease-in-out infinite;
}

.loading-dot.delay-1 { animation-delay: 0.2s; }
.loading-dot.delay-2 { animation-delay: 0.4s; }

@keyframes dot-pulse {
  0%, 80%, 100% { opacity: 0.2; transform: scale(0.8); }
  40% { opacity: 1; transform: scale(1.1); }
}
</style>
