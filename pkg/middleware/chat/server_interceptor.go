package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/process"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type ServerInterceptor struct {
	commonInterceptor
	interceptor.Interceptor
	rooms interfaces.RoomManager
}

func (i *ServerInterceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	return i.commonInterceptor.BindSocketConnection(connection, writer, reader)
}

func (i *ServerInterceptor) Init(connection interceptor.Connection) error {
	s, err := i.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while init; err: %s", err.Error())
	}

	w, ok := s.(interfaces.CanWriteMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	// TODO: SEND IDENT MESSAGE
	if err := w.Write(); err != nil {
		return fmt.Errorf("error while init; err: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(s.Ctx(), 10*time.Second)
	defer cancel()

	p := process.NewWaitUntilIdentComplete(ctx, process.WithTickerDuration(500*time.Millisecond))
	if err := p.Process(nil, s); err != nil {
		return fmt.Errorf("error while init; err: %s", err.Error())
	}

	return nil
}

func (i *ServerInterceptor) UnBindSocketConnection(connection interceptor.Connection) {

}

func (i *ServerInterceptor) Close() error {
	return nil
}

func (i *ServerInterceptor) CreateRoom(id types.RoomID, allowed []types.ClientID, ttl time.Duration) (interfaces.Room, error) {
	return i.rooms.CreateRoom(id, allowed, ttl)
}

func (i *ServerInterceptor) DeleteRoom(id types.RoomID) error {
	return i.rooms.DeleteRoom(id)
}

func (i *ServerInterceptor) WriteRoomMessage(roomid types.RoomID, msg message.Message, from types.ClientID, tos ...types.ClientID) error {
	w, ok := i.rooms.(interfaces.CanWriteRoomMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return w.WriteRoomMessage(roomid, msg, from, tos...)
}

func (i *ServerInterceptor) Add(id types.RoomID, s interfaces.State) error {
	a, ok := i.rooms.(interfaces.CanAdd)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return a.Add(id, s)
}

func (i *ServerInterceptor) Remove(id types.RoomID, s interfaces.State) error {
	r, ok := i.rooms.(interfaces.CanRemove)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return r.Remove(id, s)
}
