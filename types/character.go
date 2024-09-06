package types

import "math"

// comment
type Character interface {
	GenerateUpdateAttributes() (string, int, int, string)
	GenerateUpdateStatement() string
	GenerateInsertAttributes() (name string, currentXP int, primeReq int, level int, class string)
}

type CharacterRecord struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	CurrentXP       int    `json:"current_xp"`
	PrimeReqPercent int    `json:"prime_req"`
	Level           int    `json:"level"`
	Class           string `json:"class"`
	ClassId         int
	Status          string
	StatusId        int
	ShowAdjustedXP  func(int) int
}

func (c CharacterRecord) GenerateInsertAttributes() (name string, currentXP int, primeReq int, level int, class string) {
	return c.Name, c.CurrentXP, c.PrimeReqPercent, c.Level, c.Class
}
func (c CharacterRecord) GenerateUpdateAttributes() (string, int, int, string) {
	return c.Name, c.PrimeReqPercent, c.Level, c.Class
}

func (c CharacterRecord) GenerateUpdateStatement() string {
	return ""
}

func NewCharacter(id, currentXp, primeReq, level int, name, class string) *CharacterRecord {
	c := &CharacterRecord{
		Id:              id,
		Name:            name,
		CurrentXP:       currentXp,
		PrimeReqPercent: primeReq,
		Level:           level,
		Class:           class,
		ShowAdjustedXP:  CreateXPFunc(primeReq),
	}
	return c
}

func BlankCharacter() *CharacterRecord {
	return NewCharacter(-1, 0, 0, 1, "", "")
}

func NewCharacterById(id int) *CharacterRecord {
	char := BlankCharacter()
	char.Id = id
	return char
}
func (c *CharacterRecord) AddXP(xpGained int) {
	adjustedXPAmount := math.RoundToEven(float64(xpGained) + (float64(xpGained) * (float64(c.PrimeReqPercent) / 100)))
	c.CurrentXP += int(adjustedXPAmount)
}

func CreateXPFunc(primReq int) func(int) int {
	return func(xpGained int) int {
		adjustedXPAmount := math.RoundToEven(float64(xpGained) + (float64(xpGained) * (float64(primReq) / 100)))
		return int(adjustedXPAmount)
	}
}
func (c CharacterRecord) ApplyPrimeReq(xpGained int) int {
	adjustedXPAmount := math.RoundToEven(float64(xpGained) + (float64(xpGained) * (float64(c.PrimeReqPercent) / 100)))
	return int(adjustedXPAmount)
}

type AdventureCharacter struct {
	Id             int  `json:"id"`
	Halfshare      bool `json:"halfshare"`
	Name           string
	Preq           int
	ShowAdjustedXP func(int) int
}

func (a *AdventureCharacter) CreateXPFunc() {
	a.ShowAdjustedXP = CreateXPFunc(a.Preq)
}

func NewAdventureCharacter(halfshare bool, id int) *AdventureCharacter {
	return &AdventureCharacter{
		Id:        id,
		Halfshare: halfshare,
	}
}
