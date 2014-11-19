package parser

import (
	//"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/fsouza/go-dockerclient"
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

var expectedCommandSequence = &CommandSequence{
	Commands: []*SubSequence{
		&SubSequence{
			Metadata: &SubSequenceMetadata{
				Name:       "app",
				Dockerfile: "Dockerfile",
			},
			SubCommand: []DockerCmd{
				&BuildCmd{opts: nil, buildOpts: docker.BuildImageOptions{Name: "quay.io/rafecolton/build-runner-test:9af73a34-ab4a-4d76-593b-fda4a5d1a988", NoCache: false, SuppressOutput: false, RmTmpContainer: true, ForceRmTmpContainer: false, RawJSONStream: false, Remote: "", Auth: docker.AuthConfiguration{Username: "", Password: "", Email: "", ServerAddress: ""}, AuthConfigs: docker.AuthConfigurations{Configs: map[string]docker.AuthConfiguration{"quay.io/rafecolton": docker.AuthConfiguration{Username: "", Password: "", Email: "", ServerAddress: "quay.io/rafecolton"}}}, ContextDir: "/Users/r.colton/.gvm/pkgsets/go1.3.3/global/src/github.com/sylphon/build-runner/parser"}, origBuildOpts: []string(nil)},
				&TagCmd{TagFunc: (func(string, docker.TagImageOptions) error)(nil), Image: "", Force: false, Tag: "latest", Repo: "quay.io/rafecolton/build-runner-test", msg: ""},
			},
		},
	},
}

func TestParse(t *testing.T) {

	var opts = NewParserOptions{ContextDir: os.Getenv("PWD"), Logger: nil}
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
