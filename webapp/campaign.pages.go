package webapp

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/internal/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

const baseCampaignPath = "/campaign"

type CampaignPage struct {
	campaignService  services.CampaignService
	characterService services.CharacterService
	renderer         Renderer
	CharacterPath    path
}

func NewCampaignPage(cs services.CampaignService, ch services.CharacterService, r Renderer) *CampaignPage {
	return &CampaignPage{
		campaignService:  cs,
		characterService: ch,
		renderer:         r,
		CharacterPath:    newCampaignPathToRegister("/characters"),
	}

}
func newCampaignPathToRegister(appendedPath string, additionalPathParams ...string) path {
	path := ""
	if len(additionalPathParams) == 1 {
		path = fmt.Sprintf(baseCampaignPath+"/%s"+appendedPath+"/%s", "{campaignId}", additionalPathParams[0])
	} else {
		path = fmt.Sprintf(baseCampaignPath+"/%s"+appendedPath, "{campaignId}")
	}
	return newPath(path)
}
func newPhysicalCampaignPath(resource string, id int) path {
	path := fmt.Sprintf(baseCampaignPath+"/%d"+resource, id)
	return newPath(path)
}

func (c CampaignPage) RegisterRoutes(m *http.ServeMux) {
	mainPath := newCampaignPathToRegister("")
	characterPath := newCampaignPathToRegister("/characters")
	newCharacterPath := newCampaignPathToRegister("/new-character")

	m.HandleFunc(mainPath.Display, c.CampaignOverview)

	m.HandleFunc("GET "+characterPath.Edit, c.openCharacterEditor)
	m.HandleFunc("GET "+characterPath.Display, c.displayCharacterList)
	m.HandleFunc("POST "+characterPath.Display, c.saveCharacters)
	m.HandleFunc("DELETE "+characterPath.Display, c.DeleteCharacter)

	m.HandleFunc("GET "+newCharacterPath.Display, c.newCharacter)
	m.HandleFunc("DELETE "+newCharacterPath.Display, c.DeleteUnSavedCharacter)
}

func (ca CampaignPage) CampaignOverview(w http.ResponseWriter, r *http.Request) {
	//applyCorsHeaders(w)
	id := r.PathValue("campaignId")

	campaign, err := ca.campaignService.GetCampaign(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	campaign.Characters, err = ca.characterService.GetCharactersForCampaign(campaign)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pdata := struct {
		Name               string
		Judge              string
		Recruitment        bool
		LastAdventure      string
		Timekeeping        string
		Characters         []types.CharacterRecord
		Adventures         []types.AdventureRecord
		MainPath           path
		CharacterPath      path
		AdventurePath      path
		NumberOfCharacters int
		NumberOfAdventures int
	}{
		Name:               campaign.Name,
		Judge:              campaign.Judge,
		Recruitment:        campaign.Recruitment,
		LastAdventure:      campaign.LastAdventure.Format("2006-01-02"),
		Timekeeping:        campaign.Timekeeping,
		Characters:         campaign.Characters,
		Adventures:         campaign.Adventures,
		MainPath:           newPhysicalCampaignPath("", campaign.Id),
		CharacterPath:      newPhysicalCampaignPath("/characters", campaign.Id),
		AdventurePath:      newPhysicalAdventurePath("/adventures", campaign.Id),
		NumberOfCharacters: len(campaign.Characters),
		NumberOfAdventures: len(campaign.Adventures),
	}
	output, err := ca.renderer.RenderPage("campaignPage.html", pdata)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, output, err)
	}
	w.Write([]byte(output))
}

func (ca CampaignPage) newCharacter(w http.ResponseWriter, r *http.Request) {
	cid := ca.mustExtractCampaignId(w, r)
	classOptions, err := ca.campaignService.GetClassOptionsForCampaign(cid)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	pdata := struct {
		ClassOptions []types.CampaignClassOption
		DeletePath   path
	}{
		ClassOptions: classOptions,
		DeletePath:   newPhysicalCampaignPath("/new-character", cid),
	}
	output, err := ca.renderer.RenderPartial("characterDetailsTableRow.html", pdata)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, output, err)
	}
	w.Write([]byte(output))

}

func (ca CampaignPage) openCharacterEditor(w http.ResponseWriter, r *http.Request) {
	campaignId := ca.mustExtractCampaignId(w, r)
	characters, err := ca.characterService.GetCharactersForCampaign(types.NewCampaign(campaignId))
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	classOptions, err := ca.campaignService.GetClassOptionsForCampaign(campaignId)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	pdata := struct {
		Characters       []types.CharacterRecord
		ClassOptions     []types.CampaignClassOption
		Errors           []formError
		Path             path
		NewCharacterPath path
		DeletePath       path
	}{
		Characters:       characters,
		ClassOptions:     classOptions,
		Path:             newPhysicalCampaignPath("/characters", campaignId),
		NewCharacterPath: newPhysicalCampaignPath("/new-character", campaignId),
		Errors:           make([]formError, len(characters)),
	}

	output, err := ca.renderer.RenderEditor("characterDetailsTableEdit.html", pdata)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))

}

func (ca CampaignPage) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	charId := r.URL.Query()["char-id"][0]
	pdata := struct {
		Id string
	}{
		Id: charId,
	}
	output, err := ca.renderer.RenderPartial("deletedCharacterFormData.html", pdata)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))

}
func (ca CampaignPage) DeleteUnSavedCharacter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (ca CampaignPage) displayCharacterList(w http.ResponseWriter, r *http.Request) {
	campaignId := ca.mustExtractCampaignId(w, r)
	characters, err := ca.characterService.GetCharactersForCampaign(types.NewCampaign(campaignId))
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	pdata := struct {
		Characters         []types.CharacterRecord
		NumberOfCharacters int
		CharacterPath      path
	}{
		Characters:         characters,
		CharacterPath:      newPhysicalCampaignPath("/characters", campaignId),
		NumberOfCharacters: len(characters),
	}

	output, err := ca.renderer.RenderPartial("characterList.html", pdata)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	w.Write([]byte(output))
}

func (ca CampaignPage) saveCharacters(w http.ResponseWriter, r *http.Request) {

}

/* UTILITY FUNCTIONS */
func (ca CampaignPage) mustExtractCampaignId(w http.ResponseWriter, r *http.Request) int {
	id := r.PathValue("campaignId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	return aId
}
