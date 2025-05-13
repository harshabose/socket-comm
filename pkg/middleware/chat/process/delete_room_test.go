package process

import (
	"context"
	"errors"
	"testing"

	chaterrors "github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// MockDeleteRoomProcessor is a mock implementation of the Processor interface with CanDeleteRoom
type MockDeleteRoomProcessor struct {
	deleteRoomFunc func(types.RoomID) error
}

// DeleteRoom is a mock implementation of the CanDeleteRoom interface
func (m *MockDeleteRoomProcessor) DeleteRoom(roomID types.RoomID) error {
	if m.deleteRoomFunc != nil {
		return m.deleteRoomFunc(roomID)
	}
	return nil
}

// Process implements the CanProcess interface
func (m *MockDeleteRoomProcessor) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockDeleteRoomProcessor) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

// MockProcessorWithoutDeleteRoom is a mock implementation of the Processor interface without CanDeleteRoom
type MockProcessorWithoutDeleteRoom struct{}

// Process implements the CanProcess interface
func (m *MockProcessorWithoutDeleteRoom) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockProcessorWithoutDeleteRoom) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

func TestDeleteRoom_Process(t *testing.T) {
	tests := []struct {
		name               string
		roomID             types.RoomID
		mockDeleteRoomFunc func(types.RoomID) error
		contextTimeout     bool
		wantErr            bool
		expectedErr        error
	}{
		{
			name:   "successful delete room",
			roomID: "test-room",
			mockDeleteRoomFunc: func(roomID types.RoomID) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "error deleting room",
			roomID: "test-room",
			mockDeleteRoomFunc: func(roomID types.RoomID) error {
				return errors.New("failed to delete room")
			},
			wantErr:     true,
			expectedErr: errors.New("failed to delete room"),
		},
		{
			name:           "context cancelled",
			roomID:         "test-room",
			contextTimeout: true,
			wantErr:        true,
			expectedErr:    chaterrors.ErrContextCancelled,
		},
		{
			name:               "interface mismatch",
			roomID:             "test-room",
			mockDeleteRoomFunc: nil, // This will cause the type assertion to fail
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
			process := &DeleteRoom{
				RoomID: tt.roomID,
			}

			// Use different mock processor for interface mismatch test
			var processor interfaces.Processor
			if tt.name == "interface mismatch" {
				processor = &MockProcessorWithoutDeleteRoom{}
			} else {
				processor = &MockDeleteRoomProcessor{
					deleteRoomFunc: tt.mockDeleteRoomFunc,
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
