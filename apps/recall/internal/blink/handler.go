package blink

import (
	"blink/api/proto/pb"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var blinkCooldown = 5 * time.Second

type handler struct {
	pb.UnimplementedEvaluationServiceServer

	service *service
}

func (h *handler) EvaluateBlinkIntent(ctx context.Context, req *pb.EvaluateBlinkIntentRequest) (*pb.EvaluateBlinkIntentReply, error) {
	rep, err := h.service.evaluateBlinkIntent(ctx, req.GetNickname())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return rep, nil
}
