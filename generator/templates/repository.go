package templates

func RepositoryTemplate(data *Data) string {
	switch {
	case data.UseEnt:
		return repositoryEntTemplate()
	case data.UseGorm:
		return repositoryGormTemplate()
	case data.UseMongo:
		return repositoryMongoTemplate()
	default:
		return repositoryMemoryTemplate()
	}
}

func repositoryEntTemplate() string {
	return `package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/data/ent"
	itement "{{ .PackagePath }}/data/ent/item"
	"{{ .PackagePath }}/structs"
)

var (
	// ErrNotFound is returned when an item cannot be found.
	ErrNotFound = errors.New("item not found")
	// ErrInvalidItem is returned when an item cannot be persisted.
	ErrInvalidItem = errors.New("invalid item")
)

// RepositoryInterface defines item repository behavior.
type RepositoryInterface interface {
	Create(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Get(ctx context.Context, id string) (*structs.Item, error)
	Update(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error)
	Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error)
}

// Repository persists items with Ent.
type Repository struct {
	d *data.Data
}

// New creates an Ent-backed repository.
func New(d *data.Data) RepositoryInterface {
	return &Repository{d: d}
}

func (r *Repository) writeClient(ctx context.Context) (*ent.Client, error) {
	if r == nil || r.d == nil {
		return nil, errors.New("data layer is not configured")
	}
	client := r.d.GetEntClientWithFallback(ctx)
	if client == nil {
		return nil, errors.New("Ent write client is not configured")
	}
	return client, nil
}

func (r *Repository) readClient(ctx context.Context) (*ent.Client, error) {
	if r == nil || r.d == nil {
		return nil, errors.New("data layer is not configured")
	}
	client := r.d.GetEntClientWithFallback(ctx, true)
	if client == nil {
		return nil, errors.New("Ent read client is not configured")
	}
	return client, nil
}

// Create creates a new item.
func (r *Repository) Create(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil {
		return nil, ErrInvalidItem
	}
	if strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}

	client, err := r.writeClient(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.Status == "" {
		item.Status = "active"
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now
	extras := item.Extras
	if extras == nil {
		extras = map[string]any{}
	}

	record, err := client.Item.Create().
		SetID(item.ID).
		SetName(strings.TrimSpace(item.Name)).
		SetCode(strings.TrimSpace(item.Code)).
		SetType(strings.TrimSpace(item.Type)).
		SetStatus(strings.TrimSpace(item.Status)).
		SetDescription(item.Desc).
		SetCreatedBy(item.CreatedBy).
		SetCreatedAt(item.CreatedAt).
		SetUpdatedBy(item.UpdatedBy).
		SetUpdatedAt(item.UpdatedAt).
		SetExtras(extras).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return itemFromEnt(record), nil
}

// Get retrieves an item by ID.
func (r *Repository) Get(ctx context.Context, id string) (*structs.Item, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidItem
	}

	client, err := r.readClient(ctx)
	if err != nil {
		return nil, err
	}

	record, err := client.Item.Query().
		Where(itement.IDEQ(id), itement.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return itemFromEnt(record), nil
}

// Update updates an item.
func (r *Repository) Update(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}

	client, err := r.writeClient(ctx)
	if err != nil {
		return nil, err
	}

	extras := item.Extras
	if extras == nil {
		extras = map[string]any{}
	}
	record, err := client.Item.UpdateOneID(item.ID).
		SetName(strings.TrimSpace(item.Name)).
		SetCode(strings.TrimSpace(item.Code)).
		SetType(strings.TrimSpace(item.Type)).
		SetStatus(strings.TrimSpace(item.Status)).
		SetDescription(item.Desc).
		SetUpdatedBy(item.UpdatedBy).
		SetUpdatedAt(time.Now().UTC()).
		SetExtras(extras).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return itemFromEnt(record), nil
}

// Delete soft-deletes an item.
func (r *Repository) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidItem
	}

	client, err := r.writeClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.Item.UpdateOneID(id).
		SetDeletedAt(time.Now().UTC()).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// List lists items with filtering and pagination.
func (r *Repository) List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error) {
	client, err := r.readClient(ctx)
	if err != nil {
		return nil, 0, err
	}

	query := applyEntFilters(client.Item.Query().Where(itement.DeletedAtIsNil()), params)
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	limit, offset := pagination(params)
	records, err := query.Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	items := make([]*structs.Item, 0, len(records))
	for _, record := range records {
		items = append(items, itemFromEnt(record))
	}
	return items, int64(total), nil
}

// Count counts items with filtering.
func (r *Repository) Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error) {
	client, err := r.readClient(ctx)
	if err != nil {
		return 0, err
	}
	count, err := applyEntFilters(client.Item.Query().Where(itement.DeletedAtIsNil()), params).Count(ctx)
	return int64(count), err
}

func applyEntFilters(query *ent.ItemQuery, params *structs.ListItemsRequest) *ent.ItemQuery {
	if params == nil {
		return query
	}
	if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
		query = query.Where(itement.Or(
			itement.NameContainsFold(keyword),
			itement.CodeContainsFold(keyword),
			itement.DescriptionContainsFold(keyword),
		))
	}
	if itemType := strings.TrimSpace(params.Type); itemType != "" {
		query = query.Where(itement.TypeEQ(itemType))
	}
	if status := strings.TrimSpace(params.Status); status != "" {
		query = query.Where(itement.StatusEQ(status))
	}
	return query
}

func itemFromEnt(record *ent.Item) *structs.Item {
	if record == nil {
		return nil
	}
	return &structs.Item{
		ID:        record.ID,
		Name:      record.Name,
		Code:      record.Code,
		Type:      record.Type,
		Status:    record.Status,
		Desc:      record.Description,
		CreatedBy: record.CreatedBy,
		CreatedAt: record.CreatedAt,
		UpdatedBy: record.UpdatedBy,
		UpdatedAt: record.UpdatedAt,
		DeletedAt: record.DeletedAt,
		Extras:    record.Extras,
	}
}

func pagination(params *structs.ListItemsRequest) (limit int, offset int) {
	limit = 20
	page := 1
	if params != nil {
		if params.PageSize > 0 {
			limit = params.PageSize
		}
		if params.PageNum > 0 {
			page = params.PageNum
		}
	}
	if limit > 200 {
		limit = 200
	}
	return limit, (page - 1) * limit
}
`
}

