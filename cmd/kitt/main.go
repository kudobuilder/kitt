package main

import (
	"context"
	"os"

	"github.com/Masterminds/semver"

	"github.com/kudobuilder/kitt/pkg/cmd"
)

var (
	version = "0.0.0+dev"
)

func main() {
	ctx := context.Background()

	if err := cmd.New(*semver.MustParse(version)).ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
