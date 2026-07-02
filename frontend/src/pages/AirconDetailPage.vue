<template>
  <section v-if="item" class="page">
    <header class="page-header">
      <div>
        <h1>{{ item.code }}</h1>
        <p>{{ item.building || item.room }}</p>
      </div>
      <div class="actions">
        <RouterLink class="ghost" :to="`/aircons/${item.id}/edit`"><Pencil :size="18" />แก้ไข</RouterLink>
        <button class="primary" @click="showRecord = true"><Plus :size="18" />บันทึกล้าง</button>
      </div>
    </header>
    <div class="detail-grid">
      <section class="panel">
        <div class="panel-header">
          <h2>สถานะ</h2>
          <StatusBadge :status="item.status" />
        </div>
        <dl class="meta-grid">
          <div><dt>ล้างล่าสุด</dt><dd><DateCell :value="item.latest_cleaning_date" /></dd></div>
          <div><dt>กำหนดถัดไป</dt><dd><DateCell :value="item.next_cleaning_date" /></dd></div>
          <div><dt>แผนล้าง</dt><dd><DateCell :value="item.planned_cleaning_date" /></dd></div>
          <div><dt>ทีม</dt><dd>{{ item.responsible_team || '-' }}</dd></div>
          <div><dt>ผู้ดูแล</dt><dd>{{ item.contact_name || '-' }}</dd></div>
          <div><dt>โทร</dt><dd>{{ item.contact_phone || '-' }}</dd></div>
        </dl>
      </section>
      <section class="panel">
        <h2>ตำแหน่งและรายละเอียด</h2>
        <dl class="meta-grid">
          <div><dt>อาคาร</dt><dd>{{ item.building || '-' }}</dd></div>
          <div><dt>ชั้น</dt><dd>{{ item.floor || '-' }}</dd></div>
          <div><dt>ห้อง/พื้นที่</dt><dd>{{ item.room || '-' }}</dd></div>
          <div><dt>ยี่ห้อ</dt><dd>{{ item.brand || '-' }}</dd></div>
          <div><dt>BTU</dt><dd>{{ item.btu || '-' }}</dd></div>
          <div><dt>พื้นที่</dt><dd>{{ [item.subdistrict, item.district, item.province].filter(Boolean).join(' / ') || '-' }}</dd></div>
        </dl>
        <p class="note">{{ item.note }}</p>
      </section>
    </div>
    <section class="panel">
      <div class="panel-header"><h2>ประวัติการล้าง</h2></div>
      <table>
        <thead><tr><th>วันที่ล้าง</th><th>วันที่วางแผน</th><th>ผู้ดำเนินการ</th><th>สถานะ</th><th>หมายเหตุ</th></tr></thead>
        <tbody>
          <tr v-for="record in item.cleaning_records || []" :key="record.id">
            <td><DateCell :value="record.cleaned_date" /></td>
            <td><DateCell :value="record.planned_date" /></td>
            <td>{{ record.performed_by || '-' }}</td>
            <td>{{ record.status || '-' }}</td>
            <td>{{ record.note || '-' }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="!item.cleaning_records?.length" class="empty-state">ยังไม่มีประวัติการล้าง</div>
    </section>
    <div v-if="showRecord" class="modal-backdrop">
      <form class="modal" @submit.prevent="saveRecord">
        <h3>บันทึกการล้าง</h3>
        <label>วันที่ล้าง<input v-model="record.cleaned_date" type="date" /></label>
        <label>วันที่วางแผน<input v-model="record.planned_date" type="date" /></label>
        <label>ผู้ดำเนินการ<input v-model="record.performed_by" /></label>
        <label>หมายเหตุ<textarea v-model="record.note"></textarea></label>
        <div class="modal-actions">
          <button type="button" class="ghost" @click="showRecord = false">ยกเลิก</button>
          <button class="primary">บันทึก</button>
        </div>
      </form>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Pencil, Plus } from 'lucide-vue-next'
import { api } from '../api/client'
import type { AirConditioner } from '../api/types'
import DateCell from '../components/DateCell.vue'
import StatusBadge from '../components/StatusBadge.vue'

const route = useRoute()
const item = ref<AirConditioner | null>(null)
const showRecord = ref(false)
const record = reactive({ cleaned_date: '', planned_date: '', performed_by: '', note: '', status: 'เสร็จเรียบร้อย' })

async function load() {
  item.value = (await api.get<AirConditioner>(`/api/aircons/${route.params.id}`)).data
}
async function saveRecord() {
  await api.post(`/api/aircons/${route.params.id}/cleaning-records`, {
    ...record,
    cleaned_date: asApiDate(record.cleaned_date),
    planned_date: asApiDate(record.planned_date)
  })
  showRecord.value = false
  await load()
}
onMounted(load)

function asApiDate(value?: string) {
  return value ? `${value}T00:00:00Z` : null
}
</script>
