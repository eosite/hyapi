package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type SicLib struct {
	DllHandle                syscall.Handle
	sic_init                 uintptr
	sic_getFirstFitSmartCode uintptr
	sic_getCode              uintptr
	sic_freeString           uintptr
}

func NewSicLib(dll string) *SicLib {
	h, err := syscall.LoadLibrary(dll)
	if err != nil {
		panic(err)
	}
	lib := &SicLib{DllHandle: h}
	return lib
}
func (s *SicLib) Free() {
	syscall.FreeLibrary(s.DllHandle)
}
func (s *SicLib) Init(block int) error {
	var err error
	s.sic_init, err = syscall.GetProcAddress(s.DllHandle, "sic_init")
	if err != nil {
		return err
	}
	s.sic_getFirstFitSmartCode, err = syscall.GetProcAddress(s.DllHandle, "sic_getFirstFitSmartCode")
	if err != nil {
		return err
	}
	s.sic_getCode, err = syscall.GetProcAddress(s.DllHandle, "sic_getCode")
	if err != nil {
		return err
	}
	s.sic_freeString, err = syscall.GetProcAddress(s.DllHandle, "sic_freeString")
	if err != nil {
		return err
	}
	//阻塞调用
	_, _, err = syscall.Syscall(s.sic_init, 1, uintptr(block), 0, 0)
	return err
}
func (s *SicLib) GetFirst(desc string) string {
	r, _, _ := syscall.Syscall(s.sic_getFirstFitSmartCode, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(desc))), 0, 0)
	return fmt.Sprintf("%04d", r)
}
func (s *SicLib) GetCode(desc string, maxBookCode, maxSmartCode int) string {
	r, _, _ := syscall.Syscall(s.sic_getCode, 3, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(desc))), uintptr(maxBookCode), uintptr(maxSmartCode))
	code := UintptrToString(r)
	syscall.Syscall(s.sic_freeString, 1, r, 0, 0)
	return code
}
