package archivist

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/renderer"
	"github.com/floodedrealms/borderland-keep/types"
)

type SimpleAdventurePage struct {
	renderer renderer.Renderer
}

type AdventureForm struct {
	types.Adventure
	CharacterCount  string
	HenchmenCount   string
	CharactersError string
	HenchmenError   string
	Copper          string
	Silver          string
	Electrum        string
	Gold            string
	Platinum        string
	CopperError     string
	SilverError     string
	ElectrumError   string
	GoldError       string
	PlatinumError   string
	Loot            []LootItem
	MagicItems      []MagicItem
	Combat          []CombatItem
	Valid           bool
}

type FormInputError struct {
	Name    string
	Number  string
	XPValue string
	GPValue string
}

type LootItem struct {
	types.Loot
	Number  int
	GPValue int
	Errors  FormInputError
	Valid   bool
}

type MagicItem struct {
	types.MagicalLoot
	XPValue string
	GPValue string
	Errors  FormInputError
	Valid   bool
}

type CombatItem struct {
	types.Combat
	Number  int
	XPValue int
	Errors  FormInputError
	Valid   bool
}

func ParseAndValidateDetails(formData url.Values) *AdventureForm {
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

		_, err := strconv.Atoi(characterNumberString[0])
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
		_, err := strconv.Atoi(henchmenNumberString[0])
		if err != nil {
			henchmenError = "Henchmen Number must be numeric"
			valid = false
		}
	}
	copper, ok := formData["copper"]
	if ok {
		_, err := strconv.Atoi(copper[0])
		if err != nil {
			copperError = "Coin amount must be a number"
		}
	}
	silver, ok := formData["silver"]
	if ok {
		_, err := strconv.Atoi(silver[0])
		if err != nil {
			silverError = "Coin amount must be a number"
		}
	}
	electrum, ok := formData["electrum"]
	if ok {
		_, err := strconv.Atoi(electrum[0])
		if err != nil {
			electrumError = "Coin amount must be a number"
		}
	}
	gold, ok := formData["gold"]
	if ok {
		_, err := strconv.Atoi(gold[0])
		if err != nil {
			goldError = "Coin amount must be a number"
		}
	}
	platinum, ok := formData["platinum"]
	if ok {
		_, err := strconv.Atoi(platinum[0])
		if err != nil {
			platinumError = "Coin amount must be a number"
		}
	}

	return &AdventureForm{
		CharacterCount:  characterNumberString[0],
		HenchmenCount:   henchmenNumberString[0],
		Copper:          copper[0],
		Silver:          silver[0],
		Electrum:        electrum[0],
		Gold:            gold[0],
		Platinum:        platinum[0],
		CharactersError: characterError,
		HenchmenError:   henchmenError,
		CopperError:     copperError,
		SilverError:     silverError,
		ElectrumError:   electrumError,
		GoldError:       goldError,
		PlatinumError:   platinumError,
		Valid:           valid,
	}
}

func ParseLootItem(itemData url.Values, cur int) *LootItem {
	n, err := strconv.Atoi(itemData["loot-number"][cur])
	if err != nil {
		n = 0
	}
	gp, err := strconv.Atoi(itemData["loot-gp-value"][cur])
	if err != nil {
		n = 0
	}
	return &LootItem{
		types.GenericLoot{
			Number:    n,
			GoldValue: float64(gp),
			XPValue:   float64(gp),
		},
		n,
		gp,
		FormInputError{},
		true,
	}
}

func ParseAndValidateCombat(itemData url.Values, cur int) *CombatItem {
	valid := true
	errors := FormInputError{}
	n, err := strconv.Atoi(itemData["combat-number"][cur])
	if err != nil {
		errors.Number = "Number must be a numeric value."
		valid = false
	} else if n < 1 {
		errors.Number = "Number must be at least 1."
		valid = false
	}
	xp, err := strconv.Atoi(itemData["combat-xp-value"][0])
	if err != nil {
		errors.XPValue = "XP must be a numeric value."
		valid = false
	} else if xp < 1 {
		errors.XPValue = "XP must be at least 1."
		valid = false
	}
	return &CombatItem{
		Combat:  *types.NewMonsterGroup("", "", n, xp),
		Number:  n,
		XPValue: xp,
		Errors:  errors,
		Valid:   valid,
	}
}

