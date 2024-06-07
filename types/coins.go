package types

type Coins struct {
	Copper   XPSource `json:"copper"`
	Silver   XPSource `json:"silver"`
	Electrum XPSource `json:"electrum"`
	Gold     XPSource `json:"gold"`
	Platinum XPSource `json:"platinum"`
}

func NewCoins(c, s, e, g, p int) *Coins {
	/* For now, we are hardcoding the ACKS assumptions into the system.
	* That means that each coin demarcation is an XP source with an XPValue equal to its
	* value in gold according to ACKS and a Number equal to the normal coins included in the adventure.
	* These values are:
	* Copper = 1/100
	* Silver = 1/10
	* Electrum = 1/2
	* Gold = 1
	* Platinum = 5

	TODO: Front ends should be able to change this as they desire, i.e. a front end written to handle 7voz would be on the silver standard,
	so silver has an XP value of 1. A front end for d&d 5e would have all coin XPValues at 0, since treasure does not award XP*/
	copper := NewLoot("Coppper Coins", "The meanest of coins, still worth something", 0.001, c)
	silver := NewLoot("Silver Coins", "A little bit of real money", 0.01, s)
	electrum := NewLoot("Electrum Coins", "The best coin, despite my players protests", 0.5, e)
	gold := NewLoot("Gold Coins", "THE standard", 1.0, g)
	platinum := NewLoot("Platinum Coins", "Fancy", 5.0, p)

	return &Coins{
		Copper:   *copper,
		Silver:   *silver,
		Electrum: *electrum,
		Gold:     *gold,
		Platinum: *platinum,
	}
}
