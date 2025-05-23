package interceptor

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/message"
)

type State interface {
	Ctx() context.Context
	GetClientID() (ClientID, error)
	SetClientID(id ClientID) error
	Write(ctx context.Context, msg message.Message) error
}
