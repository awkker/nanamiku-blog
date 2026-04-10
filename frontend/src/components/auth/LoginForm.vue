<template>
  <LiquidGlassCard width="100%" maxWidth="31rem" padding="34px" class="mx-auto">
    <div class="space-y-6">
      <div class="space-y-2">
        <p class="text-xs uppercase tracking-[0.3em] text-slate-600">{{ copy.hero.consoleLabel }}</p>
        <h1 class="text-3xl font-semibold text-slate-900 sm:text-[2rem]">{{ copy.hero.title }}</h1>
        <p class="text-sm leading-7 text-slate-700">{{ copy.hero.subtitle }}</p>
      </div>

      <form
        class="space-y-4 rounded-[28px] border border-white/70 bg-white/44 p-5 shadow-[inset_0_1px_0_rgba(255,255,255,0.78)] transition duration-300 sm:p-6"
        :class="hasAuthError ? 'ring-2 ring-red-300/80 ring-offset-2 ring-offset-transparent' : ''"
        @submit.prevent="handleSubmit"
      >
        <MikuInput
          v-model="identifier"
          :label="copy.form.identifierLabel"
          :placeholder="copy.form.identifierPlaceholder"
          autocomplete="username"
          :error="errors.identifier"
          :aria-label="copy.form.identifierAria"
          @focus="clearAuthError"
          required
        />

        <MikuInput
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          :label="copy.form.passwordLabel"
          :placeholder="copy.form.passwordPlaceholder"
          autocomplete="current-password"
          :error="errors.password"
          :aria-label="copy.form.passwordAria"
          @focus="clearAuthError"
          required
        >
          <template #trailing>
            <button
              type="button"
              class="rounded-lg p-1 text-slate-600 transition duration-300 hover:bg-white/30 hover:text-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-miku/70"
              :aria-label="showPassword ? copy.form.hidePasswordAria : copy.form.showPasswordAria"
              @click="showPassword = !showPassword"
            >
              <svg v-if="showPassword" viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
                <path d="M3 3l18 18" />
                <path d="M10.58 10.58A3 3 0 0012 15a3 3 0 001.42-.36" />
                <path d="M9.88 5.09A10.94 10.94 0 0112 5c6 0 10 7 10 7a18.47 18.47 0 01-4.45 5.01" />
                <path d="M6.61 6.61A18.81 18.81 0 002 12s1.62 2.84 4.7 4.77" />
              </svg>
              <svg v-else viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
                <path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7-10-7-10-7z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
            </button>
          </template>
        </MikuInput>

        <p
          v-if="formError"
          class="flex items-center gap-2 rounded-xl border border-red-300/90 bg-red-100/80 px-3 py-2 text-sm font-medium text-slate-900"
          role="status"
          aria-live="polite"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 shrink-0 fill-none stroke-current stroke-[2]">
            <path d="M12 8v4" />
            <path d="M12 16h.01" />
            <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3l-8.47-14.14a2 2 0 00-3.42 0z" />
          </svg>
          {{ formError }}
        </p>
        <p
          v-if="successMessage"
          class="flex items-center gap-2 rounded-xl border border-emerald-300/85 bg-emerald-100/80 px-3 py-2 text-sm font-medium text-slate-900"
          role="status"
          aria-live="polite"
        >
          <svg viewBox="0 0 24 24" class="h-4 w-4 shrink-0 fill-none stroke-current stroke-[2]">
            <path d="M20 6L9 17l-5-5" />
          </svg>
          {{ successMessage }}
        </p>

        <MikuButton
          type="submit"
          :loading="submitting"
          :disabled="submitting"
          :full-width="true"
          :aria-label="copy.form.submitAria"
        >
          {{ submitting ? copy.form.submittingButton : copy.form.submitButton }}
        </MikuButton>

        <p class="text-center text-xs text-slate-700" aria-live="polite">
          {{ submitting ? copy.form.submittingHint : copy.form.idleHint }}
        </p>
      </form>

      <div class="flex flex-wrap gap-2">
        <a
          v-for="link in copy.quickLinks"
          :key="link.href"
          :href="link.href"
          :aria-label="link.ariaLabel"
          class="rounded-2xl border border-slate-300/75 bg-white/55 px-3.5 py-2 text-xs font-medium text-slate-800 transition hover:border-miku/35 hover:text-miku"
        >
          {{ link.label }}
        </a>
      </div>
    </div>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, reactive, ref } from 'vue'

import { hydrateAuth, loginWithPassword } from '../../stores/auth'
import { adminCopy } from '../../content/copy'
import { getScopeLoading, setScopeStatus } from '../../stores/loading'
import MikuButton from '../ui/MikuButton.vue'
import MikuInput from '../ui/MikuInput.vue'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

const copy = adminCopy.loginForm
const identifier = ref('')
const password = ref('')
const showPassword = ref(false)
const formError = ref('')
const successMessage = ref('')

const errors = reactive({
  identifier: '',
  password: '',
})

const submitting = useStore(getScopeLoading('loginSubmit'))
const hasAuthError = computed(() => Boolean(formError.value))

onMounted(async () => {
  const user = await hydrateAuth()
  if (user) {
    window.location.replace('/admin')
  }
})

function validate() {
  errors.identifier = identifier.value.trim() ? '' : copy.form.emptyIdentifierError
  errors.password = password.value.trim() ? '' : copy.form.emptyPasswordError
  return !errors.identifier && !errors.password
}

function clearAuthError() {
  formError.value = ''
  successMessage.value = ''

  if (errors.identifier === copy.form.invalidCredentialsError) {
    errors.identifier = ''
  }

  if (errors.password === copy.form.invalidCredentialsError) {
    errors.password = ''
  }
}

async function handleSubmit() {
  formError.value = ''
  successMessage.value = ''

  if (!validate()) {
    return
  }

  setScopeStatus('loginSubmit', 'loading')

  try {
    await loginWithPassword(identifier.value, password.value)
    setScopeStatus('loginSubmit', 'success')
    successMessage.value = copy.form.successMessage
    window.setTimeout(() => {
      window.location.assign('/admin')
    }, 250)
  } catch (error) {
    setScopeStatus('loginSubmit', 'error')
    formError.value = error instanceof Error ? error.message : copy.form.fallbackError
    errors.identifier = copy.form.invalidCredentialsError
    errors.password = copy.form.invalidCredentialsError
  }
}
</script>
