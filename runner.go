package runner

import (
	"errors"

	"github.com/Sirupsen/logrus"
	b "github.com/sylphon/builder-core/builder"
	"github.com/sylphon/builder-core/communication"
	p "github.com/sylphon/builder-core/parser"
	"github.com/sylphon/builder-core/unit-config"
)

// RunBuild runs a complete build for the provided unit config.  Currently, the
// channels argument is ignored but will be used in the future along with the
// LogMsg and StatusMsg interfaces
func RunBuild(unitConfig *unitconfig.UnitConfig, contextDir string) (comm.LogChan, comm.EventChan, comm.ExitChan) {
	var log = make(chan comm.LogEntry, 1)
	var event = make(chan comm.Event, 1)
	var exit = make(chan error)

	go func() {
		var err error

		if unitConfig == nil {
			exit <- errors.New("unit config may not be nil")
			return
		}

		parser := p.NewParser(p.NewParserOptions{
			ContextDir: contextDir,
			Log:        log,
			Event:      event,
		})
		commandSequence := parser.Parse(unitConfig)

		builder := b.NewBuilder(b.NewBuilderOptions{
			ContextDir: contextDir,
			Log:        log,
			Event:      event,
		})

		if err = builder.BuildCommandSequence(commandSequence); err != nil {
			exit <- err
			return
		}

		exit <- nil
	}()

	return log, event, exit
}

// RunBuildSynchronously - run a build, wait for it to finish, log to stdout
func RunBuildSynchronously(unitConfig *unitconfig.UnitConfig, contextDir string) error {
	var logger = logrus.New()
	logger.Level = logrus.DebugLevel
	log, _, done := RunBuild(unitConfig, contextDir)
	for {
		select {
		case e, ok := <-log:
			if !ok {
				return errors.New("log channel closed prematurely")
			}
			e.Entry().Logger = logger
			e.Entry().Debugln(e.Entry().Message)
		case err, ok := <-done:
			if !ok {
				return errors.New("exit channel closed prematurely")
			}
			return err
		}
	}
}
