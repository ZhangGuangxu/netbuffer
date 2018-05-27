package netbuffer

import (
	"bytes"
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

	buf.AppendInt64(int64(9223372036854770000))
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

	{
		buf := NewBuffer()
		buf.AppendInt8(int8(1))
		if buf.ReadableBytes() != 1 {
			t.Errorf("After buf.appendInt8(), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 1)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt16(int16(32000))
		if buf.ReadableBytes() != 2 {
			t.Errorf("After buf.appendInt16(), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 2)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt32(int32(32000))
		if buf.ReadableBytes() != 4 {
			t.Errorf("After buf.appendInt32(), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 4)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt64(int64(9223372036854770000))
		if buf.ReadableBytes() != 8 {
			t.Errorf("After buf.appendInt64(int64(9223372036854770000)), buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 8)
		}
	}
}

func TestPrepend(t *testing.T) {
	{
		buf := NewBufferWithSize(4)
		buf.PrependInt8(int8(2))
		if buf.prependableBytes() != cheapPrepend-1 {
			t.Errorf("After buf.prependInt8(int8(2)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-1)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.PrependInt16(int16(32000))
		if buf.prependableBytes() != cheapPrepend-2 {
			t.Errorf("After buf.prependInt16(int16(32000)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-2)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.PrependInt32(int32(32000))
		if buf.prependableBytes() != cheapPrepend-4 {
			t.Errorf("After buf.prependInt32(int32(32000)), buf.prependableBytes() is %d, want %d", buf.prependableBytes(), cheapPrepend-4)
		}
	}

	{
		buf := NewBufferWithSize(4)
		buf.PrependInt64(int64(32000))
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
		buf.AppendInt8(int8(1))
		if buf.PeekInt8() != int8(1) {
			t.Errorf("After buf.AppendInt8(int8(1)), buf.PeekInt8() = %d, want %d", buf.PeekInt8(), 1)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt16(int16(32000))
		if buf.PeekInt16() != int16(32000) {
			t.Errorf("After buf.AppendInt16(int16(32000)), buf.PeekInt16() = %d, want %d", buf.PeekInt16(), 32000)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt32(int32(32000))
		if buf.PeekInt32() != int32(32000) {
			t.Errorf("After buf.AppendInt32(int32(32000)), buf.PeekInt32() = %d, want %d", buf.PeekInt32(), 32000)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt64(int64(9223372036854770000))
		if buf.PeekInt64() != int64(9223372036854770000) {
			t.Errorf("After buf.AppendInt64(int64(9223372036854770000)), buf.PeekInt64() = %v, want %d", buf.PeekInt64(), 9223372036854770000)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt64(int64(9223372036854770000))
		b := buf.PeekAllAsByteSlice()
		if len(b) != 8 {
			t.Errorf("buf.PeekAllAsByteSlice() returns a byte slice with %d bytes, want 8 bytes", len(b))
		}
	}
}

func TestRetrieve(t *testing.T) {
	{
		buf := NewBuffer()
		buf.AppendInt64(int64(9223372036854770000))
		buf.RetrieveInt64()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt64(int64(9223372036854770000))
		buf.RetrieveInt16()
		if buf.ReadableBytes() != 6 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 6)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt64(int64(9223372036854770000))
		data := buf.retrieveAllAsByteSlice()
		if len(data) != 8 {
			t.Errorf("len(data) = %d, want %d", len(data), 8)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt64(int64(9223372036854770000))
		str := buf.retrieveAllAsString()
		if len(str) == 0 {
			t.Error("len(data) should be >0")
		}
	}

	{
		buf := NewBuffer()
		buf.PrependInt64(int64(9223372036854770000))
		buf.RetrieveInt64()
		if buf.ReadableBytes() != 0 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 0)
		}
	}

	{
		buf := NewBuffer()
		buf.PrependInt64(int64(9223372036854770000))
		buf.RetrieveInt32()
		if buf.ReadableBytes() != 4 {
			t.Errorf("buf.ReadableBytes() = %d, want %d", buf.ReadableBytes(), 4)
		}
	}

	{
		buf := NewBuffer()
		buf.PrependInt64(int64(9223372036854770000))
		data := buf.retrieveAllAsByteSlice()
		if len(data) != 8 {
			t.Errorf("len(data) = %d, want %d", len(data), 8)
		}
	}

	{
		buf := NewBuffer()
		buf.PrependInt64(int64(9223372036854770000))
		str := buf.retrieveAllAsString()
		if len(str) == 0 {
			t.Error("len(data) should be >0")
		}
	}
}

func TestRead(t *testing.T) {
	{
		buf := NewBuffer()
		buf.AppendInt64(int64(1))
		a := buf.ReadInt64()
		if a != 1 {
			t.Errorf("buf.readInt64() = %d, want %d", a, 1)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt32(int32(1))
		a := buf.ReadInt32()
		if a != 1 {
			t.Errorf("buf.readInt32() = %d, want %d", a, 1)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt16(int16(1))
		a := buf.ReadInt16()
		if a != 1 {
			t.Errorf("buf.readInt16() = %d, want %d", a, 1)
		}
	}

	{
		buf := NewBuffer()
		buf.AppendInt8(int8(1))
		a := buf.ReadInt8()
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
