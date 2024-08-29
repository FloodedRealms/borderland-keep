package webapp

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Renderer struct {
	pageTemplates    *template.Template
	partialTemplates *template.Template
	editorTemplates  *template.Template
}

type partialRender struct {
	Rendered string
	E        error
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.mustLoadTemplates()
	return r
}

func (r *Renderer) mustLoadTemplates() {
	wd, _ := os.Getwd()
	pageDir := filepath.Join(wd, "/internal/templates/pages/*.html")
	partialsDir := filepath.Join(wd, "/internal/templates/partials/*.html")
	partialEditorsDir := filepath.Join(wd, "/internal/templates/partials/editors/*.html")

	r.pageTemplates = template.Must(template.ParseGlob(pageDir))
	r.partialTemplates = template.Must(template.ParseGlob(partialsDir))
	r.editorTemplates = template.Must(template.ParseGlob(partialEditorsDir))
}

func (r Renderer) RenderPageWithNoData(tmpl string) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.pageTemplates.ExecuteTemplate(&renderedOutput, tmpl, nil); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

func (r Renderer) RenderPage(tmpl string, data interface{}) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.pageTemplates.ExecuteTemplate(&renderedOutput, tmpl, data); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

func (r Renderer) RenderPartial(tmpl string, data interface{}) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.partialTemplates.ExecuteTemplate(&renderedOutput, tmpl, data); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

func (r Renderer) RenderEditor(tmpl string, data interface{}) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.editorTemplates.ExecuteTemplate(&renderedOutput, tmpl, data); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

// TODO: Probably shouldn't panic on a failed error render
func (r Renderer) MustRenderErrorPage(w http.ResponseWriter, partial string, e error) {
	data := partialRender{
		Rendered: partial,
		E:        e,
	}
	if err := r.pageTemplates.ExecuteTemplate(w, "errorPage.html", data); err != nil {
		panic(err)
	}
}

func (r Renderer) TriggerErrorModal(w http.ResponseWriter, e error) {
	w.Header().Set("HX-Trigger", "gotError")
	w.Header().Set("BK-Error", e.Error())
	w.WriteHeader(http.StatusInternalServerError)
}
