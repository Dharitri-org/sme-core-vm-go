package contexts

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/wasmer"
	logger "github.com/Dharitri-org/sme-logger"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

var log = logger.GetOrCreate("core/runtime")

var _ core.RuntimeContext = (*runtimeContext)(nil)

type runtimeContext struct {
	host         core.VMHost
	instance     *wasmer.Instance
	vmInput      *vmcommon.VMInput
	scAddress    []byte
	callFunction string
	vmType       []byte
	readOnly     bool

	verifyCode bool

	stateStack    []*runtimeContext
	instanceStack []*wasmer.Instance

	maxWasmerInstances uint64

	asyncCallInfo    *core.AsyncCallInfo
	asyncContextInfo *core.AsyncContextInfo

	validator *WASMValidator

	useWarmInstance     bool
	warmInstanceAddress []byte
	warmInstance        *wasmer.Instance
}

// NewRuntimeContext creates a new runtimeContext
func NewRuntimeContext(host core.VMHost, vmType []byte, useWarmInstance bool) (*runtimeContext, error) {
	scAPINames := host.GetAPIMethods().Names()
	protocolBuiltinFunctions := host.GetProtocolBuiltinFunctions()

	context := &runtimeContext{
		host:                host,
		vmType:              vmType,
		stateStack:          make([]*runtimeContext, 0),
		instanceStack:       make([]*wasmer.Instance, 0),
		validator:           NewWASMValidator(scAPINames, protocolBuiltinFunctions),
		useWarmInstance:     useWarmInstance,
		warmInstanceAddress: nil,
		warmInstance:        nil,
	}

	context.InitState()

	return context, nil
}

func (context *runtimeContext) InitState() {
	context.vmInput = &vmcommon.VMInput{}
	context.scAddress = make([]byte, 0)
	context.callFunction = ""
	context.verifyCode = false
	context.readOnly = false
	context.asyncCallInfo = nil
	context.asyncContextInfo = &core.AsyncContextInfo{
		AsyncContextMap: make(map[string]*core.AsyncContext),
	}
}

func (context *runtimeContext) StartWasmerInstance(contract []byte, gasLimit uint64) error {
	if context.RunningInstancesCount() >= context.maxWasmerInstances {
		context.instance = nil
		return core.ErrMaxInstancesReached
	}

	scAddress := context.GetSCAddress()
	useWarm := context.useWarmInstance && context.warmInstanceAddress != nil && bytes.Equal(scAddress, context.warmInstanceAddress)
	if scAddress != nil && useWarm {
		log.Trace("Reusing the warm Wasmer instance")

		context.instance = context.warmInstance
		context.SetPointsUsed(0)
	} else {
		log.Trace("Creating a new Wasmer instance")

		gasSchedule := context.host.Metering().GasSchedule()
		options := wasmer.CompilationOptions{
			GasLimit:           gasLimit,
			UnmeteredLocals:    uint64(gasSchedule.WASMOpcodeCost.LocalsUnmetered),
			OpcodeTrace:        false,
			Metering:           true,
			RuntimeBreakpoints: true,
		}
		newInstance, err := wasmer.NewInstanceWithOptions(contract, options)
		if err != nil {
			context.instance = nil
			return err
		}

		context.instance = newInstance
		if context.useWarmInstance {
			context.warmInstanceAddress = context.GetSCAddress()
			context.warmInstance = context.instance
		}

		idContext := core.AddHostContext(context.host)
		context.instance.SetContextData(idContext)

		err = context.VerifyContractCode()
		if err != nil {
			context.CleanWasmerInstance()
			return err
		}
	}

	context.SetRuntimeBreakpointValue(core.BreakpointNone)
	return nil
}

func (context *runtimeContext) IsWarmInstance() bool {
	if context.instance != nil && context.instance == context.warmInstance {
		return true
	}

	return false
}

func (context *runtimeContext) ResetWarmInstance() {
	if context.instance == nil {
		return
	}

	core.RemoveHostContext(*context.instance.Data)
	context.instance.Clean()

	context.instance = nil
	context.warmInstanceAddress = nil
	context.warmInstance = nil
}

func (context *runtimeContext) MustVerifyNextContractCode() {
	context.verifyCode = true
}

func (context *runtimeContext) SetMaxInstanceCount(maxInstances uint64) {
	context.maxWasmerInstances = maxInstances
}

func (context *runtimeContext) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
	context.SetVMInput(&input.VMInput)
	context.scAddress = input.RecipientAddr
	context.callFunction = input.Function
	// Reset async map for initial state
	context.asyncContextInfo = &core.AsyncContextInfo{
		CallerAddr:      input.CallerAddr,
		AsyncContextMap: make(map[string]*core.AsyncContext),
	}
}

func (context *runtimeContext) SetCustomCallFunction(callFunction string) {
	context.callFunction = callFunction
}

