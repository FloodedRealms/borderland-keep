package types

// comment
type Character struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	CurrentXP       int    `json:"xp"`
	PrimeReqPercent int    `json:"percent"`
	Level           int    `json:"level"`
	Class           string `json:"class"`
}

type AdventureCharacter struct {
	Details   Character
	Halfshare bool `json:"halfshare"`
}

func NewCharacter(id, currentXp, primeReq, level int, name, class string) *Character {
	return &Character{
		ID:              id,
		Name:            name,
		CurrentXP:       currentXp,
		PrimeReqPercent: primeReq,
		Level:           level,
		Class:           class,
	}
}

func BlankCharacter() *Character {
	return NewCharacter(-1, 0, 0, 1, "", "")
}

func NewCharacterById(id int) *Character {
	char := BlankCharacter()
	char.ID = id
	return char
}

func NewAdventureCharacter(details *Character, halfshare bool) *AdventureCharacter {
	return &AdventureCharacter{
		Details:   *details,
		Halfshare: halfshare,
	}
}
