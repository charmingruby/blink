package evaluate

import (
	"blink/api/proto/pb"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var blinkCooldown = 5 * time.Second

type handler struct {
	pb.UnimplementedBlinkServiceServer

	db *sqlx.DB
}

func (h *handler) BlinkEvaluation(
	ctx context.Context,
	req *pb.BlinkEvaluationRequest,
) (*pb.BlinkEvaluationReply, error) {
	tr, err := findTracerByIP(ctx, h.db, req.GetIp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if tr.ID == "" {
		return nil, status.Error(codes.NotFound, "tracer not found")
	}

	if tr.LastBlinkAt != nil {
		timeSince := time.Since(*tr.LastBlinkAt)

		if timeSince < blinkCooldown {
			remaining := blinkCooldown - timeSince

			return &pb.BlinkEvaluationReply{
				Status:            pb.BlinkStatus_BLINK_STATUS_ON_COOLDOWN,
				RemainingCooldown: remaining.Seconds(),
			}, nil
		}
	}

	return &pb.BlinkEvaluationReply{
		Status:            pb.BlinkStatus_BLINK_STATUS_PENDING,
		RemainingCooldown: 0,
	}, nil
}
