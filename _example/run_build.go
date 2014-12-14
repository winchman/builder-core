package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/winchman/builder-core"
	"github.com/winchman/builder-core/unit-config"
)

var example = &unitconfig.UnitConfig{
	Version: 1,
	Docker: unitconfig.Docker{
		TagOpts: []string{"--force"},
	},
	ContainerArr: []*unitconfig.ContainerSection{
		&unitconfig.ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "builder-core-test",
			Tags:       []string{"latest", "{{ sha }}", "{{ tag }}", "{{ branch }}"},
			SkipPush:   true,
		},
	},
}

func main() {
	opts := runner.Options{
		UnitConfig: example,
		ContextDir: os.Getenv("GOPATH") + "/src/github.com/rafecolton/docker-builder",
		LogLevel:   logrus.InfoLevel,
	}
	if err := runner.RunBuildSynchronously(opts, runner.KeepTemporaryTag); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
