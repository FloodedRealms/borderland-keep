package services

import (
	"context"
	"strconv"
	"time"

	"github.com/kevin/adventure-archivist/repository"
	"github.com/kevin/adventure-archivist/types"
	"github.com/kevin/adventure-archivist/util"
)

type CampaignService interface {
	CreateCampaign(*types.CreateCampaignRequest) (*types.Campaign, error)
	UpdateCampaign(*types.UpdateCampaignRequest) (*types.Campaign, error)
	GetCampaign(string) (*types.Campaign, error)
	ListCampaigns() ([]*types.Campaign, error)
	DeleteCampaign(string) (bool, error)
}

type CampaignServiceImpl struct {
	repo   repository.Repository
	logger util.Logger
	ctx    context.Context
}

func NewCampaignService(repo repository.Repository, logger *util.Logger, ctx context.Context) *CampaignServiceImpl {
	return &CampaignServiceImpl{repo, *logger, ctx}
}

func (c *CampaignServiceImpl) CreateCampaign(cr *types.CreateCampaignRequest) (*types.Campaign, error) {
	cr.CreatedAt = time.Now()
	cr.UpdatedAt = cr.CreatedAt
	cr.LastAdventure = cr.CreatedAt

	return c.repo.CreateCampaign(cr)
}

func (c *CampaignServiceImpl) UpdateCampaign(ur *types.UpdateCampaignRequest) (*types.Campaign, error) {
	return nil, util.NotYetImplmented()
}

func (c *CampaignServiceImpl) GetCampaign(id string) (*types.Campaign, error) {
	campaignId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	campaign, err := c.repo.GetCampaign(campaignId)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}

func (c *CampaignServiceImpl) ListCampaigns() ([]*types.Campaign, error) {
	return c.repo.ListCampaigns()
}

func (c *CampaignServiceImpl) DeleteCampaign(id string) (bool, error) {
	campaignId, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}
	campaignToDelete := types.NewCampaign(campaignId)
	return c.repo.DeleteCampaign(campaignToDelete)
}