func repositoryGormTemplate() string {
	return `package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/data/model"
	"{{ .PackagePath }}/structs"
)

var (
	// ErrNotFound is returned when an item cannot be found.
	ErrNotFound = errors.New("item not found")
	// ErrInvalidItem is returned when an item cannot be persisted.
	ErrInvalidItem = errors.New("invalid item")
)

// RepositoryInterface defines item repository behavior.
type RepositoryInterface interface {
	Create(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Get(ctx context.Context, id string) (*structs.Item, error)
	Update(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error)
	Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error)
}

// Repository persists items with GORM.
type Repository struct {
	d *data.Data
}

// New creates a GORM-backed repository.
func New(d *data.Data) RepositoryInterface {
	return &Repository{d: d}
}

func (r *Repository) writeDB(ctx context.Context) (*gorm.DB, error) {
	if r == nil || r.d == nil {
		return nil, errors.New("data layer is not configured")
	}
	db := r.d.GetGormClientWithFallback(ctx)
	if db == nil {
		return nil, errors.New("GORM write client is not configured")
	}
	return db.WithContext(ctx), nil
}

func (r *Repository) readDB(ctx context.Context) (*gorm.DB, error) {
	if r == nil || r.d == nil {
		return nil, errors.New("data layer is not configured")
	}
	db := r.d.GetGormClientWithFallback(ctx, true)
	if db == nil {
		return nil, errors.New("GORM read client is not configured")
	}
	return db.WithContext(ctx), nil
}

// Create creates a new item.
func (r *Repository) Create(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}
	db, err := r.writeDB(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.Status == "" {
		item.Status = "active"
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	record := modelFromItem(item)
	if err := db.Create(record).Error; err != nil {
		return nil, err
	}
	return itemFromModel(record), nil
}

// Get retrieves an item by ID.
func (r *Repository) Get(ctx context.Context, id string) (*structs.Item, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidItem
	}
	db, err := r.readDB(ctx)
	if err != nil {
		return nil, err
	}

	var record model.Item
	err = db.Where("id = ? AND deleted_at IS NULL", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return itemFromModel(&record), nil
}

// Update updates an item.
func (r *Repository) Update(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}
	db, err := r.writeDB(ctx)
	if err != nil {
		return nil, err
	}

	var record model.Item
	if err := db.Where("id = ? AND deleted_at IS NULL", item.ID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	record.Name = strings.TrimSpace(item.Name)
	record.Code = strings.TrimSpace(item.Code)
	record.Type = strings.TrimSpace(item.Type)
	record.Status = strings.TrimSpace(item.Status)
	record.Description = item.Desc
	record.UpdatedBy = item.UpdatedBy
	record.UpdatedAt = time.Now().UTC()
	record.Extras = item.Extras

	if err := db.Save(&record).Error; err != nil {
		return nil, err
	}
	return itemFromModel(&record), nil
}

// Delete soft-deletes an item.
func (r *Repository) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidItem
	}
	db, err := r.writeDB(ctx)
	if err != nil {
		return err
	}

	result := db.Model(&model.Item{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": time.Now().UTC(), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List lists items with filtering and pagination.
func (r *Repository) List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error) {
	db, err := r.readDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	query := applyGormFilters(db.Model(&model.Item{}).Where("deleted_at IS NULL"), params)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit, offset := pagination(params)
	var records []*model.Item
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*structs.Item, 0, len(records))
	for _, record := range records {
		items = append(items, itemFromModel(record))
	}
	return items, total, nil
}

// Count counts items with filtering.
func (r *Repository) Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error) {
	db, err := r.readDB(ctx)
	if err != nil {
		return 0, err
	}
	var total int64
	err = applyGormFilters(db.Model(&model.Item{}).Where("deleted_at IS NULL"), params).Count(&total).Error
	return total, err
}

func applyGormFilters(query *gorm.DB, params *structs.ListItemsRequest) *gorm.DB {
	if params == nil {
		return query
	}
	if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
		like := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?", like, like, like)
	}
	if itemType := strings.TrimSpace(params.Type); itemType != "" {
		query = query.Where("type = ?", itemType)
	}
	if status := strings.TrimSpace(params.Status); status != "" {
		query = query.Where("status = ?", status)
	}
	return query
}

func modelFromItem(item *structs.Item) *model.Item {
	return &model.Item{
		ID:          item.ID,
		Name:        strings.TrimSpace(item.Name),
		Code:        strings.TrimSpace(item.Code),
		Type:        strings.TrimSpace(item.Type),
		Status:      strings.TrimSpace(item.Status),
		Description: item.Desc,
		CreatedBy:   item.CreatedBy,
		CreatedAt:   item.CreatedAt,
		UpdatedBy:   item.UpdatedBy,
		UpdatedAt:   item.UpdatedAt,
		DeletedAt:   item.DeletedAt,
		Extras:      item.Extras,
	}
}

func itemFromModel(record *model.Item) *structs.Item {
	if record == nil {
		return nil
	}
	return &structs.Item{
		ID:        record.ID,
		Name:      record.Name,
		Code:      record.Code,
		Type:      record.Type,
		Status:    record.Status,
		Desc:      record.Description,
		CreatedBy: record.CreatedBy,
		CreatedAt: record.CreatedAt,
		UpdatedBy: record.UpdatedBy,
		UpdatedAt: record.UpdatedAt,
		DeletedAt: record.DeletedAt,
		Extras:    record.Extras,
	}
}

func pagination(params *structs.ListItemsRequest) (limit int, offset int) {
	limit = 20
	page := 1
	if params != nil {
		if params.PageSize > 0 {
			limit = params.PageSize
		}
		if params.PageNum > 0 {
			page = params.PageNum
		}
	}
	if limit > 200 {
		limit = 200
	}
	return limit, (page - 1) * limit
}
`
}

