package builder

import (
	"os"
	"path/filepath"
	"regexp"
)

const (
	// DotDotSanitizeErrorMessage is the error message used in errors that occur
	// because a provided Bobfile path contains ".."
	dotDotSanitizeErrorMessage = "file path must not contain .."

	// InvalidPathSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path is invalid
	invalidPathSanitizeErrorMessage = "file path is invalid"

	// SymlinkSanitizeErrorMessage is the error message used in errors that
	// occur because a provided Bobfile path contains symlinks
	symlinkSanitizeErrorMessage = "file path must not contain symlinks"

	// DoesNotExistSanitizeErrorMessage is the error message used in cases
	// where the error results in the requested file not existing
	doesNotExistSanitizeErrorMessage = "file requested does not exist"
)

var dotDotRegex = regexp.MustCompile(`\.\.`)

// SanitizeTrustedFilePath checks for disallowed entries in the provided
// file path and returns either a sanitized version of the path or an error
func SanitizeTrustedFilePath(trustedFilePath *TrustedFilePath) (*TrustedFilePath, Error) {
	var file = trustedFilePath.File()
	var top = trustedFilePath.Top()

	if dotDotRegex.MatchString(file) {
		return nil, &sanitizeError{
			Message:  dotDotSanitizeErrorMessage,
			Filename: file,
		}
	}

	abs, err := filepath.Abs(top + "/" + file)
	if err != nil {
		return nil, &sanitizeError{
			Message:  invalidPathSanitizeErrorMessage,
			error:    err,
			Filename: file,
		}
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		msg := invalidPathSanitizeErrorMessage
		if os.IsNotExist(err) {
			msg = doesNotExistSanitizeErrorMessage
		}
		return nil, &sanitizeError{
			Message:  msg,
			error:    err,
			Filename: file,
		}
	}

	if abs != resolved {
		return nil, &sanitizeError{
			Message:  symlinkSanitizeErrorMessage,
			Filename: file,
		}
	}

	clean := filepath.Clean(abs)

	return &TrustedFilePath{
		top:  filepath.Dir(clean),
		file: filepath.Base(clean),
	}, nil
}
