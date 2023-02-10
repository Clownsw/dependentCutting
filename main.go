package main

import (
	"bytes"
	"dependentCutting/config"
	"dependentCutting/util"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// clearDir 清理目录中包含在TargetDirNameSlice中的文件夹
func clearDir(excludes []string) {
b:
	for _, dirName := range config.TargetDirNameSlice {
		for _, exclude := range excludes {
			if dirName == exclude {
				continue b
			}
		}

		_ = os.RemoveAll(util.BuildFilePath(config.FileSplit, config.JarDir, dirName))
	}
}

func finish(err error) {
	clearDir([]string{})
	panic(err)
}

// packageManifestFile 打包ManifestFile文件
func packageManifestFile(params map[string]string) []byte {
	var buffer bytes.Buffer

	for key, value := range params {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	return buffer.Bytes()
}

// unPackageManifestFile 解包ManifestFile文件
func unPackageManifestFile(content []byte) map[string]string {
	var paramMap = make(map[string]string)

	allLine := util.ReadAllLine(content)

	for _, line := range allLine {
		splitArray := strings.Split(line, ":")
		if len(splitArray) == 2 {
			key := strings.Trim(strings.Trim(splitArray[0], "\n"), " ")
			value := strings.Trim(strings.Trim(splitArray[1], "\n"), " ")
			paramMap[key] = value
		}
	}

	return paramMap
}

// handlerManifestFile 处理ManifestFile文件
func handlerManifestFile() {
	content, err := os.ReadFile(config.ManifestFile)
	if err != nil {
		finish(err)
	}

	params := unPackageManifestFile(content)
	params["Main-Class"] = config.MainClass
	if err := os.WriteFile(config.ManifestFile, packageManifestFile(params), os.FileMode(0777)); err != nil {
		finish(err)
	}
}

// handlerSnowyLib 复制以snowy-开头的依赖到单独目录
func handlerSnowyLib() error {
	_ = os.RemoveAll(config.SnowyLibDir)
	if err := os.Mkdir(config.SnowyLibDir, os.FileMode(07777)); err != nil {
		return err
	}

	return filepath.Walk(util.BuildFilePath(config.JarDir, "lib"), func(path string, fileInfo fs.FileInfo, err error) error {
		fileName := fileInfo.Name()

		if !fileInfo.IsDir() && len(fileName) >= 5 && fileInfo.Name()[0:5] == "snowy" {
			newFileFd, err := os.OpenFile(util.BuildFilePath(config.FileSplit, config.JarDir, config.SnowyLibDirName, fileInfo.Name()),
				os.O_CREATE|os.O_WRONLY, fileInfo.Mode())
			if err != nil {
				return err
			}

			defer util.Close(newFileFd)

			fileFd, err := os.Open(path)
			if err != nil {
				return err
			}

			defer util.Close(fileFd)

			if _, err := io.Copy(newFileFd, fileFd); err != nil {
				return err
			}
		}

		return nil
	})
}

//goland:noinspection GoBoolExpressions
func main() {
	fmt.Printf("input file: %s, name: %s, dir: %s\n", config.JarFile, config.JarName, config.JarDir)

	// 检查是否存在冲突目录
	clearDir([]string{})

	// 解压jar
	if err := util.Decompress(config.FileSplit, config.JarFile, config.JarDir); err != nil {
		finish(err)
	}

	// 移动lib
	if err := os.Rename(fmt.Sprintf(config.LibDirSuffix, config.JarDir), util.BuildFilePath(config.FileSplit, config.JarDir, "lib")); err != nil {
		finish(err)
	}

	// 将snowy-开头的依赖单独复制到otherLib
	if err := handlerSnowyLib(); err != nil {
		finish(err)
	}

	// 处理manifestFile
	handlerManifestFile()

	// 打包
	if err := util.Compress7z(config.SixZExeFile, func() []string {
		var params []string
		for _, targetDirName := range config.TargetDirNameSlice {
			if targetDirName == "lib" {
				continue
			}

			params = append(params, util.BuildFilePath(config.FileSplit, config.JarDir, targetDirName))
		}

		return params
	}(), util.BuildFilePath(config.FileSplit, config.JarDir, "dc.jar")); err != nil {
		finish(err)
	}

	// 收尾清理
	clearDir([]string{"lib"})
}
