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

func (a *AdventureService) UpdateAdventureCoins(adventureId int, coins types.Coins) error {
	stmt := fmt.Sprintf("UPDATE %s SET copper = ?, silver = ?, electrum = ?, gold = ?, platinum = ? WHERE id = ?;")
	_, err := a.repo.ExecuteQuery(stmt, coins.Copper.Number, coins.Silver.Number, coins.Electrum.Number, coins.Gold.Number, coins.Platinum.Number, adventureId)
	return err
}

func (a *AdventureService) DeleteAdventure(id int) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE id=?;", adventureTable)
	_, err := a.repo.ExecuteQuery(q, id)
	return err
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
		adventureToReturn.CalculateXPShares()

	} else {
		return nil, util.UnableToFindResourceWithId("adventure", id)
	}
	return adventureToReturn, nil
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
		formData.CreateXPFunc()
		if formData.Halfshare {
			paramList = []interface{}{aId, formData.Id, formData.Halfshare, formData.ShowAdjustedXP(halfshare)}
		} else {
			paramList = []interface{}{aId, formData.Id, formData.Halfshare, formData.ShowAdjustedXP(fullshare)}
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
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.ApparentValue, &cur.ActualValue)
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
		rows.Scan(&cur.Id, &trashInt, &cur.Name, &cur.Description, &cur.ApparentValue, &cur.ActualValue)
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
	stmtStr := fmt.Sprintf("SELECT atc.character_id, atc.character_name, atc.on_adventure, c.prime_req_percent FROM %s atc JOIN characters c ON atc.character_id = c.id WHERE adventure_id=? ORDER BY character_name ASC;", possibleCharactersView)
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
		rows.Scan(&cur.Id, &cur.Name, &wasThere, &cur.Preq)
		results = append(results, cur)
		if wasThere == "Yes" {
			onAdventure = append(onAdventure, true)
		} else {
			onAdventure = append(onAdventure, false)
		}
	}
	return results, onAdventure, nil
}
