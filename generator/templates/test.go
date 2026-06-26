package templates

func ExtTestTemplate(name, extType, moduleName string) string {
	return `package tests

import (
	"testing"

	module "{{ .PackagePath }}"
)

func TestExtensionMetadata(t *testing.T) {
	ext := module.New()
	if ext.Name() != "{{ .Name }}" {
		t.Fatalf("expected extension name %q, got %q", "{{ .Name }}", ext.Name())
	}
	if ext.Version() == "" {
		t.Fatal("expected extension version")
	}

	meta := ext.GetMetadata()
	if meta.Name != ext.Name() {
		t.Fatalf("expected metadata name %q, got %q", ext.Name(), meta.Name)
	}
	if meta.Version != ext.Version() {
		t.Fatalf("expected metadata version %q, got %q", ext.Version(), meta.Version)
	}
}
`
}

func HandlerTestTemplate(name, extType, moduleName string) string {
	return `package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"{{ .PackagePath }}/handler"
	"{{ .PackagePath }}/service"
	"{{ .PackagePath }}/structs"
)

func TestHandlerItemRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := service.NewWithRepository(nil, newFakeRepository())
	h := handler.New(svc)

	router := gin.New()
	h.RegisterRoutes(router.Group("/api/{{ .Name }}"))

	createBody := []byte(` + "`" + `{"name":"Handler Item","code":"HANDLER"}` + "`" + `)
	createRecorder := httptest.NewRecorder()
	createRequest := httptest.NewRequest(http.MethodPost, "/api/{{ .Name }}/items", bytes.NewReader(createBody))
	createRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createRecorder, createRequest)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRecorder.Code, createRecorder.Body.String())
	}

	var created structs.ItemResponse
	if err := json.Unmarshal(createRecorder.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" || created.Name != "Handler Item" {
		t.Fatalf("unexpected created item: %+v", created)
	}

	listRecorder := httptest.NewRecorder()
	listRequest := httptest.NewRequest(http.MethodGet, "/api/{{ .Name }}/items?page_size=10&page_num=1", nil)
	router.ServeHTTP(listRecorder, listRequest)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d: %s", http.StatusOK, listRecorder.Code, listRecorder.Body.String())
	}

	updateBody := []byte(` + "`" + `{"name":"Updated Handler Item","status":"disabled"}` + "`" + `)
	updateRecorder := httptest.NewRecorder()
	updateRequest := httptest.NewRequest(http.MethodPut, "/api/{{ .Name }}/items/"+created.ID, bytes.NewReader(updateBody))
	updateRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(updateRecorder, updateRequest)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d: %s", http.StatusOK, updateRecorder.Code, updateRecorder.Body.String())
	}

	deleteRecorder := httptest.NewRecorder()
	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/{{ .Name }}/items/"+created.ID, nil)
	router.ServeHTTP(deleteRecorder, deleteRequest)
	if deleteRecorder.Code != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d: %s", http.StatusOK, deleteRecorder.Code, deleteRecorder.Body.String())
	}
}
`
}

func ServiceTestTemplate(name, extType, moduleName string) string {
	return `package tests

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"{{ .PackagePath }}/data/repository"
	"{{ .PackagePath }}/service"
	"{{ .PackagePath }}/structs"
)

func TestServiceItemLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := service.NewWithRepository(nil, newFakeRepository())

	created, err := svc.Create(ctx, &structs.CreateItemRequest{
		Name:   "Service Item",
		Code:   "SVC",
		Status: "active",
	})
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created item ID")
	}

	updated, err := svc.Update(ctx, &structs.UpdateItemRequest{
		ID:     created.ID,
		Name:   "Updated Service Item",
		Status: "disabled",
	})
	if err != nil {
		t.Fatalf("update item: %v", err)
	}
	if updated.Name != "Updated Service Item" || updated.Status != "disabled" {
		t.Fatalf("unexpected updated item: %+v", updated)
	}

	items, total, err := svc.List(ctx, &structs.ListItemsRequest{PageSize: 10, PageNum: 1})
	if err != nil {
		t.Fatalf("list items: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected one listed item, got total=%d len=%d", total, len(items))
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete item: %v", err)
	}
	if _, err := svc.Get(ctx, created.ID); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestServiceValidation(t *testing.T) {
	svc := service.NewWithRepository(nil, newFakeRepository())
	if _, err := svc.Create(context.Background(), &structs.CreateItemRequest{}); !errors.Is(err, service.ErrInvalidRequest) {
		t.Fatalf("expected invalid request error, got %v", err)
	}
}

type fakeRepository struct {
	mu     sync.RWMutex
	nextID int
	items  map[string]*structs.Item
}

func newFakeRepository() repository.RepositoryInterface {
	return &fakeRepository{
		items: make(map[string]*structs.Item),
	}
}

func (r *fakeRepository) Create(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.Name) == "" {
		return nil, repository.ErrInvalidItem
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	clone := cloneItem(item)
	if clone.ID == "" {
		clone.ID = fmt.Sprintf("item-%d", r.nextID)
	}
	if clone.Status == "" {
		clone.Status = "active"
	}
	now := time.Now().UTC()
	if clone.CreatedAt.IsZero() {
		clone.CreatedAt = now
	}
	if clone.UpdatedAt.IsZero() {
		clone.UpdatedAt = now
	}
	r.items[clone.ID] = clone
	return cloneItem(clone), nil
}

func (r *fakeRepository) Get(ctx context.Context, id string) (*structs.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok || item.DeletedAt != nil {
		return nil, repository.ErrNotFound
	}
	return cloneItem(item), nil
}

func (r *fakeRepository) Update(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.ID) == "" {
		return nil, repository.ErrInvalidItem
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.items[item.ID]
	if !ok || current.DeletedAt != nil {
		return nil, repository.ErrNotFound
	}
	next := cloneItem(item)
	next.CreatedAt = current.CreatedAt
	next.UpdatedAt = time.Now().UTC()
	r.items[next.ID] = next
	return cloneItem(next), nil
}

func (r *fakeRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.items[id]
	if !ok || item.DeletedAt != nil {
		return repository.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	return nil
}

func (r *fakeRepository) List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*structs.Item, 0, len(r.items))
	for _, item := range r.items {
		if item.DeletedAt == nil {
			items = append(items, cloneItem(item))
		}
	}
	return items, int64(len(items)), nil
}

func (r *fakeRepository) Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var total int64
	for _, item := range r.items {
		if item.DeletedAt == nil {
			total++
		}
	}
	return total, nil
}

func cloneItem(item *structs.Item) *structs.Item {
	if item == nil {
		return nil
	}
	clone := *item
	if item.Extras != nil {
		clone.Extras = make(map[string]any, len(item.Extras))
		for key, value := range item.Extras {
			clone.Extras[key] = value
		}
	}
	return &clone
}
`
}
