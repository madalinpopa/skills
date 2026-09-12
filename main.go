package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/madalinpopa/skills/cmd"
)

func main() {
	os.Exit(run())
}

// run executes the CLI and maps the result to an exit code:
// 0 for success, 1 for runtime errors and 2 for invalid usage.
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := cmd.Execute(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	switch {
	case err == nil:
		return 0
	case errors.Is(err, cmd.ErrUsage):
		return 2
	default:
		return 1
	}
}
