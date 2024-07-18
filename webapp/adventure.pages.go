package webapp

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/internal/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

type AdventurePage struct {
	adventureService services.AdventureService
	characterService services.CharacterService
	renderer         Renderer
}

type AdventurePageModel struct {
	Id            int
	CampaignId    int
	Name          string
	TotalXPAmount int
	FullShareXP   int
	HalfShareXP   int
	AdventureDate types.ArcvhistDate
	GameDays      int
	Coins         types.Coins
	Gems          []LootPageModel
	Jewellery     []LootPageModel
	Combat        []LootPageModel
	MagicItems    []LootPageModel
	Characters    []types.AdventureCharacter
	//These will hold the paths to the various editors
	DetailsPath path
	CoinPath    path
	NewLootPath path
}

type LootPageModel struct {
	Id            int
	LootType      types.GenericLootType `json:"type"`
	DisplayType   string
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Number        int     `json:"number"`
	XPValue       float64 `json:"xp_value"`
	GoldValue     float64 `json:"gold_value"`
	TotalXPAmount int
	Path          path
}

type MagicItemPageModel struct {
	Id          int
	Name        string `json:"name"`
	Description string `json:"Description"`
	XPValue     int    `json:"magic_item_xp"`
	GoldValue   int    `json:"actual_value"`
	Path        path
}

const basepath = "/pages/adventure"

func newPhysicalAdventurePath(resource string, id int) path {
	path := fmt.Sprintf(basepath+"/%d"+resource, id)
	return newPath(path)
}
func newPhysicalAdventurePathWithResourceId(resource string, adventureId, resourceId int) path {
	path := fmt.Sprintf(basepath+"/%d"+resource+"/%d", adventureId, resourceId)
	return newPath(path)
}

func createGemPageModels(data []types.Gem, aId int) []LootPageModel {
	out := make([]LootPageModel, 0)
	for _, g := range data {
		out = append(out, NewLootPageModelFromGem(g, aId))
	}
	return out
}

func createJewelleryPageModels(data []types.Jewellery, aId int) []LootPageModel {
	out := make([]LootPageModel, 0)
	for _, g := range data {
		out = append(out, NewLootPageModelFromJewellery(g, aId))
	}
	return out
}

func createCombatPageModels(data []types.MonsterGroup, aId int) []LootPageModel {
	out := make([]LootPageModel, 0)
	for _, g := range data {
		out = append(out, NewLootPageModelFromCombat(g, aId))
	}
	return out
}

func createMagicItemPageModels(data []types.MagicItem, aId int) []LootPageModel {
	out := make([]LootPageModel, 0)
	for _, g := range data {
		out = append(out, NewLootPageModelFromMagicItem(g, aId))
	}
	return out
}

func newAdventurePathToRegister(appendedPath string, additionalPathParams ...string) path {
	path := ""
	if len(additionalPathParams) == 1 {
		path = fmt.Sprintf(basepath+"/%s"+appendedPath+"/%s", "{adventureId}", additionalPathParams[0])
	} else {
		path = fmt.Sprintf(basepath+"/%s"+appendedPath, "{adventureId}")
	}
	return newPath(path)
}

func newAdventurePageModel(a types.AdventureRecord) AdventurePageModel {
	return AdventurePageModel{
		Id:            a.Id,
		CampaignId:    a.CampaignId,
		Name:          a.Name,
		TotalXPAmount: a.TotalXPAmount(),
		FullShareXP:   a.FullShareXP,
		HalfShareXP:   a.HalfShareXP,
		AdventureDate: a.AdventureDate,
		GameDays:      a.GameDays,
		Coins:         a.Coins,
		Characters:    a.Characters,
		DetailsPath:   newPhysicalAdventurePath("", a.Id),
		CoinPath:      newPhysicalAdventurePath("/coin", a.Id),
		NewLootPath:   newPhysicalAdventurePath("/new-loot", a.Id),
	}
}

func NewLootPageModelFromGem(adata types.Gem, adventureId int) LootPageModel {
	return LootPageModel{
		Id:            adata.Id,
		LootType:      types.GemLoot,
		DisplayType:   "Gem",
		Name:          adata.Name,
		Number:        adata.Number,
		XPValue:       adata.XPValue,
		TotalXPAmount: int(adata.TotalXPAmount()),
		Path:          newPhysicalAdventurePathWithResourceId("/gems", adventureId, adata.Id),
	}
}

