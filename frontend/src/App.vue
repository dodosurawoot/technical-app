<template>
  <div v-if="authState.user" class="shell">
    <aside class="sidebar">
      <RouterLink to="/" class="brand">
        <span class="brand-mark">A</span>
        <span>AirClean Tracker</span>
      </RouterLink>
      <nav>
        <RouterLink v-for="item in navItems" :key="item.path" :to="item.path">
          <component :is="item.icon" :size="18" />
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
      <div class="account">
        <div>
          <strong>{{ authState.user.name || authState.user.email }}</strong>
          <small>{{ roleLabel(authState.user.role) }}</small>
        </div>
        <button class="icon-button" title="ออกจากระบบ" @click="logout">
          <LogOut :size="18" />
        </button>
      </div>
    </aside>
    <main class="main">
      <RouterView />
    </main>
  </div>
  <RouterView v-else />
</template>

<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { LayoutDashboard, Wind, CalendarDays, Users, Settings, LogOut } from 'lucide-vue-next'
import { authState, logout } from './stores/auth'
import { roleLabel } from './utils'

const navItems = [
  { path: '/', label: 'แดชบอร์ด', icon: LayoutDashboard },
  { path: '/aircons', label: 'รายการแอร์', icon: Wind },
  { path: '/planning', label: 'แผนล้าง', icon: CalendarDays },
  { path: '/admin/users', label: 'ผู้ใช้', icon: Users },
  { path: '/settings', label: 'ตั้งค่า', icon: Settings }
]
</script>

