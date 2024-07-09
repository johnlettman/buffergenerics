package buffergenerics

import (
	"encoding/binary"
	"golang.org/x/exp/constraints"
	"io"
	"math"
	"reflect"
)

// ReadOrderedT reads a value of type T from the given buffer starting at the specified offset,
// using the specified byte order. If the byte order is nil, it defaults to binary.NativeEndian.
// It returns the read value and any error encountered during the read operation.
func ReadOrderedT[T constraints.Integer | constraints.Float](buffer []byte, offset int, order binary.ByteOrder) (T, error) {
	if order == nil {
		order = binary.ByteOrder(binary.NativeEndian)
	}

	typ := reflect.TypeFor[T]()
	zero := *new(T)
	kind := typ.Kind()
	size := typ.Bits() / 8
	end := offset + size

	if end > len(buffer) {
		return zero, io.EOF
	}

	switch kind {
	case reflect.Int8, reflect.Uint8:
		return T(buffer[offset]), nil
	case reflect.Int16, reflect.Uint16:
		u16 := order.Uint16(buffer[offset:end])
		return T(u16), nil
	case reflect.Int32, reflect.Uint32:
		u32 := order.Uint32(buffer[offset:end])
		return T(u32), nil
	case reflect.Int64, reflect.Uint64, reflect.Uintptr:
		u64 := order.Uint64(buffer[offset:end])
		return T(u64), nil
	case reflect.Float32:
		fu32 := order.Uint32(buffer[offset:end])
		return T(math.Float32frombits(fu32)), nil
	case reflect.Float64:
		fu64 := order.Uint64(buffer[offset:end])
		return T(math.Float64frombits(fu64)), nil
	default:
		return zero, NewErrUnknownKind(kind)
	}
}

// MustReadOrderedT reads a value of type T from the given buffer starting at the specified offset,
// using the specified byte order. If the byte order is nil, it defaults to binary.NativeEndian.
// It returns the read value. If an error is encountered during the read operation, it panics with the error.
// See also: ReadOrderedT.
func MustReadOrderedT[T constraints.Integer | constraints.Float](buffer []byte, offset int, order binary.ByteOrder) T {
	val, err := ReadOrderedT[T](buffer, offset, order)
	if err != nil {
		panic(err)
	}

	return val
}

// ReadT reads a value of type T from the given buffer starting at the specified offset.
// It uses binary.NativeEndian byte order. It returns the read value and any error encountered during the read operation.
// See also: ReadOrderedT.
func ReadT[T constraints.Integer | constraints.Float](buffer []byte, offset int) (T, error) {
	return ReadOrderedT[T](buffer, offset, binary.NativeEndian)
}

// MustReadT reads a value of type T from the given buffer starting at the specified offset.
// It uses the default byte order binary.NativeEndian. If an error is encountered during
// the read operation, it panics with the error.
// See also: ReadOrderedT.
func MustReadT[T constraints.Integer | constraints.Float](buffer []byte, offset int) T {
	return MustReadOrderedT[T](buffer, offset, binary.NativeEndian)
}
