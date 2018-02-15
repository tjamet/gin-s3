package ginS3_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/tjamet/gin-s3"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

type logger struct {
	logs []string
}

func (l *logger) Printf(f string, v ...interface{}) { l.logs = append(l.logs, fmt.Sprintf(f, v...)) }

func TestIntegration(t *testing.T) {

	bucket := "test"

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "latest",
		Env:        []string{"MINIO_SECRET_KEY=EXAMPLEKEY", "MINIO_ACCESS_KEY=EXAMPLE"},
		Cmd:        []string{"server", "/tmp"},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer func() {
		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("EXAMPLE", "EXAMPLEKEY", ""),
		Endpoint:         aws.String(fmt.Sprintf("http://localhost:%s", resource.GetPort("9000/tcp"))),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("eu-west-1"),
	}
	// Create a new bucket using the CreateBucket call.
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		_, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String(bucket), // Required
	}
	_, err = s3Client.CreateBucket(cparams)
	if err != nil {
		log.Fatalf("Could create test bucket: %s", err)
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:          strings.NewReader("hello world"),
		ContentLength: aws.Int64(int64(len("hello world"))),
		ContentType:   aws.String("text/plain"),
		Key:           aws.String("/test.txt"),
		Bucket:        aws.String(bucket),
	})
	if err != nil {
		log.Fatalf("Could create test file: %s", err)
	}

	l := &logger{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test.txt", nil)

	ginS3.NewDefault(bucket,
		ginS3.WithConfig(s3Config),
		ginS3.WithLogger(l),
	)(c)
	if w.Body.String() != "hello world" {
		t.Errorf("Handler body '%s' does not match expectations: 'hello world'", w.Body.String())
	}
}
