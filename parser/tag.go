package parser

import (
	"github.com/sylphon/build-runner/git"
)

// TODO: add template-based tagging, do away with the rest of this

type tag struct {
	value string
}

/*
NewTag returns a Tag instance.  See function implementation for details on what
args to pass.
*/
func NewTag(value string) tag {
	return tag{value: value}
}

func (t tag) Evaluate(top string) string {
	var ret string

	switch t.value {
	case "git:branch":
		ret = git.Branch(top)
	case "git:rev", "git:sha":
		ret = git.Sha(top)
	case "git:short", "git:tag":
		ret = git.Tag(top)
	default:
		ret = t.value
	}
	return ret
}
