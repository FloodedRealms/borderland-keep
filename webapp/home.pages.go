package webapp

import (
	"net/http"
)

type HomePage struct {
	renderer Renderer
}

func NewHomePage(r Renderer) *HomePage {
	return &HomePage{renderer: r}
}

func (h HomePage) About(w http.ResponseWriter, r *http.Request) {
	output, err := h.renderer.RenderPageWithNoData("about.html")
	if err != nil {
		h.renderer.MustRenderErrorPage(w, "error.html", err)
	}
	w.Write([]byte(output))
}
