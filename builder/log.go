package builder

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wsxiaoys/terminal/color"
)

/*
An OutWriter is responsible for for implementing the io.Writer interface.
*/
type outWriter struct {
	*logrus.Logger
	fmtString string
}

/*
NewOutWriter accepts a logger and a format string and returns an OutWriter.
When written to, the OutWriter will take the input, split it into lines, and
print it to the logger using the provided format string.  The intended use case
of this functionality is for printing nice, colorful messages
*/
func newOutWriter(logger *logrus.Logger, fmtString string) *outWriter {
	return &outWriter{
		Logger:    logger,
		fmtString: fmtString,
	}
}

/*
Write writes the provided bytes, one line at a time, after interpolating them
into the provided format string, to the provided logger.
*/
func (ow *outWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		ow.Debug(color.Sprintf(ow.fmtString, line))
	}

	return len(p), nil
}
