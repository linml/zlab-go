package util

import "unsafe"

//ToBytes string to bytes
func ToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//ToString bytes to string
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