func repositoryMongoTemplate() string {
	return `package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/structs"
)

const (
	itemDatabase   = "{{ .Name }}"
	itemCollection = "items"
)

var (
	// ErrNotFound is returned when an item cannot be found.
	ErrNotFound = errors.New("item not found")
	// ErrInvalidItem is returned when an item cannot be persisted.
	ErrInvalidItem = errors.New("invalid item")
)

type mongoItem struct {
	ID        string         ` + "`" + `bson:"_id"` + "`" + `
	Name      string         ` + "`" + `bson:"name"` + "`" + `
	Code      string         ` + "`" + `bson:"code"` + "`" + `
	Type      string         ` + "`" + `bson:"type"` + "`" + `
	Status    string         ` + "`" + `bson:"status"` + "`" + `
	Desc      string         ` + "`" + `bson:"description"` + "`" + `
	CreatedBy string         ` + "`" + `bson:"created_by"` + "`" + `
	CreatedAt time.Time      ` + "`" + `bson:"created_at"` + "`" + `
	UpdatedBy string         ` + "`" + `bson:"updated_by"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `bson:"updated_at"` + "`" + `
	DeletedAt *time.Time     ` + "`" + `bson:"deleted_at,omitempty"` + "`" + `
	Extras    map[string]any ` + "`" + `bson:"extras,omitempty"` + "`" + `
}

