package util

import (
	"bytes"
)

// BuildFilePath 构建文件路径
func BuildFilePath(fileSplit string, value ...string) string {
	var buffer bytes.Buffer
	valueArrayLen := len(value)

	for i := 0; i < valueArrayLen; i++ {
		buffer.WriteString(value[i])

		if i+1 == valueArrayLen {
			continue
		}

		buffer.WriteString(fileSplit)
	}

	return buffer.String()
}
