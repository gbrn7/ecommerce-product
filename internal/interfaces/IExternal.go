package interfaces

import (
	"context"
	"ecommerce-product/external"
)

type IExternal interface {
	GetProfile(ctx context.Context, token string) (external.Profile, error)
}
