package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/rafecolton/go-gitutils"
)

func TestTemplateBasedTags(t *testing.T) {
	var input = ` latest {{ branch }} {{ sha }} {{ tag }} `
	var tag = NewTag(input)

	var actual = tag.Evaluate(os.Getenv("PWD"))
	top := os.Getenv("PWD")
	expected := fmt.Sprintf(` latest %s %s %s `, git.Branch(top), git.Sha(top), git.Tag(top))
	if expected != actual {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
