package repository

import "github.com/floodedrealms/adventure-archivist/types"

type Repository interface {
	CreateCampaign(*types.CreateCampaignRequest) (*types.Campaign, error)
	GetCampaign(int) (*types.Campaign, error)
	DeleteCampaign(*types.Campaign) (bool, error)
	ListCampaigns() ([]*types.Campaign, error)

	GetAdventureRecordsForCampaign(*types.Campaign) ([]*types.AdventureRecord, error)
	CreateAdventureRecordForCampaign(*types.CreateAdventureRecordRequest) (*types.AdventureRecord, error)
	GetAdventureRecordById(*types.AdventureRecord) (*types.AdventureRecord, error)
	AddGemToAdventure(*types.AdventureRecord, *types.Gem) (bool, error)
}
