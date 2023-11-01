package mock

import (
	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/crypto"
	"github.com/Dharitri-org/sme-core-vm-go/wasmer"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

var _ core.VMHost = (*VmHostMock)(nil)

type VmHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext core.BlockchainContext
	RuntimeContext    core.RuntimeContext
	OutputContext     core.OutputContext
	MeteringContext   core.MeteringContext
	StorageContext    core.StorageContext
	BigIntContext     core.BigIntContext

	SCAPIMethods  *wasmer.Imports
	IsBuiltinFunc bool
}

func (host *VmHostMock) Crypto() crypto.VMCrypto {
	return host.CryptoHook
}

func (host *VmHostMock) Blockchain() core.BlockchainContext {
	return host.BlockchainContext
}

func (host *VmHostMock) Runtime() core.RuntimeContext {
	return host.RuntimeContext
}

func (host *VmHostMock) Output() core.OutputContext {
	return host.OutputContext
}

func (host *VmHostMock) Metering() core.MeteringContext {
	return host.MeteringContext
}

func (host *VmHostMock) Storage() core.StorageContext {
	return host.StorageContext
}

func (host *VmHostMock) BigInt() core.BigIntContext {
	return host.BigIntContext
}

func (host *VmHostMock) IsCoreV2Enabled() bool {
	return true
}

func (host *VmHostMock) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	return nil, nil
}

func (host *VmHostMock) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*core.AsyncContextInfo, error) {
	return nil, nil
}

func (host *VmHostMock) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *core.AsyncContextInfo, error) {
	return nil, nil, nil
}

func (host *VmHostMock) EthereumCallData() []byte {
	return host.EthInput
}

func (host *VmHostMock) InitState() {
}

func (host *VmHostMock) PushState() {
}

func (host *VmHostMock) PopState() {
}

func (host *VmHostMock) ClearStateStack() {
}

func (host *VmHostMock) GetAPIMethods() *wasmer.Imports {
	return host.SCAPIMethods
}

func (host *VmHostMock) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

func (host *VmHostMock) IsBuiltinFunctionName(functionName string) bool {
	return host.IsBuiltinFunc
}
