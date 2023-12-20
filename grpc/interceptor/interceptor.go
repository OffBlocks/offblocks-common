package interceptor

import (
	"context"

	"github.com/google/uuid"
	"github.com/offblocks/offblocks-common/auth"
	"github.com/offblocks/offblocks-common/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientIdPropagationUnaryClientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	clientId, err := auth.Context{Context: ctx}.ClientId()
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, string(auth.ClientIdKey), clientId.String())
	return invoker(ctx, method, req, reply, cc, opts...)
}

func ClientIdPropagationStreamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	clientId, err := auth.Context{Context: ctx}.ClientId()
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, string(auth.ClientIdKey), clientId.String())
	return streamer(ctx, desc, cc, method, opts...)
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}

func ClientIdPropagationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			clientIdStr, ok := md[string(auth.ClientIdKey)]
			if !ok {
				return nil, errors.ErrUnauthorised
			}
			clientId, err := uuid.Parse(clientIdStr[0])
			if err != nil {
				return nil, errors.ErrUnauthorised
			}
			ctx = auth.WithClientId(ctx, clientId)
		}
		return handler(ctx, req)
	}
}

func ClientIdPropagationStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			clientIdStr, ok := md[string(auth.ClientIdKey)]
			if !ok {
				return errors.ErrUnauthorised
			}
			clientId, err := uuid.Parse(clientIdStr[0])
			if err != nil {
				return errors.ErrUnauthorised
			}
			ctx := auth.WithClientId(ss.Context(), clientId)
			ss = &serverStream{
				ServerStream: ss,
				ctx:          ctx,
			}
		}
		return handler(srv, ss)
	}
}
