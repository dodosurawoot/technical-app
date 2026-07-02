import axios from 'axios'
import { getAccessToken } from '../stores/auth'

export const api = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL || '',
  timeout: 20000
})

api.interceptors.request.use(async (config) => {
  const token = await getAccessToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

