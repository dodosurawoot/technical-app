import type { Role, Status } from './api/types'

export function formatDate(value?: string | null) {
  if (!value) return '-'
  return new Intl.DateTimeFormat('th-TH', { dateStyle: 'medium' }).format(new Date(value))
}

export function statusLabel(status?: Status | string) {
  const labels: Record<string, string> = {
    never_cleaned: 'ยังไม่เคยบันทึกล้าง',
    normal: 'ปกติ',
    due_soon: 'ใกล้ถึงกำหนด',
    overdue: 'เกินกำหนด',
    planned: 'วางแผนแล้ว'
  }
  return labels[status || ''] || status || '-'
}

export function statusTone(status?: Status | string) {
  const tones: Record<string, string> = {
    never_cleaned: 'muted',
    normal: 'green',
    due_soon: 'yellow',
    overdue: 'red',
    planned: 'blue'
  }
  return tones[status || ''] || 'muted'
}

export function roleLabel(role: Role | string) {
  const labels: Record<string, string> = { admin: 'ผู้ดูแลระบบ', team: 'ทีมงาน', viewer: 'ดูอย่างเดียว' }
  return labels[role] || role
}
