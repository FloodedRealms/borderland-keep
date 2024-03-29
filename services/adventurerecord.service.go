package services

import (
	"context"
	"strconv"

	"github.com/kevin/adventure-archivist/repository"
	"github.com/kevin/adventure-archivist/types"
	"github.com/kevin/adventure-archivist/util"
)

type AdventureRecordService interface {
	CreateAdventureRecordForCampaign(*types.CreateAdventureRecordRequest) (*types.AdventureRecord, error)
	UpdateAdventureRecord(*types.UpdateAdventureRecordRequest) (*types.AdventureRecord, error)
	ListAdventureRecordsForCampaign(string) ([]*types.AdventureRecord, error)
	GetAdventureRecordById(string) (*types.AdventureRecord, error)
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
	return a.repo.GetAdventureRecordsForCampaign(id)

}

func (a *AdventureRecordServiceImpl) GetAdventureRecordById(i string) (*types.AdventureRecord, error) {
	return nil, util.NotYetImplmented()
}
