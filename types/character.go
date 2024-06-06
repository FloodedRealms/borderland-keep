package types

// comment
type Character struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	CurrentXP       int    `json:"xp"`
	PrimeReqPercent int    `json:"percent"`
	Level           int    `json:"level"`
	Class           string `json:"class"`
	Halfshare       bool   `json:"halfshare"`
}
