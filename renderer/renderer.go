package renderer

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

	pageDir := filepath.Join(wd, "/renderer/templates/")
	var allFiles []string
	err := filepath.WalkDir(pageDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	r.templates = template.Must(template.ParseFiles(allFiles...))
}

func (r Renderer) RenderPageWithNoData(tmpl string) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.templates.ExecuteTemplate(&renderedOutput, tmpl, nil); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

func (r Renderer) RenderPage(tmpl string, data interface{}) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.templates.ExecuteTemplate(&renderedOutput, tmpl, data); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

func (r Renderer) RenderPartial(tmpl string, data interface{}) (string, error) {
	var renderedOutput bytes.Buffer
	if err := r.templates.ExecuteTemplate(&renderedOutput, tmpl, data); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

func (r Renderer) RenderEditor(tmpl string, data interface{}) (string, error) {
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
