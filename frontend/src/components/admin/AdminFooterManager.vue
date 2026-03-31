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
        />
        <input
          v-model="icpLink"
          type="text"
          :placeholder="copy.form.icpLinkPlaceholder"
          class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 outline-none transition focus:border-miku/50"
          :aria-label="copy.form.icpLinkLabel"
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
          />
          <button
            type="button"
            class="rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-500 transition hover:border-red-300 hover:text-red-600"
            :aria-label="copy.form.removeButtonAria"
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
        <MikuButton type="button" variant="solid" :disabled="saving" :aria-label="copy.form.saveButtonAria" @click="save">
          {{ saving ? copy.form.savingButton : copy.form.saveButton }}
        </MikuButton>

        <button
          type="button"
          class="rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-sm text-slate-600 transition hover:border-miku/40 hover:text-miku"
          :aria-label="copy.form.resetButtonAria"
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
import { onMounted, ref, watch } from 'vue'

import { adminCopy } from '../../content/copy'
import {
  hydrateSiteFooterSettings,
  resetSiteFooterSettings,
  saveSiteFooterSettings,
  siteFooterSettings,
} from '../../stores/siteFooter'
import { showToast } from '../../stores/ui'
import AdminPlainCard from '../ui/AdminPlainCard.vue'
import MikuButton from '../ui/MikuButton.vue'

const copy = adminCopy.footerManager

const settings = useStore(siteFooterSettings)

const icpText = ref('')
const icpLink = ref('')
const customTexts = ref<string[]>([])
const saving = ref(false)

function syncFromStore() {
  icpText.value = settings.value.icpText
  icpLink.value = settings.value.icpLink
  customTexts.value = [...settings.value.customTexts]
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

function save() {
  saving.value = true

  try {
    saveSiteFooterSettings({
      icpText: icpText.value.trim(),
      icpLink: icpLink.value.trim(),
      customTexts: cleanLines(customTexts.value),
    })

    syncFromStore()
    showToast(copy.toast.saveSuccess, 'success')
  } catch {
    showToast(copy.toast.saveFailed, 'error')
  } finally {
    saving.value = false
  }
}

function reset() {
  resetSiteFooterSettings()
  syncFromStore()
  showToast(copy.toast.resetSuccess, 'success')
}

onMounted(() => {
  hydrateSiteFooterSettings()
  syncFromStore()
})

watch(settings, () => {
  syncFromStore()
})
</script>
