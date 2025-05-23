package messages

import (
	"context"
	"fmt"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat"
)

const (
	IdentProtocol         message.Protocol = "room:ident"
	IdentResponseProtocol message.Protocol = "room:ident_response"
)

type Ident struct {
	interceptor.BaseMessage
}

func (m *Ident) GetProtocol() message.Protocol {
	return IdentProtocol
}

func (m *Ident) ReadProcess(ctx context.Context, _i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(*chat.ClientInterceptor)
	if !ok {
		return fmt.Errorf("error while processing 'Ident' message; err: %s", interceptor.ErrInterfaceMisMatch.Error())
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while processing 'Ident' message; err: %s", err.Error())
	}

	if err := ss.SetClientID(interceptor.ClientID(m.CurrentHeader.Sender)); err != nil {
		return fmt.Errorf("error while processing 'Ident' message; err: %s", err.Error())
	}

	if err := ss.Write(ctx, &IdentResponse{}); err != nil {
		return fmt.Errorf("error while processing 'Ident' message; err: %s", err.Error())
	}

	return nil
}

type IdentResponse struct {
	interceptor.BaseMessage
}

// TODO: ADD THE FACTORY FOR IdentResponse

func (m *IdentResponse) GetProtocol() message.Protocol {
	return IdentResponseProtocol
}

func (m *IdentResponse) ReadProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	s, ok := _i.(*chat.ServerInterceptor)
	if !ok {
		return fmt.Errorf("error while processing 'Ident' message; err: %s", interceptor.ErrInterfaceMisMatch.Error())
	}

	ss, err := s.GetState(connection)
	if err != nil {
		return fmt.Errorf("error while processing 'Ident' message; err: %s", err.Error())
	}

	if err := ss.SetClientID(interceptor.ClientID(m.CurrentHeader.Sender)); err != nil {
		return fmt.Errorf("error while processing 'Ident' message; err: %s", err.Error())
	}

	return nil
}
