# Deployment Notes

วางไฟล์ `.env` ที่ root ของโปรเจกต์ แล้วรัน:

```bash
docker compose up -d --build
```

หากต้องการ import Excel ผ่าน path ใน container ให้วางไฟล์ไว้ที่ `docs/import/ล้างแอร์.xlsx` และตั้งค่า:

```env
AUTO_IMPORT_EXCEL=true
IMPORT_EXCEL_PATH=/app/import/ล้างแอร์.xlsx
```

