package renderer

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

type PageData struct {
	Data          interface{}
	Lang          string
	User          interface{}
	HasEditAccess bool
	NaturalDate   func(time.Time) string
	YYYYMMDDDate  func(time.Time) string
	BoolToString  func(bool) string
}

func RenderNaturalDate(d time.Time) string {
	return d.Format("January 02, 2006")
}

func RenderYYYYMMDDDate(d time.Time) string {
	return d.Format("2006-01-02")
}
func RenderBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
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

func (r Renderer) RenderPageWithNoData(tmpl string, language string, user interface{}, canEdit bool) (string, error) {

	var renderedOutput bytes.Buffer
	pdata := PageData{
		Data:          nil,
		User:          user,
		Lang:          language,
		HasEditAccess: canEdit,
		NaturalDate:   RenderNaturalDate,
		YYYYMMDDDate:  RenderYYYYMMDDDate,
		BoolToString:  RenderBool,
	}

	//try to render anglish version, if can't continue to default
	if language == "anglish" {
		err := r.templates.ExecuteTemplate(&renderedOutput, "ang-"+tmpl, pdata)
		if err == nil {
			return renderedOutput.String(), nil
		}
	}
	if err := r.templates.ExecuteTemplate(&renderedOutput, tmpl, pdata); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

func (r Renderer) RenderPage(tmpl string, data interface{}, language string, user interface{}, canEdit bool) (string, error) {
	pdata := PageData{
		Data:          data,
		Lang:          language,
		User:          user,
		HasEditAccess: canEdit,
		NaturalDate:   RenderNaturalDate,
		YYYYMMDDDate:  RenderYYYYMMDDDate,
		BoolToString:  RenderBool,
	}
	var renderedOutput bytes.Buffer
	//try to render anglish version, if can't continue to default
	if language == "anglish" {

		err := r.templates.ExecuteTemplate(&renderedOutput, "ang-"+tmpl, pdata)
		if err == nil {
			return renderedOutput.String(), nil
		}
	}
	if err := r.templates.ExecuteTemplate(&renderedOutput, tmpl, pdata); err != nil {
		return renderedOutput.String(), err
	}
	return renderedOutput.String(), nil
}

// TODO: This is really hacky
func (r Renderer) RenderPartial(tmpl string, data interface{}) (string, error) {
	pdata := PageData{
		Data:          data,
		HasEditAccess: true,
		NaturalDate:   RenderNaturalDate,
		YYYYMMDDDate:  RenderYYYYMMDDDate,
		BoolToString:  RenderBool,
	}
	var renderedOutput bytes.Buffer
	if err := r.templates.ExecuteTemplate(&renderedOutput, tmpl, pdata); err != nil {
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