// RepositoryInterface defines item repository behavior.
type RepositoryInterface interface {
	Create(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Get(ctx context.Context, id string) (*structs.Item, error)
	Update(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error)
	Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error)
}

// Repository persists items with MongoDB.
type Repository struct {
	d *data.Data
}

// New creates a MongoDB-backed repository.
func New(d *data.Data) RepositoryInterface {
	return &Repository{d: d}
}

func (r *Repository) collection(readOnly bool) (*mongo.Collection, error) {
	if r == nil || r.d == nil {
		return nil, errors.New("data layer is not configured")
	}
	return r.d.Collection(itemDatabase, itemCollection, readOnly)
}

// Create creates a new item.
func (r *Repository) Create(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}
	coll, err := r.collection(false)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.Status == "" {
		item.Status = "active"
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	doc := mongoFromItem(item)
	if _, err := coll.InsertOne(ctx, doc); err != nil {
		return nil, err
	}
	return itemFromMongo(doc), nil
}

// Get retrieves an item by ID.
func (r *Repository) Get(ctx context.Context, id string) (*structs.Item, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidItem
	}
	coll, err := r.collection(true)
	if err != nil {
		return nil, err
	}

	var doc mongoItem
	if err := coll.FindOne(ctx, activeMongoFilter(bson.M{"_id": id})).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return itemFromMongo(&doc), nil
}

// Update updates an item.
func (r *Repository) Update(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}
	coll, err := r.collection(false)
	if err != nil {
		return nil, err
	}

	updates := bson.M{
		"name":        strings.TrimSpace(item.Name),
		"code":        strings.TrimSpace(item.Code),
		"type":        strings.TrimSpace(item.Type),
		"status":      strings.TrimSpace(item.Status),
		"description": item.Desc,
		"updated_by":  item.UpdatedBy,
		"updated_at":  time.Now().UTC(),
		"extras":      item.Extras,
	}
	result, err := coll.UpdateOne(ctx, activeMongoFilter(bson.M{"_id": item.ID}), bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, ErrNotFound
	}
	return r.Get(ctx, item.ID)
}

