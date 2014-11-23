package unitconfig

import (
	"reflect"
	"testing"
)

var expectedUnitConfig = &UnitConfig{
	Version: 1,
	ContainerGlobals: &ContainerSection{
		SkipPush: true,
	},
	ContainerArr: []*ContainerSection{
		&ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "docker-builder",
			Tags: []string{
				"latest",
				"{{ sha }}",
				"{{ tag }}",
				"{{ branch }}",
			},
		},
	},
}

func TestJSON(t *testing.T) {
	config, err := ReadFromFile("../_testing/.Bobfile.json", JSON)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, expectedUnitConfig) {
		t.Errorf("expected %#v, got %#v", expectedUnitConfig, config)
	}
}

func TestYAML(t *testing.T) {
	config, err := ReadFromFile("../_testing/.Bobfile.yml", YAML)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, expectedUnitConfig) {
		t.Errorf("expected %#v, got %#v", expectedUnitConfig, config)
	}
}

func TestTOML(t *testing.T) {
	config, err := ReadFromFile("../_testing/.Bobfile.toml", TOML)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, expectedUnitConfig) {
		t.Errorf("expected %#v, got %#v", expectedUnitConfig, config)
	}
}

func TestUnknown(t *testing.T) {
	_, err := ReadFromFile("read_test.go")
	if err == nil {
		t.Fatal("should not be able to parse the given file in the available formats")
	}
}
