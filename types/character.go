package types

import "math"

// comment
type Character interface {
	SQLLiteExportable
	APIObject
	GenerateUpdateAttributes() (string, int, int, string)
	GenerateInsertAttributes() (name string, currentXP int, primeReq int, level int, class string)
}

type CharacterRecord struct {
	id              int
	name            string
	CurrentXP       int
	PrimeReqPercent int
	level           int
	class           string
}

func (c CharacterRecord) GenerateInsertAttributes() (name string, currentXP int, primeReq int, level int, class string) {
	return c.name, c.CurrentXP, c.PrimeReqPercent, c.level, c.class
}
func (c CharacterRecord) GenerateUpdateAttributes() (string, int, int, string) {
	return c.name, c.PrimeReqPercent, c.level, c.class
}

func (c CharacterRecord) GenerateUpdateStatement() string {
	return ""
}

func (c CharacterRecord) Id() int {
	return c.id
}

type characterAPIResponse struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	CurrentXP       int    `json:"current_xp"`
	PrimeReqPercent int    `json:"prime_req"`
	Level           int    `json:"level"`
	Class           string `json:"class"`
}

func (c CharacterRecord) GenerateSuccessfulCreationJSON() APIResponse {
	return characterAPIResponse{
		Id:              c.id,
		Name:            c.name,
		CurrentXP:       c.CurrentXP,
		Level:           c.level,
		PrimeReqPercent: c.PrimeReqPercent,
		Class:           c.class,
	}

}
func NewCharacter(id, currentXp, primeReq, level int, name, class string) *CharacterRecord {
	return &CharacterRecord{
		id:              id,
		name:            name,
		CurrentXP:       currentXp,
		PrimeReqPercent: primeReq,
		level:           level,
		class:           class,
	}
}

func BlankCharacter() *CharacterRecord {
	return NewCharacter(-1, 0, 0, 1, "", "")
}

func NewCharacterById(id int) *CharacterRecord {
	char := BlankCharacter()
	char.id = id
	return char
}

func NewCharacterFromCreateRequest(id int, req CreateCharacterRecordRequest) *CharacterRecord {
	return NewCharacter(-1, 0, req.PrimeReqPercent, req.Level, req.Name, req.Class)
}
func NewCharacterFromUpdateRequest(id int, req UpdateCharacterRecordRequest) *CharacterRecord {
	return NewCharacter(id, 0, req.PrimeReqPercent, req.Level, req.Name, req.Class)
}

type CreateCharacterRecordRequest struct {
	Name            string `json:"name"`
	PrimeReqPercent int    `json:"percent"`
	Level           int    `json:"level"`
	Class           string `json:"class"`
}
type UpdateCharacterRecordRequest struct {
	Name            string `json:"name"`
	PrimeReqPercent int    `json:"percent"`
	Level           int    `json:"level"`
	Class           string `json:"class"`
	XpGained        int    `json:"xp_gained"`
}

type AdventureCharacter struct {
	Details   CharacterRecord
	Halfshare bool `json:"halfshare"`
	XpGained  int  `json:"xp_gained"`
}

type UpdateAdventureCharacter struct {
	ID        int  `json:"id"`
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

func (c *CharacterRecord) AddXP(xpGained int) {
	adjustedXPAmount := math.RoundToEven(float64(xpGained) + (float64(xpGained) * (float64(c.PrimeReqPercent) / 100)))
	c.CurrentXP += int(adjustedXPAmount)
}
