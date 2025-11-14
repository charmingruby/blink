package evaluate

import (
	"blink/api/proto/pb"
	"context"
)

type handler struct {
	pb.UnimplementedBlinkServiceServer
}

func (h *handler) BlinkEvaluation(
	ctx context.Context,
	req *pb.BlinkEvaluationRequest,
) (*pb.BlinkEvaluationReply, error) {
	return nil, nil
}
