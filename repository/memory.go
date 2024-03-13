package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/kevin/adventure-archivist/types"
)

const dateLayout = "2006-01-02"

type MemoryRepo struct{}

func NewMemoryRepo() (*MemoryRepo, error) {
	return &MemoryRepo{}, nil
}

func (s *MemoryRepo) CreateCampaign(c *types.CreateCampaignRequest) (*types.Campaign, error) {
	return &types.Campaign{
		ID:            1,
		Name:          c.Name,
		Recruitment:   c.Recruitment,
		Judge:         c.Judge,
		Timekeeping:   c.Timekeeping,
		Cadence:       c.Cadence,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		LastAdventure: c.LastAdventure,
	}, nil
}

func (s *MemoryRepo) GetCampaign(id int) (*types.Campaign, error) {
	if id == 1 {
		cat, _ := time.Parse(dateLayout, "2024-01-01")
		uat, _ := time.Parse(dateLayout, "2024-01-01")
		laa, _ := time.Parse(dateLayout, "2024-01-01")
		return &types.Campaign{
			ID:            1,
			Name:          "Test One",
			Recruitment:   false,
			Judge:         "Shock",
			Timekeeping:   "Pause",
			Cadence:       "weekly",
			CreatedAt:     cat,
			UpdatedAt:     uat,
			LastAdventure: laa,
		}, nil

	}
	if id == 2 {
		cat, _ := time.Parse(dateLayout, "2024-03-01")
		uat, _ := time.Parse(dateLayout, "2024-03-01")
		laa, _ := time.Parse(dateLayout, "2024-03-01")
		return &types.Campaign{
			ID:            2,
			Name:          "Test Two",
			Recruitment:   false,
			Judge:         "Shock",
			Timekeeping:   "Reference",
			Cadence:       "open",
			CreatedAt:     cat,
			UpdatedAt:     uat,
			LastAdventure: laa,
		}, nil

	}
	return nil, errors.New(fmt.Sprintf("Unable to retrieve campaign with id %d", id))
}

func (s *MemoryRepo) ListCampaigns() ([]*types.Campaign, error) {
	var ret []*types.Campaign
	c1, _ := s.GetCampaign(1)
	c2, _ := s.GetCampaign(2)
	ret = append(ret, c1)
	ret = append(ret, c2)
	return ret, nil
}

func (s *MemoryRepo) GetAdventureRecordsForCampaign(id int) ([]*types.AdventureRecord, error) {
	if id != 1 {
		return nil, errors.New("campaign id not found")
	}
	startDate, _ := time.Parse(dateLayout, "2024-01-01")
	second := startDate.AddDate(0, 0, 7)
	third := second.AddDate(0, 0, 7)
	a := types.NewAdventureRecord(1, 1, memCoins(), memGems(), memJewellery(), memMagicItem(), memMonsters(), startDate)
	b := types.NewAdventureRecord(2, 1, memCoins(), memGems(), memJewellery(), memMagicItem(), memMonsters(), second)
	c := types.NewAdventureRecord(3, 1, memCoins(), memGems(), memJewellery(), memMagicItem(), memMonsters(), third)

	return []*types.AdventureRecord{a, b, c}, nil

}

func memCoins() *types.Coins {
	return types.NewCoins(10000, 1000, 200, 100, 20)
}
func memGems() []*types.Gem {
	a := types.NewGem("Amethyst", "Shines brightly", 100)
	d := types.NewGem("Diamond", "", 1000)
	return []*types.Gem{a, d}
}
func memJewellery() []*types.Jewellery {
	a := types.NewJewellery("Necklace", "is broken", 100)
	d := types.NewJewellery("Crown", "", 1000)
	return []*types.Jewellery{a, d}
}
func memMagicItem() []*types.MagicItem {
	a := types.NewMagicItem("Rope", "is Rope-like", 1000)
	return []*types.MagicItem{a}
}

func memMonsters() []*types.MonsterGroup {
	a := types.NewMonsterGroup("Orc", 10, 14)
	b := types.NewMonsterGroup("Skeleton", 10, 9)
	return []*types.MonsterGroup{a, b}
}
