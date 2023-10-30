package dharitriapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern void getSCAddress(void *context, int32_t resultOffset);
// extern void getOwnerAddress(void *context, int32_t resultOffset);
// extern int32_t getShardOfAddress(void *context, int32_t addressOffset);
// extern int32_t isSmartContract(void *context, int32_t addressOffset);
// extern void getExternalBalance(void *context, int32_t addressOffset, int32_t resultOffset);
// extern int32_t blockHash(void *context, long long nonce, int32_t resultOffset);
// extern int32_t transferValue(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern int32_t getArgumentLength(void *context, int32_t id);
// extern int32_t getArgument(void *context, int32_t id, int32_t argOffset);
// extern int32_t getFunction(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern int32_t storageStore(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoadLength(void *context, int32_t keyOffset, int32_t keyLength );
// extern int32_t storageLoad(void *context, int32_t keyOffset, int32_t keyLength , int32_t dataOffset);
// extern void getCaller(void *context, int32_t resultOffset);
// extern int32_t callValue(void *context, int32_t resultOffset);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
// extern void returnData(void* context, int32_t dataOffset, int32_t length);
// extern void signalError(void* context, int32_t messageOffset, int32_t messageLength);
// extern long long getGasLeft(void *context);
//
// extern int32_t executeOnDestContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeOnSameContext(void *context, long long gas, int32_t addressOffset, int32_t valueOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t delegateExecution(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t executeReadOnly(void *context, long long gas, int32_t addressOffset, int32_t functionOffset, int32_t functionLength, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern int32_t createContract(void *context, int32_t valueOffset, int32_t codeOffset, int32_t length, int32_t resultOffset, int32_t numArguments, int32_t argumentsLengthOffset, int32_t dataOffset);
// extern void asyncCall(void *context, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length);
// extern void createAsyncCall(void *context, int32_t identifierOffset, int32_t identifierLength, int32_t dstOffset, int32_t valueOffset, int32_t dataOffset, int32_t length, int32_t successCallback, int32_t successLength, int32_t errorCallback, int32_t errorLength, long long gas);
// extern int32_t setAsyncContextCallback(void *context, int32_t identifierOffset, int32_t identifierLength, int32_t callback, int32_t callbackLength);
//
// extern int32_t getNumReturnData(void *context);
// extern int32_t getReturnDataSize(void *context, int32_t resultID);
// extern int32_t getReturnData(void *context, int32_t resultID, int32_t dataOffset);
//
// extern int32_t setStorageLock(void *context, int32_t keyOffset, int32_t keyLength, long long lockTimestamp);
// extern long long getStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t isStorageLocked(void *context, int32_t keyOffset, int32_t keyLength);
// extern int32_t clearStorageLock(void *context, int32_t keyOffset, int32_t keyLength);
//
// extern long long getBlockTimestamp(void *context);
// extern long long getBlockNonce(void *context);
// extern long long getBlockRound(void *context);
// extern long long getBlockEpoch(void *context);
// extern void getBlockRandomSeed(void *context, int32_t resultOffset);
// extern void getStateRootHash(void *context, int32_t resultOffset);
//
// extern long long getPrevBlockTimestamp(void *context);
// extern long long getPrevBlockNonce(void *context);
// extern long long getPrevBlockRound(void *context);
// extern long long getPrevBlockEpoch(void *context);
// extern void getPrevBlockRandomSeed(void *context, int32_t resultOffset);
// extern void getOriginalTxHash(void *context, int32_t resultOffset);
//
// extern long long int64getArgument(void *context, int32_t id);
// extern int32_t int64storageStore(void *context, int32_t keyOffset, int32_t keyLength , long long value);
// extern long long int64storageLoad(void *context, int32_t keyOffset, int32_t keyLength );
// extern void int64finish(void* context, long long value);
import "C"

