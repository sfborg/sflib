package dwca

import (
	"context"

	"github.com/sfborg/sflib/pkg/arch"
)

type Archive interface {
	arch.Packager

	Meta() *Meta

	EML() *EML

	CoreSlice(offset, limit int) ([][]string, error)

	CoreStream(
		ctx context.Context,
		chCore chan<- []string,
	) (int, error)

	ExtensionSlice(index, offset, limit int) ([][]string, error)

	ExtensionStream(
		ctx context.Context,
		index int,
		ch chan<- []string,
	) (int, error)
}
