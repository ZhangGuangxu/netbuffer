package netbuffer

import (
	"bytes"
	"encoding/binary"
)

const (
	cheapPrepend = 8
	initialSize  = 1024 // default count of byte of buffer
)

// Buffer wraps a buffer for net data.
type Buffer struct {
	buf         []byte
	readerIndex int
	writerIndex int
}

// NewBuffer returns a buffer with default length.
// Usually this function is enough.
func NewBuffer() *Buffer {
	return NewBufferWithSize(initialSize)
}

// NewBufferWithSize returns a buffer with length you specified.
func NewBufferWithSize(s int) *Buffer {
	return &Buffer{
		buf:         make([]byte, cheapPrepend+s),
		readerIndex: cheapPrepend,
		writerIndex: cheapPrepend,
	}
}

// ReadableBytes returns count of byte in this buffer.
func (b *Buffer) ReadableBytes() int {
	return b.writerIndex - b.readerIndex
}

// WritableBytes returns byte count you can write to this buffer
// without memory allocation.
func (b *Buffer) WritableBytes() int {
	return len(b.buf) - b.writerIndex
}

func (b *Buffer) prependableBytes() int {
	return b.readerIndex
}

// WritableByteSlice returns a byte slice you can write bytes of
// at most its length to it.
func (b *Buffer) WritableByteSlice() []byte {
	return b.buf[b.writerIndex:len(b.buf)]
}

// Append adds data to this buffer.
func (b *Buffer) Append(data []byte) {
	b.appendWithLen(data, len(data))
}

// appendWithLen adds length byte in data to this buffer.
func (b *Buffer) appendWithLen(data []byte, length int) {
	b.ensureWritableBytes(length)
	copy(b.buf[b.writerIndex:b.writerIndex+length], data)
	b.HasWritten(length)
}

func (b *Buffer) ensureWritableBytes(length int) {
	if b.WritableBytes() < length {
		b.makeSpace(length)
	}
}

// HasWritten add length of the content of buffer when necessary
func (b *Buffer) HasWritten(length int) {
	b.writerIndex += length
}

// AppendInt64 appends a int64 to this buffer.
func (b *Buffer) AppendInt64(x int64) error {
	return b.appendIntn(8, x)
}

// AppendInt32 appends a int32 to this buffer.
func (b *Buffer) AppendInt32(x int32) error {
	return b.appendIntn(4, x)
}

// AppendInt16 appends a int16 to this buffer.
func (b *Buffer) AppendInt16(x int16) error {
	return b.appendIntn(2, x)
}

// AppendInt8 appends a int8 to this buffer.
func (b *Buffer) AppendInt8(x int8) error {
	return b.appendIntn(1, x)
}

func (b *Buffer) appendIntn(s int, x interface{}) error {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, x); err != nil {
		return err
	}
	b.Append(buf.Bytes())
	return nil
}

// PrependInt64 prepend a int64 to this buffer.
func (b *Buffer) PrependInt64(x int64) error {
	return b.prependIntn(8, x)
}

// PrependInt32 prepend a int32 to this buffer.
func (b *Buffer) PrependInt32(x int32) error {
	return b.prependIntn(4, x)
}

// PrependInt16 prepend a int16 to this buffer.
func (b *Buffer) PrependInt16(x int16) error {
	return b.prependIntn(2, x)
}

// PrependInt8 prepend a int8 to this buffer.
func (b *Buffer) PrependInt8(x int8) error {
	return b.prependIntn(1, x)
}

func (b *Buffer) prependIntn(s int, x interface{}) error {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, x); err != nil {
		return err
	}
	b.prepend(buf.Bytes())
	return nil
}

func (b *Buffer) prepend(data []byte) {
	length := len(data)
	b.readerIndex -= length
	copy(b.buf[b.readerIndex:b.readerIndex+length], data)
}

// Retrieve removes length readable bytes.
func (b *Buffer) Retrieve(length int) {
	if length < b.ReadableBytes() {
		b.readerIndex += length
	} else {
		b.RetrieveAll()
	}
}

