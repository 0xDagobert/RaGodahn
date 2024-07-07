package process_a_injector

import (
	"syscall"
)

func CreateProcess() *syscall.ProcessInformation {
	var si syscall.StartupInfo
	var pi syscall.ProcessInformation

	commandLine, err := syscall.UTF16PtrFromString("C:\\Windows\\SysWOW64\\notepad.exe")

	if err != nil {
		panic(err)
	}

	err = syscall.CreateProcess(
		nil,
		commandLine,
		nil,
		nil,
		false,
		0,
		nil,
		nil,
		&si,
		&pi)

	if err != nil {
		panic(err)
	}

	return &pi
}
