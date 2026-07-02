import { reactive } from 'vue'
import { User, UserManager, WebStorageStateStore } from 'oidc-client-ts'
import { api } from '../api/client'
import type { User as AppUser } from '../api/types'

const oidcEnabled = Boolean(import.meta.env.VITE_AUTHENTIK_ISSUER_URL && import.meta.env.VITE_AUTHENTIK_CLIENT_ID)

const manager = oidcEnabled
  ? new UserManager({
      authority: import.meta.env.VITE_AUTHENTIK_ISSUER_URL,
      client_id: import.meta.env.VITE_AUTHENTIK_CLIENT_ID,
      redirect_uri: import.meta.env.VITE_AUTHENTIK_REDIRECT_URL || `${window.location.origin}/auth/callback`,
      post_logout_redirect_uri: window.location.origin,
      response_type: 'code',
      scope: 'openid profile email',
      userStore: new WebStorageStateStore({ store: window.localStorage })
    })
  : null

export const authState = reactive({
  loading: true,
  user: null as AppUser | null,
  oidcEnabled
})

export async function getAccessToken(): Promise<string | null> {
  if (!manager) return null
  const user = await manager.getUser()
  return user?.id_token || user?.access_token || null
}

export async function loadMe() {
  authState.loading = true
  try {
    const { data } = await api.get<AppUser>('/api/me')
    authState.user = data
  } catch {
    authState.user = null
  } finally {
    authState.loading = false
  }
}

export async function login() {
  if (!manager) {
    await loadMe()
    return
  }
  await manager.signinRedirect()
}

export async function completeLogin(): Promise<User | null> {
  if (!manager) return null
  const user = await manager.signinRedirectCallback()
  await loadMe()
  return user
}

export async function logout() {
  authState.user = null
  if (manager) await manager.signoutRedirect()
}

