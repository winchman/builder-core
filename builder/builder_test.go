package builder

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/sylphon/build-runner/parser"
	"github.com/sylphon/build-runner/unit-config"
)

var unitConfig = &unitconfig.UnitConfig{
	Version: 1,
	ContainerArr: []*unitconfig.ContainerSection{
		&unitconfig.ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "build-runner-test",
			Tags:       []string{"latest"},
			SkipPush:   true,
		},
	},
}

func TestBuilder(t *testing.T) {
	var opts = parser.NewParserOptions{ContextDir: os.Getenv("PWD"), Logger: nil}
	var p = parser.NewParser(opts)
	commandSequence := p.Parse(unitConfig)

	var logger = logrus.New()
	logger.Level = logrus.DebugLevel

	builderOpts := NewBuilderOptions{
		Logger:       logger,
		ContextDir:   os.Getenv("GOPATH") + "/src/github.com/sylphon/build-runner/_testing",
		dockerClient: &nullClient{},
	}

	builder, err := NewBuilder(builderOpts)
	if err != nil {
		t.Fatal(err)
	}

	if err := builder.BuildCommandSequence(commandSequence); err != nil {
		t.Fatal(err)
	}
}
