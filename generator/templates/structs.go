package templates

func StructsTemplate() string {
	return `package structs

import (
	"time"
)

// Item represents a generic item entity.
type Item struct {
	ID        string     ` + "`" + `json:"id"` + "`" + `
	Name      string     ` + "`" + `json:"name"` + "`" + `
	Code      string     ` + "`" + `json:"code"` + "`" + `
	Type      string     ` + "`" + `json:"type"` + "`" + `
	Status    string     ` + "`" + `json:"status"` + "`" + `
	Desc      string     ` + "`" + `json:"description"` + "`" + `
	CreatedBy string     ` + "`" + `json:"created_by"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
	UpdatedBy string     ` + "`" + `json:"updated_by"` + "`" + `
	UpdatedAt time.Time  ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt *time.Time ` + "`" + `json:"deleted_at,omitempty"` + "`" + `
	Extras    map[string]any ` + "`" + `json:"extras,omitempty"` + "`" + `
}

// CreateItemRequest represents the request to create a new item.
type CreateItemRequest struct {
	Name      string         ` + "`" + `json:"name" binding:"required"` + "`" + `
	Code      string         ` + "`" + `json:"code"` + "`" + `
	Type      string         ` + "`" + `json:"type"` + "`" + `
	Status    string         ` + "`" + `json:"status"` + "`" + `
	Desc      string         ` + "`" + `json:"description"` + "`" + `
	Extras    map[string]any ` + "`" + `json:"extras"` + "`" + `
}

// UpdateItemRequest represents the request to update an existing item.
type UpdateItemRequest struct {
	ID        string         ` + "`" + `json:"id" binding:"required"` + "`" + `
	Name      string         ` + "`" + `json:"name"` + "`" + `
	Code      string         ` + "`" + `json:"code"` + "`" + `
	Type      string         ` + "`" + `json:"type"` + "`" + `
	Status    string         ` + "`" + `json:"status"` + "`" + `
	Desc      string         ` + "`" + `json:"description"` + "`" + `
	Extras    map[string]any ` + "`" + `json:"extras"` + "`" + `
}

// ItemResponse represents the response containing item details.
type ItemResponse struct {
	ID        string         ` + "`" + `json:"id"` + "`" + `
	Name      string         ` + "`" + `json:"name"` + "`" + `
	Code      string         ` + "`" + `json:"code"` + "`" + `
	Type      string         ` + "`" + `json:"type"` + "`" + `
	Status    string         ` + "`" + `json:"status"` + "`" + `
	Desc      string         ` + "`" + `json:"description"` + "`" + `
	CreatedBy string         ` + "`" + `json:"created_by"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedBy string         ` + "`" + `json:"updated_by"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	Extras    map[string]any ` + "`" + `json:"extras,omitempty"` + "`" + `
}

// ListItemsRequest represents the request to list items with pagination and filtering.
type ListItemsRequest struct {
	PageSize int    ` + "`" + `json:"page_size" form:"page_size"` + "`" + `
	PageNum  int    ` + "`" + `json:"page_num" form:"page_num"` + "`" + `
	Keyword  string ` + "`" + `json:"keyword" form:"keyword"` + "`" + `
	Type     string ` + "`" + `json:"type" form:"type"` + "`" + `
	Status   string ` + "`" + `json:"status" form:"status"` + "`" + `
}
`
}
