package main

import (
	"os"
	"compress/gzip"
	"archive/tar"
	"io"
	"github.com/fatih/color"
	"fmt"
	"strings"
	"io/ioutil"
	"runtime"
	"github.com/alexflint/go-arg"
)

const (
	PATH_SEPARATOR = string(os.PathSeparator)
)

var total = 0

func main() {
	if len(os.Args) <= 1 ||
		(len(os.Args) == 2 && os.Args[1] == "--help") {
		writeInfo("usage: extract [--targetdir=<path>] [--desdir=<path>]")
		writeInfo("<command> [<args>]")
		return
	}

	var  args struct {
		TargetDir string
		DesDir string
	}
	arg.MustParse(&args)
	writeInfo("[Target Directory]:", args.TargetDir)
	writeInfo("[Destination Directory]:", args.DesDir)
	writeInfo()

	extractAllTarInDirectory(args.TargetDir, args.DesDir);

	writeSuccess("extract done with", total, "files extracted!")
}

func extractAllTarInDirectory(tarsDir string, desDir string) error {
	files, err := ioutil.ReadDir(tarsDir)
	if err != nil {
		return err
	}

	tarsDirFullPath := tarsDir
	if !strings.HasSuffix(tarsDirFullPath, PATH_SEPARATOR ) {
		tarsDirFullPath = tarsDirFullPath + PATH_SEPARATOR
	}

	desDirFullPath := desDir
	if !strings.HasSuffix(desDirFullPath, PATH_SEPARATOR) {
		desDirFullPath = desDirFullPath + PATH_SEPARATOR
	}

	for _, fi := range files {
		fileName := fi.Name()
		filePath := tarsDirFullPath + fileName
		if !strings.HasSuffix(filePath, "tar.gz") {
			continue
		}

		targetRootDirName := fileName[0:len(fileName) - 7];
		targetRootDir := desDirFullPath + targetRootDirName
		err = DeCompress(filePath, targetRootDir);
		if err != nil {
			writeError("decompress file '", filePath, "' failed: ", err)
			return err
		}
	}

	return nil
}

func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	if (!(strings.HasSuffix(dest, "/") ||
		strings.HasSuffix(dest, "\\"))) {
		dest = dest + PATH_SEPARATOR;
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		if runtime.GOOS == `windows` {
			filename = strings.Replace(filename, "/", "\\", -1)
		}
		writeInfo("extracting ", filename)
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		if file != nil {
			io.Copy(file, tr)
		}
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	fileName := string([]rune(name)[0:strings.LastIndex(name, PATH_SEPARATOR)]);
	if strings.HasSuffix(name, PATH_SEPARATOR) {
		// create directory
		err := os.MkdirAll(fileName, 0755)
		if err == nil {
			return nil, nil
		}
		return nil, err
	}
	total = total + 1
	// create file
	return os.Create(name)
}


func writeError(v ...interface{}) {
	content := "[ERROR]: " + fmt.Sprintln(v...)
	color.Red(content)
}

//func writeWarnning(v ...interface{}) {
//	content := "[WARNNING]: " + fmt.Sprintln(v...)
//	color.Yellow(content)
//}

func writeSuccess(v ...interface{}) {
	content := fmt.Sprintln(v...)
	color.Green(content)
}

//func writeTitle(v ...interface{}) {
//	content := fmt.Sprintln(v...)
//	color.Cyan(content)
//}

func writeInfo(v ...interface{}) {
	content := fmt.Sprintln(v...)
	color.White(content)
}

