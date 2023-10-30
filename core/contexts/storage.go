package contexts

import (
	"bytes"
	"errors"

	"github.com/Dharitri-org/sme-core-vm-go/core"
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
	storageUpdates := context.GetStorageUpdates(context.address)
	if storageUpdate, ok := storageUpdates[string(key)]; ok {
		return storageUpdate.Data
	}

	value, _ := context.blockChainHook.GetStorageData(context.address, key)
	if value != nil {
		storageUpdates[string(key)] = &vmcommon.StorageUpdate{
			Offset: key,
			Data:   value,
		}
	}

	return value
}

func (context *storageContext) isDharitriReservedKey(key []byte) bool {
	return bytes.HasPrefix(key, []byte(context.dharitriProtectedKeyPrefix))
}

func (context *storageContext) SetStorage(key []byte, value []byte) (core.StorageStatus, error) {
	if context.isDharitriReservedKey(key) {
		return core.StorageUnchanged, core.ErrStoreDharitriReservedKey
	}

	if context.host.Runtime().ReadOnly() {
		return core.StorageUnchanged, nil
	}

	metering := context.host.Metering()
	var zero []byte
	strKey := string(key)
	length := len(value)

	var oldValue []byte
	storageUpdates := context.GetStorageUpdates(context.address)
	if update, ok := storageUpdates[strKey]; !ok {
		oldValue = context.GetStorage(key)
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
