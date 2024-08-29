package webapp

import (
	"net/http"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/internal/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

type HomePage struct {
	renderer Renderer
	cs       services.CampaignService
}

func NewHomePage(r Renderer, cs services.CampaignService) *HomePage {
	return &HomePage{renderer: r, cs: cs}
}

func (h HomePage) About(w http.ResponseWriter, r *http.Request) {
	output, err := h.renderer.RenderPageWithNoData("about.html")
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}

func (h HomePage) GuildLanding(w http.ResponseWriter, r *http.Request) {
	output, err := h.renderer.RenderPageWithNoData("guild.html")
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}

func (h HomePage) TavernLanding(w http.ResponseWriter, r *http.Request) {
	output, err := h.renderer.RenderPageWithNoData("tavern.html")
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))

}

func (h HomePage) Index(w http.ResponseWriter, r *http.Request) {
	output, err := h.renderer.RenderPageWithNoData("index.html")
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))

}

func (h HomePage) Campaigns(w http.ResponseWriter, r *http.Request) {
	campaigns := h.cs.TenMostRecentlyActiveCampaigns(1)
	pdata := struct {
		Campaigns      []types.CampaignRecord
		EndOfCampaigns bool
	}{
		Campaigns:      campaigns,
		EndOfCampaigns: len(campaigns) < 10,
	}

	output, err := h.renderer.RenderPage("campaigns.html", pdata)
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}

func (h HomePage) LoadNextCampaignSet(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query()["page"][0]
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	campaigns := h.cs.TenMostRecentlyActiveCampaigns(pageNumber)
	pdata := struct {
		Campaigns      []types.CampaignRecord
		Page           int
		EndOfCampaigns bool
	}{
		Campaigns:      campaigns,
		Page:           pageNumber + 1,
		EndOfCampaigns: len(campaigns) < 10,
	}

	output, err := h.renderer.RenderPartial("campaignSections.html", pdata)
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))

}
