package netbuffer

import (
	"bytes"
	"encoding/binary"
)

const (
	cheapPrepend = 8
	initialSize  = 1024 // default count of byte of buffer
)

var crlf = []byte("\r\n")

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
func NewBufferWithSize(initialSize int) *Buffer {
	return &Buffer{
		buf:         make([]byte, cheapPrepend+initialSize),
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
	b.hasWritten(length)
}

func (b *Buffer) ensureWritableBytes(length int) {
	if b.WritableBytes() < length {
		b.makeSpace(length)
	}
}

func (b *Buffer) hasWritten(length int) {
	b.writerIndex += length
}

func (b *Buffer) appendInt64(x int64) {
	b.appendIntn(8, x)
}

func (b *Buffer) appendInt32(x int32) {
	b.appendIntn(4, x)
}

func (b *Buffer) appendInt16(x int16) {
	b.appendIntn(2, x)
}

func (b *Buffer) appendInt8(x int8) {
	b.appendIntn(1, x)
}

func (b *Buffer) appendIntn(s int, x interface{}) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, x)
	b.Append(buf.Bytes())
}

func (b *Buffer) prependInt64(x int64) {
	b.prependIntn(8, x)
}

func (b *Buffer) prependInt32(x int32) {
	b.prependIntn(4, x)
}

func (b *Buffer) prependInt16(x int16) {
	b.prependIntn(2, x)
}

func (b *Buffer) prependInt8(x int8) {
	b.prependIntn(1, x)
}

func (b *Buffer) prependIntn(s int, x interface{}) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, x)
	b.prepend(buf.Bytes())
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
		b.retrieveAll()
	}
}

func (b *Buffer) retrieveAll() {
	b.readerIndex = cheapPrepend
	b.writerIndex = cheapPrepend
}

func (b *Buffer) retrieveInt64() {
	b.Retrieve(8)
}

func (b *Buffer) retrieveInt32() {
	b.Retrieve(4)
}

func (b *Buffer) retrieveInt16() {
	b.Retrieve(2)
}

func (b *Buffer) retrieveInt8() {
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

// PeekAllAsByteSlice returns a internal byte slice with all readable bytes directly.
func (b *Buffer) PeekAllAsByteSlice() []byte {
	return b.peekAsByteSlice(b.ReadableBytes())
}

func (b *Buffer) peekAsByteSlice(length int) []byte {
	return b.buf[b.readerIndex : b.readerIndex+length]
}

func (b *Buffer) peekInt64() int64 {
	var x int64
	b.peekIntn(8, &x)
	return x
}

func (b *Buffer) peekInt32() int32 {
	var x int32
	b.peekIntn(4, &x)
	return x
}

func (b *Buffer) peekInt16() int16 {
	var x int16
	b.peekIntn(2, &x)
	return x
}

func (b *Buffer) peekInt8() int8 {
	var x int8
	b.peekIntn(1, &x)
	return x
}

func (b *Buffer) peekIntn(s int, x interface{}) {
	buf := &bytes.Buffer{}
	buf.Write(b.buf[b.readerIndex : b.readerIndex+s])
	binary.Read(buf, binary.BigEndian, x)
}

func (b *Buffer) readInt64() int64 {
	x := b.peekInt64()
	b.retrieveInt64()
	return x
}

func (b *Buffer) readInt32() int32 {
	x := b.peekInt32()
	b.retrieveInt32()
	return x
}

func (b *Buffer) readInt16() int16 {
	x := b.peekInt16()
	b.retrieveInt16()
	return x
}

func (b *Buffer) readInt8() int8 {
	x := b.peekInt8()
	b.retrieveInt8()
	return x
}

func (b *Buffer) makeSpace(length int) {
	writable := b.WritableBytes()
	if writable+b.prependableBytes() < length+cheapPrepend {
		more := length - writable
		b.buf = append(b.buf, make([]byte, more)...)
	} else {
		readable := b.ReadableBytes()
		copy(b.buf[cheapPrepend:cheapPrepend+readable], b.buf[b.readerIndex:b.writerIndex])
		b.readerIndex = cheapPrepend
		b.writerIndex = b.readerIndex + readable
	}
}
