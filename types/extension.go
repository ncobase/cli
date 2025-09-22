package types

// ExtensionType represents extension category
type ExtensionType string

const (
	CoreType     ExtensionType = "core"
	BusinessType ExtensionType = "business"
	PluginType   ExtensionType = "plugin"
)

type Extension interface {
	GetMetadata() ExtensionMetadata
	Status() string
}

type ExtensionMetadata struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Type        ExtensionType `json:"type"`
	Description string        `json:"description"`
}

// Manager handles extension lifecycle
type Manager interface {
	LoadPlugins() error
	GetExtensions() map[string]Extension
}

// Registry manages extension registration
type Registry interface {
	Register(name string, ext Extension) error
	Get(name string) (Extension, error)
}
