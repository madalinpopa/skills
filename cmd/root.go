package cmd

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"
)

var ErrUsage = errors.New("invalid usage")

type usageError struct{ err error }

func (e usageError) Error() string        { return e.err.Error() }
func (e usageError) Unwrap() error        { return e.err }
func (e usageError) Is(target error) bool { return target == ErrUsage }

func Execute(ctx context.Context, args []string, in io.Reader, out, errOut io.Writer) error {
	root := newRoot(in, out, errOut)
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
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)
	root.AddCommand(newInitCmd(), newSyncCmd())
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return usageError{err}
	})
	return root
}

func noArgs(c *cobra.Command, args []string) error {
	if err := cobra.NoArgs(c, args); err != nil {
		return usageError{err}
	}
	return nil
}
