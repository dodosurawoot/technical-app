import { createRouter, createWebHistory } from 'vue-router'
import { authState, loadMe } from '../stores/auth'
import LoginPage from '../pages/LoginPage.vue'
import DashboardPage from '../pages/DashboardPage.vue'
import AirconListPage from '../pages/AirconListPage.vue'
import AirconDetailPage from '../pages/AirconDetailPage.vue'
import AirconFormPage from '../pages/AirconFormPage.vue'
import PlanningPage from '../pages/PlanningPage.vue'
import AdminUsersPage from '../pages/AdminUsersPage.vue'
import SettingsPage from '../pages/SettingsPage.vue'
import CallbackPage from '../pages/CallbackPage.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: LoginPage, meta: { public: true } },
    { path: '/auth/callback', component: CallbackPage, meta: { public: true } },
    { path: '/', component: DashboardPage },
    { path: '/aircons', component: AirconListPage },
    { path: '/aircons/new', component: AirconFormPage },
    { path: '/aircons/:id', component: AirconDetailPage },
    { path: '/aircons/:id/edit', component: AirconFormPage },
    { path: '/history', component: AirconListPage },
    { path: '/planning', component: PlanningPage },
    { path: '/admin/users', component: AdminUsersPage },
    { path: '/settings', component: SettingsPage }
  ]
})

router.beforeEach(async (to) => {
  if (authState.loading) await loadMe()
  if (!to.meta.public && !authState.user) return '/login'
  if (to.path === '/login' && authState.user) return '/'
})

