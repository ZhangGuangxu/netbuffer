package netbuffer

import (
	"bytes"
	"math"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer()
	if cap(buf.buf) != cheapPrepend+initialSize {
		t.Errorf("cap(buf.buf) != %d", cheapPrepend+initialSize)
	}
}

func TestNewBufferWithSize(t *testing.T) {
	size := 10
	buf := NewBufferWithSize(size)
	if cap(buf.buf) != cheapPrepend+size {
		t.Errorf("cap(buf.buf) != %d", cheapPrepend+size)
	}

	err := buf.AppendInt64(int64(9223372036854770000))
	if err != nil {
		t.Errorf("buf.AppendInt64 error %v", err)
	}
	b := buf.WritableByteSlice()
	if len(b) != 2 {
		t.Errorf("len(buf.WritableByteSlice()) = %d, want 2", len(b))
	}
}

func TestAppend(t *testing.T) {
	buf := NewBuffer()

	data := []byte("abcde")

	buf.Append(data)
	if buf.ReadableBytes() != len(data) {
		t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), len(data))
	}

	writable := cap(buf.buf) - (cheapPrepend + len(data))
	if buf.WritableBytes() != writable {
		t.Errorf("buf.WritableBytes() = %d, want %d", buf.WritableBytes(), writable)
	}

	if buf.readerIndex != cheapPrepend {
		t.Errorf("buf.readerIndex = %d, not %d", buf.readerIndex, cheapPrepend)
	}
	if buf.writerIndex != cheapPrepend+len(data) {
		t.Errorf("buf.writerIndex = %d, not %d", buf.writerIndex, cheapPrepend+len(data))
	}

	data = bytes.Repeat([]byte("a"), 2048)
	buf.Append(data)
	newCap := cheapPrepend + len([]byte("abcde")) + 2048
	if len(buf.buf) != newCap {
		t.Errorf("len(buf.buf) = %d, want %d", len(buf.buf), newCap)
	}

	buf.Append(data)
	buf.Append(data)
	newCap = cheapPrepend + len([]byte("abcde")) + 2048*3
	if len(buf.buf) != newCap {
		t.Errorf("len(buf.buf) = %d, want %d", len(buf.buf), newCap)
	}

	buf.Retrieve(2048 + 1024)
	buf.Append(data)
	if buf.readerIndex != cheapPrepend {
		t.Error("buf.Retrieve, then buf.Append, then buf.readerIndex is wrong")
	}
	if buf.writerIndex != buf.readerIndex+buf.ReadableBytes() {
		t.Error("buf.Retrieve, then buf.Append, then buf.writerIndex is wrong")
	}

	for _, v := range []int8{math.MaxInt8, 0, math.MinInt8} {
		buf := NewBuffer()
		err := buf.AppendInt8(v)
		if err != nil {
			t.Errorf("buf.AppendInt8 error %v", err)
		}
		if buf.ReadableBytes() != 1 {
			t.Errorf("After buf.AppendInt8, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 1)
		}
	}

	for _, v := range []int16{math.MaxInt16, 0, math.MinInt16} {
		buf := NewBuffer()
		err := buf.AppendInt16(v)
		if err != nil {
			t.Errorf("buf.AppendInt16 error %v", err)
		}
		if buf.ReadableBytes() != 2 {
			t.Errorf("After buf.AppendInt16, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 2)
		}
	}

	for _, v := range []int32{math.MaxInt32, 0, math.MinInt32} {
		buf := NewBuffer()
		err := buf.AppendInt32(v)
		if err != nil {
			t.Errorf("buf.AppendInt32 error %v", err)
		}
		if buf.ReadableBytes() != 4 {
			t.Errorf("After buf.AppendInt32, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 4)
		}
	}

	for _, v := range []int64{math.MaxInt64, 0, math.MinInt64} {
		buf := NewBuffer()
		err := buf.AppendInt64(v)
		if err != nil {
			t.Errorf("buf.AppendInt64 error %v", err)
		}
		if buf.ReadableBytes() != 8 {
			t.Errorf("After buf.AppendInt64, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 8)
		}
	}

	for _, v := range []uint8{math.MaxUint8, 0, 100} {
		buf := NewBuffer()
		err := buf.AppendUint8(v)
		if err != nil {
			t.Errorf("buf.AppendUint8 error %v", err)
		}
		if buf.ReadableBytes() != 1 {
			t.Errorf("After buf.AppendUint8, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 1)
		}
	}

	for _, v := range []uint16{math.MaxUint16, 0, 32000} {
		buf := NewBuffer()
		err := buf.AppendUint16(v)
		if err != nil {
			t.Errorf("buf.AppendUint16 error %v", err)
		}
		if buf.ReadableBytes() != 2 {
			t.Errorf("After buf.AppendUint16, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 2)
		}
	}

	for _, v := range []uint32{math.MaxUint32, 0, 3_000_000_000}{
		buf := NewBuffer()
		err := buf.AppendUint32(v)
		if err != nil {
			t.Errorf("buf.AppendUint32 error %v", err)
		}
		if buf.ReadableBytes() != 4 {
			t.Errorf("After buf.AppendUint32, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 4)
		}
	}

	for _, v := range []uint64{math.MaxUint64, 0, 999_999_999_999_999_999}{
		buf := NewBuffer()
		err := buf.AppendUint64(v)
		if err != nil {
			t.Errorf("buf.AppendUint64 error %v", err)
		}
		if buf.ReadableBytes() != 8 {
			t.Errorf("After buf.AppendUint64, buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 8)
		}
	}
}

func TestPrepend(t *testing.T) {
	for _, v := range []int8{math.MaxInt8, 0, math.MinInt8} {
		buf := NewBufferWithSize(4)
		err := buf.PrependInt8(v)
		if err != nil {
			t.Errorf("buf.PrependInt8 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-1 {
			t.Errorf("After buf.PrependInt8, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-1)
		}
	}

	for _, v := range []int16{math.MaxInt16, 0, math.MinInt16} {
		buf := NewBufferWithSize(4)
		err := buf.PrependInt16(v)
		if err != nil {
			t.Errorf("buf.PrependInt16 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-2 {
			t.Errorf("After buf.PrependInt16, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-2)
		}
	}

	for _, v := range []int32{math.MaxInt32, 0, math.MinInt32} {
		buf := NewBufferWithSize(4)
		err := buf.PrependInt32(v)
		if err != nil {
			t.Errorf("buf.PrependInt32 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-4 {
			t.Errorf("After buf.PrependInt32, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-4)
		}
	}

	for _, v := range []int64{math.MaxInt64, 0, math.MinInt64} {
		buf := NewBufferWithSize(4)
		err := buf.PrependInt64(v)
		if err != nil {
			t.Errorf("buf.PrependInt64 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-8 {
			t.Errorf("After buf.PrependInt64, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-8)
		}
	}

	for _, v := range []uint8{math.MaxUint8, 0, 100} {
		buf := NewBufferWithSize(4)
		err := buf.PrependUint8(v)
		if err != nil {
			t.Errorf("buf.PrependUint8 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-1 {
			t.Errorf("After buf.PrependUint8, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-1)
		}
	}

	for _, v := range []uint16{math.MaxUint16, 0, 32000} {
		buf := NewBufferWithSize(4)
		err := buf.PrependUint16(v)
		if err != nil {
			t.Errorf("buf.PrependUint16 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-2 {
			t.Errorf("After buf.PrependUint16, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-2)
		}
	}

	for _, v := range []uint32{math.MaxUint32, 0, 2_000_000_000} {
		buf := NewBufferWithSize(4)
		err := buf.PrependUint32(v)
		if err != nil {
			t.Errorf("buf.PrependUint32 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-4 {
			t.Errorf("After buf.PrependUint32, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-4)
		}
	}

	for _, v := range []uint64{math.MaxUint64, 0, 321_111_111_999_999_000} {
		buf := NewBufferWithSize(4)
		err := buf.PrependUint64(v)
		if err != nil {
			t.Errorf("buf.PrependUint64 error %v", err)
		}
		if buf.prependableBytes() != cheapPrepend-8 {
			t.Errorf("After buf.PrependUint64, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-8)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.prepend(make([]byte, 5))
		if buf.prependableBytes() != cheapPrepend-5 {
			t.Errorf("After buf.prepend, buf.prependableBytes() is %d, want %d",
				buf.prependableBytes(), cheapPrepend-5)
		}
	}
}

func TestRetrieve(t *testing.T) {
	for _, v := range []int64{math.MaxInt64, 0, math.MinInt64} {
		buf := NewBuffer()
		err := buf.AppendInt64(v)
		if err != nil {
			t.Errorf("buf.AppendInt64() error %v", err)
		}
		buf.RetrieveAll()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []int64{math.MaxInt64, 0, math.MinInt64} {
		buf := NewBuffer()
		err := buf.AppendInt64(v)
		if err != nil {
			t.Errorf("buf.AppendInt64() error %v", err)
		}
		buf.RetrieveInt64()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []int32{math.MaxInt32, 0, math.MinInt32} {
		buf := NewBuffer()
		err := buf.AppendInt32(v)
		if err != nil {
			t.Errorf("buf.AppendInt32() error %v", err)
		}
		buf.RetrieveInt32()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []int16{math.MaxInt16, 0, math.MinInt16} {
		buf := NewBuffer()
		err := buf.AppendInt16(v)
		if err != nil {
			t.Errorf("buf.AppendInt16() error %v", err)
		}
		buf.RetrieveInt16()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []int8{math.MaxInt8, 0, math.MinInt8} {
		buf := NewBuffer()
		err := buf.AppendInt8(v)
		if err != nil {
			t.Errorf("buf.AppendInt8() error %v", err)
		}
		buf.RetrieveInt8()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []uint64{math.MaxUint64, 0, 999_999_999_999_999} {
		buf := NewBuffer()
		err := buf.PrependUint64(v)
		if err != nil {
			t.Errorf("buf.PrependUint64 error %v", err)
		}
		buf.RetrieveAll()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []uint64{math.MaxUint64, 0, 111_111_111_111_111} {
		buf := NewBuffer()
		err := buf.PrependUint64(v)
		if err != nil {
			t.Errorf("buf.PrependUint64 error %v", err)
		}
		buf.RetrieveUint64()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []uint32{math.MaxUint32, 0, 2_000_000_000} {
		buf := NewBuffer()
		err := buf.PrependUint32(v)
		if err != nil {
			t.Errorf("buf.PrependUint32 error %v", err)
		}
		buf.RetrieveUint32()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []uint16{math.MaxUint16, 0, 32000} {
		buf := NewBuffer()
		err := buf.PrependUint16(v)
		if err != nil {
			t.Errorf("buf.PrependUint16 error %v", err)
		}
		buf.RetrieveUint16()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	for _, v := range []uint8{math.MaxUint8, 0, 100} {
		buf := NewBuffer()
		err := buf.PrependUint8(v)
		if err != nil {
			t.Errorf("buf.PrependUint8 error %v", err)
		}
		buf.RetrieveUint8()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	{
		buf := NewBuffer()
		s := "powerful"
		buf.Append([]byte(s))
		data := buf.retrieveAllAsByteSlice()
		if len(data) != len(s) {
			t.Errorf("len(data) = %d, want %d", len(data), len(s))
		}
	}

	{
		buf := NewBuffer()
		s := "9223372036854770000"
		buf.Append([]byte(s))
		str := buf.retrieveAllAsString()
		if len(str) != len(s) {
			t.Errorf("len(str) = %d, want %d", len(str), len(s))
		}
	}

	{
		buf := NewBuffer()
		err := buf.PrependInt64(int64(math.MaxInt64))
		if err != nil {
			t.Errorf("buf.PrependInt64() error %v", err)
		}
		data := buf.retrieveAllAsByteSlice()
		if len(data) != 8 {
			t.Errorf("len(data) = %d, want %d", len(data), 8)
		}
	}

	{
		buf := NewBuffer()
		err := buf.PrependInt64(int64(math.MinInt64))
		if err != nil {
			t.Errorf("buf.PrependInt64() error %v", err)
		}
		str := buf.retrieveAllAsString()
		if len(str) == 0 {
			t.Error("len(data) should be >0")
		}
	}
}

func TestPeek(t *testing.T) {
	for _, v := range []int8{math.MaxInt8, 0, math.MinInt8} {
		buf := NewBuffer()
		err := buf.AppendInt8(v)
		if err != nil {
			t.Errorf("buf.AppendInt8 error %v", err)
		}
		x, err := buf.PeekInt8()
		if err != nil {
			t.Errorf("buf.PeekInt8 error %v", err)
		}
		if x != v {
			t.Errorf("After buf.AppendInt8, buf.PeekInt8() = %d, want %d", x, v)
		}
	}

	for _, v := range []int16{math.MaxInt16, 0, math.MinInt16} {
		buf := NewBuffer()
		err := buf.AppendInt16(v)
		if err != nil {
			t.Errorf("buf.AppendInt16 error %v", err)
		}
		x, err := buf.PeekInt16()
		if err != nil {
			t.Errorf("buf.PeekInt16 error %v", err)
		}
		if x != v {
			t.Errorf("After buf.AppendInt16, buf.PeekInt16() = %d, want %d", x, v)
		}
	}

	for _, v := range []int32{math.MaxInt32, 0, math.MinInt32} {
		buf := NewBuffer()
		err := buf.AppendInt32(v)
		if err != nil {
			t.Errorf("buf.AppendInt32 error %v", err)
		}
		x, err := buf.PeekInt32()
		if err != nil {
			t.Errorf("buf.PeekInt32 error %v", err)
		}
		if x != v {
			t.Errorf("After buf.AppendInt32, buf.PeekInt32() = %d, want %d", x, v)
		}
	}

	for _, v := range []int64{math.MaxInt64, 0, math.MinInt64} {
		buf := NewBuffer()
		err := buf.AppendInt64(v)
		if err != nil {
			t.Errorf("buf.AppendInt64 error %v", err)
		}
		x, err := buf.PeekInt64()
		if err != nil {
			t.Errorf("buf.PeekInt64 error %v", err)
		}
		if x != v {
			t.Errorf("After buf.AppendInt64, buf.PeekInt64() = %v, want %d", x, v)
		}
	}

	for _, v := range []int64{math.MaxInt64, 0, math.MinInt64} {
		buf := NewBuffer()
		err := buf.AppendInt64(v)
		if err != nil {
			t.Errorf("buf.AppendInt64 error %v", err)
		}
		b := buf.PeekAllAsByteSlice()
		if len(b) != 8 {
			t.Errorf("buf.PeekAllAsByteSlice() returns a byte slice with %d bytes, want 8 bytes", len(b))
		}
	}

	for _, v := range []uint8{math.MaxUint8, 0, 100} {
		buf := NewBuffer()
		err := buf.AppendUint8(v)
		if err != nil {
			t.Errorf("buf.AppendUint8 error %v", err)
		}
		x, err := buf.PeekUint8()
		if err != nil {
			t.Errorf("buf.PeekUint8 error %v", err)
		}
		if x != v {
			t.Errorf("buf.PeekUint8() = %v, want %d", x, v)
		}
	}

	for _, v := range []uint16{math.MaxUint16, 0, 32000} {
		buf := NewBuffer()
		err := buf.AppendUint16(v)
		if err != nil {
			t.Errorf("buf.AppendUint16 error %v", err)
		}
		x, err := buf.PeekUint16()
		if err != nil {
			t.Errorf("buf.PeekUint16 error %v", err)
		}
		if x != v {
			t.Errorf("buf.PeekUint16() = %v, want %d", x, v)
		}
	}

	for _, v := range []uint32{math.MaxUint32, 0, 2_000_000_000} {
		buf := NewBuffer()
		err := buf.AppendUint32(v)
		if err != nil {
			t.Errorf("buf.AppendUint32 error %v", err)
		}
		x, err := buf.PeekUint32()
		if err != nil {
			t.Errorf("buf.PeekUint32 error %v", err)
		}
		if x != v {
			t.Errorf("buf.PeekUint32() = %v, want %d", x, v)
		}
	}

	for _, v := range []uint64{math.MaxUint64, 0, 100_000_000_000} {
		buf := NewBuffer()
		err := buf.AppendUint64(v)
		if err != nil {
			t.Errorf("buf.AppendUint64 error %v", err)
		}
		x, err := buf.PeekUint64()
		if err != nil {
			t.Errorf("buf.PeekUint64 error %v", err)
		}
		if x != v {
			t.Errorf("buf.PeekUint64() = %v, want %d", x, v)
		}
	}

	for _, v := range []uint64{math.MaxUint64, 0, 199_990_990_990} {
		buf := NewBuffer()
		err := buf.AppendUint64(v)
		if err != nil {
			t.Errorf("buf.AppendUint64 error %v", err)
		}
		b := buf.PeekAllAsByteSlice()
		if len(b) != 8 {
			t.Errorf("buf.PeekAllAsByteSlice() returns a byte slice with %d bytes, want 8 bytes", len(b))
		}
	}
}

func TestRead(t *testing.T) {
	for _, v := range []int64{math.MaxInt64, 0, math.MinInt64} {
		buf := NewBuffer()
		err := buf.AppendInt64(v)
		if err != nil {
			t.Errorf("buf.AppendInt64 error %v", err)
		}
		a, err := buf.ReadInt64()
		if err != nil {
			t.Errorf("buf.readInt64() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readInt64() = %d, want %d", a, v)
		}
	}

	for _, v := range []int32{math.MaxInt32, 0, math.MinInt32} {
		buf := NewBuffer()
		err := buf.AppendInt32(v)
		if err != nil {
			t.Errorf("buf.AppendInt32 error %v", err)
		}
		a, err := buf.ReadInt32()
		if err != nil {
			t.Errorf("buf.readInt32() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readInt32() = %d, want %d", a, v)
		}
	}

	for _, v := range []int16{math.MaxInt16, 0, math.MinInt16} {
		buf := NewBuffer()
		err := buf.AppendInt16(v)
		if err != nil {
			t.Errorf("buf.AppendInt16 error %v", err)
		}
		a, err := buf.ReadInt16()
		if err != nil {
			t.Errorf("buf.readInt16() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readInt16() = %d, want %d", a, v)
		}
	}

	for _, v := range []int8{math.MaxInt8, 0, math.MinInt8} {
		buf := NewBuffer()
		err := buf.AppendInt8(v)
		if err != nil {
			t.Errorf("buf.AppendInt8 error %v", err)
		}
		a, err := buf.ReadInt8()
		if err != nil {
			t.Errorf("buf.readInt8() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readInt8() = %d, want %d", a, v)
		}
	}

	for _, v := range []uint64{math.MaxUint64, 0, 999_999_999_999} {
		buf := NewBuffer()
		err := buf.AppendUint64(v)
		if err != nil {
			t.Errorf("buf.AppendUint64 error %v", err)
		}
		a, err := buf.ReadUint64()
		if err != nil {
			t.Errorf("buf.readUint64() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readUint64() = %d, want %d", a, v)
		}
	}

	for _, v := range []uint32{math.MaxUint32, 0, 2_000_000_000} {
		buf := NewBuffer()
		err := buf.AppendUint32(v)
		if err != nil {
			t.Errorf("buf.AppendUint32 error %v", err)
		}
		a, err := buf.ReadUint32()
		if err != nil {
			t.Errorf("buf.readUint32() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readUint32() = %d, want %d", a, v)
		}
	}

	for _, v := range []uint16{math.MaxUint16, 0, 32000} {
		buf := NewBuffer()
		err := buf.AppendUint16(v)
		if err != nil {
			t.Errorf("buf.AppendUint16 error %v", err)
		}
		a, err := buf.ReadUint16()
		if err != nil {
			t.Errorf("buf.readUint16() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readUint16() = %d, want %d", a, v)
		}
	}

	for _, v := range []uint8{math.MaxUint8, 0, 100} {
		buf := NewBuffer()
		err := buf.AppendUint8(v)
		if err != nil {
			t.Errorf("buf.AppendUint8 error %v", err)
		}
		a, err := buf.ReadUint8()
		if err != nil {
			t.Errorf("buf.readUint8() error %v", err)
		}
		if a != v {
			t.Errorf("buf.readUint8() = %d, want %d", a, v)
		}
	}

	{
		// append
		buf := NewBuffer()
		a1 := int64(math.MaxInt64)
		err := buf.AppendInt64(a1)
		if err != nil {
			t.Errorf("buf.AppendInt64() error %v", err)
		}
		b1 := uint16(math.MaxUint16)
		err = buf.AppendUint16(b1)
		if err != nil {
			t.Errorf("buf.AppendUint16() error %v", err)
		}
		c1 := "tic tac toe"
		buf.Append([]byte(c1))
		d1 := int8(10)
		err = buf.AppendInt8(d1)
		if err != nil {
			t.Errorf("buf.AppendInt8() error %v", err)
		}
		e1 := uint32(math.MaxUint32)
		err = buf.AppendUint32(e1)
		if err != nil {
			t.Errorf("buf.AppendUint32() error %v", err)
		}

		// read
		a2, err := buf.ReadInt64()
		if err != nil {
			t.Errorf("buf.ReadInt64() error %v", err)
		}
		if a2 != a1 {
			t.Errorf("buf.ReadInt64() = %d, want %d", a2, a1)
		}
		b2, err := buf.ReadUint16()
		if err != nil {
			t.Errorf("buf.ReadUint16() error %v", err)
		}
		if b2 != b1 {
			t.Errorf("buf.ReadUint16() = %d, want %d", b2, b1)
		}
		c2 := string(buf.PeekAsByteSlice(len(c1)))
		if c2 != c1 {
			t.Errorf("buf.PeekAsByteSlice() = %s, want %s", c2, c1)
		}
		buf.Retrieve(len(c1))
		d2, err := buf.ReadInt8()
		if err != nil {
			t.Errorf("buf.ReadInt8() error %v", err)
		}
		if d2 != d1 {
			t.Errorf("buf.ReadInt8() = %d, want %d", d2, d1)
		}
		e2, err := buf.ReadUint32()
		if err != nil {
			t.Errorf("buf.ReadUint32() error %v", err)
		}
		if e2 != e1 {
			t.Errorf("buf.ReadUint32() = %d, want %d", e2, e1)
		}
	}
}

func TestCopySlice(t *testing.T) {
	buf := make([]byte, cheapPrepend+initialSize)
	data := []byte("abcde")
	copy(buf[cheapPrepend:], data)
	d2 := buf[cheapPrepend:cheapPrepend+len(data)]
	if string(d2) != string(data) {
		t.Errorf("copy, got %v, want %v", d2, data)
	}
}

func TestRetrieveToByteSlice(t *testing.T) {
	buf := NewBuffer()
	s := "hello, world"
	buf.Append([]byte(s))
	result := make([]byte, len(s))
	buf.retrieveToByteSlice(len(s), result)
	if s != string(result) {
		t.Errorf("retrieveToByteSlice, result is %+v, want %+v", result, []byte(s))
	}
}
