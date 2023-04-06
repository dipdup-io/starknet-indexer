package main

import (
	"context"
	"fmt"
	"io"
)

// Tester -
type Tester interface {
	Test(ctx context.Context) error
	io.Closer
	fmt.Stringer
}
