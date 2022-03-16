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

func (b *ByteBuffer) Clone() (target *ByteBuffer) {
	target = &ByteBuffer{}

	if b.Bytes != nil {
		target.Bytes = append(target.Bytes, b.Bytes...)
	}

	return
}

func (b *ByteBuffer) Grow(n int) {
	s := len(b.Bytes) + n

	if s <= cap(b.Bytes) {
		b.Bytes = b.Bytes[:s]
		return
	}

	temp := make([]byte, s)
	copy(temp, b.Bytes)

	b.Bytes = temp
}

func (b *ByteBuffer) Clip(n int) {
	b.Bytes = b.Bytes[:len(b.Bytes)-n]
}

func (b *ByteBuffer) Reset() {
	b.Bytes = b.Bytes[:0]
}

func (b *ByteBuffer) Close() (err error) {
	b = nil
	return
}

func (b *ByteBuffer) Copy() (target []byte) {
	target = append(target, b.Bytes...)
	return
}

func (b *ByteBuffer) StringNoCopy() (target string) {
	target = aconversion.BytesToStringNoCopy(b.Bytes)
	return
}

func (b *ByteBuffer) Set(source []byte) {
	b.Bytes = append(b.Bytes[:0], source...)
}

func (b *ByteBuffer) SetString(source string) {
	b.Bytes = append(b.Bytes[:0], source...)
}

func (b *ByteBuffer) Write(source []byte) (n int, err error) {
	n = len(source)

	b.Bytes = append(b.Bytes, source...)
	return
}

func (b *ByteBuffer) WriteByte(source byte) (err error) {
	b.Bytes = append(b.Bytes, source)
	return
}

func (b *ByteBuffer) WriteRune(source rune) (n int, err error) {
	temp := make([]byte, utf8.UTFMax)

	n = utf8.EncodeRune(temp, source)

	b.Bytes = append(b.Bytes, temp[:n]...)
	return
}

func (b *ByteBuffer) WriteString(source string) (n int, err error) {
	n = len(source)

	b.Bytes = append(b.Bytes, source...)
	return
}

func (b *ByteBuffer) ReadFrom(source io.Reader) (n int64, err error) {
	l := len(b.Bytes)
	c := cap(b.Bytes)
	r := 0

	for {
		if l == c {
			c = (c + 16) * 2

			temp := make([]byte, c)
			copy(temp, b.Bytes)

			b.Bytes = temp
		}

		r, err = source.Read(b.Bytes[l:c])

		n += int64(r)
		l += r

		b.Bytes = b.Bytes[:l]

		if err != nil || l < c {
			if err == io.EOF {
				err = nil
			}

			return
		}
	}
}
