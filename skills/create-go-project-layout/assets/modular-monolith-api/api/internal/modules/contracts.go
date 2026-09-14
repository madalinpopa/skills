package modules

import "context"

type Contracts struct{}

func (c *Contracts) Verify(_ context.Context) error {
	return nil
}
