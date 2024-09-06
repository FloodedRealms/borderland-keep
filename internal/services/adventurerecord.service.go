package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/types"
)

type AdventureService struct {
	repo repository.Repository
	Ctx  context.Context
}

const adventureTable string = "adventures"
const gemTable string = "gems"
const jewelleryTable string = "jewellery"
const magicItemTable string = "magic_items"
const combatTable string = "monster_groups"
const adventureToCharactersTable string = "adventures_to_characters"

const characterToAdventureView string = "adventures_to_character_name"
const possibleCharactersView string = "possible_characters_for_adventure"

func NewAdventureRecordService(repo repository.Repository, ctx context.Context) *AdventureService {
	return &AdventureService{repo, ctx}
}

func (a *AdventureService) CreateAdventureRecordForCampaign(r *types.AdventureRecord) (*types.AdventureRecord, error) {
	return a.repo.CreateAdventureRecordForCampaign(r)
}

func (a *AdventureService) CreateNewAdventureRecordForCampaign(campaignId int) (*types.AdventureRecord, error) {
	stmt := fmt.Sprintf("INSERT INTO %s(name, campaign_id, adventure_date, created_at, updated_at) values(?,?,?,?,?);", adventureTable)
	time := time.Now()
	result, err := a.repo.ExecuteQuery(stmt, "New Adventure", campaignId, time, time, time)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	selectNewQ := fmt.Sprintf("SELECT a.id, a.campaign_id, a.name, a.adventure_date FROM %s a where a.id = ?;", adventureTable)
	rows, err := a.repo.RunQuery(selectNewQ, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	adventures := make([]*types.AdventureRecord, 0)
	for rows.Next() {
		current := &types.AdventureRecord{}
		err := rows.Scan(&current.Id, &current.CampaignId, &current.Name, &current.AdventureDate)
		if err != nil {
			return nil, err
		}
		adventures = append(adventures, current)
	}
	return adventures[0], nil
}

func (a *AdventureService) ModifyMetadata(ad types.AdventureRecord) error {
	selectNewQ := fmt.Sprintf("UPDATE %s set name=?, adventure_date=? WHERE id=?;", adventureTable)
	_, err := a.repo.ExecuteQuery(selectNewQ, ad.Name, ad.AdventureDate, ad.Id)
	return err
}

func (a *AdventureService) DeleteAdventure(id int) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE id=?;", adventureTable)
	_, err := a.repo.ExecuteQuery(q, id)
	return err
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
	stmt := fmt.Sprintf("UPATE %s set copper=?, silver=? electrum=? gold=? platinum=? WHERE id=?", adventureTable)
	_, err := a.repo.ExecuteQuery(stmt, coins.Copper.Number, coins.Silver.Number, coins.Electrum.Number, coins.Gold.Number, coins.Platinum.Number, ad.Id)
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

func (a *AdventureService) GetAdventureRecordById(id int) (*types.AdventureRecord, error) {
	adventureToReturn := &types.AdventureRecord{}
	stmtStr := fmt.Sprintf("SELECT * FROM %s a where a.id = ?", adventureTable)

	adventureResults, err := a.repo.RunQuery(stmtStr, id)
	if err != nil {
		return nil, err
	}
	defer adventureResults.Close()
	if adventureResults.Next() {
		var (
			trashDate time.Time
			copper    int
			silver    int
			electrum  int
			gold      int
			platinum  int
		)
		err := adventureResults.Scan(&adventureToReturn.Id, &adventureToReturn.CampaignId, &adventureToReturn.Name, &adventureToReturn.AdventureDate, &trashDate, &trashDate, &copper, &silver, &electrum, &gold, &platinum, &adventureToReturn.GameDays)
		if err != nil {
			return nil, err
		}
		adventureToReturn.Coins = *types.NewCoins(copper, silver, electrum, gold, platinum)
		g, err := a.GetGemsForAdventure(id)
		if err != nil {
			return nil, err
		}
		j, err := a.GetJewelleryForAdventure(id)
		if err != nil {
			return nil, err
		}
		mi, err := a.GetMagicItemsForAdventure(id)
		if err != nil {
			return nil, err
		}
		c, err := a.GetCombatForAdventure(id)
		if err != nil {
			return nil, err
		}
		chars, err := a.GetCharactersForAdventure(id)
		if err != nil {
			return nil, err
		}
		adventureToReturn.Gems = g
		adventureToReturn.Jewellery = j
		adventureToReturn.MagicItems = mi
		adventureToReturn.Combat = c
		adventureToReturn.Characters = chars
		f, h := adventureToReturn.CalculateXPShares()
		adventureToReturn.HalfShareXP = h
		adventureToReturn.FullShareXP = f

	} else {
		return nil, util.UnableToFindResourceWithId("adventure", id)
	}
	return adventureToReturn, nil
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

func (a AdventureService) UpdateAdventureCoins(id string, data map[string]string) (*types.Coins, error) {
	i, _ := strconv.Atoi(id)
	copper, _ := stripGoodNumberValueFromFormData("copper", data)
	silver, _ := stripGoodNumberValueFromFormData("silver", data)
	electrum, _ := stripGoodNumberValueFromFormData("electrum", data)
	gold, _ := stripGoodNumberValueFromFormData("gold", data)
	platinum, _ := stripGoodNumberValueFromFormData("platinum", data)
	stmtStr := fmt.Sprintf("UPDATE %s set copper=?, silver=?, electrum=?, gold=?, platinum=? WHERE ID =?", adventureTable)
	a.repo.ExecuteQuery(stmtStr, copper, silver, electrum, gold, platinum, id)
	return a.repo.GetCoinsForAdventure(types.NewAdventureRecordById(i))
}

func stripGoodNumberValueFromFormData(field string, data map[string]string) (int, error) {
	val, ok := data[field]
	if ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, nil
}

func (a AdventureService) GetGemById(id string) (*types.Gem, error) {
	gemId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	stmtStr := fmt.Sprintf("SELECT * FROM %s WHERE id=?;", gemTable)
	rows, qErr := a.repo.RunQuery(stmtStr, gemId)
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()
	results := make([]*types.Gem, 0)
	for rows.Next() {
		cur := &types.Gem{}
		var trashInt int
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.XPValue, &cur.Number)
		cur.LootType = types.GemLoot
		results = append(results, cur)
	}
	if len(results) < 1 {
		return nil, fmt.Errorf("unable to locate gem with id %d", gemId)
	}
	return results[0], nil
}

func (a AdventureService) GetGemsForAdventure(aId int) ([]types.Gem, error) {
	stmtStr := fmt.Sprintf("SELECT * FROM %s WHERE adventure_id=?;", gemTable)
	rows, qErr := a.repo.RunQuery(stmtStr, aId)
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()
	results := make([]types.Gem, 0)
	for rows.Next() {
		cur := types.Gem{}
		var trashInt int
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.XPValue, &cur.Number)
		cur.LootType = types.GemLoot
		results = append(results, cur)
	}
	return results, nil
}

