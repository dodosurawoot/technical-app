<template>
  <main class="login-page">
    <section class="login-panel">
      <div class="brand large">
        <span class="brand-mark">A</span>
        <span>AirClean Tracker</span>
      </div>
      <h1>เข้าสู่ระบบ</h1>
      <p>จัดการรอบล้างแอร์ ดูงานใกล้ครบกำหนด และติดตามประวัติได้ในที่เดียว</p>
      <button class="primary" @click="login">
        <LogIn :size="18" />
        เข้าสู่ระบบด้วย Authentik
      </button>
      <small v-if="!authState.oidcEnabled">โหมดพัฒนาเปิดใช้งานอยู่</small>
    </section>
  </main>
</template>

<script setup lang="ts">
import { LogIn } from 'lucide-vue-next'
import { authState, login } from '../stores/auth'
</script>

