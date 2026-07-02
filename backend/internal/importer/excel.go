package importer

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"airclean-tracker/backend/internal/domain"
	"airclean-tracker/backend/internal/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Result struct {
	ImportedAirConditioners int      `json:"imported_air_conditioners"`
	ImportedRecords         int      `json:"imported_records"`
	SkippedRows             int      `json:"skipped_rows"`
	Warnings                []string `json:"warnings"`
}

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ImportPath(path string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, err
	}
	defer f.Close()
	return s.ImportReader(f)
}

func (s *Service) ImportReader(r io.Reader) (Result, error) {
	file, err := excelize.OpenReader(r)
	if err != nil {
		return Result{}, err
	}
	defer file.Close()
	sheet := file.GetSheetName(0)
	rows, err := file.GetRows(sheet)
	if err != nil {
		return Result{}, err
	}
	if len(rows) < 3 {
		return Result{}, fmt.Errorf("ไม่พบข้อมูลใน Excel")
	}
	header := headerMap(rows[1])
	result := Result{}
	for idx, row := range rows[2:] {
		code := get(row, header, "รหัส")
		if code == "" {
			result.SkippedRows++
			continue
		}
		ac := models.AirConditioner{
			Code:            code,
			Building:        get(row, header, "ชื่อหน่วยงาน"),
			Room:            get(row, header, "ชื่อหน่วยงาน"),
			ResponsibleTeam: get(row, header, "ทีมดำเนินการ"),
			ContactName:     get(row, header, "ชื่อผู้ดูแลศูนย์"),
			ContactPhone:    get(row, header, "เบอร์โทร"),
			Subdistrict:     get(row, header, "ตำบล / แขวง"),
			District:        get(row, header, "อำเภอ / เขต"),
			Province:        get(row, header, "จังหวัด"),
			Note:            joinNotes(get(row, header, "รายละเอียดซ่อมเพิ่มเติม"), get(row, header, "หมายเหตุ"), priceNote(get(row, header, "ราคางานล้างแอร์ตามรอบ 6 เดือน"))),
		}
		ac.Latitude = parseFloatPtr(get(row, header, "Lat"))
		ac.Longitude = parseFloatPtr(get(row, header, "Long"))
		ac.PlannedCleaningDate = parseDate(get(row, header, "แผนล้างแอร์"))
		ac.LatestCleaningDate = parseDate(get(row, header, "วันที่ล้าง"))
		ac.NextCleaningDate, ac.Status = domain.Recalculate(ac.LatestCleaningDate, ac.PlannedCleaningDate, time.Now())
		if ac.LatestCleaningDate == nil && strings.Contains(get(row, header, "สถานะ (Status).1"), "เสร็จ") {
			result.Warnings = append(result.Warnings, fmt.Sprintf("แถว %d: สถานะเสร็จแต่ไม่มีวันที่ล้าง", idx+3))
		}
		var existing models.AirConditioner
		err := s.db.Where("code = ?", ac.Code).First(&existing).Error
		if err == nil {
			ac.ID = existing.ID
			ac.CreatedAt = existing.CreatedAt
			if err := s.db.Model(&existing).Updates(ac).Error; err != nil {
				return result, err
			}
		} else if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&ac).Error; err != nil {
				return result, err
			}
			result.ImportedAirConditioners++
		} else {
			return result, err
		}
		if ac.PlannedCleaningDate != nil {
			plan := models.CleaningPlan{
				AirConditionerID: ac.ID,
				PlannedDate:      *ac.PlannedCleaningDate,
				Status:           models.StatusPlanned,
				ResponsibleTeam:  ac.ResponsibleTeam,
				Note:             "นำเข้าจากไฟล์ Excel",
			}
			_ = s.db.Where("air_conditioner_id = ? AND planned_date = ?", ac.ID, plan.PlannedDate).FirstOrCreate(&plan).Error
		}
		if ac.LatestCleaningDate != nil {
			reported := parseDate(get(row, header, "วันที่ส่งรายงาน"))
			record := models.CleaningRecord{
				AirConditionerID: ac.ID,
				CleanedDate:      ac.LatestCleaningDate,
				PlannedDate:      ac.PlannedCleaningDate,
				ReportedDate:     reported,
				Status:           get(row, header, "สถานะ (Status).1"),
				PerformedBy:      ac.ResponsibleTeam,
				Note:             "นำเข้าจากไฟล์ Excel",
			}
			var existingRecord models.CleaningRecord
			err := s.db.Where("air_conditioner_id = ? AND cleaned_date = ?", ac.ID, *ac.LatestCleaningDate).First(&existingRecord).Error
			if err == gorm.ErrRecordNotFound {
				if err := s.db.Create(&record).Error; err != nil {
					return result, err
				}
				result.ImportedRecords++
			}
		}
	}
	return result, nil
}

func headerMap(headers []string) map[string]int {
	out := map[string]int{}
	counts := map[string]int{}
	for i, h := range headers {
		key := strings.TrimSpace(h)
		if key == "" {
			continue
		}
		if count := counts[key]; count > 0 {
			counts[key] = count + 1
			key = fmt.Sprintf("%s.%d", key, count)
		} else {
			counts[key] = 1
		}
		out[key] = i
	}
	return out
}

func get(row []string, header map[string]int, name string) string {
	i, ok := header[name]
	if !ok || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}

func parseFloatPtr(raw string) *float64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(v) {
		return nil
	}
	return &v
}

func parseDate(raw string) *time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if serial, err := strconv.ParseFloat(raw, 64); err == nil && serial > 1000 {
		if t, err := excelize.ExcelDateToTime(serial, false); err == nil {
			return datePtr(t)
		}
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2/1/2006",
		"02/01/2006",
		"2/1/06",
		"02/01/06",
		"2/12006",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
			return datePtr(t)
		}
	}
	if strings.Contains(raw, "-") && strings.Contains(raw, "/") {
		parts := strings.Split(raw, "/")
		if len(parts) >= 2 {
			dayPart := strings.Split(parts[0], "-")[0]
			month := parts[len(parts)-2]
			year := parts[len(parts)-1]
			if len(parts) == 2 {
				month = "6"
				year = parts[1]
			}
			if len(year) == 4 {
				normalized := fmt.Sprintf("%s/%s/%s", strings.TrimSpace(dayPart), strings.TrimSpace(month), strings.TrimSpace(year))
				if t, err := time.ParseInLocation("2/1/2006", normalized, time.Local); err == nil {
					return datePtr(t)
				}
			}
		}
	}
	return nil
}

func datePtr(t time.Time) *time.Time {
	y, m, d := t.Date()
	v := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	return &v
}

func joinNotes(parts ...string) string {
	out := []string{}
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			out = append(out, strings.TrimSpace(part))
		}
	}
	return strings.Join(out, "\n")
}

func priceNote(raw string) string {
	if raw == "" {
		return ""
	}
	return "ราคางานล้างแอร์ตามรอบ 6 เดือน: " + raw
}
