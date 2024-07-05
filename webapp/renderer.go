package webapp

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Renderer struct {
	templates *template.Template
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
	dir := filepath.Join(wd, "/internal/templates/*.html")
	r.templates = template.Must(template.ParseGlob(dir))
}

func (r Renderer) Render(tmpl string, data interface{}) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.templates.ExecuteTemplate(&renderedOutput, tmpl, data); err != nil {
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
	if err := r.templates.ExecuteTemplate(w, "errorPage.html", data); err != nil {
		panic(err)
	}
}
