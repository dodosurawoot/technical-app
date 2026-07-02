<template>
  <section class="page">
    <header class="page-header">
      <div>
        <h1>รายการแอร์</h1>
        <p>ค้นหา กรอง และเปิดประวัติการล้าง</p>
      </div>
      <RouterLink class="primary" to="/aircons/new"><Plus :size="18" />เพิ่มแอร์</RouterLink>
    </header>
    <section class="panel filters">
      <input v-model="filters.q" placeholder="ค้นหารหัสหรือสถานที่" @input="load" />
      <select v-model="filters.status" @change="load">
        <option value="">ทุกสถานะ</option>
        <option value="normal">ปกติ</option>
        <option value="due_soon">ใกล้ถึงกำหนด</option>
        <option value="overdue">เกินกำหนด</option>
        <option value="planned">วางแผนแล้ว</option>
        <option value="never_cleaned">ยังไม่เคยบันทึก</option>
      </select>
      <input v-model="filters.responsible_team" placeholder="ทีมรับผิดชอบ" @input="load" />
      <input v-model="filters.start_date" type="date" @change="load" />
      <input v-model="filters.end_date" type="date" @change="load" />
    </section>
    <section class="panel">
      <AirconTable :items="items" @open="open" />
    </section>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { Plus } from 'lucide-vue-next'
import { api } from '../api/client'
import type { AirConditioner } from '../api/types'
import AirconTable from '../components/AirconTable.vue'

const router = useRouter()
const items = ref<AirConditioner[]>([])
const filters = reactive({ q: '', status: '', responsible_team: '', start_date: '', end_date: '' })
let timer: number | undefined

async function load() {
  window.clearTimeout(timer)
  timer = window.setTimeout(async () => {
    const { data } = await api.get<{ items: AirConditioner[] }>('/api/aircons', { params: filters })
    items.value = data.items
  }, 150)
}
function open(id: number) {
  router.push(`/aircons/${id}`)
}
onMounted(load)
</script>

