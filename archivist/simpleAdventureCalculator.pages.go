package archivist

import (
	"net/http"
	"strconv"

	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/renderer"
)

type SimpleAdventurePage struct {
	renderer renderer.Renderer
}

type AdventureForm struct {
	NumberCharacters string
	NumberHenchmen   string
	CharactersError  string
	HenchmenError    string
	Copper           string
	Silver           string
	Electrum         string
	Gold             string
	Platinum         string
	CopperError      string
	SilverError      string
	ElectrumError    string
	GoldError        string
	PlatinumError    string
	Loot             []LootItem
	MagicItems       []LootItem
	Combat           []LootItem
	Valid            bool
}

type FormInputError struct {
	Name    string
	Number  string
	XPValue string
	GPValue string
}

type LootItem struct {
	Name    string
	Number  string
	GPValue string
	Errors  FormInputError
	Valid   bool
}

type MagicItem struct {
	Name    string
	XPValue string
	Errors  FormInputError
	Valid   bool
}

type CombatItem struct {
	Name    string
	XPValue string
	Errors  FormInputError
	Valid   bool
}

func ParserAndValidateForm(formData map[string]string) *AdventureForm {
	valid := true
	characterError := ""
	henchmenError := ""
	copperError := ""
	silverError := ""
	electrumError := ""
	goldError := ""
	platinumError := ""

	characterNumberString, ok := formData["num-characters"]
	if !ok {
		characterError = "Character Number is required"
		valid = false
	} else {
		_, err := strconv.Atoi(characterNumberString)
		if err != nil {
			characterError = "Character Number must be numeric"
			valid = false
		}
	}
	henchmenNumberString, ok := formData["num-henchmen"]
	if !ok {
		henchmenError = "Henchmen Number is required"
		valid = false
	} else {
		_, err := strconv.Atoi(henchmenNumberString)
		if err != nil {
			henchmenError = "Henchmen Number must be numeric"
			valid = false
		}
	}
	copper, ok := formData["copper"]
	if ok {
		_, err := strconv.Atoi(copper)
		if err != nil {
			copperError = "Coin amount must be a number"
		}
	}
	silver, ok := formData["silver"]
	if ok {
		_, err := strconv.Atoi(silver)
		if err != nil {
			silverError = "Coin amount must be a number"
		}
	}
	electrum, ok := formData["electrum"]
	if ok {
		_, err := strconv.Atoi(electrum)
		if err != nil {
			electrumError = "Coin amount must be a number"
		}
	}
	gold, ok := formData["gold"]
	if ok {
		_, err := strconv.Atoi(gold)
		if err != nil {
			goldError = "Coin amount must be a number"
		}
	}
	platinum, ok := formData["platinum"]
	if ok {
		_, err := strconv.Atoi(platinum)
		if err != nil {
			platinumError = "Coin amount must be a number"
		}
	}

	return &AdventureForm{
		NumberCharacters: characterNumberString,
		NumberHenchmen:   henchmenNumberString,
		CharactersError:  characterError,
		HenchmenError:    henchmenError,
		CopperError:      copperError,
		SilverError:      silverError,
		ElectrumError:    electrumError,
		GoldError:        goldError,
		PlatinumError:    platinumError,
		Valid:            valid,
	}
}

func ParseAndValidateLoot(itemData map[string]string) *LootItem {
	valid := true
	errors := FormInputError{}
	name := itemData["loot-name"]
	if name == "" {
		errors.Name = "Name cannot be blank"
		valid = false
	}
	n, err := strconv.Atoi(itemData["loot-number"])
	if err != nil {
		errors.Number = "Number must be a numeric value."
		valid = false
	} else if n < 1 {
		errors.Number = "Number must be at least 1."
		valid = false
	}
	xp, err := strconv.Atoi(itemData["loot-xp-value"])
	if err != nil {
		errors.XPValue = "GP must be a numeric value."
		valid = false
	} else if xp < 1 {
		errors.XPValue = "GP must be at least 1."
		valid = false
	}
	return &LootItem{
		Name:    name,
		Number:  itemData["loot-number"],
		GPValue: itemData["loot-gp-value"],
		Errors:  errors,
		Valid:   valid,
	}
}

func ParseAndValidateCombat(itemData map[string]string) *CombatItem {
	valid := true
	errors := FormInputError{}
	name := itemData["combat-name"]
	if name == "" {
		errors.Name = "Name cannot be blank"
		valid = false
	}
	n, err := strconv.Atoi(itemData["combat-number"])
	if err != nil {
		errors.Number = "Number must be a numeric value."
		valid = false
	} else if n < 1 {
		errors.Number = "Number must be at least 1."
		valid = false
	}
	xp, err := strconv.Atoi(itemData["combat-xp-value"])
	if err != nil {
		errors.XPValue = "XP must be a numeric value."
		valid = false
	} else if xp < 1 {
		errors.XPValue = "XP must be at least 1."
		valid = false
	}
	return &CombatItem{
		Name:    name,
		XPValue: itemData["combat-cp-value"],
		Errors:  errors,
		Valid:   valid,
	}
}

func ParseAndValidateMagicItem(itemData map[string]string) *MagicItem {
	valid := true
	errors := FormInputError{}
	name := itemData["magicItem-name"]
	if name == "" {
		errors.Name = "Name cannot be blank"
		valid = false
	}
	xp, err := strconv.Atoi(itemData["magicItem-xp-value"])
	if err != nil {
		errors.XPValue = "Apparent value must be a numeric value."
		valid = false
	} else if xp < 1 {
		errors.XPValue = "Apparent must be at least 1."
		valid = false
	}
	return &MagicItem{
		Name:    name,
		XPValue: itemData["magicItem-xp-value"],
		Errors:  errors,
		Valid:   valid,
	}
}

func NewSimpleAdventurePage(r renderer.Renderer) *SimpleAdventurePage {
	return &SimpleAdventurePage{
		renderer: r,
	}

}

func (p SimpleAdventurePage) RegisterRoutes(r *http.ServeMux, g guardsman.Guardsman) {

}

func (p SimpleAdventurePage) AdventureCalculator(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

	}
}
