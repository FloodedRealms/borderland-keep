package archivist

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/renderer"
	"github.com/floodedrealms/borderland-keep/types"
)

const baseAdventurePath = "/adventure"
const DATE_DISPLAY = "2006-01-02"

type AdventurePage struct {
	adventureService services.AdventureService
	characterService services.CharacterService
	renderer         renderer.Renderer
}

type AdventurePageModel struct {
	Id            int
	CampaignId    int
	Name          string
	TotalXPAmount int
	FullShareXP   int
	HalfShareXP   int
	AdventureDate string
	GameDays      int
	Coins         types.Coins
	Gems          []LootPageModel
	Jewellery     []LootPageModel
	Combat        []LootPageModel
	MagicItems    []LootPageModel
	Characters    []types.AdventureCharacter
	//These will hold the paths to the various editors
	ReturnPath    string
	DetailsPath   path
	CoinPath      path
	GemPath       path
	JewelleryPath path
	CombatPath    path
	MagicItemPath path
	NewLootPath   path
	User          guardsman.WebUser
	HasEditAccess bool
}

type formError struct {
	Name    string
	Number  string
	XPValue string
}

type LootPageModel struct {
	Id            int
	LootType      types.GenericLootType `json:"type"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	Number        string                `json:"number"`
	XPValue       string                `json:"xp_value"`
	GoldValue     string                `json:"gold_value"`
	TotalXPAmount int
	Path          path
	GPPath        path
	Errors        formError
}

func newLootPageModel(aId int) *LootPageModel {
	return &LootPageModel{
		GPPath: newPhysicalAdventurePath("/gold-toggle", aId),
	}
}

type EditorPage struct {
	Items       []LootPageModel
	NewItemPath path
	Path        path
	Type        string
	DisplayType string
}

type MagicItemPageModel struct {
	Id          int
	Name        string `json:"name"`
	Description string `json:"Description"`
	XPValue     int    `json:"magic_item_xp"`
	GoldValue   int    `json:"actual_value"`
	Path        path
}

func newPhysicalAdventurePath(resource string, id int) path {
	path := fmt.Sprintf(baseAdventurePath+"/%d"+resource, id)
	return newPath(path)
}

func newPhysicalAdventurePathWithResourceId(resource string, adventureId, resourceId int) path {
	path := fmt.Sprintf(baseAdventurePath+"/%d"+resource+"/%d", adventureId, resourceId)
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
		path = fmt.Sprintf(baseAdventurePath+"/%s"+appendedPath+"/%s", "{adventureId}", additionalPathParams[0])
	} else {
		path = fmt.Sprintf(baseAdventurePath+"/%s"+appendedPath, "{adventureId}")
	}
	return newPath(path)
}

func newAdventurePageModel(a types.AdventureRecord, loggedIn, canEdit bool) AdventurePageModel {
	return AdventurePageModel{
		Id:            a.Id,
		CampaignId:    a.CampaignId,
		Name:          a.Name,
		TotalXPAmount: a.TotalXPAmount(),
		FullShareXP:   a.FullShareXP,
		HalfShareXP:   a.HalfShareXP,
		AdventureDate: a.AdventureDate.Format(DATE_DISPLAY),
		GameDays:      a.GameDays,
		Coins:         a.Coins,
		Characters:    a.Characters,
		ReturnPath:    fmt.Sprintf("/campaign/%d", a.CampaignId),
		DetailsPath:   newPhysicalAdventurePath("/adventure-details", a.Id),
		CoinPath:      newPhysicalAdventurePath("/coin", a.Id),
		GemPath:       newPhysicalAdventurePath("/gems", a.Id),
		JewelleryPath: newPhysicalAdventurePath("/jewellery", a.Id),
		CombatPath:    newPhysicalAdventurePath("/combat", a.Id),
		MagicItemPath: newPhysicalAdventurePath("/magic-item", a.Id),
		NewLootPath:   newPhysicalAdventurePath("/new-loot", a.Id),
		User:          guardsman.WebUser{LoggedIn: loggedIn},
		HasEditAccess: canEdit,
	}
}

func NewLootPageModelFromGem(adata types.Gem, adventureId int) LootPageModel {
	m := newLootPageModel(adventureId)
	m.Id = adata.Id
	m.LootType = types.GemLoot
	m.Name = adata.Name
	m.Number = intToString(adata.Number)
	m.XPValue = intToString(int(adata.XPValue))
	m.TotalXPAmount = int(adata.TotalXPAmount())
	m.Path = newPhysicalAdventurePath("/gems", adventureId)
	return *m
}

func NewLootPageModelFromJewellery(adata types.Jewellery, adventureId int) LootPageModel {
	m := newLootPageModel(adventureId)
	m.Id = adata.Id
	m.LootType = types.JewelleryLoot
	m.Name = adata.Name
	m.Number = intToString(adata.Number)
	m.XPValue = intToString(int(adata.XPValue))
	m.TotalXPAmount = int(adata.TotalXPAmount())
	m.Path = newPhysicalAdventurePath("/jewellery", adventureId)
	return *m
}

func NewLootPageModelFromCombat(adata types.MonsterGroup, adventureId int) LootPageModel {
	m := newLootPageModel(adventureId)
	m.Id = adata.Id
	m.LootType = types.CombatLoot
	m.Name = adata.Name
	m.Number = intToString(adata.NumberDefeated)
	m.XPValue = intToString(int(adata.XPPerOneKill))
	m.TotalXPAmount = int(adata.TotalXPAmount())
	m.Path = newPhysicalAdventurePath("/combat", adventureId)
	return *m
}

func NewLootPageModelFromMagicItem(adata types.MagicItem, adventureId int) LootPageModel {
	m := newLootPageModel(adventureId)
	m.Id = adata.Id
	m.LootType = types.CombatLoot
	m.Name = adata.Name
	m.Number = "1"
	m.XPValue = intToString(int(adata.XPValue))
	m.GoldValue = intToString(adata.GoldValue)
	m.TotalXPAmount = int(adata.TotalXPAmount())
	m.Path = newPhysicalAdventurePath("/combat", adventureId)
	return *m
}

func NewAdventurePage(cs services.AdventureService, ch services.CharacterService, r renderer.Renderer) *AdventurePage {
	return &AdventurePage{
		adventureService: cs,
		characterService: ch,
		renderer:         r,
	}
}

func (a AdventurePage) RegisterRoutes(m *http.ServeMux, g guardsman.Guardsman) {
	mainPath := newAdventurePathToRegister("")
	detailsPath := newAdventurePathToRegister("/details")
	coinPath := newAdventurePathToRegister("/coin")
	gemPath := newAdventurePathToRegister("/gems")
	jewelleryPath := newAdventurePathToRegister("/jewellery")
	combatPath := newAdventurePathToRegister("/combat")
	magicItemPath := newAdventurePathToRegister("/magic-item")
	newLootPath := newAdventurePathToRegister("/new-loot")
	goldTogglePath := newAdventurePathToRegister("/gold-toggle")
	adventureDetailsPath := newAdventurePathToRegister("/adventure-details")

	m.HandleFunc(mainPath.Display, g.UserLoggedInAndHasEditAccessToAdventure(a.AdventureOverview))
	m.HandleFunc("DELETE "+mainPath.Display, a.deleteAdventure)

	m.HandleFunc(detailsPath.Display, a.updateDetails)

	m.HandleFunc("GET "+newLootPath.Display, a.newLoot)

	m.HandleFunc("GET "+goldTogglePath.Edit, func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("adventureId")
		aId, err := strconv.Atoi(id)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "", err)
			return
		}
		output, _ := a.renderer.RenderPartial("goldInputTD.html", struct{ GPPath path }{GPPath: newPhysicalAdventurePath("/gold-toggle", aId)})
		w.Write([]byte(output))
	})
	m.HandleFunc("GET "+goldTogglePath.Display, func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("adventureId")
		aId, err := strconv.Atoi(id)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "", err)
			return
		}
		output, _ := a.renderer.RenderPartial("toggleGoldButtonTD.html", struct{ GPPath path }{GPPath: newPhysicalAdventurePath("/gold-toggle", aId)})
		w.Write([]byte(output))
	})

	m.HandleFunc("PUT "+coinPath.Display, a.SaveAndDisplayCoins)
	m.HandleFunc("GET "+coinPath.Display, a.CoinDisplay)
	m.HandleFunc("GET "+coinPath.Edit, a.CoinEditHandler)

	m.HandleFunc("GET "+gemPath.Edit, a.openGemEditor)
	m.HandleFunc("GET "+gemPath.Display, a.displayGemList)
	m.HandleFunc("POST "+gemPath.Display, a.saveGems)
	m.HandleFunc("DELETE "+gemPath.Display, sendEmptyResponse)

	m.HandleFunc("GET "+jewelleryPath.Edit, a.openJewelleryEditor)
	m.HandleFunc("GET "+jewelleryPath.Display, a.displayJewelleryList)
	m.HandleFunc("POST "+jewelleryPath.Display, a.saveJewellery)
	m.HandleFunc("DELETE "+jewelleryPath.Display, sendEmptyResponse)

	m.HandleFunc("GET "+combatPath.Edit, a.openCombatEditor)
	m.HandleFunc("GET "+combatPath.Display, a.displayCombatList)
	m.HandleFunc("POST "+combatPath.Display, a.saveCombat)
	m.HandleFunc("DELETE "+combatPath.Display, sendEmptyResponse)

	m.HandleFunc("GET "+magicItemPath.Edit, a.openMagicItemEditor)
	m.HandleFunc("GET "+magicItemPath.Display, a.displayMagicItemList)
	m.HandleFunc("POST "+magicItemPath.Display, a.saveMagicItems)
	m.HandleFunc("DELETE "+magicItemPath.Display, sendEmptyResponse)

	m.HandleFunc("GET "+adventureDetailsPath.Edit, a.openAdventureEditor)
	m.HandleFunc("GET "+adventureDetailsPath.Display, a.updateDetails)
	m.HandleFunc("POST "+adventureDetailsPath.Display, a.saveAdventureDetails)

}

func (a AdventurePage) AdventureOverview(w http.ResponseWriter, r *http.Request) {
	id, _ := a.extractAdventureId(r)
	data, err := a.adventureService.GetAdventureRecordById(id)
	l, c := ExtractGuardsmanHeaders(r)
	language := util.ExtractLangageCookie(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	model := newAdventurePageModel(*data, l, c)
	model.Gems = createGemPageModels(data.Gems, data.Id)
	model.Jewellery = createJewelleryPageModels(data.Jewellery, data.Id)
	model.Combat = createCombatPageModels(data.Combat, data.Id)
	model.MagicItems = createMagicItemPageModels(data.MagicItems, data.Id)
	output, err := a.renderer.RenderPage("adventurePage.html", model, language, l, c)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, output, err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) renderAdventurePage(w http.ResponseWriter, data types.AdventureRecord, loggedin, canedit bool) {

}

func (a AdventurePage) updateDetails(w http.ResponseWriter, r *http.Request) {
	aid, _ := a.extractAdventureId(r)
	adventure, _ := a.adventureService.GetAdventureRecordById(aid)
	pdata := AdventurePageModel{
		DetailsPath:   newPhysicalAdventurePath("/adventure-details", aid),
		TotalXPAmount: adventure.TotalXPAmount(),
		FullShareXP:   adventure.FullShareXP,
		HalfShareXP:   adventure.HalfShareXP,
		Characters:    adventure.Characters,
		Name:          adventure.Name,
		AdventureDate: adventure.AdventureDate.Format(DATE_DISPLAY),
	}
	output, err := a.renderer.RenderPartial("adventureOverview.html", pdata)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) CoinDisplay(w http.ResponseWriter, r *http.Request) {
	id, _ := a.extractAdventureId(r)
	adata, err := a.adventureService.GetAdventureRecordById(id)
	l, c := ExtractGuardsmanHeaders(r)
	adventure := newAdventurePageModel(*adata, l, c)
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
	id, _ := a.extractAdventureId(r)
	adata, err := a.adventureService.GetAdventureRecordById(id)
	l, c := ExtractGuardsmanHeaders(r)
	adventure := newAdventurePageModel(*adata, l, c)
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
	coins, err := a.adventureService.UpdateAdventureCoins(id, data)
	aid, _ := strconv.Atoi(id)
	pageData := struct {
		types.Coins
		CoinPath path
	}{
		*coins,
		newPhysicalAdventurePath("coins", aid),
	}
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.RenderPartial("coins.html", pageData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}

	w.Header().Set("HX-Trigger", "updateOverview")
	w.Write([]byte(output))
}

func (a AdventurePage) newLoot(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "error.html", fmt.Errorf("could not convert %s into adventure id", aId))
	}
	t := r.URL.Query()["type"][0]
	pdata := newLootPageModel(aId)
	switch t {

	case string(types.GemLoot):
		pdata.Path = newPhysicalAdventurePath("/gems", aId)
		outString, err := a.renderer.RenderPartial("lootFormRow.html", pdata)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
			return
		}
		output := []byte(outString)
		w.Write(output)
		return

	case string(types.JewelleryLoot):
		pdata.Path = newPhysicalAdventurePath("/jewellery", aId)
		outString, err := a.renderer.RenderPartial("lootFormRow.html", pdata)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
			return
		}
		output := []byte(outString)
		w.Write(output)
		return
	case string(types.CombatLoot):
		pdata.Path = newPhysicalAdventurePath("/jewellery", aId)
		outString, err := a.renderer.RenderPartial("lootFormRow.html", pdata)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
			return
		}
		output := []byte(outString)
		w.Write(output)
		return
	case string(types.MagicItemLoot):
		pdata.Path = newPhysicalAdventurePath("/magic-item", aId)
		outString, err := a.renderer.RenderPartial("lootFormRow.html", pdata)
		if err != nil {
			a.renderer.MustRenderErrorPage(w, "error.html", err)
			return
		}
		output := []byte(outString)
		w.Write(output)
		return
	}
}

func (a AdventurePage) openGemEditor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetGemsForAdventure(aId)
	gems := createGemPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	pData := EditorPage{
		Items:       gems,
		NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
		Path:        newPhysicalAdventurePath("/gems", aId),
		Type:        string(types.GemLoot),
		DisplayType: "Gems",
	}

	output, err := a.renderer.RenderEditor("lootTableEdit.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) openJewelleryEditor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetJewelleryForAdventure(aId)
	jewllery := createJewelleryPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	pData := EditorPage{
		Items:       jewllery,
		NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
		Path:        newPhysicalAdventurePath("/jewellery", aId),
		Type:        string(types.JewelleryLoot),
		DisplayType: "Jewellery",
	}

	output, err := a.renderer.RenderEditor("lootTableEdit.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) openCombatEditor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetCombatForAdventure(aId)
	combats := createCombatPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	pData := EditorPage{
		Items:       combats,
		NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
		Path:        newPhysicalAdventurePath("/combat", aId),
		Type:        string(types.CombatLoot),
		DisplayType: "Combat",
	}

	output, err := a.renderer.RenderEditor("combatTableEdit.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) openMagicItemEditor(w http.ResponseWriter, r *http.Request) {
	aId, err := a.extractAdventureId(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetMagicItemsForAdventure(aId)
	jewllery := createMagicItemPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	pData := EditorPage{
		Items:       jewllery,
		NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
		Path:        newPhysicalAdventurePath("/magic-item", aId),
		Type:        string(types.MagicItemLoot),
		DisplayType: "Magic Items",
	}

	output, err := a.renderer.RenderEditor("lootTableEdit.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) openAdventureEditor(w http.ResponseWriter, r *http.Request) {
	aId, err := a.extractAdventureId(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	possibleCharacters, attendance, err := a.adventureService.GetPossibleCharactersForAdventure(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adventureInfo, err := a.adventureService.GetAdventureRecordById(aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	pData := struct {
		AdventureId int
		Name        string
		Date        string
		Characters  []types.AdventureCharacter
		Attendance  []bool
		Path        path
		DisplayType string
	}{
		AdventureId: adventureInfo.Id,
		Name:        adventureInfo.Name,
		Date:        adventureInfo.AdventureDate.Format("2006-01-02"),
		Characters:  possibleCharacters,
		Attendance:  attendance,
		Path:        newPhysicalAdventurePath("/adventure-details", aId),
		DisplayType: "Characters",
	}

	output, err := a.renderer.RenderEditor("adventureDetailsEdit.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Write([]byte(output))
}

func (a AdventurePage) displayGemList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetGemsForAdventure(aId)
	gems := createGemPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	pData := EditorPage{
		Items:       gems,
		Path:        newPhysicalAdventurePath("/gems", aId),
		DisplayType: "Gemstones",
	}

	output, err := a.renderer.RenderPartial("lootList.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Header().Set("HX-Trigger", "updateOverview")

	w.Write([]byte(output))
}

func (a AdventurePage) displayJewelleryList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetJewelleryForAdventure(aId)
	jewellery := createJewelleryPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	pData := EditorPage{
		Items:       jewellery,
		Path:        newPhysicalAdventurePath("/jewellery", aId),
		DisplayType: "Jewellery",
	}

	output, err := a.renderer.RenderPartial("lootList.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Header().Set("HX-Trigger", "updateOverview")

	w.Write([]byte(output))
}

func (a AdventurePage) displayCombatList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetCombatForAdventure(aId)
	combat := createCombatPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	pData := EditorPage{
		Items:       combat,
		Path:        newPhysicalAdventurePath("/combat", aId),
		DisplayType: "Combat",
	}

	output, err := a.renderer.RenderPartial("combatList.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Header().Set("HX-Trigger", "updateOverview")

	w.Write([]byte(output))
}

func (a AdventurePage) displayMagicItemList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adata, err := a.adventureService.GetMagicItemsForAdventure(aId)
	jewellery := createMagicItemPageModels(adata, aId)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	pData := EditorPage{
		Items:       jewellery,
		Path:        newPhysicalAdventurePath("/magic-item", aId),
		DisplayType: "Magic Items",
	}

	output, err := a.renderer.RenderPartial("lootList.html", pData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.Header().Set("HX-Trigger", "updateOverview")

	w.Write([]byte(output))
}

func (a AdventurePage) saveGems(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	data, _ := parseForm(r)
	valid, loot := validateLootForm(data, types.GemLoot, aId)
	if !valid {
		pdata := EditorPage{
			Items:       loot,
			NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
			Path:        newPhysicalAdventurePath("/gems", aId),
			Type:        string(types.GemLoot),
			DisplayType: "Gems",
		}
		output, _ := a.renderer.RenderEditor("lootTableEdit.html", pdata)
		w.Write([]byte(output))
		return
	}
	a.adventureService.ModifyGems(aId, data)
	a.displayGemList(w, r)
}

func (a AdventurePage) saveJewellery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	data, _ := parseForm(r)
	valid, loot := validateLootForm(data, types.JewelleryLoot, aId)

	if !valid {
		pdata := EditorPage{
			Items:       loot,
			NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
			Path:        newPhysicalAdventurePath("/jewellery", aId),
			Type:        string(types.JewelleryLoot),
			DisplayType: "Jewellery",
		}
		output, _ := a.renderer.RenderEditor("lootTableEdit.html", pdata)
		w.Write([]byte(output))
		return
	}
	a.adventureService.ModifyJewellery(aId, data)
	a.displayJewelleryList(w, r)
}

func (a AdventurePage) saveMagicItems(w http.ResponseWriter, r *http.Request) {
	aId, err := a.extractAdventureId(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	data, _ := parseForm(r)
	valid, loot := validateLootForm(data, types.MagicItemLoot, aId)

	if !valid {
		pdata := EditorPage{
			Items:       loot,
			NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
			Path:        newPhysicalAdventurePath("/magic-item", aId),
			Type:        string(types.MagicItemLoot),
			DisplayType: "Magic Items",
		}
		output, _ := a.renderer.RenderEditor("lootTableEdit.html", pdata)
		w.Write([]byte(output))
		return
	}
	a.adventureService.ModifyMagicItems(aId, data)
	a.displayMagicItemList(w, r)
}

func (a AdventurePage) saveCombat(w http.ResponseWriter, r *http.Request) {
	aId, err := a.extractAdventureId(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	data, _ := parseForm(r)
	valid, loot := validateLootForm(data, types.CombatLoot, aId)
	if !valid {
		pdata := EditorPage{
			Items:       loot,
			NewItemPath: newPhysicalAdventurePath("/new-loot", aId),
			Path:        newPhysicalAdventurePath("/combat", aId),
			Type:        string(types.CombatLoot),
			DisplayType: "Combat",
		}
		output, _ := a.renderer.RenderEditor("combatTableEdit.html", pdata)
		w.Write([]byte(output))
		return
	}
	a.adventureService.ModifyCombat(aId, data)
	a.displayCombatList(w, r)
}

func (a AdventurePage) saveAdventureDetails(w http.ResponseWriter, r *http.Request) {
	id, _ := a.extractAdventureId(r)
	formErr := r.ParseForm()
	if formErr != nil {
		a.renderer.MustRenderErrorPage(w, "", formErr)
		return
	}
	adventureData, err := parseAdventureMetaData(*r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	err = a.adventureService.ModifyMetadata(*adventureData)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	charData, err := parseCharacterForm(*r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	adv, err := a.adventureService.GetAdventureRecordById(adventureData.Id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	f, h := adv.CalculateXPShares()
	err = a.adventureService.ModifyCharacters(id, charData, h, f)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	a.updateDetails(w, r)
}

func (a AdventurePage) deleteAdventure(w http.ResponseWriter, r *http.Request) {
	id, err := a.extractAdventureId(r)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	err = a.adventureService.DeleteAdventure(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func parseForm(r *http.Request) ([]map[string]string, error) {
	var data = []map[string]string{}
	formErr := r.ParseForm()
	if formErr != nil {
		return nil, formErr
	}
	for key, values := range r.Form {
		for i, value := range values {
			if i >= len(data) {
				m := map[string]string{}
				m[key] = value
				data = append(data, m)
			} else {
				data[i][key] = value
			}
		}
	}

	return data, nil
}

func parseCharacterForm(r http.Request) ([]types.AdventureCharacter, error) {

	var characters = make([]types.AdventureCharacter, 0)
	for _, id := range r.Form["character-id"] {
		_, ok := r.Form[fmt.Sprintf("on-adventure-%s", id)]
		if !ok {
			continue
		}
		charId, _ := strconv.Atoi(id)
		halfShare := false
		xpType, ok := r.Form[fmt.Sprintf("character-type-%s", id)]
		if !ok {
			return nil, fmt.Errorf("unable to get character type")
		}
		if xpType[0] == "henchmen" {
			halfShare = true
		}
		characters = append(characters, *types.NewAdventureCharacter(halfShare, charId))

	}
	return characters, nil
}

func parseAdventureMetaData(r http.Request) (*types.AdventureRecord, error) {
	idStr := r.Form["adventure-id"][0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	name := r.Form["adventure-name"][0]
	dateString := r.Form["adventure-date"][0]
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return &types.AdventureRecord{
			Id:            id,
			Name:          name,
			AdventureDate: date,
		}, err
	}
	return &types.AdventureRecord{
		Id:            id,
		Name:          name,
		AdventureDate: date,
	}, nil
}

func validateLootForm(data []map[string]string, lt types.GenericLootType, aId int) (bool, []LootPageModel) {
	loot := make([]LootPageModel, len(data))
	valid := true
	for i, itemData := range data {

		errors := formError{}
		n, err := strconv.Atoi(itemData["number"])
		if itemData["name"] == "" {
			errors.Name = "Name cannot be blank"
			valid = false
		}
		if err != nil {
			errors.Number = "Number must be a numeric value."
			valid = false
		} else if n < 1 {
			errors.Number = "Number must be at least 1."
			valid = false
		}
		xp, err := strconv.Atoi(itemData["xp-value"])
		if err != nil {
			errors.XPValue = "XP must be a numeric value."
			valid = false
		} else if xp < 1 {
			errors.XPValue = "XP must be at least 1."
			valid = false
		}
		item := LootPageModel{
			Name:    itemData["name"],
			Number:  itemData["number"],
			XPValue: itemData["xp-value"],
			Errors:  errors,
		}
		switch lt {
		case types.GemLoot:
			item.Path = newPhysicalAdventurePath("/gems", aId)
		case types.JewelleryLoot:
			item.Path = newPhysicalAdventurePath("/gems", aId)
		case types.CombatLoot:
			item.Path = newPhysicalAdventurePath("/gems", aId)

		}
		loot[i] = item
	}
	return valid, loot
}

func (a AdventurePage) extractAdventureId(r *http.Request) (int, error) {
	id := r.PathValue("adventureId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return aId, nil
}

func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}
