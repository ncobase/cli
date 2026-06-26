package templates

func ServiceTemplate(name, extType, moduleName string) string {
	return `package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ncobase/ncore/config"
	"github.com/ncobase/ncore/logging/logger"
	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/data/repository"
	"{{ .PackagePath }}/structs"
)

var (
	// ErrInvalidRequest is returned when a request cannot be processed.
	ErrInvalidRequest = errors.New("invalid item request")
	// ErrNotFound is returned when the requested item does not exist.
	ErrNotFound = repository.ErrNotFound
)

// ServiceInterface defines item service behavior.
type ServiceInterface interface {
	Create(ctx context.Context, req *structs.CreateItemRequest) (*structs.ItemResponse, error)
	Get(ctx context.Context, id string) (*structs.ItemResponse, error)
	Update(ctx context.Context, req *structs.UpdateItemRequest) (*structs.ItemResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, req *structs.ListItemsRequest) ([]*structs.ItemResponse, int64, error)
}

// Service coordinates item business logic.
type Service struct {
	conf *config.Config
	repo repository.RepositoryInterface
	now  func() time.Time
}

// New creates a service backed by the generated repository.
func New(conf *config.Config, d *data.Data) ServiceInterface {
	return NewWithRepository(conf, repository.New(d))
}

// NewWithRepository creates a service with an explicit repository.
func NewWithRepository(conf *config.Config, repo repository.RepositoryInterface) ServiceInterface {
	return &Service{
		conf: conf,
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

// Create creates a new item.
func (s *Service) Create(ctx context.Context, req *structs.CreateItemRequest) (*structs.ItemResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request body is required", ErrInvalidRequest)
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidRequest)
	}
	if s.repo == nil {
		return nil, errors.New("item repository is not configured")
	}

	now := s.now()
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}
	item := &structs.Item{
		Name:      name,
		Code:      strings.TrimSpace(req.Code),
		Type:      strings.TrimSpace(req.Type),
		Status:    status,
		Desc:      req.Desc,
		Extras:    req.Extras,
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdItem, err := s.repo.Create(ctx, item)
	if err != nil {
		logger.Errorf(ctx, "Failed to create item: %v", err)
		return nil, err
	}
	return toResponse(createdItem), nil
}

// Get retrieves an item by ID.
func (s *Service) Get(ctx context.Context, id string) (*structs.ItemResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidRequest)
	}
	if s.repo == nil {
		return nil, errors.New("item repository is not configured")
	}

	item, err := s.repo.Get(ctx, id)
	if err != nil {
		logger.Errorf(ctx, "Failed to get item: %v", err)
		return nil, err
	}
	return toResponse(item), nil
}

// Update updates an existing item.
func (s *Service) Update(ctx context.Context, req *structs.UpdateItemRequest) (*structs.ItemResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request body is required", ErrInvalidRequest)
	}
	if strings.TrimSpace(req.ID) == "" {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidRequest)
	}
	if s.repo == nil {
		return nil, errors.New("item repository is not configured")
	}

	item, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if name := strings.TrimSpace(req.Name); name != "" {
		item.Name = name
	}
	if code := strings.TrimSpace(req.Code); code != "" {
		item.Code = code
	}
	if itemType := strings.TrimSpace(req.Type); itemType != "" {
		item.Type = itemType
	}
	if status := strings.TrimSpace(req.Status); status != "" {
		item.Status = status
	}
	if req.Desc != "" {
		item.Desc = req.Desc
	}
	if req.Extras != nil {
		item.Extras = req.Extras
	}
	item.UpdatedAt = s.now()

	updatedItem, err := s.repo.Update(ctx, item)
	if err != nil {
		logger.Errorf(ctx, "Failed to update item: %v", err)
		return nil, err
	}
	return toResponse(updatedItem), nil
}

// Delete deletes an item by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidRequest)
	}
	if s.repo == nil {
		return errors.New("item repository is not configured")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Errorf(ctx, "Failed to delete item: %v", err)
		return err
	}
	return nil
}

// List lists items.
func (s *Service) List(ctx context.Context, req *structs.ListItemsRequest) ([]*structs.ItemResponse, int64, error) {
	if s.repo == nil {
		return nil, 0, errors.New("item repository is not configured")
	}
	if req == nil {
		req = &structs.ListItemsRequest{}
	}
	if req.PageSize < 0 || req.PageNum < 0 {
		return nil, 0, fmt.Errorf("%w: pagination values must be positive", ErrInvalidRequest)
	}

	items, count, err := s.repo.List(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to list items: %v", err)
		return nil, 0, err
	}

	responses := make([]*structs.ItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toResponse(item))
	}
	return responses, count, nil
}

func toResponse(item *structs.Item) *structs.ItemResponse {
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
`
}
