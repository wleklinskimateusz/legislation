package application

import (
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/act"
)

// CreateActService creates new draft Acts.
type CreateActService struct {
	ActRepo act.ActRepository
}

// CreateAct creates a new Act in Draft status and persists it.
func (s *CreateActService) CreateAct(actID string) error {
	a := act.NewDraftAct(actID)
	return s.ActRepo.Save(a)
}
