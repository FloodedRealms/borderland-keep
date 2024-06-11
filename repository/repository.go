package repository

import "github.com/floodedrealms/adventure-archivist/types"

type Repository interface {
	CreateCampaign(*types.CreateCampaignRequest) (*types.Campaign, error)
	GetCampaign(int) (*types.Campaign, error)
	DeleteCampaign(*types.Campaign) (bool, error)
	ListCampaigns() ([]*types.Campaign, error)

	GetAdventureRecordsForCampaign(*types.Campaign) ([]*types.Adventure, error)
	CreateAdventureRecordForCampaign(*types.CreateAdventureRecordRequest) (*types.Adventure, error)
	GetAdventureRecordById(*types.Adventure) (*types.Adventure, error)
	AddCoinsToAdventure(a *types.Adventure, c *types.Coins) (bool, error)
	AddGemToAdventure(*types.Adventure, *types.Gem) (bool, error)
	AddJewelleryToAdventure(*types.Adventure, *types.Jewellery) (bool, error)
	AddMagicItemToAdventure(*types.Adventure, *types.MagicItem) (bool, error)
	AddCombatToAdventure(*types.Adventure, *types.MonsterGroup) (bool, error)

	CreateCharacterForCampaign(*types.Campaign) (*types.Character, error)
	GetCharactersForCampaign(*types.Campaign) ([]types.Character, error)
	AddCharacterToAdventure(*types.Adventure, *types.Character, bool) (bool, error)
	RemoveCharacterFromAdventure(*types.Adventure, *types.Character) (bool, error)
	ChangeCharacterShares(*types.Adventure, *types.Character, bool) (bool, error)
}
