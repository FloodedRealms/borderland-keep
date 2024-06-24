package types

type CampaignActivity struct {
	Id                         int    `json:"id"`
	Character_id               int    `json:"character_id"`
	Name                       string `json:"name"`
	LevelAtWhichActionWasTaken int    `json:"action_level"`
	XPgained                   int    `json:"xp_gained"`
	XPThreshold                int    `json:"xp_threshold"`
}
