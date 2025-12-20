package templates

import "fmt"

func HandlerTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/net/resp"
	"{{ .PackagePath }}/service"
	"{{ .PackagePath }}/structs"
)

// Handler represents the %s handler.
type Handler struct {
	s service.ServiceInterface
}

// New creates a new handler.
func New(s service.ServiceInterface) *Handler {
	return &Handler{
		s: s,
	}
}

// RegisterRoutes registers the HTTP routes for this handler.
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	items := r.Group("/items")
	{
		items.POST("", h.Create)
		items.GET("/:id", h.Get)
		items.PUT("", h.Update)
		items.DELETE("/:id", h.Delete)
		items.GET("", h.List)
	}
}

// Create handles the creation of a new item.
func (h *Handler) Create(c *gin.Context) {
	var req structs.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
		return
	}

	res, err := h.s.Create(c.Request.Context(), &req)
	if err != nil {
		resp.Fail(c.Writer, resp.InternalServerError(err.Error()))
		return
	}

	resp.Success(c.Writer, res)
}

// Get handles retrieving an item by ID.
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		resp.Fail(c.Writer, resp.BadRequest("ID is required"))
		return
	}

	res, err := h.s.Get(c.Request.Context(), id)
	if err != nil {
		resp.Fail(c.Writer, resp.InternalServerError(err.Error()))
		return
	}

	resp.Success(c.Writer, res)
}

// Update handles updating an item.
func (h *Handler) Update(c *gin.Context) {
	var req structs.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
		return
	}

	res, err := h.s.Update(c.Request.Context(), &req)
	if err != nil {
		resp.Fail(c.Writer, resp.InternalServerError(err.Error()))
		return
	}

	resp.Success(c.Writer, res)
}

// Delete handles deleting an item by ID.
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		resp.Fail(c.Writer, resp.BadRequest("ID is required"))
		return
	}

	if err := h.s.Delete(c.Request.Context(), id); err != nil {
		resp.Fail(c.Writer, resp.InternalServerError(err.Error()))
		return
	}

	resp.Success(c.Writer, nil)
}

// List handles listing items.
func (h *Handler) List(c *gin.Context) {
	var req structs.ListItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.Fail(c.Writer, resp.BadRequest(err.Error()))
		return
	}

	items, count, err := h.s.List(c.Request.Context(), &req)
	if err != nil {
		resp.Fail(c.Writer, resp.InternalServerError(err.Error()))
		return
	}

	resp.Success(c.Writer, map[string]interface{}{
		"items": items,
		"total": count,
	})
}
`, name)
}
