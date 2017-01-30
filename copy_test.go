package chunked

import (
	"bytes"
	"io"
	"testing"
)

func TestCopy(t *testing.T) {
	src := new(bufWrap)
	dst := new(bufWrap)
	exp := bb[:26]

	src.Write(exp)
	Copy(dst, src)

	if out := dst.Bytes(); !bytes.Equal(exp, out) {
		t.Errorf("expected %q, got %q", exp, out)
	}
}

func TestCopyN(t *testing.T) {
	src := new(bufWrap)
	dst := new(bufWrap)
	exp := bb[:26]

	src.Write(bb)
	CopyN(dst, src, 26)

	if out := dst.Bytes(); !bytes.Equal(exp, out) {
		t.Errorf("expected %q, got %q", exp, out)
	}
}

func TestCopy_ReadFrom(t *testing.T) {
	src := new(bufWrap)
	dst := new(Buffer) // implemets io.ReaderFrom
	exp := bb[:26]

	src.Write(exp)
	Copy(dst, src)

	if out := dst.Bytes(); !bytes.Equal(exp, out) {
		t.Errorf("expected %q, got %q", exp, out)
	}
}

func TestCopy_WriteTo(t *testing.T) {
	src := new(Buffer) // implemets io.WriterTo
	dst := new(bufWrap)
	exp := bb[:26]

	src.Write(exp)
	Copy(dst, src)

	if out := dst.Bytes(); !bytes.Equal(exp, out) {
		t.Errorf("expected %q, got %q", exp, out)
	}
}

type bufWrap struct {
	Buffer
	io.ReaderFrom
	io.WriterTo
}