import (
	"math/big"
	"unsafe"

	twos "github.com/Dharitri-org/sme-components-big-int/twos-complement"
	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/wasmer"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

// DharitriEiImports creates a new wasmer.Imports populated with the DharitriEi API methods
func DharitriEiImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()
	imports = imports.Namespace("env")

	imports, err := imports.Append("getSCAddress", getSCAddress, C.getSCAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getOwnerAddress", getOwnerAddress, C.getOwnerAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getShardOfAddress", getShardOfAddress, C.getShardOfAddress)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isSmartContract", isSmartContract, C.isSmartContract)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getExternalBalance", getExternalBalance, C.getExternalBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockHash", blockHash, C.blockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("transferValue", transferValue, C.transferValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("asyncCall", asyncCall, C.asyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("createAsyncCall", createAsyncCall, C.createAsyncCall)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("setAsyncContextCallback", setAsyncContextCallback, C.setAsyncContextCallback)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgumentLength", getArgumentLength, C.getArgumentLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgument", getArgument, C.getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getFunction", getFunction, C.getFunction)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumArguments", getNumArguments, C.getNumArguments)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStore", storageStore, C.storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadLength", storageLoadLength, C.storageLoadLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoad", storageLoad, C.storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStorageLock", getStorageLock, C.getStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("setStorageLock", setStorageLock, C.setStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("isStorageLocked", isStorageLocked, C.isStorageLocked)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("clearStorageLock", clearStorageLock, C.clearStorageLock)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCaller", getCaller, C.getCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getCallValue", callValue, C.callValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", writeLog, C.writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("finish", returnData, C.returnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", signalError, C.signalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockNonce", getBlockNonce, C.getBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRound", getBlockRound, C.getBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockEpoch", getBlockEpoch, C.getBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockRandomSeed", getBlockRandomSeed, C.getBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getStateRootHash", getStateRootHash, C.getStateRootHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockTimestamp", getPrevBlockTimestamp, C.getPrevBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockNonce", getPrevBlockNonce, C.getPrevBlockNonce)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRound", getPrevBlockRound, C.getPrevBlockRound)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockEpoch", getPrevBlockEpoch, C.getPrevBlockEpoch)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getPrevBlockRandomSeed", getPrevBlockRandomSeed, C.getPrevBlockRandomSeed)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getGasLeft", getGasLeft, C.getGasLeft)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnDestContext", executeOnDestContext, C.executeOnDestContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("executeOnSameContext", executeOnSameContext, C.executeOnSameContext)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("delegateExecution", delegateExecution, C.delegateExecution)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("createContract", createContract, C.createContract)
	if err != nil {
		return nil, err
	}

	// TODO: Add extra function, upgradeContract()

	imports, err = imports.Append("executeReadOnly", executeReadOnly, C.executeReadOnly)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumReturnData", getNumReturnData, C.getNumReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnDataSize", getReturnDataSize, C.getReturnDataSize)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getReturnData", getReturnData, C.getReturnData)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64getArgument", int64getArgument, C.int64getArgument)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageStore", int64storageStore, C.int64storageStore)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64storageLoad", int64storageLoad, C.int64storageLoad)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("int64finish", int64finish, C.int64finish)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export getGasLeft
func getGasLeft(context unsafe.Pointer) int64 {
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetGasLeft
	metering.UseGas(gasToUse)

	return int64(metering.GasLeft())
}

//export getSCAddress
func getSCAddress(context unsafe.Pointer, resultOffset int32) {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	owner := runtime.GetSCAddress()
	err := runtime.MemStore(resultOffset, owner)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.GetSCAddress
	metering.UseGas(gasToUse)
}

//export getOwnerAddress
func getOwnerAddress(context unsafe.Pointer, resultOffset int32) {
	blockchain := core.GetBlockchainContext(context)
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	owner, err := blockchain.GetOwnerAddress()
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	err = runtime.MemStore(resultOffset, owner)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.GetOwnerAddress
	metering.UseGas(gasToUse)
}

