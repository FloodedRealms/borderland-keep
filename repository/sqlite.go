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
	return &SqliteRepo{db: db}, nil
}

func (s SqliteRepo) CreateCampaign(c types.Campaign) *types.Campaign {
	id, err := s.insertCampaign(c)
	util.CheckErr(err)
	return s.selectCampaignById(id)
}

func (s SqliteRepo) insertCampaign(c types.Campaign) (int, error) {
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

func (s SqliteRepo) GetCampaign(id int) *types.Campaign {
	return s.selectCampaignById(id)
}

func (s SqliteRepo) selectCampaignById(id int) *types.Campaign {
	tableq := fmt.Sprintf("SELECT * FROM %s c where c.id = ?", campaignTable)
	rows, err := s.runQuery(tableq, id)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processCampaignRows(rows)

	return results[0]
}

func (s SqliteRepo) ListCampaigns() []*types.Campaign {
	return s.selectAllCampaigns()
}
func (s SqliteRepo) selectAllCampaigns() []*types.Campaign {
	tableq := fmt.Sprintf("SELECT * FROM %s", campaignTable)
	rows, err := s.runQuery(tableq)
	defer rows.Close()
	util.CheckErr(err)

	results := s.processCampaignRows(rows)

	return results
}

func (s SqliteRepo) runQuery(q string, params ...interface{}) (*sql.Rows, error) {
	stmt, err := s.db.Prepare(q)
	util.CheckErr(err)
	return stmt.Query(params...)
}

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
