// Ref. https://github.com/stackup-wallet/stackup-bundler/blob/main/pkg/tracer/load.go

// Package tracer provides custom tracing capabilities to comply with EIP-4337 specifications for
// forbidden opcodes.
package tracer

import (
	"embed"
	"io/fs"
	"regexp"
	"strings"
)

//go:embed *EventTracer.js
var files embed.FS
var (
	commentRegex    = regexp.MustCompile("(?m)^.*//.*$[\r\n]+")
	whiteSpaceRegex = regexp.MustCompile(`\B\s+|\s+\B`)
	constInitStr    = "var tracer ="
	endLineChar     = ";"
)

// parse takes the raw tracer from file and removes all non-essential code.
func parse(code string) string {
	m := commentRegex.ReplaceAllString(code, "")
	m = strings.TrimSpace(m)
	m = strings.TrimPrefix(m, constInitStr)
	m = strings.TrimSuffix(m, endLineChar)
	m = whiteSpaceRegex.ReplaceAllString(m, "")
	return m
}

type Tracers struct {
	BundlerCollectorTracer string
	BundlerExecutionTracer string
	EventTracer string
}

func NewTracers() (*Tracers, error) {
	var et string
	err := fs.WalkDir(
		files,
		"EventTracer.js",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			b, err := fs.ReadFile(files, path)
			if err != nil {
				return err
			}

			et = parse(string(b))
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &Tracers{
		EventTracer: et,
	}, nil
}