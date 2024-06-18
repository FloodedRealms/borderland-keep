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
const apiUserTable string = "api_users"

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
func (s SqliteRepo) executeQuery(q string, params ...interface{}) (sql.Result, error) {
	stmt, err := s.db.Prepare(q)
	util.CheckErr(err)
	return stmt.Exec(params...)
}

// Campaigns
func (s SqliteRepo) processCampaignRows(r *sql.Rows) []*types.CampaignRecord {
	campaigns := make([]*types.CampaignRecord, 0)
	for r.Next() {
		current := &types.CampaignRecord{}
		err := r.Scan(&current.ID, &current.Name, &current.Recruitment, &current.Judge, &current.Timekeeping, &current.Cadence, &current.CreatedAt, &current.UpdatedAt, &current.LastAdventure, &current.ClientId)
		util.CheckErr(err)
		current.Characters, err = s.getCampaignCharacters(current.ID)
		util.CheckErr(err)
		campaigns = append(campaigns, current)
	}
	return campaigns
}

func (s SqliteRepo) selectAllCampaigns() []*types.CampaignRecord {
	tableq := fmt.Sprintf("SELECT * FROM %s", campaignTable)
	rows, err := s.runQuery(tableq)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processCampaignRows(rows)

	return results
}

func (s SqliteRepo) insertCampaign(c types.CampaignRecord) (int, error) {
	s.logger.Debug("tried to insert campaign")
	stmtString := fmt.Sprintf("INSERT INTO %s(name, recruitment, judge, timekeeping, cadence, created_at, updated_at, last_adventure, api_user_id) values(?, ?, ?, ?, ?, ?, ?, ?, ?) ;", campaignTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	res, err := stmt.Exec(c.Name, c.Recruitment, c.Judge, c.Timekeeping, c.Cadence, c.CreatedAt, c.UpdatedAt, c.LastAdventure, c.ClientId)
	util.CheckErr(err)
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s SqliteRepo) selectCampaignById(id int) (*types.CampaignRecord, error) {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", campaignTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processCampaignRows(rows)

	return results[0], nil
}

func (s SqliteRepo) CreateCampaign(c *types.CampaignRecord) (*types.CampaignRecord, error) {
	id, err := s.insertCampaign(*c)
	util.CheckErr(err)
	return s.selectCampaignById(id)
}

func (s SqliteRepo) UpdateCampaign(c *types.CampaignRecord) (*types.CampaignRecord, error) {
	stmtString := fmt.Sprintf("UPDATE %s SET name=?, recruitment=?, judge=?, timekeeping=?, cadence=?, created_at=?, updated_at=? WHERE id=?;", campaignTable)
	s.logger.Print(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(c.Name, c.Recruitment, c.Judge, c.Timekeeping, c.Cadence, c.CreatedAt, c.UpdatedAt, c.ID)
	if err != nil {
		return nil, err
	}
	return s.selectCampaignById(c.ID)

}

func (s SqliteRepo) GetCampaign(id int) (*types.CampaignRecord, error) {
	return s.selectCampaignById(id)
}

func (s SqliteRepo) ListCampaigns() ([]*types.CampaignRecord, error) {
	return s.selectAllCampaigns(), nil
}

func (s SqliteRepo) ListCampaignsForClient(clientId string) ([]*types.CampaignRecord, error) {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.api_user_id = ?", campaignTable)
	rows, err := s.runQuery(tableq, clientId)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processCampaignRows(rows)

	return results, nil

}
func (s SqliteRepo) DeleteCampaign(c *types.CampaignRecord) (bool, error) {
	return false, util.NotYetImplmented()
}

// Adventures
func (s SqliteRepo) insertAdventureRecord(a types.CreateAdventureRequest) (int, error) {
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
func (s SqliteRepo) selectAdventureById(id int) *types.AdventureRecord {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", adventureTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processAdventureRows(rows)

	if len(results) == 0 {
		return nil
	}
	return results[0]
}
func (s SqliteRepo) selectAdventureByCampaignId(c *types.CampaignRecord) ([]*types.AdventureRecord, error) {

	tableq := fmt.Sprintf("SELECT * FROM %s c where c.campaign_id = ?", adventureTable)
	rows, err := s.runQuery(tableq, c.ID)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processAdventureRows(rows)

	return results, nil
}
func (s SqliteRepo) processAdventureRows(r *sql.Rows) []*types.AdventureRecord {
	adventures := make([]*types.AdventureRecord, 0)
	for r.Next() {
		var id int
		var campaignId int
		var duration int
		var name string
		var createdDate time.Time
		var updatedDate time.Time
		var adventureDate time.Time
		copper, silver, electrum, gold, platinum := 0, 0, 0, 0, 0
		err := r.Scan(&id, &campaignId, &name, &adventureDate, &createdDate, &updatedDate, &copper, &silver, &electrum, &gold, &platinum, &duration)

		currentCoins := types.NewCoins(copper, silver, electrum, gold, platinum)
		currentGems := s.getGemsForAdventure(id)
		currentJewellery := s.getJewelleryForAdventure(id)
		currentMagicItems := s.getMagicItemsForAdventure(id)
		currentCombat := s.getCombatForAdventure(id)
		currentCharacters := s.getCharactersForAdventure(id)
		current := types.NewAdventureRecord(id, campaignId, duration, *currentCoins, currentGems, currentJewellery, currentCombat, currentMagicItems, currentCharacters, name, createdDate, updatedDate, adventureDate)
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
		id, actualValue := 0, 0
		name, desc := "", ""
		value := 0.0
		err := rows.Scan(&id, &trashInt, &name, &desc, &value, &actualValue)
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

func (s SqliteRepo) CreateAdventureRecordForCampaign(a *types.CreateAdventureRequest) (*types.AdventureRecord, error) {
	id, err := s.insertAdventureRecord(*a)
	util.CheckErr(err)
	return s.selectAdventureById(id), nil
}

func (s SqliteRepo) GetAdventureRecordsForCampaign(c *types.CampaignRecord) ([]*types.AdventureRecord, error) {
	return s.selectAdventureByCampaignId(c)
}

func (s SqliteRepo) GetAdventureRecordById(a *types.AdventureRecord) (*types.AdventureRecord, error) {
	return s.selectAdventureById(a.ID), nil
}

func (s SqliteRepo) UpdateCoinsForAdventure(a *types.AdventureRecord, c *types.Coins) (bool, error) {
	stmt_string := fmt.Sprintf("UPDATE %s SET copper = ?, silver = ?, electrum = ?, gold = ?, platinum = ? WHERE id=?;", adventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(c.Copper.NumberOfItem, c.Silver.NumberOfItem, c.Electrum.NumberOfItem, c.Gold.NumberOfItem, c.Platinum.NumberOfItem, a.ID)
	if resErr != nil {
		return false, resErr
	}
	return true, nil

}

func (s SqliteRepo) DeleteGemsForAdventure(a *types.AdventureRecord) error {
	stmt_string := fmt.Sprintf("DELETE FROM %s WHERE adventure_id=?;", gemTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID)
	if resErr != nil {
		return resErr
	}
	return nil
}
func (s SqliteRepo) DeleteJewelleryForAdventure(a *types.AdventureRecord) error {
	stmt_string := fmt.Sprintf("DELETE FROM %s WHERE adventure_id=?;", jewelleryTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID)
	if resErr != nil {
		return resErr
	}
	return nil
}
func (s SqliteRepo) DeleteMagicItemsForAdventure(a *types.AdventureRecord) error {
	stmt_string := fmt.Sprintf("DELETE FROM %s WHERE adventure_id=?;", magicItemsTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID)
	if resErr != nil {
		return resErr
	}
	return nil
}
func (s SqliteRepo) DeleteCombatForAdventure(a *types.AdventureRecord) error {
	stmt_string := fmt.Sprintf("DELETE FROM %s WHERE adventure_id=?;", combatTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID)
	if resErr != nil {
		return resErr
	}
	return nil
}
func (s SqliteRepo) DeleteCharactersForAdventure(a *types.AdventureRecord) error {
	stmt_string := fmt.Sprintf("DELETE FROM %s WHERE adventure_id=?;", characterToAdventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID)
	if resErr != nil {
		return resErr
	}
	return nil
}

func (s SqliteRepo) AddGemToAdventure(a *types.AdventureRecord, g *types.Gem) (bool, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?, ?, ?, ?, ?)", gemTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID, g.Loot.Name, g.Loot.Description, g.Loot.XPValueOfOne, g.Loot.NumberOfItem)
	if resErr != nil {
		return false, resErr
	}
	return true, nil
}

func (s SqliteRepo) AddJewelleryToAdventure(a *types.AdventureRecord, g *types.Jewellery) (bool, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(adventure_id, name, description, value, total) values(?, ?, ?, ?, ?)", jewelleryTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID, g.Loot.Name, g.Loot.Description, g.Loot.XPValueOfOne, g.Loot.NumberOfItem)
	if resErr != nil {
		return false, resErr
	}
	return true, nil
}

func (s SqliteRepo) AddMagicItemToAdventure(a *types.AdventureRecord, g *types.MagicItem) (bool, error) {
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

func (s SqliteRepo) AddCombatToAdventure(a *types.AdventureRecord, g *types.MonsterGroup) (bool, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(adventure_id, monster_name, number_defeated, xp_per_monster, total_xp) values(?, ?, ?, ?, ?)", combatTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.ID, g.XP.Name, g.XP.XPValueOfOne, g.XP.XPValueOfOne, g.XP.TotalXPAmount())
	if resErr != nil {
		return false, resErr
	}
	return true, nil
}

// Characters
func (s SqliteRepo) insertCharacter(campaignId int, char types.Character) (int, error) {
	s.logger.Debug("tried to insert characer")
	stmtString := fmt.Sprintf("INSERT INTO %s(campaign_id, name, current_xp, prime_req_percent, character_level, character_class, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?) ;", characterTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	name, xp, primReq, level, class := char.GenerateInsertAttributes()
	res, err := stmt.Exec(campaignId, name, xp, primReq, level, class, time.Now(), time.Now())
	util.CheckErr(err)
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil

}

func (s SqliteRepo) selectCharacterById(id int) *types.CharacterRecord {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", characterTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processCharacterRows(rows)

	return results[0]
}

func (s SqliteRepo) updateCharacter(char types.Character) error {
	tableq := char.GenerateUpdateStatement()
	stmt, err := s.db.Prepare(tableq)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(char.GenerateUpdateAttributes())
	if err != nil {
		return err
	}
	return nil

}

func (s SqliteRepo) getCampaignCharacters(campaignId int) ([]types.CharacterRecord, error) {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", characterTable)
	rows, err := s.runQuery(tableq, campaignId)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processCharacterRowsForCampaign(rows)

	return results, nil
}

func (s SqliteRepo) processCharacterRowsForCampaign(r *sql.Rows) []types.CharacterRecord {
	characters := make([]types.CharacterRecord, 0)
	for r.Next() {
		name, class := "", ""
		id, xp, primeReq, level := 0, 0, 0, 0
		err := r.Scan(&id, &trashInt, &name, &xp, &primeReq, &level, &class, &trashDate, &trashDate)
		util.CheckErr(err)
		current := types.NewCharacter(id, xp, primeReq, level, name, class)
		characters = append(characters, *current)
	}
	return characters
}

func (s SqliteRepo) processCharacterRows(r *sql.Rows) []*types.CharacterRecord {
	characters := make([]*types.CharacterRecord, 0)
	for r.Next() {
		name, class := "", ""
		id, xp, primeReq, level := 0, 0, 0, 0
		err := r.Scan(&id, &trashInt, &name, &xp, &primeReq, &level, &class, &trashDate, &trashDate)
		util.CheckErr(err)
		current := types.NewCharacter(id, xp, primeReq, level, name, class)
		characters = append(characters, current)

	}
	return characters
}

func (s SqliteRepo) CreateCharacterForCampaign(campaign *types.CampaignRecord, character types.Character) (*types.CharacterRecord, error) {
	id, err := s.insertCharacter(campaign.ID, character)
	if err != nil {
		return nil, err
	}
	return s.selectCharacterById(id), nil
}

func (s SqliteRepo) GetCharacterById(char types.CharacterRecord) *types.CharacterRecord {
	return s.selectCharacterById(char.Id())
}

func (s SqliteRepo) UpdateCharacter(character types.Character) (*types.CharacterRecord, error) {
	err := s.updateCharacter(character)
	if err != nil {
		return nil, err
	}
	return s.selectCharacterById(character.Id()), nil
}
func (s SqliteRepo) AddCharacterToAdventure(ad *types.AdventureRecord, char *types.AdventureCharacter) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share) values(?, ?, ?, ?) ;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	if char.Halfshare {
		_, err := stmt.Exec(ad.ID, char.Details.Id(), 1)
		if err != nil {
			return false, err
		}
	} else {
		_, err := stmt.Exec(ad.ID, char.Details.Id(), 0)
		if err != nil {
			return false, err
		}
	}
	return true, nil

}
func (s SqliteRepo) AddHalfshareCharacterToAdventure(ad *types.AdventureRecord, char *types.AdventureCharacter, shareAmount int) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share) values(?, ?, ?, ?) ;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	_, err = stmt.Exec(ad.ID, char.Details.Id(), 1, shareAmount)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (s SqliteRepo) AddFullshareCharacterToAdventure(ad *types.AdventureRecord, char *types.AdventureCharacter, shareAmount int) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share) values(?, ?, ?, ?) ;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	_, err = stmt.Exec(ad.ID, char.Details.Id(), 0, shareAmount)
	if err != nil {
		return false, err
	}
	return true, nil

}
func (s SqliteRepo) RemoveCharacterFromAdventure(ad *types.AdventureRecord, char *types.CharacterRecord) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("DELETE FROM %s WHERE adventure_id = ? AND character_id = ?;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	_, err = stmt.Exec(ad.ID, char.Id())
	if err != nil {
		return false, err
	}

	return true, nil
}
func (s SqliteRepo) ChangeCharacterShares(ad *types.AdventureRecord, char *types.CharacterRecord, isGettingHalfShare bool) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("UPDATE %s set half_share=? WHERE adventure_id = ? AND character_id = ?;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	if isGettingHalfShare {
		_, err = stmt.Exec(ad.ID, char.Id(), 1)
		if err != nil {
			return false, err
		}
	} else {
		_, err = stmt.Exec(ad.ID, char.Id(), 0)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s SqliteRepo) GetCharactersForCampaign(camp *types.CampaignRecord) ([]types.CharacterRecord, error) {
	return s.getCampaignCharacters(camp.ID)
}

// Users
func (s SqliteRepo) GetApiUserById(providedClientId, providedAPIKey string) (*types.APIUser, error) {
	tableq := fmt.Sprintf("SELECT * FROM %s u where u.id = ?", apiUserTable)
	s.logger.Debug(tableq)
	s.logger.Debug(fmt.Printf("Looking for id: %s", providedClientId))
	rows, err := s.runQuery(tableq, providedClientId)
	util.CheckErr(err)
	defer rows.Close()
	users := make([]*types.APIUser, 0)
	for rows.Next() {
		s.logger.Debug("Trying to get user")
		id, hashedKey, name, salt := "", "", "", ""
		err := rows.Scan(&id, &hashedKey, &name, &trashInt, &salt)
		util.CheckErr(err)
		current, err := types.LoadUser(id, providedAPIKey, name, hashedKey, salt)
		if err != nil {
			return nil, err
		}
		users = append(users, current)
	}
	return users[0], nil

}

func (s SqliteRepo) SaveApiUser(user types.User, campaignNumberLimited bool) error {
	stmtString := fmt.Sprintf("INSERT INTO %s(id,api_key,friendly_name,campaign_number_limited,salt) values(?,?,?,?,?);", apiUserTable)
	stmt, err := s.db.Prepare(stmtString)
	if err != nil {
		return err
	}

	if campaignNumberLimited {
		_, err := stmt.Exec(user.DisplayUUID(), user.RetreiveHash(), user.DisplayUserName(), 1, user.RetreiveSalt())
		if err != nil {
			return err
		}
	} else {
		_, err := stmt.Exec(user.DisplayUUID(), user.RetreiveHash(), user.DisplayUserName(), 1, user.RetreiveSalt())
		if err != nil {
			return err
		}
	}
	return nil
}
