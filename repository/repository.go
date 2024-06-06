package repository

import "github.com/kevin/adventure-archivist/types"

type Repository interface {
	CreateCampaign(*types.CreateCampaignRequest) (*types.Campaign, error)
	GetCampaign(int) (*types.Campaign, error)
	DeleteCampaign(*types.Campaign) (bool, error)
	ListCampaigns() ([]*types.Campaign, error)

	GetAdventureRecordsForCampaign(*types.Campaign) ([]*types.AdventureRecord, error)
	CreateAdventureRecordForCampaign(*types.CreateAdventureRecordRequest) (*types.AdventureRecord, error)
}
