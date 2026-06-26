package templates

func HandlerTemplate(name, extType, moduleName string) string {
	return `package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/net/resp"
	"{{ .PackagePath }}/service"
	"{{ .PackagePath }}/structs"
)

// Handler handles item HTTP requests.
type Handler struct {
	s service.ServiceInterface
}

// New creates a handler.
func New(s service.ServiceInterface) *Handler {
	return &Handler{s: s}
}

// RegisterRoutes registers item routes.
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	items := r.Group("/items")
	items.POST("", h.Create)
	items.GET("", h.List)
	items.GET("/:id", h.Get)
	items.PUT("/:id", h.Update)
	items.PATCH("/:id", h.Update)
	items.DELETE("/:id", h.Delete)
}

// Create handles item creation.
func (h *Handler) Create(c *gin.Context) {
	var req structs.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
		return
	}

	res, err := h.s.Create(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	resp.WithStatusCode(c.Writer, http.StatusCreated, res)
}

// Get handles item retrieval.
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	res, err := h.s.Get(c.Request.Context(), id)
	if err != nil {
		h.respondError(c, err)
		return
	}

	resp.Success(c.Writer, res)
}

// Update handles item updates.
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req structs.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
		return
	}
	if req.ID == "" {
		req.ID = id
	} else if id != "" && req.ID != id {
		resp.Fail(c.Writer, resp.BadRequest("path id and body id do not match"))
		return
	}

	res, err := h.s.Update(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	resp.Success(c.Writer, res)
}

// Delete handles item deletion.
func (h *Handler) Delete(c *gin.Context) {
	if err := h.s.Delete(c.Request.Context(), c.Param("id")); err != nil {
		h.respondError(c, err)
		return
	}

	resp.Success(c.Writer, gin.H{"deleted": true})
}

// List handles item listing.
func (h *Handler) List(c *gin.Context) {
	var req structs.ListItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
		return
	}

	items, count, err := h.s.List(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	resp.Success(c.Writer, gin.H{
		"items": items,
		"total": count,
	})
}

func (h *Handler) respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidRequest):
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
	case errors.Is(err, service.ErrNotFound):
		resp.Fail(c.Writer, resp.NotFound("item not found"))
	default:
		resp.Fail(c.Writer, resp.InternalServer(err.Error()))
	}
}
`
}
