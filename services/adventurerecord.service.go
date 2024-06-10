package services

import (
	"context"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type AdventureRecordService interface {
	CreateAdventureRecordForCampaign(*types.CreateAdventureRecordRequest) (*types.Adventure, error)
	UpdateAdventureRecord(*types.UpdateAdventureRecordRequest) (*types.Adventure, error)
	ListAdventureRecordsForCampaign(string) ([]*types.Adventure, error)
	GetAdventureRecordById(string) (*types.Adventure, error)
	AddCoinsAdventure(*types.Adventure, *types.Coins) (bool, error)
	AddGemLootToAdventure(*types.Adventure, *types.Gem) (bool, error)
	AddJewelleryLootToAdventure(*types.Adventure, *types.Jewellery) (bool, error)
	AddMagicItemToAdventure(*types.Adventure, *types.MagicItem) (bool, error)
	AddCombatToAdventure(*types.Adventure, *types.MonsterGroup) (bool, error)
}

type AdventureRecordServiceImpl struct {
	repo repository.Repository
	Ctx  context.Context
}

func NewAdventureRecordService(repo repository.Repository, ctx context.Context) *AdventureRecordServiceImpl {
	return &AdventureRecordServiceImpl{repo, ctx}
}

func (a *AdventureRecordServiceImpl) CreateAdventureRecordForCampaign(r *types.CreateAdventureRecordRequest) (*types.Adventure, error) {
	return a.repo.CreateAdventureRecordForCampaign(r)
}

func (a *AdventureRecordServiceImpl) UpdateAdventureRecord(r *types.UpdateAdventureRecordRequest) (*types.Adventure, error) {
	return nil, util.NotYetImplmented()
}
func (a *AdventureRecordServiceImpl) ListAdventureRecordsForCampaign(i string) ([]*types.Adventure, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	campaign := types.NewCampaign(id)
	return a.repo.GetAdventureRecordsForCampaign(campaign)

}

func (a *AdventureRecordServiceImpl) GetAdventureRecordById(i string) (*types.Adventure, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	ad, err2 := a.repo.GetAdventureRecordById(types.NewAdventureRecordById(id))
	if err2 != nil {
		return nil, util.NotYetImplmented()
	}
	return ad, nil
}

func (a *AdventureRecordServiceImpl) AddCoinsAdventure(ad *types.Adventure, c *types.Coins) (bool, error) {
	return a.repo.AddCoinsToAdventure(ad, c)
}

func (a *AdventureRecordServiceImpl) AddGemLootToAdventure(ad *types.Adventure, g *types.Gem) (bool, error) {
	return a.repo.AddGemToAdventure(ad, g)
}
func (a *AdventureRecordServiceImpl) AddJewelleryLootToAdventure(ad *types.Adventure, j *types.Jewellery) (bool, error) {
	return a.repo.AddJewelleryToAdventure(ad, j)
}
func (a *AdventureRecordServiceImpl) AddMagicItemToAdventure(ad *types.Adventure, j *types.MagicItem) (bool, error) {
	return a.repo.AddMagicItemToAdventure(ad, j)
}
func (a *AdventureRecordServiceImpl) AddCombatToAdventure(ad *types.Adventure, j *types.MonsterGroup) (bool, error) {
	return a.repo.AddCombatToAdventure(ad, j)
}
