package archivist

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/renderer"
	"github.com/floodedrealms/borderland-keep/types"
)

const baseCampaignPath = "/campaign"

type CampaignPage struct {
	campaignService  services.CampaignService
	characterService services.CharacterService
	adventureService services.AdventureService
	renderer         renderer.Renderer
	CharacterPath    path
}

func NewCampaignPage(cs services.CampaignService, ch services.CharacterService, as services.AdventureService, r renderer.Renderer) *CampaignPage {
	return &CampaignPage{
		campaignService:  cs,
		characterService: ch,
		adventureService: as,
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

func (c CampaignPage) RegisterRoutes(m *http.ServeMux, g guardsman.Guardsman) {
	mainPath := newCampaignPathToRegister("")
	characterPath := newCampaignPathToRegister("/characters")
	newCharacterPath := newCampaignPathToRegister("/new-character")
	adventurePath := newCampaignPathToRegister("/adventures")

	m.HandleFunc(mainPath.Display, g.UserLoggedInAndHasEditAccessToCampaign(c.CampaignOverview))
	m.HandleFunc(mainPath.Edit, c.HandleEditCampaign)
	m.HandleFunc("PATCH "+mainPath.Display, g.UserLoggedInAndHasEditAccessToCampaign(c.HandleEditCampaign))

	m.HandleFunc("GET "+characterPath.Edit, g.UserMustBeLoggedIn(c.openCharacterEditor))
	m.HandleFunc("GET "+characterPath.Display, c.displayCharacterList)
	m.HandleFunc("POST "+characterPath.Display, g.UserMustBeLoggedIn(c.saveCharacters))
	m.HandleFunc("DELETE "+characterPath.Display, g.UserMustBeLoggedIn(c.DeleteCharacter))

	m.HandleFunc("GET "+newCharacterPath.Display, c.newCharacter)
	m.HandleFunc("DELETE "+newCharacterPath.Display, c.DeleteUnSavedCharacter)

	m.HandleFunc("POST "+adventurePath.Display, c.createNewAdventure)
}

func (ca CampaignPage) CampaignPageForUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	switch r.Method {
	case http.MethodPost:
		campaign, _ := ca.campaignService.CreateCampaignForUser(userId)
		newPage := fmt.Sprintf("/user/1/campaign/%d?edit=true", campaign.Id)
		w.Header().Set("HX-Redirect", newPage)
		w.Header().Set("location", newPage)
		w.WriteHeader(http.StatusNoContent)

	case http.MethodGet:
		campaignId, _ := strconv.Atoi(r.PathValue("campaignId"))
		editor := r.URL.Query()["edit"][0]
		openEditor := false
		if editor == "true" {
			openEditor = true
		}
		campaign, _ := ca.campaignService.CampaignSummary(campaignId)
		pageData := struct {
			types.CampaignRecord
			OpenCampaignEditor bool
			UserId             string
			CampaignId         int
		}{
			*campaign,
			openEditor,
			userId,
			campaign.Id,
		}
		lang := util.ExtractLangageCookie(r)
		loggedIn, edit := ExtractGuardsmanHeaders(r)
		output, _ := ca.renderer.RenderPage("campaignPage.html", pageData, lang, loggedIn, edit)
		w.Write([]byte(output))

	case http.MethodPatch:
		campaignId := ca.mustExtractCampaignId(w, r)
		r.ParseForm()
		name := r.Form["campaign-name"][0]
		judge := r.Form["campaign-judge"][0]
		rectuitingString := r.Form["campaign-recruiting"][0]
		isRecruiting := false
		if rectuitingString == "yes" {
			isRecruiting = true
		}
		timekeepingString := r.Form["campaign-timekeeping"][0]
		err := ca.campaignService.UpdateCampaignDetails(campaignId, name, judge, timekeepingString, isRecruiting)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		campaign, err := ca.campaignService.GetCampaign(campaignId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		output, err := ca.renderer.RenderPartial("campaignDetails.html", campaign)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))

	}
}

func (ca CampaignPage) HandleEditCampaign(w http.ResponseWriter, r *http.Request) {
	loggedIn := r.Header.Get(http.CanonicalHeaderKey(guardsman.LoggedInHeader)) == "true"
	canEdit := r.Header.Get(http.CanonicalHeaderKey(guardsman.EditAccessHeader)) == "true"
	switch r.Method {
	case http.MethodGet:
		campaignId, _ := util.ExtractCampaignId(r)
		campaign, _ := ca.campaignService.CampaignSummary(campaignId)
		pageData := struct {
			types.CampaignRecord
			CampaignId int
			MainPath   path
		}{
			*campaign,
			campaignId,
			newPhysicalCampaignPath("", campaignId),
		}
		lang := util.ExtractLangageCookie(r)
		loggedIn, edit := ExtractGuardsmanHeaders(r)
		output, err := ca.renderer.RenderPage("campaignEditor.html", pageData, lang, loggedIn, edit)
		if err != nil {
			ca.renderer.MustRenderErrorPage(w, output, err)
		}
		w.Write([]byte(output))

	case http.MethodPatch:
		campaignId := ca.mustExtractCampaignId(w, r)
		r.ParseForm()
		name := r.Form["campaign-name"][0]
		judge := r.Form["campaign-judge"][0]
		rectuitingString := r.Form["campaign-recruiting"][0]
		isRecruiting := false
		if rectuitingString == "yes" {
			isRecruiting = true
		}
		timekeepingString := r.Form["campaign-timekeeping"][0]
		err := ca.campaignService.UpdateCampaignDetails(campaignId, name, judge, timekeepingString, isRecruiting)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		campaign, err := ca.campaignService.GetCampaign(campaignId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pdata := struct {
			types.CampaignRecord
			CampaignId    int
			MainPath      path
			HasEditAccess bool
			User          guardsman.WebUser
		}{
			*campaign,
			campaignId,
			newPhysicalCampaignPath("", campaignId),
			canEdit,
			guardsman.WebUser{LoggedIn: loggedIn},
		}
		output, err := ca.renderer.RenderPartial("campaignDetails.html", pdata)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))

	}
}

