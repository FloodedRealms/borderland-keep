package webapp

import (
	"net/http"

	"github.com/floodedrealms/adventure-archivist/types"
)

type Modals struct {
	r Renderer
}

func NewModals(r Renderer) *Modals {
	return &Modals{r: r}
}

func (m Modals) LootModal(w http.ResponseWriter, r *http.Request) {
	in := r.URL.Query()["type"][0]
	l := types.GenericLootType(in)
	var (
		rendered string
	)
	switch l {
	case types.GemLoot:
		rendered, _ = m.r.Render("lootModal.html", types.Gem{Number: 0, XPValue: 0, Description: "", LootType: types.GemLoot})
	}
	w.Write([]byte(rendered))
}
