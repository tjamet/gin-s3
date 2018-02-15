# gin-s3

[![Build Status](https://travis-ci.org/tjamet/gin-s3.svg)](https://travis-ci.org/tjamet/gin-s3)
[![codecov](https://codecov.io/gh/tjamet/gin-s3/branch/master/graph/badge.svg)](https://codecov.io/gh/tjamet/gin-s3)
[![Go Report Card](https://goreportcard.com/badge/github.com/tjamet/gin-s3)](https://goreportcard.com/report/github.com/tjamet/gin-s3)
[![GoDoc](https://godoc.org/github.com/tjamet/gin-s3?status.svg)](https://godoc.org/github.com/tjamet/gin-s3)

a gin handler to fetch files from s3

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/tjamet/gin-s3
```

Import it in your code:

```go
import "github.com/tjamet/gin-s3"
```

### Canonical example:

```go
package main

import (
	"time"

	"github.com/tjamet/gin-s3"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Gets credentials from the environment, the config files or the amazon instance
	router.Use(ginS3.NewDefault("test-bucket"))
	router.Run()
}
```

### Using Specific access keys

```go
func main() {
	router := gin.Default()
	router.Use(ginS3.NewDefault(
        "test-bucket",
        ginS3.AddProvider(
            &credentials.StaticProvider{
                Value: credentials.Value{
                    AccessKeyID:     "EXAMPLE",
                    SecretAccessKey: "EXAMPLEKEY",
                },
            })
        )
    )
	router.Run()
}
```