// Delete soft-deletes an item.
func (r *Repository) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidItem
	}
	coll, err := r.collection(false)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	result, err := coll.UpdateOne(ctx, activeMongoFilter(bson.M{"_id": id}), bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// List lists items with filtering and pagination.
func (r *Repository) List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error) {
	coll, err := r.collection(true)
	if err != nil {
		return nil, 0, err
	}

	filter := mongoFilter(params)
	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	limit, offset := pagination(params)
	cursor, err := coll.Find(ctx, filter, options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{bson.E{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []mongoItem
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	items := make([]*structs.Item, 0, len(docs))
	for i := range docs {
		items = append(items, itemFromMongo(&docs[i]))
	}
	return items, total, nil
}

// Count counts items with filtering.
func (r *Repository) Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error) {
	coll, err := r.collection(true)
	if err != nil {
		return 0, err
	}
	return coll.CountDocuments(ctx, mongoFilter(params))
}

func mongoFilter(params *structs.ListItemsRequest) bson.M {
	filter := activeMongoFilter(bson.M{})
	clauses := make([]bson.M, 0, 4)
	if baseOr, ok := filter["$or"].([]bson.M); ok {
		clauses = append(clauses, baseOr...)
	}
	delete(filter, "$or")

	andClauses := make([]bson.M, 0, 4)
	andClauses = append(andClauses, bson.M{"$or": clauses})
	if params != nil {
		if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
			andClauses = append(andClauses, bson.M{"$or": []bson.M{
				{"name": bson.M{"$regex": keyword, "$options": "i"}},
				{"code": bson.M{"$regex": keyword, "$options": "i"}},
				{"description": bson.M{"$regex": keyword, "$options": "i"}},
			}})
		}
		if itemType := strings.TrimSpace(params.Type); itemType != "" {
			andClauses = append(andClauses, bson.M{"type": itemType})
		}
		if status := strings.TrimSpace(params.Status); status != "" {
			andClauses = append(andClauses, bson.M{"status": status})
		}
	}
	filter["$and"] = andClauses
	return filter
}

func activeMongoFilter(filter bson.M) bson.M {
	filter["$or"] = []bson.M{
		{"deleted_at": bson.M{"$exists": false}},
		{"deleted_at": nil},
	}
	return filter
}

func mongoFromItem(item *structs.Item) *mongoItem {
	return &mongoItem{
		ID:        item.ID,
		Name:      strings.TrimSpace(item.Name),
		Code:      strings.TrimSpace(item.Code),
		Type:      strings.TrimSpace(item.Type),
		Status:    strings.TrimSpace(item.Status),
		Desc:      item.Desc,
		CreatedBy: item.CreatedBy,
		CreatedAt: item.CreatedAt,
		UpdatedBy: item.UpdatedBy,
		UpdatedAt: item.UpdatedAt,
		DeletedAt: item.DeletedAt,
		Extras:    item.Extras,
	}
}

func itemFromMongo(doc *mongoItem) *structs.Item {
	if doc == nil {
		return nil
	}
	return &structs.Item{
		ID:        doc.ID,
		Name:      doc.Name,
		Code:      doc.Code,
		Type:      doc.Type,
		Status:    doc.Status,
		Desc:      doc.Desc,
		CreatedBy: doc.CreatedBy,
		CreatedAt: doc.CreatedAt,
		UpdatedBy: doc.UpdatedBy,
		UpdatedAt: doc.UpdatedAt,
		DeletedAt: doc.DeletedAt,
		Extras:    doc.Extras,
	}
}

func pagination(params *structs.ListItemsRequest) (limit int, offset int) {
	limit = 20
	page := 1
	if params != nil {
		if params.PageSize > 0 {
			limit = params.PageSize
		}
		if params.PageNum > 0 {
			page = params.PageNum
		}
	}
	if limit > 200 {
		limit = 200
	}
	return limit, (page - 1) * limit
}
`
}

func repositoryMemoryTemplate() string {
	return `package repository

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/structs"
)

