package builder

import "fmt"

type Builder struct {
	cfg Config
}

func New(config *Config) *Builder {
	return &Builder{
		cfg: *config,
	}
}

func (g *Builder) Build(unit string) error {
	return nil
}

type BuildKind uint32

const (
	// debug-friendly optimizations + debug information in binaries + safety checks
	BuildDebug = iota + 1

	// moderate-level optimizations + debug information in binaries + safety checks
	BuildTest

	// most optimizations enabled + safety checks
	BuildSafe

	// all optimizations enabled + disabled safety checks
	BuildFast
)

var buildKindText = [...]string{
	0: "<nil>",

	BuildDebug: "debug",
	BuildTest:  "test",
	BuildSafe:  "safe",
	BuildFast:  "fast",
}

func (k BuildKind) String() string {
	return buildKindText[k]
}

func ParseKind(s string) (BuildKind, error) {
	switch s {
	case "":
		return 0, fmt.Errorf("empty build kind")
	case "debug":
		return BuildDebug, nil
	case "test":
		return BuildTest, nil
	case "safe":
		return BuildSafe, nil
	case "fast":
		return BuildFast, nil
	default:
		return 0, fmt.Errorf("unknown build kind: %s", s)
	}
}

type Config struct {
	BaseOutputDir string
	BaseCacheDir  string

	BuildKind BuildKind
}
