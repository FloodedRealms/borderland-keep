package webapp

import (
	"net/http"

	"github.com/floodedrealms/adventure-archivist/internal/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

type AdventurePage struct {
	adventureService services.AdventureService
	characterService services.CharacterService
	renderer         Renderer
}

func NewAdventurePage(cs services.AdventureService, ch services.CharacterService, r Renderer) *AdventurePage {
	return &AdventurePage{
		adventureService: cs,
		characterService: ch,
		renderer:         r,
	}

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
	output, err := a.renderer.Render("adventurePage.html", data)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, output, err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) CoinSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	adventure, err := a.adventureService.GetAdventureRecordById(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.Render("coins.html", adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) CoinEditHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	adventure, err := a.adventureService.GetAdventureRecordById(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.Render("coinsEdit.html", adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (a AdventurePage) SaveAndDisplayCOind(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("adventureId")
	formErr := r.ParseForm()
	if formErr != nil {
		a.renderer.MustRenderErrorPage(w, "", formErr)
	}
	var data = map[string]string{}
	for key, value := range r.Form {
		data[key] = value[0]
	}
	adventure, err := a.adventureService.UpdateAdventureCoins(id, data)

	//coins, err := a.adventureService.GetCoinsForAdventure(id)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	output, err := a.renderer.Render("coins.html", adventure)
	if err != nil {
		a.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}
