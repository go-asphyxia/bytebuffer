package bytebuffer

import (
	"io"
	"unicode/utf8"

	aconversion "github.com/go-asphyxia/conversion"
)

type (
	ByteBuffer struct {
		Bytes []byte
	}
)

func (bb *ByteBuffer) Clone() (target *ByteBuffer) {
	target = &ByteBuffer{
		Bytes: make([]byte, len(bb.Bytes)),
	}

	copy(target.Bytes, bb.Bytes)
	return
}

func (bb *ByteBuffer) Grow(n int) {
	s := len(bb.Bytes) + n

	if s <= cap(bb.Bytes) {
		bb.Bytes = bb.Bytes[:s]
		return
	}

	temp := make([]byte, s)
	copy(temp, bb.Bytes)

	bb.Bytes = temp
}

func (bb *ByteBuffer) Clip(n int) {
	bb.Bytes = bb.Bytes[:len(bb.Bytes)-n]
}

func (bb *ByteBuffer) Reset() {
	bb.Bytes = bb.Bytes[:0]
}

func (bb *ByteBuffer) Close() (err error) {
	bb.Bytes = nil
	return
}

func (bb *ByteBuffer) Copy() (target []byte) {
	target = make([]byte, len(bb.Bytes))

	copy(target, bb.Bytes)
	return
}

func (bb *ByteBuffer) StringNoCopy() (target string) {
	target = aconversion.BytesToStringNoCopy(bb.Bytes)
	return
}

func (bb *ByteBuffer) Set(source []byte) {
	bb.Bytes = append(bb.Bytes[:0], source...)
}

func (bb *ByteBuffer) SetString(source string) {
	bb.Bytes = append(bb.Bytes[:0], source...)
}

func (bb *ByteBuffer) Write(source []byte) (n int, err error) {
	n = len(source)

	bb.Bytes = append(bb.Bytes, source...)
	return
}

func (bb *ByteBuffer) WriteByte(source byte) (err error) {
	bb.Bytes = append(bb.Bytes, source)
	return
}

func (bb *ByteBuffer) WriteRune(source rune) (n int, err error) {
	temp := make([]byte, utf8.UTFMax)

	n = utf8.EncodeRune(temp, source)

	bb.Bytes = append(bb.Bytes, temp[:n]...)
	return
}

func (bb *ByteBuffer) WriteString(source string) (n int, err error) {
	n = len(source)

	bb.Bytes = append(bb.Bytes, source...)
	return
}

func (bb *ByteBuffer) ReadFrom(source io.Reader) (n int64, err error) {
	l := len(bb.Bytes)
	c := cap(bb.Bytes)
	r := 0

	for {
		if l == c {
			c = (c + 16) * 2

			temp := make([]byte, c)
			copy(temp, bb.Bytes)

			bb.Bytes = temp
		}

		r, err = source.Read(bb.Bytes[l:c])

		n += int64(r)
		l += r

		bb.Bytes = bb.Bytes[:l]

		if err != nil || l < c {
			if err == io.EOF {
				err = nil
			}

			return
		}
	}
}
