// Package images provides metadata extraction for image file formats including JPEG, PNG, and GIF.
package images

import (
	"image"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// extractImageDimensions extracts width and height from an image file.
// This helper function eliminates code duplication across image extractors..
func extractImageDimensions(filePath, format string, logCmd models.CommandLogger) (width, height int, err error) {
	startTime := time.Now()
	file, openErr := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	endTime := time.Now()
	exitCode := 0
	if openErr != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "os.Open", []string{filePath, format + " metadata"}, startTime, endTime, exitCode, openErr, "", filePath)
	}
	if openErr != nil {
		return 0, 0, openErr
	}
	defer func() {
		startTime := time.Now()
		closeErr := file.Close()
		endTime := time.Now()
		exitCode := 0
		if closeErr != nil {
			exitCode = 1
		}
		if logCmd != nil {
			logCmd(utils.GenerateRandomID(), "file.Close", []string{filePath, format + " metadata"}, startTime, endTime, exitCode, closeErr, "", filePath)
		}
	}()

	startTime = time.Now()
	config, _, decodeErr := image.DecodeConfig(file)
	endTime = time.Now()
	exitCode = 0
	if decodeErr != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "image.DecodeConfig", []string{filePath, format}, startTime, endTime, exitCode, decodeErr, "", filePath)
	}
	if decodeErr == nil {
		width = config.Width
		height = config.Height
	}

	return width, height, decodeErr
}
