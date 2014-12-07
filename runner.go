package runner

import (
	"errors"

	"github.com/Sirupsen/logrus"
	b "github.com/winchman/builder-core/builder"
	"github.com/winchman/builder-core/communication"
	p "github.com/winchman/builder-core/parser"
	"github.com/winchman/builder-core/unit-config"
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
	logger.Level = logrus.InfoLevel
	log, status, done := RunBuild(unitConfig, contextDir)
	for {
		select {
		case e, ok := <-log:
			if !ok {
				return errors.New("log channel closed prematurely")
			}
			e.Entry().Logger = logger
			switch e.Entry().Level {
			case logrus.PanicLevel:
				e.Entry().Panicln(e.Entry().Message)
			case logrus.FatalLevel:
				e.Entry().Fatalln(e.Entry().Message)
			case logrus.ErrorLevel:
				e.Entry().Errorln(e.Entry().Message)
			case logrus.WarnLevel:
				e.Entry().Warnln(e.Entry().Message)
			case logrus.InfoLevel:
				e.Entry().Infoln(e.Entry().Message)
			default:
				e.Entry().Debugln(e.Entry().Message)
			}
		case event, ok := <-status:
			if !ok {
				return errors.New("status channel closed prematurely")
			}
			logger.WithFields(event.Data()).Debugf("status event (type %s)", event.EventType())
		case err, ok := <-done:
			if !ok {
				return errors.New("exit channel closed prematurely")
			}
			return err
		}
	}
}
