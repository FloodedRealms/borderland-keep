package webapp

import (
	"fmt"
	"net/http"

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
	MagicItems    []types.MagicItem
	Characters    []types.AdventureCharacter
	//These will hold the paths to the various editors
	DetailsPath    path
	CoinPath       path
	CharactersPath path
}

type LootPageModel struct {
	Id          int
	LootType    types.GenericLootType `json:"type"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Number      int                   `json:"number"`
	XPValue     float64               `json:"xp_value"`
	GoldValue   float64               `json:"gold_value"`
	Path        path
}

const basepath = "/pages/adventure"

func newAdventurePath(appendedPath string, id int) path {
	path := fmt.Sprintf(basepath+"/%d"+appendedPath, id)
	return newPath(path)
}

func newAdventurePathToRegister(appendedPath string) path {
	path := fmt.Sprintf(basepath+"/%s"+appendedPath, "{adventureId}")
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
		MagicItems:    a.MagicItems,
		Characters:    a.Characters,
		DetailsPath:   newAdventurePath("", a.Id),
		CoinPath:      newAdventurePath("/coin", a.Id),
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
	m.HandleFunc(mainPath.Display, a.AdventureOverview)
	m.HandleFunc("PUT "+coinPath.Display, a.SaveAndDisplayCoins)
	m.HandleFunc("GET "+coinPath.Display, a.CoinSummary)
	m.HandleFunc("GET "+coinPath.Edit, a.CoinEditHandler)

}

func (a AdventurePage) AdventureOverview(w http.ResponseWriter, r *http.Request) {
	//applyCorsHeaders(w)
	id := r.PathValue("adventureId")

	adventure, err := a.adventureService.GetAdventureRecordById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	print(adventure)

	a.renderAdventurePage(w, *adventure)
}

func (a AdventurePage) renderAdventurePage(w http.ResponseWriter, data types.AdventureRecord) {
	model := newAdventurePageModel(data)
	output, err := a.renderer.RenderPage("adventurePage.html", model)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, output, err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) CoinSummary(w http.ResponseWriter, r *http.Request) {
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
