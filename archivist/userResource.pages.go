package archivist

import (
	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/renderer"
)

type UserResourcePage struct {
	renderer          renderer.Renderer
	guard             guardsman.Guardsman
	campaignListPath  path
	campaignPath      path
	adventureListPath path
}
