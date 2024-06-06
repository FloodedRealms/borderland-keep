package repository

import (
	"database/sql"
	"fmt"

	"github.com/kevin/adventure-archivist/types"
	"github.com/kevin/adventure-archivist/util"
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

type SqliteRepo struct {
	db *sql.DB
}

func NewSqliteRepo(f string) (*SqliteRepo, error) {
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
	return &SqliteRepo{db: db}, nil
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
		err := r.Scan(current.ID, current.Name, current.Recruitment, current.Judge, current.Timekeeping, current.Cadence, current.Cadence, current.UpdatedAt, current.LastAdventure)
		util.CheckErr(err)
		campaigns = append(campaigns, current)
	}
	return campaigns
}

func (s SqliteRepo) selectAllCampaigns() []*types.Campaign {
	tableq := fmt.Sprintf("SELECT * FROM %s", campaignTable)
	rows, err := s.runQuery(tableq)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processCampaignRows(rows)

	return results
}

func (s SqliteRepo) insertCampaign(c types.CreateCampaignRequest) (int, error) {
	stmt, err := s.db.Prepare(fmt.Sprintf("INSERT INTO %s (name, recruitment, judge, timekeeping, cadence, created_at, updated_at, last_adventure) (?, ?, ?, ?, ?, ?, ?, ?)", campaignTable))
	util.CheckErr(err)
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
	stmt, err := s.db.Prepare(fmt.Sprintf("INSERT INTO %s (campaign_id, name, adventure_date, copper, silver, electrum, gold, platinum created_at, updated_at) (?, ?, ?, ?, ?)", adventureTable))
	util.CheckErr(err)
	res, err := stmt.Exec(a.CampaignID, a.Name, a.AdventureDate, a.Copper, a.Silver, a.Electrum, a.Gold, a.Platinum, a.CreatedAt, a.UpdatedAt)
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
	defer rows.Close()
	util.CheckErr(err)

	results := s.processAdventureRows(rows)

	return results[0]
}
func (s SqliteRepo) selectAdventureByCampaignId(c *types.Campaign) ([]*types.AdventureRecord, error) {

	tableq := fmt.Sprintf("SELECT * FROM %s c where c.campaign_id = ?", adventureTable)
	rows, err := s.runQuery(tableq, c.ID)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processAdventureRows(rows)

	return results, nil
}
func (s SqliteRepo) processAdventureRows(r *sql.Rows) []*types.AdventureRecord {
	adventures := make([]*types.AdventureRecord, 0)
	for r.Next() {
		current := &types.AdventureRecord{}
		err := r.Scan(current.ID, current.CampaignID, current.Name, current.AdventureDate, current.Copper, current.Silver, current.Electrum, current.Gold, current.Electrum)
		current.Gems = s.getGemsForAdventure(current.ID)
		current.Jewellery = s.getJewelleryForAdventure(current.ID)
		current.MagicItems = s.getMagicItemsForAdventure(current.ID)
		current.Combat = s.getCombatForAdventure(current.ID)
		current.Characters = s.getCharactersForAdventure(current.ID)
		util.CheckErr(err)
		adventures = append(adventures, current)
	}
	return adventures
}

func (s SqliteRepo) getGemsForAdventure(id int) []*types.Gem {
	tableq := fmt.Sprintf("SELECT * FROM %s g where c.adventure_id = ?", gemTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)
	gems := make([]*types.Gem, 0)
	for rows.Next() {
		current := &types.Gem{}
		err := rows.Scan(nil, nil, current.Name, current.Description, current.Value, current.Total)
		util.CheckErr(err)
		gems = append(gems, current)
	}
	return gems
}
func (s SqliteRepo) getJewelleryForAdventure(id int) []*types.Jewellery {
	tableq := fmt.Sprintf("SELECT * FROM %s g where c.adventure_id = ?", jewelleryTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)
	jewellery := make([]*types.Jewellery, 0)
	for rows.Next() {
		current := &types.Jewellery{}
		err := rows.Scan(nil, nil, current.Name, current.Description, current.Value, current.Total)
		util.CheckErr(err)
		jewellery = append(jewellery, current)
	}
	return jewellery
}

func (s SqliteRepo) getMagicItemsForAdventure(id int) []*types.MagicItem {
	tableq := fmt.Sprintf("SELECT * FROM %s g where c.adventure_id = ?", magicItemsTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)
	magicItems := make([]*types.MagicItem, 0)
	for rows.Next() {
		current := &types.MagicItem{}
		err := rows.Scan(nil, nil, current.Name, current.Description, current.ApparentValue)
		util.CheckErr(err)
		magicItems = append(magicItems, current)
	}
	return magicItems
}

func (s SqliteRepo) getCombatForAdventure(id int) []*types.MonsterGroup {
	tableq := fmt.Sprintf("SELECT * FROM %s g where c.adventure_id = ?", combatTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)
	combat := make([]*types.MonsterGroup, 0)
	for rows.Next() {
		current := &types.MonsterGroup{}
		err := rows.Scan(nil, nil, current.MonsterName, current.NumberDefeated, current.XP, current.TotalXP)
		util.CheckErr(err)
		combat = append(combat, current)
	}
	return combat
}
func (s SqliteRepo) getCharactersForAdventure(id int) []*types.Character {
	tableq := fmt.Sprintf("SELECT * FROM %s g where c.adventure_id = ?", characterToAdventureTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)
	characters := make([]*types.Character, 0)
	for rows.Next() {
		var rid int
		var aid int
		var chid int
		var halfStatus int
		err := rows.Scan(&rid, &aid, &chid, &halfStatus)
		current := s.selectCharacterById(chid)
		current.Halfshare = false
		if halfStatus == 1 {
			current.Halfshare = true
		}
		util.CheckErr(err)
		characters = append(characters, current)
	}
	return characters

}

func (s SqliteRepo) CreateAdventureRecordForCampaign(a *types.CreateAdventureRecordRequest) (*types.AdventureRecord, error) {
	id, err := s.insertAdventureRecord(*a)
	util.CheckErr(err)
	return s.selectAdventureById(id), nil
}

func (s SqliteRepo) GetAdventureRecordsForCampaign(c *types.Campaign) ([]*types.AdventureRecord, error) {
	return s.selectAdventureByCampaignId(c)
}

// Characters
func (s SqliteRepo) selectCharacterById(id int) *types.Character {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", characterTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processCharacterRows(rows)

	return results[0]
}

func (s SqliteRepo) processCharacterRows(r *sql.Rows) []*types.Character {
	characters := make([]*types.Character, 0)
	for r.Next() {
		current := &types.Character{}
		err := r.Scan(current.ID, nil, current.Name, current.Name, current.PrimeReqPercent, current.Level, current.Class, nil, nil)
		util.CheckErr(err)
		characters = append(characters, current)
	}
	return characters
}