//export getShardOfAddress
func getShardOfAddress(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := core.GetBlockchainContext(context)
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	address, err := runtime.MemLoad(addressOffset, core.AddressLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 0
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.GetShardOfAddress
	metering.UseGas(gasToUse)

	return int32(blockchain.GetShardOfAddress(address))
}

//export isSmartContract
func isSmartContract(context unsafe.Pointer, addressOffset int32) int32 {
	blockchain := core.GetBlockchainContext(context)
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	address, err := runtime.MemLoad(addressOffset, core.AddressLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 0
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.IsSmartContract
	metering.UseGas(gasToUse)

	isSmartContract := blockchain.IsSmartContract(address)
	return int32(core.BooleanToInt(isSmartContract))
}

//export signalError
func signalError(context unsafe.Pointer, messageOffset int32, messageLength int32) {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}
	runtime.SignalUserError(string(message))

	gasToUse := metering.GasSchedule().DharitriAPICost.SignalError
	metering.UseGas(gasToUse)
}

//export getExternalBalance
func getExternalBalance(context unsafe.Pointer, addressOffset int32, resultOffset int32) {
	blockchain := core.GetBlockchainContext(context)
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	address, err := runtime.MemLoad(addressOffset, core.AddressLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	balance := blockchain.GetBalance(address)

	err = runtime.MemStore(resultOffset, balance)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.GetExternalBalance
	metering.UseGas(gasToUse)
}

//export blockHash
func blockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	blockchain := core.GetBlockchainContext(context)
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	hash := blockchain.BlockHash(nonce)
	err := runtime.MemStore(resultOffset, hash)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export transferValue
func transferValue(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) int32 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)
	output := core.GetOutputContext(context)

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(destOffset, core.AddressLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, core.BalanceLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.TransferValue
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(gasToUse)

	err = output.Transfer(dest, send, 0, big.NewInt(0).SetBytes(value), data)
	if err != nil {
		return 1
	}

	return 0
}

//export createAsyncCall
func createAsyncCall(context unsafe.Pointer,
	asyncContextIdentifier int32,
	identifierLength int32,
	destOffset int32,
	valueOffset int32,
	dataOffset int32,
	length int32,
	successOffset int32,
	successLength int32,
	errorOffset int32,
	errorLength int32,
	gas int64,
) {
	runtime := core.GetRuntimeContext(context)

	acIdentifier, err := runtime.MemLoad(asyncContextIdentifier, identifierLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	calledSCAddress, err := runtime.MemLoad(destOffset, core.AddressLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, core.BalanceLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	successFunc, err := runtime.MemLoad(successOffset, successLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	errorFunc, err := runtime.MemLoad(errorOffset, errorLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	err = runtime.AddAsyncContextCall(acIdentifier, &core.AsyncGeneratedCall{
		Destination:     calledSCAddress,
		Data:            data,
		ValueBytes:      value,
		SuccessCallback: string(successFunc),
		ErrorCallback:   string(errorFunc),
		ProvidedGas:     uint64(gas),
	})
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

}

//export setAsyncContextCallback
func setAsyncContextCallback(context unsafe.Pointer,
	asyncContextIdentifier int32,
	identifierLength int32,
	callback int32,
	callbackLength int32,
) int32 {
	runtime := core.GetRuntimeContext(context)

	acIdentifier, err := runtime.MemLoad(asyncContextIdentifier, identifierLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	asyncContext, err := runtime.GetAsyncContext(acIdentifier)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	callbackFunc, err := runtime.MemLoad(callback, callbackLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	asyncContext.Callback = string(callbackFunc)

	return 0
}

//export asyncCall
func asyncCall(context unsafe.Pointer, destOffset int32, valueOffset int32, dataOffset int32, length int32) {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	calledSCAddress, err := runtime.MemLoad(destOffset, core.AddressLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	value, err := runtime.MemLoad(valueOffset, core.BalanceLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	data, err := runtime.MemLoad(dataOffset, length)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.DharitriAPICost.AsyncCallStep
	gasToUse += gasSchedule.BaseOperationCost.DataCopyPerByte * uint64(length)
	metering.UseGas(gasToUse)

	gasLimit := metering.GasLeft()

	minAsyncCallCost := 2*gasSchedule.DharitriAPICost.AsyncCallStep + gasSchedule.DharitriAPICost.AsyncCallbackGasLock
	if gasLimit < minAsyncCallCost {
		runtime.SetRuntimeBreakpointValue(core.BreakpointOutOfGas)
		return
	}

	// Set up the async call as if it is not known whether the called SC
	// is in the same shard with the caller or not. This will be later resolved
	// in the handler for BreakpointAsyncCall.
	runtime.SetAsyncCallInfo(&core.AsyncCallInfo{
		Destination: calledSCAddress,
		Data:        data,
		GasLimit:    gasLimit,
		ValueBytes:  value,
	})

	// Instruct Wasmer to interrupt the execution of the caller SC.
	runtime.SetRuntimeBreakpointValue(core.BreakpointAsyncCall)
}

//export getArgumentLength
func getArgumentLength(context unsafe.Pointer, id int32) int32 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	return int32(len(args[id]))
}

//export getArgument
func getArgument(context unsafe.Pointer, id int32, argOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || int32(len(args)) <= id {
		return -1
	}

	err := runtime.MemStore(argOffset, args[id])
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(args[id]))
}

//export getFunction
func getFunction(context unsafe.Pointer, functionOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetFunction
	metering.UseGas(gasToUse)

	function := runtime.Function()
	err := runtime.MemStore(functionOffset, []byte(function))
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(function))
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetNumArguments
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	return int32(len(args))
}

//export storageStore
func storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32, dataLength int32) int32 {
	runtime := core.GetRuntimeContext(context)
	storage := core.GetStorageContext(context)
	metering := core.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	data, err := runtime.MemLoad(dataOffset, dataLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.StorageStore
	metering.UseGas(gasToUse)

	storageStatus, err := storage.SetStorage(key, data)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export storageLoadLength
func storageLoadLength(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	runtime := core.GetRuntimeContext(context)
	storage := core.GetStorageContext(context)
	metering := core.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorage(key)

	gasToUse := metering.GasSchedule().DharitriAPICost.StorageLoad
	metering.UseGas(gasToUse)

	return int32(len(data))
}

//export storageLoad
func storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32, dataOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	storage := core.GetStorageContext(context)
	metering := core.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	data := storage.GetStorage(key)

	gasToUse := metering.GasSchedule().DharitriAPICost.StorageLoad
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(data))
	metering.UseGas(gasToUse)

	err = runtime.MemStore(dataOffset, data)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(data))
}

//export setStorageLock
func setStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32, lockTimestamp int64) int32 {
	runtime := core.GetRuntimeContext(context)
	storage := core.GetStorageContext(context)
	metering := core.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	timeLockKey := core.CustomStorageKey(core.TimeLockKeyPrefix, key)
	gasToUse := metering.GasSchedule().DharitriAPICost.Int64StorageStore
	metering.UseGas(gasToUse)

	bigTimestamp := big.NewInt(0).SetInt64(lockTimestamp)
	storageStatus, err := storage.SetStorage(timeLockKey, bigTimestamp.Bytes())
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}
	return int32(storageStatus)
}

//export getStorageLock
func getStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)
	storage := core.GetStorageContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	timeLockKey := core.CustomStorageKey(core.TimeLockKeyPrefix, key)
	gasToUse := metering.GasSchedule().DharitriAPICost.StorageLoad
	metering.UseGas(gasToUse)

	data := storage.GetStorage(timeLockKey)
	timeLock := big.NewInt(0).SetBytes(data).Int64()

	return timeLock
}

