package blink

import (
	"blink/api/proto/pb"
	"blink/lib/telemetry"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	pb.UnimplementedEvaluationServiceServer

	service *service
}

func (h *handler) EvaluateBlinkIntent(ctx context.Context, req *pb.EvaluateBlinkIntentRequest) (*pb.EvaluateBlinkIntentReply, error) {
	ctx, span := telemetry.StartSpan(ctx, "blink.handler.EvaluateBlinkIntent")
	defer span.End()

	rep, err := h.service.evaluateBlinkIntent(ctx, req.GetNickname())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return rep, nil
}
