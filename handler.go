package ginS3

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// Logger defines the interface a type must implement to log errors
type Logger interface {
	Printf(pattern string, v ...interface{})
}

// Client defines the methods of s3.S3 used by the middleware.
// It is provided as an interface to allow mocking for testing purposes
type Client interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// S3 holds the configuration for the s3 handler
type S3 struct {
	// Client holds
	Client Client
	Bucket string
	Logger Logger
}

// Handle implements a gin handler fetching the sources from a s3 bucket
func (s *S3) Handle(c *gin.Context) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(c.Request.URL.Path),
	}
	object, err := s.Client.GetObject(input)
	if err != nil {
		if s.Logger != nil {
			s.Logger.Printf("error fetching object %s from bucket %s: %s\n", c.Request.URL.Path, s.Bucket, err.Error())
		}
		return
	}
	if object.Body != nil {
		defer object.Body.Close()
	}
	// If there is no content length, it is a directory
	if object.ContentLength == nil {
		return
	}
	c.Header("Content-Type", *object.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", *object.ContentLength))
	c.Status(http.StatusOK)
	io.Copy(c.Writer, object.Body)

}
