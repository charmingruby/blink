package blink

import (
	"blink/api/proto/pb"
	"blink/lib/core"
	"blink/lib/http/grpcx"
	"context"
	"errors"
)

var ErrBlinkOnCooldown = errors.New("blink is on cooldown")

type service struct {
	evaluationClient pb.EvaluationServiceClient
}

func newService(evaluationClient pb.EvaluationServiceClient) *service {
	return &service{
		evaluationClient: evaluationClient,
	}
}

func (s *service) emitBlinkIntent(ctx context.Context, nickname string) (float64, error) {
	rep, err := s.evaluationClient.EvaluateBlinkIntent(ctx, &pb.EvaluateBlinkIntentRequest{
		Nickname: nickname,
	})
	if err != nil {
		if err := grpcx.TranslateErr(err); err != nil {
			if err.Is(grpcx.ErrNotFound) {
				return 0, core.NewNotFoundError(err.Error())
			}

			return 0, core.NewUnknowClientError(err.Error())
		}
	}

	isOnCooldown := rep.RemainingCooldown > 0 &&
		rep.Status == pb.BlinkStatus_BLINK_STATUS_ON_COOLDOWN

	if isOnCooldown {
		return rep.GetRemainingCooldown(), ErrBlinkOnCooldown
	}

	return 0, nil
}
