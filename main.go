package main

import (
	"os"

	"github.com/knqyf263/labeler/cmd"
	"github.com/knqyf263/labeler/logs"
	"github.com/tonglil/versioning"
)

var version string

func init() {
	versioning.Set(version)

	logs.Output = os.Stdout
}

func main() {
	cmd.Execute()
}