func (b *Buffer) RetrieveAll() {
	b.readerIndex = cheapPrepend
	b.writerIndex = cheapPrepend
}

// RetrieveInt64 removes a int64(8 bytes) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveInt64() {
	b.Retrieve(8)
}

// RetrieveInt32 removes a int32(4 bytes) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveInt32() {
	b.Retrieve(4)
}

// RetrieveInt16 removes a int16(2 bytes) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveInt16() {
	b.Retrieve(2)
}

// RetrieveInt8 removes a int8(1 byte) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveInt8() {
	b.Retrieve(1)
}

func (b *Buffer) retrieveAllAsByteSlice() []byte {
	return b.retrieveAsByteSlice(b.ReadableBytes())
}

func (b *Buffer) retrieveAsByteSlice(length int) []byte {
	result := make([]byte, 0, length)
	result = append(result, b.buf[b.readerIndex:b.readerIndex+length]...)
	b.Retrieve(length)
	return result
}

func (b *Buffer) retrieveAllAsString() string {
	return b.retrieveAsString(b.ReadableBytes())
}

func (b *Buffer) retrieveAsString(length int) string {
	result := string(b.buf[b.readerIndex : b.readerIndex+length])
	b.Retrieve(length)
	return result
}

// PeekAllAsByteSlice returns a byte slice with all readable bytes of this buffer.
// You MUST NOT modify the content of the returned slice.
func (b *Buffer) PeekAllAsByteSlice() []byte {
	return b.PeekAsByteSlice(b.ReadableBytes())
}

// PeekAsByteSlice returns a byte slice which contains length count bytes.
// You MUST NOT modify the content of the returned slice.
func (b *Buffer) PeekAsByteSlice(length int) []byte {
	return b.buf[b.readerIndex : b.readerIndex+length]
}

// PeekInt64 parses a int64 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt64() (int64, error) {
	var x int64
	err := b.peekIntn(8, &x)
	return x, err
}

// PeekInt32 parses a int32 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt32() (int32, error) {
	var x int32
	err := b.peekIntn(4, &x)
	return x, err
}

// PeekInt16 parses a int16 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt16() (int16, error) {
	var x int16
	err := b.peekIntn(2, &x)
	return x, err
}

// PeekInt8 parses a int8 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt8() (int8, error) {
	var x int8
	err := b.peekIntn(1, &x)
	return x, err
}

func (b *Buffer) peekIntn(s int, x interface{}) error {
	buf := &bytes.Buffer{}
	if _, err := buf.Write(b.buf[b.readerIndex : b.readerIndex+s]); err != nil {
		return err
	}
	return binary.Read(buf, binary.BigEndian, x)
}

// ReadInt64 parses a int64 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt64() (int64, error) {
	x, err := b.PeekInt64()
	if err != nil {
		return 0, err
	}
	b.RetrieveInt64()
	return x, nil
}

// ReadInt32 parses a int32 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt32() (int32, error) {
	x, err := b.PeekInt32()
	if err != nil {
		return 0, err
	}
	b.RetrieveInt32()
	return x, nil
}

// ReadInt16 parses a int16 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt16() (int16, error) {
	x, err := b.PeekInt16()
	if err != nil {
		return 0, err
	}
	b.RetrieveInt16()
	return x, nil
}

// ReadInt8 parses a int8 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt8() (int8, error) {
	x, err := b.PeekInt8()
	if err != nil {
		return 0, err
	}
	b.RetrieveInt8()
	return x, nil
}

func (b *Buffer) makeSpace(length int) {
	writable := b.WritableBytes()
	if writable+b.prependableBytes() >= length+cheapPrepend {
		readable := b.ReadableBytes()
		copy(b.buf[cheapPrepend:cheapPrepend+readable], b.buf[b.readerIndex:b.writerIndex])
		b.readerIndex = cheapPrepend
		b.writerIndex = b.readerIndex + readable
	} else {
		more := length - writable
		b.buf = append(b.buf, make([]byte, more)...)
	}
}
