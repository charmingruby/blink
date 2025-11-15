package commit

import (
	"blink/api/proto/pb"
	"context"

	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/proto"
)

type handler struct {
	db *sqlx.DB
}

func (h *handler) OnBlinkEvaluated(ctx context.Context, body []byte) error {
	var event pb.BlinkEvaluatedEvent

	if err := proto.Unmarshal(body, &event); err != nil {
		return err
	}

	println(event.GetNickname())
	println(event.GetStatus().String())

	return nil
}
