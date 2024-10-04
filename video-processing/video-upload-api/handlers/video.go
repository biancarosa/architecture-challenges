package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func UploadVideo(c echo.Context) error {
	// Read the uploaded file from the request
	file, err := c.FormFile("video")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read the uploaded file",
		})
	}

	// Open the file for writing
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open the file for writing",
		})
	}
	defer src.Close()

	// Create a new file on the server
	dst, err := os.Create("uploads/" + file.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create the file on the server",
		})
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to copy the file",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("File %s uploaded successfully", file.Filename),
	})
}