func (context *runtimeContext) PushState() {
	newState := &runtimeContext{
		vmInput:          context.vmInput,
		scAddress:        context.scAddress,
		callFunction:     context.callFunction,
		readOnly:         context.readOnly,
		asyncCallInfo:    context.asyncCallInfo,
		asyncContextInfo: context.asyncContextInfo,
	}

	context.stateStack = append(context.stateStack, newState)
}

func (context *runtimeContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.vmInput = prevState.vmInput
	context.scAddress = prevState.scAddress
	context.callFunction = prevState.callFunction
	context.readOnly = prevState.readOnly
	context.asyncCallInfo = prevState.asyncCallInfo
	context.asyncContextInfo = prevState.asyncContextInfo
}

func (context *runtimeContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	context.stateStack = context.stateStack[:stateStackLen-1]
}

func (context *runtimeContext) ClearStateStack() {
	context.stateStack = make([]*runtimeContext, 0)
}

func (context *runtimeContext) PushInstance() {
	context.instanceStack = append(context.instanceStack, context.instance)
}

func (context *runtimeContext) PopInstance() {
	instanceStackLen := len(context.instanceStack)
	prevInstance := context.instanceStack[instanceStackLen-1]
	context.instanceStack = context.instanceStack[:instanceStackLen-1]

	context.CleanWasmerInstance()
	context.instance = prevInstance
}

func (context *runtimeContext) RunningInstancesCount() uint64 {
	return uint64(len(context.instanceStack))
}

func (context *runtimeContext) ClearInstanceStack() {
	for _, instance := range context.instanceStack {
		instance.Clean()
	}
	context.instanceStack = make([]*wasmer.Instance, 0)
}

func (context *runtimeContext) GetVMType() []byte {
	return context.vmType
}

func (context *runtimeContext) GetVMInput() *vmcommon.VMInput {
	return context.vmInput
}

func (context *runtimeContext) SetVMInput(vmInput *vmcommon.VMInput) {
	if !context.host.IsCoreV2Enabled() || vmInput == nil {
		context.vmInput = vmInput
		return
	}

	context.vmInput = &vmcommon.VMInput{
		CallType:    vmInput.CallType,
		GasPrice:    vmInput.GasPrice,
		GasProvided: vmInput.GasProvided,
		CallValue:   big.NewInt(0),
	}

	if vmInput.CallValue != nil {
		context.vmInput.CallValue.Set(vmInput.CallValue)
	}

	if len(vmInput.CallerAddr) > 0 {
		context.vmInput.CallerAddr = make([]byte, len(vmInput.CallerAddr))
		copy(context.vmInput.CallerAddr, vmInput.CallerAddr)
	}

	if len(vmInput.OriginalTxHash) > 0 {
		context.vmInput.OriginalTxHash = make([]byte, len(vmInput.OriginalTxHash))
		copy(context.vmInput.OriginalTxHash, vmInput.OriginalTxHash)
	}

	if len(vmInput.CurrentTxHash) > 0 {
		context.vmInput.CurrentTxHash = make([]byte, len(vmInput.CurrentTxHash))
		copy(context.vmInput.CurrentTxHash, vmInput.CurrentTxHash)
	}

	if len(vmInput.Arguments) > 0 {
		context.vmInput.Arguments = make([][]byte, len(vmInput.Arguments))
		for i, arg := range vmInput.Arguments {
			context.vmInput.Arguments[i] = make([]byte, len(arg))
			copy(context.vmInput.Arguments[i], arg)
		}
	}
}

func (context *runtimeContext) GetSCAddress() []byte {
	return context.scAddress
}

func (context *runtimeContext) SetSCAddress(scAddress []byte) {
	context.scAddress = scAddress
}

func (context *runtimeContext) GetCurrentTxHash() []byte {
	return context.vmInput.CurrentTxHash
}

func (context *runtimeContext) GetOriginalTxHash() []byte {
	return context.vmInput.OriginalTxHash
}

func (context *runtimeContext) Function() string {
	return context.callFunction
}

func (context *runtimeContext) Arguments() [][]byte {
	return context.vmInput.Arguments
}

func (context *runtimeContext) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	const numMinUpgradeArguments = 2

	arguments := context.vmInput.Arguments
	if len(arguments) < numMinUpgradeArguments {
		return nil, nil, core.ErrInvalidUpgradeArguments
	}

	code := arguments[0]
	codeMetadata := arguments[1]
	context.vmInput.Arguments = context.vmInput.Arguments[numMinUpgradeArguments:]
	return code, codeMetadata, nil
}

func (context *runtimeContext) FailExecution(err error) {
	context.host.Output().SetReturnCode(vmcommon.ExecutionFailed)

	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = "execution failed"
	}

	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(core.BreakpointExecutionFailed)
}

func (context *runtimeContext) SignalUserError(message string) {
	context.host.Output().SetReturnCode(vmcommon.UserError)
	context.host.Output().SetReturnMessage(message)
	context.SetRuntimeBreakpointValue(core.BreakpointSignalError)
}

func (context *runtimeContext) SetRuntimeBreakpointValue(value core.BreakpointValue) {
	context.instance.SetBreakpointValue(uint64(value))
}

