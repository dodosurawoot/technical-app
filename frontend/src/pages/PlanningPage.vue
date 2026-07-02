<template>
  <section class="page">
    <header class="page-header">
      <div>
        <h1>แผนล้าง</h1>
        <p>งานที่มีแผน ใกล้ถึงกำหนด หรือเกินกำหนด</p>
      </div>
    </header>
    <section class="panel filters">
      <input v-model="start" type="date" @change="load" />
      <input v-model="end" type="date" @change="load" />
      <input v-model="plannedDate" type="date" />
      <button class="primary" :disabled="!selected.length || !plannedDate" @click="bulkUpdate">
        <CalendarCheck :size="18" />อัปเดตแผนที่เลือก
      </button>
    </section>
    <section class="panel">
      <AirconTable v-model:selected="selected" selectable :items="items" @open="open" />
    </section>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { CalendarCheck } from 'lucide-vue-next'
import { api } from '../api/client'
import type { AirConditioner } from '../api/types'
import AirconTable from '../components/AirconTable.vue'

const router = useRouter()
const items = ref<AirConditioner[]>([])
const selected = ref<number[]>([])
const plannedDate = ref('')
const start = ref('')
const end = ref('')

async function load() {
  const { data } = await api.get<{ items: AirConditioner[] }>('/api/plans', { params: { start_date: start.value, end_date: end.value } })
  items.value = data.items
}
async function bulkUpdate() {
  await api.post('/api/plans/bulk-update', { ids: selected.value, planned_date: `${plannedDate.value}T00:00:00Z` })
  selected.value = []
  plannedDate.value = ''
  await load()
}
function open(id: number) {
  router.push(`/aircons/${id}`)
}
onMounted(load)
</script>

