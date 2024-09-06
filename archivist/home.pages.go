package archivist

import (
	"net/http"
	"strconv"

	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/renderer"

	"github.com/floodedrealms/borderland-keep/types"
)

const sessionCookie = "session_token"

type HomePage struct {
	renderer renderer.Renderer
	cs       services.CampaignService
	guard    guardsman.Guardsman
}

func NewHomePage(r renderer.Renderer, cs services.CampaignService, g guardsman.Guardsman) *HomePage {
	return &HomePage{renderer: r, cs: cs, guard: g}
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
	u := h.guard.SimpleLoginCheck(r)
	campaigns := h.cs.TenMostRecentlyActiveCampaigns(1)
	pdata := struct {
		Campaigns      []types.CampaignRecord
		Page           int
		EndOfCampaigns bool
		User           guardsman.WebUser
		IsIndex        bool
	}{
		Campaigns:      campaigns,
		Page:           2,
		EndOfCampaigns: len(campaigns) < 10,
		User:           u,
		IsIndex:        true,
	}
	output, err := h.renderer.RenderPage("index.html", pdata)
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
		User           guardsman.WebUser
	}{
		Campaigns:      campaigns,
		EndOfCampaigns: len(campaigns) < 10,
		User:           guardsman.WebUser{LoggedIn: false},
	}

	output, err := h.renderer.RenderPage("campaigns.html", pdata)
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}

func (h HomePage) MyCampaigns(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	campaigns := h.cs.CampaignsForUser(userId)
	pdata := struct {
		Campaigns      []types.CampaignRecord
		EndOfCampaigns bool
		User           guardsman.WebUser
		IsIndex        bool
	}{
		Campaigns:      campaigns,
		EndOfCampaigns: len(campaigns) < 10,
		User:           guardsman.WebUser{Id: userId, LoggedIn: true},
		IsIndex:        false,
	}

	output, err := h.renderer.RenderPage("myCampaigns.html", pdata)
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
