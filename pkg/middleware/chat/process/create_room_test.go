package process

import (
	"context"
	"errors"
	"testing"
	"time"

	chaterrors "github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/room"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// MockCreateRoomProcessor is a mock implementation of the Processor interface with CanCreateRoom
type MockCreateRoomProcessor struct {
	createRoomFunc func(types.RoomID, []types.ClientID, time.Duration) (*room.Room, error)
}

// CreateRoom is a mock implementation of the CanCreateRoom interface
func (m *MockCreateRoomProcessor) CreateRoom(roomID types.RoomID, allowed []types.ClientID, ttl time.Duration) (*room.Room, error) {
	if m.createRoomFunc != nil {
		return m.createRoomFunc(roomID, allowed, ttl)
	}
	return &room.Room{}, nil
}

// Process implements the CanProcess interface
func (m *MockCreateRoomProcessor) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockCreateRoomProcessor) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

// MockProcessorWithoutCreateRoom is a mock implementation of the Processor interface without CanCreateRoom
type MockProcessorWithoutCreateRoom struct{}

// Process implements the CanProcess interface
func (m *MockProcessorWithoutCreateRoom) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockProcessorWithoutCreateRoom) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

func TestCreateRoom_Process(t *testing.T) {
	tests := []struct {
		name               string
		roomID             types.RoomID
		allowed            []types.ClientID
		ttl                time.Duration
		mockCreateRoomFunc func(types.RoomID, []types.ClientID, time.Duration) (*room.Room, error)
		contextTimeout     bool
		wantErr            bool
		expectedErr        error
	}{
		{
			name:    "successful create room",
			roomID:  "test-room",
			allowed: []types.ClientID{"client1", "client2"},
			ttl:     time.Hour,
			mockCreateRoomFunc: func(roomID types.RoomID, allowed []types.ClientID, ttl time.Duration) (*room.Room, error) {
				return &room.Room{}, nil
			},
			wantErr: false,
		},
		{
			name:    "error creating room",
			roomID:  "test-room",
			allowed: []types.ClientID{"client1", "client2"},
			ttl:     time.Hour,
			mockCreateRoomFunc: func(roomID types.RoomID, allowed []types.ClientID, ttl time.Duration) (*room.Room, error) {
				return nil, errors.New("failed to create room")
			},
			wantErr:     true,
			expectedErr: errors.New("failed to create room"),
		},
		{
			name:           "context cancelled",
			roomID:         "test-room",
			allowed:        []types.ClientID{"client1", "client2"},
			ttl:            time.Hour,
			contextTimeout: true,
			wantErr:        true,
			expectedErr:    chaterrors.ErrContextCancelled,
		},
		{
			name:               "interface mismatch",
			roomID:             "test-room",
			allowed:            []types.ClientID{"client1", "client2"},
			ttl:                time.Hour,
			mockCreateRoomFunc: nil, // This will cause the type assertion to fail
			wantErr:            true,
			expectedErr:        chaterrors.ErrInterfaceMisMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a context that can be cancelled if needed
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Cancel the context if the test requires it
			if tt.contextTimeout {
				cancel()
			}

			// Create a dummy state
			s := &state.State{}

			// Create the process
			process := NewCreateRoom(tt.roomID, tt.allowed, tt.ttl)

			// Use different mock processor for interface mismatch test
			var processor interfaces.Processor
			if tt.name == "interface mismatch" {
				processor = &MockProcessorWithoutCreateRoom{}
			} else {
				processor = &MockCreateRoomProcessor{
					createRoomFunc: tt.mockCreateRoomFunc,
				}
			}

			// Execute the process
			err := process.Process(ctx, processor, s)

			// Check if the error matches expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Process() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
