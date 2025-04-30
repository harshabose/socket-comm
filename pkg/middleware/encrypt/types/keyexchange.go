package types

type SessionState int

const (
	SessionStateNotStart SessionState = iota
	SessionStateInitial
	SessionStateInProgress
	SessionStateCompleted
	SessionStateError
)
