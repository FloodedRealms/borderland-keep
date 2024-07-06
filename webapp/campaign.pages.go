package webapp

import (
	"net/http"

	"github.com/floodedrealms/adventure-archivist/internal/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

type CampaignPage struct {
	campaignService  services.CampaignService
	characterService services.CharacterService
	renderer         Renderer
}

func NewCampaignPage(cs services.CampaignService, ch services.CharacterService, r Renderer) *CampaignPage {
	return &CampaignPage{
		campaignService:  cs,
		characterService: ch,
		renderer:         r,
	}

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
	ca.renderCamapaignPage(w, *campaign)
	//sendGoodResponseWithObject(w, campaign)
}

func (ca CampaignPage) renderCamapaignPage(w http.ResponseWriter, data types.CampaignRecord) {
	output, err := ca.renderer.RenderPage("campaignPage.html", data)
	if err != nil {
		ca.renderer.MustRenderErrorPage(w, output, err)
	}
	w.Write([]byte(output))
}
