// Package cmd builds the Cobra command tree for the skills CLI.
package cmd

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"
)

// ErrUsage marks errors caused by invalid command syntax or arguments.
// Callers map it to a usage exit code instead of a runtime failure.
var ErrUsage = errors.New("invalid usage")

// usageError wraps a Cobra parse or argument error so it matches ErrUsage.
type usageError struct{ err error }

func (e usageError) Error() string        { return e.err.Error() }
func (e usageError) Unwrap() error        { return e.err }
func (e usageError) Is(target error) bool { return target == ErrUsage }

// Execute runs the command tree with the given arguments and streams.
// It prints errors to errOut and returns them so the caller can pick the
// exit code. It never exits the process.
func Execute(ctx context.Context, args []string, in io.Reader, out, errOut io.Writer) error {
	root := newRoot(in, out, errOut)
	// Cobra reads os.Args when args is nil, which breaks tests and callers
	// that mean "no arguments".
	if args == nil {
		args = []string{}
	}
	root.SetArgs(args)

	failed, err := root.ExecuteContextC(ctx)
	if err == nil {
		return nil
	}

	failed.PrintErrln(failed.ErrPrefix(), err.Error())
	if errors.Is(err, ErrUsage) {
		failed.PrintErr(failed.UsageString())
	}
	return err
}

func newRoot(in io.Reader, out, errOut io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:   "skills",
		Short: "Install agent skills from a shared store",
		Args:  noArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return c.Help()
		},
		// Execute prints errors and usage itself, so Cobra must stay quiet.
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return usageError{err}
	})
	return root
}

// noArgs rejects positional arguments as a usage error. On the root command
// this is how an unknown subcommand surfaces.
func noArgs(c *cobra.Command, args []string) error {
	if err := cobra.NoArgs(c, args); err != nil {
		return usageError{err}
	}
	return nil
}
