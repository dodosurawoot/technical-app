package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"airclean-tracker/backend/internal/auth"
	"airclean-tracker/backend/internal/domain"
	"airclean-tracker/backend/internal/importer"
	"airclean-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type airconInput struct {
	Code                string     `json:"code" binding:"required"`
	Building            string     `json:"building"`
	Floor               string     `json:"floor"`
	Room                string     `json:"room"`
	Brand               string     `json:"brand"`
	BTU                 int        `json:"btu"`
	ResponsibleTeam     string     `json:"responsible_team"`
	LatestCleaningDate  *time.Time `json:"latest_cleaning_date"`
	PlannedCleaningDate *time.Time `json:"planned_cleaning_date"`
	Note                string     `json:"note"`
	ContactName         string     `json:"contact_name"`
	ContactPhone        string     `json:"contact_phone"`
	Subdistrict         string     `json:"subdistrict"`
	District            string     `json:"district"`
	Province            string     `json:"province"`
	Latitude            *float64   `json:"latitude"`
	Longitude           *float64   `json:"longitude"`
}

type cleaningInput struct {
	CleanedDate  *time.Time `json:"cleaned_date"`
	PlannedDate  *time.Time `json:"planned_date"`
	ReportedDate *time.Time `json:"reported_date"`
	Status       string     `json:"status"`
	Note         string     `json:"note"`
	EvidenceURL  string     `json:"evidence_url"`
	PerformedBy  string     `json:"performed_by"`
}

func (s *Server) me(c *gin.Context) {
	c.JSON(http.StatusOK, auth.CurrentUser(c))
}

func (s *Server) dashboard(c *gin.Context) {
	now := time.Now()
	var total, cleanedMonth, dueSoon, overdue, planned int64
	s.db.Model(&models.AirConditioner{}).Count(&total)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := monthStart.AddDate(0, 1, 0)
	s.db.Model(&models.CleaningRecord{}).Where("cleaned_date >= ? AND cleaned_date < ?", monthStart, nextMonth).Count(&cleanedMonth)
	s.db.Model(&models.AirConditioner{}).Where("status = ?", models.StatusDueSoon).Count(&dueSoon)
	s.db.Model(&models.AirConditioner{}).Where("status = ?", models.StatusOverdue).Count(&overdue)
	s.db.Model(&models.AirConditioner{}).Where("planned_cleaning_date IS NOT NULL").Count(&planned)

	var upcoming []models.AirConditioner
	s.db.Where("next_cleaning_date IS NOT NULL").Order("next_cleaning_date ASC").Limit(8).Find(&upcoming)
	c.JSON(http.StatusOK, gin.H{
		"total_air_conditioners": total,
		"cleaned_this_month":     cleanedMonth,
		"due_soon":               dueSoon,
		"overdue":                overdue,
		"planned_jobs":           planned,
		"upcoming":               upcoming,
	})
}

func (s *Server) listAircons(c *gin.Context) {
	var items []models.AirConditioner
	q := s.db.Model(&models.AirConditioner{})
	if search := c.Query("q"); search != "" {
		like := "%" + search + "%"
		q = q.Where("code ILIKE ? OR building ILIKE ? OR room ILIKE ? OR note ILIKE ?", like, like, like, like)
	}
	for _, f := range []struct{ query, col string }{
		{"location", "building"},
		{"room", "room"},
		{"status", "status"},
		{"responsible_team", "responsible_team"},
	} {
		if v := c.Query(f.query); v != "" {
			q = q.Where(fmt.Sprintf("%s = ?", f.col), v)
		}
	}
	if start := parseQueryDate(c.Query("start_date")); start != nil {
		q = q.Where("next_cleaning_date >= ?", *start)
	}
	if end := parseQueryDate(c.Query("end_date")); end != nil {
		q = q.Where("next_cleaning_date <= ?", *end)
	}
	if err := q.Order("next_cleaning_date ASC NULLS LAST, code ASC").Find(&items).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) createAircon(c *gin.Context) {
	var input airconInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ac := input.toModel()
	ac.NextCleaningDate, ac.Status = domain.Recalculate(ac.LatestCleaningDate, ac.PlannedCleaningDate, time.Now())
	if err := s.db.Create(&ac).Error; err != nil {
		errorJSON(c, err)
		return
	}
	audit(s.db, auth.CurrentUser(c), "create", "air_conditioner", ac.ID, ac)
	c.JSON(http.StatusCreated, ac)
}

