package mock

import (
	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/wasmer"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

var _ core.RuntimeContext = (*RuntimeContextMock)(nil)

type RuntimeContextMock struct {
	Err                    error
	VmInput                *vmcommon.VMInput
	SCAddress              []byte
	CallFunction           string
	VmType                 []byte
	ReadOnlyFlag           bool
	VerifyCode             bool
	CurrentBreakpointValue core.BreakpointValue
	PointsUsed             uint64
	InstanceCtxID          int
	MemLoadResult          []byte
	FailCryptoAPI          bool
	FailDharitriAPI        bool
	FailBigIntAPI          bool
	AsyncCallInfo          *core.AsyncCallInfo
	RunningInstances       uint64
	CurrentTxHash          []byte
	OriginalTxHash         []byte
}

func (r *RuntimeContextMock) InitState() {
}

func (r *RuntimeContextMock) StartWasmerInstance(contract []byte, gasLimit uint64) error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (r *RuntimeContextMock) InitStateFromContractCallInput(input *vmcommon.ContractCallInput) {
}

func (r *RuntimeContextMock) PushState() {
}

func (r *RuntimeContextMock) PopSetActiveState() {
}

func (r *RuntimeContextMock) PopDiscard() {
}

func (r *RuntimeContextMock) MustVerifyNextContractCode() {
}

func (r *RuntimeContextMock) ClearStateStack() {
}

func (r *RuntimeContextMock) PushInstance() {
}

func (r *RuntimeContextMock) PopInstance() {
}

func (r *RuntimeContextMock) IsWarmInstance() bool {
	return false
}

func (r *RuntimeContextMock) ResetWarmInstance() {
}

func (r *RuntimeContextMock) RunningInstancesCount() uint64 {
	return r.RunningInstances
}

func (r *RuntimeContextMock) SetMaxInstanceCount(uint64) {
}

func (r *RuntimeContextMock) ClearInstanceStack() {
}

func (r *RuntimeContextMock) GetVMType() []byte {
	return r.VmType
}

func (r *RuntimeContextMock) GetVMInput() *vmcommon.VMInput {
	return r.VmInput
}

func (r *RuntimeContextMock) SetVMInput(vmInput *vmcommon.VMInput) {
	r.VmInput = vmInput
}

func (r *RuntimeContextMock) GetSCAddress() []byte {
	return r.SCAddress
}

func (r *RuntimeContextMock) SetSCAddress(scAddress []byte) {
	r.SCAddress = scAddress
}

func (r *RuntimeContextMock) Function() string {
	return r.CallFunction
}

func (r *RuntimeContextMock) Arguments() [][]byte {
	return r.VmInput.Arguments
}

func (r *RuntimeContextMock) GetCurrentTxHash() []byte {
	return r.CurrentTxHash
}

func (r *RuntimeContextMock) GetOriginalTxHash() []byte {
	return r.OriginalTxHash
}

func (r *RuntimeContextMock) ExtractCodeUpgradeFromArgs() ([]byte, []byte, error) {
	arguments := r.VmInput.Arguments
	if len(arguments) < 2 {
		panic("ExtractCodeUpgradeFromArgs: bad test setup")
	}

	return r.VmInput.Arguments[0], r.VmInput.Arguments[1], nil
}

func (r *RuntimeContextMock) SignalExit(_ int) {
}

func (r *RuntimeContextMock) SignalUserError(_ string) {
}

func (r *RuntimeContextMock) SetRuntimeBreakpointValue(value core.BreakpointValue) {
}

func (r *RuntimeContextMock) GetRuntimeBreakpointValue() core.BreakpointValue {
	return r.CurrentBreakpointValue
}

func (r *RuntimeContextMock) VerifyContractCode() error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (r *RuntimeContextMock) GetPointsUsed() uint64 {
	return r.PointsUsed
}

func (r *RuntimeContextMock) SetPointsUsed(gasPoints uint64) {
	r.PointsUsed = gasPoints
}

func (r *RuntimeContextMock) ReadOnly() bool {
	return r.ReadOnlyFlag
}

func (r *RuntimeContextMock) SetReadOnly(readOnly bool) {
	r.ReadOnlyFlag = readOnly
}

func (r *RuntimeContextMock) GetInstanceExports() wasmer.ExportsMap {
	return nil
}

func (r *RuntimeContextMock) CleanWasmerInstance() {
}

func (r *RuntimeContextMock) GetFunctionToCall() (wasmer.ExportedFunctionCallback, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return nil, nil
}

func (r *RuntimeContextMock) GetInitFunction() wasmer.ExportedFunctionCallback {
	return nil
}

func (r *RuntimeContextMock) MemLoad(offset int32, length int32) ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return r.MemLoadResult, nil
}

func (r *RuntimeContextMock) MemStore(offset int32, data []byte) error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (r *RuntimeContextMock) DharitriAPIErrorShouldFailExecution() bool {
	return r.FailDharitriAPI
}

func (r *RuntimeContextMock) CryptoAPIErrorShouldFailExecution() bool {
	return r.FailCryptoAPI
}

func (r *RuntimeContextMock) BigIntAPIErrorShouldFailExecution() bool {
	return r.FailBigIntAPI
}

func (r *RuntimeContextMock) FailExecution(err error) {
}

func (r *RuntimeContextMock) GetAsyncCallInfo() *core.AsyncCallInfo {
	return r.AsyncCallInfo
}

func (r *RuntimeContextMock) SetAsyncCallInfo(asyncCallInfo *core.AsyncCallInfo) {
	r.AsyncCallInfo = asyncCallInfo
}

func (r *RuntimeContextMock) AddAsyncContextCall(_ []byte, _ *core.AsyncGeneratedCall) error {
	return nil
}

func (r *RuntimeContextMock) GetAsyncContextInfo() *core.AsyncContextInfo {
	return nil
}

func (r *RuntimeContextMock) GetAsyncContext(_ []byte) (*core.AsyncContext, error) {
	return nil, nil
}

func (r *RuntimeContextMock) SetCustomCallFunction(callFunction string) {

}
