package version

// Dependency versions
const (
	DefaultGoVersion  = "1.25"
	GinVersion        = "v1.12.0"
	UUIDVersion       = "v1.6.0"
	EntVersion        = "v0.14.6"
	GormVersion       = "v1.31.2"
	MongoVersion      = "v2.7.0"
	PrometheusVersion = "v1.23.2"
	LogrusVersion     = "v1.9.4"
)

// NcoreVersion returns the released version for a ncore module.
func NcoreVersion(module string) string {
	if version, ok := ncoreVersions[module]; ok {
		return version
	}
	return "v0.2.3"
}

// GetNcoreDeps returns the base ncore modules required by generated services.
func GetNcoreDeps() []string {
	return []string{
		"github.com/ncobase/ncore/config",
		"github.com/ncobase/ncore/consts",
		"github.com/ncobase/ncore/ctxutil",
		"github.com/ncobase/ncore/logging",
		"github.com/ncobase/ncore/ecode",
		"github.com/ncobase/ncore/net",
		"github.com/ncobase/ncore/extension",
		"github.com/ncobase/ncore/data",
	}
}

var ncoreVersions = map[string]string{
	"github.com/ncobase/ncore/config":             "v0.2.3",
	"github.com/ncobase/ncore/consts":             "v0.2.3",
	"github.com/ncobase/ncore/ctxutil":            "v0.2.3",
	"github.com/ncobase/ncore/data":               "v0.2.4",
	"github.com/ncobase/ncore/data/cache":         "v0.2.3",
	"github.com/ncobase/ncore/data/elasticsearch": "v0.2.3",
	"github.com/ncobase/ncore/data/kafka":         "v0.2.3",
	"github.com/ncobase/ncore/data/meilisearch":   "v0.2.3",
	"github.com/ncobase/ncore/data/mongodb":       "v0.2.3",
	"github.com/ncobase/ncore/data/mysql":         "v0.2.3",
	"github.com/ncobase/ncore/data/opensearch":    "v0.2.3",
	"github.com/ncobase/ncore/data/postgres":      "v0.2.3",
	"github.com/ncobase/ncore/data/rabbitmq":      "v0.2.3",
	"github.com/ncobase/ncore/data/redis":         "v0.2.3",
	"github.com/ncobase/ncore/data/sqlite":        "v0.2.3",
	"github.com/ncobase/ncore/ecode":              "v0.2.3",
	"github.com/ncobase/ncore/extension":          "v0.2.3",
	"github.com/ncobase/ncore/logging":            "v0.2.3",
	"github.com/ncobase/ncore/net":                "v0.2.3",
	"github.com/ncobase/ncore/oss":                "v0.2.4",
	"github.com/ncobase/ncore/security":           "v0.2.3",
	"github.com/ncobase/ncore/utils":              "v0.2.3",
	"github.com/ncobase/ncore/validation":         "v0.2.3",
}
