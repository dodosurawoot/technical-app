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
  </section>
</template>

<script setup lang="ts">
const issuer = import.meta.env.VITE_AUTHENTIK_ISSUER_URL
const clientId = import.meta.env.VITE_AUTHENTIK_CLIENT_ID
const backend = import.meta.env.VITE_BACKEND_URL
const redirect = import.meta.env.VITE_AUTHENTIK_REDIRECT_URL
const windowOrigin = window.location.origin
</script>

