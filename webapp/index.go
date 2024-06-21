package webapp

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/floodedrealms/adventure-archivist/services"
)

type IndexPage struct {
	cs services.CampaignService
}

func NewIndexPage(cs services.CampaignService) *IndexPage {
	return &IndexPage{cs: cs}
}

func (i IndexPage) Index(w http.ResponseWriter, r *http.Request) {

	campaigns, _ := i.cs.ListCampaigns()
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
