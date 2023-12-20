package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/offblocks/offblocks-common/errors"
)

type key string

const (
	ClientIdKey key = "client-id"
)

type Context struct {
	context.Context
}

func (c Context) ClientId() (uuid.UUID, error) {
	clientId := c.Value(ClientIdKey)
	if clientId == nil {
		return uuid.Nil, errors.ErrUnauthorised
	}
	return clientId.(uuid.UUID), nil
}

func WithClientId(ctx context.Context, clientId uuid.UUID) Context {
	return Context{context.WithValue(ctx, ClientIdKey, clientId)}
}
