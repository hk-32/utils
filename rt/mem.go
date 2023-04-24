package rt

import "unsafe"

// Allocate memory that is argument size number of bytes
func Alloc(size int) unsafe.Pointer {
	block := make([]byte, size)
	return unsafe.Pointer(&block[0])
}

// Allocate memory that is big enough to store T
func AllocFor[T any]() unsafe.Pointer {
	block := make([]byte, SizeOfType[T]())
	return unsafe.Pointer(&block[0])
}

// Copies num bytes from src to the memory block pointed to by dst
func MemCpyRaw(dst unsafe.Pointer, src unsafe.Pointer, num int) {
	copy(
		unsafe.Slice((*byte)(dst), num),
		unsafe.Slice((*byte)(src), num),
	)
}

// Copies num bytes from src to the memory block pointed to by dst
func MemCpy[T1, T2 any](dst *T1, src *T2, num int) {
	MemCpyRaw(unsafe.Pointer(dst), unsafe.Pointer(src), num)
}

// Return the size of T
func SizeOfType[T any]() int {
	var item T
	return int(unsafe.Sizeof(item))
}

// Return the size of item
func SizeOf[T any](item T) int {
	return int(unsafe.Sizeof(item))
}
