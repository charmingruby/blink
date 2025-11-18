package blink

import (
	"blink/api/proto/pb"
	"blink/lib/core"
	"blink/lib/http/grpcx"
	"blink/lib/telemetry"
	"context"
	"errors"
)

var (
	ErrBlinkOnCooldown = errors.New("blink is on cooldown")
	ErrProcessingBlink = errors.New("processing blink")
)

type service struct {
	evaluationClient pb.EvaluationServiceClient
}

func newService(evaluationClient pb.EvaluationServiceClient) *service {
	return &service{
		evaluationClient: evaluationClient,
	}
}

func (s *service) emitBlinkIntent(ctx context.Context, nickname string) (float64, error) {
	ctx, span := telemetry.StartSpan(ctx, "blink.service.emitBlinkIntent")
	defer span.End()

	rep, err := s.evaluationClient.EvaluateBlinkIntent(ctx, &pb.EvaluateBlinkIntentRequest{
		Nickname: nickname,
	})
	if err != nil {
		if err := grpcx.TranslateErr(err); err != nil {
			telemetry.RecordError(ctx, err)

			if err.Is(grpcx.ErrNotFound) {
				return 0, core.NewNotFoundError(err.Error())
			}

			return 0, core.NewUnknowClientError(err.Error())
		}
	}

	if rep.Status == pb.BlinkStatus_BLINK_STATUS_PROCESSING {
		return 0, ErrProcessingBlink
	}

	isOnCooldown := rep.RemainingCooldown > 0 &&
		rep.Status == pb.BlinkStatus_BLINK_STATUS_ON_COOLDOWN

	if isOnCooldown {
		telemetry.RecordError(ctx, ErrBlinkOnCooldown)
		return rep.GetRemainingCooldown(), ErrBlinkOnCooldown
	}

	return 0, nil
}
