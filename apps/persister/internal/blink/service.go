package blink

import (
	"blink/api/proto/pb"
	"blink/lib/core"
	"blink/lib/lock"
	"blink/lib/telemetry"
	"context"
	"time"
)

type service struct {
	txManager *tracerBlinkTransactionManager
	lock      *lock.RedisLock
}

func newService(txManager *tracerBlinkTransactionManager, redisLock *lock.RedisLock) *service {
	return &service{
		txManager: txManager,
		lock:      redisLock,
	}
}

func (s *service) bootstrapTracer(ctx context.Context, evt *pb.BlinkEvaluatedEvent) error {
	ctx, span := telemetry.StartSpan(ctx, "blink.service.bootstrapTracer")
	defer span.End()

	idempotencyKey := "processed:" + evt.GetIdempotencyKey()
	retryKey := "retry:" + evt.GetIdempotencyKey()

	processed, err := s.lock.CheckIdempotency(ctx, idempotencyKey)
	if err != nil {
		return err
	}

	if processed {
		return nil
	}

	retryCount, err := s.lock.IncrementRetry(ctx, retryKey, 24*time.Hour)
	if err != nil {
		return err
	}

	if retryCount > 3 {
		return nil
	}

	err = s.txManager.executeInTransaction(ctx, func(tracerRepo *tracerRepository, blinkRepo *blinkRepository) error {
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

	if err != nil {
		return err
	}

	return s.lock.MarkIdempotent(ctx, idempotencyKey, 24*time.Hour)
}

func (s *service) createBlink(ctx context.Context, evt *pb.BlinkEvaluatedEvent) error {
	ctx, span := telemetry.StartSpan(ctx, "blink.service.createBlink")
	defer span.End()

	idempotencyKey := "processed:" + evt.GetIdempotencyKey()
	retryKey := "retry:" + evt.GetIdempotencyKey()

	processed, err := s.lock.CheckIdempotency(ctx, idempotencyKey)
	if err != nil {
		return err
	}

	if processed {
		return nil
	}

	retryCount, err := s.lock.IncrementRetry(ctx, retryKey, 24*time.Hour)
	if err != nil {
		return err
	}

	if retryCount > 3 {
		return nil
	}

	err = s.txManager.executeInTransaction(ctx, func(tracerRepo *tracerRepository, blinkRepo *blinkRepository) error {
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

	if err != nil {
		return err
	}

	return s.lock.MarkIdempotent(ctx, idempotencyKey, 24*time.Hour)
}
