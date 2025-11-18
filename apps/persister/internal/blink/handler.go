package blink

import (
	"blink/api/proto/pb"
	"blink/lib/telemetry"
	"context"
	"errors"

	"google.golang.org/protobuf/proto"
)

var ErrUnspecifiedEvent = errors.New("unspecified event")

type handler struct {
	log     *telemetry.Logger
	service *service
}

func (h *handler) onBlinkEvaluated(ctx context.Context, body []byte) error {
	ctx, span := telemetry.StartSpan(ctx, "blink.handler.onBlinkEvaluated")
	defer span.End()

	var evt pb.BlinkEvaluatedEvent
	if err := proto.Unmarshal(body, &evt); err != nil {
		return err
	}

	switch evt.GetStatus() {
	case pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_BOOTSTRAP:
		return h.service.bootstrapTracer(ctx, &evt)
	case pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_CREATE:
		return h.service.createBlink(ctx, &evt)
	default:
		return ErrUnspecifiedEvent
	}
}
