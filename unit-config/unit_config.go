package unitconfig

/*
UnitConfig is a struct representation of what is expected to be inside a
Builderfile for a single build/tag/push sequence.
*/
type UnitConfig struct {
	Version          int                 `toml:"version"`
	Docker           Docker              `toml:"docker"`
	ContainerArr     []*ContainerSection `toml:"container"`
	ContainerGlobals *ContainerSection   `toml:"container_globals"`
}

/*
Docker is a struct representation of the "docker" section of a Builderfile.
*/
type Docker struct {
	BuildOpts []string `toml:"build_opts"`
	TagOpts   []string `toml:"tag_opts"`
}

/*
ContainerSection is a struct representation of an individual member of the  "containers"
section of a Builderfile. Each of these sections defines a docker container to
be built and other related options.
*/
type ContainerSection struct {
	Name       string   `toml:"name"`
	Dockerfile string   `toml:"Dockerfile"`
	Registry   string   `toml:"registry"`
	Project    string   `toml:"project"`
	Tags       []string `toml:"tags"`
	SkipPush   bool     `toml:"skip_push"`
	CfgUn      string   `toml:"dockercfg_un"`
	CfgPass    string   `toml:"dockercfg_pass"`
	CfgEmail   string   `toml:"dockercfg_email"`
}