func ParseAndValidateMagicItem(itemData url.Values, cur int) *MagicItem {
	valid := true
	errors := FormInputError{}
	xp, err := strconv.Atoi(itemData["magic-item-gp-value"][cur])
	if err != nil {
		errors.XPValue = "Apparent value must be a numeric value."
		valid = false
	} else if xp < 1 {
		errors.XPValue = "Apparent must be at least 1."
		valid = false
	}
	s, ok := itemData["magic-item-sold"]
	sold := false
	if ok && s[0] == "on" {
		sold = true
	}
	soldAmount, ok := itemData["magic-item-gp-sold-value"]
	soldAmountString := ""
	gp := 0
	if ok && sold {
		gp, _ = strconv.Atoi(soldAmount[cur])
		soldAmountString = soldAmount[cur]
		if err != nil {
			errors.XPValue = "Apparent value must be a numeric value."
			valid = false
		} else if gp < 1 {
			errors.XPValue = "Apparent must be at least 1."
			valid = false
		}
	}

	return &MagicItem{
		MagicalLoot: types.NewMagicItem("", "", xp, gp, sold),
		XPValue:     itemData["magic-item-gp-value"][cur],
		GPValue:     soldAmountString,
		Errors:      errors,
		Valid:       valid,
	}
}

func ParseAndValidateForm(r *http.Request) *AdventureForm {
	r.ParseForm()
	form := ParseAndValidateDetails(r.Form)
	lootData := make([]LootItem, 0)
	for i := range r.Form["loot-id"] {
		item := ParseLootItem(r.Form, i)
		lootData = append(lootData, *item)
	}
	// using gems for now. Might split this back out, but it's probably to just combine the types since they are indentical
	g := make([]types.Gem, 0)
	for _, item := range lootData {
		g = append(g, types.Gem{Number: item.Number, XPValue: float64(item.DisplayXPAmount()), GoldValue: float64(item.DisplayGPAmount())})
	}

	combatData := make([]CombatItem, 0)
	for i := range r.Form["combat-id"] {
		item := ParseAndValidateCombat(r.Form, i)
		combatData = append(combatData, *item)
	}
	// using gems for now. Might split this back out, but it's probably to just combine the types since they are indentical
	m := make([]types.MonsterGroup, 0)
	for _, item := range combatData {
		m = append(m, *types.NewMonsterGroup("", "", item.Number, item.XPValue))
	}
	magicItemData := make([]MagicItem, 0)
	for i := range r.Form["magic-item-id"] {
		item := ParseAndValidateMagicItem(r.Form, i)
		magicItemData = append(magicItemData, *item)
	}
	mi := make([]types.MagicItem, 0)
	for _, item := range magicItemData {
		mi = append(mi, *types.NewMagicItem("", "", item.DisplayApparentValue(), item.DisplayActualValue(), item.WasSold()))
	}
	numChars, _ := strconv.Atoi(form.CharacterCount)
	numHenchmen, _ := strconv.Atoi(form.HenchmenCount)
	adventure := types.NewAdventureRecordByNumberOfCharacters(numChars, numHenchmen)
	copper, _ := strconv.Atoi(form.Copper)
	silver, _ := strconv.Atoi(form.Silver)
	electrum, _ := strconv.Atoi(form.Electrum)
	gold, _ := strconv.Atoi(form.Gold)
	platinum, _ := strconv.Atoi(form.Platinum)
	adventure.Coins = *types.NewCoins(copper, silver, electrum, gold, platinum)
	adventure.Gems = g
	adventure.Combat = m
	adventure.MagicItems = mi
	form.Adventure = adventure
	return form
}

func NewSimpleAdventurePage(r renderer.Renderer) *SimpleAdventurePage {
	return &SimpleAdventurePage{
		renderer: r,
	}
}

func (p SimpleAdventurePage) RegisterRoutes(r *http.ServeMux, g guardsman.Guardsman) {
	r.HandleFunc("/calculator", p.Calculator)
	r.HandleFunc("/calculator/coins", p.Coins)
	r.HandleFunc("/calculator/loot", p.Loot)
	r.HandleFunc("/calculator/attendance", p.Attendance)
	r.HandleFunc("/calculator/combat", p.Combat)
	r.HandleFunc("/calculator/magic-item", p.MagicItem)
}

