<template>
  <main class="center-page">
    <div class="spinner"></div>
    <p>กำลังเข้าสู่ระบบ...</p>
  </main>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { completeLogin } from '../stores/auth'

const router = useRouter()
onMounted(async () => {
  await completeLogin()
  router.replace('/')
})
</script>

