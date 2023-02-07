package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// Decompress 解压ZIP文件到targetDir目录
func Decompress(sourceFile, targetDir string) error {
	reader, err := zip.OpenReader(sourceFile)

	if err != nil {
		return err
	}

	defer func() {
		_ = reader.Close()
	}()

	for _, file := range reader.File {
		filePath := fmt.Sprintf("%s\\%s", targetDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}

			continue
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		reader, err := file.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, reader); err != nil {
			return err
		}

		_ = reader.Close()
		_ = dstFile.Close()
	}

	return nil
}
