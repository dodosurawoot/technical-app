package models

import "time"

const (
	RoleAdmin  = "admin"
	RoleTeam   = "team"
	RoleViewer = "viewer"

	StatusNeverCleaned = "never_cleaned"
	StatusNormal       = "normal"
	StatusDueSoon      = "due_soon"
	StatusOverdue      = "overdue"
	StatusPlanned      = "planned"
)

type User struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ProviderSubject string    `gorm:"index" json:"provider_subject"`
	Email           string    `gorm:"uniqueIndex;size:320;not null" json:"email"`
	Name            string    `gorm:"size:255" json:"name"`
	Username        string    `gorm:"size:120" json:"username"`
	Role            string    `gorm:"size:40;default:viewer;index" json:"role"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AirConditioner struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	Code                string     `gorm:"uniqueIndex;size:80;not null" json:"code"`
	Building            string     `gorm:"size:255;index" json:"building"`
	Floor               string     `gorm:"size:80" json:"floor"`
	Room                string     `gorm:"size:255;index" json:"room"`
	Brand               string     `gorm:"size:120" json:"brand"`
	BTU                 int        `json:"btu"`
	ResponsibleTeam     string     `gorm:"size:160;index" json:"responsible_team"`
	LatestCleaningDate  *time.Time `gorm:"index" json:"latest_cleaning_date"`
	NextCleaningDate    *time.Time `gorm:"index" json:"next_cleaning_date"`
	PlannedCleaningDate *time.Time `gorm:"index" json:"planned_cleaning_date"`
	Status              string     `gorm:"size:40;index" json:"status"`
	Note                string     `gorm:"type:text" json:"note"`
	ContactName         string     `gorm:"size:255" json:"contact_name"`
	ContactPhone        string     `gorm:"size:80" json:"contact_phone"`
	Subdistrict         string     `gorm:"size:160" json:"subdistrict"`
	District            string     `gorm:"size:160" json:"district"`
	Province            string     `gorm:"size:160" json:"province"`
	Latitude            *float64   `json:"latitude"`
	Longitude           *float64   `json:"longitude"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CleaningRecords     []CleaningRecord
	Plans               []CleaningPlan
}

type CleaningRecord struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	AirConditionerID uint       `gorm:"index;not null" json:"air_conditioner_id"`
	CleanedDate      *time.Time `gorm:"index" json:"cleaned_date"`
	PlannedDate      *time.Time `gorm:"index" json:"planned_date"`
	ReportedDate     *time.Time `json:"reported_date"`
	Status           string     `gorm:"size:80" json:"status"`
	Note             string     `gorm:"type:text" json:"note"`
	EvidenceURL      string     `gorm:"size:500" json:"evidence_url"`
	PerformedBy      string     `gorm:"size:160" json:"performed_by"`
	CreatedByUserID  *uint      `json:"created_by_user_id"`
	CreatedByUser    *User
	AirConditioner   AirConditioner
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CleaningPlan struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	AirConditionerID uint      `gorm:"index;not null" json:"air_conditioner_id"`
	PlannedDate      time.Time `gorm:"index;not null" json:"planned_date"`
	Status           string    `gorm:"size:40;default:planned;index" json:"status"`
	Note             string    `gorm:"type:text" json:"note"`
	ResponsibleTeam  string    `gorm:"size:160;index" json:"responsible_team"`
	CreatedByUserID  *uint     `json:"created_by_user_id"`
	CreatedByUser    *User
	AirConditioner   AirConditioner
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	Action    string    `gorm:"size:120;index" json:"action"`
	Entity    string    `gorm:"size:120;index" json:"entity"`
	EntityID  string    `gorm:"size:120;index" json:"entity_id"`
	Details   string    `gorm:"type:text" json:"details"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
