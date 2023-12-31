package contexts

import (
	"bytes"
	"errors"

	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-logger/check"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

const lockKeyContext string = "timelock"

type storageContext struct {
	host                       core.VMHost
	blockChainHook             vmcommon.BlockchainHook
	address                    []byte
	stateStack                 [][]byte
	dharitriProtectedKeyPrefix []byte
}

// NewStorageContext creates a new storageContext
func NewStorageContext(
	host core.VMHost,
	blockChainHook vmcommon.BlockchainHook,
	dharitriProtectedKeyPrefix []byte,
) (*storageContext, error) {
	if len(dharitriProtectedKeyPrefix) == 0 {
		return nil, errors.New("dharitriProtectedKeyPrefix cannot be empty")
	}
	context := &storageContext{
		host:                       host,
		blockChainHook:             blockChainHook,
		stateStack:                 make([][]byte, 0),
		dharitriProtectedKeyPrefix: dharitriProtectedKeyPrefix,
	}

	return context, nil
}

func (context *storageContext) InitState() {
}

func (context *storageContext) PushState() {
	context.stateStack = append(context.stateStack, context.address)
}

func (context *storageContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	prevAddress := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	context.address = prevAddress
}

func (context *storageContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	context.stateStack = context.stateStack[:stateStackLen-1]
}

func (context *storageContext) ClearStateStack() {
	context.stateStack = make([][]byte, 0)
}

func (context *storageContext) SetAddress(address []byte) {
	context.address = address
}

func (context *storageContext) GetStorageUpdates(address []byte) map[string]*vmcommon.StorageUpdate {
	account, _ := context.host.Output().GetOutputAccount(address)
	return account.StorageUpdates
}

func (context *storageContext) GetStorage(key []byte) []byte {
	metering := context.host.Metering()

	extraBytes := len(key) - core.AddressLen
	if extraBytes > 0 {
		gasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(extraBytes)
		metering.UseGas(gasToUse)
	}

	value := context.GetStorageUnmetered(key)

	gasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(value))
	metering.UseGas(gasToUse)

	return value
}

func (context *storageContext) GetStorageFromAddress(address []byte, key []byte) []byte {
	metering := context.host.Metering()

	extraBytes := len(key) - core.AddressLen
	if extraBytes > 0 {
		gasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(extraBytes)
		metering.UseGas(gasToUse)
	}

	if !bytes.Equal(address, context.address) {
		userAcc, err := context.blockChainHook.GetUserAccount(address)
		if err != nil || check.IfNil(userAcc) {
			return nil
		}

		metadata := vmcommon.CodeMetadataFromBytes(userAcc.GetCodeMetadata())
		if !metadata.Readable {
			return nil
		}
	}

	value := context.getStorageFromAddressUnmetered(address, key)

	gasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(len(value))
	metering.UseGas(gasToUse)

	return value
}

func (context *storageContext) getStorageFromAddressUnmetered(address []byte, key []byte) []byte {
	var value []byte

	storageUpdates := context.GetStorageUpdates(address)
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		value = storageUpdate.Data
	} else {
		value, _ = context.blockChainHook.GetStorageData(address, key)
		storageUpdates[string(key)] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   value,
		}
	}

	return value
}

func (context *storageContext) GetStorageUnmetered(key []byte) []byte {
	return context.getStorageFromAddressUnmetered(context.address, key)
}

func (context *storageContext) isDharitriReservedKey(key []byte) bool {
	return bytes.HasPrefix(key, context.dharitriProtectedKeyPrefix)
}

func (context *storageContext) SetStorage(key []byte, value []byte) (core.StorageStatus, error) {
	if context.isDharitriReservedKey(key) {
		return core.StorageUnchanged, core.ErrStoreDharitriReservedKey
	}

	if context.host.Runtime().ReadOnly() {
		return core.StorageUnchanged, nil
	}

	metering := context.host.Metering()

	extraBytes := len(key) - core.AddressLen
	if extraBytes > 0 {
		gasToUse := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(extraBytes)
		metering.UseGas(gasToUse)
	}

	var zero []byte
	strKey := string(key)
	length := len(value)

	var oldValue []byte
	storageUpdates := context.GetStorageUpdates(context.address)
	if update, ok := storageUpdates[strKey]; !ok {
		oldValue = context.GetStorageUnmetered(key)
		storageUpdates[strKey] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   oldValue,
		}
	} else {
		oldValue = update.Data
	}

	lengthOldValue := len(oldValue)
	if bytes.Equal(oldValue, value) {
		useGas := metering.GasSchedule().BaseOperationCost.DataCopyPerByte * uint64(length)
		metering.UseGas(useGas)
		return core.StorageUnchanged, nil
	}

	newUpdate := &vmcommon.StorageUpdate{
		Offset: key,
		Data:   make([]byte, length),
	}
	copy(newUpdate.Data[:length], value[:length])
	storageUpdates[strKey] = newUpdate

	if bytes.Equal(oldValue, zero) {
		useGas := metering.GasSchedule().BaseOperationCost.StorePerByte * uint64(length)
		metering.UseGas(useGas)
		return core.StorageAdded, nil
	}
	if bytes.Equal(value, zero) {
		freeGas := metering.GasSchedule().BaseOperationCost.ReleasePerByte * uint64(lengthOldValue)
		metering.FreeGas(freeGas)
		return core.StorageDeleted, nil
	}

	newValueExtraLength := length - lengthOldValue
	if newValueExtraLength > 0 {
		useGas := metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(lengthOldValue)
		useGas += metering.GasSchedule().BaseOperationCost.StorePerByte * uint64(newValueExtraLength)
		metering.UseGas(useGas)
	}
	if newValueExtraLength < 0 {
		newValueExtraLength = -newValueExtraLength

		useGas := metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(length)
		metering.UseGas(useGas)

		freeGas := metering.GasSchedule().BaseOperationCost.ReleasePerByte * uint64(newValueExtraLength)
		metering.FreeGas(freeGas)
	}

	return core.StorageModified, nil
}
