package chat

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/config"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type Interceptor struct {
	interceptor.NoOpInterceptor
	localMessageRegistry message.Registry
	states               interfaces.StateManager
	rooms                interfaces.RoomManager
	config               config.Config
}

func (i *Interceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	ctx, cancel := context.WithCancel(i.Ctx())

	if err := i.states.SetState(connection, state.NewState(ctx, cancel, i.config, connection, writer, reader)); err != nil {
		cancel()
		return nil, nil, err
	}

	return writer, reader, nil
}

func (i *Interceptor) Init(connection interceptor.Connection) error {
	// TODO: ADD DEFAULT ROOM CREATION / JOIN ROOM
	return nil
}

func (i *Interceptor) InterceptSocketWriter(writer interceptor.Writer) interceptor.Writer {
	return interceptor.WriterFunc(func(connection interceptor.Connection, message message.Message) error {
		return writer.Write(connection, message)
	})
}

func (i *Interceptor) InterceptSocketReader(reader interceptor.Reader) interceptor.Reader {
	return interceptor.ReaderFunc(func(connection interceptor.Connection) (message.Message, error) {
		msg, err := reader.Read(connection)
		if err != nil {
			return msg, err
		}

		if !i.localMessageRegistry.Check(msg.GetProtocol()) {
			return msg, nil
		}

		m, ok := msg.(interceptor.Message)
		if !ok {
			return msg, errors.ErrInterfaceMisMatch
		}

		if err := m.ReadProcess(i, connection); err != nil {
			return nil, err
		}

		return m.GetNext(i.GetMessageRegistry())
	})
}

func (i *Interceptor) UnBindSocketConnection(connection interceptor.Connection) {

}

func (i *Interceptor) Close() error {
	return nil
}

func (i *Interceptor) GetState(connection interceptor.Connection) (interfaces.State, error) {
	return i.states.GetState(connection)
}

func (i *Interceptor) CreateRoom(id types.RoomID, allowed []types.ClientID, ttl time.Duration) (interfaces.Room, error) {
	return i.rooms.CreateRoom(id, allowed, ttl)
}

func (i *Interceptor) DeleteRoom(id types.RoomID) error {
	return i.rooms.DeleteRoom(id)
}

func (i *Interceptor) WriteMessage(roomid types.RoomID, msg message.Message, from types.ClientID, tos ...types.ClientID) error {
	w, ok := i.rooms.(interfaces.CanWriteRoomMessage)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return w.WriteMessage(roomid, msg, from, tos...)
}

func (i *Interceptor) Add(id types.RoomID, s interfaces.State) error {
	a, ok := i.rooms.(interfaces.CanAdd)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return a.Add(id, s)
}

func (i *Interceptor) Remove(id types.RoomID, s interfaces.State) error {
	r, ok := i.rooms.(interfaces.CanRemove)
	if !ok {
		return errors.ErrInterfaceMisMatch
	}

	return r.Remove(id, s)
}
