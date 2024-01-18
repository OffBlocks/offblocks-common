package temporal

import (
	"context"

	"github.com/offblocks/offblocks-common/auth"
	"github.com/offblocks/offblocks-common/errors"
	"github.com/offblocks/offblocks-common/types"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

// clientIdPropagator implements the clientId context propagator
type clientIdPropagator struct{}

// NewContextPropagator returns a context propagator that propagates a set of
// string key-value pairs across a workflow
func NewContextPropagator() workflow.ContextPropagator {
	return &clientIdPropagator{}
}

// Inject injects values from context into headers for propagation
func (s *clientIdPropagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	clientId, err := auth.Context{Context: ctx}.ClientId()
	if err != nil {
		return err
	}

	payload, err := converter.GetDefaultDataConverter().ToPayload(clientId)
	if err != nil {
		return err
	}

	writer.Set(string(auth.ClientIdKey), payload)
	return nil
}

// InjectFromWorkflow injects values from context into headers for propagation
func (s *clientIdPropagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	clientId := ctx.Value(auth.ClientIdKey)
	if clientId == nil {
		return errors.ErrUnauthorised
	}

	payload, err := converter.GetDefaultDataConverter().ToPayload(clientId)
	if err != nil {
		return err
	}
	writer.Set(string(auth.ClientIdKey), payload)
	return nil
}

// Extract extracts values from headers and puts them into context
func (s *clientIdPropagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if value, ok := reader.Get(string(auth.ClientIdKey)); ok {
		var clientId types.UUID
		if err := converter.GetDefaultDataConverter().FromPayload(value, &clientId); err != nil {
			return ctx, nil
		}
		ctx = auth.WithClientId(ctx, clientId)
	}

	return ctx, nil
}

// ExtractToWorkflow extracts values from headers and puts them into context
func (s *clientIdPropagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if value, ok := reader.Get(string(auth.ClientIdKey)); ok {
		var clientId types.UUID
		if err := converter.GetDefaultDataConverter().FromPayload(value, &clientId); err != nil {
			return ctx, nil
		}
		ctx = workflow.WithValue(ctx, auth.ClientIdKey, clientId)
	}

	return ctx, nil
}