//export isStorageLocked
func isStorageLocked(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	timeLock := getStorageLock(context, keyOffset, keyLength)
	if timeLock < 0 {
		return -1
	}

	currentTimestamp := getBlockTimestamp(context)
	if timeLock <= currentTimestamp {
		return 0
	}

	return 1
}

//export clearStorageLock
func clearStorageLock(context unsafe.Pointer, keyOffset int32, keyLength int32) int32 {
	return setStorageLock(context, keyOffset, keyLength, 0)
}

//export getCaller
func getCaller(context unsafe.Pointer, resultOffset int32) {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	caller := runtime.GetVMInput().CallerAddr

	err := runtime.MemStore(resultOffset, caller)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.GetCaller
	metering.UseGas(gasToUse)
}

//export callValue
func callValue(context unsafe.Pointer, resultOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	value := runtime.GetVMInput().CallValue.Bytes()
	value = core.PadBytesLeft(value, core.BalanceLen)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetCallValue
	metering.UseGas(gasToUse)

	err := runtime.MemStore(resultOffset, value)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(len(value))
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	runtime := core.GetRuntimeContext(context)
	output := core.GetOutputContext(context)
	metering := core.GetMeteringContext(context)

	log, err := runtime.MemLoad(pointer, length)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	topics, err := core.GuardedMakeByteSlice2D(numTopics)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	for i := int32(0); i < numTopics; i++ {
		topics[i], err = runtime.MemLoad(topicPtr+i*core.HashLen, core.HashLen)
		if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
			return
		}
	}

	output.WriteLog(runtime.GetSCAddress(), topics, log)

	gasToUse := metering.GasSchedule().DharitriAPICost.Log
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(numTopics*core.HashLen+length)
	metering.UseGas(gasToUse)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentTimeStamp())
}

