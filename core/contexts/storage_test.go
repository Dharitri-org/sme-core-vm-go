package contexts

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/Dharitri-org/sme-core-vm-go/config"
	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/mock"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
	"github.com/stretchr/testify/require"
)

var dharitriReservedTestPrefix = []byte("RESERVED")

func TestNewStorageContext(t *testing.T) {
	t.Parallel()

	host := &mock.VmHostMock{}
	mockBlockchain := &mock.BlockchainHookMock{}

	storageContext, err := NewStorageContext(host, mockBlockchain, dharitriReservedTestPrefix)
	require.Nil(t, err)
	require.NotNil(t, storageContext)
}

func TestStorageContext_SetAddress(t *testing.T) {
	t.Parallel()

	addressA := []byte("accountA")
	addressB := []byte("accountB")
	stubOutput := &mock.OutputContextStub{}
	accountA := &vmcommon.OutputAccount{
		Address:        addressA,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	accountB := &vmcommon.OutputAccount{
		Address:        addressB,
		Nonce:          0,
		BalanceDelta:   big.NewInt(0),
		Balance:        big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
	stubOutput.GetOutputAccountCalled = func(address []byte) (*vmcommon.OutputAccount, bool) {
		if bytes.Equal(address, addressA) {
			return accountA, false
		}
		if bytes.Equal(address, addressB) {
			return accountB, false
		}
		return nil, false
	}

	mockRuntime := &mock.RuntimeContextMock{}
	mockMetering := &mock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &mock.VmHostMock{
		OutputContext:   stubOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &mock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, dharitriReservedTestPrefix)

	keyA := []byte("keyA")
	valueA := []byte("valueA")

	storageContext.SetAddress(addressA)
	storageStatus, err := storageContext.SetStorage(keyA, valueA)
	require.Nil(t, err)
	require.Equal(t, core.StorageAdded, storageStatus)
	require.Equal(t, valueA, storageContext.GetStorage(keyA))
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 0)

	keyB := []byte("keyB")
	valueB := []byte("valueB")
	storageContext.SetAddress(addressB)
	storageStatus, err = storageContext.SetStorage(keyB, valueB)
	require.Nil(t, err)
	require.Equal(t, core.StorageAdded, storageStatus)
	require.Equal(t, valueB, storageContext.GetStorage(keyB))
	require.Len(t, storageContext.GetStorageUpdates(addressA), 1)
	require.Len(t, storageContext.GetStorageUpdates(addressB), 1)
	require.Equal(t, []byte(nil), storageContext.GetStorage(keyA))
}

func TestStorageContext_GetStorageUpdates(t *testing.T) {
	t.Parallel()

	mockOutput := &mock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount([]byte("account"))
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	account.StorageUpdates["update"] = &vmcommon.StorageUpdate{
		Offset: []byte("update"),
		Data:   []byte("some data"),
	}

	host := &mock.VmHostMock{
		OutputContext: mockOutput,
	}

	mockBlockchainHook := &mock.BlockchainHookMock{}
	storageContext, _ := NewStorageContext(host, mockBlockchainHook, dharitriReservedTestPrefix)

	storageUpdates := storageContext.GetStorageUpdates([]byte("account"))
	require.Equal(t, 1, len(storageUpdates))
	require.Equal(t, []byte("update"), storageUpdates["update"].Offset)
	require.Equal(t, []byte("some data"), storageUpdates["update"].Data)
}

func TestStorageContext_SetStorage(t *testing.T) {
	t.Parallel()

	address := []byte("account")
	mockOutput := &mock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(address)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &mock.RuntimeContextMock{}
	mockMetering := &mock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &mock.VmHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}
	bcHook := &mock.BlockchainHookStub{}

	storageContext, _ := NewStorageContext(host, bcHook, dharitriReservedTestPrefix)
	storageContext.SetAddress(address)

	key := []byte("key")
	value := []byte("value")

	storageStatus, err := storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, core.StorageAdded, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, core.StorageModified, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, core.StorageUnchanged, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	value = nil
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, core.StorageDeleted, storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(true)
	value = []byte("newValue")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, core.StorageUnchanged, storageStatus)
	require.Equal(t, []byte{}, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 1)

	mockRuntime.SetReadOnly(false)
	key = []byte("other_key")
	value = []byte("other_value")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Nil(t, err)
	require.Equal(t, core.StorageAdded, storageStatus)
	require.Equal(t, value, storageContext.GetStorage(key))
	require.Len(t, storageContext.GetStorageUpdates(address), 2)

	key = []byte("RESERVEDkey")
	value = []byte("doesn't matter")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Equal(t, core.ErrStoreDharitriReservedKey, err)

	key = []byte("RESERVED")
	value = []byte("doesn't matter")
	storageStatus, err = storageContext.SetStorage(key, value)
	require.Equal(t, core.ErrStoreDharitriReservedKey, err)
}

func TestStorageContext_GetStorageFromAddress(t *testing.T) {
	t.Parallel()

	scAddress := []byte("account")
	mockOutput := &mock.OutputContextMock{}
	account := mockOutput.NewVMOutputAccount(scAddress)
	mockOutput.OutputAccountMock = account
	mockOutput.OutputAccountIsNew = false

	mockRuntime := &mock.RuntimeContextMock{}
	mockMetering := &mock.MeteringContextMock{}
	mockMetering.SetGasSchedule(config.MakeGasMapForTests())
	mockMetering.BlockGasLimitMock = uint64(15000)

	host := &mock.VmHostMock{
		OutputContext:   mockOutput,
		MeteringContext: mockMetering,
		RuntimeContext:  mockRuntime,
	}

	readable := []byte("readable")
	nonreadable := []byte("nonreadable")
	internalData := []byte("internalData")

	bcHook := &mock.BlockchainHookStub{
		GetUserAccountCalled: func(address []byte) (vmcommon.UserAccountHandler, error) {
			if bytes.Equal(readable, address) {
				return &mock.AccountMock{CodeMetadata: []byte{4, 0}}, nil
			}
			if bytes.Equal(nonreadable, address) || bytes.Equal(scAddress, address) {
				return &mock.AccountMock{CodeMetadata: []byte{0, 0}}, nil
			}
			return nil, nil
		},
		GetStorageDataCalled: func(accountsAddress []byte, index []byte) ([]byte, error) {
			return internalData, nil
		},
	}

	storageContext, _ := NewStorageContext(host, bcHook, dharitriReservedTestPrefix)
	storageContext.SetAddress(scAddress)

	key := []byte("key")
	data := storageContext.GetStorageFromAddress(scAddress, key)
	require.Equal(t, data, internalData)

	data = storageContext.GetStorageFromAddress(readable, key)
	require.Equal(t, data, internalData)

	data = storageContext.GetStorageFromAddress(nonreadable, key)
	require.Nil(t, data)
}

func TestStorageContext_LoadGasStoreGasPerKey(t *testing.T) {
	// TODO
}

func TestStorageContext_StoreGasPerKey(t *testing.T) {
	// TODO
}