func (a AdventureService) SaveGem(gemId string, data map[string]string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	amount, err := stripGoodNumberValueFromFormData("number", data)
	if err != nil {
		return err
	}
	xpValue, err := stripGoodNumberValueFromFormData("value", data)
	if err != nil {
		return err
	}
	desc, dOk := data["description"]
	name, nOk := data["name"]
	if !dOk {
		desc = ""
	}
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("UPDATE %s set name=?, description=?, value=?, total=? WHERE ID =?", gemTable)
	_, err = a.repo.ExecuteQuery(stmtStr, name, desc, xpValue, amount, id)
	return err
}

func (a AdventureService) SaveJewellery(gemId string, data map[string]string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	amount, err := stripGoodNumberValueFromFormData("number", data)
	if err != nil {
		return err
	}
	xpValue, err := stripGoodNumberValueFromFormData("value", data)
	if err != nil {
		return err
	}
	desc, dOk := data["description"]
	name, nOk := data["name"]
	if !dOk {
		desc = ""
	}
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("UPDATE %s set name=?, description=?, value=?, total=? WHERE ID =?", jewelleryTable)
	_, err = a.repo.ExecuteQuery(stmtStr, name, desc, xpValue, amount, id)
	return err
}

func (a AdventureService) SaveCombat(gemId string, data map[string]string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	amount, err := stripGoodNumberValueFromFormData("number", data)
	if err != nil {
		return err
	}
	xpValue, err := stripGoodNumberValueFromFormData("value", data)
	if err != nil {
		return err
	}
	name, nOk := data["name"]
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("UPDATE %s set monster_name=?, xp_per_monster=?, number_defeated=? WHERE ID =?", combatTable)
	_, err = a.repo.ExecuteQuery(stmtStr, name, xpValue, amount, id)
	return err
}

