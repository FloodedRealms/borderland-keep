package archivist

import (
	"net/http"
	"strconv"

	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
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
	lang := util.ExtractLangageCookie(r)
	loggedIn, edit := ExtractGuardsmanHeaders(r)
	output, err := h.renderer.RenderPageWithNoData("about.html", lang, loggedIn, edit)
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}

func (h HomePage) GuildLanding(w http.ResponseWriter, r *http.Request) {
	lang := util.ExtractLangageCookie(r)
	loggedIn, edit := ExtractGuardsmanHeaders(r)
	output, err := h.renderer.RenderPageWithNoData("guild.html", lang, loggedIn, edit)
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}

func (h HomePage) TavernLanding(w http.ResponseWriter, r *http.Request) {
	lang := util.ExtractLangageCookie(r)
	loggedIn, edit := ExtractGuardsmanHeaders(r)
	output, err := h.renderer.RenderPageWithNoData("tavern.html", lang, loggedIn, edit)
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
	lang := util.ExtractLangageCookie(r)
	loggedIn, edit := ExtractGuardsmanHeaders(r)
	output, err := h.renderer.RenderPage("index.html", pdata, lang, loggedIn, edit)
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
	lang := util.ExtractLangageCookie(r)
	loggedIn, edit := ExtractGuardsmanHeaders(r)
	output, err := h.renderer.RenderPage("campaigns.html", pdata, lang, loggedIn, edit)
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}

func (h HomePage) MyCampaigns(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	lang := util.ExtractLangageCookie(r)
	loggedIn, edit := ExtractGuardsmanHeaders(r)
	campaigns := h.cs.CampaignsForUser(userId)
	pdata := struct {
		Campaigns      []types.CampaignRecord
		EndOfCampaigns bool
		IsIndex        bool
	}{
		Campaigns:      campaigns,
		EndOfCampaigns: len(campaigns) < 10,
		IsIndex:        false,
	}

	output, err := h.renderer.RenderPage("myCampaigns.html", pdata, lang, loggedIn, edit)
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

func (g HomePage) DisplayLoginPage(w http.ResponseWriter, r *http.Request) {
	lang := util.ExtractLangageCookie(r)
	login, edt := false, false
	output, err := g.renderer.RenderPage("login.html", nil, lang, login, edt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}

func (g HomePage) ErrorPage(w http.ResponseWriter, r *http.Request) {
	msg := r.Header.Get("x-borderland-keep-error")
	pdata := struct {
		Message string
	}{
		msg,
	}
	output, err := g.renderer.RenderPage("errorPage.html", pdata, "english", nil, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}
