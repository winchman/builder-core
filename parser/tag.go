package parser

import (
	"github.com/rafecolton/go-gitutils"
)

// TODO: add template-based tagging, do away with the rest of this

// Tag is for tagging
type Tag struct {
	value string
}

/*
NewTag returns a Tag instance.  See function implementation for details on what
args to pass.
*/
func NewTag(value string) Tag {
	return Tag{value: value}
}

// Evaluate evaluates any git-based tags
func (t Tag) Evaluate(top string) string {
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
