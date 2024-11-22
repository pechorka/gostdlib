package stringx

import "unsafe"

func ToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func FromBytes(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
