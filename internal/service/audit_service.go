package service

import "golang-auth-api-mysql/internal/domain"

type AuditService interface {
	Create(audit *domain.AuditLog) error
	GetAll(query domain.AuditQuery) (*domain.AuditListResponse, error)
}