func (s *Server) getAircon(c *gin.Context) {
	var ac models.AirConditioner
	if err := s.db.Preload("CleaningRecords", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("cleaned_date DESC NULLS LAST, created_at DESC")
	}).First(&ac, c.Param("id")).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, ac)
}

func (s *Server) updateAircon(c *gin.Context) {
	var ac models.AirConditioner
	if err := s.db.First(&ac, c.Param("id")).Error; err != nil {
		errorJSON(c, err)
		return
	}
	var input airconInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	next, status := domain.Recalculate(input.LatestCleaningDate, input.PlannedCleaningDate, time.Now())
	ac.Code = input.Code
	ac.Building = input.Building
	ac.Floor = input.Floor
	ac.Room = input.Room
	ac.Brand = input.Brand
	ac.BTU = input.BTU
	ac.ResponsibleTeam = input.ResponsibleTeam
	ac.LatestCleaningDate = input.LatestCleaningDate
	ac.PlannedCleaningDate = input.PlannedCleaningDate
	ac.Note = input.Note
	ac.ContactName = input.ContactName
	ac.ContactPhone = input.ContactPhone
	ac.Subdistrict = input.Subdistrict
	ac.District = input.District
	ac.Province = input.Province
	ac.Latitude = input.Latitude
	ac.Longitude = input.Longitude
	ac.NextCleaningDate = next
	ac.Status = status
	if err := s.db.Save(&ac).Error; err != nil {
		errorJSON(c, err)
		return
	}
	audit(s.db, auth.CurrentUser(c), "update", "air_conditioner", ac.ID, ac)
	c.JSON(http.StatusOK, ac)
}

func (s *Server) deleteAircon(c *gin.Context) {
	if err := s.db.Delete(&models.AirConditioner{}, c.Param("id")).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) listCleaningRecords(c *gin.Context) {
	var records []models.CleaningRecord
	if err := s.db.Where("air_conditioner_id = ?", c.Param("id")).Order("cleaned_date DESC NULLS LAST, created_at DESC").Find(&records).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": records})
}

func (s *Server) createCleaningRecord(c *gin.Context) {
	airconID, _ := strconv.Atoi(c.Param("id"))
	var input cleaningInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := auth.CurrentUser(c)
	record := models.CleaningRecord{
		AirConditionerID: uint(airconID),
		CleanedDate:      input.CleanedDate,
		PlannedDate:      input.PlannedDate,
		ReportedDate:     input.ReportedDate,
		Status:           input.Status,
		Note:             input.Note,
		EvidenceURL:      input.EvidenceURL,
		PerformedBy:      input.PerformedBy,
	}
	if user != nil {
		record.CreatedByUserID = &user.ID
	}
	if err := s.db.Create(&record).Error; err != nil {
		errorJSON(c, err)
		return
	}
	if err := s.applyRecordToAircon(uint(airconID), input); err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusCreated, record)
}

func (s *Server) updateCleaningRecord(c *gin.Context) {
	var record models.CleaningRecord
	if err := s.db.First(&record, c.Param("id")).Error; err != nil {
		errorJSON(c, err)
		return
	}
	var input cleaningInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record.CleanedDate = input.CleanedDate
	record.PlannedDate = input.PlannedDate
	record.ReportedDate = input.ReportedDate
	record.Status = input.Status
	record.Note = input.Note
	record.EvidenceURL = input.EvidenceURL
	record.PerformedBy = input.PerformedBy
	if err := s.db.Save(&record).Error; err != nil {
		errorJSON(c, err)
		return
	}
	if err := s.applyRecordToAircon(record.AirConditionerID, input); err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, record)
}