//export getBlockNonce
func getBlockNonce(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentNonce())
}

//export getBlockRound
func getBlockRound(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentRound())
}

//export getBlockEpoch
func getBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.CurrentEpoch())
}

//export getBlockRandomSeed
func getBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := core.GetRuntimeContext(context)
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.CurrentRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution())
}

//export getStateRootHash
func getStateRootHash(context unsafe.Pointer, pointer int32) {
	runtime := core.GetRuntimeContext(context)
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetStateRootHash
	metering.UseGas(gasToUse)

	stateRootHash := blockchain.GetStateRootHash()
	err := runtime.MemStore(pointer, stateRootHash)
	core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution())
}

//export getPrevBlockTimestamp
func getPrevBlockTimestamp(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockTimeStamp
	metering.UseGas(gasToUse)

	return int64(blockchain.LastTimeStamp())
}

//export getPrevBlockNonce
func getPrevBlockNonce(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockNonce
	metering.UseGas(gasToUse)

	return int64(blockchain.LastNonce())
}

//export getPrevBlockRound
func getPrevBlockRound(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockRound
	metering.UseGas(gasToUse)

	return int64(blockchain.LastRound())
}

//export getPrevBlockEpoch
func getPrevBlockEpoch(context unsafe.Pointer) int64 {
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockEpoch
	metering.UseGas(gasToUse)

	return int64(blockchain.LastEpoch())
}

//export getPrevBlockRandomSeed
func getPrevBlockRandomSeed(context unsafe.Pointer, pointer int32) {
	runtime := core.GetRuntimeContext(context)
	blockchain := core.GetBlockchainContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockRandomSeed
	metering.UseGas(gasToUse)

	randomSeed := blockchain.LastRandomSeed()
	err := runtime.MemStore(pointer, randomSeed)
	core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution())
}

//export returnData
func returnData(context unsafe.Pointer, pointer int32, length int32) {
	runtime := core.GetRuntimeContext(context)
	output := core.GetOutputContext(context)
	metering := core.GetMeteringContext(context)

	data, err := runtime.MemLoad(pointer, length)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return
	}

	output.Finish(data)
	gasToUse := metering.GasSchedule().DharitriAPICost.Finish
	gasToUse += metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
	metering.UseGas(gasToUse)
}

//export int64getArgument
func int64getArgument(context unsafe.Pointer, id int32) int64 {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.Int64GetArgument
	metering.UseGas(gasToUse)

	args := runtime.Arguments()
	if id < 0 || id >= int32(len(args)) {
		core.WithFault(core.ErrArgIndexOutOfRange, context, runtime.DharitriAPIErrorShouldFailExecution())
		return 0
	}

	arg := args[id]
	argBigInt := twos.SetBytes(big.NewInt(0), arg)
	if !argBigInt.IsInt64() {
		core.WithFault(core.ErrArgOutOfRange, context, runtime.DharitriAPIErrorShouldFailExecution())
		return 0
	}
	return argBigInt.Int64()
}

