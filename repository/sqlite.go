package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
	_ "github.com/mattn/go-sqlite3"
)

const campaignTable string = "campaigns"
const createCampaignTable string = `
CREATE TABLE IF NOT EXISTS campaigns (
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT,
	recruitment INTEGER,
	judge TEXT,
	timekeeping TEXT,
	cadence TEXT,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	last_adventure DATETIME
	);
`

const adventureTable string = "adventures"
const createAdventureTable string = `
CREATE TABLE IF NOT EXISTS adventures (
	id INTEGER NOT NULL PRIMARY KEY,
	campaign_id INTEGER,
	name TEXT,
	adventure_date DATETIME,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	FOREIGN KEY(campaign_id) REFERENCES campaigns(id)	
	);
`

const gemTable string = "gems"
const jewelleryTable string = "jewellery"
const magicItemsTable string = "magic_items"
const combatTable string = "monster_groups"
const characterToAdventureTable string = "adventures_to_characters"
const characterTable string = "characters"

var trashInt int = 0
var trashDate time.Time = time.Now()

type SqliteRepo struct {
	db     *sql.DB
	logger *util.Logger
}

func NewSqliteRepo(f string, logger *util.Logger) (*SqliteRepo, error) {
	db, err := sql.Open("sqlite3", f)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createCampaignTable); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createAdventureTable); err != nil {
		return nil, err
	}
	return &SqliteRepo{db: db, logger: logger}, nil
}

func (s SqliteRepo) runQuery(q string, params ...interface{}) (*sql.Rows, error) {
	stmt, err := s.db.Prepare(q)
	util.CheckErr(err)
	return stmt.Query(params...)
}

// Campaigns
func (s SqliteRepo) processCampaignRows(r *sql.Rows) []*types.Campaign {
	campaigns := make([]*types.Campaign, 0)
	for r.Next() {
		current := &types.Campaign{}
		err := r.Scan(&current.ID, &current.Name, &current.Recruitment, &current.Judge, &current.Timekeeping, &current.Cadence, &current.CreatedAt, &current.UpdatedAt, &current.LastAdventure)
		util.CheckErr(err)
		current.Characters, err = s.getCampaignCharacters(current.ID)
		util.CheckErr(err)
		campaigns = append(campaigns, current)
	}
	return campaigns
}

func (s SqliteRepo) selectAllCampaigns() []*types.Campaign {
	tableq := fmt.Sprintf("SELECT * FROM %s", campaignTable)
	rows, err := s.runQuery(tableq)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processCampaignRows(rows)

	return results
}

