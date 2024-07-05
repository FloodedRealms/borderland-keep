package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/floodedrealms/adventure-archivist/internal/util"
	"github.com/floodedrealms/adventure-archivist/types"
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
const campaignActivityTable string = "campaign_activities"

var trashInt int = 0
var trashDate time.Time = time.Now()
var trashString string = ""

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
func (s SqliteRepo) ExecuteQuery(q string, params ...interface{}) (sql.Result, error) {
	stmt, err := s.db.Prepare(q)
	util.CheckErr(err)
	return stmt.Exec(params...)
}

// Campaigns
func (s SqliteRepo) processCampaignRows(r *sql.Rows) []*types.CampaignRecord {
	campaigns := make([]*types.CampaignRecord, 0)
	for r.Next() {
		var (
			current *types.CampaignRecord
		)
		err := r.Scan(&current.Id, &current.Name, &current.Recruitment, &current.Judge, &current.Timekeeping, &current.Cadence, &current.CreatedAt, &current.UpdatedAt, &current.LastAdventure, &current.ClientId, &trashInt, &trashString, &trashString)
		util.CheckErr(err)
		current.Characters, err = s.getCampaignCharacters(current.Id)
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
	stmtString := fmt.Sprintf("INSERT INTO %s(name, recruitment, judge, timekeeping, cadence, created_at, updated_at, last_adventure, api_user_id, password, salt) values(?, ?, ?, ?, ?, ?, ?, ?, ?, \"\", \"\") ;", campaignTable)
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

func (s SqliteRepo) UpdateCampaignPassword(id int, pass types.Password) error {
	stmtString := fmt.Sprintf("UPDATE %s set password=?, salt=? WHERE id=?;", campaignTable)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	_, err = stmt.Exec(pass.Hash.Hash, pass.Hash.Salt, id)
	return nil
}

func (s SqliteRepo) selectCampaignById(id int) (*types.CampaignRecord, error) {
	tableq := fmt.Sprintf("SELECT c.*, a.id, a.name, a.adventure_date FROM %s c JOIN %s a ON a.campaign_id = c.id where c.id =?;", campaignTable, adventureTable)
	//tableq1 := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", campaignTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()
	var (
		campaignRows []*types.CampaignRecord
		adventures   []types.AdventureRecord
	)
	for rows.Next() {
		var (
			current types.CampaignRecord
			adId    int
			adName  string
			aDate   time.Time
		)
		err := rows.Scan(&current.Id, &current.Name, &current.Recruitment, &current.Judge, &current.Timekeeping, &current.Cadence, &current.CreatedAt, &current.UpdatedAt, &current.LastAdventure, &current.ClientId, &trashInt, &trashString, &trashString, &adId, &adName, &aDate)
		util.CheckErr(err)
		adventures = append(adventures, types.AdventureRecord{Id: adId, Name: adName, AdventureDate: types.ArcvhistDate(aDate)})
		current.Characters, err = s.getCampaignCharacters(current.Id)
		util.CheckErr(err)

		campaignRows = append(campaignRows, &current)
	}
	result := campaignRows[0]
	result.Adventures = adventures
	return result, nil
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
	_, err = stmt.Exec(c.Name, c.Recruitment, c.Judge, c.Timekeeping, c.Cadence, c.CreatedAt, c.UpdatedAt, c.Id)
	if err != nil {
		return nil, err
	}
	return s.selectCampaignById(c.Id)

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
func (s SqliteRepo) insertAdventureRecord(a types.AdventureRecord) (int, error) {
	stmt_string := fmt.Sprintf("INSERT INTO %s(campaign_id, name, created_at, updated_at, adventure_date) values(?, ?, ?, ?, ?)", adventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	res, err := stmt.Exec(a.CampaignId, a.Name, time.Now(), time.Now(), a.AdventureDate.Date())
	util.CheckErr(err)
	stmt2, err := s.db.Prepare(fmt.Sprintf("UPDATE %s SET last_adventure=? WHERE id=?", campaignTable))
	util.CheckErr(err)
	_, err = stmt2.Exec(a.AdventureDate.Date(), a.CampaignId)
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
	rows, err := s.runQuery(tableq, c.Id)
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
		currentCharacters, _ := s.getCharactersForAdventure(id)
		current := types.NewAdventureRecord(id, campaignId, duration, *currentCoins, currentGems, currentJewellery, currentCombat, currentMagicItems, currentCharacters, name, types.ArcvhistDate(adventureDate))
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
		current := types.NewGem(number, name, desc, value, value)
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
		current := types.NewJewellery(number, name, desc, value, value)

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
		current := types.NewMagicItem(name, desc, int(value), int(actualValue))
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
		current := types.NewMonsterGroup(name, name, number, int(value))

		combat = append(combat, *current)
	}
	return combat
}
func (s SqliteRepo) getCharactersForAdventure(id int) ([]types.AdventureCharacter, error) {
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
		var xp int
		err := rows.Scan(&rowId, &adventureId, &characterId, &halfStatus, &xp)
		util.CheckErr(err)
		halfshare := false
		if halfStatus == 1 {
			halfshare = true
		}
		currentDetails, err := s.selectCharacterById(characterId)
		if err != nil {
			return nil, err
		}
		adventureCharacter := types.NewAdventureCharacter(halfshare, currentDetails.Id)
		characters = append(characters, *adventureCharacter)
	}
	return characters, nil

}

func (s SqliteRepo) CreateAdventureRecordForCampaign(a *types.AdventureRecord) (*types.AdventureRecord, error) {
	id, err := s.insertAdventureRecord(*a)
	util.CheckErr(err)
	return s.selectAdventureById(id), nil
}

func (s SqliteRepo) GetAdventureRecordsForCampaign(c *types.CampaignRecord) ([]*types.AdventureRecord, error) {
	return s.selectAdventureByCampaignId(c)
}

func (s SqliteRepo) GetAdventureRecordById(a *types.AdventureRecord) (*types.AdventureRecord, error) {
	return s.selectAdventureById(a.Id), nil
}

func (s SqliteRepo) UpdateCoinsForAdventure(a *types.AdventureRecord, c *types.Coins) (bool, error) {
	stmt_string := fmt.Sprintf("UPDATE %s SET copper = ?, silver = ?, electrum = ?, gold = ?, platinum = ? WHERE id=?;", adventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(c.Copper.Number, c.Silver.Number, c.Electrum.Number, c.Gold.Number, c.Platinum.Number, a.Id)
	if resErr != nil {
		return false, resErr
	}
	return true, nil

}
func (s SqliteRepo) UpdateAdventureName(a *types.AdventureRecord, n string) error {
	stmt_string := fmt.Sprintf("UPDATE %s SET name = ? WHERE id=?;", adventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(n, a.Id)
	if resErr != nil {
		return resErr
	}
	return nil
}
func (s SqliteRepo) UpdateCharacterTotalXP(c types.CharacterRecord) error {
	stmt_string := fmt.Sprintf("UPDATE %s SET total_xp = ? WHERE id=?;", characterTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(c.CurrentXP, c.Id)
	if resErr != nil {
		return resErr
	}
	return nil
}
func (s SqliteRepo) UpdateAdventureDate(a *types.AdventureRecord, d types.ArcvhistDate) error {
	stmt_string := fmt.Sprintf("UPDATE %s SET adventure_date = ? WHERE id=?;", adventureTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(d.Date(), a.Id)
	if resErr != nil {
		return resErr
	}
	return nil
}

func (s SqliteRepo) DeleteGemsForAdventure(a *types.AdventureRecord) error {
	stmt_string := fmt.Sprintf("DELETE FROM %s WHERE adventure_id=?;", gemTable)
	s.logger.Debug(stmt_string)
	stmt, err := s.db.Prepare(stmt_string)
	util.CheckErr(err)
	_, resErr := stmt.Exec(a.Id)
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
	_, resErr := stmt.Exec(a.Id)
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
	_, resErr := stmt.Exec(a.Id)
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
	_, resErr := stmt.Exec(a.Id)
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
	_, resErr := stmt.Exec(a.Id)
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
	_, resErr := stmt.Exec(a.Id, g.Name, g.Description, g.XPValue, g.Number)
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
	_, resErr := stmt.Exec(a.Id, g.Name, g.Description, g.XPValue, g.Number)
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
	_, resErr := stmt.Exec(a.Id, g.Name, g.Description, g.GoldValue, g.XPValue)
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
	_, resErr := stmt.Exec(a.Id, g.Name, g.NumberDefeated, g.XPPerOneKill, g.TotalXPAmount())
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

func (s SqliteRepo) selectCharacterById(id int) (*types.CharacterRecord, error) {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", characterTable)
	rows, err := s.runQuery(tableq, id)
	util.CheckErr(err)
	defer rows.Close()

	results := s.processCharacterRows(rows)

	if len(results) < 1 {
		return nil, errors.New(fmt.Sprintf("unable able to find characted with id %d", id))
	}

	return results[0], nil
}

func (s SqliteRepo) updateCharacter(char types.CharacterRecord) error {
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
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.campaign_id = ?", characterTable)

	rows, err := s.runQuery(tableq, campaignId)
	util.CheckErr(err)
	defer rows.Close()

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
	id, err := s.insertCharacter(campaign.Id, character)
	if err != nil {
		return nil, err
	}
	return s.selectCharacterById(id)
}

func (s SqliteRepo) GetCharacterById(char types.CharacterRecord) (*types.CharacterRecord, error) {
	return s.selectCharacterById(char.Id)
}

func (s SqliteRepo) UpdateCharacter(character types.CharacterRecord) (*types.CharacterRecord, error) {
	err := s.updateCharacter(character)
	if err != nil {
		return nil, err
	}
	return s.selectCharacterById(character.Id)
}
func (s SqliteRepo) AddCharacterToAdventure(ad *types.AdventureRecord, char *types.AdventureCharacter) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share, xp_gained) values(?, ?, ?, ?) ;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	if char.Halfshare {
		_, err := stmt.Exec(ad.Id, char.Id, 1)
		if err != nil {
			return false, err
		}
	} else {
		_, err := stmt.Exec(ad.Id, char.Id, 0)
		if err != nil {
			return false, err
		}
	}
	return true, nil

}
func (s SqliteRepo) AddHalfshareCharacterToAdventure(ad *types.AdventureRecord, char *types.AdventureCharacter, shareAmount int) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share, xp_gained) values(?, ?, ?, ?) ;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	_, err = stmt.Exec(ad.Id, char.Id, 1, shareAmount)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (s SqliteRepo) AddFullshareCharacterToAdventure(ad *types.AdventureRecord, char *types.AdventureCharacter, shareAmount int) (bool, error) {
	s.logger.Debug("tried to insert characer to adventure mapping")
	stmtString := fmt.Sprintf("INSERT INTO %s(adventure_id, character_id, half_share, xp_gained) values(?, ?, ?, ?) ;", characterToAdventureTable)
	s.logger.Debug("String is:")
	s.logger.Debug(stmtString)
	stmt, err := s.db.Prepare(stmtString)
	util.CheckErr(err)
	s.logger.Debug("statement successully prepared.")
	_, err = stmt.Exec(ad.Id, char.Id, 0, shareAmount)
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
	_, err = stmt.Exec(ad.Id, char.Id)
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
		_, err = stmt.Exec(ad.Id, char.Id, 1)
		if err != nil {
			return false, err
		}
	} else {
		_, err = stmt.Exec(ad.Id, char.Id, 0)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s SqliteRepo) GetCharactersForCampaign(camp *types.CampaignRecord) ([]types.CharacterRecord, error) {
	return s.getCampaignCharacters(camp.Id)
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

func (s SqliteRepo) GetCharacterXPGains(c types.CharacterRecord) ([]int, error) {
	aXp, aErr := s.getAdventureXP(c)
	if aErr != nil {
		return nil, aErr
	}
	cXp, cErr := s.getCampaignXP(c)
	if cErr != nil {
		return nil, cErr
	}
	// Check if both have values
	if len(aXp) > 0 && len(cXp) > 0 {
		return append(aXp, cXp...), nil

	} else if len(aXp) > 0 {
		// if either is empty, check if aXp has any values
		return aXp, nil
	} // if not, return cXp which is either an empty array, or the only thing with XP values
	return cXp, nil
}

func (s SqliteRepo) getAdventureXP(c types.CharacterRecord) ([]int, error) {
	tableq := fmt.Sprintf("SELECT u.xp_gained FROM %s u where u.character_id = ?", characterToAdventureTable)

	s.logger.Debug(tableq)
	rows, err := s.runQuery(tableq, c.Id)
	util.CheckErr(err)
	defer rows.Close()
	xpGains := make([]int, 0)
	for rows.Next() {
		xp := 0
		err := rows.Scan(&xp)
		util.CheckErr(err)
		if err != nil {
			return nil, err
		}
		xpGains = append(xpGains, xp)
	}
	return xpGains, nil
}
func (s SqliteRepo) getCampaignXP(c types.CharacterRecord) ([]int, error) {
	tableq := fmt.Sprintf("SELECT u.xp_gained FROM %s u where u.character_id = ?", campaignActivityTable)

	s.logger.Debug(tableq)
	rows, err := s.runQuery(tableq, c.Id)
	util.CheckErr(err)
	defer rows.Close()
	xpGains := make([]int, 0)
	for rows.Next() {
		xp := 0
		err := rows.Scan(&xp)
		util.CheckErr(err)
		if err != nil {
			return nil, err
		}
		xpGains = append(xpGains, xp)
	}
	return xpGains, nil
}
func (s SqliteRepo) GetLevelForXP(camp types.CampaignRecord, c types.CharacterRecord) int {
	tableq := "SELECT clt.xp_level, clt.xp_amount FROM classes cl \n" +
		"INNER JOIN class_level_thresholds clt ON clt.class_id=cl.id WHERE cl.system_id = ? AND cl.class_name = ? AND clt.xp_amount <= ? \n" +
		"ORDER BY clt.xp_amount  DESC;"
	s.logger.Print(tableq)
	rows, err := s.runQuery(tableq, 1, c.Class, c.CurrentXP)
	util.CheckErr(err)
	defer rows.Close()
	// This always returns the current level as the first result
	// so if rows has no results, somethine failed.
	// TODO: Actual error handling and surfacing here
	level := -1
	if !rows.Next() {
		return level
	}
	rows.Scan(&level, &trashInt)
	return level

}
func (s SqliteRepo) AddCampaignActivityForCharacter(a types.CampaignActivity) error {
	stmtString := fmt.Sprintf("INSERT INTO %s(id,character_id,name,taken_at_level,xp_gained) values(?,?,?,?,?);", apiUserTable)
	stmt, err := s.db.Prepare(stmtString)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(a.Id, a.Character_id, a.Name, a.LevelAtWhichActionWasTaken, a.XPgained)
	if err != nil {
		return err
	}

	return nil
}

func (s SqliteRepo) GetCoinsForAdventure(a *types.AdventureRecord) (*types.Coins, error) {
	stmtStr := fmt.Sprintf("SELECT a.copper, a.silver, a.electrum, a.gold, a.platinum FROM %s a WHERE a.id=?", adventureTable)
	rows, err := s.runQuery(stmtStr, a.Id)
	if err != nil {
		return nil, err
	}

	var (
		c  int
		si int
		e  int
		g  int
		p  int
	)
	for rows.Next() {
		rows.Scan(&c, &si, &e, &g, &p)
	}
	return types.NewCoins(c, si, e, g, p), nil
}
