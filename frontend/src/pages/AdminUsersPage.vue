<template>
  <section class="page">
    <header class="page-header">
      <div>
        <h1>ผู้ใช้และสิทธิ์</h1>
        <p>กำหนดบทบาทจากบัญชีที่เข้าสู่ระบบแล้ว</p>
      </div>
    </header>
    <section class="panel">
      <table>
        <thead><tr><th>ชื่อ</th><th>อีเมล</th><th>ผู้ใช้</th><th>บทบาท</th></tr></thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.name || '-' }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.username || '-' }}</td>
            <td>
              <select :value="user.role" @change="updateRole(user.id, ($event.target as HTMLSelectElement).value)">
                <option value="admin">ผู้ดูแลระบบ</option>
                <option value="team">ทีมงาน</option>
                <option value="viewer">ดูอย่างเดียว</option>
              </select>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="!users.length" class="empty-state">ยังไม่มีผู้ใช้</div>
    </section>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '../api/client'
import type { User } from '../api/types'

const users = ref<User[]>([])
async function load() {
  users.value = (await api.get<{ items: User[] }>('/api/users')).data.items
}
async function updateRole(id: number, role: string) {
  await api.put(`/api/users/${id}/role`, { role })
  await load()
}
onMounted(load)
</script>

