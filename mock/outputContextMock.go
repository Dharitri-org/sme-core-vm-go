package mock

import (
	"math/big"

	"github.com/Dharitri-org/sme-core-vm-go/core"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

var _ core.OutputContext = (*OutputContextMock)(nil)

type OutputContextMock struct {
	OutputStateMock    *vmcommon.VMOutput
	ReturnDataMock     [][]byte
	ReturnCodeMock     vmcommon.ReturnCode
	ReturnMessageMock  string
	GasRemaining       uint64
	GasRefund          *big.Int
	OutputAccounts     map[string]*vmcommon.OutputAccount
	DeletedAccounts    [][]byte
	TouchedAccounts    [][]byte
	Logs               []*vmcommon.LogEntry
	OutputAccountMock  *vmcommon.OutputAccount
	OutputAccountIsNew bool
	Err                error
	TransferResult     error
}

func (o *OutputContextMock) AddToActiveState(_ *vmcommon.VMOutput) {
}

func (o *OutputContextMock) InitState() {
}

func (o *OutputContextMock) NewVMOutputAccount(address []byte) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        address,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

func (o *OutputContextMock) NewVMOutputAccountFromMockAccount(account *AccountMock) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        account.Address,
		Nonce:          account.Nonce,
		BalanceDelta:   big.NewInt(0),
		Balance:        account.Balance,
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

func (o *OutputContextMock) PushState() {
}

func (o *OutputContextMock) PopSetActiveState() {
}

func (o *OutputContextMock) PopMergeActiveState() {
}

func (o *OutputContextMock) PopDiscard() {
}

func (o *OutputContextMock) ClearStateStack() {
}

func (o *OutputContextMock) CopyTopOfStackToActiveState() {
}

func (o *OutputContextMock) CensorVMOutput() {
}

func (o *OutputContextMock) ResetGas() {
}

func (o *OutputContextMock) GetOutputAccount(_ []byte) (*vmcommon.OutputAccount, bool) {
	return o.OutputAccountMock, o.OutputAccountIsNew
}

func (o *OutputContextMock) DeleteOutputAccount(_ []byte) {
}

func (o *OutputContextMock) GetRefund() uint64 {
	return uint64(o.GasRefund.Int64())
}

func (o *OutputContextMock) SetRefund(refund uint64) {
	o.GasRefund = big.NewInt(int64(refund))
}

func (o *OutputContextMock) ReturnData() [][]byte {
	return o.ReturnDataMock
}

func (o *OutputContextMock) ReturnCode() vmcommon.ReturnCode {
	return o.ReturnCodeMock
}

func (o *OutputContextMock) SetReturnCode(returnCode vmcommon.ReturnCode) {
	o.ReturnCodeMock = returnCode
}

func (o *OutputContextMock) ReturnMessage() string {
	return o.ReturnMessageMock
}

func (o *OutputContextMock) SetReturnMessage(returnMessage string) {
	o.ReturnMessageMock = returnMessage
}

func (o *OutputContextMock) ClearReturnData() {
	o.ReturnDataMock = make([][]byte, 0)
}

func (o *OutputContextMock) SelfDestruct(_ []byte, _ []byte) {
	panic("not implemented")
}

func (o *OutputContextMock) Finish(data []byte) {
	o.ReturnDataMock = append(o.ReturnDataMock, data)
}

func (o *OutputContextMock) WriteLog(_ []byte, _ [][]byte, _ []byte) {
}

func (o *OutputContextMock) TransferValueOnly(_ []byte, _ []byte, _ *big.Int) error {
	return o.TransferResult
}

func (o *OutputContextMock) Transfer(_ []byte, _ []byte, _ uint64, _ *big.Int, _ []byte, _ vmcommon.CallType) error {
	return o.TransferResult
}

func (o *OutputContextMock) AddTxValueToAccount(_ []byte, _ *big.Int) {
}

func (o *OutputContextMock) GetVMOutput() *vmcommon.VMOutput {
	return o.OutputStateMock
}

func (o *OutputContextMock) DeployCode(_ core.CodeDeployInput) {
}

func (o *OutputContextMock) CreateVMOutputInCaseOfError(_ error) *vmcommon.VMOutput {
	return o.OutputStateMock
}
