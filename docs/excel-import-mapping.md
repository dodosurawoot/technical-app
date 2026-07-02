# Excel Import Mapping

ไฟล์ที่ตรวจสอบ: `/Users/admin/Downloads/ล้างแอร์.xlsx`

Sheet: `Progress ล้างแอร์_รอบที่ 4`

ใช้ row 2 เป็น header และเริ่มอ่านข้อมูลตั้งแต่ row 3 มีข้อมูลรหัสแอร์ 386 แถว

| Excel column | Database field |
| --- | --- |
| `รหัส` | `air_conditioners.code` |
| `ชื่อหน่วยงาน` | `building`, `room` |
| `ทีมดำเนินการ` | `responsible_team` |
| `ตำบล / แขวง` | `subdistrict` |
| `อำเภอ / เขต` | `district` |
| `จังหวัด` | `province` |
| `Lat`, `Long` | `latitude`, `longitude` |
| `ชื่อผู้ดูแลศูนย์`, `เบอร์โทร` | `contact_name`, `contact_phone` |
| `แผนล้างแอร์` | `planned_cleaning_date`, `cleaning_plans.planned_date` |
| `วันที่ล้าง` | `latest_cleaning_date`, `cleaning_records.cleaned_date` |
| `วันที่ส่งรายงาน` | `cleaning_records.reported_date` |
| `สถานะ (Status).1` | `cleaning_records.status` |
| `รายละเอียดซ่อมเพิ่มเติม`, `ราคางานล้างแอร์ตามรอบ 6 เดือน`, `หมายเหตุ` | `note` |

Assumptions:

- Excel ไม่มี column อาคาร/ชั้น/ห้องแยก จึงใช้ `ชื่อหน่วยงาน` เป็นทั้ง `building` และ `room`.
- Excel ไม่มี brand และ BTU จึงเว้นว่างไว้เพื่อให้ทีมเติมภายหลัง.
- วันที่แบบช่วง เช่น `22-26/6/2026` จะนำเข้าวันแรกของช่วงเป็น planned date.
- วันที่ล้างล่าสุด + 6 เดือน ถูกใช้เป็นกำหนดล้างครั้งถัดไป.
- ไฟล์ต้นฉบับไม่ถูกแก้ไข การ import อ่านอย่างเดียว.

