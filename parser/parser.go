package parser

import (
	"github.com/Sirupsen/logrus"

	"github.com/sylphon/build-runner/unit-config"
)

/*
Parser is a struct that contains a Builderfile and knows how to parse it both
as raw text and to convert toml to a Builderfile struct.  It also knows how to
tell if the Builderfile is valid (openable) or nat.
*/
type Parser struct {
	*logrus.Logger
	top string
}

// NewParserOptions encapsulates all of the options necessary when creating a new parser
type NewParserOptions struct {
	ContextDir string
	Logger     *logrus.Logger
}

/*
NewParser returns an initialized Parser.  Not currently necessary, as no
default values are assigned to a new Parser, but useful to have in case we need
to change this.
*/
func NewParser(opts NewParserOptions) *Parser {
	unitconfig.Logger(opts.Logger)
	return &Parser{
		Logger: opts.Logger,
		top:    opts.ContextDir,
	}
}
