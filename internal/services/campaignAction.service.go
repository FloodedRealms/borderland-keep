package services

import (
	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/types"
)

type CampaignActionService struct {
	repo repository.Repository
}

func NewCampaignActionService(r repository.Repository) *CampaignActionService {
	return &CampaignActionService{repo: r}
}

func (cas CampaignActionService) AddNewCampaignActionToCharacter(a types.CampaignActivity) error {
	cas.repo.AddCampaignActivityForCharacter(a)
	return nil
}
