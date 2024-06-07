package services

import (
	"context"
	"log"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type AdventureRecordService interface {
	CreateAdventureRecordForCampaign(*types.CreateAdventureRecordRequest) (*types.AdventureRecord, error)
	UpdateAdventureRecord(*types.UpdateAdventureRecordRequest) (*types.AdventureRecord, error)
	ListAdventureRecordsForCampaign(string) ([]*types.AdventureRecord, error)
	GetAdventureRecordById(string) (*types.AdventureRecord, error)
	AddGemLootToAdventure(*types.AdventureRecord, *types.Gem) (bool, error)
}

type AdventureRecordServiceImpl struct {
	repo repository.Repository
	Ctx  context.Context
}

func NewAdventureRecordService(repo repository.Repository, ctx context.Context) *AdventureRecordServiceImpl {
	return &AdventureRecordServiceImpl{repo, ctx}
}

func (a *AdventureRecordServiceImpl) CreateAdventureRecordForCampaign(r *types.CreateAdventureRecordRequest) (*types.AdventureRecord, error) {
	return a.repo.CreateAdventureRecordForCampaign(r)
}

func (a *AdventureRecordServiceImpl) UpdateAdventureRecord(r *types.UpdateAdventureRecordRequest) (*types.AdventureRecord, error) {
	return nil, util.NotYetImplmented()
}
func (a *AdventureRecordServiceImpl) ListAdventureRecordsForCampaign(i string) ([]*types.AdventureRecord, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	campaign := types.NewCampaign(id)
	return a.repo.GetAdventureRecordsForCampaign(campaign)

}

func (a *AdventureRecordServiceImpl) GetAdventureRecordById(i string) (*types.AdventureRecord, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	ad, err2 := a.repo.GetAdventureRecordById(types.NewAdventureRecordById(id))
	if err2 != nil {
		return nil, util.NotYetImplmented()
	}
	return ad, nil
}

func (a *AdventureRecordServiceImpl) AddGemLootToAdventure(ad *types.AdventureRecord, g *types.Gem) (bool, error) {
	log.Print("hit")
	return a.repo.AddGemToAdventure(ad, g)
}
