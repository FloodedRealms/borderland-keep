package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/internal/repository"
	"github.com/floodedrealms/adventure-archivist/internal/util"
	"github.com/floodedrealms/adventure-archivist/types"
)

type AdventureService struct {
	repo repository.Repository
	Ctx  context.Context
}

const adventureTable string = "adventures"

func NewAdventureRecordService(repo repository.Repository, ctx context.Context) *AdventureService {
	return &AdventureService{repo, ctx}
}

func (a *AdventureService) CreateAdventureRecordForCampaign(r *types.AdventureRecord) (*types.AdventureRecord, error) {
	return a.repo.CreateAdventureRecordForCampaign(r)
}

func (a *AdventureService) UpdateAdventureRecord(r *types.AdventureRecord) (*types.AdventureRecord, error) {
	adventureToUpdate, err := a.repo.GetAdventureRecordById(r)
	if err != nil {
		return nil, err
	}
	charactersInCampaign, _ := a.repo.GetCharactersForCampaign(types.NewCampaign(adventureToUpdate.CampaignId))
	fullShare, halfShare := r.CalculateXPShares()
	if r.Name != "" && r.Name != adventureToUpdate.Name {
		err := a.repo.UpdateAdventureName(adventureToUpdate, r.Name)
		if err != nil {
			return nil, util.UnableToUpdateAdventure("Name", err.Error())
		}
	}
	if r.AdventureDate != types.NewAdventureRecordById(r.Id).AdventureDate && r.AdventureDate != adventureToUpdate.AdventureDate {
		err := a.repo.UpdateAdventureDate(adventureToUpdate, r.AdventureDate)
		if err != nil {
			return nil, util.UnableToUpdateAdventure("DATE", err.Error())
		}
	}
	err = a.updateAdventureCoins(adventureToUpdate, &r.Coins)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Coins", err.Error())
	}
	err = a.updateAdventureGems(adventureToUpdate, r.Gems)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Gems", err.Error())
	}
	err = a.updateAdventureJewellery(adventureToUpdate, r.Jewellery)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Jewellery", err.Error())
	}
	err = a.updateAdventureMagicItems(adventureToUpdate, r.MagicItems)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Magic Items", err.Error())
	}
	err = a.updateAdventureCombat(adventureToUpdate, r.Combat)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Combats", err.Error())
	}
	err = a.updateAdventureCharacters(adventureToUpdate, r.Characters, fullShare, halfShare, charactersInCampaign)
	if err != nil {
		return nil, util.UnableToUpdateAdventure("Characters", err.Error())
	}

	updatedRecorded, err := a.repo.GetAdventureRecordById(adventureToUpdate)
	return updatedRecorded, err
}

func (a AdventureService) updateAdventureCoins(ad *types.AdventureRecord, coins *types.Coins) error {
	_, err := a.repo.UpdateCoinsForAdventure(ad, coins)
	return err
}
func (a AdventureService) updateAdventureGems(ad *types.AdventureRecord, gems []types.Gem) error {
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
func (a AdventureService) updateAdventureJewellery(ad *types.AdventureRecord, jewellery []types.Jewellery) error {
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
func (a AdventureService) updateAdventureMagicItems(ad *types.AdventureRecord, gems []types.MagicItem) error {
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
func (a AdventureService) updateAdventureCombat(ad *types.AdventureRecord, gems []types.MonsterGroup) error {
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
func (a AdventureService) updateAdventureCharacters(ad *types.AdventureRecord, chars []types.AdventureCharacter, fullShareAmount, halfShareAmount int, campChars []types.CharacterRecord) error {
	charMap := map[int]types.CharacterRecord{}
	for _, c := range campChars {
		charMap[c.Id] = c
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
		c := charMap[char.Id]
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

func (a *AdventureService) ListAdventureRecordsForCampaign(i string) ([]*types.AdventureRecord, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	campaign := types.NewCampaign(id)
	return a.repo.GetAdventureRecordsForCampaign(campaign)

}

func (a *AdventureService) GetAdventureRecordById(i string) (*types.AdventureRecord, error) {
	id, err := strconv.Atoi(i)
	util.CheckErr(err)
	ad, err2 := a.repo.GetAdventureRecordById(types.NewAdventureRecordById(id))
	if err2 != nil {
		return nil, util.NotYetImplmented()
	}
	return ad, nil
}

func (a AdventureService) GetCoinsForAdventure(i string) (*types.Coins, error) {
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, err
	}
	coins, err := a.repo.GetCoinsForAdventure(types.NewAdventureRecordById(id))
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func (a AdventureService) UpdateAdventureCoins(id string, data map[string]string) (*types.AdventureRecord, error) {
	i, _ := strconv.Atoi(id)
	copper, _ := stripGoodValueFromCoinForm("copper", data)
	silver, _ := stripGoodValueFromCoinForm("silver", data)
	electrum, _ := stripGoodValueFromCoinForm("electrum", data)
	gold, _ := stripGoodValueFromCoinForm("gold", data)
	platinum, _ := stripGoodValueFromCoinForm("platinum", data)
	stmtStr := fmt.Sprintf("UPDATE %s set copper=?, silver=?, electrum=?, gold=?, platinum=? WHERE ID =?", adventureTable)
	a.repo.ExecuteQuery(stmtStr, copper, silver, electrum, gold, platinum, id)
	return a.repo.GetAdventureRecordById(types.NewAdventureRecordById(i))
}

func stripGoodValueFromCoinForm(coin string, data map[string]string) (int, error) {
	val, ok := data[coin]
	if ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, nil
}

/*
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
