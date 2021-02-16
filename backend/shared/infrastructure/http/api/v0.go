package api

import (
	cgt_routes "emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/http/routes"
	identity_routes "emersonargueta/m/v1/modules/identity/infrastructure/http/routes"

	"net/http"
)

// BaseHandler is a collection of all the api handlers.
type BaseHandler struct {
	BasePath string
	*http.ServeMux
}

// NewBaseHandler with basePath
func NewBaseHandler(basePath string) *BaseHandler {
	bh := new(BaseHandler)
	bh.BasePath = basePath

	mux := http.NewServeMux()
	mux.Handle(basePath, bh)

	productFinderHandler := cgt_routes.NewCommunityGoalTrackerHandler(basePath)
	mux.Handle(basePath+cgt_routes.CommunitygoalTrackerURLPrefix+"/", productFinderHandler)

	identityHandler := identity_routes.NewIdentityHandler(basePath)
	mux.Handle(basePath+identity_routes.IdentityURLPrefix+"/", identityHandler)

	bh.ServeMux = mux

	return bh

}
