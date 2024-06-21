package webapp

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/floodedrealms/adventure-archivist/services"
)

type CampaignPage struct {
	ad  services.AdventureService
	chs services.CharacterService
	cs  services.CampaignService
}

func NewCampaignPage(cs services.CampaignService, chs services.CharacterService, ad services.AdventureService) *CampaignPage {
	return &CampaignPage{
		ad:  ad,
		cs:  cs,
		chs: chs,
	}
}

func (p CampaignPage) Page(w http.ResponseWriter, r *http.Request) {

	campaigns, _ := p.cs.ListCampaigns()
	funcMap := template.FuncMap{
		"link": func(i int) string { return fmt.Sprintf("/campaigns/%d", i) },
	}
	tmplFile := "./webapp/views/index.html"
	tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index", campaigns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
