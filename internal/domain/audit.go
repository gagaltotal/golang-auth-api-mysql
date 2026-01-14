package domain

import "time"

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     string    `gorm:"size:50;not null" json:"action"`
	Resource   string    `gorm:"size:50" json:"resource"`
	ResourceID uint      `json:"resource_id"`
	Method     string    `gorm:"size:10" json:"method"`
	Path       string    `gorm:"size:255" json:"path"`
	IPAddress  string    `gorm:"size:50" json:"ip_address"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	Payload    string    `gorm:"type:text" json:"payload"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuditQuery struct {
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=20"`
	UserID   uint   `form:"user_id"`
	Action   string `form:"action"`
	Resource string `form:"resource"`
}

type AuditListResponse struct {
	Data       []AuditLog `json:"data"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalItems int64      `json:"total_items"`
	TotalPages int        `json:"total_pages"`
}
