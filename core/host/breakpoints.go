package host

import (
	"github.com/Dharitri-org/sme-core-vm-go/core"
)

func (host *vmHost) handleBreakpointIfAny(executionErr error) error {
	if executionErr == nil {
		return nil
	}

	runtime := host.Runtime()
	breakpointValue := runtime.GetRuntimeBreakpointValue()

	if breakpointValue != core.BreakpointNone {
		executionErr = host.handleBreakpoint(breakpointValue)
	}

	return executionErr
}

func (host *vmHost) handleBreakpoint(breakpointValue core.BreakpointValue) error {
	if breakpointValue == core.BreakpointAsyncCall {
		return host.handleAsyncCallBreakpoint()
	}
	if breakpointValue == core.BreakpointExecutionFailed {
		return core.ErrExecutionFailed
	}
	if breakpointValue == core.BreakpointSignalError {
		return core.ErrSignalError
	}
	if breakpointValue == core.BreakpointOutOfGas {
		return core.ErrNotEnoughGas
	}

	return core.ErrUnhandledRuntimeBreakpoint
}
