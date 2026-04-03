<template>
  <section class="space-y-5">
    <AdminPlainCard padding="24px">
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">{{ copy.title }}</h1>
          <p class="mt-1 text-sm text-slate-600">{{ copy.subtitle }}</p>
        </div>
      </div>
    </AdminPlainCard>

    <div class="grid gap-4 md:grid-cols-2">
      <AdminPlainCard padding="20px">
        <h2 class="text-lg font-semibold text-slate-900">{{ copy.jsonCardTitle }}</h2>
        <p class="mt-1 text-sm text-slate-600">{{ copy.jsonCardDescription }}</p>
        <div class="mt-4">
          <MikuButton variant="solid" :disabled="loadingFormat !== null" @click="download('json')">
            {{ loadingFormat === 'json' ? copy.exporting : copy.jsonButton }}
          </MikuButton>
        </div>
      </AdminPlainCard>

      <AdminPlainCard padding="20px">
        <h2 class="text-lg font-semibold text-slate-900">{{ copy.sqlCardTitle }}</h2>
        <p class="mt-1 text-sm text-slate-600">{{ copy.sqlCardDescription }}</p>
        <div class="mt-4">
          <MikuButton variant="solid" :disabled="loadingFormat !== null" @click="download('sql')">
            {{ loadingFormat === 'sql' ? copy.exporting : copy.sqlButton }}
          </MikuButton>
        </div>
      </AdminPlainCard>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { ApiError } from '../../lib/api'
import { fetchWithSessionRetry } from '../../lib/auth-session'
import { adminCopy } from '../../content/copy'
import { showToast } from '../../stores/ui'
import AdminPlainCard from '../ui/AdminPlainCard.vue'
import MikuButton from '../ui/MikuButton.vue'

type BackupFormat = 'json' | 'sql'

const copy = adminCopy.backup
const loadingFormat = ref<BackupFormat | null>(null)

function resolveDownloadFilename(contentDisposition: string | null, fallback: string): string {
  if (!contentDisposition) return fallback
  const match = contentDisposition.match(/filename="?([^";]+)"?/i)
  if (!match || !match[1]) return fallback
  return match[1]
}

async function download(format: BackupFormat) {
  loadingFormat.value = format
  try {
    const res = await fetchWithSessionRetry(`/api/v1/admin/backup/export?format=${format}`, {
      method: 'GET',
    })

    if (!res.ok) {
      let message = copy.failed
      try {
        const body = await (res.json() as Promise<{ message?: string }>)
        message = body.message || message
      } catch {
        // ignore json parse error
      }
      throw new ApiError(message, -1, res.status)
    }

    const blob = await res.blob()
    const fallback = `miku-backup.${format}`
    const filename = resolveDownloadFilename(res.headers.get('Content-Disposition'), fallback)

    const url = window.URL.createObjectURL(blob)
    const anchor = document.createElement('a')
    anchor.href = url
    anchor.download = filename
    document.body.appendChild(anchor)
    anchor.click()
    anchor.remove()
    window.URL.revokeObjectURL(url)

    showToast(copy.success, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.failed
    showToast(message, 'error')
  } finally {
    loadingFormat.value = null
  }
}
</script>
