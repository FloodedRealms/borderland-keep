package services

import (
	"context"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type AdventureService interface {
	CreateAdventureRecordForCampaign(*types.CreateAdventureRequest) (*types.AdventureRecord, error)
	UpdateAdventureRecord(*types.UpdateAdventureRequest) (*types.AdventureRecord, error)
	ListAdventureRecordsForCampaign(string) ([]*types.AdventureRecord, error)
	GetAdventureRecordById(string) (*types.AdventureRecord, error)
	//AddCoinsAdventure(*types.Adventure, *types.Coins) (bool, error)
	//AddGemLootToAdventure(*types.Adventure, *types.Gem) (bool, error)
	//AddJewelleryLootToAdventure(*types.Adventure, *types.Jewellery) (bool, error)
	//AddMagicItemToAdventure(*types.Adventure, *types.MagicItem) (bool, error)
	//AddCombatToAdventure(*types.Adventure, *types.MonsterGroup) (bool, error)
}

type AdventureServiceImpl struct {
	repo repository.Repository
	Ctx  context.Context
}

func NewAdventureRecordService(repo repository.Repository, ctx context.Context) *AdventureServiceImpl {
	return &AdventureServiceImpl{repo, ctx}
}

func (a *AdventureServiceImpl) CreateAdventureRecordForCampaign(r *types.CreateAdventureRequest) (*types.AdventureRecord, error) {
	return a.repo.CreateAdventureRecordForCampaign(r)
}

func (a *AdventureServiceImpl) UpdateAdventureRecord(r *types.UpdateAdventureRequest) (*types.AdventureRecord, error) {
	adventureToUpdate, _ := a.repo.GetAdventureRecordById(types.NewAdventureRecordById(r.ID))
	charactersInCampaign, _ := a.repo.GetCharactersForCampaign(types.NewCampaign(adventureToUpdate.CampaignId))
	coinsToAdd := types.NewCoins(r.Copper, r.Silver, r.Electrum, r.Gold, r.Platinum)
	gemsToUpdate := r.GenerateGemList()
	jewelleryToUpdate := r.GenerateJewelleryList()
	magicItemsToUpdate := r.GenerateMagicItemList()
	combatsToUpdate := r.GenerateCombatList()
	charactersToUpdate := r.GenerateCharacterList()
	updatedAdventure := types.NewAdventureRecord(r.ID, r.CampaignID, 0, *coinsToAdd, gemsToUpdate, jewelleryToUpdate, combatsToUpdate, magicItemsToUpdate, charactersToUpdate, r.Name, adventureToUpdate.CreatedAt, adventureToUpdate.UpdatedAt, types.ArcvhistDate(r.AdventureDate))
	//fullShare, halfShare := updatedAdventure.CalculateXPShares()
	if updatedAdventure.Name != "" && updatedAdventure.Name != adventureToUpdate.Name {
		err := a.repo.UpdateAdventureName(adventureToUpdate, updatedAdventure.Name)
		if err != nil {
			return nil, util.UnableToUpdateAdventure("Name", err.Error())
		}
	}
	if updatedAdventure.AdventureDate != types.NewAdventureRecordById(r.ID).AdventureDate && updatedAdventure.AdventureDate != adventureToUpdate.AdventureDate {
		err := a.repo.UpdateAdventureDate(adventureToUpdate, r.AdventureDate)
		if err != nil {
			return nil, util.UnableToUpdateAdventure("DATE", err.Error())
		}
	}
	err := a.updateAdventureCoins(adventureToUpdate, coinsToAdd)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Coins", err.Error())
	}
	err = a.updateAdventureGems(adventureToUpdate, gemsToUpdate)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Gems", err.Error())
	}
	err = a.updateAdventureJewellery(adventureToUpdate, jewelleryToUpdate)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Jewellery", err.Error())
	}
	err = a.updateAdventureMagicItems(adventureToUpdate, magicItemsToUpdate)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Magic Items", err.Error())
	}
	err = a.updateAdventureCombat(adventureToUpdate, combatsToUpdate)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Combats", err.Error())
	}
	err = a.updateAdventureCharacters(adventureToUpdate, charactersToUpdate, updatedAdventure.FullShareXP, updatedAdventure.HalfShareXP, charactersInCampaign)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Characters", err.Error())
	}

	updatedRecorded, err := a.repo.GetAdventureRecordById(adventureToUpdate)
	return updatedRecorded, err
}

