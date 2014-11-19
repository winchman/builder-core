package builder

import (
	"os"
	"reflect"
)

// Error is an interface for any error types returned by the builder
// package / during the building process
type Error interface {
	// Error returns the error message to satisfy the error interface
	Error() string

	// ExitCode returns the code that should be used when exiting as a result
	// of this error
	ExitCode() int
}

// SanitizeError is used for errors related to sanitizing a given Bobfile path
type sanitizeError struct {
	Message  string
	Filename string
	error
}

// IsSanitizeFileNotExist returns true if the error provided is a SanitizeError
// resulting from an "os" IsNotExist PathError
func isSanitizeFileNotExist(err Error) bool {
	if isSanitizeError(err) {
		return os.IsNotExist(err.(*sanitizeError).error)
	}
	return false
}

// IsSanitizeError returns true if the error provided is of the type SanitizeError
func isSanitizeError(err Error) bool {
	return reflect.TypeOf(err).ConvertibleTo(reflect.TypeOf(&sanitizeError{}))
}

// Error returns the error message for a SanitizeError.  It is expected to be
// set at the time that the struct instance is created
func (err *sanitizeError) Error() string {
	return err.Message
}

// ExitCode returns the exit code for errors related to sanitizing the Bobfile
// path.  It is the same value for all sanitize errors.
func (err *sanitizeError) ExitCode() int {
	return 67
}

// ParserRelatedError is used for errors encounted during the building process
// that are related to parsing the Bobfile
type parserRelatedError struct {
	Message string
	Code    int
}

// Error returns the error message for a ParserRelatedError.  It is expected to
// be set at the time that the struct instance is created
func (err *parserRelatedError) Error() string {
	return err.Message
}

// ExitCode returns the exit code for errors related to parsing the Bobfile
// path.  It is expected to be set during the time that struct instance is
// created
func (err *parserRelatedError) ExitCode() int {
	return err.Code
}

// BuildRelatedError is used for build-related errors produced by the builder package
// that are encountered during the build process
type buildRelatedError struct {
	Message string
	Code    int
}

// Error returns the error message for a build-related error.  It is expected
// to be set at the time that the struct instance is created
func (err *buildRelatedError) Error() string {
	return err.Message
}

// ExitCode returns the exit code for errors related to the build process.  It
// is expected to be set during the time that struct instance is created
func (err *buildRelatedError) ExitCode() int {
	return err.Code
}
