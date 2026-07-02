CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  provider_subject TEXT,
  email TEXT NOT NULL UNIQUE,
  name TEXT,
  username TEXT,
  role TEXT NOT NULL DEFAULT 'viewer',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS air_conditioners (
  id BIGSERIAL PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  building TEXT,
  floor TEXT,
  room TEXT,
  brand TEXT,
  btu BIGINT,
  responsible_team TEXT,
  latest_cleaning_date TIMESTAMPTZ,
  next_cleaning_date TIMESTAMPTZ,
  planned_cleaning_date TIMESTAMPTZ,
  status TEXT,
  note TEXT,
  contact_name TEXT,
  contact_phone TEXT,
  subdistrict TEXT,
  district TEXT,
  province TEXT,
  latitude NUMERIC,
  longitude NUMERIC,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cleaning_records (
  id BIGSERIAL PRIMARY KEY,
  air_conditioner_id BIGINT NOT NULL REFERENCES air_conditioners(id) ON DELETE CASCADE,
  cleaned_date TIMESTAMPTZ,
  planned_date TIMESTAMPTZ,
  reported_date TIMESTAMPTZ,
  status TEXT,
  note TEXT,
  evidence_url TEXT,
  performed_by TEXT,
  created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cleaning_plans (
  id BIGSERIAL PRIMARY KEY,
  air_conditioner_id BIGINT NOT NULL REFERENCES air_conditioners(id) ON DELETE CASCADE,
  planned_date TIMESTAMPTZ NOT NULL,
  status TEXT NOT NULL DEFAULT 'planned',
  note TEXT,
  responsible_team TEXT,
  created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  action TEXT,
  entity TEXT,
  entity_id TEXT,
  details TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_air_conditioners_status ON air_conditioners(status);
CREATE INDEX IF NOT EXISTS idx_air_conditioners_next_cleaning_date ON air_conditioners(next_cleaning_date);
CREATE INDEX IF NOT EXISTS idx_air_conditioners_planned_cleaning_date ON air_conditioners(planned_cleaning_date);
CREATE INDEX IF NOT EXISTS idx_air_conditioners_building ON air_conditioners(building);
CREATE INDEX IF NOT EXISTS idx_air_conditioners_room ON air_conditioners(room);
CREATE INDEX IF NOT EXISTS idx_air_conditioners_responsible_team ON air_conditioners(responsible_team);
CREATE INDEX IF NOT EXISTS idx_cleaning_records_cleaned_date ON cleaning_records(cleaned_date);
CREATE INDEX IF NOT EXISTS idx_cleaning_records_air_conditioner_id ON cleaning_records(air_conditioner_id);
CREATE INDEX IF NOT EXISTS idx_cleaning_plans_planned_date ON cleaning_plans(planned_date);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

