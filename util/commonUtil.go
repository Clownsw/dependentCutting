package util

import (
	"io"
	"os"
	"path/filepath"
)

// Close a io.Closer implement
func Close[T io.Closer](t T) {
	_ = t.Close()
}

// GetCurrentExecuteDir 获取当前执行目录
func GetCurrentExecuteDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(ex)
}
