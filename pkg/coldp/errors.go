package coldp

import "fmt"

type ErrorFileMissing struct {
	Path string
}

func (e *ErrorFileMissing) Error() string {
	return fmt.Sprintf("file not found '%s'", e.Path)
}

// ErrBadDir is returned when the directory is in an unknown state, or is not a
// directory.
type ErrBadDir struct {
	Dir string
}

func (e *ErrBadDir) Error() string {
	return fmt.Sprintf("bad directory '%s'", e.Dir)
}

// ErrUnzip is returned when the extraction of the ColDP file fails.
type ErrUnzip struct {
	Path string
	Err  error
}

func (e *ErrUnzip) Error() string {
	return fmt.Sprintf("extracting '%s' failed: %v", e.Path, e.Err)
}
