package services

import (
	"context"
	"fmt"
	"time"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/types"
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

func (c *CampaignService) UpdateCampaign(ur *types.CampaignRecord) (*types.CampaignRecord, error) {
	ur.UpdatedAt = time.Now()

	return c.repo.UpdateCampaign(ur)
}

func (c *CampaignService) GetCampaign(id int) (*types.CampaignRecord, error) {
	tableq := fmt.Sprintf("SELECT c.name, c.judge, c.timekeeping, c.recruitment FROM %s c where c.id =?;", campaignTable)
	//tableq1 := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", campaignTable)
	rows, err := c.repo.RunQuery(tableq, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		campaignRows []*types.CampaignRecord
	)
	for rows.Next() {
		var current types.CampaignRecord
		err := rows.Scan(&current.Name, &current.Judge, &current.Timekeeping, &current.Recruitment)
		if err != nil {
			return nil, err
		}
		campaignRows = append(campaignRows, &current)
	}
	if len(campaignRows) < 1 {
		return nil, fmt.Errorf("campaign %d not found", id)
	}
	return campaignRows[0], nil
}

func (c *CampaignService) CampaignSummary(id int) (*types.CampaignRecord, error) {
	tableq := fmt.Sprintf("SELECT c.id, c.name, c.recruitment, c.judge, c.timekeeping, c.cadence, c.last_adventure FROM %s c where c.id =?;", campaignTable)
	//tableq1 := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", campaignTable)
	rows, err := c.repo.RunQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()
	var (
		campaignRows []*types.CampaignRecord
	)
	for rows.Next() {
		var current types.CampaignRecord
		err := rows.Scan(&current.Id, &current.Name, &current.Recruitment, &current.Judge, &current.Timekeeping, &current.Cadence, &current.LastAdventure)
		util.CheckErr(err)
		util.CheckErr(err)
		campaignRows = append(campaignRows, &current)
	}
	result := campaignRows[0]
	return result, nil

}

func (c *CampaignService) CampaignAdventuresSummary(id int) ([]types.AdventureRecord, error) {
	tableq := fmt.Sprintf("SELECT c.id, c.name, c.adventure_date FROM %s c where c.campaign_id =?;", adventureTable)
	//tableq1 := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", campaignTable)
	rows, err := c.repo.RunQuery(tableq, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []types.AdventureRecord
	for rows.Next() {
		var current types.AdventureRecord
		err := rows.Scan(&current.Id, &current.Name, &current.AdventureDate)
		if err != nil {
			return nil, err
		}
		results = append(results, current)
	}
	return results, nil

}

func (c *CampaignService) ListCampaigns() ([]*types.CampaignRecord, error) {
	return c.repo.ListCampaigns()
}

func (c *CampaignService) ListCampaignsForClient(clientId string) ([]*types.CampaignRecord, error) {
	return c.repo.ListCampaignsForClient(clientId)
}

func (c *CampaignService) DeleteCampaign(id string) (bool, error) {
	idList := []interface{}{id}
	deleteAtcRecords := fmt.Sprintf("DELETE FROM %s WHERE adventure_id IN (SELECT a.id FROM adventures a WHERE a.campaign_id = ?);", adventureToCharactersTable)
	deleteAdventureRecords := fmt.Sprintf("DELETE FROM %s WHERE campaign_id = ?;", adventureTable)
	deleteCharacters := fmt.Sprintf("DELETE FROM %s WHERE campaign_id = ?;", adventureTable)
	deleteCampaigns := fmt.Sprintf("DELETE FROM %s WHERE id = ?;", campaignTable)
	statements := []string{deleteAtcRecords, deleteAdventureRecords, deleteCharacters, deleteCampaigns}
	params := [][]interface{}{idList, idList, idList, idList}

	err := c.repo.DoTransaction(statements, params)
	if err != nil {
		return false, err
	}
	return true, nil
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

func (c *CampaignService) CampaignsForUser(userId string) []types.CampaignRecord {
	stmtStr := fmt.Sprintf("SELECT id, name, recruitment, judge, timekeeping, cadence, last_adventure FROM %s WHERE user_id=?;", campaignTable)
	rows, err := c.repo.RunQuery(stmtStr, userId)
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

func (c *CampaignService) CreateCampaignForUser(userId string) (*types.CampaignRecord, error) {
	stmt := fmt.Sprintf("INSERT INTO %s(user_id, name, created_at, updated_at, last_adventure) values(?, ?, ?, ?, ?)", campaignTable)
	time := time.Now()
	results, err := c.repo.ExecuteQuery(stmt, userId, "New Campaign", time, time, time)
	if err != nil {
		return nil, err
	}
	id, _ := results.LastInsertId()
	return c.CampaignSummary(int(id))

}

func (c *CampaignService) UpdateCampaignDetails(id int, name, judge, timekeeping string, isRecruiting bool) error {
	stms := fmt.Sprintf("UPDATE %s set name=?, judge=?, timekeeping=?, recruitment=?, updated_at=? WHERE id=?;", campaignTable)
	time := time.Now()
	_, err := c.repo.ExecuteQuery(stms, name, judge, timekeeping, isRecruiting, time, id)
	return err
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