func (s *Server) deleteCleaningRecord(c *gin.Context) {
	if err := s.db.Delete(&models.CleaningRecord{}, c.Param("id")).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) listPlans(c *gin.Context) {
	var items []models.AirConditioner
	q := s.db.Where("planned_cleaning_date IS NOT NULL OR status IN ?", []string{models.StatusDueSoon, models.StatusOverdue})
	if start := parseQueryDate(c.Query("start_date")); start != nil {
		q = q.Where("COALESCE(planned_cleaning_date, next_cleaning_date) >= ?", *start)
	}
	if end := parseQueryDate(c.Query("end_date")); end != nil {
		q = q.Where("COALESCE(planned_cleaning_date, next_cleaning_date) <= ?", *end)
	}
	if err := q.Order("COALESCE(planned_cleaning_date, next_cleaning_date) ASC").Find(&items).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) bulkUpdatePlans(c *gin.Context) {
	var input struct {
		IDs         []uint    `json:"ids" binding:"required"`
		PlannedDate time.Time `json:"planned_date" binding:"required"`
		Note        string    `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range input.IDs {
			var ac models.AirConditioner
			if err := tx.First(&ac, id).Error; err != nil {
				return err
			}
			ac.PlannedCleaningDate = &input.PlannedDate
			ac.NextCleaningDate, ac.Status = domain.Recalculate(ac.LatestCleaningDate, ac.PlannedCleaningDate, time.Now())
			if err := tx.Save(&ac).Error; err != nil {
				return err
			}
			plan := models.CleaningPlan{AirConditionerID: id, PlannedDate: input.PlannedDate, Status: models.StatusPlanned, ResponsibleTeam: ac.ResponsibleTeam, Note: input.Note}
			if err := tx.Create(&plan).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": len(input.IDs)})
}

func (s *Server) importExcel(c *gin.Context) {
	service := importer.New(s.db)
	file, err := c.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			errorJSON(c, err)
			return
		}
		defer src.Close()
		result, err := service.ImportReader(src)
		if err != nil {
			errorJSON(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}
	if s.cfg.ImportExcelPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ส่งไฟล์ field=file หรือกำหนด IMPORT_EXCEL_PATH"})
		return
	}
	result, err := service.ImportPath(s.cfg.ImportExcelPath)
	if err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) listUsers(c *gin.Context) {
	var users []models.User
	if err := s.db.Order("email ASC").Find(&users).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": users})
}

func (s *Server) updateUserRole(c *gin.Context) {
	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Role != models.RoleAdmin && input.Role != models.RoleTeam && input.Role != models.RoleViewer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role ไม่ถูกต้อง"})
		return
	}
	if err := s.db.Model(&models.User{}).Where("id = ?", c.Param("id")).Update("role", input.Role).Error; err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"role": input.Role})
}

func (s *Server) applyRecordToAircon(id uint, input cleaningInput) error {
	var ac models.AirConditioner
	if err := s.db.First(&ac, id).Error; err != nil {
		return err
	}
	if input.CleanedDate != nil {
		ac.LatestCleaningDate = input.CleanedDate
	}
	if input.PlannedDate != nil {
		ac.PlannedCleaningDate = input.PlannedDate
	}
	ac.NextCleaningDate, ac.Status = domain.Recalculate(ac.LatestCleaningDate, ac.PlannedCleaningDate, time.Now())
	return s.db.Save(&ac).Error
}

func (i airconInput) toModel() models.AirConditioner {
	return models.AirConditioner{
		Code:                i.Code,
		Building:            i.Building,
		Floor:               i.Floor,
		Room:                i.Room,
		Brand:               i.Brand,
		BTU:                 i.BTU,
		ResponsibleTeam:     i.ResponsibleTeam,
		LatestCleaningDate:  i.LatestCleaningDate,
		PlannedCleaningDate: i.PlannedCleaningDate,
		Note:                i.Note,
		ContactName:         i.ContactName,
		ContactPhone:        i.ContactPhone,
		Subdistrict:         i.Subdistrict,
		District:            i.District,
		Province:            i.Province,
		Latitude:            i.Latitude,
		Longitude:           i.Longitude,
	}
}

func parseQueryDate(raw string) *time.Time {
	if raw == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil
	}
	return &t
}

func errorJSON(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	}
	c.JSON(status, gin.H{"error": err.Error()})
}

func audit(db *gorm.DB, user *models.User, action, entity string, entityID uint, details any) {
	var uid *uint
	if user != nil {
		uid = &user.ID
	}
	raw, _ := json.Marshal(details)
	db.Create(&models.AuditLog{
		UserID:   uid,
		Action:   action,
		Entity:   entity,
		EntityID: strconv.Itoa(int(entityID)),
		Details:  string(raw),
	})
}
