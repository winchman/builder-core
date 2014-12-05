package builder

import (
	"os"
	"testing"

	"github.com/rafecolton/go-dockerclient-quick"
	"github.com/winchman/builder-core/parser"
	"github.com/winchman/builder-core/unit-config"
)

var unitConfig = &unitconfig.UnitConfig{
	Version: 1,
	ContainerArr: []*unitconfig.ContainerSection{
		&unitconfig.ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "builder-core-test",
			Tags:       []string{"latest"},
			SkipPush:   true,
		},
	},
}

func TestBuildCommandSequence(t *testing.T) {
	var opts = parser.NewParserOptions{ContextDir: os.Getenv("PWD")}
	var p = parser.NewParser(opts)
	commandSequence := p.Parse(unitConfig)

	builder := NewBuilder(NewBuilderOptions{
		ContextDir:   os.Getenv("GOPATH") + "/src/github.com/winchman/builder-core/_testing",
		dockerClient: dockerclient.FakeClient(),
	})

	if err := builder.BuildCommandSequence(commandSequence); err != nil {
		t.Fatal(err)
	}
}
