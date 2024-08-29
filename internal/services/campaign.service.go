package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/floodedrealms/adventure-archivist/internal/repository"
	"github.com/floodedrealms/adventure-archivist/internal/util"
	"github.com/floodedrealms/adventure-archivist/types"
)

type CampaignService struct {
	repo   repository.Repository
	logger util.Logger
	ctx    context.Context
}

const campaignTable = "campaigns"
const PAGE_SIZE = 10

func NewCampaignService(repo repository.Repository, logger *util.Logger, ctx context.Context) *CampaignService {
	return &CampaignService{repo, *logger, ctx}
}

func (c *CampaignService) CreateCampaign(cr types.CampaignRecord, pass types.Password) (*types.CampaignRecord, error) {
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

func (c *CampaignService) UpdateCampaign(ur *types.CampaignRecord) (*types.CampaignRecord, error) {
	ur.UpdatedAt = time.Now()

	return c.repo.UpdateCampaign(ur)
}

func (c *CampaignService) GetCampaign(id string) (*types.CampaignRecord, error) {
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

func (c *CampaignService) ListCampaigns() ([]*types.CampaignRecord, error) {
	return c.repo.ListCampaigns()
}

func (c *CampaignService) ListCampaignsForClient(clientId string) ([]*types.CampaignRecord, error) {
	return c.repo.ListCampaignsForClient(clientId)
}

func (c *CampaignService) DeleteCampaign(id string) (bool, error) {
	campaignId, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}
	campaignToDelete := types.NewCampaign(campaignId)
	return c.repo.DeleteCampaign(campaignToDelete)
}

func (c *CampaignService) UpdateCampaignPassword(id, password string) (string, error) {
	campaignId, err := strconv.Atoi(id)
	if err != nil {
		return "Password update failed", err
	}
	hashedPassword, _ := types.NewPassword(password)
	return password, c.repo.UpdateCampaignPassword(campaignId, *hashedPassword)

}

func (c *CampaignService) TenMostRecentlyActiveCampaigns(page int) []types.CampaignRecord {
	pageToOffeset := (page - 1) * PAGE_SIZE
	stmtStr := fmt.Sprintf("SELECT id, name, recruitment, judge, timekeeping, cadence, last_adventure FROM %s ORDER BY last_adventure DESC LIMIT %d OFFSET %d ;", campaignTable, PAGE_SIZE, pageToOffeset)
	rows, err := c.repo.RunQuery(stmtStr)
	if err != nil {
		return nil
	}
	results := make([]types.CampaignRecord, 0)
	defer rows.Close()
	for rows.Next() {
		cur := types.CampaignRecord{}
		rows.Scan(&cur.Id, &cur.Name, &cur.Recruitment, &cur.Judge, &cur.Timekeeping, &cur.Cadence, &cur.LastAdventure)
		results = append(results, cur)
	}
	return results
}

func (c CampaignService) GetClassOptionsForCampaign(id int) ([]types.CampaignClassOption, error) {
	stmtStr := fmt.Sprintf("SELECT cl.class_id, cl.class_name FROM %s cl WHERE campaign_id = %d", "campaign_to_class_options", id)
	rows, err := c.repo.RunQuery(stmtStr)
	if err != nil {
		return nil, err
	}
	results := make([]types.CampaignClassOption, 0)
	defer rows.Close()
	for rows.Next() {
		cur := types.CampaignClassOption{}
		rows.Scan(&cur.ClassId, &cur.ClassName)
		results = append(results, cur)
	}
	return results, err
}
