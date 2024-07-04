package types

import "math"

// comment
type Character interface {
	GenerateUpdateAttributes() (string, int, int, string)
	GenerateUpdateStatement() string
	GenerateInsertAttributes() (name string, currentXP int, primeReq int, level int, class string)
}

type CharacterRecord struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	CurrentXP       int    `json:"current_xp"`
	PrimeReqPercent int    `json:"prime_req"`
	Level           int    `json:"level"`
	Class           string `json:"class"`
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

func (c CharacterRecord) Id() int {
	return c.ID
}

func NewCharacter(id, currentXp, primeReq, level int, name, class string) *CharacterRecord {
	return &CharacterRecord{
		ID:              id,
		Name:            name,
		CurrentXP:       currentXp,
		PrimeReqPercent: primeReq,
		Level:           level,
		Class:           class,
	}
}

func BlankCharacter() *CharacterRecord {
	return NewCharacter(-1, 0, 0, 1, "", "")
}

func NewCharacterById(id int) *CharacterRecord {
	char := BlankCharacter()
	char.ID = id
	return char
}
func (c *CharacterRecord) AddXP(xpGained int) {
	adjustedXPAmount := math.RoundToEven(float64(xpGained) + (float64(xpGained) * (float64(c.PrimeReqPercent) / 100)))
	c.CurrentXP += int(adjustedXPAmount)
}
func (c CharacterRecord) ApplyPrimeReq(xpGained int) int {
	adjustedXPAmount := math.RoundToEven(float64(xpGained) + (float64(xpGained) * (float64(c.PrimeReqPercent) / 100)))
	return int(adjustedXPAmount)
}

type AdventureCharacter struct {
	Details   CharacterRecord
	Halfshare bool `json:"halfshare"`
	XpGained  int  `json:"xp_gained"`
}

func NewAdventureCharacter(details *CharacterRecord, halfshare bool, xp int) *AdventureCharacter {
	return &AdventureCharacter{
		Details:   *details,
		Halfshare: halfshare,
		XpGained:  xp,
	}
}
