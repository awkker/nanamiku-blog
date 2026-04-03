<template>
  <section class="space-y-5">
    <AdminPlainCard padding="24px">
      <div>
        <h1 class="text-2xl font-semibold text-slate-900">{{ copy.page.title }}</h1>
        <p class="mt-1 text-sm text-slate-600">{{ copy.page.subtitle }}</p>
      </div>
    </AdminPlainCard>

    <AdminPlainCard padding="20px">
      <div>
        <p class="text-sm font-semibold text-slate-900">{{ copy.siteProfile.title }}</p>
        <p class="mt-1 text-xs text-slate-500">{{ copy.siteProfile.subtitle }}</p>
      </div>

      <div class="mt-5">
        <p class="text-sm font-semibold text-slate-900">{{ copy.siteProfile.brandTitle }}</p>
        <p class="mt-1 text-xs text-slate-500">{{ copy.siteProfile.brandHint }}</p>
      </div>

      <div class="mt-4 grid gap-3 md:grid-cols-2">
        <input
          v-model="brandText"
          type="text"
          :placeholder="copy.siteProfile.brandTextPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.siteProfile.brandTextLabel"
          :disabled="loading || savingSiteProfile"
        />
        <input
          v-model="siteTitle"
          type="text"
          :placeholder="copy.siteProfile.siteTitlePlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.siteProfile.siteTitleLabel"
          :disabled="loading || savingSiteProfile"
        />
        <input
          v-model="logoAlt"
          type="text"
          :placeholder="copy.siteProfile.logoAltPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50 md:col-span-2"
          :aria-label="copy.siteProfile.logoAltLabel"
          :disabled="loading || savingSiteProfile"
        />
      </div>

      <div class="mt-6">
        <p class="text-sm font-semibold text-slate-900">{{ copy.siteProfile.seoTitle }}</p>
        <p class="mt-1 text-xs text-slate-500">{{ copy.siteProfile.seoHint }}</p>
      </div>

      <div class="mt-4 grid gap-3">
        <input
          v-model="siteUrl"
          type="text"
          :placeholder="copy.siteProfile.siteUrlPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.siteProfile.siteUrlLabel"
          :disabled="loading || savingSiteProfile"
        />
        <textarea
          v-model="defaultDescription"
          rows="3"
          :placeholder="copy.siteProfile.defaultDescriptionPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.siteProfile.defaultDescriptionLabel"
          :disabled="loading || savingSiteProfile"
        />
        <input
          v-model="defaultSocialImage"
          type="text"
          :placeholder="copy.siteProfile.defaultSocialImagePlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.siteProfile.defaultSocialImageLabel"
          :disabled="loading || savingSiteProfile"
        />
      </div>

      <div class="mt-5 flex flex-wrap items-center gap-3">
        <MikuButton type="button" variant="solid" :disabled="loading || savingSiteProfile" :aria-label="copy.siteProfile.saveButtonAria" @click="saveSiteProfile">
          {{ savingSiteProfile ? copy.siteProfile.savingButton : copy.siteProfile.saveButton }}
        </MikuButton>

        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-sm text-slate-600 transition hover:border-miku/40 hover:text-miku"
          :aria-label="copy.siteProfile.resetButtonAria"
          :disabled="loading || savingSiteProfile"
          @click="resetSiteProfile"
        >
          {{ copy.siteProfile.resetButton }}
        </button>
      </div>
    </AdminPlainCard>

    <AdminPlainCard padding="20px">
      <div>
        <p class="text-sm font-semibold text-slate-900">{{ copy.form.icpTitle }}</p>
        <p class="mt-1 text-xs text-slate-500">{{ copy.form.icpHint }}</p>
      </div>

      <div class="mt-4 grid gap-3 md:grid-cols-[1fr_1fr]">
        <input
          v-model="icpText"
          type="text"
          :placeholder="copy.form.icpPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.form.icpLabel"
          :disabled="loading || saving"
        />
        <input
          v-model="icpLink"
          type="text"
          :placeholder="copy.form.icpLinkPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.form.icpLinkLabel"
          :disabled="loading || saving"
        />
      </div>
    </AdminPlainCard>

    <AdminPlainCard padding="20px">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p class="text-sm font-semibold text-slate-900">{{ copy.form.customTitle }}</p>
          <p class="mt-1 text-xs text-slate-500">{{ copy.form.customHint }}</p>
        </div>

        <MikuButton
          type="button"
          variant="solid"
          :disabled="loading || saving"
          :aria-label="copy.form.addButtonAria"
          @click="addLine"
        >
          {{ copy.form.addButton }}
        </MikuButton>
      </div>

      <div class="mt-4 space-y-2.5">
        <div
          v-for="(line, index) in customTexts"
          :key="`footer-custom-${index}`"
          class="grid gap-2 md:grid-cols-[1fr_auto]"
        >
          <input
            v-model="customTexts[index]"
            type="text"
            :placeholder="copy.form.customLinePlaceholder"
            class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
            :aria-label="`${copy.form.customLineLabelPrefix}${index + 1}`"
            :disabled="loading || saving"
          />
          <button
            type="button"
            class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-500 transition hover:border-red-300 hover:text-red-600"
            :aria-label="copy.form.removeButtonAria"
            :disabled="loading || saving"
            @click="removeLine(index)"
          >
            {{ copy.form.removeButton }}
          </button>
        </div>

        <p v-if="customTexts.length === 0" class="rounded-xl border border-dashed border-slate-200 bg-white/60 px-3 py-3 text-xs text-slate-500">
          {{ copy.form.emptyCustomHint }}
        </p>
      </div>
    </AdminPlainCard>

    <AdminPlainCard padding="20px">
      <div class="flex flex-wrap items-center gap-3">
        <MikuButton type="button" variant="solid" :disabled="loading || saving" :aria-label="copy.form.saveButtonAria" @click="save">
          {{ saving ? copy.form.savingButton : copy.form.saveButton }}
        </MikuButton>

        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-sm text-slate-600 transition hover:border-miku/40 hover:text-miku"
          :aria-label="copy.form.resetButtonAria"
          :disabled="loading || saving"
          @click="reset"
        >
          {{ copy.form.resetButton }}
        </button>
      </div>
    </AdminPlainCard>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { onMounted, ref } from 'vue'