func (a AdventureService) SaveMagicItem(gemId string, data map[string]string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	xp, err := stripGoodNumberValueFromFormData("xp_value", data)
	if err != nil {
		return err
	}
	gold, err := stripGoodNumberValueFromFormData("gold_value", data)
	if err != nil {
		return err
	}
	desc, dOk := data["description"]
	name, nOk := data["name"]
	if !dOk {
		desc = ""
	}
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("UPDATE %s set name=?, description=?, apparent_value=?, actual_value=? WHERE ID =?", magicItemTable)
	_, err = a.repo.ExecuteQuery(stmtStr, name, desc, xp, gold, id)
	return err
}

func (a AdventureService) SaveNewGem(adventureId int, data map[string]string) error {
	amount, err := stripGoodNumberValueFromFormData("number", data)
	if err != nil {
		return err
	}
	xpValue, err := stripGoodNumberValueFromFormData("xp-value", data)
	if err != nil {
		return err
	}
	desc, dOk := data["description"]
	name, nOk := data["name"]
	if !dOk {
		desc = ""
	}
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?,?,?,?,?)", gemTable)
	_, err = a.repo.ExecuteQuery(stmtStr, adventureId, name, desc, xpValue, amount)
	return err
}

func (a AdventureService) SaveNewJewellery(adventureId int, data map[string]string) error {
	amount, err := stripGoodNumberValueFromFormData("number", data)
	if err != nil {
		return err
	}
	xpValue, err := stripGoodNumberValueFromFormData("value", data)
	if err != nil {
		return err
	}
	desc, dOk := data["description"]
	name, nOk := data["name"]
	if !dOk {
		desc = ""
	}
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?,?,?,?,?)", jewelleryTable)
	_, err = a.repo.ExecuteQuery(stmtStr, adventureId, name, desc, xpValue, amount)
	return err
}

func (a AdventureService) SaveNewCombat(adventureId int, data map[string]string) error {
	amount, err := stripGoodNumberValueFromFormData("number_defeated", data)
	if err != nil {
		return err
	}
	xpValue, err := stripGoodNumberValueFromFormData("xp_value", data)
	if err != nil {
		return err
	}
	name, nOk := data["name"]
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("INSERT INTO %s(adventure_id, monster_name, xp_per_monster, number_defeated) values(?,?,?,?)", combatTable)
	_, err = a.repo.ExecuteQuery(stmtStr, adventureId, name, xpValue, amount)
	return err
}

func (a AdventureService) SaveNewMagicItem(adventureId int, data map[string]string) error {
	amount, err := stripGoodNumberValueFromFormData("gold_value", data)
	if err != nil {
		return err
	}
	xpValue, err := stripGoodNumberValueFromFormData("xp_value", data)
	if err != nil {
		return err
	}
	desc, dOk := data["description"]
	name, nOk := data["name"]
	if !dOk {
		desc = ""
	}
	if !nOk {
		name = ""
	}
	stmtStr := fmt.Sprintf("INSERT INTO %s(adventure_id, name, apparent_value, actual_value, total) values(?,?,?,?,?)", magicItemTable)
	_, err = a.repo.ExecuteQuery(stmtStr, adventureId, name, desc, xpValue, amount)
	return err
}