//export int64storageStore
func int64storageStore(context unsafe.Pointer, keyOffset int32, keyLength int32, value int64) int32 {
	runtime := core.GetRuntimeContext(context)
	storage := core.GetStorageContext(context)
	metering := core.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	data := big.NewInt(value)

	gasToUse := metering.GasSchedule().DharitriAPICost.Int64StorageStore
	metering.UseGas(gasToUse)

	storageStatus, err := storage.SetStorage(key, data.Bytes())
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return -1
	}

	return int32(storageStatus)
}

//export int64storageLoad
func int64storageLoad(context unsafe.Pointer, keyOffset int32, keyLength int32) int64 {
	runtime := core.GetRuntimeContext(context)
	storage := core.GetStorageContext(context)
	metering := core.GetMeteringContext(context)

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 0
	}

	data := storage.GetStorage(key)

	bigInt := big.NewInt(0).SetBytes(data)

	gasToUse := metering.GasSchedule().DharitriAPICost.Int64StorageLoad
	metering.UseGas(gasToUse)

	return bigInt.Int64()
}

//export int64finish
func int64finish(context unsafe.Pointer, value int64) {
	output := core.GetOutputContext(context)
	metering := core.GetMeteringContext(context)

	valueBytes := twos.ToBytes(big.NewInt(0).SetInt64(value))
	output.Finish(valueBytes)

	gasToUse := metering.GasSchedule().DharitriAPICost.Int64Finish
	metering.UseGas(gasToUse)
}

//export executeOnSameContext
func executeOnSameContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := core.GetVmContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(addressOffset, core.AddressLen)
	if core.WithFault(err, context, false) {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, core.BalanceLen)
	if core.WithFault(err, context, false) {
		return 1
	}

	function, data, actualLen, err := getArgumentsFromMemory(
		context,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if core.WithFault(err, context, false) {
		return 1
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.ExecuteOnSameContext
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   data,
			CallValue:   bigIntVal,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      function,
	}

	_, err = host.ExecuteOnSameContext(contractCallInput)
	if core.WithFault(err, context, false) {
		return 1
	}

	return 0
}