func (a AdventureServiceImpl) updateAdventureCoins(ad *types.AdventureRecord, coins *types.Coins) error {
	_, err := a.repo.UpdateCoinsForAdventure(ad, coins)
	return err
}
func (a AdventureServiceImpl) updateAdventureGems(ad *types.AdventureRecord, gems []types.Gem) error {
	err := a.repo.DeleteGemsForAdventure(ad)
	if err != nil {
		return err
	}
	for _, gem := range gems {
		_, err := a.repo.AddGemToAdventure(ad, &gem)
		if err != nil {
			return err
		}
	}
	return err
}
func (a AdventureServiceImpl) updateAdventureJewellery(ad *types.AdventureRecord, jewellery []types.Jewellery) error {
	err := a.repo.DeleteJewelleryForAdventure(ad)
	if err != nil {
		return err
	}
	for _, gem := range jewellery {
		_, err := a.repo.AddJewelleryToAdventure(ad, &gem)
		if err != nil {
			return err
		}
	}
	return err
}
func (a AdventureServiceImpl) updateAdventureMagicItems(ad *types.AdventureRecord, gems []types.MagicItem) error {
	err := a.repo.DeleteMagicItemsForAdventure(ad)
	if err != nil {
		return err
	}
	for _, gem := range gems {
		_, err := a.repo.AddMagicItemToAdventure(ad, &gem)
		if err != nil {
			return err
		}
	}
	return err
}
func (a AdventureServiceImpl) updateAdventureCombat(ad *types.AdventureRecord, gems []types.MonsterGroup) error {
	err := a.repo.DeleteCombatForAdventure(ad)
	if err != nil {
		return err
	}
	for _, gem := range gems {
		_, err := a.repo.AddCombatToAdventure(ad, &gem)
		if err != nil {
			return err
		}
	}
	return err
}
func (a AdventureServiceImpl) updateAdventureCharacters(ad *types.AdventureRecord, chars []types.AdventureCharacter, fullShareAmount, halfShareAmount int, campChars []types.CharacterRecord) error {
	charMap := map[int]types.CharacterRecord{}
	for _, c := range campChars {
		charMap[c.ID] = c
	}
	err := a.repo.DeleteCharactersForAdventure(ad)
	if err != nil {
		return err
	}
	for _, char := range chars {
		xpToGain := fullShareAmount
		if char.Halfshare {
			xpToGain = halfShareAmount
		}
		c := charMap[char.Details.ID]
		adjustedAmount := c.ApplyPrimeReq(xpToGain)
		if char.Halfshare {
			_, err := a.repo.AddHalfshareCharacterToAdventure(ad, &char, adjustedAmount)
			if err != nil {
				return err
			}
		} else {
			_, err := a.repo.AddFullshareCharacterToAdventure(ad, &char, adjustedAmount)

			if err != nil {
				return err
			}
		}
	}
	return err
}

func (a *AdventureServiceImpl) ListAdventureRecordsForCampaign(i string) ([]*types.AdventureRecord, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	campaign := types.NewCampaign(id)
	return a.repo.GetAdventureRecordsForCampaign(campaign)

}

func (a *AdventureServiceImpl) GetAdventureRecordById(i string) (*types.AdventureRecord, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	ad, err2 := a.repo.GetAdventureRecordById(types.NewAdventureRecordById(id))
	if err2 != nil {
		return nil, util.NotYetImplmented()
	}
	return ad, nil
}

/*
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
}*/
