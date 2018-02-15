package main

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/gin-gonic/gin"
	"github.com/tjamet/gin-s3"
)

const usage = `Usage: example --bucket=<bucket>

Starts a server fetching files from s3`

type logger struct{}

func (l logger) Printf(f string, v ...interface{}) { fmt.Printf(f, v...) }

func main() {

	arg, err := docopt.Parse(usage, os.Args[1:], true, "example", false, true)
	if err != nil {
		panic(err)
	}
	r := gin.Default()

	r.Use(ginS3.NewDefault(arg["--bucket"].(string), ginS3.WithLogger(logger{})))
	r.Run(":8080")
}