func (a AdventureService) DeleteGem(gemId string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE ID =?", gemTable)
	_, err = a.repo.ExecuteQuery(stmtStr, id)
	return err
}

func (a AdventureService) DeleteGems(aId string) error {
	id, err := strconv.Atoi(aId)
	if err != nil {
		return err
	}
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE adventure_id =?", gemTable)
	_, err = a.repo.ExecuteQuery(stmtStr, id)
	return err
}

func (a AdventureService) ModifyGems(aId int, data []map[string]string) error {
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE adventure_id =?", gemTable)
	queries := []string{stmtStr}
	firstParamList := []interface{}{aId}
	params := [][]interface{}{firstParamList}
	for _, formData := range data {
		amount, err := stripGoodNumberValueFromFormData("number", formData)
		if err != nil {
			return err
		}
		xpValue, err := stripGoodNumberValueFromFormData("xp-value", formData)
		if err != nil {
			return err
		}
		desc, dOk := formData["description"]
		name, nOk := formData["name"]
		if !dOk {
			desc = ""
		}
		if !nOk {
			name = ""
		}
		queries = append(queries, fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?,?,?,?,?)", gemTable))
		paramList := []interface{}{aId, name, desc, xpValue, amount}
		params = append(params, paramList)

	}
	err := a.repo.DoTransaction(queries, params)
	return err
}

func (a AdventureService) ModifyJewellery(aId int, data []map[string]string) error {
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE adventure_id =?", jewelleryTable)
	queries := []string{stmtStr}
	firstParamList := []interface{}{aId}
	params := [][]interface{}{firstParamList}
	for _, formData := range data {
		amount, err := stripGoodNumberValueFromFormData("number", formData)
		if err != nil {
			return err
		}
		xpValue, err := stripGoodNumberValueFromFormData("xp-value", formData)
		if err != nil {
			return err
		}
		desc, dOk := formData["description"]
		name, nOk := formData["name"]
		if !dOk {
			desc = ""
		}
		if !nOk {
			name = ""
		}
		queries = append(queries, fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?,?,?,?,?)", jewelleryTable))
		paramList := []interface{}{aId, name, desc, xpValue, amount}
		params = append(params, paramList)

	}
	err := a.repo.DoTransaction(queries, params)
	return err
}

func (a AdventureService) ModifyCombat(aId int, data []map[string]string) error {
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE adventure_id =?", combatTable)
	queries := []string{stmtStr}
	firstParamList := []interface{}{aId}
	params := [][]interface{}{firstParamList}
	for _, formData := range data {
		amount, err := stripGoodNumberValueFromFormData("number", formData)
		if err != nil {
			return err
		}
		xpValue, err := stripGoodNumberValueFromFormData("xp-value", formData)
		if err != nil {
			return err
		}
		//desc, dOk := formData["description"]
		name, nOk := formData["name"]
		/*	if !dOk {
			desc = ""
		}*/
		if !nOk {
			name = ""
		}
		c := types.NewMonsterGroup(name, "", amount, xpValue)
		queries = append(queries, fmt.Sprintf("INSERT INTO %s(adventure_id, monster_name, number_defeated, xp_per_monster, total_xp) values(?,?,?,?,?)", combatTable))
		paramList := []interface{}{aId, name, amount, xpValue, int(c.TotalXPAmount())}
		params = append(params, paramList)

	}
	err := a.repo.DoTransaction(queries, params)
	return err
}

