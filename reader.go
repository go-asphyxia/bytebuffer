package bytebuffer

import (
	"io"
	"unicode/utf8"

	aconversion "github.com/go-asphyxia/conversion"
)

type (
	Reader struct {
		ByteBuffer *ByteBuffer

		Point int
	}
)

func (r *Reader) Read(target []byte) (n int, err error) {
	if len(r.ByteBuffer.Bytes) <= r.Point {
		err = io.EOF
		return
	}

	n = copy(target, r.ByteBuffer.Bytes[r.Point:])
	r.Point += n
	return
}

func (r *Reader) ReadByte() (target byte, err error) {
	if len(r.ByteBuffer.Bytes) <= r.Point {
		err = io.EOF
		return
	}

	target = r.ByteBuffer.Bytes[r.Point]
	r.Point++
	return
}

func (r *Reader) ReadRune() (target rune, n int, err error) {
	if len(r.ByteBuffer.Bytes) <= r.Point {
		err = io.EOF
		return
	}

	target, n = utf8.DecodeRune(r.ByteBuffer.Bytes[r.Point:])
	r.Point += n
	return
}

func (r *Reader) ReadString(target string) (n int, err error) {
	if len(r.ByteBuffer.Bytes) <= r.Point {
		err = io.EOF
		return
	}

	n = copy(aconversion.StringToBytesNoCopy(target), r.ByteBuffer.Bytes[r.Point:])
	r.Point += n
	return
}

func (r *Reader) WriteTo(target io.Writer) (n int64, err error) {
	if len(r.ByteBuffer.Bytes) <= r.Point {
		err = io.EOF
		return
	}

	wrote, err := target.Write(r.ByteBuffer.Bytes[r.Point:])
	r.Point += wrote
	n = int64(wrote)
	return
}
