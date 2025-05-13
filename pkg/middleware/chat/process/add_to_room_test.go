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

// MockProcessor is a mock implementation of the Processor interface
type MockProcessor struct {
	addFunc func(types.RoomID, *state.State) error
}

// Add is a mock implementation of the CanAdd interface
func (m *MockProcessor) Add(roomID types.RoomID, s *state.State) error {
	if m.addFunc != nil {
		return m.addFunc(roomID, s)
	}
	return nil
}

// Process implements the CanProcess interface
func (m *MockProcessor) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockProcessor) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

// MockProcessorWithoutAdd is a mock implementation of the Processor interface without CanAdd
type MockProcessorWithoutAdd struct{}

// Process implements the CanProcess interface
func (m *MockProcessorWithoutAdd) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockProcessorWithoutAdd) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

func TestAddToRoom_Process(t *testing.T) {
	tests := []struct {
		name           string
		roomID         types.RoomID
		mockAddFunc    func(types.RoomID, *state.State) error
		contextTimeout bool
		wantErr        bool
		expectedErr    error
	}{
		{
			name:   "successful add to room",
			roomID: "test-room",
			mockAddFunc: func(roomID types.RoomID, s *state.State) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "error adding to room",
			roomID: "test-room",
			mockAddFunc: func(roomID types.RoomID, s *state.State) error {
				return errors.New("failed to add to room")
			},
			wantErr:     true,
			expectedErr: errors.New("failed to add to room"),
		},
		{
			name:           "context cancelled",
			roomID:         "test-room",
			contextTimeout: true,
			wantErr:        true,
			expectedErr:    chaterrors.ErrContextCancelled,
		},
		{
			name:        "interface mismatch",
			roomID:      "test-room",
			mockAddFunc: nil, // This will cause the type assertion to fail
			wantErr:     true,
			expectedErr: chaterrors.ErrInterfaceMisMatch,
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
			process := NewAddToRoom(tt.roomID)

			// Use different mock processor for interface mismatch test
			var processor interfaces.Processor
			if tt.name == "interface mismatch" {
				processor = &MockProcessorWithoutAdd{}
			} else {
				processor = &MockProcessor{
					addFunc: tt.mockAddFunc,
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
