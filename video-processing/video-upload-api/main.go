package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	s3Session *s3.S3
	bucket    = "videos"
)

func init() {
	// Initialize the S3 session
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(os.Getenv("AWS_REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		S3ForcePathStyle: aws.Bool(true), // Use path-style addressing
	})
	if err != nil {
		logrus.Fatalf("Failed to create S3 session: %v", err)
	}
	s3Session = s3.New(sess)
}

func main() {
	e := echo.New()

	// Route to handle video uploads
	e.POST("/upload", handleVideoUpload)

	// Start the server
	logrus.Info("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func handleVideoUpload(c echo.Context) error {
	// Read the uploaded file
	file, err := c.FormFile("video")
	if err != nil {
		logrus.Errorf("Failed to read form file: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Failed to read form file: %v", err)})
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		logrus.Errorf("Failed to open uploaded file: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to open uploaded file: %v", err)})
	}
	defer src.Close()

	// Upload the file to S3
	uploadResult, err := s3Session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Filename),
		Body:   src,
	})
	if err != nil {
		logrus.Errorf("Failed to upload file to S3: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to upload file to S3: %v", err)})
	}

	// fake time sleep since we are using localstack and its too fast
	time.Sleep(5 * time.Minute)

	logrus.Infof("File %s uploaded successfully to S3 with ETag %s", file.Filename, *uploadResult.ETag)
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("File %s uploaded successfully to S3", file.Filename)})
}
