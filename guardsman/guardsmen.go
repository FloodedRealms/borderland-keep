package guardsman

import (
	"net/http"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/renderer"
)

type Guardsman struct {
	repo        repository.Repository
	userService services.UserService
	renderer    *renderer.Renderer
}

// Ya'll are lucky I didn't name this "RecruitGuardsman"
// renderer can be a pointet because it might not be needed for a given application
func NewGuardsman(r repository.Repository, s services.UserService, renderer *renderer.Renderer) *Guardsman {
	return &Guardsman{
		repo:        r,
		userService: s,
		renderer:    renderer,
	}
}

func (g Guardsman) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/gate", g.DisplayLoginPage)
	router.HandleFunc("POST /gate", g.HandleLogin)
	router.HandleFunc("POST /depart", g.Logout)
}
