<template>
  <section class="page">
    <header class="page-header">
      <div>
        <h1>ตั้งค่า</h1>
        <p>รอบล้างเริ่มต้น 6 เดือน และแจ้งใกล้ครบกำหนดภายใน 30 วัน</p>
      </div>
    </header>
    <section class="panel settings-grid">
      <div>
        <small>Issuer</small>
        <strong>{{ issuer || 'โหมดพัฒนา' }}</strong>
      </div>
      <div>
        <small>Client ID</small>
        <strong>{{ clientId || '-' }}</strong>
      </div>
      <div>
        <small>Backend</small>
        <strong>{{ backend || windowOrigin }}</strong>
      </div>
      <div>
        <small>Redirect URI</small>
        <strong>{{ redirect || `${windowOrigin}/auth/callback` }}</strong>
      </div>
    </section>

    <section class="panel import-panel">
      <div class="panel-header">
        <div>
          <h2>นำเข้า Excel</h2>
          <small>รองรับไฟล์ .xlsx จากแบบฟอร์มล้างแอร์เดิม</small>
        </div>
      </div>
      <form class="upload-row" @submit.prevent="uploadExcel">
        <input ref="fileInput" type="file" accept=".xlsx,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" @change="onFileChange" />
        <button class="primary" :disabled="!selectedFile || uploading">
          <Upload :size="18" />
          {{ uploading ? 'กำลังนำเข้า...' : 'นำเข้าไฟล์' }}
        </button>
      </form>
      <p v-if="message" class="notice" :class="messageTone">{{ message }}</p>
      <div v-if="result" class="import-result">
        <div>
          <small>เพิ่มแอร์ใหม่</small>
          <strong>{{ result.imported_air_conditioners.toLocaleString('th-TH') }}</strong>
        </div>
        <div>
          <small>เพิ่มประวัติการล้าง</small>
          <strong>{{ result.imported_records.toLocaleString('th-TH') }}</strong>
        </div>
        <div>
          <small>ข้ามแถว</small>
          <strong>{{ result.skipped_rows.toLocaleString('th-TH') }}</strong>
        </div>
      </div>
      <ul v-if="result?.warnings?.length" class="warning-list">
        <li v-for="warning in result.warnings.slice(0, 8)" :key="warning">{{ warning }}</li>
      </ul>
    </section>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Upload } from 'lucide-vue-next'
import { api } from '../api/client'

interface ImportResult {
  imported_air_conditioners: number
  imported_records: number
  skipped_rows: number
  warnings: string[]
}

const issuer = import.meta.env.VITE_AUTHENTIK_ISSUER_URL
const clientId = import.meta.env.VITE_AUTHENTIK_CLIENT_ID
const backend = import.meta.env.VITE_BACKEND_URL
const redirect = import.meta.env.VITE_AUTHENTIK_REDIRECT_URL
const windowOrigin = window.location.origin
const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const uploading = ref(false)
const message = ref('')
const messageTone = ref<'success' | 'error'>('success')
const result = ref<ImportResult | null>(null)

function onFileChange(event: Event) {
  selectedFile.value = (event.target as HTMLInputElement).files?.[0] || null
  result.value = null
  message.value = selectedFile.value ? `เลือกไฟล์: ${selectedFile.value.name}` : ''
  messageTone.value = 'success'
}

async function uploadExcel() {
  if (!selectedFile.value) return
  uploading.value = true
  message.value = ''
  result.value = null
  try {
    const form = new FormData()
    form.append('file', selectedFile.value)
    const { data } = await api.post<ImportResult>('/api/import/excel', form)
    result.value = data
    message.value = 'นำเข้า Excel สำเร็จ'
    messageTone.value = 'success'
    selectedFile.value = null
    if (fileInput.value) fileInput.value.value = ''
  } catch (error: any) {
    message.value = error?.response?.data?.error || 'นำเข้าไฟล์ไม่สำเร็จ'
    messageTone.value = 'error'
  } finally {
    uploading.value = false
  }
}
</script>
