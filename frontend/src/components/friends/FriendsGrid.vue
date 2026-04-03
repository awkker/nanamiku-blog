<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex items-center gap-2.5">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#c084fc]/10">
          <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-[#c084fc] stroke-[2]" aria-hidden="true">
            <path d="M10 13a5 5 0 007.07 0l2.12-2.12a5 5 0 00-7.07-7.07L10 5M14 11a5 5 0 00-7.07 0L4.81 13.12a5 5 0 007.07 7.07L14 19" />
          </svg>
        </div>
        <div>
          <h2 class="text-base font-semibold text-slate-800">{{ copy.title }}</h2>
          <p class="text-xs text-slate-400">{{ linkSummary }}</p>
        </div>
      </div>
      <MikuButton
        variant="solid"
        class="w-full justify-center sm:w-auto"
        :aria-label="copy.applyButtonAria"
        @click="toggleApplicationForm"
      >
        <svg viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]" aria-hidden="true">
          <path d="M12 5v14M5 12h14" />
        </svg>
        {{ showApplicationForm ? copy.closeButton : copy.applyButton }}
      </MikuButton>
    </div>

    <div
      v-if="showApplicationForm"
      class="rounded-[28px] border border-white/65 bg-white/70 p-5 shadow-[0_18px_50px_rgba(15,23,42,0.08)] backdrop-blur-xl sm:p-6"
    >
      <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <p class="text-xs uppercase tracking-[0.22em] text-[#c084fc]">{{ copy.form.title }}</p>
          <h3 class="mt-2 text-xl font-semibold text-slate-900">{{ copy.form.title }}</h3>
          <p class="mt-2 max-w-2xl text-sm leading-relaxed text-slate-500">{{ copy.form.description }}</p>
        </div>
      </div>

      <form class="mt-5 space-y-4" @submit.prevent="submitApplication">
        <div class="grid gap-4 md:grid-cols-2">
          <label class="space-y-1.5">
            <span class="text-sm font-medium text-slate-700">{{ copy.form.siteNameLabel }}</span>
            <input
              v-model="applicationForm.site_name"
              type="text"
              :placeholder="copy.form.siteNamePlaceholder"
              class="w-full rounded-2xl border border-slate-200/80 bg-white/85 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-2 focus:ring-miku/20"
            />
          </label>

          <label class="space-y-1.5">
            <span class="text-sm font-medium text-slate-700">{{ copy.form.siteUrlLabel }}</span>
            <input
              v-model="applicationForm.site_url"
              type="text"
              :placeholder="copy.form.siteUrlPlaceholder"
              class="w-full rounded-2xl border border-slate-200/80 bg-white/85 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-2 focus:ring-miku/20"
            />
          </label>

          <label class="space-y-1.5">
            <span class="text-sm font-medium text-slate-700">{{ copy.form.contactEmailLabel }}</span>
            <input
              v-model="applicationForm.contact_email"
              type="email"
              :placeholder="copy.form.contactEmailPlaceholder"
              class="w-full rounded-2xl border border-slate-200/80 bg-white/85 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-2 focus:ring-miku/20"
            />
          </label>

          <label class="space-y-1.5">
            <span class="text-sm font-medium text-slate-700">{{ copy.form.avatarUrlLabel }}</span>
            <input
              v-model="applicationForm.avatar_url"
              type="text"
              :placeholder="copy.form.avatarUrlPlaceholder"
              class="w-full rounded-2xl border border-slate-200/80 bg-white/85 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-miku/50 focus:ring-2 focus:ring-miku/20"
            />
          </label>
        </div>

        <label class="block space-y-1.5">
          <span class="text-sm font-medium text-slate-700">{{ copy.form.descriptionLabel }}</span>
          <textarea
            v-model="applicationForm.description"
            rows="4"
            :placeholder="copy.form.descriptionPlaceholder"
            class="w-full rounded-[24px] border border-slate-200/80 bg-white/85 px-4 py-3 text-sm leading-relaxed text-slate-900 outline-none transition focus:border-miku/50 focus:ring-2 focus:ring-miku/20"
          />
        </label>

        <label class="block space-y-1.5">
          <span class="text-sm font-medium text-slate-700">{{ copy.form.contactNoteLabel }}</span>
          <textarea
            v-model="applicationForm.contact_note"
            rows="3"
            :placeholder="copy.form.contactNotePlaceholder"
            class="w-full rounded-[24px] border border-slate-200/80 bg-white/85 px-4 py-3 text-sm leading-relaxed text-slate-900 outline-none transition focus:border-miku/50 focus:ring-2 focus:ring-miku/20"
          />
        </label>

        <div class="flex flex-col gap-3 border-t border-slate-200/70 pt-4 sm:flex-row sm:items-center sm:justify-between">
          <p class="text-xs leading-relaxed text-slate-400">{{ copy.form.footnote }}</p>
          <div class="flex items-center gap-3">
            <button
              type="button"
              class="rounded-xl px-3 py-2 text-sm text-slate-500 transition hover:bg-white/80 hover:text-slate-700"
              @click="closeApplicationForm"
            >
              {{ copy.form.cancelButton }}
            </button>
            <MikuButton type="submit" :loading="submittingApplication" :disabled="submittingApplication">
              {{ submittingApplication ? copy.form.submittingButton : copy.form.submitButton }}
            </MikuButton>
          </div>
        </div>
      </form>
    </div>

    <ErrorState
      v-if="status === 'error'"
      :description="error || copy.errorFallback"
      @retry="loadLinks"
    />

    <div v-else>
      <div v-if="status === 'loading'" class="grid gap-5 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <SkeletonCard v-for="item in 6" :key="item" />
      </div>

      <EmptyState
        v-else-if="links.length === 0"
        :title="copy.emptyTitle"
        :description="copy.emptyDescription"
      />

      <div v-else class="grid gap-5 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <FriendLinkCard
          v-for="item in links"
          :key="item.id"
          :friend="item"
        />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref } from 'vue'

