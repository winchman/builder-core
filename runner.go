package runner

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/sylphon/builder-core/builder"
	"github.com/sylphon/builder-core/parser"
	"github.com/sylphon/builder-core/unit-config"
)

// Stream corresponds to a file stream (stdout/stderr)
type Stream int

const (
	stdin Stream = iota
	// Stdout indicates the LogMsg's message should be printed, if at all, to stdout
	Stdout

	// Stderr indicates the LogMsg's message should be printed, if at all, to stderr
	Stderr
)

// LogMsg is the tentative interface for the structs returned on the log messages channel
type LogMsg interface {
	BuildID() string
	Level() int // type may change
	Msg() string
	Stream() Stream
}

// StatusMsg is the tentative interface for the structs returned on the status channel
type StatusMsg interface {
	BuildID() int
	Status() int // type may change
	Msg() string
	Error() error // should be checked for non-nil
}

// RunBuild runs a complete build for the provided unit config.  Currently, the
// channels argument is ignored but will be used in the future along with the
// LogMsg and StatusMsg interfaces
func RunBuild(unitConfig *unitconfig.UnitConfig, contextDir string, channels ...chan interface{}) error {
	var err error
	var logger = logrus.New()

	logger.Level = logrus.DebugLevel

	if unitConfig == nil {
		return errors.New("unit config may not be nil")
	}

	var p *parser.Parser
	opts := parser.NewParserOptions{ContextDir: contextDir, Logger: logger}
	p = parser.NewParser(opts)

	commandSequence := p.Parse(unitConfig)

	var bob *builder.Builder
	bobOpts := builder.NewBuilderOptions{ContextDir: contextDir, Logger: logger}
	bob, err = builder.NewBuilder(bobOpts)
	if err != nil {
		return err
	}

	if buildErr := bob.BuildCommandSequence(commandSequence); buildErr != nil {
		return buildErr
	}

	return nil
}
