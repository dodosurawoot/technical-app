<template>
  <section class="page">
    <header class="page-header">
      <div>
        <h1>{{ isEdit ? 'แก้ไขแอร์' : 'เพิ่มแอร์' }}</h1>
        <p>{{ form.code || 'กรอกข้อมูลแอร์และรอบล้าง' }}</p>
      </div>
    </header>
    <form class="panel form-grid" @submit.prevent="save">
      <label>รหัส / เลขครุภัณฑ์<input v-model="form.code" required /></label>
      <label>อาคาร<input v-model="form.building" /></label>
      <label>ชั้น<input v-model="form.floor" /></label>
      <label>ห้อง / พื้นที่<input v-model="form.room" /></label>
      <label>ยี่ห้อ<input v-model="form.brand" /></label>
      <label>BTU<input v-model.number="form.btu" type="number" min="0" /></label>
      <label>ทีมรับผิดชอบ<input v-model="form.responsible_team" /></label>
      <label>วันที่ล้างล่าสุด<input v-model="form.latest_cleaning_date" type="date" /></label>
      <label>วันที่วางแผน<input v-model="form.planned_cleaning_date" type="date" /></label>
      <label>ผู้ดูแล<input v-model="form.contact_name" /></label>
      <label>เบอร์โทร<input v-model="form.contact_phone" /></label>
      <label>ตำบล/แขวง<input v-model="form.subdistrict" /></label>
      <label>อำเภอ/เขต<input v-model="form.district" /></label>
      <label>จังหวัด<input v-model="form.province" /></label>
      <label>Lat<input v-model.number="form.latitude" type="number" step="0.000001" /></label>
      <label>Long<input v-model.number="form.longitude" type="number" step="0.000001" /></label>
      <label class="wide">หมายเหตุ<textarea v-model="form.note"></textarea></label>
      <div class="form-actions">
        <RouterLink class="ghost" to="/aircons">ยกเลิก</RouterLink>
        <button class="primary">บันทึก</button>
      </div>
    </form>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { api } from '../api/client'
import type { AirConditioner } from '../api/types'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => Boolean(route.params.id))
const form = reactive({
  code: '',
  building: '',
  floor: '',
  room: '',
  brand: '',
  btu: 0,
  responsible_team: '',
  latest_cleaning_date: '',
  planned_cleaning_date: '',
  note: '',
  contact_name: '',
  contact_phone: '',
  subdistrict: '',
  district: '',
  province: '',
  latitude: undefined as number | undefined,
  longitude: undefined as number | undefined
})

onMounted(async () => {
  if (!isEdit.value) return
  const { data } = await api.get<AirConditioner>(`/api/aircons/${route.params.id}`)
  Object.assign(form, {
    ...data,
    latest_cleaning_date: asInputDate(data.latest_cleaning_date),
    planned_cleaning_date: asInputDate(data.planned_cleaning_date)
  })
})

async function save() {
  const payload = {
    ...form,
    latest_cleaning_date: asApiDate(form.latest_cleaning_date),
    planned_cleaning_date: asApiDate(form.planned_cleaning_date)
  }
  if (isEdit.value) await api.put(`/api/aircons/${route.params.id}`, payload)
  else await api.post('/api/aircons', payload)
  router.push('/aircons')
}

function asInputDate(value?: string | null) {
  return value ? value.slice(0, 10) : ''
}

function asApiDate(value?: string) {
  return value ? `${value}T00:00:00Z` : null
}
</script>

