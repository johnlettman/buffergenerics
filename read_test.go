package buffergenerics

import (
	"encoding/binary"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"testing"
)

func TestReadOrderedT(t *testing.T) {
	t.Run("it should return an EOF error for out-of-bounds reads", func(t *testing.T) {
		buf := []byte{0xDE, 0xAD, 0xCA, 0xFE}
		_, err := ReadOrderedT[byte](buf, len(buf)+gofakeit.Int(), binary.LittleEndian)

		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("it should return an EOF error for too-large-type reads", func(t *testing.T) {
		buf := []byte{0xDE, 0xAD, 0xCA, 0xFE}
		_, err := ReadOrderedT[int64](buf, 0, binary.LittleEndian)

		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("it should assume binary.NativeEndian if no order is provided", func(t *testing.T) {
		want := gofakeit.Int64()
		buf := make([]byte, 8)
		binary.NativeEndian.PutUint64(buf, uint64(want))

		i64 := MustReadOrderedT[int64](buf, 0, nil)

		assert.Equal(t, want, i64)
	})
}

func TestMustReadOrderedT(t *testing.T) {
	t.Run("it should panic with EOF for out-of-bounds reads", func(t *testing.T) {
		assert.PanicsWithError(t, "EOF", func() {
			buf := []byte{0xDE, 0xAD, 0xCA, 0xFE}
			_ = MustReadOrderedT[byte](buf, len(buf)+gofakeit.Int(), binary.LittleEndian)
		})
	})

	t.Run("it should panic with EOF for too-large-type reads", func(t *testing.T) {
		assert.PanicsWithError(t, "EOF", func() {
			buf := []byte{0xDE, 0xAD, 0xCA, 0xFE}
			_ = MustReadOrderedT[int64](buf, 0, binary.LittleEndian)
		})
	})

	t.Run("it should otherwise passthrough to ReadOrderedT", func(t *testing.T) {
		assert.NotPanics(t, func() {
			order := binary.NativeEndian
			want := gofakeit.Int64()
			buf := make([]byte, 8)
			order.PutUint64(buf, uint64(want))

			i64 := MustReadOrderedT[int64](buf, 0, order)

			assert.Equal(t, want, i64)
		})
	})
}

func TestReadT(t *testing.T) {
	t.Run("it should passthrough to ReadOrderedT using binary.NativeEndian order", func(t *testing.T) {
		want := gofakeit.Int64()
		buf := make([]byte, 8)
		binary.NativeEndian.PutUint64(buf, uint64(want))

		i64, err := ReadT[int64](buf, 0)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, i64)
	})
}

func TestMustReadT(t *testing.T) {
	t.Run("it should passthrough to MustReadOrderedT using binary.NativeEndian order", func(t *testing.T) {
		assert.NotPanics(t, func() {
			want := gofakeit.Int64()
			buf := make([]byte, 8)
			binary.NativeEndian.PutUint64(buf, uint64(want))

			i64 := MustReadT[int64](buf, 0)

			assert.Equal(t, want, i64)
		})
	})
}

func TestReadOrderedT_SingleByte(t *testing.T) {
	t.Run("it should handle uint8 reads", func(t *testing.T) {
		want := gofakeit.Uint8()
		buf := []byte{want}

		u8, err := ReadOrderedT[uint8](buf, 0, binary.LittleEndian)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, u8)
	})

	t.Run("it should handle int8 reads", func(t *testing.T) {
		want := int8(math.MinInt8)
		buf := []byte{byte(want)}

		i8, err := ReadOrderedT[int8](buf, 0, binary.LittleEndian)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, i8)
	})

	t.Run("it should handle byte reads", func(t *testing.T) {
		want := gofakeit.Uint8()
		buf := []byte{want}

		b, err := ReadOrderedT[byte](buf, 0, binary.LittleEndian)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, b)
	})

	t.Run("it should handle custom single byte types", func(t *testing.T) {
		type myType byte
		want := myType(gofakeit.Uint8())
		buf := []byte{byte(want)}

		b, err := ReadOrderedT[myType](buf, 0, binary.LittleEndian)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, b)
	})
}

func doTestReadOrderedT_Order(t *testing.T, order binary.ByteOrder) {
	name := order.String()

	t.Run("it should handle custom "+name+" multibyte types", func(t *testing.T) {
		type myType int32
		want := myType(gofakeit.Int32())
		buf := make([]byte, 4)
		order.PutUint32(buf, uint32(want))

		i32, err := ReadOrderedT[myType](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, i32)
	})

	t.Run("it should handle uintptr "+name+" reads", func(t *testing.T) {
		want := uintptr(gofakeit.Uint64())
		buf := make([]byte, 8)
		order.PutUint64(buf, uint64(want))

		pointer, err := ReadOrderedT[uintptr](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, pointer)
	})

	t.Run("it should handle uint16 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Uint16()
		buf := make([]byte, 2)
		order.PutUint16(buf, want)

		u16, err := ReadOrderedT[uint16](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, u16)
	})

	t.Run("it should handle int16 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Int16()
		buf := make([]byte, 2)
		order.PutUint16(buf, uint16(want))

		i16, err := ReadOrderedT[int16](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, i16)
	})

	t.Run("it should handle uint32 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Uint32()
		buf := make([]byte, 4)
		order.PutUint32(buf, want)

		u32, err := ReadOrderedT[uint32](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, u32)
	})

	t.Run("it should handle int32 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Int32()
		buf := make([]byte, 4)
		order.PutUint32(buf, uint32(want))

		i32, err := ReadOrderedT[int32](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, i32)
	})

	t.Run("it should handle uint64 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Uint64()
		buf := make([]byte, 8)
		order.PutUint64(buf, want)

		u64, err := ReadOrderedT[uint64](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, u64)
	})

	t.Run("it should handle int64 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Int64()
		buf := make([]byte, 8)
		order.PutUint64(buf, uint64(want))

		i64, err := ReadOrderedT[int64](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, i64)
	})

	t.Run("it should handle float32 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Float32()
		buf := make([]byte, 4)
		order.PutUint32(buf, math.Float32bits(want))

		f32, err := ReadOrderedT[float32](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, f32)
	})

	t.Run("it should handle float64 "+name+" reads", func(t *testing.T) {
		want := gofakeit.Float64()
		buf := make([]byte, 8)
		order.PutUint64(buf, math.Float64bits(want))

		f64, err := ReadOrderedT[float64](buf, 0, order)

		assert.NoError(t, err, "it should not return an error")
		assert.Equal(t, want, f64)
	})
}

func TestReadOrderedT_BigEndian(t *testing.T) {
	doTestReadOrderedT_Order(t, binary.BigEndian)
}

func TestReadOrderedT_LittleEndian(t *testing.T) {
	doTestReadOrderedT_Order(t, binary.LittleEndian)
}

func TestReadOrderedT_NativeEndian(t *testing.T) {
	doTestReadOrderedT_Order(t, binary.NativeEndian)
}
