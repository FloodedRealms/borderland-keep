package services

import (
	"context"
	"strconv"
	"time"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type CampaignService interface {
	CreateCampaign(types.CampaignRecord, types.Password) (*types.CampaignRecord, error)
	UpdateCampaign(*types.CampaignRecord) (*types.CampaignRecord, error)
	GetCampaign(string) (*types.CampaignRecord, error)
	ListCampaigns() ([]*types.CampaignRecord, error)
	ListCampaignsForClient(string) ([]*types.CampaignRecord, error)
	DeleteCampaign(string) (bool, error)
	UpdateCampaignPassword(string, string) (string, error)
}

type CampaignServiceImpl struct {
	repo   repository.Repository
	logger util.Logger
	ctx    context.Context
}

func NewCampaignService(repo repository.Repository, logger *util.Logger, ctx context.Context) *CampaignServiceImpl {
	return &CampaignServiceImpl{repo, *logger, ctx}
}

func (c *CampaignServiceImpl) CreateCampaign(cr types.CampaignRecord, pass types.Password) (*types.CampaignRecord, error) {
	ca, err := c.repo.CreateCampaign(&cr)
	if err != nil {
		return nil, err
	}
	err = c.repo.UpdateCampaignPassword(ca.Id, pass)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (c *CampaignServiceImpl) UpdateCampaign(ur *types.CampaignRecord) (*types.CampaignRecord, error) {
	ur.UpdatedAt = time.Now()

	return c.repo.UpdateCampaign(ur)
}

func (c *CampaignServiceImpl) GetCampaign(id string) (*types.CampaignRecord, error) {
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

func (c *CampaignServiceImpl) ListCampaigns() ([]*types.CampaignRecord, error) {
	return c.repo.ListCampaigns()
}

func (c *CampaignServiceImpl) ListCampaignsForClient(clientId string) ([]*types.CampaignRecord, error) {
	return c.repo.ListCampaignsForClient(clientId)
}

func (c *CampaignServiceImpl) DeleteCampaign(id string) (bool, error) {
	campaignId, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}
	campaignToDelete := types.NewCampaign(campaignId)
	return c.repo.DeleteCampaign(campaignToDelete)
}

func (c *CampaignServiceImpl) UpdateCampaignPassword(id, password string) (string, error) {
	campaignId, err := strconv.Atoi(id)
	if err != nil {
		return "Password update failed", err
	}
	hashedPassword, _ := types.NewPassword(password)
	return password, c.repo.UpdateCampaignPassword(campaignId, *hashedPassword)

}
