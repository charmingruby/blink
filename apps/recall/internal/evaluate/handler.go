package evaluate

import (
	"blink/api/proto/pb"
	"blink/lib/queue"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var blinkCooldown = 5 * time.Second

type handler struct {
	pb.UnimplementedBlinkServiceServer

	repo   *tracerRepo
	pubsub *queue.RabbitMQPubSub
}

func (h *handler) BlinkEvaluation(
	ctx context.Context,
	req *pb.BlinkEvaluationRequest,
) (*pb.BlinkEvaluationReply, error) {
	tr, err := h.repo.findTracerByIP(ctx, req.GetIp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if tr.ID == "" {
		msg := pb.BlinkEvaluatedEvent{
			Ip:     req.GetIp(),
			Status: pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_BOOTSTRAP,
		}

		protoMsg, err := proto.Marshal(&msg)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if err := h.pubsub.Publish(ctx, "tracers", protoMsg); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &pb.BlinkEvaluationReply{
			Status:            pb.BlinkStatus_BLINK_STATUS_PENDING,
			RemainingCooldown: 0,
		}, nil
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

	msg := pb.BlinkEvaluatedEvent{
		Ip:     req.GetIp(),
		Status: pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_BLINK,
	}

	protoMsg, err := proto.Marshal(&msg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := h.pubsub.Publish(ctx, "tracers", protoMsg); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BlinkEvaluationReply{
		Status:            pb.BlinkStatus_BLINK_STATUS_PENDING,
		RemainingCooldown: 0,
	}, nil
}