var (
	// ErrNotFound is returned when an item cannot be found.
	ErrNotFound = errors.New("item not found")
	// ErrInvalidItem is returned when an item cannot be persisted.
	ErrInvalidItem = errors.New("invalid item")
)

// RepositoryInterface defines item repository behavior.
type RepositoryInterface interface {
	Create(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Get(ctx context.Context, id string) (*structs.Item, error)
	Update(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error)
	Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error)
}

// Repository persists items in process memory.
type Repository struct {
	d     *data.Data
	mu    sync.RWMutex
	items map[string]*structs.Item
}

// New creates an in-memory repository.
func New(d *data.Data) RepositoryInterface {
	return &Repository{
		d:     d,
		items: make(map[string]*structs.Item),
	}
}

// Create creates a new item.
func (r *Repository) Create(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}
	now := time.Now().UTC()
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.Status == "" {
		item.Status = "active"
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	r.mu.Lock()
	defer r.mu.Unlock()
	clone := cloneItem(item)
	r.items[clone.ID] = clone
	return cloneItem(clone), nil
}

// Get retrieves an item by ID.
func (r *Repository) Get(ctx context.Context, id string) (*structs.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok || item.DeletedAt != nil {
		return nil, ErrNotFound
	}
	return cloneItem(item), nil
}

// Update updates an item.
func (r *Repository) Update(ctx context.Context, item *structs.Item) (*structs.Item, error) {
	if item == nil || strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidItem
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.items[item.ID]
	if !ok || current.DeletedAt != nil {
		return nil, ErrNotFound
	}
	next := cloneItem(item)
	next.CreatedAt = current.CreatedAt
	next.CreatedBy = current.CreatedBy
	next.UpdatedAt = time.Now().UTC()
	r.items[next.ID] = next
	return cloneItem(next), nil
}

// Delete soft-deletes an item.
func (r *Repository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok || item.DeletedAt != nil {
		return ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	return nil
}

// List lists items with filtering and pagination.
func (r *Repository) List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*structs.Item, 0, len(r.items))
	for _, item := range r.items {
		if item.DeletedAt == nil && matches(item, params) {
			filtered = append(filtered, cloneItem(item))
		}
	}

	limit, offset := pagination(params)
	total := int64(len(filtered))
	if offset >= len(filtered) {
		return []*structs.Item{}, total, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

// Count counts items with filtering.
func (r *Repository) Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var total int64
	for _, item := range r.items {
		if item.DeletedAt == nil && matches(item, params) {
			total++
		}
	}
	return total, nil
}

func matches(item *structs.Item, params *structs.ListItemsRequest) bool {
	if params == nil {
		return true
	}
	if keyword := strings.ToLower(strings.TrimSpace(params.Keyword)); keyword != "" {
		haystack := strings.ToLower(item.Name + " " + item.Code + " " + item.Desc)
		if !strings.Contains(haystack, keyword) {
			return false
		}
	}
	if itemType := strings.TrimSpace(params.Type); itemType != "" && item.Type != itemType {
		return false
	}
	if status := strings.TrimSpace(params.Status); status != "" && item.Status != status {
		return false
	}
	return true
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

func pagination(params *structs.ListItemsRequest) (limit int, offset int) {
	limit = 20
	page := 1
	if params != nil {
		if params.PageSize > 0 {
			limit = params.PageSize
		}
		if params.PageNum > 0 {
			page = params.PageNum
		}
	}
	if limit > 200 {
		limit = 200
	}
	return limit, (page - 1) * limit
}
`
}
