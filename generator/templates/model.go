package templates

func GormItemModelTemplate() string {
	return `package model

import "time"

// Item represents the persisted item model.
type Item struct {
	ID          string         ` + "`" + `json:"id" gorm:"primaryKey;size:64"` + "`" + `
	Name        string         ` + "`" + `json:"name" gorm:"size:255;not null;index"` + "`" + `
	Code        string         ` + "`" + `json:"code" gorm:"size:128;index"` + "`" + `
	Type        string         ` + "`" + `json:"type" gorm:"size:64;index"` + "`" + `
	Status      string         ` + "`" + `json:"status" gorm:"size:64;index"` + "`" + `
	Description string         ` + "`" + `json:"description" gorm:"column:description;type:text"` + "`" + `
	CreatedBy   string         ` + "`" + `json:"created_by" gorm:"size:64;index"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `json:"created_at" gorm:"not null;index"` + "`" + `
	UpdatedBy   string         ` + "`" + `json:"updated_by" gorm:"size:64"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `json:"updated_at" gorm:"not null"` + "`" + `
	DeletedAt   *time.Time     ` + "`" + `json:"deleted_at,omitempty" gorm:"index"` + "`" + `
	Extras      map[string]any ` + "`" + `json:"extras,omitempty" gorm:"serializer:json"` + "`" + `
}

// TableName returns the database table name.
func (Item) TableName() string {
	return "items"
}
`
}
