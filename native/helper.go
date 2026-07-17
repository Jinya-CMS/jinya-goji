package native

import "unsafe"

func CStringToGoString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	bytes := make([]byte, 0)

	for {
		b := *(*byte)(unsafe.Pointer(ptr))
		if b == 0 {
			break
		}

		bytes = append(bytes, b)
		ptr++
	}

	return string(bytes)
}