//export executeOnDestContext
func executeOnDestContext(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	valueOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := core.GetVmContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	send := runtime.GetSCAddress()
	dest, err := runtime.MemLoad(addressOffset, core.AddressLen)
	if core.WithFault(err, context, false) {
		return 1
	}

	value, err := runtime.MemLoad(valueOffset, core.BalanceLen)
	if core.WithFault(err, context, false) {
		return 1
	}

	function, data, actualLen, err := getArgumentsFromMemory(
		context,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if core.WithFault(err, context, false) {
		return 1
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.ExecuteOnDestContext
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	bigIntVal := big.NewInt(0).SetBytes(value)
	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  send,
			Arguments:   data,
			CallValue:   bigIntVal,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: dest,
		Function:      function,
	}

	_, _, err = host.ExecuteOnDestContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

func getArgumentsFromMemory(
	context unsafe.Pointer,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) (string, [][]byte, int32, error) {
	runtime := core.GetRuntimeContext(context)

	function, err := runtime.MemLoad(functionOffset, functionLength)
	if err != nil {
		return "", nil, 0, err
	}

	argumentsLengthData, err := runtime.MemLoad(argumentsLengthOffset, numArguments*4)
	if err != nil {
		return "", nil, 0, err
	}

	currOffset := dataOffset
	data, err := core.GuardedMakeByteSlice2D(numArguments)
	if err != nil {
		return "", nil, 0, err
	}

	for i := int32(0); i < numArguments; i++ {
		currArgLenData := argumentsLengthData[i*4 : i*4+4]
		actualLen := bytesToInt32(currArgLenData)

		data[i], err = runtime.MemLoad(currOffset, actualLen)
		if err != nil {
			return "", nil, 0, err
		}

		currOffset += actualLen
	}

	return string(function), data, currOffset - dataOffset, nil
}

//export delegateExecution
func delegateExecution(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := core.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	address, err := runtime.MemLoad(addressOffset, core.HashLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	function, data, actualLen, err := getArgumentsFromMemory(
		context,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	value := runtime.GetVMInput().CallValue
	sender := runtime.GetVMInput().CallerAddr

	gasToUse := metering.GasSchedule().DharitriAPICost.DelegateExecution
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	err = output.Transfer(address, sender, 0, value, nil)
	if err != nil {
		return 1
	}

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      function,
	}

	_, err = host.ExecuteOnSameContext(contractCallInput)
	if err != nil {
		return 1
	}

	return 0
}

func bytesToInt32(data []byte) int32 {
	actualLen := int32(0)
	for i := len(data) - 1; i >= 0; i-- {
		actualLen = (actualLen << 8) + int32(data[i])
	}

	return actualLen
}

//export executeReadOnly
func executeReadOnly(
	context unsafe.Pointer,
	gasLimit int64,
	addressOffset int32,
	functionOffset int32,
	functionLength int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := core.GetVmContext(context)
	runtime := host.Runtime()
	output := host.Output()
	metering := host.Metering()

	address, err := runtime.MemLoad(addressOffset, core.HashLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	function, data, actualLen, err := getArgumentsFromMemory(
		context,
		functionOffset,
		functionLength,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	value := runtime.GetVMInput().CallValue
	sender := runtime.GetVMInput().CallerAddr

	gasToUse := metering.GasSchedule().DharitriAPICost.ExecuteReadOnly
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)

	err = output.Transfer(address, sender, 0, value, nil)
	if err != nil {
		return 1
	}

	runtime.SetReadOnly(true)

	contractCallInput := &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   value,
			GasPrice:    0,
			GasProvided: metering.BoundGasLimit(gasLimit),
		},
		RecipientAddr: address,
		Function:      function,
	}

	_, err = host.ExecuteOnSameContext(contractCallInput)
	runtime.SetReadOnly(false)
	if err != nil {
		return 1
	}

	return 0
}

//export createContract
func createContract(
	context unsafe.Pointer,
	valueOffset int32,
	codeOffset int32,
	length int32,
	resultOffset int32,
	numArguments int32,
	argumentsLengthOffset int32,
	dataOffset int32,
) int32 {
	host := core.GetVmContext(context)
	runtime := host.Runtime()
	metering := host.Metering()

	sender := runtime.GetSCAddress()
	value, err := runtime.MemLoad(valueOffset, core.BalanceLen)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	code, err := runtime.MemLoad(codeOffset, length)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	_, data, actualLen, err := getArgumentsFromMemory(
		context,
		0,
		0,
		numArguments,
		argumentsLengthOffset,
		dataOffset,
	)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse := metering.GasSchedule().DharitriAPICost.CreateContract
	gasToUse += metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(actualLen)
	metering.UseGas(gasToUse)
	gasLimit := metering.GasLeft()

	contractCreate := &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  sender,
			Arguments:   data,
			CallValue:   big.NewInt(0).SetBytes(value),
			GasPrice:    0,
			GasProvided: gasLimit,
		},
		ContractCode: code,
		// TODO: Receive code metadata as argument
		ContractCodeMetadata: []byte{1, 0},
	}

	newAddress, err := host.CreateNewContract(contractCreate)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, newAddress)
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export getNumReturnData
func getNumReturnData(context unsafe.Pointer) int32 {
	output := core.GetOutputContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetNumReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	return int32(len(returnData))
}

//export getReturnDataSize
func getReturnDataSize(context unsafe.Pointer, resultID int32) int32 {
	output := core.GetOutputContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetReturnDataSize
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) {
		return 0
	}

	return int32(len(returnData[resultID]))
}

//export getReturnData
func getReturnData(context unsafe.Pointer, resultID int32, dataOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	output := core.GetOutputContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetReturnData
	metering.UseGas(gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) {
		return 0
	}

	err := runtime.MemStore(dataOffset, returnData[resultID])
	if core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution()) {
		return 0
	}

	return int32(len(returnData[resultID]))
}

//export getOriginalTxHash
func getOriginalTxHash(context unsafe.Pointer, dataOffset int32) {
	runtime := core.GetRuntimeContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().DharitriAPICost.GetBlockHash
	metering.UseGas(gasToUse)

	err := runtime.MemStore(dataOffset, runtime.GetOriginalTxHash())
	_ = core.WithFault(err, context, runtime.DharitriAPIErrorShouldFailExecution())
}
