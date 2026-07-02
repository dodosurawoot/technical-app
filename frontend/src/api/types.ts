export type Role = 'admin' | 'team' | 'viewer'
export type Status = 'never_cleaned' | 'normal' | 'due_soon' | 'overdue' | 'planned'

export interface User {
  id: number
  email: string
  name: string
  username: string
  role: Role
}

export interface AirConditioner {
  id: number
  code: string
  building: string
  floor: string
  room: string
  brand: string
  btu: number
  responsible_team: string
  latest_cleaning_date: string | null
  next_cleaning_date: string | null
  planned_cleaning_date: string | null
  status: Status
  note: string
  contact_name: string
  contact_phone: string
  subdistrict: string
  district: string
  province: string
  latitude?: number
  longitude?: number
  cleaning_records?: CleaningRecord[]
}

export interface CleaningRecord {
  id: number
  air_conditioner_id: number
  cleaned_date: string | null
  planned_date: string | null
  reported_date: string | null
  status: string
  note: string
  evidence_url: string
  performed_by: string
  created_at: string
}

export interface Dashboard {
  total_air_conditioners: number
  cleaned_this_month: number
  due_soon: number
  overdue: number
  planned_jobs: number
  upcoming: AirConditioner[]
}

