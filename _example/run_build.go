package main

import (
	"fmt"
	"os"

	"github.com/sylphon/builder-core"
	"github.com/sylphon/builder-core/unit-config"
)

var example = &unitconfig.UnitConfig{
	Version: 1,
	ContainerArr: []*unitconfig.ContainerSection{
		&unitconfig.ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "builder-core-test",
			Tags:       []string{"latest", "git:sha", "git:tag", "git:branch"},
			SkipPush:   true,
		},
	},
}

func main() {
	if err := runner.RunBuild(example, os.Getenv("GOPATH")+"/src/github.com/rafecolton/docker-builder"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
