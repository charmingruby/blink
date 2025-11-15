package commit

import (
	"blink/api/proto/pb"
	"blink/lib/core"
	"blink/lib/telemetry"
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
)

var ErrUnspecifiedEvent = errors.New("unspecified event")

type handler struct {
	log       *telemetry.Logger
	txManager *tracerBlinkTransactionManager
}

func (h *handler) onBlinkEvaluated(ctx context.Context, body []byte) error {
	var evt pb.BlinkEvaluatedEvent
	if err := proto.Unmarshal(body, &evt); err != nil {
		return err
	}

	switch evt.Status {
	case pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_BOOTSTRAP:
		return h.handleBootstrapEvent(ctx, &evt)
	case pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_CREATE:
		return h.handleCreateEvent(ctx, &evt)
	default:
		return ErrUnspecifiedEvent
	}
}

func (h *handler) handleBootstrapEvent(ctx context.Context, evt *pb.BlinkEvaluatedEvent) error {
	return h.txManager.executeInTransaction(ctx, func(tracerRepo *tracerRepository, blinkRepo *blinkRepository) error {
		tr := core.Tracer{
			ID:          core.GenerateID(),
			Nickname:    evt.GetNickname(),
			TotalBlinks: 0,
			CreatedAt:   time.Now(),
			LastBlinkAt: nil,
			UpdatedAt:   nil,
		}

		bl := core.Blink{
			ID:        core.GenerateID(),
			TracerID:  tr.ID,
			CreatedAt: time.Now(),
		}

		if err := tracerRepo.create(ctx, tr); err != nil {
			return err
		}

		if err := blinkRepo.create(ctx, bl); err != nil {
			return err
		}

		now := time.Now()
		tr.TotalBlinks += 1
		tr.UpdatedAt = &now
		tr.LastBlinkAt = &now

		if err := tracerRepo.update(ctx, tr); err != nil {
			return err
		}

		return nil
	})
}

func (h *handler) handleCreateEvent(ctx context.Context, evt *pb.BlinkEvaluatedEvent) error {
	return h.txManager.executeInTransaction(ctx, func(tracerRepo *tracerRepository, blinkRepo *blinkRepository) error {
		blinkID := core.GenerateID()

		tr := core.Tracer{
			ID:          evt.GetTracerId(),
			Nickname:    evt.GetNickname(),
			TotalBlinks: int(evt.GetCurrentBlinksCount()),
		}

		bl := core.Blink{
			ID:        blinkID,
			TracerID:  tr.ID,
			CreatedAt: time.Now(),
		}

		if err := blinkRepo.create(ctx, bl); err != nil {
			return err
		}

		now := time.Now()
		tr.TotalBlinks += 1
		tr.UpdatedAt = &now
		tr.LastBlinkAt = &now

		if err := tracerRepo.update(ctx, tr); err != nil {
			return err
		}

		return nil
	})
}
