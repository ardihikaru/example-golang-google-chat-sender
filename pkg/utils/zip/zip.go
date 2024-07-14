package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
)

// CheckFor7ZipPackage checks whether the container / server has 7zip installed or not
// apt install p7zip-full
func CheckFor7ZipPackage(log *logger.Logger) error {
	_, err := exec.LookPath("7za")
	if err != nil {
		log.Warn(fmt.Sprintf("Make sure 7zip is install and include your path: %s", err.Error()))

		return err
	}

	return nil
}

func CreateZipFileWithPassword(log *logger.Logger, zipFilePath, extractedFilePath, zipPasswd, encryptionType string) error {
	log.Debug(fmt.Sprintf("creating zip file: %s", zipFilePath))
	commandString := fmt.Sprintf(`7za a %s %s -p"%s" -mem=%s`, zipFilePath, extractedFilePath, zipPasswd,
		encryptionType)
	commandSlice := strings.Fields(commandString)
	log.Debug(commandString)
	c := exec.Command(commandSlice[0], commandSlice[1:]...)
	err := c.Run()

	return err
}

func ExtractZipFileOrFolderWithPassword(log *logger.Logger, zipFileOrFolderPath, extractedFileOrFolderPath, zipPasswd string) error {
	log.Debug(fmt.Sprintf("Unzipping `%s` to directory `%s`", zipFileOrFolderPath, extractedFileOrFolderPath))
	commandString := fmt.Sprintf(`7za e %s -o%s -p%s -aoa`, zipFileOrFolderPath, extractedFileOrFolderPath, zipPasswd)
	commandSlice := strings.Fields(commandString)
	log.Debug(commandString)
	c := exec.Command(commandSlice[0], commandSlice[1:]...)
	err := c.Run()

	return err
}

// ZipSourceFolder zips source folder
func ZipSourceFolder(source, targetZipFilePath string) error {
	// 1. creates a ZIP file and zip.Writer
	f, err := os.Create(targetZipFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. goes through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. creates a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// sets compression
		header.Method = zip.Deflate

		// 4. sets relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. creates writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(filepath.Clean(path))
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
