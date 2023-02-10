package main

import (
	"bytes"
	"dependentCutting/util"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	LibDirSuffix       = "%s\\BOOT-INF\\lib"
	ManifestFileSuffix = "%s\\META-INF\\MANIFEST.MF"
	MainClass          = "org.springframework.boot.loader.PropertiesLauncher"
)

var (
	NotFoundFileError             = errors.New("not found file")
	SixZExeFile                   string
	jarFile, jarDir, manifestFile string
	TargetDirNameSlice            = [...]string{"BOOT-INF", "META-INF", "org", "lib"}
	SnowyLibDirName               = "otherLib"
	snowyLibDir                   string
)

// clearDir 清理目录中包含在TargetDirNameSlice中的文件夹
func clearDir(excludes []string) {
b:
	for _, dirName := range TargetDirNameSlice {
		for _, exclude := range excludes {
			if dirName == exclude {
				continue b
			}
		}

		_ = os.RemoveAll(fmt.Sprintf("%s\\%s", jarDir, dirName))
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
	content, err := os.ReadFile(manifestFile)
	if err != nil {
		finish(err)
	}

	params := unPackageManifestFile(content)
	params["Main-Class"] = MainClass
	if err := os.WriteFile(manifestFile, packageManifestFile(params), os.FileMode(0777)); err != nil {
		finish(err)
	}
}

// handlerSnowyLib 复制以snowy-开头的依赖到单独目录
func handlerSnowyLib() error {
	_ = os.RemoveAll(snowyLibDir)
	if err := os.Mkdir(snowyLibDir, os.FileMode(07777)); err != nil {
		return err
	}

	return filepath.Walk(fmt.Sprintf("%s\\lib", jarDir), func(path string, fileInfo fs.FileInfo, err error) error {
		fileName := fileInfo.Name()

		if !fileInfo.IsDir() && len(fileName) >= 5 && fileInfo.Name()[0:5] == "snowy" {
			newFileFd, err := os.OpenFile(fmt.Sprintf("%s\\%s\\%s", jarDir, SnowyLibDirName, fileInfo.Name()), os.O_CREATE|os.O_WRONLY, fileInfo.Mode())
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

func main() {
	SixZExeFile = fmt.Sprintf("%s\\7z.exe", util.GetCurrentExecuteDir())

	if len(os.Args) == 1 {
		panic(NotFoundFileError)
	}

	jarFile = os.Args[1]
	index := strings.LastIndex(jarFile, "\\")
	jarDir = jarFile[:index]
	jarName := jarFile[index:]
	manifestFile = fmt.Sprintf(ManifestFileSuffix, jarDir)
	snowyLibDir = fmt.Sprintf("%s\\%s", jarDir, SnowyLibDirName)

	fmt.Printf("input file: %s, name: %s, dir: %s\n", jarFile, jarName, jarDir)

	// 检查是否存在冲突目录
	clearDir([]string{})

	// 解压jar
	if err := util.Decompress(jarFile, jarDir); err != nil {
		finish(err)
	}

	// 移动lib
	if err := os.Rename(fmt.Sprintf(LibDirSuffix, jarDir), fmt.Sprintf("%s\\lib", jarDir)); err != nil {
		finish(err)
	}

	// 将snowy-开头的依赖单独复制到otherLib
	if err := handlerSnowyLib(); err != nil {
		finish(err)
	}

	// 处理manifestFile
	handlerManifestFile()

	// 打包
	if err := util.Compress7z(SixZExeFile, func() []string {
		var params []string
		for _, targetDirName := range TargetDirNameSlice {
			if targetDirName == "lib" {
				continue
			}

			params = append(params, fmt.Sprintf("%s\\%s", jarDir, targetDirName))
		}

		return params
	}(), fmt.Sprintf("%s\\dc.jar", jarDir)); err != nil {
		finish(err)
	}

	// 收尾清理
	clearDir([]string{"lib"})
}
