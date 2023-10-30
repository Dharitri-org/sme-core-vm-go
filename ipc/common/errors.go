package common

import (
	"fmt"
)

// CriticalError signals a critical error
type CriticalError struct {
	InnerErr error
}

// WrapCriticalError wraps an error
func WrapCriticalError(err error) *CriticalError {
	return &CriticalError{InnerErr: err}
}

func (err *CriticalError) Error() string {
	return fmt.Sprintf("critical error: %v", err.InnerErr)
}

// Unwrap unwraps the inner error
func (err *CriticalError) Unwrap() error {
	return err.InnerErr
}

// IsCriticalError returns whether the error is critical
func IsCriticalError(err error) bool {
	_, ok := err.(*CriticalError)
	return ok
}

// ErrBadCoreArguments signals a critical error
var ErrBadCoreArguments = &CriticalError{InnerErr: fmt.Errorf("bad arguments passed to core")}

// ErrCoreClosed signals a critical error
var ErrCoreClosed = &CriticalError{InnerErr: fmt.Errorf("core closed")}

// ErrCoreTimeExpired signals a critical error
var ErrCoreTimeExpired = &CriticalError{InnerErr: fmt.Errorf("core time expired")}

// ErrCoreNotFound signals a critical error
var ErrCoreNotFound = &CriticalError{InnerErr: fmt.Errorf("core binary not found")}

// ErrInvalidMessageNonce signals a critical error
var ErrInvalidMessageNonce = &CriticalError{InnerErr: fmt.Errorf("invalid dialogue nonce")}

// ErrStopPerNodeRequest signals a critical error
var ErrStopPerNodeRequest = &CriticalError{InnerErr: fmt.Errorf("core will stop, as requested")}

// ErrBadRequestFromNode signals a critical error
var ErrBadRequestFromNode = &CriticalError{InnerErr: fmt.Errorf("bad message from node")}

// ErrBadMessageFromCore signals a critical error
var ErrBadMessageFromCore = &CriticalError{InnerErr: fmt.Errorf("bad message from core")}

// ErrCannotSendContractRequest signals a critical error
var ErrCannotSendContractRequest = &CriticalError{InnerErr: fmt.Errorf("cannot send contract request")}

// ErrCannotSendHookCallResponse signals a critical error
var ErrCannotSendHookCallResponse = &CriticalError{InnerErr: fmt.Errorf("cannot send hook call response")}

// ErrCannotSendHookCallRequest signals a critical error
var ErrCannotSendHookCallRequest = &CriticalError{InnerErr: fmt.Errorf("cannot send hook call request")}

// ErrCannotReceiveHookCallResponse signals a critical error
var ErrCannotReceiveHookCallResponse = &CriticalError{InnerErr: fmt.Errorf("cannot receive hook call response")}

// ErrBadHookResponseFromNode signals a critical error
var ErrBadHookResponseFromNode = &CriticalError{InnerErr: fmt.Errorf("bad hook response from node")}

const (
	// ErrCodeSuccess signals success
	ErrCodeSuccess = iota
	// ErrCodeCannotCreateFile signals a critical error
	ErrCodeCannotCreateFile
	// ErrCodeInit signals a critical error
	ErrCodeInit
	// ErrCodeTerminated signals a critical error
	ErrCodeTerminated
)