func (s SqliteRepo) insertCampaign(c types.CreateCampaignRequest) (int, error) {
	s.logger.Debug("tried to insert campaign")
	stmtString := fmt.Sprintf("INSERT INTO %s(name, recruitment, judge, timekeeping, cadence, created_at, updated_at, last_adventure) values(?, ?, ?, ?, ?, ?, ?, ?) ;", campaignTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	res, err := stmt.Exec(c.Name, c.Recruitment, c.Judge, c.Timekeeping, c.Cadence, c.CreatedAt, c.UpdatedAt, c.LastAdventure)
	util.CheckErr(err)
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s SqliteRepo) selectCampaignById(id int) (*types.Campaign, error) {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", campaignTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processCampaignRows(rows)

	return results[0], nil
}

func (s SqliteRepo) CreateCampaign(c *types.CreateCampaignRequest) (*types.Campaign, error) {
	id, err := s.insertCampaign(*c)
	util.CheckErr(err)
	return s.selectCampaignById(id)
}

func (s SqliteRepo) GetCampaign(id int) (*types.Campaign, error) {
	return s.selectCampaignById(id)
}

func (s SqliteRepo) ListCampaigns() ([]*types.Campaign, error) {
	return s.selectAllCampaigns(), nil
}
func (s SqliteRepo) DeleteCampaign(c *types.Campaign) (bool, error) {
	return false, util.NotYetImplmented()
}

// Adventures
func (s SqliteRepo) insertAdventureRecord(a types.CreateAdventureRecordRequest) (int, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(campaign_id, name, created_at, updated_at, adventure_date) values(?, ?, ?, ?, ?)", adventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	res, err := stmt.Exec(a.CampaignID, a.Name, time.Now(), time.Now(), a.AdventureDate)
	util.CheckErr(err)
	stmt2, err := s.db.Prepare(fmt.Sprintf("UPDATE %s SET last_adventure=? WHERE id=?", campaignTable))
	util.CheckErr(err)
	_, err = stmt2.Exec(a.AdventureDate, a.CampaignID)
	util.CheckErr(err)

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}
func (s SqliteRepo) selectAdventureById(id int) *types.Adventure {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", adventureTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processAdventureRows(rows)

	return results[0]
}
func (s SqliteRepo) selectAdventureByCampaignId(c *types.Campaign) ([]*types.Adventure, error) {

	tableq := fmt.Sprintf("SELECT * FROM %s c where c.campaign_id = ?", adventureTable)
	rows, err := s.runQuery(tableq, c.ID)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processAdventureRows(rows)

	return results, nil
}
func (s SqliteRepo) processAdventureRows(r *sql.Rows) []*types.Adventure {
	adventures := make([]*types.Adventure, 0)
	for r.Next() {
		var id int
		var campaignId int
		var name string
		var createdDate time.Time
		var updatedDate time.Time
		var adventureDate time.Time
		copper, silver, electrum, gold, platinum := 0, 0, 0, 0, 0
		err := r.Scan(&id, &campaignId, &name, &adventureDate, &createdDate, &updatedDate, &copper, &silver, &electrum, &gold, &platinum)

		currentCoins := types.NewCoins(copper, silver, electrum, gold, platinum)
		currentGems := s.getGemsForAdventure(id)
		currentJewellery := s.getJewelleryForAdventure(id)
		currentMagicItems := s.getMagicItemsForAdventure(id)
		currentCombat := s.getCombatForAdventure(id)
		currentCharacters := s.getCharactersForAdventure(id)
		current := types.NewAdventureRecord(id, campaignId, *currentCoins, currentGems, currentJewellery, currentCombat, currentMagicItems, currentCharacters, name, createdDate, updatedDate, adventureDate)
		util.CheckErr(err)
		adventures = append(adventures, current)
	}
	return adventures
}

func (s SqliteRepo) getGemsForAdventure(id int) []types.Gem {
	tableq := fmt.Sprintf("SELECT * FROM %s g where g.adventure_id = ?", gemTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()
	gems := make([]types.Gem, 0)
	for rows.Next() {
		s.logger.Debug("Trying to get gem")
		id, number := 0, 0
		name, desc := "", ""
		value := 0.0
		err := rows.Scan(&id, &trashInt, &name, &desc, &value, &number)
		util.CheckErr(err)
		current := types.NewGem(name, desc, value, number, id)
		gems = append(gems, *current)
	}
	return gems
}
func (s SqliteRepo) getJewelleryForAdventure(id int) []types.Jewellery {
	tableq := fmt.Sprintf("SELECT * FROM %s j where j.adventure_id = ?", jewelleryTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()
	jewellery := make([]types.Jewellery, 0)
	for rows.Next() {
		s.logger.Debug("Trying to get jewellery")
		id, number := 0, 0
		name, desc := "", ""
		value := 0.0
		err := rows.Scan(&id, &trashInt, &name, &desc, &value, &number)
		util.CheckErr(err)
		current := types.NewJewellery(name, desc, value, number, id)

		jewellery = append(jewellery, *current)
	}
	return jewellery
}

func (s SqliteRepo) getMagicItemsForAdventure(id int) []types.MagicItem {
	tableq := fmt.Sprintf("SELECT * FROM %s m where m.adventure_id = ?", magicItemsTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()
	magicItems := make([]types.MagicItem, 0)
	for rows.Next() {
		s.logger.Debug("Trying to get magic item")
		id, number, actualValue := 0, 0, 0
		name, desc := "", ""
		value := 0.0
		err := rows.Scan(&id, &trashInt, &name, &desc, &value, &number, &actualValue)
		util.CheckErr(err)
		current := types.NewMagicItem(name, desc, value, actualValue, id)
		magicItems = append(magicItems, *current)
	}
	return magicItems
}

func (s SqliteRepo) getCombatForAdventure(id int) []types.MonsterGroup {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.adventure_id = ?", combatTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()
	combat := make([]types.MonsterGroup, 0)
	for rows.Next() {
		s.logger.Debug("Trying to get combat")
		id, number := 0, 0
		name := ""
		value := 0.0
		err := rows.Scan(&id, &trashInt, &name, &number, &value, &trashInt)
		util.CheckErr(err)
		current := types.NewMonsterGroup(name, number, id, value)

		combat = append(combat, *current)
	}
	return combat
}
func (s SqliteRepo) getCharactersForAdventure(id int) []types.AdventureCharacter {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.adventure_id = ?", characterToAdventureTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()
	characters := make([]types.AdventureCharacter, 0)
	for rows.Next() {
		var rowId int
		var adventureId int
		var characterId int
		var halfStatus int
		err := rows.Scan(&rowId, &adventureId, &characterId, &halfStatus)
		util.CheckErr(err)
		halfshare := false
		if halfStatus == 1 {
			halfshare = true
		}
		currentDetails := s.selectCharacterById(characterId)
		adventureCharacter := types.NewAdventureCharacter(currentDetails, halfshare)
		characters = append(characters, *adventureCharacter)
	}
	return characters

}

func (s SqliteRepo) CreateAdventureRecordForCampaign(a *types.CreateAdventureRecordRequest) (*types.Adventure, error) {
	id, err := s.insertAdventureRecord(*a)
	util.CheckErr(err)
	return s.selectAdventureById(id), nil
}

func (s SqliteRepo) GetAdventureRecordsForCampaign(c *types.Campaign) ([]*types.Adventure, error) {
	return s.selectAdventureByCampaignId(c)
}

func (s SqliteRepo) GetAdventureRecordById(a *types.Adventure) (*types.Adventure, error) {
	return s.selectAdventureById(a.ID), nil
}

func (s SqliteRepo) AddCoinsToAdventure(a *types.Adventure, c *types.Coins) (bool, error) {
	stmt_string := fmt.Sprintf("UPDATE %s SET copper = ?, silver = ?, electrum = ?, gold = ?, platinum = ? WHERE id=?;", adventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(c.Copper.Number, c.Silver.Number, c.Electrum.Number, c.Gold.Number, c.Platinum.Number, a.ID)
	if resErr != nil {
		return false, resErr
	}
	return true, nil

}

func (s SqliteRepo) AddGemToAdventure(a *types.Adventure, g *types.Gem) (bool, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?, ?, ?, ?, ?)", gemTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID, g.Loot.Name, g.Loot.Description, g.Loot.XPValue, g.Loot.Number)
	if resErr != nil {
		return false, resErr
	}
	return true, nil
}

func (s SqliteRepo) AddJewelleryToAdventure(a *types.Adventure, g *types.Jewellery) (bool, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?, ?, ?, ?, ?)", jewelleryTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID, g.Loot.Name, g.Loot.Description, g.Loot.XPValue, g.Loot.Number)
	if resErr != nil {
		return false, resErr
	}
	return true, nil
}

func (s SqliteRepo) AddMagicItemToAdventure(a *types.Adventure, g *types.MagicItem) (bool, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, apparent_value, actual_value) values(?, ?, ?, ?, ?)", magicItemsTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID, g.Loot.Name, g.Loot.Description, g.ApparentValue(), g.ActualValue)
	if resErr != nil {
		return false, resErr
	}
	return true, nil
}

func (s SqliteRepo) AddCombatToAdventure(a *types.Adventure, g *types.MonsterGroup) (bool, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(adventure_id, monster_name, number_defeated, xp_per_monster, total_xp) values(?, ?, ?, ?, ?)", combatTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID, g.XP.Name, g.XP.Number, g.XP.XPValue, g.TotalXPAmount)
	if resErr != nil {
		return false, resErr
	}
	return true, nil
}

// Characters
func (s SqliteRepo) insertCharacter(campaignId int) (int, error) {
	characterToInsert := types.BlankCharacter()
	s.logger.Debug("tried to insert characer")
	stmtString := fmt.Sprintf("INSERT INTO %s(campaign_id, name, current_xp, prime_req_percent, character_level, character_class, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?) ;", characterTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	res, err := stmt.Exec(campaignId, characterToInsert.Name, characterToInsert.CurrentXP, characterToInsert.PrimeReqPercent, characterToInsert.Level, characterToInsert.Class, time.Now(), time.Now())
	util.CheckErr(err)
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil

}

func (s SqliteRepo) selectCharacterById(id int) *types.Character {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", characterTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processCharacterRows(rows)

	return results[0]
}

func (s SqliteRepo) getCampaignCharacters(campaignId int) ([]types.Character, error) {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", characterTable)
	rows, err := s.runQuery(tableq, campaignId)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processCharacterRowsForCampaign(rows)

	return results, nil
}

func (s SqliteRepo) processCharacterRowsForCampaign(r *sql.Rows) []types.Character {
	characters := make([]types.Character, 0)
	for r.Next() {
		current := types.Character{}
		err := r.Scan(&current.ID, &trashInt, &current.Name, &current.CurrentXP, &current.PrimeReqPercent, &current.Level, &current.Class, &trashDate, &trashDate)
		util.CheckErr(err)
		characters = append(characters, current)
	}
	return characters
}

func (s SqliteRepo) processCharacterRows(r *sql.Rows) []*types.Character {
	characters := make([]*types.Character, 0)
	for r.Next() {
		current := &types.Character{}
		err := r.Scan(&current.ID, &trashInt, &current.Name, &current.CurrentXP, &current.PrimeReqPercent, &current.Level, &current.Class, &trashDate, &trashDate)
		util.CheckErr(err)
		characters = append(characters, current)
	}
	return characters
}

func (s SqliteRepo) CreateCharacterForCampaign(campaign *types.Campaign) (*types.Character, error) {
	id, err := s.insertCharacter(campaign.ID)
	if err != nil {
		return nil, err
	}
	return s.selectCharacterById(id), nil
}
func (s SqliteRepo) AddCharacterToAdventure(ad *types.Adventure, char *types.Character, isGettingHalfshare bool) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share) values(?, ?, ?) ;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	if isGettingHalfshare {
		_, err := stmt.Exec(ad.ID, char.ID, 1)
		if err != nil {
			return false, err
		}
	} else {
		_, err := stmt.Exec(ad.ID, char.ID, 0)
		if err != nil {
			return false, err
		}
	}
	return true, nil

}
func (s SqliteRepo) RemoveCharacterFromAdventure(ad *types.Adventure, char *types.Character) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("DELETE FROM %s WHERE adventure_id = ? AND character_id = ?;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	_, err = stmt.Exec(ad.ID, char.ID)
	if err != nil {
		return false, err
	}

	return true, nil
}
func (s SqliteRepo) ChangeCharacterShares(ad *types.Adventure, char *types.Character, isGettingHalfShare bool) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("UPDATE %s set half_share=? WHERE adventure_id = ? AND character_id = ?;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	if isGettingHalfShare {
		_, err = stmt.Exec(ad.ID, char.ID, 1)
		if err != nil {
			return false, err
		}
	} else {
		_, err = stmt.Exec(ad.ID, char.ID, 0)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s SqliteRepo) GetCharactersForCampaign(camp *types.Campaign) ([]types.Character, error) {
	return s.getCampaignCharacters(camp.ID)
}
