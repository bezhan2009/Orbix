package PasswordAlgoritm

import (
	"fmt"
	"strings"
	"unsafe"
)

func GetRawBytes[T any](val *T) []byte {
	// unsafe.Sizeof gets the size of the type in bytes
	// unsafe.Slice converts the raw pointer into a readable []byte slice
	return unsafe.Slice((*byte)(unsafe.Pointer(val)), unsafe.Sizeof(*val))
}

func GetBinaryString[T any](val *T) string {
	// 1. Get the size of the type in bytes
	size := unsafe.Sizeof(*val)

	// 2. Convert the pointer to a slice of raw bytes directly from RAM
	rawBytes := unsafe.Slice((*byte)(unsafe.Pointer(val)), size)

	// 3. Convert each byte into 8 bits of 0s and 1s
	var bitStrings []string

	// We read backwards (from last byte to first byte)
	// because modern CPUs (Intel/AMD/ARM) use "Little Endian" format.
	for i := int(size) - 1; i >= 0; i-- {
		bitStrings = append(bitStrings, fmt.Sprintf("%08b", rawBytes[i]))
	}

	// Join the bytes with a space for easy reading
	return strings.Join(bitStrings, " ")
}

func Klg(password string, recourse int) string {
	for i := 0; i < recourse; i++ {
		password += string(GetRawBytes(&password)) + GetBinaryString(&password)
	}

	return password
}