func (ca CampaignPage) CampaignOverview(w http.ResponseWriter, r *http.Request) {
	//applyCorsHeaders(w)
	id := ca.mustExtractCampaignId(w, r)
	user, edit := ExtractGuardsmanHeaders(r)
	campaign, err := ca.campaignService.CampaignSummary(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	campaign.Adventures, err = ca.campaignService.CampaignAdventuresSummary(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	campaign.Characters, err = ca.characterService.GetCharacterCampaignSummary(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pdata := struct {
		types.CampaignRecord
		OpenCampaignEditor bool
		HasEditAccess      bool
		IsIndex            bool
		User               guardsman.WebUser
		MainPath           path
		CharacterPath      path
		AdventurePath      path
	}{
		CampaignRecord:     *campaign,
		OpenCampaignEditor: false,
		HasEditAccess:      edit,
		IsIndex:            false,
		MainPath:           newPhysicalCampaignPath("", campaign.Id),
		CharacterPath:      newPhysicalCampaignPath("/characters", campaign.Id),
		AdventurePath:      newPhysicalCampaignPath("/adventures", campaign.Id),
		User:               user,
	}
	lang := util.ExtractLangageCookie(r)
	output, err := ca.renderer.RenderPage("campaignPage.html", pdata, lang, user, edit)
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
	characters, err := ca.characterService.GetCharactersForCampaign(campaignId)
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
	characters, err := ca.characterService.GetCharacterCampaignSummary(campaignId)
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
	cid := ca.mustExtractCampaignId(w, r)
	cta, ctu, ctd, _, err := extractCharacterDataFromForm(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	err = ca.characterService.CreateCharactersForCampaign(cid, cta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	err = ca.characterService.DeleteCharacters(ctd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	err = ca.characterService.UpdateCharacters(ctu)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	ca.displayCharacterList(w, r)
}

func (ca CampaignPage) createNewAdventure(w http.ResponseWriter, r *http.Request) {
	campaignId := ca.mustExtractCampaignId(w, r)
	na, err := ca.adventureService.CreateNewAdventureRecordForCampaign(campaignId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	newPage := fmt.Sprintf("/adventure/%d", na.Id)
	w.Header().Set("HX-Redirect", newPage)
	w.Header().Set("location", newPage)
	w.WriteHeader(http.StatusNoContent)
}

/* UTILITY FUNCTIONS */

func extractCharacterDataFromForm(r *http.Request) (charactersToAdd, charactersToUpdate []types.CharacterRecord, deletedIds []string, formErrors []formError, err error) {
	err = r.ParseForm()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	charactersToAdd, _, _ = extractNewCharacters(*r)
	charactersToUpdate, _, _ = extractCharactersToUpdate(*r)
	deletedIds = extractCharactersToDelete(*r)
	return charactersToAdd, charactersToUpdate, deletedIds, formErrors, err
}

func extractNewCharacters(r http.Request) (charsToAdd []types.CharacterRecord, fErrs []formError, err error) {
	newNames := r.Form["new-character-name"]
	newClasses := r.Form["new-character-class"]
	newPReqs := r.Form["new-character-preq"]
	for i, name := range newNames {
		if name == "" {
			fErrs = append(fErrs, formError{Name: "Name must be more than 0 characters"})
			continue
		}
		ClassId, _ := strconv.Atoi(newClasses[i])
		preqPercent, _ := strconv.Atoi(newPReqs[i])
		charsToAdd = append(charsToAdd, types.CharacterRecord{
			Name:            name,
			ClassId:         ClassId,
			PrimeReqPercent: preqPercent,
		})
		fErrs = append(fErrs, formError{}) // keep arrays parrallel

	}
	return charsToAdd, fErrs, nil
}

func extractCharactersToDelete(r http.Request) (charsToDelete []string) {
	return r.Form["deleted-character-id"]
}

func extractCharactersToUpdate(r http.Request) (charsToUpdate []types.CharacterRecord, fErrs []formError, err error) {
	ids := r.Form["character-id"]
	for _, id := range ids {
		fieldname := "character-status-" + id
		statusId, _ := strconv.Atoi(r.Form[fieldname][0])
		idInt, _ := strconv.Atoi(id)
		charsToUpdate = append(charsToUpdate, types.CharacterRecord{
			Id:       idInt,
			StatusId: statusId,
		})
		fErrs = append(fErrs, formError{}) // keep arrays parrallel

	}
	return charsToUpdate, fErrs, nil
}

func (ca CampaignPage) mustExtractCampaignId(w http.ResponseWriter, r *http.Request) int {
	id := r.PathValue("campaignId")
	aId, err := strconv.Atoi(id)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, "", err)
	}
	return aId
}