func NewLootPageModelFromJewellery(adata types.Jewellery, adventureId int) LootPageModel {
	return LootPageModel{
		Id:            adata.Id,
		LootType:      types.JewelleryLoot,
		DisplayType:   "Jewellery",
		Name:          adata.Name,
		Number:        adata.Number,
		XPValue:       adata.XPValue,
		TotalXPAmount: int(adata.TotalXPAmount()),
		Path:          newPhysicalAdventurePathWithResourceId("/jewellery", adventureId, adata.Id),
	}
}

func NewLootPageModelFromCombat(adata types.MonsterGroup, adventureId int) LootPageModel {
	return LootPageModel{
		Id:            adata.Id,
		LootType:      types.CombatLoot,
		DisplayType:   "Combat",
		Name:          adata.Name,
		Number:        adata.NumberDefeated,
		XPValue:       float64(adata.XPPerOneKill),
		TotalXPAmount: int(adata.TotalXPAmount()),
		Path:          newPhysicalAdventurePathWithResourceId("/combat", adventureId, adata.Id),
	}
}

func NewLootPageModelFromMagicItem(adata types.MagicItem, adventureId int) LootPageModel {
	return LootPageModel{
		Id:            adata.Id,
		LootType:      types.GemLoot,
		DisplayType:   "Magic Item",
		Name:          adata.Name,
		XPValue:       float64(adata.XPValue),
		GoldValue:     float64(adata.GoldValue),
		TotalXPAmount: int(adata.TotalXPAmount()),
		Path:          newPhysicalAdventurePathWithResourceId("/magic-items", adventureId, adata.Id),
	}
}

func NewAdventurePage(cs services.AdventureService, ch services.CharacterService, r Renderer) *AdventurePage {
	return &AdventurePage{
		adventureService: cs,
		characterService: ch,
		renderer:         r,
	}

}

func (a AdventurePage) RegisterRoutes(m *http.ServeMux) {
	mainPath := newAdventurePathToRegister("")
	coinPath := newAdventurePathToRegister("/coin")
	gemPath := newAdventurePathToRegister("/gems/{gemId}")
	jewelleryPath := newAdventurePathToRegister("/jewellery/{jewelleryId}")
	combatPath := newAdventurePathToRegister("/combat/{combatId}")
	magicItemPath := newAdventurePathToRegister("/magic-item/{magicItemId}")
	newLootPath := newAdventurePathToRegister("/new-loot")
	//characterPath := newAdventurePathToRegister("/characters")

	m.HandleFunc(mainPath.Display, a.AdventureOverview)
	m.HandleFunc("PUT "+coinPath.Display, a.SaveAndDisplayCoins)
	m.HandleFunc("GET "+coinPath.Display, a.CoinDisplay)
	m.HandleFunc("GET "+coinPath.Edit, a.CoinEditHandler)

	m.HandleFunc("PUT "+gemPath.Display, a.saveGem)
	m.HandleFunc("GET "+gemPath.Display, a.displayGem)
	m.HandleFunc("DELETE "+gemPath.Display, a.deleteGem)
	m.HandleFunc("GET "+gemPath.Edit, a.editGem)

	m.HandleFunc("PUT "+jewelleryPath.Display, a.saveJewellery)
	m.HandleFunc("GET "+jewelleryPath.Display, a.displayJewellery)
	m.HandleFunc("DELETE "+jewelleryPath.Display, a.deleteJewellery)
	m.HandleFunc("GET "+jewelleryPath.Edit, a.editJewellery)

	m.HandleFunc("PUT "+combatPath.Display, a.saveCombat)
	m.HandleFunc("GET "+combatPath.Display, a.displayCombat)
	m.HandleFunc("DELETE "+combatPath.Display, a.deleteCombat)
	m.HandleFunc("GET "+combatPath.Edit, a.editCombat)

	m.HandleFunc("PUT "+magicItemPath.Display, a.saveMagicItem)
	m.HandleFunc("GET "+magicItemPath.Display, a.displayMagicItem)
	m.HandleFunc("DELETE "+magicItemPath.Display, a.deleteMagicItem)
	m.HandleFunc("GET "+magicItemPath.Edit, a.editMagicItem)

	m.HandleFunc("PUT "+newLootPath.Display, a.saveNewLoot)
	m.HandleFunc("GET "+newLootPath.Edit, a.newLoot)

	//m.HandleFunc("GET "+characterPath.Display, a.displayCharacters)

}

func (a AdventurePage) AdventureOverview(w http.ResponseWriter, r *http.Request) {
	//applyCorsHeaders(w)
	id := r.PathValue("adventureId")

	adventure, err := a.adventureService.GetAdventureRecordById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.renderAdventurePage(w, *adventure)
}

