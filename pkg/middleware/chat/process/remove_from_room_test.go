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

// MockRemoveProcessor is a mock implementation of the Processor interface with CanRemove
type MockRemoveProcessor struct {
	removeFunc func(types.RoomID, *state.State) error
}

// Remove is a mock implementation of the CanRemove interface
func (m *MockRemoveProcessor) Remove(roomID types.RoomID, s *state.State) error {
	if m.removeFunc != nil {
		return m.removeFunc(roomID, s)
	}
	return nil
}

// Process implements the CanProcess interface
func (m *MockRemoveProcessor) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockRemoveProcessor) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

// MockProcessorWithoutRemove is a mock implementation of the Processor interface without CanRemove
type MockProcessorWithoutRemove struct{}

// Process implements the CanProcess interface
func (m *MockProcessorWithoutRemove) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockProcessorWithoutRemove) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

func TestRemoveFromRoom_Process(t *testing.T) {
	tests := []struct {
		name           string
		roomID         types.RoomID
		mockRemoveFunc func(types.RoomID, *state.State) error
		wantErr        bool
		expectedErr    error
	}{
		{
			name:   "successful remove from room",
			roomID: "test-room",
			mockRemoveFunc: func(roomID types.RoomID, s *state.State) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "error removing from room",
			roomID: "test-room",
			mockRemoveFunc: func(roomID types.RoomID, s *state.State) error {
				return errors.New("failed to remove from room")
			},
			wantErr:     true,
			expectedErr: errors.New("failed to remove from room"),
		},
		{
			name:           "interface mismatch",
			roomID:         "test-room",
			mockRemoveFunc: nil, // This will cause the type assertion to fail
			wantErr:        true,
			expectedErr:    chaterrors.ErrInterfaceMisMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy state
			s := &state.State{}

			// Create the process
			process := &RemoveFromRoom{
				RoomID: tt.roomID,
			}

			// Use different mock processor for interface mismatch test
			var processor interfaces.Processor
			if tt.name == "interface mismatch" {
				processor = &MockProcessorWithoutRemove{}
			} else {
				processor = &MockRemoveProcessor{
					removeFunc: tt.mockRemoveFunc,
				}
			}

			// Execute the process
			err := process.Process(context.Background(), processor, s)

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
