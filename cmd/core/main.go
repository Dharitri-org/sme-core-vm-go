package main

import (
	"fmt"
	"os"

	"github.com/Dharitri-org/sme-core-vm-go/ipc/common"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/corepart"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/marshaling"
	"github.com/Dharitri-org/sme-logger/pipes"
)

const (
	fileDescriptorCoreInit       = 3
	fileDescriptorNodeToCore     = 4
	fileDescriptorCoreToNode     = 5
	fileDescriptorReadLogProfile = 6
	fileDescriptorLogToNode      = 7
)

func main() {
	errCode, errMessage := doMain()
	if errCode != common.ErrCodeSuccess {
		fmt.Fprintln(os.Stderr, errMessage)
		os.Exit(errCode)
	}
}

// doMain returns (error code, error message)
func doMain() (int, string) {
	coreInitFile := getPipeFile(fileDescriptorCoreInit)
	if coreInitFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [coreInitFile]"
	}

	nodeToCoreFile := getPipeFile(fileDescriptorNodeToCore)
	if nodeToCoreFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [nodeToCoreFile]"
	}

	coreToNodeFile := getPipeFile(fileDescriptorCoreToNode)
	if coreToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [coreToNodeFile]"
	}

	readLogProfileFile := getPipeFile(fileDescriptorReadLogProfile)
	if readLogProfileFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [readLogProfileFile]"
	}

	logToNodeFile := getPipeFile(fileDescriptorLogToNode)
	if logToNodeFile == nil {
		return common.ErrCodeCannotCreateFile, "Cannot get pipe file: [logToNodeFile]"
	}

	coreArguments, err := common.GetCoreArguments(coreInitFile)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot receive gasSchedule: %v", err)
	}

	messagesMarshalizer := marshaling.CreateMarshalizer(coreArguments.MessagesMarshalizer)
	logsMarshalizer := marshaling.CreateMarshalizer(coreArguments.LogsMarshalizer)

	logsPart, err := pipes.NewChildPart(readLogProfileFile, logToNodeFile, logsMarshalizer)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot create logs part: %v", err)
	}

	err = logsPart.StartLoop()
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot start logs loop: %v", err)
	}

	defer logsPart.StopLoop()

	part, err := corepart.NewCorePart(
		nodeToCoreFile,
		coreToNodeFile,
		&coreArguments.VMHostParameters,
		messagesMarshalizer,
	)
	if err != nil {
		return common.ErrCodeInit, fmt.Sprintf("Cannot create CorePart: %v", err)
	}

	err = part.StartLoop()
	if err != nil {
		return common.ErrCodeTerminated, fmt.Sprintf("Ended Core loop: %v", err)
	}

	// This is never reached, actually. Core is supposed to run an infinite message loop.
	return common.ErrCodeSuccess, ""
}

func getPipeFile(fileDescriptor uintptr) *os.File {
	file := os.NewFile(fileDescriptor, fmt.Sprintf("/proc/self/fd/%d", fileDescriptor))
	return file
}
