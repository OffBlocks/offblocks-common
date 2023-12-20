package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/offblocks/offblocks-common/errors"
	"github.com/offblocks/offblocks-common/types"
)

type key string

const (
	ClientIdKey key = "client-id"
)

type Context struct {
	context.Context
}

func (c Context) ClientId() (types.UUID, error) {
	clientId := c.Value(ClientIdKey)
	if clientId == nil {
		return types.UUID{UUID: uuid.Nil}, errors.ErrUnauthorised
	}
	return clientId.(types.UUID), nil
}

func WithClientId(ctx context.Context, clientId types.UUID) Context {
	return Context{context.WithValue(ctx, ClientIdKey, clientId)}
}
