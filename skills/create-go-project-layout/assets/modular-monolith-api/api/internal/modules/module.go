package modules

import "context"

type Name string

type Module interface {
	Name() Name
	Init(ctx context.Context) error
	RegisterContracts(ctx context.Context, c *Contracts) error
	Connect(ctx context.Context, c *Contracts) error
	RegisterHTTP(ctx context.Context, e EchoRouter) error
}
