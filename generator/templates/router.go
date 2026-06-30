package templates

func RouterTemplate() string {
	return `package router

import (
	"github.com/gin-gonic/gin"

	"{{ .PackagePath }}/handler"
)

// Router registers HTTP routes for the module.
type Router struct {
	h *handler.Handler
}

// New creates a router.
func New(h *handler.Handler) *Router {
	return &Router{h: h}
}

// RegisterRoutes registers module routes under the module path.
func (r *Router) RegisterRoutes(root *gin.RouterGroup) {
	if r == nil || r.h == nil {
		return
	}
	r.h.RegisterRoutes(root.Group("/{{ .Name }}"))
}
`
}
