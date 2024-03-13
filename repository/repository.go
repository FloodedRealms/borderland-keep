package repository

import "github.com/kevin/adventure-archivist/types"

type Repository interface {
	CreateCampaign(*types.CreateCampaignRequest) (*types.Campaign, error)
	GetCampaign(int) (*types.Campaign, error)
	ListCampaigns() ([]*types.Campaign, error)
	GetAdventureRecordsForCampaign(int) ([]*types.AdventureRecord, error)
}
