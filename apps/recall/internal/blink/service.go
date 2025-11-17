package blink

import (
	"blink/api/proto/pb"
	"blink/lib/core"
	"blink/lib/queue"
	"context"
	"time"

	"google.golang.org/protobuf/proto"
)

type service struct {
	tracerRepo *tracerRepository
	pubsub     *queue.RabbitMQPubSub
	queueName  string
}

func newService(tracerRepo *tracerRepository, pubsub *queue.RabbitMQPubSub, queueName string) *service {
	return &service{
		tracerRepo: tracerRepo,
		pubsub:     pubsub,
		queueName:  queueName,
	}
}

func (s *service) evaluateBlinkIntent(ctx context.Context, nickname string) (*pb.EvaluateBlinkIntentReply, error) {
	tr, err := s.tracerRepo.findByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}

	if tr.ID == "" {
		return s.dispatchBootstrapTracerEvent(ctx, nickname)
	}

	if tr.LastBlinkAt != nil {
		timeSince := time.Since(*tr.LastBlinkAt)

		if timeSince < blinkCooldown {
			remaining := blinkCooldown - timeSince

			return &pb.EvaluateBlinkIntentReply{
				Status:            pb.BlinkStatus_BLINK_STATUS_ON_COOLDOWN,
				RemainingCooldown: remaining.Seconds(),
			}, nil
		}
	}

	return s.dispatchCreateBlinkEvent(ctx, tr)
}

func (s *service) dispatchBootstrapTracerEvent(ctx context.Context, nickname string) (*pb.EvaluateBlinkIntentReply, error) {
	msg := pb.BlinkEvaluatedEvent{
		Nickname:           nickname,
		Status:             pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_BOOTSTRAP,
		CurrentBlinksCount: 0,
		TracerId:           nil,
	}

	protoMsg, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	if err := s.pubsub.Publish(ctx, s.queueName, protoMsg); err != nil {
		return nil, err
	}

	return &pb.EvaluateBlinkIntentReply{
		Status:            pb.BlinkStatus_BLINK_STATUS_SUCCESS,
		RemainingCooldown: 0,
	}, nil
}

func (s *service) dispatchCreateBlinkEvent(ctx context.Context, tr core.Tracer) (*pb.EvaluateBlinkIntentReply, error) {
	msg := pb.BlinkEvaluatedEvent{
		Nickname:           tr.Nickname,
		Status:             pb.BlinkEvaluationStatus_BLINK_EVALUATION_STATUS_CREATE,
		CurrentBlinksCount: int32(tr.TotalBlinks),
		TracerId:           &tr.ID,
	}

	protoMsg, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	if err := s.pubsub.Publish(ctx, s.queueName, protoMsg); err != nil {
		return nil, err
	}

	return &pb.EvaluateBlinkIntentReply{
		Status:            pb.BlinkStatus_BLINK_STATUS_SUCCESS,
		RemainingCooldown: 0,
	}, nil
}
