package util

import "io"

// Close close a io.Closer implement
func Close[T io.Closer](t T) {
	_ = t.Close()
}
