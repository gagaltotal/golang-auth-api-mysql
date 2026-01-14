package repository

import (
	"golang-auth-api-mysql/internal/domain"

	"gorm.io/gorm"
)

type auditRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepositoryImpl{db: db}
}

func (r *auditRepositoryImpl) Create(audit *domain.AuditLog) error {
	return r.db.Create(audit).Error
}

func (r *auditRepositoryImpl) FindAll(query domain.AuditQuery) ([]domain.AuditLog, int64, error) {
	var audits []domain.AuditLog
	var total int64

	offset := (query.Page - 1) * query.Limit
	db := r.db.Model(&domain.AuditLog{}).Preload("User")

	// Apply filters
	if query.UserID > 0 {
		db = db.Where("user_id = ?", query.UserID)
	}
	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}
	if query.Resource != "" {
		db = db.Where("resource = ?", query.Resource)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Order("created_at DESC").Offset(offset).Limit(query.Limit).Find(&audits).Error
	return audits, total, err
}
