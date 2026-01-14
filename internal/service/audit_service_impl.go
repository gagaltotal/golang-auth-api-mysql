package service

import (
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/repository"
	"math"
)

type auditServiceImpl struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) AuditService {
	return &auditServiceImpl{auditRepo: auditRepo}
}

func (s *auditServiceImpl) Create(audit *domain.AuditLog) error {
	return s.auditRepo.Create(audit)
}

func (s *auditServiceImpl) GetAll(query domain.AuditQuery) (*domain.AuditListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}

	audits, total, err := s.auditRepo.FindAll(query)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &domain.AuditListResponse{
		Data:       audits,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}
