package templates

import "fmt"

// ExtTestTemplate generates extension test template
func ExtTestTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package tests

import (
	"testing"
	"github.com/ncobase/ncore/config"
	ext "github.com/ncobase/ncore/extension/types"
	"{{ .PackagePath }}"
)

func TestModuleLifecycle(t *testing.T) {
	m := %s.New()

	t.Run("initialization", func(t *testing.T) {
		// Test Pre-Init
		if err := m.PreInit(); err != nil {
			t.Errorf("PreInit failed: %%v", err)
		}

		// Test Init
		conf := &config.Config{}
		em := &ext.ManagerInterface{}
		if err := m.Init(conf, em); err != nil {
			t.Errorf("Init failed: %%v", err)
		}

		// Test Post-Init
		if err := m.PostInit(); err != nil {
			t.Errorf("PostInit failed: %%v", err)
		}
	})

	t.Run("metadata", func(t *testing.T) {
		meta := m.GetMetadata()
		if meta.Name != "%s" {
			t.Errorf("want name %%s, got %%s", "%s", meta.Name)
		}
	})

	t.Run("cleanup", func(t *testing.T) {
		if err := m.Cleanup(); err != nil {
			t.Errorf("Cleanup failed: %%v", err)
		}
	})
}`, name, name, name)
}

// HandlerTestTemplate generates handler test template
func HandlerTestTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package tests

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"{{ .PackagePath }}/handler"
	"{{ .PackagePath }}/service"
)

func TestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	s := service.New(nil, nil)
	h := handler.New(s)

	// Register routes
	// h.RegisterRoutes(r.Group("/api/v1"))

	t.Run("list items", func(t *testing.T) {
		// Setup route
		r.GET("/items", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/items", nil)

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("want status 200, got %%d", w.Code)
		}
	})
}`)
}

// ServiceTestTemplate generates service test template
func ServiceTestTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package tests

import (
	"testing"
	"context"
	"{{ .PackagePath }}/service"
	"{{ .PackagePath }}/structs"
)

func TestService(t *testing.T) {
	ctx := context.Background()
	s := service.New(nil, nil)

	t.Run("create item", func(t *testing.T) {
		req := &structs.CreateItemRequest{
			Name: "Test Item",
			Code: "TEST001",
		}

		resp, err := s.Create(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %%v", err)
		}
		if resp == nil {
			t.Error("response should not be nil")
		}
		if resp.Name != req.Name {
			t.Errorf("want name %%s, got %%s", req.Name, resp.Name)
		}
	})
}`)
}
