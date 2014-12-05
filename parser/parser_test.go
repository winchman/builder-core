package parser

import (
	"os"
	"reflect"
	"testing"

	"github.com/fsouza/go-dockerclient"
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

var expectedCommandSequence = &CommandSequence{
	Commands: []*SubSequence{
		&SubSequence{
			Metadata: &SubSequenceMetadata{
				Name:       "app",
				Dockerfile: "Dockerfile",
			},
			SubCommand: []DockerCmd{
				&BuildCmd{
					buildOpts: docker.BuildImageOptions{
						Name:           "quay.io/rafecolton/builder-core-test:9af73a34-ab4a-4d76-593b-fda4a5d1a988",
						RmTmpContainer: true,
						AuthConfigs: docker.AuthConfigurations{
							Configs: map[string]docker.AuthConfiguration{
								"quay.io/rafecolton": docker.AuthConfiguration{
									ServerAddress: "quay.io/rafecolton",
								},
							},
						},
						ContextDir: os.Getenv("GOPATH") + "/src/github.com/winchman/builder-core/parser"},
				},
				&TagCmd{
					Tag:  "latest",
					Repo: "quay.io/rafecolton/builder-core-test",
				},
			},
		},
	},
}

func TestParse(t *testing.T) {

	var opts = NewParserOptions{ContextDir: os.Getenv("PWD")}
	var p = NewParser(opts)
	commandSequence := p.Parse(unitConfig)

	got := commandSequence.Commands[0].Metadata
	expected := expectedCommandSequence.Commands[0].Metadata
	got.UUID = expected.UUID // monkey-patch

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %#v, expected %#v", got, expected)
	}

	buildCmd := commandSequence.Commands[0].SubCommand[0]
	expectedBuildCmd := expectedCommandSequence.Commands[0].SubCommand[0]
	buildCmd.(*BuildCmd).buildOpts.Name = expectedBuildCmd.(*BuildCmd).buildOpts.Name // monkey-patch for uuid

	if !reflect.DeepEqual(buildCmd, expectedBuildCmd) {
		t.Errorf("got %#v, expected %#v\n\n", buildCmd, expectedBuildCmd)
	}

	tagCmd := commandSequence.Commands[0].SubCommand[1]
	expectedTagCmd := expectedCommandSequence.Commands[0].SubCommand[1]
	if !reflect.DeepEqual(tagCmd, expectedTagCmd) {
		t.Errorf("got %#v, expected %#v\n\n", tagCmd, expectedTagCmd)
	}
}
