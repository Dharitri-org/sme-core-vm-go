package main

import (
	"github.com/Dharitri-org/sme-core-vm-go/coredebug"
	"github.com/urfave/cli"
)

type cliArguments struct {
	// Common arguments
	ServerAddress string
	Database      string
	World         string
	Outcome       string
	// For contract-related actions
	Impersonated    string
	ContractAddress string
	Action          string
	Function        string
	Arguments       cli.StringSlice
	Code            string
	CodePath        string
	CodeMetadata    string
	Value           string
	GasLimit        uint64
	GasPrice        uint64
	// For blockchain-related action
	AccountAddress string
	AccountBalance string
	AccountNonce   uint64
}

func (args *cliArguments) toDeployRequest() coredebug.DeployRequest {
	request := &coredebug.DeployRequest{}
	args.populateDeployRequest(request)

	return *request
}

func (args *cliArguments) populateDeployRequest(request *coredebug.DeployRequest) {
	args.populateContractRequestBase(&request.ContractRequestBase)

	request.CodeHex = args.Code
	request.CodePath = args.CodePath
	request.CodeMetadata = args.CodeMetadata
	request.ArgumentsHex = args.Arguments
}

func (args *cliArguments) populateContractRequestBase(request *coredebug.ContractRequestBase) {
	args.populateRequestBase(&request.RequestBase)

	request.ImpersonatedHex = args.Impersonated
	request.Value = args.Value
	request.GasLimit = args.GasLimit
	request.GasPrice = args.GasPrice
}

func (args *cliArguments) populateRequestBase(request *coredebug.RequestBase) {
	request.DatabasePath = args.Database
	request.World = args.World
	request.Outcome = args.Outcome
}

func (args *cliArguments) toUpgradeRequest() coredebug.UpgradeRequest {
	request := &coredebug.UpgradeRequest{}
	args.populateDeployRequest(&request.DeployRequest)

	request.ContractAddressHex = args.ContractAddress
	return *request
}

func (args *cliArguments) toRunRequest() coredebug.RunRequest {
	request := &coredebug.RunRequest{}
	args.populateRunRequest(request)

	return *request
}

func (args *cliArguments) populateRunRequest(request *coredebug.RunRequest) {
	args.populateContractRequestBase(&request.ContractRequestBase)

	request.ContractAddressHex = args.ContractAddress
	request.Function = args.Function
	request.ArgumentsHex = args.Arguments
}

func (args *cliArguments) toQueryRequest() coredebug.QueryRequest {
	request := &coredebug.QueryRequest{}
	args.populateRunRequest(&request.RunRequest)

	return *request
}

func (args *cliArguments) toCreateAccountRequest() coredebug.CreateAccountRequest {
	request := &coredebug.CreateAccountRequest{}
	args.populateRequestBase(&request.RequestBase)

	request.AddressHex = args.AccountAddress
	request.Balance = args.AccountBalance
	request.Nonce = args.AccountNonce
	return *request
}