func (a AdventureService) ModifyMagicItems(aId int, data []map[string]string) error {
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE adventure_id =?", magicItemTable)
	queries := []string{stmtStr}
	firstParamList := []interface{}{aId}
	params := [][]interface{}{firstParamList}
	for _, formData := range data {
		xpValue, err := stripGoodNumberValueFromFormData("xp-value", formData)
		if err != nil {
			return err
		}
		//desc, dOk := formData["description"]
		name, nOk := formData["name"]
		/*	if !dOk {
			desc = ""
		}*/
		if !nOk {
			name = ""
		}
		queries = append(queries, fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, apparent_value, actual_value) values(?,?,?,?,?)", magicItemTable))
		paramList := []interface{}{aId, name, "", xpValue, xpValue}
		params = append(params, paramList)

	}
	err := a.repo.DoTransaction(queries, params)
	return err
}

func (a AdventureService) ModifyCharacters(aId int, data []types.AdventureCharacter, halfshare, fullshare int) error {
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE adventure_id =?", adventureToCharactersTable)
	queries := []string{stmtStr}
	firstParamList := []interface{}{aId}
	params := [][]interface{}{firstParamList}
	for _, formData := range data {
		var paramList []interface{}
		queries = append(queries, fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share, xp_gained) values(?,?,?,?)", adventureToCharactersTable))
		if formData.Halfshare {
			paramList = []interface{}{aId, formData.Id, formData.Halfshare, halfshare}
		} else {
			paramList = []interface{}{aId, formData.Id, formData.Halfshare, fullshare}
		}
		params = append(params, paramList)
	}
	err := a.repo.DoTransaction(queries, params)
	return err
}

func (a AdventureService) DeleteJewellery(gemId string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE ID =?", jewelleryTable)
	_, err = a.repo.ExecuteQuery(stmtStr, id)
	return err
}

func (a AdventureService) DeleteCombat(gemId string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE ID =?", combatTable)
	_, err = a.repo.ExecuteQuery(stmtStr, id)
	return err
}

func (a AdventureService) DeleteMagicItem(gemId string) error {
	id, err := strconv.Atoi(gemId)
	if err != nil {
		return err
	}
	stmtStr := fmt.Sprintf("DELETE FROM %s WHERE ID =?", magicItemTable)
	_, err = a.repo.ExecuteQuery(stmtStr, id)
	return err
}

func (a AdventureService) GetJewelleryById(id string) (*types.Jewellery, error) {
	jewelleryId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	stmtStr := fmt.Sprintf("SELECT * FROM %s WHERE id=?;", jewelleryTable)
	rows, qErr := a.repo.RunQuery(stmtStr, jewelleryId)
	if qErr != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]*types.Jewellery, 0)
	for rows.Next() {
		cur := &types.Jewellery{}
		var trashInt int
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.XPValue, &cur.Number)
		cur.LootType = types.JewelleryLoot
		results = append(results, cur)
	}
	if len(results) < 1 {
		return nil, fmt.Errorf("unable to locate jewellery with id %d", jewelleryId)
	}
	return results[0], nil
}

func (a AdventureService) GetJewelleryForAdventure(id int) ([]types.Jewellery, error) {
	results := make([]types.Jewellery, 0)
	stmtStr := fmt.Sprintf("SELECT * FROM %s WHERE adventure_id=?;", jewelleryTable)
	rows, qErr := a.repo.RunQuery(stmtStr, id)
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()
	for rows.Next() {
		cur := types.Jewellery{}
		var trashInt int
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.XPValue, &cur.Number)
		cur.LootType = types.GemLoot
		results = append(results, cur)
	}
	return results, nil
}

func (a AdventureService) GetCombatById(id string) (*types.MonsterGroup, error) {
	combatId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	stmtStr := fmt.Sprintf("SELECT id, monster_name, number_defeated, xp_per_monster FROM %s WHERE id=?;", combatTable)
	rows, qErr := a.repo.RunQuery(stmtStr, combatId)
	if qErr != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]*types.MonsterGroup, 0)
	for rows.Next() {
		cur := &types.MonsterGroup{}
		rows.Scan(&cur.Id, &cur.Name, &cur.NumberDefeated, &cur.XPPerOneKill)
		results = append(results, cur)
	}
	if len(results) < 1 {
		return nil, fmt.Errorf("unable to locate combat with id %d", combatId)
	}
	return results[0], nil
}

