package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-toast/toast"
	"github.com/sqweek/dialog"
)

//go:embed assets/icon.ico
var iconData []byte

func main() {
	iconPath := getIconPath()

	overlayImagePath, err := getImagePath()
	if err != nil {
		if err != dialog.ErrCancelled {
			notifyUser(iconPath, "File Selection Error", err.Error())
		}
		return
	}

	fileExportPath, err := setFileExportPath(overlayImagePath)
	if err != nil {
		notifyUser(iconPath, "Setting Output Directory Error", err.Error())
		return
	}

	notifyUser(iconPath, "Image Processing Started", "")

	padding := 40
	r := 0
	g := 0
	b := 0
	runGimpPlugin(padding, r, g, b, overlayImagePath, fileExportPath)

	notifyUser(iconPath, "Image Processing Finished", "Saved to: "+fileExportPath)
}

func getImagePath() (string, error) {
	argumentsPassed := len(os.Args)

	if argumentsPassed > 2 {
		return "", errors.New("Dropping multiple files is not supported")
	}

	if argumentsPassed > 1 {
		// File drop
		return os.Args[1], nil
	}

	// File dialog
	imagePath, err := dialog.File().Filter("Images", "jpg", "jpeg").Load()

	if err != nil {
		return "", err
	}

	return imagePath, nil
}

func setFileExportPath(overlayImagePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	outputFolder := filepath.Join(homeDir, "Desktop", "InstaPictures")

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err = os.MkdirAll(outputFolder, 0755)
		if err != nil {
			return "", err
		}
	}

	overlayImageExt := filepath.Ext(overlayImagePath)
	baseName := strings.TrimSuffix(filepath.Base(overlayImagePath), overlayImageExt)
	outputImageName := baseName + "-output" + overlayImageExt
	exportPath := filepath.Join(outputFolder, outputImageName)
	return ifPathExistsReturnNewValidOne(exportPath), nil
}

func ifPathExistsReturnNewValidOne(path string) string {
	counter := 1
	pathToCheck := path

	for {
		_, err := os.Stat(pathToCheck)
		if os.IsNotExist(err) {
			return pathToCheck
		}

		fileExtension := filepath.Ext(path)
		pathToCheck = fmt.Sprintf("%s-%d%s", strings.TrimSuffix(path, fileExtension), counter, fileExtension)
		counter++
	}
}

func runGimpPlugin(padding, r, g, b int, imageOverlayPath, exportPath string) {
	batchCmd := fmt.Sprintf(`pdb.python_fu_create_new_insta_image(%d, %d, %d, %d, "%s", "%s")`,
		padding, r, g, b,
		strings.ReplaceAll(imageOverlayPath, `\`, `/`),
		strings.ReplaceAll(exportPath, `\`, `/`))

	cmd := exec.Command("gimp-2.10", "--batch-interpreter=python-fu-eval", "-i", "-b", batchCmd)

	// Hide GIMP console window on Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	// For debugging GIMP
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	cmd.Run()
}

func notifyUser(iconPath, title, message string) {
	notification := toast.Notification{
		AppID:   "InstaSquare",
		Title:   title,
		Message: message,
		Icon:    iconPath,
	}
	notification.Push()
}

func getIconPath() string {
	tempDirectory := os.TempDir()
	tempIconFilePath := filepath.Join(tempDirectory, "embedded_icon.ico")
	os.WriteFile(tempIconFilePath, iconData, 0644)

	return tempIconFilePath
}
