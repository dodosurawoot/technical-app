# AirClean Tracker

เว็บแอปสำหรับติดตามรอบล้างแอร์ทุก 6 เดือน พร้อมแดชบอร์ด แผนงาน ประวัติการล้าง และสิทธิ์ผู้ใช้ผ่าน Authentik OAuth2/OIDC

## โครงสร้าง

- `frontend` Vue 3 + Vite + TypeScript
- `backend` Go REST API + GORM + PostgreSQL
- `deploy` บันทึกการ deploy
- `docs` เอกสาร mapping และไฟล์ import

## รันด้วย Docker

```bash
cp .env.example .env
docker compose up -d --build
```

Frontend: `http://localhost:5173`

Backend health check: `http://localhost:8080/healthz`

## ตั้งค่า Authentik

สร้าง OAuth2/OpenID Provider ใน Authentik แล้วตั้งค่า:

- Redirect URI: `http://localhost:5173/auth/callback`
- Scopes: `openid profile email`
- Signing key: ใช้ค่าเริ่มต้นของ Authentik ได้
- Groups หรือ claim สำหรับบทบาท:
  - กลุ่มที่มีคำว่า `admin` จะได้ role `admin`
  - กลุ่มที่มีคำว่า `team` หรือ `technician` จะได้ role `team`
  - กลุ่มที่มีคำว่า `viewer` จะได้ role `viewer`

ตั้งค่า `.env`:

```env
AUTHENTIK_ISSUER_URL=https://auth.example.com/application/o/airclean-tracker/
AUTHENTIK_CLIENT_ID=...
AUTHENTIK_CLIENT_SECRET=...
AUTHENTIK_REDIRECT_URL=http://localhost:5173/auth/callback
VITE_AUTHENTIK_ISSUER_URL=https://auth.example.com/application/o/airclean-tracker/
VITE_AUTHENTIK_CLIENT_ID=...
VITE_AUTHENTIK_REDIRECT_URL=http://localhost:5173/auth/callback
DEV_AUTH=false
```

สำหรับ local development สามารถตั้ง `DEV_AUTH=true` และเว้นค่า OIDC เพื่อใช้บัญชี admin จำลอง

## Import Excel

ไฟล์ต้นฉบับ: `/Users/admin/Downloads/ล้างแอร์.xlsx`

วิธีผ่านหน้า API:

```bash
curl -X POST http://localhost:8080/api/import/excel \
  -H "Authorization: Bearer <id-token>" \
  -F "file=@/Users/admin/Downloads/ล้างแอร์.xlsx"
```

วิธี auto import ใน Docker:

```bash
mkdir -p docs/import
cp /Users/admin/Downloads/ล้างแอร์.xlsx docs/import/
```

ตั้งค่า `.env`:

```env
AUTO_IMPORT_EXCEL=true
IMPORT_EXCEL_PATH=/app/import/ล้างแอร์.xlsx
```

รายละเอียด mapping อยู่ที่ `docs/excel-import-mapping.md`

## บทบาทผู้ใช้

- `admin` จัดการข้อมูลทั้งหมด import Excel และแก้ role ผู้ใช้
- `team` เพิ่ม/แก้ข้อมูลแอร์ บันทึกการล้าง และวางแผนงาน
- `viewer` ดูข้อมูลอย่างเดียว

## API Overview

- `GET /api/me`
- `GET /api/dashboard`
- `GET /api/aircons`
- `POST /api/aircons`
- `GET /api/aircons/:id`
- `PUT /api/aircons/:id`
- `DELETE /api/aircons/:id`
- `GET /api/aircons/:id/cleaning-records`
- `POST /api/aircons/:id/cleaning-records`
- `PUT /api/cleaning-records/:id`
- `DELETE /api/cleaning-records/:id`
- `GET /api/plans`
- `POST /api/plans/bulk-update`
- `POST /api/import/excel`
- `GET /api/users`
- `PUT /api/users/:id/role`

## Local development

Backend:

```bash
cd backend
go test ./...
go run ./cmd/server
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

## สถานะการล้าง

- ไม่มีวันที่ล้างล่าสุด: `ยังไม่เคยบันทึกล้าง`
- วันที่ครบกำหนดก่อนวันนี้: `เกินกำหนด`
- ครบกำหนดภายใน 30 วัน: `ใกล้ถึงกำหนด`
- มี planned date: `วางแผนแล้ว`
- อื่น ๆ: `ปกติ`

