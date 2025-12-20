package templates

import "fmt"

func RepositoryTemplate(data *Data) string {
	imports := fmt.Sprintf(`"context"
	"fmt"
	"{{ .PackagePath }}/data"
	"{{ .PackagePath }}/structs"
`)

	// Add specific imports based on DB
	if data.UseEnt {
		imports += fmt.Sprintf(`	"{{ .PackagePath }}/data/ent"
	"{{ .PackagePath }}/data/ent/user"
`)
	}
	if data.UseMongo {
		imports += `	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
`
	}
	if data.UseGorm {
		imports += `	"gorm.io/gorm"
`
	}

	content := fmt.Sprintf(`package repository

import (
%s
)

// RepositoryInterface represents the repository interface.
type RepositoryInterface interface {
	Create(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Get(ctx context.Context, id string) (*structs.Item, error)
	Update(ctx context.Context, item *structs.Item) (*structs.Item, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error)
	Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error)
}

// Repository represents the %s repository.
type Repository struct {
	d *data.Data
}

// New creates a new repository.
func New(d *data.Data) RepositoryInterface {
	return &Repository{
		d: d,
	}
}
`, imports, data.Name)

	// Create Method
	content += `
// Create creates a new item.
func (r *Repository) Create(ctx context.Context, item *structs.Item) (*structs.Item, error) {
`
	if data.UseEnt {
		content += `	// Using Ent
	client := r.d.GetEntClientWithFallback(ctx)
	// Example: assumes User schema for demo
	// u, err := client.User.Create().SetName(item.Name).Save(ctx)
	_ = client
	fmt.Printf("Creating item with Ent: %%+v\n", item)
	// Return dummy
	item.ID = "generated_id"
	return item, nil
`
	} else if data.UseMongo {
		content += `	// Using Mongo
	coll := r.d.GetMongoCollection("database_name", "items", false)
	res, err := coll.InsertOne(ctx, bson.M{
		"name": item.Name,
		"code": item.Code,
	})
	if err != nil {
		return nil, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		item.ID = oid.Hex()
	}
	return item, nil
`
	} else if data.UseGorm {
		content += `	// Using Gorm
	db := r.d.GetGormClient()
	if err := db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
`
	} else {
		content += `	// Default Mock Implementation
	fmt.Printf("Creating item: %%+v\n", item)
	item.ID = "mock_id"
	return item, nil
`
	}
	content += `}
`
	// Get Method
	content += `
// Get retrieves an item by ID.
func (r *Repository) Get(ctx context.Context, id string) (*structs.Item, error) {
`
	if data.UseEnt {
		content += `	// Using Ent
	// client := r.d.GetEntClientWithFallback(ctx, true)
	// u, err := client.User.Query().Where(user.IDEQ(id)).Only(ctx)
	return &structs.Item{ID: id, Name: "Ent Item"}, nil
`
	} else if data.UseMongo {
		content += `	// Using Mongo
	coll := r.d.GetMongoCollection("database_name", "items", true)
	var item structs.Item
	objID, _ := primitive.ObjectIDFromHex(id)
	err := coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}
	return &item, nil
`
	} else if data.UseGorm {
		content += `	// Using Gorm
	db := r.d.GetGormClientRead()
	var item structs.Item
	if err := db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
`
	} else {
		content += `	// Default Mock Implementation
	return &structs.Item{ID: id, Name: "Mock Item"}, nil
`
	}
	content += `}
`

	// Update Method
	content += `
// Update updates an existing item.
func (r *Repository) Update(ctx context.Context, item *structs.Item) (*structs.Item, error) {
`
	if data.UseEnt {
		content += `	// Using Ent
	return item, nil
`
	} else if data.UseMongo {
		content += `	// Using Mongo
	coll := r.d.GetMongoCollection("database_name", "items", false)
	objID, _ := primitive.ObjectIDFromHex(item.ID)
	_, err := coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{
		"name": item.Name,
		"updated_at": item.UpdatedAt,
	}})
	if err != nil {
		return nil, err
	}
	return item, nil
`
	} else if data.UseGorm {
		content += `	// Using Gorm
	db := r.d.GetGormClient()
	if err := db.WithContext(ctx).Save(item).Error; err != nil {
		return nil, err
	}
	return item, nil
`
	} else {
		content += `	// Default Mock Implementation
	fmt.Printf("Updating item: %%+v\n", item)
	return item, nil
`
	}
	content += `}
`

	// Delete Method
	content += `
// Delete deletes an item by ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
`
	if data.UseEnt {
		content += `	// Using Ent
	return nil
`
	} else if data.UseMongo {
		content += `	// Using Mongo
	coll := r.d.GetMongoCollection("database_name", "items", false)
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := coll.DeleteOne(ctx, bson.M{"_id": objID})
	return err
`
	} else if data.UseGorm {
		content += `	// Using Gorm
	db := r.d.GetGormClient()
	return db.WithContext(ctx).Delete(&structs.Item{}, "id = ?", id).Error
`
	} else {
		content += `	// Default Mock Implementation
	fmt.Printf("Deleting item: %%s\n", id)
	return nil
`
	}
	content += `}
`

	// List Method
	content += `
// List lists items based on parameters.
func (r *Repository) List(ctx context.Context, params *structs.ListItemsRequest) ([]*structs.Item, int64, error) {
`
	if data.UseEnt {
		content += `	// Using Ent
	return []*structs.Item{}, 0, nil
`
	} else if data.UseMongo {
		content += `	// Using Mongo
	coll := r.d.GetMongoCollection("database_name", "items", true)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var items []*structs.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}
	return items, int64(len(items)), nil
`
	} else if data.UseGorm {
		content += `	// Using Gorm
	db := r.d.GetGormClientRead()
	var items []*structs.Item
	var count int64
	if err := db.WithContext(ctx).Find(&items).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return items, count, nil
`
	} else {
		content += `	// Default Mock Implementation
	items := []*structs.Item{
		{ID: "1", Name: "Item 1"},
		{ID: "2", Name: "Item 2"},
	}
	return items, 2, nil
`
	}
	content += `}
`

	// Count Method
	content += `
// Count counts items based on parameters.
func (r *Repository) Count(ctx context.Context, params *structs.ListItemsRequest) (int64, error) {
	// Implement counting logic
	return 0, nil
}
`

	return content
}
