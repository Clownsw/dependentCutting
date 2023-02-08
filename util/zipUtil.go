package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Compress 压缩ZIP文件夹到targetFile
func Compress(sourceDirSlice []string, targetFile string) error {
	_ = os.Remove(targetFile)

	zipFile, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer Close(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer Close(zipWriter)

	for _, sourceDir := range sourceDirSlice {
		err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if path != sourceDir && !info.IsDir() {
				zipFilePath := path[strings.LastIndex(sourceDir, "\\")+1:]

				fileFd, err := os.Open(path)
				if err != nil {
					return err
				}
				defer Close(fileFd)

				header, err := zip.FileInfoHeader(info)
				if err != nil {
					return err
				}

				header.Method = zip.Deflate
				header.Name = zipFilePath
				headerWriter, err := zipWriter.CreateHeader(header)
				if err != nil {
					return err
				}

				if _, err := io.Copy(headerWriter, fileFd); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Decompress 解压ZIP文件到targetDir目录
func Decompress(sourceFile, targetDir string) error {
	reader, err := zip.OpenReader(sourceFile)

	if err != nil {
		return err
	}

	defer Close(reader)

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

// Compress7z 使用7z ZIP压缩
func Compress7z(sixZExeFile string, sourceDirSlice []string, targetFile string) error {
	cmd := exec.Command(sixZExeFile, append([]string{"a", targetFile}, sourceDirSlice...)...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
