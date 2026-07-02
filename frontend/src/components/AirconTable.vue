<template>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th v-if="selectable"><input type="checkbox" :checked="allSelected" @change="toggleAll" /></th>
          <th>รหัส</th>
          <th>สถานที่</th>
          <th>ทีมรับผิดชอบ</th>
          <th>ล้างล่าสุด</th>
          <th>กำหนดถัดไป</th>
          <th>แผน</th>
          <th>สถานะ</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id" @click="$emit('open', item.id)">
          <td v-if="selectable" @click.stop>
            <input type="checkbox" :checked="selected.includes(item.id)" @change="toggle(item.id)" />
          </td>
          <td><strong>{{ item.code }}</strong></td>
          <td>
            <div>{{ item.building || item.room || '-' }}</div>
            <small>{{ [item.subdistrict, item.district, item.province].filter(Boolean).join(' / ') }}</small>
          </td>
          <td>{{ item.responsible_team || '-' }}</td>
          <td><DateCell :value="item.latest_cleaning_date" /></td>
          <td><DateCell :value="item.next_cleaning_date" /></td>
          <td><DateCell :value="item.planned_cleaning_date" /></td>
          <td><StatusBadge :status="item.status" /></td>
        </tr>
      </tbody>
    </table>
    <div v-if="!items.length" class="empty-state">ยังไม่มีข้อมูลแอร์</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AirConditioner } from '../api/types'
import DateCell from './DateCell.vue'
import StatusBadge from './StatusBadge.vue'

const props = defineProps<{ items: AirConditioner[]; selected?: number[]; selectable?: boolean }>()
const emit = defineEmits<{ open: [id: number]; 'update:selected': [ids: number[]] }>()
const selected = computed(() => props.selected || [])
const allSelected = computed(() => props.items.length > 0 && props.items.every((item) => selected.value.includes(item.id)))

function toggle(id: number) {
  const next = selected.value.includes(id) ? selected.value.filter((item) => item !== id) : [...selected.value, id]
  emit('update:selected', next)
}

function toggleAll(event: Event) {
  const checked = (event.target as HTMLInputElement).checked
  emit('update:selected', checked ? props.items.map((item) => item.id) : [])
}
</script>