import { adminCopy } from '../../content/copy'
import { ApiError } from '../../lib/api'
import {
  hydrateSiteFooterSettings,
  resetSiteFooterSettings,
  saveSiteFooterSettings,
  siteFooterSettings,
} from '../../stores/siteFooter'
import {
  hydrateSiteProfileSettings,
  resetSiteProfileSettings,
  saveSiteProfileSettings,
  siteProfileSettings,
} from '../../stores/siteProfile'
import { showToast } from '../../stores/ui'
import AdminPlainCard from '../ui/AdminPlainCard.vue'
import MikuButton from '../ui/MikuButton.vue'

const copy = adminCopy.footerManager

const settings = useStore(siteFooterSettings)
const siteProfile = useStore(siteProfileSettings)

const icpText = ref('')
const icpLink = ref('')
const customTexts = ref<string[]>([])
const brandText = ref('')
const siteTitle = ref('')
const logoAlt = ref('')
const siteUrl = ref('')
const defaultDescription = ref('')
const defaultSocialImage = ref('')
const saving = ref(false)
const savingSiteProfile = ref(false)
const loading = ref(false)

function syncFromStore() {
  icpText.value = settings.value.icpText
  icpLink.value = settings.value.icpLink
  customTexts.value = [...settings.value.customTexts]
}

function syncSiteProfileFromStore() {
  brandText.value = siteProfile.value.brandText
  siteTitle.value = siteProfile.value.siteTitle
  logoAlt.value = siteProfile.value.logoAlt
  siteUrl.value = siteProfile.value.siteUrl
  defaultDescription.value = siteProfile.value.defaultDescription
  defaultSocialImage.value = siteProfile.value.defaultSocialImage
}

function addLine() {
  customTexts.value.push('')
}

function removeLine(index: number) {
  customTexts.value = customTexts.value.filter((_, lineIndex) => lineIndex !== index)
}

function cleanLines(lines: string[]): string[] {
  return lines
    .map((line) => line.trim())
    .filter((line) => line.length > 0)
}

async function load() {
  loading.value = true
  try {
    await Promise.all([
      hydrateSiteFooterSettings(),
      hydrateSiteProfileSettings(),
    ])
  } finally {
    syncFromStore()
    syncSiteProfileFromStore()
    loading.value = false
  }
}

function validateSiteProfile() {
  return [
    brandText.value,
    siteTitle.value,
    siteUrl.value,
    defaultDescription.value,
    defaultSocialImage.value,
  ].every((value) => value.trim().length > 0)
}

async function saveSiteProfile() {
  if (!validateSiteProfile()) {
    showToast(copy.siteProfile.validationRequired, 'error')
    return
  }

  savingSiteProfile.value = true
  try {
    await saveSiteProfileSettings({
      brandText: brandText.value.trim(),
      siteTitle: siteTitle.value.trim(),
      logoAlt: logoAlt.value.trim(),
      siteUrl: siteUrl.value.trim(),
      defaultDescription: defaultDescription.value.trim(),
      defaultSocialImage: defaultSocialImage.value.trim(),
    })

    syncSiteProfileFromStore()
    showToast(copy.siteProfileToast.saveSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : copy.siteProfileToast.saveFailed
    showToast(msg, 'error')
  } finally {
    savingSiteProfile.value = false
  }
}

async function resetSiteProfile() {
  savingSiteProfile.value = true
  try {
    await resetSiteProfileSettings()
    syncSiteProfileFromStore()
    showToast(copy.siteProfileToast.resetSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : copy.siteProfileToast.saveFailed
    showToast(msg, 'error')
  } finally {
    savingSiteProfile.value = false
  }
}

async function save() {
  saving.value = true

  try {
    await saveSiteFooterSettings({
      icpText: icpText.value.trim(),
      icpLink: icpLink.value.trim(),
      customTexts: cleanLines(customTexts.value),
    })

    syncFromStore()
    showToast(copy.toast.saveSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : copy.toast.saveFailed
    showToast(msg, 'error')
  } finally {
    saving.value = false
  }
}

async function reset() {
  saving.value = true
  try {
    await resetSiteFooterSettings()
    syncFromStore()
    showToast(copy.toast.resetSuccess, 'success')
  } catch (err) {
    const msg = err instanceof ApiError ? err.message : copy.toast.saveFailed
    showToast(msg, 'error')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  void load()
})
</script>
