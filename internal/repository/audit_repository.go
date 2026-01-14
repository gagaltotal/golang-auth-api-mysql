package repository

import "golang-auth-api-mysql/internal/domain"

type AuditRepository interface {
	Create(audit *domain.AuditLog) error
	FindAll(query domain.AuditQuery) ([]domain.AuditLog, int64, error)
}
