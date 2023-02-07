package main

import (
	"bytes"
	"dependentCutting/util"
	"fmt"
	"os"
	"strings"
)

const (
	LibDirSuffix       = "%s\\BOOT-INF\\lib"
	ManifestFileSuffix = "%s\\META-INF\\MANIFEST.MF"
	MainClass          = "org.springframework.boot.loader.PropertiesLauncher"
)

var (
	jarFile, jarDir, manifestFile string
	TargetDirNameSlice            = [...]string{"BOOT-INF", "META-INF", "org", "lib"}
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

func readAllLine(content []byte) []string {
	var start = 0
	var result []string

	for i, c := range content {
		if c == '\n' {
			// start - i
			v := string(content[start:i])
			result = append(result, v)
			start = i
		}
	}

	return result
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

	allLine := readAllLine(content)

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

func main() {
	//if len(os.Args) == 1 {
	//	panic(NotFoundFileError)
	//}
	//
	//filePath := os.Args[1]
	//if filePath == "" {
	//	panic(NotFoundFileError)
	//}

	jarFile = "C:\\Users\\Administrator\\Desktop\\1\\2\\Cloud-KernelApp.jar"
	index := strings.LastIndex(jarFile, "\\")
	jarDir = jarFile[:index]
	jarName := jarFile[index:]
	manifestFile = fmt.Sprintf(ManifestFileSuffix, jarDir)

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

	// 处理manifestFile
	handlerManifestFile()

	// 收尾清理
	clearDir([]string{"lib"})
}