func (a AdventureService) GetCombatForAdventure(id int) ([]types.MonsterGroup, error) {
	results := make([]types.MonsterGroup, 0)
	stmtStr := fmt.Sprintf("SELECT id, monster_name, number_defeated, xp_per_monster FROM %s WHERE adventure_id=?;", combatTable)
	rows, qErr := a.repo.RunQuery(stmtStr, id)
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()
	for rows.Next() {
		cur := types.MonsterGroup{}
		rows.Scan(&cur.Id, &cur.Name, &cur.NumberDefeated, &cur.XPPerOneKill)
		results = append(results, cur)
	}
	return results, nil
}

func (a AdventureService) GetMagicItemById(id string) (*types.MagicItem, error) {
	jewelleryId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	stmtStr := fmt.Sprintf("SELECT * FROM %s WHERE id=?;", magicItemTable)
	rows, qErr := a.repo.RunQuery(stmtStr, jewelleryId)
	if qErr != nil {
		return nil, qErr
	}
	defer rows.Close()
	results := make([]*types.MagicItem, 0)
	for rows.Next() {
		cur := &types.MagicItem{}
		var trashInt int
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.XPValue, &cur.GoldValue)
		results = append(results, cur)
	}
	if len(results) < 1 {
		return nil, fmt.Errorf("unable to locate Magic Item with id %d", jewelleryId)
	}
	return results[0], nil
}

func (a AdventureService) GetMagicItemsForAdventure(id int) ([]types.MagicItem, error) {
	stmtStr := fmt.Sprintf("SELECT * FROM %s WHERE adventure_id=?;", magicItemTable)
	rows, err := a.repo.RunQuery(stmtStr, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]types.MagicItem, 0)
	for rows.Next() {
		cur := types.MagicItem{}
		var trashInt int
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.XPValue, &cur.GoldValue)
		results = append(results, cur)
	}
	return results, nil
}

func (a AdventureService) GetCharactersForAdventure(id int) ([]types.AdventureCharacter, error) {
	stmtStr := fmt.Sprintf("SELECT atc.character_id, atc.half_share, atc.name, c.prime_req_percent FROM %s atc RIGHT JOIN characters c ON c.id =atc.character_id WHERE adventure_id=?;", characterToAdventureView)
	rows, err := a.repo.RunQuery(stmtStr, id)
	if err != nil {
		return nil, err
	}
	results := make([]types.AdventureCharacter, 0)
	defer rows.Close()
	for rows.Next() {
		cur := types.AdventureCharacter{}
		rows.Scan(&cur.Id, &cur.Halfshare, &cur.Name, &cur.Preq)
		cur.CreateXPFunc()
		results = append(results, cur)
	}
	return results, nil
}

func (a AdventureService) GetPossibleCharactersForAdventure(id int) ([]types.AdventureCharacter, []bool, error) {
	stmtStr := fmt.Sprintf("SELECT atc.character_id, atc.character_name, atc.on_adventure FROM %s atc WHERE adventure_id=? ORDER BY character_name ASC;", possibleCharactersView)
	rows, err := a.repo.RunQuery(stmtStr, id)
	if err != nil {
		return nil, nil, err
	}
	results := make([]types.AdventureCharacter, 0)
	onAdventure := make([]bool, 0)
	defer rows.Close()
	for rows.Next() {
		cur := types.AdventureCharacter{}
		wasThere := ""
		rows.Scan(&cur.Id, &cur.Name, &wasThere)
		results = append(results, cur)
		if wasThere == "Yes" {
			onAdventure = append(onAdventure, true)
		} else {
			onAdventure = append(onAdventure, false)
		}
	}
	return results, onAdventure, nil
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
