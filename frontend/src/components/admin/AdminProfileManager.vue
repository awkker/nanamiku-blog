<template>
  <section class="space-y-5">
    <AdminPlainCard padding="24px">
      <div>
        <h1 class="text-2xl font-semibold text-slate-900">{{ copy.page.title }}</h1>
        <p class="mt-1 text-sm text-slate-600">{{ copy.page.subtitle }}</p>
      </div>
    </AdminPlainCard>

    <AdminPlainCard padding="20px">
      <div class="grid gap-4 lg:grid-cols-[auto_1fr_auto] lg:items-center">
        <img
          :src="profileAvatarPreview"
          :alt="copy.profile.avatarAlt"
          class="h-16 w-16 rounded-full border border-slate-200 object-cover"
        />
        <div>
          <p class="text-sm font-semibold text-slate-900">{{ copy.profile.title }}</p>
          <p class="mt-1 text-xs text-slate-500">{{ copy.profile.subtitle }}</p>
        </div>
      </div>

      <div class="mt-4 grid gap-3 md:grid-cols-[220px_1fr_auto]">
        <input
          v-model="profileDisplayName"
          type="text"
          :placeholder="copy.profile.displayNamePlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.profile.displayNameLabel"
        />
        <input
          v-model="profileAvatarURL"
          type="text"
          :placeholder="copy.profile.avatarPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.profile.avatarLabel"
        />
        <MikuButton type="button" variant="solid" :disabled="profileSaving" @click="saveProfile">
          {{ profileSaving ? copy.profile.savingButton : copy.profile.saveButton }}
        </MikuButton>
      </div>
    </AdminPlainCard>

    <AdminPlainCard padding="20px">
      <div>
        <p class="text-sm font-semibold text-slate-900">{{ copy.account.title }}</p>
        <p class="mt-1 text-xs text-slate-500">{{ copy.account.subtitle }}</p>
      </div>

      <div class="mt-4 grid gap-3 md:grid-cols-[1fr_1fr]">
        <input
          v-model="accountUsername"
          type="text"
          :placeholder="copy.account.usernamePlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.account.usernameLabel"
        />
        <input
          v-model="accountEmail"
          type="email"
          :placeholder="copy.account.emailPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.account.emailLabel"
        />
      </div>

      <div class="mt-3 grid gap-3 md:grid-cols-[1fr_auto]">
        <input
          v-model="accountPassword"
          type="password"
          :placeholder="copy.account.passwordPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.account.passwordLabel"
          autocomplete="new-password"
        />
        <MikuButton type="button" variant="solid" :disabled="accountSaving" @click="saveAccount">
          {{ accountSaving ? copy.account.savingButton : copy.account.saveButton }}
        </MikuButton>
      </div>
    </AdminPlainCard>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref } from 'vue'

import { adminCopy } from '../../content/copy'
import { api, ApiError } from '../../lib/api'
import { authState, hydrateAuth, updateMyAccount, updateMyProfile } from '../../stores/auth'
import { showToast } from '../../stores/ui'
import AdminPlainCard from '../ui/AdminPlainCard.vue'
import MikuButton from '../ui/MikuButton.vue'

interface AdminProfilePayload {
  username: string
  email: string
  display_name?: string
  avatar_url?: string
}

const copy = adminCopy.profileManager
const auth = useStore(authState)

const profileDisplayName = ref('')
const profileAvatarURL = ref('/picture/author.jpg')
const profileSaving = ref(false)
const accountUsername = ref('')
const accountEmail = ref('')
const accountPassword = ref('')
const accountSaving = ref(false)

const profileAvatarPreview = computed(() => {
  return (profileAvatarURL.value || '').trim() || '/picture/author.jpg'
})

async function loadProfile() {
  hydrateAuth()
  const current = auth.value.user
  if (current) {
    profileDisplayName.value = current.name || current.username || ''
    profileAvatarURL.value = current.avatar || '/picture/author.jpg'
    accountUsername.value = current.username || ''
    accountEmail.value = current.email || ''
  }

  try {
    const me = await api.get<AdminProfilePayload>('/auth/me')
    profileDisplayName.value = (me.display_name || '').trim() || me.username
    profileAvatarURL.value = (me.avatar_url || '').trim() || '/picture/author.jpg'
    accountUsername.value = (me.username || '').trim()
    accountEmail.value = (me.email || '').trim()
  } catch {
    // keep store fallback
  }
}

async function saveProfile() {
  const displayName = profileDisplayName.value.trim()
  if (!displayName) {
    showToast(copy.profile.emptyNameError, 'error')
    return
  }

  profileSaving.value = true
  try {
    await updateMyProfile(displayName, profileAvatarURL.value.trim())
    await loadProfile()
    showToast(copy.profile.saveSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : copy.profile.saveFailed
    showToast(msg, 'error')
  } finally {
    profileSaving.value = false
  }
}

async function saveAccount() {
  const username = accountUsername.value.trim()
  const email = accountEmail.value.trim()
  const newPassword = accountPassword.value.trim()

  if (!username) {
    showToast(copy.account.emptyUsernameError, 'error')
    return
  }
  if (!email) {
    showToast(copy.account.emptyEmailError, 'error')
    return
  }

  accountSaving.value = true
  try {
    await updateMyAccount(username, email, newPassword)
    accountPassword.value = ''
    await loadProfile()
    showToast(copy.account.saveSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : copy.account.saveFailed
    showToast(msg, 'error')
  } finally {
    accountSaving.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>
