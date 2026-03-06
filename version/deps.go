package version

// Dependency versions
const (
	DefaultGoVersion = "1.25"
	GinVersion       = "v1.11.0"
	UUIDVersion      = "v1.6.0"
	EntVersion       = "v0.14.5"
	GormVersion      = "v1.25.12"
	MongoVersion     = "v1.17.6"
	NcoreVersion     = "v0.2.3"
)

// GetNcoreDeps returns list of ncore dependencies
func GetNcoreDeps() []string {
	return []string{
		"github.com/ncobase/ncore/config",
		"github.com/ncobase/ncore/logging",
		"github.com/ncobase/ncore/ecode",
		"github.com/ncobase/ncore/net",
		"github.com/ncobase/ncore/extension",
	}
}