import { siteCopy } from '../../content/copy'
import { api, ApiError } from '../../lib/api'
import { friendError, friendFetchStatus, friendLinks, loadFriendLinks } from '../../stores/friends'
import { showToast } from '../../stores/ui'
import EmptyState from '../ui/EmptyState.vue'
import ErrorState from '../ui/ErrorState.vue'
import FriendLinkCard from './FriendLinkCard.vue'
import MikuButton from '../ui/MikuButton.vue'
import SkeletonCard from '../ui/SkeletonCard.vue'

interface FriendApplicationForm {
  site_name: string
  site_url: string
  avatar_url: string
  description: string
  contact_email: string
  contact_note: string
}

const copy = siteCopy.friendsPage.grid
const links = useStore(friendLinks)
const status = useStore(friendFetchStatus)
const error = useStore(friendError)
const showApplicationForm = ref(false)
const submittingApplication = ref(false)
const applicationForm = ref<FriendApplicationForm>(createEmptyApplicationForm())

const linkSummary = computed(() => (
  links.value.length > 0
    ? `${copy.countPrefix}${links.value.length}${copy.countSuffix}`
    : copy.buildingStatus
))

function createEmptyApplicationForm(): FriendApplicationForm {
  return {
    site_name: '',
    site_url: '',
    avatar_url: '',
    description: '',
    contact_email: '',
    contact_note: '',
  }
}

async function loadLinks() {
  await loadFriendLinks()
}

function toggleApplicationForm() {
  showApplicationForm.value = !showApplicationForm.value
}

function closeApplicationForm() {
  showApplicationForm.value = false
}

async function submitApplication() {
  if (
    !applicationForm.value.site_name.trim()
    || !applicationForm.value.site_url.trim()
    || !applicationForm.value.description.trim()
    || !applicationForm.value.contact_email.trim()
  ) {
    showToast(copy.validationRequired, 'info')
    return
  }

  submittingApplication.value = true
  try {
    await api.post('/friends/applications', {
      ...applicationForm.value,
      site_name: applicationForm.value.site_name.trim(),
      site_url: applicationForm.value.site_url.trim(),
      avatar_url: applicationForm.value.avatar_url.trim(),
      description: applicationForm.value.description.trim(),
      contact_email: applicationForm.value.contact_email.trim(),
      contact_note: applicationForm.value.contact_note.trim(),
    })
    applicationForm.value = createEmptyApplicationForm()
    closeApplicationForm()
    showToast(copy.applySuccess, 'success')
  } catch (err) {
    if (err instanceof ApiError && err.status === 409) {
      showToast(copy.applyDuplicate, 'info')
    } else {
      showToast(copy.applyFailed, 'error')
    }
  } finally {
    submittingApplication.value = false
  }
}

onMounted(async () => {
  await loadLinks()
})
</script>
