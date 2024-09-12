package services

import (
	"github.com/floodedrealms/borderland-keep/internal/repository"
)

type CampaignActionService struct {
	repo repository.Repository
}

func NewCampaignActionService(r repository.Repository) *CampaignActionService {
	return &CampaignActionService{repo: r}
}