func (p SimpleAdventurePage) Calculator(w http.ResponseWriter, r *http.Request) {
	user, canEdit := ExtractGuardsmanHeaders(r)
	language := util.ExtractLangageCookie(r)
	switch r.Method {
	case http.MethodGet:
		form := AdventureForm{
			Adventure:      *types.NewAdventureRecordByNumberOfCharacters(0, 0),
			CharacterCount: "0",
			HenchmenCount:  "0",
		}
		output, err := p.renderer.RenderPage("simpleAdventureCalculator.html", form, language, user, canEdit)
		if err != nil {
			RedirectToErrorPage(w, err, true)
			return
		}
		w.Write([]byte(output))
	case http.MethodPatch:
		r.ParseForm()
		form := ParseAndValidateForm(r)
		output, err := p.renderer.RenderPage("simpleAdventureDetails.html", form, language, user, canEdit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))

	}
}

func (p SimpleAdventurePage) Attendance(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		w.Header().Add("HX-Trigger", "updateOverview")
		w.WriteHeader(http.StatusOK)
	}
}

func (p SimpleAdventurePage) Coins(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		w.Header().Add("HX-Trigger", "updateOverview")
		w.WriteHeader(http.StatusOK)
	}
}

func (p SimpleAdventurePage) Loot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	case http.MethodPost:
		lootData := make([]LootItem, 0)
		for i := range r.Form["loot-id"] {
			item := ParseLootItem(r.Form, i)
			lootData = append(lootData, *item)
		}
		lootData = append(lootData, LootItem{types.GenericLoot{}, 0, 0, FormInputError{}, true})
		formData := AdventureForm{
			Loot: lootData,
		}
		output, err := p.renderer.RenderPartial("simpleLootForm.html", formData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))
	case http.MethodDelete:
		w.Header().Add("HX-Trigger", "updateOverview")
		w.WriteHeader(http.StatusOK)
	case http.MethodPatch:
		id, _ := strconv.Atoi(r.URL.Query()["item"][0])
		item := ParseLootItem(r.Form, id)

		output, err := p.renderer.RenderPartial("simpleLootDetails.html", item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("HX-Trigger", "updateOverview")
		w.Write([]byte(output))

	}
}

func (p SimpleAdventurePage) Combat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	case http.MethodPost:
		combatData := make([]CombatItem, 0)
		for i := range r.Form["combat-id"] {
			item := ParseAndValidateCombat(r.Form, i)
			combatData = append(combatData, *item)
		}
		combatData = append(combatData, CombatItem{types.MonsterGroup{}, 0, 0, FormInputError{}, true})
		formData := AdventureForm{
			Combat: combatData,
		}
		output, err := p.renderer.RenderPartial("simpleCombatForm.html", formData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))
	case http.MethodDelete:
		w.Header().Add("HX-Trigger", "updateOverview")
		w.WriteHeader(http.StatusOK)
	case http.MethodPatch:
		w.Header().Add("HX-Trigger", "updateOverview")
		w.WriteHeader(http.StatusOK)

	}
}

func (p SimpleAdventurePage) MagicItem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Method {
	case http.MethodPost:
		itemData := make([]MagicItem, 0)
		for i := range r.Form["magic-item-id"] {
			item := ParseAndValidateMagicItem(r.Form, i)
			itemData = append(itemData, *item)
		}
		itemData = append(itemData, MagicItem{types.MagicItem{}, "0", "0", FormInputError{}, true})
		formData := AdventureForm{
			MagicItems: itemData,
		}
		output, err := p.renderer.RenderPartial("simpleMagicItemForm.html", formData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))
	case http.MethodPut:
		itemData := make([]MagicItem, 0)
		for i := range r.Form["magic-item-id"] {
			item := ParseAndValidateMagicItem(r.Form, i)
			itemData = append(itemData, *item)
		}
		formData := AdventureForm{
			MagicItems: itemData,
		}
		output, err := p.renderer.RenderPartial("simpleMagicItemForm.html", formData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("HX-Trigger", "updateOverview")
		w.Write([]byte(output))
	case http.MethodDelete:
		w.Header().Add("HX-Trigger", "updateOverview")
		w.WriteHeader(http.StatusOK)
	case http.MethodPatch:
		w.Header().Add("HX-Trigger", "updateOverview")
		w.WriteHeader(http.StatusOK)
	}
}
