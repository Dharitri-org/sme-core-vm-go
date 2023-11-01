package cryptoapi

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
//
// extern int32_t sha256(void* context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t keccak256(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t ripemd160(void *context, int32_t dataOffset, int32_t length, int32_t resultOffset);
// extern int32_t verifyBLS(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t verifyEd25519(void *context, int32_t keyOffset, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
// extern int32_t verifySecp256k1(void *context, int32_t keyOffset, int32_t keyLength, int32_t messageOffset, int32_t messageLength, int32_t sigOffset);
import "C"

import (
	"unsafe"

	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/wasmer"
)

const BlsPublicKeyLength = 96
const BlsSignatureLength = 48
const Ed25519PublicKeyLength = 32
const Ed25519SignatureLength = 64
const Secp256k1CompressedPublicKeyLength = 33
const Secp256k1UncompressedPublicKeyLength = 65
const Secp256k1SignatureLength = 64

func CryptoImports(imports *wasmer.Imports) (*wasmer.Imports, error) {
	imports = imports.Namespace("env")
	imports, err := imports.Append("sha256", sha256, C.sha256)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("keccak256", keccak256, C.keccak256)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("ripemd160", ripemd160, C.ripemd160)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyBLS", verifyBLS, C.verifyBLS)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifyEd25519", verifyEd25519, C.verifyEd25519)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("verifySecp256k1", verifySecp256k1, C.verifySecp256k1)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

//export sha256
func sha256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	crypto := core.GetCryptoContext(context)
	metering := core.GetMeteringContext(context)

	memLoadGas := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	gasToUse := metering.GasSchedule().CryptoAPICost.SHA256 + memLoadGas
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Sha256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export keccak256
func keccak256(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	crypto := core.GetCryptoContext(context)
	metering := core.GetMeteringContext(context)

	memLoadGas := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	gasToUse := metering.GasSchedule().CryptoAPICost.Keccak256 + memLoadGas
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Keccak256(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export ripemd160
func ripemd160(context unsafe.Pointer, dataOffset int32, length int32, resultOffset int32) int32 {
	runtime := core.GetRuntimeContext(context)
	crypto := core.GetCryptoContext(context)
	metering := core.GetMeteringContext(context)

	memLoadGas := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
	gasToUse := metering.GasSchedule().CryptoAPICost.Ripemd160 + memLoadGas
	metering.UseGas(gasToUse)

	data, err := runtime.MemLoad(dataOffset, length)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	result, err := crypto.Ripemd160(data)
	if err != nil {
		return 1
	}

	err = runtime.MemStore(resultOffset, result)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	return 0
}

//export verifyBLS
func verifyBLS(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := core.GetRuntimeContext(context)
	crypto := core.GetCryptoContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyBLS
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, BlsPublicKeyLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(messageLength)
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, BlsSignatureLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyBLS(key, message, sig)
	if invalidSigErr != nil {
		return -1
	}

	return 0
}

//export verifyEd25519
func verifyEd25519(
	context unsafe.Pointer,
	keyOffset int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := core.GetRuntimeContext(context)
	crypto := core.GetCryptoContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifyEd25519
	metering.UseGas(gasToUse)

	key, err := runtime.MemLoad(keyOffset, Ed25519PublicKeyLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(messageLength)
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, Ed25519SignatureLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifyEd25519(key, message, sig)
	if invalidSigErr != nil {
		return -1
	}

	return 0
}

//export verifySecp256k1
func verifySecp256k1(
	context unsafe.Pointer,
	keyOffset int32,
	keyLength int32,
	messageOffset int32,
	messageLength int32,
	sigOffset int32,
) int32 {
	runtime := core.GetRuntimeContext(context)
	crypto := core.GetCryptoContext(context)
	metering := core.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().CryptoAPICost.VerifySecp256k1
	metering.UseGas(gasToUse)

	if keyLength != Secp256k1CompressedPublicKeyLength && keyLength != Secp256k1UncompressedPublicKeyLength {
		core.WithFault(core.ErrInvalidPublicKeySize, context, runtime.DharitriAPIErrorShouldFailExecution())
		return 1
	}

	key, err := runtime.MemLoad(keyOffset, keyLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(messageLength)
	metering.UseGas(gasToUse)

	message, err := runtime.MemLoad(messageOffset, messageLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	sig, err := runtime.MemLoad(sigOffset, Secp256k1SignatureLength)
	if core.WithFault(err, context, runtime.CryptoAPIErrorShouldFailExecution()) {
		return 1
	}

	invalidSigErr := crypto.VerifySecp256k1(key, message, sig)
	if invalidSigErr != nil {
		return -1
	}

	return 0
}
