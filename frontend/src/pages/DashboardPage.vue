<template>
  <section class="page">
    <header class="page-header">
      <div>
        <h1>แดชบอร์ด</h1>
        <p>{{ today }}</p>
      </div>
      <RouterLink class="primary" to="/aircons/new"><Plus :size="18" />เพิ่มแอร์</RouterLink>
    </header>

    <div class="summary-grid">
      <SummaryCard label="แอร์ทั้งหมด" :value="data?.total_air_conditioners || 0" :icon="Wind" tone="blue" />
      <SummaryCard label="ล้างเดือนนี้" :value="data?.cleaned_this_month || 0" :icon="CheckCircle2" tone="green" />
      <SummaryCard label="ใกล้ถึงกำหนด" :value="data?.due_soon || 0" :icon="Clock3" tone="yellow" />
      <SummaryCard label="เกินกำหนด" :value="data?.overdue || 0" :icon="AlertTriangle" tone="red" />
      <SummaryCard label="มีแผนล้าง" :value="data?.planned_jobs || 0" :icon="CalendarCheck" tone="blue" />
    </div>

    <section class="panel">
      <div class="panel-header">
        <h2>งานที่ควรติดตาม</h2>
        <RouterLink to="/planning">ดูแผนทั้งหมด</RouterLink>
      </div>
      <AirconTable :items="data?.upcoming || []" @open="open" />
    </section>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { AlertTriangle, CalendarCheck, CheckCircle2, Clock3, Plus, Wind } from 'lucide-vue-next'
import { api } from '../api/client'
import type { Dashboard } from '../api/types'
import SummaryCard from '../components/SummaryCard.vue'
import AirconTable from '../components/AirconTable.vue'

const router = useRouter()
const data = ref<Dashboard | null>(null)
const today = new Intl.DateTimeFormat('th-TH', { dateStyle: 'full' }).format(new Date())
onMounted(async () => {
  data.value = (await api.get<Dashboard>('/api/dashboard')).data
})
function open(id: number) {
  router.push(`/aircons/${id}`)
}
</script>

