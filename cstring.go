package main

import (
	"unicode/utf16"
	"unsafe"
)

func UintptrToString(cstr uintptr) string {
	if cstr != 0 {
		us := make([]uint16, 0, 256)
		for p := cstr; ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				return string(utf16.Decode(us))
			}
			us = append(us, u)
		}
	}
	return ""
}
