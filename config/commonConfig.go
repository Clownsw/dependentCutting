package config

import (
	"dependentCutting/util"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	MainClass = "org.springframework.boot.loader.PropertiesLauncher"
)

var (
	NotFoundFileError = errors.New("not found file")

	LibDirSuffix       string
	ManifestFileSuffix string
	SixZExeFile        string
	FileSplit          string
	JarDir             string
	JarName            string
	ManifestFile       string
	JarFile            string
	TargetDirNameSlice = [...]string{"BOOT-INF", "META-INF", "org", "lib"}
	SnowyLibDirName    = "otherLib"
	SnowyLibDir        string
)

//goland:noinspection GoBoolExpressions
func init() {
	if runtime.GOOS == "windows" {
		SixZExeFile = util.BuildFilePath(FileSplit, util.GetCurrentExecuteDir(), "\\7z.exe")
		FileSplit = "\\"
	} else {
		SixZExeFile = "7z"
		FileSplit = "/"
	}

	if len(os.Args) == 1 {
		panic(NotFoundFileError)
	}

	JarFile = os.Args[1]
	index := strings.LastIndex(JarFile, FileSplit)

	if index == -1 {
		panic("error jar file path")
	}

	LibDirSuffix = util.BuildFilePath(FileSplit, "%s", "BOOT-INF", "lib")
	ManifestFileSuffix = util.BuildFilePath(FileSplit, "%s", "META-INF", "MANIFEST.MF")
	JarDir = JarFile[:index]
	JarName = JarFile[index:]
	ManifestFile = fmt.Sprintf(ManifestFileSuffix, JarDir)
	SnowyLibDir = util.BuildFilePath(FileSplit, JarDir, SnowyLibDirName)
}
