package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/madalinpopa/skills/cmd"
)

var version string

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return cmd.ExitCode(cmd.Execute(ctx, version, os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
