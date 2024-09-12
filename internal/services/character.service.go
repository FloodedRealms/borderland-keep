package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/types"
)

type CharacterService struct {
	repo   repository.Repository
	logger *util.Logger
	ctx    context.Context
}

const characterTable = "characters"
const characterCampaignSummaryView = "character_campaign_view"

func NewCharacterService(repo repository.Repository, logger *util.Logger, ctx context.Context) *CharacterService {
	return &CharacterService{
		repo:   repo,
		logger: logger,
		ctx:    ctx,
	}
}

func (s CharacterService) CreateCharactersForCampaign(campaignId int, charactersToCreate []types.CharacterRecord) error {
	queries := make([]string, 0)
	paramList := make([][]interface{}, 0)
	for _, character := range charactersToCreate {
		queries = append(queries, fmt.Sprintf("INSERT INTO %s(campaign_id, name, status_id, prime_req_percent, class_id, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?) ;", characterTable))
		t := time.Now()
		params := []interface{}{campaignId, character.Name, 1, character.PrimeReqPercent, character.ClassId, t, t}
		paramList = append(paramList, params)

	}
	s.repo.DoTransaction(queries, paramList)
	return nil
}

func (s CharacterService) UpdateCharacters(charactersToUpdate []types.CharacterRecord) error {
	q := make([]string, 0)
	p := make([][]interface{}, 0)
	for _, c := range charactersToUpdate {
		stmtstring := fmt.Sprintf("UPDATE %s set status_id=? WHERE id=?;", characterTable)
		params := []interface{}{c.StatusId, c.Id}
		q = append(q, stmtstring)
		p = append(p, params)
	}
	return s.repo.DoTransaction(q, p)
}

func (s CharacterService) DeleteCharacters(charactersToDelete []string) error {
	stmtString := fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", characterTable, strings.Join(charactersToDelete, ", "))
	_, err := s.repo.ExecuteQuery(stmtString)
	return err
}

func (s CharacterService) GetCharactersForCampaign(campaignId int) (characters []types.CharacterRecord, err error) {
	stmtstr := fmt.Sprintf("SELECT t.id, t.name, t.status_id, t.class_id, t.prime_req_percent FROM %s t where t.campaign_id = ?;", characterTable)
	rows, err := s.repo.RunQuery(stmtstr, campaignId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := types.CharacterRecord{}
		err = rows.Scan(&c.Id, &c.Name, &c.StatusId, &c.ClassId, &c.PrimeReqPercent)
		if err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}
	return characters, nil
}

func (s CharacterService) GetCharacterCampaignSummary(campaignId int) (characters []types.CharacterRecord, err error) {
	stmtstr := fmt.Sprintf("SELECT t.name, t.status, t.class_name, t.preq, t.total_xp, t.level FROM %s t where t.campaign_id = ?;", characterCampaignSummaryView)
	rows, err := s.repo.RunQuery(stmtstr, campaignId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := types.CharacterRecord{}
		err = rows.Scan(&c.Name, &c.Status, &c.Class, &c.PrimeReqPercent, &c.CurrentXP, &c.Level)
		if err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}
	return characters, nil
}

func (s CharacterService) GetPossibleCharactersForAdventure(aId int) ([]types.CharacterRecord, error) {
	chars := make([]types.CharacterRecord, 0)
	stmtString := ("SELECT c.id, c.name from adventures a " +
		"JOIN \"characters\" c ON c.campaign_id = a.campaign_id " +
		"WHERE a.id = ?;")

	result, err := s.repo.RunQuery(stmtString, aId)
	if err != nil {
		return chars, err
	}
	for result.Next() {
		c := types.CharacterRecord{}
		result.Scan(&c.Id, &c.Name)
		chars = append(chars, c)
	}
	return chars, nil
}

func (s CharacterService) sumXP(a []int) int {
	xp := 0
	for _, x := range a {
		xp += x
	}
	return xp
}
