package process

import (
	"context"
	"errors"
	"testing"
	"time"

	chaterrors "github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/health"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

// MockUpdateProcessor is a mock implementation of the Processor interface with CanUpdate
type MockUpdateProcessor struct {
	updateFunc func(types.RoomID, types.ClientID, *health.Stat) error
}

// Update is a mock implementation of the CanUpdate interface
func (m *MockUpdateProcessor) Update(roomID types.RoomID, clientID types.ClientID, stat *health.Stat) error {
	if m.updateFunc != nil {
		return m.updateFunc(roomID, clientID, stat)
	}
	return nil
}

// Process implements the CanProcess interface
func (m *MockUpdateProcessor) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockUpdateProcessor) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return nil
}

// MockProcessorWithoutUpdate is a mock implementation of the Processor interface without CanUpdate
type MockProcessorWithoutUpdate struct{}

// Process implements the CanProcess interface
func (m *MockProcessorWithoutUpdate) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// Call the Process method on the CanBeProcessed interface
	return p.Process(ctx, m, s)
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockProcessorWithoutUpdate) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// Return the CanBeProcessedBackground interface
	return p
}

// MockBackgroundProcessorWithoutUpdate is a mock implementation of the Processor interface without CanUpdate
// specifically for background processing tests
type MockBackgroundProcessorWithoutUpdate struct{}

// Process implements the CanProcess interface
func (m *MockBackgroundProcessorWithoutUpdate) Process(ctx context.Context, p interfaces.CanBeProcessed, s *state.State) error {
	// This method should not be called in our tests
	return errors.New("Process method should not be called in tests")
}

// ProcessBackground implements the CanProcessBackground interface
func (m *MockBackgroundProcessorWithoutUpdate) ProcessBackground(ctx context.Context, p interfaces.CanBeProcessedBackground, s *state.State) interfaces.CanBeProcessedBackground {
	// This method should not be called in our tests
	return p
}

func TestUpdateHealthStat_Process(t *testing.T) {
	tests := []struct {
		name           string
		roomID         types.RoomID
		stat           health.Stat
		mockUpdateFunc func(types.RoomID, types.ClientID, *health.Stat) error
		contextTimeout bool
		wantErr        bool
		expectedErr    error
	}{
		{
			name:   "successful update health",
			roomID: "test-room",
			stat:   health.Stat{},
			mockUpdateFunc: func(roomID types.RoomID, clientID types.ClientID, stat *health.Stat) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "error updating health",
			roomID: "test-room",
			stat:   health.Stat{},
			mockUpdateFunc: func(roomID types.RoomID, clientID types.ClientID, stat *health.Stat) error {
				return errors.New("failed to update health")
			},
			wantErr:     true,
			expectedErr: errors.New("failed to update health"),
		},
		{
			name:           "context cancelled",
			roomID:         "test-room",
			stat:           health.Stat{},
			contextTimeout: true,
			wantErr:        true,
			expectedErr:    chaterrors.ErrContextCancelled,
		},
		{
			name:           "interface mismatch",
			roomID:         "test-room",
			stat:           health.Stat{},
			mockUpdateFunc: nil, // This will cause the type assertion to fail
			wantErr:        true,
			expectedErr:    chaterrors.ErrInterfaceMisMatch,
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

			// Create a dummy state with a client id
			s := &state.State{}
			s.SetClientID("test-client")

			// Create the process
			process := &UpdateHealthStat{
				RoomID: tt.roomID,
				Stat:   tt.stat,
			}

			// Use different mock processor for interface mismatch test
			var processor interfaces.Processor
			if tt.name == "interface mismatch" {
				processor = &MockProcessorWithoutUpdate{}
			} else {
				processor = &MockUpdateProcessor{
					updateFunc: tt.mockUpdateFunc,
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

func TestUpdateHealthStat_ProcessBackground(t *testing.T) {
	tests := []struct {
		name           string
		roomID         types.RoomID
		stat           health.Stat
		mockUpdateFunc func(types.RoomID, types.ClientID, *health.Stat) error
		contextTimeout bool
		wantErr        bool
		expectedErr    error
	}{
		{
			name:   "successful background update health",
			roomID: "test-room",
			stat:   health.Stat{},
			mockUpdateFunc: func(roomID types.RoomID, clientID types.ClientID, stat *health.Stat) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:   "error updating health in background",
			roomID: "test-room",
			stat:   health.Stat{},
			mockUpdateFunc: func(roomID types.RoomID, clientID types.ClientID, stat *health.Stat) error {
				return errors.New("failed to update health")
			},
			wantErr:     true,
			expectedErr: errors.New("failed to update health"),
		},
		{
			name:           "context cancelled in background",
			roomID:         "test-room",
			stat:           health.Stat{},
			contextTimeout: true,
			wantErr:        true,
			expectedErr:    chaterrors.ErrContextCancelled,
		},
		{
			name:           "interface mismatch in background",
			roomID:         "test-room",
			stat:           health.Stat{},
			mockUpdateFunc: nil, // This will cause the type assertion to fail
			wantErr:        true,
			expectedErr:    chaterrors.ErrInterfaceMisMatch,
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

			// Create a dummy state with a client id
			s := &state.State{}
			s.SetClientID("test-client")

			// Create the process
			process := &UpdateHealthStat{
				RoomID: tt.roomID,
				Stat:   tt.stat,
			}

			process.AsyncProcess = AsyncProcess{
				CanBeProcessed: process,
			}

			// Use different mock processor for interface mismatch test
			var processor interfaces.Processor
			if tt.name == "interface mismatch in background" {
				processor = &MockBackgroundProcessorWithoutUpdate{}
			} else {
				processor = &MockUpdateProcessor{
					updateFunc: tt.mockUpdateFunc,
				}
			}

			// Execute the process in background
			backgroundProcess := process.ProcessBackground(ctx, processor, s)

			// Wait for the background process to complete
			err := backgroundProcess.Wait()

			// Check if the error matches expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessBackground() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("ProcessBackground() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestUpdateHealthStat_ProcessBackground_Stop(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a dummy state with a client id
	s := &state.State{}
	s.SetClientID("test-client")

	// Create a mock processor that sleeps to simulate a long-running process
	processor := &MockUpdateProcessor{
		updateFunc: func(roomID types.RoomID, clientID types.ClientID, stat *health.Stat) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		},
	}

	// Create the process
	process := &UpdateHealthStat{
		RoomID: "test-room",
		Stat:   health.Stat{},
	}

	process.AsyncProcess = AsyncProcess{
		CanBeProcessed: process,
	}

	// Execute the process in background
	backgroundProcess := process.ProcessBackground(ctx, processor, s)

	// Stop the background process immediately
	backgroundProcess.Stop()

	// Wait for the background process to complete
	err := backgroundProcess.Wait()

	// The process should have been cancelled
	if err != chaterrors.ErrContextCancelled {
		t.Errorf("ProcessBackground() after Stop() error = %v, want %v", err, chaterrors.ErrContextCancelled)
	}
}