func (a AdventurePage) renderAdventurePage(w http.ResponseWriter, data types.AdventureRecord) {
	model := newAdventurePageModel(data)
	model.Gems = createGemPageModels(data.Gems, data.Id)
	model.Jewellery = createJewelleryPageModels(data.Jewellery, data.Id)
	model.Combat = createCombatPageModels(data.Combat, data.Id)
	model.MagicItems = createMagicItemPageModels(data.MagicItems, data.Id)
	output, err := a.renderer.RenderPage("adventurePage.html", model)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, output, err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) CoinDisplay(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	adata, err := a.adventureService.GetAdventureRecordById(id)
	adventure := newAdventurePageModel(*adata)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderPartial("coins.html", adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) CoinEditHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	adata, err := a.adventureService.GetAdventureRecordById(id)
	adventure := newAdventurePageModel(*adata)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderEditor("coinsEdit.html", adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) SaveAndDisplayCoins(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	formErr := r.ParseForm()
	if formErr != nil {
		a.renderer.MustRenderErrorPage(w, "", formErr)
	}
	var data = map[string]string{}
	for key, value := range r.Form {
		data[key] = value[0]
	}
	adata, err := a.adventureService.UpdateAdventureCoins(id, data)
	adventure := newAdventurePageModel(*adata)

	//coins, err := a.adventureService.GetCoinsForAdventure(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderPartial("coins.html", adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) saveGem(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("gemId")
	formErr := r.ParseForm()
	if formErr != nil {
		a.renderer.MustRenderErrorPage(w, "", formErr)
	}
	var data = map[string]string{}
	for key, value := range r.Form {
		data[key] = value[0]
	}
	err := a.adventureService.SaveGem(gemId, data)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	a.displayGem(w, r)
}

func (a AdventurePage) saveJewellery(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("jewelleryId")
	formErr := r.ParseForm()
	if formErr != nil {
		a.renderer.MustRenderErrorPage(w, "", formErr)
	}
	var data = map[string]string{}
	for key, value := range r.Form {
		data[key] = value[0]
	}
	err := a.adventureService.SaveJewellery(gemId, data)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	a.displayGem(w, r)
}

func (a AdventurePage) saveCombat(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("combatId")
	formErr := r.ParseForm()
	if formErr != nil {
		a.renderer.MustRenderErrorPage(w, "", formErr)
	}
	var data = map[string]string{}
	for key, value := range r.Form {
		data[key] = value[0]
	}
	err := a.adventureService.SaveCombat(gemId, data)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	a.displayGem(w, r)
}

func (a AdventurePage) saveMagicItem(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("magicItemId")
	formErr := r.ParseForm()
	if formErr != nil {
		a.renderer.MustRenderErrorPage(w, "", formErr)
	}
	var data = map[string]string{}
	for key, value := range r.Form {
		data[key] = value[0]
	}
	err := a.adventureService.SaveMagicItem(gemId, data)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	a.displayGem(w, r)
}

func (a AdventurePage) deleteGem(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("gemId")
	err := a.adventureService.DeleteGem(gemId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write(make([]byte, 0))
}

func (a AdventurePage) deleteJewellery(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("jewelleryId")
	err := a.adventureService.DeleteJewellery(gemId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write(make([]byte, 0))
}

func (a AdventurePage) deleteCombat(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("combatId")
	err := a.adventureService.DeleteCombat(gemId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write(make([]byte, 0))
}

func (a AdventurePage) deleteMagicItem(w http.ResponseWriter, r *http.Request) {
	gemId := r.PathValue("magicItemId")
	err := a.adventureService.DeleteMagicItem(gemId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write(make([]byte, 0))
}

func (a AdventurePage) displayGem(w http.ResponseWriter, r *http.Request) {

	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	id := r.PathValue("gemId")
	adata, err := a.adventureService.GetGemById(id)
	pageData := NewLootPageModelFromGem(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderPartial("loot.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) displayJewellery(w http.ResponseWriter, r *http.Request) {

	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	id := r.PathValue("jewelleryId")
	adata, err := a.adventureService.GetJewelleryById(id)
	pageData := NewLootPageModelFromJewellery(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderPartial("loot.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) displayCombat(w http.ResponseWriter, r *http.Request) {

	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	id := r.PathValue("combatId")
	adata, err := a.adventureService.GetCombatById(id)
	pageData := NewLootPageModelFromCombat(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderPartial("loot.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) displayMagicItem(w http.ResponseWriter, r *http.Request) {

	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	id := r.PathValue("magicItem")
	adata, err := a.adventureService.GetMagicItemById(id)
	pageData := NewLootPageModelFromMagicItem(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderPartial("magicItem.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) displayAdventureCharacters(w http.ResponseWriter, r *http.Request) {
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	//adata, err := a.characterService.GetCharactersForCampaign(types.NewCampaign(campaign))
	pageData, err := a.adventureService.GetCharactersForAdventure(adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}

	output, err := a.renderer.RenderPartial("characters.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))

}

func (a AdventurePage) editGem(w http.ResponseWriter, r *http.Request) {
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
		return
	}
	id := r.PathValue("gemId")
	adata, err := a.adventureService.GetGemById(id)
	pageData := NewLootPageModelFromGem(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	output, err := a.renderer.RenderEditor("lootEdit.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) editJewellery(w http.ResponseWriter, r *http.Request) {
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
		return
	}
	id := r.PathValue("jewelleryId")
	adata, err := a.adventureService.GetJewelleryById(id)
	pageData := NewLootPageModelFromJewellery(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	output, err := a.renderer.RenderEditor("lootEdit.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) editCombat(w http.ResponseWriter, r *http.Request) {
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
		return
	}
	id := r.PathValue("combatId")
	adata, err := a.adventureService.GetCombatById(id)
	pageData := NewLootPageModelFromCombat(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	output, err := a.renderer.RenderEditor("lootEdit.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) editMagicItem(w http.ResponseWriter, r *http.Request) {
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
		return
	}
	id := r.PathValue("magicItemId")
	adata, err := a.adventureService.GetMagicItemById(id)
	pageData := NewLootPageModelFromMagicItem(*adata, adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	output, err := a.renderer.RenderEditor("lootEdit.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) newLoot(w http.ResponseWriter, r *http.Request) {
	output := make([]byte, 0)
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", fmt.Errorf("could not convert %s into adventure id", aId))
	}
	t := r.URL.Query()["type"][0]
	pageData := LootPageModel{}
	switch t {
	case string(types.GemLoot):
		pageData.Name = "New Gem"
		pageData.LootType = types.GemLoot
		pageData.DisplayType = "Gem"
		pageData.XPValue = 0.0
		pageData.Number = 0.0
		pageData.TotalXPAmount = 0
		pageData.Path = newPhysicalAdventurePath("/new-loot", adventure)

		outString, err := a.renderer.RenderEditor("newLootModalEdit.html", pageData)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
		}
		output = []byte(outString)
		w.Write(output)
		return

	case string(types.JewelleryLoot):
		pageData.Name = "New Jewellery"
		pageData.LootType = types.JewelleryLoot
		pageData.DisplayType = "Jewellery"
		pageData.XPValue = 0.0
		pageData.Number = 0.0
		pageData.TotalXPAmount = 0
		pageData.Path = newPhysicalAdventurePath("/new-loot", adventure)

		outString, err := a.renderer.RenderEditor("newLootModalEdit.html", pageData)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
		}
		output = []byte(outString)
		w.Write(output)
	case string(types.CombatLoot):
		pageData.Name = "New Combat"
		pageData.LootType = types.CombatLoot
		pageData.DisplayType = "Combat"
		pageData.XPValue = 0.0
		pageData.Number = 0.0
		pageData.TotalXPAmount = 0
		pageData.Path = newPhysicalAdventurePath("/new-loot", adventure)

		outString, err := a.renderer.RenderEditor("newCombatModalEdit.html", pageData)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
		}
		output = []byte(outString)
		w.Write(output)
		return
	case string(types.MagicItemLoot):
		pageData.Name = "New Magic Item"
		pageData.LootType = types.MagicItemLoot
		pageData.DisplayType = "Magic Item"
		pageData.XPValue = 0.0
		pageData.GoldValue = 0
		pageData.TotalXPAmount = 0
		pageData.Path = newPhysicalAdventurePath("/new-loot", adventure)

		outString, err := a.renderer.RenderEditor("newMagicItemModalEdit.html", pageData)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
		}
		output = []byte(outString)
		w.Write(output)
		return

	}
}

func (a AdventurePage) saveNewLoot(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query()["type"][0]
	switch t {
	case string(types.GemLoot):
		a.saveNewGem(w, r)
	case string(types.JewelleryLoot):
		a.saveNewJewellery(w, r)
	case string(types.CombatLoot):
		a.saveNewCombat(w, r)
	case string(types.MagicItemLoot):
		a.saveNewMagicItem(w, r)
	}
}

func (a AdventurePage) saveNewGem(w http.ResponseWriter, r *http.Request) {
	//	output := make([]byte, 0)
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", fmt.Errorf("could not convert %s into adventure id", aId))
	}
	data, err := parseHTMLFormIntoMap(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	err = a.adventureService.SaveNewGem(adventure, data)
	//output, err = a.renderGemList(adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	// w.Write(output)
}

func (a AdventurePage) saveNewJewellery(w http.ResponseWriter, r *http.Request) {
	output := make([]byte, 0)
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", fmt.Errorf("could not convert %s into adventure id", aId))
	}
	data, err := parseHTMLFormIntoMap(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	err = a.adventureService.SaveNewJewellery(adventure, data)
	//output, err = a.renderJewelleryList(adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write(output)
}

func (a AdventurePage) saveNewCombat(w http.ResponseWriter, r *http.Request) {
	output := make([]byte, 0)
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", fmt.Errorf("could not convert %s into adventure id", aId))
	}
	data, err := parseHTMLFormIntoMap(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	err = a.adventureService.SaveNewCombat(adventure, data)
	//output, err = a.renderCombatList(adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write(output)
}

func (a AdventurePage) saveNewMagicItem(w http.ResponseWriter, r *http.Request) {
	output := make([]byte, 0)
	aId := r.PathValue("adventureId")
	adventure, err := strconv.Atoi(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", fmt.Errorf("could not convert %s into adventure id", aId))
	}
	data, err := parseHTMLFormIntoMap(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	err = a.adventureService.SaveNewMagicItem(adventure, data)
	//output, err = a.renderMagicItemList(adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write(output)
}

/*
	func (a AdventurePage) renderGemList(adventure int) ([]byte, error) {
		var renderData struct {
			DisplayType string
			Loot        []LootPageModel
		}
		gems, err := a.adventureService.GetGemsForAdventure(adventure)
		if err != nil {
			return make([]byte, 0), err
		}
		renderData.Loot = make([]LootPageModel, 0)
		for _, gem := range gems {
			data := NewLootPageModelFromGem(gem, adventure)
			renderData.DisplayType = data.DisplayType
			renderData.Loot = append(renderData.Loot, data)
		}
		str, err := a.renderer.RenderPartial("lootRange.html", renderData)
		return []byte(str), err
	}

	func (a AdventurePage) renderJewelleryList(adventure int) ([]byte, error) {
		var renderData struct {
			DisplayType string
			Loot        []LootPageModel
		}
		jewellery, err := a.adventureService.GetJewelleryForAdventure(adventure)
		if err != nil {
			return make([]byte, 0), err
		}
		renderData.Loot = make([]LootPageModel, 0)
		for _, j := range jewellery {
			data := NewLootPageModelFromJewellery(j, adventure)
			renderData.DisplayType = data.DisplayType
			renderData.Loot = append(renderData.Loot, data)
		}
		str, err := a.renderer.RenderPartial("lootRange.html", renderData)
		return []byte(str), err
	}

	func (a AdventurePage) renderCombatList(adventure int) ([]byte, error) {
		var renderData struct {
			DisplayType string
			Loot        []LootPageModel
		}
		combats, err := a.adventureService.GetCombatForAdventure(adventure)
		if err != nil {
			return make([]byte, 0), err
		}
		renderData.Loot = make([]LootPageModel, 0)
		for _, combat := range combats {
			data := NewLootPageModelFromCombat(combat, adventure)
			renderData.DisplayType = data.DisplayType
			renderData.Loot = append(renderData.Loot, data)
		}
		str, err := a.renderer.RenderPartial("combatRange.html", renderData)
		return []byte(str), err
	}

	func (a AdventurePage) renderMagicItemList(adventure int) ([]byte, error) {
		var renderData struct {
			DisplayType string
			Loot        []LootPageModel
		}
		magicItems, err := a.adventureService.GetMagicItemsForAdventure(adventure)
		if err != nil {
			return make([]byte, 0), err
		}
		renderData.Loot = make([]LootPageModel, 0)
		for _, mi := range magicItems {
			data := NewLootPageModelFromMagicItem(mi, adventure)
			renderData.DisplayType = data.DisplayType
			renderData.Loot = append(renderData.Loot, data)
		}
		str, err := a.renderer.RenderPartial("magicItemRange.html", renderData)
		return []byte(str), err
	}
*/
func parseHTMLFormIntoMap(r *http.Request) (map[string]string, error) {
	var data = map[string]string{}
	formErr := r.ParseForm()
	if formErr != nil {
		return nil, formErr
	}
	for key, value := range r.Form {
		data[key] = value[0]
	}
	return data, nil
}
