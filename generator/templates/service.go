package templates

import "fmt"

func ServiceTemplate(name, extType, moduleName string) string {
	return fmt.Sprintf(`package service

import (
	"context"
	"time"

	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/logging/logger"
	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/data/repository"
	"{{ .PackagePath }}/structs"
)

// ServiceInterface represents the service interface.
type ServiceInterface interface {
	Create(ctx context.Context, req *structs.CreateItemRequest) (*structs.ItemResponse, error)
	Get(ctx context.Context, id string) (*structs.ItemResponse, error)
	Update(ctx context.Context, req *structs.UpdateItemRequest) (*structs.ItemResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, req *structs.ListItemsRequest) ([]*structs.ItemResponse, int64, error)
}

// Service represents the %s service.
type Service struct {
	conf *config.Config
	repo repository.RepositoryInterface
}

// New creates a new service.
func New(conf *config.Config, d *data.Data) ServiceInterface {
	return &Service{
		conf: conf,
		repo: repository.New(d),
	}
}

// Create creates a new item.
func (s *Service) Create(ctx context.Context, req *structs.CreateItemRequest) (*structs.ItemResponse, error) {
	item := &structs.Item{
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		Status:    req.Status,
		Desc:      req.Desc,
		Extras:    req.Extras,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdItem, err := s.repo.Create(ctx, item)
	if err != nil {
		logger.Errorf(ctx, "Failed to create item: %%v", err)
		return nil, err
	}

	return s.toResponse(createdItem), nil
}

// Get retrieves an item by ID.
func (s *Service) Get(ctx context.Context, id string) (*structs.ItemResponse, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "Failed to get item: %%v", err)
		return nil, err
	}

	return s.toResponse(item), nil
}

// Update updates an existing item.
func (s *Service) Update(ctx context.Context, req *structs.UpdateItemRequest) (*structs.ItemResponse, error) {
	item, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Code != "" {
		item.Code = req.Code
	}
	// ... update other fields ...
	item.UpdatedAt = time.Now()

	updatedItem, err := s.repo.Update(ctx, item)
	if err != nil {
		logger.Errorf(ctx, "Failed to update item: %%v", err)
		return nil, err
	}

	return s.toResponse(updatedItem), nil
}

// Delete deletes an item by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Errorf(ctx, "Failed to delete item: %%v", err)
		return err
	}
	return nil
}

// List lists items.
func (s *Service) List(ctx context.Context, req *structs.ListItemsRequest) ([]*structs.ItemResponse, int64, error) {
	items, count, err := s.repo.List(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to list items: %%v", err)
		return nil, 0, err
	}

	responses := make([]*structs.ItemResponse, len(items))
	for i, item := range items {
		responses[i] = s.toResponse(item)
	}

	return responses, count, nil
}

func (s *Service) toResponse(item *structs.Item) *structs.ItemResponse {
	if item == nil {
		return nil
	}
	return &structs.ItemResponse{
		ID:        item.ID,
		Name:      item.Name,
		Code:      item.Code,
		Type:      item.Type,
		Status:    item.Status,
		Desc:      item.Desc,
		CreatedBy: item.CreatedBy,
		CreatedAt: item.CreatedAt,
		UpdatedBy: item.UpdatedBy,
		UpdatedAt: item.UpdatedAt,
		Extras:    item.Extras,
	}
}
`, name)
}
