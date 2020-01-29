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
	return b.appendInteger(x)
}

// AppendInt32 appends a int32 to this buffer.
func (b *Buffer) AppendInt32(x int32) error {
	return b.appendInteger(x)
}

// AppendInt16 appends a int16 to this buffer.
func (b *Buffer) AppendInt16(x int16) error {
	return b.appendInteger(x)
}

// AppendInt8 appends a int8 to this buffer.
func (b *Buffer) AppendInt8(x int8) error {
	return b.appendInteger(x)
}

// AppendUint64 appends a uint64 to this buffer.
func (b *Buffer) AppendUint64(x uint64) error {
	return b.appendInteger(x)
}

// AppendUint32 appends a uint32 to this buffer.
func (b *Buffer) AppendUint32(x uint32) error {
	return b.appendInteger(x)
}

// AppendUint16 appends a uint16 to this buffer.
func (b *Buffer) AppendUint16(x uint16) error {
	return b.appendInteger(x)
}

// AppendUint8 appends a uint8 to this buffer.
func (b *Buffer) AppendUint8(x uint8) error {
	return b.appendInteger(x)
}

func (b *Buffer) appendInteger(x interface{}) error {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, x); err != nil {
		return err
	}
	b.Append(buf.Bytes())
	return nil
}

// PrependInt64 prepend a int64 to this buffer.
func (b *Buffer) PrependInt64(x int64) error {
	return b.prependInteger(x)
}

// PrependInt32 prepend a int32 to this buffer.
func (b *Buffer) PrependInt32(x int32) error {
	return b.prependInteger(x)
}

// PrependInt16 prepend a int16 to this buffer.
func (b *Buffer) PrependInt16(x int16) error {
	return b.prependInteger(x)
}

// PrependInt8 prepend a int8 to this buffer.
func (b *Buffer) PrependInt8(x int8) error {
	return b.prependInteger(x)
}

// PrependUint64 prepend a uint64 to this buffer.
func (b *Buffer) PrependUint64(x uint64) error {
	return b.prependInteger(x)
}

// PrependUint32 prepend a uint32 to this buffer.
func (b *Buffer) PrependUint32(x uint32) error {
	return b.prependInteger(x)
}

// PrependUint16 prepend a uint16 to this buffer.
func (b *Buffer) PrependUint16(x uint16) error {
	return b.prependInteger(x)
}

// PrependUint8 prepend a uint8 to this buffer.
func (b *Buffer) PrependUint8(x uint8) error {
	return b.prependInteger(x)
}

func (b *Buffer) prependInteger(x interface{}) error {
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

// RetrieveUint64 removes a uint64(8 bytes) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveUint64() {
	b.Retrieve(8)
}

// RetrieveUint32 removes a uint32(4 bytes) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveUint32() {
	b.Retrieve(4)
}

// RetrieveUint16 removes a uint16(2 bytes) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveUint16() {
	b.Retrieve(2)
}

// RetrieveUint8 removes a uint8(1 byte) from the beginning of
// the readable bytes of this buffer.
func (b *Buffer) RetrieveUint8() {
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

// len(result) == length
func (b *Buffer) retrieveToByteSlice(length int, result []byte) {
	copy(result, b.buf[b.readerIndex:b.readerIndex+length])
	b.Retrieve(length)
}

func (b *Buffer) retrieveAllAsString() string {
	return b.retrieveAsString(b.ReadableBytes())
}

func (b *Buffer) retrieveAsString(length int) string {
	result := string(b.buf[b.readerIndex:b.readerIndex+length])
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
	return b.buf[b.readerIndex:b.readerIndex+length]
}

// PeekInt64 parses a int64 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt64() (x int64, err error) {
	err = b.peekInteger(8, &x)
	return
}

// PeekInt32 parses a int32 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt32() (x int32, err error) {
	err = b.peekInteger(4, &x)
	return
}

// PeekInt16 parses a int16 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt16() (x int16, err error) {
	err = b.peekInteger(2, &x)
	return
}

// PeekInt8 parses a int8 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekInt8() (x int8, err error) {
	err = b.peekInteger(1, &x)
	return
}

// PeekUint64 parses a uint64 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekUint64() (x uint64, err error) {
	err = b.peekInteger(8, &x)
	return
}

// PeekUint32 parses a uint32 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekUint32() (x uint32, err error) {
	err = b.peekInteger(4, &x)
	return
}

// PeekUint16 parses a uint16 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekUint16() (x uint16, err error) {
	err = b.peekInteger(2, &x)
	return
}

// PeekUint8 parses a uint8 from the beginning of the readable bytes of this buffer.
// This function does not modify this buffer.
func (b *Buffer) PeekUint8() (x uint8, err error) {
	err = b.peekInteger(1, &x)
	return
}

func (b *Buffer) peekInteger(s int, x interface{}) error {
	buf := &bytes.Buffer{}
	if _, err := buf.Write(b.buf[b.readerIndex:b.readerIndex+s]); err != nil {
		return err
	}
	return binary.Read(buf, binary.BigEndian, x)
}

// ReadInt64 parses a int64 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt64() (x int64, err error) {
	x, err = b.PeekInt64()
	if err != nil {
		return
	}
	b.RetrieveInt64()
	return
}

// ReadInt32 parses a int32 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt32() (x int32, err error) {
	x, err = b.PeekInt32()
	if err != nil {
		return
	}
	b.RetrieveInt32()
	return
}

// ReadInt16 parses a int16 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt16() (x int16, err error) {
	x, err = b.PeekInt16()
	if err != nil {
		return
	}
	b.RetrieveInt16()
	return
}

// ReadInt8 parses a int8 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadInt8() (x int8, err error) {
	x, err = b.PeekInt8()
	if err != nil {
		return
	}
	b.RetrieveInt8()
	return
}

// ReadUint64 parses a uint64 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadUint64() (x uint64, err error) {
	x, err = b.PeekUint64()
	if err != nil {
		return
	}
	b.RetrieveUint64()
	return
}

// ReadUnt32 parses a uint32 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadUint32() (x uint32, err error) {
	x, err = b.PeekUint32()
	if err != nil {
		return
	}
	b.RetrieveUint32()
	return
}

// ReadUint16 parses a uint16 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadUint16() (x uint16, err error) {
	x, err = b.PeekUint16()
	if err != nil {
		return
	}
	b.RetrieveUint16()
	return
}

// ReadUint8 parses a uint8 from the beginning of the readable bytes of this buffer and
// changes readable bytes of this buffer.
func (b *Buffer) ReadUint8() (x uint8, err error) {
	x, err = b.PeekUint8()
	if err != nil {
		return
	}
	b.RetrieveUint8()
	return
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
