package blink

import (
	"blink/api/proto/pb"
	"blink/lib/core"
	"context"
	"time"
)

type service struct {
	txManager *tracerBlinkTransactionManager
}

func newService(txManager *tracerBlinkTransactionManager) *service {
	return &service{
		txManager: txManager,
	}
}

func (s *service) bootstrapTracer(ctx context.Context, evt *pb.BlinkEvaluatedEvent) error {
	return s.txManager.executeInTransaction(ctx, func(tracerRepo *tracerRepository, blinkRepo *blinkRepository) error {
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

func (s *service) createBlink(ctx context.Context, evt *pb.BlinkEvaluatedEvent) error {
	return s.txManager.executeInTransaction(ctx, func(tracerRepo *tracerRepository, blinkRepo *blinkRepository) error {
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
