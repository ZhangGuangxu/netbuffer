package netbuffer

import (
	"testing"
	"bytes"
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
}

func TestAppend(t *testing.T) {
	buf := NewBuffer()

	data := []byte("abcde")

	buf.Append(data)
	if buf.ReadableBytes() != len(data) {
		t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), len(data))
	}

	writable := cap(buf.buf)-(cheapPrepend + len(data))
	if buf.WritableBytes() != writable {
		t.Errorf("buf.WritableBytes() = %d, want %d", buf.WritableBytes(), writable)
	}

	if buf.readerIndex != cheapPrepend {
		t.Errorf("buf.readerIndex = %d, not %d", buf.readerIndex, cheapPrepend)
	}
	if buf.writerIndex != cheapPrepend + len(data) {
		t.Errorf("buf.writerIndex = %d, not %d", buf.writerIndex, cheapPrepend + len(data))
	}

	data = bytes.Repeat([]byte("a"), 2048)
	buf.Append(data)
	newCap := cheapPrepend+len([]byte("abcde"))+2048
	if len(buf.buf) != newCap {
		t.Errorf("len(buf.buf) = %d, want %d", len(buf.buf), newCap)
	}

	{
		buf := NewBuffer()
		buf.appendInt8(int8(1))
		if buf.ReadableBytes() != 1 {
			t.Errorf("After buf.appendInt8(), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 1)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt16(int16(32000))
		if buf.ReadableBytes() != 2 {
			t.Errorf("After buf.appendInt16(), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 2)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt32(int32(32000))
		if buf.ReadableBytes() != 4 {
			t.Errorf("After buf.appendInt32(), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 4)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt64(int64(9223372036854770000))
		if buf.ReadableBytes() != 8 {
			t.Errorf("After buf.appendInt64(int64(9223372036854770000)), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 8)
		}
	}
}

func TestPrepend(t *testing.T) {
	{
		buf := NewBufferWithSize(4)
		buf.prependInt8(int8(2))
		if buf.prependableBytes() != cheapPrepend-1 {
			t.Errorf("After buf.prependInt8(int8(2)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-1)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.prependInt16(int16(32000))
		if buf.prependableBytes() != cheapPrepend-2 {
			t.Errorf("After buf.prependInt16(int16(32000)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-2)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.prependInt32(int32(32000))
		if buf.prependableBytes() != cheapPrepend-4 {
			t.Errorf("After buf.prependInt32(int32(32000)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-4)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.prependInt64(int64(32000))
		if buf.prependableBytes() != cheapPrepend-8 {
			t.Errorf("After buf.prependInt64(int64(32000)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-8)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.prepend(make([]byte, 5))
		if buf.prependableBytes() != cheapPrepend-5 {
			t.Errorf("After buf.prepend(make([]byte, 5)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-5)
		}
	}
}

func TestPeek(t *testing.T) {
	{
		buf := NewBuffer()
		buf.appendInt8(int8(1))
		if buf.peekInt8() != int8(1) {
			t.Errorf("After buf.appendInt8(int8(1)), buf.peekInt8() = %d, want %d", buf.peekInt8(), 1)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt16(int16(32000))
		if buf.peekInt16() != int16(32000) {
			t.Errorf("After buf.appendInt16(int16(32000)), buf.peekInt16() = %d, want %d", buf.peekInt16(), 32000)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt32(int32(32000))
		if buf.peekInt32() != int32(32000) {
			t.Errorf("After buf.appendInt32(int32(32000)), buf.peekInt32() = %d, want %d", buf.peekInt32(), 32000)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt64(int64(9223372036854770000))
		if buf.peekInt64() != int64(9223372036854770000) {
			t.Errorf("After buf.appendInt64(int64(9223372036854770000)), buf.peekInt64() = %v, want %d", buf.peekInt64(), 9223372036854770000)
		}
	}
}

func TestRetrieve(t *testing.T) {
	{
		buf := NewBuffer()
		buf.appendInt64(int64(9223372036854770000))
		buf.retrieveInt64()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt64(int64(9223372036854770000))
		buf.retrieveInt16()
		if buf.ReadableBytes() != 6 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 6)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt64(int64(9223372036854770000))
		data := buf.retrieveAllAsByteSlice()
		if len(data) != 8 {
			t.Errorf("len(data) = %d, want %d", len(data), 8)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt64(int64(9223372036854770000))
		str := buf.retrieveAllAsString()
		if len(str) == 0 {
			t.Error("len(data) should be >0")
		}
	}

	{
		buf := NewBuffer()
		buf.prependInt64(int64(9223372036854770000))
		buf.retrieveInt64()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	{
		buf := NewBuffer()
		buf.prependInt64(int64(9223372036854770000))
		buf.retrieveInt32()
		if buf.ReadableBytes() != 4 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 4)
		}
	}

	{
		buf := NewBuffer()
		buf.prependInt64(int64(9223372036854770000))
		data := buf.retrieveAllAsByteSlice()
		if len(data) != 8 {
			t.Errorf("len(data) = %d, want %d", len(data), 8)
		}
	}

	{
		buf := NewBuffer()
		buf.prependInt64(int64(9223372036854770000))
		str := buf.retrieveAllAsString()
		if len(str) == 0 {
			t.Error("len(data) should be >0")
		}
	}
}

func TestRead(t *testing.T) {
	{
		buf := NewBuffer()
		buf.appendInt64(int64(1))
		a := buf.readInt64()
		if a != 1 {
			t.Errorf("buf.readInt64() = %d, want %d", a, 1)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt32(int32(1))
		a := buf.readInt32()
		if a != 1 {
			t.Errorf("buf.readInt32() = %d, want %d", a, 1)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt16(int16(1))
		a := buf.readInt16()
		if a != 1 {
			t.Errorf("buf.readInt16() = %d, want %d", a, 1)
		}
	}

	{
		buf := NewBuffer()
		buf.appendInt8(int8(1))
		a := buf.readInt8()
		if a != 1 {
			t.Errorf("buf.readInt8() = %d, want %d", a, 1)
		}
	}
}

func TestCopySlice(t *testing.T) {
	buf := make([]byte, cheapPrepend+initialSize, cheapPrepend+initialSize)
	wIdx := cheapPrepend
	data := []byte("abcde")
	copy(buf[wIdx:], data)
}
