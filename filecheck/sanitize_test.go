package filecheck

import (
	//"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	dotDotPath       = "_testing/../_testing/Dockerfile"
	symlinkPath      = "_testing/dockerfile-symlink"
	bogusPath        = "foobarbaz"
	validPath        = "_testing/Dockerfile"
	absValidPath, _  = filepath.Abs("../" + validPath)
	cleanedValidPath = filepath.Clean(absValidPath)
)

func TestDotDotPath(t *testing.T) {
	opts := NewTrustedFilePathOptions{File: dotDotPath, Top: ".."}

	file, err := NewTrustedFilePath(opts)
	if err != nil {
		t.Fatal(err)
	}

	file.Sanitize()
	if file.State == OK {
		t.Errorf("state should not be OK")
	}

	if file.State != NotOK {
		t.Errorf("state should be NotOK")
	}
}

func TestSymlinkPath(t *testing.T) {
	opts := NewTrustedFilePathOptions{File: symlinkPath, Top: ".."}

	file, err := NewTrustedFilePath(opts)
	if err != nil {
		t.Fatal(err)
	}

	file.Sanitize()
	if file.State != NotOK {
		t.Errorf("state should be NotOK")
	}
}

func TestBogusPath(t *testing.T) {
	opts := NewTrustedFilePathOptions{File: bogusPath, Top: ".."}

	file, err := NewTrustedFilePath(opts)
	if err != nil {
		t.Fatal(err)
	}

	file.Sanitize()
	if file.State != NotOK {
		t.Errorf("state should be NotOK")
	}
}

func TestValidPath(t *testing.T) {
	var path = os.Getenv("GOPATH") + "/src/github.com/winchman/builder-core/_testing"
	var filename = "Dockerfile"
	var fullpath = path + "/" + filename

	opts := NewTrustedFilePathOptions{File: validPath, Top: ".."}
	file, err := NewTrustedFilePath(opts)
	if err != nil {
		t.Fatal(err)
	}

	file.Sanitize()
	if file.State != OK {
		t.Errorf("state should be OK")
	}

	if file.File() != filename {
		t.Errorf("expected file %q, got file %q", filename, file.File())
	}

	if file.Top() != path {
		t.Errorf("expected path to file %q, got path to file %q", path, file.Top())
	}

	if file.FullPath() != fullpath {
		t.Errorf("expected full path %q, got full path %q", fullpath, file.FullPath())
	}
}
