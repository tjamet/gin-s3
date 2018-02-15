package ginS3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func checkExistingBucket(c *s3.S3, bucket string) error {
	output, err := c.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	for _, b := range output.Buckets {
		if bucket == *b.Name {
			return nil
		}
	}
	return fmt.Errorf("failed to find bucket %s", bucket)
}

// NewDefault instanciates a new S3 handler with the default S3 client
func NewDefault(bucket string, modifiers ...buildFunc) gin.HandlerFunc {
	s, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	b := builder{
		region:              "eu-west-1",
		credentialProviders: []credentials.Provider{},
		config:              aws.NewConfig(),
	}

	for _, f := range modifiers {
		b = f(b)
	}

	if b.config.Credentials == nil {
		if len(b.credentialProviders) == 0 {
			b.credentialProviders = []credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
				&ec2rolecreds.EC2RoleProvider{
					Client: ec2metadata.New(s),
				},
			}
		}
		creds := credentials.NewChainCredentials(b.credentialProviders)
		b.config.WithCredentials(creds)
	}

	if b.config.Region == nil {
		region, err := s3manager.GetBucketRegion(aws.BackgroundContext(), s, bucket, b.region)
		if err != nil {
			panic(err)
		}
		b.config.WithRegion(region)
	}

	c := s3.New(s, b.config)

	err = checkExistingBucket(c, bucket)
	if err != nil {
		panic(err)
	}

	handler := &S3{
		Client: c,
		Bucket: bucket,
		Logger: b.logger,
	}
	return handler.Handle
}

type buildFunc func(builder) builder

type builder struct {
	credentialProviders []credentials.Provider
	region              string
	logger              Logger
	config              *aws.Config
}

// AddProvider adds a credential provider to the s3 client
func AddProvider(provider credentials.Provider) buildFunc {
	return func(b builder) builder {
		r := b
		r.credentialProviders = append(r.credentialProviders, provider)
		return r
	}
}

// WithLogger sets the logger for the handler
func WithLogger(logger Logger) buildFunc {
	return func(b builder) builder {
		r := b
		r.logger = logger
		return r
	}
}

// WithRegion sets the default s3 region
func WithRegion(region string) buildFunc {
	return func(b builder) builder {
		r := b
		r.region = region
		return r
	}
}

// WithConfig sets the s3 endpoint
func WithConfig(config *aws.Config) buildFunc {
	return func(b builder) builder {
		r := b
		r.config = config
		return r
	}
}