func (context *runtimeContext) GetRuntimeBreakpointValue() core.BreakpointValue {
	return core.BreakpointValue(context.instance.GetBreakpointValue())
}

func (context *runtimeContext) VerifyContractCode() error {
	if !context.verifyCode {
		return nil
	}

	context.verifyCode = false

	err := context.validator.verifyMemoryDeclaration(context.instance)
	if err != nil {
		return err
	}

	err = context.validator.verifyFunctions(context.instance)
	if err != nil {
		return err
	}

	return nil
}

func (context *runtimeContext) DharitriAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) BigIntAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) CryptoAPIErrorShouldFailExecution() bool {
	return true
}

func (context *runtimeContext) GetPointsUsed() uint64 {
	return context.instance.GetPointsUsed()
}

func (context *runtimeContext) SetPointsUsed(gasPoints uint64) {
	context.instance.SetPointsUsed(gasPoints)
}

func (context *runtimeContext) ReadOnly() bool {
	return context.readOnly
}

func (context *runtimeContext) SetReadOnly(readOnly bool) {
	context.readOnly = readOnly
}

func (context *runtimeContext) GetInstanceExports() wasmer.ExportsMap {
	return context.instance.Exports
}

func (context *runtimeContext) CleanWasmerInstance() {
	if context.instance == nil || context.IsWarmInstance() {
		return
	}

	core.RemoveHostContext(*context.instance.Data)
	context.instance.Clean()
	context.instance = nil
}

func (context *runtimeContext) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	exports := context.instance.Exports
	if function, ok := exports[context.callFunction]; ok {
		return function, nil
	}

	if context.callFunction == core.CallbackDefault {
		return nil, core.ErrNilCallbackFunction
	}

	return nil, core.ErrFuncNotFound
}

func (context *runtimeContext) GetInitFunction() wasmer.ExportedFunctionCallback {
	exports := context.instance.Exports
	if init, ok := exports[core.InitFunctionName]; ok {
		return init
	}

	return nil
}

func (context *runtimeContext) SetAsyncCallInfo(asyncCallInfo *core.AsyncCallInfo) {
	context.asyncCallInfo = asyncCallInfo
}

func (context *runtimeContext) AddAsyncContextCall(contextIdentifier []byte, asyncCall *core.AsyncGeneratedCall) error {
	_, ok := context.asyncContextInfo.AsyncContextMap[string(contextIdentifier)]
	currentContextMap := context.asyncContextInfo.AsyncContextMap
	if !ok {
		currentContextMap[string(contextIdentifier)] = &core.AsyncContext{
			AsyncCalls: make([]*core.AsyncGeneratedCall, 0),
		}
	}

	currentContextMap[string(contextIdentifier)].AsyncCalls =
		append(currentContextMap[string(contextIdentifier)].AsyncCalls, asyncCall)

	return nil
}

func (context *runtimeContext) GetAsyncContextInfo() *core.AsyncContextInfo {
	return context.asyncContextInfo
}

func (context *runtimeContext) GetAsyncContext(contextIdentifier []byte) (*core.AsyncContext, error) {
	asyncContext, ok := context.asyncContextInfo.AsyncContextMap[string(contextIdentifier)]
	if !ok {
		return nil, core.ErrAsyncContextDoesNotExist
	}

	return asyncContext, nil
}

func (context *runtimeContext) GetAsyncCallInfo() *core.AsyncCallInfo {
	return context.asyncCallInfo
}

func (context *runtimeContext) MemLoad(offset int32, length int32) ([]byte, error) {
	if length == 0 {
		return []byte{}, nil
	}

	memory := context.instance.InstanceCtx.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + length)
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > memoryLength
	isRequestedEndTooLarge := requestedEnd > memoryLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge {
		return nil, fmt.Errorf("mem load: %w", core.ErrBadBounds)
	}
	if isLengthNegative {
		return nil, fmt.Errorf("mem load: %w", core.ErrNegativeLength)
	}

	result := make([]byte, length)
	if isRequestedEndTooLarge {
		copy(result, memoryView[offset:])
	} else {
		copy(result, memoryView[offset:requestedEnd])
	}

	return result, nil
}

func (context *runtimeContext) MemStore(offset int32, data []byte) error {
	dataLength := int32(len(data))
	if dataLength == 0 {
		return nil
	}

	memory := context.instance.InstanceCtx.Memory()
	memoryView := memory.Data()
	memoryLength := memory.Length()
	requestedEnd := uint32(offset + dataLength)
	isOffsetTooSmall := offset < 0
	isNewPageNecessary := requestedEnd > memoryLength

	if isOffsetTooSmall {
		return core.ErrBadLowerBounds
	}
	if isNewPageNecessary {
		err := memory.Grow(1)
		if err != nil {
			return err
		}

		memoryView = memory.Data()
		memoryLength = memory.Length()
	}

	isRequestedEndTooLarge := requestedEnd > memoryLength
	if isRequestedEndTooLarge {
		return core.ErrBadUpperBounds
	}

	copy(memoryView[offset:requestedEnd], data)
	return nil
}
